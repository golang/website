---
title: Goroutine Leak Profiles
date: 2026-03-03
by:
- Vlad Saioc
tags:
- pprof
- profiles
- performance
summary: Go 1.26 includes new experimental goroutine leak profiles.
---

## Background

Goroutine leaks are a common concurrency error that occurs when misusing
blocking concurrency features in Go.
A goroutine is _leaked_ once it is permanently blocked at an operation over one or
more concurrency primitives, i.e., it cannot be unblocked regardless of the
execution moving forward.

This behavior is clearly undesirable, especially in long-running systems,
where the accumulation of goroutine leaks degrades performance.
It consumes a significant amount of memory, in the form of the goroutine stacks
and any heap resources they may transitively point to, up to the point of
out-of-memory exceptions.
CPU performance is likewise degraded, through the unnecessary burdening of the
garbage collector with the task of marking wasted memory.

While the Go runtime is equipped to signal global deadlocks if all
goroutines are simultaneously blocked,
goroutine leaks have, so far, been difficult to detect and diagnose.
Existing goroutine profiles take a snapshot of all goroutines, but
do not distinguish between leaked goroutines and those which are
blocked legitimately.

To compensate, Go 1.26 introduces experimental, specialized goroutine leak profiles.
In the following article, we showcase their use, examples of leaks that they can detect
and the underlying implementation and trade-offs.
While we will explain nuances, some familiarity with basic Go
[concurrency features](/tour/concurrency/1) is expected.

## Example: A common goroutine leak

Let's look at a realistic example of code that contains a goroutine leak.
Consider a function that processes work items in parallel:

```go
type result struct {
	res workResult
	err error
}

func processWorkItems(ws []workItem) ([]workResult, error) {
	// Process work items in parallel, aggregating results in ch.
	ch := make(chan result)
	for _, w := range ws {
		go func() {
			res, err := processWorkItem(w)
			ch <- result{res, err}
		}()
	}

	// Collect the results from ch, or return an error if one is found.
	var results []workResult
	for range len(ws) {
		r := <-ch
		if r.err != nil {
			// This early return may cause goroutine leaks.
			return nil, r.err
		}
		results = append(results, r.res)
	}
	return results, nil
}
```
Because `ch` is an unbuffered channel, each worker goroutine blocks when sending its result until the main goroutine receives from the channel.
If `processWorkItems` returns early due to an error, the receiving loop terminates, and all remaining worker goroutines will block forever trying to send to `ch`.

## Enabling goroutine leak profiles

Goroutine leak profiles are available as an experiment in Go 1.26.
To enable it, build your program with:

```
$ GOEXPERIMENT=goroutineleakprofile go build .
```

Once enabled, the profile becomes available through the [`runtime/pprof`](/pkg/runtime/pprof) package, as the `goroutineleak` profile type, or by exposing an HTTP endpoint with [`net/http/pprof`](/pkg/net/http/pprof).

### Example set up

In the following section, we demonstrate how to use the goroutine leak profile to
detect the leak in the example above.

Create a simple program that exhibits a leak in `main.go`:

```go
package main

import (
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type workItem int
type workResult int

func processWorkItem(w workItem) (workResult, error) {
	time.Sleep(10 * time.Millisecond)
	if w == 5 {
		return 0, errors.New("simulated error")
	}
	return workResult(w * 2), nil
}

type result struct {
	res workResult
	err error
}

func processWorkItems(ws []workItem) ([]workResult, error) {
	ch := make(chan result)
	for _, w := range ws {
		w := w // capture for closure
		go func() {
			res, err := processWorkItem(w)
			ch <- result{res, err}
		}()
	}

	var results []workResult
	for range len(ws) {
		r := <-ch
		if r.err != nil {
			return nil, r.err
		}
		results = append(results, r.res)
	}
	return results, nil
}

func main() {
	// Start pprof server
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Repeatedly trigger the leak
	for {
		items := []workItem{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		_, err := processWorkItems(items)
		if err != nil {
			log.Printf("Error processing items: %v", err)
		}

		time.Sleep(time.Second)
	}
}
```

Build and run with the experiment enabled:

```
$ GOEXPERIMENT=goroutineleakprofile go build -o leaky
$ ./leaky
```

### Collecting the profile

It won't take long for the program to start accumulating
leaks, which you can then view by using the web UI
at http://localhost:6060/debug/pprof.

