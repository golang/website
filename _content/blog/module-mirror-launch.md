---
title: Module Mirror and Checksum Database Launched
date: 2019-08-29
by:
- Katie Hockman
tags:
- tools
- versioning
summary: The Go module mirror and checksum database provide faster, verified downloads of your Go dependencies.
---


We are excited to share that our module [mirror](https://proxy.golang.org),
[index](https://index.golang.org), and
[checksum database](https://sum.golang.org) are now production ready! The `go` command
will use the module mirror and checksum database by default for
[Go 1.13 module users](/doc/go1.13#introduction).  See
[proxy.golang.org/privacy](https://proxy.golang.org/privacy) for privacy
information about these services and the
[go command documentation](/cmd/go/#hdr-Module_downloading_and_verification)
for configuration details, including how to disable the use of these servers or
use different ones.  If you depend on non-public modules, see the
[documentation for configuring your environment](/cmd/go/#hdr-Module_configuration_for_non_public_modules).

This post will describe these services and the benefits of using them, and
summarizes some of the points from the
[Go Module Proxy: Life of a Query](https://youtu.be/KqTySYYhPUE) talk at Gophercon 2019.
See the [recording](https://youtu.be/KqTySYYhPUE) if you are interested in the full talk.

## Module Mirror

[Modules](https://blog.golang.org/versioning-proposal) are sets of Go packages
that are versioned together, and the contents of each version are immutable.
That immutability provides new opportunities for caching and authentication.
When `go get` runs in module mode, it must fetch the module containing the
requested packages, as well as any new dependencies introduced by that module,
updating your
[go.mod](/cmd/go/#hdr-The_go_mod_file) and
[go.sum](/cmd/go/#hdr-Module_downloading_and_verification)
files as needed. Fetching modules from version control can be expensive in terms
of latency and storage in your system: the `go` command may be forced to pull down
the full commit history of a repository containing a transitive dependency, even
one that isn’t being built, just to resolve its version.

The solution is to use a module proxy, which speaks an API that is better suited
to the `go` command’s needs (see `go help goproxy`). When `go get` runs in
module mode with a proxy, it will work faster by only asking for the specific
module metadata or source code it needs, and not worrying about the rest. Below is
an example of how the `go` command may use a proxy with `go get` by requesting the list
of versions, then the info, mod, and zip file for the latest tagged version.

{{image "module-mirror-launch/proxy-protocol.png" 800}}

A module mirror is a special kind of module proxy that caches metadata and
source code in its own storage system, allowing the mirror to continue to serve
source code that is no longer available from the original locations. This can
speed up downloads and protect you from disappearing dependencies. See
[Go Modules in 2019](https://blog.golang.org/modules2019) for more information.

The Go team maintains a module mirror, served at
[proxy.golang.org](https://proxy.golang.org), which the `go` command will use by
default for module users as of Go 1.13. If you are running an earlier version of the `go`
command, then you can use this service by setting
`GOPROXY=https://proxy.golang.org` in your local environment.

## Checksum Database

Modules introduced the `go.sum` file, which is a list of SHA-256 hashes of the
source code and `go.mod` files of each dependency when it was first downloaded.
The `go` command can use the hashes to detect misbehavior by an origin server or
proxy that gives you different code for the same version.

The limitation of this `go.sum` file is that it works entirely by trust on _your_
first use. When you add a version of a dependency that you’ve never seen before
to your module (possibly by upgrading an existing dependency), the `go` command
fetches the code and adds lines to the `go.sum` file on the fly. The problem is
that those `go.sum` lines aren’t being checked against anyone else’s: they might
be different from the `go.sum` lines that the `go` command just generated for
someone else, perhaps because a proxy intentionally served malicious code
targeted to you.

Go's solution is a global source of `go.sum` lines, called a
[checksum database](https://go.googlesource.com/proposal/+/master/design/25530-sumdb.md#checksum-database),
which ensures that the `go` command always adds the same lines to everyone's
`go.sum` file. Whenever the `go` command receives new source code, it can verify the
hash of that code against this global database to make sure the hashes match,
ensuring that everyone is using the same code for a given version.

The checksum database is served by [sum.golang.org](https://sum.golang.org), and
is built on a [Transparent Log](https://research.swtch.com/tlog) (or “Merkle
tree”) of hashes backed by [Trillian](https://github.com/google/trillian). The
main advantage of a Merkle tree is that it is tamper proof and has properties
that don’t allow for misbehavior to go undetected, which makes it more
trustworthy than a simple database. The `go` command uses this tree to check
“inclusion” proofs (that a specific record exists in the log) and “consistency”
proofs (that the tree hasn’t been tampered with) before adding new `go.sum` lines
to your module’s `go.sum` file. Below is an example of such a tree.

{{image "module-mirror-launch/tree.png" 800}}

The checksum database supports
[a set of endpoints](https://go.googlesource.com/proposal/+/master/design/25530-sumdb.md#checksum-database)
used by the `go` command to request and verify `go.sum` lines. The `/lookup`
endpoint provides a “signed tree head” (STH) and the requested `go.sum` lines. The
`/tile` endpoint provides chunks of the tree called _tiles_ which the `go` command
can use for proofs. Below is an example of how the `go` command may
interact with the checksum database by doing a `/lookup` of a module version, then
requesting the tiles required for the proofs.

{{image "module-mirror-launch/sumdb-protocol.png" 800}}

This checksum database allows the `go` command to safely use an otherwise
untrusted proxy. Because there is an auditable security layer sitting on top of
it, a proxy or origin server can’t intentionally, arbitrarily, or accidentally
start giving you the wrong code without getting caught. Even the author of a
module can’t move their tags around or otherwise change the bits associated with
a specific version from one day to the next without the change being detected.

If you are using Go 1.12 or earlier, you can manually check a `go.sum` file
against the checksum database with
[gosumcheck](https://godoc.org/golang.org/x/mod/gosumcheck):

	$ go get golang.org/x/mod/gosumcheck
	$ gosumcheck /path/to/go.sum

In addition to verification done by the `go` command, third-party
auditors can hold the checksum database accountable by iterating over the log
looking for bad entries. They can work together and gossip about the state of
the tree as it grows to ensure that it remains uncompromised, and we hope that
the Go community will run them.

## Module Index

The module index is served by [index.golang.org](https://index.golang.org), and
is a public feed of new module versions that become available through
[proxy.golang.org](https://proxy.golang.org). This is particularly useful for
tool developers that want to keep their own cache of what’s available in
[proxy.golang.org](https://proxy.golang.org), or keep up-to-date on some of the
newest modules that people are using.

## Feedback or bugs

We hope these services improve your experience with modules, and encourage you
to [file issues](https://github.com/golang/go/issues/new?title=proxy.golang.org) if you run into
problems or have feedback!
