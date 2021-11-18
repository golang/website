---
title: Share Memory By Communicating
date: 2010-07-13
by:
- Andrew Gerrand
tags:
- concurrency
- technical
summary: A preview of the new Go codelab, Share Memory by Communicating.
---


Traditional threading models (commonly used when writing Java,
C++, and Python programs, for example) require the programmer to communicate
between threads using shared memory.
Typically, shared data structures are protected by locks,
and threads will contend over those locks to access the data.
In some cases, this is made easier by the use of thread-safe data structures
such as Python's Queue.

Go's concurrency primitives - goroutines and channels - provide an elegant
and distinct means of structuring concurrent software.
(These concepts have an [interesting history](https://swtch.com/~rsc/thread/) that begins with C.
A. R. Hoare's [Communicating Sequential Processes](http://www.usingcsp.com/).)
Instead of explicitly using locks to mediate access to shared data,
Go encourages the use of channels to pass references to data between goroutines.
This approach ensures that only one goroutine has access to the data at a given time.
The concept is summarized in the document [Effective Go](/doc/effective_go.html)
(a must-read for any Go programmer):

_Do not communicate by sharing memory; instead, share memory by communicating._

Consider a program that polls a list of URLs.
In a traditional threading environment, one might structure its data like so:

	type Resource struct {
	    url        string
	    polling    bool
	    lastPolled int64
	}

	type Resources struct {
	    data []*Resource
	    lock *sync.Mutex
	}

And then a Poller function (many of which would run in separate threads) might look something like this:

{{raw `
	func Poller(res *Resources) {
	    for {
	        // get the least recently-polled Resource
	        // and mark it as being polled
	        res.lock.Lock()
	        var r *Resource
	        for _, v := range res.data {
	            if v.polling {
	                continue
	            }
	            if r == nil || v.lastPolled < r.lastPolled {
	                r = v
	            }
	        }
	        if r != nil {
	            r.polling = true
	        }
	        res.lock.Unlock()
	        if r == nil {
	            continue
	        }

	        // poll the URL

	        // update the Resource's polling and lastPolled
	        res.lock.Lock()
	        r.polling = false
	        r.lastPolled = time.Nanoseconds()
	        res.lock.Unlock()
	    }
	}
`}}

This function is about a page long, and requires more detail to make it complete.
It doesn't even include the URL polling logic (which,
itself, would only be a few lines), nor will it gracefully handle exhausting
the pool of Resources.

Let's take a look at the same functionality implemented using Go idiom.
In this example, Poller is a function that receives Resources to be polled
from an input channel,
and sends them to an output channel when they're done.

{{raw `
	type Resource string

	func Poller(in, out chan *Resource) {
	    for r := range in {
	        // poll the URL

	        // send the processed Resource to out
	        out <- r
	    }
	}
`}}

The delicate logic from the previous example is conspicuously absent,
and our Resource data structure no longer contains bookkeeping data.
In fact, all that's left are the important parts.
This should give you an inkling as to the power of these simple language features.

There are many omissions from the above code snippets.
For a walkthrough of a complete, idiomatic Go program that uses these ideas,
see the Codewalk [_Share Memory By Communicating_](/doc/codewalk/sharemem/).
