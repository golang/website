// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Golangorg serves the golang.org web sites.
package main

import (
	"archive/zip"
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/build/repos"
	"golang.org/x/website"
	"golang.org/x/website/internal/blog"
	"golang.org/x/website/internal/codewalk"
	"golang.org/x/website/internal/dl"
	"golang.org/x/website/internal/gitfs"
	"golang.org/x/website/internal/history"
	"golang.org/x/website/internal/memcache"
	"golang.org/x/website/internal/pkgdoc"
	"golang.org/x/website/internal/play"
	"golang.org/x/website/internal/redirect"
	"golang.org/x/website/internal/short"
	"golang.org/x/website/internal/talks"
	"golang.org/x/website/internal/tour"
	"golang.org/x/website/internal/web"
	"golang.org/x/website/internal/webtest"
)

var (
	httpAddr   = flag.String("http", "localhost:6060", "HTTP service address")
	verbose    = flag.Bool("v", false, "verbose mode")
	goroot     = flag.String("goroot", runtime.GOROOT(), "Go root directory")
	contentDir = flag.String("content", "", "path to _content directory")

	runningOnAppEngine = false // os.Getenv("PORT") != ""

	tipFlag = flag.Bool("tip", runningOnAppEngine, "load git content for tip.golang.org")

	googleAnalytics string
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: golangorg\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	// Running locally, find the local _content directory,
	// so that updates to those files appear on the local dev instance without restarting.
	// On App Engine, leave contentDir empty, so we use the embedded copy,
	// which is much faster to access than the simulated file system.
	// if *contentDir == "" && !runningOnAppEngine {
	// 	repoRoot := "../.."
	// 	if _, err := os.Stat("_content"); err == nil {
	// 		repoRoot = "."
	// 	}
	// 	*contentDir = filepath.Join(repoRoot, "_content")
	// }

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
		testdataFS, "testdata/*.txt")

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

//go:embed testdata
var testdataFS embed.FS

