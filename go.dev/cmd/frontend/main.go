// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/website/internal/backport/html/template"
	"golang.org/x/website/internal/backport/osfs"
	"golang.org/x/website/internal/web"
	"golang.org/x/website/internal/webtest"
)

var discoveryHosts = map[string]string{
	"":               "pkg.go.dev",
	"dev.go.dev":     "dev-pkg.go.dev",
	"staging.go.dev": "staging-pkg.go.dev",
}

func main() {
	dir := "../../_content"
	if _, err := os.Stat("go.dev/_content/events.yaml"); err == nil {
		// Running in repo root.
		dir = "go.dev/_content"
	}

	h, err := NewHandler(dir)
	if err != nil {
		log.Fatal(err)
	}

	h = webtest.HandlerWithCheck(h, "/_readycheck",
		filepath.Join(dir, "cmd/frontend/testdata/*.txt"))

	addr := ":" + listenPort()
	if addr == ":0" {
		addr = "localhost:0"
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("net.Listen(%q, %q) = _, %v", "tcp", addr, err)
	}
	defer l.Close()
	log.Printf("Listening on http://%v/\n", l.Addr().String())
	log.Print(http.Serve(l, h))
}

func NewHandler(dir string) (http.Handler, error) {
	godev := web.NewSite(osfs.DirFS(dir))
	godev.Funcs(template.FuncMap{
		"newest":  newest,
		"section": section,
	})
	mux := http.NewServeMux()
	mux.Handle("/", addCSP(godev))
	mux.Handle("/explore/", http.StripPrefix("/explore/", redirectHosts(discoveryHosts)))
	mux.Handle("learn.go.dev/", http.HandlerFunc(redirectLearn))
	return mux, nil
}

// newest returns the pages sorted newest first,
// breaking ties by .linkTitle or else .title.
func newest(pages []web.Page) []web.Page {
	out := make([]web.Page, len(pages))
	copy(out, pages)

	sort.Slice(out, func(i, j int) bool {
		pi := out[i]
		pj := out[j]
		di, _ := pi["date"].(time.Time)
		dj, _ := pj["date"].(time.Time)
		if !di.Equal(dj) {
			return di.After(dj)
		}
		ti, _ := pi["linkTitle"].(string)
		if ti == "" {
			ti, _ = pi["title"].(string)
		}
		tj, _ := pj["linkTitle"].(string)
		if tj == "" {
			tj, _ = pj["title"].(string)
		}
		if ti != tj {
			return ti < tj
		}
		return false
	})
	return out
}

// section returns the site section for the given Page,
// defined as the first path element, or else an empty string.
// For example if p's URL is /x/y/z then section is "x".
func section(p web.Page) string {
	u, _ := p["URL"].(string)
	if !strings.HasPrefix(u, "/") {
		return ""
	}
	i := strings.Index(u[1:], "/")
	if i < 0 {
		return ""
	}
	return u[:1+i+1]
}

func redirectLearn(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://go.dev/learn/"+strings.TrimPrefix(r.URL.Path, "/"), http.StatusMovedPermanently)
}

func listenPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "0"
}

type redirectHosts map[string]string

func (rh redirectHosts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := &url.URL{Scheme: "https", Path: r.URL.Path, RawQuery: r.URL.RawQuery}
	if h, ok := rh[r.Host]; ok {
		u.Host = h
	} else if h, ok := rh[""]; ok {
		u.Host = h
	} else {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, u.String(), http.StatusFound)
}
