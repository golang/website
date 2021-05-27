// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command blog is a web server for the Go blog that can run on App Engine or
// as a stand-alone HTTP server.
package main

import (
	"flag"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/blog"
	_ "golang.org/x/tools/playground"
	"golang.org/x/website"
	"golang.org/x/website/internal/backport/httpfs"
)

var (
	httpAddr = flag.String("http", "localhost:8080", "HTTP listen address")
	reload   = flag.Bool("reload", false, "reload content on each page load")

	runningOnAppEngine = os.Getenv("GAE_ENV") != ""
	blogRoot           = "./"
)

const hostname = "blog.golang.org" // default hostname for blog server

var config = blog.Config{
	Hostname:      hostname,
	BaseURL:       "https://" + hostname,
	GodocURL:      "https://golang.org",
	HomeArticles:  5,  // articles to display on the home page
	FeedArticles:  10, // articles to include in Atom and JSON feeds
	PlayEnabled:   true,
	FeedTitle:     "The Go Programming Language Blog",
	ContentPath:   "_content/",
	TemplatePath:  "_template/",
	AnalyticsHTML: template.HTML(os.Getenv("BLOG_ANALYTICS")),
}

func main() {
	if runningOnAppEngine {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		*httpAddr = ":" + port
	}

	flag.Parse()

	if _, err := os.Stat("blog/_content/10years"); err == nil {
		blogRoot = "blog/"
	}
	config.ContentPath = blogRoot + config.ContentPath
	config.TemplatePath = blogRoot + config.TemplatePath

	h, err := blogHandler()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", h)

	ln, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on addr", *httpAddr)
	log.Fatal(http.Serve(ln, nil))
}

func blogHandler() (http.Handler, error) {
	var h http.Handler
	if *reload {
		h = http.HandlerFunc(reloadingBlogServer)
	} else {
		s, err := blog.NewServer(config)
		if err != nil {
			return nil, err
		}
		h = s
	}
	h = maybeStatic(h)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if runningOnAppEngine {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; preload")
		}
		h.ServeHTTP(w, r)
	})

	// Redirect "/blog/" to "/", because the menu bar link is to "/blog/"
	// but we're serving from the root.
	redirect := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	mux.HandleFunc("/blog", redirect)
	mux.HandleFunc("/blog/", redirect)

	// Keep these static file handlers in sync with app.yaml.
	static := http.FileServer(http.Dir(blogRoot + "static"))
	mux.Handle("/favicon.ico", static)
	mux.Handle("/fonts.css", static)
	mux.Handle("/fonts/", static)
	mux.Handle("/lib/godoc/", http.FileServer(httpfs.FS(website.Content)))

	// Redirects in redirect.go.
	for old, new := range redirects {
		new := new // for closure
		mux.HandleFunc(old, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/"+new, http.StatusMovedPermanently)
		})
	}

	return mux, nil
}

// maybeStatic serves from _static if possible,
// or else defers to the fallback handler.
func maybeStatic(fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if strings.Contains(p, ".") && !strings.HasSuffix(p, "/") {
			f := filepath.Join("_static", p)
			if _, err := os.Stat(f); err == nil {
				http.ServeFile(w, r, f)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}
}

// reloadingBlogServer is an handler that restarts the blog server on each page
// view. Inefficient; don't enable by default. Handy when editing blog content.
func reloadingBlogServer(w http.ResponseWriter, r *http.Request) {
	s, err := blog.NewServer(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.ServeHTTP(w, r)
}
