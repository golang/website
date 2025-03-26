---
title: What's in an (Alias) Name?
date: 2024-09-17
by:
- Robert Griesemer
tags:
- type aliases
- type parameters
- generics
summary: A description of generic alias types, a planned feature for Go 1.24
---

This post is about generic alias types, what they are, and why we need them.


## Background

Go was designed for programming at scale.
Programming at scale means dealing with large amounts of data, but
also large codebases, with many engineers working on those codebases
over long periods of time.

Go's organization of code into packages enables programming at scale
by splitting up large codebases into smaller, more manageable pieces,
often written by different people, and connected through public
APIs.
In Go, these APIs consist of the identifiers exported by a package:
the exported constants, types, variables, and functions.
This includes the exported fields of structs and methods of types.

As software projects evolve over time or requirements change,
the original organization of the code into packages may turn out to be
inadequate and require _refactoring_.
Refactoring may involve moving exported identifiers and their respective
declarations from an old package to a new package.
This also requires that any references to the moved declarations must be
updated so that they refer to the new location.
In large codebases it may be unpractical or infeasible to make such
a change atomically; or in other words, to do the move and update all clients
in a single change.
Instead, the change must happen incrementally: for instance, to "move"
a function `F`, we add its declaration in a new package without
deleting the original declaration in the old package.
This way, clients can be updated incrementally, over time.
Once all callers refer to `F` in the new package, the original declaration
of `F` may be safely deleted (unless it must be retained indefinitely, for
backward compatibility).
Russ Cox describes refactoring in detail in his 2016 article on
[Codebase Refactoring (with help from Go)](/talks/2016/refactor.article).

Moving a function `F` from one package to another while also retaining it
in the original package is easy: a wrapper function is all that's needed.
To move `F` from `pkg1` to `pkg2`, `pkg2` declares a new function `F`
(the wrapper function) with the same signature as `pkg1.F`, and `pkg2.F`
calls `pkg1.F`.
New callers may call `pkg2.F`, old callers may call `pkg1.F`, yet in both
cases the function eventually called is the same.

Moving constants is similarly straightforward.
Variables take a bit more work: one may have to introduce a pointer to the
original variable in the new package or perhaps use accessor functions.
This is less ideal, but at least it is workable.
The point here is that for constants, variables, and functions,
existing language features exist that permit incremental refactoring as
described above.

But what about moving a type?

