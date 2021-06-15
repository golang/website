// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"bytes"
	"fmt"
	"html"
	"path"
	"strings"

	"golang.org/x/website/internal/backport/html/template"
)

var siteFuncs = template.FuncMap{
	// various helpers
	"basename": path.Base,

	// Number operation
	"multiply": func(a, b int) int { return a * b },
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
