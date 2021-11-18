---
title: Proposals for Go 1.15
date: 2020-01-28
by:
- Robert Griesemer, for the Go team
tags:
- go1.15
- proposals
- community
- language
- vet
summary: For Go 1.15, we propose three minor language cleanup changes.
---

## Status

We are close to the Go 1.14 release, planned for February assuming all goes
well, with an RC1 candidate almost ready. Per the process outlined in the
[Go 2, here we come!](https://blog.golang.org/go2-here-we-come) blog post,
it is again the time in our development and release cycle to consider if and
what language or library changes we might want to include for our next release,
Go 1.15, scheduled for August of this year.

The primary goals for Go remain package and version management, better error
handling support, and generics. Module support is in good shape and getting
better with each day, and we are also making progress on the generics front
(more on that later this year). Our attempt seven months ago at providing a
better error handling mechanism, the
[`try` proposal](/issue/32437), met good support
but also strong opposition and we decided to abandon it. In its aftermath there
were many follow-up proposals, but none of them seemed convincing enough,
clearly superior to the `try` proposal, or less likely to cause similar
controversy. Thus, we have not further pursued changes to error handling
for now. Perhaps some future insight will help us to improve upon the status
quo.

## Proposals

Given that modules and generics are actively being worked on, and with error
handling changes out of the way for the time being, what other changes should
we pursue, if any? There are some perennial favorites such as requests for
enums and immutable types, but none of those ideas are sufficiently developed
yet, nor are they urgent enough to warrant a lot of attention by the Go team,
especially when also considering the cost of making a language change.

After reviewing all potentially viable proposals, and more importantly, because
we don’t want to incrementally add new features without a long-term plan, we
concluded that it is better to hold off with major changes this time. Instead
we concentrate on a couple of new `vet` checks and a minor adjustment to the
language. We have selected the following three proposals:

[\#32479](/issue/32479).
Diagnose `string(int)` conversion in `go vet`.

We were planning to get this done for the upcoming Go 1.14 release but we didn’t
get around to it, so here it is again. The `string(int)` conversion was introduced
early in Go for convenience, but it is confusing to newcomers (`string(10)` is
`"\n"` not `"10"`) and not justified anymore now that the conversion is available
in the `unicode/utf8` package.
Since [removing this conversion](/issue/3939) is
not a backwards-compatible change, we propose to start with a `vet` error instead.

[\#4483](/issue/4483).
Diagnose impossible interface-interface type assertions in `go vet`.

Currently, Go permits any type assertion `x.(T)` (and corresponding type switch case)
where the type of `x` and `T` are interfaces. Yet, if both `x` and `T` have a method
with the same name but different signatures it is impossible for any value assigned
to `x` to also implement `T`; such type assertions will always fail at runtime
(panic or evaluate to `false`). Since we know this at compile time, the compiler
might as well report an error. Reporting a compiler error in this case is not a
backwards-compatible change, thus we also propose to start with a `vet` error
instead.

[\#28591](/issue/28591).
Constant-evaluate index and slice expressions with constant strings and indices.

Currently, indexing or slicing a constant string with a constant index, or indices,
produces a non-constant `byte` or `string` value, respectively. But if all operands
are constant, the compiler can constant-evaluate such expressions and produce a
constant (possibly untyped) result. This is a fully backward-compatible change
and we propose to make the necessary adjustments to the spec and compilers.

(Correction: We found out after posting that this change is not backward-compatible;
see [comment](/issue/28591#issuecomment-579993684) for details.)

## Timeline

We believe that none of these three proposals are controversial but there’s
always a chance that we missed something important. For that reason we plan
to have the proposals implemented at the beginning of the Go 1.15 release cycle
(at or shortly after the Go 1.14 release) so that there is plenty of time to
gather experience and provide feedback. Per the
[proposal evaluation process](https://blog.golang.org/go2-here-we-come),
the final decision will be made at the end of the development cycle, at the
beginning of May, 2020.

## And one more thing...

We receive many more language change proposals
([issues labeled LanguageChange](https://github.com/golang/go/labels/LanguageChange))
than we can review thoroughly. For instance, just for error handling alone,
there are 57 issues, of which five are currently still open. Since the cost
of making a language change, no matter how small, is high and the benefits
are often unclear, we must err on the side of caution. Consequently, most
language change proposals get rejected sooner or later, sometimes with minimal
feedback. This is unsatisfactory for all parties involved. If you have spent a
lot of time and effort outlining your idea in detail, it would be nice to not
have it immediately rejected. On the flip side, because the general
[proposal process](https://github.com/golang/proposal/blob/master/README.md)
is deliberately simple, it is very easy to create language change proposals
that are only marginally explored, causing the review committee significant
amounts of work. To improve this experience for everybody we are adding a new
[questionnaire](https://github.com/golang/proposal/blob/master/go2-language-changes.md)
for language changes: filling out that template will help reviewers evaluate
proposals more efficiently because they don’t need to try to answer those
questions themselves. And hopefully it will also provide better guidance for
proposers by setting expectations right from the start. This is an experiment
that we will refine over time as needed.

Thank you for helping us improve the Go experience!