Alternatively, you can collect the goroutine
leak profile using `curl`, and then examine it with `go tool pprof`:
```
$ curl http://localhost:6060/debug/pprof/goroutineleak > leak.prof
$ go tool pprof leak.prof
Type: goroutineleak
Time: 2026-03-01 13:19:49 UTC
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list processWorkItems
Total: 116
ROUTINE ======================== main.processWorkItems.func1 in .../main.go
         0        116 (flat, cum)   100% of Total
         .          .     31:           go func() {
         .          .     32:                   res, err := processWorkItem(w)
         .        116     33:                   ch <- result{res, err}
         .          .     34:           }()
```
The profile clearly shows that 116 goroutines are leaked at
`ch <- result{res, err}` (line 33), pinpointing the culprit operation.
Notably, the longer the program is running, the larger the number of leaked
goroutines.

### Addressing the leak

This leak can be simply fixed by giving `ch` a **buffer**:
```go
ch := make(chan result, len(ws))
```
This allows all the work item goroutines to send a message without blocking
in the event of a premature return of `processWorkItems`.

## Other examples

Goroutine leaks come in various forms, so in the following section,
we present a few common examples of coding patterns that lead to leaks
observed in industrial-scale codebases and open source projects,
in ascending order of complexity.

You can quickly test drive the goroutine leak detector on them in
[the Go playground](/play/p/3C71z4Dpav-?v=gotip), and even
experiment with your own leaks.

### Double send

Some of the simplest leaks occur when more messages
are sent over a channel than expected.
Below, a goroutine sends messages to the main goroutine
over an unbuffered channel.
However, in case of an error, the return statement is missing,
causing two messages to be sent, and leading to the leak.
```go
func DoubleSend() {
	ch := make(chan any)
	go func(err error) {
		if err != nil {
			// In case of an error, send nil.
			ch <- nil
			// Return statement is missing.
		}
		// Otherwise, continue with normal behaviour.
		// This send is still executed, which causes a leak in the error case.
		ch <- struct{}{}
	}(fmt.Errorf("error"))
	// Retrieve only one message.
	<-ch
}
```
While the profile will not explicitly highlight the cause,
it directs you to the send operation.
```
(pprof) list DoubleSend
Total: 1
ROUTINE ======================== main.DoubleSend.func1 in .../main.go
         0          1 (flat, cum)   100% of Total
         .          .    118:   go func(err error) {
         .          .    119:           if err != nil {
         .          .    121:                   ch <- nil
         .          .    123:           }
         .          1    126:           ch <- struct{}{}
         .          .    127:   }(fmt.Errorf("error"))
         .          .    129:   <-ch
```
This leak can be addressed simply by adding a `return` statement after the
send operation in the error case.

### Early return

The inverse situation is just as common, where the receiver
omits communication on some control flow paths.
This is effectively a simplification of the introductory example.
```go
// Incoming error simulates an error produced internally.
func EarlyReturn(err error) {
	ch := make(chan any)

	// Create a worker goroutine.
	go func() {
		// Send something to the channel.
		// Leaks if the parent goroutine terminates early.
		ch <- struct{}{}
	}()

	if err != nil {
		// Interrupt evaluation of parent early in case of error.
		// Sender leaks.
		return
	}

	// Only receive if there is no error.
	<-ch
}
```
The goroutine leak is exposed by the profile:
```
ROUTINE ======================== main.EarlyReturn.func1 in .../main.go
         0          1 (flat, cum)   100% of Total
         .          .    140:   go func() {
         .          1    143:           ch <- struct{}{}
         .          .    144:   }()
         .          .    145:
         .          .    146:   if err != nil {
```
The leak can be addressed by increasing the buffer size of `ch` to 1.

### Timeout

A variation of the **Early return** pattern above involves contexts
and non-deterministic choice (`select` statements):
```go
func Timeout(ctx context.Context) {
	// An unbuffered channel is used to coordinate
	// a worker and parent thread
	ch := make(chan any)

	// Create worker goroutine
	go func() {
		// Perform some work then signal to the parent thread.
		ch <- struct{}{}
	}()

	// Wait for message from worker or context
	// to be cancelled or timed out.
	select {
	case <-ch: // Receive message from worker
	case <-ctx.Done():
		// Sender leaks because there is no
		// future rendez-vous over the channel.
	}
}
```
If the context is cancelled before the sender synchronizes with the parent,
the sender will leak:
```
(pprof) list Timeout
Total: 10
ROUTINE ======================== main.Timeout.func1.1 in .../main.go
         0         10 (flat, cum)   100% of Total
         .          .    198:           go func() {
         .         10    201:                   ch <- struct{}{}
         .          .    202:           }()
```
As in the previous example, the fix is to increase the channel
buffer: `ch = make(chan any, 1)`.

