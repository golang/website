---
title: Go 2, here we come!
date: 2018-11-29
by:
- Robert Griesemer
tags:
- go2
- proposals
- community
summary: How Go 2 proposals will be evaluated, selected, and shipped.
---

## Background

At GopherCon 2017, Russ Cox officially started the thought process on the
next big version of Go with his talk [The Future of Go](https://www.youtube.com/watch?v=0Zbh_vmAKvk)
([blog post](https://blog.golang.org/toward-go2)). We have
called this future language informally Go 2, even though we understand now
that it will arrive in incremental steps rather than with a big bang and a
single major release. Still, Go 2 is a useful moniker, if only to have a way
to talk about that future language, so let’s keep using it for now.

A major difference between Go 1 and Go 2 is who is going to influence the
design and how decisions are made. Go 1 was a small team effort with modest
outside influence; Go 2 will be much more community-driven.
After almost 10 years of exposure, we have
learned a lot about the language and libraries that we didn’t know in the
beginning, and that was only possible through feedback from the Go community.

In 2015 we introduced the [proposal process](/s/proposal)
to gather a specific kind of feedback: proposals for language and library
changes. A committee composed of senior Go team members has been reviewing,
categorizing, and deciding on incoming proposals on a regular basis. That
has worked pretty well, but as part of that process we have ignored all
proposals that are not backward-compatible, simply labeling them Go 2 instead.
In 2017 we also stopped making any kind of incremental backward-compatible
language changes, however small, in favor of a more comprehensive plan that
takes the bigger picture of Go 2 into account.

It is now time to act on the Go 2 proposals, but to do this we first need a plan.

## Status

At the time of writing, there are around 120
[open issues labeled Go 2 proposal](https://github.com/golang/go/issues?page=1&q=is%3Aissue+is%3Aopen+label%3Aproposal+label%3AGo2&utf8=%E2%9C%93).
Each of them proposes a significant library or language change, often one
that does not satisfy the existing
[Go 1 compatibility guarantee](/doc/go1compat).
Ian Lance Taylor and I
have been working through these proposals and categorized them
([Go2Cleanup](https://github.com/golang/go/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+label%3Aproposal+label%3AGo2+label%3AGo2Cleanup),
[NeedsDecision](https://github.com/golang/go/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+label%3Aproposal+label%3AGo2+label%3ANeedsDecision),
etc.) to get an idea of what’s there and to make it easier to
proceed with them. We also merged related proposals and closed the ones which
seemed clearly out of the scope of Go, or were otherwise unactionable.

Ideas from the remaining proposals will likely influence Go 2’s libraries
and languages. Two major themes have emerged early on: support for better
error handling, and generics. [Draft designs](https://blog.golang.org/go2draft)
for these two areas have been
published at this year’s GopherCon, and more exploration is needed.

But what about the rest? We are [constrained](https://blog.golang.org/toward-go2)
by the fact that we now have
millions of Go programmers and a large body of Go code, and we need to
bring it all along, lest we risk a split ecosystem. That means we cannot
make many changes, and the changes we are going to make need to be chosen
carefully. To make progress, we are implementing a new proposal evaluation
process for these significant potential changes.

## Proposal evaluation process

The purpose of the proposal evaluation process is to collect feedback on
a small number of select proposals such that a final decision can be made.
The process runs more or less in parallel to a release cycle and consists
of the following steps:

1. _Proposal selection_. The Go team selects a small number of
[Go 2 proposals](https://github.com/golang/go/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+label%3AGo2+label%3AProposal)
that seem worth considering for acceptance, without making a final decision.
See below for more on the selection criteria.

2. _Proposal feedback_. The Go team sends out an announcement listing the selected
proposals. The announcement explains to the community the tentative intent to
move forward with the selected proposals and to collect feedback for each
of them. This gives the community a chance to make suggestions and express
concerns.

3. _Implementation_. Based on that feedback, the proposals are implemented.
The target for these significant language and library changes is to have
them ready to submit on day 1 of an upcoming release cycle.

4. _Implementation feedback_. During the development cycle, the Go team and
community have a chance to experiment with the new features and collect
further feedback.

5. _Launch decision_. At the end of the three month
[development cycle](https://github.com/golang/go/wiki/Go-Release-Cycle)
(just when starting the three month repo freeze before a release), and
based on the experience and feedback gathered during the release cycle,
the Go team makes the final decision about whether to ship each change.
This provides an opportunity to consider whether the change has delivered
the expected benefits or created any unexpected costs. Once shipped, the
changes become part of the language and libraries. Excluded proposals may
go back to the drawing board or may be declined for good.

With two rounds of feedback, this process is slanted towards declining
proposals, which will hopefully prevent feature creep and help with
keeping the language small and clean.

We can’t go through this process for each of the open Go 2
proposals, there are simply too many of them. That’s where the selection
criteria come into play.

## Proposal selection criteria

A proposal must at the very least:

1. _address an important issue for many people_,

2. _have minimal impact on everybody else_, and

3. _come with a clear and well-understood solution_.

Requirement 1 ensures that any changes we make help as many Go developers
as possible (make their code more robust, easier to write, more likely to
be correct, and so on), while requirement 2 ensures we are careful to hurt
as few developers as possible, whether by breaking their programs or causing
other churn. As a rule of thumb, we should aim to help at least ten times as
many developers as we hurt with a given change. Changes that don't affect
real Go usage are a net zero benefit put up against a significant
implementation cost and should be avoided.

Without requirement 3 we don’t have an implementation of the proposal.
For instance, we believe that some form of genericity might solve an important
issue for a lot of people, but we don’t yet have a clear and well-understood
solution. That’s fine, it just means that the proposal needs to go back to
the drawing board before it can be considered.

## Proposals

We feel that this is a good plan that should serve us well but it is important
to understand that this is only a starting point. As the process is used we will
discover the ways in which it fails to work well and we will refine it as needed.
The critical part is that until we use it in practice we won't know how to improve it.

A safe place to start is with a small number of backward-compatible language
proposals. We haven’t done language changes for a long time, so this gets us
back into that mode. Also, the changes won’t require us worrying about
breaking existing code, and thus they serve as a perfect trial balloon.

With all that said, we propose the following selection of Go 2 proposals for
the Go 1.13 release (step 1 in the proposal evaluation process):

1. [_\#20706_](https://github.com/golang/go/issues/20706) _General Unicode identifiers based on_ [_Unicode TR31_](http://unicode.org/reports/tr31/):
This addresses an important issue for Go programmers using non-Western alphabets
and should have little if any impact on anyone else. There are normalization
questions which we need to answer and where community feedback will be
important, but after that the implementation path is well understood.
Note that identifier export rules will not be affected by this.

2. [_\#19308_](https://github.com/golang/go/issues/19308), [_\#28493_](https://github.com/golang/go/issues/28493) _Binary integer literals and support for \_ in number literals_:
These are relatively minor changes that seem hugely popular among many
programmers. They may not quite reach the threshold of solving an
“important issue” (hexadecimal numbers have worked well so far) but they
bring Go up to par with most other languages in this respect and relieve
a pain point for some programmers. They have minimal impact on others who
don’t care about binary integer literals or number formatting, and the
implementation is well understood.

3. [_\#19113_](https://github.com/golang/go/issues/19113) _Permit signed integers as shift counts_:
An estimated 38% of all non-constant shifts require an (artificial) uint
conversion (see the issue for a more detailed break-down). This proposal
will clean up a lot of code, get shift expressions better in sync with index
expressions and the built-in functions cap and len. It will mostly have a
positive impact on code. The implementation is well understood.

## Next steps

With this blog post we have executed the first step and started the second
step of the proposal evaluation process. It’s now up to you, the
Go community, to provide feedback on the issues listed above.

For each proposal for which we have clear and approving feedback, we will
move forward with the implementation (step 3 in the process). Because we
want the changes implemented on the first day of the next release cycle
(tentatively Feb. 1, 2019) we may start the implementation a bit early
this time to leave time for two full months of feedback (Dec. 2018,
Jan. 2019).

For the 3-month development cycle (Feb. to May 2019) the chosen features
are implemented and available at tip and everybody will have a chance to
gather experience with them. This provides another opportunity for feedback
(step 4 in the process).

Finally, shortly after the repo freeze (May 1, 2019), the Go team makes the
final decision whether to keep the new features for good (and include them
in the Go 1 compatibility guarantee), or whether to abandon them (final
step in the process).

(Since there is a real chance that a feature may need to be removed just
when we freeze the repo, the implementation will need to be such that the
feature can be disabled without destabilizing the rest of the system.
For language changes that may mean that all feature-related code is
guarded by an internal flag.)

This will be the first time that we have followed this process, hence the
repo freeze will also be a good moment to reflect on the process and to
adjust it if necessary. Let’s see how it goes.

Happy evaluating!
