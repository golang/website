// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"testing/fstest"
	"text/template"
)

// TestIgnoredGoFiles tests the scenario where a folder has no .go or .c files,
// but has an ignored go file.
func TestIgnoredGoFiles(t *testing.T) {
	packagePath := "github.com/package"
	packageComment := "main is documented in an ignored .go file"

	fs := fstest.MapFS{
		"src/" + packagePath + "/ignored.go": {Data: []byte(`// +build ignore

// ` + packageComment + `
package main`)},
	}
	d := NewDocTree(fs)
	pInfo := d.GetPageInfo("/src/"+packagePath, packagePath, NoFiltering, "linux", "amd64")

	if pInfo.PDoc == nil {
		t.Error("pInfo.PDoc = nil; want non-nil.")
	} else {
		if got, want := pInfo.PDoc.Doc, packageComment+"\n"; got != want {
			t.Errorf("pInfo.PDoc.Doc = %q; want %q.", got, want)
		}
		if got, want := pInfo.PDoc.Name, "main"; got != want {
			t.Errorf("pInfo.PDoc.Name = %q; want %q.", got, want)
		}
		if got, want := pInfo.PDoc.ImportPath, packagePath; got != want {
			t.Errorf("pInfo.PDoc.ImportPath = %q; want %q.", got, want)
		}
	}
	if pInfo.FSet == nil {
		t.Error("pInfo.FSet = nil; want non-nil.")
	}
}

func TestIssue5247(t *testing.T) {
	const packagePath = "example.com/p"
	fs := fstest.MapFS{
		"src/" + packagePath + "/p.go": {Data: []byte(`package p

//line notgen.go:3
// F doc //line 1 should appear
// line 2 should appear
func F()
//line foo.go:100`)}, // No newline at end to check corner cases.
	}

	d := NewDocTree(fs)
	pInfo := d.GetPageInfo("/src/"+packagePath, packagePath, 0, "linux", "amd64")
	if got, want := pInfo.PDoc.Funcs[0].Doc, "F doc //line 1 should appear\nline 2 should appear\n"; got != want {
		t.Errorf("pInfo.PDoc.Funcs[0].Doc = %q; want %q", got, want)
	}
}

func testServeBody(t *testing.T, p *Presentation, path, body string) {
	t.Helper()
	r := &http.Request{URL: &url.URL{Path: path}}
	rw := httptest.NewRecorder()
	p.ServeFile(rw, r)
	if rw.Code != 200 || !strings.Contains(rw.Body.String(), body) {
		t.Fatalf("GET %s: expected 200 w/ %q: got %d w/ body:\n%s",
			path, body, rw.Code, rw.Body)
	}
}

func TestRedirectAndMetadata(t *testing.T) {
	c := NewCorpus(fstest.MapFS{
		"doc/y/index.html": {Data: []byte("Hello, y.")},
		"doc/x/index.html": {Data: []byte(`<!--{
		"Path": "/doc/x/"
}-->

Hello, x.
`)},
	})
	c.updateMetadata()
	p := &Presentation{
		Corpus:    c,
		GodocHTML: template.Must(template.New("").Parse(`{{printf "%s" .Body}}`)),
	}

	// Test that redirect is sent back correctly.
	// Used to panic. See golang.org/issue/40665.
	for _, elem := range []string{"x", "y"} {
		dir := "/doc/" + elem + "/"

		r := &http.Request{URL: &url.URL{Path: dir + "index.html"}}
		rw := httptest.NewRecorder()
		p.ServeFile(rw, r)
		loc := rw.Result().Header.Get("Location")
		if rw.Code != 301 || loc != dir {
			t.Errorf("GET %s: expected 301 -> %q, got %d -> %q", r.URL.Path, dir, rw.Code, loc)
		}

		testServeBody(t, p, dir, "Hello, "+elem)
	}
}

func TestMarkdown(t *testing.T) {
	p := &Presentation{
		Corpus: NewCorpus(fstest.MapFS{
			"doc/test.md":  {Data: []byte("**bold**")},
			"doc/test2.md": {Data: []byte(`{{"*template*"}}`)},
		}),
		GodocHTML: template.Must(template.New("").Parse(`{{printf "%s" .Body}}`)),
	}

	testServeBody(t, p, "/doc/test.html", "<strong>bold</strong>")
	testServeBody(t, p, "/doc/test2.html", "<em>template</em>")
}
