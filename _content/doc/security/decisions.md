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

### net/http: Redirects

The `net/http` package's HTTP client handles redirects.
It implements security-relevant behavior in redirect handling.
For example, it strips the "Authorization" header when following
a redirect to a domain that is not a subdomain or exact match
for the initial request's domain.

Header stripping is a defense-in-depth measure,
avoiding the case where a misconfigured or compromised server
inadvertently forwards a client request containing sensitive headers
to an untrusted destination.
Failure to strip headers on redirect does not, by itself, permit
an attacker to acquire credentials passed in header.
However, it may be combined with other vulnerabilities
(for example, a server with an open redirect vulnerability)
to do so.

Changing the HTTP client's behavior runs a high risk of breaking
existing users who depend on the current behavior.
For example, the client's same origin policy currently permits subdomains
(a redirect from example.com to www.example.com will preserve headers),
while the WHATWG Fetch standard does not.
Aligning the client with the standard may be worthwhile,
but doing so in a security release is more likely to cause
pain to existing users than it is to address real vulnerabilities.

Since redirect sanitization is a defense-in-depth measure,
and making changes to it is risky,
we consider all aspects of the HTTP client redirects
to be out of scope for the security bug process.
