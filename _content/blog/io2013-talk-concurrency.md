---
title: Advanced Go Concurrency Patterns
date: 2013-05-23
by:
- Andrew Gerrand
tags:
- talk
- video
- concurrency
summary: Watch Sameer Ajmani's talk, “Advanced Go Concurrency Patterns,” from Google I/O 2013.
---


At Google I/O a year ago Rob Pike presented [_Go Concurrency Patterns_](/talks/2012/concurrency.slide),
an introduction to Go's concurrency model.
Last week, at I/O 2013, Go team member Sameer Ajmani continued the story
with [_Advanced Go Concurrency Patterns_](http://go.dev/talks/2013/advconc.slide),
an in-depth look at a real concurrent programming problem.
The talk shows how to detect and avoid deadlocks and race conditions,
and demonstrates the implementation of deadlines,
cancellation, and more.
For those who want to take their Go programming to the next level, this is a must-see.

{{video "https://www.youtube.com/embed/QDDwwePbDtw?rel=0"}}

The slides are [available here](/talks/2013/advconc.slide)
(use the left and right arrows to navigate).

The slides were produced with [the present tool](https://godoc.org/golang.org/x/tools/present),
and the runnable code snippets are powered by the [Go Playground](http://play.golang.org/).
The source code for this talk is in [the go.talks sub-repository](https://github.com/golang/talks/tree/master/content/2013/advconc).
