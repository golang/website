---
title: "Go: one year ago today"
date: 2010-11-10
by:
- Andrew Gerrand
tags:
- birthday
summary: Happy 1st birthday, Go!
---


On the 10th of November 2009 we launched the Go project:
an open-source programming language with a focus on simplicity and efficiency.
The intervening year has seen a great many developments both in the Go project
itself and in its community.

We set out to build a language for systems programming - the kinds of programs
one might typically write in C or C++ - and we were surprised by Go’s
utility as a general purpose language.
We had anticipated interest from C, C++, and Java programmers,
but the flurry of interest from users of dynamically-typed languages like
Python and JavaScript was unexpected.
Go’s combination of native compilation,
static typing, memory management, and lightweight syntax seemed to strike
a chord with a broad cross-section of the programming community.

That cross-section grew to become a dedicated community of enthusiastic Go coders.
Our [mailing list](http://groups.google.com/group/golang-nuts) has over 3,800 members,
with around 1,500 posts each month.
The project has over 130 contributors
(people who have submitted code or documentation),
and of the 2,800 commits since launch almost one third were contributed
by programmers outside the core team.
To get all that code into shape, nearly 14,000 emails were exchanged on
our [development mailing list](http://groups.google.com/group/golang-dev).

Those numbers reflect a labor whose fruits are evident in the project’s code base.
The compilers have improved substantially,
with faster and more efficient code generation,
more than one hundred reported bugs fixed,
and support for a widening range of operating systems and architectures.
The Windows port is approaching completion thanks to a dedicated group of
contributors (one of whom became our first non-Google committer to the project).
The ARM port has also made great progress,
recently reaching the milestone of passing all tests.

The Go tool set has been expanded and improved.
The Go documentation tool, [godoc](/cmd/godoc/),
now supports the documentation of other source trees (you can browse and
search your own code) and provides a ["code walk"](/doc/codewalk/)
interface for presenting tutorial materials (among many more improvements).
[Goinstall](/cmd/goinstall/) ,
a new package management tool, allows users to install and update external
packages with a single command.
[Gofmt](/cmd/gofmt/),
the Go pretty-printer, now makes syntactic simplifications where possible.
[Goplay](/misc/goplay/),
a web-based “compile-as-you-type” tool,
is a convenient way to experiment with Go for those times when you don’t
have access to the [Go Playground](/doc/play/).

The standard library has grown by over 42,000 lines of code and includes
20 new [packages](/pkg/).
Among the additions are the [jpeg](/pkg/image/jpeg/),
[jsonrpc](/pkg/rpc/jsonrpc/),
[mime](/pkg/mime/), [netchan](/pkg/netchan/),
and [smtp](/pkg/smtp/) packages,
as well as a slew of new [cryptography](/pkg/crypto/) packages.
More generally, the standard library has been continuously refined and revised
as our understanding of Go’s idioms deepens.

The debugging story has gotten better, too.
Recent improvements to the DWARF output of the gc compilers make the GNU debugger,
GDB, useful for Go binaries, and we’re actively working on making that
debugging information more complete.
(See the [ recent blog post](https://blog.golang.org/2010/11/debugging-go-code-status-report.html) for details.)

It’s now easier than ever to link against existing libraries written in
languages other than Go.
Go support is in the most recent [SWIG](http://www.swig.org/) release,
version 2.0.1, making it easier to link against C and C++ code,
and our [cgo](/cmd/cgo/) tool has seen many fixes and improvements.

[Gccgo](/doc/install/gccgo),
the Go front end for the GNU C Compiler, has kept pace with the gc compiler
as a parallel Go implementation.
It now has a working garbage collector, and has been accepted into the GCC core.
We’re now working toward making [gofrontend](http://code.google.com/p/gofrontend/)
available as a BSD-licensed Go compiler front end,
fully decoupled from GCC.

Outside the Go project itself Go is starting to be used to build real software.
There are more than 200 Go programs and libraries listed on our [Project dashboard](http://godashboard.appspot.com/project),
and hundreds more on [Google Code](http://code.google.com/hosting/search?q=label:Go)
and [GitHub](https://github.com/search?q=language:Go).
On our mailing list and IRC channel you can find coders from around the
world who use Go for their programming projects.
(See our [guest blog post](https://blog.golang.org/2010/10/real-go-projects-smarttwitter-and-webgo.html)
from last month for a real-world example.) Internally at Google there are
several teams that choose Go for building production software,
and we have received reports from other companies that are developing sizable systems in Go.
We have also been in touch with several educators who are using Go as a teaching language.

The language itself has grown and matured, too.
In the past year we have received many feature requests.
But Go is a small language, and we’ve worked hard to ensure that any new
feature strikes the right compromise between simplicity and utility.
Since the launch we have made a number of language changes,
many of which were driven by feedback from the community.

  - Semicolons are now optional in almost all instances. [spec](/doc/go_spec.html#Semicolons)
  - The new built-in functions `copy` and `append` make management of slices
    more efficient and straightforward.
    [spec](/doc/go_spec.html#Appending_and_copying_slices)
  - The upper and lower bounds may be omitted when making a sub-slice.
    This means that `s[:]` is shorthand for `s[0:len(s)]`.
    [spec](/doc/go_spec.html#Slices)
  - The new built-in function `recover` complements `panic` and `defer` as
    an error handling mechanism.
    [blog](https://blog.golang.org/2010/08/defer-panic-and-recover.html),
    [spec](/doc/go_spec.html#Handling_panics)
  - The new complex number types (`complex`,
    `complex64`, and `complex128`) simplify certain mathematical operations.
    [spec](/doc/go_spec.html#Complex_numbers),
    [spec](/doc/go_spec.html#Imaginary_literals)
  - The composite literal syntax permits the omission of redundant type information
    (when specifying two-dimensional arrays, for example).
    [release.2010-10-27](/doc/devel/release.html#2010-10-27),
    [spec](/doc/go_spec.html#Composite_literals)
  - A general syntax for variable function arguments (`...T`) and their propagation
    (`v...`) is now specified.
    [spec](/doc/go_spec.html#Function_Types),
    [ spec](/doc/go_spec.html#Passing_arguments_to_..._parameters),
    [release.2010-09-29](/doc/devel/release.html#2010-09-29)

Go is certainly ready for production use,
but there is still room for improvement.
Our focus for the immediate future is making Go programs faster and more
efficient in the context of high performance systems.
This means improving the garbage collector,
optimizing generated code, and improving the core libraries.
We’re also exploring some further additions to the type system to make
generic programming easier.
A lot has happened in a year; it’s been both thrilling and satisfying.
We hope that this coming year will be even more fruitful than the last.

_If you’ve been meaning to get [back] into Go, now is a great time to do so! Check out the_
[_Documentation_](/doc/docs.html) _and_ [_Getting Started_](/doc/install.html)
_pages for more information, or just go nuts in the_ [_Go Playground_](/doc/play/).
