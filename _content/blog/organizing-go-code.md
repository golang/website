---
title: Organizing Go code
date: 2012-08-16
by:
- Andrew Gerrand
tags:
- godoc
- gopath
- interface
- libraries
- tools
- technical
summary: How to name and package the parts of your Go program to best serve your users.
---

## Introduction

Go code is organized differently to that of other languages.
This post discusses how to name and package the elements of your Go program
to best serve its users.

## Choose good names

The names you choose affect how you think about your code,
so take care when naming your package and its exported identifiers.

A package's name provides context for its contents.
For instance, the [bytes package](/pkg/bytes/) from
the standard library exports the `Buffer` type.
On its own, the name `Buffer` isn't very descriptive,
but when combined with its package name its meaning becomes clear: `bytes.Buffer`.
If the package had a less descriptive name,
like `util`, the buffer would likely acquire the longer and clumsier name `util.BytesBuffer`.

Don't be shy about renaming things as you work.
As you spend time with your program you will better understand how its pieces fit together and,
therefore, what their names should be.
There's no need to lock yourself into early decisions.
(The [gofmt command](/cmd/gofmt/) has a `-r` flag that
provides a syntax-aware search and replace,
making large-scale refactoring easier.)

A good name is the most important part of a software interface:
the name is the first thing every client of the code will see.
A well-chosen name is therefore the starting point for good documentation.
Many of the following practices result organically from good naming.

## Choose a good import path (make your package "go get"-able)

An import path is the string with which users import a package.
It specifies the directory (relative to `$GOROOT/src/pkg` or `$GOPATH/src`)
in which the package's source code resides.

Import paths should be globally unique, so use the path of your source repository as its base.
For instance, the `websocket` package from the `go.net` sub-repository has
an import path of `"golang.org/x/net/websocket"`.
The Go project owns the path `"github.com/golang"`,
so that path cannot be used by another author for a different package.
Because the repository URL and import path are one and the same,
the `go get` command can fetch and install the package automatically.

If you don't use a hosted source repository,
choose some unique prefix such as a domain,
company, or project name.
As an example, the import path of all Google's internal Go code starts with
the string `"google"`.

The last element of the import path is typically the same as the package name.
For instance, the import path `"net/http"` contains package `http`.
This is not a requirement - you can make them different if you like - but
you should follow the convention for predictability's sake:
a user might be surprised that import `"foo/bar"` introduces the identifier
`quux` into the package name space.

Sometimes people set `GOPATH` to the root of their source repository and
put their packages in directories relative to the repository root,
such as `"src/my/package"`.
On one hand, this keeps the import paths short (`"my/package"` instead of
`"github.com/me/project/my/package"`),
but on the other it breaks `go get` and forces users to re-set their `GOPATH`
to use the package. Don't do this.

## Minimize the exported interface

Your code is likely composed of many small pieces of useful code,
and so it is tempting to expose much of that functionality in your package's
exported interface. Resist that urge!

The larger the interface you provide, the more you must support.
Users will quickly come to depend on every type,
function, variable, and constant you export,
creating an implicit contract that you must honor in perpetuity or risk
breaking your users' programs.
In preparing Go 1 we carefully reviewed the standard library's exported
interfaces and removed the parts we weren't ready to commit to.
You should take similar care when distributing your own libraries.

If in doubt, leave it out!

## What to put into a package

It is easy to just throw everything into a "grab bag" package,
but this dilutes the meaning of the package name (as it must encompass a
lot of functionality) and forces the users of small parts of the package
to compile and link a lot of unrelated code.

On the other hand, it is also easy to go overboard in splitting your code
into small packages,
in which case you will likely become bogged down in interface design,
rather than just getting the job done.

Look to the Go standard libraries as a guide.
Some of its packages are large and some are small.
For instance, the [http package](/pkg/net/http/) comprises
17 go source files (excluding tests) and exports 109 identifiers,
and the [hash package](/pkg/hash/) consists of one file
that exports just three declarations.
There is no hard and fast rule; both approaches are appropriate given their context.

With that said, package main is often larger than other packages.
Complex commands contain a lot of code that is of little use outside the
context of the executable,
and often it's simpler to just keep it all in the one place.
For instance, the go tool is more than 12000 lines spread across [34 files](/src/cmd/go/).

## Document your code

Good documentation is an essential quality of usable and maintainable code.
Read the [Godoc: documenting Go code](/doc/articles/godoc_documenting_go_code.html)
article to learn how to write good doc comments.
