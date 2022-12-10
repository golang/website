// A trivial redirector for google.golang.org.
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var repoMap = map[string]*repoImport{
	"api": {
		VCS: "git",
		URL: "https://github.com/googleapis/google-api-go-client",
		Src: github("googleapis/google-api-go-client"),
	},
	"appengine": {
		VCS: "git",
		URL: "https://github.com/golang/appengine",
		Src: github("golang/appengine"),
	},
	"cloud": {
		// This repo is now at "cloud.google.com/go", but still specifying the repo
		// here gives nicer errors in the go tool.
		VCS: "git",
		URL: "https://github.com/googleapis/google-cloud-go",
		Src: github("googleapis/google-cloud-go"),
	},
	"genproto": {
		VCS: "git",
		URL: "https://github.com/googleapis/go-genproto",
		Src: github("googleapis/go-genproto"),
	},
	"grpc": {
		VCS: "git",
		URL: "https://github.com/grpc/grpc-go",
		Src: github("grpc/grpc-go"),
	},
	"protobuf": {
		VCS: "git",
		URL: "https://go.googlesource.com/protobuf",
		Src: github("protocolbuffers/protobuf-go"),
	},
}

// repoImport represents an import meta-tag, as per
// https://golang.org/cmd/go/#hdr-Import_path_syntax
type repoImport struct {
	VCS string
	URL string
	Src *src
}

// src represents a pkg.go.dev source redirect.
// https://github.com/golang/gddo/search?utf8=%E2%9C%93&q=sourceMeta
type src struct {
	URL     string
	DirTpl  string
	FileTpl string
}

// github returns the *src representing a github repo.
func github(base string) *src {
	return &src{
		URL:     fmt.Sprintf("https://github.com/%s", base),
		DirTpl:  fmt.Sprintf("https://github.com/%s/tree/master{/dir}", base),
		FileTpl: fmt.Sprintf("https://github.com/%s/tree/master{/dir}/{file}#L{line}", base),
	}
}

func googsource(repo, base string) *src {
	return &src{
		URL:     fmt.Sprintf("https://%s.googlesource.com/%s", repo, base),
		DirTpl:  fmt.Sprintf("https://%s.googlesource.com/%s/+/master{/dir}", repo, base),
		FileTpl: fmt.Sprintf("https://%s.googlesource.com/%s/+/master{/dir}/{file}#{line}", repo, base),
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
	godoc := "https://pkg.go.dev/google.golang.org/" + head + tail
	// For users visiting in a browser, redirect straight to pkg.go.dev.
	if isBrowser := r.FormValue("go-get") == ""; isBrowser {
		http.Redirect(w, r, godoc, http.StatusFound)
		return
	}
	data := struct {
		Head, GoDoc string
		Repo        *repoImport
	}{head, godoc, repo}
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
<meta http-equiv="refresh" content="0; url={{.GoDoc}}">
</head>
<body>
Nothing to see here. Please <a href="{{.GoDoc}}">move along</a>.
</body>
</html>
`))
