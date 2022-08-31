---
title: Go Vulnerability Management
layout: article
---

## Overview

This page describes the Go vulnerability management system.

_This project is a work in progress._

## Architecture

<div class="image">
  <center>
    <img src="architecture.svg" alt="Go Vulnerability Management Architecture"></img>
  </center>
</div>

The Go vulnerability management system consists of the following high-level
pieces:

1. A **data pipeline** that populates the vulnerability database. Data about
   new vulnerabilities come directly from Go package maintainers or sources such as
   MITRE and GitHub. Reports are curated by the Go Security team.

2. A **vulnerability database** that stores all information presented by
   govulncheck and can be consumed by other clients.

3. A **client library**
   ([golang.org/x/vuln/client](https://pkg.go.dev/golang.org/x/vuln/client)), which reads data
   from the Go vulnerability database. This is also used by pkg.go.dev to surface
   vulnerabilities.

4. A **vulncheck API**
   ([golang.org/x/vuln/vulncheck](https://pkg.go.dev/golang.org/x/vuln/vulncheck)), which is
   used to find vulnerabilities affecting Go packages and perform static analysis.
   This API is made available for clients that do not want to run the govulncheck
   binary, such as VS Code Go.

5. The **govulncheck command**
   ([golang.org/x/vuln/cmd/govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck),
   a wrapper around the vulncheck library for use on the command line.

6. A **web portal** that presents information about vulnerabilities, hosted at
   [pkg.go.dev/vuln](https://pkg.go.dev/vuln).


## References

### [Go Vulnerability Database API](https://go.dev/security/vuln/database)

Documentation on the Go vulnerability database API.

### [Vulnerability Detection For Go](https://go.dev/security/vulncheck)

An explanation of the features of vulncheck. Reference documentation is
at
[pkg.go.dev/golang.org/x/vuln/vulncheck](https://pkg.go.dev/golang.org/x/vuln/vulncheck)

### [Command govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)

Documentation on the CLI tool govulncheck.

### [Go CNA Policy](https://go.dev/security/vuln/cna)

Documentation on the Go CNA policy.
