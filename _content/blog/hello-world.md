---
title: "Go: What's New in March 2010"
date: 2010-03-18
by:
- Andrew Gerrand
summary: First post!
---


Welcome to the official Go Blog. We, the Go team,
hope to use this blog to keep the world up-to-date on the development of
the Go programming language and the growing ecosystem of libraries and applications surrounding it.

It's been a few months since we launched (November last year),
so let's talk about what's been happening in Go World since then.

The core team at Google has continued to develop the language,
compilers, packages, tools, and documentation.
The compilers now produce code that is in some cases between 2x and an order
of magnitude faster than at release.
We have put together some graphs of a selection of [Benchmarks](http://godashboard.appspot.com/benchmarks),
and the [Build Status](http://godashboard.appspot.com/) page tracks the
reliability of each changeset submitted to the repository.

We have made syntax changes to make the language more concise,
regular, and flexible.
Semicolons have been [almost entirely removed](http://groups.google.com/group/golang-nuts/t/5ee32b588d10f2e9) from the language.
The [...T syntax](/doc/go_spec.html#Function_types)
makes it simpler to handle an arbitrary number of typed function parameters.
The syntax x[lo:] is now shorthand for x[lo:len(x)].
Go also now natively supports complex numbers.
See the [release notes](/doc/devel/release.html) for more.

[Godoc](/cmd/godoc/) now provides better support for
third-party libraries,
and a new tool - [goinstall](/cmd/goinstall) - has been
released to make it easy to install them.
Additionally, we've started working on a package tracking system to make
it easier to find what you need.
You can view the beginnings of this on the [Packages page](http://godashboard.appspot.com/package).

More than 40,000 lines of code have been added to [the standard library](/pkg/),
including many entirely new packages, a sizable portion written by external contributors.

Speaking of third parties, since launch a vibrant community has flourished
on our [mailing list](http://groups.google.com/group/golang-nuts/) and
irc channel (#go-nuts on freenode).
We have officially added more than 50 people to the project.
Their contributions range from bug fixes and documentation corrections to
core packages and support for additional operating systems (Go is now supported under FreeBSD,
and a [Windows port](http://code.google.com/p/go/wiki/WindowsPort) is underway).
We regard these community contributions our greatest success so far.

We've received some good reviews, too.  This [recent article in PC World](http://www.pcworld.idg.com.au/article/337773/google_go_captures_developers_imaginations/)
summarized the enthusiasm surrounding the project.
Several bloggers have begun documenting their experiences in the language
(see [here](http://golang.tumblr.com/),
[here](http://www.infi.nl/blog/view/id/47),
and [here](http://freecella.blogspot.com/2010/01/gospecify-basic-setup-of-projects.html)
for example)  The general reaction of our users has been very positive;
one first-timer remarked ["I came away extremely impressed. Go walks an elegant line between simplicity and power."](https://groups.google.com/group/golang-nuts/browse_thread/thread/5fabdd59f8562ed2)

As to the future: we have listened to the myriad voices telling us what they need,
and are now focused on getting Go ready for the prime time.
We are improving the garbage collector, runtime scheduler,
tools, and standard libraries, as well as exploring new language features.
2010 will be an exciting year for Go, and we look forward to collaborating
with the community to make it a successful one.
