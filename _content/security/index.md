---
title: Go Security
layout: article
---

## Overview

This page provides information on writing secure and reliable software in Go.

## Go Security

### [Go Security Policy](https://go.dev/security/policy)

An explanation of how to report security issues in the Go standard library and
sub-repositories to the Go team.

### [Go Security Releases](https://go.dev/doc/devel/release)

Release notes for past security problems. Per the
[release policy](https://go.dev/doc/devel/release#policy), the two most recent
major Go releases are supported.

## Go Vulnerability Management

_This project is a work in progress._

### [Go Vulnerability Management](https://go.dev/security/vulndb)

The main documentation page for the Go vulnerability management system.

### Go Vulnerability Database

A list of vulnerabilities in the Go vulnerability database can be found at
[vuln.go.dev/index.json](https://vuln.go.dev/index.json).
[See protocol documentation](https://go.dev/security/vulndb/api) for more
information.

### Vulnerability Detection For Go

An overview of the Go vulnerability detection package,
[golang.org/x/vuln/vulncheck](https://golang.org/x/vuln/vulncheck), which
enables Go developers to scan dependencies in their Go projects for public
vulnerabilities.

## Go Fuzzing

### [Go Fuzzing](https://go.dev/doc/fuzz)

The main documentation page for Go native fuzzing.

### [Tutorial: Getting started with fuzzing](https://go.dev/doc/tutorial/fuzz)

Tutorial introducing the basics of fuzzing in Go.

### [Integrating with OSS-Fuzz](https://google.github.io/oss-fuzz/getting-started/new-project-guide/go-lang/#native-go-fuzzing-support)

Documentation on running native Go fuzz tests with OSS-Fuzz.

## Go Cryptography

### Cryptography libraries

The Go cryptography libraries are the [crypto/...](https://pkg.go.dev/crypto)
and [golang.org/x/crypto/...](https://golang.org/x/crypto) packages in the Go
standard library and subrepos.

### [Cryptography Principles](https://go.googlesource.com/proposal/+/master/design/cryptography-principles.md)

Goals and principles for the Go cryptography libraries.
