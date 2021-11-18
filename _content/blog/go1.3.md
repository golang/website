---
title: Go 1.3 is released
date: 2014-06-18
by:
- Andrew Gerrand
summary: Go 1.3 adds better performance, static analysis in godoc, and more.
---


Today we are happy to announce the release of [Go 1.3](/doc/go1.3).
This release comes six months after our last major release and provides better
performance, improved tools, support for running Go in new environments, and more.
All Go users should upgrade to Go 1.3.
You can grab the release from our [downloads page](/dl/) and
find the full list of improvements and fixes in the
[release notes](/doc/go1.3).
What follows are some highlights.

[Godoc](https://godoc.org/code.google.com/p/go.tools/cmd/godoc),
the Go documentation server, now performs static analysis.
When enabled with the -analysis flag, analysis results are presented
in both the source and package documentation views, making it easier
than ever to navigate and understand Go programs.
See [the documentation](/lib/godoc/analysis/help.html) for the details.

The gc toolchain now supports the Native Client (NaCl) execution sandbox on the
32- and 64-bit Intel architectures.
This permits the safe execution of untrusted code, useful in environments such as the
[Playground](https://blog.golang.org/playground).
To set up NaCl on your system see the [NativeClient wiki page](/wiki/NativeClient).

Also included in this release is experimental support for the DragonFly BSD,
Plan 9, and Solaris operating systems. To use Go on these systems you must
[install from source](/doc/install/source).

Changes to the runtime have improved the
[performance](/doc/go1.3#performance) of Go binaries,
with an improved garbage collector, a new
["contiguous" goroutine stack management strategy](/s/contigstacks),
a faster race detector, and improvements to the regular expression engine.

As part of the general [overhaul](/s/go13linker) of the Go
linker, the compilers and linkers have been refactored. The instruction
selection phase that was part of the linker has been moved to the compiler.
This can speed up incremental builds for large projects.

The [garbage collector](/doc/go1.3#garbage_collector) is now
precise when examining stacks (collection of the heap has been precise since Go
1.1), meaning that a non-pointer value such as an integer will never be
mistaken for a pointer and prevent unused memory from being reclaimed. This
change affects code that uses package unsafe; if you have unsafe code you
should read the [release notes](/doc/go1.3#garbage_collector)
carefully to see if your code needs updating.

We would like to thank the many people who contributed to this release;
it would not have been possible without your help.

So, what are you waiting for?
Head on over to the [downloads page](/dl/) and start hacking.
