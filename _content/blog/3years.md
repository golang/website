---
title: Go turns three
date: 2012-11-10
by:
- Russ Cox
tags:
- community
- birthday
summary: Happy 3rd birthday, Go!
---


The Go open source project is
[three years old today](http://google-opensource.blogspot.com/2009/11/hey-ho-lets-go.html).

It's great to look at how far Go has come in those three years.
When we launched, Go was an idea backed by two implementations that worked on Linux and OS X.
The syntax, semantics, and libraries changed regularly as we reacted to feedback from users
and experience with the language.

Since the open source launch,
we've been joined by
hundreds of external contributors,
who have extended and improved Go in myriad ways,
including writing a Windows port from scratch.
We added a package management system
[goinstall](https://groups.google.com/d/msg/golang-nuts/8JFwR3ESjjI/cy7qZzN7Lw4J),
which eventually became the
[go command](/cmd/go/).
We also added
[support for Go on App Engine](https://blog.golang.org/2011/07/go-for-app-engine-is-now-generally.html).
Over the past year we've also given [many talks](/doc/#talks), created an [interactive introductory tour](/tour/)
and recently we added support for [executable examples in package documentation](/pkg/strings/#pkg-examples).

Perhaps the most important development in the past year
was the launch of the first stable version,
[Go 1](https://blog.golang.org/2012/03/go-version-1-is-released.html).
People who write Go 1 programs can now be confident that their programs will
continue to compile and run without change, in many environments,
on a time scale of years.
As part of the Go 1 launch we spent months cleaning up the
[language and libraries](/doc/go1.html)
to make it something that will age well.

We're working now toward the release of Go 1.1 in 2013. There will be some
new functionality, but that release will focus primarily on making Go perform
even better than it does today.

We're especially happy about the community that has grown around Go:
the mailing list and IRC channels seem like they are overflowing with discussion,
and a handful of Go books were published this year. The community is thriving.
Use of Go in production environments has also taken off, especially since Go 1.

We use Go at Google in a variety of ways, many of them invisible to the outside world.
A few visible ones include
[serving Chrome and other downloads](https://groups.google.com/d/msg/golang-nuts/BNUNbKSypE0/E4qSfpx9qI8J),
[scaling MySQL database at YouTube](http://code.google.com/p/vitess/),
and of course running the
[Go home page](/)
on [App Engine](https://developers.google.com/appengine/docs/go/overview).
Last year's
[Thanksgiving Doodle](https://blog.golang.org/2011/12/from-zero-to-go-launching-on-google.html)
and the recent
[Jam with Chrome](http://www.jamwithchrome.com/technology)
site are also served by Go programs.

Other companies and projects are using Go too, including
[BBC Worldwide](http://www.quora.com/Go-programming-language/Is-Google-Go-ready-for-production-use/answer/Kunal-Anand),
[Canonical](http://dave.cheney.net/wp-content/uploads/2012/08/august-go-meetup.pdf),
[CloudFlare](http://blog.cloudflare.com/go-at-cloudflare),
[Heroku](https://blog.golang.org/2011/04/go-at-heroku.html),
[Novartis](https://plus.google.com/114945221884326152379/posts/d1SVaqkRyTL),
[SoundCloud](http://backstage.soundcloud.com/2012/07/go-at-soundcloud/),
[SmugMug](http://sorcery.smugmug.com/2012/04/06/deriving-json-types-in-go/),
[StatHat](https://blog.golang.org/2011/12/building-stathat-with-go.html),
[Tinkercad](https://tinkercad.com/about/jobs),
and
[many others](/wiki/GoUsers).

Here's to many more years of productive programming in Go.
