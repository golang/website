// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package esbuild transforms TypeScript code into
// JavaScript code.
package esbuild

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"golang.org/x/website/internal/web"
)

const cacheHeader = "X-Go-Dev-Cache-Hit"

type server struct {
	fsys  fs.FS
	site  *web.Site
	cache sync.Map // TypeScript filepath -> JavaScript output
}

// NewServer returns a new server for handling TypeScript files.
func NewServer(fsys fs.FS, site *web.Site) http.Handler {
	return &server{fsys, site, sync.Map{}}
}

type JSOut struct {
	output []byte
	stat   fs.FileInfo // stat for file when page was loaded
}

// Handler for TypeScript files. Transforms TypeScript code into
// JavaScript code before serving them.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := path.Clean(r.URL.Path)[1:]
	if cjs, ok := s.cache.Load(filename); ok {
		js := cjs.(*JSOut)
		info, err := fs.Stat(s.fsys, filename)
		if err == nil && info.ModTime().Equal(js.stat.ModTime()) {
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			w.Header().Set(cacheHeader, "true")
			http.ServeContent(w, r, filename, info.ModTime(), bytes.NewReader(js.output))
			return
		}
	}
	file, err := s.fsys.Open(filename)
	if err != nil {
		s.site.ServeError(w, r, err)
		return
	}
	var contents bytes.Buffer
	_, err = io.Copy(&contents, file)
	if err != nil {
		s.site.ServeError(w, r, err)
		return
	}
	result := api.Transform(contents.String(), api.TransformOptions{
		Loader: api.LoaderTS,
	})
	var buf bytes.Buffer
	for _, v := range result.Errors {
		fmt.Fprintln(&buf, v.Text)
	}
	if buf.Len() > 0 {
		s.site.ServeError(w, r, errors.New(buf.String()))
		return
	}
	info, err := file.Stat()
	if err != nil {
		s.site.ServeError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	http.ServeContent(w, r, filename, info.ModTime(), bytes.NewReader(result.Code))
	s.cache.Store(filename, &JSOut{
		output: result.Code,
		stat:   info,
	})
}
