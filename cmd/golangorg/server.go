// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Golangorg serves the golang.org web sites.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
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
	"golang.org/x/website/internal/webtest"

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
	if *contentDir == "" {
		repoRoot := "../.."
		if _, err := os.Stat("_content"); err == nil {
			repoRoot = "."
		}
		*contentDir = filepath.Join(repoRoot, "_content")
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

	handler := NewHandler(*contentDir, *goroot)

	handler = webtest.HandlerWithCheck(handler, "/_readycheck",
		filepath.Join(*contentDir, "../cmd/golangorg/testdata/*.txt"))

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

// NewHandler returns the http.Handler for the web site,
// given the directory where the content can be found
// (can be "", in which case an internal copy is used)
// and the directory or zip file of the GOROOT.
func NewHandler(contentDir, goroot string) http.Handler {
	// Serve files from _content, falling back to GOROOT.
	var content fs.FS
	if contentDir != "" {
		content = osfs.DirFS(contentDir)
	} else {
		content = website.Content
	}

	var gorootFS fs.FS
	if strings.HasSuffix(goroot, ".zip") {
		z, err := zip.OpenReader(goroot)
		if err != nil {
			log.Fatal(err)
		}
		gorootFS = &seekableFS{z}
	} else {
		gorootFS = osfs.DirFS(goroot)
	}
	fsys := unionFS{content, gorootFS}

	site, err := web.NewSite(fsys)
	if err != nil {
		log.Fatalf("NewSite: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", site)
	mux.Handle("/doc/codewalk/", codewalk.NewServer(fsys, site))
	mux.Handle("/fmt", http.HandlerFunc(fmtHandler))
	mux.Handle("/x/", http.HandlerFunc(xHandler))
	redirect.Register(mux)

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "User-agent: *\nDisallow: /search\n")
	})

	if err := redirect.LoadChangeMap(filepath.Join(contentDir, "../cmd/golangorg/hg-git-mapping.bin")); err != nil {
		log.Fatalf("LoadChangeMap: %v", err)
	}

	if runningOnAppEngine {
		appEngineSetup(site, mux)
	} else {
		// Register a redirect handler for /dl/ to the golang.org download page.
		mux.Handle("/dl/", http.RedirectHandler("https://golang.org/dl/", http.StatusFound))
	}

	var h http.Handler = mux
	if env.EnforceHosts() {
		h = hostEnforcerHandler(h)
	}
	h = hostPathHandler(h)
	return h
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
	proxy.RegisterHandlers(mux)

	log.Println("AppEngine initialization complete")
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

var validHosts = map[string]bool{
	"golang.org":       true,
	"golang.google.cn": true,
	"tip.golang.org":   true,
}

// hostEnforcerHandler redirects http://foo.golang.org/bar to https://golang.org/bar.
// It permits golang.google.cn as an alias for golang.org, for use in China.
func hostEnforcerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHTTPS := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" || r.URL.Scheme == "https"
		isValidHost := validHosts[strings.ToLower(r.Host)]

		if !isHTTPS || !isValidHost {
			r.URL.Scheme = "https"
			if isValidHost {
				r.URL.Host = r.Host
			} else {
				r.URL.Host = "golang.org"
			}
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
			return
		}
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		h.ServeHTTP(w, r)
	})
}

// hostPathHandler infers the host from the first element of the URL path
// when the actual host is a testing domain (localhost or *.appspot.com).
// It also rewrites the output HTML in that case to link back to URLs on
// the test site.
func hostPathHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "localhost" && !strings.HasPrefix(r.Host, "localhost:") && !strings.HasSuffix(r.Host, ".appspot.com") {
			h.ServeHTTP(w, r)
			return
		}

		elem, rest := strings.TrimPrefix(r.URL.Path, "/"), ""
		if i := strings.Index(elem, "/"); i >= 0 {
			elem, rest = elem[:i], elem[i+1:]
		}
		if !validHosts[elem] {
			u := "/golang.org" + r.URL.EscapedPath()
			if r.URL.RawQuery != "" {
				u += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, u, http.StatusTemporaryRedirect)
			return
		}

		r.Host = elem
		r.URL.Scheme = "https"
		r.URL.Host = elem
		r.URL.Path = "/" + rest
		lw := &linkRewriter{ResponseWriter: w, host: r.Host}
		h.ServeHTTP(lw, r)
		lw.Flush()
	})
}

// A linkRewriter is a ResponseWriter that rewrites links in HTML output.
// It rewrites relative links /foo to be /host/foo, and it rewrites any link
// https://h/foo, where h is in validHosts, to be /h/foo. This corrects the
// links to have the right form for the test server.
type linkRewriter struct {
	http.ResponseWriter
	host string
	buf  []byte
	ct   string // content-type
}

func (r *linkRewriter) Write(data []byte) (int, error) {
	if r.ct == "" {
		ct := r.Header().Get("Content-Type")
		if ct == "" {
			// Note: should use first 512 bytes, but first write is fine for our purposes.
			ct = http.DetectContentType(data)
		}
		r.ct = ct
	}
	if !strings.HasPrefix(r.ct, "text/html") {
		return r.ResponseWriter.Write(data)
	}
	r.buf = append(r.buf, data...)
	return len(data), nil
}

func (r *linkRewriter) Flush() {
	repl := []string{
		`href="/`, `href="/` + r.host + `/`,
	}
	for host := range validHosts {
		repl = append(repl, `href="https://`+host, `href="/`+host)
	}
	strings.NewReplacer(repl...).WriteString(r.ResponseWriter, string(r.buf))
	r.buf = nil
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s\t%s", r.RemoteAddr, r.URL)
		h.ServeHTTP(w, r)
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

// A seekableFS is an FS wrapper that makes every file seekable
// by reading it entirely into memory when it is opened and then
// serving read operations (including seek) from the memory copy.
type seekableFS struct {
	fs fs.FS
}

func (s *seekableFS) Open(name string) (fs.File, error) {
	f, err := s.fs.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if info.IsDir() {
		return f, nil
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	var sf seekableFile
	sf.File = f
	sf.Reset(data)
	return &sf, nil
}

// A seekableFile is a fs.File augmented by an in-memory copy of the file data to allow use of Seek.
type seekableFile struct {
	bytes.Reader
	fs.File
}

// Read calls f.Reader.Read.
// Both f.Reader and f.File have Read methods - a conflict - so f inherits neither.
// This method calls the one we want.
func (f *seekableFile) Read(b []byte) (int, error) {
	return f.Reader.Read(b)
}
