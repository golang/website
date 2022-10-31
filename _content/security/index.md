---
title: Go Security
layout: article
---

## Overview

This page provides information on writing secure and reliable software in Go.

## Go Security

### Go Security Policy

The [Go Security Policy](/security/policy) explains how to report security
issues in the Go standard library and sub-repositories to the Go team.

### Go Security Releases

The [Go Release History](/doc/devel/release) includes release notes for past
security problems. Per the [release policy](/doc/devel/release#policy), we
issue security fixes to the two most recent major releases of Go.

## Go Vulnerability Management

[Go's vulnerability management](/security/vuln) support helps developers find
known public vulnerabilities that may affect their Go projects.

## Go Fuzzing

[Go native fuzzing](/security/fuzz) provides a type of automated testing which
continuously manipulates inputs to a program to find bugs.

Go supports fuzzing in its standard toolchain beginning in Go 1.18.
Native Go fuzz tests are
[supported by OSS-Fuzz](https://google.github.io/oss-fuzz/getting-started/new-project-guide/go-lang/#native-go-fuzzing-support).
Try out [the tutorial for fuzzing with Go](/doc/tutorial/fuzz).

## Go Cryptography

The Go cryptography libraries are the [crypto/…](https://pkg.go.dev/crypto)
and [golang.org/x/crypto/…](https://pkg.go.dev/golang.org/x/crypto) packages
in the Go standard library and subrepos,
and developed following [these principles](https://go.googlesource.com/proposal/+/master/design/cryptography-principles.md).
