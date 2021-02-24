// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// Page describes the contents of the top-level godoc webpage.
type Page struct {
	Title    string
	Tabtitle string
	Subtitle string
	SrcPath  string
	Query    string
	Body     []byte
	GoogleCN bool // page is being served from golang.google.cn

	// filled in by ServePage
	Version         string
	GoogleAnalytics string
}

func (p *Presentation) ServePage(w http.ResponseWriter, page Page) {
	if page.Tabtitle == "" {
		page.Tabtitle = page.Title
	}
	page.Version = runtime.Version()
	page.GoogleAnalytics = p.GoogleAnalytics
	applyTemplateToResponseWriter(w, p.GodocHTML, page)
}

func (p *Presentation) ServeError(w http.ResponseWriter, r *http.Request, relpath string, err error) {
	w.WriteHeader(http.StatusNotFound)
	if perr, ok := err.(*os.PathError); ok {
		rel, err := filepath.Rel(runtime.GOROOT(), perr.Path)
		if err != nil {
			perr.Path = "REDACTED"
		} else {
			perr.Path = filepath.Join("$GOROOT", rel)
		}
	}
	p.ServePage(w, Page{
		Title:           "File " + relpath,
		Subtitle:        relpath,
		Body:            applyTemplate(p.ErrorHTML, "errorHTML", err),
		GoogleCN:        p.googleCN(r),
		GoogleAnalytics: p.GoogleAnalytics,
	})
}
