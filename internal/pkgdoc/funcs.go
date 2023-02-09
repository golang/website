// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkgdoc

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/printer"
	"go/token"
	"html/template"
	"io"
	"log"
	"path"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/website/internal/api"
	"golang.org/x/website/internal/texthtml"
)

var slashSlash = []byte("//")

// Node formats the given AST node as HTML.
// Identifiers in the rendered node
// are turned into links to their documentation.
func (p *Page) Node(node interface{}) template.HTML {
	var buf1 bytes.Buffer
	p.docs.writeNode(&buf1, p, p.fset, node)

	var buf2 bytes.Buffer
	n, _ := node.(ast.Node)
	buf2.Write(texthtml.Format(buf1.Bytes(), texthtml.Config{
		AST:        n,
		GoComments: true,
		OldDocs:    p.OldDocs,
	}))
	return template.HTML(buf2.String())
}

// NodeTOC formats the given AST node as HTML
// for inclusion in the table of contents.
func (p *Page) NodeTOC(node interface{}) template.HTML {
	var buf1 bytes.Buffer
	p.docs.writeNode(&buf1, p, p.fset, node)

	var buf2 bytes.Buffer
	buf2.Write(texthtml.Format(buf1.Bytes(), texthtml.Config{
		GoComments: true,
		OldDocs:    p.OldDocs,
	}))

	return sanitize(template.HTML(buf2.String()))
}

const tabWidth = 4

// writeNode writes the AST node x to w.
//
// The provided fset must be non-nil. The pageInfo is optional. If
// present, the pageInfo is used to add comments to struct fields to
// say which version of Go introduced them.
func (d *docs) writeNode(w io.Writer, pageInfo *Page, fset *token.FileSet, x interface{}) {
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
		gd.Tok == token.TYPE && len(gd.Specs) != 0 {
		pkgName = pageInfo.PDoc.ImportPath
		if ts, ok := gd.Specs[0].(*ast.TypeSpec); ok {
			if _, ok := ts.Type.(*ast.StructType); ok {
				structName = ts.Name.Name
			}
		}
		apiInfo = d.api[pkgName]
	}

	var out = w
	var buf bytes.Buffer
	if structName != "" {
		out = &buf
	}

	mode := printer.TabIndent | printer.UseSpaces
	err := (&printer.Config{Mode: mode, Tabwidth: tabWidth}).Fprint(tabSpacer(out, tabWidth), fset, x)
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

// Comment formats the given documentation comment as HTML.
func (p *Page) Comment(comment string) template.HTML {
	return template.HTML(p.PDoc.HTML(comment))
}

// sanitize sanitizes the argument src by replacing newlines with
// blanks, removing extra blanks, and by removing trailing whitespace
// and commas before closing parentheses.
func sanitize(src template.HTML) template.HTML {
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
	return template.HTML(buf[:j])
}

// Since reports the Go version that introduced the API feature
// identified by kind, receiver, name.
func (p *Page) Since(kind, receiver, name string) string {
	pkg := p.PDoc.ImportPath
	return p.docs.api.Func(pkg, kind, receiver, name)
}

type Example struct {
	Page   *Page
	Name   string
	Doc    string
	Code   template.HTML
	Play   string
	Output string
}

// Example renders the examples for the given function name as HTML.
func (p *Page) FmtExamples(funcName string) []*Example {
	var list []*Example
	for _, eg := range p.Examples {
		name := trimExampleSuffix(eg.Name)

		if name != funcName {
			continue
		}

		// print code
		cnode := &printer.CommentedNode{Node: eg.Code, Comments: eg.Comments}
		code := p.Node(cnode)
		out := eg.Output
		wholeFile := true

		// Additional formatting if this is a function body.
		if n := len(code); n >= 2 && code[0] == '{' && code[n-1] == '}' {
			wholeFile = false
			// remove surrounding braces
			code = code[1 : n-1]
			// unindent
			code = template.HTML(replaceLeadingIndentation(string(code), strings.Repeat(" ", tabWidth), ""))
			// remove output comment
			if loc := exampleOutputRx.FindStringIndex(string(code)); loc != nil {
				code = template.HTML(strings.TrimSpace(string(code)[:loc[0]]))
			}
		}

		// Write out the playground code in standard Go style
		// (use tabs, no comment highlight, etc).
		play := ""
		if eg.Play != nil {
			var buf bytes.Buffer
			eg.Play.Comments = filterOutBuildAnnotations(eg.Play.Comments)
			if err := format.Node(&buf, p.fset, eg.Play); err != nil {
				log.Print(err)
			} else {
				play = buf.String()
			}
		}

		// Drop output, as the output comment will appear in the code.
		if wholeFile && play == "" {
			out = ""
		}

		list = append(list, &Example{
			Page:   p,
			Name:   eg.Name,
			Doc:    eg.Doc,
			Code:   code,
			Play:   play,
			Output: out,
		})
	}
	return list
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

// ExampleName takes an example function name and returns its display
// name. For example, "Foo_Bar_quux" becomes "Foo.Bar (Quux)".
func (*Page) ExampleName(s string) string {
	name, suffix := splitExampleName(s)
	// replace _ with . for method names
	name = strings.Replace(name, "_", ".", 1)
	// use "Package" if no name provided
	if name == "" {
		name = "Package"
	}
	return name + suffix
}

// ExampleSuffix takes an example function name and returns its suffix in
// parenthesized form. For example, "Foo_Bar_quux" becomes " (Quux)".
func (*Page) ExampleSuffix(name string) string {
	_, suffix := splitExampleName(name)
	return suffix
}

// SrcPosLink returns a link to the specific source code position containing n,
// which must be either an ast.Node or a *doc.Note.
func (p *Page) SrcPosLink(n interface{}) template.HTML {
	// n must be an ast.Node or a *doc.Note
	var pos, end token.Pos

	switch n := n.(type) {
	case ast.Node:
		pos = n.Pos()
		end = n.End()
	case *doc.Note:
		pos = n.Pos
		end = n.End
	default:
		panic(fmt.Sprintf("wrong type for SrcPosLink template formatter: %T", n))
	}

	var relpath string
	var line int
	var low, high int // selection offset range

	if pos.IsValid() {
		xp := p.fset.Position(pos)
		relpath = xp.Filename
		line = xp.Line
		low = xp.Offset
	}
	if end.IsValid() {
		high = p.fset.Position(end).Offset
	}

	return srcPosLink(relpath, line, low, high)
}

func srcPosLink(s string, line, low, high int) template.HTML {
	s = path.Clean("/" + s)
	if !strings.HasPrefix(s, "/src/") {
		s = "/src" + s
	}
	var buf bytes.Buffer
	template.HTMLEscape(&buf, []byte(s))
	// selection ranges are of form "s=low:high"
	if low < high {
		fmt.Fprintf(&buf, "?s=%d:%d", low, high) // no need for URL escaping
		// if we have a selection, position the page
		// such that the selection is a bit below the top
		line -= 10
		if line < 1 {
			line = 1
		}
	}
	// line id's in html-printed source are of the
	// form "L%d" where %d stands for the line number
	if line > 0 {
		fmt.Fprintf(&buf, "#L%d", line) // no need for URL escaping
	}
	return template.HTML(buf.String())
}
