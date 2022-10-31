---
title: Go, Open Source, Community
date: 2015-07-08
by:
- Russ Cox
tags:
- community
summary: Why is Go open source, and how can we strengthen our open-source community?
---

## Welcome

[This is the text of my opening keynote at Gophercon 2015.
[The video is available here](https://www.youtube.com/watch?v=XvZOdpd_9tc).]

Thank you all for traveling to Denver to be here,
and thank you to everyone watching on video.
If this is your first Gophercon, welcome.
If you were here last year, welcome back.
Thank you to the organizers
for all the work it takes
to make a conference like this happen.
I am thrilled to be here and to be able to talk to all of you.

I am the tech lead for the Go project
and the Go team at Google.
I share that role with Rob Pike.
In that role, I spend a lot of time thinking about
the overall Go open source project,
in particular the way it runs,
what it means to be open source,
and the interaction between
contributors inside and outside Google.
Today I want to share with you
how I see the Go project as a whole
and then based on that explain
how I see the Go open source project
evolving.

## Why Go?

To get started,
we have to go back to the beginning.
Why did we start working on Go?

Go is an attempt to make programmers more productive.
We wanted to improve the software development process
at Google,
but the problems Google has
are not unique to Google.

There were two overarching goals.

The first goal is to make a better language
to meet the challenges of scalable concurrency.
By scalable concurrency I mean
software that deals with many concerns simultaneously,
such as coordinating a thousand back end servers
by sending network traffic back and forth.

Today, that kind of software has a shorter name:
we call it cloud software.
It's fair to say that Go was designed for the cloud
before clouds ran software.

The larger goal is to make a better environment
to meet the challenges of scalable software development,
software worked on and used by many people,
with limited coordination between them,
and maintained for years.
At Google we have thousands of engineers
writing and sharing their code with each other,
trying to get their work done,
reusing the work of others as much as possible,
and working in a code base with a history
dating back over ten years.
Engineers often work on or at least look at
code originally written by someone else,
or that they wrote years ago,
which often amounts to the same thing.

That situation inside Google
has a lot in common with
large scale, modern open source development
as practiced on sites like GitHub.
Because of this,
Go is a great fit for open source projects,
helping them accept and manage
contributions from a large community
over a long period of time.

I believe much of Go's success is explained by the fact that
Go is a great fit for cloud software,
Go is a great fit for open source projects,
and, serendipitously, both of those are
growing in popularity and importance
in the software industry.

Other people have made similar observations.
Here are two.
Last year, on RedMonk.com, Donnie Berkholz
wrote about
“[Go as the emerging language of cloud infrastructure](http://redmonk.com/dberkholz/2014/03/18/go-the-emerging-language-of-cloud-infrastructure/),”
observing that
“[Go's] marquee projects ... are cloud-centric or otherwise
made for dealing with distributed systems
or transient environments.”

This year, on Texlution.com, the author
wrote an article titled
“[Why Golang is doomed to succeed](https://texlution.com/post/why-go-is-doomed-to-succeed/),”
pointing out that this focus on large-scale development
was possibly even better suited to open source than
to Google itself: “This open source fitness is why I think
you are about to see more and more Go around ...”

## The Go Balance

How does Go accomplish those things?

How does it make scalable concurrency
and scalable software development easier?

Most people answer this question by talking about
channels and goroutines, and interfaces, and fast builds,
and the go command, and good tool support.
Those are all important parts of the answer,
but I think there is a broader idea behind them.

I think of that idea as Go's balance.
There are competing concerns in any software design,
and there is a very natural tendency to try to solve
all the problems you foresee.
In Go, we have explicitly tried not to solve everything.
Instead, we've tried to do just enough that you can build
your own custom solutions easily.

The way I would summarize Go's chosen balance is this: **Do Less. Enable More.**

Do less, but enable more.

Go can't do everything.
We shouldn't try.
But if we work at it,
Go can probably do
a few things well.
If we select those things carefully,
we can lay a foundation
on which developers can _easily_ build
the solutions and tools they need,
and ideally can interoperate with
the solutions and tools built by others.

### Examples

Let me illustrate this with some examples.

First, the size of the Go language itself.
We worked hard to put in as few concepts as possible,
to avoid the problem of mutually incomprehensible dialects
forming in different parts of a large developer community.
No idea went into Go until
it had been simplified to its essence
and then had clear benefits
that justified the complexity being added.

In general, if we have 100 things
we want Go to do well,
we can't make 100 separate changes.
Instead, we try to research and understand
the design space
and then identify a few changes
that work well together
and that enable maybe 90 of those things.
We're willing to sacrifice the remaining 10
to avoid bloating the language,
to avoid adding complexity
only to address specific use cases
that seem important today
but might be gone tomorrow.

Keeping the language small
enables more important goals.
Being small makes Go
easier to learn,
easier to understand,
easier to implement,
easier to reimplement,
easier to debug,
easier to adjust,
and easier to evolve.
Doing less enables more.

I should point out that
this means we say no
to a lot of other people's ideas,
but I assure you
we've said no
to even more of our own ideas.

Next, channels and goroutines.
How should we structure and coordinate
concurrent and parallel computations?
Mutexes and condition variables are very general
but so low-level that they're difficult to use correctly.
Parallel execution frameworks like OpenMP are so high-level
that they can only be used to solve a narrow range of problems.
Channels and goroutines sit between these two extremes.
By themselves, they aren't a solution to much.
But they are powerful enough to be easily arranged
to enable solutions to many common problems
in concurrent software.
Doing less—really doing just enough—enables more.

Next, types and interfaces.
Having static types enables useful compile-time checking,
something lacking in dynamically-typed languages
like Python or Ruby.
At the same time,
Go's static typing avoids
much of the repetition
of traditional statically typed languages,
making it feel more lightweight,
more like the dynamically-typed languages.
This was one of the first things people noticed,
and many of Go's early adopters came from
dynamically-typed languages.

Go's interfaces are a key part of that.
In particular,
omitting the ``implements'' declarations
of Java or other languages with static hierarchy
makes interfaces lighter weight and more flexible.
Not having that rigid hierarchy
enables idioms such as test interfaces that describe
existing, unrelated production implementations.
Doing less enables more.

Next, testing and benchmarking.
Is there any shortage of testing
and benchmarking frameworks in most languages?
Is there any agreement between them?

Go's testing package is not meant
to address every possible facet of these topics.
Instead, it is meant to provide
the basic concepts necessary
for most higher-level tooling.
Packages have test cases that pass, fail, or are skipped.
Packages have benchmarks that run and can be measured
by various metrics.

Doing less here is an attempt
to reduce these concepts to their essence,
to create a shared vocabulary
so that richer tools can interoperate.
That agreement enables higher-level testing software
like Miki Tebeka's go2xunit converter,
or the benchcmp and benchstat
benchmark analysis tools.

Because there _is_ agreement
about the representation of the basic concepts,
these higher-level tools work for all Go packages,
not just ones that make the effort to opt in,
and they interoperate with each other,
in that using, say, go2xunit
does not preclude also using benchstat,
the way it would if these tools were, say,
plugins for competing testing frameworks.
Doing less enables more.

Next, refactoring and program analysis.
Because Go is for large code bases,
we knew it would need to support automatic
maintenance and updating of source code.
We also knew that this topic was too large
to build in directly.
But we knew one thing that we had to do.
In our experience attempting
automated program changes in other settings,
the most significant barrier we hit
was actually writing the modified program out
in a format that developers can accept.

In other languages,
it's common for different teams to use
different formatting conventions.
If an edit by a program uses the wrong convention,
it either writes a section of the source file that looks nothing
like the rest of the file, or it reformats the entire file,
causing unnecessary and unwanted diffs.

Go does not have this problem.
We designed the language to make gofmt possible,
we worked hard
to make gofmt's formatting acceptable
for all Go programs,
and we made sure gofmt was there
from day one of the original public release.
Gofmt imposes such uniformity that
automated changes blend into the rest of the file.
You can't tell whether a particular change
was made by a person or a computer.
We didn't build explicit refactoring support.
Establishing an agreed-upon formatting algorithm
was enough of a shared base
for independent tools to develop and to interoperate.
Gofmt enabled gofix, goimports, eg, and other tools.
I believe the work here is only just getting started.
Even more can be done.

Last, building and sharing software.
In the run up to Go 1, we built goinstall,
which became what we all know as "go get".
That tool defined a standard zero-configuration way
to resolve import paths on sites like github.com,
and later a way to resolve paths on other sites
by making HTTP requests.
This agreed-upon resolution algorithm
enabled other tools that work in terms of those paths,
most notably Gary Burd's creation of godoc.org.
In case you haven't used it,
you go to godoc.org/the-import-path
for any valid "go get" import path,
and the web site will fetch the code
and show you the documentation for it.
A nice side effect of this has been that
godoc.org serves as a rough master list
of the Go packages publicly available.
All we did was give import paths a clear meaning.
Do less, enable more.

You'll notice that many of these tooling examples
are about establishing a shared convention.
Sometimes people refer to this as Go being “opinionated,”
but there's something deeper going on.
Agreeing to the limitations
of a shared convention
is a way to enable
a broad class of tools that interoperate,
because they all speak the same base language.
This is a very effective way
to do less but enable more.
Specifically, in many cases
we can do the minimum required
to establish a shared understanding
of a particular concept, like remote imports,
or the proper formatting of a source file,
and thereby enable
the creation of packages and tools
that work together
because they all agree
about those core details.

I'm going to return to that idea later.

## Why is Go open source?

But first, as I said earlier,
I want to explain how I see
the balance of Do Less and Enable More
guiding our work
on the broader
Go open source project.
To do that, I need to start with
why Go is open source at all.

Google pays me and others to work on Go, because,
if Google's programmers are more productive,
Google can build products faster,
maintain them more easily,
and so on.
But why open source Go?
Why should Google share this benefit with the world?

Of course, many of us
worked on open source projects before Go,
and we naturally wanted Go
to be part of that open source world.
But our preferences are not a business justification.
The business justification is that
Go is open source
because that's the only way
that Go can succeed.
We, the team that built Go within Google,
knew this from day one.
We knew that Go had to be made available
to as many people as possible
for it to succeed.

Closed languages die.

A language needs large, broad communities.

A language needs lots of people writing lots of software,
so that when you need a particular tool or library,
there's a good chance it has already been written,
by someone who knows the topic better than you,
and who spent more time than you have to make it great.

A language needs lots of people reporting bugs,
so that problems are identified and fixed quickly.
Because of the much larger user base,
the Go compilers are much more robust and spec-compliant
than the Plan 9 C compilers they're loosely based on ever were.

A language needs lots of people using it
for lots of different purposes,
so that the language doesn't overfit to one use case
and end up useless when the technology landscape changes.

A language needs lots of people who want to learn it,
so that there is a market for people to write books
or teach courses,
or run conferences like this one.

None of this could have happened
if Go had stayed within Google.
Go would have suffocated inside Google,
or inside any single company
or closed environment.

Fundamentally,
Go must be open,
and Go needs you.
Go can't succeed without all of you,
without all the people using Go
for all different kinds of projects
all over the world.

In turn, the Go team at Google
could never be large enough
to support the entire Go community.
To keep scaling,
we
need to enable all this ``more''
while doing less.
Open source is a huge part of that.

## Go's open source

What does open source mean?
The minimum requirement is to open the source code,
making it available under an open source license,
and we've done that.

But we also opened our development process:
since announcing Go,
we've done all our development in public,
on public mailing lists open to all.
We accept and review
source code contributions from anyone.
The process is the same
whether you work for Google or not.
We maintain our bug tracker in public,
we discuss and develop proposals for changes in public,
and we work toward releases in public.
The public source tree is the authoritative copy.
Changes happen there first.
They are only brought into
Google's internal source tree later.
For Go, being open source means
that this is a collective effort
that extends beyond Google, open to all.

Any open source project starts with a few people,
often just one, but with Go it was three:
Robert Griesemer, Rob Pike, and Ken Thompson.
They had a vision of
what they wanted Go to be,
what they thought Go could do better
than existing languages, and
Robert will talk more about that tomorrow morning.
I was the next person to join the team,
and then Ian Taylor,
and then, one by one,
we've ended up where we are today,
with hundreds of contributors.

Thank You
to the many people who have contributed
code
or ideas
or bug reports
to the Go project so far.
We tried to list everyone we could
in our space in the program today.
If your name is not there,
I apologize,
but thank you.

I believe
the hundreds of contributors so far
are working toward a shared vision
of what Go can be.
It's hard to put words to these things,
but I did my best
to explain one part of the vision
earlier:
Do Less, Enable More.

## Google's role

A natural question is:
What is the role
of the Go team at Google,
compared to other contributors?
I believe that role
has changed over time,
and it continues to change.
The general trend is that
over time
the Go team at Google
should be doing less
and enabling more.

In the very early days,
before Go was known to the public,
the Go team at Google
was obviously working by itself.
We wrote the first draft of everything:
the specification,
the compiler,
the runtime,
the standard library.

Once Go was open sourced, though,
our role began to change.
The most important thing
we needed to do
was communicate our vision for Go.
That's difficult,
and we're still working at it.
The initial implementation
was an important way
to communicate that vision,
as was the development work we led
that resulted in Go 1,
and the various blog posts,
and articles,
and talks we've published.

But as Rob said at Gophercon last year,
"the language is done."
Now we need to see how it works,
to see how people use it,
to see what people build.
The focus now is on
expanding the kind of work
that Go can help with.

Google's primarily role is now
to enable the community,
to coordinate,
to make sure changes work well together,
and to keep Go true to the original vision.

Google's primary role is:
Do Less. Enable More.

I mentioned earlier
that we'd rather have a small number of features
that enable, say, 90% of the target use cases,
and avoid the orders of magnitude
more features necessary
to reach 99 or 100%.
We've been successful in applying that strategy
to the areas of software that we know well.
But if Go is to become useful in many new domains,
we need experts in those areas
to bring their expertise
to our discussions,
so that together
we can design small adjustments
that enable many new applications for Go.

This shift applies not just to design
but also to development.
The role of the Go team at Google
continues to shift
more to one of guidance
and less of pure development.
I certainly spend much more time
doing code reviews than writing code,
more time processing bug reports
than filing bug reports myself.
We need to do less and enable more.

As design and development shift
to the broader Go community,
one of the most important things
we
the original authors of Go
can offer
is consistency of vision,
to help keep Go
Go.
The balance that we must strike
is certainly subjective.
For example, a mechanism for extensible syntax
would be a way to
enable more
ways to write Go code,
but that would run counter to our goal
of having a consistent language
without different dialects.

We have to say no sometimes,
perhaps more than in other language communities,
but when we do,
we aim to do so
constructively and respectfully,
to take that as an opportunity
to clarify the vision for Go.

Of course, it's not all coordination and vision.
Google still funds Go development work.
Rick Hudson is going to talk later today
about his work on reducing garbage collector latency,
and Hana Kim is going to talk tomorrow
about her work on bringing Go to mobile devices.
But I want to make clear that,
as much as possible,
we aim to treat
development funded by Google
as equal to
development funded by other companies
or contributed by individuals using their spare time.
We do this because we don't know
where the next great idea will come from.
Everyone contributing to Go
should have the opportunity to be heard.

### Examples

I want to share some evidence for this claim
that, over time,
the original Go team at Google
is focusing more on
coordination than direct development.

First, the sources of funding
for Go development are expanding.
Before the open source release,
obviously Google paid for all Go development.
After the open source release,
many individuals started contributing their time,
and we've slowly but steadily
been growing the number of contributors
supported by other companies
to work on Go at least part-time,
especially as it relates to
making Go more useful for those companies.
Today, that list includes
Canonical, Dropbox, Intel, Oracle, and others.
And of course Gophercon and the other
regional Go conferences are organized
entirely by people outside Google,
and they have many corporate sponsors
besides Google.

Second, the conceptual depth
of Go development
done outside the original team
is expanding.

Immediately after the open source release,
one of the first large contributions
was the port to Microsoft Windows,
started by Hector Chu
and completed by Alex Brainman and others.
More contributors ported Go
to other operating systems.
Even more contributors
rewrote most of our numeric code
to be faster or more precise or both.
These were all important contributions,
and very much appreciated,
but
for the most part
they did not involve new designs.

More recently,
a group of contributors led by Aram Hăvărneanu
ported Go to the ARM 64 architecture,
This was the first architecture port
by contributors outside Google.
This is significant, because
in general
support for a new architecture
requires more design work
than support for a new operating system.
There is more variation between architectures
than between operating systems.

Another example is the introduction
over the past few releases
of preliminary support
for building Go programs using shared libraries.
This feature is important for many Linux distributions
but not as important for Google,
because we deploy static binaries.
We have been helping guide the overall strategy,
but most of the design
and nearly all of the implementation
has been done by contributors outside Google,
especially Michael Hudson-Doyle.

My last example is the go command's
approach to vendoring.
I define vendoring as
copying source code for external dependencies
into your tree
to make sure that they don't disappear
or change underfoot.

Vendoring is not a problem Google suffers,
at least not the way the rest of the world does.
We copy open source libraries we want to use
into our shared source tree,
record what version we copied,
and only update the copy
when there is a need to do so.
We have a rule
that there can only be one version
of a particular library in the source tree,
and it's the job of whoever wants to upgrade that library
to make sure it keeps working as expected
by the Google code that depends on it.
None of this happens often.
This is the lazy approach to vendoring.

In contrast, most projects outside Google
take a more eager approach,
importing and updating code
using automated tools
and making sure that they are
always using the latest versions.

Because Google has relatively little experience
with this vendoring problem,
we left it to users outside Google to develop solutions.
Over the past five years,
people have built a series of tools.
The main ones in use today are
Keith Rarick's godep,
Owen Ou's nut,
and the gb-vendor plugin for Dave Cheney's gb,

There are two problems with the current situation.
The first is that these tools
are not compatible
out of the box
with the go command's "go get".
The second is that the tools
are not even compatible with each other.
Both of these problems
fragment the developer community by tool.

Last fall, we started a public design discussion
to try to build consensus on
some basics about
how these tools all operate,
so that they can work alongside "go get"
and each other.

Our basic proposal was that all tools agree
on the approach of rewriting import paths during vendoring,
to fit with "go get"'s model,
and also that all tools agree on a file format
describing the source and version of the copied code,
so that the different vendoring tools
can be used together
even by a single project.
If you use one today,
you should still be able to use another tomorrow.

Finding common ground in this way
was very much in the spirit of Do Less, Enable More.
If we could build consensus
about these basic semantic aspects,
that would enable "go get" and all these tools to interoperate,
and it would enable switching between tools,
the same way that
agreement about how Go programs
are stored in text files
enables the Go compiler and all text editors to interoperate.
So we sent out our proposal for common ground.

Two things happened.

First, Daniel Theophanes
started a vendor-spec project on GitHub
with a new proposal
and took over coordination and design
of the spec for vendoring metadata.

Second, the community spoke
with essentially one voice
to say that
rewriting import paths during vendoring
was not tenable.
Vendoring works much more smoothly
if code can be copied without changes.

Keith Rarick posted an alternate proposal
for a minimal change to the go command
to support vendoring without rewriting import paths.
Keith's proposal was configuration-free
and fit in well with the rest of the go command's approach.
That proposal will ship
as an experimental feature in Go 1.5
and likely enabled by default in Go 1.6.
And I believe that the various vendoring tool authors
have agreed to adopt Daniel's spec once it is finalized.

The result
is that at the next Gophercon
we should have broad interoperability
between vendoring tools and the go command,
and the design to make that happen
was done entirely by contributors
outside the original Go team.

Not only that,
the Go team's proposal for how to do this
was essentially completely wrong.
The Go community told us that
very clearly.
We took that advice,
and now there's a plan for vendoring support
that I believe
everyone involved is happy with.

This is also a good example
of our general approach to design.
We try not to make any changes to Go
until we feel there is broad consensus
on a well-understood solution.
For vendoring,
feedback and design
from the Go community
was critical to reaching that point.

This general trend
toward both code and design
coming from the broader Go community
is important for Go.
You, the broader Go community,
know what is working
and what is not
in the environments where you use Go.
We at Google don't.
More and more,
we will rely on your expertise,
and we will try to help you develop
designs and code
that extend Go to be useful in more settings
and fit well with Go's original vision.
At the same time,
we will continue to wait
for broad consensus
on well-understood solutions.

This brings me to my last point.

## Code of Conduct

I've argued that Go must be open,
and that Go needs your help.

But in fact Go needs everyone's help.
And everyone isn't here.

Go needs ideas from as many people as possible.

To make that a reality,
the Go community needs to be
as inclusive,
welcoming,
helpful,
and respectful as possible.

The Go community is large enough now that,
instead of assuming that everyone involved
knows what is expected,
I and others believe that it makes sense
to write down those expectations explicitly.
Much like the Go spec
sets expectations for all Go compilers,
we can write a spec
setting expectations for our behavior
in online discussions
and in offline meetings like this one.

Like any good spec,
it must be general enough
to allow many implementations
but specific enough
that it can identify important problems.
When our behavior doesn't meet the spec,
people can point that out to us,
and we can fix the problem.
At the same time,
it's important to understand that
this kind of spec
cannot be as precise as a language spec.
We must start with the assumption
that we will all be reasonable in applying it.

This kind of spec
is often referred to as
a Code of Conduct.
Gophercon has one,
which we've all agreed to follow
by being here,
but the Go community does not.
I and others
believe the Go community
needs a Code of Conduct.

But what should it say?

I believe
the most important
overall statement we can make
is that
if you want to use or discuss Go,
then you are welcome here,
in our community.
That is the standard
I believe we aspire to.

If for no other reason
(and, to be clear, there are excellent other reasons),
Go needs as large a community as possible.
To the extent that behavior
limits the size of the community,
it holds Go back.
And behavior can easily
limit the size of the community.

The tech community in general
and the Go community in particular
is skewed toward people who communicate bluntly.
I don't believe this is fundamental.
I don't believe this is necessary.
But it's especially easy to do
in online discussions like email and IRC,
where plain text is not supplemented
by the other cues and signals we have
in face-to-face interactions.

For example, I have learned
that when I am pressed for time
I tend to write fewer words,
with the end result that
my emails seem not just hurried
but blunt, impatient, even dismissive.
That's not how I feel,
but it's how I can come across,
and that impression can be enough
to make people think twice
about using or contributing
to Go.
I realized I was doing this
when some Go contributors
sent me private email to let me know.
Now, when I am pressed for time,
I pay extra attention to what I'm writing,
and I often write more than I naturally would,
to make sure
I'm sending the message I intend.

I believe
that correcting the parts
of our everyday interactions,
intended or not,
that drive away potential users and contributors
is one of the most important things
we can all do
to make sure the Go community
continues to grow.
A good Code of Conduct can help us do that.

We have no experience writing a Code of Conduct,
so we have been reading existing ones,
and we will probably adopt an existing one,
perhaps with minor adjustments.
The one I like the most is the Django Code of Conduct,
which originated with another project called SpeakUp!
It is structured as an elaboration of a list of
reminders for everyday interaction.

"Be friendly and patient.
Be welcoming.
Be considerate.
Be respectful.
Be careful in the words that you choose.
When we disagree, try to understand why."

I believe this captures the tone we want to set,
the message we want to send,
the environment we want to create
for new contributors.
I certainly want to be
friendly,
patient,
welcoming,
considerate,
and respectful.
I won't get it exactly right all the time,
and I would welcome a helpful note
if I'm not living up to that.
I believe most of us
feel the same way.

I haven't mentioned
active exclusion based on
or disproportionately affecting
race, gender, disability,
or other personal characteristics,
and I haven't mentioned harassment.
For me,
it follows from what I just said
that exclusionary behavior
or explicit harassment
is absolutely unacceptable,
online and offline.
Every Code of Conduct says this explicitly,
and I expect that ours will too.
But I believe the SpeakUp! reminders
about everyday interactions
are an equally important statement.
I believe that
setting a high standard
for those everyday interactions
makes extreme behavior
that much clearer
and easier to deal with.

I have no doubts that
the Go community can be
one of the most
friendly,
welcoming,
considerate,
and
respectful communities
in the tech industry.
We can make that happen,
and it will be
a benefit and credit to us all.

Andrew Gerrand
has been leading the effort
to adopt an appropriate Code of Conduct
for the Go community.
If you have suggestions,
or concerns,
or experience with Codes of Conduct,
or want to be involved,
please find Andrew or me
during the conference.
If you'll still be here on Friday,
Andrew and I are going to block off
some time for Code of Conduct discussions
during Hack Day.

Again, we don't know
where the next great idea will come from.
We need all the help we can get.
We need a large, diverse Go community.

## Thank You

I consider the many people
releasing software for download using “go get,”
sharing their insights via blog posts,
or helping others on the mailing lists or IRC
to be part of this broad open source effort,
part of the Go community.
Everyone here today is also part of that community.

Thank you in advance
to the presenters
who over the next few days
will take time to share their experiences
using and extending Go.

Thank you in advance
to all of you in the audience
for taking the time to be here,
to ask questions,
and to let us know
how Go is working for you.
When you go back home,
please continue to share what you've learned.
Even if you don't use Go
for daily work,
we'd love to see what's working for Go
adopted in other contexts,
just as we're always looking for good ideas
to bring back into Go.

Thank you all again
for making the effort to be here
and for being part of the Go community.

For the next few days, please:
tell us what we're doing right,
tell us what we're doing wrong,
and help us all work together
to make Go even better.

Remember to
be friendly,
patient,
welcoming,
considerate,
and respectful.

Above all, enjoy the conference.
