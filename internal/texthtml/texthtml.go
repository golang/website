// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package texthtml formats text files to HTML.
package texthtml

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/scanner"
	"go/token"
	"io"
	"regexp"
	"text/template"
)

// A Span describes a text span [start, end).
// The zero value of a Span is an empty span.
type Span struct {
	Start, End int
}

func (s *Span) isEmpty() bool { return s.Start >= s.End }

// A Selection is an "iterator" function returning a text span.
// Repeated calls to a selection return consecutive, non-overlapping,
// non-empty spans, followed by an infinite sequence of empty
// spans. The first empty span marks the end of the selection.
type Selection func() Span

// A Config configures how to format text as HTML.
type Config struct {
	Line       int       // if >= 1, number lines beginning with number Line, with <span class="ln">
	GoComments bool      // mark comments in Go text with <span class="comment">
	Playground bool      // format for playground sample
	Highlight  string    // highlight matches for this regexp with <span class="highlight">
	HL         string    // highlight lines that end with // HL (x/tools/present convention)
	Selection  Selection // mark selected spans with <span class="selection">
	AST        ast.Node  // link uses to declarations, assuming text is formatting of AST
	OldDocs    bool      // emit links to ?m=old docs
}

// Format formats text to HTML according to the configuration cfg.
func Format(text []byte, cfg Config) (html []byte) {
	var comments, highlights Selection
	if cfg.GoComments {
		comments = tokenSelection(text, token.COMMENT)
	}
	if cfg.Highlight != "" {
		highlights = regexpSelection(text, cfg.Highlight)
	}
	if cfg.HL != "" {
		highlights = hlSelection(text, cfg.HL)
	}

	var buf bytes.Buffer
	var idents Selection = Spans()
	var goLinks []goLink
	if cfg.AST != nil {
		idents = tokenSelection(text, token.IDENT)
		goLinks = goLinksFor(cfg.AST)
		if cfg.OldDocs {
			for i := range goLinks {
				goLinks[i].oldDocs = true
			}
		}
	}

	formatSelections(&buf, text, goLinks, comments, highlights, cfg.Selection, idents)

	if cfg.AST != nil {
		postFormatAST(&buf, cfg.AST)
	}

	trimSpaces(&buf)

	if cfg.Line > 0 {
		// Add line numbers in a separate pass.
		old := buf.Bytes()
		buf = bytes.Buffer{}
		n := cfg.Line
		for _, line := range bytes.Split(old, []byte("\n")) {
			// The line numbers are inserted into the document via a CSS ::before
			// pseudo-element. This prevents them from being copied when users
			// highlight and copy text.
			// ::before is supported in 98% of browsers: https://caniuse.com/#feat=css-gencontent
			// This is also the trick Github uses to hide line numbers.
			//
			// The first tab for the code snippet needs to start in column 9, so
			// it indents a full 8 spaces, hence the two nbsp's. Otherwise the tab
			// character only indents a short amount.
			//
			// Due to rounding and font width Firefox might not treat 8 rendered
			// characters as 8 characters wide, and subsequently may treat the tab
			// character in the 9th position as moving the width from (7.5 or so) up
			// to 8. See
			// https://github.com/webcompat/web-bugs/issues/17530#issuecomment-402675091
			// for a fuller explanation. The solution is to add a CSS class to
			// explicitly declare the width to be 8 characters.
			if cfg.Playground {
				fmt.Fprintf(&buf, `<span class="number">%2d&nbsp;&nbsp;</span>`, n)
			} else {
				fmt.Fprintf(&buf, `<span id="L%d" class="ln">%6d&nbsp;&nbsp;</span>`, n, n)
			}
			n++
			buf.Write(line)
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
}

// formatSelections takes a text and writes it to w using link and span
// writers lw and sw as follows: lw is invoked for consecutive span starts
// and ends as specified through the links selection, and sw is invoked for
// consecutive spans of text overlapped by the same selections as specified
// by selections.
func formatSelections(w io.Writer, text []byte, goLinks []goLink, selections ...Selection) {
	// compute the sequence of consecutive span changes
	changes := newMerger(selections)

	// The i'th bit in bitset indicates that the text
	// at the current offset is covered by selections[i].
	bitset := 0
	lastOffs := 0

	// Text spans are written in a delayed fashion
	// such that consecutive spans belonging to the
	// same selection can be combined (peephole optimization).
	// last describes the last span which has not yet been written.
	var last struct {
		begin, end int // valid if begin < end
		bitset     int
	}

	// flush writes the last delayed text span
	flush := func() {
		if last.begin < last.end {
			selectionTag(w, text[last.begin:last.end], last.bitset)
		}
		last.begin = last.end // invalidate last
	}

	// span runs the span [lastOffs, end) with the selection
	// indicated by bitset through the span peephole optimizer.
	span := func(end int) {
		if lastOffs < end { // ignore empty spans
			if last.end != lastOffs || last.bitset != bitset {
				// the last span is not adjacent to or
				// differs from the new one
				flush()
				// start a new span
				last.begin = lastOffs
			}
			last.end = end
			last.bitset = bitset
		}
	}

	linkEnd := ""
	for {
		// get the next span change
		index, offs, start := changes.next()
		if index < 0 || offs > len(text) {
			// no more span changes or the next change
			// is past the end of the text - we're done
			break
		}

		// format the previous selection span, determine
		// the new selection bitset and start a new span
		span(offs)
		if index == 3 { // Go link
			flush()
			if start {
				if len(goLinks) > 0 {
					start, end := goLinks[0].tags()
					io.WriteString(w, start)
					linkEnd = end
					goLinks = goLinks[1:]
				}
			} else {
				if linkEnd != "" {
					io.WriteString(w, linkEnd)
					linkEnd = ""
				}
			}
		} else {
			mask := 1 << uint(index)
			if start {
				bitset |= mask
			} else {
				bitset &^= mask
			}
		}
		lastOffs = offs
	}
	span(len(text))
	flush()
}

// A merger merges a slice of Selections and produces a sequence of
// consecutive span change events through repeated next() calls.
type merger struct {
	selections []Selection
	spans      []Span // spans[i] is the next span of selections[i]
}

const infinity int = 2e9

func newMerger(selections []Selection) *merger {
	spans := make([]Span, len(selections))
	for i, sel := range selections {
		spans[i] = Span{infinity, infinity}
		if sel != nil {
			if seg := sel(); !seg.isEmpty() {
				spans[i] = seg
			}
		}
	}
	return &merger{selections, spans}
}

// next returns the next span change: index specifies the Selection
// to which the span belongs, offs is the span start or end offset
// as determined by the start value. If there are no more span changes,
// next returns an index value < 0.
func (m *merger) next() (index, offs int, start bool) {
	// find the next smallest offset where a span starts or ends
	offs = infinity
	index = -1
	for i, seg := range m.spans {
		switch {
		case seg.Start < offs:
			offs = seg.Start
			index = i
			start = true
		case seg.End < offs:
			offs = seg.End
			index = i
			start = false
		}
	}
	if index < 0 {
		// no offset found => all selections merged
		return
	}
	// offset found - it's either the start or end offset but
	// either way it is ok to consume the start offset: set it
	// to infinity so it won't be considered in the following
	// next call
	m.spans[index].Start = infinity
	if start {
		return
	}
	// end offset found - consume it
	m.spans[index].End = infinity
	// advance to the next span for that selection
	seg := m.selections[index]()
	if !seg.isEmpty() {
		m.spans[index] = seg
	}
	return
}

// lineSelection returns the line spans for text as a Selection.
func lineSelection(text []byte) Selection {
	i, j := 0, 0
	return func() (seg Span) {
		// find next newline, if any
		for j < len(text) {
			j++
			if text[j-1] == '\n' {
				break
			}
		}
		if i < j {
			// text[i:j] constitutes a line
			seg = Span{i, j}
			i = j
		}
		return
	}
}

// tokenSelection returns, as a selection, the sequence of
// consecutive occurrences of token sel in the Go src text.
func tokenSelection(src []byte, sel token.Token) Selection {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)
	return func() (seg Span) {
		for {
			pos, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			offs := file.Offset(pos)
			if tok == sel {
				seg = Span{offs, offs + len(lit)}
				break
			}
		}
		return
	}
}

// Spans is a helper function to make a Selection from a slice of spans.
// Empty spans are discarded.
func Spans(spans ...Span) Selection {
	i := 0
	return func() Span {
		for i < len(spans) {
			s := spans[i]
			i++
			if s.Start < s.End {
				// non-empty
				return s
			}
		}
		return Span{}
	}
}

var hlRE = regexp.MustCompile(`(?m)\s*(.+)(\s+// (HL[a-zA-Z0-9_]*))$`)

// hlSelection returns the Selection for lines ending in // hl in text,
// also removing any // HLxxx from the text (overwriting with spaces)
func hlSelection(text []byte, hl string) Selection {
	lines := bytes.SplitAfter(text, []byte("\n"))
	off := 0
	var spans []Span
	for _, line := range lines {
		if m := hlRE.FindSubmatchIndex(line); m != nil {
			if string(line[m[6]:m[7]]) == hl {
				spans = append(spans, Span{off + m[2], off + m[3]})
			}
			for i := m[4]; i < m[5]; i++ {
				line[i] = ' '
			}
		}
		off += len(line)
	}
	return Spans(spans...)
}

// regexpSelection computes the Selection for the regular expression expr in text.
func regexpSelection(text []byte, expr string) Selection {
	var matches [][]int
	if rx, err := regexp.Compile(expr); err == nil {
		matches = rx.FindAllIndex(text, -1)
	}
	var spans []Span
	for _, m := range matches {
		spans = append(spans, Span{m[0], m[1]})
	}
	return Spans(spans...)
}

// Span tags for all the possible selection combinations that may
// be generated by FormatText. Selections are indicated by a bitset,
// and the value of the bitset specifies the tag to be used.
//
// bit 0: comments
// bit 1: highlights
// bit 2: selections
//
var startTags = [][]byte{
	/* 000 */ []byte(``),
	/* 001 */ []byte(`<span class="comment">`),
	/* 010 */ []byte(`<span class="highlight">`),
	/* 011 */ []byte(`<span class="highlight-comment">`),
	/* 100 */ []byte(`<span class="selection">`),
	/* 101 */ []byte(`<span class="selection-comment">`),
	/* 110 */ []byte(`<span class="selection-highlight">`),
	/* 111 */ []byte(`<span class="selection-highlight-comment">`),
}

var endTag = []byte(`</span>`)

func selectionTag(w io.Writer, text []byte, selections int) {
	if selections < len(startTags) {
		if tag := startTags[selections]; len(tag) > 0 {
			w.Write(tag)
			template.HTMLEscape(w, text)
			w.Write(endTag)
			return
		}
	}
	template.HTMLEscape(w, text)
}

// trimSpaces removes trailing spaces at the end of each line in buf.
func trimSpaces(buf *bytes.Buffer) {
	data := buf.Bytes()
	out := data[:0]
	for len(data) > 0 {
		j := bytes.IndexByte(data, '\n')
		if j < 0 {
			j = len(data)
		}
		var line []byte
		line, data = data[:j], data[j:]
		for len(line) > 0 && (line[len(line)-1] == ' ' || line[len(line)-1] == '\t') {
			line = line[:len(line)-1]
		}
		out = append(out, line...)
		if len(data) > 0 {
			out = append(out, '\n')
			data = data[1:]
		}
	}
	buf.Truncate(len(out))
}
