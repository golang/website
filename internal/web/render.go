// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"golang.org/x/website/internal/tmplfunc"
)

// RenderContent returns the HTML rendering for the page using the named base template
// (the standard base template is "site.tmpl").
func (site *Site) RenderContent(p Page, tmpl string) (template.HTML, error) {
	html, err := site.renderHTML(p, tmpl, &http.Request{URL: &url.URL{Path: "/missingurl"}})
	if err != nil {
		return "", err
	}
	return template.HTML(html), nil
}

// renderHTML renders and returns the Content and framed HTML for the page.
func (site *Site) renderHTML(p Page, tmpl string, r *http.Request) ([]byte, error) {
	// Clone p, because we are going to set its Content key-value pair.
	p2 := make(Page)
	for k, v := range p {
		p2[k] = v
	}
	p = p2

	url, ok := p["URL"].(string)
	if !ok {
		// Set URL - caller did not.
		p["URL"] = r.URL.Path
	}
	file, _ := p["File"].(string)
	data, _ := p["FileData"].(string)

	// Load base template.
	base, err := site.readFile(".", tmpl)
	if err != nil {
		return nil, err
	}

	dir := strings.Trim(path.Dir(url), "/")
	if dir == "" {
		dir = "."
	}
	sd := &siteDir{site, dir}

	t := template.New("site.tmpl").Funcs(template.FuncMap{
		"add":          func(a, b int) int { return a + b },
		"sub":          func(a, b int) int { return a - b },
		"mul":          func(a, b int) int { return a * b },
		"div":          func(a, b int) int { return a / b },
		"code":         sd.code,
		"data":         sd.data,
		"page":         sd.page,
		"pages":        sd.pages,
		"play":         sd.play,
		"request":      func() *http.Request { return r },
		"path":         func() pkgPath { return pkgPath{} },
		"strings":      func() pkgStrings { return pkgStrings{} },
		"file":         sd.file,
		"first":        first,
		"markdown":     markdown,
		"raw":          raw,
		"yaml":         yamlFn,
		"presentStyle": presentStyle,
	})
	t.Funcs(site.funcs)

	if err := tmplfunc.Parse(t, string(base)); err != nil {
		return nil, err
	}

	// Load page-specific layout template.
	layout, _ := p["layout"].(string)
	if layout == "" {
		l, ok := site.findLayout(dir, "default")
		if ok {
			layout = l
		} else {
			layout = "none"
		}
	} else if path.IsAbs(layout) {
		layout = strings.TrimLeft(path.Clean(layout+".tmpl"), "/")
	} else if strings.Contains(layout, "/") {
		layout = path.Join(dir, layout+".tmpl")
	} else if layout != "none" {
		l, ok := site.findLayout(dir, layout)
		if !ok {
			return nil, fmt.Errorf("cannot find layout %q", layout)
		}
		layout = l
	}

	if layout != "none" {
		ldata, err := site.readFile(".", layout)
		if err != nil {
			return nil, err
		}
		if err := tmplfunc.Parse(t.New(layout), string(ldata)); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer
	if _, ok := p["Content"]; !ok && data != "" {
		// Load actual Markdown content (also a template).
		tf := t.New(file)
		if err := tmplfunc.Parse(tf, data); err != nil {
			return nil, err
		}
		if err := tf.Execute(&buf, p); err != nil {
			return nil, err
		}
		if strings.HasSuffix(file, ".md") {
			html, err := markdownToHTML(buf.String())
			if err != nil {
				return nil, err
			}
			p["Content"] = html
		} else {
			p["Content"] = template.HTML(buf.String())
		}
		buf.Reset()
	}

	if err := t.Execute(&buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// findLayout searches the start directory and parent directories for a template with the given base name.
func (site *Site) findLayout(dir, name string) (string, bool) {
	name += ".tmpl"
	for {
		abs := path.Join(dir, name)
		if _, err := fs.Stat(site.fs, abs); err == nil {
			return abs, true
		}
		if dir == "." {
			return "", false
		}
		dir = path.Dir(dir)
	}
}

// markdownToHTML converts Markdown to HTML.
// The Markdown source may contain raw HTML,
// but Go templates have already been processed.
func markdownToHTML(markdown string) (template.HTML, error) {
	// parser.WithHeadingAttribute allows custom ids on headings.
	// html.WithUnsafe allows use of raw HTML, which we need for tables.
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithHeadingAttribute(),
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(util.Prioritized(mdTransformFunc(mdLink), 1)),
		),
		goldmark.WithRendererOptions(html.WithUnsafe()),
		goldmark.WithExtensions(
			extension.NewTypographer(),
			extension.NewLinkify(
				extension.WithLinkifyAllowedProtocols([][]byte{[]byte("http"), []byte("https")}),
				extension.WithLinkifyEmailRegexp(regexp.MustCompile(`[^\x00-\x{10FFFF}]`)), // impossible
			),
			extension.DefinitionList,
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(replaceTabs([]byte(markdown)), &buf); err != nil {
		return "", err
	}
	return template.HTML(buf.Bytes()), nil
}

// mdTransformFunc is a func implementing parser.ASTTransformer.
type mdTransformFunc func(*ast.Document, text.Reader, parser.Context)

func (f mdTransformFunc) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	f(node, reader, pc)
}

// mdLink walks doc, adding rel=noreferrer target=_blank to non-relative links.
func mdLink(doc *ast.Document, _ text.Reader, _ parser.Context) {
	mdLinkWalk(doc)
}

func mdLinkWalk(n ast.Node) {
	switch n := n.(type) {
	case *ast.Link:
		dest := string(n.Destination)
		if strings.HasPrefix(dest, "https://") || strings.HasPrefix(dest, "http://") {
			n.SetAttributeString("rel", []byte("noreferrer"))
			n.SetAttributeString("target", []byte("_blank"))
		}
		return
	case *ast.AutoLink:
		// All autolinks are non-relative.
		n.SetAttributeString("rel", []byte("noreferrer"))
		n.SetAttributeString("target", []byte("_blank"))
		return
	}

	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		mdLinkWalk(child)
	}
}

// replaceTabs replaces all tabs in text with spaces up to a 4-space tab stop.
//
// In Markdown, tabs used for indentation are required to be interpreted as
// 4-space tab stops. See https://spec.commonmark.org/0.30/#tabs.
// Go also renders nicely and more compactly on the screen with 4-space
// tab stops, while browsers often use 8-space.
// And Goldmark crashes in some inputs that mix spaces and tabs.
// Fix the crashes and make the Go code consistently compact across browsers,
// all while staying Markdown-compatible, by expanding to 4-space tab stops.
//
// This function does not handle multi-codepoint Unicode sequences correctly.
func replaceTabs(text []byte) []byte {
	var buf bytes.Buffer
	col := 0
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]

		switch r {
		case '\n':
			buf.WriteByte('\n')
			col = 0

		case '\t':
			buf.WriteByte(' ')
			col++
			for col%4 != 0 {
				buf.WriteByte(' ')
				col++
			}

		default:
			buf.WriteRune(r)
			col++
		}
	}
	return buf.Bytes()
}
