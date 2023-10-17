---
title: When To Use Generics
date: 2022-04-12
by:
- Ian Lance Taylor
tags:
- go2
- generics
summary: When to use generics when writing Go code, and when not to use them.
---

## Introduction

This is the blog post version of my talks at Google Open Source Live:

{{video "https://www.youtube.com/embed/nr8EpUO9jhw"}}

and GopherCon 2021:

{{video "https://www.youtube.com/embed/Pa_e9EeCdy8?start=1244"}}

The Go 1.18 release adds a major new language feature: support for
generic programming.
In this article I'm not going to describe what generics are nor how to
use them.
This article is about when to use generics in Go code, and when not to
use them.

To be clear, I'll provide general guidelines, not hard and fast
rules.
Use your own judgement.
But if you aren't sure, I recommend using the guidelines discussed
here.

## Write code

Let's start with a general guideline for programming Go: write Go
programs by writing code, not by defining types.
When it comes to generics, if you start writing your program by
defining type parameter constraints, you are probably on the wrong
path.
Start by writing functions.
It's easy to add type parameters later when it's clear that they will
be useful.

## When are type parameters useful?

That said, let's look at cases for which type parameters can be
useful.

### When using language-defined container types

One case is when writing functions that operate on the special
container types that are defined by the language: slices, maps, and
channels.
If a function has parameters with those types, and the function code
doesn't make any particular assumptions about the element types, then
it may be useful to use a type parameter.

For example, here is a function that returns a slice of all the keys
in a map of any type:

{{raw `
	// MapKeys returns a slice of all the keys in m.
	// The keys are not returned in any particular order.
	func MapKeys[Key comparable, Val any](m map[Key]Val) []Key {
		s := make([]Key, 0, len(m))
		for k := range m {
			s = append(s, k)
		}
		return s
	}
`}}

This code doesn't assume anything about the map key type, and it
doesn't use the map value type at all.
It works for any map type.
That makes it a good candidate for using type parameters.

The alternative to type parameters for this kind of function is
typically to use reflection, but that is a more awkward programming
model, is not statically typechecked at build time, and is often slower
at run time.

### General purpose data structures

Another case where type parameters can be useful is for general
purpose data structures.
A general purpose data structure is something like a slice or map, but
one that is not built into the language, such as a linked list, or a
binary tree.

Today, programs that need such data structures typically do one of two
things: write them with a specific element type, or use an interface
type.
Replacing a specific element type with a type parameter can produce a
more general data structure that can be used in other parts of the
program, or by other programs.
Replacing an interface type with a type parameter can permit data to
be stored more efficiently, saving memory resources; it can also
permit the code to avoid type assertions, and to be fully type checked
at build time.

For example, here is part of what a binary tree data structure might
look like using type parameters:

{{raw `
	// Tree is a binary tree.
	type Tree[T any] struct {
		cmp  func(T, T) int
		root *node[T]
	}

	// A node in a Tree.
	type node[T any] struct {
		left, right  *node[T]
		val          T
	}

	// find returns a pointer to the node containing val,
	// or, if val is not present, a pointer to where it
	// would be placed if added.
	func (bt *Tree[T]) find(val T) **node[T] {
		pl := &bt.root
		for *pl != nil {
			switch cmp := bt.cmp(val, (*pl).val); {
			case cmp < 0:
				pl = &(*pl).left
		   	case cmp > 0:
				pl = &(*pl).right
			default:
				return pl
			}
		}
		return pl
	}

	// Insert inserts val into bt if not already there,
	// and reports whether it was inserted.
	func (bt *Tree[T]) Insert(val T) bool {
		pl := bt.find(val)
		if *pl != nil {
			return false
		}
		*pl = &node[T]{val: val}
		return true
	}
`}}

Each node in the tree contains a value of the type parameter `T`.
When the tree is instantiated with a particular type argument, values
of that type will be stored directly in the nodes.
They will not be stored as interface types.

This is a reasonable use of type parameters because the `Tree` data
structure, including the code in the methods, is largely independent
of the element type `T`.

The `Tree` data structure does need to know how to compare values of
the element type `T`; it uses a passed-in comparison function for
that.
You can see this on the fourth line of the `find` method, in the call
to `bt.cmp`.
Other than that, the type parameter doesn't matter at all.

### For type parameters, prefer functions to methods

The `Tree` example illustrates another general guideline: when you
need something like a comparison function, prefer a function to a
method.

We could have defined the `Tree` type such that the element type is
required to have a `Compare` or `Less` method.
This would be done by writing a constraint that requires the method,
meaning that any type argument used to instantiate the `Tree` type
would need to have that method.

A consequence would be that anybody who wants to use `Tree` with a
simple data type like `int` would have to define their own integer
type and write their own comparison method.
If we define `Tree` to take a comparison function, as in the code
shown above, then it is easy to pass in the desired function.
It's just as easy to write that comparison function as it is to write
a method.

If the `Tree` element type happens to already have a `Compare` method,
then we can simply use a method expression like `ElementType.Compare`
as the comparison function.

