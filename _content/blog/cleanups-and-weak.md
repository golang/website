---
title: "From unique to cleanups and weak: new low-level tools for efficiency"
date: 2025-03-06
by:
- Michael Knyszek
tags:
- weak
- cleanup
- finalizer
summary: Weak pointers and better finalization in Go 1.24.
---

In [last year's blog post](/blog/unique) about the `unique` package, we alluded
to some new features then in proposal review, and we're excited to share that as
of Go 1.24 they are now available to all Go developers.
These new features are [the `runtime.AddCleanup`
function](https://pkg.go.dev/runtime#AddCleanup), which queues up a function to
run when an object is no longer reachable, and [the `weak.Pointer`
type](https://pkg.go.dev/weak#Pointer), which safely points to an object without
preventing it from being garbage collected.
Together, these two features are powerful enough to build your own `unique`
package!
Let's dig into what makes these features useful, and when to use them.

Note: these new features are advanced features of the garbage collector.
If you're not already familiar with basic garbage collection concepts, we
strongly recommend reading the introduction of our [garbage collector
guide](/doc/gc-guide#Introduction).

## Cleanups

If you've ever used a finalizer, then the concept of a cleanup will be
familiar.
A finalizer is a function, associated with an allocated object by [calling
`runtime.SetFinalizer`](https://pkg.go.dev/runtime#SetFinalizer), that is later
called by the garbage collector some time after the object becomes unreachable.
At a high level, cleanups work the same way.

Let's consider an application that makes use of a memory-mapped file, and see
how cleanups can help.

```
//go:build unix

type MemoryMappedFile struct {
	data []byte
}

func NewMemoryMappedFile(filename string) (*MemoryMappedFile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get the file's info; we need its size.
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// Extract the file descriptor.
	conn, err := f.SyscallConn()
	if err != nil {
		return nil, err
	}
	var data []byte
	connErr := conn.Control(func(fd uintptr) {
		// Create a memory mapping backed by this file.
		data, err = syscall.Mmap(int(fd), 0, int(fi.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	})
	if connErr != nil {
		return nil, connErr
	}
	if err != nil {
		return nil, err
	}
	mf := &MemoryMappedFile{data: data}
	cleanup := func(data []byte) {
		syscall.Munmap(data) // ignore error
	}
	runtime.AddCleanup(mf, cleanup, data)
	return mf, nil
}
```

A memory-mapped file has its contents mapped to memory, in this case, the
underlying data of a byte slice.
Thanks to operating-system magic, reads and writes to the byte slice directly
access the contents of the file.
With this code, we can pass around a `*MemoryMappedFile`, and when it's
no longer referenced, the memory mapping we created will get cleaned up.

Notice that `runtime.AddCleanup` takes three arguments: the address of a
variable to attach the cleanup to, the cleanup function itself, and an argument
to the cleanup function.
A key difference between this function and `runtime.SetFinalizer` is that the
cleanup function takes a different argument than the object we're attaching the
cleanup to.
This change fixes some problems with finalizers.

It's no secret that [finalizers
are difficult to use correctly](/doc/gc-guide#Common_finalizer_issues).
For example, objects to which finalizers are attached must not be involved
in any reference cycles (even a pointer to itself is too much!), otherwise the
object will never be reclaimed and the finalizer will never run, causing a
leak.
Finalizers also significantly delay reclamation of memory.
It takes at a minimum two full garbage collection cycles to reclaim the
memory for a finalized object: one to determine that it's unreachable, and
another to determine that it's still unreachable after the finalizer
executes.

The problem is that finalizers [resurrect the objects they're attached
to](https://en.wikipedia.org/wiki/Object_resurrection).
The finalizer doesn't run until the object is unreachable, at which point it is
considered "dead."
But since the finalizer is called with a pointer to the object, the garbage
collector must prevent the reclamation of that object's memory, and instead must
generate a new reference for the finalizer, making it reachable, or "live," once
more.
That reference may even remain after the finalizer returns, for example if the
finalizer writes it to a global variable or sends it across a channel.
Object resurrection is problematic because it means the object, and everything
it points to, and everything those objects point to, and so on, is reachable,
even if it would otherwise have been collected as garbage.

We solve both of these problems by not passing the original object to the
cleanup function.
First, the values the object refers to don't need to be kept specially
reachable by the garbage collector, so the object can still be reclaimed even
if it's involved in a cycle.
Second, since the object is not needed for the cleanup, its memory can be
reclaimed immediately.

## Weak pointers

Returning to our memory-mapped file example, suppose we notice that our program
frequently maps the same files over and over, from different goroutines that are
unaware of each other.
This is fine from a memory perspective, since all these mappings will share
physical memory, but it results in lots of unnecessary system calls to map and
unmap the file.
This is especially bad if each goroutine reads only a small section of each
file.

So, let's deduplicate the mappings by filename.
(Let's assume that our program only reads from the mappings, and the files
themselves are never modified or renamed once created.
Such assumptions are reasonable for system font files, for example.)

We could maintain a map from filename to memory mapping, but then it becomes
unclear when it's safe to remove entries from that map.
We could *almost* use a cleanup, if it weren't for the fact that the map entry
itself will keep the memory-mapped file object alive.

Weak pointers solve this problem.
A weak pointer is a special kind of pointer that the garbage collector ignores
when deciding whether an object is reachable.
Go 1.24's [new weak pointer type,
`weak.Pointer`](https://pkg.go.dev/weak#Pointer), has a `Value` method that
returns either a real pointer if the object is still reachable, or `nil` if it
is not.

If we instead maintain a map that only *weakly* points to the memory-mapped
file, we can clean up the map entry when nobody's using it anymore!
Let's see what this looks like.

```
var cache sync.Map // map[string]weak.Pointer[MemoryMappedFile]

func NewCachedMemoryMappedFile(filename string) (*MemoryMappedFile, error) {
	var newFile *MemoryMappedFile
	for {
		// Try to load an existing value out of the cache.
		value, ok := cache.Load(filename)
		if !ok {
			// No value found. Create a new mapped file if needed.
			if newFile == nil {
				var err error
				newFile, err = NewMemoryMappedFile(filename)
				if err != nil {
					return nil, err
				}
			}

			// Try to install the new mapped file.
			wp := weak.Make(newFile)
			var loaded bool
			value, loaded = cache.LoadOrStore(filename, wp)
			if !loaded {
				runtime.AddCleanup(newFile, func(filename string) {
					// Only delete if the weak pointer is equal. If it's not, someone
					// else already deleted the entry and installed a new mapped file.
					cache.CompareAndDelete(filename, wp)
				}, filename)
				return newFile, nil
			}
			// Someone got to installing the file before us.
			//
			// If it's still there when we check in a moment, we'll discard newFile
			// and it'll get cleaned up by garbage collector.
		}

		// See if our cache entry is valid.
		if mf := value.(weak.Pointer[MemoryMappedFile]).Value(); mf != nil {
			return mf, nil
		}

		// Discovered a nil entry awaiting cleanup. Eagerly delete it.
		cache.CompareAndDelete(filename, value)
	}
}
```

This example is a little complicated, but the gist is simple.
We start with a global concurrent map of all the mapped files we made.
`NewCachedMemoryMappedFile` consults this map for an existing mapped
file, and if that fails, creates and tries to insert a new mapped file.
This could of course fail as well since we're racing with other insertions, so
we need to be careful about that too, and retry.
(This design has a flaw in that we might wastefully map the same file multiple
times in a race, and we'll have to throw it away via the cleanup added by
`NewMemoryMappedFile`.
This is probably not a big deal most of the time.
Fixing it is left as an exercise for the reader.)

Let's look at some useful properties of weak pointers and cleanups exploited by
this code.

Firstly, notice that weak pointers are comparable.
Not only that, weak pointers have a stable and independent identity, which
remains even after the objects they point to are long gone.
This is why it is safe for the cleanup function to call `sync.Map`'s
`CompareAndDelete`, which compares the `weak.Pointer`, and a crucial reason
this code works at all.

Secondly, observe that we can add multiple independent cleanups to a single
`MemoryMappedFile` object.
This allows us to use cleanups in a composable way and use them to build
generic data structures.
In this particular example, it might be more efficient to combine
`NewCachedMemoryMappedFile` with `NewMemoryMappedFile` and
have them share a cleanup.
However, the advantage of the code we wrote above is that it can be rewritten
in a generic way!

```
type Cache[K comparable, V any] struct {
	create func(K) (*V, error)
	m     sync.Map
}

func NewCache[K comparable, V any](create func(K) (*V, error)) *Cache[K, V] {
	return &Cache[K, V]{create: create}
}

func (c *Cache[K, V]) Get(key K) (*V, error) {
	var newValue *V
	for {
		// Try to load an existing value out of the cache.
		value, ok := cache.Load(key)
		if !ok {
			// No value found. Create a new mapped file if needed.
			if newValue == nil {
				var err error
				newValue, err = c.create(key)
				if err != nil {
					return nil, err
				}
			}

			// Try to install the new mapped file.
			wp := weak.Make(newValue)
			var loaded bool
			value, loaded = cache.LoadOrStore(key, wp)
			if !loaded {
				runtime.AddCleanup(newValue, func(key K) {
					// Only delete if the weak pointer is equal. If it's not, someone
					// else already deleted the entry and installed a new mapped file.
					cache.CompareAndDelete(key, wp)
				}, key)
				return newValue, nil
			}
		}

		// See if our cache entry is valid.
		if mf := value.(weak.Pointer[V]).Value(); mf != nil {
			return mf, nil
		}

		// Discovered a nil entry awaiting cleanup. Eagerly delete it.
		cache.CompareAndDelete(key, value)
	}
}
```

## Caveats and future work

Despite our best efforts, cleanups and weak pointers can still be error-prone.
To guide those considering using finalizers, cleanups, and weak pointers, we
recently updated the [guide to the garbage
collector](/doc/gc-guide#Finalizers_cleanups_and_weak_pointers) with some advice
about using these features.
Take a look next time you reach for them, but also carefully consider whether
you need to use them at all.
These are advanced tools with subtle semantics and, as the guide says, most
Go code benefits from these features indirectly, not from using them directly.
Stick to the use-cases where these features shine, and you'll be alright.

For now, we'll call out some of the issues that you are more likely to run into.

First, the object the cleanup is attached to must be reachable from neither
the cleanup function (as a captured variable) nor the argument to the cleanup
function.
Both of these situations result in the cleanup never executing.
(In the special case of the cleanup argument being exactly the pointer passed
to `runtime.AddCleanup`, `runtime.AddCleanup` will panic, as a signal to the
caller that they should not use cleanups the same way as finalizers.)

Second, when weak pointers are used as map keys, the weakly referenced object
must not be reachable from the corresponding map value, otherwise the object
will continue to remain live.
This may seem obvious when deep inside of a blog post about weak pointers, but
it's an easy subtlety to miss.
This problem inspired the entire concept of an
[ephemeron](https://en.wikipedia.org/wiki/Ephemeron) to resolve it, which is a
potential future direction.

Thirdly, a common pattern with cleanups is that a wrapper object is needed, like
we see here with our `MemoryMappedFile` example.
In this particular case, you could imagine the garbage collector directly
tracking the mapped memory region and passing around the inner `[]byte`.
Such functionality is possible future work, and an API for it has been recently
[proposed](/issue/70224).

Lastly, both weak pointers and cleanups are inherently non-deterministic, their
behavior depending intimately on the design and dynamics of the garbage
collector.
The documentation for cleanups even permits the garbage collector never to run
cleanups at all.
Effectively testing code that uses them can be tricky, but [it is
possible](/doc/gc-guide#Testing_object_death).

## Why now?

Weak pointers have been brought up as a feature for Go since nearly the
beginning, but for years were not prioritized by the Go team.
One reason for that is that they are subtle, and the design space of weak
pointers is a minefield of decisions that can make them even harder to use.
Another is that weak pointers are a niche tool, while simultaneously adding
complexity to the language.
We already had experience with how painful `SetFinalizer` could be to use.
But there are some useful programs that are not expressible without them, and
the `unique` package and the reasons for its existence really emphasized that.

With generics, the hindsight of finalizers, and insights from all the great work
since done by teams in other languages like C# and Java, the designs for weak
pointers and cleanups came together quickly.
The desire to use weak pointers with finalizers raised additional questions,
and so the design for `runtime.AddCleanup` quickly came together as well.

## Acknowledgements

I want to thank everyone in the community who contributed feedback on the
proposal issues and filed bugs when the features became available.
I also want to thank David Chase for thoroughly thinking through weak
pointer semantics with me, and I want to thank him, Russ Cox, and Austin
Clements for their help with the design of `runtime.AddCleanup`.
I want to thank Carlos Amedee for his work on getting `runtime.AddCleanup`
implemented, polished, landed for Go 1.24.
And finally I want to thank Carlos Amedee and Ian Lance Taylor for their work
replacing `runtime.SetFinalizer` with `runtime.AddCleanup` throughout the
standard library for Go 1.25.
