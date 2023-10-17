---
title: Go 1.4 is released
date: 2014-12-10
by:
- Andrew Gerrand
summary: Go 1.4 adds support for Android, go generate, optimizations, and more.
---


Today we announce Go 1.4, the fifth major stable release of Go, arriving six
months after our previous major release [Go 1.3](https://blog.golang.org/go1.3).
It contains a small language change, support for more operating systems
and processor architectures, and improvements to the tool chain and libraries.
As always, Go 1.4 keeps the promise of compatibility, and almost everything
will continue to compile and run without change when moved to 1.4.
For the full details, see the [Go 1.4 release notes](/doc/go1.4).

The most notable new feature in this release is official support for Android.
Using the support in the core and the libraries in the
[golang.org/x/mobile](https://godoc.org/golang.org/x/mobile) repository,
it is now possible to write simple Android apps using only Go code.
At this stage, the support libraries are still nascent and under heavy development.
Early adopters should expect a bumpy ride, but we welcome the community to get involved.

The language change is a tweak to the syntax of for-range loops.
You may now write "for range s {" to loop over each item from s,
without having to assign the value, loop index, or map key.
See the [release notes](/doc/go1.4#forrange) for details.

The go command has a new subcommand, go generate, to automate the running of
tools to generate source code before compilation.
For example, it can be used to automate the generation of String methods for
typed constants using the
[new stringer tool](https://godoc.org/golang.org/x/tools/cmd/stringer/).
For more information, see the [design document](/s/go1.4-generate).

Most programs will run about the same speed or slightly faster in 1.4 than in
1.3; some will be slightly slower.
There are many changes, making it hard to be precise about what to expect.
See the [release notes](/doc/go1.4#performance) for more discussion.

And, of course, there are many more improvements and bug fixes.

In case you missed it, a few weeks ago the sub-repositories were moved to new locations.
For example, the go.tools packages are now imported from "golang.org/x/tools".
See the [announcement post](https://groups.google.com/d/msg/golang-announce/eD8dh3T9yyA/HDOEU_ZSmvAJ) for details.

This release also coincides with the project's move from Mercurial to Git (for
source control), Rietveld to Gerrit (for code review), and Google Code to
GitHub (for issue tracking and wiki).
The move affects the core Go repository and its sub-repositories.
You can find the canonical Git repositories at
[go.googlesource.com](https://go.googlesource.com),
and the issue tracker and wiki at the
[golang/go GitHub repo](https://github.com/golang/go).

While development has already moved over to the new infrastructure,
for the 1.4 release we still recommend that users who
[install from source](/doc/install/source)
use the Mercurial repositories.

For App Engine users, Go 1.4 is now available for beta testing.
See [the announcement](https://groups.google.com/d/msg/google-appengine-go/ndtQokV3oFo/25wV1W9JtywJ) for details.

From all of us on the Go team, please enjoy Go 1.4, and have a happy holiday season.
