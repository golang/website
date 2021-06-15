// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkgdoc

import "io"

var spaces = []byte("                                ") // 32 spaces seems like a good number

const (
	indenting = iota
	collecting
)

// tabSpacer returns a writer that passes writes through to w,
// expanding tabs to one or more spaces ending at a width-spaces-aligned boundary.
func tabSpacer(w io.Writer, width int) io.Writer {
	return &tconv{output: w, tabWidth: width}
}

// A tconv is an io.Writer filter for converting leading tabs into spaces.
type tconv struct {
	output   io.Writer
	state    int // indenting or collecting
	indent   int // valid if state == indenting
	tabWidth int
}

func (t *tconv) writeIndent() (err error) {
	i := t.indent
	for i >= len(spaces) {
		i -= len(spaces)
		if _, err = t.output.Write(spaces); err != nil {
			return
		}
	}
	// i < len(spaces)
	if i > 0 {
		_, err = t.output.Write(spaces[0:i])
	}
	return
}

func (t *tconv) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return
	}
	pos := 0 // valid if p.state == collecting
	var b byte
	for n, b = range data {
		switch t.state {
		case indenting:
			switch b {
			case '\t':
				t.indent += t.tabWidth
			case '\n':
				t.indent = 0
				if _, err = t.output.Write(data[n : n+1]); err != nil {
					return
				}
			case ' ':
				t.indent++
			default:
				t.state = collecting
				pos = n
				if err = t.writeIndent(); err != nil {
					return
				}
			}
		case collecting:
			if b == '\n' {
				t.state = indenting
				t.indent = 0
				if _, err = t.output.Write(data[pos : n+1]); err != nil {
					return
				}
			}
		}
	}
	n = len(data)
	if pos < n && t.state == collecting {
		_, err = t.output.Write(data[pos:])
	}
	return
}
