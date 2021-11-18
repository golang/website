---
title: Next steps for pkg.go.dev
date: 2020-01-31
by:
- Julie Qiu
summary: What the Go team is planning for pkg.go.dev in 2020.
---

## Introduction

In 2019, we launched [go.dev](/), a new hub for Go
developers.

As part of the site, we also launched [pkg.go.dev](https://pkg.go.dev), a
central source of information about Go packages and modules. Like
[godoc.org](https://godoc.org), pkg.go.dev serves Go
documentation. However, it also understands modules and has information about
past versions of a package!

Throughout this year, we will be adding features to
[pkg.go.dev](https://pkg.go.dev) to help our users better understand their
dependencies and help them make better decisions around what libraries to
import.

## Redirecting godoc.org requests to pkg.go.dev

To minimize confusion about which site to use, later this year we are planning
to redirect traffic from [godoc.org](https://godoc.org) to the corresponding
page on [pkg.go.dev](https://pkg.go.dev). We need your help to ensure that
pkg.go.dev addresses all of our users' needs. We encourage everyone to begin
using pkg.go.dev today for all of their needs and provide feedback.

Your feedback will inform our transition plan, with the goal of making
[pkg.go.dev](https://pkg.go.dev) our primary source of information and
documentation for packages and modules. We’re sure there are things that you
want to see on pkg.go.dev, and we want to hear from you
about what those features are!

You can share your feedback with us on these channels:

  - Post on the [Go issue tracker](/s/discovery-feedback).
  - Email [go-discovery-feedback@google.com](mailto:go-discovery-feedback@google.com).
  - Click “Share Feedback” or “Report an Issue” in the go.dev footer.

As part of this transition, we will also be discussing plans for API access to
[pkg.go.dev](https://pkg.go.dev). We will be posting updates on
[Go issue 33654](/s/discovery-updates).

## Frequently asked questions

Since our launch in November, we’ve received tons of great feedback about
[pkg.go.dev](https://pkg.go.dev) from Go users. For the remainder of this post, we thought it would
be helpful to answer some frequently asked questions.

### My package doesn’t show up on pkg.go.dev! How do I add it?

We monitor the [Go Module Index](https://index.golang.org/index) regularly for
new packages to add to [pkg.go.dev](https://pkg.go.dev). If you don’t see a
package on pkg.go.dev, you can add it by fetching the module version from
[proxy.golang.org](https://proxy.golang.org). See
[go.dev/about](/about) for instructions.

### My package has license restrictions. What’s wrong with it?

We understand it can be a frustrating experience to not be able to see the
package you want in its entirety on [pkg.go.dev](https://pkg.go.dev). We
appreciate your patience as we improve our license detection algorithm.

Since our launch in November, we've made the following improvements:

  - Updated our [license policy](https://pkg.go.dev/license-policy) to include the list of licenses that we detect and recognize
  - Worked with the [licensecheck](https://github.com/google/licensecheck) team to improve detection for copyright notices
  - Established a manual review process for special cases

As always, our license policy is at
[pkg.go.dev/license-policy](https://pkg.go.dev/license-policy). If you are
having issues, feel free to [file an issue on the Go issue tracker](/s/discovery-feedback), or email
[go-discovery-feedback@google.com](mailto:go-discovery-feedback@google.com)
so that we can work with you directly!

### Will pkg.go.dev be open-sourced so I can run it at work for my private code?

We understand that corporations with private code want to run a documentation
server that provides module support. We want to help meet that need, but we
feel we don’t yet understand it as well as we need to.

We’ve heard from users that running the [godoc.org](https://godoc.org) server
is more complex than it should be, because it is designed for serving at public
internet scale instead of just within a company. We believe the current
[pkg.go.dev](https://pkg.go.dev) server would have the same problem.

We think a new server is more likely to be the right answer for use with
private code, instead of exposing every company to the complexity of running
the internet-scale [pkg.go.dev](https://pkg.go.dev) codebase. In addition to
serving documentation, a new server could also serve information to
[goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports?tab=doc) and
[gopls](https://pkg.go.dev/golang.org/x/tools/gopls).

If you want to run such a server, please fill out this
[**3-5 minute survey**](https://google.qualtrics.com/jfe/form/SV_6FHmaLveae6d8Bn)
to help us better understand your needs. This survey will be available until
March 1st, 2020.

We’re excited about the future of [pkg.go.dev](https://pkg.go.dev) in 2020,
and we hope you are too! We look forward to hearing your feedback and working
with the Go community on this transition.
