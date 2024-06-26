---
title: Strings, bytes, runes and characters in Go
date: 2013-10-23
by:
- Rob Pike
tags:
- strings
- bytes
- runes
- characters
summary: How strings work in Go, and how to use them.
---

## Introduction

The [previous blog post](/blog/slices) explained how slices
work in Go, using a number of examples to illustrate the mechanism behind
their implementation.
Building on that background, this post discusses strings in Go.
At first, strings might seem too simple a topic for a blog post, but to use
them well requires understanding not only how they work,
but also the difference between a byte, a character, and a rune,
the difference between Unicode and UTF-8,
the difference between a string and a string literal,
and other even more subtle distinctions.

One way to approach this topic is to think of it as an answer to the frequently
asked question, "When I index a Go string at position _n_, why don't I get the
_nth_ character?"
As you'll see, this question leads us to many details about how text works
in the modern world.

An excellent introduction to some of these issues, independent of Go,
is Joel Spolsky's famous blog post,
[The Absolute Minimum Every Software Developer Absolutely, Positively Must Know About Unicode and Character Sets (No Excuses!)](http://www.joelonsoftware.com/articles/Unicode.html).
Many of the points he raises will be echoed here.

## What is a string?

Let's start with some basics.

In Go, a string is in effect a read-only slice of bytes.
If you're at all uncertain about what a slice of bytes is or how it works,
please read the [previous blog post](/blog/slices);
we'll assume here that you have.

It's important to state right up front that a string holds _arbitrary_ bytes.
It is not required to hold Unicode text, UTF-8 text, or any other predefined format.
As far as the content of a string is concerned, it is exactly equivalent to a
slice of bytes.

Here is a string literal (more about those soon) that uses the
`\xNN` notation to define a string constant holding some peculiar byte values.
(Of course, bytes range from hexadecimal values 00 through FF, inclusive.)

{{code "strings/basic.go" `/const sample/`}}

## Printing strings

Because some of the bytes in our sample string are not valid ASCII, not even
valid UTF-8, printing the string directly will produce ugly output.
The simple print statement

{{code "strings/basic.go" `/println/` `/println/`}}

produces this mess (whose exact appearance varies with the environment):

	��=� ⌘

To find out what that string really holds, we need to take it apart and examine the pieces.
There are several ways to do this.
The most obvious is to loop over its contents and pull out the bytes
individually, as in this `for` loop:

{{code "strings/basic.go" `/byte loop/` `/byte loop/`}}

As implied up front, indexing a string accesses individual bytes, not
characters. We'll return to that topic in detail below. For now, let's
stick with just the bytes.
This is the output from the byte-by-byte loop:

	bd b2 3d bc 20 e2 8c 98

Notice how the individual bytes match the
hexadecimal escapes that defined the string.

A shorter way to generate presentable output for a messy string
is to use the `%x` (hexadecimal) format verb of `fmt.Printf`.
It just dumps out the sequential bytes of the string as hexadecimal
digits, two per byte.

{{code "strings/basic.go" `/percent x/` `/percent x/`}}

Compare its output to that above:

	bdb23dbc20e28c98

A nice trick is to use the "space" flag in that format, putting a
space between the `%` and the `x`. Compare the format string
used here to the one above,

{{code "strings/basic.go" `/percent space x/` `/percent space x/`}}

and notice how the bytes come
out with spaces between, making the result a little less imposing:

	bd b2 3d bc 20 e2 8c 98

There's more. The `%q` (quoted) verb will escape any non-printable
byte sequences in a string so the output is unambiguous.

{{code "strings/basic.go" `/percent q/` `/percent q/`}}

This technique is handy when much of the string is
intelligible as text but there are peculiarities to root out; it produces:

	"\xbd\xb2=\xbc ⌘"

If we squint at that, we can see that buried in the noise is one ASCII equals sign,
along with a regular space, and at the end appears the well-known Swedish "Place of Interest"
symbol.
That symbol has Unicode value U+2318, encoded as UTF-8 by the bytes
after the space (hex value `20`): `e2` `8c` `98`.

If we are unfamiliar or confused by strange values in the string,
we can use the "plus" flag to the `%q` verb. This flag causes the output to escape
not only non-printable sequences, but also any non-ASCII bytes, all
while interpreting UTF-8.
The result is that it exposes the Unicode values of properly formatted UTF-8
that represents non-ASCII data in the string:

{{code "strings/basic.go" `/percent plus q/` `/percent plus q/`}}

With that format, the Unicode value of the Swedish symbol shows up as a
`\u` escape:

	"\xbd\xb2=\xbc \u2318"

These printing techniques are good to know when debugging
the contents of strings, and will be handy in the discussion that follows.
It's worth pointing out as well that all these methods behave exactly the
same for byte slices as they do for strings.

Here's the full set of printing options we've listed, presented as
a complete program you can run (and edit) right in the browser:

{{play "strings/basic.go" `/package/` `/^}/`}}

[Exercise: Modify the examples above to use a slice of bytes
instead of a string. Hint: Use a conversion to create the slice.]

[Exercise: Loop over the string using the `%q` format on each byte.
What does the output tell you?]

## UTF-8 and string literals

As we saw, indexing a string yields its bytes, not its characters: a string is just a
bunch of bytes.
That means that when we store a character value in a string,
we store its byte-at-a-time representation.
Let's look at a more controlled example to see how that happens.

Here's a simple program that prints a string constant with a single character
three different ways, once as a plain string, once as an ASCII-only quoted
string, and once as individual bytes in hexadecimal.
To avoid any confusion, we create a "raw string", enclosed by back quotes,
so it can contain only literal text. (Regular strings, enclosed by double
quotes, can contain escape sequences as we showed above.)

{{play "strings/utf8.go" `/^func/` `/^}/`}}

The output is:

	plain string: ⌘
	quoted string: "\u2318"
	hex bytes: e2 8c 98

which reminds us that the Unicode character value U+2318, the "Place
of Interest" symbol ⌘, is represented by the bytes `e2` `8c` `98`, and
that those bytes are the UTF-8 encoding of the hexadecimal
value 2318.

It may be obvious or it may be subtle, depending on your familiarity with
UTF-8, but it's worth taking a moment to explain how the UTF-8 representation
of the string was created.
The simple fact is: it was created when the source code was written.

Source code in Go is _defined_ to be UTF-8 text; no other representation is
allowed. That implies that when, in the source code, we write the text

	`⌘`

the text editor used to create the program places the UTF-8 encoding
of the symbol ⌘ into the source text.
When we print out the hexadecimal bytes, we're just dumping the
data the editor placed in the file.

In short, Go source code is UTF-8, so
_the source code for the string literal is UTF-8 text_.
If that string literal contains no escape sequences, which a raw
string cannot, the constructed string will hold exactly the
source text  between the quotes.
Thus by definition and
by construction the raw string will always contain a valid UTF-8
representation of its contents.
Similarly, unless it contains UTF-8-breaking escapes like those
from the previous section, a regular string literal will also always
contain valid UTF-8.

Some people think Go strings are always UTF-8, but they
are not: only string literals are UTF-8.
As we showed in the previous section, string _values_ can contain arbitrary
bytes;
as we showed in this one, string _literals_ always contain UTF-8 text
as long as they have no byte-level escapes.

To summarize, strings can contain arbitrary bytes, but when constructed
from string literals, those bytes are (almost always) UTF-8.

## Code points, characters, and runes

We've been very careful so far in how we use the words "byte" and "character".
That's partly because strings hold bytes, and partly because the idea of "character"
is a little hard to define.
The Unicode standard uses the term "code point" to refer to the item represented
by a single value.
The code point U+2318, with hexadecimal value 2318, represents the symbol ⌘.
(For lots more information about that code point, see
[its Unicode page](http://unicode.org/cldr/utility/character.jsp?a=2318).)

To pick a more prosaic example, the Unicode code point U+0061 is the lower
case Latin letter 'A': a.

But what about the lower case grave-accented letter 'A', à?
That's a character, and it's also a code point (U+00E0), but it has other
representations.
For example we can use the "combining" grave accent code point, U+0300,
and attach it to the lower case letter a, U+0061, to create the same character à.
In general, a character may be represented by a number of different
sequences of code points, and therefore different sequences of UTF-8 bytes.

The concept of character in computing is therefore ambiguous, or at least
confusing, so we use it with care.
To make things dependable, there are _normalization_ techniques that guarantee that
a given character is always represented by the same code points, but that
subject takes us too far off the topic for now.
A later blog post will explain how the Go libraries address normalization.

"Code point" is a bit of a mouthful, so Go introduces a shorter term for the
concept: _rune_.
The term appears in the libraries and source code, and means exactly
the same as "code point", with one interesting addition.

The Go language defines the word `rune` as an alias for the type `int32`, so
programs can be clear when an integer value represents a code point.
Moreover, what you might think of as a character constant is called a
_rune constant_ in Go.
The type and value of the expression

	'⌘'

is `rune` with integer value `0x2318`.

To summarize, here are the salient points:

  - Go source code is always UTF-8.
  - A string holds arbitrary bytes.
  - A string literal, absent byte-level escapes, always holds valid UTF-8 sequences.
  - Those sequences represent Unicode code points, called runes.
  - No guarantee is made in Go that characters in strings are normalized.

## Range loops

Besides the axiomatic detail that Go source code is UTF-8,
there's really only one way that Go treats UTF-8 specially, and that is when using
a `for` `range` loop on a string.

We've seen what happens with a regular `for` loop.
A `for` `range` loop, by contrast, decodes one UTF-8-encoded rune on each
iteration.
Each time around the loop, the index of the loop is the starting position of the
current rune, measured in bytes, and the code point is its value.
Here's an example using yet another handy `Printf` format, `%#U`, which shows
the code point's Unicode value and its printed representation:

{{play "strings/range.go" `/const/` `/}/`}}

The output shows how each code point occupies multiple bytes:

	U+65E5 '日' starts at byte position 0
	U+672C '本' starts at byte position 3
	U+8A9E '語' starts at byte position 6

[Exercise: Put an invalid UTF-8 byte sequence into the string. (How?)
What happens to the iterations of the loop?]

## Libraries

Go's standard library provides strong support for interpreting UTF-8 text.
If a `for` `range` loop isn't sufficient for your purposes,
chances are the facility you need is provided by a package in the library.

The most important such package is
[`unicode/utf8`](/pkg/unicode/utf8/),
which contains
helper routines to validate, disassemble, and reassemble UTF-8 strings.
Here is a program equivalent to the `for` `range` example above,
but using the `DecodeRuneInString` function from that package to
do the work.
The return values from the function are the rune and its width in
UTF-8-encoded bytes.

{{play "strings/encoding.go" `/const/` `/}/`}}

Run it to see that it performs the same.
The `for` `range` loop and `DecodeRuneInString` are defined to produce
exactly the same iteration sequence.

Look at the
[documentation](/pkg/unicode/utf8/)
for the `unicode/utf8` package to see what
other facilities it provides.

## Conclusion

To answer the question posed at the beginning: Strings are built from bytes
so indexing them yields bytes, not characters.
A string might not even hold characters.
In fact, the definition of "character" is ambiguous and it would
be a mistake to try to resolve the ambiguity by defining that strings are made
of characters.

There's much more to say about Unicode, UTF-8, and the world of multilingual
text processing, but it can wait for another post.
For now, we hope you have a better understanding of how Go strings behave
and that, although they may contain arbitrary bytes, UTF-8 is a central part
of their design.
