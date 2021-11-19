---
title: Go 1.9 is released
date: 2017-08-24
by:
- Francesc Campoy
summary: Go 1.9 adds type aliases, bit intrinsics, optimizations, and more.
---


Today the Go team is happy to announce the release of Go 1.9.
You can get it from the [download page](/dl/).
There are many changes to the language, standard library, runtime, and tooling.
This post covers the most significant visible ones.
Most of the engineering effort put into this release went to improvements of the runtime and tooling,
which makes for a less exciting announcement, but nonetheless a great release.

The most important change to the language is the introduction of type aliases: a feature
created to support gradual code repair. A type alias declaration has the form:

	type T1 = T2

This declaration introduces an alias name `T1` for the type `T2`, in the same way that `byte` has
always been an alias for `uint8`.
The [type alias design document](/design/18130-type-alias) and
[an article on refactoring](/talks/2016/refactor.article) cover this addition in more detail.

The new [math/bits](/pkg/math/bits) package provides bit counting and manipulation functions
for unsigned integers, implemented by special CPU instructions when possible.
For example, on x86-64 systems, `bits.TrailingZeros(x)` uses the
[BSF](https://pdos.csail.mit.edu/6.828/2010/readings/i386/BSF.htm) instruction.

The `sync` package has added a new [Map](/pkg/sync#Map) type, safe for concurrent access.
You can read more about it from its documentation and learn more about why it was created from this
[GopherCon 2017 lightning talk](https://www.youtube.com/watch?v=C1EtfDnsdDs)
([slides](https://github.com/gophercon/2017-talks/blob/master/lightningtalks/BryanCMills-AnOverviewOfSyncMap/An%20Overview%20of%20sync.Map.pdf)).
It is not a general replacement for Go's map type; please see the documentation to learn when it should be used.

The `testing` package also has an addition. The new `Helper` method, added to both
[testing.T](/pkg/testing#T.Helper) and [testing.B](/pkg/testing#B.Helper),
marks the calling function as a test helper function.
When the testing package prints file and line information, it shows the location of the call to a helper function
instead of a line in the helper function itself.

For example, consider this test:

{{code "go1.9/helper_test.go" `/package p/` `$`}}

Because `failure` identifies itself as a test helper, the error message printed during `Test` will indicate line 11,
where `failure` is called, instead of line 7, where `failure` calls `t.Fatal`.

The `time` package now transparently tracks monotonic time in each `Time` value,
making computing durations between two `Time` values a safe operation in the presence of wall clock adjustments.
For example, this code now computes the right elapsed time even across a leap second clock reset:

	start := time.Now()
	f()
	elapsed := time.Since(start)

See the [package docs](http://beta.golang.org/pkg/time/#hdr-Monotonic_Clocks) and
[design document](https://github.com/golang/proposal/blob/master/design/12914-monotonic.md) for details.

Finally, as part of the efforts to make the Go compiler faster, Go 1.9 compiles functions in a package concurrently.

Go 1.9 includes many more additions, improvements, and fixes. Find the complete set of changes,
and more information about the improvements listed above, in the
[Go 1.9 release notes](/doc/go1.9).

To celebrate the release, Go User Groups around the world are holding
[release parties](https://github.com/golang/cowg/blob/master/events/2017-08-go1.9-release-party.md).
