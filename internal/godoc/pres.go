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
	cmdHandler handlerServer
	pkgHandler handlerServer

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
	p.cmdHandler = handlerServer{
		p:       p,
		c:       c,
		pattern: "/cmd/",
		fsRoot:  "/src",
	}
	p.pkgHandler = handlerServer{
		p:           p,
		c:           c,
		pattern:     "/pkg/",
		stripPrefix: "pkg/",
		fsRoot:      "/src",
		exclude:     []string{"/src/cmd"},
	}
	p.cmdHandler.registerWithMux(p.mux)
	p.pkgHandler.registerWithMux(p.mux)
	p.mux.HandleFunc("/", p.ServeFile)
	return p
}

func (p *Presentation) FileServer() http.Handler {
	return p.fileServer
}

func (p *Presentation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mux.ServeHTTP(w, r)
}

func (p *Presentation) PkgFSRoot() string {
	return p.pkgHandler.fsRoot
}

func (p *Presentation) CmdFSRoot() string {
	return p.cmdHandler.fsRoot
}

// TODO(bradfitz): move this to be a method on Corpus. Just moving code around for now,
// but this doesn't feel right.
func (p *Presentation) GetPkgPageInfo(abspath, relpath string, mode PageInfoMode) *PageInfo {
	return p.pkgHandler.GetPageInfo(abspath, relpath, mode, "", "")
}

// TODO(bradfitz): move this to be a method on Corpus. Just moving code around for now,
// but this doesn't feel right.
func (p *Presentation) GetCmdPageInfo(abspath, relpath string, mode PageInfoMode) *PageInfo {
	return p.cmdHandler.GetPageInfo(abspath, relpath, mode, "", "")
}

func (p *Presentation) googleCN(r *http.Request) bool {
	return p.GoogleCN != nil && p.GoogleCN(r)
}
