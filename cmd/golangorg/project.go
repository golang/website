// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/website/internal/history"
)

// projectHandler serves The Go Project page on /project/.
type projectHandler struct {
	ReleaseHistory []MajorRelease // Pre-computed release history to display.
}

func (h projectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/project/" {
		pres.ServeHTTP(w, req) // 404
		return
	}

	const relPath = "doc/contrib.html"

	src, err := vfs.ReadFile(fs, relPath)
	if err != nil {
		log.Printf("reading template %s: %v", relPath, err)
		pres.ServeError(w, req, relPath, err)
		return
	}

	meta, src, err := extractMetadata(src)
	if err != nil {
		log.Printf("decoding metadata %s: %v", relPath, err)
		pres.ServeError(w, req, relPath, err)
		return
	}
	if !meta.Template {
		err := fmt.Errorf("got non-template, want template")
		log.Printf("unexpected metadata %s: %v", relPath, err)
		pres.ServeError(w, req, relPath, err)
		return
	}

	page := godoc.Page{
		Title:    meta.Title,
		Subtitle: meta.Subtitle,
		GoogleCN: googleCN(req),
	}
	data := projectTemplateData{
		Major: h.ReleaseHistory,
	}

	// Evaluate as HTML template.
	tmpl, err := template.New("").Parse(string(src))
	if err != nil {
		log.Printf("parsing template %s: %v", relPath, err)
		pres.ServeError(w, req, relPath, err)
		return
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("executing template %s: %v", relPath, err)
		pres.ServeError(w, req, relPath, err)
		return
	}
	src = buf.Bytes()

	page.Body = src
	pres.ServePage(w, page)
}

// sortMajorReleases returns a sorted list of major Go releases,
// suitable to be displayed on the Go project page.
func sortMajorReleases(rs map[history.GoVer]history.Release) []MajorRelease {
	var major []MajorRelease
	for v, r := range rs {
		if !v.IsMajor() {
			continue
		}
		major = append(major, MajorRelease{ver: v, rel: r})
	}
	sort.Slice(major, func(i, j int) bool {
		if major[i].ver.X != major[j].ver.X {
			return major[i].ver.X > major[j].ver.X
		}
		return major[i].ver.Y > major[j].ver.Y
	})
	return major
}

type projectTemplateData struct {
	Major []MajorRelease
}

// MajorRelease represents a major Go release entry as displayed on the Go project page.
type MajorRelease struct {
	ver history.GoVer
	rel history.Release
}

// V returns the Go release version string, like "1.14", "1.14.1", "1.14.2", etc.
func (r MajorRelease) V() string {
	return r.ver.String()
}

// Date returns the date of the release, formatted for display on the Go project page.
func (r MajorRelease) Date() string {
	d := r.rel.Date
	return fmt.Sprintf("%s %d", d.Month, d.Year)
}
