// +build !appengine

package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/tools/blog"
)

func TestServer(t *testing.T) {
	mux, err := newServer(false, "/static", blog.Config{
		TemplatePath: "./template",
	})
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Errorf("GET /: code = %d; want 200", w.Code)
	}
	want := "The Go Programming Language Blog"
	if !strings.Contains(w.Body.String(), want) {
		t.Errorf("GET /: want to find %q, got\n\n%q", want, w.Body.String())
	}
	if hdr := w.Header().Get("Content-Type"); hdr != "text/html; charset=utf-8" {
		t.Errorf("GET /: want text/html content-type, got %q", hdr)
	}
}
