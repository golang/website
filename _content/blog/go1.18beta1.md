---
title: Go 1.18 Beta 1 is available, with generics
date: 2021-12-14
by:
- Russ Cox, for the Go team
summary: Go 1.18 Beta 1 is our first preview of Go 1.18. Please try it and let us know if you find problems.
---

We have just released Go 1.18 Beta 1,
which you can get by visiting the [download page](/dl/#go1.18beta1).

The official Go 1.18 release won't happen for a couple of months yet.
This is the first preview release of Go 1.18, to let you kick the tires,
take it for a spin, and let us know what problems you encounter.
Go 1.18 Beta 1 represents an enormous amount of work
by the entire Go team at Google and Go contributors around the world,
and we're excited to hear what you think.

Go 1.18 Beta 1 is the first preview release containing
Go's new support for [generic code using parameterized types](/blog/why-generics).
Generics are the most significant change to Go since the release of Go 1,
and certainly the largest single language change we've ever made.
With any large, new feature, it is common for new users to discover new bugs,
and we don’t expect generics to be an exception to this rule;
be sure to approach them with appropriate caution.
Also, certain subtle cases, such as specific kinds of recursive generic types,
have been postponed to future releases.
That said, we know of early adopters who have been quite happy,
and if you have use cases that you think are particularly suited to generics,
we hope you will give them a try.
We've published a
[brief tutorial about how to get started with generics](/doc/tutorial/generics)
and gave a
[talk at GopherCon last week](https://www.youtube.com/watch?v=35eIxI_n5ZM&t=1755s).
You can even try it on the
[Go playground in Go dev branch mode](/play/?v=gotip).

Go 1.18 Beta 1 adds built-in support for writing
[fuzzing-based tests](/blog/fuzz-beta),
to automatically find inputs that cause your program to crash or return invalid answers.

Go 1.18 Beta 1 adds a new “[Go workspace mode](/design/45713-workspace)”,
which lets you work with multiple Go modules simultaneously,
an important use case for larger projects.

Go 1.18 Beta 1 contains an expanded `go version -m` command,
which now records build details such as compiler flags.
A program can query its own build details using
[debug.ReadBuildInfo](https://pkg.go.dev/runtime/debug@master#BuildInfo),
and it can now read build details from other binaries using the new
[debug/buildinfo](https://pkg.go.dev/debug/buildinfo@master) package.
This functionality is meant to be the foundation
for any tool that needs to produce a software bill of materials (SBOM) for Go binaries.

Earlier this year, Go 1.17 added a new register-based
calling convention to speed up Go code on x86-64 systems.
Go 1.18 Beta 1 expands that feature to ARM64 and PPC64,
resulting in as much as 20% speed-ups.

Thanks to everyone who contributed to this beta release,
and especially to the team here at Google who has been
working tirelessly for years on making generics a reality.
It's been a long road, we're very happy with the result,
and we hope you like it too.

See the full [draft release notes for Go 1.18](https://tip.golang.org/doc/go1.18) for more details.

As always, especially for beta releases, if you notice any problems,
please [file an issue](/issue/new).

We hope you enjoy testing the beta,
and we hope you all have a restful remainder of 2021.
Happy holidays!


