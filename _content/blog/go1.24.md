---
title: Go 1.24 is released!
date: 2025-02-11
by:
- Junyang Shao, on behalf of the Go team
summary:
  Go 1.24 brings generic type aliases, map performance improvements, FIPS 140
  compliance and more.
---

Today the Go team is excited to release Go 1.24,
which you can get by visiting the [download page](/dl/).

Go 1.24 comes with many improvements over Go 1.23. Here are some of the notable
changes; for the full list, refer to the [release notes](/doc/go1.24).

## Language changes

<!-- go.dev/issue/46477 -->
Go 1.24 now fully supports [generic type aliases](/issue/46477): a type alias
may be parameterized like a defined type.
See the [language spec](/ref/spec#Alias_declarations) for details.

## Performance improvements

<!-- go.dev/issue/54766, go.dev/cl/614795, go.dev/issue/68578 -->
Several performance improvements in the runtime have decreased CPU overhead
by 2–3% on average across a suite of representative benchmarks. These
improvements include a new builtin `map` implementation based on
[Swiss Tables](https://abseil.io/about/design/swisstables), more efficient
memory allocation of small objects, and a new runtime-internal mutex
implementation.

## Tool improvements

- <!-- go.dev/issue/48429 -->
  The `go` command now provides a mechanism for tracking tool dependencies for a
  module. Use `go get -tool` to add a `tool` directive to the current module. Use
  `go tool [tool name]` to run the tools declared with the `tool` directive.
  Read more on the [go command](/doc/go1.24#go-command) in the release notes.
- <!-- go.dev/issue/44251 -->
  The new `test` analyzer in `go vet` subcommand reports common mistakes in
  declarations of tests, fuzzers, benchmarks, and examples in test packages.
  Read more on [vet](/doc/go1.24#vet) in the release notes.

## Standard library additions

- The standard library now includes [a new set of mechanisms to facilitate
  FIPS 140-3 compliance](/doc/security/fips140). Applications require no source code
  changes to use the new mechanisms for approved algorithms. Read more
  on [FIPS 140-3 compliance](/doc/go1.24#fips140) in the release notes.
  Apart from FIPS 140, several packages that were previously in the
  [x/crypto](/pkg/golang.org/x/crypto) module are now available in the
  [standard library](/doc/go1.24#crypto-mlkem).

- Benchmarks may now use the faster and less error-prone
  [`testing.B.Loop`](/pkg/testing#B.Loop) method to perform benchmark iterations
  like `for b.Loop() { ... }` in place of the typical loop structures involving
  `b.N` like `for range b.N`. Read more on
  [the new benchmark function](/doc/go1.24#new-benchmark-function) in the
  release notes.

- The new [`os.Root`](/pkg/os#Root) type provides the ability to perform
  filesystem operations isolated under a specific directory. Read more on
  [filesystem access](/doc/go1.24#directory-limited-filesystem-access) in the
  release notes.

- The runtime provides a new finalization mechanism,
  [`runtime.AddCleanup`](/pkg/runtime#AddCleanup), that is more flexible,
  more efficient, and less error-prone than
  [`runtime.SetFinalizer`](/pkg/runtime#SetFinalizer). Read more on
  [cleanups](/doc/go1.24#improved-finalizers) in the release notes.

## Improved WebAssembly support

<!-- go.dev/issue/65199, CL 603055 -->
Go 1.24 adds a new `go:wasmexport` directive for Go programs to export
functions to the WebAssembly host, and supports building a Go program as a WASI
[reactor/library](https://github.com/WebAssembly/WASI/blob/63a46f61052a21bfab75a76558485cf097c0dbba/legacy/application-abi.md#current-unstable-abi).
Read more on [WebAssembly](/doc/go1.24#wasm) in the release notes.

---


Please read the [Go 1.24 release notes](/doc/go1.24) for the complete and
detailed information. Don't forget to watch for follow-up blog posts that
will go in more depth on some of the topics mentioned here!

Thank you to everyone who contributed to this release by writing code and
documentation, reporting bugs, sharing feedback, and testing the release
candidates. Your efforts helped to ensure that Go 1.24 is as stable as possible.
As always, if you notice any problems, please [file an issue](/issue/new).

Enjoy Go 1.24!
