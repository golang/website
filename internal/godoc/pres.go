// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"io/fs"
	"net/http"
	"text/template"

	"golang.org/x/website/internal/api"
	"golang.org/x/website/internal/pkgdoc"
)

// Presentation generates output from a file system.
type Presentation struct {
	fs  fs.FS
	api api.DB

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

	DocFuncs  template.FuncMap
	SiteFuncs template.FuncMap
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
	p.mux.HandleFunc("/", p.ServeFile)
	p.initFuncMap()

	if p.DirlistHTML, err = p.ReadTemplate("dirlist.html"); err != nil {
		return nil, err
	}
	if p.ErrorHTML, err = p.ReadTemplate("error.html"); err != nil {
		return nil, err
	}
	if p.ExampleHTML, err = p.ReadTemplate("example.html"); err != nil {
		return nil, err
	}
	if p.GodocHTML, err = p.ReadTemplate("godoc.html"); err != nil {
		return nil, err
	}
	if p.PackageHTML, err = p.ReadTemplate("package.html"); err != nil {
		return nil, err
	}
	if p.PackageRootHTML, err = p.ReadTemplate("packageroot.html"); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Presentation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mux.ServeHTTP(w, r)
}

func (p *Presentation) googleCN(r *http.Request) bool {
	return p.GoogleCN != nil && p.GoogleCN(r)
}

func (p *Presentation) ReadTemplate(name string) (*template.Template, error) {
	data, err := fs.ReadFile(p.fs, "lib/godoc/"+name)
	if err != nil {
		return nil, err
	}
	t, err := template.New(name).Funcs(p.SiteFuncs).Parse(string(data))
	if err != nil {
		return nil, err
	}
	return t, nil
}
