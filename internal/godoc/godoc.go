// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"path"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/website/internal/pkgdoc"
)

// FuncMap defines template functions used in godoc templates.
//
// Convention: template function names ending in "_html" or "_url" produce
//             HTML- or URL-escaped strings; all other function results may
//             require explicit escaping in the template.
func (p *Presentation) FuncMap() template.FuncMap {
	p.initFuncMapOnce.Do(p.initFuncMap)
	return p.funcMap
}

func (p *Presentation) TemplateFuncs() template.FuncMap {
	p.initFuncMapOnce.Do(p.initFuncMap)
	return p.templateFuncs
}

func (p *Presentation) initFuncMap() {
	if p.Corpus == nil {
		panic("nil Presentation.Corpus")
	}
	p.templateFuncs = template.FuncMap{
		"code": p.code,
	}
	p.funcMap = template.FuncMap{
		// various helpers
		"filename": filenameFunc,
		"since":    p.Corpus.pkgAPIInfo.Func,

		// formatting of AST nodes
		"node":         p.nodeFunc,
		"node_html":    p.node_htmlFunc,
		"comment_html": comment_htmlFunc,
		"sanitize":     sanitizeFunc,

		// support for URL attributes
		"pkgLink":       pkgLinkFunc,
		"srcLink":       srcLinkFunc,
		"posLink_url":   posLink_urlFunc,
		"docLink":       docLinkFunc,
		"queryLink":     queryLinkFunc,
		"srcBreadcrumb": srcBreadcrumbFunc,
		"srcToPkgLink":  srcToPkgLinkFunc,

		// formatting of Examples
		"example_html":   p.example_htmlFunc,
		"example_name":   p.example_nameFunc,
		"example_suffix": p.example_suffixFunc,

		// Number operation
		"multiply": multiply,

		// formatting of PageInfoMode query string
		"modeQueryString": modeQueryString,
	}
}

func multiply(a, b int) int { return a * b }

func filenameFunc(name string) string {
	_, localname := path.Split(name)
	return localname
}

func pkgLinkFunc(path string) string {
	// because of the irregular mapping under goroot
	// we need to correct certain relative paths
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "src/")
	path = strings.TrimPrefix(path, "pkg/")
	return "pkg/" + path
}

// srcToPkgLinkFunc builds an <a> tag linking to the package
// documentation of relpath.
func srcToPkgLinkFunc(relpath string) string {
	relpath = pkgLinkFunc(relpath)
	relpath = path.Dir(relpath)
	if relpath == "pkg" {
		return `<a href="/pkg">Index</a>`
	}
	return fmt.Sprintf(`<a href="/%s">%s</a>`, relpath, relpath[len("pkg/"):])
}

// srcBreadcrumbFun converts each segment of relpath to a HTML <a>.
// Each segment links to its corresponding src directories.
func srcBreadcrumbFunc(relpath string) string {
	segments := strings.Split(relpath, "/")
	var buf bytes.Buffer
	var selectedSegment string
	var selectedIndex int

	if strings.HasSuffix(relpath, "/") {
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
			strings.Join(segments[:i+1], "/"),
			segments[i],
		))
	}

	buf.WriteString(`<span class="text-muted">`)
	buf.WriteString(selectedSegment)
	buf.WriteString(`</span>`)
	return buf.String()
}

func posLink_urlFunc(info *pkgdoc.Page, n interface{}) string {
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
		panic(fmt.Sprintf("wrong type for posLink_url template formatter: %T", n))
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

	return srcPosLinkFunc(relpath, line, low, high)
}

func srcPosLinkFunc(s string, line, low, high int) string {
	s = srcLinkFunc(s)
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
	return buf.String()
}

func srcLinkFunc(s string) string {
	s = path.Clean("/" + s)
	if !strings.HasPrefix(s, "/src/") {
		s = "/src" + s
	}
	return s
}

// queryLinkFunc returns a URL for a line in a source file with a highlighted
// query term.
// s is expected to be a path to a source file.
// query is expected to be a string that has already been appropriately escaped
// for use in a URL query.
func queryLinkFunc(s, query string, line int) string {
	url := path.Clean("/"+s) + "?h=" + query
	if line > 0 {
		url += "#L" + strconv.Itoa(line)
	}
	return url
}

func docLinkFunc(s string, ident string) string {
	return path.Clean("/pkg/"+s) + "/#" + ident
}
