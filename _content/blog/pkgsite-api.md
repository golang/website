---
title: Introducing the pkg.go.dev API
date: 2026-05-21
by:
- Ethan Lee
- Hana Kim
- Jonathan Amsterdam
summary: Introducing the new programmatic API for pkg.go.dev, allowing developers to fetch package and module data directly.
---

Since its inception, [pkg.go.dev](https://pkg.go.dev) has established itself as
the Go community's primary resource for package documentation and discovery.
While we initially prioritized creating a comprehensive and highly accessible
web interface for users, the need for programmatic access has become
increasingly clear. Developers building tools, IDE integrations, and automated
workflows have historically relied on fragile workarounds like web scraping to
access this data. To better address these evolving requirements, we are now
expanding our platform to provide robust, direct access to the information our
community needs.

Today, we are excited to introduce the official
[pkg.go.dev API](https://pkg.go.dev/api) — a service interface for querying
metadata about published Go modules. This launch is a direct response to
years of community feedback. Furthermore, the need for a formalized
interface has become even more acute with the rise of AI-assisted coding.
Tools can now access the specific, high-fidelity context needed to reason
about the Go ecosystem with greater precision.

## The service interface

Built for stability and efficient caching, the API uses a stateless, GET-only
architecture. Primary endpoints are currently hosted under the `/v1beta` path.
Following a period of community feedback and confirmed stability, we intend to
transition toward a formal `v1` release.

For a complete interactive reference of all endpoints, query parameters, and
response shapes, see the [pkg.go.dev/api specification](https://pkg.go.dev/api).
The machine-readable API contract is also published directly as an [OpenAPI
specification](https://pkg.go.dev/v1beta/openapi.yaml).

### Core endpoints

| Endpoint | Description |
| :--- | :--- |
| `/v1beta/package/{path}` | Information about the package at `{path}`. |
| `/v1beta/module/{path}` | Information about the module at `{path}`. |
| `/v1beta/versions/{path}` | Versions of the module at `{path}`. |
| `/v1beta/packages/{path}` | Information about packages of the module at `{path}`. |
| `/v1beta/search?q={query}` | Search results for a given query. |
| `/v1beta/symbols/{path}` | List of symbols declared by the package at `{path}`. |
| `/v1beta/imported-by/{path}` | Paths of packages importing the package at `{path}`. |
| `/v1beta/vulns/{path}` | Vulnerabilities of the module or package at `{path}`. |

One design principle for this API is "precision over convenience." For context,
when `go mod tidy` encounters an import of a package that isn't provided by an
existing dependency of the main module, it applies the "longest module path"
rule to determine which module is needed. (The fact that two or more modules
could provide the package is what makes it possible to later carve out a
submodule without breaking existing programs.) The
[pkg.go.dev](https://pkg.go.dev) web interface follows a similar convention
when choosing which package to display for a given package path.
By contrast, the [pkg.go.dev](https://pkg.go.dev) API requires the module to be
specified unambiguously. If a package path is ambiguous because it exists in
multiple modules, the API returns a list of candidates and reports an error
asking the client to be more specific.

For example, a package imported as `example.com/a/b/c` could be provided by
module `example.com/a` or by `example.com/a/b`. While the
[pkg.go.dev](https://pkg.go.dev) web interface will automatically resolve the
"longest module path" (`example.com/a/b`), a client querying the API must
specify the module explicitly to avoid an ambiguous resolution error.

### Specifying versions

For endpoints that retrieve package, module, or symbol information, you can
specify the desired version using the optional `version` query parameter. The
API returns information about the latest version of the module or package by
default. The parameter supports:

* **Semantic Versions:** Retrieve data for a specific release tag (e.g.,
  `?version=v1.2.3` or `?version=v0.6.0`).
* **Branch Names:** Reference default development branches—specifically `master`
  or `main` (e.g., `?version=master`). The API will automatically resolve the
  branch to its corresponding pseudo-version. Note that custom or arbitrary
  branch names are not supported.

If the `version` parameter is omitted, the API defaults to resolving the
request against the latest tagged version of the package or module.

### Example: raw API request

To retrieve structured metadata for a specific package directly (using `jq` for
formatting):

```console
$ curl https://pkg.go.dev/v1beta/package/github.com/google/go-cmp/cmp | jq .
{
  "modulePath": "github.com/google/go-cmp",
  "version": "v0.7.0",
  "isLatest": true,
  "isStandardLibrary": false,
  "goos": "all",
  "goarch": "all",
  "path": "github.com/google/go-cmp/cmp",
  "name": "cmp",
  "synopsis": "Package cmp determines equality of values.",
  "isRedistributable": true
}
```

To query a specific branch version (like `master`) and see it resolve
automatically to its corresponding pseudo-version:

```console
$ curl -s "https://pkg.go.dev/v1beta/package/github.com/google/go-cmp/cmp?version=master" | jq '{path, version}'
{
  "path": "github.com/google/go-cmp/cmp",
  "version": "v0.7.1-0.20260310220054-34c9473539b8"
}
```

## The pkgsite-cli reference implementation

To demonstrate how to interact with our API, we are providing a reference
client implementation:
[pkgsite-cli](https://github.com/golang/pkgsite/tree/master/cmd/internal/pkgsite-cli).
This implementation serves as a practical example for developers looking to
build their own integrations, showing how to handle the data directly from the
terminal. Please be aware that as the API continues to evolve, the interface and
behavior of this command may change.

To get started, install the command:

```bash
$ go install golang.org/x/pkgsite/cmd/internal/pkgsite-cli@latest
```

To search for packages:

```
$ pkgsite-cli search "uuid"
github.com/google/uuid
  Module:   github.com/google/uuid@v1.6.0
  Synopsis: Package uuid generates and inspects UUIDs.
... more
```

To inspect a specific package:

```
$ pkgsite-cli package github.com/google/go-cmp/cmp
github.com/google/go-cmp/cmp
  Name:      cmp
  Module:    github.com/google/go-cmp
  Version:   v0.7.0 (latest)
  Synopsis:  Package cmp determines equality of values.
```

To see which packages import a specific package:

```
$ pkgsite-cli package --imported-by github.com/google/go-cmp/cmp
github.com/google/go-cmp/cmp
  Name:     cmp
  Module:   github.com/google/go-cmp
  Version:  v0.7.0 (latest)
  Synopsis: Package cmp determines equality of values.

Imported by:
  cloud.google.com/go/internal/testutil
  cuelang.org/go/internal/cuetxtar
  chainguard.dev/apko/pkg/build/types
  ... more
```

To list symbols declared by a package:

```
$ pkgsite-cli package --symbols github.com/google/go-cmp/cmp
github.com/google/go-cmp/cmp
  Name:     cmp
  Module:   github.com/google/go-cmp
  Version:  v0.7.0 (latest)
  Synopsis: Package cmp determines equality of values.

Symbols:
  type Indirect struct{}
  type MapIndex struct{}
  type Option interface{}
  ... more
```

To list versions of a module:

```
$ pkgsite-cli module -versions github.com/google/go-cmp
github.com/google/go-cmp
  Version:          v0.7.0 (latest)
  Repository:       https://github.com/google/go-cmp
  Has go.mod:       yes
  Redistributable:  yes

Versions:
  v0.7.0
  v0.6.0
  v0.5.9
  ... more
```

To list both versions and packages of a module:

```
$ pkgsite-cli module -packages -versions github.com/google/go-cmp
github.com/google/go-cmp
  Version:          v0.7.0 (latest)
  Repository:       https://github.com/google/go-cmp
  Has go.mod:       yes
  Redistributable:  yes

Versions:
  v0.7.0
  v0.6.0
  v0.5.9
  ... more

Packages:
  github.com/google/go-cmp/cmp             Package cmp determines equality of values.
  github.com/google/go-cmp/cmp/cmpopts     Package cmpopts provides common options for the cmp package.
  ... more
```

The command handles pagination and formatting, allowing you to focus on the
data you need for your scripts or manual investigation. To learn more, please
visit [pkgsite-cli's
documentation](https://pkg.go.dev/golang.org/x/pkgsite/cmd/internal/pkgsite-cli).

## Stability and the future

This concludes our brief tour of the [pkg.go.dev](https://pkg.go.dev) API. While
we plan to expand the interface's capabilities over time, we are committed to
maintaining backward compatibility so that existing integrations continue to
function seamlessly. (Note that command line interface of the `pkgsite-cli` reference client is not yet stable.) We welcome your feedback via our [issue
tracker](https://github.com/golang/go/issues), and we look forward to seeing the
new tools and workflows the community will build.
