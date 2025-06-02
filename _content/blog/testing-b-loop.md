---
title: "More predictable benchmarking with testing.B.Loop"
date: 2025-04-02
by:
- Junyang Shao
tags:
- benchmark
- testing
- compile
summary: Better benchmark looping in Go 1.24.
---

Go developers who have written benchmarks using the
[`testing`](https://pkg.go.dev/testing) package might have encountered some of
its various pitfalls. Go 1.24 introduces a new way to write benchmarks that's just
as easy to use, but at the same time far more robust:
[`testing.B.Loop`](https://pkg.go.dev/testing#B.Loop).

Traditionally, Go benchmarks are written using a loop from 0 to `b.N`:
```
func Benchmark(b *testing.B) {
  for range b.N {
    ... code to measure ...
  }
}
```
Using `b.Loop` instead is a trivial change:
```
func Benchmark(b *testing.B) {
  for b.Loop() {
    ... code to measure ...
  }
}
```

`testing.B.Loop` has many benefits:
* It prevents unwanted compiler optimizations within the benchmark loop.
* It automatically excludes setup and cleanup code from benchmark timing.
* Code can't accidentally depend on the total number of iterations or the current
iteration.

These were all easy mistakes to make with `b.N`-style benchmarks that would
silently result in bogus benchmark results. As an added bonus, `b.Loop`-style
benchmarks even complete in less time!

Let's explore the advantages of `testing.B.Loop` and how to effectively utilize it.

## Old benchmark loop problems

Before Go 1.24, while the basic structure of a benchmark was simple, more sophisticated
benchmarks required more care:
```
func Benchmark(b *testing.B) {
  ... setup ...
  b.ResetTimer() // if setup may be expensive
  for range b.N {
    ... code to measure ...
    ... use sinks or accumulation to prevent dead-code elimination ...
  }
  b.StopTimer() // if cleanup or reporting may be expensive
  ... cleanup ...
  ... report ...
}
```
If setup or cleanup are non-trivial, the developer needs to surround the benchmark loop
with `ResetTimer` and/or `StopTimer` calls. These are easy to forget, and even if the
developer remembers they may be necessary, it can be difficult to judge whether setup or
cleanup are "expensive enough" to require them.

Without these, the `testing` package can only time the entire benchmark function. If a
benchmark function omits them, the setup and cleanup code will be included in the overall
time measurement, silently skewing the final benchmark result.


There is another, more subtle pitfall that requires deeper understanding:
([Example source](https://eli.thegreenplace.net/2023/common-pitfalls-in-go-benchmarking/))

```
func isCond(b byte) bool {
  if b%3 == 1 && b%7 == 2 && b%17 == 11 && b%31 == 9 {
    return true
  }
  return false
}

func BenchmarkIsCondWrong(b *testing.B) {
  for range b.N {
    isCond(201)
  }
}
```
In this example, the user might observe `isCond` executing in sub-nanosecond
time. CPUs are fast, but not that fast! This seemingly anomalous result stems
from the fact that `isCond` is inlined, and since its result is never used, the
compiler eliminates it as dead code. As a result, this benchmark doesn't measure `isCond`
at all; it measures how long it takes to do nothing. In this case, the sub-nanosecond
result is a clear red flag, but in more complex benchmarks, partial dead-code elimination
can lead to results that look reasonable but still aren't measuring what was intended.

## How `testing.B.Loop` helps

Unlike a `b.N`-style benchmark, `testing.B.Loop` is able to track when it is first called
in a benchmark when the final iteration ends. The `b.ResetTimer` at the loop's start
and `b.StopTimer` at its end are integrated into `testing.B.Loop`, eliminating the need
to manually manage the benchmark timer for setup and cleanup code.

Furthermore, the Go compiler now detects loops where the condition is just a call to
`testing.B.Loop` and prevents dead code elimination within the loop. In Go 1.24, this is
implemented by disallowing inlining into the body of such a loop, but we plan to
[improve](/issue/73137) this in the future.

Another nice feature of `testing.B.Loop` is its one-shot ramp-up approach. With a `b.N`-style
benchmark, the testing package must call the benchmark function several times with different
values of `b.N`, ramping up until the measured time reached a threshold. In contrast, `b.Loop`
can simply run the benchmark loop until it reaches the time threshold, and only needs to call
the benchmark function once. Internally, `b.Loop` still uses a ramp-up process to amortize
measurement overhead, but this is hidden from the caller and can be more efficient.

Certain constraints of the `b.N`-style loop still apply to the `b.Loop`-style
loop. It remains the user's responsibility to manage the timer within the benchmark loop,
when necessary:
([Example source](https://eli.thegreenplace.net/2023/common-pitfalls-in-go-benchmarking/))

```
func BenchmarkSortInts(b *testing.B) {
  ints := make([]int, N)
  for b.Loop() {
    b.StopTimer()
    fillRandomInts(ints)
    b.StartTimer()
    slices.Sort(ints)
  }
}
```
In this example, to benchmark the in-place sorting performance of `slices.Sort`, a
randomly initialized array is required for each iteration. The user must still
manually manage the timer in such cases.

Furthermore, there still needs to be exactly one such loop in the benchmark function body
(a `b.N`-style loop cannot coexist with a `b.Loop`-style loop), and every iteration of the
loop should do the same thing.

## When to use

The `testing.B.Loop` method is now the preferred way to write benchmarks:
```
func Benchmark(b *testing.B) {
  ... setup ...
  for b.Loop() {
    // optional timer control for in-loop setup/cleanup
    ... code to measure ...
  }
  ... cleanup ...
}
```

`testing.B.Loop` offers faster, more accurate, and
more intuitive benchmarking.

## Acknowledgements

A huge thank you to everyone in the community who provided feedback on the proposal
issue and reported bugs as this feature was released! I'm also grateful to Eli
Bendersky for his helpful blog summaries. And finally a big thank you to Austin Clements,
Cherry Mui and Michael Pratt for their review, thoughtful work on the design options and
documentation improvements. Thank you all for your contributions!