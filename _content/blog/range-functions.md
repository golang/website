---
title: Range Over Function Types
date: 2024-08-20
by:
- Ian Lance Taylor
tags:
- iterators
summary: A description of range over function types, a new feature in Go 1.23.
---

## Introduction

This is the blog post version of my talk at GopherCon 2024.

{{video "https://www.youtube.com/embed/i9zwUT9dlVc"}}

Range over function types is a new language feature in the Go 1.23
release.
This blog post will explain why we are adding this new feature, what
exactly it is, and how to use it.

## Why?

Since Go 1.18 we've had the ability to write new generic container
types in Go.
For example, let's consider this very simple `Set` type, a generic
type implemented on top of a map.

```
// Set holds a set of elements.
type Set[E comparable] struct {
	m map[E]struct{}
}

// New returns a new [Set].
func New[E comparable]() *Set[E] {
	return &Set[E]{m: make(map[E]struct{})}
}
```

Naturally a set type has a way to add elements and a way to check
whether elements are present.  The details here don't matter.

```
// Add adds an element to a set.
func (s *Set[E]) Add(v E) {
	s.m[v] = struct{}{}
}

// Contains reports whether an element is in a set.
func (s *Set[E]) Contains(v E) bool {
	_, ok := s.m[v]
	return ok
}
```

And among other things we will want a function to return the union of
two sets.

```
// Union returns the union of two sets.
func Union[E comparable](s1, s2 *Set[E]) *Set[E] {
	r := New[E]()
	// Note for/range over internal Set field m.
	// We are looping over the maps in s1 and s2.
	for v := range s1.m {
		r.Add(v)
	}
	for v := range s2.m {
		r.Add(v)
	}
	return r
}
```

Let's look at this implementation of the `Union` function for a
minute.
In order to compute the union of two sets, we need a way to get all
the elements that are in each set.
In this code we use a for/range statement over an unexported field of
the set type.
That only works if the `Union` function is defined in the set package.

But there are a lot of reasons why someone might want to loop over all
the elements in a set.
This set package has to provide some way for its users to do that.

How should that work?

### Push Set elements

One approach is to provide a `Set` method that takes a function, and
to call that function with every element in the Set.
We'll call this `Push`, because the `Set` pushes every value to the
function.
Here if the function returns false, we stop calling it.

```
func (s *Set[E]) Push(f func(E) bool) {
	for v := range s.m {
		if !f(v) {
			return
		}
	}
}
```

