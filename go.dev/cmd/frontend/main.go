package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

var discoveryHosts = map[string]string{
	"":               "pkg.go.dev",
	"dev.go.dev":     "dev-pkg.go.dev",
	"staging.go.dev": "staging-pkg.go.dev",
}

func main() {
	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/", fs)
	http.Handle("/explore/", http.StripPrefix("/explore/", redirectHosts(discoveryHosts)))
	http.Handle("/learn/", http.StripPrefix("/learn/", redirectHosts(map[string]string{"": "learn.go.dev"})))

	p := listenPort()
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("net.Listen(%q, %q) = _, %v", "tcp", p, err)
	}
	defer l.Close()
	log.Printf("Listening on http://%v/\n", l.Addr().String())
	log.Print(http.Serve(l, nil))
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
