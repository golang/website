---
title: Debugging Go programs with the GNU Debugger
date: 2011-10-30
by:
- Andrew Gerrand
tags:
- debug
- gdb
- technical
summary: Announcing a new article about debugging Go programs with GDB.
---


Last year we [reported](https://blog.golang.org/2010/11/debugging-go-code-status-report.html)
that Go's [gc](/cmd/gc/)/[ld](/cmd/6l/)
toolchain produces DWARFv3 debugging information that can be read by the GNU Debugger (GDB).
Since then, work has continued steadily on improving support for debugging Go code with GDB.
Among the improvements are the ability to inspect goroutines and to print
native Go data types,
including structs, slices, strings, maps,
interfaces, and channels.

To learn more about Go and GDB, see the [Debugging with GDB](/doc/debugging_with_gdb.html) article.
