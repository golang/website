---
title: Go 1.21 Release Candidate
date: 2023-06-21
by:
- Eli Bendersky, on behalf of the Go team
summary: Go 1.21 RC brings language improvements, new standard library packages, PGO GA, backward and forward compatibility in the toolchain and faster builds.
---

The Go 1.21 first Release Candidate (RC) is available today on the [download
page](/dl/#go1.21rc2)! Go 1.21 is packed with new features and improvements.
Getting the RC (release candidate) allows you to experiment with it early, try
it on your workloads, and report any issues before the final release (scheduled
for August). Here are some notable changes and features in Go 1.21; for the full
list, refer to the [full release notes](https://tip.golang.org/doc/go1.21).

*(Please note that the first RC for Go 1.21 is called `go1.21rc2`
because a bug was found and fixed after tagging `go1.21rc1`)*

## Tool improvements

- The Profile Guided Optimization (PGO) feature we [announced for preview in
  1.20](/blog/pgo-preview) is now generally available! If a file named
  `default.pgo` is present in the main package’s directory, the `go` command
  will use it to enable a PGO build. See the [PGO documentation](/doc/pgo) for
  more details. We’ve measured the impact of PGO on a wide set of Go programs and
  see performance improvements of 2-7%.
- The [`go` tool](/cmd/go) now supports [backward](https://tip.golang.org/doc/godebug)
  and [forward](/doc/toolchain) language compatibility.

## Language changes

- New built-in functions: [min, max](https://tip.golang.org/ref/spec#Min_and_max)
  and [clear](https://tip.golang.org/ref/spec#Clear).
- Several improvements to type inference for generic functions. The description of
  [type inference in the spec](https://tip.golang.org/ref/spec#Type_inference)
  has been expanded and clarified.
- In a future version of Go we’re planning to address one of the most common
  gotchas of Go programming:
  [loop variable capture](https://go.dev/wiki/CommonMistakes).
  Go 1.21 comes with a preview of this feature that you can enable in your code
  using an environment variable. See [this LoopvarExperiment wiki
  page](https://go.dev/wiki/LoopvarExperiment) for more details.

## Standard library additions

- New [log/slog](https://tip.golang.org/pkg/log/slog) package for structured logging.
- New [slices](https://tip.golang.org/pkg/slices) package for common operations
  on slices of any element type. This includes sorting functions that are generally
  faster and more ergonomic than the [sort](https://tip.golang.org/pkg/sort) package.
- New [maps](https://tip.golang.org/pkg/maps) package for common operations on maps
  of any key or element type.
- New [cmp](https://tip.golang.org/pkg/cmp) package with new utilities for comparing
  ordered values.

## Improved performance

In addition to the performance improvements when enabling PGO:

- The Go compiler itself has been rebuilt with PGO enabled for 1.21, and as a
  result it builds Go programs 2-4% faster, depending on the host architecture.
- Due to tuning of the garbage collector, some applications may see up to a 40%
  reduction in tail latency.
- Collecting traces with [runtime/trace](https://pkg.go.dev/runtime/trace) now
  incurs a substantially smaller CPU cost on amd64 and arm64.

## A new port to WASI

Go 1.21 adds an experimental port for [WebAssembly System Interface (WASI)](https://wasi.dev/),
Preview 1 (`GOOS=wasip1`, `GOARCH=wasm`).

To facilitate writing more general WebAssembly (WASM) code, the compiler also
supports a new directive for importing functions from the WASM host:
`go:wasmimport`.

Please [download the Go 1.21 RC](/dl/#go1.21rc2) and try it! If you notice any
problems, please [file an issue](/issue/new).