### Range over channel without closing

One slightly esoteric concurrency feature is
[iterating over channels](/tour/concurrency/4) by using `range`.
This allows you to repeatedly receive values from a channel in a loop,
automatically extracting each sent value until the channel is closed and
all buffered values have been received, upon which the loop exits.

Importantly, **if the channel is never closed**, a `range` loop will block
the executing goroutine forever.
Omitting the `close` operation is a common mistake, as below:
```go
// Incoming list of items and the number of workers.
func noCloseRange(list []any, workers int) {
	// Create a channel that distributes work items.
	ch := make(chan any)

	// Create the worker goroutines.
	for i := 0; i < workers; i++ {
		go func() {
			// Each worker pulls items from the channel
			// and then processes it.
			for item := range ch {
				// Process each item
				_ = item
			}
		}()
	}

	// Queue items to the workers by using the channel.
	for _, item := range list {
		// The parent leaks by sending an item if workers == 0
		// or if all the workers panic, but the panic is recovered.
		ch <- item
	}
	// Otherwise, the channel is never closed, so workers
	// leak once there are no more items left to process.
}

...
// Example calls
go noCloseRange([]any{1, 2, 3}, 3) // Leaks all 3 workers
go noCloseRange([]any{1, 2, 3}, 0) // Leaks caused by 0 workers
```
A goroutine leak profile for such a program would include the following:
```
Type: goroutineleak
(pprof) list noCloseRange.func1
Total: 4
ROUTINE ======================== main.noCloseRange.func1 in .../main.go
         0          3 (flat, cum) 75.00% of Total
         .          .     82:           go func() {
         .          3     84:                   for item := range ch {
         .          .     86:                           _ = item
         .          .     87:                   }
         .          .     88:           }()
```
We see the 3 workers blocked at the `range ch` operation, which
gives an ample hint as to the cause of the leak.

There is a bonus leak scenario in this case,
if the number of workers is mistakenly set to zero,
wherein the parent sender will leak.
This is also captured by the profile:
```
(pprof) list noCloseRange$
Total: 4
ROUTINE ======================== main.noCloseRange in .../main.go
         0          1 (flat, cum) 25.00% of Total
         .          .     76:func noCloseRange(list []any, workers int) {
...
         .          .     92:   for _, item := range list {
         .          1     95:           ch <- item
         .          .     96:   }
```

The `range` leak can be addressed by simply closing the channel once
all messages have been sent:
```go
	for _, item := range list {
		ch <- item
	}
	// All items have been sent. It is now safe to close.
	close(ch)
```
While `workers > 0` can be assumed to be an invariant in a realistic production system,
goroutine leak profiles can nevertheless be used to implicitly monitor for off-chance
violations without conservative `workers <= 0` checks.

### Method contract violations

The patterns seen so far have been relatively constrained in their lexical scope.
However, as functionality is spread out across functions, methods and packages, and
implementations are obfuscated by interfaces, the difficulty of manually diagnosing
leaks drastically increases.

Such a case is exemplified in this section, with the custom `worker` type that embeds two channel
fields, `ch` and `done` and creates a looping goroutine with its `Start` method that
reads from both channels with a `select` statement.
Said goroutine can only be terminated by receiving a message through the `done` channel,
which is closed by the `Stop` method.

The `Start` method can be invoked any number of times, but if it is invoked
at least once, `Stop` should eventually be called.

