---
title: Go 1.6 is released
date: 2016-02-17
by:
- Andrew Gerrand
summary: Go 1.6 adds HTTP/2, template blocks, and more.
---


Today we release [Go version 1.6](/doc/go1.6),
the seventh major stable release of Go.
You can grab it right now from the [download page](/dl/).
Although [the release of Go 1.5](https://blog.golang.org/go1.5) six months ago
contained dramatic implementation changes,
this release is more incremental.

The most significant change is support for [HTTP/2](https://http2.github.io/)
in the [net/http package](/pkg/net/http/).
HTTP/2 is a new protocol, a follow-on to HTTP that has already seen
widespread adoption by browser vendors and major websites.
In Go 1.6, support for HTTP/2 is [enabled by default](/doc/go1.6#http2)
for both servers and clients when using HTTPS,
bringing [the benefits](https://http2.github.io/faq/) of the new protocol
to a wide range of Go projects,
such as the popular [Caddy web server](https://caddyserver.com/download).

The template packages have learned some new tricks,
with support for [trimming spaces around template actions](/pkg/text/template/#hdr-Text_and_spaces)
to produce cleaner template output,
and the introduction of the {{raw "[`{{block}}` action]"}}(/pkg/text/template/#hdr-Actions)
that can be used to create templates that build on other templates.
A [new template example program](https://cs.opensource.google/go/x/example/+/master:template) demonstrates these new features.

Go 1.5 introduced [experimental support](/s/go15vendor)
for a “vendor” directory that was enabled by an environment variable.
In Go 1.6, the feature is now [enabled by default](/doc/go1.6#go_command).
Source trees that contain a directory named “vendor” that is not used in accordance with the new feature
will require changes to avoid broken builds (the simplest fix is to rename the directory).

The runtime has added lightweight, best-effort detection of concurrent misuse of maps.
As always, if one goroutine is writing to a map, no other goroutine should be reading or writing the map concurrently.
If the runtime detects this condition, it prints a diagnosis and crashes the program.
The best way to find out more about the problem is to run it under the
[race detector](https://blog.golang.org/race-detector),
which will more reliably identify the race and give more detail.

The runtime has also changed how it prints program-ending panics.
It now prints only the stack of the panicking goroutine, rather than all existing goroutines.
This behavior can be configured using the
[GOTRACEBACK](/pkg/runtime/#hdr-Environment_Variables) environment variable
or by calling the [debug.SetTraceback](/pkg/runtime/debug/#SetTraceback) function.

Users of cgo should be aware of major changes to the rules for sharing pointers between Go and C code.
The rules are designed to ensure that such C code can coexist with Go's garbage collector
and are checked during program execution, so code may require changes to avoid crashes.
See the [release notes](/doc/go1.6#cgo) and
[cgo documentation](/cmd/cgo/#hdr-Passing_pointers) for the details.

The compiler, linker, and go command have a new `-msan` flag
analogous to `-race` and only available on linux/amd64,
that enables interoperation with the
[Clang MemorySanitizer](http://clang.llvm.org/docs/MemorySanitizer.html).
This is useful for testing a program containing suspect C or C++ code.
You might like to try it while testing your cgo code with the new pointer rules.

Performance of Go programs built with Go 1.6 remains similar to those built with Go 1.5.
Garbage-collection pauses are even lower than with Go 1.5,
but this is particularly noticeable for programs using large amounts of memory.
With regard to the performance of the compiler tool chain,
build times should be similar to those of Go 1.5.

The algorithm inside [sort.Sort](/pkg/sort/#Sort)
was improved to run about 10% faster,
but the change may break programs that expect a specific ordering
of equal but distinguishable elements.
Such programs should refine their `Less` methods to indicate the desired ordering
or use [sort.Stable](/pkg/sort/#Stable)
to preserve the input order for equal values.

And, of course, there are many more additions, improvements, and fixes.
You can find them all in the comprehensive [release notes](/doc/go1.6).

To celebrate the release,
[Go User Groups around the world](https://github.com/golang/go/wiki/Go-1.6-release-party)
are holding release parties on the 17th of February.
Online, the Go contributors are hosting a question and answer session
on the [golang subreddit](https://reddit.com/r/golang) for the next 24 hours.
If you have questions about the project, the release, or just Go in general,
then please [join the discussion](https://www.reddit.com/r/golang/comments/46bd5h/ama_we_are_the_go_contributors_ask_us_anything/).

Thanks to everyone that contributed to the release.
Happy hacking.
