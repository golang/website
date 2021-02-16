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

// Some pages are being transitioned from $GOROOT to content/doc.
// See golang.org/issue/29206 and golang.org/issue/33637.

// +build go1.16
// +build !golangorg

package main

import (
	_ "expvar" // to serve /debug/vars
	"flag"
	"fmt"
	"go/build"
	"log"
	"net/http"
	_ "net/http/pprof" // to serve /debug/pprof/*
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/gatefs"
	"golang.org/x/website"
)

const defaultAddr = "localhost:6060" // default webserver address

var (
	// network
	httpAddr = flag.String("http", defaultAddr, "HTTP service address")

	verbose = flag.Bool("v", false, "verbose mode")

	// file system roots
	// TODO(gri) consider the invariant that goroot always end in '/'
	goroot = flag.String("goroot", findGOROOT(), "Go root directory")

	// layout control
	autoFlag       = flag.Bool("a", false, "update templates automatically")
	showTimestamps = flag.Bool("timestamps", false, "show timestamps with directory listings")
	templateDir    = flag.String("templates", "", "load templates/JS/CSS from disk in this directory (usually /path-to-website/content)")
	showPlayground = flag.Bool("play", false, "enable playground")
	declLinks      = flag.Bool("links", true, "link identifiers to their declarations")

	// source code notes
	notesRx = flag.String("notes", "BUG", "regular expression matching note markers to show")
)

func getFullPath(relPath string) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath + relPath
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: golangorg -http="+defaultAddr+"\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s\t%s", req.RemoteAddr, req.URL)
		h.ServeHTTP(w, req)
	})
}

func initCorpus(corpus *godoc.Corpus) {
	err := corpus.Init()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// Find templates in -a mode.
	if *autoFlag {
		if *templateDir != "" {
			fmt.Fprintln(os.Stderr, "Cannot use -a and -templates together.")
			usage()
		}
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			fmt.Fprintln(os.Stderr, "runtime.Caller failed: cannot find templates for -a mode.")
			os.Exit(2)
		}
		dir := filepath.Join(file, "../../../_content")
		if _, err := os.Stat(filepath.Join(dir, "godoc.html")); err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, "Cannot find templates for -a mode.")
			os.Exit(2)
		}
		*templateDir = dir
	}

	playEnabled = *showPlayground

	// Check usage.
	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "Unexpected arguments.")
		usage()
	}
	if *httpAddr == "" {
		fmt.Fprintln(os.Stderr, "-http must be set")
		usage()
	}

	// Set the resolved goroot.
	vfs.GOROOT = *goroot

	fsGate := make(chan bool, 20)

	// Determine file system to use.
	rootfs := gatefs.New(vfs.OS(*goroot), fsGate)
	fs.Bind("/", rootfs, "/", vfs.BindReplace)
	// Try serving files in /doc from a local copy before trying the main
	// go repository. This lets us update some documentation outside the
	// Go release cycle. This includes root.html, which redirects to "/".
	// See golang.org/issue/29206.
	if *templateDir != "" {
		fs.Bind("/doc", vfs.OS(*templateDir), "/doc", vfs.BindBefore)
		fs.Bind("/lib/godoc", vfs.OS(*templateDir), "/", vfs.BindBefore)
	} else {
		fs.Bind("/doc", vfs.FromFS(website.Content), "/doc", vfs.BindBefore)
		fs.Bind("/lib/godoc", vfs.FromFS(website.Content), "/", vfs.BindReplace)
	}

	// Bind $GOPATH trees into Go root.
	for _, p := range filepath.SplitList(build.Default.GOPATH) {
		fs.Bind("/src", gatefs.New(vfs.OS(p), fsGate), "/src", vfs.BindAfter)
	}

	webroot := getFullPath("/src/golang.org/x/website")
	fs.Bind("/robots.txt", gatefs.New(vfs.OS(webroot), fsGate), "/robots.txt", vfs.BindBefore)
	fs.Bind("/favicon.ico", gatefs.New(vfs.OS(webroot), fsGate), "/favicon.ico", vfs.BindBefore)

	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = *verbose

	go initCorpus(corpus)

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
	registerHandlers(pres)

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
