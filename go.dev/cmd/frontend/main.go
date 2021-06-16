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
	"strings"

	"golang.org/x/website/go.dev/cmd/internal/site"
)

var discoveryHosts = map[string]string{
	"":               "pkg.go.dev",
	"dev.go.dev":     "dev-pkg.go.dev",
	"staging.go.dev": "staging-pkg.go.dev",
}

func main() {
	dir := "../.."
	if _, err := os.Stat("go.dev/_content/events.yaml"); err == nil {
		// Running in repo root.
		dir = "go.dev"
	}
	godev, err := site.Load(dir)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", addCSP(http.FileServer(godev)))
	http.Handle("/explore/", http.StripPrefix("/explore/", redirectHosts(discoveryHosts)))
	http.Handle("learn.go.dev/", http.HandlerFunc(redirectLearn))

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
	log.Print(http.Serve(l, nil))
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
