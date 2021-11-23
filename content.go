// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package website exports the static content as an embed.FS.
package website

import (
	"embed"
	"io/fs"
)

// Content returns the go.dev website's static content.
func Content() fs.FS {
	return subdir(embedded, "_content")
}

// TourOnly returns the content needed only for the standalone tour.
func TourOnly() fs.FS {
	return subdir(tourOnly, "_content")
}

//go:embed _content
var embedded embed.FS

//go:embed _content/favicon.ico
//go:embed _content/images/go-logo-white.svg
//go:embed _content/js/playground.js
//go:embed _content/tour
var tourOnly embed.FS

func subdir(fsys fs.FS, path string) fs.FS {
	s, err := fs.Sub(fsys, path)
	if err != nil {
		panic(err)
	}
	return s
}
