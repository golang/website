---
title: Six years of Go
date: 2015-11-10
by:
- Andrew Gerrand
summary: Happy 6th birthday, Go!
---


Six years ago today the Go language was released as an open source project.
Since then, more than 780 contributors have made over 30,000 commits to the
project's 22 repositories. The ecosystem continues to grow, with GitHub
reporting more than 90,000 Go repositories. And, offline, we see new Go events
and user groups pop up [around](https://blog.golang.org/gophercon2015)
[the](https://blog.golang.org/gouk15)
[world](https://blog.golang.org/gopherchina) with regularity.

{{image "6years/6years-gopher.png"}}

In August we [released Go 1.5](https://blog.golang.org/go1.5), the most
significant release since Go 1. It features a completely
[redesigned garbage collector](/doc/go1.5#gc) that makes
the language more suitable for latency-sensitive applications; it marks the
transition from a C-based compiler tool chain to one
[written entirely in Go](/doc/go1.5#c); and it includes
ports to [new architectures](/doc/go1.5#ports), with better
support for ARM processors (the chips that power most smartphones).
These improvements make Go better suited to a broader range of tasks, a trend
that we hope will continue over the coming years.

Improvements to tools continue to boost developer productivity.
We introduced the [execution tracer](/cmd/trace/) and the
"[go doc](/cmd/go/#hdr-Show_documentation_for_package_or_symbol)"
command, as well as more enhancements to our various
[static analysis tools](/talks/2014/static-analysis.slide).
We are also working on an
[official Go plugin for Sublime Text](https://groups.google.com/forum/#!topic/Golang-nuts/8oCSjAiKXUQ),
with better support for other editors in the pipeline.

Early next year we will release more improvements in Go 1.6, including
HTTP/2 support for [net/http](/pkg/net/http/) servers and
clients, an official package vendoring mechanism, support for blocks in text
and HTML templates, a memory sanitizer that checks both Go and C/C++ code, and
the usual assortment of other improvements and fixes.

This is the sixth time we have had the pleasure of writing a birthday blog post
for Go, and we would not be doing so if not for the wonderful and passionate
people in our community. The Go team would like to thank everyone who has
contributed code, written an open source library, authored a blog post, helped
a new gopher, or just given Go a try. Without you, Go would not be as complete,
useful, or successful as it is today. Thank you, and celebrate!
