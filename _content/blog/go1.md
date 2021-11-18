---
title: Go version 1 is released
date: 2012-03-28
by:
- Andrew Gerrand
tags:
- release
- go1
summary: "A major milestone: announcing Go 1, the first stable version of Go."
---


{{image "go1/gophermega.jpg"}}

Today marks a major milestone in the development of the Go programming language.
We're announcing Go version 1, or Go 1 for short,
which defines a language and a set of core libraries to provide a stable
foundation for creating reliable products,
projects, and publications.

Go 1 is the first release of Go that is available in supported binary distributions.
They are available for Linux, FreeBSD, Mac OS X and,
we are thrilled to announce, Windows.

The driving motivation for Go 1 is stability for its users.
People who write Go 1 programs can be confident that those programs will
continue to compile and run without change,
in many environments, on a time scale of years.
Similarly, authors who write books about Go 1 can be sure that their examples
and explanations will be helpful to readers today and into the future.

Forward compatibility is part of stability.
Code that compiles in Go 1 should, with few exceptions,
continue to compile and run throughout the lifetime of that version,
even as we issue updates and bug fixes such as Go version 1.1, 1.2, and so on.
The [Go 1 compatibility document](/doc/go1compat.html)
explains the compatibility guidelines in more detail.

Go 1 is a representation of Go as it is used today,
not a major redesign.
In its planning, we focused on cleaning up problems and inconsistencies
and improving portability.
There had long been many changes to Go that we had designed and prototyped
but not released because they were backwards-incompatible.
Go 1 incorporates these changes, which provide significant improvements
to the language and libraries but sometimes introduce incompatibilities for old programs.
Fortunately, the [go fix](/cmd/go/#Run_go_tool_fix_on_packages)
tool can automate much of the work needed to bring programs up to the Go 1 standard.

Go 1 introduces changes to the language (such as new types for [Unicode characters](/doc/go1.html#rune)
and [errors](/doc/go1.html#errors)) and the standard
library (such as the new [time package](/doc/go1.html#time)
and renamings in the [strconv package](/doc/go1.html#strconv)).
Also, the package hierarchy has been rearranged to group related items together,
such as moving the networking facilities,
for instance the [rpc package](/pkg/net/rpc/),
into subdirectories of net.
A complete list of changes is documented in the [Go 1 release notes](/doc/go1.html).
That document is an essential reference for programmers migrating code from
earlier versions of Go.

We also restructured the Go tool suite around the new [go command](/doc/go1.html#cmd_go),
a program for fetching, building, installing and maintaining Go code.
The go command eliminates the need for Makefiles to write Go code because
it uses the Go program source itself to derive the build instructions.
No more build scripts!

Finally, the release of Go 1 triggers a new release of the [Google App Engine SDK](https://developers.google.com/appengine/docs/go).
A similar process of revision and stabilization has been applied to the
App Engine libraries,
providing a base for developers to build programs for App Engine that will run for years.

Go 1 is the result of a major effort by the core Go team and our many contributors
from the open source community.
We thank everyone who helped make this happen.

There has never been a better time to be a Go programmer.
Everything you need to get started is at [golang.org](/).
