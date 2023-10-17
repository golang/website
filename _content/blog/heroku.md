---
title: Go at Heroku
date: 2011-04-21
by:
- Keith Rarick
- Blake Mizerany
tags:
- guest
summary: Two Heroku system engineers discuss their experiences using Go.
---


_This week’s blog post is written by_ [_Keith Rarick_](http://xph.us/)
_and_ [_Blake Mizerany_](http://itsbonus.heroku.com/),
_systems engineers at_ [Heroku](http://www.heroku.com/).
_In their own words, they "eat, drink, and sleep distributed systems." Here they discuss their experiences using Go._

A big problem that comes with building distributed systems is the coordination
of physical servers.
Each server needs to know various facts about the system as a whole.
This critical data includes locks, configuration data,
and so on, and it must be consistent and available even during data store failures,
so we need a data store with solid consistency guarantees.
Our solution to this problem is [Doozer](http://xph.us/2011/04/13/introducing-doozer.html),
a new, consistent, highly-available data store written in Go.

At Doozer's core is [Paxos](http://en.wikipedia.org/wiki/Paxos_(computer_science)),
a family of protocols for solving consensus in an unreliable network of unreliable nodes.
While Paxos is essential to running a fault-tolerant system,
it is notorious for being difficult to implement.
Even example implementations that can be found online are complex and hard to follow,
despite being simplified for educational purposes.
Existing production systems have a reputation for being worse.

Fortunately, Go's concurrency primitives made the task much easier.
Paxos is defined in terms of independent,
concurrent processes that communicate via passing messages.
In Doozer, these processes are implemented as goroutines,
and their communications as channel operations.
In the same way that garbage collectors improve upon malloc and free,
we found that [goroutines and channels](https://blog.golang.org/2010/07/share-memory-by-communicating.html)
improve upon the lock-based approach to concurrency.
These tools let us avoid complex bookkeeping and stay focused on the problem at hand.
We are still amazed at how few lines of code it took to achieve something
renowned for being difficult.

The standard packages in Go were another big win for Doozer.
The Go team is very pragmatic about what goes into them.
For instance, a package we quickly found useful was [websocket](/pkg/websocket/).
Once we had a working data store, we needed an easy way to introspect it
and visualize activity.
Using the websocket package, Keith was able to add the web viewer on his
train ride home and without requiring external dependencies.
This is a real testament to how well Go mixes systems and application programming.

One of our favorite productivity gains was provided by Go's source formatter:
[gofmt](/cmd/gofmt/).
We never argued over where to put a curly-brace,
tabs vs. spaces, or if we should align assignments.
We simply agreed that the buck stopped at the default output from gofmt.

Deploying Doozer was satisfyingly simple.
Go builds statically linked binaries which means Doozer has no external dependencies;
it's a single file that can be copied to any machine and immediately launched
to join a cluster of running Doozers.

Finally, Go's maniacal focus on simplicity and orthogonality aligns with
our view of software engineering.
Like the Go team, we are pragmatic about what features go into Doozer.
We sweat the details, preferring to change an existing feature instead of
introducing a new one.
In this sense, Go is a perfect match for Doozer.

We already have future projects in mind for Go. Doozer is just the start of much bigger system.
