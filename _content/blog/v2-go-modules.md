---
title: "Go Modules: v2 and Beyond"
date: 2019-11-07
by:
- Jean de Klerk
- Tyler Bui-Palsulich
tags:
- tools
- versioning
summary: How to release major version 2 of your module.
---

## Introduction

This post is part 4 in a series.

  - Part 1 — [Using Go Modules](/blog/using-go-modules)
  - Part 2 — [Migrating To Go Modules](/blog/migrating-to-go-modules)
  - Part 3 — [Publishing Go Modules](/blog/publishing-go-modules)
  - **Part 4 — Go Modules: v2 and Beyond** (this post)
  - Part 5 — [Keeping Your Modules Compatible](/blog/module-compatibility)

**Note:** For documentation on developing modules, see
[Developing and publishing modules](/doc/modules/developing).

As a successful project matures and new requirements are added, past features
and design decisions might stop making sense. Developers may want to integrate
lessons they've learned by removing deprecated functions, renaming types, or
splitting complicated packages into manageable pieces. These kinds of changes
require effort by downstream users to migrate their code to the new API, so they
should not be made without careful consideration that the benefits outweigh the
costs.

For projects that are still experimental — at major version `v0` — occasional
breaking changes are expected by users. For projects which are declared stable
— at major version `v1` or higher — breaking changes must be done in a new major
version. This post explores major version semantics, how to create and publish a new
major version, and how to maintain multiple major versions of a module.

## Major versions and module paths

Modules formalized an important principle in Go, the
[**import compatibility rule**](https://research.swtch.com/vgo-import):

	If an old package and a new package have the same import path,
	the new package must be backwards compatible with the old package.

By definition, a new major version of a package is not backwards compatible with
the previous version. This means a new major version of a module must have a
different module path than the previous version. Starting with `v2`, the major
version must appear at the end of the module path (declared in the `module`
statement in the `go.mod` file). For example, when the authors of the module
`github.com/googleapis/gax-go` developed `v2`, they used the new module path
`github.com/googleapis/gax-go/v2`. Users who wanted to use `v2` had to change
their package imports and module requirements to `github.com/googleapis/gax-go/v2`.

The need for major version suffixes is one of the ways Go modules differs from
most other dependency management systems. Suffixes are needed to solve
the [diamond dependency problem](https://research.swtch.com/vgo-import#dependency_story).
Before Go modules, [gopkg.in](http://gopkg.in) allowed package maintainers to
follow what we now refer to as the import compatibility rule. With gopkg.in, if
you depend on a package that imports `gopkg.in/yaml.v1` and another package that
imports `gopkg.in/yaml.v2`, there is no conflict because the two `yaml` packages
have different import paths — they use a version suffix, as with Go modules.
Since gopkg.in shares the same version suffix methodology as Go modules, the Go
command accepts the `.v2` in `gopkg.in/yaml.v2` as a valid major version suffix.
This is a special case for compatibility with gopkg.in: modules hosted at other
domains need a slash suffix like `/v2`.

## Major version strategies

The recommended strategy is to develop `v2+` modules in a directory named after
the major version suffix.

	github.com/googleapis/gax-go @ master branch
	/go.mod    → module github.com/googleapis/gax-go
	/v2/go.mod → module github.com/googleapis/gax-go/v2

This approach is compatible with tools that aren't aware of modules: file paths
within the repository match the paths expected by `go get` in `GOPATH` mode.
This strategy also allows all major versions to be developed together in
different directories.

Other strategies may keep major versions on separate branches. However, if
`v2+` source code is on the repository's default branch (usually `master`),
tools that are not version-aware — including the `go` command in `GOPATH` mode
— may not distinguish between major versions.

The examples in this post will follow the major version subdirectory strategy,
since it provides the most compatibility. We recommend that module authors
follow this strategy as long as they have users developing in `GOPATH` mode.

## Publishing v2 and beyond

This post uses `github.com/googleapis/gax-go` as an example:

	$ pwd
	/tmp/gax-go
	$ ls
	CODE_OF_CONDUCT.md  call_option.go  internal
	CONTRIBUTING.md     gax.go          invoke.go
	LICENSE             go.mod          tools.go
	README.md           go.sum          RELEASING.md
	header.go
	$ cat go.mod
	module github.com/googleapis/gax-go

	go 1.9

	require (
		github.com/golang/protobuf v1.3.1
		golang.org/x/exp v0.0.0-20190221220918-438050ddec5e
		golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3
		golang.org/x/tools v0.0.0-20190114222345-bf090417da8b
		google.golang.org/grpc v1.19.0
		honnef.co/go/tools v0.0.0-20190102054323-c2f93a96b099
	)
	$

To start development on `v2` of `github.com/googleapis/gax-go`, we'll create a
new `v2/` directory and copy our package into it.

	$ mkdir v2
	$ cp -v *.go v2
	'call_option.go' -> 'v2/call_option.go'
	'gax.go' -> 'v2/gax.go'
	'header.go' -> 'v2/header.go'
	'invoke.go' -> 'v2/invoke.go'
	$

Now, let's create a v2 `go.mod` file by copying the current `go.mod` file and
adding a `/v2` suffix to the module path:

	$ cp go.mod v2/go.mod
	$ go mod edit -module github.com/googleapis/gax-go/v2 v2/go.mod
	$

Note that the `v2` version is treated as a separate module from the `v0 / v1`
versions: both may coexist in the same build. So, if your `v2+` module has
multiple packages, you should update them to use the new `/v2` import path:
otherwise, your `v2+` module will depend on your `v0 / v1` module. For example,
to update all `github.com/my/project` references to `github.com/my/project/v2`,
you can use `find` and `sed`:

	$ find . -type f \
		-name '*.go' \
		-exec sed -i -e 's,github.com/my/project,github.com/my/project/v2,g' {} \;
	$

Now we have a `v2` module, but we want to experiment and make changes before
publishing a release. Until we release `v2.0.0` (or any version without a
pre-release suffix), we can develop and make breaking changes as we decide on
the new API. If we want users to be able to experiment with the new API before
we officially make it stable, we can publish a `v2` pre-release version:

	$ git tag v2.0.0-alpha.1
	$ git push origin v2.0.0-alpha.1
	$

Once we are happy with our `v2` API and are sure we don't need any other breaking
changes, we can tag `v2.0.0`:

	$ git tag v2.0.0
	$ git push origin v2.0.0
	$

At that point, there are now two major versions to maintain. Backwards
compatible changes and bug fixes will lead to new minor and patch releases
(for example, `v1.1.0`, `v2.0.1`, etc.).

## Conclusion

Major version changes result in development and maintenance overhead and
require investment from downstream users to migrate. The larger the project,
the larger these overheads tend to be. A major version change should only come
after identifying a compelling reason. Once a compelling reason has been
identified for a breaking change, we recommend developing multiple major
versions in the master branch because it is compatible with a wider variety of
existing tools.

Breaking changes to a `v1+` module should always happen in a new, `vN+1` module.
When a new module is released, it means additional work for the maintainers and
for the users who need to migrate to the new package. Maintainers should
therefore validate their APIs before making a stable release, and consider
carefully whether breaking changes are really necessary beyond `v1`.
