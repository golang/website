// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tour

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/tools/playground/socket"
	"golang.org/x/website/internal/webtest"
)

const (
	socketPath = "/socket"
)

var (
	httpListen  *string
	openBrowser *bool
	httpAddr    string
)

func Main() {
	httpListen = flag.String("http", "127.0.0.1:3999", "host:port to listen on")
	openBrowser = flag.Bool("openbrowser", true, "open browser automatically")

	flag.Parse()

	host, port, err := net.SplitHostPort(*httpListen)
	if err != nil {
		log.Fatal(err)
	}
	if host == "" {
		host = "localhost"
	}
	if host != "127.0.0.1" && host != "localhost" {
		log.Print(localhostWarning)
	}
	httpAddr = host + ":" + port

	if err := initTour(http.DefaultServeMux, "SocketTransport"); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/_/fmt", fmtHandler)
	fs := http.FileServer(http.FS(contentTour))
	http.Handle("/favicon.ico", fs)
	http.Handle("/images/", fs)

	origin := &url.URL{Scheme: "http", Host: host + ":" + port}
	http.Handle(socketPath, socket.NewHandler(origin))

	h := webtest.HandlerWithCheck(http.DefaultServeMux, "/_readycheck",
		os.DirFS("."), "tour/testdata/*.txt")

	go func() {
		url := "http://" + httpAddr
		if waitServer(url) && *openBrowser && startBrowser(url) {
			log.Printf("A browser window should open. If not, please visit %s", url)
		} else {
			log.Printf("Please open your web browser and visit %s", url)
		}
	}()

	log.Fatal(http.ListenAndServe(httpAddr, &logging{h}))
}

type logging struct {
	h http.Handler
}

func (l *logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	println(r.URL.Path)
	l.h.ServeHTTP(w, r)
}

// rootHandler returns a handler for all the requests except the ones for lessons.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/tour/", http.StatusFound)
		return
	}
	if !validTourPath(r.URL.Path) {
		// Serve the tour UI, which renders its own not found page,
		// but report the status so that crawlers do not treat
		// nonexistent lessons as valid pages.
		w.WriteHeader(http.StatusNotFound)
	}
	if err := renderUI(w); err != nil {
		log.Println(err)
	}
}

// validTourPath reports whether path names a tour page:
// the tour itself, the lesson index, or a page of an existing lesson.
func validTourPath(path string) bool {
	switch path {
	case "/tour", "/tour/", "/tour/list", "/tour/list/":
		return true
	}
	rest, ok := strings.CutPrefix(path, "/tour/")
	if !ok {
		return false
	}
	name, page, _ := strings.Cut(rest, "/")
	numPages, ok := lessonPages[name]
	if !ok {
		return false
	}
	// The page number is optional, but if present it must be a number.
	page = strings.TrimSuffix(page, "/")
	if page == "" {
		return true
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 || pageNum > numPages || strconv.Itoa(pageNum) != page {
		// It must be a positive number in the range of the lesson's pages.
		return false
	}
	return true
}

// lessonHandler handler the HTTP requests for lessons.
func lessonHandler(w http.ResponseWriter, r *http.Request) {
	lesson := strings.TrimPrefix(r.URL.Path, "/tour/lesson/")
	if err := writeLesson(lesson, w); err != nil {
		if err == lessonNotFound {
			http.NotFound(w, r)
		} else {
			log.Println(err)
		}
	}
}

const localhostWarning = `
WARNING!  WARNING!  WARNING!

The tour server appears to be listening on an address that is
not localhost and is configured to run code snippets locally.
Anyone with access to this address and port will have access
to this machine as the user running gotour.

If you don't understand this message, hit Control-C to terminate this process.

WARNING!  WARNING!  WARNING!
`

// waitServer waits some time for the http Server to start
// serving url. The return value reports whether it starts.
func waitServer(url string) bool {
	tries := 20
	for tries > 0 {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
		tries--
	}
	return false
}

// startBrowser tries to open the URL in a browser, and returns
// whether it succeed.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

// prepContent for the local tour simply returns the content as-is.
var prepContent = func(r io.Reader) io.Reader { return r }

// socketAddr returns the WebSocket handler address.
var socketAddr = func() string { return "ws://" + httpAddr + socketPath }

// analyticsHTML is optional analytics HTML to insert at the beginning of <head>.
var analyticsHTML template.HTML
