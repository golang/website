---
title: A Proposal for Package Versioning in Go
date: 2018-03-26
by:
- Russ Cox
tags:
- tools
- versioning
summary: Proposing official support for package versioning in Go, using Go modules.
---

## Introduction

Eight years ago, the Go team introduced `goinstall`
(which led to `go get`)
and with it the decentralized, URL-like import paths
that Go developers are familiar with today.
After we released `goinstall`, one of the first questions people asked
was how to incorporate version information.
We admitted we didn’t know.
For a long time, we believed that the problem of package versioning
would be best solved by an add-on tool,
and we encouraged people to create one.
The Go community created many tools with different approaches.
Each one helped us all better understand the problem,
but by mid-2016 it was clear that there were now too many solutions.
We needed to adopt a single, official tool.

After a community discussion started at GopherCon in July 2016 and continuing into the fall,
we all believed the answer would be to follow the package versioning approach
exemplified by Rust’s Cargo, with tagged semantic versions,
a manifest, a lock file, and a
[SAT solver](https://research.swtch.com/version-sat) to decide which versions to use.
Sam Boyer led a team to create Dep, which followed this rough plan,
and which we intended to serve as the model for `go` command integration.
But as we learned more about the implications of the Cargo/Dep approach,
it became clear to me that Go would benefit from changing
some of the details, especially concerning backwards compatibility.

## The Impact of Compatibility

The most important new feature of
[Go 1](https://blog.golang.org/preview-of-go-version-1)
was not a language feature.
It was Go 1’s emphasis on backwards compatibility.
Until that point we’d issued stable release
snapshots approximately monthly,
each with significant incompatible changes.
We observed significant acceleration in interest and adoption
immediately after the release of Go 1.
We believe that the
[promise of compatibility](/doc/go1compat.html)
made developers feel much more comfortable relying on
Go for production use
and is a key reason that Go is popular today.
Since 2013 the
[Go FAQ](/doc/faq#get_version)
has encouraged package developers to provide their own
users with similar expectations of compatibility.
We call this the _import compatibility rule_:
“If an old package and a new package have the same import path,
the new package must be backwards compatible with the old package.”

Independently,
[semantic versioning](http://semver.org/)
has become the _de facto_
standard for describing software versions in many language communities,
including the Go community.
Using semantic versioning, later versions are expected to be
backwards-compatible with earlier versions,
but only within a single major version:
v1.2.3 must be compatible with v1.2.1 and v1.1.5,
but v2.3.4 need not be compatible with any of those.

If we adopt semantic versioning for Go packages,
as most Go developers expect,
then the import compatibility rule requires that
different major versions must use different import paths.
This observation led us to _semantic import versioning_,
in which versions starting at v2.0.0 include the major
version in the import path: `my/thing/v2/sub/pkg`.

A year ago I strongly believed that whether to include
version numbers in import paths was largely a matter of taste,
and I was skeptical that having them was particularly elegant.
But the decision turns out to be a matter not of taste but of logic:
import compatibility and semantic versioning together require
semantic import versioning.
When I realized this, the logical necessity surprised me.

I was also surprised to realize that
there is a second, independent logical route to
semantic import versioning:
[gradual code repair](/talks/2016/refactor.article)
or partial code upgrades.
In a large program, it’s unrealistic to expect all packages in the program
to update from v1 to v2 of a particular dependency at the same time.
Instead, it must be possible for some of the program to keep using v1
while other parts have upgraded to v2.
But then the program’s build, and the program’s final binary,
must include both v1 and v2 of the dependency.
Giving them the same import path would lead to confusion,
violating what we might call the _import uniqueness rule_:
different packages must have different import paths.
The only way to have
partial code upgrades, import uniqueness, _and_ semantic versioning
is to adopt
semantic import versioning as well.

It is of course possible to build systems that use semantic versioning
without semantic import versioning,
but only by giving up either partial code upgrades or import uniqueness.
Cargo allows partial code upgrades by
giving up import uniqueness:
a given import path can have different meanings
in different parts of a large build.
Dep ensures import uniqueness by
giving up partial code upgrades:
all packages involved in a large build must find
a single agreed-upon version of a given dependency,
raising the possibility that large programs will be unbuildable.
Cargo is right to insist on partial code upgrades,
which are critical to large-scale software development.
Dep is equally right to insist on import uniqueness.
Complex uses of Go’s current vendoring support can violate import uniqueness.
When they have, the resulting problems have been quite challenging
for both developers and tools to understand.
Deciding between partial code upgrades
and import uniqueness
requires predicting which will hurt more to give up.
Semantic import versioning lets us avoid the choice
and keep both instead.

I was also surprised to discover how much
import compatibility simplifies version selection,
which is the problem of deciding which package versions to use for a given build.
The constraints of Cargo and Dep make version selection
equivalent to
[solving Boolean satisfiability](https://research.swtch.com/version-sat),
meaning it can be very expensive to determine whether
a valid version configuration even exists.
And then there may be many valid configurations,
with no clear criteria for choosing the “best” one.
Relying on import compatibility can instead let Go use
a trivial, linear-time algorithm
to find the single best configuration, which always exists.
This algorithm,
which I call
[_minimal version selection_](https://research.swtch.com/vgo-mvs),
in turn eliminates the need for separate lock and manifest files.
It replaces them with a single, short configuration file,
edited directly by both developers and tools,
that still supports reproducible builds.

Our experience with Dep demonstrates the impact of compatibility.
Following the lead of Cargo and earlier systems,
we designed Dep to give up import compatibility
as part of adopting semantic versioning.
I don’t believe we decided this deliberately;
we just followed those other systems.
The first-hand experience of using Dep helped us
better understand exactly how much complexity
is created by permitting incompatible import paths.
Reviving the import compatibility rule
by introducing semantic import versioning
eliminates that complexity,
leading to a much simpler system.

## Progress, a Prototype, and a Proposal

Dep was released in January 2017.
Its basic model—code tagged with
semantic versions, along with a configuration file that
specified dependency requirements—was
a clear step forward from most of the Go vendoring tools,
and converging on Dep itself was also a clear step forward.
I wholeheartedly encouraged its adoption,
especially to help developers get used to thinking about Go package versions,
both for their own code and their dependencies.
While Dep was clearly moving us in the right direction, I had lingering concerns
about the complexity devil in the details.
I was particularly concerned about Dep
lacking support for gradual code upgrades in large programs.
Over the course of 2017, I talked to many people,
including Sam Boyer and the rest of the
package management working group,
but none of us could see any clear way to reduce the complexity.
(I did find many approaches that added to it.)
Approaching the end of the year,
it still seemed like SAT solvers and unsatisfiable builds
might be the best we could do.

In mid-November, trying once again to work through
how Dep could support gradual code upgrades,
I realized that our old advice about import compatibility
implied semantic import versioning.
That seemed like a real breakthrough.
I wrote a first draft of what became my
[semantic import versioning](https://research.swtch.com/vgo-import)
blog post,
concluding it by suggesting that Dep adopt the convention.
I sent the draft to the people I’d been talking to,
and it elicited very strong responses:
everyone loved it or hated it.
I realized that I needed to work out more of the
implications of semantic import versioning
before circulating the idea further,
and I set out to do that.

In mid-December, I discovered that import compatibility
and semantic import versioning together allowed
cutting version selection down to [minimal version selection](https://research.swtch.com/vgo-mvs).
I wrote a basic implementation to be sure I understood it,
I spent a while learning the theory behind why it was so simple,
and I wrote a draft of the post describing it.
Even so, I still wasn’t sure the approach would be practical
in a real tool like Dep.
It was clear that a prototype was needed.

In January, I started work on a simple `go` command wrapper
that implemented semantic import versioning
and minimal version selection.
Trivial tests worked well.
Approaching the end of the month,
my simple wrapper could build Dep,
a real program that made use of many versioned packages.
The wrapper still had no command-line interface—the fact that
it was building Dep was hard-coded in a few string constants—but
the approach was clearly viable.

I spent the first three weeks of February turning the
wrapper into a full versioned `go` command, `vgo`;
writing drafts of a
[blog post series introducing `vgo`](https://research.swtch.com/vgo);
and discussing them with
Sam Boyer, the package management working group,
and the Go team.
And then I spent the last week of February finally
sharing `vgo` and the ideas behind it with the whole Go community.

In addition to the core ideas of import compatibility,
semantic import versioning, and minimal version selection,
the `vgo` prototype introduces a number of smaller
but significant changes motivated by eight years of
experience with `goinstall` and `go get`:
the new concept of a [Go module](https://research.swtch.com/vgo-module),
which is a collection of packages versioned as a unit;
[verifiable and verified builds](https://research.swtch.com/vgo-repro);
and
[version-awareness throughout the `go` command](https://research.swtch.com/vgo-cmd),
enabling work outside `$GOPATH`
and the elimination of (most) `vendor` directories.

The result of all of this is the [official Go proposal](/design/24301-versioned-go),
which I filed last week.
Even though it might look like a complete implementation,
it’s still just a prototype,
one that we will all need to work together to complete.
You can download and try the `vgo` prototype from [golang.org/x/vgo](https://golang.org/x/vgo),
and you can read the
[Tour of Versioned Go](https://research.swtch.com/vgo-tour)
to get a sense of what using `vgo` is like.

## The Path Forward

The proposal I filed last week is exactly that: an initial proposal.
I know there are problems with it that the Go team and I can’t see,
because Go developers use Go in many clever ways that we don’t know about.
The goal of the proposal feedback process is for us all to work together
to identify and address the problems in the current proposal,
to make sure that the final implementation that ships in a future
Go release works well for as many developers as possible.
Please point out problems on the [proposal discussion issue](/issue/24301).
I will keep the
[discussion summary](/issue/24301#issuecomment-371228742)
and
[FAQ](/issue/24301#issuecomment-371228664)
updated as feedback arrives.

For this proposal to succeed, the Go ecosystem as a
whole—and in particular today’s major Go projects—will need to
adopt the import compatibility rule and semantic import versioning.
To make sure that can happen smoothly,
we will also be conducting user feedback sessions
by video conference with projects that have questions about
how to incorporate the new versioning proposal into their code bases
or have feedback about their experiences.
If you are interested in participating in such a session,
please email Steve Francia at spf@golang.org.

We’re looking forward to (finally!) providing the Go community with a single, official answer
to the question of how to incorporate package versioning into `go get`.
Thanks to everyone who helped us get this far, and to everyone who will help us going forward.
We hope that, with your help, we can ship something that Go developers will love.
