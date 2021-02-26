// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package web

import (
	"log"
	"net/http"
	"path"
	"strings"

	"golang.org/x/website/internal/pkgdoc"
)

// docServer serves a package doc tree (/cmd or /pkg).
type docServer struct {
	p *Site
	d *pkgdoc.Docs
}

func (h *docServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if maybeRedirect(w, r) {
		return
	}

	// TODO(rsc): URL should be clean already.
	relpath := path.Clean(strings.TrimPrefix(r.URL.Path, "/pkg"))
	relpath = strings.TrimPrefix(relpath, "/")

	mode := pkgdoc.ParseMode(r.FormValue("m"))
	if relpath == "builtin" {
		// The fake built-in package contains unexported identifiers,
		// but we want to show them. Also, disable type association,
		// since it's not helpful for this fake package (see issue 6645).
		mode |= pkgdoc.ModeAll | pkgdoc.ModeBuiltin
	}
	info := pkgdoc.Doc(h.d, "src/"+relpath, mode, r.FormValue("GOOS"), r.FormValue("GOARCH"))
	if info.Err != nil {
		log.Print(info.Err)
		h.p.ServeError(w, r, info.Err)
		return
	}

	var tabtitle, title, subtitle string
	switch {
	case info.PDoc != nil:
		tabtitle = info.PDoc.Name
	default:
		tabtitle = info.Dirname
		title = "Directory "
	}
	if title == "" {
		if info.IsMain {
			// assume that the directory name is the command name
			_, tabtitle = path.Split(relpath)
			title = "Command "
		} else {
			title = "Package "
		}
	}
	title += tabtitle

	// special cases for top-level package/command directories
	switch tabtitle {
	case "/src":
		title = "Packages"
		tabtitle = "Packages"
	case "/src/cmd":
		title = "Commands"
		tabtitle = "Commands"
	}

	name := "package.html"
	if info.Dirname == "src" {
		name = "packageroot.html"
	}
	h.p.ServePage(w, r, Page{
		Title:    title,
		TabTitle: tabtitle,
		Subtitle: subtitle,
		Template: name,
		Data:     info,
	})
}

// ModeQuery returns the "?m=..." query for the current page.
// The page's Data must be a *pkgdoc.Page (to find the mode).
func (p *Page) ModeQuery() string {
	m := p.Data.(*pkgdoc.Page).Mode
	s := m.String()
	if s == "" {
		return ""
	}
	return "?m=" + s
}
