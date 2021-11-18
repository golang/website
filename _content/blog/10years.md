---
title: Go Turns 10
date: 2019-11-08
by:
- Russ Cox, for the Go team
summary: Happy 10th birthday, Go!
---


Happy birthday, Go!

This weekend we celebrate the 10th anniversary of
[the Go release](https://opensource.googleblog.com/2009/11/hey-ho-lets-go.html),
marking the 10th birthday of Go as an open-source programming language
and ecosystem for building modern networked software.

To mark the occasion,
[Renee French](https://twitter.com/reneefrench),
the creator of the
[Go gopher](https://blog.golang.org/gopher),
painted this delightful scene:

<a href="10years/gopher10th-large.jpg">
{{image "10years/gopher10th-small.jpg" 850}}
</a>

Celebrating 10 years of Go makes me think back to early November 2009,
when we were getting ready to share Go with the world.
We didn’t know what kind of reaction to expect,
whether anyone would care about this little language.
I hoped that even if no one ended up using Go,
we would at least have drawn attention to some good ideas,
especially Go’s approach to concurrency and interfaces,
that could influence follow-on languages.

Once it became clear that people were excited about Go,
I looked at the history of popular languages
like C, C++, Perl, Python, and Ruby,
examining how long each took to gain widespread adoption.
For example, Perl seemed to me to have appeared fully-formed
in the mid-to-late 1990s, with CGI scripts and the web,
but it was first released in 1987.
This pattern repeated for almost every language I looked at:
it seems to take roughly a decade of quiet, steady improvement
and dissemination before a new language really takes off.

I wondered: where would Go be after a decade?

Today, we can answer that question:
Go is everywhere, used by at least [a million developers worldwide](https://research.swtch.com/gophercount).

Go’s original target was networked system infrastructure,
what we now call cloud software.
Every major cloud provider today uses core cloud infrastructure written in Go,
such as Docker, Etcd, Istio, Kubernetes, Prometheus, and Terraform;
the majority of the
[Cloud Native Computing Foundation’s projects](https://www.cncf.io/projects/)
are written in Go.
Countless companies are using Go to move their own work to the cloud as well,
from startups building from scratch
to enterprises modernizing their software stack.
Go has also found adoption well beyond its original cloud target,
with uses ranging
from
controlling tiny embedded systems with
[GoBot](https://gobot.io) and [TinyGo](https://tinygo.org/)
to detecting cancer with
[massive big data analysis and machine learning at GRAIL](https://medium.com/grail-eng/bigslice-a-cluster-computing-system-for-go-7e03acd2419b),
and everything in between.

All this is to say that Go has succeeded beyond our wildest dreams.
And Go’s success isn’t just about the language.
It’s about the language, the ecosystem, and especially the community working together.

In 2009, the language was a good idea with a working sketch of an implementation.
The `go` command did not exist:
we ran commands like `6g` to compile and `6l` to link binaries,
automated with makefiles.
We typed semicolons at the ends of statements.
The entire program stopped during garbage collection,
which then struggled to make good use of two cores.
Go ran only on Linux and Mac, on 32- and 64-bit x86 and 32-bit ARM.

Over the last decade, with the help of Go developers all over the world,
we have evolved this idea and sketch into a productive language
with fantastic tooling,
a production-quality implementation,
a
[state-of-the-art garbage collector](https://blog.golang.org/ismmkeynote),
and [ports to 12 operating systems and 10 architectures](/doc/install/source#introduction).

Any programming language needs the support of a thriving ecosystem.
The open source release was the seed for that ecosystem,
but since then, many people have contributed their time and talent
to fill the Go ecosystem with great tutorials, books, courses, blog posts,
podcasts, tools, integrations, and of course reusable Go packages importable with `go` `get`.
Go could never have succeeded without the support of this ecosystem.

Of course, the ecosystem needs the support of a thriving community.
In 2019 there are dozens of Go conferences all over the world,
along with
[over 150 Go meetup groups with over 90,000 members](https://www.meetup.com/pro/go).
[GoBridge](https://golangbridge.org)
and
[Women Who Go](https://medium.com/@carolynvs/www-loves-gobridge-ccb26309f667)
help bring new voices into the Go community,
through mentoring, training, and conference scholarships.
This year alone, they have taught
hundreds of people from traditionally underrepresented groups
at workshops where community members teach and mentor those new to Go.

There are
[over a million Go developers](https://research.swtch.com/gophercount)
worldwide,
and companies all over the globe are looking to hire more.
In fact, people often tell us that learning Go
helped them get their first jobs in the tech industry.
In the end, what we’re most proud of about Go
is not a well-designed feature or a clever bit of code
but the positive impact Go has had in so many people’s lives.
We aimed to create a language that would help us be better developers,
and we are thrilled that Go has helped so many others.

As
[\#GoTurns10](https://twitter.com/search?q=%23GoTurns10),
I hope everyone will take a moment to celebrate
the Go community and all we have achieved.
On behalf of the entire Go team at Google,
thank you to everyone who has joined us over the past decade.
Let’s make the next one even more incredible!

<div>
<center>
<a href="10years/gopher10th-pin-large.jpg">
{{image "10years/gopher10th-pin-small.jpg" 150}}
</center>
</div>
