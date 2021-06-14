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

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/datastore"
	"golang.org/x/website"
	"golang.org/x/website/internal/backport/archive/zip"
	"golang.org/x/website/internal/backport/io/fs"
	"golang.org/x/website/internal/backport/osfs"
	"golang.org/x/website/internal/dl"
	"golang.org/x/website/internal/memcache"
	"golang.org/x/website/internal/proxy"
	"golang.org/x/website/internal/redirect"
	"golang.org/x/website/internal/short"
	"golang.org/x/website/internal/web"

	// Registers "/compile" handler that redirects to play.golang.org/compile.
	// If we are in prod we will register "golang.org/compile" separately,
	// which will get used instead.
	_ "golang.org/x/tools/playground"
)

var (
	httpAddr   = flag.String("http", "localhost:6060", "HTTP service address")
	verbose    = flag.Bool("v", false, "verbose mode")
	goroot     = flag.String("goroot", runtime.GOROOT(), "Go root directory")
	contentDir = flag.String("content", "", "path to _content directory")

	runningOnAppEngine = os.Getenv("PORT") != ""
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
	repoRoot := "../.."
	if _, err := os.Stat("_content"); err == nil {
		repoRoot = "."
	}

	if runningOnAppEngine {
		log.Print("golang.org server starting")
		*goroot = "_goroot.zip"
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		port := "8080"
		if p := os.Getenv("PORT"); p != "" {
			port = p
		}
		*httpAddr = ":" + port
	} else {
		if *contentDir == "" {
			*contentDir = filepath.Join(repoRoot, "_content")
		}
	}

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

	// Serve files from _content, falling back to GOROOT.
	var content fs.FS
	if *contentDir != "" {
		content = osfs.DirFS(*contentDir)
	} else {
		content = website.Content
	}

	var gorootFS fs.FS
	if strings.HasSuffix(*goroot, ".zip") {
		z, err := zip.OpenReader(*goroot)
		if err != nil {
			log.Fatal(err)
		}
		defer z.Close()
		gorootFS = z
	} else {
		gorootFS = osfs.DirFS(*goroot)
	}
	fsys := unionFS{content, gorootFS}

	site, err := web.NewSite(fsys)
	if err != nil {
		log.Fatalf("NewSite: %v", err)
	}
	site.GoogleCN = googleCN

	mux := registerHandlers(fsys, site)

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "User-agent: *\nDisallow: /search\n")
	})

	if err := redirect.LoadChangeMap(filepath.Join(repoRoot, "cmd/golangorg/hg-git-mapping.bin")); err != nil {
		log.Fatalf("LoadChangeMap: %v", err)
	}

	if runningOnAppEngine {
		appEngineSetup(site, mux)
	} else {
		// Register a redirect handler for /dl/ to the golang.org download page.
		mux.Handle("/dl/", http.RedirectHandler("https://golang.org/dl/", http.StatusFound))
	}

	var handler http.Handler = http.DefaultServeMux
	if *verbose {
		log.Printf("golang.org server:")
		log.Printf("\tversion = %s", runtime.Version())
		log.Printf("\taddress = %s", *httpAddr)
		log.Printf("\tgoroot = %s", *goroot)
		handler = loggingHandler(handler)
	}

	// Start http server.
	fmt.Fprintf(os.Stderr, "serving http://%s\n", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, handler); err != nil {
		log.Fatalf("ListenAndServe %s: %v", *httpAddr, err)
	}
}

var _ fs.ReadDirFS = unionFS{}

// A unionFS is an FS presenting the union of the file systems in the slice.
// If multiple file systems provide a particular file, Open uses the FS listed earlier in the slice.
// If multiple file systems provide a particular directory, ReadDir presents the
// concatenation of all the directories listed in the slice (with duplicates removed).
type unionFS []fs.FS

func (fsys unionFS) Open(name string) (fs.File, error) {
	var errOut error
	for _, sub := range fsys {
		f, err := sub.Open(name)
		if err == nil {
			// Note: Should technically check for directory
			// and return a synthetic directory that merges
			// reads from all the matching directories,
			// but all the directory reads in internal/godoc
			// come from fsys.ReadDir, which does that for us.
			// So we can ignore direct f.ReadDir calls.
			return f, nil
		}
		if errOut == nil {
			errOut = err
		}
	}
	return nil, errOut
}

func (fsys unionFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var all []fs.DirEntry
	var seen map[string]bool // seen[name] is true if name is listed in all; lazily initialized
	var errOut error
	for _, sub := range fsys {
		list, err := fs.ReadDir(sub, name)
		if err != nil {
			errOut = err
		}
		if len(all) == 0 {
			all = append(all, list...)
		} else {
			if seen == nil {
				// Initialize seen only after we get two different directory listings.
				seen = make(map[string]bool)
				for _, d := range all {
					seen[d.Name()] = true
				}
			}
			for _, d := range list {
				name := d.Name()
				if !seen[name] {
					seen[name] = true
					all = append(all, d)
				}
			}
		}
	}
	if len(all) > 0 {
		return all, nil
	}
	return nil, errOut
}

func appEngineSetup(site *web.Site, mux *http.ServeMux) {
	site.GoogleAnalytics = os.Getenv("GOLANGORG_ANALYTICS")

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

	dl.RegisterHandlers(mux, site, datastoreClient, memcacheClient)
	short.RegisterHandlers(mux, datastoreClient, memcacheClient)

	// Register /compile and /share handlers against the default serve mux
	// so that other app modules can make plain HTTP requests to those
	// hosts. (For reasons, HTTPS communication between modules is broken.)
	proxy.RegisterHandlers(http.DefaultServeMux)

	log.Println("AppEngine initialization complete")
}
