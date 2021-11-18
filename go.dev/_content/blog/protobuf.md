---
title: "Third-party libraries: goprotobuf and beyond"
date: 2010-04-20
by:
- Andrew Gerrand
tags:
- protobuf
- community
summary: Announcing Go support for Protocol Buffers, Google's data interchange format.
---


On March 24, Rob Pike announced [goprotobuf](http://code.google.com/p/goprotobuf/),
the Go bindings of Google's data interchange format [Protocol Buffers](http://code.google.com/apis/protocolbuffers/docs/overview.html),
called protobufs for short.
With this announcement, Go joins C++, Java,
and Python as languages providing official protobuf implementations.
This marks an important milestone in enabling the interoperability between
existing systems and those built in Go.

The goprotobuf project consists of two parts:
a 'protocol compiler plugin' that generates Go source files that,
once compiled, can access and manage protocol buffers;
and a Go package that implements run-time support for encoding (marshaling),
decoding (unmarshaling), and accessing protocol buffers.

To use goprotobuf, you first need to have both Go and [protobuf](http://code.google.com/p/protobuf/) installed.
You can then install the 'proto' package with [goinstall](/cmd/goinstall/):

	goinstall goprotobuf.googlecode.com/hg/proto

And then install the protobuf compiler plugin:

	cd $GOROOT/src/pkg/goprotobuf.googlecode.com/hg/compiler
	make install

For more detail see the project's [README](http://code.google.com/p/goprotobuf/source/browse/README) file.

This is one of a growing list of third-party [Go projects](http://godashboard.appspot.com/package).
Since the announcement of goprotobuf, the X Go bindings have been spun off
from the standard library to the [x-go-binding](http://code.google.com/p/x-go-binding/) project,
and work has begun on a [Freetype](http://www.freetype.org/) port,
[freetype-go](http://code.google.com/p/freetype-go/).
Other popular third-party projects include the lightweight web framework
[web.go](http://github.com/hoisie/web.go),
and the Go GTK bindings [gtk-go](http://github.com/mattn/go-gtk).

We wish to encourage the development of other useful packages by the open source community.
If you're working on something, don't keep it to yourself - let us know
through our mailing list [golang-nuts](http://groups.google.com/group/golang-nuts).
