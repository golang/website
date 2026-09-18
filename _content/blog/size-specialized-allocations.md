---
title: Size-Specialized Memory Allocation
date: 2026-09-16
by:
- Michael Matloob
summary: Go 1.27 improves performance of small allocations using size-specialized allocation functions.
---

Go 1.27 includes faster memory allocation for allocations of 80 bytes or fewer.
Allocations can be up to 20-30% faster, making allocation-heavy programs up to 1% faster.
The Go runtime improves the performance of those allocations by adding specialized
functions that are used to allocate certain sizes. These specialized functions
can then make certain assumptions that make them faster and easier to optimize.
This blog post will explain how this works and how it makes your programs faster.

Heap allocations are created by the runtime's `mallocgc` function, which requires the size
of the allocation and whether it contains pointers. When the compiler determines that an
object escapes to the heap or otherwise needs to be dynamically
allocated, it inserts a call to `newobject`, which is a simple wrapper function that extracts
the size of the object and whether it contains pointers and passes those to `mallocgc`.

These two pieces of information determine
most of the work that the allocator needs to do. The size is important because the allocator
defines ranges of sizes called "size classes". For most allocations that are not too large
and not too small, it will return a block of memory from a list of free objects
that are all sized to the maximum size of the size class. So, for instance, size class 3
is 17-24 bytes. So whether you allocate 17 bytes or 24 bytes, the allocator will give
you the next free 24 byte object available in its list of 24 byte objects.
Below is a table of the size ranges of each of the size classes up to 80 bytes.
There are separate sets of free lists, which we call spans, depending on whether the allocation
has pointers, because of the bookkeeping we need to do for the garbage collector.

| Size class | Range of sizes |
| :--- | :--- |
| `1` | 1-8 bytes |
| `2` | 9-16 bytes |
| `3` | 17-24 bytes |
| `4` | 25-32 bytes |
| `5` | 33-48 bytes |
| `6` | 49-64 bytes |
| `7` | 65-80 bytes |

Because the size class and whether the allocation contains pointers determine which span
we allocate from, the allocator has a type called the span class, whose value encodes both:
it's defined as `sizeClass<<1 | noPointers`. When implementing size-specialized malloc
it was clear that it made the most sense to specialize on these span classes because much of
the behavior of the allocator was determined by the span class.

We generate a specialized `mallocgc` variant for every span class. For instance, the function
that allocates pointer-free objects in size class 3 is called `mallocgcSmallNoScanSC3`.
`Small` means not tiny or large, `NoScan` means it has no pointers, and `SC3` means size class 3
(17-24 bytes). Specialized functions are designed to be as simple as possible and cannot handle
every corner case during allocation. When they detect such a case, such as when the GC is active,
they fall back to a more generic allocation routine.

So we end up with a new specialized `mallocgc` function for the tiny case (always no pointers)
and one for each non-tiny span class: the non-pointer spans of size class 2 and above, and
the pointer spans of size class 1 and above. These specialized `mallocgc` variants themselves
are faster than calling `mallocgc`, but the benefits of specializing decrease as the sizes
get larger: the allocation often gets dominated by needing to clear the memory, and at some point
the rest of the work of allocation becomes negligible. And it's not enough to just be faster than
`mallocgc`! With size-specialized allocation, when the compiler knows which span class it needs to
allocate for, it can directly insert a call to the specialized function instead of the call to
`newobject`. But the compiler often doesn't know which span class is being allocated when code
is generated: think slices of dynamic length.
Since the compiler can't determine the size of the allocation, it will keep the call to `mallocgc`,
and `mallocgc` itself has to determine whether a specialized function is available, which one it is,
and then needs to call it. So the performance of the specialized function has to be high enough that even
with the overhead of a dynamic call it's still faster.

