---
title: Half a decade with Go
date: 2014-11-10
by:
- Andrew Gerrand
summary: Happy 5th birthday, Go!
---


Five years ago we launched the Go project. It seems like only yesterday that we
were preparing the initial public release: our
[website](https://web.archive.org/web/20091112094121/http://golang.org/) was
a lovely shade of yellow, we were calling Go a "systems language", and you had
to terminate statements with a semicolon and write Makefiles to build your
code. We had no idea how Go would be received. Would people share our vision
and goals? Would people find Go useful?

At launch, there was a flurry of attention. Google had produced a new
programming language, and everyone was eager to check it out. Some programmers
were turned off by Go's conservative feature set—at first glance they saw
"nothing to see here"—but a smaller group saw the beginnings of an ecosystem
tailored to their needs as working software engineers. These few would form the
kernel of the Go community.

{{image "5years/gophers5th.jpg" 850}}

[_Gopher_](/blog/gopher) _illustration by_ [_Renee French_](http://reneefrench.blogspot.com.au/)

After the initial release, it took us a while to properly communicate the
goals and design ethos behind Go. Rob Pike did so eloquently in his 2012 essay
[_Go at Google: Language Design in the Service of Software Engineering_](/talks/2012/splash.article) and
more personally in his blog post
[_Less is exponentially more_](https://commandcenter.blogspot.com.au/2012/06/less-is-exponentially-more.html).
Andrew Gerrand's
[_Code that grows with grace_](http://vimeo.com/53221560)
([slides](/talks/2012/chat.slide)) and
[_Go for Gophers_](https://www.youtube.com/watch?v=dKGmK_Z1Zl0)
([slides](/talks/2014/go4gophers.slide)) give a
more in-depth, technical take on Go's design philosophy.

Over time, the few became many. The turning point for the project was the
release of Go 1 in March 2012, which provided a stable language and standard
library that developers could trust. By 2014, the project had hundreds of core
contributors, the ecosystem had countless [libraries and tools](https://godoc.org/)
maintained by thousands of developers, and the greater community had
many passionate members (or, as we call them, "gophers"). Today, by our current
metrics, the Go community is growing faster than we believed possible.

Where can those gophers be found? They are at the many Go events that are
popping up around the world. This year we saw several dedicated Go conferences:
the inaugural [GopherCon](https://blog.golang.org/gophercon) and
[dotGo](http://www.dotgo.eu/) conferences in Denver and Paris, the
[Go DevRoom at FOSDEM](https://blog.golang.org/fosdem14) and two more
instances of the biannual [GoCon](https://github.com/GoCon/GoCon) conference
in Tokyo. At each event, gophers from around the globe eagerly presented their
Go projects. For the Go team, it is very satisfying to meet so many programmers
that share our vision and excitement.

{{image "5years/conferences.jpg"}}

_More than 1,200 gophers attended GopherCon in Denver and dotGo in Paris._

There are also dozens of community-run
[Go User Groups](/wiki/GoUserGroups) spread across cities
worldwide. If you haven't visited your local group, consider going along. And
if there isn't a group in your area, maybe you should
[start one](https://blog.golang.org/getthee-to-go-meetup)?

Today, Go has found a home in the cloud. Go arrived as the industry underwent a
tectonic shift toward cloud computing, and we were thrilled to see it quickly
become an important part of that movement. Its simplicity, efficiency, built-in
concurrency primitives, and modern standard library make it a great fit for
cloud software development (after all, that's what it was designed for).
Significant open source cloud projects like
[Docker](https://www.docker.com/) and
[Kubernetes](https://github.com/GoogleCloudPlatform/kubernetes) have been
written in Go, and infrastructure companies like Google, CloudFlare, Canonical,
Digital Ocean, GitHub, Heroku, and Microsoft are now using Go to do some heavy
lifting.

So, what does the future hold? We think that 2015 will be Go's biggest year yet.

Go 1.4—in addition to its [new features and fixes](/doc/go1.4)—lays
the groundwork for a new low-latency garbage collector and support for running
Go on mobile devices. It is due to be released on December 1st 2014.
We expect the new GC to be available in Go 1.5, due June 1st 2015, which will
make Go appealing for a broader range of applications.
We can't wait to see where people take it.

And there will be more great events, with [GothamGo](http://gothamgo.com/) in
New York (15 Nov), another Go DevRoom at FOSDEM in Brussels (Jan 31 and Feb 1;
[get involved!](https://groups.google.com/d/msg/golang-nuts/1xgBazQzs1I/hwrZ5ni8cTEJ)),
[GopherCon India](http://www.gophercon.in/) in Bengaluru (19-21 Feb),
the original [GopherCon](http://gophercon.com/) back at Denver in July, and
[dotGo](http://www.dotgo.eu/) on again at Paris in November.

The Go team would like to extend its thanks to all the gophers out there.
Here's to the next five years.

_To celebrate 5 years of Go, over the coming month the_
[_Gopher Academy_](http://blog.gopheracademy.com/)
_will publish a series of articles by prominent Go users. Be sure to check out_
[_their blog_](http://blog.gopheracademy.com/)
_for more Go action._
