// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/website/internal/history"
)

// releaseHandler serves the Release History page.
type releaseHandler struct {
	ReleaseHistory []Major // Pre-computed release history to display.
}

func (h releaseHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	const relPath = "doc/devel/release.html"

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
	data := releaseTemplateData{
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

// sortReleases returns a sorted list of Go releases, suitable to be
// displayed on the Release History page. Releases are arranged into
// major releases, each with minor revisions.
func sortReleases(rs map[history.GoVer]history.Release) []Major {
	var major []Major
	byMajorVersion := make(map[history.GoVer]Major)
	for v, r := range rs {
		switch {
		case v.IsMajor():
			m := byMajorVersion[v]
			m.Release = Release{ver: v, rel: r}
			byMajorVersion[v] = m
		case v.IsMinor():
			m := byMajorVersion[majorOf(v)]
			m.Minor = append(m.Minor, Release{ver: v, rel: r})
			byMajorVersion[majorOf(v)] = m
		}
	}
	for _, m := range byMajorVersion {
		sort.Slice(m.Minor, func(i, j int) bool { return m.Minor[i].ver.Z < m.Minor[j].ver.Z })
		major = append(major, m)
	}
	sort.Slice(major, func(i, j int) bool {
		if major[i].ver.X != major[j].ver.X {
			return major[i].ver.X > major[j].ver.X
		}
		return major[i].ver.Y > major[j].ver.Y
	})
	return major
}

// majorOf takes a Go version like 1.5, 1.5.1, 1.5.2, etc.,
// and returns the corresponding major version like 1.5.
func majorOf(v history.GoVer) history.GoVer {
	return history.GoVer{X: v.X, Y: v.Y, Z: 0}
}

type releaseTemplateData struct {
	Major []Major
}

// Major represents a major Go release and its minor revisions
// as displayed on the release history page.
type Major struct {
	Release
	Minor []Release
}

// Release represents a Go release entry as displayed on the release history page.
type Release struct {
	ver history.GoVer
	rel history.Release
}

// V returns the Go release version string, like "1.14", "1.14.1", "1.14.2", etc.
func (r Release) V() string {
	return r.ver.String()
}

// Date returns the date of the release, formatted for display on the release history page.
func (r Release) Date() string {
	d := r.rel.Date
	return fmt.Sprintf("%04d/%02d/%02d", d.Year, d.Month, d.Day)
}

// Released reports whether release r has been released.
func (r Release) Released() bool {
	return !r.rel.Future
}

func (r Release) Summary() (template.HTML, error) {
	var buf bytes.Buffer
	err := releaseSummaryHTML.Execute(&buf, releaseSummaryTemplateData{
		V:                     r.V(),
		Security:              r.rel.Security,
		Released:              r.Released(),
		Quantifier:            r.rel.Quantifier,
		ComponentsAndPackages: joinComponentsAndPackages(r.rel),
		More:                  r.rel.More,
		CustomSummary:         r.rel.CustomSummary,
	})
	return template.HTML(buf.String()), err
}

type releaseSummaryTemplateData struct {
	V                     string        // Go release version string, like "1.14", "1.14.1", "1.14.2", etc.
	Security              bool          // Security release.
	Released              bool          // Whether release has been released.
	Quantifier            string        // Optional quantifier. Empty string for unspecified amount of fixes (typical), "a" for a single fix, "two", "three" for multiple fixes, etc.
	ComponentsAndPackages template.HTML // Components and packages involved.
	More                  template.HTML // Additional release content.
	CustomSummary         template.HTML // CustomSummary, if non-empty, replaces the entire release content summary with custom HTML.
}

var releaseSummaryHTML = template.Must(template.New("").Parse(`
{{if not .CustomSummary}}
	{{if .Released}}includes{{else}}will include{{end}}
	{{.Quantifier}}
	{{if .Security}}security{{end}}
	{{if eq .Quantifier "a"}}fix{{else}}fixes{{end -}}
	{{with .ComponentsAndPackages}} to {{.}}{{end}}.
	{{.More}}

	{{if .Released}}
	See the
	<a href="https://github.com/golang/go/issues?q=milestone%3AGo{{.V}}+label%3ACherryPickApproved">Go
	{{.V}} milestone</a> on our issue tracker for details.
	{{end}}
{{else}}
	{{.CustomSummary}}
{{end}}
`))

// joinComponentsAndPackages joins components and packages involved
// in a Go release for the purposes of being displayed on the
// release history page, keeping English grammar rules in mind.
//
// The different special cases are:
//
// 	c1
// 	c1 and c2
// 	c1, c2, and c3
//
// 	the p1 package
// 	the p1 and p2 packages
// 	the p1, p2, and p3 packages
//
// 	c1 and [1 package]
// 	c1, and [2 or more packages]
// 	c1, c2, and [1 or more packages]
//
func joinComponentsAndPackages(r history.Release) template.HTML {
	var buf strings.Builder

	// List components, if any.
	for i, comp := range r.Components {
		if len(r.Packages) == 0 {
			// No packages, so components are joined with more rules.
			switch {
			case i != 0 && len(r.Components) == 2:
				buf.WriteString(" and ")
			case i != 0 && len(r.Components) >= 3 && i != len(r.Components)-1:
				buf.WriteString(", ")
			case i != 0 && len(r.Components) >= 3 && i == len(r.Components)-1:
				buf.WriteString(", and ")
			}
		} else {
			// When there are packages, all components are comma-separated.
			if i != 0 {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(string(comp))
	}

	// Join components and packages using a comma and/or "and" as needed.
	if len(r.Components) > 0 && len(r.Packages) > 0 {
		if len(r.Components)+len(r.Packages) >= 3 {
			buf.WriteString(",")
		}
		buf.WriteString(" and ")
	}

	// List packages, if any.
	if len(r.Packages) > 0 {
		buf.WriteString("the ")
	}
	for i, pkg := range r.Packages {
		switch {
		case i != 0 && len(r.Packages) == 2:
			buf.WriteString(" and ")
		case i != 0 && len(r.Packages) >= 3 && i != len(r.Packages)-1:
			buf.WriteString(", ")
		case i != 0 && len(r.Packages) >= 3 && i == len(r.Packages)-1:
			buf.WriteString(", and ")
		}
		buf.WriteString("<code>" + html.EscapeString(pkg) + "</code>")
	}
	switch {
	case len(r.Packages) == 1:
		buf.WriteString(" package")
	case len(r.Packages) >= 2:
		buf.WriteString(" packages")
	}

	return template.HTML(buf.String())
}
