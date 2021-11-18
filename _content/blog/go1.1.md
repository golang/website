---
title: Go 1.1 is released
date: 2013-05-13
by:
- Andrew Gerrand
tags:
- release
summary: Go 1.1 is faster, less picky about return statements, and adds method expressions.
---


It is our great pleasure to announce the release of Go 1.1.

{{image "go1.1/gopherbiplane5.jpg"}}

In March last year we released Go 1.0, and since then we have released three
minor "point releases".
The point releases were made to fix only critical issues,
so the Go 1.0.3 you use today is still, in essence,
the Go 1.0 we released in March 2012.

Go 1.1 includes many improvements over 1.0.

The most significant improvements are performance-related.
We have made optimizations in the compiler and linker,
garbage collector, goroutine scheduler, map implementation,
and parts of the standard library.
It is likely that your Go code will run noticeably faster when built with Go 1.1.

There are some minor changes to the language itself,
two of which are worth singling out here:
the [changes to return requirements](/doc/go1.1#return) will
lead to more succinct and correct programs,
and the introduction of [method values](/doc/go1.1#method_values) provides
an expressive way to bind a method to its receiver as a function value.

Concurrent programming is safer in Go 1.1 with the addition of a race
detector for finding memory synchronization errors in your programs.
We will discuss the race detector more in an upcoming article,
but for now [the manual](/doc/articles/race_detector.html) is
a great place to get started.

The tools and standard library have been improved and expanded.
You can read the full story in the [release notes](/doc/go1.1).

As per our [compatibility guidelines](/doc/go1compat.html),
Go 1.1 remains compatible with Go 1.0 and we recommend all Go users upgrade to the new release.

All this would not have been possible without the help of our contributors from
the open source community.
Since Go 1.0, the core received more than 2600 commits from 161 people outside Google.
Thank you everyone for your time and effort.
In particular, we would like to thank Shenghou Ma,
Rémy Oudompheng, Dave Cheney, Mikio Hara,
Alex Brainman, Jan Ziak, and Daniel Morsing for their outstanding contributions.

To grab the new release, follow the usual [installation instructions](/doc/install). Happy hacking!

_Thanks to Renée French for the gopher!_