In the Go standard library, we see this general pattern used for cases
like the [`sync.Map.Range`](https://pkg.go.dev/sync#Map.Range) method,
the [`flag.Visit`](https://pkg.go.dev/flag#Visit) function, and the
[`filepath.Walk`](https://pkg.go.dev/path/filepath#Walk) function.
This is a general pattern, not an exact one; as it happens, none of
those three examples work quite the same way.

This is what it looks like to use the `Push` method to print all the
elements of a set: you call `Push` with a function that does what you
want with the element.

```
func PrintAllElementsPush[E comparable](s *Set[E]) {
	s.Push(func(v E) bool {
		fmt.Println(v)
		return true
	})
}
```

### Pull Set elements

Another approach to looping over the elements of a `Set` is to return
a function.
Each time the function is called, it will return a value from the
`Set`, along with a boolean that reports whether the value is valid.
The boolean result will be false when the loop has gone through all
the elements.
In this case we also need a stop function that can be called when no
more values are needed.

This implementation uses a pair of channels, one for values in the
set and one to stop returning values.
We use a goroutine to send values on the channel.
The `next` function returns an element from the set by reading from
the element channel, and the `stop` function tells the goroutine to
exit by closing the stop channel.
We need the `stop` function to make sure that the goroutine exits when
no more values are needed.

{{raw `
	// Pull returns a next function that returns each
	// element of s with a bool for whether the value
	// is valid. The stop function should be called
	// when finished calling the next function.
	func (s *Set[E]) Pull() (func() (E, bool), func()) {
		ch := make(chan E)
		stopCh := make(chan bool)

		go func() {
			defer close(ch)
			for v := range s.m {
				select {
				case ch <- v:
				case <-stopCh:
					return
				}
			}
		}()

		next := func() (E, bool) {
			v, ok := <-ch
			return v, ok
		}

		stop := func() {
			close(stopCh)
		}

		return next, stop
	}
`}}

Nothing in the standard library works exactly this way.  Both
[`runtime.CallersFrames`](https://pkg.go.dev/runtime#CallersFrames)
and
[`reflect.Value.MapRange`](https://pkg.go.dev/reflect#Value.MapRange)
are similar, though they return values with methods rather than
returning functions directly.

This is what it looks like to use the `Pull` method to print all the
elements of a `Set`.
You call `Pull` to get a function, and you repeatedly call that
function in a for loop.

```
func PrintAllElementsPull[E comparable](s *Set[E]) {
	next, stop := s.Pull()
	defer stop()
	for v, ok := next(); ok; v, ok = next() {
		fmt.Println(v)
	}
}
```

## Standardize the approach

We've now seen two different approaches to looping over all the
elements of a set.
Different Go packages use these approaches and several others.
That means that when you start using a new Go container package you
may have to learn a new looping mechanism.
It also means that we can't write one function that works with several
different types of containers, as the container types will handle
looping differently.

We want to improve the Go ecosystem by developing standard approaches
for looping over containers.

### Iterators

This is, of course, an issue that arises in many programming
languages.

The popular [Design Patterns
book](https://en.wikipedia.org/wiki/Design_Patterns), first published
in 1994, describes this as the iterator pattern.
You use an iterator to "provide a way to access the elements of an
aggregate object sequentially without exposing its underlying
representation."
What this quote calls an aggregate object is what I've been calling a
container.
An aggregate object, or container, is just a value that holds other
values, like the `Set` type we've been discussing.

Like many ideas in programming, iterators date back to Barbara
Liskov's [CLU language](https://en.wikipedia.org/wiki/CLU_(programming_language)),
developed in the 1970's.

Today many popular languages provide iterators one way or another,
including, among others, C++, Java, Javascript, Python, and Rust.

However, Go before version 1.23 did not.

### For/range

As we all know, Go has container types that are built in to the
language: slices, arrays, and maps.
And it has a way to access the elements of those values without
exposing the underlying representation: the for/range statement.
The for/range statement works for Go's built-in container types (and
also for strings, channels, and, as of Go 1.22, int).

The for/range statement is iteration, but it is not iterators as they
appear in today's popular languages.
Still, it would be nice to be able to use for/range to iterate over a
user-defined container like the `Set` type.

However, Go before version 1.23 did not support this.

### Improvements in this release

For Go 1.23 we've decided to support both for/range over user-defined
container types, and a standardized form of iterators.

We extended the for/range statement to support ranging over function
types.
We'll see below how this helps loop over user-defined containers.

We also added standard library types and functions to support using
function types as iterators.
A standard definition of iterators lets us write functions that work
smoothly with different container types.

### Range over (some) function types

The improved for/range statement doesn't support arbitrary function
types.
As of Go 1.23 it now supports ranging over functions that take a
single argument.
The single argument must itself be a function that takes zero to two
arguments and returns a bool; by convention, we call it the yield
function.

```
func(yield func() bool)

func(yield func(V) bool)

func(yield func(K, V) bool)
```

When we speak of an iterator in Go, we mean a function with one of
these three types.
As we'll discuss below, there is another kind of iterator in the
standard library: a pull iterator.
When it is necessary to distinguish between standard iterators and
pull iterators, we call the standard iterators push iterators.
That is because, as we will see, they push out a sequence of values by
calling a yield function.

### Standard (push) iterators

To make iterators easier to use, the new standard library package iter
defines two types: `Seq` and `Seq2`.
These are names for the iterator function types, the types that can be
used with the for/range statement.
The name `Seq` is short for sequence, as iterators loop through a
sequence of values.

```
package iter

type Seq[V any] func(yield func(V) bool)

type Seq2[K, V any] func(yield func(K, V) bool)

// for now, no Seq0
```

The difference between `Seq` and `Seq2` is just that `Seq2` is a
sequence of pairs, such as a key and a value from a map.
In this post we'll focus on `Seq` for simplicity, but most of what we
say covers `Seq2` as well.

It's easiest to explain how iterators work with an example.
Here the `Set` method `All` returns a function.
The return type of `All` is `iter.Seq[E]`, so we know that it returns
an iterator.

```
// All is an iterator over the elements of s.
func (s *Set[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
	}
}
```

The iterator function itself takes another function, the yield
function, as an argument.
The iterator calls the yield function with every value in the set.
In this case the iterator, the function returned by `Set.All`, is a
lot like the `Set.Push` function we saw earlier.

This shows how iterators work: for some sequence of values, they call
a yield function with each value in the sequence.
If the yield function returns false, no more values are needed, and
the iterator can just return, doing any cleanup that may be required.
If the yield function never returns false, the iterator can just
return after calling yield with all the values in the sequence.

That's how they work, but let's acknowledge that the first time you
see one of these, your first reaction is probably "there are a lot of
functions flying around here."
You're not wrong about that.
Let's focus on two things.

The first is that once you get past the first line of this function's
code, the actual implementation of the iterator is pretty simple: call
yield with every element of the set, stopping if yield returns false.

```
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
```

The second is that using this is really easy.
You call `s.All` to get an iterator, and then you use for/range to
loop over all the elements in `s`.
The for/range statement supports any iterator, and this shows how easy
that is to use.

```
func PrintAllElements[E comparable](s *Set[E]) {
	for v := range s.All() {
		fmt.Println(v)
	}
}
```

In this kind of code `s.All` is a method that returns a function.
We are calling `s.All`, and then using for/range to range over the
function that it returns.
In this case we could have made `Set.All` be an iterator function
itself, rather than having it return an iterator function.
However, in some cases that won't work, such as if the function that
returns the iterator needs to take an argument, or needs to do some
set up work.
As a matter of convention, we encourage all container types to provide
an `All` method that returns an iterator, so that programmers don't
have to remember whether to range over `All` directly or whether to
call `All` to get a value they can range over.
They can always do the latter.

If you think about it, you'll see that the compiler must be adjusting
the loop to create a yield function to pass to the iterator returned
by `s.All`.
There's a fair bit of complexity in the Go compiler and runtime to
make this efficient, and to correctly handle things like `break` or
`panic` in the loop.
We're not going to cover any of that in this blog post.
Fortunately the implementation details are not important when it comes
to actually using this feature.

### Pull iterators

We've now seen how to use iterators in a for/range loop.
But a simple loop is not the only way to use an iterator.
For example, sometimes we may need to iterate over two containers in
parallel.
How do we do that?

The answer is that we use a different kind of iterator: a pull
iterator.
We've seen that a standard iterator, also known as a push iterator, is
a function that takes a yield function as an argument and pushes each
value in a sequence by calling the yield function.

A pull iterator works the other way around: it is a function that is
written such that each time you call it, it returns the next value in
the sequence.

We'll repeat the difference between the two types of iterators to help
you remember:
- A push iterator pushes each value in a sequence to a yield
  function.
  Push iterators are standard iterators in the Go standard library,
  and are supported directly by the for/range statement.
- A pull iterator works the other way around.
  Each time you call a pull iterator, it pulls another value from a
  sequence and returns it.
  Pull iterators are _not_ supported directly by the for/range
  statement; however, it's straightforward to write an ordinary for
  statement that loops through a pull iterator.
  In fact, we saw an example earlier when we looked at using the
  `Set.Pull` method.

You could write a pull iterator yourself, but normally you don't have
to.
The new standard library function
[`iter.Pull`](https://pkg.go.dev/iter#Pull) takes a standard iterator,
that is to say a function that is a push iterator, and returns a pair
of functions.
The first is a pull iterator: a function that returns the next value
in the sequence each time it is called.
The second is a stop function that should be called when we are done
with the pull iterator.
This is like the `Set.Pull` method we saw earlier.

The first function returned by `iter.Pull`, the pull iterator, returns
a value and a boolean that reports whether that value is valid.
The boolean will be false at the end of the sequence.

`iter.Pull` returns a stop function in case we don't read through the
sequence to the end.
In the general case the push iterator, the argument to `iter.Pull`,
may start goroutines, or build new data structures that need to be
cleaned up when iteration is complete.
The push iterator will do any cleanup when the yield function returns
false, meaning that no more values are required.
When used with a for/range statement, the for/range statement will
ensure that if the loop exits early, through a `break` statement or
for any other reason, then the yield function will return false.
With a pull iterator, on the other hand, there is no way to force the
yield function to return false, so the stop function is needed.

Another way to say this is that calling the stop function will cause
the yield function to return false when it is called by the push
iterator.

Strictly speaking you don't need to call the stop function if the pull
iterator returns false to indicate that it has reached the end of the
sequence, but it's usually simpler to just always call it.

Here is an example of using pull iterators to walk through two
sequences in parallel.
This function reports whether two arbitrary sequences contain the same
elements in the same order.

```
// EqSeq reports whether two iterators contain the same
// elements in the same order.
func EqSeq[E comparable](s1, s2 iter.Seq[E]) bool {
	next1, stop1 := iter.Pull(s1)
	defer stop1()
	next2, stop2 := iter.Pull(s2)
	defer stop2()
	for {
		v1, ok1 := next1()
		v2, ok2 := next2()
		if !ok1 {
			return !ok2
		}
		if ok1 != ok2 || v1 != v2 {
			return false
		}
	}
}
```

The function uses `iter.Pull` to convert the two push iterators, `s1`
and `s2`, into pull iterators.
It uses `defer` statements to make sure that the pull iterators are
stopped when we are done with them.

Then the code loops, calling the pull iterators to retrieve values.
If the first sequence is done, it returns true if the second sequence
is also done, or false if it isn't.
If the values are different, it returns false.
Then it loops to pull the next two values.

As with push iterators, there is some complexity in the Go runtime to
make pull iterators efficient, but this does not affect code that
actually uses the `iter.Pull` function.

## Iterating on iterators

Now you know everything there is to know about range over function
types and about iterators.
We hope you enjoy using them!

Still, there are a few more things worth mentioning.

### Adapters

An advantage of a standard definition of iterators is the ability to
write standard adapter functions that use them.

For example, here is a function that filters a sequence of values,
returning a new sequence.
This `Filter` function takes an iterator as an argument and returns a
new iterator.
The other argument is a filter function that decides which values
should be in the new iterator that `Filter` returns.

```
// Filter returns a sequence that contains the elements
// of s for which f returns true.
func Filter[V any](f func(V) bool, s iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range s {
			if f(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}
```

As with the earlier example, the function signatures look complicated
when you first see them.
Once you get past the signatures, the implementation is
straightforward.

```
		for v := range s {
			if f(v) {
				if !yield(v) {
					return
				}
			}
		}
```

The code ranges over the input iterator, checks the filter function,
and calls yield with the values that should go into the output
iterator.

We'll show an example of using `Filter` below.

(There is no version of `Filter` in the Go standard library today, but
one may be added in future releases.)

### Binary tree

As an example of how convenient a push iterator can be to loop over a
container type, let's consider this simple binary tree type.

```
// Tree is a binary tree.
type Tree[E any] struct {
	val         E
	left, right *Tree[E]
}
```

We won't show the code to insert values into the tree, but naturally
there should be some way to range over all the values in the tree.

It turns out that the iterator code is easier to write if it returns a
bool.
Since the function types supported by for/range don't return anything,
the `All` method here returns a small function literal that calls the
iterator itself, here called `push`, and ignores the bool result.

```
// All returns an iterator over the values in t.
func (t *Tree[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		t.push(yield)
	}
}

// push pushes all elements to the yield function.
func (t *Tree[E]) push(yield func(E) bool) bool {
	if t == nil {
		return true
	}
	return t.left.push(yield) &&
		yield(t.val) &&
		t.right.push(yield)
}
```

The `push` method uses recursion to walk over the whole tree, calling
yield on each element.
If the yield function returns false, the method returns false all the
way up the stack.
Otherwise it just returns once the iteration is complete.

This shows how straightforward it is to use this iterator approach to
loop over even complex data structures.
There is no need to maintain a separate stack to record the position
within the tree; we can just use the goroutine call stack to do that
for us.

### New iterator functions.

Also new in Go 1.23 are functions in the slices and maps packages that
work with iterators.

Here are the new functions in the slices package.
`All` and `Values` are functions that return iterators over the
elements of a slice.
`Collect` fetches the values out of an iterator and returns a slice
holding those values.
See the docs for the others.

- [`All([]E) iter.Seq2[int, E]`](https://pkg.go.dev/slices#All)
- [`Values([]E) iter.Seq[E]`](https://pkg.go.dev/slices#Values)
- [`Collect(iter.Seq[E]) []E`](https://pkg.go.dev/slices#Collect)
- [`AppendSeq([]E, iter.Seq[E]) []E`](https://pkg.go.dev/slices#AppendSeq)
- [`Backward([]E) iter.Seq2[int, E]`](https://pkg.go.dev/slices#Backward)
- [`Sorted(iter.Seq[E]) []E`](https://pkg.go.dev/slices#Sorted)
- [`SortedFunc(iter.Seq[E], func(E, E) int) []E`](https://pkg.go.dev/slices#SortedFunc)
- [`SortedStableFunc(iter.Seq[E], func(E, E) int) []E`](https://pkg.go.dev/slices#SortedStableFunc)
- [`Repeat([]E, int) []E`](https://pkg.go.dev/slices#Repeat)
- [`Chunk([]E, int) iter.Seq([]E)`](https://pkg.go.dev/slices#Chunk)

Here are the new functions in the maps package.
`All`, `Keys`, and `Values` returns iterators over the map contents.
`Collect` fetches the keys and values out of an iterator and returns a
new map.

- [`All(map[K]V) iter.Seq2[K, V]`](https://pkg.go.dev/maps#All)
- [`Keys(map[K]V) iter.Seq[K]`](https://pkg.go.dev/maps#Keys)
- [`Values(map[K]V) iter.Seq[V]`](https://pkg.go.dev/maps#Values)
- [`Collect(iter.Seq2[K, V]) map[K, V]`](https://pkg.go.dev/maps#Collect)
- [`Insert(map[K, V], iter.Seq2[K, V])`](https://pkg.go.dev/maps#Insert)

### Standard library iterator example

Here is an example of how you might use these new functions along with
the `Filter` function we saw earlier.
This function takes a map from int to string and returns a slice
holding just the values in the map that are longer than some argument `n`.

```
// LongStrings returns a slice of just the values
// in m whose length is n or more.
func LongStrings(m map[int]string, n int) []string {
	isLong := func(s string) bool {
		return len(s) >= n
	}
	return slices.Collect(Filter(isLong, maps.Values(m)))
}
```

The `maps.Values` function returns an iterator over the values in `m`.
`Filter` reads that iterator and returns a new iterator that only
contains the long strings.
`slices.Collect` reads from that iterator into a new slice.

Of course, you could write a loop to do this easily enough, and in
many cases a loop will be clearer.
We don't want to encourage everybody to write code in this style all
the time.
That said, the advantage of using iterators is that this kind of
function works the same way with any sequence.
In this example, notice how Filter is using a map as an input and a
slice as an output, without having to change the code in Filter at
all.

### Looping over lines in a file

Although most of the examples we've seen have involved containers,
iterators are flexible.

Consider this simple code, which doesn't use iterators, to loop over
the lines in a byte slice.
This is easy to write and fairly efficient.

```
	nl := []byte{'\n'}
	// Trim a trailing newline to avoid a final empty blank line.
	for _, line := range bytes.Split(bytes.TrimSuffix(data, nl), nl) {
		handleLine(line)
	}
```

However, `bytes.Split` does allocate and return a slice of byte slices
to hold the lines.
The garbage collector will have to do a bit of work to eventually free
that slice.

Here is a function that returns an iterator over the lines of some
byte slice.
After the usual iterator signatures, the function is pretty simple.
We keep picking lines out of data until there is nothing left, and we
pass each line to the yield function.

```
// Lines returns an iterator over lines in data.
func Lines(data []byte) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for len(data) > 0 {
			line, rest, _ := bytes.Cut(data, []byte{'\n'})
			if !yield(line) {
				return
			}
			data = rest
		}
	}
}
```

Now our code to loop over the lines of a byte slice looks like this.

```
	for line := range Lines(data) {
		handleLine(line)
	}
```

This is just as easy to write as the earlier code, and is a bit more
efficient because it doesn't have allocate a slice of lines.

### Passing a function to a push iterator

For our final example, we'll see that you don't have to use a push
iterator in a range statement.

Earlier we saw a `PrintAllElements` function that prints out each
element of a set.
Here is another way to print all the elements of a set: call `s.All`
to get an iterator, then pass in a hand-written yield function.
This yield function just prints a value and returns true.
Note that there are two function calls here: we call `s.All` to get an
iterator which is itself a function, and we call that function with
our hand-written yield function.


```
func PrintAllElements[E comparable](s *Set[E]) {
	s.All()(func(v E) bool {
		fmt.Println(v)
		return true
	})
}
```

There's no particular reason to write this code this way.
This is just an example to show that the yield function isn't magic.
It can be any function you like.

## Update go.mod

A final note: every Go module specifies the language version that it
uses.
That means that in order to use new language features in an existing
module you may need to update that version.
This is true for all new language features; it's not something
specific to range over function types.
As range over function types is new in the Go 1.23 release, using it
requires specifying at least Go language version 1.23.

There are (at least) four ways to set the language version:
- On the command line, run `go get go@1.23` (or `go mod edit -go=1.23`
  to only edit the `go` directive).
- Manually edit the `go.mod` file and change the `go` line.
- Keep the older language version for the module as a whole, but use a
  `//go:build go1.23` build tag to permit using range over function
  types in a specific file.
