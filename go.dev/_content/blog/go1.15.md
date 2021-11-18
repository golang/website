---
title: Go 1.15 is released
date: 2020-08-11
by:
- Alex Rakoczy
summary: Go 1.15 adds a new linker, X.509 changes, runtime improvements, compiler improvements, GOPROXY improvements, and more.
---


Today the Go team is very happy to announce the release of Go 1.15. You can get it from the [download page](/dl).

Some of the highlights include:

  - [Substantial improvements to the Go linker](/doc/go1.15#linker)
  - [Improved allocation for small objects at high core counts](/doc/go1.15#runtime)
  - [X.509 CommonName deprecation](/doc/go1.15#commonname)
  - [GOPROXY supports skipping proxies that return errors](/doc/go1.15#go-command)
  - [New embedded tzdata package](/doc/go1.15#time/tzdata)
  - [A number of Core Library improvements](/doc/go1.15#library)

For the complete list of changes and more information about the improvements above, see the [**Go 1.15 release notes**](/doc/go1.15).

We want to thank everyone who contributed to this release by writing code, filing bugs, providing feedback, and/or testing the beta and release candidates.
Your contributions and diligence helped to ensure that Go 1.15 is as stable as possible.
That said, if you notice any problems, please [file an issue](/issue/new).

We hope you enjoy the new release!
