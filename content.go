// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

// Package website exports the static content as an embed.FS.
package website

import (
	"embed"
	"io/fs"
)

// Content is the website's static content.
var Content = subdir(embedded, "_content")

//go:embed _content
var embedded embed.FS

func subdir(fsys fs.FS, path string) fs.FS {
	s, err := fs.Sub(fsys, path)
	if err != nil {
		panic(err)
	}
	return s
}
