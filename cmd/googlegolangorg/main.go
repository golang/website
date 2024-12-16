// A trivial redirector for google.golang.org.
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var repoMap = map[string]Repo{
	"api": {
		VCS: "git",
		URL: "https://github.com/googleapis/google-api-go-client",
	},
	"appengine": {
		VCS: "git",
		URL: "https://github.com/golang/appengine",
	},
	"cloud": {
		// This repo is now at "cloud.google.com/go", but still specifying the repo
		// here gives nicer errors in the go tool.
		VCS: "git",
		URL: "https://github.com/googleapis/google-cloud-go",
	},
	"genai": {
		VCS: "git",
		URL: "https://github.com/googleapis/go-genai",
	},
	"genproto": {
		VCS: "git",
		URL: "https://github.com/googleapis/go-genproto",
	},
	"grpc": {
		VCS: "git",
		URL: "https://github.com/grpc/grpc-go",
	},
	"protobuf": {
		VCS: "git",
		URL: "https://go.googlesource.com/protobuf",
		Src: github("protocolbuffers/protobuf-go"),
	},
	"open2opaque": {
		VCS: "git",
		URL: "https://go.googlesource.com/open2opaque",
		Src: github("golang/open2opaque"),
	},
}

// Repo represents a repository containing Go code.
type Repo struct {
	// VCS and URL set the go-import meta-tag,
	// as per https://go.dev/ref/mod#vcs-find.
	VCS string
	URL string

	// Src sets additional control over where to
	// link to for viewing source code. Optional.
	Src *src
}

// src represents a pkg.go.dev source redirect.
// See https://cs.opensource.google/go/x/pkgsite/+/master:internal/source/meta-tags.go;l=19;drc=19794c8aeb90c0a8f17c5ee1ed187bd005a1fd40?q=sourceMeta&ss=go%2Fx%2Fpkgsite.
type src struct {
	URL     string
	DirTpl  string
	FileTpl string
}

// github returns the *src representing a repo on github.com.
func github(base string) *src {
	return &src{
		URL:     fmt.Sprintf("https://github.com/%s", base),
		DirTpl:  fmt.Sprintf("https://github.com/%s/tree/master{/dir}", base),
		FileTpl: fmt.Sprintf("https://github.com/%s/tree/master{/dir}/{file}#L{line}", base),
	}
}

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Printf("Defaulting to port %s\n", port)
	}

	fmt.Printf("Listening on port %s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		fmt.Fprintf(os.Stderr, "http.ListenAndServe: %v\n", err)
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	head, tail := strings.TrimPrefix(r.URL.Path, "/"), ""
	if i := strings.Index(head, "/"); i != -1 {
		head, tail = head[:i], head[i:]
	}
	if head == "" {
		http.Redirect(w, r, "https://cloud.google.com/go/google.golang.org", http.StatusFound)
		return
	}
	repo, ok := repoMap[head]
	if !ok {
		http.NotFound(w, r)
		return
	}
	docURL := "https://pkg.go.dev/google.golang.org/" + head + tail
	// For users visiting in a browser, redirect straight to pkg.go.dev.
	if isBrowser := r.FormValue("go-get") == ""; isBrowser {
		http.Redirect(w, r, docURL, http.StatusFound)
		return
	}
	data := struct {
		Head, DocURL string
		Repo         Repo
	}{head, docURL, repo}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Fprintf(os.Stderr, "tmpl.Execute: %v\n", err)
	}
}

var tmpl = template.Must(template.New("redir").Parse(`<!DOCTYPE html>
<html>
<head>
<meta name="go-import" content="google.golang.org/{{.Head}} {{.Repo.VCS}} {{.Repo.URL}}">
{{if .Repo.Src}}
<meta name="go-source" content="google.golang.org/{{.Head}} {{.Repo.Src.URL}} {{.Repo.Src.DirTpl}} {{.Repo.Src.FileTpl}}">
{{end}}
<meta http-equiv="refresh" content="0; url={{.DocURL}}">
</head>
<body>
<a href="{{.DocURL}}">Redirecting to documentation...</a>
</body>
</html>
`))
