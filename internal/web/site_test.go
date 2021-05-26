// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/website/internal/backport/testing/fstest"
)

func testServeBody(t *testing.T, p *Site, path, body string) {
	t.Helper()
	r := &http.Request{URL: &url.URL{Path: path}}
	rw := httptest.NewRecorder()
	p.serveFile(rw, r)
	if rw.Code != 200 || !strings.Contains(rw.Body.String(), body) {
		t.Fatalf("GET %s: expected 200 w/ %q: got %d w/ body:\n%s",
			path, body, rw.Code, rw.Body)
	}
}

func TestRedirectAndMetadata(t *testing.T) {
	fsys := fstest.MapFS{
		"doc/x/index.html":    {Data: []byte("Hello, x.")},
		"lib/godoc/site.html": {Data: []byte(`{{.Data}}`)},
	}
	p, err := NewSite(fsys)
	if err != nil {
		t.Fatal(err)
	}

	// Test that redirect is sent back correctly.
	// Used to panic. See golang.org/issue/40665.
	dir := "/doc/x/"

	r := &http.Request{URL: &url.URL{Path: dir + "index.html"}}
	rw := httptest.NewRecorder()
	p.serveFile(rw, r)
	loc := rw.Result().Header.Get("Location")
	if rw.Code != 301 || loc != dir {
		t.Errorf("GET %s: expected 301 -> %q, got %d -> %q", r.URL.Path, dir, rw.Code, loc)
	}

	testServeBody(t, p, dir, "Hello, x")
}

func TestMarkdown(t *testing.T) {
	p, err := NewSite(fstest.MapFS{
		"doc/test.md":         {Data: []byte("**bold**")},
		"doc/test2.md":        {Data: []byte(`{{"*template*"}}`)},
		"lib/godoc/site.html": {Data: []byte(`{{.Data}}`)},
	})
	if err != nil {
		t.Fatal(err)
	}

	testServeBody(t, p, "/doc/test", "<strong>bold</strong>")
	testServeBody(t, p, "/doc/test2", "<em>template</em>")
}
