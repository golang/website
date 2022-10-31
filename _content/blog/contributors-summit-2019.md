---
title: Contributors Summit 2019
date: 2019-08-15
by:
- Carmen Andoh and contributors
tags:
- community
summary: Reporting from the Go Contributor Summit at GopherCon 2019.
---

## Introduction

For the third year in a row, the Go team and contributors convened
the day before GopherCon to discuss and plan for the future of the Go project.
The event included self-organizing into breakout groups,
a town-hall style discussion about the proposal process in the morning,
and afternoon break-out roundtable discussions
based on topics our contributors chose.
We asked five contributors to write about their experience
in various discussions at this year’s summit.

{{image "contributors-summit-2019/group.jpg" 800}}

_(Photo by Steve Francia.)_

## Compiler and Runtime (report by Lynn Boger)

The Go contributors summit was a great opportunity
to meet and discuss topics and ideas with others who also contribute to Go.

The day started out with a time to meet everyone in the room.
There was a good mix of the core Go team
and others who actively contribute to Go.
From there we decided what topics were of interest
and how to split the big group into smaller groups.
My area of interest is the compiler, so I joined that group
and stayed with them for most of the time.

At our first meeting, a long list of topics were brought up
and as a result the compiler group decided to keep meeting throughout the day.
I had a few topics of interest that I shared and many that others suggested
were also of interest to me.
Not all items on the list were discussed in detail;
here is my list of those topics which had the most interest and discussion,
followed by some brief comments that were made on other topics.

**Binary size**.
There was a concern expressed about binary size,
especially that it continues to grow with each release.
Some possible reasons were identified such as increased inlining and other optimizations.
Most likely there is a set of users who want small binaries,
and another group who wants the best performance possible and maybe some don’t care.
This led to the topic of TinyGo, and it was noted that TinyGo was not a full implementation of Go
and that it is important to keep TinyGo from diverging from Go and splitting the user base.
More investigation is required to understand the need among users and the exact reasons
contributing to the current size.
If there are opportunities to reduce the size without affecting performance,
those changes could be made, but if performance were affected
some users would prefer better performance.

**Vector assembly**.
How to leverage vector assembly in Go was discussed for a while
and has been a topic of interest in the past.
I have split this into three separate possibilities, since they all relate to the use of vector instructions,
but the way they are used are different, starting with the topic of vector assembly.
This is another case of a compiler trade off.

For most targets, there are critical functions in standard packages
such as crypto, hash, math and others, where the use of assembly is necessary
to get the best possible performance; however having large functions
written in assembly makes them difficult to support and maintain
and could require different implementations for each target platform.
One solution is to make use of macro assembly or other
high-level generation techniques to make the vector assembly easier to read and understand.

Another side to this question is whether the Go compiler can
directly generate SIMD vector instructions when compiling a Go source file,
by enhancing the Go compiler to transform code sequences to “simdize”
the code to make use of vector instructions.
Implementing SIMD in the Go compiler would add complexity and compile time,
and might not always result in code that performs better.
The way the code is transformed could in some cases depend
on the target platform so that would not be ideal.

Another way to leverage vector instructions in Go is to provide a way
to make it easier to make use of vector instructions from within the Go source code.
Topics discussed were intrinsics, or implementations that exist in other compilers like Rust.
In gcc some platforms provide inline asm, and Go possibly could provide this capability,
but I know from experience that intermixing inline asm with Go code adds complexity
to the compiler in terms of tracking register use and debugging.
It allows the user to do things the compiler might not expect or want,
and it does add an extra level of complexity.
It could be inserted in places that are not ideal.

In summary, it is important to provide a way to leverage
the available vector instructions, and make it easier and safer to write.
Where possible, functions use as much Go code as possible,
and potentially find a way to use high level assembly.
There was some discussion of designing an experimental vector package
to try and implement some of these ideas.

**New calling convention**.
Several people were interested in the topic of the
[ABI changes to provide a register based calling convention](/issue/18597).
The current status was reported with details.
There was discussion on what remained to be done before it could be used.
The ABI specification needs to be written first and it was not clear when that would be done.
I know this will benefit some target platforms more than others
and a register calling convention is used in most compilers for other platforms.

**General optimizations**.
Certain optimizations that are more beneficial for some platforms other than x86 were discussed.
In particular, loop optimizations such as hoisting of invariants and strength reduction could be done
and provide more benefit on some platforms.
Potential solutions were discussed, and implementation would probably be
up to the targets that find those improvements important.

