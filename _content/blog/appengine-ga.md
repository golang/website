---
title: Go for App Engine is now generally available
date: 2011-07-21
by:
- Andrew Gerrand
tags:
- appengine
- release
summary: You can use Go on App Engine now!
---


The Go and App Engine teams are excited to announce that the Go runtime
for App Engine is now generally available.
This means you can take that Go app you've been working on (or meaning to
work on) and deploy it to App Engine right now with the new [1.5.2 SDK](http://code.google.com/appengine/downloads.html).

Since we announced the Go runtime at Google I/O we have continued to [improve and extend](http://code.google.com/p/googleappengine/wiki/SdkForGoReleaseNotes)
Go support for the App Engine APIs and have added the Channels API.
The Go Datastore API now supports transactions and ancestor queries, too.
See the [Go App Engine documentation](https://code.google.com/appengine/docs/go/)
for all the details.

For those who have been using the Go SDK already,
please note that the 1.5.2 release introduces `api_version` 2.
This is because the new SDK is based on Go `release.r58.1` (the current
stable version of Go) and is not backwards compatible with the previous release.
Existing apps may require changes as per the [r58 release notes](/doc/devel/release.html#r58).
Once you've updated your code, you should redeploy your app with the line
`api_version: 2` in its `app.yaml` file.
Apps written against `api_version` 1 will stop working after the 18th of August.

Finally, we owe a huge thanks to our trusted testers and their many bug reports.
Their help was invaluable in reaching this important milestone.

_The fastest way to get started with Go on App Engine is with the_ [_Getting Started guide_](http://code.google.com/appengine/docs/go/gettingstarted/).

_Note that the Go runtime is still considered experimental; it is not as well-supported as the Python and Java runtimes._
