---
title: "Go at I/O: Frequently Asked Questions"
date: 2010-05-27
by:
- Andrew Gerrand
tags:
- appengine
summary: Q&A about Go from Google I/O 2010.
---


Among the high-profile product launches at Google I/O last week,
our small team gave presentations to packed rooms and met many present and
future Go programmers.
It was especially gratifying to meet with so many people who,
after learning a bit about Go, were excited by the potential benefits (both
immediate and long-term) they could gain from using it.

We were asked a lot of good questions during I/O, and in this post I'd like to recap and expand upon some of them.

How suitable is Go for production systems?
Go is ready and stable now. We are pleased to report that Google is using
Go for some production systems,
and they are performing well.
Of course there is still room for improvement - that's why we're continuing
to work on the language,
libraries, tools, and runtime.

Do you have plans to implement generics?
Many proposals for generics-like features have been mooted both publicly and internally,
but as yet we haven't found a proposal that is consistent with the rest of the language.
We think that one of Go's key strengths is its simplicity,
so we are wary of introducing new features that might make the language
more difficult to understand.
Additionally, the more Go code we write (and thus the better we learn how
to write Go code ourselves),
the less we feel the need for such a language feature.

Do you have any plans to support GPU programming?
We don't have any immediate plans to do this,
but as Go is architecture-agnostic it's quite possible.
The ability to launch a goroutine that runs on a different processor architecture,
and to use channels to communicate between goroutines running on separate architectures,
seem like good ideas.

Are there plans to support Go under App Engine?
Both the Go and App Engine teams would like to see this happen.
As always, it is a question of resources and priorities as to if and when
it will become a reality.

Are there plans to support Go under Android?
Both Go compilers support ARM code generation, so it is possible.
While we think Go would be a great language for writing mobile applications,
Android support is not something that's being actively worked on.

What can I use Go for?
Go was designed with systems programming in mind.
Servers, clients, databases, caches, balancers,
distributors - these are applications Go is obviously useful for,
and  this is how we have begun to use it within Google.
However, since Go's open-source release, the community has found a diverse
range of applications for the language.
From web apps to games to graphics tools,
Go promises to shine as a general-purpose programming language.
The potential is only limited by library support,
which is improving at a tremendous rate.
Additionally, educators have expressed interest in using Go to teach programming,
citing its succinct syntax and consistency as well-suited to the task.

Thanks to everyone who attended our presentations,
or came to talk with us at Office Hours.
We hope to see you again at future events.

The video of Rob and Russ' talk is [available on YouTube](https://youtu.be/jgVhBThJdXc).
