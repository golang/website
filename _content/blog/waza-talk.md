---
title: Concurrency is not parallelism
date: 2013-01-16
by:
- Andrew Gerrand
tags:
- concurrency
- talk
- video
summary: Watch Rob Pike's talk, _Concurrency is not parallelism._
---


If there's one thing most people know about Go,
is that it is designed for concurrency.
No introduction to Go is complete without a demonstration of its goroutines and channels.

But when people hear the word _concurrency_ they often think of _parallelism_,
a related but quite distinct concept.
In programming, concurrency is the _composition_ of independently executing processes,
while parallelism is the simultaneous _execution_ of (possibly related) computations.
Concurrency is about _dealing with_ lots of things at once.
Parallelism is about _doing_ lots of things at once.

To clear up this conflation, Rob Pike gave a talk at [Heroku](http://heroku.com/)'s
Waza conference entitled
[_Concurrency is not parallelism_](https://blog.heroku.com/concurrency_is_not_parallelism),
and a video recording of the talk was released a few months ago.

{{video "https://www.youtube.com/embed/oV9rvDllKEg" 500 281}}

The slides are available at [go.dev/talks](/talks/2012/waza.slide)
(use the left and right arrow keys to navigate).

To learn about Go's concurrency primitives,
watch [Go concurrency patterns](http://www.youtube.com/watch?v=f6kdp27TYZs)
([slides](/talks/2012/concurrency.slide)).
