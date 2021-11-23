---
title: "Tidying up the Go web experience"
date: 2021-08-18
by:
- Russ Cox
summary: Consolidating our web sites onto go.dev.
---

In 2019, which seems like a decade ago, we [launched go.dev](/blog/go.dev),
a new hub for Go developers, along with the [companion site pkg.go.dev](https://pkg.go.dev/),
providing information about Go packages and modules.

The go.dev web site contains useful information for people evaluating Go,
but golang.org continued to serve distribution downloads, documentation,
and a package reference for the standard library.
Other sites — blog.golang.org, play.golang.org, talks.golang.org,
and tour.golang.org — hold additional material.
It's all a bit fragmented and confusing.

Over the next month or two we will be merging
the golang.org sites into
a single coherent web presence, here on go.dev.
You may have already noticed that links to the package reference docs
for the standard library on golang.org/pkg now redirect to
their [equivalents on pkg.go.dev](https://pkg.go.dev/std),
which is a better experience today and will continue to improve.
As the next step, the Go blog has moved to go.dev/blog,
starting with the post you are reading right now.
(Of course, all the old blog posts are here too.)

As we move the content to its new home on go.dev,
rest assured that all existing URLs will redirect to their new homes:
no links will be broken.

We are excited to have a single coherent web site
where everyone can find what they need to know about Go.
It's a small detail, but one long overdue.

If you have any ideas or suggestions, or you run into problems,
please let us know via the “Report an Issue” link at the bottom of every page.
Thanks!
