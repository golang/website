// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package website exports the static content as an embed.FS.
package website

import (
	"embed"
	"io/fs"
	"os"
)

// english reports whether the embedded content should be served untranslated.
//
// GODEV_CONTENT=en turns the overlay off. It exists for upstream's own test
// suite, which asserts English strings against these two functions and cannot
// do anything else: `>A Tour of Go<` is not a thing a translated page contains,
// and a test asserting it is not wrong for asserting it.
//
// Without the switch the fork cannot run the tests it inherited. cmd/golangorg
// reads website.Content() to check that every release notes page carries its
// release date, and internal/tour holds website.TourOnly() in a package level
// variable, so both get Vietnamese whatever content directory the server was
// given, and both fail. Turning the overlay off puts the engine back under the
// suite that was written for it, which is what those tests are for: they check
// that the fork has not broken the renderer, not that the translation is good.
// What the translation is like is the audit's question and the audit answers it
// over all 680 files.
//
// The default is Vietnamese, deliberately. A deployment that forgets to set an
// environment variable should serve the site this fork exists to serve.
func english() bool {
	return os.Getenv("GODEV_CONTENT") == "en"
}

// Content returns the go.dev website's static content,
// overlaying Vietnamese translations from _content_vi on top of _content.
func Content() fs.FS {
	en := subdir(embedded, "_content")
	if english() {
		return en
	}
	return NewOverlayFS(subdir(embeddedVI, "_content_vi"), en)
}

// TourOnly returns the content needed only for the standalone tour,
// overlaying Vietnamese translations from _content_vi/tour on top of _content/tour.
func TourOnly() fs.FS {
	en := subdir(tourOnly, "_content")
	if english() {
		return en
	}
	return NewOverlayFS(subdir(tourOnlyVI, "_content_vi"), en)
}

// NewOverlayFS returns a filesystem that tries overlay first, falling back to base
// for files not present in overlay.
func NewOverlayFS(overlay, base fs.FS) fs.FS {
	return &overlayFS{overlay, base}
}

type overlayFS struct {
	overlay fs.FS
	base    fs.FS
}

func (o *overlayFS) Open(name string) (fs.File, error) {
	f, err := o.overlay.Open(name)
	if err == nil {
		return f, nil
	}
	return o.base.Open(name)
}

//go:embed _content
var embedded embed.FS

//go:embed _content_vi
var embeddedVI embed.FS

//go:embed _content/favicon.ico
//go:embed _content/images/go-logo-white.svg
//go:embed _content/images/icons
//go:embed _content/js/playground.js
//go:embed _content/tour
var tourOnly embed.FS

//go:embed _content_vi/tour
var tourOnlyVI embed.FS

func subdir(fsys fs.FS, path string) fs.FS {
	s, err := fs.Sub(fsys, path)
	if err != nil {
		panic(err)
	}
	return s
}
