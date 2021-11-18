---
title: go fmt your code
date: 2013-01-23
by:
- Andrew Gerrand
tags:
- gofix
- gofmt
- technical
summary: How and why to format your Go code using gofmt.
---

## Introduction

[Gofmt](/cmd/gofmt/) is a tool that automatically formats Go source code.

Gofmt'd code is:

  - easier to **write**: never worry about minor formatting concerns while hacking away,

  - easier to **read**: when all code looks the same you need not mentally convert
    others' formatting style into something you can understand.

  - easier to **maintain**: mechanical changes to the source don't cause unrelated
    changes to the file's formatting;
    diffs show only the real changes.

  - **uncontroversial**: never have a debate about spacing or brace position ever again!

## Format your code

We recently conducted a survey of Go packages in the wild and found that
about 70% of them are formatted according to gofmt's rules.
This was more than expected - and thanks to everyone who uses gofmt - but
it would be great to close the gap.

To format your code, you can use the gofmt tool directly:

	gofmt -w yourcode.go

Or you can use the "[go fmt](/cmd/go/#hdr-Gofmt__reformat__package_sources)" command:

	go fmt path/to/your/package

To help keep your code in the canonical style,
the Go repository contains hooks for editors and version control systems
that make it easy to run gofmt on your code.

For Vim users, the [Vim plugin for Go](https://github.com/fatih/vim-go)
includes the :Fmt command that runs gofmt on the current buffer.

For emacs users, [go-mode.el](https://github.com/dominikh/go-mode.el)
provides a gofmt-before-save hook that can be installed by adding this line
to your .emacs file:

	(add-hook 'before-save-hook #'gofmt-before-save)

For Eclipse or Sublime Text users, the [GoClipse](https://github.com/GoClipse/goclipse)
and [GoSublime](https://github.com/DisposaBoy/GoSublime) projects add
a gofmt facility to those editors.

And for Git aficionados, the [misc/git/pre-commit script](https://github.com/golang/go/blob/release-branch.go1.1/misc/git/pre-commit)
is a pre-commit hook that prevents incorrectly-formatted Go code from being committed.
If you use Mercurial, the [hgstyle plugin](https://bitbucket.org/fhs/hgstyle/overview)
provides a gofmt pre-commit hook.

## Mechanical source transformation

One of the greatest virtues of machine-formatted code is that it can be
transformed mechanically without generating unrelated formatting noise in the diffs.
Mechanical transformation is invaluable when working with large code bases,
as it is both more comprehensive and less error prone than making wide-sweeping changes by hand.
Indeed, when working at scale (like we do at Google) it often isn't practical
to make these kinds of changes manually.

The easiest way to mechanically manipulate Go code is with gofmt's -r flag.
The flag specifies a rewrite rule of the form

	pattern -> replacement

where both pattern and replacement are valid Go expressions.
In the pattern, single-character lowercase identifiers serve as wildcards
matching arbitrary sub-expressions,
and those expressions are substituted for the same identifiers in the replacement.

For example, this[ recent change](/cl/7038051) to the
Go core rewrote some uses of [bytes.Compare](/pkg/bytes/#Compare)
to use the more efficient [bytes.Equal](/pkg/bytes/#Equal).
The contributor made the change using just two gofmt invocations:

	gofmt -r 'bytes.Compare(a, b) == 0 -> bytes.Equal(a, b)'
	gofmt -r 'bytes.Compare(a, b) != 0 -> !bytes.Equal(a, b)'

Gofmt also enables [gofix](/cmd/fix/),
which can make arbitrarily complex source transformations.
Gofix was an invaluable tool during the early days when we regularly made
breaking changes to the language and libraries.
For example, before Go 1 the built-in error interface didn't exist and the
convention was to use the os.Error type.
When we [introduced error](/doc/go1.html#errors),
we provided a gofix module that rewrote all references to os.Error and its
associated helper functions to use error and the new [errors package](/pkg/errors/).
It would have been daunting to attempt by hand,
but with the code in a standard format it was relatively easy to prepare,
execute, and review this change which touched almost all Go code in existence.

For more about gofix, see [this article](https://blog.golang.org/introducing-gofix).
