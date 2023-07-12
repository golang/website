---
title: Security Best Practices for Go Developers
layout: article
---

[Back to Go Security](/security)

This page provides Go developers with best practices for prioritizing the
security of their projects. From automating testing with fuzzing to easily
checking for race conditions, these tips can help make your codebase more
secure and reliable.

## Scan source code and binaries for vulnerabilities

Regularly scanning your code and binaries for vulnerabilities helps identify
potential security risks early.
You can use [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck),
backed by the [Go vulnerability database](https://pkg.go.dev),
to scan your code for vulnerabilities and analyze which ones actually affect you.
Get started with [the govulncheck tutorial](https://go.dev/doc/tutorial/govulncheck).

Govulncheck can also be integrated into CI/CD flows.
The Go team provides a
[GitHub Action for govulncheck](https://github.com/marketplace/actions/golang-govulncheck-action)
on the GitHub Marketplace.
Govulncheck also supports a `-json` flag to help developers integrate vulnerability
scanning with other CI/CD systems.

You can also scan for vulnerabilities directly in your code editor by using
the [Go extension for Visual Studio Code](https://go.dev/security/vuln/editor).
Get started with [this tutorial](https://go.dev/doc/tutorial/govulncheck-ide).

## Keep your Go version and dependencies up to date

Keeping your [Go version up-to-date](https://go.dev/doc/install) offers
access to the latest language features,
performance improvements and patches for known security vulnerabilities.
An updated Go version also ensures compatibility with newer versions of dependencies,
helping to avoid potential integration issues.
Review the [Go release history](https://go.dev/doc/devel/release) to see
what changes have been made to Go between releases.
The Go team issues point releases throughout the release cycle to address security bugs.
Be sure to update to the latest minor Go version to ensure you have the
latest security fixes.

Maintaining up-to-date third-party dependencies is also crucial for software security,
performance, and compliance with the latest standards in the Go ecosystem.
However, updating to the latest versions without thorough review
[can also be risky](https://research.swtch.com/npm-colors),
potentially introducing new bugs, incompatible changes,
or even malicious code.
Therefore, while it's essential to update dependencies for the latest security
patches and improvements,
each update should be carefully reviewed and tested.

## Test with fuzzing to uncover edge-case exploits

[Fuzzing](https://go.dev/security/fuzz) is a type of automated testing that
uses coverage guidance to manipulate random inputs and walk through code
to find and report potential vulnerabilities like SQL injections,
buffer overflows, denial or service and cross-site scripting attacks.
Fuzzing can often reach edge cases that programmers miss,
or deem too improbable to test.
Get started with [this tutorial](https://go.dev/doc/tutorial/fuzz).

## Check for race conditions with Go’s race detector

Race conditions occur when two or more [goroutines](https://go.dev/tour/concurrency/1)
access the same resource concurrently,
and at least one of those accesses is a write.
This can lead to unpredictable, difficult-to-diagnose issues in your software.
Identify potential race conditions in your Go code using the built-in
[race detector](https://go.dev/doc/articles/race_detector),
which can help you ensure the safety and reliability of your concurrent programs.
The race detector finds races that occur at runtime,
however, so it will not find races in code paths that are not executed.

To use the race detector, add the `-race` flag when running your tests or
building your application,
for example, `go test -race`.
This will compile your code with the race detector enabled and report any
race conditions it detects at runtime.
When the race detector finds a data race in the program, it will
[print a report](https://go.dev/doc/articles/race_detector#report-format)
containing stack traces for conflicting accesses,
and stacks where the involved goroutines were created.

## Use Vet to examine suspicious constructs

Go’s [vet command](https://pkg.go.dev/cmd/vet) is designed to analyze
your source code and flag potential issues that might not necessarily be syntax errors,
but could lead to problems during runtime.
These include suspicious constructs, such as unreachable code,
unused variables, and common mistakes around goroutines.
By catching these issues early in the development process,
go vet helps maintain code quality, reduces debugging time,
and enhances overall software reliability.
To run go vet for a specified project, run:

```
go vet ./...
```

## Subscribe to golang-announce for notification of security releases

Go releases containing security fixes are pre-announced to the low-volume
mailing list [golang-announce@googlegroups.com](https://groups.google.com/group/golang-announce).
If you want to know when security fixes to Go itself are on the way, subscribe.
