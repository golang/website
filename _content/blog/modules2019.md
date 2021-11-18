---
title: Go Modules in 2019
date: 2018-12-19
by:
- Russ Cox
tags:
- tools
- versioning
summary: What the Go team is planning for Go modules in 2019.
---

## What a year!

2018 was a great year for the Go ecosystem, with
package management as one of our major focuses.
In February, we started a community-wide discussion about how to integrate
package management directly into the Go toolchain,
and in August we delivered the first rough implementation of that feature,
called Go modules, in Go 1.11.
The migration to Go modules will be the most far-reaching change
for the Go ecosystem since Go 1.
Converting the entire ecosystem—code, users, tools, and so on—from
GOPATH to modules will require work in many different areas.
The module system will in turn help us deliver
better authentication and build speeds to the Go ecosystem.

This post is a preview of what the Go team is planning
relating to modules in 2019.

## Releases

Go 1.11, released in August 2018, introduced [preliminary support for modules](/doc/go1.11#modules).
For now, module support is maintained alongside the
traditional GOPATH-based mechanisms.
The `go` command defaults to module mode when run
in directory trees outside GOPATH/src and
marked by `go.mod` files in their roots.
This setting can be overridden by setting the transitional
environment variable `$GO111MODULE` to `on` or `off`;
the default behavior is `auto` mode.
We’ve already seen significant adoption of modules across the Go community,
along with many helpful suggestions and bug reports
to help us improve modules.

Go 1.12, scheduled for February 2019, will refine module support
but still leave it in `auto` mode by default.
In addition to many bug fixes and other minor improvements,
perhaps the most significant change in Go 1.12
is that commands like `go` `run` `x.go`
or `go` `get` `rsc.io/2fa@v1.1.0`
can now operate in `GO111MODULE=on` mode without an explicit `go.mod` file.

Our aim is for Go 1.13, scheduled for August 2019, to enable module mode by
default (that is, to change the default from `auto` to `on`)
and deprecate GOPATH mode.
In order to do that, we’ve been working on better tooling support
along with better support for the open-source module ecosystem.

## Tooling & IDE Integration

In the eight years that we’ve had GOPATH,
an incredible amount of tooling has been created
that assumes Go source code is stored in GOPATH.
Moving to modules requires updating all code that makes
that assumption.
We’ve designed a new package,
[golang.org/x/tools/go/packages](https://godoc.org/golang.org/x/tools/go/packages),
that abstracts the operation of finding and loading information
about the Go source code for a given target.
This new package adapts automatically to both
GOPATH and modules mode and is also extensible
to tool-specific code layouts, such as the one
used by Bazel.
We’ve been working with tool authors throughout the Go community
to help them adopt golang.org/x/tools/go/packages in their tools.

As part of this effort, we’ve also been working to
unify the various source code querying tools
like gocode, godef, and go-outline
into a single tool that can be used from the
command line and also supports
the [language server protocol](https://langserver.org/)
used by modern IDEs.

The transition to modules and the changes in package loading
also prompted a significant change to Go program analysis.
As part of reworking `go` `vet` to support modules,
we introduced a generalized framework for incremental
analysis of Go programs,
in which an analyzer is invoked for one package at a time.
In this framework, the analysis of one package can write out facts
made available to analyses of other packages that import the first.
For example, `go` `vet`’s analysis of the [log package](/pkg/log/)
determines and records the fact that `log.Printf` is a `fmt.Printf` wrapper.
Then `go` `vet` can check printf-style format strings in other packages
that call `log.Printf`.
This framework should enable many new, sophisticated
program analysis tools to help developers find bugs earlier
and understand code better.

## Module Index

One of the most important parts of the original design for `go` `get`
was that it was _decentralized_:
we believed then—and we still believe today—that
anyone should be able to publish their code on any server,
in contrast to central registries
such as Perl’s CPAN, Java’s Maven, or Node’s NPM.
Placing domain names at the start of the `go` `get` import space
reused an existing decentralized system
and avoided needing to solve anew the problems of
deciding who can use which names.
It also allowed companies to import code on private servers
alongside code from public servers.
It is critical to preserve this decentralization as we shift to Go modules.

Decentralization of Go’s dependencies has had many benefits,
but it also brought a few significant drawbacks.
The first is that it’s too hard to find all the publicly-available Go packages.
Every site that wants to deliver information about
packages has
to do its own crawling, or else wait until a user asks
about a particular package before fetching it.

We are working on a new service, the Go Module Index,
that will provide a public log of packages entering the Go ecosystem.
Sites like godoc.org and goreportcard.com will be able to watch this log
for new entries instead of each independently implementing code
to find new packages.
We also want the service to allow looking up packages
using simple queries, to allow `goimports` to add
imports for packages that have not yet been downloaded to the local system.

## Module Authentication

Today, `go` `get` relies on connection-level authentication (HTTPS or SSH)
to check that it is talking to the right server to download code.
There is no additional check of the code itself,
leaving open the possibility of man-in-the-middle attacks
if the HTTPS or SSH mechanisms are compromised in some way.
Decentralization means that the code for a build is fetched
from many different servers, which means the build depends on
many systems to serve correct code.

The Go modules design improves code authentication by storing
a `go.sum` file in each module;
that file lists the cryptographic hash
of the expected file tree for each of the module’s dependencies.
When using modules, the `go` command uses `go.sum` to verify
that dependencies are bit-for-bit identical to the expected versions
before using them in a build.
But the `go.sum` file only lists hashes for the specific dependencies
used by that module.
If you are adding a new dependency
or updating dependencies with `go` `get` `-u`,
there is no corresponding entry in `go.sum` and therefore
no direct authentication of the downloaded bits.

For publicly-available modules, we intend to run a service we call a _notary_
that follows the module index log,
downloads new modules,
and cryptographically signs statements of the form
“module M at version V has file tree hash H.”
The notary service will publish all these notarized hashes
in a queryable, [Certificate Transparency](https://www.certificate-transparency.org/)-style
[tamper-proof log](http://static.usenix.org/event/sec09/tech/full_papers/crosby.pdf),
so that anyone can verify that the notary is behaving correctly.
This log will serve as a public, global `go.sum` file
that `go` `get` can use to authenticate modules
when adding or updating dependencies.

We are aiming to have the `go` command check notarized hashes
for publicly-available modules not already in `go.sum`
starting in Go 1.13.

## Module Mirrors

Because the decentralized `go` `get` fetches code from multiple origin servers,
fetching code is only as fast and reliable as the slowest,
least reliable server.
The only defense available before modules was to vendor
dependencies into your own repositories.
While vendoring will continue to be supported,
we’d prefer a solution that works for all modules—not just the ones you’re already using—and
that does not require duplicating a dependency into every
repository that uses it.

The Go module design introduces the idea of a module proxy,
which is a server that the `go` command asks for modules,
instead of the origin servers.
One important kind of proxy is a _module mirror_,
which answers requests for modules by fetching them
from origin servers and then caching them for use in
future requests.
A well-run mirror should be fast and reliable
even when some origin servers have gone down.
We are planning to launch a mirror service for publicly-available modules in 2019.
Other projects, like GoCenter and Athens, are planning mirror services too.
(We anticipate that companies will have multiple options for running
their own internal mirrors as well, but this post is focusing on public mirrors.)

One potential problem with mirrors is that they are
precisely man-in-the-middle servers,
making them a natural target for attacks.
Go developers need some assurance that the mirrors
are providing the same bits that the origin servers would.
The notary process we described in the previous section
addresses exactly this concern, and it
will apply to downloads using mirrors
as well as downloads using origin servers.
The mirrors themselves need not be trusted.

We are aiming to have the Google-run module mirror
ready to be used by default in the `go` command starting in Go 1.13.
Using an alternate mirror, or no mirror at all, will be trivial
to configure.

## Module Discovery

Finally, we mentioned earlier that the module index will make it easier to
build sites like godoc.org.
Part of our work in 2019 will be a major revamp of godoc.org
to make it more useful for developers who need to
discover available modules
and then decide whether to rely on a given module or not.

## Big Picture

This diagram shows how module source code
moves through the design in this post.

{{image "modules2019/code.png" 374}}

Before, all consumers of Go source code—the `go` command
and any sites like godoc.org—fetched code directly from each code host.
Now they can fetch cached code from a fast, reliable mirror,
while still authenticating that the downloaded bits are correct.
And the index service makes it easy for mirrors, godoc.org,
and any other similar sites to keep up with all the great new
code being added to the Go ecosystem every day.

We’re excited about the future of Go modules in 2019,
and we hope you are too. Happy New Year!