**Feedback-directed optimizations**.
This was discussed and debated as a possible future enhancement.
In my experience, it is hard to find meaningful programs to use for
collecting performance data that can later be used to optimize code.
It increases compile time and takes a lot of space to save the data
which might only be meaningful for a small set of programs.

**Pending submissions**.
A few members in the group mentioned changes they had been working on
and plan to submit soon, including improvements to makeslice, and a rewrite of rulegen.

**Compile time concerns**.
Compile time was discussed briefly. It was noted that phase timing was added to the GOSSAFUNC output.

**Compiler contributor communication**.
Someone asked if there was a need for a Go compiler mailing list.
It was suggested that we use golang-dev for that purpose,
adding compiler to the subject line to identify it.
If there is too much traffic on golang-dev, then a compiler-specific mailing list
can be considered at some later point in time.

**Community**.
I found the day very beneficial in terms of connecting with people
who have been active in the community and have similar areas of interest.
I was able to meet many people who I’ve only known by the user name
appearing in issues or mailing lists or CLs.
I was able to discuss some topics and existing issues
and get direct interactive feedback instead of waiting for online responses.
I was encouraged to write issues on problems I have seen.
These connections happened not just during this day but while
running into others throughout the conference,
having been introduced on this first day, which led to many interesting discussions.
Hopefully these connections will lead to more effective communication
and improved handling of issues and code changes in the future.

## Tools (report by Paul Jolly)