As a result, `Start` and `Stop` form an implicit contract that dictates the order
in which, and number of times each method should be invoked.
Breaking that contract can lead to undesirable behavior,
in this case, goroutine leaks:
```go
func MethodContractViolation() {
	items := make([]any, 10)
	// Create a new worker
	w := NewWorker()

	// Start worker
	w.Start()

	// Operate on worker
	for _, item := range items {
		w.AddToQueue(item)
	}
	// Exits without calling ’Stop’.
}

type worker struct {
	once *sync.Once

	ch   chan any
	done chan any
}

type Worker interface {
	Start()
	Stop()
	AddToQueue(item any)
}

func NewWorker() Worker {
	return &worker{
		once: &sync.Once{},

		ch:   make(chan any),
		done: make(chan any),
	}
}

// Start spawns a background goroutine that extracts items pushed to the queue.
func (w *worker) Start() {
	go func() {
		for {
			select {
			case <-w.ch: // Normal workflow
			case <-w.done:
				return // Shut down
			}
		}
	}()
}

func (w *worker) Stop() {
	// Allows goroutine created by Start to terminate
	w.once.Do(func() {
		close(w.done)
	})
}

func (w *worker) AddToQueue(item any) {
	w.ch <- item
}
```
This issue is further exacerbated in practice, where such custom types are exposed
as APIs through interfaces, in this case, through the non-descript
`Worker` type.
Clients may not even be aware of the underlying implementation and,
consequently, violate the implicit contract without realizing.

Fortunately, soliciting a goroutine leak profile can reveal the defect:
```
(pprof) list Start
Total: 1
ROUTINE ======================== main.(*worker).Start.func1 in .../main.go
         0          1 (flat, cum)   100% of Total
         .          .    266:   go func() {
         .          .    267:           for {
         .          1    268:                   select {
         .          .    269:                   case <-w.ch:
         .          .    270:                   case <-w.done:
         .          .    271:                           return
```
Naturally, the fix involves following the trail to the `Start` call
and adding an invocation of `Stop`.

### Cockroach/584

