// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"log"
	"path"
	"strings"
)

var (
	doctype   = []byte("<!DOCTYPE ")
	jsonStart = []byte("<!--{")
	jsonEnd   = []byte("}-->")
)

type file struct {
	// Copied from document metadata directives
	Title    string
	Subtitle string
	Template bool

	Path     string // URL path
	FilePath string // filesystem path relative to goroot
	Body     []byte // content after metadata
}

type metaJSON struct {
	Title    string
	Subtitle string
	Template bool
	Redirect string // if set, redirect to other URL
}

// extractMetadata extracts the metaJSON from a byte slice.
// It returns the metadata and the remaining text.
// If no metadata is present, it returns an empty metaJSON and the original text.
func extractMetadata(b []byte) (meta metaJSON, tail []byte, err error) {
	tail = b
	if !bytes.HasPrefix(b, jsonStart) {
		return
	}
	end := bytes.Index(b, jsonEnd)
	if end < 0 {
		return
	}
	b = b[len(jsonStart)-1 : end+1] // drop leading <!-- and include trailing }
	if err = json.Unmarshal(b, &meta); err != nil {
		return
	}
	tail = tail[end+len(jsonEnd):]
	return
}

var join = path.Join

// open returns the file for a given absolute path or nil if none exists.
func open(fsys fs.FS, path string) *file {
	// Strip trailing .html or .md or /; it all names the same page.
	if strings.HasSuffix(path, ".html") {
		path = strings.TrimSuffix(path, ".html")
	} else if strings.HasSuffix(path, ".md") {
		path = strings.TrimSuffix(path, ".md")
	} else if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	files := []string{path + ".html", path + ".md", join(path, "index.html"), join(path, "index.md")}
	var filePath string
	var b []byte
	var err error
	for _, filePath = range files {
		b, err = fs.ReadFile(fsys, toFS(filePath))
		if err == nil {
			break
		}
	}

	// Special case for memory model and spec, which live
	// in the main Go repo's doc directory and therefore have not
	// been renamed to their serving paths.
	// We wait until the ReadFiles above have failed so that the
	// code works if these are ever moved to /ref/spec and /ref/mem.
	if err != nil && path == "/ref/spec" {
		return open(fsys, "/doc/go_spec")
	}
	if err != nil && path == "/ref/mem" {
		return open(fsys, "/doc/go_mem")
	}

	if err != nil {
		return nil
	}

	// Special case for memory model and spec, continued.
	switch path {
	case "/doc/go_spec":
		path = "/ref/spec"
	case "/doc/go_mem":
		path = "/ref/mem"
	}

	// If we read an index.md or index.html, the canonical path is without the index.md/index.html suffix.
	if strings.HasSuffix(filePath, "/index.md") || strings.HasSuffix(filePath, "/index.html") {
		path = filePath[:strings.LastIndex(filePath, "/")+1]
	}

	js, body, err := extractMetadata(b)
	if err != nil {
		log.Printf("extractMetadata %s: %v", path, err)
		return nil
	}

	f := &file{
		Title:    js.Title,
		Subtitle: js.Subtitle,
		Template: js.Template,
		Path:     path,
		FilePath: filePath,
		Body:     body,
	}
	if js.Redirect != "" {
		// Allow (placeholder) documents to declare a redirect.
		f.Path = js.Redirect
	}

	return f
}