// NewHandler returns the http.Handler for the web site,
// given the directory where the content can be found
// (can be "", in which case an internal copy is used)
// and the directory or zip file of the GOROOT.
func NewHandler(contentDir, goroot string) http.Handler {
	mux := http.NewServeMux()

	// Serve files from _content, falling back to GOROOT.

	// Use explicit contentDir if specified, otherwise embedded copy.
	var contentFS fs.FS
	if contentDir != "" {
		contentFS = os.DirFS(contentDir)
	} else {
		contentFS = website.Content()
	}

	var gorootFS fs.FS
	if strings.HasSuffix(goroot, ".zip") {
		z, err := zip.OpenReader(goroot)
		if err != nil {
			log.Fatal(err)
		}
		gorootFS = &seekableFS{z}
	} else {
		gorootFS = os.DirFS(goroot)
	}

	// tip.golang.org serves content from the very latest Git commit
	// of the main Go repo, instead of the one the app is bundled with.
	var tipGoroot atomicFS
	if _, err := newSite(mux, "tip.golang.org", contentFS, &tipGoroot); err != nil {
		log.Fatalf("loading tip site: %v", err)
	}
	if *tipFlag {
		go watchTip(&tipGoroot)
	}

	// beta.golang.org is an old name for tip.
	mux.Handle("beta.golang.org/", redirectPrefix("https://tip.golang.org/"))

	// By default, golang.org/foo redirects to go.dev/foo.
	// All the user-facing golang.org subdomains have moved to go.dev subdirectories.
	// There are some exceptions below, like for golang.org/x.
	mux.Handle("golang.org/", redirectPrefix("https://go.dev/"))
	mux.Handle("blog.golang.org/", redirectPrefix("https://go.dev/blog/"))
	mux.Handle("learn.go.dev/", redirectPrefix("https://go.dev/learn/"))
	mux.Handle("talks.golang.org/", redirectPrefix("https://go.dev/talks/"))
	mux.Handle("tour.golang.org/", redirectPrefix("https://go.dev/tour/"))

	// Redirect subdirectory-like domains to the actual subdirectories,
	// for people whose fingers learn to type go.dev instead of golang.org
	// but not the rest of the URL schema change.
	// Note that these domains have to be listed in knownHosts below as well.
	mux.Handle("blog.go.dev/", redirectPrefix("https://go.dev/blog/"))
	mux.Handle("play.go.dev/", redirectPrefix("https://go.dev/play/"))
	mux.Handle("talks.go.dev/", redirectPrefix("https://go.dev/talks/"))
	mux.Handle("tour.go.dev/", redirectPrefix("https://go.dev/tour/"))

	// m.golang.org is an old shortcut for golang.org mail.
	// Gmail itself can serve this redirect, but only on HTTP (not HTTPS).
	// Golang.org's HSTS header tells browsers to use HTTPS for all subdomains,
	// which broke the redirect.
	mux.Handle("m.golang.org/", http.RedirectHandler("https://mail.google.com/a/golang.org/", http.StatusMovedPermanently))

	// Register a redirect handler for tip.golang.org/dl/ to the golang.org download page.
	// (golang.org/dl and golang.google.cn/dl are registered separately.)
	mux.Handle("tip.golang.org/dl/", http.RedirectHandler("https://go.dev/dl/", http.StatusFound))

	// TODO(rsc): The unionFS is a hack until we move the files in a followup CL.
	siteMux := http.NewServeMux()
	godevSite, err := newSite(siteMux, "", contentFS, gorootFS)
	if err != nil {
		log.Fatalf("newSite go.dev: %v", err)
	}
	chinaSite, err := newSite(siteMux, "golang.google.cn", contentFS, gorootFS)
	if err != nil {
		log.Fatalf("newSite golang.google.cn: %v", err)
	}
	if runningOnAppEngine {
		appEngineSetup(mux)
	}
	dl.RegisterHandlers(siteMux, godevSite, "", datastoreClient, memcacheClient)
	dl.RegisterHandlers(siteMux, chinaSite, "golang.google.cn", datastoreClient, memcacheClient)
	mux.Handle("/", siteMux)

	play.RegisterHandlers(mux, godevSite, chinaSite)

	mux.Handle("/explore/", http.StripPrefix("/explore/", redirectPrefix("https://pkg.go.dev/")))
	if err := blog.RegisterFeeds(mux, "", godevSite); err != nil {
		log.Fatalf("blog: %v", err)
	}

	// Note: Only golang.org/x/, no go.dev/x/.
	mux.Handle("golang.org/x/", http.HandlerFunc(xHandler))
	mux.Handle("golang.org/toolchain", http.HandlerFunc(toolchainHandler))

	redirect.Register(mux)

	// Note: Using godevSite (non-China) for global mux registration because there's no sharing in talks.
	// Don't need the hassle of two separate registrations for different domains in siteMux.
	if err := talks.RegisterHandlers(mux, godevSite, contentFS); err != nil {
		log.Fatalf("talks: %v", err)
	}
	if err := tour.RegisterHandlers(mux); err != nil {
		log.Fatalf("tour: %v", err)
	}

	var h http.Handler = mux
	h = addCSP(mux)
	h = hostEnforcerHandler(h)
	h = hostPathHandler(h)
	return h
}

var gorebuild = NewCachedURL("https://gorebuild.storage.googleapis.com/gorebuild.json", 5*time.Minute)

// newSite creates a new site for a given content and goroot file system pair
// and registers it in mux to handle requests for host.
// If host is the empty string, the registrations are for the wildcard host.
func newSite(mux *http.ServeMux, host string, content, goroot fs.FS) (*web.Site, error) {
	fsys := unionFS{content, &hideRootMDFS{&fixSpecsFS{goroot}}}
	site := web.NewSite(fsys)
	site.Funcs(template.FuncMap{
		"googleAnalytics": func() string { return googleAnalytics },
		"googleCN":        func() bool { return host == "golang.google.cn" },
		"gorebuild":       gorebuild.Get,
		"json":            jsonUnmarshal,
		"newest":          newest,
		"now":             func() time.Time { return time.Now() },
		"releases":        func() []*history.Major { return history.Majors },
		"rfc3339":         parseRFC3339,
		"section":         section,
		"version":         func() string { return runtime.Version() },
	})
	docs, err := pkgdoc.NewServer(fsys, site, googleCN)
	if err != nil {
		return nil, err
	}

	mux.Handle(host+"/", site)
	mux.Handle(host+"/cmd/", docs)
	mux.Handle(host+"/pkg/", docs)
	mux.Handle(host+"/doc/codewalk/", codewalk.NewServer(fsys, site))
	return site, nil
}

func parseRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// watchTip is a background goroutine that watches the main Go repo for updates.
// When a new commit is available, watchTip downloads the new tree and calls
// tipGoroot.Set to install the new file system.
func watchTip(tipGoroot *atomicFS) {
	for {
		// watchTip1 runs until it panics (hopefully never).
		// If that happens, sleep 5 minutes and try again.
		watchTip1(tipGoroot)
		time.Sleep(5 * time.Minute)
	}
}

// watchTip1 does the actual work of watchTip and recovers from panics.
func watchTip1(tipGoroot *atomicFS) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("watchTip panic: %v\n%s", e, debug.Stack())
		}
	}()

	var r *gitfs.Repo
	for {
		var err error
		r, err = gitfs.NewRepo("https://go.googlesource.com/go")
		if err != nil {
			log.Printf("tip: %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}
		break
	}

	var h gitfs.Hash
	for {
		var fsys fs.FS
		var err error
		h, fsys, err = r.Clone("HEAD")
		if err != nil {
			log.Printf("tip: %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}
		tipGoroot.Set(fsys)
		break
	}

	for {
		time.Sleep(5 * time.Minute)
		h2, err := r.Resolve("HEAD")
		if err != nil {
			log.Printf("tip: %v", err)
			continue
		}
		if h2 != h {
			fsys, err := r.CloneHash(h2)
			if err != nil {
				log.Printf("tip: %v", err)
				time.Sleep(1 * time.Minute)
				continue
			}
			tipGoroot.Set(fsys)
			h = h2
		}
	}
}

var datastoreClient *datastore.Client
var memcacheClient *memcache.Client

func appEngineSetup(mux *http.ServeMux) {
	googleAnalytics = os.Getenv("GOLANGORG_ANALYTICS")

	ctx := context.Background()

	var err error
	datastoreClient, err = datastore.NewClient(ctx, "")
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
	memcacheClient = memcache.New(redisAddr)

	short.RegisterHandlers(mux, "", datastoreClient, memcacheClient)

	log.Println("AppEngine initialization complete")
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
	"beta.golang.org":  true,
	"blog.golang.org":  true,
	"m.golang.org":     true,
	"talks.golang.org": true,
	"tip.golang.org":   true,
	"tour.golang.org":  true,

	"go.dev":       true,
	"blog.go.dev":  true,
	"learn.go.dev": true,
	"play.go.dev":  true,
	"talks.go.dev": true,
	"tour.go.dev":  true,
}

// hostEnforcerHandler redirects http://foo.golang.org/bar to https://golang.org/bar.
// It also forces all requests coming from China for golang.org to use golang.google.cn.
func hostEnforcerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHTTPS := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" || r.URL.Scheme == "https"
		defaultHost := "go.dev"
		host := strings.ToLower(r.Host)
		isValidHost := validHosts[host]

		if googleCN(r) && !strings.HasSuffix(host, "google.cn") {
			// golang.google.cn is the only web site in China.
			defaultHost = "golang.google.cn"
			isValidHost = strings.ToLower(r.Host) == defaultHost
		}

		if !isHTTPS || !isValidHost {
			r.URL.Scheme = "https"
			if isValidHost {
				r.URL.Host = r.Host
			} else {
				r.URL.Host = defaultHost
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
// It also rewrites the output HTML and Location headers in that case to
// link back to URLs on the test site.
func hostPathHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "localhost" && !strings.HasPrefix(r.Host, "localhost:") && !strings.HasSuffix(r.Host, ".appspot.com") {
			h.ServeHTTP(w, r)
			return
		}

		elem, rest := strings.TrimPrefix(r.URL.Path, "/"), ""
		if i := strings.Index(elem, "/"); i >= 0 {
			if elem[:i] == "tour" {
				// The Angular router serving /tour/ fails badly when it sees /go.dev/tour/.
				// Just take http://localhost/tour/ as meaning /go.dev/tour/ instead of redirecting.
				elem, rest = "go.dev", elem
			} else {
				elem, rest = elem[:i], elem[i+1:]
			}
		}
		if !validHosts[elem] {
			u := "/go.dev" + r.URL.EscapedPath()
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

		log.Print(r.URL.String())

		lw := &linkRewriter{ResponseWriter: w, host: r.Host, tour: strings.HasPrefix(r.URL.Path, "/tour/")}
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
	tour bool // is this go.dev/tour/?
	buf  []byte
	ct   string // content-type
}

func (r *linkRewriter) WriteHeader(code int) {
	loc := r.Header().Get("Location")
	delete(r.Header(), "Content-Length") // we might change the content
	if strings.HasPrefix(loc, "/") && !strings.HasPrefix(loc, "/tour/") {
		r.Header().Set("Location", "/"+r.host+loc)
	} else if u, _ := url.Parse(loc); u != nil && validHosts[u.Host] {
		r.Header().Set("Location", "/"+u.Host+"/"+strings.TrimPrefix(u.Path, "/")+u.RawQuery)
	}
	r.ResponseWriter.WriteHeader(code)
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
	var repl []string
	if !r.tour {
		repl = []string{
			`href="/`, `href="/` + r.host + `/`,
			`src="/`, `src="/` + r.host + `/`,
		}
	}
	for host := range validHosts {
		repl = append(repl, `href="https://`+host, `href="/`+host)
		repl = append(repl, `src="https://`+host, `src="/`+host)
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
		http.Redirect(w, r, "https://pkg.go.dev/golang.org/x", http.StatusTemporaryRedirect)
		return
	}
	proj, suffix := strings.TrimPrefix(r.URL.Path, "/x/"), ""
	if i := strings.Index(proj, "/"); i != -1 {
		proj, suffix = proj[:i], proj[i:]
	}
	if proj == "" {
		http.Redirect(w, r, "https://pkg.go.dev/golang.org/x", http.StatusTemporaryRedirect)
		return
	}
	repo, ok := repos.ByGerritProject[proj]
	if !ok || !strings.HasPrefix(repo.ImportPath, "golang.org/x/") {
		http.NotFound(w, r)
		return
	}
	if suffix == "/info/refs" && strings.HasPrefix(r.URL.Query().Get("service"), "git-") && repo.GoGerritProject != "" {
		// Someone is running 'git clone https://golang.org/x/repo'.
		// We want the eventual git checkout to have the right origin (go.googlesource.com)
		// and not just keep hitting golang.org all the time.
		// A redirect would work for this git command but not record the result.
		// Instead, print a useful error for the user.
		http.Error(w, fmt.Sprintf("Use 'git clone https://go.googlesource.com/%s' instead.", repo.GoGerritProject), http.StatusNotFound)
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
<html lang="en">
<title>The Go Programming Language</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="golang.org/x/{{.Proj}} git https://go.googlesource.com/{{.Proj}}">
<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/golang.org/x/{{.Proj}}{{.Suffix}}">
</head>
<body>
<a href="https://pkg.go.dev/golang.org/x/{{.Proj}}{{.Suffix}}">Redirecting to documentation...</a>
</body>
</html>
`))

func toolchainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/toolchain" {
		// Shouldn't happen if handler is registered correctly.
		http.NotFound(w, r)
		return
	}
	w.Write(toolchainPage)
}

var toolchainPage = []byte(`<!DOCTYPE html>
<html lang="en">
<head>
<title>The Go Programming Language</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="golang.org/toolchain mod https://go.dev/dl/mod">
<meta http-equiv="refresh" content="0; url=https://go.dev/dl/">
</head>
<body>
golang.org/toolchain is the module form of the Go toolchain releases.
<a href="https://go.dev/dl/">Redirecting to Go toolchain download page...</a>
</body>
</html>
`)

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

// A fixSpecsFS is an FS mapping /ref/mem.html and /ref/spec.html to
// /doc/go_mem.html and /doc/go_spec.html.
var _ fs.FS = &fixSpecsFS{}

type fixSpecsFS struct {
	fs fs.FS
}

func (fsys fixSpecsFS) Open(name string) (fs.File, error) {
	switch name {
	case "ref/mem.html", "ref/spec.html":
		if f, err := fsys.fs.Open(name); err == nil {
			// Let Go distribution win if they move.
			return f, nil
		}
		// Otherwise fall back to doc/go_*.html
		name = "doc/go_" + strings.TrimPrefix(name, "ref/")
		return fsys.fs.Open(name)

	case "doc/go_mem.html", "doc/go_spec.html":
		data := []byte("<!--{\n\t\"Redirect\": \"/ref/" + strings.TrimPrefix(strings.TrimSuffix(name, ".html"), "doc/go_") + "\"\n}-->\n")
		return &memFile{path.Base(name), bytes.NewReader(data)}, nil
	}

	return fsys.fs.Open(name)
}

// A hideRootMDFS is an FS that hides *.md files in the root directory.
// We use this to hide the Go repository's CONTRIBUTING.md,
// README.md, and SECURITY.md. The last is particularly problematic
// when running locally on a Mac, because it can be opened as
// security.md, which takes priority over _content/security.html.
type hideRootMDFS struct {
	fs fs.FS
}

func (fsys hideRootMDFS) Open(name string) (fs.File, error) {
	if !strings.Contains(name, "/") && strings.HasSuffix(name, ".md") {
		return nil, errors.New(".md file not available")
	}
	return fsys.fs.Open(name)
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

// A memFile is an fs.File implementation backed by in-memory data.
type memFile struct {
	name string
	*bytes.Reader
}

func (f *memFile) Stat() (fs.FileInfo, error) { return f, nil }
func (f *memFile) Name() string               { return f.name }
func (*memFile) Mode() fs.FileMode            { return 0444 }
func (*memFile) ModTime() time.Time           { return time.Time{} }
func (*memFile) IsDir() bool                  { return false }
func (*memFile) Sys() interface{}             { return nil }
func (*memFile) Close() error                 { return nil }

// An atomicFS is an fs.FS value safe for reading from multiple goroutines
// as well as updating (assigning a different fs.FS to use in future read requests).
type atomicFS struct {
	v atomic.Value
}

// Set sets the file system used by future calls to Open.
func (a *atomicFS) Set(fsys fs.FS) {
	a.v.Store(&fsys)
}

// Open returns fsys.Open(name) where fsys is the file system passed to the most recent call to Set.
// If there has been no call to Set, Open returns an error with text “no file system”.
func (a *atomicFS) Open(name string) (fs.File, error) {
	fsys, _ := a.v.Load().(*fs.FS)
	if fsys == nil {
		return nil, &fs.PathError{Path: name, Op: "open", Err: fmt.Errorf("no file system")}
	}
	return (*fsys).Open(name)
}

func redirectPrefix(prefix string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := strings.TrimSuffix(prefix, "/") + "/" + strings.TrimPrefix(r.URL.Path, "/")
		if r.URL.RawQuery != "" {
			url += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})
}

type CachedURL struct {
	url     string
	timeout time.Duration

	mu      sync.Mutex
	data    []byte
	err     error
	etag    string
	updated time.Time
}

func NewCachedURL(url string, timeout time.Duration) *CachedURL {
	return &CachedURL{url: url, timeout: timeout}
}

func (c *CachedURL) Get() (data []byte, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Since(c.updated) < c.timeout {
		return c.data, c.err
	}
	defer func() {
		c.updated = time.Now()
		c.data, c.err = data, err
	}()

	cli := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		return nil, err
	}
	if c.etag != "" {
		req.Header.Set("If-None-Match", c.etag)
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("loading rebuild report JSON: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		// Unmodified.
		log.Printf("checked %s - unmodified", c.url)
		return c.data, c.err
	}
	log.Printf("reloading %s", c.url)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("loading rebuild report JSON: %v", resp.Status)
	}
	c.etag = resp.Header.Get("Etag")
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("loading rebuild report JSON: %v", err)
	}
	return data, nil
}

func jsonUnmarshal(data []byte) (any, error) {
	var x any
	err := json.Unmarshal(data, &x)
	return x, err
}
