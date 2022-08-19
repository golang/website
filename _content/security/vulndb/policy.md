---
title: Go CNA Policy
layout: article
---

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

To report potential new vulnerabilities in the Go project, refer to
[go.dev/security/policy](https://go.dev/security/policy).

## Requesting a CVE ID for a public vulnerability

**IMPORTANT**: The form linked below creates a public issue on the issue tracker, and therefore
*must not* be used to report undisclosed vulnerabilites in Go (see our
[security policy](https://go.dev/security/policy) for instructions on reporting
undisclosed issues).

To request a CVE ID for an existing PUBLIC vulnerability in the Go ecosystem,
[submit a request via this form](https://github.com/golang/vulndb/issues/new?assignees=&labels=Needs+Triage%2CDirect+External+Report&template=new_third_party_vuln.yml&title=x%2Fvulndb%3A+potential+Go+vuln+in+%3Cpackage%3E).

A vulnerability is considered public if it has already been disclosed publicly, or it exists in a
package you maintain (and you are ready to disclose it publicly).

This is a new feature and is still in development; please give any feedback
regarding the process to security@golang.org.

## Contact

For more information, email security@golang.org.
