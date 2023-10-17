---
title: Why Generics?
date: 2019-07-31
by:
- Ian Lance Taylor
tags:
- go2
- proposals
- generics
summary: Why should we add generics to Go, and what might they look like?
---

## Introduction

This is the blog post version of my talk last week at Gophercon 2019.

{{video "https://www.youtube.com/embed/WzgLqE-3IhY?rel=0"}}

This article is about what it would mean to add generics to Go, and
why I think we should do it.
I'll also touch on an update to a possible design for
adding generics to Go.

Go was released on November 10, 2009.
Less than 24 hours later we saw the
[first comment about generics](https://groups.google.com/d/msg/golang-nuts/70-pdwUUrbI/onMsQspcljcJ).
(That comment also mentions exceptions, which we added to the
language, in the form of `panic` and `recover`, in early 2010.)

In three years of Go surveys, lack of generics has always been listed
as one of the top three problems to fix in the language.

## Why generics?

But what does it mean to add generics, and why would we want it?

To paraphrase
[Jazayeri, et al](https://www.dagstuhl.de/en/program/calendar/semhp/?semnr=98171):
generic programming enables the representation of functions and data
structures in a generic form, with types factored out.

What does that mean?

For a simple example, let's assume we want to reverse the elements in
a slice.
It's not something that many programs need to do, but it's
not all that unusual.

Let's say it's a slice of int.

{{raw `
	func ReverseInts(s []int) {
		first := 0
		last := len(s)
		for first < last {
			s[first], s[last] = s[last], s[first]
			first++
			last--
		}
	}
`}}

Pretty simple, but even for a simple function like that you'd want to
write a few test cases.
In fact, when I did, I found a bug.
I'm sure many readers have spotted it already.

{{raw `
	func ReverseInts(s []int) {
		first := 0
		last := len(s) - 1
		for first < last {
			s[first], s[last] = s[last], s[first]
			first++
			last--
		}
	}
`}}

We need to subtract 1 when we set the variable last.

Now let's reverse a slice of string.

{{raw `
	func ReverseStrings(s []string) {
		first := 0
		last := len(s) - 1
		for first < last {
			s[first], s[last] = s[last], s[first]
			first++
			last--
		}
	}
`}}

If you compare `ReverseInts` and `ReverseStrings`, you'll see that the
two functions are exactly the same, except for the type of the parameter.
I don't think any reader is surprised by that.

What some people new to Go find surprising is that there is no way to
write a simple `Reverse` function that works for a slice of any type.

Most other languages do let you write that kind of function.

In a dynamically typed language like Python or JavaScript you can
simply write the function, without bothering to specify the element
type.  This doesn't work in Go because Go is statically typed, and
requires you to write down the exact type of the slice and the type of
the slice elements.

Most other statically typed languages, like C++ or Java or Rust or
Swift, support generics to address exactly this kind of issue.

## Go generic programming today

So how do people write this kind of code in Go?

In Go you can write a single function that works for different slice
types by using an interface type, and defining a method on the slice
types you want to pass in.
That is how the standard library's `sort.Sort` function works.

In other words, interface types in Go are a form of generic
programming.
They let us capture the common aspects of different types and express
them as methods.
We can then write functions that use those interface types, and those
functions will work for any type that implements those methods.

But this approach falls short of what we want.
With interfaces you have to write the methods yourself.
It's awkward to have to define a named type with a couple of methods
just to reverse a slice.
And the methods you write are exactly the same for each slice type, so
in a sense we've just moved and condensed the duplicate code, we
haven't eliminated it.
Although interfaces are a form of generics, they don’t give us
everything we want from generics.

A different way of using interfaces for generics, which could get around
the need to write the methods yourself, would be to have the language
define methods for some kinds of types.
That isn't something the language supports today, but, for example,
the language could define that every slice type has an Index method
that returns an element.
But in order to use that method in practice it would have to return an
empty interface type, and then we lose all the benefits of static
typing.
More subtly, there would be no way to define a generic function that
takes two different slices with the same element type, or that takes a
map of one element type and returns a slice of the same element type.
Go is a statically typed language because that makes it easier to
write large programs; we don’t want to lose the benefits of static
typing in order to gain the benefits of generics.

Another approach would be to write a generic `Reverse` function using
the reflect package, but that is so awkward to write and slow to run
that few people do that.
That approach also requires explicit type assertions and has no static
type checking.

Or, you could write a code generator that takes a type and generates a
`Reverse` function for slices of that type.
There are several code generators out there that do just that.
But this adds another step to every package that needs `Reverse`,
it complicates the build because all the different copies have to be
compiled, and fixing a bug in the master source requires re-generating
all the instances, some of which may be in different projects
entirely.

All these approaches are awkward enough that I think most
people who have to reverse a slice in Go just write the function for
the specific slice type that they need.
Then they'll need to write test cases for the function, to make sure
they didn't make a simple mistake like the one I made initially.
And they'll need to run those tests routinely.

However we do it, it means a lot of extra work just for a function that
looks exactly the same except for the element type.
It's not that it can't be done.
It clearly can be done, and Go programmers are doing it.
It's just that there ought to be a better way.

For a statically typed language like Go, that better way is generics.
What I wrote earlier is that generic programming enables the
representation of functions and data structures in a generic form,
with types factored out.
That's exactly what we want here.

## What generics can bring to Go

The first and most important thing we want from generics in Go is to
be able to write functions like `Reverse` without caring about the
element type of the slice.
We want to factor out that element type.
Then we can write the function once, write the tests once, put them in
a go-gettable package, and call them whenever we want.

Even better, since this is an open source world, someone else can
write `Reverse` once, and we can use their implementation.

At this point I should say that “generics” can mean a lot of different
things.
In this article, what I mean by “generics” is what I just described.
In particular, I don’t mean templates as found in the C++ language,
which support quite a bit more than what I’ve written here.

I went through `Reverse` in detail, but there are many other functions
that we could write generically, such as:

  - Find smallest/largest element in slice
  - Find average/standard deviation of slice
  - Compute union/intersection of maps
  - Find shortest path in node/edge graph
  - Apply transformation function to slice/map, returning new slice/map

These examples are available in most other languages.
In fact, I wrote this list by glancing at the C++ standard template
library.

There are also examples that are specific to Go with its strong
support for concurrency.

  - Read from a channel with a timeout
  - Combine two channels into a single channel
  - Call a list of functions in parallel, returning a slice of results
  - Call a list of functions, using a Context, return the result of the first function to finish, canceling and cleaning up extra goroutines

I've seen all of these functions written out many times with different
types.
It's not hard to write them in Go.
But it would be nice to be able to reuse an efficient and debugged
implementation that works for any value type.

To be clear, these are just examples.
There are many more general purpose functions that could be written
more easily and safely using generics.

Also, as I wrote earlier, it's not just functions.
It's also data structures.

Go has two general purpose generic data structures built into the
language: slices and maps.
Slices and maps can hold values of any data type, with static type
checking for values stored and retrieved.
The values are stored as themselves, not as interface types.
That is, when I have a `[]int`, the slice holds ints directly, not
ints converted to an interface type.

Slices and maps are the most useful generic data structures, but they
aren’t the only ones.
Here are some other examples.

  - Sets
  - Self-balancing trees, with efficient insertion and traversal in sorted order
  - Multimaps, with multiple instances of a key
  - Concurrent hash maps, supporting parallel insertions and lookups with no single lock

If we can write generic types, we can define new data structures, like
these, that have the same type-checking advantages as slices and maps:
the compiler can statically type-check the types of the values that
they hold, and the values can be stored as themselves, not as
interface types.

It should also be possible to take algorithms like the ones mentioned
earlier and apply them to generic data structures.

These examples should all be just like `Reverse`: generic functions
and data structures written once, in a package, and reused whenever
they are needed.
They should work like slices and maps, in that they shouldn't store
values of empty interface type, but should store specific types, and
those types should be checked at compile time.

So that's what Go can gain from generics.
Generics can give us powerful building blocks that let us share code
and build programs more easily.

I hope I’ve explained why this is worth looking into.

## Benefits and costs

But generics don't come from the
[Big Rock Candy Mountain](https://mainlynorfolk.info/folk/songs/bigrockcandymountain.html),
the land where the sun shines every day over the
[lemonade springs](http://www.lat-long.com/Latitude-Longitude-773297-Montana-Lemonade_Springs.html).
Every language change has a cost.
There's no doubt that adding generics to Go will make the language
more complicated.
As with any change to the language, we need to talk about maximizing
the benefit and minimizing the cost.

In Go, we’ve aimed to reduce complexity through independent, orthogonal
language features that can be combined freely.
We reduce complexity by making the individual features simple, and we
maximize the benefit of the features by permitting their free
combination.
We want to do the same with generics.

To make this more concrete I’m going to list a few guidelines we
should follow.

### Minimize new concepts

We should add as few new concepts to the language as possible.
That means a minimum of new syntax and a minimum of new keywords and
other names.

### Complexity falls on the writer of generic code, not the user

As much as possible the complexity should fall on the programmer
writing the generic package.
We don't want the user of the package to have to worry about generics.
This means that it should be possible to call generic functions in a
natural way, and it means that any errors in using a generic package
should be reported in a way that is easy to understand and to fix.
It should also be easy to debug calls into generic code.

### Writer and user can work independently

Similarly, we should make it easy to separate the concerns of the
writer of the generic code and its user, so that they can develop their
code independently.
They shouldn't have to worry about what the other is doing, any more
than the writer and caller of a normal function in different packages
have to worry.
This sounds obvious, but it's not true of generics in every other
programming language.

### Short build times, fast execution times

Naturally, as much as possible, we want to keep the short build times
and fast execution time that Go gives us today.
Generics tend to introduce a tradeoff between fast builds and fast
execution.
As much as possible, we want both.

### Preserve clarity and simplicity of Go

Most importantly, Go today is a simple language.
Go programs are usually clear and easy to understand.
A major part of our long process of exploring this space has been
trying to understand how to add generics while preserving that clarity
and simplicity.
We need to find mechanisms that fit well into the existing language,
without turning it into something quite different.

These guidelines should apply to any generics implementation in Go.
That’s the most important message I want to leave you with today:
**generics can bring a significant benefit to the language, but they are only worth doing if Go still feels like Go**.

## Draft design

Fortunately, I think it can be done.
To finish up this article I’m going to shift from discussing why we
want generics, and what the requirements on them are, to briefly
discuss a design for how we think we can add them to the language.

Note added January 2022: This blog post was written in 2019 and does
not describe the version of generics that was finally adopted.
For updated information please see the description of type parameters in
[the language spec](https://go.dev/ref/spec) and
[the generics design document](https://go.dev/design/43651-type-parameters).

At this year's Gophercon Robert Griesemer and I published
[a design draft](https://github.com/golang/proposal/blob/master/design/go2draft-contracts.md)
for adding generics to Go.
See the draft for full details.
I'll go over some of the main points here.

Here is the generic Reverse function in this design.

{{raw `
	func Reverse (type Element) (s []Element) {
		first := 0
		last := len(s) - 1
		for first < last {
			s[first], s[last] = s[last], s[first]
			first++
			last--
		}
	}
`}}

You'll notice that the body of the function is exactly the same.
Only the signature has changed.

The element type of the slice has been factored out.
It's now named `Element` and has become what we call a
_type parameter_.
Instead of being part of the type of the slice parameter, it's now a
separate, additional, type parameter.

To call a function with a type parameter, in the general case you pass
a type argument, which is like any other argument except that it's a
type.

	func ReverseAndPrint(s []int) {
		Reverse(int)(s)
		fmt.Println(s)
	}

That is the `(int)` seen after `Reverse` in this example.

Fortunately, in most cases, including this one, the compiler can
deduce the type argument from the types of the regular arguments, and
you don't need to mention the type argument at all.

Calling a generic function just looks like calling any other function.

	func ReverseAndPrint(s []int) {
		Reverse(s)
		fmt.Println(s)
	}

In other words, although the generic `Reverse` function is slightly
more complex than `ReverseInts` and `ReverseStrings`, that complexity
falls on the writer of the function, not the caller.

### Contracts

Since Go is a statically typed language, we have to talk about the
type of a type parameter.
This _meta-type_ tells the compiler what sorts of type arguments are
permitted when calling a generic function, and what sorts of
operations the generic function can do with values of the type
parameter.

The `Reverse` function can work with slices of any type.
The only thing it does with values of type `Element` is assignment,
which works with any type in Go.
For this kind of generic function, which is a very common case, we
don't need to say anything special about the type parameter.

Let's take a quick look at a different function.

{{raw `
	func IndexByte (type T Sequence) (s T, b byte) int {
		for i := 0; i < len(s); i++ {
			if s[i] == b {
				return i
			}
		}
		return -1
	}
`}}

Currently both the bytes package and the strings package in the
standard library have an `IndexByte` function.
This function returns the index of `b` in the sequence `s`, where `s`
is either a `string` or a `[]byte`.
We could use this single generic function to replace the two functions
in the bytes and strings packages.
In practice we may not bother doing that, but this is a useful simple
example.

Here we need to know that the type parameter `T` acts like a `string`
or a `[]byte`.
We can call `len` on it, and we can index to it, and we can compare
the result of the index operation to a byte value.

To let this compile, the type parameter `T` itself needs a type.
It's a meta-type, but because we sometimes need to describe multiple
related types, and because it describes a relationship between the
implementation of the generic function and its callers, we actually
call the type of `T` a contract.
Here the contract is named `Sequence`.
It appears after the list of type parameters.

This is how the Sequence contract is defined for this example.

	contract Sequence(T) {
		T string, []byte
	}

It's pretty simple, since this is a simple example: the type parameter
`T` can be either `string` or `[]byte`.
Here `contract` may be a new keyword, or a special identifier
recognized in package scope; see the design draft for details.

Anybody who remembers [the design we presented at Gophercon 2018](https://github.com/golang/proposal/blob/4a530dae40977758e47b78fae349d8e5f86a6c0a/design/go2draft-contracts.md)
will see that this way of writing a contract is a lot simpler.
We got a lot of feedback on that earlier design that contracts were
too complicated, and we've tried to take that into account.
The new contracts are much simpler to write, and to read, and to
understand.

They let you specify the underlying type of a type parameter, and/or
list the methods of a type parameter.
They also let you describe the relationship between different type
parameters.

### Contracts with methods

Here is another simple example, of a function that uses the String
method to return a `[]string` of the string representation of all the
elements in `s`.

	func ToStrings (type E Stringer) (s []E) []string {
		r := make([]string, len(s))
		for i, v := range s {
			r[i] = v.String()
		}
		return r
	}

It's pretty straightforward: walk through the slice, call the `String`
method on each element, and return a slice of the resulting strings.

This function requires that the element type implement the `String`
method.
The Stringer contract ensures that.

	contract Stringer(T) {
		T String() string
	}

The contract simply says that `T` has to implement the `String`
method.

You may notice that this contract looks like the `fmt.Stringer`
interface, so it's worth pointing out that the argument of the
`ToStrings` function is not a slice of `fmt.Stringer`.
It's a slice of some element type, where the element type implements
`fmt.Stringer`.
The memory representation of a slice of the element type and a slice
of `fmt`.Stringer are normally different, and Go does not support
direct conversions between them.
So this is worth writing, even though `fmt.Stringer` exists.

### Contracts with multiple types

Here is an example of a contract with multiple type parameters.

	type Graph (type Node, Edge G) struct { ... }

	contract G(Node, Edge) {
		Node Edges() []Edge
		Edge Nodes() (from Node, to Node)
	}

	func New (type Node, Edge G) (nodes []Node) *Graph(Node, Edge) {
		...
	}

	func (g *Graph(Node, Edge)) ShortestPath(from, to Node) []Edge {
		...
	}

Here we're describing a graph, built from nodes and edges.
We're not requiring a particular data structure for the graph.
Instead, we're saying that the `Node` type has to have an `Edges`
method that returns the list of edges that connect to the `Node`.
And the `Edge` type has to have a `Nodes` method that returns the two
`Nodes` that the `Edge` connects.

I've skipped the implementation, but this shows the signature of a
`New` function that returns a `Graph`, and the signature of a
`ShortestPath` method on `Graph`.

The important takeaway here is that a contract isn't just about a
single type.  It can describe the relationships between two or more
types.

### Ordered types

One surprisingly common complaint about Go is that it doesn't have a
`Min` function.
Or, for that matter, a `Max` function.
That's because a useful `Min` function should work for any ordered
type, which means that it has to be generic.

While `Min` is pretty trivial to write yourself, any useful generics
implementation should let us add it to the standard library.
This is what it looks like with our design.

{{raw `
	func Min (type T Ordered) (a, b T) T {
		if a < b {
			return a
		}
		return b
	}
`}}

The `Ordered` contract says that the type T has to be an ordered type,
which means that it supports operators like less than, greater than,
and so forth.

	contract Ordered(T) {
		T int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64, uintptr,
			float32, float64,
			string
	}

The `Ordered` contract is just a list of all the ordered types that
are defined by the language.
This contract accepts any of the listed types, or any named type whose
underlying type is one of those types.
Basically, any type you can use with the less than operator.

It turns out that it's much easier to simply enumerate the types that
support the less than operator than it is to invent a new notation
that works for all operators.
After all, in Go, only built-in types support operators.

This same approach can be used for any operator, or more generally
to write a contract for any generic function intended to work with
builtin types.
It lets the writer of the generic function specify clearly the set of
types the function is expected to be used with.
It lets the caller of the generic function clearly see whether the
function is applicable for the types being used.

In practice this contract would probably go into the standard library,
and so really the `Min` function (which will probably also be in the
standard library somewhere) will look like this.
Here we're just referring to the contract `Ordered` defined in the
contracts package.

{{raw `
	func Min (type T contracts.Ordered) (a, b T) T {
		if a < b {
			return a
		}
		return b
	}
`}}

### Generic data structures

Finally, let's look at a simple generic data structure, a binary
tree.  In this example the tree has a comparison function, so there
are no requirements on the element type.

	type Tree (type E) struct {
		root    *node(E)
		compare func(E, E) int
	}

	type node (type E) struct {
		val         E
		left, right *node(E)
	}

Here is how to create a new binary tree.
The comparison function is passed to the `New` function.

	func New (type E) (cmp func(E, E) int) *Tree(E) {
		return &Tree(E){compare: cmp}
	}

An unexported method returns a pointer either to the slot holding v,
or to the location in the tree where it should go.

{{raw `
	func (t *Tree(E)) find(v E) **node(E) {
		pn := &t.root
		for *pn != nil {
			switch cmp := t.compare(v, (*pn).val); {
			case cmp < 0:
				pn = &(*pn).left
			case cmp > 0:
				pn = &(*pn).right
			default:
				return pn
			}
		}
		return pn
	}
`}}

The details here don't really matter, especially since I haven't
tested this code.
I'm just trying to show what it looks like to write a simple generic
data structure.

This is the code for testing whether the tree contains a value.

	func (t *Tree(E)) Contains(v E) bool {
		return *t.find(e) != nil
	}

This is the code for inserting a new value.

	func (t *Tree(E)) Insert(v E) bool {
		pn := t.find(v)
		if *pn != nil {
			return false
		}
		*pn = &node(E){val: v}
		return true
	}

Notice that the type `node` has a type argument `E`.
This is what it looks like to write a generic data structure.
As you can see, it looks like writing ordinary Go code, except that
some type arguments are sprinkled in here and there.

Using the tree is pretty simple.

	var intTree = tree.New(func(a, b int) int { return a - b })

	func InsertAndCheck(v int) {
		intTree.Insert(v)
		if !intTree.Contains(v) {
			log.Fatalf("%d not found after insertion", v)
		}
	}

That's as it should be.
It's a bit harder to write a generic data structure, because you often
have to explicitly write out type arguments for supporting types, but
as much as possible using one is no different from using an ordinary
non-generic data structure.

### Next steps

We are working on actual implementations to allow us to experiment
with this design.
It's important to be able to try out the design in practice, to make
sure that we can write the kinds of programs we want to write.
It hasn't gone as fast as we'd hoped, but we'll send out more detail
on these implementations as they become available.

Robert Griesemer has written a
[preliminary CL](/cl/187317)
that modifies the go/types package.
This permits testing whether code using generics and contracts can
type check.
It’s incomplete right now, but it mostly works for a single package,
and we’ll keep working on it.

What we'd like people to do with this and future implementations is to
try writing and using generic code and see what happens.
We want to make sure that people can write the code they need, and
that they can use it as expected.
Of course not everything is going to work at first, and as we explore
this space we may have to change things.
And, to be clear, we're much more interested in feedback on the
semantics than on details of the syntax.

I’d like to thank everyone who commented on the earlier design, and
everyone who has discussed what generics can look like in Go.
We’ve read all of the comments, and we greatly appreciate the work
that people have put into this.
We would not be where we are today without that work.

Our goal is to arrive at a design that makes it possible to write the
kinds of generic code I’ve discussed today, without making the
language too complex to use or making it not feel like Go anymore.
We hope that this design is a step toward that goal, and we expect to
continue to adjust it as we learn, from our experiences and yours,
what works and what doesn’t.
If we do reach that goal, then we’ll have something that we can
propose for future versions of Go.
