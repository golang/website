// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package web

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/printer"
	"html/template"
	"log"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/website/internal/pkgdoc"
)

// Example renders the examples for the given function name as HTML.
// The current package is deduced from p.Data, which must be a *pkgdoc.Page.
func (p *Page) Example(funcName string) template.HTML {
	info := p.Data.(*pkgdoc.Page)
	var buf bytes.Buffer
	for _, eg := range info.Examples {
		name := pkgdoc.TrimExampleSuffix(eg.Name)

		if name != funcName {
			continue
		}

		// print code
		cnode := &printer.CommentedNode{Node: eg.Code, Comments: eg.Comments}
		code := string(p.Node(cnode))
		out := eg.Output
		wholeFile := true

		// Additional formatting if this is a function body.
		if n := len(code); n >= 2 && code[0] == '{' && code[n-1] == '}' {
			wholeFile = false
			// remove surrounding braces
			code = code[1 : n-1]
			// unindent
			code = replaceLeadingIndentation(code, strings.Repeat(" ", tabWidth), "")
			// remove output comment
			if loc := exampleOutputRx.FindStringIndex(code); loc != nil {
				code = strings.TrimSpace(code[:loc[0]])
			}
		}

		// Write out the playground code in standard Go style
		// (use tabs, no comment highlight, etc).
		play := ""
		if eg.Play != nil {
			var buf bytes.Buffer
			eg.Play.Comments = filterOutBuildAnnotations(eg.Play.Comments)
			if err := format.Node(&buf, info.FSet, eg.Play); err != nil {
				log.Print(err)
			} else {
				play = buf.String()
			}
		}

		// Drop output, as the output comment will appear in the code.
		if wholeFile && play == "" {
			out = ""
		}

		t := p.site.Templates.Lookup("example.html")
		if t == nil {
			return ""
		}

		newPage := *p
		newPage.Data = struct {
			Name, Doc, Code, Play, Output string
		}{
			eg.Name, eg.Doc, code, play, out,
		}
		err := t.Execute(&buf, &newPage)
		if err != nil {
			log.Print(err)
		}
	}
	return template.HTML(buf.String())
}

// replaceLeadingIndentation replaces oldIndent at the beginning of each line
// with newIndent. This is used for formatting examples. Raw strings that
// span multiple lines are handled specially: oldIndent is not removed (since
// go/printer will not add any indentation there), but newIndent is added
// (since we may still want leading indentation).
func replaceLeadingIndentation(body, oldIndent, newIndent string) string {
	// Handle indent at the beginning of the first line. After this, we handle
	// indentation only after a newline.
	var buf bytes.Buffer
	if strings.HasPrefix(body, oldIndent) {
		buf.WriteString(newIndent)
		body = body[len(oldIndent):]
	}

	// Use a state machine to keep track of whether we're in a string or
	// rune literal while we process the rest of the code.
	const (
		codeState = iota
		runeState
		interpretedStringState
		rawStringState
	)
	searchChars := []string{
		"'\"`\n", // codeState
		`\'`,     // runeState
		`\"`,     // interpretedStringState
		"`\n",    // rawStringState
		// newlineState does not need to search
	}
	state := codeState
	for {
		i := strings.IndexAny(body, searchChars[state])
		if i < 0 {
			buf.WriteString(body)
			break
		}
		c := body[i]
		buf.WriteString(body[:i+1])
		body = body[i+1:]
		switch state {
		case codeState:
			switch c {
			case '\'':
				state = runeState
			case '"':
				state = interpretedStringState
			case '`':
				state = rawStringState
			case '\n':
				if strings.HasPrefix(body, oldIndent) {
					buf.WriteString(newIndent)
					body = body[len(oldIndent):]
				}
			}

		case runeState:
			switch c {
			case '\\':
				r, size := utf8.DecodeRuneInString(body)
				buf.WriteRune(r)
				body = body[size:]
			case '\'':
				state = codeState
			}

		case interpretedStringState:
			switch c {
			case '\\':
				r, size := utf8.DecodeRuneInString(body)
				buf.WriteRune(r)
				body = body[size:]
			case '"':
				state = codeState
			}

		case rawStringState:
			switch c {
			case '`':
				state = codeState
			case '\n':
				buf.WriteString(newIndent)
			}
		}
	}
	return buf.String()
}

var exampleOutputRx = regexp.MustCompile(`(?i)//[[:space:]]*(unordered )?output:`)

func filterOutBuildAnnotations(cg []*ast.CommentGroup) []*ast.CommentGroup {
	if len(cg) == 0 {
		return cg
	}

	for i := range cg {
		if !strings.HasPrefix(cg[i].Text(), "+build ") {
			// Found the first non-build tag, return from here until the end
			// of the slice.
			return cg[i:]
		}
	}

	// There weren't any non-build tags, return an empty slice.
	return []*ast.CommentGroup{}
}
