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
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/datastore"
	"golang.org/x/build/repos"
	"golang.org/x/website"
	"golang.org/x/website/internal/backport/archive/zip"
	"golang.org/x/website/internal/backport/html/template"
	"golang.org/x/website/internal/backport/io/fs"
	"golang.org/x/website/internal/backport/osfs"
	"golang.org/x/website/internal/codewalk"
	"golang.org/x/website/internal/dl"
	"golang.org/x/website/internal/env"
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

	mux := http.NewServeMux()
	mux.Handle("/", site)
	mux.Handle("/doc/codewalk/", codewalk.NewServer(fsys, site))
	mux.Handle("/fmt", http.HandlerFunc(fmtHandler))
	mux.Handle("/x/", http.HandlerFunc(xHandler))
	redirect.Register(mux)
	http.Handle("/", hostEnforcerHandler{mux})

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

// googleCN reports whether request r is considered
// to be served from golang.google.cn.
// TODO: This is duplicated within internal/proxy. Move to a common location.
func googleCN(r *http.Request) bool {
	if r.FormValue("googlecn") != "" {
		return true
	}
	if strings.HasSuffix(r.Host, ".cn") {
		return true
	}
	if !env.CheckCountry() {
		return false
	}
	switch r.Header.Get("X-Appengine-Country") {
	case "", "ZZ", "CN":
		return true
	}
	return false
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://blog.golang.org"+strings.TrimPrefix(r.URL.Path, "/blog"), http.StatusFound)
}

type fmtResponse struct {
	Body  string
	Error string
}

// fmtHandler takes a Go program in its "body" form value, formats it with
// standard gofmt formatting, and writes a fmtResponse as a JSON object.
func fmtHandler(w http.ResponseWriter, r *http.Request) {
	resp := new(fmtResponse)
	body, err := format.Source([]byte(r.FormValue("body")))
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Body = string(body)
	}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

// hostEnforcerHandler redirects http://foo.golang.org/bar to https://golang.org/bar.
// It permits golang.google.cn for China and *-dot-golang-org.appspot.com for testing.
type hostEnforcerHandler struct {
	h http.Handler
}

func (h hostEnforcerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !env.EnforceHosts() {
		h.h.ServeHTTP(w, r)
		return
	}
	if !h.isHTTPS(r) || !h.validHost(r.Host) {
		r.URL.Scheme = "https"
		if h.validHost(r.Host) {
			r.URL.Host = r.Host
		} else {
			r.URL.Host = "golang.org"
		}
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
		return
	}
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	h.h.ServeHTTP(w, r)
}

func (h hostEnforcerHandler) isHTTPS(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}

func (h hostEnforcerHandler) validHost(host string) bool {
	switch strings.ToLower(host) {
	case "golang.org", "golang.google.cn":
		return true
	}
	if strings.HasSuffix(host, "-dot-golang-org.appspot.com") {
		// staging/test
		return true
	}
	return false
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s\t%s", req.RemoteAddr, req.URL)
		h.ServeHTTP(w, req)
	})
}

func xHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/x/") {
		// Shouldn't happen if handler is registered correctly.
		http.Redirect(w, r, "https://pkg.go.dev/search?q=golang.org/x", http.StatusTemporaryRedirect)
		return
	}
	proj, suffix := strings.TrimPrefix(r.URL.Path, "/x/"), ""
	if i := strings.Index(proj, "/"); i != -1 {
		proj, suffix = proj[:i], proj[i:]
	}
	if proj == "" {
		http.Redirect(w, r, "https://pkg.go.dev/search?q=golang.org/x", http.StatusTemporaryRedirect)
		return
	}
	repo, ok := repos.ByGerritProject[proj]
	if !ok || !strings.HasPrefix(repo.ImportPath, "golang.org/x/") {
		http.NotFound(w, r)
		return
	}
	data := struct {
		Proj   string // Gerrit project ("net", "sys", etc)
		Suffix string // optional "/path" for requests like /x/PROJ/path
	}{proj, suffix}
	if err := xTemplate.Execute(w, data); err != nil {
		log.Println("xHandler:", err)
	}
}

var xTemplate = template.Must(template.New("x").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="golang.org/x/{{.Proj}} git https://go.googlesource.com/{{.Proj}}">
<meta name="go-source" content="golang.org/x/{{.Proj}} https://github.com/golang/{{.Proj}}/ https://github.com/golang/{{.Proj}}/tree/master{/dir} https://github.com/golang/{{.Proj}}/blob/master{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/golang.org/x/{{.Proj}}{{.Suffix}}">
</head>
<body>
<a href="https://pkg.go.dev/golang.org/x/{{.Proj}}{{.Suffix}}">Redirecting to documentation...</a>
</body>
</html>
`))

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
