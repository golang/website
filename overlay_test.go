// Copyright 2026 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package website

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

// The overlay is the whole fork. Everything else in this repository is
// upstream's, so this is the one behaviour with no test anywhere upstream and
// the one that every test that does exist would miss.
func TestOverlayPrefersTheOverlayAndFallsBack(t *testing.T) {
	base := fstest.MapFS{
		"only-english.md": {Data: []byte("english")},
		"both.md":         {Data: []byte("english")},
	}
	overlay := fstest.MapFS{"both.md": {Data: []byte("tiếng Việt")}}
	o := NewOverlayFS(overlay, base)

	if got := read(t, o, "both.md"); got != "tiếng Việt" {
		t.Errorf("both.md = %q, want the overlay's copy", got)
	}
	// The fallback is what makes a partial translation servable, and it is also
	// what makes a missing translation invisible on the site. Both of those
	// follow from this one line, which is why the audit exists.
	if got := read(t, o, "only-english.md"); got != "english" {
		t.Errorf("only-english.md = %q, want the base's copy", got)
	}
	if _, err := o.Open("neither.md"); err == nil {
		t.Error("a file in neither filesystem must not open")
	}
}

// GODEV_CONTENT=en turns the overlay off, and upstream's test suite depends on
// it. cmd/golangorg asserts English release note text read through Content, and
// internal/tour holds TourOnly in a package level variable and asserts
// `>A Tour of Go<`. Neither goes through the server's -content flag, so without
// the switch neither package can pass in this fork.
//
// The assertion is that the switch changes what comes back, not that a
// particular Vietnamese sentence is in a particular file. A test that pins
// prose fails the next time somebody improves the prose.
func TestGODEVCONTENTServesEnglish(t *testing.T) {
	const name = "doc/install.html"

	t.Setenv("GODEV_CONTENT", "en")
	english := read(t, Content(), name)

	// The default is Vietnamese, because a deployment that forgets to set an
	// environment variable should serve the site this fork exists to serve.
	t.Setenv("GODEV_CONTENT", "")
	translated := read(t, Content(), name)

	if english == translated {
		t.Fatalf("GODEV_CONTENT made no difference to %s, so either the switch "+
			"is not wired up or that page lost its translation", name)
	}
	if want := read(t, subdir(embedded, "_content"), name); english != want {
		t.Errorf("GODEV_CONTENT=en did not serve the embedded English copy of %s", name)
	}
}

func read(t *testing.T, fsys fs.FS, name string) string {
	t.Helper()
	b, err := fs.ReadFile(fsys, name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(b)
}
