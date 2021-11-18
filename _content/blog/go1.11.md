---
title: Go 1.11 is released
date: 2018-08-24
by:
- Andrew Bonventre
summary: Go 1.11 adds preliminary support for Go modules, WebAssembly, and more.
---


Who says releasing on Friday is a bad idea?

Today the Go team is happy to announce the release of Go 1.11.
You can get it from the [download page](/dl/).

There are many changes and improvements to the toolchain,
runtime, and libraries, but two features stand out as being especially exciting:
modules and WebAssembly support.

This release adds preliminary support for a [new concept called “modules,”](/doc/go1.11#modules)
an alternative to GOPATH with integrated support for versioning and package distribution.
Module support is considered experimental,
and there are still a few rough edges to smooth out,
so please make liberal use of the [issue tracker](/issue/new).

Go 1.11 also adds an experimental port to [WebAssembly](/doc/go1.11#wasm) (`js/wasm`).
This allows programmers to compile Go programs to a binary format compatible with four major web browsers.
You can read more about WebAssembly (abbreviated “Wasm”) at [webassembly.org](https://webassembly.org/)
and see [this wiki page](/wiki/WebAssembly) on how to
get started with using Wasm with Go.
Special thanks to [Richard Musiol](https://github.com/neelance) for contributing the WebAssembly port!

We also want to thank everyone who contributed to this release by writing code,
filing bugs, providing feedback, and/or testing the betas and release candidates.
Your contributions and diligence helped to ensure that Go 1.11 is as bug-free as possible.
That said, if you do notice any problems, please [file an issue](/issues/new).

For more detail about the changes in Go 1.11, see the [release notes](/doc/go1.11).

Have a wonderful weekend and enjoy the release!
