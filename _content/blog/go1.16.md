---
title: Go 1.16 is released
date: 2021-02-16
by:
- Matt Pearring
- Dmitri Shuralyov
summary: Go 1.16 adds embedded files, Apple Silicon support, and more.
---


Today the Go team is very happy to announce the release of Go 1.16.
You can get it from the [download page](/dl/).

The new
[embed package](/doc/go1.16#library-embed)
provides access to files embedded at compile time using the new `//go:embed` directive.
Now it is easy to bundle supporting data files into your Go programs,
making developing with Go even smoother.
You can get started using the
[embed package documentation](https://pkg.go.dev/embed).
Carl Johnson has also written a nice tutorial,
“[How to use Go embed](https://blog.carlmjohnson.net/post/2021/how-to-use-go-embed/)”.

Go 1.16 also adds
[macOS ARM64 support](/doc/go1.16#darwin)
(also known as Apple silicon).
Since Apple’s announcement of their new arm64 architecture, we have been working closely with them to ensure Go is fully supported; see our blog post
“[Go on ARM and Beyond](https://blog.golang.org/ports)”
for more.

Note that Go 1.16
[requires use of Go modules by default](/doc/go1.16#modules),
now that, according to our 2020 Go Developer Survey,
96% of Go developers have made the switch.
We recently added official documentation for [developing and publishing modules](/doc/modules/developing).

Finally, there are many other improvements and bug fixes,
including builds that are up to 25% faster and use as much as 15% less memory.
For the complete list of changes and more information about the improvements above,
see the
[Go 1.16 release notes](/doc/go1.16).

We want to thank everyone who contributed to this release by writing code
filing bugs, providing feedback, and testing the beta and release candidate.

Your contributions and diligence helped to ensure that Go 1.16 is as stable as possible.
That said, if you notice any problems, please
[file an issue](/issue/new).

We hope you enjoy the new release!
