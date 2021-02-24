// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/doc"
	"go/printer"
	"go/token"
	"io"
	"log"
	"unicode"

	"golang.org/x/website/internal/api"
	"golang.org/x/website/internal/texthtml"
)

var slashSlash = []byte("//")

func (p *Presentation) nodeFunc(info *PageInfo, node interface{}) string {
	var buf bytes.Buffer
	p.writeNode(&buf, info, info.FSet, node)
	return buf.String()
}

func (p *Presentation) node_htmlFunc(info *PageInfo, node interface{}, linkify bool) string {
	var buf1 bytes.Buffer
	p.writeNode(&buf1, info, info.FSet, node)

	var buf2 bytes.Buffer
	var n ast.Node
	if linkify {
		n, _ = node.(ast.Node)
	}
	buf2.Write(texthtml.Format(buf1.Bytes(), texthtml.Config{
		AST:        n,
		GoComments: true,
	}))
	return buf2.String()
}

const TabWidth = 4

// writeNode writes the AST node x to w.
//
// The provided fset must be non-nil. The pageInfo is optional. If
// present, the pageInfo is used to add comments to struct fields to
// say which version of Go introduced them.
func (p *Presentation) writeNode(w io.Writer, pageInfo *PageInfo, fset *token.FileSet, x interface{}) {
	// convert trailing tabs into spaces using a tconv filter
	// to ensure a good outcome in most browsers (there may still
	// be tabs in comments and strings, but converting those into
	// the right number of spaces is much harder)
	//
	// TODO(gri) rethink printer flags - perhaps tconv can be eliminated
	//           with an another printer mode (which is more efficiently
	//           implemented in the printer than here with another layer)

	var pkgName, structName string
	var apiInfo api.PkgDB
	if gd, ok := x.(*ast.GenDecl); ok && pageInfo != nil && pageInfo.PDoc != nil &&
		p.Corpus != nil &&
		gd.Tok == token.TYPE && len(gd.Specs) != 0 {
		pkgName = pageInfo.PDoc.ImportPath
		if ts, ok := gd.Specs[0].(*ast.TypeSpec); ok {
			if _, ok := ts.Type.(*ast.StructType); ok {
				structName = ts.Name.Name
			}
		}
		apiInfo = p.Corpus.pkgAPIInfo[pkgName]
	}

	var out = w
	var buf bytes.Buffer
	if structName != "" {
		out = &buf
	}

	mode := printer.TabIndent | printer.UseSpaces
	err := (&printer.Config{Mode: mode, Tabwidth: TabWidth}).Fprint(TabSpacer(out, TabWidth), fset, x)
	if err != nil {
		log.Print(err)
	}

	// Add comments to struct fields saying which Go version introduced them.
	if structName != "" {
		fieldSince := apiInfo.Field[structName]
		typeSince := apiInfo.Type[structName]
		// Add/rewrite comments on struct fields to note which Go version added them.
		var buf2 bytes.Buffer
		buf2.Grow(buf.Len() + len(" // Added in Go 1.n")*10)
		bs := bufio.NewScanner(&buf)
		for bs.Scan() {
			line := bs.Bytes()
			field := firstIdent(line)
			var since string
			if field != "" {
				since = fieldSince[field]
				if since != "" && since == typeSince {
					// Don't highlight field versions if they were the
					// same as the struct itself.
					since = ""
				}
			}
			if since == "" {
				buf2.Write(line)
			} else {
				if bytes.Contains(line, slashSlash) {
					line = bytes.TrimRight(line, " \t.")
					buf2.Write(line)
					buf2.WriteString("; added in Go ")
				} else {
					buf2.Write(line)
					buf2.WriteString(" // Go ")
				}
				buf2.WriteString(since)
			}
			buf2.WriteByte('\n')
		}
		w.Write(buf2.Bytes())
	}
}

// firstIdent returns the first identifier in x.
// This actually parses "identifiers" that begin with numbers too, but we
// never feed it such input, so it's fine.
func firstIdent(x []byte) string {
	x = bytes.TrimSpace(x)
	i := bytes.IndexFunc(x, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) })
	if i == -1 {
		return string(x)
	}
	return string(x[:i])
}

func comment_htmlFunc(comment string) string {
	var buf bytes.Buffer
	// TODO(gri) Provide list of words (e.g. function parameters)
	//           to be emphasized by ToHTML.
	doc.ToHTML(&buf, comment, nil) // does html-escaping
	return buf.String()
}

// sanitizeFunc sanitizes the argument src by replacing newlines with
// blanks, removing extra blanks, and by removing trailing whitespace
// and commas before closing parentheses.
func sanitizeFunc(src string) string {
	buf := make([]byte, len(src))
	j := 0      // buf index
	comma := -1 // comma index if >= 0
	for i := 0; i < len(src); i++ {
		ch := src[i]
		switch ch {
		case '\t', '\n', ' ':
			// ignore whitespace at the beginning, after a blank, or after opening parentheses
			if j == 0 {
				continue
			}
			if p := buf[j-1]; p == ' ' || p == '(' || p == '{' || p == '[' {
				continue
			}
			// replace all whitespace with blanks
			ch = ' '
		case ',':
			comma = j
		case ')', '}', ']':
			// remove any trailing comma
			if comma >= 0 {
				j = comma
			}
			// remove any trailing whitespace
			if j > 0 && buf[j-1] == ' ' {
				j--
			}
		default:
			comma = -1
		}
		buf[j] = ch
		j++
	}
	// remove trailing blank, if any
	if j > 0 && buf[j-1] == ' ' {
		j--
	}
	return string(buf[:j])
}
