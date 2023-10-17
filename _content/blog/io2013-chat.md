---
title: A conversation with the Go team
date: 2013-06-06
summary: At Google I/O 2013, several members of the Go team hosted a "Fireside chat."
---


At Google I/O 2013, several members of the Go team hosted a "Fireside chat."
Robert Griesemer, Rob Pike, David Symonds, Andrew Gerrand, Ian Lance Taylor,
Sameer Ajmani, Brad Fitzpatrick, and Nigel Tao took questions from the audience
and people around the world about various aspects of the Go project.

{{video "https://www.youtube.com/embed/p9VUCp98ay4"}}

We also hosted a similar session at I/O last year:
[_Meet the Go team_](http://www.youtube.com/watch?v=sln-gJaURzk).

There were many more questions from Google Moderator than we were able to
answer in the short 40 minute session.
Here we answer some of those we missed in the live session.

_Linking speed (and memory usage) for the gc toolchain are a known problem._
_Are there any plans to address this during the 1.2 cycle?_

**Rob:** Yes. We are always thinking about ways to improve performance of the
tools as well as the language and libraries.

_I have been very pleased to see how quickly Go appears to be gaining traction._
_Can you talk about the reactions you have experienced working with other_
_developers inside and outside Google? Are there any major sticking points remaining?_

**Robert:** A lot of developers that seriously tried Go are very happy with it.
Many of them report a much smaller, more readable and thus maintainable code
base: A 50% code size reduction or more when coming from C++ seems common.
Developers that switched to Go from Python are invariably pleased with the
performance gain. The typical complaints are about small inconsistencies in the
language (some of which we might iron out at some point). What surprises me is
that almost nobody complains about the lack of generics.

_When will Go be a first-class language for Android development?_

**Andrew:** This would be great, but we don't have anything to announce.

_Is there a roadmap for the next version of Go?_

**Andrew:** We have no feature roadmap as such. The contributors tend to work on
what interests them. Active areas of development include the gc and gccgo
compilers, the garbage collector and runtime, and many others. We expect the
majority of exciting new additions will be in the form of improvements to our
tools. You can find design discussions and code reviews on the
[golang-dev mailing list](http://groups.google.com/group/golang-dev).

As for the timeline, we do have
[concrete plans](https://docs.google.com/document/d/106hMEZj58L9nq9N9p7Zll_WKfo-oyZHFyI6MttuZmBU/edit?usp=sharing):
we expect to release Go 1.2 on December 1, 2013.

_Where do you guys want to see Go used externally?_
_What would you consider a big win for Go adoption outside Google?_
_Where do you think Go has the potential to make a significant impact?_

**Rob:** Where Go is deployed is up to its users, not to us. We're happy to see
it gain traction anywhere it helps. It was designed with server-side software
in mind, and is showing promise there, but has also shown strengths in many
other areas and the story is really just beginning. There are many surprises to
come.

**Ian:** It’s easier for startups to use Go, because they don’t have an
entrenched code base that they need to work with. So I see two future big wins
for Go. One would be a significant use of Go by an existing large software
company other than Google. Another would be a significant IPO or acquisition
of a startup that primarily uses Go. These are both indirect: clearly choice
of programming language is a very small factor in the success of a company.
But it would be another way to show that Go can be part of a successful
software system.

_Have you thought any (more) about the potential of dynamically loading_
_Go packages or objects and how it could work in Go?_
_I think this could enable some really interesting and expressive constructs,_
_especially coupled with interfaces._

**Rob:** This is an active topic of discussion. We appreciate how powerful the
concept can be and hope we can find a way to implement it before too long.
There are serious challenges in the design approach to take and the need to
make it work portably.

_There was a discussion a while ago about collecting some best-of-breed_
`database/sql` _drivers in a more central place._
_Some people had strong opinions to the contrary though._
_Where is_ `database/sql` _and its drivers going in the next year?_

**Brad:** While we could create an official subrepo (“go.db”) for database
drivers, we fear that would unduly bless certain drivers. At this point we’d
still rather see healthy competition between different drivers. The
[SQLDrivers wiki page](/wiki/SQLDrivers)
lists some good ones.

The `database/sql` package didn’t get much attention for a while, due to lack of
drivers. Now that drivers exist, usage of the package is increasing and
correctness and performance bugs are now being reported (and fixed). Fixes will
continue, but no major changes to the interface of `database/sql` are planned.
 There might be small extensions here and there as needed for performance or to
assist some drivers.

_What is the status of versioning?_
_Is importing some code from GitHub a best practice recommended by the Go team?_
_What happens when we publish our code that is dependent on a GitHub repo and_
_the API of the dependee changes?_

**Ian:** This is frequently discussed on the mailing list. What we do internally
is take a snapshot of the imported code, and update that snapshot from time to
time. That way, our code base won't break unexpectedly if the API changes.
But we understand that that approach doesn’t work very well for people who are
themselves providing a library. We’re open to good suggestions in this area.
Remember that this is an aspect of the tools that surround the language rather
than the language itself; the place to fix this is in the tools, not the
language.

_What about Go and Graphical User Interfaces?_

**Rob:** This is a subject close to my heart. Newsqueak, a very early precursor
language, was designed specifically for writing graphics programs (that's what
we used to call apps). The landscape has changed a lot but I think Go's
concurrency model has much to offer in the field of interactive graphics.

**Andrew:** There are many
[bindings for existing graphics libraries](/wiki/Projects#Graphics_and_Audio)
out there, and a few Go-specific projects. One of the more promising ones is
[go.uik](https://github.com/skelterjohn/go.uik), but it's still in its early
days. I think there's a lot of potential for a great Go-specific UI toolkit for
writing native applications (consider handling user events by receiving from a
channel), but developing a production-quality package is a significant
undertaking. I have no doubt one will come in time.

In the meantime, the web is the most broadly available platform for user
interfaces. Go provides great support for building web apps, albeit only on the
back end.

_In the mailing lists Adam Langley has stated that the TLS code has not been_
_reviewed by outside groups, and thus should not be used in production._
_Are there plans to have the code reviewed?_
_A good secure implementation of concurrent TLS would be very nice._

**Adam**: Cryptography is notoriously easy to botch in subtle and surprising ways
and I’m only human. I don’t feel that I can warrant that Go’s TLS code is
flawless and I wouldn’t want to misrepresent it.

There are a couple of places where the code is known to have side-channel
issues: the RSA code is blinded but not constant time, elliptic curves other
than P-224 are not constant time and the Lucky13 attack might work. I hope to
address the latter two in the Go 1.2 timeframe with a constant-time P-256
implementation and AES-GCM.

Nobody has stepped forward to do a review of the TLS stack however and I’ve not
investigated whether we could get Matasano or the like to do it. That depends
on whether Google wishes to fund it.

_What do you think about_ [_GopherCon 2014_](http://www.gophercon.com/)_?_
_Does anyone from the team plan to attend?_

**Andrew:** It's very exciting. I'm sure some of us will be there.
