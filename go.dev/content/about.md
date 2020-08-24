---
title: 'About'
date: 2019-06-25T17:51:23-04:00
draft: false
---

[Go.dev](https://go.dev) is a companion website to [golang.org](https://golang.org). [Golang.org](https://golang.org) is the home of the open source project and distribution, while [go.dev](https://go.dev) is the hub for Go users providing centralized and curated resources from across the Go ecosystem.

{{% gopher gopher=pink align=right %}}
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

Pkg.go.dev displays the GOOS and GOARCH for the documentation displayed at the bottom of the documentation page.

## Creating a badge

The pkg.go.dev badge provides a way for Go users to learn about the pkg.go.dev page associated with a given Go package or module. You can create a badge using the [badge generation tool](https://pkg.go.dev/badge). The tool will generate html and markdown snippets that you can use on your project website or in a README file.

[![PkgGoDev](https://pkg.go.dev/badge/golang.org/x/pkgsite)](https://pkg.go.dev/golang.org/x/pkgsite)

## Keyboard Shortcuts

There are keyboard shortcuts for navigating package documentation pages. Type '?' on a package page for help.

## Bookmarklet

The pkg.go.dev bookmarklet navigates from pages on source code hosts, such as GitHub, Bitbucket, Launchpad etc., to the package documentation. To install the bookmarklet, click and drag the following link to your bookmark bar: <a href="javascript:(function(){ const pathRegex = window.location.pathname.match(/([^\/]+)(?:\/([^\/]+))?/); const host = window.location.hostname; if (pathRegex) { window.location='https://pkg.go.dev/'+host+'/'+pathRegex[0]; } else { alert('There was an error navigating to pkg.go.dev!'); } })()">Pkg.go.dev Doc</a>

## License policy
Information for a given package or module may be limited if we are not able to detect a suitable license. See our [license policy](https://pkg.go.dev/license-policy) for more information.

## Feedback

Share your ideas, feature requests, and bugs on the [Go Issue Tracker](https://golang.org/s/discovery-feedback) For questions, please post on the #tools slack channel on the [Gophers Slack](https://invite.slack.golangbridge.org/), or email the [golang-dev mailing list](https://groups.google.com/group/golang-dev).
