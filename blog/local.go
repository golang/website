// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine

// This file implements a stand-alone blog server.

package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"golang.org/x/tools/blog"
)

var (
	httpAddr     = flag.String("http", "localhost:8080", "HTTP listen address")
	contentPath  = flag.String("content", "content/", "path to content files")
	templatePath = flag.String("template", "template/", "path to template files")
	staticPath   = flag.String("static", "static/", "path to static files")
	reload       = flag.Bool("reload", false, "reload content on each page load")
)

func newServer(reload bool, staticPath string, config blog.Config) (http.Handler, error) {
	mux := http.NewServeMux()
	if reload {
		mux.HandleFunc("/", reloadingBlogServer)
	} else {
		s, err := blog.NewServer(config)
		if err != nil {
			return nil, err
		}
		mux.Handle("/", s)
	}
	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	return mux, nil
}

func main() {
	flag.Parse()
	config.ContentPath = *contentPath
	config.TemplatePath = *templatePath
	mux, err := newServer(*reload, *staticPath, config)
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on addr", *httpAddr)
	log.Fatal(http.Serve(ln, mux))
}

// reloadingBlogServer is an handler that restarts the blog server on each page
// view. Inefficient; don't enable by default. Handy when editing blog content.
func reloadingBlogServer(w http.ResponseWriter, r *http.Request) {
	s, err := blog.NewServer(config)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	s.ServeHTTP(w, r)
}
