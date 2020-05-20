// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the handlers that serve go-import redirects for Go
// sub-repositories. It specifies the mapping from import paths like
// "golang.org/x/tools" to the actual repository locations.

package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"golang.org/x/build/repos"
)

const xPrefix = "/x/"

func init() {
	http.HandleFunc(xPrefix, xHandler)
}

func xHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, xPrefix) {
		// Shouldn't happen if handler is registered correctly.
		http.Redirect(w, r, "https://godoc.org/-/subrepo", http.StatusTemporaryRedirect)
		return
	}
	proj, suffix := strings.TrimPrefix(r.URL.Path, xPrefix), ""
	if i := strings.Index(proj, "/"); i != -1 {
		proj, suffix = proj[:i], proj[i:]
	}
	if proj == "" {
		http.Redirect(w, r, "https://godoc.org/-/subrepo", http.StatusTemporaryRedirect)
		return
	}
	repo, ok := repos.ByGerritProject[proj]
	if !ok || !strings.HasPrefix(repo.ImportPath, "golang.org/x/") {
		http.NotFound(w, r)
		return
	}
	docSite := "godoc.org"
	if repo.UsePkgGoDev() {
		docSite = "pkg.go.dev"
	}
	data := struct {
		DocSite string // Website providing documentation, either "godoc.org" or "pkg.go.dev".
		Proj    string // Gerrit project ("net", "sys", etc)
		Suffix  string // optional "/path" for requests like /x/PROJ/path
	}{docSite, proj, suffix}
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
<meta http-equiv="refresh" content="0; url=https://{{.DocSite}}/golang.org/x/{{.Proj}}{{.Suffix}}">
</head>
<body>
Nothing to see here; <a href="https://{{.DocSite}}/golang.org/x/{{.Proj}}{{.Suffix}}">move along</a>.
</body>
</html>
`))
