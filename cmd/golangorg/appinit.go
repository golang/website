// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.16
// +build golangorg

package main

// This file replaces main.go when running golangorg under App Engine.
// See README.md for details.

import (
	"context"
	"go/build"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/gatefs"
	"golang.org/x/website"
	"golang.org/x/website/internal/dl"
	"golang.org/x/website/internal/proxy"
	"golang.org/x/website/internal/redirect"
	"golang.org/x/website/internal/short"

	"cloud.google.com/go/datastore"
	"golang.org/x/website/internal/memcache"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	playEnabled = true

	log.Println("initializing golang.org server ...")

	fsGate := make(chan bool, 20)

	rootfs := gatefs.New(vfs.OS(runtime.GOROOT()), fsGate)
	fs.Bind("/", rootfs, "/", vfs.BindReplace)

	// Try serving files in /doc from a local copy before trying the main
	// go repository. This lets us update some documentation outside the
	// Go release cycle. This includes root.html, which redirects to "/".
	// See golang.org/issue/29206.
	fs.Bind("/doc", vfs.FromFS(website.Content), "/doc", vfs.BindBefore)
	fs.Bind("/lib/godoc", vfs.FromFS(website.Content), "/", vfs.BindReplace)

	webroot := getFullPath("/src/golang.org/x/website")
	fs.Bind("/favicon.ico", gatefs.New(vfs.OS(webroot), fsGate), "/favicon.ico", vfs.BindBefore)

	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = false
	corpus.MaxResults = 10000 // matches flag default in main.go
	corpus.IndexEnabled = false
	if err := corpus.Init(); err != nil {
		log.Fatal(err)
	}
	corpus.InitVersionInfo()

	pres = godoc.NewPresentation(corpus)
	pres.ShowPlayground = true
	pres.DeclLinks = true
	pres.NotesRx = regexp.MustCompile("BUG")
	pres.GoogleAnalytics = os.Getenv("GOLANGORG_ANALYTICS")

	readTemplates(pres)

	datastoreClient, memcacheClient := getClients()

	// NOTE(cbro): registerHandlers registers itself against DefaultServeMux.
	// The mux returned has host enforcement, so it's important to register
	// against this mux and not DefaultServeMux.
	mux := registerHandlers(pres)
	dl.RegisterHandlers(mux, datastoreClient, memcacheClient)
	short.RegisterHandlers(mux, datastoreClient, memcacheClient)

	// Register /compile and /share handlers against the default serve mux
	// so that other app modules can make plain HTTP requests to those
	// hosts. (For reasons, HTTPS communication between modules is broken.)
	proxy.RegisterHandlers(http.DefaultServeMux)

	http.HandleFunc("/_ah/health", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "User-agent: *\nDisallow: /search\n")
	})

	if err := redirect.LoadChangeMap("hg-git-mapping.bin"); err != nil {
		log.Fatalf("LoadChangeMap: %v", err)
	}

	log.Println("godoc initialization complete")

	// TODO(cbro): add instrumentation via opencensus.
	port := "8080"
	if p := os.Getenv("PORT"); p != "" { // PORT is set by GAE flex.
		port = p
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getFullPath(relPath string) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath + relPath
}

func getClients() (*datastore.Client, *memcache.Client) {
	ctx := context.Background()

	datastoreClient, err := datastore.NewClient(ctx, "")
	if err != nil {
		if strings.Contains(err.Error(), "missing project") {
			log.Fatalf("Missing datastore project. Set the DATASTORE_PROJECT_ID env variable. Use `gcloud beta emulators datastore` to start a local datastore.")
		}
		log.Fatalf("datastore.NewClient: %v.", err)
	}

	redisAddr := os.Getenv("GOLANGORG_REDIS_ADDR")
	if redisAddr == "" {
		log.Fatalf("Missing redis server for golangorg in production mode. set GOLANGORG_REDIS_ADDR environment variable.")
	}
	memcacheClient := memcache.New(redisAddr)
	return datastoreClient, memcacheClient
}
