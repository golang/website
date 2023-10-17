---
title: "Go Doc Comments"
layout: article
date: 2022-06-01T00:00:00Z
---

“Doc comments” are comments that appear immediately before top-level package,
const, func, type, and var declarations with no intervening newlines.
Every exported (capitalized) name should have a doc comment.

The [go/doc](/pkg/go/doc) and [go/doc/comment](/pkg/go/doc/comment) packages
provide the ability to extract documentation from Go source code,
and a variety of tools make use of this functionality.
The [`go` `doc`](/cmd/go#hdr-Show_documentation_for_package_or_symbol)
command looks up and prints the doc comment for a given package or symbol.
(A symbol is a top-level const, func, type, or var.)
The web server [pkg.go.dev](https://pkg.go.dev/) shows the documentation
for public Go packages (when their licenses permit that use).
The program serving that site is
[golang.org/x/pkgsite/cmd/pkgsite](https://pkg.go.dev/golang.org/x/pkgsite/cmd/pkgsite),
which can also be run locally to view documentation for private modules
or without an internet connection.
The language server [gopls](https://pkg.go.dev/golang.org/x/tools/gopls)
provides documentation when editing Go source files in IDEs.

The rest of this page documents how to write Go doc comments.

## Packages {#package}

Every package should have a package comment introducing the package.
It provides information relevant to the package as a whole
and generally sets expectations for the package.
Especially in large packages, it can be helpful for the package comment
to give a brief overview of the most important parts of the API,
linking to other doc comments as needed.

If the package is simple, the package comment can be brief.
For example:

	// Package path implements utility routines for manipulating slash-separated
	// paths.
	//
	// The path package should only be used for paths separated by forward
	// slashes, such as the paths in URLs. This package does not deal with
	// Windows paths with drive letters or backslashes; to manipulate
	// operating system paths, use the [path/filepath] package.
	package path

The square brackets in `[path/filepath]` create a [documentation link](#links).

As can be seen in this example, Go doc comments use complete sentences.
For a package comment, that means the [first sentence](/pkg/go/doc/#Package.Synopsis)
begins with “Package <name>”.

For multi-file packages, the package comment should only be in one source file.
If multiple files have package comments, they are concatenated to form one
large comment for the entire package.

## Commands {#cmd}

A package comment for a command is similar, but it describes the behavior
of the program rather than the Go symbols in the package.
The first sentence conventionally begins with the name of the program itself,
capitalized because it is at the start of a sentence.
For example, here is an abridged version of the package comment for [gofmt](/cmd/gofmt):

	/*
	Gofmt formats Go programs.
	It uses tabs for indentation and blanks for alignment.
	Alignment assumes that an editor is using a fixed-width font.

	Without an explicit path, it processes the standard input. Given a file,
	it operates on that file; given a directory, it operates on all .go files in
	that directory, recursively. (Files starting with a period are ignored.)
	By default, gofmt prints the reformatted sources to standard output.

	Usage:

		gofmt [flags] [path ...]

	The flags are:

		-d
			Do not print reformatted sources to standard output.
			If a file's formatting is different than gofmt's, print diffs
			to standard output.
		-w
			Do not print reformatted sources to standard output.
			If a file's formatting is different from gofmt's, overwrite it
			with gofmt's version. If an error occurred during overwriting,
			the original file is restored from an automatic backup.

	When gofmt reads from standard input, it accepts either a full Go program
	or a program fragment. A program fragment must be a syntactically
	valid declaration list, statement list, or expression. When formatting
	such a fragment, gofmt preserves leading indentation as well as leading
	and trailing spaces, so that individual sections of a Go program can be
	formatted by piping them through gofmt.
	*/
	package main

The beginning of the comment is written using
[semantic linefeeds](https://rhodesmill.org/brandon/2012/one-sentence-per-line/),
in which each new sentence or long phrase is on a line by itself,
which can make diffs easier to read as code and comments evolve.
The later paragraphs happen not to follow this convention
and have been wrapped by hand.
Whatever is best for your code base is fine.
Either way, `go` `doc` and `pkgsite` rewrap doc comment text when printing it.
For example:

	$ go doc gofmt
	Gofmt formats Go programs. It uses tabs for indentation and blanks for
	alignment. Alignment assumes that an editor is using a fixed-width font.

	Without an explicit path, it processes the standard input. Given a file, it
	operates on that file; given a directory, it operates on all .go files in that
	directory, recursively. (Files starting with a period are ignored.) By default,
	gofmt prints the reformatted sources to standard output.

	Usage:

		gofmt [flags] [path ...]

	The flags are:

		-d
			Do not print reformatted sources to standard output.
			If a file's formatting is different than gofmt's, print diffs
			to standard output.
	...

The indented lines are treated as preformatted text:
they are not rewrapped and are printed in code font
in HTML and Markdown presentations.
(The [Syntax](#syntax) section below gives the details.)

## Types {#type}

A type's doc comment should explain what each instance of that type represents or provides.
If the API is simple, the doc comment can be quite short.
For example:

	package zip

	// A Reader serves content from a ZIP archive.
	type Reader struct {
		...
	}

By default, programmers should expect that a type is safe for use only by
a single goroutine at a time.
If a type provides stronger guarantees, the doc comment should state them.
For example:

	package regexp

	// Regexp is the representation of a compiled regular expression.
	// A Regexp is safe for concurrent use by multiple goroutines,
	// except for configuration methods, such as Longest.
	type Regexp struct {
		...
	}

Go types should also aim to make the zero value have a useful meaning.
If it isn't obvious, that meaning should be documented. For example:

	package bytes

	// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
	// The zero value for Buffer is an empty buffer ready to use.
	type Buffer struct {
		...
	}

For a struct with exported fields, either the doc comment or per-field comments
should explain the meaning of each exported field.
For example, this type's doc comment explains the fields:

{{raw `
	package io

	// A LimitedReader reads from R but limits the amount of
	// data returned to just N bytes. Each call to Read
	// updates N to reflect the new amount remaining.
	// Read returns EOF when N <= 0.
	type LimitedReader struct {
		R   Reader // underlying reader
		N   int64  // max bytes remaining
	}
`}}

In contrast, this type's doc comment leaves the explanations to per-field comments:

{{raw `
	package comment

	// A Printer is a doc comment printer.
	// The fields in the struct can be filled in before calling
	// any of the printing methods
	// in order to customize the details of the printing process.
	type Printer struct {
		// HeadingLevel is the nesting level used for
		// HTML and Markdown headings.
		// If HeadingLevel is zero, it defaults to level 3,
		// meaning to use <h3> and ###.
		HeadingLevel int
		...
	}
`}}

As with packages (above) and funcs (below), doc comments for types
start with complete sentences naming the declared symbol.
An explicit subject often makes the wording clearer,
and it makes the text easier to search, whether on a web page
or a command line.
For example:

	$ go doc -all regexp | grep pairs
	pairs within the input string: result[2*n:2*n+2] identifies the indexes
	    FindReaderSubmatchIndex returns a slice holding the index pairs identifying
	    FindStringSubmatchIndex returns a slice holding the index pairs identifying
	    FindSubmatchIndex returns a slice holding the index pairs identifying the
	$

## Funcs {#func}

A function's doc comment should explain what the function returns
or, for functions called for side effects, what it does.
Named parameters and results can be referred to directly in
the comment, without any special syntax like backquotes.
(A consequence of this convention is that names like `a`,
which might be mistaken for ordinary words, are typically avoided.)
For example:

	package strconv

	// Quote returns a double-quoted Go string literal representing s.
	// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100)
	// for control characters and non-printable characters as defined by IsPrint.
	func Quote(s string) string {
		...
	}

And:

	package os

	// Exit causes the current program to exit with the given status code.
	// Conventionally, code zero indicates success, non-zero an error.
	// The program terminates immediately; deferred functions are not run.
	//
	// For portability, the status code should be in the range [0, 125].
	func Exit(code int) {
		...
	}

Doc comments typically use the phrase “reports whether”
to describe functions that return a boolean.
The phrase “or not” is unnecessary.
For example:

	package strings

	// HasPrefix reports whether the string s begins with prefix.
	func HasPrefix(s, prefix string) bool

If a doc comment needs to explain multiple results,
naming the results can make the doc comment more understandable,
even if the names are not used in the body of the function.
For example:

	package io

	// Copy copies from src to dst until either EOF is reached
	// on src or an error occurs. It returns the total number of bytes
	// written and the first error encountered while copying, if any.
	//
	// A successful Copy returns err == nil, not err == EOF.
	// Because Copy is defined to read from src until EOF, it does
	// not treat an EOF from Read as an error to be reported.
	func Copy(dst Writer, src Reader) (n int64, err error) {
		...
	}

Conversely, when the results don't need to be named in the doc comment,
they are usually omitted in the code as well, like in the `Quote` example above,
to avoid cluttering the presentation.

These rules all apply both to plain functions and to methods.
For methods, using the same receiver name avoids needless
variation when listing all the methods of a type:

	$ go doc bytes.Buffer
	package bytes // import "bytes"

	type Buffer struct {
		// Has unexported fields.
	}
	    A Buffer is a variable-sized buffer of bytes with Read and Write methods.
	    The zero value for Buffer is an empty buffer ready to use.

	func NewBuffer(buf []byte) *Buffer
	func NewBufferString(s string) *Buffer
	func (b *Buffer) Bytes() []byte
	func (b *Buffer) Cap() int
	func (b *Buffer) Grow(n int)
	func (b *Buffer) Len() int
	func (b *Buffer) Next(n int) []byte
	func (b *Buffer) Read(p []byte) (n int, err error)
	func (b *Buffer) ReadByte() (byte, error)
	...

This example also shows that top-level functions returning a type `T` or pointer `*T`,
perhaps with an additional error result,
are shown alongside the type `T` and its methods,
under the assumption that they are `T`'s constructors.

By default, programmers can assume that a top-level function
is safe to call from multiple goroutines;
this fact need not be stated explicitly.

On the other hand, as noted in the previous section,
using an instance of a type in any way,
including calling a method, is typically assumed
to be restricted to a single goroutine at a time.
If the methods that are safe for concurrent use
are not documented in the type's doc comment,
they should be documented in per-method comments.
For example:

	package sql

	// Close returns the connection to the connection pool.
	// All operations after a Close will return with ErrConnDone.
	// Close is safe to call concurrently with other operations and will
	// block until all other operations finish. It may be useful to first
	// cancel any used context and then call Close directly after.
	func (c *Conn) Close() error {
		...
	}

Note that function and method doc comments focus on
what the operation returns or does,
detailing what the caller needs to know.
Special cases can be particularly important to document.
For example:

{{raw `
	package math

	// Sqrt returns the square root of x.
	//
	// Special cases are:
	//
	//	Sqrt(+Inf) = +Inf
	//	Sqrt(±0) = ±0
	//	Sqrt(x < 0) = NaN
	//	Sqrt(NaN) = NaN
	func Sqrt(x float64) float64 {
		...
	}
`}}

Doc comments should not explain internal details
such as the algorithm used in the current implementation.
Those are best left to comments inside the function body.
It may be appropriate to give asymptotic time or space bounds
when that detail is particularly important to callers.
For example:

	package sort

	// Sort sorts data in ascending order as determined by the Less method.
	// It makes one call to data.Len to determine n and O(n*log(n)) calls to
	// data.Less and data.Swap. The sort is not guaranteed to be stable.
	func Sort(data Interface) {
		...
	}

Because this doc comment makes no mention of which sorting algorithm is used,
it is easier to change the implementation to use a different algorithm in the future.

## Consts {#const}

Go's declaration syntax allows grouping of declarations,
in which case a single doc comment can introduce a group of related constants,
with individual constants only documented by short end-of-line comments.
For example:

	package scanner // import "text/scanner"

	// The result of Scan is one of these tokens or a Unicode character.
	const (
		EOF = -(iota + 1)
		Ident
		Int
		Float
		Char
		...
	)

Sometimes the group needs no doc comment at all. For example:

	package unicode // import "unicode"

	const (
		MaxRune         = '\U0010FFFF' // maximum valid Unicode code point.
		ReplacementChar = '\uFFFD'     // represents invalid code points.
		MaxASCII        = '\u007F'     // maximum ASCII value.
		MaxLatin1       = '\u00FF'     // maximum Latin-1 value.
	)

On the other hand, ungrouped constants typically warrant a full
doc comment starting with a complete sentence. For example:

	package unicode

	// Version is the Unicode edition from which the tables are derived.
	const Version = "13.0.0"

Typed constants are displayed next to the declaration of their type
and as a result often omit a const group doc comment in favor of
the type's doc comment.
For example:

	package syntax

	// An Op is a single regular expression operator.
	type Op uint8

	const (
		OpNoMatch        Op = 1 + iota // matches no strings
		OpEmptyMatch                   // matches empty string
		OpLiteral                      // matches Runes sequence
		OpCharClass                    // matches Runes interpreted as range pair list
		OpAnyCharNotNL                 // matches any character except newline
		...
	)

(See [pkg.go.dev/regexp/syntax#Op](https://pkg.go.dev/regexp/syntax#Op) for the HTML presentation.)

## Vars {#var}

The conventions for variables are the same as those for constants.
For example, here is a set of grouped variables:

	package fs

	// Generic file system errors.
	// Errors returned by file systems can be tested against these errors
	// using errors.Is.
	var (
		ErrInvalid    = errInvalid()    // "invalid argument"
		ErrPermission = errPermission() // "permission denied"
		ErrExist      = errExist()      // "file already exists"
		ErrNotExist   = errNotExist()   // "file does not exist"
		ErrClosed     = errClosed()     // "file already closed"
	)

And a single variable:

	package unicode

	// Scripts is the set of Unicode script tables.
	var Scripts = map[string]*RangeTable{
		"Adlam":                  Adlam,
		"Ahom":                   Ahom,
		"Anatolian_Hieroglyphs":  Anatolian_Hieroglyphs,
		"Arabic":                 Arabic,
		"Armenian":               Armenian,
		...
	}

## Syntax {#syntax}

Go doc comments are written in a simple syntax that supports
paragraphs, headings, links, lists, and preformatted code blocks.
To keep comments lightweight and readable in source files,
there is no support for complex features like font changes or raw HTML.
Markdown aficionados can view the syntax as a simplified subset of Markdown.

The standard formatter [gofmt](/cmd/gofmt) reformats doc comments
to use a canonical formatting for each of these features.
Gofmt aims for readability and user control over how comments
are written in source code but will adjust presentation to make
the semantic meaning of a particular comment clearer,
analogous to reformatting `1+2 * 3` to `1 + 2*3` in ordinary source code.

Directive comments such as `//go:generate` are not
considered part of a doc comment and are omitted from
rendered documentation.
Gofmt moves directive comments to the end of the doc comment,
preceded by a blank line.
For example:

	package regexp

	// An Op is a single regular expression operator.
	//
	//go:generate stringer -type Op -trimprefix Op
	type Op uint8

A directive comment is a line matching the regular expression
`//(line |extern |export |[a-z0-9]+:[a-z0-9])`.
Tools that define their own directives should use the form
`//toolname:directive`.

Gofmt removes leading and trailing blank lines in doc comments.
If all lines in a doc comment begin with the same sequence of
spaces and tabs, gofmt removes that prefix.

### Paragraphs {#paragraphs}

A paragraph is a span of unindented non-blank lines.
We've already seen many examples of paragraphs.

A pair of consecutive backticks (\` U+0060)
is interpreted as a Unicode left quote (“ U+201C),
and a pair of consecutive single quotes (\' U+0027)
is interpreted as a Unicode right quote (” U+201D).

Gofmt preserves line breaks in paragraph text: it does not rewrap the text.
This allows the use of [semantic linefeeds](https://rhodesmill.org/brandon/2012/one-sentence-per-line/),
as seen earlier.
Gofmt replaces duplicated blank lines between paragraphs
with a single blank line.
Gofmt also reformats consecutive backticks or single quotes
to their Unicode interpretations.

### Headings {#headings}

A heading is a line beginning with a number sign (U+0023) and then a space and the heading text.
To be recognized as a heading, the line must be unindented and set off from adjacent paragraph text
by blank lines.

For example:

	// Package strconv implements conversions to and from string representations
	// of basic data types.
	//
	// # Numeric Conversions
	//
	// The most common numeric conversions are [Atoi] (string to int) and [Itoa] (int to string).
	...
	package strconv

On the other hand:

	// #This is not a heading, because there is no space.
	//
	// # This is not a heading,
	// # because it is multiple lines.
	//
	// # This is not a heading,
	// because it is also multiple lines.
	//
	// The next paragraph is not a heading, because there is no additional text:
	//
	// #
	//
	// In the middle of a span of non-blank lines,
	// # this is not a heading either.
	//
	//     # This is not a heading, because it is indented.

The # syntax was added in Go 1.19.
Before Go 1.19, headings were identified implicitly by single-line paragraphs
satisfying certain conditions, most notably the lack of any terminating punctuation.

Gofmt reformats [lines treated as implicit headings](https://github.com/golang/proposal/blob/master/design/51082-godocfmt.md#headings)
by earlier versions of Go to use # headings instead.
If the reformatting is not appropriate—that is, if the line was not meant to be a heading—the easiest
way to make it a paragraph is to introduce terminating punctuation
such as a period or colon, or to break it into two lines.

### Links {#links}

A span of unindented non-blank lines defines link targets
when every line is of the form “[Text]: URL”.
In other text in the same doc comment,
“[Text]” represents a link to URL using the given text—in HTML,
\<a href="URL">Text\</a>.
For example:

	// Package json implements encoding and decoding of JSON as defined in
	// [RFC 7159]. The mapping between JSON and Go values is described
	// in the documentation for the Marshal and Unmarshal functions.
	//
	// For an introduction to this package, see the article
	// “[JSON and Go].”
	//
	// [RFC 7159]: https://tools.ietf.org/html/rfc7159
	// [JSON and Go]: https://golang.org/doc/articles/json_and_go.html
	package json

By keeping URLs in a separate section,
this format only minimally interrupts the flow of the actual text.
It also roughly matches the Markdown
[shortcut reference link format](https://spec.commonmark.org/0.30/#shortcut-reference-link),
without the optional title text.

If there is no corresponding URL declaration,
then (except for doc links, described in the next section)
“[Text]” is not a hyperlink, and the square brackets are preserved
when displayed.
Each doc comment is considered independently:
link target definitions in one comment do not affect other comments.

Although link target definition blocks may be interleaved with
ordinary paragraphs, gofmt moves all link target definitions to
the end of the doc comment,
in up to two blocks: first a block containing all the link targets
that are referenced in the comment, and then a block
containing all the targets _not_ referenced in the comment.
The separate block makes unused targets easy
to notice and fix (in case the links or the definitions have typos)
or to delete (in case the definitions are no longer needed).

Plain text that is recognized as a URL is automatically linked in HTML renderings.

### Doc links {#doclinks}

Doc links are links of the form “[Name1]” or “[Name1.Name2]” to refer
to exported identifiers in the current package, or “[pkg]”,
“[pkg.Name1]”, or “[pkg.Name1.Name2]” to refer to identifiers in other
packages.

For example:

	package bytes

	// ReadFrom reads data from r until EOF and appends it to the buffer, growing
	// the buffer as needed. The return value n is the number of bytes read. Any
	// error except [io.EOF] encountered during the read is also returned. If the
	// buffer becomes too large, ReadFrom will panic with [ErrTooLarge].
	func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
		...
	}

The bracketed text for a symbol link
can include an optional leading star, making it easy to refer to
pointer types, such as \[\*bytes.Buffer\].

When referring to other packages, “pkg” can be either a full import path
or the assumed package name of an existing import. The assumed package
name is either the identifier in a renamed import or else
[the name assumed by
goimports](https://pkg.go.dev/golang.org/x/tools/internal/imports#ImportPathToAssumedName).
(Goimports inserts renamings when that assumption is not correct, so
this rule should work for essentially all Go code.)
For example, if the current package imports encoding/json,
then “[json.Decoder]” can be written in place of “[encoding/json.Decoder]”
to link to the docs for encoding/json's Decoder.
If different source files in a package import different packages using the same name,
then the shorthand is ambiguous and cannot be used.

A “pkg” is only
assumed to be a full import path if it starts with a domain name (a
path element with a dot) or is one of the packages from the standard
library (“[os]”, “[encoding/json]”, and so on).
For example, `[os.File]` and `[example.com/sys.File]` are documentation links
(the latter will be a broken link),
but `[os/sys.File]` is not, because there is no os/sys package in the standard library.

To avoid problems with
maps, generics, and array types, doc links must be both preceded and
followed by punctuation, spaces, tabs, or the start or end of a line.
For example, the text “map[ast.Expr]TypeAndValue” does not contain
a doc link.

### Lists {#lists}

A list is a span of indented or blank lines
(which would otherwise be a code block,
as described in the next section)
in which the first indented line begins with
a bullet list marker or a numbered list marker.

A bullet list marker is a star, plus, dash, or Unicode bullet
(*, +, -, •; U+002A, U+002B, U+002D, U+2022)
followed by a space or tab and then text.
In a bullet list, each line beginning with a bullet list
marker starts a new list item.

For example:

	package url

	// PublicSuffixList provides the public suffix of a domain. For example:
	//   - the public suffix of "example.com" is "com",
	//   - the public suffix of "foo1.foo2.foo3.co.uk" is "co.uk", and
	//   - the public suffix of "bar.pvt.k12.ma.us" is "pvt.k12.ma.us".
	//
	// Implementations of PublicSuffixList must be safe for concurrent use by
	// multiple goroutines.
	//
	// An implementation that always returns "" is valid and may be useful for
	// testing but it is not secure: it means that the HTTP server for foo.com can
	// set a cookie for bar.com.
	//
	// A public suffix list implementation is in the package
	// golang.org/x/net/publicsuffix.
	type PublicSuffixList interface {
		...
	}

A numbered list marker is a decimal number of any length
followed by a period or right parenthesis, then a space or tab, and then text.
In a numbered list, each line beginning with a number list marker starts a new list item.
Item numbers are left as is, never renumbered.

For example:

	package path

	// Clean returns the shortest path name equivalent to path
	// by purely lexical processing. It applies the following rules
	// iteratively until no further processing can be done:
	//
	//  1. Replace multiple slashes with a single slash.
	//  2. Eliminate each . path name element (the current directory).
	//  3. Eliminate each inner .. path name element (the parent directory)
	//     along with the non-.. element that precedes it.
	//  4. Eliminate .. elements that begin a rooted path:
	//     that is, replace "/.." by "/" at the beginning of a path.
	//
	// The returned path ends in a slash only if it is the root "/".
	//
	// If the result of this process is an empty string, Clean
	// returns the string ".".
	//
	// See also Rob Pike, “[Lexical File Names in Plan 9].”
	//
	// [Lexical File Names in Plan 9]: https://9p.io/sys/doc/lexnames.html
	func Clean(path string) string {
		...
	}

List items only contain paragraphs, not code blocks or nested lists.
This avoids any space-counting subtlety as well as questions about
how many spaces a tab counts for in inconsistent indentation.

Gofmt reformats bullet lists to use a dash as the bullet marker,
two spaces of indentation before the dash,
and four spaces of indentation for continuation lines.

Gofmt reformats numbered lists to use a single space before the number,
a period after the number, and again
four spaces of indentation for continuation lines.

Gofmt preserves but does not require a blank line between a list and the preceding paragraph.
It inserts a blank line between a list and the following paragraph or heading.

### Code blocks {#code}

A code block is a span of indented or blank lines
not starting with a bullet list marker or numbered list marker.
It is rendered as preformatted text (a \<pre> block in HTML).

Code blocks often contain Go code. For example:

{{raw `
	package sort

	// Search uses binary search...
	//
	// As a more whimsical example, this program guesses your number:
	//
	//	func GuessingGame() {
	//		var s string
	//		fmt.Printf("Pick an integer from 0 to 100.\n")
	//		answer := sort.Search(100, func(i int) bool {
	//			fmt.Printf("Is your number <= %d? ", i)
	//			fmt.Scanf("%s", &s)
	//			return s != "" && s[0] == 'y'
	//		})
	//		fmt.Printf("Your number is %d.\n", answer)
	//	}
	func Search(n int, f func(int) bool) int {
		...
	}
`}}

Of course, code blocks also often contain preformatted text besides code. For example:

{{raw `
	package path

	// Match reports whether name matches the shell pattern.
	// The pattern syntax is:
	//
	//	pattern:
	//		{ term }
	//	term:
	//		'*'         matches any sequence of non-/ characters
	//		'?'         matches any single non-/ character
	//		'[' [ '^' ] { character-range } ']'
	//		            character class (must be non-empty)
	//		c           matches character c (c != '*', '?', '\\', '[')
	//		'\\' c      matches character c
	//
	//	character-range:
	//		c           matches character c (c != '\\', '-', ']')
	//		'\\' c      matches character c
	//		lo '-' hi   matches character c for lo <= c <= hi
	//
	// Match requires pattern to match all of name, not just a substring.
	// The only possible returned error is [ErrBadPattern], when pattern
	// is malformed.
	func Match(pattern, name string) (matched bool, err error) {
		...
	}
`}}

Gofmt indents all lines in a code block by a single tab,
replacing any other indentation the non-blank lines have in common.
Gofmt also inserts a blank line before and after each code block,
distinguishing the code block clearly from the surrounding paragraph text.

## Common mistakes and pitfalls {#mistakes}

The rule that any span of indented or blank lines
in a doc comment is rendered as a code block
dates to the earliest days of Go.
Unfortunately, the lack of support for doc comments in gofmt
has led to many existing comments that use indentation
without meaning to create a code block.

For example, this unindented list has always been interpreted
by godoc as a three-line paragraph followed by a one-line code block:

	package http

	// cancelTimerBody is an io.ReadCloser that wraps rc with two features:
	// 1) On Read error or close, the stop func is called.
	// 2) On Read failure, if reqDidTimeout is true, the error is wrapped and
	//    marked as net.Error that hit its timeout.
	type cancelTimerBody struct {
		...
	}

This always rendered in `go` `doc` as:

	cancelTimerBody is an io.ReadCloser that wraps rc with two features:
	1) On Read error or close, the stop func is called. 2) On Read failure,
	if reqDidTimeout is true, the error is wrapped and

	    marked as net.Error that hit its timeout.

Similarly, the command in this comment is a one-line paragraph
followed by a one-line code block:

	package smtp

	// localhostCert is a PEM-encoded TLS cert generated from src/crypto/tls:
	//
	// go run generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com \
	//     --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h
	var localhostCert = []byte(`...`)

This rendered in `go` `doc` as:

	localhostCert is a PEM-encoded TLS cert generated from src/crypto/tls:

	go run generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com \

	    --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h

And this comment is a two-line paragraph (the second line is “{”),
followed by a six-line indented code block and a one-line paragraph (“}”).

	// On the wire, the JSON will look something like this:
	// {
	//	"kind":"MyAPIObject",
	//	"apiVersion":"v1",
	//	"myPlugin": {
	//		"kind":"PluginA",
	//		"aOption":"foo",
	//	},
	// }

And this rendered in `go` `doc` as:

	On the wire, the JSON will look something like this: {

	    "kind":"MyAPIObject",
	    "apiVersion":"v1",
	    "myPlugin": {
	    	"kind":"PluginA",
	    	"aOption":"foo",
	    },

	}

Another common mistake was an unindented Go function definition
or block statement, similarly bracketed by “{” and “}”.

The introduction of doc comment reformatting in Go 1.19's gofmt makes mistakes
like these more visible by adding blank lines around the code blocks.

Analysis in 2022 found that only 3% of doc comments in public Go modules
were reformatted at all by the draft Go 1.19 gofmt.
Limiting ourselves to those comments, about 87% of gofmt's reformattings
preserved the structure that a person would infer from reading the comment;
about 6% were tripped up by these kinds of unindented lists,
unindented multiline shell commands, and unindented brace-delimited code blocks.

Based on this analysis, the Go 1.19 gofmt applies a few heuristics to merge
unindented lines into an adjacent indented list or code block.
With those adjustments, the Go 1.19 gofmt reformats the above examples to:

	// cancelTimerBody is an io.ReadCloser that wraps rc with two features:
	//  1. On Read error or close, the stop func is called.
	//  2. On Read failure, if reqDidTimeout is true, the error is wrapped and
	//     marked as net.Error that hit its timeout.

	// localhostCert is a PEM-encoded TLS cert generated from src/crypto/tls:
	//
	//	go run generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com \
	//	    --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h

	// On the wire, the JSON will look something like this:
	//
	//	{
	//		"kind":"MyAPIObject",
	//		"apiVersion":"v1",
	//		"myPlugin": {
	//			"kind":"PluginA",
	//			"aOption":"foo",
	//		},
	//	}

This reformatting makes the meaning clearer as well as making the doc comments
render correctly in earlier versions of Go.
If the heuristic ever makes a bad decision, it can be overridden by inserting
a blank line to clearly separate the paragraph text from non-paragraph text.

Even with these heuristics, other existing comments will need manual
adjustment to correct their rendering.
The most common mistake is indenting a wrapped unindented line of text.
For example:

	// TODO Revisit this design. It may make sense to walk those nodes
	//      only once.

	// According to the document:
	// "The alignment factor (in bytes) that is used to align the raw data of sections in
	//  the image file. The value should be a power of 2 between 512 and 64 K, inclusive."

In both of these, the last line is indented, making it a code block.
The fix is to unindent the lines.

Another common mistake is not indenting a wrapped indented line of a list or code block.
For example:

	// Uses of this error model include:
	//
	//   - Partial errors. If a service needs to return partial errors to the
	// client,
	//     it may embed the `Status` in the normal response to indicate the
	// partial
	//     errors.
	//
	//   - Workflow errors. A typical workflow has multiple steps. Each step
	// may
	//     have a `Status` message for error reporting.

The fix is to indent the wrapped lines.

Go doc comments do not support nested lists, so gofmt reformats

	// Here is a list:
	//
	//  - Item 1.
	//    * Subitem 1.
	//    * Subitem 2.
	//  - Item 2.
	//  - Item 3.

to

	// Here is a list:
	//
	//  - Item 1.
	//  - Subitem 1.
	//  - Subitem 2.
	//  - Item 2.
	//  - Item 3.

Rewriting the text to avoid nested lists usually
improves the documentation and is the best solution.
Another potential workaround is to mix list markers,
since bullet markers do not introduce list items in a numbered list,
nor vice versa.
For example:

	// Here is a list:
	//
	//  1. Item 1.
	//
	//     - Subitem 1.
	//
	//     - Subitem 2.
	//
	//  2. Item 2.
	//
	//  3. Item 3.
