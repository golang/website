// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"bytes"
	"encoding/json"
	"log"
	"path"
	"strings"

	"golang.org/x/website/internal/backport/io/fs"
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

type fileJSON struct {
	Title    string
	Subtitle string
	Template bool
	Redirect string // if set, redirect to other URL
}

// open returns the *file for a given relative path or nil if none exists.
func open(fsys fs.FS, relpath string) *file {
	// Strip trailing .html or .md or /; it all names the same page.
	if strings.HasSuffix(relpath, ".html") {
		relpath = strings.TrimSuffix(relpath, ".html")
	} else if strings.HasSuffix(relpath, ".md") {
		relpath = strings.TrimSuffix(relpath, ".md")
	} else if strings.HasSuffix(relpath, "/") {
		relpath = strings.TrimSuffix(relpath, "/")
	}

	// Check md before html to work correctly when x/website is layered atop Go 1.15 goroot during Go 1.15 tests.
	// Want to find x/website's debugging_with_gdb.md not Go 1.15's debuging_with_gdb.html.
	files := []string{relpath + ".md", relpath + ".html", path.Join(relpath, "index.md"), path.Join(relpath, "index.html")}
	var filePath string
	var b []byte
	var err error
	for _, filePath = range files {
		b, err = fs.ReadFile(fsys, filePath)
		if err == nil {
			break
		}
	}

	// Special case for memory model and spec, which live
	// in the main Go repo's doc directory and therefore have not
	// been renamed to their serving relpaths.
	// We wait until the ReadFiles above have failed so that the
	// code works if these are ever moved to /ref/spec and /ref/mem.
	if err != nil && relpath == "ref/spec" {
		return open(fsys, "doc/go_spec")
	}
	if err != nil && relpath == "ref/mem" {
		return open(fsys, "doc/go_mem")
	}

	if err != nil {
		return nil
	}

	// Special case for memory model and spec, continued.
	switch relpath {
	case "doc/go_spec":
		relpath = "ref/spec"
	case "doc/go_mem":
		relpath = "ref/mem"
	}

	// If we read an index.md or index.html, the canonical relpath is without the index.md/index.html suffix.
	if name := path.Base(filePath); name == "index.html" || name == "index.md" {
		relpath, _ = path.Split(filePath)
	}

	js, body, err := parseFile(b)
	if err != nil {
		log.Printf("extractMetadata %s: %v", relpath, err)
		return nil
	}

	f := &file{
		Title:    js.Title,
		Subtitle: js.Subtitle,
		Template: js.Template,
		Path:     "/" + relpath,
		FilePath: filePath,
		Body:     body,
	}
	if js.Redirect != "" {
		// Allow (placeholder) documents to declare a redirect.
		f.Path = js.Redirect
	}

	return f
}

var (
	jsonStart = []byte("<!--{")
	jsonEnd   = []byte("}-->")
)

// parseFile extracts the metaJSON from a byte slice.
// It returns the metadata and the remaining text.
// If no metadata is present, it returns an empty metaJSON and the original text.
func parseFile(b []byte) (meta fileJSON, tail []byte, err error) {
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
