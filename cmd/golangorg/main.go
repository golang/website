// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// golangorg: The Go Website (golang.org)

// Web server tree:
//
//	https://golang.org/			main landing page
//	https://golang.org/doc/	serve from content/doc, then $GOROOT/doc. spec, mem, etc.
//	https://golang.org/src/	serve files from $GOROOT/src; .go gets pretty-printed
//	https://golang.org/cmd/	serve documentation about commands
//	https://golang.org/pkg/	serve documentation about packages
//				(idea is if you say import "compress/zlib", you go to
//				https://golang.org/pkg/compress/zlib)
//

//go:build go1.16
// +build go1.16

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/gatefs"
	"golang.org/x/website"
)

var (
	httpAddr       = flag.String("http", "localhost:6060", "HTTP service address")
	verbose        = flag.Bool("v", false, "verbose mode")
	goroot         = flag.String("goroot", runtime.GOROOT(), "Go root directory")
	showTimestamps = flag.Bool("timestamps", false, "show timestamps with directory listings")
	templateDir    = flag.String("templates", "", "load templates/JS/CSS from disk in this directory (usually /path-to-website/content)")
	showPlayground = flag.Bool("play", true, "enable playground")
	declLinks      = flag.Bool("links", true, "link identifiers to their declarations")
	notesRx        = flag.String("notes", "BUG", "regular expression matching note markers to show")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: golangorg\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s\t%s", req.RemoteAddr, req.URL)
		h.ServeHTTP(w, req)
	})
}

func main() {
	earlySetup()

	flag.Usage = usage
	flag.Parse()

	// Check usage.
	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "Unexpected arguments.")
		usage()
	}
	if *httpAddr == "" {
		fmt.Fprintln(os.Stderr, "-http must be set")
		usage()
	}

	fsGate := make(chan bool, 20)

	// Determine file system to use.
	rootfs := gatefs.New(vfs.OS(*goroot), fsGate)
	fs.Bind("/", rootfs, "/", vfs.BindReplace)

	// Try serving files from _content before trying the main
	// go repository. This lets us update some documentation outside the
	// Go release cycle. This includes root.html, which redirects to "/".
	// See golang.org/issue/29206.
	if *templateDir != "" {
		fs.Bind("/", vfs.OS(*templateDir), "/", vfs.BindBefore)
	} else {
		fs.Bind("/", vfs.FromFS(website.Content), "/", vfs.BindBefore)
	}

	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = *verbose
	corpus.IndexEnabled = false
	if err := corpus.Init(); err != nil {
		log.Fatal(err)
	}
	// Initialize the version info before readTemplates, which saves
	// the map value in a method value.
	corpus.InitVersionInfo()

	pres = godoc.NewPresentation(corpus)
	pres.ShowTimestamps = *showTimestamps
	pres.ShowPlayground = *showPlayground
	pres.DeclLinks = *declLinks
	if *notesRx != "" {
		pres.NotesRx = regexp.MustCompile(*notesRx)
	}

	readTemplates(pres)
	mux := registerHandlers(pres)
	lateSetup(mux)

	var handler http.Handler = http.DefaultServeMux
	if *verbose {
		log.Printf("golang.org server:")
		log.Printf("\tversion = %s", runtime.Version())
		log.Printf("\taddress = %s", *httpAddr)
		log.Printf("\tgoroot = %s", *goroot)
		fs.Fprint(os.Stderr)
		handler = loggingHandler(handler)
	}

	// Start http server.
	fmt.Fprintf(os.Stderr, "serving http://%s\n", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, handler); err != nil {
		log.Fatalf("ListenAndServe %s: %v", *httpAddr, err)
	}
}
