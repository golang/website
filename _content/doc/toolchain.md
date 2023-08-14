---
title: "Go Toolchains"
layout: article
---

## Introduction {#intro}

Starting in Go 1.21, the Go distribution consists of a `go` command and a bundled Go toolchain,
which is the standard library as well as the compiler, assembler, and other tools.
The `go` command can use its bundled Go toolchain as well as other versions
that it finds in the local `PATH` or downloads as needed.

The choice of Go toolchain being used depends on the `GOTOOLCHAIN` environment setting
and the `go` and `toolchain` lines in the main module's `go.mod` file or the current workspace's `go.work` file.
As you move between different main modules and workspaces,
the toolchain version being used can vary, just as module dependency versions do.

In the standard configuration, the `go` command uses its own bundled toolchain
when that toolchain is at least as new as the `go` or `toolchain` lines in the main module or workspace.
For example, when using the `go` command bundled with Go 1.21.3 in a main module that says `go 1.21.0`,
the `go` command uses Go 1.21.3.
When the `go` or `toolchain` line is newer than the bundled toolchain,
the `go` command runs the newer toolchain instead.
For example, when using the `go` command bundled with Go 1.21.3 in a main module that says `go 1.21.9`,
the `go` command finds and runs Go 1.21.9 instead.
It first looks in the PATH for a program named `go1.21.9` and otherwise downloads and caches
a copy of the Go 1.21.9 toolchain.
This automatic toolchain switching can be disabled, but in that case,
for more precise forwards compatibility,
the `go` command will refuse to run in a main module or workspace in which the `go` line
requires a newer version of Go.
That is, the `go` line sets the minimum required Go version necessary to use a module or workspace.

Modules that are dependencies of other modules may need to set a minimum Go version requirement
lower than the preferred toolchain to use when working in that module directly.
In this case, the `toolchain` line in `go.mod` or `go.work` sets a preferred toolchain
that takes precedence over the `go` line when the `go` command is deciding
which toolchain to use.

The `go` and `toolchain` lines can be thought of as specifying the version requirements
for the module's dependency on the Go toolchain itself, just as the `require` lines in `go.mod`
specify the version requirements for dependencies on other modules.
The `go get` command manages the Go toolchain dependency just as it
manages dependencies on other modules.
For example, `go get go@latest` updates the module to require the latest released Go toolchain.

The `GOTOOLCHAIN` environment setting can force a specific Go version, overriding
the `go` and `toolchain` lines. For example, to test a package with Go 1.21rc3:

	GOTOOLCHAIN=go1.21rc3 go test

The default `GOTOOLCHAIN` setting is `auto`, which enables the toolchain switching described earlier.
The alternate form `<name>+auto` sets the default toolchain to use before deciding whether to
switch further. For example `GOTOOLCHAIN=go1.21.3+auto` directs the `go` command to
begin its decision with a default of using Go 1.21.3 but still use a newer toolchain if
directed by `go` and `toolchain` lines.
Because the default `GOTOOLCHAIN` setting can be changed with `go env -w`,
if you have Go 1.21.0 or later installed, then

	go env -w GOTOOLCHAIN=go1.21.3+auto

is equivalent to replacing your Go 1.21.0 installation with Go 1.21.3.

The rest of this document explains how Go toolchains are versioned, chosen, and managed in more detail.

## Go versions {#version}

Released versions of Go use the version syntax ‘1.*N*.*P*’, denoting the *P*th release of Go 1.*N*.
The initial release is 1.*N*.0, like in ‘1.21.0’. Later releases like 1.*N*.9 are often referred to as patch releases.

Go 1.*N* release candidates, which are issued before 1.*N*.0, use the version syntax ‘1.*N*rc*R*’.
The first release candidate for Go 1.*N* has version 1.*N*rc1, like in `1.23rc1`.

The syntax ‘1.*N*’ is called a “language version”. It denotes the overall family of Go releases
implementing that version of the Go language and standard library.

The language version for a Go version is the result of truncating everything after the *N*:
1.21, 1.21rc2, and 1.21.3 all implement language version 1.21.

