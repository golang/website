// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command blog is a web server for the Go blog that can run on App Engine or
// as a stand-alone HTTP server.
package main

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/tools/blog"
	"golang.org/x/tools/godoc/static"

	_ "golang.org/x/tools/playground"
)

const hostname = "blog.golang.org" // default hostname for blog server

var config = blog.Config{
	Hostname:     hostname,
	BaseURL:      "//" + hostname,
	GodocURL:     "//golang.org",
	HomeArticles: 5,  // articles to display on the home page
	FeedArticles: 10, // articles to include in Atom and JSON feeds
	PlayEnabled:  true,
	FeedTitle:    "The Go Programming Language Blog",
}

func init() {
	// Redirect "/blog/" to "/", because the menu bar link is to "/blog/"
	// but we're serving from the root.
	redirect := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	http.HandleFunc("/blog", redirect)
	http.HandleFunc("/blog/", redirect)

	http.Handle("/lib/godoc/", http.StripPrefix("/lib/godoc/", http.HandlerFunc(staticHandler)))
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	b, ok := static.Files[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, name, time.Time{}, strings.NewReader(b))
}
