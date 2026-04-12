---
title: Go 1.25 is released
date: 2025-08-12
by:
- Dmitri Shuralyov, on behalf of the Go team
summary: Go 1.25 adds container-aware GOMAXPROCS, testing/synctest package, experimental GC, experimental encoding/json/v2, and more.
---

Today the Go team is pleased to release Go 1.25.
You can find its binary archives and installers on the [download page](/dl/).

Go 1.25 comes with improvements over Go 1.24 across
its [tools](/doc/go1.25#tools),
the [runtime](/doc/go1.25#runtime),
[compiler](/doc/go1.25#compiler),
[linker](/doc/go1.25#linker),
and the [standard library](/doc/go1.25#library),
including the addition of one [new package](/doc/go1.25#new-testingsynctest-package).
There are [port-specific](/doc/go1.25#ports) changes
and [`GODEBUG` settings](/doc/godebug#go-125) updates.

Some of the additions in Go 1.25 are in an experimental stage
and become exposed only when you explicitly opt in.
Notably, a [new experimental garbage collector](/doc/go1.25#new-experimental-garbage-collector),
and a [new experimental `encoding/json/v2` package](/doc/go1.25#json_v2)
are available for you to try ahead of time and provide your feedback.
It really helps if you're able to do that!

Please refer to the [Go 1.25 Release Notes](/doc/go1.25) for the complete list
of additions, changes and improvements in Go 1.25.

Over the next few weeks, follow-up blog posts will cover some of the topics
relevant to Go 1.25 in more detail. Check back in later to read those posts.

Thanks to everyone who contributed to this release by writing code, filing bugs,
trying out experimental additions, sharing feedback, and testing the release candidates.
Your efforts helped make Go 1.25 as stable as possible.
As always, if you notice any problems, please [file an issue](/issue/new).

We hope you enjoy using the new release!
