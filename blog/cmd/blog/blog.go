// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command blog is a web server for the Go blog that can run on App Engine or
// as a stand-alone HTTP server.
package main

import (
	"code.google.com/p/go.blog/pkg/blog"
	_ "code.google.com/p/go.tools/godoc/playground"
)

const hostname = "blog.golang.org" // default hostname for blog server

var config = &blog.Config{
	Hostname:     hostname,
	BaseURL:      "http://" + hostname,
	HomeArticles: 5,  // articles to display on the home page
	FeedArticles: 10, // articles to include in Atom and JSON feeds
}
