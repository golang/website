---
title: Go CNA
layout: article
---

**This page is a work in progress.**

## Overview
The Go CNA is a
[CVE Numbering Authority](https://www.cve.org/ProgramOrganization/CNAs), which issues
[CVE IDs](https://www.cve.org/ResourcesSupport/Glossary?activeTerm=glossaryCVEID) and publishes
[CVE Records](https://www.cve.org/ResourcesSupport/Glossary?activeTerm=glossaryRecord)
for public vulnerabilities in the Go ecosystem. It is a sub-CNA of the Google CNA.

## Scope
The Go CNA covers vulnerabilities in the Go project (the Go
[standard library](https://pkg.go.dev/std) and
[sub-repositories](https://pkg.go.dev/golang.org/x)) and public vulnerabilities
in importable Go modules that are not already covered by another CNA.

This scope is intended to explicitly exclude vulnerabilities in applications or
packages written in Go that are not importable (for example, anything in
package `main` or an `internal/` directory).

To report vulnerabilities in the Go project, refer to
[go.dev/security/policy](https://go.dev/security/policy).

## Requesting a CVE

TODO: add instructions

## Contact

For more information, email security@golang.org.
