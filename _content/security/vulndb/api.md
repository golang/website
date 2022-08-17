---
title: Go Vulnerability Database API
layout: article
---

## Overview

The Go vulnerability database ([https://vuln.go.dev](https://vuln.go.dev))
serves Go vulnerability information in the
[Open Source Vulnerability (OSV)](https://ossf.github.io/osv-schema/) format.
We recommend using
[client.Client](https://pkg.go.dev/golang.org/x/vuln/client#Client) to read
data from the Go vulnerability database.

Do not rely on the contents of the x/vulndb repository. The YAML files in that
repository are maintained using an internal format that is intended to change
without warning.

## API

The endpoints in the table below are supported. For each path:

- `$base` is the path portion of a Go vulnerability database URL ([https://vuln.go.dev](https://vuln.go.dev)).
- `$module` is a module path
- `$vuln` is a Go vulnerability ID (for example, `GO-2021-1234`)

<table>
  <tr>
    <td>Path</td>
    <td>Description</td>
  </tr>
  <tr>
    <td>`$base/index.json`</td>
    <td>List of module paths in the database mapped to its last modified timestamp (<a href="https://vuln.go.dev/index.json">link</a>).</td>
  </tr>
  <tr>
    <td>`$base/$module.json`</td>
    <td>List of vulnerability entries for that module ([example](https://vuln.go.dev/golang.org/x/crypto.json).</td>
  </tr>
  <tr>
    <td>`$base/ID/index.json`</td>
    <td>List of all the vulnerability entries in the database. </td>
  </tr>
  <tr>
    <td>`$base/ID/$vuln.json`</td>
    <td>An individual Go vulnerability report.</td>
  </tr>
</table>

Note that this project is under active development, and it is possible for
these endpoints to change.

## Schema

Reports are written following the
[Open Source Vulnerability (OSV)](https://ossf.github.io/osv-schema/) format.
The fields below are specific to the Go vulnerability database:

### id

The id field is a unique identifier for the vulnerability entry. It is a string
of the format GO-&lt;YEAR>-&lt;ENTRYID>.

### affected

The [affected](https://ossf.github.io/osv-schema/#affected-fields) field is a
JSON array containing objects that describes the module versions that contain
the vulnerability.

#### affected[].package

The
[affected[].package](https://ossf.github.io/osv-schema/#affectedpackage-field)
field is a JSON object identifying the affected _module._ The object has two
required fields:

- **ecosystem**: this will always be "Go"
- **name**: this is the Go module path
  - Importable packages in the standard library will have the name _stdlib_.
  - The go command will have the name _toolchain_.

**affected[].ecosystem_specific**

The
[affected[].ecosystem_specific](https://ossf.github.io/osv-schema/#affectedecosystem_specific-field)
field is a JSON object with additional information about the vulnerability,
which is used by [package
vulncheck](https://pkg.go.dev/golang.org/x/vuln/vulncheck).

For now, ecosystem specific will always be an object with a single field,
`imports`.

#### affected[].ecosystem_specific.imports

The `affected[].ecosystem_specific.imports` field is a JSON array containing
the packages and symbols affected by the vulnerability. Each object in the
array will have these two fields:

- **path:** a string with the import path of the package containing the vulnerability
- **symbols:** a string array with the names of the symbols (function or method) that contains the vulnerability

For information on other fields in the schema, refer to the [OSV spec](https://ossf.github.io/osv-schema).

## Examples

TODO: add an example once changes have been made in x/vulndb and x/vuln.

## Excluded Reports

The reports in the Go vulnerability database are collected from different
sources and curated by the Go Security team. We may come across a vulnerability
(for example, a CVE or GHSA) and choose to exclude it for a variety of reasons.
In these cases, a minimal report will be created in the x/vulndb repository,
under
[x/vulndb/data/excluded](https://github.com/golang/vulndb/tree/master/data/excluded).

Reports may be excluded for these reasons:

- `NOT_GO_CODE`: The vulnerability is not in a Go package, and cannot affect any
  Go packages. (For example, a vulnerability in  a C++ library.)
- `NOT_IMPORTABLE`: The vulnerability occurs in package `main`, an `internal/`
  package only imported by package `main`, or some  other location which can
  never be imported by another module.
- `EFFECTIVELY_PRIVATE`: While the vulnerability occurs in a Go package which
  can be imported by another module, the package is not intended for external
  use and is not likely to ever be imported outside the module in which it is
  defined.
- `DEPENDENT_VULNERABILITY`: This vulnerability is a subset of another
  vulnerability in the database. For example, if package A contains a
  vulnerability, package B depends on package A, and there are separate CVEs
  for packages A and B, we might mark the report for B as a dependent
  vulnerability entirely superseded by the report for A.
- `NOT_A_VULNERABILITY`: While a CVE or GHSA has been assigned, there is no
  known vulnerability associated with it.

At the moment, excluded reports are not served via
[vuln.go.dev](https://vuln.go.dev) API.  excluded reports. However, if you have
a specific use case and it would be helpful to have access to this information
through the API,
[please let us know](https://golang.org/s/govulncheck-feedback).
