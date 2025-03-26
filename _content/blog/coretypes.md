---
title: Goodbye core types - Hello Go as we know and love it!
date: 2025-03-26
by:
- Robert Griesemer
summary: Go 1.25 simplifies the language spec by removing the notion of core types
---

The Go 1.18 release introduced generics and with that a number of new features, including type parameters, type constraints, and new concepts such as type sets.
It also introduced the notion of a _core type_.
While the former provide concrete new functionality, a core type is an abstract construct that was introduced
for expediency and to simplify dealing with generic operands (operands whose types are type parameters).
In the Go compiler, code that in the past relied on the [underlying type](/ref/spec/#Underlying_types) of an operand,
now instead had to call a function computing the operand's core type.
In the language spec, in many places we just needed to replace "underlying type" with "core type".
What's not to like?

Quite a few things, as it turns out!
To understand how we got here, it's useful to briefly revisit how type parameters and type constraints work.

## Type parameters and type constraints

A type parameter is a placeholder for a future type argument;
it acts like a _type variable_ whose value is known at compile time,
similar to how a named constant stands for a number, string, or bool whose value is known at compile time.
Like ordinary variables, type parameters have a type.
That type is described by their _type constraint_ which determines
what operations are permitted on operands whose type is the respective type parameter.

Any concrete type that instantiates a type parameter must satisfy the type parameter's constraint.
This ensures that an operand whose type is a type parameter possesses all of the respective type constraint's properties,
no matter what concrete type is used to instantiate the type parameter.

In Go, type constraints are described through a mixture of method and type requirements which together
define a _type set_: this is the set of all the types that satisfy all the requirements. Go uses a
generalized form of interfaces for this purpose. An interface enumerates a set of methods and types,
and the type set described by such an interface consists of all the types that implement those methods
and that are included in the enumerated types.

For instance, the type set described by the interface

```Go
type Constraint interface {
	~[]byte | ~string
	Hash() uint64
}
```

consists of all the types whose representation is `[]byte` or `string` and whose method set includes the `Hash` method.

With this we can now write down the rules that govern operations on generic operands.
For instance, the [rules for index expressions](/ref/spec#Index_expressions) state that (among other things)
for an operand `a` of type parameter type `P`:

> The index expression `a[x]` must be valid for values of all types in `P`'s type set.
> The element types of all types in `P`'s type set must be identical.
  (In this context, the element type of a string type is `byte`.)

These rules make it possible to index the generic variable `s` below ([playground](/play/p/M1LYKm3x3IB)):

```Go
func at[bytestring Constraint](s bytestring, i int) byte {
	return s[i]
}
```

The indexing operation `s[i]` is permitted because the type of `s` is `bytestring`, and the type constraint (type set) of
`bytestring` contains `[]byte` and `string` types for which indexing with `i` is valid.

## Core types

This type set-based approach is very flexible and in line with the intentions of the
[original generics proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md):
an operation involving operands of generic type should be valid if it is valid for any type permitted by the respective
type constraint.
To simplify matters with respect to the implementation, knowing that we would be able to relax rules later,
this approach was _not_ chosen universally.
Instead, for instance, for [Send statements](/ref/spec#Send_statements), the spec states that

> The channel expression's _core type_ must be a channel, the channel direction must permit send operations,
> and the type of the value to be sent must be assignable to the channel's element type.

These rules are based on the notion of a core type which is defined roughly as follows:

- If a type is not a type parameter, its core type is just its [underlying type](/ref/spec#Underlying_types).
- If the type is a type parameter, the core type is the single underlying type of all the types in the type parameter's type set.
  If the type set has _different_ underlying types, the core type doesn't exist.

For instance, `interface{ ~[]int }` has a core type (`[]int`), but the `Constraint` interface above does not have a core type.
To make things more complicated, when it comes to channel operations and certain built-in calls (`append`, `copy`) the above definition
of core types is too restrictive.
The actual rules have adjustments that allow for differing channel directions and type sets containing both `[]byte` and `string` types.

There are various problems with this approach:

- Because the definition of core type must lead to sound type rules for different language features,
it is overly restrictive for specific operations.
For instance, the Go 1.24 rules for [slice expressions](/ref/spec#Slice_expressions) do rely on core types,
and as a consequence slicing an operand of type `S` constrained by `Constraint` is not permitted, even though
it could be valid.

- When trying to understand a specific language feature, one may have to learn the intricacies of
core types even when considering non-generic code.
Again, for slice expressions, the language spec talks about the core type of the sliced operand,
rather than just stating that the operand must be an array, slice, or string.
The latter is more direct, simpler, and clearer, and doesn't require knowing another concept that may be
irrelevant in the concrete case.

- Because the notion of core types exists, the rules for index expressions, and `len` and `cap` (and others),
which all eschew core types, appear as exceptions in the language rather than the norm.
In turn, core types cause proposals such as [issue #48522](/issue/48522) which would permit a selector
`x.f` to access a field `f` shared by all elements of `x`'s type set, to appear to add more exceptions to the
language.
Without core types, that feature becomes a natural and useful consequence of the ordinary rules for non-generic
field access.

## Go 1.25

For the upcoming Go 1.25 release (August 2025) we decided to remove the notion of core types from the
language spec in favor of explicit (and equivalent!) prose where needed.
This has multiple benefits:

- The Go spec presents fewer concepts, making it easier to learn the language.
- The behavior of non-generic code can be understood without reference to generics concepts.
- The individualized approach (specific rules for specific operations) opens the door for more flexible rules.
We already mentioned [issue #48522](/issue/48522), but there are also ideas for more powerful
slice operations, and [improved type inference](/issue/69153).

The respective [proposal issue #70128](/issue/70128) was recently approved and the relevant changes
are already implemented.
Concretely this means that a lot of prose in the language spec was reverted to its original,
pre-generics form, and new paragraphs were added where needed to explain the rules as they
pertain to generic operands. Importantly, no behavior was changed.
The entire section on core types was removed.
The compiler's error messages were updated to not mention "core type" anymore, and in many
cases error messages are now more specific by pointing out exactly which type in a type set
is causing a problem.

Here is a sample of the changes made. For the built-in function `close`,
starting with Go 1.18 the spec began as follows:

> For an argument `ch` with core type that is a channel,
> the built-in function `close` records that no more values will be sent on the channel.

A reader who simply wanted to know how `close` works, had to first learn about core types.
Starting with Go 1.25, this section will again begin the same way it began before Go 1.18:

> For a channel `ch`, the built-in function `close(ch)`
> records that no more values will be sent on the channel.

This is shorter and easier to understand.
Only when the reader is dealing with a generic operand will they have to contemplate
the newly added paragraph:

> If the type of the argument to `close` is a type parameter
> all types in its type set must be channels with the same element type.
> It is an error if any of those channels is a receive-only channel.

We made similar changes to each place that mentioned core types.
In summary, although this spec update does not affect any current Go program, it opens the
door to future language improvements while making the language as it is today easier to
learn and its spec simpler.
