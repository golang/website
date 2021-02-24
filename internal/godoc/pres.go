// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"net/http"
	"sync"
	"text/template"
)

// Presentation generates output from a corpus.
type Presentation struct {
	Corpus *Corpus

	mux        *http.ServeMux
	fileServer http.Handler

	DirlistHTML,
	ErrorHTML,
	ExampleHTML,
	GodocHTML,
	PackageHTML,
	PackageRootHTML *template.Template

	// GoogleCN reports whether this request should be marked GoogleCN.
	// If the function is nil, no requests are marked GoogleCN.
	GoogleCN func(*http.Request) bool

	// GoogleAnalytics optionally adds Google Analytics via the provided
	// tracking ID to each page.
	GoogleAnalytics string

	initFuncMapOnce sync.Once
	funcMap         template.FuncMap
	templateFuncs   template.FuncMap
}

// NewPresentation returns a new Presentation from a corpus.
func NewPresentation(c *Corpus) *Presentation {
	if c == nil {
		panic("nil Corpus")
	}
	p := &Presentation{
		Corpus:     c,
		mux:        http.NewServeMux(),
		fileServer: http.FileServer(http.FS(c.fs)),
	}
	docs := &docServer{
		p: p,
		d: NewDocTree(c.fs),
	}
	p.mux.Handle("/cmd/", docs)
	p.mux.Handle("/pkg/", docs)
	p.mux.HandleFunc("/", p.ServeFile)
	return p
}

func (p *Presentation) FileServer() http.Handler {
	return p.fileServer
}

func (p *Presentation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mux.ServeHTTP(w, r)
}

func (p *Presentation) googleCN(r *http.Request) bool {
	return p.GoogleCN != nil && p.GoogleCN(r)
}
