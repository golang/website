---
title: Fourteen Years of Go
date: 2023-11-10
by:
- Russ Cox, for the Go team
summary: Happy Birthday, Go!
---

<img src="/doc/gopher/gopherdrink.png" height="219" width="223" align="right" style="margin: 0 0 1em 1em;">

Today we celebrate the fourteenth birthday of the Go open source release!
Go has had a great year, with two feature-filled releases and other important milestones.

We released [Go 1.20 in February](/blog/go1.20)
and [Go 1.21 in August](/blog/go1.21),
focusing more on implementation improvements
than new language changes.

Profile-guided optimization (PGO),
[previewed in Go 1.20](/blog/pgo-preview)
and
[released in Go 1.21](/blog/pgo),
allows the Go compiler to read a profile of your program
and then spend more time optimizing the parts
of your program that run most often.
In Go 1.21, workloads typically get between
2% and 7% CPU usage improvements from enabling PGO.
See “[Profile-guided optimization in Go 1.21](/blog/pgo)” for an overview
and the [profile-guided optimization user guide](/doc/pgo)
for complete documentation.

Go has provided support for gathering coverage profiles during `go test`
[since Go 1.2](/blog/cover).
Go 1.20 added support for gathering coverage profiles in binaries
built by `go build`,
allowing you to gather coverage during larger integration tests as well.
See “[Code coverage for Go integration tests](/blog/integration-test-coverage)” for details.

Compatibility has been an important part of Go since
“[Go 1 and the Future of Go Programs](/doc/go1compat)”.
Go 1.21 improved compatibility further
by expanding the conventions for use of GODEBUG
in situations where we need to make a change,
such as an important bug fix,
that must be permitted but may still break existing programs.
See the blog post
“[Backward Compatibility, Go 1.21, and Go 2](/blog/compat)”
for an overview and
the documentation
“[Go, Backwards Compatibility, and GODEBUG](/doc/godebug)” for details.

Go 1.21 also shipped support for built-in toolchain management,
allowing you to change which version of the
Go toolchain you use in a specific module
as easily as you change the versions of other dependencies.
See the blog post
“[Forward Compatibility and Toolchain Management in Go 1.21](/blog/toolchain)”
for an overview and the documentation
“[Go Toolchains](/doc/toolchain)”
for details.

Another important tooling achievement was the
integration of on-disk indexes into
gopls, the Go LSP server.
This cut gopls's startup latency and memory usage by 3-5X
in typical use cases.
“[Scaling gopls for the growing Go ecosystem](/blog/gopls-scalability)”
explains the technical details.
You can make sure you're running the latest gopls by running:

```
go install golang.org/x/tools/gopls@latest
```

Go 1.21 introduced new
[cmp](/pkg/cmp/),
[maps](/pkg/maps/),
and
[slices](/pkg/slices/)
packages — Go’s first generic standard libraries —
as well as expanding the set of comparable types.
For details about that, see the blog post
“[All your comparable types](/blog/comparable)”.

Overall, we continue to refine generics
and to write talks and blog posts explaining
important details.
Two notable posts this year were
“[Deconstructing Type Parameters](/blog/deconstructing-type-parameters)”,
and
“[Everything You Always Wanted to Know About Type Inference – And a Little Bit More](/blog/type-inference)”.

Another important new package in Go 1.21 is
[log/slog](/pkg/log/slog/),
which adds an official API for
structured logging to the standard library.
See “[Structured logging with slog](/blog/slog)” for an overview.

For the WebAssembly (Wasm) port, Go 1.21 shipped support
for running on WebAssembly System Interface (WASI) preview 1.
WASI preview 1 is a new “operating system” interface for Wasm
that is supported by most server-side Wasm environments.
See “[WASI support in Go](/blog/wasi)” for a walkthrough.

On the security side, we are continuing to make sure
Go leads the way in helping developers understand their
dependencies and vulnerabilities,
with [Govulncheck 1.0 launching in July](/blog/govulncheck).
If you use VS Code, you can run govulncheck directly in your
editor using the Go extension:
see [this tutorial](/doc/tutorial/govulncheck-ide) to get started.
And if you use GitHub, you can run govulncheck as part of
your CI/CD, with the
[GitHub Action for govulncheck](https://github.com/marketplace/actions/golang-govulncheck-action).
For more about checking your dependencies for vulnerability problems,
see this year's Google I/O talk,
“[Build more secure apps with Go and Google](https://www.youtube.com/watch?v=HSt6FhsPT8c&ab_channel=TheGoProgrammingLanguage)”.)

Another important security milestone was
Go 1.21's highly reproducible toolchain builds.
See “[Perfectly Reproducible, Verified Go Toolchains](/blog/rebuild)” for details,
including a demonstration of reproducing an Ubuntu Linux Go toolchain
on a Mac without using any Linux tools at all.

It has been a busy year!

In Go's 15th year, we'll keep working to make Go the best environment
for software engineering at scale.
One change we're particularly excited about is
redefining for loop `:=` semantics to remove the
potential for accidental aliasing bugs.
See “[Fixing For Loops in Go 1.22](/blog/loopvar-preview)”
for details,
including instructions for previewing this change in Go 1.21.

## Thank You!

The Go project has always been far more than just us on the Go team at Google.
Thank you to all our contributors and everyone in the Go community for
making Go what it is today.
We wish you all the best in the year ahead.

