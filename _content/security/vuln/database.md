---
title: Go Vulnerability Database
layout: article
---

[Back to Go Vulnerability Management](/security/vuln)

## Overview

The Go vulnerability database ([https://vuln.go.dev](https://vuln.go.dev))
serves Go vulnerability information in the
[Open Source Vulnerability (OSV) schema](https://ossf.github.io/osv-schema/).

We recommend using
[golang.org/x/vuln/client](https://pkg.go.dev/golang.org/x/vuln/client) to read
data from the Go vulnerability database. You can also browse vulnerabilities
in the database at [pkg.go.dev/vuln](https://pkg.go.dev/vuln).

**Do not** rely on the contents of the x/vulndb Git repository. The YAML files in that
repository are maintained using an internal format that may change
without warning.

## Contributing

We would love for all Go package maintainers to [contribute](https://go.dev/s/vulndb-report-new)
information about public vulnerabilities in their own projects,
and [update](https://go.dev/s/vulndb-report-feedback) existing information about vulnerabilities
in their Go packages.

We aim to make reporting a low friction process,
so feel free to [send us your suggestions](https://go.dev/s/vuln-feedback).

Please **do not** use the forms above to report a vulnerability in the Go
standard library or sub-repositories.
Instead, follow the process at [go.dev/security/policy](/security/policy)
for vulnerabilities about the Go project.

## API

The vulnerability database is an HTTP server that can respond to GET requests
for paths specified below.
The requests have no query parameters, and no specific headers are required,
so even a site serving from a fixed file system (including a file:// URL)
can be a vulnerability database.

By default, govulncheck uses the Go vulnerability database at [vuln.go.dev](https://vuln.go.dev).
The command can be configured to contact vulnerability databases using the
GOVULNDB environment variable,
which accepts a list of a comma-separated list of vulnerability database URLs,
with http://, https://, or
file:// protocols. Entries from multiple databases are merged.

The vulnerability database supports the endpoints listed below. For each path:

- `$base` is the path portion of a Go vulnerability database URL, by default, "https://vuln.go.dev".
- `$module` is a module path
- `$vuln` is a Go vulnerability ID (for example, GO-2021-1234)

To avoid various character set issues, the `$module` element is encoded
using [module.EncodePath](https://pkg.go.dev/golang.org/x/mod/module/#EncodePath).

The endpoints are:

- `$base/index.json`
   List of module paths in the database mapped to its last modified timestamp ([link](https://vuln.go.dev/index.json)).

- `$base/$module.json`
   List of vulnerability entries for that module ([example](https://vuln.go.dev/golang.org/x/crypto.json)).

- `$base/ID/$vuln.json`
   An individual Go vulnerability report ([example](https://vuln.go.dev/ID/GO-2022-0191.json)).

## Schema

Reports use the
[Open Source Vulnerability (OSV) schema](https://ossf.github.io/osv-schema/).
The Go vulnerability database assigns the following meanings to the fields:

### id

The id field is a unique identifier for the vulnerability entry. It is a string
of the format GO-\<YEAR>-\<ENTRYID>.

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

#### affected[].ecosystem_specific

The
[affected[].ecosystem_specific](https://ossf.github.io/osv-schema/#affectedecosystem_specific-field)
field is a JSON object with additional information about the vulnerability,
which is used by [package
vulncheck](https://pkg.go.dev/golang.org/x/vuln/vulncheck).

For now, ecosystem specific will always be an object with a single field,
`imports`.

##### affected[].ecosystem_specific.imports

The `affected[].ecosystem_specific.imports` field is a JSON array containing
the packages and symbols affected by the vulnerability. Each object in the
array will have these two fields:

- **path:** a string with the import path of the package containing the vulnerability
- **symbols:** a string array with the names of the symbols (function or method) that contains the vulnerability
- **goos**: a string array with the execution operating system where the symbols appear, if known
- **goarch**: a string array with the architecture where the symbols appear, if known

For information on other fields in the schema, refer to the [OSV spec](https://ossf.github.io/osv-schema).

## Examples

All vulnerabilities in the Go vulnerability database use the OSV schema
described above.

See the links below for examples of different Go vulnerabilities:

- **Go standard library vulnerability** (GO-2022-0191):
  [JSON](https://vuln.go.dev/ID/GO-2022-0191.json),
  [HTML](https://pkg.go.dev/vuln/GO-2022-0191)
- **Go toolchain vulnerability** (GO-2022-0189):
  [JSON](https://vuln.go.dev/ID/GO-2022-0189.json),
  [HTML](https://pkg.go.dev/vuln/GO-2022-0189)
- **Vulnerability in Go module** (GO-2020-0015):
  [JSON](https://vuln.go.dev/ID/GO-2020-0015.json),
  [HTML](https://pkg.go.dev/vuln/GO-2020-0015)

## Excluded Reports

The reports in the Go vulnerability database are collected from different
sources and curated by the Go Security team. We may come across a vulnerability advisory
(for example, a CVE or GHSA) and choose to exclude it for a variety of reasons.
In these cases, a minimal report will be created in the x/vulndb repository,
under
[x/vulndb/data/excluded](https://github.com/golang/vulndb/tree/master/data/excluded).

Reports may be excluded for these reasons:

- `NOT_GO_CODE`: The vulnerability is not in a Go package,
  but it was marked as a security advisory for the Go ecosystem by another source.
  This vulnerability cannot affect any
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
  vulnerability, package B depends on package A, and there are separate CVE IDs
  for packages A and B, we might mark the report for B as a dependent
  vulnerability entirely superseded by the report for A.
- `NOT_A_VULNERABILITY`: While a CVE ID or GHSA has been assigned, there is no
  known vulnerability associated with it.

At the moment, excluded reports are not served via
[vuln.go.dev](https://vuln.go.dev) API. However, if you have
a specific use case and it would be helpful to have access to this information
through the API,
[please let us know](https://go.dev/s/govulncheck-feedback).
