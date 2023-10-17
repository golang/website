---
title: "Godoc: documenting Go code"
date: 2011-03-31
by:
- Andrew Gerrand
tags:
- godoc
- technical
summary: How and why to document your Go packages.
---

[_**Note, June 2022**: For updated guidelines about documenting Go code,
see “[Go Doc Comments](/doc/comment).”_]

The Go project takes documentation seriously.
Documentation is a huge part of making software accessible and maintainable.
Of course it must be well-written and accurate,
but it also must be easy to write and to maintain.
Ideally, it should be coupled to the code itself so the documentation evolves
along with the code.
The easier it is for programmers to produce good documentation,
the better for everyone.

To that end, we have developed the [godoc](/cmd/godoc/) documentation tool.
This article describes godoc's approach to documentation,
and explains how you can use our conventions and tools to write good documentation
for your own projects.

Godoc parses Go source code - including comments - and produces documentation
as HTML or plain text.
The end result is documentation tightly coupled with the code it documents.
For example, through godoc's web interface you can navigate from a function's
[documentation](/pkg/strings/#HasPrefix) to its [implementation](/src/strings/strings.go?s=11163:11200#L434) with one click.

Godoc is conceptually related to Python's [Docstring](https://www.python.org/dev/peps/pep-0257/)
and Java's [Javadoc](https://www.oracle.com/java/technologies/javase/javadoc-tool.html)
but its design is simpler.
The comments read by godoc are not language constructs (as with Docstring)
nor must they have their own machine-readable syntax (as with Javadoc).
Godoc comments are just good comments, the sort you would want to read even
if godoc didn't exist.

The convention is simple: to document a type,
variable, constant, function, or even a package,
write a regular comment directly preceding its declaration,
with no intervening blank line.
Godoc will then present that comment as text alongside the item it documents.
For example, this is the documentation for the `fmt` package's [`Fprint`](/pkg/fmt/#Fprint) function:

	// Fprint formats using the default formats for its operands and writes to w.
	// Spaces are added between operands when neither is a string.
	// It returns the number of bytes written and any write error encountered.
	func Fprint(w io.Writer, a ...interface{}) (n int, err error) {

Notice this comment is a complete sentence that begins with the name of
the element it describes.
This important convention allows us to generate documentation in a variety of formats,
from plain text to HTML to UNIX man pages,
and makes it read better when tools truncate it for brevity,
such as when they extract the first line or sentence.

Comments on package declarations should provide general package documentation.
These comments can be short, like the [`sort`](/pkg/sort/)
package's brief description:

	// Package sort provides primitives for sorting slices and user-defined
	// collections.
	package sort

They can also be detailed like the [gob package](/pkg/encoding/gob/)'s overview.
That package uses another convention for packages that need large amounts
of introductory documentation:
the package comment is placed in its own file,
[doc.go](/src/pkg/encoding/gob/doc.go),
which contains only those comments and a package clause.

When writing a package comment of any size,
keep in mind that its [first sentence](/pkg/go/doc/#Package.Synopsis)
will appear in godoc's [package list](/pkg/).

Comments that are not adjacent to a top-level declaration are omitted from godoc's output,
with one notable exception.
Top-level comments that begin with the word `"BUG(who)”` are recognized as known bugs,
and included in the "Bugs” section of the package documentation.
The "who” part should be the user name of someone who could provide more information.
For example, this is a known issue from the [bytes package](/pkg/bytes/#pkg-note-BUG):

	// BUG(r): The rule Title uses for word boundaries does not handle Unicode punctuation properly.

Sometimes a struct field, function, type, or even a whole package becomes
redundant or unnecessary, but must be kept for compatibility with existing
programs.
To signal that an identifier should not be used, add a paragraph to its doc
comment that begins with "Deprecated:" followed by some information about the
deprecation.

There are a few formatting rules that Godoc uses when converting comments to HTML:

  - Subsequent lines of text are considered part of the same paragraph;
    you must leave a blank line to separate paragraphs.

  - Pre-formatted text must be indented relative to the surrounding comment
    text (see gob's [doc.go](/src/pkg/encoding/gob/doc.go) for an example).

  - URLs will be converted to HTML links; no special markup is necessary.

Note that none of these rules requires you to do anything out of the ordinary.

In fact, the best thing about godoc's minimal approach is how easy it is to use.
As a result, a lot of Go code, including all of the standard library,
already follows the conventions.

Your own code can present good documentation just by having comments as described above.
Any Go packages installed inside `$GOROOT/src/pkg` and any `GOPATH` work
spaces will already be accessible via godoc's command-line and HTTP interfaces,
and you can specify additional paths for indexing via the `-path` flag or
just by running `"godoc ."` in the source directory.
See the [godoc documentation](/cmd/godoc/) for more details.
