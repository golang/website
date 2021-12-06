---
title: Go 1.8 is released
date: 2017-02-16
by:
- Chris Broadfoot
summary: Go 1.8 adds faster non-x86 compiled code, sub-millisecond garbage collection pauses, HTTP/2 push, and more.
---


Today the Go team is happy to announce the release of Go 1.8.
You can get it from the [download page](/dl/).
There are significant performance improvements and changes across the standard library.

The compiler back end introduced in [Go 1.7](https://blog.golang.org/go1.7) for 64-bit x86 is now used
on all architectures, and those architectures should see significant [performance improvements](/doc/go1.8#compiler).
For instance, the CPU time required by our benchmark programs was reduced by 20-30% on 32-bit ARM systems.
There are also some modest performance improvements in this release for 64-bit x86 systems.
The compiler and linker have been made faster.
Compile times should be improved by about 15% over Go 1.7.
There is still more work to be done in this area: expect faster compilation speeds in future releases.

Garbage collection pauses should be [significantly shorter](/doc/go1.8#gc),
usually under 100 microseconds and often as low as 10 microseconds.

The HTTP server adds support for [HTTP/2 Push](/doc/go1.8#h2push),
allowing servers to preemptively send responses to a client.
This is useful for minimizing network latency by eliminating roundtrips.
The HTTP server also adds support for [graceful shutdown](/doc/go1.8#http_shutdown),
allowing servers to minimize downtime by shutting down only after serving all requests that are in flight.

[Contexts](/pkg/context/) (added to the standard library in Go 1.7)
provide a cancellation and timeout mechanism.
Go 1.8 [adds](/doc/go1.8#more_context) support for contexts in more parts of the standard library,
including the [`database/sql`](/pkg/database/sql) and [`net`](/pkg/net) packages
and [`Server.Shutdown`](http://beta.golang.org/pkg/net/http/#Server.Shutdown) in the `net/http` package.

It's now much simpler to sort slices using the newly added [`Slice`](/pkg/sort/#Slice)
function in the `sort` package. For example, to sort a slice of structs by their `Name` field:

{{raw `
	sort.Slice(s, func(i, j int) bool { return s[i].Name < s[j].Name })
`}}

Go 1.8 includes many more additions, improvements, and fixes.
Find the complete set of changes, and more information about the improvements listed above, in the
[Go 1.8 release notes](/doc/go1.8.html).

To celebrate the release, Go User Groups around the world are holding [release parties](https://github.com/golang/go/wiki/Go-1.8-release-party) this week.
Release parties have become a tradition in the Go community, so if you missed out this time, keep an eye out when 1.9 nears.

Thank you to over 200 contributors who helped with this release.
