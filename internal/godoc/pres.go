// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"html/template"
	"io/fs"
	"net/http"

	"golang.org/x/website/internal/api"
	"golang.org/x/website/internal/pkgdoc"
)

// Presentation is a website served from a file system.
type Presentation struct {
	fs  fs.FS
	api api.DB

	mux        *http.ServeMux
	fileServer http.Handler

	Templates *template.Template

	// GoogleCN reports whether this request should be marked GoogleCN.
	// If the function is nil, no requests are marked GoogleCN.
	GoogleCN func(*http.Request) bool

	// GoogleAnalytics optionally adds Google Analytics via the provided
	// tracking ID to each page.
	GoogleAnalytics string

	docFuncs template.FuncMap
}

// NewPresentation returns a new Presentation from a file system.
func NewPresentation(fsys fs.FS) (*Presentation, error) {
	apiDB, err := api.Load(fsys)
	if err != nil {
		return nil, err
	}
	p := &Presentation{
		fs:         fsys,
		api:        apiDB,
		mux:        http.NewServeMux(),
		fileServer: http.FileServer(http.FS(fsys)),
	}
	docs := &docServer{
		p: p,
		d: pkgdoc.NewDocs(fsys),
	}
	p.mux.Handle("/cmd/", docs)
	p.mux.Handle("/pkg/", docs)
	p.mux.HandleFunc("/", p.serveFile)
	p.initFuncMap()

	t, err := template.New("").Funcs(siteFuncs).ParseFS(fsys, "lib/godoc/*.html")
	if err != nil {
		return nil, err
	}
	p.Templates = t

	return p, nil
}

// ServeHTTP implements http.Handler, dispatching the request appropriately.
func (p *Presentation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mux.ServeHTTP(w, r)
}

func (p *Presentation) googleCN(r *http.Request) bool {
	return p.GoogleCN != nil && p.GoogleCN(r)
}