This overhead issue isn't the only problem. If it were, then we could create larger specialized
functions, and have them reserved for the compiler to insert them when the size is known at compile-time.
And if we didn't call them dynamically we wouldn't have to pay the dynamic overhead. But each specialized
function that's added increases the sizes of the executables produced by the compiler, and, more importantly,
takes up precious instruction cache space. The single `mallocgc` function is often in the instruction
cache because of how frequently allocations happen. If we have too many specialized functions and
they're not available in the cache, the overhead to retrieve the specialized code into the cache
can cancel out any benefits. And the more specialized allocation code that's in the icache, the more
it crowds out the user code, making it slower to fetch and run user code. Doing a bunch of benchmarks
stopping at different size classes, we determined that stopping at 80 bytes was the sweet spot.

Of course, adding more specialized functions means more code to maintain. If each function were
hand-written, the code in the specialized functions could easily drift apart and go out of sync.
To help with this, we broke out the common parts of the specialized
functions and wrote an inliner using the standard library [`go/ast`](/pkg/go/ast) package to parse and format and
the [`golang.org/x/tools/go/ast/astutil`](/pkg/golang.org/x/tools/go/ast/astutil) package to manipulate the ASTs. The common parts of the functions
are all written in standard Go code that is built and typechecked with the rest of the runtime to
allow our tools to catch issues, but they are mostly just stubs for the inliner.

So we were able to measure improvements in the size-specialized allocation functions, but why are they actually
faster? The most apparent optimization is in clearing memory. The memory returned by `mallocgc`
does not always need to be zeroed, but it often does, and doing so can often take up most of the
time of the allocation. The memory clearing function
`memclrNoHeapPointers` is written in highly-optimized assembly, but we could do better
for the smaller clears. In a specialized function, if the size of the clear is constant, the
compiler can replace calls to `memclrNoHeapPointers` with code to directly produce the instructions
to clear the memory. The allocation can then skip a function call and some branches. This makes
a difference for the really small allocations, but as they get bigger, the function call overhead
becomes negligible.

While faster memory clearing provides the biggest improvement, there are some other tricks that are also possible
in the specialized functions: Since they are specialized per span-class, the function doesn't need
to calculate the span class when retrieving the span. And because the size of the allocation
is a constant, the compiler is able to do some optimizations to speed up the bookkeeping necessary for
the allocation. One such case is with marking where the pointers are in the allocated memory.
The specialized functions could also manually inline several of their helper functions. The Go compiler
can inline code, but it avoids inlining functions that it considers too large. We can override that
in the generated code by inserting the function bodies into the callers. With the generator we can produce
copies of each of the bodies without needing to worry about each of the copies drifting. And we could move code that
handled less common cases, such as the runtime debugging flags, into the slow path functions to
make the specialized functions smaller.

While we hope this explanation of size-specialized allocation is interesting, you as a Go programmer
don't need to think about any of this when writing your code. Memory allocations will just be
a little faster, with the biggest benefits going to some of the most common allocation sizes,
especially the 16 and 24 byte allocations. These allocations are some of the most common because
they consist of two or three 64-bit values, so they include allocations for things such as interface
values and strings which have two values, or slices which have three. We spent a lot of time tuning
the behavior of size specialized malloc and making sure the impact of the instruction cache effects
will be minimal: we were originally planning to release size-specialized allocation in Go 1.26 but
decided to wait an extra release to do extra tuning and cut down the additional code size as much
as we could.

All you need to do to get the improved performance from size specialized allocations in your
programs is to build them with Go 1.27. If you're interested in more concrete actions you can
take to improve memory allocation and garbage collection performance, please read the
[Go Garbage Collector Optimization Guide](/doc/gc-guide#Optimization_guide).

Although we are confident that size-specialized allocation should not cause regressions in your
code, if necessary, you can build your programs using `GOEXPERIMENT=nosizespecializedmalloc` to disable it.
If you do need to do this to solve an issue you're experiencing with size-specialized allocation,
please file an issue at [go.dev/issue/new](/issue/new) so we can investigate it.

