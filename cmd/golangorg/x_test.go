// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestXHandler(t *testing.T) {
	type check func(t *testing.T, rec *httptest.ResponseRecorder)
	status := func(v int) check {
		return func(t *testing.T, rec *httptest.ResponseRecorder) {
			t.Helper()
			if rec.Code != v {
				t.Errorf("response status = %v; want %v", rec.Code, v)
			}
		}
	}
	substr := func(s string) check {
		return func(t *testing.T, rec *httptest.ResponseRecorder) {
			t.Helper()
			if !strings.Contains(rec.Body.String(), s) {
				t.Errorf("missing expected substring %q in value: %#q", s, rec.Body)
			}
		}
	}
	hasHeader := func(k, v string) check {
		return func(t *testing.T, rec *httptest.ResponseRecorder) {
			t.Helper()
			if got := rec.HeaderMap.Get(k); got != v {
				t.Errorf("header[%q] = %q; want %q", k, got, v)
			}
		}
	}

	tests := []struct {
		name   string
		path   string
		checks []check
	}{
		{
			name: "net",
			path: "/x/net",
			checks: []check{
				status(200),
				substr(`<meta name="go-import" content="golang.org/x/net git https://go.googlesource.com/net">`),
				substr(`http-equiv="refresh" content="0; url=https://pkg.go.dev/golang.org/x/net">`),
			},
		},
		{
			name: "net-suffix",
			path: "/x/net/suffix",
			checks: []check{
				status(200),
				substr(`<meta name="go-import" content="golang.org/x/net git https://go.googlesource.com/net">`),
				substr(`http-equiv="refresh" content="0; url=https://pkg.go.dev/golang.org/x/net/suffix">`),
			},
		},
		{
			name: "pkgsite",
			path: "/x/pkgsite",
			checks: []check{
				status(200),
				substr(`<meta name="go-import" content="golang.org/x/pkgsite git https://go.googlesource.com/pkgsite">`),
				substr(`Nothing to see here; <a href="https://pkg.go.dev/golang.org/x/pkgsite">move along</a>.`),
				substr(`http-equiv="refresh" content="0; url=https://pkg.go.dev/golang.org/x/pkgsite">`),
			},
		},
		{
			name:   "notexist",
			path:   "/x/notexist",
			checks: []check{status(404)},
		},
		{
			name: "empty",
			path: "/x/",
			checks: []check{
				status(307),
				hasHeader("Location", "https://pkg.go.dev/search?q=golang.org/x"),
			},
		},
		{
			name:   "invalid",
			path:   "/x/In%20Valid,X",
			checks: []check{status(404)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rec := httptest.NewRecorder()
			xHandler(rec, req)
			for _, check := range tt.checks {
				check(t, rec)
			}
		})
	}
}
