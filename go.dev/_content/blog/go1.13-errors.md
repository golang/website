---
title: Working with Errors in Go 1.13
date: 2019-10-17
by:
- Damien Neil and Jonathan Amsterdam
tags:
- errors
- technical
summary: How to use the new Go 1.13 error interfaces and functions.
---

## Introduction

Go’s treatment of [errors as values](https://blog.golang.org/errors-are-values)
has served us well over the last decade. Although the standard library’s support
for errors has been minimal—just the `errors.New` and `fmt.Errorf` functions,
which produce errors that contain only a message—the built-in `error` interface
allows Go programmers to add whatever information they desire. All it requires
is a type that implements an `Error` method:

	type QueryError struct {
		Query string
		Err   error
	}

	func (e *QueryError) Error() string { return e.Query + ": " + e.Err.Error() }

Error types like this one are ubiquitous, and the information they store varies
widely, from timestamps to filenames to server addresses. Often, that
information includes another, lower-level error to provide additional context.

The pattern of one error containing another is so pervasive in Go code that,
after [extensive discussion](/issue/29934), Go 1.13 added
explicit support for it. This post describes the additions to the standard
library that provide that support: three new functions in the `errors` package,
and a new formatting verb for `fmt.Errorf`.

Before describing the changes in detail, let's review how errors are examined
and constructed in previous versions of the language.

## Errors before Go 1.13

### Examining errors

Go errors are values. Programs make decisions based on those values in a few
ways. The most common is to compare an error to `nil` to see if an operation
failed.

	if err != nil {
		// something went wrong
	}

Sometimes we compare an error to a known _sentinel_ value, to see if a specific error has occurred.

	var ErrNotFound = errors.New("not found")

	if err == ErrNotFound {
		// something wasn't found
	}

An error value may be of any type which satisfies the language-defined `error`
interface. A program can use a type assertion or type switch to view an error
value as a more specific type.

	type NotFoundError struct {
		Name string
	}

	func (e *NotFoundError) Error() string { return e.Name + ": not found" }

	if e, ok := err.(*NotFoundError); ok {
		// e.Name wasn't found
	}

### Adding information

Frequently a function passes an error up the call stack while adding information
to it, like a brief description of what was happening when the error occurred. A
simple way to do this is to construct a new error that includes the text of the
previous one:

	if err != nil {
		return fmt.Errorf("decompress %v: %v", name, err)
	}

Creating a new error with `fmt.Errorf` discards everything from the original
error except the text. As we saw above with `QueryError`, we may sometimes want
to define a new error type that contains the underlying error, preserving it for
inspection by code. Here is `QueryError` again:

	type QueryError struct {
		Query string
		Err   error
	}

Programs can look inside a `*QueryError` value to make decisions based on the
underlying error. You'll sometimes see this referred to as "unwrapping" the
error.

	if e, ok := err.(*QueryError); ok && e.Err == ErrPermission {
		// query failed because of a permission problem
	}

The `os.PathError` type in the standard library is another example of one error which contains another.

## Errors in Go 1.13

### The Unwrap method

Go 1.13 introduces new features to the `errors` and `fmt` standard library
packages to simplify working with errors that contain other errors. The most
significant of these is a convention rather than a change: an error which
contains another may implement an `Unwrap` method returning the underlying
error. If `e1.Unwrap()` returns `e2`, then we say that `e1` _wraps_ `e2`, and
that you can _unwrap_ `e1` to get `e2`.

Following this convention, we can give the `QueryError` type above an `Unwrap`
method that returns its contained error:

	func (e *QueryError) Unwrap() error { return e.Err }

The result of unwrapping an error may itself have an `Unwrap` method; we call
the sequence of errors produced by repeated unwrapping the _error chain_.

### Examining errors with Is and As

The Go 1.13 `errors` package includes two new functions for examining errors: `Is` and `As`.

The `errors.Is` function compares an error to a value.

	// Similar to:
	//   if err == ErrNotFound { … }
	if errors.Is(err, ErrNotFound) {
		// something wasn't found
	}

The `As` function tests whether an error is a specific type.

	// Similar to:
	//   if e, ok := err.(*QueryError); ok { … }
	var e *QueryError
	// Note: *QueryError is the type of the error.
	if errors.As(err, &e) {
		// err is a *QueryError, and e is set to the error's value
	}

In the simplest case, the `errors.Is` function behaves like a comparison to a
sentinel error, and the `errors.As` function behaves like a type assertion. When
operating on wrapped errors, however, these functions consider all the errors in
a chain. Let's look again at the example from above of unwrapping a `QueryError`
to examine the underlying error:

	if e, ok := err.(*QueryError); ok && e.Err == ErrPermission {
		// query failed because of a permission problem
	}

Using the `errors.Is` function, we can write this as:

	if errors.Is(err, ErrPermission) {
		// err, or some error that it wraps, is a permission problem
	}

The `errors` package also includes a new `Unwrap` function which returns the
result of calling an error's `Unwrap` method, or `nil` when the error has no
`Unwrap` method. It is usually better to use `errors.Is` or `errors.As`,
however, since these functions will examine the entire chain in a single call.

Note: although it may feel odd to take a pointer to a pointer, in this case it
is correct. Think of it instead as taking a pointer to a value of the error
type; it so happens in this case that the returned error is a pointer type.

### Wrapping errors with %w

As mentioned earlier, it is common to use the `fmt.Errorf` function to add additional information to an error.

	if err != nil {
		return fmt.Errorf("decompress %v: %v", name, err)
	}

In Go 1.13, the `fmt.Errorf` function supports a new `%w` verb. When this verb
is present, the error returned by `fmt.Errorf` will have an `Unwrap` method
returning the argument of `%w`, which must be an error. In all other ways, `%w`
is identical to `%v`.

	if err != nil {
		// Return an error which unwraps to err.
		return fmt.Errorf("decompress %v: %w", name, err)
	}

Wrapping an error with `%w` makes it available to `errors.Is` and `errors.As`:

	err := fmt.Errorf("access denied: %w", ErrPermission)
	...
	if errors.Is(err, ErrPermission) ...

### Whether to Wrap

When adding additional context to an error, either with `fmt.Errorf` or by
implementing a custom type, you need to decide whether the new error should wrap
the original. There is no single answer to this question; it depends on the
context in which the new error is created. Wrap an error to expose it to
callers. Do not wrap an error when doing so would expose implementation details.

As one example, imagine a `Parse` function which reads a complex data structure
from an `io.Reader`. If an error occurs, we wish to report the line and column
number at which it occurred. If the error occurs while reading from the
`io.Reader`, we will want to wrap that error to allow inspection of the
underlying problem. Since the caller provided the `io.Reader` to the function,
it makes sense to expose the error produced by it.

In contrast, a function which makes several calls to a database probably should
not return an error which unwraps to the result of one of those calls. If the
database used by the function is an implementation detail, then exposing these
errors is a violation of abstraction. For example, if the `LookupUser` function
of your package `pkg` uses Go's `database/sql` package, then it may encounter a
`sql.ErrNoRows` error. If you return that error with
`fmt.Errorf("accessing DB: %v", err)`
then a caller cannot look inside to find the `sql.ErrNoRows`. But if
the function instead returns `fmt.Errorf("accessing DB: %w", err)`, then a
caller could reasonably write

	err := pkg.LookupUser(...)
	if errors.Is(err, sql.ErrNoRows) …

At that point, the function must always return `sql.ErrNoRows` if you don't want
to break your clients, even if you switch to a different database package. In
other words, wrapping an error makes that error part of your API. If you don't
want to commit to supporting that error as part of your API in the future, you
shouldn't wrap the error.

It’s important to remember that whether you wrap or not, the error text will be
the same. A _person_ trying to understand the error will have the same information
either way; the choice to wrap is about whether to give _programs_ additional
information so they can make more informed decisions, or to withhold that
information to preserve an abstraction layer.

## Customizing error tests with Is and As methods

The `errors.Is` function examines each error in a chain for a match with a
target value. By default, an error matches the target if the two are
[equal](/ref/spec#Comparison_operators). In addition, an
error in the chain may declare that it matches a target by implementing an `Is`
_method_.

As an example, consider this error inspired by the
[Upspin error package](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html)
which compares an error against a template, considering only fields which are
non-zero in the template:

	type Error struct {
		Path string
		User string
	}

	func (e *Error) Is(target error) bool {
		t, ok := target.(*Error)
		if !ok {
			return false
		}
		return (e.Path == t.Path || t.Path == "") &&
			   (e.User == t.User || t.User == "")
	}

	if errors.Is(err, &Error{User: "someuser"}) {
		// err's User field is "someuser".
	}

The `errors.As` function similarly consults an `As` method when present.

## Errors and package APIs

A package which returns errors (and most do) should describe what properties of
those errors programmers may rely on. A well-designed package will also avoid
returning errors with properties that should not be relied upon.

The simplest specification is to say that operations either succeed or fail,
returning a nil or non-nil error value respectively. In many cases, no further
information is needed.

If we wish a function to return an identifiable error condition, such as "item
not found," we might return an error wrapping a sentinel.

	var ErrNotFound = errors.New("not found")

	// FetchItem returns the named item.
	//
	// If no item with the name exists, FetchItem returns an error
	// wrapping ErrNotFound.
	func FetchItem(name string) (*Item, error) {
		if itemNotFound(name) {
			return nil, fmt.Errorf("%q: %w", name, ErrNotFound)
		}
		// ...
	}

There are other existing patterns for providing errors which can be semantically
examined by the caller, such as directly returning a sentinel value, a specific
type, or a value which can be examined with a predicate function.

In all cases, care should be taken not to expose internal details to the user.
As we touched on in "Whether to Wrap" above, when you return
an error from another package you should convert the error to a form that does
not expose the underlying error, unless you are willing to commit to returning
that specific error in the future.

	f, err := os.Open(filename)
	if err != nil {
		// The *os.PathError returned by os.Open is an internal detail.
		// To avoid exposing it to the caller, repackage it as a new
		// error with the same text. We use the %v formatting verb, since
		// %w would permit the caller to unwrap the original *os.PathError.
		return fmt.Errorf("%v", err)
	}

If a function is defined as returning an error wrapping some sentinel or type,
do not return the underlying error directly.

	var ErrPermission = errors.New("permission denied")

	// DoSomething returns an error wrapping ErrPermission if the user
	// does not have permission to do something.
	func DoSomething() error {
		if !userHasPermission() {
			// If we return ErrPermission directly, callers might come
			// to depend on the exact error value, writing code like this:
			//
			//     if err := pkg.DoSomething(); err == pkg.ErrPermission { … }
			//
			// This will cause problems if we want to add additional
			// context to the error in the future. To avoid this, we
			// return an error wrapping the sentinel so that users must
			// always unwrap it:
			//
			//     if err := pkg.DoSomething(); errors.Is(err, pkg.ErrPermission) { ... }
			return fmt.Errorf("%w", ErrPermission)
		}
		// ...
	}

## Conclusion

Although the changes we’ve discussed amount to just three functions and a
formatting verb, we hope they will go a long way toward improving how errors are
handled in Go programs. We expect that wrapping to provide additional context
will become commonplace, helping programs to make better decisions and helping
programmers to find bugs more quickly.

As Russ Cox said in his [GopherCon 2019 keynote](https://blog.golang.org/experiment),
on the path to Go 2 we experiment, simplify and ship. Now that we’ve
shipped these changes, we look forward to the experiments that will follow.