To put it another way, it is much simpler to turn a method into a
function than it is to add a method to a type.
So for general purpose data types, prefer a function rather than
writing a constraint that requires a method.

### Implementing a common method

Another case where type parameters can be useful is when different
types need to implement some common method, and the implementations
for the different types all look the same.

For example, consider the standard library's `sort.Interface`.
It requires that a type implement three methods: `Len`, `Swap`, and
`Less`.

Here is an example of a generic type `SliceFn` that implements
`sort.Interface` for any slice type:

{{raw `
	// SliceFn implements sort.Interface for a slice of T.
	type SliceFn[T any] struct {
		s    []T
		less func(T, T) bool
	}

	func (s SliceFn[T]) Len() int {
		return len(s.s)
	}
	func (s SliceFn[T]) Swap(i, j int) {
		s.s[i], s.s[j] = s.s[j], s.s[i]
	}
	func (s SliceFn[T]) Less(i, j int) bool {
		return s.less(s.s[i], s.s[j])
	}
`}}

For any slice type, the `Len` and `Swap` methods are exactly the same.
The `Less` method requires a comparison, which is the `Fn` part of the
name `SliceFn`.
As with the earlier `Tree` example, we will pass in a function when we
create a `SliceFn`.

Here is how to use `SliceFn` to sort any slice using a comparison
function:

{{raw `
	// SortFn sorts s in place using a comparison function.
	func SortFn[T any](s []T, less func(T, T) bool) {
		sort.Sort(SliceFn[T]{s, less})
	}
`}}

This is similar to the standard library function `sort.Slice`, but the
comparison function is written using values rather than slice
indexes.

Using type parameters for this kind of code is appropriate because the
methods look exactly the same for all slice types.

(I should mention that Go 1.19--not 1.18--will most likely include a
generic function to sort a slice using a comparison function, and that
generic function will most likely not use `sort.Interface`.
See [proposal #47619](https://go.dev/issue/47619).
But the general point is still true even if this specific example will
most likely not be useful: it's reasonable to use type parameters when
you need to implement methods that look the same for all the relevant
types.)

## When are type parameters not useful?

Now let's talk about the other side of the question: when not to use
type parameters.

### Don't replace interface types with type parameters

As we all know, Go has interface types.
Interface types permit a kind of generic programming.

For example, the widely used `io.Reader` interface provides a generic
mechanism for reading data from any value that contains information
(for example, a file) or that produces information (for example, a
random number generator).
If all you need to do with a value of some type is call a method on
that value, use an interface type, not a type parameter.
`io.Reader` is easy to read, efficient, and effective.
There is no need to use a type parameter to read data from a value by
calling the `Read` method.

For example, it might be tempting to change the first function
signature here, which uses just an interface type, into the second
version, which uses a type parameter.

{{raw `
	func ReadSome(r io.Reader) ([]byte, error)

	func ReadSome[T io.Reader](r T) ([]byte, error)
`}}

Don't make that kind of change.
Omitting the type parameter makes the function easier to write, easier
to read, and the execution time will likely be the same.

It's worth emphasizing the last point.
While it's possible to implement generics in several different ways,
and implementations will change and improve over time, the
implementation used in Go 1.18 will in many cases treat values whose
type is a type parameter much like values whose type is an interface
type.
What this means is that using a type parameter will generally not be
faster than using an interface type.
So don't change from interface types to type parameters just for
speed, because it probably won't run any faster.

### Don't use type parameters if method implementations differ

When deciding whether to use a type parameter or an interface type,
consider the implementation of the methods.
Earlier we said that if the implementation of a method is the same for
all types, use a type parameter.
Inversely, if the implementation is different for each type, then use
an interface type and write different method implementations, don't
use a type parameter.

For example, the implementation of `Read` from a file is nothing like
the implementation of `Read` from a random number generator.
That means that we should write two different `Read` methods, and
use an interface type like `io.Reader`.

### Use reflection where appropriate

Go has [run time reflection](https://pkg.go.dev/reflect).
Reflection permits a kind of generic programming, in that it permits
you to write code that works with any type.

If some operation has to support even types that don't have methods
(so that interface types don't help), and if the operation is
different for each type (so that type parameters aren't appropriate),
use reflection.

An example of this is the
[encoding/json](https://pkg.go.dev/encoding/json) package.
We don't want to require that every type that we encode have a
`MarshalJSON` method, so we can't use interface types.
But encoding an interface type is nothing like encoding a struct type,
so we shouldn't use type parameters.
Instead, the package uses reflection.
The code is not simple, but it works.
For details, see [the source
code](https://go.dev/src/encoding/json/encode.go).

## One simple guideline

In closing, this discussion of when to use generics can be reduced to
one simple guideline.

If you find yourself writing the exact same code multiple times, where
the only difference between the copies is that the code uses different
types, consider whether you can use a type parameter.

Another way to say this is that you should avoid type parameters until
you notice that you are about to write the exact same code multiple
times.
