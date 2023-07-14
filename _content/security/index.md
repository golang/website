---
title: Security
layout: article
---

This page provides resources for Go developers to improve security for their
projects.

(See also: [Security Best Practices for Go Developers](https://go.dev/security/best-practices).)

## Find and fix known vulnerabilities

Go’s vulnerability detection aims to provide low-noise, reliable tools for
developers to learn about known vulnerabilities that may affect their projects.
For an overview, start at [this summary and FAQ page](https://go.dev/security/vuln)
about Go’s vulnerability management architecture. For an applied approach,
explore the tools below.

### Scan code for vulnerabilities with govulncheck

Developers can use the govulncheck tool to determine whether any known
vulnerabilities affect their code and prioritize next steps based on which vulnerable
functions and methods are actually called.

- [View the govulncheck documentation](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
- [Tutorial: Get started with govulncheck](https://go.dev/doc/tutorial/govulncheck)

### Detect vulnerabilities from your editor

The VS Code Go extension checks third-party dependencies and surfaces relevant vulnerabilities.

- [User documentation](https://go.dev/security/vuln/editor)
- [Download VS Code Go](https://marketplace.visualstudio.com/items?itemName=golang.go)
- [Tutorial: Get started with VS Code Go](https://go.dev/doc/tutorial/govulncheck-ide)

### Find Go modules to build upon

[Pkg.go.dev](https://pkg.go.dev/) is a website for discovering, evaluating and
learning more about Go packages and modules. When discovering and evaluating
packages on pkg.go.dev, you will
[see a banner on the top of a page](https://pkg.go.dev/golang.org/x/text@v0.3.7/language)
if there are vulnerabilities in that version. Additionally, you can see the
[vulnerabilities impacting each version of a package](https://pkg.go.dev/golang.org/x/text@v0.3.7/language?tab=versions)
on the version history page.

### Browse the vulnerability database

The Go vulnerability database collects data directly from Go package
maintainers as well as from outside sources such as [MITRE](https://www.cve.org/) and [GitHub](https://github.com/). Reports
are curated by the Go Security team.

- [Browse reports in the Go vulnerability database](https://pkg.go.dev/vuln/)
- [View the Go Vulnerability Database documentation](https://go.dev/security/vuln/database)
- [Contribute a public vulnerability to the database](https://go.dev/s/vulndb-report-new)


## Report security bugs in the Go project

### [Security Policy](https://go.dev/security/policy)

Consult the Security Policy for instructions on how to
[report a vulnerability in the Go project](https://go.dev/security/policy#reporting-a-security-bug).
The page also details the Go security team’s process of tracking issues and
disclosing them to the public. See the
[release history](https://go.dev/doc/devel/release) for details about past security
fixes. Per the [release policy](https://go.dev/doc/devel/release#policy),
we issue security fixes to the two most recent major releases of Go.

## Test unexpected inputs with fuzzing

Go native fuzzing provides a type of automated testing which continuously
manipulates inputs to a program to find bugs. Go supports fuzzing in its
standard toolchain beginning in Go 1.18.  Native Go fuzz tests are
[supported by OSS-Fuzz](https://google.github.io/oss-fuzz/getting-started/new-project-guide/go-lang/#native-go-fuzzing-support).

- [Review the basics of fuzzing](https://go.dev/security/fuzz)
- [Tutorial: Get started with fuzzing](https://go.dev/doc/tutorial/fuzz)

## Secure services with Go's cryptography libraries

Go’s cryptography libraries aim to help developers build secure applications.
See documentation for the [crypto packages](https://pkg.go.dev/golang.org/x/crypto)
and [golang.org/x/crypto/](https://pkg.go.dev/golang.org/x/crypto).
