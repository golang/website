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
	"strings"
)

var (
	doctype   = []byte("<!DOCTYPE ")
	jsonStart = []byte("<!--{")
	jsonEnd   = []byte("}-->")
)

type Metadata struct {
	// Copied from document metadata
	Title    string
	Subtitle string
	Template bool

	Path     string // URL path
	FilePath string // filesystem path relative to goroot
}

type MetaJSON struct {
	Title    string
	Subtitle string
	Template bool
	Redirect string // if set, redirect to other URL
}

// extractMetadata extracts the MetaJSON from a byte slice.
// It returns the Metadata value and the remaining data.
// If no metadata is present the original byte slice is returned.
//
func extractMetadata(b []byte) (meta MetaJSON, tail []byte, err error) {
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

// MetadataFor returns the *Metadata for a given absolute path
// or nil if none exists.
func (c *Corpus) MetadataFor(path string) *Metadata {
	// Strip any .html or .md; it all names the same page.
	if strings.HasSuffix(path, ".html") {
		path = strings.TrimSuffix(path, ".html")
	} else if strings.HasSuffix(path, ".md") {
		path = strings.TrimSuffix(path, ".md")
	}

	file := path + ".html"
	b, err := fs.ReadFile(c.fs, toFS(file))
	if err != nil {
		file = path + ".md"
		b, err = fs.ReadFile(c.fs, toFS(file))
	}
	if err != nil {
		// Special case for memory model and spec, which live
		// in the main Go repo's doc directory and therefore have not
		// been renamed to their serving paths.
		// We wait until the ReadFiles above have failed so that the
		// code works if these are ever moved to /ref/spec and /ref/mem.
		switch path {
		case "/ref/spec":
			if m := c.MetadataFor("/doc/go_spec"); m != nil {
				return m
			}
		case "/ref/mem":
			if m := c.MetadataFor("/doc/go_mem"); m != nil {
				return m
			}
		}
		return nil
	}

	js, _, err := extractMetadata(b)
	if err != nil {
		log.Printf("MetadataFor %s: %v", path, err)
		return nil
	}

	meta := &Metadata{
		Title:    js.Title,
		Subtitle: js.Subtitle,
		Template: js.Template,
		Path:     path,
		FilePath: file,
	}
	if js.Redirect != "" {
		// Allow (placeholder) documents to declare a redirect.
		meta.Path = js.Redirect
	}

	// Special case for memory model and spec, continued.
	switch path {
	case "/doc/go_spec":
		meta.Path = "/ref/spec"
	case "/doc/go_mem":
		meta.Path = "/ref/mem"
	}

	return meta
}
