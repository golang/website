// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"html"
	"path"
	"strings"

	"golang.org/x/website/internal/backport/html/template"
	"golang.org/x/website/internal/pkgdoc"
)

var siteFuncs = template.FuncMap{
	// various helpers
	"basename": path.Base,

	// formatting of Examples
	"example_name":   example_name,
	"example_suffix": example_suffix,

	// Number operation
	"multiply": func(a, b int) int { return a * b },
}

// example_name takes an example function name and returns its display
// name. For example, "Foo_Bar_quux" becomes "Foo.Bar (Quux)".
func example_name(s string) string {
	name, suffix := pkgdoc.SplitExampleName(s)
	// replace _ with . for method names
	name = strings.Replace(name, "_", ".", 1)
	// use "Package" if no name provided
	if name == "" {
		name = "Package"
	}
	return name + suffix
}

// example_suffix takes an example function name and returns its suffix in
// parenthesized form. For example, "Foo_Bar_quux" becomes " (Quux)".
func example_suffix(name string) string {
	_, suffix := pkgdoc.SplitExampleName(name)
	return suffix
}

func srcToPkg(path string) string {
	// because of the irregular mapping under goroot
	// we need to correct certain relative paths
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "src/")
	path = strings.TrimPrefix(path, "pkg/")
	return "pkg/" + path
}

// SrcPkgLink builds an <a> tag linking to the package documentation
// for p.SrcPath.
func (p *Page) SrcPkgLink() template.HTML {
	dir := path.Dir(srcToPkg(p.SrcPath))
	if dir == "pkg" {
		return `<a href="/pkg">Index</a>`
	}
	dir = html.EscapeString(dir)
	return template.HTML(fmt.Sprintf(`<a href="/%s">%s</a>`, dir, dir[len("pkg/"):]))
}

// SrcBreadcrumb converts each segment of p.SrcPath to a HTML <a>.
// Each segment links to its corresponding src directories.
func (p *Page) SrcBreadcrumb() template.HTML {
	segments := strings.Split(p.SrcPath, "/")
	var buf bytes.Buffer
	var selectedSegment string
	var selectedIndex int

	if strings.HasSuffix(p.SrcPath, "/") {
		// relpath is a directory ending with a "/".
		// Selected segment is the segment before the last slash.
		selectedIndex = len(segments) - 2
		selectedSegment = segments[selectedIndex] + "/"
	} else {
		selectedIndex = len(segments) - 1
		selectedSegment = segments[selectedIndex]
	}

	for i := range segments[:selectedIndex] {
		buf.WriteString(fmt.Sprintf(`<a href="/%s">%s</a>/`,
			html.EscapeString(strings.Join(segments[:i+1], "/")),
			html.EscapeString(segments[i]),
		))
	}

	buf.WriteString(`<span class="text-muted">`)
	buf.WriteString(html.EscapeString(selectedSegment))
	buf.WriteString(`</span>`)
	return template.HTML(buf.String())
}

// SrcPosLink returns a link to the specific source code position containing n,
// which must be either an ast.Node or a *doc.Note.
// The current package is deduced from p.Data, which must be a *pkgdoc.Page.
func (p *Page) SrcPosLink(n interface{}) template.HTML {
	info := p.Data.(*pkgdoc.Page)
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
		p := info.FSet.Position(pos)
		relpath = p.Filename
		line = p.Line
		low = p.Offset
	}
	if end.IsValid() {
		high = info.FSet.Position(end).Offset
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
