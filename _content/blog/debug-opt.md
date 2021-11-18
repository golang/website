---
title: Debugging what you deploy in Go 1.12
date: 2019-03-21
by:
- David Chase
tags:
- debug
- technical
summary: Go 1.12 improves support for debugging optimized binaries.
---

## Introduction

Go 1.11 and Go 1.12 make significant progress toward allowing developers
to debug the same optimized binaries that they deploy to production.

As the Go compiler has become increasingly aggressive in producing faster binaries,
we've lost ground in debuggability.
In Go 1.10, users needed to disable optimizations entirely in order to have
a good debugging experience from interactive tools like Delve.
But users shouldn’t have to trade performance for debuggability,
especially when running production services.
If your problem is occurring in production,
you need to debug it in production, and that shouldn’t require deploying
unoptimized binaries.

For Go 1.11 and 1.12, we focused on improving the debugging experience on
optimized binaries (the default setting of the Go compiler).
Improvements include

  - More accurate value inspection, in particular for arguments at function entry;
  - More precisely identifying statement boundaries so that stepping is less
    jumpy and breakpoints more often land where the programmer expects;
  - And preliminary support for Delve to call Go functions (goroutines and
    garbage collection make this trickier than it is in C and C++).

## Debugging optimized code with Delve

[Delve](https://github.com/go-delve/delve) is a debugger for Go on x86
supporting both Linux and macOS.
Delve is aware of goroutines and other Go features and provides one of the
best Go debugging experiences.
Delve is also the debugging engine behind [GoLand](https://www.jetbrains.com/go/),
[VS Code](https://code.visualstudio.com/),
and [Vim](https://github.com/fatih/vim-go).

Delve normally rebuilds the code it is debugging with `-gcflags "all=-N -l"`,
which disables inlining and most optimizations.
To debug optimized code with delve, first build the optimized binary,
then use `dlv exec your_program` to debug it.
Or, if you have a core file from a crash,
you can examine it with `dlv core your_program your_core`.
With 1.12 and the latest Delve releases, you should be able to examine many variables,
even in optimized binaries.

## Improved value inspection

When debugging optimized binaries produced by Go 1.10,
variable values were usually completely unavailable.
In contrast, starting with Go 1.11, variables can usually be examined even
in optimized binaries,
unless they’ve been optimized away completely.
In Go 1.11 the compiler began emitting DWARF location lists so debuggers
can track variables as they move in and out of registers and reconstruct
complex objects that are split across different registers and stack slots.

## Improved stepping

This shows an example of stepping through a simple function in a debugger in 1.10,
with flaws (skipped and repeated lines) highlighted by red arrows.

{{image "debug-opt/stepping.svg" 450}}

Flaws like this make it easy to lose track of where you are when stepping
through a program and interfere with hitting breakpoints.

Go 1.11 and 1.12 record statement boundary information and do a better job
of tracking source line numbers through optimizations and inlining.
As a result, in Go 1.12, stepping through this code stops on every line
and does so in the order you would expect.

## Function calls

Function call support in Delve is still under development, but simple cases work.  For example:

	(dlv) call fib(6)
	> main.main() ./hello.go:15 (PC: 0x49d648)
	Values returned:
		~r1: 8

## The path forward

Go 1.12 is a step toward a better debugging experience for optimized binaries
and we have plans to improve it even further.

There are fundamental tradeoffs between debuggability and performance,
so we’re focusing on the highest-priority debugging defects,
and working to collect automated metrics to monitor our progress and catch regressions.

We’re focusing on generating correct information for debuggers about variable locations,
so if a variable can be printed, it is printed correctly.
We’re also looking at making variable values available more of the time,
particularly at key points like call sites,
though in many cases improving this would require slowing down program execution.
Finally, we’re working on improving stepping:
we’re focusing on the order of stepping with panics,
the order of stepping around loops, and generally trying to follow source
order where possible.

## A note on macOS support

Go 1.11 started compressing debug information to reduce binary sizes.
This is natively supported by Delve, but neither LLDB nor GDB support compressed
debug info on macOS.
If you are using LLDB or GDB, there are two workarounds:
build binaries with `-ldflags=-compressdwarf=false`,
or use [splitdwarf](https://godoc.org/golang.org/x/tools/cmd/splitdwarf)
(`go get golang.org/x/tools/cmd/splitdwarf`) to decompress the debug information
in an existing binary.
