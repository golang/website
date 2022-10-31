// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"syscall"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
)

func testServeBody(t *testing.T, p *Site, path, body string) {
	t.Helper()
	r := &http.Request{URL: &url.URL{Path: path}}
	rw := httptest.NewRecorder()
	p.ServeHTTP(rw, r)
	if rw.Code != 200 || !strings.Contains(rw.Body.String(), body) {
		t.Fatalf("GET %s: expected 200 w/ %q: got %d w/ body:\n%s",
			path, body, rw.Code, rw.Body)
	}
}

func TestRedirectAndMetadata(t *testing.T) {
	fsys := fstest.MapFS{
		"site.tmpl":           {Data: []byte(`{{.Content}}`)},
		"doc/x/index.html":    {Data: []byte("Hello, x.")},
		"lib/godoc/site.html": {Data: []byte(`{{.Data}}`)},
	}
	site := NewSite(fsys)

	// Test that redirect is sent back correctly.
	// Used to panic. See golang.org/issue/40665.
	dir := "/doc/x/"

	r := &http.Request{URL: &url.URL{Path: dir + "index.html"}}
	rw := httptest.NewRecorder()
	site.ServeHTTP(rw, r)
	loc := rw.Result().Header.Get("Location")
	if rw.Code != 301 || loc != dir {
		t.Errorf("GET %s: expected 301 -> %q, got %d -> %q", r.URL.Path, dir, rw.Code, loc)
	}

	testServeBody(t, site, dir, "Hello, x")
}

func TestMarkdown(t *testing.T) {
	site := NewSite(fstest.MapFS{
		"site.tmpl":           {Data: []byte(`{{.Content}}`)},
		"doc/test.md":         {Data: []byte("**bold**")},
		"doc/test2.md":        {Data: []byte(`{{"*template*"}}`)},
		"lib/godoc/site.html": {Data: []byte(`{{.Data}}`)},
	})

	testServeBody(t, site, "/doc/test", "<strong>bold</strong>")
	testServeBody(t, site, "/doc/test2", "<em>template</em>")
}

func TestTypeScript(t *testing.T) {
	exampleOut, err := os.ReadFile("testdata/example.js")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name            string
		path            string
		wantCode        int
		wantBody        string
		wantCacheHeader bool
	}{
		{
			name:     "example code",
			path:     "/example.ts",
			wantCode: 200,
			wantBody: string(exampleOut),
		},
		{
			name:            "example code cached",
			path:            "/example.ts",
			wantCode:        200,
			wantBody:        string(exampleOut),
			wantCacheHeader: true,
		},
		{
			name:     "file not found",
			path:     "/notfound.ts",
			wantCode: 500,
			wantBody: syscall.ENOENT.Error(),
		},
		{
			name:     "syntax error",
			path:     "/error.ts",
			wantCode: 500,
			wantBody: "\n\nExpected identifier but found &#34;function&#34;\n\n",
		},
	}
	fsys := os.DirFS("testdata")
	site := NewSite(fsys)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			got := httptest.NewRecorder()
			site.serveTypeScript(got, req)
			gotHeader := got.Header().Get(cacheHeader) == "true"
			if got.Code != tt.wantCode {
				t.Errorf("got status %d but wanted %d", got.Code, http.StatusOK)
			}
			if (tt.wantCacheHeader && !gotHeader) || (!tt.wantCacheHeader && gotHeader) {
				t.Errorf("got cache hit %v but wanted %v", gotHeader, tt.wantCacheHeader)
			}
			if !strings.Contains(got.Body.String(), tt.wantBody) {
				t.Errorf("ServeHTTP() mismatch (-want +got):\n%s", cmp.Diff(tt.wantBody, got.Body.String()))
			}
		})
	}
}
