---
title: "Real Go Projects: SmartTwitter and web.go"
date: 2010-10-19
by:
- Michael Hoisie
tags:
- guest
summary: How Michael Hoisie used Go to build SmartTwitter and web.go.
---


_This week's article is written by_ [_Michael Hoisie_](http://www.hoisie.com/).
_A programmer based in San Francisco, he is one of Go's early adopters and the author of several popular Go libraries. He describes his experiences using Go:_

I was introduced to Go by a post on [Hacker News](http://news.ycombinator.com/).
About an hour later I was hooked. At the time I was working at a web start-up,
and had been developing internal testing apps in Python.
Go offered speed, better concurrency support,
and sane Unicode handling, so I was keen to port my programs to the language.
At that time there wasn't an easy way to write web apps in Go,
so I decided to build a simple web framework,
[web.go](http://github.com/hoisie/web.go).
It was modeled after a popular Python framework,
[web.py](http://webpy.org/), which I had worked with previously.
While working on web.go I got involved in the Go community,
submitted a bunch of bug reports, and hacked on some standard library packages
(mainly [http](/pkg/http/) and [json](/pkg/json/)).

After a few weeks I noticed that web.go was getting attention at GitHub.
This was surprising because I'd never really promoted the project.
I think there's a niche for simple, fast web applications,
and I think Go can fill it.

One weekend I decided to write a simple Facebook application:
it would re-post your Twitter status updates to your Facebook profile.
There is an official Twitter application to do this,
but it re-posts everything, creating noise in your Facebook feed.
My application allowed you to filter retweets,
mentions, hashtags, replies, and more.
This turned into [Smart Twitter](http://www.facebook.com/apps/application.php?id=135488932982),
which currently has nearly 90,000 users.

The entire program is written in Go, and uses [Redis](https://redis.io/)
as its storage back-end.
It is very fast and robust. It currently processes about two dozen tweets per second,
and makes heavy use of Go's channels.
It runs on a single Virtual Private Server instance with 2GB of RAM,
which has no problem handling the load.
Smart Twitter uses very little CPU time, and is almost entirely memory-bound
as the entire database is kept in memory.
At any given time there are around 10 goroutines running concurrently:
one accepting HTTP connections, another reading from the Twitter Streaming API,
a couple for error handling, and the rest either processing web requests
or re-posting incoming tweets.

Smart Twitter also spawned other open-source Go projects:
[mustache.go](http://github.com/hoisie/mustache.go),
[redis.go](http://github.com/hoisie/redis.go),
and [twitterstream](http://github.com/hoisie/twitterstream).

I see a lot of work left to do on web.go.
For instance, I'd like to add better support for streaming connections,
websockets, route filters, better support in shared hosts,
and improving the documentation.
I recently left the start-up to do software freelancing,
and I'm planning to use Go where possible.
This means I'll probably use it as a back end for personal apps,
as well as for clients that like working with cutting edge technology.

Finally, I'd like to thank the Go team for all their effort.
Go is a wonderful platform and I think it has a bright future.
I hope to see the language grow around the needs of the community.
There's a lot of interesting stuff happening in the community,
and I look forward to seeing what people can hack together with the language.
