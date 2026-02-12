---
title: Allocating on the Stack
date: 2026-02-27
by:
- Keith Randall
summary: A description of some of the recent changes to do allocations on the stack instead of the heap.
template: true
---

We're always looking for ways to make Go programs faster. In the last
2 releases, we have concentrated on mitigating a particular source of
slowness, heap allocations. Each time a Go program allocates memory
from the heap, there's a fairly large chunk of code that needs to run
to satisfy that allocation. In addition, heap allocations present
additional load on the garbage collector.  Even with recent
enhancements like [Green Tea](/blog/greenteagc), the garbage collector
still incurs substantial overhead.

So we've been working on ways to do more allocations on the stack
instead of the heap.  Stack allocations are considerably cheaper to
perform (sometimes completely free).  Moreover, they present no load
to the garbage collector, as stack allocations can be collected
automatically together with the stack frame itself. Stack allocations
also enable prompt reuse, which is very cache friendly.

## Stack allocation of constant-sized slices

Consider the task of building a slice of tasks to process:
{{raw `
	func process(c chan task) {
		var tasks []task
		for t := range c {
			tasks = append(tasks, t)
		}
		processAll(tasks)
	}
`}}

Let's walk through what happens at runtime when pulling tasks from the
channel `c` and adding them to the slice `tasks`.

On the first loop iteration, there is no backing store for `tasks`, so
`append` has to allocate one. Because it doesn't know how big the
slice will eventually be, it can't be too aggressive. Currently, it
allocates a backing store of size 1.

On the second loop iteration, the backing store now exists, but it is
full. `append` again has to allocate a new backing store, this time of
size 2. The old backing store of size 1 is now garbage.

On the third loop iteration, the backing store of size 2 is
full. `append` *again* has to allocate a new backing store, this time
of size 4. The old backing store of size 2 is now garbage.

On the fourth loop iteration, the backing store of size 4 has only 3
items in it. `append` can just place the item in the existing backing
store and bump up the slice length. Yay! No call to the allocator for
this iteration.

On the fifth loop iteration, the backing store of size 4 is full, and
`append` again has to allocate a new backing store, this time of size
8.

And so on. We generally double the size of the allocation each time it
fills up, so we can eventually append most new tasks to the slice
without allocation. But there is a fair amount of overhead in the
"startup" phase when the slice is small. During this startup phase we
spend a lot of time in the allocator, and produce a bunch of garbage,
which seems pretty wasteful. And it may be that in your program, the
slice never really gets large. This startup phase may be all you ever
encounter.

If this code was a really hot part of your program, you might be
tempted to start the slice out at a larger size, to avoid all of these
allocations.

{{raw `
	func process2(c chan task) {
		tasks := make([]task, 0, 10) // probably at most 10 tasks
		for t := range c {
			tasks = append(tasks, t)
		}
		processAll(tasks)
	}
`}}

This is a reasonable optimization to do. It is never incorrect; your
program still runs correctly. If the guess is too small, you get
allocations from `append` as before. If the guess is too large, you
waste some memory.

If your guess for the number of tasks was a good one, then there's
only one allocation site in this program. The `make` call allocates a
slice backing store of the correct size, and `append` never has to do
any reallocation.

The surprising thing is that if you benchmark this code with 10
elements in the channel, you'll see that you didn't reduce the number
of allocations to 1, you reduced the number of allocations to 0!

