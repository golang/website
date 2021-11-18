---
title: Redirecting godoc.org requests to pkg.go.dev
date: 2020-12-15
by:
- Julie Qiu
summary: The plan for moving from godoc.org to pkg.go.dev.
---


With the introduction of Go modules and the growth of the Go ecosystem,
[pkg.go.dev](https://pkg.go.dev) was
[launched in 2019](https://blog.golang.org/go.dev) to provide a central place
where developers can discover and evaluate Go packages and modules. Like
godoc.org, pkg.go.dev serves Go documentation, but it also supports modules,
better search functionality, and signals to help Go users to find the right
packages.

As [we shared in January 2020](https://blog.golang.org/pkg.go.dev-2020), our
goal is to eventually redirect traffic from godoc.org to the corresponding page
on pkg.go.dev. We’ve also made it possible for users to opt in to redirecting
their own requests from godoc.org to pkg.go.dev.

We’ve received a lot of great feedback this year, which has been tracked and
resolved through the
[pkgsite/godoc.org-redirect](https://github.com/golang/go/milestone/157?closed=1)
and [pkgsite/design-2020](https://github.com/golang/go/milestone/159?closed=1)
milestones on the Go issue tracker. Your feedback resulted in support for
popular feature requests on pkg.go.dev,
[open sourcing pkgsite](https://blog.golang.org/pkgsite), and most recently, a
[redesign of pkg.go.dev](https://blog.golang.org/pkgsite-redesign).

## Next Steps

The next step in this migration is to redirect all requests from godoc.org to
the corresponding page on pkg.go.dev.

This will happen in early 2021, once the work tracked at the
[pkgsite/godoc.org-redirect milestone](https://github.com/golang/go/milestone/157)
is complete.

During this migration, updates will be posted to
[Go issue 43178](/issue/43178).

We encourage everyone to begin using pkg.go.dev today. You can do so by
visiting [godoc.org?redirect=on](https://godoc.org?redirect=on), or clicking
“Always use pkg.go.dev” in the top right corner of any godoc.org page.

## FAQs

**Will godoc.org URLs continue to work?**

Yes! We will redirect all requests arriving at godoc.org to the equivalent page
on pkg.go.dev, so all your bookmarks and links will continue to take you to the
documentation you need.

**What will happen to the golang/gddo repository?**

The [gddo repository](http://go.googlesource.com/gddo) will remain available
for anyone who wants to keep running it themselves, or even fork and improve
it. We will mark it archived to make clear that we will no longer accept
contributions. However, you will be able to continue forking the repository.

**Will api.godoc.org continue to work?**

This transition will have no impact on api.godoc.org. Until an API is available
for pkg.go.dev, api.godoc.org will continue to serve traffic. See
[Go issue 36785](/issue/36785) for updates on an API for
pkg.go.dev.

**Will my godoc.org badges keep working?**

Yes! Badge URLs will redirect to the equivalent URL on pkg.go.dev too. Your
page will automatically get a new pkg.go.dev badge. You can also generate a new
badge at [pkg.go.dev/badge](https://pkg.go.dev/badge) if you would like to
update your badge link.

## Feedback

As always, feel free to [file an issue](/s/pkgsite-feedback)
on the Go issue tracker for any feedback.

## Contributing

Pkg.go.dev is an [open source project](https://go.googlesource.com/pkgsite). If
you’re interested in contributing to the pkgsite project, check out the
[contribution guidelines](https://go.googlesource.com/pkgsite/+/refs/heads/master/CONTRIBUTING.md)
or join the
[#pkgsite channel](https://gophers.slack.com/messages/pkgsite) on Gophers Slack
to learn more.