In Go, the [(qualified) identifier](/ref/spec#Qualified_identifiers),
or just _name_ for short, determines the _identity_ of types:
a type `T` [defined](/ref/spec#Type_definitions) and exported by a package
`pkg1` is [different](/ref/spec#Type_identity) from an _otherwise identical_
type definition of a type `T` exported by a package `pkg2`.
This property complicates a move of `T` from one package to another while
retaining a copy of it in the original package.
For instance, a value of type `pkg2.T` is not [assignable](/ref/spec#Assignability)
to a variable of type `pkg1.T` because their type names and thus their
type identities are different.
During an incremental update phase, clients may have values and variables
of both types, even though the programmer's intent is for them to have the
same type.

To solve this problem, [Go 1.9](/doc/go1.9) introduced the notion of a
[_type alias_](/ref/spec#Alias_declarations).
A type alias provides a new name for an existing type without introducing
a new type that has a different identity.

In contrast to a regular [type definition](/ref/spec#Type_definitions)

```
type T T0
```

which declares a _new type_ that is never identical to the type on the right-hand side
of the declaration, an [alias declaration](/ref/spec#Alias_declarations)

```
type A = T  // the "=" indicates an alias declaration
```
declares only a _new name_ `A` for the type on the right-hand side:
here, `A` and `T` denote the same and thus identical type `T`.

Alias declarations make it possible to provide a new name (in a new package!)
for a given type while retaining type identity:

```
package pkg2

import "path/to/pkg1"

type T = pkg1.T
```

The type name has changed from `pkg1.T` to `pkg2.T` but values
of type `pkg2.T` have the same type as variables of type `pkg1.T`.


## Generic alias types

[Go 1.18](/doc/go1.18) introduced generics.
Since that release, type definitions and function
declarations can be customized through type parameters.
For technical reasons, alias types didn't gain the same ability at that time.
Obviously, there were also no large codebases exporting generic
types and requiring refactoring.

Today, generics have been around for a couple of years, and large codebases
are making use of generic features.
Eventually the need will arise to refactor these codebases, and with that the
need to migrate generic types from one package to another.

To support incremental refactorings involving generic types, the future Go 1.24 release,
planned for early February 2025, will fully support type parameters on alias types
in accordance with proposal [#46477](/issue/46477).
The new syntax follows the same pattern as it does for type definitions and function declarations,
with an optional type parameter list following the identifier (the alias name) on the left-hand side.
Before this change one could only write:

```
type Alias = someType
```

but now we can also declare type parameters with the alias declaration:

```
type Alias[P1 C1, P2 C2] = someType
```

Consider the previous example, now with generic types.
The original package `pkg1` declared and exported a generic type `G` with a type parameter `P`
that is suitably constrained:

```
package pkg1

type Constraint      someConstraint
type G[P Constraint] someType
```

If the need arises to provide access to the same type `G` from a new package `pkg2`,
a generic alias type is just the ticket [(playground)](/play/p/wKOf6NbVtdw?v=gotip):

```
package pkg2

import "path/to/pkg1"

type Constraint      = pkg1.Constraint  // pkg1.Constraint could also be used directly in G
type G[P Constraint] = pkg1.G[P]
```

Note that one **cannot** simply write

```
type G = pkg1.G
```

for a couple of reasons:

1) Per [existing spec rules](/ref/spec#Type_definitions), generic
types must be [instantiated](/ref/spec#Instantiations) when they
are _used_.
The right-hand side of the alias declaration uses the type `pkg1.G` and
therefore type arguments must be provided.
Not doing so would require an exception for this case, making the spec more
complicated.
It is not obvious that the minor convenience is worth the complication.

2) If the alias declaration doesn't need to declare its own type parameters and
instead simply "inherits" them from the aliased type `pkg1.G`, the declaration of
`G` provides no indication that it is a generic type.
Its type parameters and constraints would have to be retrieved from the declaration
of `pkg1.G` (which itself might be an alias).
Readability will suffer, yet readable code is one of the primary aims of the Go project.

Writing down an explicit type parameter list may seem like an unnecessary burden
at first, but it also provides additional flexibility.
For one, the number of type parameters declared by the alias type doesn't have to
match the number of type parameters of the aliased type.
Consider a generic map type:

```
type Map[K comparable, V any] mapImplementation
```

If uses of `Map` as sets are common, the alias

```
type Set[K comparable] = Map[K, bool]
```

might be useful [(playground)](/play/p/IxeUPGCztqf?v=gotip).
Because it is an alias, types such as `Set[int]` and `Map[int, bool]` are
identical.
This would not be the case if `Set` were a [defined](/ref/spec#Type_definitions)
(non-alias) type.

Furthermore, the type constraints of a generic alias type don't have to match the
constraints of the aliased type, they only have to
[satisfy](/ref/spec#Satisfying_a_type_constraint) them.
For instance, reusing the set example above, one could define
an `IntSet` as follows:

```
type integers interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 }
type IntSet[K integers] = Set[K]
```

This map can be instantiated with any key type that satisfies the `integers`
constraint [(playground)](/play/p/0f7hOAALaFb?v=gotip).
Because `integers` satisfies `comparable`, the type parameter `K` may be used
as type argument for the `K` parameter of `Set`, following the usual
instantiation rules.

Finally, because an alias may also denote a type literal, parameterized aliases
make it possible to create generic type literals
[(playground)](/play/p/wql3NJaUs0o?v=gotip):

```
type Point3D[E any] = struct{ x, y, z E }
```

To be clear, none of these examples are "special cases" or somehow require
additional rules in the spec. They follow directly from the application of
the existing rules put in place for generics. The only thing that changed in the
spec is the ability to declare type parameters in an alias declaration.


## An interlude about type names

Before the introduction of alias types, Go had only one form of type declarations:

```
type TypeName existingType
```

This declaration creates a new and different type from an existing type
and gives that new type a name.
It was natural to call such types _named types_ as they have a _type name_
in contrast to unnamed [type literals](/ref/spec#Types) such as
`struct{ x, y int }`.

With the introduction of alias types in Go 1.9 it became possible to give
a name (an alias) to type literals, too. For instance, consider:

```
type Point2D = struct{ x, y int }
```

Suddenly, the notion of a _named type_ describing something that is different from
a type literal didn't make that much sense anymore, since an alias name clearly is
a name for a type, and thus the denoted type (which might be a type literal, not a type name!)
arguably could be called a "named type".

Because (proper) named types have special properties (one can bind methods to them,
they follow different assignment rules, etc.), it seemed prudent to use a new
term in order to avoid confusions.
Thus, since Go 1.9, the spec calls the types formerly called named types _defined types_:
only defined types have properties (methods, assignability restrictions, etc) that are
tied to their names.
Defined types are introduced through type definitions, and alias types are
introduced through alias declarations.
In both cases, names are given to types.

The introduction of generics in Go 1.18 made things more complicated.
Type parameters are types, too, they have a name, and they share rules
with defined types.
For instance, like defined types, two differently named type parameters
denote different types.
In other words, type parameters are named types, and furthermore, they
behave similarly to Go's original named types in some ways.

To top things off, Go's predeclared types (`int`, `string` and so on)
can only be accessed through their names, and like defined types and
type parameters, are different if their names are different
(ignoring for a moment the `byte` and `rune` alias types).
The predeclared types truly are named types.

Therefore, with Go 1.18, the spec came full circle and formally
re-introduced the notion of a [named type](/ref/spec#Types) which now
comprises "predeclared types, defined types, and type parameters".
To correct for alias types denoting type literals the spec says:
"An alias denotes a named type if the type given in the alias declaration
is a named type."

Stepping back and outside the box of Go nomenclature for a moment, the correct
technical term for a named type in Go is probably
[_nominal type_](https://en.wikipedia.org/wiki/Nominal_type_system).
A nominal type's identity is explicitly tied to its name which is exactly what
Go's named types (now using the 1.18 terminology) are all about.
A nominal type's behavior is in contrast to a _structural type_ which has
behavior that only depends on its structure and not its name
(if it has one in the first place).
Putting it all together, Go's predeclared, defined, and type parameter types are
all nominal types, while Go's type literals and aliases denoting type literals
are structural types.
Both nominal and structural types can have names, but having a name
doesn't mean the type is nominal, it just means it is named.

None of this matters for day-to-day use of Go and in practice the details
can safely be ignored.
But precise terminology matters in the spec because it makes it easier
to describe the rules governing the language.
So should the spec change its terminology one more time?
It is probably not worth the churn: it is not just the spec that would need
to be updated, but also a lot of supporting documentation.
A fair number of books written on Go might become inaccurate.
Furthermore, "named", while less precise, is probably intuitively clearer
than "nominal" for most people.
It also matches the original terminology used in the spec, even if it now
requires an exception for alias types denoting type literals.


## Availability

Implementing generic type aliases has taken longer than expected:
the necessary changes required adding a new exported `Alias` type
to [`go/types`](/pkg/go/types) and then adding the ability to record type parameters
with that type.
On the compiler side, the analogous changes also required modifications to
the export data format, the file format that describes a package's
exports, which now needs to be able to describe type parameters for aliases.
The impact of these changes is not confined to the compiler, but affects
clients of `go/types` and thus many third-party packages.
This was very much a change affecting a large code base; to avoid
breaking things, an incremental roll-out over several releases was necessary.

After all this work, generic alias types will finally be available by default in Go 1.24.

To allow third-party clients to get their code ready, starting with
Go 1.23, support for generic type aliases can be enabled by setting
`GOEXPERIMENT=aliastypeparams` when invoking the `go` tool.
However, be aware that support for exported generic aliases is still
missing for that version.

Full support (including export) is implemented at tip, and the default
setting for `GOEXPERIMENT` will soon be switched so that generic type
aliases are enabled by default.
Thus, another option is to experiment with the latest version of Go
at tip.

As always, please let us know if you encounter any problems by filing an
[issue](/issue/new);
the better we test a new feature, the smoother the general roll-out.

Thanks and happy refactoring!
