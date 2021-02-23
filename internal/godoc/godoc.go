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
	pathpkg "path"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Fake relative package path for built-ins. Documentation for all globals
// (not just exported ones) will be shown for packages in this directory,
// and there will be no association of consts, vars, and factory functions
// with types (see issue 6645).
const builtinPkgPath = "builtin"

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
		"since":    p.Corpus.pkgAPIInfo.sinceVersionFunc,

		// formatting of AST nodes
		"node":         p.nodeFunc,
		"node_html":    p.node_htmlFunc,
		"comment_html": comment_htmlFunc,
		"sanitize":     sanitizeFunc,

		// support for URL attributes
		"pkgLink":       pkgLinkFunc,
		"srcLink":       srcLinkFunc,
		"posLink_url":   newPosLink_urlFunc(srcPosLinkFunc),
		"docLink":       docLinkFunc,
		"queryLink":     queryLinkFunc,
		"srcBreadcrumb": srcBreadcrumbFunc,
		"srcToPkgLink":  srcToPkgLinkFunc,

		// formatting of Examples
		"example_html":   p.example_htmlFunc,
		"example_name":   p.example_nameFunc,
		"example_suffix": p.example_suffixFunc,

		// formatting of Notes
		"noteTitle": noteTitle,

		// Number operation
		"multiply": multiply,

		// formatting of PageInfoMode query string
		"modeQueryString": modeQueryString,

		// check whether to display third party section or not
		"hasThirdParty": hasThirdParty,
	}
	if p.URLForSrc != nil {
		p.funcMap["srcLink"] = p.URLForSrc
	}
	if p.URLForSrcPos != nil {
		p.funcMap["posLink_url"] = newPosLink_urlFunc(p.URLForSrcPos)
	}
	if p.URLForSrcQuery != nil {
		p.funcMap["queryLink"] = p.URLForSrcQuery
	}
}

func multiply(a, b int) int { return a * b }

func filenameFunc(path string) string {
	_, localname := pathpkg.Split(path)
	return localname
}

type PageInfo struct {
	Dirname  string // directory containing the package
	Err      error  // error or nil
	GoogleCN bool   // page is being served from golang.google.cn

	Mode PageInfoMode // display metadata from query string

	// package info
	FSet       *token.FileSet         // nil if no package documentation
	PDoc       *doc.Package           // nil if no package documentation
	Examples   []*doc.Example         // nil if no example code
	Notes      map[string][]*doc.Note // nil if no package Notes
	PAst       map[string]*ast.File   // nil if no AST with package exports
	IsMain     bool                   // true for package main
	IsFiltered bool                   // true if results were filtered

	// directory info
	Dirs    *DirList  // nil if no directory information
	DirTime time.Time // directory time stamp
	DirFlat bool      // if set, show directory in a flat (non-indented) manner
}

func (info *PageInfo) IsEmpty() bool {
	return info.Err != nil || info.PAst == nil && info.PDoc == nil && info.Dirs == nil
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
	relpath = pathpkg.Dir(relpath)
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

func newPosLink_urlFunc(srcPosLinkFunc func(s string, line, low, high int) string) func(info *PageInfo, n interface{}) string {
	// n must be an ast.Node or a *doc.Note
	return func(info *PageInfo, n interface{}) string {
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
	s = pathpkg.Clean("/" + s)
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
	url := pathpkg.Clean("/"+s) + "?h=" + query
	if line > 0 {
		url += "#L" + strconv.Itoa(line)
	}
	return url
}

func docLinkFunc(s string, ident string) string {
	return pathpkg.Clean("/pkg/"+s) + "/#" + ident
}

func noteTitle(note string) string {
	return strings.Title(strings.ToLower(note))
}
