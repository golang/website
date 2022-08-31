---
title: Go Security
layout: article
---

## Overview

This page provides information on writing secure and reliable software in Go.

## Go Security

**[Go Security Policy](https://go.dev/security/policy)**

An explanation of how to report security issues in the Go standard library and
sub-repositories to the Go team.

**[Go Security Releases](https://go.dev/doc/devel/release)**

Release notes for past security problems. Per the [release
policy](https://go.dev/doc/devel/release#policy), we issue security fixes to
the two most recent major releases of Go.

## [Go Vulnerability Management](https://go.dev/security/vulndb)

The main documentation page for the Go vulnerability management system, which
helps developers find vulnerabilities in their Go projects.

## [Go Fuzzing](https://go.dev/doc/fuzz)

The main documentation page for Go native fuzzing.

Fuzzing is a type of automated testing which continuously manipulates inputs to
a program to find bugs. Go supports fuzzing in its standard toolchain beginning
in Go 1.18. Native Go fuzz tests are [supported by
OSS-Fuzz](https://google.github.io/oss-fuzz/getting-started/new-project-guide/go-lang/#native-go-fuzzing-support).

## Go Cryptography

The Go cryptography libraries are the [crypto/…](https://pkg.go.dev/crypto) and
[goland.org/x/crypto/…](https://pkg.go.dev/crypto) packages in the Go standard
library and subrepos, and developed following [these
principles](https://go.googlesource.com/proposal/+/master/design/cryptography-principles.md).
