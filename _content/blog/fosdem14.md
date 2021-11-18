---
title: Go talks at FOSDEM 2014
date: 2014-02-24
by:
- Andrew Gerrand
tags:
- fosdem
- youtube
- talk
summary: Reporting from the Go Devroom at FOSDEM 2014.
---

## Introduction

At [FOSDEM](http://fosdem.org/) on the 2nd of February 2014 members of the Go
community presented a series of talks in the Go Devroom. The day was a huge
success, with 13 great talks presented to a consistently jam-packed room.

Video recordings of the talks are now available, and a selection of these
videos are presented below.

The complete series of talks is available
[as a YouTube playlist](http://www.youtube.com/playlist?list=PLtLJO5JKE5YDKG4WcaNts3IVZqhDmmuBH).
(You can also get them directly at the
[FOSDEM video archive](http://video.fosdem.org/2014/K4601/Sunday/).)

## Scaling with Go: YouTube's Vitess

Google Engineer Sugu Sougoumarane described how he and his
team built [Vitess](https://github.com/youtube/vitess) in Go to help scale
[YouTube](https://youtube.com).

Vitess is a set of servers and tools primarily developed in Go.
It helps scale MySQL databases for the web, and is currently used as a
fundamental component of YouTube's MySQL infrastructure.

The talk covers some history about how and why the team chose Go, and how it
paid off.
Sugu also talks abou tips and techniques used to scale Vitess using Go.

{{video "https://www.youtube.com/embed/qATTTSg6zXk"}}

The slides for the talk are [available here](https://github.com/youtube/vitess/blob/master/doc/Vitess2014.pdf?raw=true).

## Camlistore

[Camlistore](http://camlistore.org/) is designed to be "your personal storage
system for life, putting you in control, and designed to last." It's open
source, under nearly 4 years of active development, and extremely flexible.  In
this talk, Brad Fitzpatrick and Mathieu Lonjaret explain why they built it,
what it does, and talk about its design.

{{video "https://www.youtube.com/embed/yvjeIZgykiA"}}

## Write your own Go compiler

Elliot Stoneham explains the potential for Go as a portable language and
reviews the Go tools that make that such an exciting possibility.

He said: "Based on my experiences writing an experimental Go to Haxe
translator, I'll talk about the practical issues of code generation and runtime
emulation required. I'll compare some of my design decisions with those of two
other Go compiler/translators that build on the go.tools library. My aim is to
encourage you to try one of these new 'mutant' Go compilers. I hope some of you
will be inspired to contribute to one of them or even to write a new one of
your own."

{{video "https://www.youtube.com/embed/Qe8Dq7V3hXY"}}

## More

There were many more great talks, so please check out the complete series
[as a YouTube playlist](http://www.youtube.com/playlist?list=PLtLJO5JKE5YDKG4WcaNts3IVZqhDmmuBH).
In particular, the [lightning talks](http://www.youtube.com/watch?v=cwpI5ONWGxc&list=PLtLJO5JKE5YDKG4WcaNts3IVZqhDmmuBH&index=7) were a lot of fun.

I would like to give my personal thanks to the excellent speakers, Mathieu
Lonjaret for managing the video gear, and to the FOSDEM staff for making all
this possible.