Released Go toolchains such as Go 1.21.0 and Go 1.21rc1 report that specific version
(for example, `go1.21.0` or `go1.21rc1`)
from `go version` and [`runtime.Version`](/pkg/runtime/#Version).
Unreleased (still in development) Go toolchains built from the Go development repository
instead report only the language version (for example, `go1.21`).

Any two Go versions can be compared to decide whether one is less than, greater than,
or equal to the other. If the language versions are different, that decides the comparison:
1.21.9 < 1.22. Within a language version, the ordering from least to greatest is:
the language version itself, then release candidates ordered by *R*, then releases ordered by *P*.

For example, 1.21 < 1.21rc1 < 1.21rc2 < 1.21.0 < 1.21.1 < 1.21.2.

Before Go 1.21, the initial release of a Go toolchain was version 1.*N*, not 1.*N*.0,
so for *N* < 21, the ordering is adjusted to place 1.*N* after the release candidates.

For example, 1.20rc1 < 1.20rc2 < 1.20rc3 < 1.20 < 1.20.1.

Earlier versions of Go had beta releases, with versions like 1.18beta2.
Beta releases are placed immediately before release candidates in the version ordering.

For example, 1.18beta1 < 1.18beta2 < 1.18rc1 < 1.18 < 1.18.1.

<!-- Unpublished note: the download page also lists Go 1.9.2rc2, which does not respect
this version syntax. That was created as a test of some potential release automation
before Go 1.9.2 but is not considered a “real” toolchain. -->

## Go toolchain names {#name}

The standard Go toolchains are named <code>go<i>V</i></code> where *V* is a Go version
denoting a beta release, release candidate, or release.
For example, `go1.21rc1` and `go1.21.0` are toolchain names;
`go1.21` and `go1.22` are not (the initial releases are `go1.21.0` and `go1.22.0`),
but `go1.20` and `go1.19` are.

Non-standard toolchains use names of the form <code>go<i>V</i>-<i>suffix</i></code>
for any suffix.

Toolchains are compared by comparing the version <code><i>V</i></code> embedded in the name
(dropping the initial `go` and discarding off any suffix beginning with `-`).
For example, `go1.21.0` and `go1.21.0-custom` compare equal for ordering purposes.

## Module and workspace configuration {#config}

Go modules and workspaces specify version-related configuration
in their `go.mod` or `go.work` files.

The `go` line declares the minimum required Go version for using
the module or workspace.
For compatibility reasons, if the `go` line is omitted from a `go.mod` file,
the module is considered to have an implicit `go 1.16` line,
and if the `go` line is omitted from a `go.work` file,
the workspace is considered to have an implicit `go 1.18` line.

The `toolchain` line declares a suggested toolchain to use with
the module or workspace.
As described in “[Go toolchain selection](#select)” below,
the `go` command may run this specific toolchain when operating
in that module or workspace
if the default toolchain's version is less than the suggested toolchain's version.
If the `toolchain` line is omitted,
the module or workspace is considered to have an implicit
<code>toolchain go<i>V</i></code> line,
where *V* is the Go version from the `go` line.

For example, a `go.mod` that says `go 1.21.0` with no `toolchain` line
is interpreted as if it had a `toolchain go1.21.0` line.

The Go toolchain refuses to load a module or workspace that declares
a minimum required Go version greater than the toolchain's own version.

For example, Go 1.21.2 will refuse to load a module or workspace
with a `go 1.21.3` or `go 1.22` line.

A module's `go` line must declare a version greater than or equal to
the `go` version declared by each of the modules listed in `require` statements.
A workspace's `go` line must declare a version greater than or equal to
the `go` version declared by each of the modules listed in `use` statements.

For example, if module *M* requires a dependency *D* with a `go.mod`
that declares `go 1.22.0`, then *M*'s `go.mod` cannot say `go 1.21.3`.

The `go` line for each module sets the language version the compiler
enforces when compiling packages in that module.
The language version can be changed on a per-file basis by using a
[build constraint](/cmd/go#hdr-Build_constraints).

For example, a module containing code that uses the Go 1.21 language version
should have a `go.mod` file with a `go` line such as `go 1.21` or `go 1.21.3`.
If a specific source file should be compiled only when using a newer Go toolchain,
adding `//go:build go1.22` to that source file both ensures that only Go 1.22 and
newer toolchains will compile the file and also changes the language version in that
file to Go 1.22.

The `go` and `toolchain` lines are most conveniently and safely modified
by using `go get`; see the [section dedicated to `go get` below](#get).

Before Go 1.21, Go toolchains treated the `go` line as an advisory requirement:
if builds succeeded the toolchain assumed everything worked,
and if not it printed a note about the potential version mismatch.
Go 1.21 changed the `go` line to be a mandatory requirement instead.
This behavior is partly backported to earlier language versions:
Go 1.19 releases starting at Go 1.19.13 and Go 1.20 releases starting at Go 1.20.8,
refuse to load workspaces or modules declaring version Go 1.22 or later.

Before Go 1.21, toolchains did not require a module
or workspace to have a `go` line greater than or equal to the
`go` version required by each of its dependency modules.

## The `GOTOOLCHAIN` setting {#GOTOOLCHAIN}

The `go` command selects the Go toolchain to use based on the `GOTOOLCHAIN` setting.
To find the `GOTOOLCHAIN` setting, the `go` command uses the standard rules for any
Go environment setting:

 - If `GOTOOLCHAIN` is set to a non-empty value in the process environment
   (as queried by [`os.Getenv`](/pkg/os/#Getenv)), the `go` command uses that value.

 - Otherwise, if `GOTOOLCHAIN` is set in the user's environment default file
   (managed with
   [`go env -w` and `go env -u`](/cmd/go/#hdr-Print_Go_environment_information)),
   the `go` command uses that value.

 - Otherwise, if `GOTOOLCHAIN` is set in the bundled Go toolchain's environment
   default file (`$GOROOT/go.env`), the `go` command uses that value.

In standard Go toolchains, the `$GOROOT/go.env` file sets the default `GOTOOLCHAIN=auto`,
but repackaged Go toolchains may change this value.

If the `$GOROOT/go.env` file is missing or does not set a default, the `go` command
assumes `GOTOOLCHAIN=local`.

Running `go env GOTOOLCHAIN` prints the `GOTOOLCHAIN` setting.

## Go toolchain selection {#select}

At startup, the `go` command selects which Go toolchain to use.
It consults the `GOTOOLCHAIN` setting,
which takes the form `<name>`, `<name>+auto`, or `<name>+path`.
`GOTOOLCHAIN=auto` is shorthand for `GOTOOLCHAIN=local+auto`;
similarly, `GOTOOLCHAIN=path` is shorthand for `GOTOOLCHAIN=local+path`.
The `<name>` sets the default Go toolchain:
`local` indicates the bundled Go toolchain
(the one that shipped with the `go` command being run), and otherwise `<name>` must
be a specific Go toolchain name, such as `go1.21.0`.
The `go` command prefers to run the default Go toolchain.
As noted above, starting in Go 1.21, Go toolchains refuse to run in
workspaces or modules that require newer Go versions.
Instead, they report an error and exit.

When `GOTOOLCHAIN` is set to `local`, the `go` command always runs the bundled Go toolchain.

When `GOTOOLCHAIN` is set to `<name>` (for example, `GOTOOLCHAIN=go1.21.0`),
the `go` command always runs that specific Go toolchain.
If a binary with that name is found in the system PATH, the `go` command uses it.
Otherwise the `go` command uses a Go toolchain it downloads and verifies.

When `GOTOOLCHAIN` is set to `<name>+auto` or `<name>+path` (or the shorthands `auto` or `path`),
the `go` command selects and runs a newer Go version as needed.
Specifically, it consults the `toolchain` and `go` lines in the current workspace's
`go.work` file or, when there is no workspace,
the main module's `go.mod` file.
If the `go.work` or `go.mod` file has a `toolchain <tname>` line
and `<tname>` is newer than the default Go toolchain,
then the `go` command runs `<tname>` instead.
If the file has a `toolchain default` line,
then the `go` command runs the default Go toolchain,
disabling any attempt at updating beyond `<name>`.
Otherwise, if the file has a `go <version>` line
and `<version>` is newer than the default Go toolchain,
then the `go` command runs `go<version>` instead.

To run a toolchain other than the bundled Go toolchain,
the `go` command searches the process's executable path
(`$PATH` on Unix and Plan 9, `%PATH%` on Windows)
for a program with the given name (for example, `go1.21.3`) and runs that program.
If no such program is found, the `go` command
[downloads and runs the specified Go toolchain](#download).
Using the `GOTOOLCHAIN` form `<name>+path` disables the download fallback,
causing the `go` command to stop after searching the executable path.

Running `go version` prints the selected Go toolchain's version
(by running the selected toolchain's implementation of `go version`).

Running `GOTOOLCHAIN=local go version` prints the bundled Go toolchain's version.

## Go toolchain switches {#switch}

For most commands, the workspace's `go.work` or the main module's `go.mod`
will have a `go` line that is at least as new as the `go` line in any module dependency,
due to the version ordering [configuration requirements](#config).
In this case, the startup toolchain selection runs a new enough Go toolchain
to complete the command.

Some commands incorporate new module versions as part of their operation:
`go get` adds new module dependencies to the main module;
`go work use` adds new local modules to the workspace;
`go work sync` resynchronizes a workspace with local modules that may have been updated
since the workspace was created;
`go install package@version` and `go run package@version`
effectively run in an empty main module and add `package@version` as a new dependency.
All these commands may encounter a module with a `go.mod` `go` line
requiring a newer Go version than the currently executed Go version.

When a command encounters a module requiring a newer Go version
and `GOTOOLCHAIN` permits running different toolchains
(it is one of the `auto` or `path` forms),
the `go` command chooses and switches to an appropriate newer toolchain
to continue executing the current command.

Any time the `go` command switches toolchains after startup toolchain selection,
it prints a message explaining why. For example:

	go: module example.com/widget@v1.2.3 requires go >= 1.24rc1; switching to go 1.27.9

As shown in the example, the `go` command may switch to a toolchain
newer than the discovered requirement.
In general the `go` command aims to switch to a supported Go toolchain.

To choose the toolchain, the `go` command first obtains a list of available toolchains.
For the `auto` form, the `go` command downloads a list of available toolchains.
For the `path` form, the `go` command scans the PATH for any executables
named for valid toolchains and uses a list of all the toolchains it finds.
Using that list of toolchains, the `go` command identifies up to three candidates:

 - the latest release candidate of an unreleased Go language version (1.*N*₃rc*R*₃),
 - the latest patch release of the most recently released Go language version (1.*N*₂.*P*₂), and
 - the latest patch release of the previous Go language version (1.*N*₁.*P*₁).

These are the supported Go releases according to Go's [release policy](/doc/devel/release#policy).
Consistent with [minimal version selection](https://research.swtch.com/vgo-mvs),
the `go` command then conservatively uses the candidate with the _minimum_ (oldest)
version that satisfies the new requirement.

For example, suppose `example.com/widget@v1.2.3` requires Go 1.24rc1 or later.
The `go` command obtains the list of available toolchains
and finds that the latest patch releases of the two most recent Go toolchains are
Go 1.28.3 and Go 1.27.9,
and the release candidate Go 1.29rc2 is also available.
In this situation, the `go` command will choose Go 1.27.9.
If `widget` had required Go 1.28 or later, the `go` command would choose Go 1.28.3,
because Go 1.27.9 is too old.
If `widget` had required Go 1.29 or later, the `go` command would choose Go 1.29rc2,
because both Go 1.27.9 and Go 1.28.3 are too old.

Commands that incorporate new module versions that require new Go versions
write the new minimum `go` version requirement to the current workspace's `go.work` file
or the main module's `go.mod` file, updating the `go` line.
For [repeatability](https://research.swtch.com/vgo-principles#repeatability),
any command that updates the `go` line also updates the `toolchain` line
to record its own toolchain name.
The next time the `go` command runs in that workspace or module,
it will use that updated `toolchain` line during [toolchain selection](#select).

For example, `go get example.com/widget@v1.2.3` may print a switching notice
like above and switch to Go 1.27.9.
Go 1.27.9 will complete the `go get` and update the `toolchain` line
to say `toolchain go1.27.9`.
The next `go` command run in that module or workspace will select `go1.27.9`
during startup and will not print any switching message.

In general, if any `go` command is run twice, if the first prints a switching
message, the second will not, because the first also updated `go.work` or `go.mod`
to select the right toolchain at startup.
The exception is the `go install package@version` and `go run package@version` forms,
which run in no workspace or main module and cannot write a `toolchain` line.
They print a switching message every time they need to switch
to a newer toolchain.

## Downloading toolchains {#download}

When using `GOTOOLCHAIN=auto` or `GOTOOLCHAIN=<name>+auto`, the Go command
downloads newer toolchains as needed.
These toolchains are packaged as special modules
with module path `golang.org/toolchain`
and version <code>v0.0.1-go<i>VERSION</i>.<i>GOOS</i>-<i>GOARCH</i></code>.
Toolchains are downloaded like any other module,
meaning that toolchain downloads can be proxied by setting `GOPROXY`
and have their checksums checked by the Go checksum database.
Because the specific toolchain used depends on the system's own
default toolchain as well as the local operating system and architecture (GOOS and GOARCH),
it is not practical to write toolchain module checksums to `go.sum`.
Instead, toolchain downloads fail for lack of verification if `GOSUMDB=off`.
`GOPRIVATE` and `GONOSUMDB` patterns do not apply to the toolchain downloads.

## Managing Go version module requirements with `go get` {#get}

In general the `go` command treats the `go` and `toolchain` lines
as declaring versioned toolchain dependencies of the main module.
The `go get` command can manage these lines just as it manages
the `require` lines that specify versioned module dependencies.

For example, `go get go@1.22.1 toolchain@1.24rc1` changes the main module's
`go.mod` file to read `go 1.22.1` and `toolchain go1.24rc1`.

The `go` command understands that the `go` dependency requires a `toolchain` dependency
with a greater or equal Go version.

Continuing the example, a later `go get go@1.25.0` will update
the toolchain to `go1.25.0` as well.
When the toolchain matches the `go` line exactly, it can be
omitted and implied, so this `go get` will delete the `toolchain` line.

The same requirement applies in reverse when downgrading:
if the `go.mod` starts at `go 1.22.1` and `toolchain go1.24rc1`,
then `go get toolchain@go1.22.9` will update only the `toolchain` line,
but `go get toolchain@go1.21.3` will downgrade the `go` line to
`go 1.21.3` as well.
The effect will be to leave just `go 1.21.3` with no `toolchain` line.

The special form `toolchain@none` means to remove any `toolchain` line,
as in `go get toolchain@none` or `go get go@1.25.0 toolchain@none`.

The `go` command understands the version syntax for
`go` and `toolchain` dependencies as well as queries.

For example, just as `go get example.com/widget@v1.2` uses
the latest `v1.2` version of `example.com/widget` (perhaps `v1.2.3`),
`go get go@1.22` uses the latest available release of the Go 1.22 language version
(perhaps `1.22rc3`, or perhaps `1.22.3`).
The same applies to `go get toolchain@go1.22`.

The `go get` and `go mod tidy` commands maintain the `go` line to
be greater than or equal to the `go` line of any required dependency module.

For example, if the main module has `go 1.22.1` and we run
`go get example.com/widget@v1.2.3` which declares `go 1.24rc1`,
then `go get` will update the main module's `go` line to `go 1.24rc1`.

Continuing the example, a later `go get go@1.22.1` will
downgrade `example.com/widget` to a version compatible with Go 1.22.1
or else remove the requirement entirely,
just as it would when downgrading any other dependency of `example.com/widget`.

Before Go 1.21, the suggested way to update a module to a new Go version (say, Go 1.22)
was `go mod tidy -go=1.22`, to make sure that any adjustments
specific to Go 1.22 were made to the `go.mod` at the same time that the
`go` line is updated.
That form is still valid, but the simpler `go get go@1.22` is now preferred.

When `go get` is run in a module in a directory contained in a workspace root,
`go get` mostly ignores the workspace,
but it does update the `go.work` file to upgrade the `go` line
when the workspace would otherwise be left with too old a `go` line.

## Managing Go version workspace requirements with `go work` {#work}

As noted in the previous section, `go get` run in a directory
inside a workspace root will take care to update the `go.work` file's `go` line
as needed to be greater than or equal to any module inside that root.
However, workspaces can also refer to modules outside the root directory;
running `go get` in those directories may result in an invalid workspace
configuration, one in which the `go` version declared in `go.work` is less
than one or more of the modules in the `use` directives.

The command `go work use`, which adds new `use` directives, also checks
that the `go` version in the `go.work` file is new enough for all the
existing `use` directives.
To update a workspace that has gotten its `go` version out of sync
with its modules, run `go work use` with no arguments.

The commands `go work init` and `go work sync` also update the `go`
version as needed.

To remove the `toolchain` line from a `go.work` file, use
`go work edit -toolchain=none`.

