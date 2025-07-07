---
title: Generic interfaces
date: 2025-07-07
by:
- Axel Wagner
tags:
- type parameters
- generics
- interfaces
summary: Adding type parameters to interface types is surprisingly powerful
---

There is an idea that is not obvious until you hear about it for the first time: as interfaces are types themselves, they too can have type parameters.
This idea proves to be surprisingly powerful when it comes to expressing constraints on generic functions and types.
In this post, we'll demonstrate it, by discussing the use of interfaces with type parameters in a couple of common scenarios.

## A simple tree set

As a motivating example, assume we need a generic version of a [binary search tree](https://en.wikipedia.org/wiki/Binary_search_tree).
The elements stored in such a tree need to be ordered, so our type parameter needs a constraint that determines the ordering to use.
A simple option is to use the [cmp.Ordered](/pkg/cmp#Ordered) constraint, introduced in Go 1.21.
It restricts a type parameter to ordered types (strings and numbers) and allows methods of the type to use the built-in ordering operators.

{{raw `
    // The zero value of a Tree is a ready-to-use empty tree.
    type Tree[E cmp.Ordered] struct {
        root *node[E]
    }

    func (t *Tree[E]) Insert(element E) {
        t.root = t.root.insert(element)
    }

    type node[E cmp.Ordered] struct {
        value E
        left  *node[E]
        right *node[E]
    }

    func (n *node[E]) insert(element E) *node[E] {
        if n == nil {
            return &node[E]{value: element}
        }
        switch {
        case element < n.value:
            n.left = n.left.insert(element)
        case element > n.value:
            n.right = n.right.insert(element)
        }
        return n
    }
`}}

([playground](/play/p/H7-n33X7P2h))

However, this approach has the disadvantage that it only works on basic types for which <code>&lt;</code> is defined;
you cannot insert struct types, like [time.Time](/pkg/time#Time).

We can remedy that by requiring the user to provide a comparison function:

{{raw `
    // A FuncTree must be created with NewTreeFunc.
    type FuncTree[E any] struct {
        root *funcNode[E]
        cmp  func(E, E) int
    }

    func NewFuncTree[E any](cmp func(E, E) int) *FuncTree[E] {
        return &FuncTree[E]{cmp: cmp}
    }

    func (t *FuncTree[E]) Insert(element E) {
        t.root = t.root.insert(t.cmp, element)
    }

    type funcNode[E any] struct {
        value E
        left  *funcNode[E]
        right *funcNode[E]
    }

    func (n *funcNode[E]) insert(cmp func(E, E) int, element E) *funcNode[E] {
        if n == nil {
            return &funcNode[E]{value: element}
        }
        sign := cmp(element, n.value)
        switch {
        case sign < 0:
            n.left = n.left.insert(cmp, element)
        case sign > 0:
            n.right = n.right.insert(cmp, element)
        }
        return n
    }
`}}

([playground](/play/p/tiEjuxCHtFF))

This works, but it also comes with downsides.
We can no longer use the zero value of our container type, because it needs to have an explicitly initialized comparison function.
And the use of a function field makes it harder for the compiler to inline the comparison calls, which can introduce a significant runtime overhead.

Using a method on the element type can solve these issues, because methods are directly associated with a type.
A method does not have to be explicitly passed and the compiler can see the target of the call and may be able to inline it.
But how can we express the constraint to require that element types provide the necessary method?

## Using the receiver in constraints

The first approach we might try is to define a plain old interface with a `Compare` method:

{{raw `
    type Comparer interface {
        Compare(Comparer) int
    }
`}}

However, we quickly realize that this does not work well.
To implement this interface, the method's parameter must itself be `Comparer`.
Not only does that mean that the implementation of this method must type-assert the parameter to its own type, it also requires that every type must explicitly refer to our package with the `Comparer` type by name (otherwise the method signatures would not be identical).
That is not very orthogonal.

A better approach is to make the `Comparer` interface itself generic:

{{raw `
    type Comparer[T any] interface {
        Compare(T) int
    }
`}}

This `Comparer` now describes a whole family of interfaces, one for each type that `Comparer` may be instantiated with.
A type that implements `Comparer[T]` declares "I can compare myself to a `T`".
For instance, `time.Time` naturally implements `Comparer[time.Time]` because [it has a matching `Compare` method](/pkg/time#Time.Compare):

{{raw `
    // Implements Comparer[Time]
    func (t Time) Compare(u Time) int
`}}

This is better, but not enough.
What we really want is a constraint that says that a type parameter can be compared to *itself*: we want the constraint to be self-referential.
The subtle insight is that the self-referential aspect does not have to be part of the interface definition itself; specifically, the constraint for `T` in the `Comparer` type is just `any`.
Instead, it is a consequence of how we use `Comparer` as a constraint for the type parameter of `MethodTree`:

{{raw `
    // The zero value of a MethodTree is a ready-to-use empty tree.
    type MethodTree[E Comparer[E]] struct {
        root *methodNode[E]
    }

    func (t *MethodTree[E]) Insert(element E) {
        t.root = t.root.insert(element)
    }

    type methodNode[E Comparer[E]] struct {
        value E
        left  *methodNode[E]
        right *methodNode[E]
    }

    func (n *methodNode[E]) insert(element E) *methodNode[E] {
        if n == nil {
            return &methodNode[E]{value: element}
        }
        sign := element.Compare(n.value)
        switch {
        case sign < 0:
            n.left = n.left.insert(element)
        case sign > 0:
            n.right = n.right.insert(element)
        }
        return n
    }
`}}

([playground](/play/p/LuhzYej_2SP))

Because `time.Time` implements `Comparer[time.Time]` it is now a valid type argument for this container, and we can still use the zero value as an empty container:

{{raw `
    var t MethodTree[time.Time]
    t.Insert(time.Now())
`}}

For full flexibility, a library can provide all three API versions.
If we want to minimize repetition, all versions could use a shared implementation.
We could use the function version for that, as it is the most general:

{{raw `
    type node[E any] struct {
        value E
        left  *node[E]
        right *node[E]
    }

    func (n *node[E]) insert(cmp func(E, E) int, element E) *node[E] {
        if n == nil {
            return &node[E]{value: element}
        }
        sign := cmp(element, n.value)
        switch {
        case sign < 0:
            n.left = n.left.insert(cmp, element)
        case sign > 0:
            n.right = n.right.insert(cmp, element)
        }
        return n
    }

    // Insert inserts element into the tree, if E implements cmp.Ordered.
    func (t *Tree[E]) Insert(element E) {
        t.root = t.root.insert(cmp.Compare[E], element)
    }

    // Insert inserts element into the tree, using the provided comparison function.
    func (t *FuncTree[E]) Insert(element E) {
        t.root = t.root.insert(t.cmp, element)
    }

    // Insert inserts element into the tree, if E implements Comparer[E].
    func (t *MethodTree[E]) Insert(element E) {
        t.root = t.root.insert(E.Compare, element)
    }
`}}

([playground](/play/p/jzmoaH5eaIv))

An important observation here is that the shared implementation (the function-based variant) is not constrained in any way.
It must remain maximally flexible to serve as a common core.
We also do not store the comparison function in a struct field.
Instead, we pass it as a parameter because function arguments are easier for the compiler to analyze than struct fields.

There is still some amount of boilerplate involved, of course.
All the exported implementations need to replicate the full API with slightly different call patterns.
But this part is straightforward to write and to read.

## Combining methods and type sets

We can use our new tree data structure to implement an ordered set, providing element lookup in logarithmic time.
Let's now imagine we need to make lookup run in constant time; we might try to do this by maintaining an ordinary Go map alongside the tree:


{{raw `
    type OrderedSet[E Comparer[E]] struct {
        tree     MethodTree[E] // for efficient iteration in order
        elements map[E]bool    // for (near) constant time lookup
    }

    func (s *OrderedSet[E]) Has(e E) bool {
        return s.elements[e]
    }

    func (s *OrderedSet[E]) Insert(e E) {
        if s.elements == nil {
            s.elements = make(map[E]bool)
        }
        if s.elements[e] {
            return
        }
        s.elements[e] = true
        s.tree.Insert(e)
    }

    func (s *OrderedSet[E]) All() iter.Seq[E] {
        return func(yield func(E) bool) {
            s.tree.root.all(yield)
        }
    }

    func (n *node[E]) all(yield func(E) bool) bool {
        return n == nil || (n.left.all(yield) && yield(n.value) && n.right.all(yield))
    }
`}}

([playground](/play/p/TANUnnSnDqf))

However, compiling this code will produce an error:

> invalid map key type E (missing comparable constraint)

The error message tells us that we need to further constrain our type parameter to be able to use it as a map key.
The `comparable` constraint is a special predeclared constraint that is satisfied by all types for which the equality operators `==` and `!=` are defined.
In Go, that is also the set of types which can be used as keys for the built-in `map` type.

We have three options to add this constraint to our type parameter, all with different tradeoffs:

1.  We can [embed](/ref/spec#Embedded_interfaces) `comparable` into our original `Comparer` definition ([playground](/play/p/g8NLjZCq97q)):

    {{raw `
        type Comparer[E any] interface {
            comparable
            Compare(E) int
        }
    `}}

    This has the downside that it would also make our `Tree` types only usable with types that are `comparable`.
    In general, we do not want to unnecessarily restrict generic types.
2.  We can add a new constraint definition ([playground](/play/p/Z2eg4X8xK5Z)).

    {{raw `
        type Comparer[E any] interface {
            Compare(E) int
        }

        type ComparableComparer[E any] interface {
            comparable
            Comparer[E]
        }
    `}}

    This is tidy, but it introduces a new identifier (`ComparableComparer`) into our API, and naming is hard.
3.  We can add the constraint inline into our more constrained type ([playground](/play/p/ZfggVma_jNc)):

    {{raw `
        type OrderedSet[E interface {
            comparable
            Comparer[E]
        }] struct {
            tree     Tree[E]
            elements map[E]struct{}
        }
    `}}

    This can become a bit hard to read, especially if it needs to happen often.
    It also makes it harder to reuse the constraint in other places.

Which of these to use is a style choice and ultimately up to personal preference.

## (Not) constraining generic interfaces

At this point it is worth discussing constraints on generic interfaces.
You might want to define an interface for a generic container type.
For example, say you have an algorithm that requires a set data structure.
There are many different kinds of set implementations with different tradeoffs.
Defining an interface for the set operations you require can add flexibility to your package, leaving the decision of what tradeoffs are right for the specific application to the user:

{{raw `
    type Set[E any] interface {
        Insert(E)
        Delete(E)
        Has(E) bool
        All() iter.Seq[E]
    }
`}}

A natural question here is what the constraint on this interface should be.
If possible, type parameters on generic interfaces should use `any` as a constraint, allowing arbitrary types.

From our discussions above, the reasons should be clear:
Different concrete implementations might require different constraints.
All the `Tree` types we have examined above, as well as the `OrderedSet` type, can implement `Set` for their element types, even though these types have different constraints.

The point of defining an interface is to leave the implementation up to the user.
Since one cannot predict what kinds of constraints a user may want to impose on their implementation, try to leave any constraints (stronger than `any`) to concrete implementations, not the interfaces.

## Pointer receivers

Let us try to use the `Set` interface in an example.
Consider a function that removes duplicate elements in a sequence:

{{raw `
    // Unique removes duplicate elements from the input sequence, yielding only
    // the first instance of any element.
    func Unique[E comparable](input iter.Seq[E]) iter.Seq[E] {
        return func(yield func(E) bool) {
            seen := make(map[E]bool)
            for v := range input {
                if seen[v] {
                    continue
                }
                if !yield(v) {
                    return
                }
                seen[v] = true
            }
        }
    }
`}}

([playground](/play/p/hsYoFjkU9kA))

This uses a `map[E]bool` as a simple set of `E` elements.
Consequently, it works only for types that are `comparable` and which therefore define the built-in equality operators.
If we want to generalize this to arbitrary types, we need to replace that with a generic set:

{{raw `
    // Unique removes duplicate elements from the input sequence, yielding only
    // the first instance of any element.
    func Unique[E any](input iter.Seq[E]) iter.Seq[E] {
        return func(yield func(E) bool) {
            var seen Set[E]
            for v := range input {
                if seen.Has(v) {
                    continue
                }
                if !yield(v) {
                    return
                }
                seen.Insert(v)
            }
        }
    }
`}}

([playground](/play/p/FZYPNf56nnY))

However, this does not work.
`Set[E]` is an interface type, and the `seen` variable will be initialized to `nil`.
We need to use a concrete implementation of the `Set[E]` interface.
But as we have seen in this post, there is no general implementation of a set that works for `any` element type.

We have to ask the user to provide a concrete implementation we can use, as an extra type parameter:

{{raw `
    // Unique removes duplicate elements from the input sequence, yielding only
    // the first instance of any element.
    func Unique[E any, S Set[E]](input iter.Seq[E]) iter.Seq[E] {
        return func(yield func(E) bool) {
            var seen S
            for v := range input {
                if seen.Has(v) {
                    continue
                }
                if !yield(v) {
                    return
                }
                seen.Insert(v)
            }
        }
    }
`}}

([playground](/play/p/kjkGy5cNz8T))

However, if we instantiate this with our set implementation, we run into another problem:

{{raw `
    // OrderedSet[E] does not satisfy Set[E] (method All has pointer receiver)
    Unique[E, OrderedSet[E]](slices.Values(s))
    // panic: invalid memory address or nil pointer dereference
    Unique[E, *OrderedSet[E]](slices.Values(s))
`}}

The first problem is clear from the error message: Our type constraint says that the type argument for `S` needs to implement the `Set[E]` interface.
And as the methods on `OrderedSet` use a pointer receiver, the type argument also has to be the pointer type.

When trying to do that, we run into the second problem.
This stems from the fact that we declare a variable in the implementation:

{{raw `
    var seen S
`}}

If `S` is `*OrderedSet[E]`, the variable is initialized with `nil`, as before.
Calling `seen.Insert` panics.

If we only have the pointer type, we cannot get a valid variable of the value type.
And if we only have the value type, we cannot call pointer-methods on it.
The consequence is that we need both the value *and* the pointer type.
So we have to introduce an additional type parameter `PS` with a new constraint `PtrToSet`:

{{raw `
    // PtrToSet is implemented by a pointer type implementing the Set[E] interface.
    type PtrToSet[S, E any] interface {
        *S
        Set[E]
    }

    // Unique removes duplicate elements from the input sequence, yielding only
    // the first instance of any element.
    func Unique[E, S any, PS PtrToSet[S, E]](input iter.Seq[E]) iter.Seq[E] {
        return func(yield func(E) bool) {
            // We convert to PS, as only that is constrained to have the methods.
            // The conversion is allowed, because the type set of PS only contains *S.
            seen := PS(new(S))
            for v := range input {
                if seen.Has(v) {
                    continue
                }
                if !yield(v) {
                    return
                }
                seen.Insert(v)
            }
        }
    }
`}}

([playground](/play/p/Kp1jJRVjmYa))

The trick here is the connection of the two type parameters in the function signature via the extra type parameter on the `PtrToSet` interface.
`S` itself is unconstrained, but `PS` must have type `*S` and it must have the methods we need.
So effectively, we are restricting `S` to have some methods, but those methods need to use a pointer receiver.

While the definition of a function with this kind of constraint requires an additional type parameter, importantly this does not complicate code using it:
as long as this extra type parameter is at the end of the type parameter list, it [can be inferred](/blog/type-inference):

{{raw `
    // The third type argument is inferred to be *OrderedSet[int]
    Unique[int, OrderedSet[int]](slices.Values(s))
`}}

This is a general pattern, and worth remembering: for when you encounter it in someone else's work, or when you want to use it in your own.

{{raw `
    func SomeFunction[T any, PT interface{ *T; SomeMethods }]()
`}}

If you have two type parameters, where one is constrained to be a pointer to the other, the constraint ensures that the relevant methods use a pointer receiver.

## Should you constrain to pointer receivers?

At this point, you might feel pretty overwhelmed.
This is rather complicated and it seems unreasonable to expect every Go programmer to understand what is going on in this function signature.
We also had to introduce yet more names into our API.
When people cautioned against adding generics to Go in the first place, this is one of the things they were worried about.

So if you find yourself entangled in these problems, it is worth taking a step back.
We can often avoid this complexity by thinking about our problem in a different way.
In this example, we built a function that takes an `iter.Seq[E]` and returns an `iter.Seq[E]` with the unique elements.
But to do the deduplication, we needed to collect the unique elements into a set.
And as this requires us to allocate the space for the entire result, we do not really benefit from representing the result as a stream.

If we rethink this problem, we can avoid the extra type parameter altogether by using `Set[E]` as a regular interface value:

{{raw `
    // InsertAll adds all unique elements from seq into set.
    func InsertAll[E any](set Set[E], seq iter.Seq[E]) {
        for v := range seq {
            set.Insert(v)
        }
    }
`}}

([playground](/play/p/woZcHodAgaa))

By using `Set` as a plain interface type, it is clear that the caller has to provide a valid value of their concrete implementation.
This is a very common pattern.
And if they need an `iter.Seq[E]`, they can simply call `All()` on the `set` to obtain one.

This complicates things for callers slightly, but it has another advantage over the constraint to pointer receivers:
remember that we started with a `map[E]bool` as a simple set type.
It is easy to implement the `Set[E]` interface on that basis:

{{raw `
    type HashSet[E comparable] map[E]bool

    func (s HashSet[E]) Insert(v E)       { s[v] = true }
    func (s HashSet[E]) Delete(v E)       { delete(s, v) }
    func (s HashSet[E]) Has(v E) bool     { return s[v] }
    func (s HashSet[E]) All() iter.Seq[E] { return maps.Keys(s) }
`}}

([playground](/play/p/KPPpWa7M93d))

This implementation does not use pointer receivers.
So while this is perfectly valid, it would not be usable with the complicated constraint to pointer receivers.
But it works fine with our `InsertAll` version.
As with many constraints, enforcing that methods use a pointer receiver might actually be overly restrictive for many practical use cases.

## Conclusion

I hope this illustrates some of the patterns and trade-offs that type parameters on interfaces enable.
It is a powerful tool, but it also comes with a cost.
The primary take-aways are:

1. Use generic interfaces to express constraints on the receiver by using them self-referentially.
2. Use them to create constrained relationships between different type parameters.
3. Use them to abstract over different implementations with different kinds of constraints.
4. When you find yourself in a situation where you need to constrain to pointer receivers, consider whether you can refactor your code to avoid the extra complexity. See ["Should you constrain to pointer receivers?"](#should-you-constrain-to-pointer-receivers).

As always, do not over-engineer things: a less flexible but simpler and more readable solution may ultimately be the wiser choice.
