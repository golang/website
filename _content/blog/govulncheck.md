---
title: Govulncheck v1.0.0 is released!
date: 2023-07-13
by:
- Julie Qiu, for the Go security team
summary: Version v1.0.0 of golang.org/x/vuln has been released, introducing a new API and other improvements.
---

We are excited to announce that govulncheck v1.0.0 has been released,
along with v1.0.0 of the API for integrating scanning into other tools!

Go's support for vulnerability management was [first announced](https://go.dev/blog/vuln) last September.
We have made several changes since then, culminating in today's release.

This post describes Go's updated vulnerability tooling, and how to get started
using it. We also recently published a
[security best practices guide](https://go.dev/security/best-practices)
to help you prioritize security in your Go projects.

## Govulncheck

[Govulncheck](https://golang.org/x/vuln/cmd/govulncheck)
is a command-line tool that helps Go users find known vulnerabilities in
their project dependencies.
The tool can analyze both codebases and binaries,
and it reduces noise by prioritizing vulnerabilities in functions that your
code is actually calling.

You can install the latest version of govulncheck using
[go install](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies):

```
go install golang.org/x/vuln/cmd/govulncheck@latest
```

Then, run govulncheck inside your module:
```
govulncheck ./...
```

See the [govulncheck tutorial](https://go.dev/doc/tutorial/govulncheck)
for additional information on how to get started with using the tool.

As of this release, there is now a stable API available,
which is described at [golang.org/x/vuln/scan](https://golang.org/x/vuln/scan).
This API provides the same functionality as the govulncheck command,
enabling developers to integrate security scanners and other tools with govulncheck.
As an example, see the
[osv-scanner integration with govulncheck](https://github.com/google/osv-scanner/blob/d93d6b73e90ae392fe2b1b64a33bda6976b65b2d/internal/sourceanalysis/go.go#L20).

## Database

Govulncheck is powered by the Go vulnerability database, [https://vuln.go.dev](https://vuln.go.dev),
which provides a comprehensive source of information about known vulnerabilities
in public Go modules.
You can browse the entries in the database at [pkg.go.dev/vuln](https://pkg.go.dev/vuln).

Since the initial release, we have updated the [database API](https://go.dev/security/vuln/database#api)
to improve performance and ensure long-term extensibility.
An experimental tool to generate your own vulnerability database index is
provided at [golang.org/x/vulndb/cmd/indexdb](https://golang.org/x/vulndb/cmd/indexdb).

If you are a Go package maintainer, we encourage you to
[contribute information](https://go.dev/s/vulndb-report-new)
about public vulnerabilities in your projects.

For more information about the Go vulnerability database,
see [go.dev/security/vuln/database](https://go.dev/security/vuln/database).

## Integrations

Vulnerability detection is now integrated into a suite of tools that are
already part of many Go developers' workflows.

Data from the Go vulnerability database can be browsed at
[pkg.go.dev/vuln](https://pkg.go.dev/vuln).
Vulnerability information is also surfaced on the search and package pages
of pkg.go.dev. For example,
[the versions page of golang.org/x/text/language](https://pkg.go.dev/golang.org/x/text/language?tab=versions)
shows vulnerabilities in older versions of the module.

You can also run govulncheck directly in your editor using the Go extension
for Visual Studio Code.
See [the tutorial](https://go.dev/doc/tutorial/govulncheck-ide) to get started.

Lastly, we know that many developers will want to run govulncheck as part
of their CI/CD systems.
As a starting point, we have provided a
[GitHub Action for govulncheck](https://github.com/marketplace/actions/golang-govulncheck-action)
for integration with your projects.

## Video Walkthrough

If you are interested in a demo of the integrations described above,
we presented a walkthrough of these tools at Google I/O this year, in our talk,
[Build more secure apps with Go and Google](https://www.youtube.com/watch?v=HSt6FhsPT8c&ab_channel=TheGoProgrammingLanguage).

## Feedback

As always, we welcome your feedback! See details on
[how to contribute and help us make improvements](https://go.dev/security/vuln/#feedback).

We hope you’ll find the latest release of Go’s support for vulnerability
management useful and work with us to build a more secure and reliable Go
ecosystem.
