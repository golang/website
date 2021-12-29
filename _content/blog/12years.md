---
title: Twelve Years of Go
date: 2021-11-10
by:
- Russ Cox, for the Go team
summary: Happy Birthday, Go!
---


Today we celebrate the twelfth birthday of the Go open source release.
We have had an eventful year and have a lot to look forward to next year.

The most visible change here on the blog is our
[new home on go.dev](/blog/tidy-web),
part of consolidating all our Go web sites into a single, coherent site.
Another part of that consolidation was
[replacing godoc.org with pkg.go.dev](/blog/godoc.org-redirect).

In February, the [Go 1.16 release](/blog/go1.16)
added [macOS ARM64 support](/blog/ports),
added [a file system interface](/pkg/io/fs) and [embedded files](/pkg/embed),
and [enabled modules by default](/blog/go116-module-changes),
along with the usual assortment of improvements and optimizations.

In August, the [Go 1.17 release](/blog/go1.17)
added Windows ARM64 support,
made [TLS cipher suite decisions easier and more secure](/blog/tls-cipher-suites),
introduced [pruned module graphs](/doc/go1.17#go-command)
to make modules even more efficient in large projects,
and added
[new, more readable build constraint syntax](https://pkg.go.dev/cmd/go#hdr-Build_constraints).
Under the hood, Go 1.17 also switched to a register-based calling convention for Go functions
on x86-64, improving performance in CPU-bound applications by 5–15%.

Over the course of the year, we
published [many new tutorials](/doc/tutorial/),
a [guide to databases in Go](/doc/database/),
a [guide to developing modules](/doc/#developing-modules),
and a [Go modules reference](/ref/mod).
One highlight is the new tutorial
“[Developing a RESTful API with Go and Gin](/doc/tutorial/web-service-gin)”,
which is also available in
[interactive form using Google Cloud Shell](/s/cloud-shell-web-tutorial).

We've been busy on the IDE side,
[enabling gopls by default in VS Code Go](/blog/gopls-vscode-go)
and delivering countless improvements to both `gopls` and VS Code Go,
including a [powerful debugging experience](https://github.com/golang/vscode-go/blob/master/docs/debugging.md)
powered by Delve.

We also launched the [Go fuzzing beta](/blog/fuzz-beta)
and [officially proposed adding generics to Go](/blog/generics-proposal),
both of which are now expected in Go 1.18.

Continuing to adapt to “virtual-first”, the Go team hosted our second annual
[Go day at Google Open Source Live](https://opensourcelive.withgoogle.com/events/go-day-2021).
You can watch the talks on YouTube:

- “[Using Generics in Go](https://www.youtube.com/watch?v=nr8EpUO9jhw)”,
  by Ian Lance Taylor, introduces generics and how to use them effectively.

- “[Modern Enterprise Applications](https://www.youtube.com/watch?v=5fgG1qZaV4w)”,
  by Steve Francia, shows how Go plays a role in enterprise modernization.

- “[Building Better Projects with the Go Editor](https://www.youtube.com/watch?v=jMyzsp2E_0U)”,
  by Suzy Mueller, demonstrates how VS Code Go's integrated tooling
  helps you navigate code, debug tests, and more.

- “[From Proof of Concept to Production](https://www.youtube.com/watch?v=e7PtBOsTpXE)”,
  by Benjamin Cane, a Distinguished Engineer at American Express,
  explains how American Express came to use Go for its payments and rewards platforms.

## Going Forward

We’re incredibly excited about what’s in store for Go’s 13th year.
Next month, we will have two talks at [GopherCon 2021](https://www.gophercon.com/),
along with [many talented speakers from across the Go community](https://www.gophercon.com/agenda).
Register for free and mark your calendars!

- “Why and How to Use Go Generics”,
  by Robert Griesemer and Ian Lance Taylor,
  who led the design and implementation of this new feature. \
  [Dec 8, 11:00 AM (US Eastern)](https://www.gophercon.com/agenda/session/593015).

- “Debugging Go Code Using the Debug Adapter Protocol (DAP)”,
  by Suzy Mueller,
  show how to use VS Code Go's advanced debugging features with Delve. \
  [Dec 9, 3:20 PM (US Eastern)](https://www.gophercon.com/agenda/session/593029).

In February, the Go 1.18 release will expand the new
register-based calling convention to non-x86 architectures,
bringing dramatic performance improvements with it.
It will include the new Go fuzzing support.
And it will be the first release to include support for generics.

Generics will be one of our focuses for 2022.
The initial release in Go 1.18 is only the beginning.
We need to spend time using generics and learning what works
and what doesn't, so that we can write best practices
and decide what should be added to the standard library
and other libraries.
We expect that Go 1.19 (expected in August 2022)
and later releases will further refine the design and implementation
of generics as well as integrating them further into the overall Go experience.

Another focus for 2022 is supply chain security.
We have been talking for years about the
[problems of dependencies](https://research.swtch.com/deps).
The design of Go modules provides
[reproducible, verifiable, verified builds](https://research.swtch.com/vgo-repro),
but there is still more work to be done.
Starting in Go 1.18, the `go` command will embed more information in binaries
about their build configurations, both to make reproducibility easier
and to help projects that need to
[generate an SBOM](https://en.wikipedia.org/wiki/Software_bill_of_materials) for Go binaries.
We have also started work on a
[Go vulnerability database](https://pkg.go.dev/golang.org/x/vuln)
and an associated tool to report vulnerabilities in a program's dependencies.
One of our goals in this work is to significantly improve the signal-to-noise ratio
of this kind of tool:
if a program doesn't use the vulnerable function, we don't want to report that.
Over the course of 2022 we plan to make this available as a standalone tool
but also to add it to existing tooling, including `gopls` and VS Code Go, and to [pkg.go.dev](https://pkg.go.dev).
There is also more to do to improve other aspects of Go's supply chain security posture.
Stay tuned for details.

Overall, we expect 2022 to be an eventful year for Go,
and we will continue to deliver the timely releases and improvements
you've come to expect.

## Thank You!

Go is far more than just us on the Go team at Google.
Thank you for your help making Go a success
and joining us on this adventure.
We hope you are all staying safe and wish you all the best.

