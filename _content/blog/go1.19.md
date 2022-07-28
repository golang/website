---
title: Go 1.19 is released!
date: 2022-08-02
by:
- The Go Team
summary: Go 1.19 adds richer doc comments, performance improvements, and more.
---

Today the Go team is thrilled to release Go 1.19,
which you can get by visiting the [download page](/dl/).

Go 1.19 refines and improves our massive [Go 1.18 release](/blog/go1.18) earlier this year.
We focused Go 1.19's generics development on addressing the subtle issues
and corner cases reported to us by the community,
as well as important performance improvements (up to 20% for some generic programs).

Doc comments now support [links, lists, and clearer heading syntax](/doc/comment).
This change helps users write clearer, more navigable doc comments,
especially in packages with large APIs.
As part of this change `gofmt` now reformats doc comments to apply a
standard formatting to uses of these features.
See “[Go Doc Comments](/doc/comment)” for all the details.

[Go's memory model](/ref/mem) now explicitly defines
the behavior of the [sync/atomic package](/pkg/sync/atomic/).
The formal definition of the happens-before relation has been revised
to align with the memory models used by C, C++, Java, JavaScript, Rust, and Swift.
Existing programs are unaffected.
Along with the memory model update, there are
[new types in the sync/atomic package](/doc/go1.19#atomic_types),
such as [atomic.Int64](/pkg/sync/atomic/#Int64) and [atomic.Pointer[T]](/pkg/sync/atomic/#Pointer),
to make it easier to use atomic values.

For [security reasons](/blog/path-security), the os/exec package
no longer respects relative paths in PATH lookups.
See the [package documentation](/pkg/os/exec/#hdr-Executables_in_the_current_directory)
for details.
Existing uses of [golang.org/x/sys/execabs](https://pkg.go.dev/golang.org/x/sys/execabs)
can be moved back to os/exec in programs that only build using Go 1.19 or later.

The garbage collector has added support for a soft memory limit,
discussed in detail in [the new garbage collection guide](/doc/gc-guide#Memory_limit).
The limit can be particularly helpful for optimizing Go programs to
run as efficiently as possible in containers with dedicated amounts of memory.

The new build constraint `unix` is satisfied when the target operating system (`GOOS`)
is any Unix-like system.
Today, Unix-like means all of
Go's target operating systems except `js`, `plan9`, `windows`, and `zos`.

Finally, Go 1.19 includes a wide variety of performance and implementation improvements, including
dynamic sizing of initial goroutine stacks to reduce stack copying,
automatic use of additional file descriptors on most Unix systems,
jump tables for large switch statements on x86-64 and ARM64,
support for debugger-injected function calls on ARM64,
register ABI support on RISC-V,
and experimental support for
Linux running on Loongson 64-bit architecture LoongArch (`GOARCH=loong64`).

Thanks to everyone who contributed to this release by writing code, filing bugs, sharing feedback,
and testing the beta and release candidates.
Your efforts helped to ensure that Go 1.19 is as stable as possible.
As always, if you notice any problems, please [file an issue](/issue/new).

Enjoy Go 1.19!
