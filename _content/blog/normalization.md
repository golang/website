---
title: Text normalization in Go
date: 2013-11-26
by:
- Marcel van Lohuizen
tags:
- strings
- bytes
- runes
- characters
summary: How and why to normalize UTF-8 text in Go.
---

## Introduction

An earlier [post](https://blog.golang.org/strings) talked about strings, bytes
and characters in Go. I've been working on various packages for multilingual
text processing for the go.text repository. Several of these packages deserve a
separate blog post, but today I want to focus on
[go.text/unicode/norm](https://pkg.go.dev/golang.org/x/text/unicode/norm),
which handles normalization, a topic touched in the
[strings article](https://blog.golang.org/strings) and the subject of this
post. Normalization works at a higher level of abstraction than raw bytes.

To learn pretty much everything you ever wanted to know about normalization
(and then some), [Annex 15 of the Unicode Standard](http://unicode.org/reports/tr15/)
is a good read. A more approachable article is the corresponding
[Wikipedia page](http://en.wikipedia.org/wiki/Unicode_equivalence). Here we
focus on how normalization relates to Go.

## What is normalization?

There are often several ways to represent the same string. For example, an é
(e-acute) can be represented in a string as a single rune ("\u00e9") or an 'e'
followed by an acute accent ("e\u0301"). According to the Unicode standard,
these two are "canonically equivalent" and should be treated as equal.

Using a byte-to-byte comparison to determine equality would clearly not give
the right result for these two strings. Unicode defines a set of normal forms
such that if two strings are canonically equivalent and are normalized to the
same normal form, their byte representations are the same.

Unicode also defines a "compatibility equivalence" to equate characters that
represent the same characters, but may have a different visual appearance. For
example, the superscript digit '⁹' and the regular digit '9' are equivalent in
this form.

For each of these two equivalence forms, Unicode defines a composing and
decomposing form. The former replaces runes that can combine into a single rune
with this single rune. The latter breaks runes apart into their components.
This table shows the names, all starting with NF, by which the Unicode
Consortium identifies these forms:

{{raw (file "normalization/table1.html")}}

## Go's approach to normalization

As mentioned in the strings blog post, Go does not guarantee that characters in
a string are normalized. However, the go.text packages can compensate. For
example, the
[collate](https://pkg.go.dev/golang.org/x/text/collate) package, which
can sort strings in a language-specific way, works correctly even with
unnormalized strings. The packages in go.text do not always require normalized
input, but in general normalization may be necessary for consistent results.

Normalization isn't free but it is fast, particularly for collation and
searching or if a string is either in NFD or in NFC and can be converted to NFD
by decomposing without reordering its bytes. In practice,
[99.98%](http://www.macchiato.com/unicode/nfc-faq#TOC-How-much-text-is-already-NFC-) of
the web's HTML page content is in NFC form (not counting markup, in which case
it would be more). By far most NFC can be decomposed to NFD without the need
for reordering (which requires allocation). Also, it is efficient to detect
when reordering is necessary, so we can save time by doing it only for the rare
segments that need it.

To make things even better, the collation package typically does not use the
norm package directly, but instead uses the norm package to interleave
normalization information with its own tables. Interleaving the two problems
allows for reordering and normalization on the fly with almost no impact on
performance. The cost of on-the-fly normalization is compensated by not having
to normalize text beforehand and ensuring that the normal form is maintained
upon edits. The latter can be tricky. For instance, the result of concatenating
two NFC-normalized strings is not guaranteed to be in NFC.

Of course, we can also avoid the overhead outright if we know in advance that a
string is already normalized, which is often the case.

## Why bother?

After all this discussion about avoiding normalization, you might ask why it's
worth worrying about at all. The reason is that there are cases where
normalization is required and it is important to understand what those are, and
in turn how to do it correctly.

Before discussing those, we must first clarify the concept of 'character'.

## What is a character?

As was mentioned in the strings blog post, characters can span multiple runes.
For example, an 'e' and '◌́' (acute "\u0301") can combine to form 'é' ("e\u0301"
in NFD).  Together these two runes are one character. The definition of a
character may vary depending on the application. For normalization we will
define it as a sequence of runes that starts with a starter, a rune that does
not modify or combine backwards with any other rune, followed by possibly empty
sequence of non-starters, that is, runes that do (typically accents). The
normalization algorithm processes one character at a time.

Theoretically, there is no bound to the number of runes that can make up a
Unicode character. In fact, there are no restrictions on the number of
modifiers that can follow a character and a modifier may be repeated, or
stacked. Ever seen an 'e' with three acutes? Here you go: 'é́́'. That is a
perfectly valid 4-rune character according to the standard.

As a consequence, even at the lowest level, text needs to be processed in
increments of unbounded chunk sizes. This is especially awkward with a
streaming approach to text processing, as used by Go's standard Reader and
Writer interfaces, as that model potentially requires any intermediate buffers
to have unbounded size as well. Also, a straightforward implementation of
normalization will have a O(n²) running time.

There are really no meaningful interpretations for such large sequences of
modifiers for practical applications. Unicode defines a Stream-Safe Text
format, which allows capping the number of modifiers (non-starters) to at most
30, more than enough for any practical purpose. Subsequent modifiers will be
placed after a freshly inserted Combining Grapheme Joiner (CGJ or U+034F). Go
adopts this approach for all normalization algorithms. This decision gives up a
little conformance but gains a little safety.

## Writing in normal form

Even if you don't need to normalize text within your Go code, you might still
want to do so when communicating to the outside world. For example, normalizing
to NFC might compact your text, making it cheaper to send down a wire. For some
languages, like Korean, the savings can be substantial. Also, some external
APIs might expect text in a certain normal form. Or you might just want to fit
in and output your text as NFC like the rest of the world.

To write your text as NFC, use the
[unicode/norm](https://pkg.go.dev/golang.org/x/text/unicode/norm) package
to wrap your `io.Writer` of choice:

	wc := norm.NFC.Writer(w)
	defer wc.Close()
	// write as before...

If you have a small string and want to do a quick conversion, you can use this
simpler form:

	norm.NFC.Bytes(b)

Package norm provides various other methods for normalizing text.
Pick the one that suits your needs best.

## Catching look-alikes

Can you tell the difference between 'K' ("\u004B") and 'K' (Kelvin sign
"\u212A") or 'Ω' ("\u03a9") and 'Ω' (Ohm sign "\u2126")? It is easy to overlook
the sometimes minute differences between variants of the same underlying
character. It is generally a good idea to disallow such variants in identifiers
or anything where deceiving users with such look-alikes can pose a security
hazard.

The compatibility normal forms, NFKC and NFKD, will map many visually nearly
identical forms to a single value. Note that it will not do so when two symbols
look alike, but are really from two different alphabets. For example the Latin
'o', Greek 'ο', and Cyrillic 'о' are still different characters as defined by
these forms.

## Correct text modifications

The norm package might also come to the rescue when one needs to modify text.
Consider a case where you want to search and replace the word "cafe" with its
plural form "cafes".  A code snippet could look like this.

	s := "We went to eat at multiple cafe"
	cafe := "cafe"
	if p := strings.Index(s, cafe); p != -1 {
		p += len(cafe)
		s = s[:p] + "s" + s[p:]
	}
	fmt.Println(s)

This prints "We went to eat at multiple cafes" as desired and expected. Now
consider our text contains the French spelling "café" in NFD form:

	s := "We went to eat at multiple cafe\u0301"

Using the same code from above, the plural "s" would still be inserted after
the 'e', but before the acute, resulting in  "We went to eat at multiple
cafeś".  This behavior is undesirable.

The problem is that the code does not respect the boundaries between multi-rune
characters and inserts a rune in the middle of a character.  Using the norm
package, we can rewrite this piece of code as follows:

	s := "We went to eat at multiple cafe\u0301"
	cafe := "cafe"
	if p := strings.Index(s, cafe); p != -1 {
		p += len(cafe)
		if bp := norm.FirstBoundary(s[p:]); bp > 0 {
			p += bp
		}
		s = s[:p] + "s" + s[p:]
	}
	fmt.Println(s)

This may be a contrived example, but the gist should be clear. Be mindful of
the fact that characters can span multiple runes. Generally these kinds of
problems can be avoided by using search functionality that respects character
boundaries (such as the planned go.text/search package.)

## Iteration

Another tool provided by the norm package that may help dealing with character
boundaries is its iterator,
[`norm.Iter`](https://pkg.go.dev/golang.org/x/text/unicode/norm#Iter).
It iterates over characters one at a time in the normal form of choice.

## Performing magic

As mentioned earlier, most text is in NFC form, where base characters and
modifiers are combined into a single rune whenever possible.  For the purpose
of analyzing characters, it is often easier to handle runes after decomposition
into their smallest components. This is where the NFD form comes in handy. For
example, the following piece of code creates a `transform.Transformer` that
decomposes text into its smallest parts, removes all accents, and then
recomposes the text into NFC:

	import (
		"unicode"

		"golang.org/x/text/transform"
		"golang.org/x/text/unicode/norm"
	)

	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

The resulting `Transformer` can be used to remove accents from an `io.Reader`
of choice as follows:

	r = transform.NewReader(r, t)
	// read as before ...

This will, for example, convert any mention of "cafés" in the text to "cafes",
regardless of the normal form in which the original text was encoded.

## Normalization info

As mentioned earlier, some packages precompute normalizations into their tables
to minimize the need for normalization at run time. The type `norm.Properties`
provides access to the per-rune information needed by these packages, most
notably the Canonical Combining Class and decomposition information. Read the
[documentation](https://pkg.go.dev/golang.org/x/text/unicode/norm#Properties)
for this type if you want to dig deeper.

## Performance

To give an idea of the performance of normalization, we compare it against the
performance of strings.ToLower. The sample in the first row is both lowercase
and NFC and can in every case be returned as is. The second sample is neither
and requires writing a new version.

{{raw (file "normalization/table2.html")}}

The column with the results for the iterator shows both the measurement with
and without initialization of the iterator, which contain buffers that don't
need to be reinitialized upon reuse.

As you can see, detecting whether a string is normalized can be quite
efficient. A lot of the cost of normalizing in the second row is for the
initialization of buffers, the cost of which is amortized when one is
processing larger strings. As it turns out, these buffers are rarely needed, so
we may change the implementation at some point to speed up the common case for
small strings even further.

## Conclusion

If you're dealing with text inside Go, you generally do not have to use the
unicode/norm package to normalize your text. The package may still be useful
for things like ensuring that strings are normalized before sending them out or
to do advanced text manipulation.

This article briefly mentioned the existence of other go.text packages as well
as multilingual text processing and it may have raised more questions than it
has given answers. The discussion of these topics, however, will have to wait
until another day.