The reason is that the compiler decided to allocate the backing store
on the stack. Because it knows what size it needs to be (10 times the
size of a task) it can allocate storage for it in the stack frame of
`process2` instead of on the heap[<sup>1</sup>](#footnotes).  Note
that this depends on the fact that the backing store does not [escape
to the heap](/doc/gc-guide#Escape_analysis) inside of `processAll`.

## Stack allocation of variable-sized slices

But of course, hard coding a size guess is a bit rigid.
Maybe we can pass in an estimated length?

{{raw `
	func process3(c chan task, lengthGuess int) {
		tasks := make([]task, 0, lengthGuess)
		for t := range c {
			tasks = append(tasks, t)
		}
		processAll(tasks)
	}
`}}

This lets the caller pick a good size for the `tasks` slice, which may
vary depending on where this code is being called from.

Unfortunately, in Go 1.24 the non-constant size of the backing store
means the compiler can no longer allocate the backing store on the
stack.  It will end up on the heap, converting our 0-allocation code
to 1-allocation code. Still better than having `append` do all the
intermediate allocations, but unfortunate.

But never fear, Go 1.25 is here!

Imagine you decide to do the following, to get the stack allocation
only in cases where the guess is small:

{{raw `
	func process4(c chan task, lengthGuess int) {
		var tasks []task
		if lengthGuess <= 10 {
			tasks = make([]task, 0, 10)
		} else {
			tasks = make([]task, 0, lengthGuess)
		}
		for t := range c {
			tasks = append(tasks, t)
		}
		processAll(tasks)
	}
`}}

Kind of ugly, but it would work. When the guess is small, you use a
constant size `make` and thus a stack-allocated backing store, and
when the guess is larger you use a variable size `make` and allocate
the backing store from the heap.

But in Go 1.25, you don't need to head down this ugly road. The Go
1.25 compiler does this transformation for you!  For certain slice
allocation locations, the compiler automatically allocates a small
(currently 32-byte) slice backing store, and uses that backing store
for the result of the `make` if the size requested is small
enough. Otherwise, it uses a heap allocation as normal.

In Go 1.25, `process3` performs zero heap allocations, if
`lengthGuess` is small enough that a slice of that length fits into 32
bytes. (And of course that `lengthGuess` is a correct guess for how
many items are in `c`.)

We're always improving the performance of Go, so upgrade to the latest
Go release and [be
surprised](https://youtu.be/FUm0pfgWehI?si=QRTt_JYwr-cRHDNJ&t=960) by
how much faster and memory efficient your program becomes!

## Stack allocation of append-allocated slices

Ok, but you still don't want to have to change your API to add this
weird length guess. Anything else you could do?

Upgrade to Go 1.26!

{{raw `
	func process(c chan task) {
		var tasks []task
		for t := range c {
			tasks = append(tasks, t)
		}
		processAll(tasks)
	}
`}}

In Go 1.26, we allocate the same kind of small, speculative backing
store on the stack, but now we can use it directly at the `append`
site.

On the first loop iteration, there is no backing store for `tasks`, so
`append` uses a small, stack-allocated backing store as the first
allocation. If, for instance, we can fit 4 `task`s in that backing store,
the first `append` allocates a backing store of length 4 from the stack.

The next 3 loop iterations append directly to the stack backing store,
requiring no allocation.

On the 4th iteration, the stack backing store is finally full and we
have to go to the heap for more backing store. But we have avoided
almost all of the startup overhead described earlier in this article.
No heap allocations of size, 1, 2, and 4, and none of the garbage that
they eventually become. If your slices are small, maybe you will never
have a heap allocation.

## Stack allocation of append-allocated escaping slices

Ok, this is all good when the `tasks` slice doesn't escape. But what if
I'm returning the slice? Then it can't be allocated on the stack, right?

Right! The backing store for the slice returned by `extract` below
can't be allocated on the stack, because the stack frame for `extract`
disappears when `extract` returns.

{{raw `
	func extract(c chan task) []task {
		var tasks []task
		for t := range c {
			tasks = append(tasks, t)
		}
		return tasks
	}
`}}

But you might think, the *returned* slice can't be allocated on the
stack. But what about all those intermediate slices that just become
garbage? Maybe we can allocate those on the stack?

{{raw `
	func extract2(c chan task) []task {
		var tasks []task
		for t := range c {
			tasks = append(tasks, t)
		}
		tasks2 := make([]task, len(tasks))
		copy(tasks2, tasks)
		return tasks2
	}
`}}

Then the `tasks` slice never escapes `extract2`. It can benefit from
all of the optimizations described above. Then at the very end of
`extract2`, when we know the final size of the slice, we do one heap
allocation of the required size, copy our `task`s into it, and return
the copy.

But do you really want to write all that additional code? It seems
error prone. Maybe the compiler can do this transformation for us?

In Go 1.26, it can!

For escaping slices, the compiler will transform the original `extract`
code to something like this:

{{raw `
	func extract3(c chan task) []task {
		var tasks []task
		for t := range c {
			tasks = append(tasks, t)
		}
		tasks = runtime.move2heap(tasks)
		return tasks
	}
`}}

`runtime.move2heap` is a special compiler+runtime function that is the
identity function for slices that are already allocated in the heap.
For slices that are on the stack, it allocates a new slice on the
heap, copies the stack-allocated slice to the heap copy, and returns
the heap copy.

This ensures that for our original `extract` code, if the number of
items fits in our small stack-allocated buffer, we perform exactly 1
allocation of exactly the right size. If the number of items exceeds
the capacity our small stack-allocated buffer, we do our normal
doubling-allocation once the stack-allocated buffer overflows.

The optimization that Go 1.26 does is actually better than the
hand-optimized code, because it does not require the extra
allocation+copy that the hand-optimized code always does at the end.
It requires the allocation+copy only in the case that we've exclusively
operated on a stack-backed slice up to the return point.

We do pay the cost for a copy, but that cost is almost completely
offset by the copies in the startup phase that we no longer have to
do. (In fact, the the new scheme at worst has to copy one more element
than the old scheme.)

## Wrapping up

Hand optimization can still be beneficial, especially if you have a
good estimate of the slice size ahead of time. But hopefully the
compiler will now catch a lot of the simple cases for you and allow
you to focus on the remaining ones that really matter.

There are a lot of details that the compiler needs to ensure to get
all these optimizations right. If you think that one of these
optimizations is causing correctness or (negative) performance issues
for you, you can turn them off with
`-gcflags=all=-d=variablemakehash=n`. If turning these optimizations
off helps, please [file an issue](/issue/new) so we can investigate.

## Footnotes

<sup>1</sup> Go stacks do not have any `alloca`-style mechanism for
dynamically-sized stack frames. All Go stack frames are constant
sized.
