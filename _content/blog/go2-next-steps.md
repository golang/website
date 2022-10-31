---
title: Next steps toward Go 2
date: 2019-06-26
by:
- Robert Griesemer, for the Go team
tags:
- go2
- proposals
- community
summary: What Go 2 language changes should we include in Go 1.14?
---

## Status

We’re well on the way towards the release of Go 1.13,
hopefully in early August of this year.
This is the first release that will include concrete changes
to the language (rather than just minor adjustments to the spec),
after a longer moratorium on any such changes.

To arrive at these language changes,
we started out with a small set of viable proposals,
selected from the much larger list of
[Go 2 proposals](https://github.com/golang/go/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+label%3AGo2+label%3AProposal),
per the new proposal evaluation process outlined
in the
“[Go 2, here we come!](https://blog.golang.org/go2-here-we-come)” blog post.
We wanted our initial selection of proposals
to be relatively minor and mostly uncontroversial,
to have a reasonably high chance of having them
make it through the process.
The proposed changes had to be backward-compatible
to be minimally disruptive since
[modules](https://blog.golang.org/using-go-modules),
which eventually will allow module-specific language version selection,
are not the default build mode quite yet.
In short, this initial round of changes was more about
getting the ball rolling again and gaining experience
with the new process, rather than tackling big issues.

Our
[original list of proposals](https://blog.golang.org/go2-here-we-come) –
[general Unicode identifiers](/issue/20706),
[binary integer literals](/issue/19308),
[separators for number literals](/issue/28493),
[signed integer shift counts](/issue/19113) –
got both trimmed and expanded.
The general Unicode identifiers didn’t make the cut
as we didn’t have a concrete design document in place in time.
The proposal for binary integer literals was expanded significantly
and led to a comprehensive overhaul and modernization of
[Go’s number literal syntax](/design/19308-number-literals).
And we added the Go 2 draft design proposal on
[error inspection](/design/go2draft-error-inspection),
which has been
[partially accepted](/issue/29934#issuecomment-489682919).

With these initial changes in place for Go 1.13,
it’s now time to look forward to Go 1.14
and determine what we want to tackle next.

## Proposals for Go 1.14

The goals we have for Go today are the same as in 2007: to
[make software development scale](https://blog.golang.org/toward-go2).
The three biggest hurdles on this path to improved scalability for Go are
package and version management,
better error handling support,
and generics.

With Go module support getting increasingly stronger,
support for package and version management is being addressed.
This leaves better error handling support and generics.
We have been working on both of these and presented
[draft designs](/design/go2draft)
at last year’s GopherCon in Denver.
Since then we have been iterating those designs.
For error handling, we have published a concrete,
significantly revised and simplified proposal (see below).
For generics, we are making progress, with a talk
(“Generics in Go” by Ian Lance Taylor)
[coming up](https://www.gophercon.com/agenda/session/49028)
at this year’s GopherCon in San Diego,
but we have not reached the concrete proposal stage yet.

We also want to continue with smaller
improvements to the language.
For Go 1.14, we have selected the following proposals:

[\#32437](/issue/32437).
A built-in Go error check function, “try”
([design doc](/design/32437-try-builtin)).

This is our concrete proposal for improved error handling.
While the proposed, fully backwards-compatible language extension
is minimal, we expect an outsize impact on error handling code.
This proposal has already attracted an enormous amount of comments,
and it’s not easy to follow up.
We recommend starting with the
[initial comment](/issue/32437#issue-452239211)
for a quick outline and then to read the detailed design doc.
The initial comment contains a couple of links leading to summaries
of the feedback so far.
Please follow the feedback recommendations
(see the “Next steps” section below) before posting.

[\#6977](/issue/6977).
Allow embedding overlapping interfaces
([design doc](/design/6977-overlapping-interfaces)).

This is an old, backwards-compatible proposal for making interface embedding more tolerant.

[\#32479](/issue/32479) Diagnose `string(int)` conversion in `go vet`.

The `string(int)` conversion was introduced early in Go for convenience,
but it is confusing to newcomers (`string(10)` is `"\n"` not `"10"`)
and not justified anymore now that the conversion is available
in the `unicode/utf8` package.
Since removing this conversion is not a backwards-compatible change,
we propose to start with a `vet` error instead.

[\#32466](/issue/32466) Adopt crypto principles
([design doc](/design/cryptography-principles)).

This is a request for feedback on a set of design principles for
cryptographic libraries that we would like to adopt.
See also the related
[proposal to remove SSLv3 support](/issue/32716)
from `crypto/tls`.

## Next steps

We are actively soliciting feedback on all these proposals.
We are especially interested in fact-based evidence
illustrating why a proposal might not work well in practice,
or problematic aspects we might have missed in the design.
Convincing examples in support of a proposal are also very helpful.
On the other hand, comments containing only personal opinions
are less actionable:
we can acknowledge them but we can’t address them
in any constructive way.
Before posting, please take the time to read the detailed
design docs and prior feedback or feedback summaries.
Especially in long discussions, your concern may have already
been raised and discussed in earlier comments.

Unless there are strong reasons to not even proceed into the
experimental phase with a given proposal,
we are planning to have all these implemented at the
start of the
[Go 1.14 cycle](/wiki/Go-Release-Cycle)
(beginning of August, 2019)
so that they can be evaluated in practice.
Per the
[proposal evaluation process](https://blog.golang.org/go2-here-we-come),
the final decision will be
made at the end of the development cycle (beginning of November, 2019).

Thank you for helping make Go a better language!
