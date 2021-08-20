// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package website exports the static content as an embed.FS.
package website

import (
	"io/fs"
	"os"
)

// Content is the website's static content.
var Content = findContent()

// TODO: Use with Go 1.16 in place of findContent call above.
// var Content = subdir(embedded, "_content")
// //go:embed _content
// var embedded embed.FS

func findContent() fs.FS {
	// Walk parent directories looking for _content.
	dir := "_content"
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(dir + "/lib/godoc/godocs.js"); err == nil {
			return os.DirFS(dir)
		}
		dir = "../" + dir
	}
	panic("cannot find _content")
}

func subdir(fsys fs.FS, path string) fs.FS {
	s, err := fs.Sub(fsys, path)
	if err != nil {
		panic(err)
	}
	return s
}