The tools breakout session during the contributor summit took an extended form,
with two further sessions on the main conference days organized by the
[golang-tools](https://github.com/golang/go/wiki/golang-tools) group.
This summary is broken down into two parts: the tools session at the contributor workshop,
and a combined report from the golang-tools sessions on the main conference days.

**Contributor summit**.
The tools session started with introductions from ~25 folks gathered,
followed by a brainstorming of topics, including:
gopls, ARM 32-bit, eval, signal, analysis, go/packages api, refactoring, pprof,
module experience, mono repo analysis, go mobile, dependencies, editor integrations,
compiler opt decisions, debugging, visualization, documentation.
A lot of people with lots of interest in lots of tools!

The session focused on two areas (all that time allowed): gopls and visualizations.
[Gopls](/wiki/gopls) (pronounced: “go please”) is an implementation of the
[Language Server Protocol (LSP)](https://langserver.org) server for Go.
Rebecca Stambler, the gopls lead author, and the rest of the Go tools team were interested
in hearing people’s experiences with gopls: stability, missing features, integrations in editors working, etc?
The general feeling was that gopls was in really good shape and working extremely well for the majority of use cases.
Integration test coverage needs to be improved, but this is a hard problem to get “right” across all editors.
We discussed a better means of users reporting gopls errors they encounter via their editor,
telemetry/diagnostics, gopls performance metrics, all subjects that got more detailed coverage
in golang-tools sessions that followed on the main conference days (see below).
A key area of discussion was how to extend gopls, e.g., in the form of
additional go/analysis vet-like checks, lint checks, refactoring, etc.
Currently there is no good solution, but it’s actively under investigation.
Conversation shifted to the very broad topic of visualizations, with a
demo-based introduction from Anthony Starks (who, incidentally, gave an excellent talk about
[Go for information displays](https://www.youtube.com/watch?v=NyDNJnioWhI) at GopherCon 2018).

**Conference days**.
The golang-tools sessions on the main conference days were a continuation of the
[monthly calls](/wiki/golang-tools) that have been happening since the group’s inception at GopherCon 2018.
Full notes are available for the
[day 1](https://docs.google.com/document/d/1-RVyttQ0ncjCpR_sRwizf-Ubedkr0Emwmk2LhnsUOmE/edit) and
[day 2](https://docs.google.com/document/d/1ZI_WqpLCB8DO6teJ3aBuXTeYD2iZZZlkDptmcY6Ja60/edit#heading=h.x9lkytc2gxmg) sessions.
These sessions were again well attended with 25-30 people at each session.
The Go tools team was there in strength (a good sign of the support being put behind this area), as was the Uber platform team.
In contrast to the contributor summit, the goal from these sessions was to come away with specific action items.

**Gopls**.
Gopls “readiness” was a major focus for both sessions.
This answer effectively boiled down to determining when it makes sense to tell
editor integrators “we have a good first cut of gopls” and then compiling a
list of “blessed” editor integrations/plugins known to work with gopls.
Central to this “certification” of editor integrations/plugins is a well-defined process
by which users can report problems they experience with gopls.
Performance and memory are not blockers for this initial “release”.
The conversation about how to extend gopls, started in the
contributor summit the day before, continued in earnest.
Despite the many obvious benefits and attractions to extending gopls
(custom go/analysis checks, linter support, refactoring, code generation…),
there isn’t a clear answer on how to implement this in a scalable way.
Those gathered agreed that this should not be seen as a blocker for the
initial “release”, but should continue to be worked on.
In the spirit of gopls and editor integrations,
Heschi Kreinick from the Go tools team brought up the topic of debugging support.
Delve has become the de facto debugger for Go and is in good shape;
now the state of debugger-editor integration needs to be established,
following a process similar to that of gopls and the “blessed” integrations.

**Go Discovery Site**.
The second golang-tools session started with an excellent introduction to
the Go Discovery Site by Julie Qiu from the Go tools team, along with a quick demo.
Julie talked about the plans for the Discovery Site: open sourcing the project,
what signals are used in search ranking, how [godoc.org](http://godoc.org/) will ultimately be replaced,
how submodules should work, how users can discover new major versions.

**Build Tags**.
Conversation then moved to build tag support within gopls.
This is an area that clearly needs to be better understood
(use cases are currently being gathered in [issue 33389](/issue/33389)).
In light of this conversation, the session wrapped up with
Alexander Zolotov from the JetBrains GoLand team suggesting that the gopls and
GoLand teams should share experience in this and more areas, given GoLand
has already gained lots of experience.

**Join Us!**
We could easily have talked about tools-related topics for days!
The good news is that the golang-tools calls will continue for the foreseeable future.
Anyone interested in Go tooling is very much encouraged to join: [the wiki](/wiki/golang-tools) has more details.

## Enterprise Use (report by Daniel Theophanes)

Actively asking after the needs of less vocal developers will be the largest challenge,
and greatest win, for the Go language. There is a large segment of programmers
who don’t actively participate in the Go community.
Some are business associates, marketers, or quality assurance who also do development.
Some will wear management hats and make hiring or technology decisions.
Others just do their job and return to their families.
And lastly, many times these developers work in businesses with strict IP protection contracts.
Even though most of these developers won’t end up directly participating in open source
or the Go community proposals, their ability to use Go depends on both.

The Go community and Go proposals need to understand the needs of these less vocal developers.
Go proposals can have a large impact on what is adopted and used.
For instance, the vendor folder and later the Go modules proxy are incredibly important
for businesses that strictly control source code and
typically have fewer direct conversations with the Go community.
Having these mechanisms allow these organizations to use Go at all.
It follows that we must not only pay attention to current Go users,
but also to developers and organizations who have considered Go,
but have chosen against it.
We need to understand these reasons.

Similarly, should the Go community pay attention to “enterprise”
environments it would unlock many additional organizations who can utilize Go.
By ensuring active directory authentication works, users who would
be forced to use a different ecosystem can keep Go on the table.
By ensuring WSDL just works, a section of users can pick Go up as a tool.
No one suggested blindly making changes to appease non-Go users.
But rather we should be aware of untapped potential and unrecognized
hindrances in the Go language and ecosystem.

While several different possibilities to actively solicit this information
from the outside were discussed, this is a problem we fundamentally need your help.
If you are in an organization that doesn’t use Go even though it was considered,
let us know why Go wasn’t chosen.
If you are in an organization where Go is only used for a subsection of programming tasks,
but not others, why isn’t it used for more? Are there specific blockers to adoption?

## Education (report by Andy Walker)

One of the roundtables I was involved in at the Contributors Summit
this year was on the topic of Go education,
specifically what kind of resources we make available
to the new Go programmer, and how we can improve them.
Present were a number of very passionate organizers, engineers and educators,
each of whom had a unique perspective on the subject,
either through tools they’d designed,
documents they’d written or workshops they’d given to developers of all stripes.

Early on, talk turned to whether or not Go makes a good first programming language.
I wasn’t sure, and advocated against it.
Go isn’t a good first language, I argued, because it isn’t intended to be.
As Rob Pike [wrote back in 2012](/talks/2012/splash.article),
“the language was designed by and for people who write—and read and debug and maintain—large software systems”.
To me, this guiding ethos is clear: Go is a deliberate response to perceived flaws
in the processes used by experienced engineers, not an attempt to create an ideal
programming language, and as such a certain basic familiarity with programming concepts is assumed.

This is evident in the official documentation at [golang.org/doc](/doc/).
It jumps right into how to install the language before passing the user on to the
[tour](/tour/), which is geared towards programmers
who are already familiar with a C-like language.
From there, they are taken to [How to Write Go Code](/doc/code.html),
which provides a very basic introduction to the classic non-module Go workspace,
before moving immediately on to writing libraries and testing.
Finally, we have [Effective Go](/doc/effective_go.html),
and a series of references including the [spec](/ref/spec),
rounded out by some examples.
These are all decent resources if you’re already familiar with a C-like language,
but they still leave a lot to be desired, and there’s nothing to be found
for the raw beginner or even someone coming directly from a language like Python.

As an accessible, interactive starting point, the tour is a natural first target
towards making the language more beginner friendly,
and I think a lot of headway can be made targeting that alone.
First, it should be the first link in the documentation,
if not the first link in the bar at the top of golang.org, front and center.
We should encourage the curious user to jump right in and start playing with the language.
We should also consider including optional introductory sections on coming
from other common languages, and the differences they are
likely to encounter in Go, with interactive exercises.
This would go a long way to helping new Go programmers in mapping
the concepts they are already familiar with onto Go.

For experienced programmers, an optional, deeper treatment should be given
to most sections in the tour, allowing them to drill down into more
detailed documentation or interactive exercises enumerating the
design decisions principles of good architecture in Go.
They should find answers to questions like:

  - Why are there so many integer types when I am encouraged to use `int` most of the time?
  - Is there ever a good reason to pick a value receiver?
  - Why is there a plain `int`, but no plain `float`?
  - What are send- and receive-only channels, and when would I use them?
  - How do I effectively compose concurrency primitives, and when would I _not_ want to use channels?
  - What is `uint` good for? Should I use it to restrict my user to positive values? Why not?

The tour should be someplace they can revisit upon finishing the first run-through
to dive more deeply into some of the more interesting choices in language design.

But we can do more. Many people seek out programming as a way to design
applications or scratch a particular itch, and they are most likely to want
to target the interface they are most familiar with: the browser.
Go does not have a good front-end story yet.
Javascript is still the only language that really provides
both a frontend and a backend environment,
but WASM is fast becoming a first-order platform,
and there are so many places we could go with that.
We could provide something like [vecty](https://github.com/gopherjs/vecty)
in [The Go Play Space](https://goplay.space/),
or perhaps [Gio](https://gioui.org/), targeting WASM, for people to get
started programming in the browser right away, inspiring their imagination,
and provide them a migration path out of our playground into
a terminal and onto GitHub.

So, is Go a good first language?
I honestly don’t know, but it’s certainly true there are a significant
number of people entering the programming profession
with Go as their starting point, and I am very interested in talking to them,
learning about their journey and their process,
and shaping the future of Go education with their input.

## Learning Platforms (report by Ronna Steinberg)

We discussed what a learning platform for Go should look like
and how we can combine global resources to effectively teach the language.
We generally agreed that teaching and learning is easier with visualization
and that a REPL is very gratifying.
We also overviewed some existing solutions for visualization with Go:
templates, Go WASM, GopherJS as well as SVG and GIFs generation.

Compiler errors not making sense to the new developer was also brought up
and we considered ideas of how to handle it, perhaps a bank of errors and how they would be useful.
One idea was a wrapper for the compiler that explains your errors to you, with examples and solutions.

A new group convened for a second round later and we focused more on
what UX should the Go learning platform have,
and if and how we can take existing materials (talks, blog posts, podcasts, etc)
from the community and organize them into a program people can learn from.
Should such a platform link to those external resources?
Embed them?
Cite them?
We agreed that a portal-like-solution (of external links to resources)
makes navigation difficult and takes away from the learning experience,
which led us to the conclusion that such contribution cannot be passive,
and contributors will likely have to opt in to have their material on the platform.
There was then much excitement around the idea of adding a voting mechanism to the platform,
effectively turning the learners into contributors, too,
and incentivizing the contributors to put their materials on the platform.

(If you are interested in helping in educational efforts for Go,
please email Carmen Andoh candoh@google.com.)

## Thank You!

Thanks to all the attendees for the excellent discussions on contributor day,
and thanks especially to Lynn, Paul, Daniel, Andy, and Ronna
for taking the time to write these reports.
