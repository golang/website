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
