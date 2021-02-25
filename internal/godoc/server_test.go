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
		"doc/x/index.html": {Data: []byte("Hello, x.")},
	})
	p := &Presentation{
		Corpus:    c,
		GodocHTML: template.Must(template.New("").Parse(`{{printf "%s" .Body}}`)),
	}

	// Test that redirect is sent back correctly.
	// Used to panic. See golang.org/issue/40665.
	dir := "/doc/x/"

	r := &http.Request{URL: &url.URL{Path: dir + "index.html"}}
	rw := httptest.NewRecorder()
	p.ServeFile(rw, r)
	loc := rw.Result().Header.Get("Location")
	if rw.Code != 301 || loc != dir {
		t.Errorf("GET %s: expected 301 -> %q, got %d -> %q", r.URL.Path, dir, rw.Code, loc)
	}

	testServeBody(t, p, dir, "Hello, x")
}

func TestMarkdown(t *testing.T) {
	p := &Presentation{
		Corpus: NewCorpus(fstest.MapFS{
			"doc/test.md":  {Data: []byte("**bold**")},
			"doc/test2.md": {Data: []byte(`{{"*template*"}}`)},
		}),
		GodocHTML: template.Must(template.New("").Parse(`{{printf "%s" .Body}}`)),
	}

	testServeBody(t, p, "/doc/test", "<strong>bold</strong>")
	testServeBody(t, p, "/doc/test2", "<em>template</em>")
}
