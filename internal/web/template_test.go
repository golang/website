// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"testing"

	"golang.org/x/website/internal/backport/html/template"
)

func TestSrcToPkg(t *testing.T) {
	for _, tc := range []struct {
		path string
		want string
	}{
		{"/src/fmt", "pkg/fmt"},
		{"src/fmt", "pkg/fmt"},
		{"/fmt", "pkg/fmt"},
		{"fmt", "pkg/fmt"},
		{"src/pkg/fmt", "pkg/fmt"},
		{"/src/pkg/fmt", "pkg/fmt"},
	} {
		if got := srcToPkg(tc.path); got != tc.want {
			t.Errorf("srcToPkg(%v) = %v; want %v", tc.path, got, tc.want)
		}
	}
}

func TestSrcBreadcrumbFunc(t *testing.T) {
	for _, tc := range []struct {
		path string
		want template.HTML
	}{
		{"src/", `<span class="text-muted">src/</span>`},
		{"src/fmt/", `<a href="/src">src</a>/<span class="text-muted">fmt/</span>`},
		{"src/fmt/print.go", `<a href="/src">src</a>/<a href="/src/fmt">fmt</a>/<span class="text-muted">print.go</span>`},
	} {
		if got := (&Page{SrcPath: tc.path}).SrcBreadcrumb(); got != tc.want {
			t.Errorf("srcBreadcrumbFunc(%v) = %v; want %v", tc.path, got, tc.want)
		}
	}
}

func TestSrcPkgLink(t *testing.T) {
	for _, tc := range []struct {
		path string
		want template.HTML
	}{
		{"src/", `<a href="/pkg">Index</a>`},
		{"src/fmt/", `<a href="/pkg/fmt">fmt</a>`},
		{"pkg/", `<a href="/pkg">Index</a>`},
		{"pkg/LICENSE", `<a href="/pkg">Index</a>`},
	} {
		if got := (&Page{SrcPath: tc.path}).SrcPkgLink(); got != tc.want {
			t.Errorf("srcToPkgLinkFunc(%v) = %v; want %v", tc.path, got, tc.want)
		}
	}
}
