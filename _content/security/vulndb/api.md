---
title: Go Vulnerability Database API
layout: article
---

## Overview

The Go vulnerability database is rooted at `https://vuln.go.dev` and
provides data as JSON. We recommend using
[client.Client](https://pkg.go.dev/golang.org/x/vuln/client#Client) to read
data from the Go vulnerability database.

Do not rely on the contents of the x/vulndb repository. The YAML files in that
repository are maintained using an internal format that is subject to change
without warning.

## API

The endpoints in the table below are supported. For each path:

- `$base` is the path portion of a Go vulnerability database URL (`https://vuln.go.dev`).
- `$module` is a module path
- `$vuln` is a Go vulnerability ID (for example, `GO-2021-1234`)

<table>
  <thead>
    <tr>
      <th>Path</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>$base/index.json</code></td>
      <td>
        List of module paths in the database mapped to its last modified
        timestamp (<a href="https://vuln.go.dev/index.json">link</a>).
      </td>
    </tr>
    <tr>
      <td><code>$base/$module.json</code></td>
      <td>
        List of vulnerability entries for that module (<a href="https://vuln.go.dev/golang.org/x/crypto.json">example</a>).
      </td>
    </tr>
    <tr>
      <td><code>$base/ID/index.json</code></td>
      <td>
        List of all the vulnerability entries in the database.
      </td>
    </tr>
    <tr>
      <td><code>$base/ID/$vuln.json</code></td>
      <td>
        An individual Go vulnerability report.
      </td>
    </tr>
  </tbody>
</table>

Note that these paths and format are provisional and likely to change until an
approved proposal.
