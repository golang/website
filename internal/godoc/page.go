// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"net/http"
	"runtime"
)

// A Page describes the contents of a webpage to be served.
//
// A Page's Methods are for use by the templates rendering the page.
type Page struct {
	Title    string // <h1>
	TabTitle string // prefix in <title>; defaults to Title
	Subtitle string // subtitle (date for spec, memory model)
	SrcPath  string // path to file in /src for text view

	// Template and Data describe the data to be
	// rendered into the overall site frame template.
	// If Template is empty, then Data should be a template.HTML
	// holding raw HTML to render into the site frame.
	// Otherwise, Template should be the name of a template file
	// in _content/lib/godoc (for example, "package.html"),
	// and that template will be executed
	// (with the *Page as its data argument) to produce HTML.
	//
	// The overall site template site.html is also invoked with
	// the *Page as its data argument. It is what arranges to call Template.
	Template string      // template to apply to data (empty string when Data is raw template.HTML)
	Data     interface{} // data to be rendered into page frame

	// Filled in automatically by ServePage
	GoogleCN        bool   // page is being served from golang.google.cn
	GoogleAnalytics string // Google Analytics tag
	Version         string // current Go version

	pres *Presentation
}

// fullPage returns a copy of page with the “automatic” fields filled in.
func (p *Presentation) fullPage(r *http.Request, page Page) Page {
	if page.TabTitle == "" {
		page.TabTitle = page.Title
	}
	page.Version = runtime.Version()
	page.GoogleCN = p.googleCN(r)
	page.GoogleAnalytics = p.GoogleAnalytics
	page.pres = p
	return page
}

// ServePage responds to the request with the content described by page.
func (p *Presentation) ServePage(w http.ResponseWriter, r *http.Request, page Page) {
	page = p.fullPage(r, page)
	applyTemplateToResponseWriter(w, p.Templates.Lookup("site.html"), &page)
}

// ServeError responds to the request with the given error.
func (p *Presentation) ServeError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusNotFound)
	p.ServePage(w, r, Page{
		Title:    r.URL.Path,
		Template: "error.html",
		Data:     err,
	})
}