The following real-world example is taken from the open-source
project [cockroachdb](https://github.com/cockroachdb/cockroach/pull/584/files).
It involves acquiring and releasing a lock in a loop, but omitting the release
upon executing a `break` statement:
```go
type Gossip struct {
	mu     sync.Mutex // L1
	closed bool
}

func Cockroach584() {
	g := &Gossip{
		closed: true,
	}
	// ...
	g.bootstrap()
	g.manage()
}

func (g *Gossip) bootstrap() {
	for {
		g.mu.Lock()
		if g.closed {
			// Missing g.mu.Unlock
			break
		}
		g.mu.Unlock()
	}
}

func (g *Gossip) manage() {
	for {
		g.mu.Lock()
		if g.closed {
			// Missing g.mu.Unlock
			break
		}
		g.mu.Unlock()
	}
}
```
It is easy to see how one of the two goroutines in this scenario
will eventually leak by failing to acquire the lock.
```
(pprof) list Gossip
Total: 1
ROUTINE ======================== main.(*Gossip).bootstrap in .../main.go
         0          1 (flat, cum)   100% of Total
         .          .    165:func (g *Gossip) bootstrap() {
         .          .    166:   for {
         .          1    167:           g.mu.Lock()
         .          .    168:           if g.closed {
         .          .    170:                   break
         .          .    171:           }
         .          .    172:           g.mu.Unlock()
```


## Implementation {#implementation}

This section is intended for those interested in the underlying machinations
of the leak detector.
If you are, instead, interested in capabilities and performance,
skip ahead to [limitations](/blog/goroutine-leak-profiles#limitations).

### Core concept

Let's start with an initial observation: if a goroutine
is blocked over some concurrency primitive that no other goroutine has access to,
then it is obviously leaked.

This observation, while simplified, already gives us a strong lead on how to
reliably detect goroutine leaks at runtime.
Our goal now is to achieve it in practice, and generalize and expand upon
this observation into a definition that encompasses more leak
scenarios.
We, therefore, define for goroutines the property of _maybe-liveness_,
i.e., whether a goroutine could eventually be unblocked
(not to be confused with liveness in memory, as determined by the GC).

Maybe-liveness is an inductive property, and is defined thusly:
> A goroutine is maybe-live if:
> 1. it is not in the blocked state, or
> 2. at least one concurrency primitive that blocks it is accessible through a
	maybe-live goroutine.

In the first case, goroutines which are not blocked are obviously
not leaked.
In the second case, the assumption is that the primitive may still be operated
upon in the future by the other goroutine in order to unblock our goroutine.

As a corollary, any goroutine which is not maybe-live is definitely leaked.

Our goal is, therefore, to determine which goroutines in the system
satisfy maybe-liveness.
The core strategy is to initially assume maybe-liveness only
for non-blocked goroutines.
Blocked goroutines are then incrementally determined to satisfy
the property, depending on
what concurrency primitives existing maybe-live goroutines have
access to.
Since access to a concurrency primitive coincides with holding a reference
to it, the problem reduces to memory reachability.

Fortunately for us, a mechanism for computing memory reachability already
exists in the Go runtime:
the [garbage collector](/doc/gc-guide) (GC).
The Go runtime uses a concurrent tri-color mark-and-sweep garbage collector,
now with the [Green Tea](/blog/greenteagc) variant!

With the task laid out before us, we endeavored to adapt the GC to suit our purposes.
You can quickly compare the two GCs with the following diagram:

<div class="centered">
<button type="button" id="greentea-prev" class="scroll-button scroll-button-left" hidden disabled>← Prev</button>
<button type="button" id="greentea-next" class="scroll-button scroll-button-right" hidden>Next →</button>
<div id="goroutineleakgc" class="carousel">
	<figure class="carouselitem">
		<img src="goroutine-leak-profiles/gc-original.svg" />
	</figure>
	<figure class="carouselitem">
		<img src="goroutine-leak-profiles/gc-modified.svg" />
	</figure>
</div>
</div>

As it happens, a complete overhaul of the GC is not necessary, as
its MO already neatly aligns with our goals.
Only a few key modifications were needed:
1. In the initial phases, the regular GC uses global data and all goroutines
	as mark roots. That wouldn't suit our purposes, so we change it to initially
	only include non-blocked goroutines, aligning with the base case of the
	inductive definition of maybe-liveness.
2. Marking carries on as in the usual GC, with the added benefit that now only
	memory that is reachable through maybe-live goroutines is marked. The algorithm
	assumes that any concurrency primitives marked in this step might still be used
	in the future.
3. Finalizing the marking phase now involves an additional task: checking
	whether any blocked goroutines are now "maybe-live" owing to one of their
	blocking concurrency primitives having been marked. If such a goroutine is found,
	it is added as a mark root, and the GC resumes the marking phase to mark
	its reachable memory. This coincides with the inductive step in the definition
	of maybe-liveness.
4. Once all maybe-live goroutines have been discovered, the only remaining
	goroutines are obviously leaked, so they can be reported. In our case, their
	status is internally set to leaked, such that they may be included in the
	goroutine leak profile.
5. The marking phase should now resume once again with all the leaked goroutines
	added	as mark roots, such that the remaining memory can be marked, to re-align
	GC behavior with the regular runtime.

For the sake of correctness, goroutine leak detection is strictly limited
to Go first-class concurrency primitives, which includes:
-	channel-based concurrency, such as send and receive operations, including
	over `nil` channels, as well as blocking `select` statements (i.e., without a
	`default` case), including with no cases, and
- specific members of the [`sync`](/pkg/sync) package, including `Mutex`,
	`RWMutex`, `WaitGroup` and `Cond`.

### Limitations {#limitations}

The examples above demonstrate the usefulness of goroutine leak profiles.
Nevertheless, the reliance of the detection mechanism on the garbage
collector does impose some limitations that may lead
it to miss leaks:

1. **Memory overreach**: if a concurrency primitive is
consistently reachable through **global variables** or **runnable goroutines**,
then goroutines blocking on it are never reported as leaked, even if
that is the case in practice.
	* This can be alleviated by better delineating the lifecycle of
		concurrency resources, and more strictly regimenting which
		parties may acquire their references.

2. **Non-standard blocking**: goroutines blocked for any reason that
	does not involve first-class concurrency primitives, e.g.,
	netpollers or semaphores internal to the runtime, are never considered as leaking.
	This likewise applies for custom, user-defined concurrency operations
	(e.g., spin locks), unless they rely on the primitives outlined above.

3. **Non-determinism**: leaks can only be detected after
they have occurred, not predicted.
		 Reproducing and diagnosing leaks in flaky programs
		continues to be a challenge.
		For the best results, we encourage mixing approaches, with
		goroutine leak profiles at various levels, as well as comprehensive
		test suites instrumented with [goleak](https://github.com/uber-go/goleak).

### Performance impact {#performance}

While the goroutine leak detection mechanism is carefully designed to minimize
performance impact, some costs are still incurred.

#### Memory

Leak detection attempts to minimize memory overhead by only adding
constant-sized components for the sake of book-keeping.

The only introduced scaling factor is: _maybe-traceable pointers_.
To ensure that the GC behaves according to our specifications,
we must prevent some references from being traced prematurely.

This is where maybe-traceable pointers come in, which are objects
that carry two references: a) an untraceable pointer-as-an-integer value,
`vu`, and b) the same pointer as an actual reference that is understood
as a pointer by the GC, `vp`.
A maybe-traceable pointer is, therefore, double the size of
its regular counterpart.

Untraceable pointers come in 3 valid states:
1. `vu` and `vp` are unset, which is analogous to a `nil` pointer,
2. `vu` and `vp` are set (and equal), which is analogous to a regular reference
that can be traced by the GC,
3. `vu` is set, but `vp` is unset (`nil`), which preserves the reference,
	but "hides" it from the GC.

Maybe-traceable pointers are relevant for `sudog`s,
objects which pair individual goroutines and concurrency primitive.
One concurrency primitive can block multiple goroutines,
and one goroutine can be blocked on multiple concurrency primitives (because of `select`
statements).
Therefore, the maximum number of active `sudog` objects at any given point is
the product of the number of goroutines and concurrency primitives.

Each `sudog` holds, among other things, references to its
blocking concurrency primitive.
However, `sudog`s are also globally reachable through the `sudog`
cache, which exposes these references to the GC during the marking phase.
This goes against our goal of only tracing these references when
they are reachable from a maybe-live goroutine.
Therefore, in order to prevent the GC from tracing them, we update
these references in `sudog` to be maybe-traceable pointers.
Maybe-traceable pointers are set as untraceable at the start of
goroutine leak detection, and only updated to traceable once
the goroutine paired to the same `sudog` is scheduled for marking.

While the asymptotic complexity remains unchanged, a modest cost
is nevertheless incurred by doubling the size of these references.
The worst case scenario occurs only when every goroutine is blocked on every
concurrency primitive in the system, which is unrealistic for most Go
programs.

#### Computational overhead

In its current implementation, goroutine leak detection is more
computationally expensive than the regular GC.
This is best illustrated by looking at a pathological case we
will call the "daisy-chain":
<img src="goroutine-leak-profiles/daisy-chain.svg" />
In this example without leaks, runnable goroutine G₀ has a
reference to primitive P₁ which blocks G₁, and so on in a daisy chain
pattern.

This implies that to prove maybe-liveness for some Pᵢ₊₁,
we must first prove maybe-liveness for Pᵢ, which introduces
two costs:
1. The marking phase is effectively serialized relative to the
	order in which goroutines can be scanned, as all the memory reachable
	from some Pᵢ must be marked before Pᵢ₊₁ can be added as a root.
2. The inspection currently traverses the entire tail of blocked goroutines
	at the end of each marking round, which takes n² steps.

The second point can be addressed over time, but the first point of contention
is an intrinsic limitation that cannot be circumvented.

Regardless, if a goroutine leak can be observed at some point in time, then it
can also be observed at any future point in the same Go program's execution.

Periodic profiling infrastructures can therefore tune profiling frequency,
e.g., every 4 hours, for virtually no loss in leak detection capabilities,
while simultaneously only sporadically incurring the execution
overhead.


## Next steps

The goroutine leak profile is available as an experiment in Go 1.26, enabled with `GOEXPERIMENT=goroutineleakprofile`.
We encourage developers to try it in testing, continuous integration, and production environments.

The implementation is production-ready; the experimental status is solely to gather feedback on the API design.
We plan to enable goroutine leak profiles by default in Go 1.27, making automatic leak detection available to all Go programs without any build flags.

Please share your experiences and feedback on the [proposal issue](/issue/74609)!

## Acknowledgements

The goroutine leak detection is the result of a research collaboration between
Aarhus University, Washington University in St. Louis, and Uber, as presented in
["Dynamic Partial Deadlock Detection and Recovery via Garbage Collection"](https://dl.acm.org/doi/pdf/10.1145/3676641.3715990)
(Saioc et al., ASPLOS 2025).

The transition from academic prototype to actual Go feature was made possible
with the guidance of Michael Knyszek and Michael Pratt, in the Go team at Google, and
[@thepudds](https://github.com/thepudds).
This kind of cross-company collaboration continues to make Go better for everyone.

<script src="greenteagc/carousel.js"></script>