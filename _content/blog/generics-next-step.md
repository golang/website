---
title: The Next Step for Generics
date: 2020-06-16
by:
- Ian Lance Taylor
- Robert Griesemer
tags:
- go2
- proposals
- generics
summary: An updated generics design draft, and a translation tool for experimentation
---

## Introduction

It’s been almost a year since we [last wrote about the possibility of
adding generics to Go](https://blog.golang.org/why-generics).
It’s time for an update.

## Updated design

We’ve been continuing to refine the [generics design
draft](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-contracts.md).
We’ve written a type checker for it: a program that can parse Go code
that uses generics as described in the design draft and report any
type errors.
We’ve written example code.
And we’ve collected feedback from many, many people&mdash;thanks for
providing it!

Based on what we’ve learned, we’re releasing an [updated design
draft](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md).
The biggest change is that we are dropping the idea of contracts.
The difference between contracts and interface types was confusing, so
we’re eliminating that difference.
Type parameters are now constrained by interface types.
Interface types are now permitted to include type lists, though only
when used as constraints; in the previous design draft type lists were
a feature of contracts.
More complex cases will use a parameterized interface type.

We hope that people will find this design draft simpler and easier to
understand.

## Experimentation tool

To help decide how to further refine the design draft, we are
releasing a translation tool.
This is a tool that permits people to type check and run code written
using the version of generics described in the design draft.
It works by translating generic code into ordinary Go code.
This translation process imposes some limitations, but we hope that it
will be good enough for people to get a feel for what generic Go code
might look like.
The real implementation of generics, if they are accepted into the
language, will work differently.
(We have only just begun to sketch out what a direct compiler
implementation would look like.)

The tool is available on a variant of the Go playground at
[https://go2goplay.golang.org](https://go2goplay.golang.org).
This playground works just like the usual Go playground, but it
supports generic code.

You can also build and use the tool yourself.
It is available in a branch of the master Go repo.
Follow the [instructions on installing Go from
source](/doc/install/source).
Where those instructions direct you to check out the latest release
tag, instead run `git checkout dev.go2go`.
Then build the Go toolchain as directed.

The translation tool is documented in
[README.go2go](https://go.googlesource.com/go/+/refs/heads/dev.go2go/README.go2go.md).

## Next steps

We hope that the tool will give the Go community a chance to
experiment with generics.
There are two main things that we hope to learn.

First, does generic code make sense?
Does it feel like Go?
What surprises do people encounter?
Are the error messages useful?

Second, we know that many people have said that Go needs generics, but
we don’t necessarily know exactly what that means.
Does this draft design address the problem in a useful way?
If there is a problem that makes you think “I could solve this if Go
had generics,” can you solve the problem when using this tool?

We will use the feedback we gather from the Go community to decide how
to move forward.
If the draft design is well received and doesn’t need significant
changes, the next step would be a [formal language change
proposal](/s/proposal).
To set expectations, if everybody is completely happy with the design
draft and it does not require any further adjustments, the earliest
that generics could be added to Go would be the Go 1.17 release,
scheduled for August 2021.
In reality, of course, there may be unforeseen problems, so this is an
optimistic timeline; we can’t make any definite prediction.

## Feedback

The best way to provide feedback for the language changes will be on
the mailing list `golang-nuts@googlegroups.com`.
Mailing lists are imperfect, but they seem like our best option for
initial discussion.
When writing about the design draft, please put `[generics]` at the
start of the Subject line and to start different threads for different
specific topics.

If you find bugs in the generics type checker or the translation tool,
they should be filed in the standard Go issue tracker at
[https://golang.org/issue](/issue).
Please start the issue title with `cmd/go2go:`.
Note that the issue tracker is not the best place to discuss changes
to the language, because it does not provide threading and it is not
well suited to lengthy conversations.

We look forward to your feedback.

## Acknowledgements

We’re not finished, but we’ve come a long way.
We would not be here without a lot of help.

We’d like to thank Philip Wadler and his collaborators for thinking
formally about generics in Go and helping us clarify the theoretical
aspects of the design.
Their paper [Featherweight Go](https://arxiv.org/abs/2005.11710)
analyzes generics in a restricted version of Go, and they have
developed a prototype [on GitHub](https://github.com/rhu1/fgg).

We would also like to thank [the
people](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md#acknowledgements)
who provided detailed feedback on an earlier version of the design
draft.

And last but definitely not least, we’d like to thank many people on
the Go team, many contributors to the Go issue tracker, and everybody
else who shared ideas and feedback on earlier design drafts.
We read all of it, and we’re grateful.  We wouldn’t be here without
you.
