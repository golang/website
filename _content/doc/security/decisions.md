---
title: Go Security Decisions
layout: article
breadcrumb: true
---

## Overview

This document includes decisions the Go Security team has made about
various commonly-reported issues. It mostly serves as a reference
for things we do not consider to be a vulnerability.

This list is not comprehensive.

## Vulnerabilities

### Remote Code Execution

A scenario which permits an attacker to execute code in a situation
where code execution is not expected is a PRIVATE-track vulnerability.
This supersedes all other decisions.

This decision does not cover functions which are expected to execute code.

Covered:

- A parse function which executes code on a malicious input.
- A malicious request causing the HTTP server to execute code.

Not covered:

- The `go test` command runs tests. Running `go test` on attacker-controlled
  tests is not within our threat model.

## Non-Vulnerabilities

### image, x/image: Large images

Parsing a large image can allocate a large amount of memory.
For example, a 65536x65536 32-bit color image requires 16MiB
to store uncompressed.

Many image compression formats can reduce a large, simple image
to a very small file size. Decoding the small file may allocate
a large amount of memory.

Users parsing untrusted images should verify the image size prior
to parsing, using a function such as
[image.DecodeConfig](https://pkg.go.dev/image#DecodeConfig).

We do not consider it to be a vulnerability for an image parsing
function to decode a large, well-compressed image.
