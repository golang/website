---
title: "Go runtime: 4 years later"
date: 2022-09-26
by:
- Michael Knyszek
summary: A check-in on the status of Go runtime development
---

Since our [last blog post about the Go GC in 2018](/blog/ismmkeynote) the
Go GC, and the Go runtime more broadly, has been steadily improving.
We've tackled some large projects, motivated by real-world Go programs and real
challenges facing Go users.
Let's catch you up on the highlights!

### What's new?

- `sync.Pool`, a GC-aware tool for reusing memory, has a [lower latency
  impact](https://go.dev/cl/166960) and [recycles memory much more
  effectively](https://go.dev/cl/166961) than before.
  (Go 1.13)

- The Go runtime returns unneeded memory back to the operating system [much
  more proactively](https://go.dev/issue/30333), reducing excess memory
  consumption and the chance of out-of-memory errors.
  This reduces idle memory consumption by up to 20%.
  (Go 1.13 and 1.14)

- The Go runtime is able to preempt goroutines more readily in many cases,
  reducing stop-the-world latencies up to 90%.
  [Watch the talk from Gophercon
  2020 here.](https://www.youtube.com/watch?v=1I1WmeSjRSw)
  (Go 1.14)

- The Go runtime [manages timers more efficiently than
  before](https://go.dev/cl/171883), especially on machines with many CPU cores.
  (Go 1.14)

- Function calls that have been deferred with the `defer` statement now cost as
  little as a regular function call in most cases.
  [Watch the talk from Gophercon 2020
  here.](https://www.youtube.com/watch?v=DHVeUsrKcbM)
  (Go 1.14)

- The memory allocator's slow path [scales](https://go.dev/issue/35112)
  [better](https://go.dev/issue/37487) with CPU cores, increasing throughput up
  to 10% and decreasing tail latencies up to 30%, especially in highly-parallel
  programs.
  (Go 1.14 and 1.15)

- Go memory statistics are now accessible in a more granular, flexible, and
  efficient API, the [runtime/metrics](https://pkg.go.dev/runtime/metrics)
  package.
  This reduces latency of obtaining runtime statistics by two orders of
  magnitude (milliseconds to microseconds).
  (Go 1.16)

- The Go scheduler spends up to [30% less CPU time spinning to find new
  work](https://go.dev/issue/43997).
  (Go 1.17)

- Go code now follows a [register-based calling
  convention](https://go.dev/issues/40724) on amd64, arm64, and ppc64, improving
  CPU efficiency by up to 15%.
  (Go 1.17 and Go 1.18)

- The Go GC's internal accounting and scheduling has been
  [redesigned](https://go.dev/issue/44167), resolving a variety of long-standing
  issues related to efficiency and robustness.
  This results in a significant decrease in application tail latency (up to 66%)
  for applications where goroutines stacks are a substantial portion of memory
  use.
  (Go 1.18)

- The Go GC now limits [its own CPU use when the application is
  idle](https://go.dev/issue/44163).
  This results in 75% lower CPU utilization during a GC cycle in very idle
  applications, reducing CPU spikes that can confuse job shapers.
  (Go 1.19)

These changes have been mostly invisible to users: the Go code they've come to
know and love runs better, just by upgrading Go.

### A new knob

With Go 1.19 comes an long-requested feature that requires a little extra work
to use, but carries a lot of potential: [the Go runtime's soft memory
limit](https://pkg.go.dev/runtime/debug#SetMemoryLimit).

For years, the Go GC has had only one tuning parameter: `GOGC`.
`GOGC` lets the user adjust [the trade-off between CPU overhead and memory
overhead made by the Go GC](https://pkg.go.dev/runtime/debug#SetGCPercent).
For years, this "knob" has served the Go community well, capturing a wide
variety of use-cases.

The Go runtime team has been reluctant to add new knobs to the Go runtime,
with good reason: every new knob represents a new _dimension_ in the space of
configurations that we need to test and maintain, potentially forever.
The proliferation of knobs also places a burden on Go developers to understand
and use them effectively, which becomes more difficult with more knobs.
Hence, the Go runtime has always leaned into behaving reasonably with minimal
configuration.

So why add a memory limit knob?

Memory is not as fungible as CPU time.
With CPU time, there's always more of it in the future, if you just wait a bit.
But with memory, there's a limit to what you have.

The memory limit solves two problems.

The first is that when the peak memory use of an application is unpredictable,
`GOGC` alone offers virtually no protection from running out of memory.
With just `GOGC`, the Go runtime is simply unaware of how much memory it has
available to it.
Setting a memory limit enables the runtime to be robust against transient,
recoverable load spikes by making it aware of when it needs to work harder to
reduce memory overhead.

The second is that to avoid out-of-memory errors without using the memory limit,
`GOGC` must be tuned according to peak memory, resulting in higher GC CPU
overheads to maintain low memory overheads, even when the application is not at
peak memory use and there is plenty of memory available.
This is especially relevant in our containerized world, where programs are
placed in boxes with specific and isolated memory reservations; we might as
well make use of them!
By offering protection from load spikes, setting a memory limit allows for
`GOGC` to be tuned much more aggressively with respect to CPU overheads.

The memory limit is designed to be easy to adopt and robust.
For example, it's a limit on the whole memory footprint of the Go parts of an
application, not just the Go heap, so users don't have to worry about accounting
for Go runtime overheads.
The runtime also adjusts its memory scavenging policy in response to the memory
limit so it returns memory to the OS more proactively in response to memory
pressure.

But while the memory limit is a powerful tool, it must still be used with some
care.
One big caveat is that it opens up your program to GC thrashing: a state in
which a program spends too much time running the GC, resulting in not enough
time spent making meaningful progress.
For example, a Go program might thrash if the memory limit is set too low for
how much memory the program actually needs.
GC thrashing is something that was unlikely previously, unless `GOGC` was
explicitly tuned heavily in favor of memory use.
We chose to favor running out of memory over thrashing, so as a mitigation, the
runtime will limit the GC to 50% of total CPU time, even if this means exceeding
the memory limit.

All of this is a lot to consider, so as a part of this work, we released [a
shiny new GC guide](/doc/gc-guide), complete with interactive visualizations to
help you understand GC costs and how to manipulate them.

### Conclusion

Try out the memory limit!
Use it in production!
Read the [GC guide](/doc/gc-guide)!

We're always looking for feedback on how to improve Go, but it also helps to
hear about when it just works for you.
[Send us feedback](https://groups.google.com/g/golang-dev)!
