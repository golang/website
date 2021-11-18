---
title: "Two Go Talks: \"Lexical Scanning in Go\" and \"Cuddle: an App Engine Demo\""
date: 2011-09-01
by:
- Andrew Gerrand
tags:
- appengine
- lexer
- talk
- video
summary: "Two talks about Go from the Sydney GTUG: Rob Pike explains lexical scanning, and Andrew Gerrand builds a simple real-time chat using App Engine."
---


On Tuesday night Rob Pike and Andrew Gerrand each presented at the [Sydney Google Technology User Group](http://www.sydney-gtug.org/).

Rob's talk, "[Lexical Scanning in Go](http://www.youtube.com/watch?v=HxaD_trXwRE)",
discusses the design of  a particularly interesting and idiomatic piece of Go code,
the lexer component of the new [template package.](/pkg/exp/template/)

{{video "https://www.youtube.com/embed/HxaD_trXwRE"}}

The slides are [available here](http://cuddle.googlecode.com/hg/talk/lex.html).
The new template package is available as [exp/template](/pkg/exp/template/) in Go release r59.
In a future release it will replace the old template package.

Andrew's talk, "[Cuddle: an App Engine Demo](http://www.youtube.com/watch?v=HQtLRqqB-Kk)",
describes the construction of a simple real-time chat application that uses
App Engine's [Datastore](http://code.google.com/appengine/docs/go/datastore/overview.html),
[Channel](http://code.google.com/appengine/docs/go/channel/overview.html),
and [Memcache](http://code.google.com/appengine/docs/go/datastore/memcache.html) APIs.
It also includes a question and answer session that covers [Go for App Engine](http://code.google.com/appengine/docs/go/gettingstarted/)
and Go more generally.

{{video "https://www.youtube.com/embed/HQtLRqqB-Kk"}}

The slides are [available here](http://cuddle.googlecode.com/hg/talk/index.html).
The code is available at the [cuddle Google Code project](http://code.google.com/p/cuddle/).
