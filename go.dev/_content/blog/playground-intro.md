---
title: Introducing the Go Playground
date: 2010-09-15
by:
- Andrew Gerrand
tags:
- playground
summary: "Announcing the Go Playground, https://play.golang.org/."
---


If you visit [golang.org](/) today you'll see our new look.
We have given the site a new coat of paint and reorganized its content to
make it easier to find.
These changes are also reflected in the web interface of [godoc](/cmd/godoc/),
the Go documentation tool.
But the real news is a prominent new feature: the [Go Playground](/).

{{image "playground-intro/screenshot.png"}}

The Playground allows anyone with a web browser to write Go code that we
immediately compile,
link, and run on our servers.
There are a few example programs to get you started (see the "Examples" drop-down).
We hope that this will give curious programmers an opportunity to try the
language before [installing it](/doc/install.html),
and experienced Go users a convenient place in which to experiment.
Beyond the front page, this functionality has the potential to make our
reference and tutorial materials more engaging.
We hope to extend its use in the near future.

Of course, there are some limitations to the kinds of programs you can run in the Playground.
We can't simply accept arbitrary code and run it on our servers without restrictions.
The programs build and run in a sandbox with a reduced standard library;
the only communication your program has to the outside world is via standard output,
and there are limits to CPU and memory use.
As such, consider this just a taste of the wonderful world of Go;
to have the full experience you'll need to [download it yourself](/doc/install.html).
If you've been meaning to try Go but never got around to it,
why not visit [golang.org](/) to try it right now?
