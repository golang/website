---
title: All your comparable types
date: 2023-02-17
by:
- Robert Griesemer
summary: type parameters, type sets, comparable types, constraint satisfaction
---

On February 1 we released our latest Go version, 1.20,
which included a few language changes.
Here we'll discuss one of those changes: the predeclared `comparable` type constraint
is now satisfied by all [comparable types](/ref/spec#Comparison_operators).
Surprisingly, before Go 1.20, some comparable types did not satisfy `comparable`!

If you're confused, you've come to the right place.
Consider the valid map declaration

```Go
var lookupTable map[any]string
```

where the map's key type is `any` (which is a
[comparable type](/ref/spec#Comparison_operators)).
This works perfectly fine in Go.
On the other hand, before Go 1.20, the seemingly equivalent generic map type

```Go
type genericLookupTable[K comparable, V any] map[K]V
```

could be used just like a regular map type, but produced a compile-time error when
`any` was used as the key type:

```Go
var lookupTable genericLookupTable[any, string] // ERROR: any does not implement comparable (Go 1.18 and Go 1.19)
```

Starting with Go 1.20 this code will compile just fine.

The pre-Go 1.20 behavior of `comparable` was particularly annoying because it
prevented us from writing the kind of generic libraries we were hoping to write with
generics in the first place.
The proposed [`maps.Clone`](/issue/57436) function

```Go
func Clone[M ~map[K]V, K comparable, V any](m M) M { … }
```
can be written but could not be used for a map such as `lookupTable` for the same reason
our `genericLookupTable` could not be used with `any` as key type.

In this blog post, we hope to shine some light on the language mechanics behind all this.
In order to do so, we start with a bit of background information.

## Type parameters and constraints

Go 1.18 introduced generics and, with that,
[_type parameters_](/ref/spec#Type_parameter_declarations)
as a new language construct.

In an ordinary function, a parameter ranges over a set of values that is restricted by its type.
Analogously, in a generic function (or type), a type parameter ranges over a set of types that is restricted
by its [_type constraint_](/ref/spec#Type_constraints).
Thus, a type constraint defines the _set of types_ that are permissible
as type arguments.

Go 1.18 also changed how we view interfaces: while in the past an
interface defined a set of methods, now an interface defines a set of types.
This new view is completely backward compatible:
for any given set of methods defined by an interface, we can imagine the (infinite)
set of all types that implement those methods.
For instance, given an [`io.Writer`](/pkg/io#Writer) interface,
we can imagine the infinite set of all types that have a `Write` method
with the appropriate signature.
All of these types _implement_ the interface because they all have the
required `Write` method.

But the new type set view is more powerful than the old method set one:
we can describe a set of types explicitly, not only indirectly through methods.
This gives us new ways to control a type set.
Starting with Go 1.18, an interface may embed not just other interfaces,
but any type, a union of types, or an infinite set of types that share the
same [underlying type](/ref/spec#Underlying_types). These types are then included in the
[type set computation](/ref/spec#General_interfaces):
the union notation `A|B` means "type `A` or type `B`",
and the `~T` notation stands for "all types that have the underlying type `T`".
For instance, the interface

```Go
interface {
	~int | ~string
	io.Writer
}
```

defines the set of all types whose underlying types are either `int` or `string`
and that also implement `io.Writer`'s `Write` method.

Such generalized interfaces can't be used as variable types.
But because they describe type sets they are used as type constraints, which
are sets of types.
For instance, we can write a generic `min` function

```Go
func min[P interface{ ~int64 | ~float64 }](x, y P) P
```

which accepts any `int64` or `float64` argument.
(Of course, a more realistic implementation would use a constraint that
enumerates all basic types with an <code>&lt;</code> operator.)

As an aside, because enumerating explicit types without methods is common,
a little bit of [syntactic sugar](https://en.wikipedia.org/wiki/Syntactic_sugar)
allows us to [omit the enclosing `interface{}`](/ref/spec#General_interfaces),
leading to the compact and more idiomatic

```Go
func min[P ~int64 | ~float64](x, y P) P { … }
```

With the new type set view we also need a new way to explain what it means
to [_implement_](/ref/spec#Implementing_an_interface) an interface.
We say that a (non-interface) type `T` implements
an interface `I` if `T` is an element of the interface's type set.
If `T` is an interface itself, it describes a type set. Every single type in that set
must also be in the type set of `I`, otherwise `T` would contain types that do not implement `I`.
Thus, if `T` is an interface, it implements interface `I` if the type
set of `T` is a subset of the type set of `I`.

Now we have all the ingredients in place to understand constraint satisfaction.
As we have seen earlier, a type constraint describes the set of acceptable argument
types for a type parameter. A type argument satisfies the corresponding type parameter
constraint if the type argument is in the set described by the constraint interface.
This is another way of saying that the type argument implements the constraint.
In Go 1.18 and Go 1.19, constraint satisfaction meant constraint implementation.
As we'll see in a bit, in Go 1.20 constraint satisfaction is not quite constraint
implementation anymore.

## Operations on type parameter values

A type constraint does not just specify what type arguments are acceptable for a type parameter,
it also determines the operations that are possible on values of a type parameter.
As we would expect, if a constraint defines a method such as `Write`,
the `Write` method can be called on a value of the respective type parameter.
More generally, an operation such as `+` or `*` that is supported by all types in the type set
defined by a constraint is permitted with values of the corresponding type parameter.

For instance, given the `min` example, in the function body any operation that is supported by
`int64` and `float64` types is permitted on values of the type parameter `P`.
That includes all the basic arithmetic operations, but also comparisons such as <code>&lt;</code>.
But it does not include bitwise operations such as `&` or `|`
because those operations are not defined on `float64` values.

## Comparable types

In contrast to other unary and binary operations, `==` is defined on not just
a limited set of
[predeclared types](/ref/spec#Types), but on an infinite variety of types,
including arrays, structs, and interfaces.
It is impossible to enumerate all these types in a constraint.
We need a different mechanism to express that a type parameter must support `==`
(and `!=`, of course) if we care about more than predeclared types.

We solve this problem through the predeclared type
[`comparable`](/ref/spec#Predeclared_identifiers), introduced with Go 1.18.
`comparable` is
an interface type whose type set is the infinite set of comparable types, and that
may be used as a constraint whenever we require a type argument to support `==`.

Yet, the set of types comprised by `comparable` is not the same
as the set of all [comparable types](/ref/spec#Comparison_operators) defined by the Go spec.
[By construction](/ref/spec#Interface_types), a type set specified by an interface
(including `comparable`) does not contain the interface itself (or any other interface).
Thus, an interface such as `any` is not included in `comparable`,
even though all interfaces support `==`.
What gives?

Comparison of interfaces (and of composite types containing them) may panic at run time:
this happens when the dynamic type, the type of the actual value stored in the
interface variable, is not comparable.
Consider our original `lookupTable` example: it accepts arbitrary values as keys.
But if we try to enter a value with a key that does not support `==`, say
a slice value, we get a run-time panic:

```Go
lookupTable[[]int{}] = "slice"  // PANIC: runtime error: hash of unhashable type []int
```

By contrast, `comparable` contains only types that the compiler guarantees will not panic with `==`.
We call these types _strictly comparable_.

Most of the time this is exactly what we want: it's comforting to know that `==` in a generic
function won't panic if the operands are constrained by `comparable`, and it is what we
would intuitively expect.

Unfortunately, this definition of `comparable` together with the rules for
constraint satisfaction prevented us from writing useful
generic code, such as the `genericLookupTable` type shown earlier:
for `any` to be an acceptable argument type, `any` must satisfy (and therefore implement) `comparable`.
But the type set of `any` is larger than (not a subset of) the type set of `comparable`
and therefore does not implement `comparable`.

```Go
var lookupTable GenericLookupTable[any, string] // ERROR: any does not implement comparable (Go 1.18 and Go 1.19)
```

Users recognized the problem early on and filed a multitude of issues and proposals in short order
([#51338](/issue/51338),
[#52474](/issue/52474),
[#52531](/issue/52531),
[#52614](/issue/52614),
[#52624](/issue/52624),
[#53734](/issue/53734),
etc).
Clearly this was a problem we needed to address.

The "obvious" solution was simply to include even non-strictly comparable types in the
`comparable` type set.
But this leads to inconsistencies with the type set model.
Consider the following example:

```Go
func f[Q comparable]() { … }

func g[P any]() {
        _ = f[int] // (1) ok: int implements comparable
        _ = f[P]   // (2) error: type parameter P does not implement comparable
        _ = f[any] // (3) error: any does not implement comparable (Go 1.18, Go.19)
}
```

Function `f` requires a type argument that is strictly comparable.
Obviously it is ok to instantiate `f` with `int`: `int` values never panic on `==`
and thus `int` implements `comparable` (case 1).
On the other hand, instantiating `f` with `P` is not permitted: `P`'s type set is defined
by its constraint `any`, and `any` stands for the set of all possible types.
This set includes types that are not comparable at all.
Hence, `P` doesn't implement `comparable` and thus cannot be used to instantiate `f`
(case 2).
And finally, using the type `any` (rather than a type parameter constrained by `any`)
doesn't work either, because of exactly the same problem (case 3).

Yet, we do want to be able to use the type `any` as type argument in this case.
The only way out of this dilemma was to change the language somehow.
But how?

## Interface implementation vs constraint satisfaction

As mentioned earlier, constraint satisfaction is interface implementation:
a type argument `T` satisfies a constraint `C` if `T` implements `C`.
This makes sense: `T` must be in the type set expected by `C` which is
exactly the definition of interface implementation.

But this is also the problem because it prevents us from using non-strictly comparable
types as type arguments for `comparable`.

So for Go 1.20, after almost a year of publicly discussing numerous alternatives
(see the issues mentioned above), we decided to introduce an exception for just this case.
To avoid the inconsistency, rather than changing what `comparable` means,
we differentiated between _interface implementation_,
which is relevant for passing values to variables, and _constraint satisfaction_,
which is relevant for passing type arguments to type parameters.
Once separated, we could give each of those concepts (slightly) different
rules, and that is exactly what we did with proposal [#56548](/issue/56548).

The good news is that the exception is quite localized in the
[spec](/ref/spec#Satisfying_a_type_constraint).
Constraint satisfaction remains almost the same as interface implementation, with a caveat:

> A type `T` satisfies a constraint `C` if
>
> - `T` implements `C`; or
> - `C` can be written in the form `interface{ comparable; E }`, where `E` is a basic interface
>   and `T` is [comparable](/ref/spec#Comparison_operators) and implements `E`.

The second bullet point is the exception.
Without going too much into the formalism of the spec, what the exception says is the following:
a constraint `C` that expects strictly comparable types (and which may also have other requirements
such as methods `E`) is satisfied by any type argument `T` that supports `==`
(and which also implements the methods in `E`, if any).
Or even shorter: a type that supports `==` also satisfies `comparable`
(even though it may not implement it).

We can immediately see that this change is backward-compatible:
before Go 1.20, constraint satisfaction was the same as interface implementation, and we still
have that rule (1st bullet point).
All code that relied on that rule continues to work as before.
Only if that rule fails do we need to consider the exception.

Let's revisit our previous example:

```Go
func f[Q comparable]() { … }

func g[P any]() {
        _ = f[int] // (1) ok: int satisfies comparable
        _ = f[P]   // (2) error: type parameter P does not satisfy comparable
        _ = f[any] // (3) ok: satisfies comparable (Go 1.20)
}
```

Now, `any` does satisfy (but not implement!) `comparable`.
Why?
Because Go permits `==` to be used with values of type `any`
(which corresponds to the type `T` in the spec rule),
and because the constraint `comparable` (which corresponds to the constraint `C` in the rule)
can be written as `interface{ comparable; E }` where `E` is simply the empty interface
in this example (case 3).

Interestingly, `P` still does not satisfy `comparable` (case 2).
The reason is that `P` is a type parameter constrained by `any` (it _is not_ `any`).
The operation `==` is _not_ available with all types in the type set of `P`
and thus not available on `P`;
it is not a [comparable type](/ref/spec#Comparison_operators).
Therefore the exception doesn't apply.
But this is ok: we do like to know that `comparable`, the strict comparability
requirement, is enforced most of the time. We just need an exception for
Go types that support `==`, essentially for historical reasons:
we always had the ability to compare non-strictly comparable types.

## Consequences and remedies

We gophers take pride in the fact that language-specific behavior
can be explained and reduced to a fairly compact set of rules, spelled out
in the language spec.
Over the years we have refined these rules, and when possible made them
simpler and often more general.
We also have been careful to keep the rules orthogonal,
always on the lookout for unintended and unfortunate consequences.
Disputes are resolved by consulting the spec, not by decree.
That is what we have aspired to since the inception of Go.

_One does not simply add an exception to a carefully crafted type
system without consequences!_

So where's the catch?
There's an obvious (if mild) drawback, and a less obvious (and more
severe) one.
Obviously, we now have a more complex rule for constraint satisfaction
which is arguably less elegant than what we had before.
This is unlikely to affect our day-to-day work in any significant way.

But we do pay a price for the exception: in Go 1.20, generic functions
that rely on `comparable` are not statically type-safe anymore.
The `==` and `!=` operations may panic if applied to operands of
`comparable` type parameters, even though the declaration says
that they are strictly comparable.
A single non-comparable value may sneak its way through
multiple generic functions or types by way of a single non-strictly
comparable type argument and cause a panic.
In Go 1.20 we can now declare

```Go
var lookupTable genericLookupTable[any, string]
```

without compile-time error, but we will get a run-time panic
if we ever use a non-strictly comparable key type in this case, exactly like we would
with the built-in `map` type.
We have given up static type safety for a run-time check.

There may be situations where this is not good enough,
and where we want to enforce strict comparability.
The following observation allows us to do exactly that, at least in limited
form: type parameters do not benefit from the exception that we added to the
constraint satisfaction rule.
For instance, in our earlier example, the type parameter `P` in the function
`g` is constrained by `any` (which by itself is comparable but not strictly comparable)
and so `P` does not satisfy `comparable`.
We can use this knowledge to craft a compile-time assertion of sorts for
a given type `T`:

```Go
type T struct { … }
```

We want to assert that `T` is strictly comparable.
It's tempting to write something like:

```Go
// isComparable may be instantiated with any type that supports ==
// including types that are not strictly comparable because of the
// exception for constraint satisfaction.
func isComparable[_ comparable]() {}

// Tempting but not quite what we want: this declaration is also
// valid for types T that are not strictly comparable.
var _ = isComparable[T] // compile-time error if T does not support ==
```

The dummy (blank) variable declaration serves as our "assertion".
But because of the exception in the constraint satisfaction rule,
`isComparable[T]` only fails if `T` is not comparable at all;
it will succeed if `T` supports `==`.
We can work around this problem by using `T` not as a type argument,
but as a type constraint:

```Go
func _[P T]() {
	_ = isComparable[P] // P supports == only if T is strictly comparable
}
```

Here is a [passing](/play/p/9i9iEto3TgE) and [failing](/play/p/5d4BeKLevPB) playground example
illustrating this mechanism.

## Final observations

Interestingly, until two months before the Go 1.18
release, the compiler implemented constraint satisfaction exactly as we do
now in Go 1.20.
But because at that time constraint satisfaction meant interface implementation,
we did have an implementation that was inconsistent with the language specification.
We were alerted to this fact with [issue #50646](/issue/50646).
We were extremely close to the release and had to make a decision quickly.
In the absence of a convincing solution, it seemed safest to make the
implementation consistent with the spec.
A year later, and with plenty of time to consider different approaches,
it seems that the implementation we had was the implementation we wanted in the first place.
We have come full circle.

As always, please let us know if anything doesn't work as expected
by filing issues at [https://go.dev/issue/new](/issue/new).

Thank you!
