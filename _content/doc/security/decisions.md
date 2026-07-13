---
title: Go Security Decisions
layout: article
breadcrumb: true
---

## Overview {#overview}

This document includes decisions the Go Security team has made about
various commonly-reported issues. It mostly serves as a reference
for things we do not consider to be a vulnerability.

This list is not comprehensive.

## Vulnerabilities {#vulns}

### Remote Code Execution {#rce}

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

### Panics {#panics}

A panic when processing attacker-controlled input may be a vulnerability.

A panic in a server expected to handle attacker-controlled requests,
such as a `net/http` server, is usually a PRIVATE-track vulnerability.

A panic in a client, such as a `net/http` client, is usually a PUBLIC-track
vulnerability.

A panic in a parse function in a package which handles plausibly
malicious input, such as `archive/zip` or `image/png`,
is usually a PUBLIC-track vulnerablity.

An input which causes a panic due to its natural size,
such as a very large image, is usually not classified as a vulnerability.

Invalid inputs to functions which are not clearly designed to parse
potentially malicious data are not in our threat model and
generally out of scope as security bugs. For example, image *parsers*
are expected to defend against invalid inputs, but a panic in
an image *encoder* might be a bug but would not be handled
as a vulnerability.

### Excessive resource consumption {#quadratic}

We generally treat excessive CPU or memory consumption,
such as a function with a runtime that is O(n²) in terms of its input size,
as equivalent to a panic.

### Building malicious code {#malicious-build}

Building an attacker-controlled program should be safe, even when
running it is not. For example, we intend it to be safe to build
untrusted code which is then run in a sandbox.
We strongly recommend that anyone building untrusted code consider
this a defense-in-depth measure, and that builds of untrusted code
also be performed in an unprivileged sandbox environment.

Data exfiltation is not in our threat model.
Building an attacker-controlled program may produce output which
contains the contents of arbitrary local files.

*Executing* malicious code is also not in our threat model.
Running a malicious program obviously permits it to do malicious things.

A vulnerability, usually PRIVATE track:

- A malicious module can cause "go build" to execute arbitrary code.

Not handled as a security issue:

- "go build" produces an error containing text from an arbitrary local file.
- A compiled executable contains the contents from an arbitrary local file.
- A malicious program can corrupt the runtime's state and execute arbitrary code.
- A compilation bug allows obfuscating malicious code.

## Non-Vulnerabilities {#non-vuln}

### Attacker-controlled environment {#attacker-control}

If an attack relies on the attacker having control over the environment
a program runs in, it is not a vulnerability.

This includes, but is not restricted to, an attacker with the ability to
add programs to `$PATH` or set arbitrary environment variables.

### image, x/image: Large images {#large-image}

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

### net/http: Redirects {#http-redirect}

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
