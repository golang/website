---
title: Generic Methods
date: 2026-08-26
by:
- Mark Freeman
summary: Go 1.27 adds generic methods—a highly desired language feature.
template: true
---

Generics was a profound change for Go—it fundamentally expanded the kinds of
programs that gophers could write by adding type parameters to the language.

One could now express generic types and functions in Go. These constructs
decreased the need for specialized data types and functions, reducing program
verbosity and improving the ergonomics of the language in those cases.

To illustrate, different kinds of linked lists could now be condensed into a
single type definition:

```go
// before Go 1.18
type ListOfInts struct {
	elem int
	next *ListOfInts
}
type ListOfStrings struct {
	elem string
	next *ListOfStrings
}
// ...

// after Go 1.18
type List[E any] struct {
	elem E
	next *List[E]
}
```

Likewise, sorting different kinds of ordered data could be expressed as a single
function:

```go
// before Go 1.18
func SortInts(s []int) { /* ... */ }
func SortStrings(s []string) { /* ... */ }
// ...

// after Go 1.18
func Sort[E cmp.Ordered](s []E) { /* ... */ }
```

However, no such capability came to methods.

The generics [proposal](/issue/43651) reasoned that because generic interface
methods are difficult to implement efficiently (discussed
[later](#the-trouble-with-generic-interface-methods)), it didn't make sense to
add type parameters to non-interface (or "concrete") methods either. In that
view, methods are primarily the means of implementing an interface.

Go 1.27 takes a different view and consequently adds generic methods to Go. In
this post, we'll explain this change of view and demonstrate some uses of this
new feature.

# Methods for organization

Methods enable organizing functionality around types. To illustrate, let's
reconsider the linked list from above (adding some conveniences):

```go
type List[E any] struct {
	elem E
	next *List[E]
}
func NewList[E any](elems ...E) List[E] { /* ... */ }
func (List[E]) String() string { /* ... */ }
```

Suppose that one wanted to map this structure to hold a value of some other
type, like a `string`. A method like `ToString` would do:

```go
func (List[E]) ToString(f func(E) string) List[string] { /* ... */ }
```

By providing a “transform” function `f`, `ToString` can be customized with
“off-the-shelf” routines.

For example, `strconv.Itoa` might work for a `List[int]`:

```go
func main() {
	fmt.Println(NewList(1, 2, 3).ToString(strconv.Itoa)) // [1 2 3]
}
```

For a `List[[]byte]`, more options come to mind (depending on the use case):

```go
func main() {
	l := NewList([]byte("Hallo Welt"), []byte("Helló világ"))
	fmt.Println(l).ToString(hex.EncodeToString)					// [48616c6c6f2057656c74 48656c6cc3b32076696cc3a167]
	fmt.Println(l).ToString(base64.StdEncoding.EncodeToString)	// [SGFsbG8gV2VsdA== SGVsbMOzIHZpbMOhZw==]
}
```

Note that parameterizing `List` allowed generalization of the *source* type, but
not the *destination* type, since that depends on the transformation being
applied. This might be reasonable with just a single destination type, but what
if there were many?

In Go 1.18, a generic function could be used for this:

```go
// after Go 1.18
func MapList[E, R any](l List[E], f func(E) R) List[R] { /* ... */ }
```

But using a function has the disadvantage of moving `MapList` into the package
scope—if many data types support mapping operations, things could get crowded.
Furthermore, chained calls must be written “inside out”:

```go
func main() {
	fmt.Println(MapList(MapList(NewList(0, 2, 4), add(2)), divideBy(2))) // [1 2 3]
}
```

Because Go 1.18 didn’t support generic methods, we needed a generic function as
a workaround. This is now remedied in Go 1.27:

```go
// after Go 1.27
func (List[E]) Map[R any](f func(E) R) List[R] { /* ... */ }
```

This is compact, expressive, *and* locally scoped. Additionally, chained calls
are more readable, as they can be written more naturally from left to right:

```go
func main() {
	fmt.Println(NewList(0, 2, 4).Map(add(2)).Map(divideBy(2))) // [1 2 3]
}
```

Prefer the “inside out” form? A method expression can convert any method,
including a generic one, into its function equivalent. Thus, a function call
structure can be recovered if desired:

```go
func main() {
	f := List[int].Map[int]
	fmt.Println(f(f(NewList(0, 2, 4), add(2)), divideBy(2))) // [1 2 3]
}
```

Generic methods, like other generics in Go, must be instantiated (either
explicitly or implicitly) before they are used (called or converted to a
function).

# Decoupling things

If one views methods *also* as an organizational tool, then Go 1.18's reasoning
appears overly restrictive. While a generic concrete method can't help implement
an interface (without generic interface methods), it can still be useful for
code organization. In other words, these concerns can be decoupled.

To elaborate, let's consider a simple interface:

```go
type I interface {
	M()
}
type T struct{}                   // a struct, not an interface
func (T) M[P any]() { /* ... */ } // a generic *concrete* method
```

Here, `T.M` is parameterized by `P`. Any instantiation of `T.M` would result in
a method signature identical to that of `I.M`:

```go
func main() {
	T{}.M[int]()
}
```

Both `T.M[int]` and `I.M` have the signature `func M()`. However, this doesn’t
mean `T` implements `I`. Interface implementation is a property of a (possibly
instantiated) type, not of any particular method. Importantly, `T` does not
declare `T.M[int]`, but rather its generic (and thus uninstantiated) counterpart
`T.M`.

For `T.M` to participate in interface implementation, there would need to exist
a suitable generic method on `I.M`:

```go
type I interface {
	M[P any]() // a generic *interface* method
}
type T struct{}
func (T) M[P any]() { /* ... */ }
```

Since the Go 1.27 syntax doesn’t allow interface methods to declare type
parameters, this code is impossible to write.

But *why* can’t Go have generic interface methods? To understand, a bit of a
detour is needed.

# The trouble with generic interface methods

An interface value is a box that can hold some other value—the boxed value can
be of any type, as long as that type implements the interface. This means that
the declared methods on the type are a superset of those declared on the
interface.

To illustrate, `T` implements `I` below:

```go
type I interface {
	M()
}
type T struct{}
func (T) M() { /* ... */ }
```

One can call methods on interface values:

```go
func main() {
	F(T{})
}

func F(i I) {
	i.M() // "backed" by T.M
}
```

Above, it’s easy to observe that the call to `i.M` will be routed to `T.M`. But
in more complex programs, the relation between an interface value method call
and its *possible* targets is significantly less obvious. This relation
typically spans across packages, making it difficult or impossible to deduce in
the general case.

To illustrate, let’s introduce a package boundary:

```go
// -- package main --
import "p"

func main() {
	p.F(T{})
}
type T struct{}
func (T) M() { /* ... */ }

// -- package p --
func F(i I) {
	i.M()
}
type I interface {
	M()
}
```

Because the two packages are compiled separately, the `main` package won’t know
how `p.F` uses `T{}`. To ensure the existence of any method that `p.F` *might*
call, the compiler generates code for *all* non-generic methods of a type at its
declaration (or instantiation). That way, regardless of what `p.F` does with
`T{}`, the necessary code will exist at runtime.

Now, suppose that `I.M` were generic:

```go
// -- package main --
import "p"

func main() {
	p.F(T{})
}
type T struct{}
func (T) M[P any]() { /* ... */ }

// -- package p --
func F(i I) {
	i.M[int]()
}
type I interface {
	M[P any]()
}
```

Again, the `main` package won’t know how `p.F` uses `T{}`—including how it might
instantiate `T.M`. To handle *any* usage of `T{}` inside `p.F`, the compiler
would need to instantiate `T.M` with *every possible type argument*.

This is impractical with Go’s approach to instantiation, wherein the compiler
generates specific code for each method based on the type arguments. If method
arguments were “boxed” (i.e. passed as values of their constraint interfaces),
each instantiation would share the same code. This would avoid having too many
instantiations at the cost of indirect call overhead—even for direct calls to
instantiated methods.

# Conclusion

Go 1.27 introduces type parameters on concrete methods—this feature has been
highly desired by the Go community, as it permits more ergonomic and readable
code. Although we can’t support generic interface methods, we decided that
allowing type parameters for concrete methods was still worthwhile.

We hope you enjoy using Go’s generic methods and find helpful ways to use them
in your projects!
