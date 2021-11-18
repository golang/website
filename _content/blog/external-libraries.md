---
title: Spotlight on external Go libraries
date: 2011-06-03
by:
- Andrew Gerrand
tags:
- community
- libraries
summary: Some popular Go libraries and how to use them.
---


While the Go authors have been working hard at improving Go's standard library,
the greater community has created a growing ecosystem of external libraries.
In this post we look at some popular Go libraries and how they can be used.

[Mgo](http://labix.org/mgo) (pronounced "mango") is a MongoDB database driver.
[MongoDB](http://www.mongodb.org/) is a [document-oriented database](http://en.wikipedia.org/wiki/Document-oriented_database)
with a long list of features suitable for [a broad range of uses](http://www.mongodb.org/display/DOCS/Use%2BCases).
The mgo package provides a rich, idiomatic Go API for working with MongoDB,
from basic operations such as inserting and updating records to the more
advanced [MapReduce](http://www.mongodb.org/display/DOCS/MapReduce) and
[GridFS](http://www.mongodb.org/display/DOCS/GridFS) features.
Mgo has a bunch of cool features including automated cluster discovery and
result pre-fetching - see the [mgo homepage](http://labix.org/mgo) for
details and example code.
For working with large data sets Go, MongoDB,
and mgo are a powerful combination.

[Authcookie](https://github.com/dchest/authcookie) is a web library for
generating and verifying user authentication cookies.
It allows web servers to hand out cryptographically secure tokens tied to
a specific user that will expire after a specified time period.
It has a simple API that makes it straightforward to add authentication
to existing web applications.
See the [README file](https://github.com/dchest/authcookie/blob/master/README.md)
for details and example code.

[Go-charset](http://code.google.com/p/go-charset) provides support for
converting between Go's standard UTF-8 encoding and a variety of character sets.
The go-charset package implements a translating io.Reader and io.Writer
so you can wrap existing Readers and Writers (such as network connections
or file descriptors),
making it easy to communicate with systems that use other character encodings.

[Go-socket.io](https://github.com/madari/go-socket.io) is a Go implementation
of [Socket.IO](http://socket.io/),
a client/server API that allows web servers to push messages to web browsers.
Depending on the capabilities of the user's browser,
Socket.IO uses the best transport for the connection,
be it modern websockets, AJAX long polling,
or some [other mechanism](http://socket.io/#transports).
Go-socket.io bridges the gap between Go servers and rich JavaScript clients
for a wide range of browsers.
To get a feel for go-socket.io see the [chat server example](https://github.com/madari/go-socket.io/blob/master/example/example.go).

It's worth mentioning that these packages are [goinstallable](/cmd/goinstall/).
With an up-to-date Go [installation](/doc/install.html)
you can install them all with a single command:

	goinstall launchpad.net/mgo \
	    github.com/dchest/authcookie \
	    go-charset.googlecode.com/hg/charset \
	    github.com/madari/go-socket.io

Once goinstalled, the packages can be imported using those same paths:

	import (
	    "launchpad.net/mgo"
	    "github.com/dchest/authcookie"
	    "go-charset.googlecode.com/hg/charset"
	    "github.com/madari/go-socket.io"
	)

Also, as they are now a part of the local Go system,
we can inspect their documentation with [godoc](/cmd/godoc/):

	godoc launchpad.net/mgo Database # see docs for Database type

Of course, this is just the tip of the iceberg;
there are more great Go libraries listed on the [package dashboard](http://godashboard.appspot.com/package)
and many more to come.
