---
title: Go 1.23 is released
date: 2024-08-13
by:
- Dmitri Shuralyov, on behalf of the Go team
summary: Go 1.23 adds iterators, continues loop enhancements, improves compatibility, and more.
---

Today the Go team is happy to release Go 1.23,
which you can get by visiting the [download page](/dl/).

If you already have Go 1.22 or Go 1.21 installed on your machine,
you can also try `go get toolchain@go1.23.0` in an existing module.
This will download the new toolchain and let you begin using it
in your module right away. At some later point, you can follow up
with `go get go@1.23.0` when you're ready to fully switch to Go 1.23
and have that be your module's minimum required Go version.
See [Managing Go version module requirements with go get](/doc/toolchain#get)
for more information on this functionality.

Go 1.23 comes with many improvements over Go 1.22. Some of the highlights include:

## Language changes

-	<!-- go.dev/issue/61405, go.dev/issue/61897, go.dev/issue/61899, go.dev/issue/61900 -->
	Range expressions in a "for-range" loop may now be iterator functions,
	such as `func(func(K) bool)`.
	This supports user-defined iterators over arbitrary sequences.
	There are several additions to the standard `slices` and `maps`
	packages that work with iterators, as well as a new `iter` package.
	As an example, if you wish to collect the keys of a map `m` into a slice
	and then sort its values, you can do that in Go 1.23 with `slices.Sorted(maps.Keys(m))`.

	Go 1.23 also includes preview support for generic type aliases.

	Read more about [language changes](/doc/go1.23#language) and [iterators](/doc/go1.23#iterators)
	in the release notes.

## Tool improvements

-	<!-- go.dev/issue/58894 -->
	Starting with Go 1.23, it's possible for the Go toolchain to collect usage and breakage
	statistics to help understand how the Go toolchain is used, and how well it is working.
	This is Go telemetry, an _opt-in system_. Please consider opting in to help us keep Go
	working well and better understand Go usage.
	Read more on [Go telemetry](/doc/go1.23#telemetry) in the release notes.
-	The `go` command has new conveniences. For example, running `go env -changed` makes it easier to
	see only those settings whose effective value differs from the default value, and
	`go mod tidy -diff` helps determine the necessary changes to the go.mod and go.sum files
	without modifying them.
	Read more on the [Go command](/doc/go1.23#go-command) in the release notes.
-	The `go vet` subcommand now reports symbols that are too new for the intended Go version.
	Read more on [tools](/doc/go1.23#tools) in the release notes.

## Standard library improvements

-	Go 1.23 improves the implementation of `time.Timer` and `time.Ticker`.
	Read more on [timer changes](/doc/go1.23#timer-changes) in the release notes.
- 	There are a total of 3 new packages in the Go 1.23 standard library: `iter`, `structs`, and `unique`.
	Package `iter` is mentioned above.
	Package `structs` defines marker types to modify the properties of a struct.
	Package `unique` provides facilities for canonicalizing ("interning") comparable
	values.
	Read more on [new standard library packages](/doc/go1.23#new-unique-package)
	in the release notes.
-	There are many improvements and additions to the standard library enumerated
	in the [minor changes to the library](/doc/go1.23#minor_library_changes)
	section of the release notes.
	The “Go, Backwards Compatibility, and GODEBUG” documentation
	enumerates [new to Go 1.23 GODEBUG settings](/doc/godebug#go-123).
-	<!-- go.dev/issue/65573 -->
	Go 1.23 supports the new `godebug` directive in `go.mod` and `go.work` files to
	allow separate control of the default GODEBUGs and the “go” directive of `go.mod`,
	in addition to `//go:debug` directive comments made available two releases ago (Go 1.21).
	See the updated documentation on [Default GODEBUG Values](/doc/godebug#default).

## More improvements and changes

-	Go 1.23 adds experimental support for OpenBSD on 64-bit RISC-V (`openbsd/riscv64`).
	There are several minor changes relevant to Linux, macOS, ARM64, RISC-V, and WASI.
	Read more on [ports](/doc/go1.23#ports) in the release notes.
-	Build time when using profile-guided optimization (PGO) is reduced, and performance
	with PGO on 386 and amd64 architectures is improved.
	Read more on [runtime, compiler, and linker](/doc/go1.23#runtime) in the release notes.

We encourage everyone to read the [Go 1.23 release notes](/doc/go1.23) for the
complete and detailed information on these changes, and everything else that's
new to Go 1.23.

Over the next few weeks, look out for follow-up blog posts that will go in more depth
on some of the topics mentioned here, including “range-over-func”, the new `unique` package,
Go 1.23 timer implementation changes, and more.

---

Thank you to everyone who contributed to this release by writing code and
documentation, reporting bugs, sharing feedback, and testing the release
candidates. Your efforts helped to ensure that Go 1.23 is as stable as possible.
As always, if you notice any problems, please [file an issue](/issue/new).

Enjoy Go 1.23!
