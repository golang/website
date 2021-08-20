// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package website exports the static content as an embed.FS.
package website

import (
	"embed"
	"io/fs"
)

// Golang is the golang.org website's static content.
var Golang fs.FS = subdir(embedded, "_content")

// Godev is the go.dev website's static content.
var Godev fs.FS = subdir(embedded, "go.dev/_content")

//go:embed _content go.dev/_content
var embedded embed.FS

func subdir(fsys fs.FS, path string) fs.FS {
	s, err := fs.Sub(fsys, path)
	if err != nil {
		panic(err)
	}
	return s
}
