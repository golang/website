// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements a stand-alone blog server.

package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/tools/blog"
)

var (
	httpAddr     = flag.String("http", "localhost:8080", "HTTP listen address")
	contentPath  = flag.String("content", "content/", "path to content files")
	templatePath = flag.String("template", "template/", "path to template files")
	staticPath   = flag.String("static", "static/", "path to static files")
	godocPath    = flag.String("godoc", defaultGodocPath(), "path to lib/godoc static files")
	reload       = flag.Bool("reload", false, "reload content on each page load")
)

func defaultGodocPath() string {
	out, err := exec.Command("go", "list", "-f", "{{.Dir}}", "golang.org/x/tools/godoc").CombinedOutput()
	if err != nil {
		log.Printf("warning: locating -godoc directory: %v", err)
		return ""
	}
	dir := strings.TrimSpace(string(out))
	return filepath.Join(dir, "static")
}

// maybeStatic serves from one of the two static directories
// (-static and -godoc) if possible, or else defers to the fallback handler.
func maybeStatic(fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if strings.Contains(p, ".") && !strings.HasSuffix(p, "/") {
			f := filepath.Join(*staticPath, p)
			if _, err := os.Stat(f); err == nil {
				http.ServeFile(w, r, f)
				return
			}
		}
		if strings.HasPrefix(p, "/lib/godoc/") {
			f := filepath.Join(*godocPath, p[len("/lib/godoc/"):])
			if _, err := os.Stat(f); err == nil {
				http.ServeFile(w, r, f)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}
}

func newServer(reload bool, staticPath string, config blog.Config) (http.Handler, error) {
	mux := http.NewServeMux()
	var h http.Handler
	if reload {
		h = http.HandlerFunc(reloadingBlogServer)
	} else {
		s, err := blog.NewServer(config)
		if err != nil {
			return nil, err
		}
		h = s
	}
	mux.Handle("/", maybeStatic(h))
	return mux, nil
}

func main() {
	flag.Parse()

	if os.Getenv("GAE_ENV") == "standard" {
		log.Println("running in App Engine Standard mode")
		gaeMain()
		return
	}

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.ServeHTTP(w, r)
}
