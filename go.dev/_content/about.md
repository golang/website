---
title: 'About'
layout: article
date: 2019-06-25T17:51:23-04:00
---

[Go.dev](https://go.dev) is a companion website to [golang.org](https://golang.org). [Golang.org](https://golang.org) is the home of the open source project and distribution, while [go.dev](https://go.dev) is the hub for Go users providing centralized and curated resources from across the Go ecosystem.

{{gopher `
  color: pink
  align: right
`}}
Go.dev provides:

1. Centralized information for Go packages and modules published on index.golang.org.
2. Essential learning resources
3. Critical use cases & case studies

Go.dev is currently in [MVP](https://en.wikipedia.org/wiki/Minimum_viable_product) status. We're proud of what we've built and excited to share it with the community. We hope you find value and joy in using go.dev. Go.dev only has a small portion of features we intend to build, and we are actively seeking feedback. If you have any ideas, suggestions or issues, please let us know.

## Adding a package
Data for the site is downloaded from [proxy.golang.org](https://proxy.golang.org/). We monitor the [Go Module Index](https://index.golang.org/index) regularly for new packages to add to pkg.go.dev. If you don’t see a package on pkg.go.dev, you can add it by doing one of the following:

* Visiting that page on pkg.go.dev, and clicking the "Request" button. For example: <br /> https://<span></span>pkg.go.dev/example.com/my/module

*  Making a request to proxy.golang.org for the module version, to any endpoint specified by the [Module proxy protocol](https://golang.org/cmd/go/#hdr-Module_proxy_protocol). For example: <br /> https://<span></span>proxy.golang.org/example.com/my/module/@v/v1.0.0.info

*  Downloading the package via the [go command](https://golang.org/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them). For example:  <br /> GOPROXY=https://<span></span>proxy.golang.org GO111MODULE=on go get example.com/my/module@v1.0.0

## Removing a package
If you are the author of a package and would like to have it removed from pkg.go.dev, please [file an issue](https://golang.org/s/pkgsite-feedback) on the Go Issue Tracker with the path that you want to remove.

Note that we can only remove a module entirely from the site. We cannot remove it just for specific versions.

## Documentation

Documentation is generated based on Go source code downloaded from the Go Module Mirror at `proxy.golang.org/<module>/@v/<version>.zip`. New module versions are fetched from index.golang.org and added to pkg.go.dev site every few minutes.

The [guidelines for writing documentation](https://blog.golang.org/godoc) for the godoc tool apply to pkg.go.dev.

It’s important to write a good summary of the package in the first sentence of the package comment. The go.dev site indexes the first sentence and displays it in search results.

### Build Context

Most Go packages look and behave the same regardless of the machine architecture
or operating system. But some have different documentation, even different
exported symbols, for different architectures or OSes. Some packages may not even
exist for some architectures.

Go calls an OS/architecture pair a "build context" and writes it with a slash,
like `linux/amd64`. You may also see the terms `GOOS` and `GOARCH` for the OS
and architecture respectively, because those are the names of the environment
variables that the go command uses. (See the [go command
documentation](https://golang.org/cmd/go) for more information.)

If a package exists at only one build context, pkg.go.dev displays that build
context at the upper right corner of the documentation. For example,
https://pkg.go.dev/syscall/js displays "js/wasm".

If a package is different in different build contexts, then pkg.go.dev will
display one by default and provide a dropdown control at the upper right so you
can select a different one.

For packages that are the same across all build contexts, pkg.go.dev does not
display any build context information.

Although there are many possible OS/architecture pairs, pkg.go.dev considers
only a
[handful](https://go.googlesource.com/pkgsite/+/master/internal/build_context.go#29)
of them. So if a package only exists for unsupported build contexts, pkg.go.dev
will not display documentation for it.

### Source Links

Most of the time, pkg.go.dev can determine the location of a package's source
files, and provide links from symbols in the documentation to their definitions
in the source. If your package's source is not linked, try one of the following
two approaches.

If pkg.go.dev finds a `go-source` meta tag on your site that follows the
[specified format](https://github.com/golang/gddo/wiki/Source-Code-Links), it
can often determine the right links, even though the format doesn't take
versioning into account.

If that doesn't work, you will need to add your repo or code-hosting site to
pkg.go.dev's list of patterns (see  [Go Issue 40477](https://golang.org/issues/40477) for context).
Read about how to [contribute to pkg.go.dev](https://go.googlesource.com/pkgsite#contributing),
then produce a CL that adds a pattern to the
[`internal/source`](https://go.googlesource.com/pkgsite/+/refs/heads/master/internal/source/source.go)
package.

## Best practices

Pkg.go.dev surfaces details about Go packages and modules in order to help provide guidelines for best practices with Go.

Here are the details we surface:

* Has go.mod file
  * The Go module system was introduced in Go 1.11 and is the official dependency management solution for Go. A module version is defined by a tree of source files, with a go.mod file in its root. [More information about the go.mod file](https://golang.org/cmd/go/#hdr-The_go_mod_file).

* Redistributable license
  * Redistributable licenses place minimal restrictions on how software can be used, modified, and redistributed. For more information on how pkg.go.dev determines if a license is redistributable, see our [license policy](http://pkg.go.dev/license-policy).

* Tagged version
  * When the go get command resolves modules by default it prioritizes tagged versions. When no tagged versions exist, go get looks up the latest known commit. Modules with tagged versions give importers more predictable builds. See [semver.org](https://semver.org) and [Keeping Your Modules Compatible](https://blog.golang.org/module-compatibility) for more information.

* Stable version
  * Projects at v0 are assumed to be experimental. When a project reaches a stable version — major version v1 or higher — breaking changes must be done in a new major version. Stable versions give developers the confidence that breaking changes won’t occur when they upgrade a package to the latest minor version. See [Go Modules: v2 and Beyond](https://blog.golang.org/v2-go-modules) for more information.

## Creating a badge

The pkg.go.dev badge provides a way for Go users to learn about the pkg.go.dev page associated with a given Go package or module. You can create a badge using the [badge generation tool](https://pkg.go.dev/badge). The tool will generate html and markdown snippets that you can use on your project website or in a README file.

[![PkgGoDev](https://pkg.go.dev/badge/golang.org/x/pkgsite)](https://pkg.go.dev/golang.org/x/pkgsite)

## Adding links

You can add links to your README files and package documentation that will be
shown on the right side of the pkg.go.dev page. For details, see [this
issue](https://golang.org/issue/42968).

## Keyboard Shortcuts

There are keyboard shortcuts for navigating package documentation pages. Type '?' on a package page for help.

## Bookmarklet

The pkg.go.dev bookmarklet navigates from pages on source code hosts, such as GitHub, Bitbucket, Launchpad etc., to the package documentation. To install the bookmarklet, click and drag the following link to your bookmark bar: <a href="javascript:(function(){ const pathRegex = window.location.pathname.match(/([^\/]+)(?:\/([^\/]+))?/); const host = window.location.hostname; if (pathRegex) { window.location='https://pkg.go.dev/'+host+'/'+pathRegex[0]; } else { alert('There was an error navigating to pkg.go.dev!'); } })()">Pkg.go.dev Doc</a>

## License policy
Information for a given package or module may be limited if we are not able to detect a suitable license. See our [license policy](https://pkg.go.dev/license-policy) for more information.

## Feedback

Share your ideas, feature requests, and bugs on the [Go Issue Tracker](https://golang.org/s/discovery-feedback) For questions, please post on the #tools slack channel on the [Gophers Slack](https://invite.slack.golangbridge.org/), or email the [golang-dev mailing list](https://groups.google.com/group/golang-dev).
