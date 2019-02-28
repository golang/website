// Copyright 2018 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine

package main

import (
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/tools/blog"
)

func TestServer(t *testing.T) {
	if runtime.GOOS == "android" {
		t.Skip("skipping on android; can't run go tool")
	}
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
