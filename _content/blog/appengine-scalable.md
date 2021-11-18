---
title: Writing scalable App Engine applications
date: 2011-11-01
by:
- David Symonds
tags:
- appengine
- optimization
summary: How to build scalable web applications using Go with Google App Engine.
---


Back in May, we [announced](https://blog.golang.org/2011/05/go-and-google-app-engine.html)
the Go runtime for App Engine.
Since then, we've opened it up for everyone to use,
added many new APIs, and improved performance.
We have been thrilled by all the interesting ways that people are using Go on App Engine.
One of the key benefits of the Go runtime,
apart from working in a fantastic language,
is that it has high performance.
Go applications compile to native code, with no interpreter or virtual machine
getting between your program and the machine.

Making your web application fast is important because it is well known that
a web site's latency has a measurable impact on user happiness,
and [Google web search uses it as a ranking factor](https://googlewebmastercentral.blogspot.com/2010/04/using-site-speed-in-web-search-ranking.html).
Also announced in May was that App Engine would be [leaving its Preview status](http://googleappengine.blogspot.com/2011/05/year-ahead-for-google-app-engine.html)
and transitioning to a [new pricing model](https://www.google.com/enterprise/cloud/appengine/pricing.html),
providing another reason to write efficient App Engine applications.

To make it easier for Go developers using App Engine to write highly efficient,
scalable applications, we recently updated some existing App Engine articles
to include snippets of Go source code and to link to relevant Go documentation.

  - [Best practices for writing scalable applications](http://code.google.com/appengine/articles/scaling/overview.html)

  - [Managing Your App's Resource Usage](http://code.google.com/appengine/articles/managing-resources.html)
