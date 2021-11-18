---
title: Pkg.go.dev is open source!
date: 2020-06-15
by:
- Julie Qiu
---


We’re excited to announce that the codebase for
[pkg.go.dev](https://pkg.go.dev) is now open source.

The repository lives at
[go.googlesource.com/pkgsite](https://go.googlesource.com/pkgsite)
and is mirrored to
[github.com/golang/pkgsite](https://github.com/golang/pkgsite).
We will continue using the Go issue tracker to track
[feedback](https://github.com/golang/go/labels/go.dev)
related to pkg.go.dev.

## Contributing

If you are interested in contributing to any
[issues related to pkg.go.dev](https://github.com/golang/go/labels/go.dev),
check out our
[contribution guidelines](https://go.googlesource.com/pkgsite/+/refs/heads/master/CONTRIBUTING.md).
We also encourage you to continue
[filing issues](/s/discovery-feedback)
if you run into problems or have feedback.

## What’s Next

We really appreciate all the feedback we’ve received so far. It has been a big
help in shaping our
[roadmap](https://go.googlesource.com/pkgsite#roadmap) for the coming year.
Now that pkg.go.dev is open source, here’s what we’ll be working on next:

- We have some design changes planned for pkg.go.dev,
  to address
  [UX feedback](https://github.com/golang/go/issues?q=is%3Aissue+is%3Aopen+label%3Ago.dev+label%3AUX)
  that we have received. You can expect a more cohesive search and navigation
  experience. We plan to share these designs for feedback once they are ready.

- We know that there are features available on godoc.org that users
  want to see on pkg.go.dev. We’ve been keeping track of them on
  [Go issue #39144](/issue/39144),
  and will prioritize adding them in the next few months. We also plan to
  continue improving our license detection algorithm based on feedback.

- We’ll be improving our search experience based on feedback in
  [Go issue #37810](/issue/37810),
  to make it easier for users to find the dependencies they are looking for and
  make better decisions around which ones to import.

Thanks for being patient with us in the process of open sourcing pkg.go.dev.
We’re looking forward to receiving your contributions and working with you on
the future of the project.
