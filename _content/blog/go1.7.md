---
title: Go 1.7 is released
date: 2016-08-15
by:
- Chris Broadfoot
summary: Go 1.7 adds faster x86 compiled code, context in the standard library, and more.
---


Today we are happy to announce the release of Go 1.7.
You can get it from the [download page](/dl/).
There are several significant changes in this release: a port for
[Linux on IBM z Systems](https://en.wikipedia.org/wiki/IBM_System_z) (s390x),
compiler improvements, the addition of the [context](/pkg/context/) package,
and support for [hierarchical tests and benchmarks](/pkg/testing/#hdr-Subtests_and_Sub_benchmarks).

A new compiler back end, based on [static single-assignment](https://en.wikipedia.org/wiki/Static_single_assignment_form) form (SSA),
has been under development for the past year.
By representing a program in SSA form, a compiler may perform advanced optimizations more easily.
This new back end generates more compact, more efficient code that includes
optimizations like
[bounds check elimination](https://en.wikipedia.org/wiki/Bounds-checking_elimination) and
[common subexpression elimination](https://en.wikipedia.org/wiki/Common_subexpression_elimination).
We observed a 5–35% speedup across our [benchmarks](/test/bench/go1/).
For now, the new backend is only available for the 64-bit x86 platform ("amd64"),
but we’re planning to convert more architecture backends to SSA in future releases.

The compiler front end uses a new, more compact export data format, and
processes import declarations more efficiently.
While these [changes across the compiler toolchain](/doc/go1.7#compiler) are mostly invisible,
users have [observed](http://dave.cheney.net/2016/04/02/go-1-7-toolchain-improvements)
a significant speedup in compile time and a reduction in binary size by as much as 20–30%.

Programs should run a bit faster due to speedups in the garbage collector and optimizations in the standard library.
Programs with many idle goroutines will experience much shorter garbage collection pauses than in Go 1.6.

Over the past few years, the [golang.org/x/net/context](https://godoc.org/golang.org/x/net/context/)
package has proven to be essential to many Go applications.
Contexts are used to great effect in applications related to networking, infrastructure, and microservices
(such as [Kubernetes](http://kubernetes.io/) and [Docker](https://www.docker.com/)).
They make it easy to enable cancellation, timeouts, and passing request-scoped data.
To make use of contexts within the standard library and to encourage more extensive use,
the package has been moved from the [x/net](https://godoc.org/golang.org/x/net/context/) repository
to the standard library as the [context](/pkg/context/) package.
Support for contexts has been added to the
[net](/pkg/net/),
[net/http](/pkg/net/http/), and
[os/exec](/pkg/os/exec/) packages.
For more information about contexts, see the [package documentation](/pkg/context)
and the Go blog post [_Go Concurrency Patterns: Context_](https://blog.golang.org/context).

Go 1.5 introduced experimental support for a ["vendor" directory](/cmd/go/#hdr-Vendor_Directories),
enabled by the `GO15VENDOREXPERIMENT` environment variable.
Go 1.6 enabled this behavior by default, and in Go 1.7, this switch has been removed and the "vendor" behavior is always enabled.

Go 1.7 includes many more additions, improvements, and fixes.
Find the complete set of changes, and details of the points above, in the
[Go 1.7 release notes](/doc/go1.7.html).

Finally, the Go team would like thank everyone who contributed to the release.
170 people contributed to this release, including 140 from the Go community.
These contributions ranged from changes to the compiler and linker, to the standard library, to documentation, and code reviews.
We welcome contributions; if you'd like to get involved, check out the
[contribution guidelines](/doc/contribute.html).
