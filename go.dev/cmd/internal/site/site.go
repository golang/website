// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package site implements generation of content for serving from go.dev.
// It is meant to support a transition from being a Hugo-based web site
// to being a site compatible with x/website.
package site

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/website/internal/backport/html/template"
	"gopkg.in/yaml.v3"
)

// A Site holds metadata about the entire site.
type Site struct {
	URL   string
	Title string

	pagesByID map[string]*page
	dir       string
	base      *template.Template
}

// Load loads and returns the site in the directory rooted at dir.
func Load(dir string) (*Site, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	site := &Site{
		dir:       dir,
		pagesByID: make(map[string]*page),
	}
	if err := site.initTemplate(); err != nil {
		return nil, err
	}

	// Read site config.
	data, err := ioutil.ReadFile(site.file("_content/site.yaml"))
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &site); err != nil {
		return nil, fmt.Errorf("parsing _content/site.yaml: %v", err)
	}

	// Load site pages from md files.
	err = filepath.Walk(site.file("_content"), func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(name, ".md") {
			_, err := site.loadPage(name[len(site.file("."))+1:])
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("loading pages: %v", err)
	}

	// Now that all pages are loaded and set up, can render all.
	// (Pages can refer to other pages.)
	for _, p := range site.pagesByID {
		if err := site.renderHTML(p); err != nil {
			return nil, err
		}
	}

	return site, nil
}

// file returns the full path to the named file within the site.
func (site *Site) file(name string) string { return filepath.Join(site.dir, name) }

// newPage returns a new page belonging to site.
func (site *Site) newPage(short string) *page {
	p := &page{
		id:     short,
		params: make(tPage),
	}
	site.pagesByID[p.id] = p
	return p
}

// data parses the named yaml file and returns its structured data.
func (site *Site) data(name string) (interface{}, error) {
	data, err := ioutil.ReadFile(site.file("_content/" + name + ".yaml"))
	if err != nil {
		return nil, err
	}
	var d interface{}
	if err := yaml.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// pageByID returns the page with a given path.
func (site *Site) pageByPath(path string) (tPage, error) {
	p := site.pagesByID[strings.Trim(path, "/")]
	if p == nil {
		return nil, fmt.Errorf("no such page with path %q", path)
	}
	return p.params, nil
}

// pagesGlob returns the pages with IDs matching glob.
func (site *Site) pagesGlob(glob string) ([]tPage, error) {
	_, err := path.Match(glob, "")
	if err != nil {
		return nil, err
	}
	glob = strings.Trim(glob, "/")
	var out []tPage
	for _, p := range site.pagesByID {
		if ok, _ := path.Match(glob, p.id); ok {
			out = append(out, p.params)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i]["Path"].(string) < out[j]["Path"].(string)
	})
	return out, nil
}

// newest returns the pages sorted newest first,
// breaking ties by .linkTitle or else .title.
func newest(pages []tPage) []tPage {
	out := make([]tPage, len(pages))
	copy(out, pages)

	sort.Slice(out, func(i, j int) bool {
		pi := out[i]
		pj := out[j]
		di, _ := pi["Date"].(time.Time)
		dj, _ := pj["Date"].(time.Time)
		if !di.Equal(dj) {
			return di.After(dj)
		}
		ti, _ := pi["linkTitle"].(string)
		tj, _ := pj["linkTitle"].(string)
		if ti != tj {
			return ti < tj
		}
		return false
	})
	return out
}

// Open returns the content to serve at the given path.
// This function makes Site an http.FileServer, for easy HTTP serving.
func (site *Site) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(name, "/")
	switch ext := path.Ext(name); ext {
	case ".css", ".jpeg", ".jpg", ".js", ".png", ".svg", ".txt":
		if f, err := os.Open(site.file("_content/" + name)); err == nil {
			return f, nil
		}

	case ".html":
		id := strings.TrimSuffix(name, "/index.html")
		if name == "index.html" {
			id = ""
		}
		if p := site.pagesByID[id]; p != nil {
			if redir, ok := p.params["redirect"].(string); ok {
				s := fmt.Sprintf(redirectFmt, redir)
				return &httpFile{strings.NewReader(s), int64(len(s))}, nil
			}
			return &httpFile{bytes.NewReader(p.html), int64(len(p.html))}, nil
		}
	}

	if !strings.HasSuffix(name, ".html") {
		if f, err := site.Open(name + "/index.html"); err == nil {
			size, err := f.Seek(0, io.SeekEnd)
			f.Close()
			if err == nil {
				return &httpDir{httpFileInfo{"index.html", size, false}, 0}, nil
			}
		}
	}

	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

type httpFile struct {
	io.ReadSeeker
	size int64
}

func (*httpFile) Close() error                 { return nil }
func (f *httpFile) Stat() (os.FileInfo, error) { return &httpFileInfo{".", f.size, false}, nil }
func (*httpFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("readdir not available")
}

const redirectFmt = `<!DOCTYPE html><html><head><title>%s</title><link rel="canonical" href="%[1]s"/><meta name="robots" content="noindex"><meta charset="utf-8" /><meta http-equiv="refresh" content="0; url=%[1]s" /></head></html>`

type httpDir struct {
	info httpFileInfo
	off  int // 0 or 1
}

func (*httpDir) Close() error                   { return nil }
func (*httpDir) Read([]byte) (int, error)       { return 0, fmt.Errorf("read not available") }
func (*httpDir) Seek(int64, int) (int64, error) { return 0, fmt.Errorf("seek not available") }
func (*httpDir) Stat() (os.FileInfo, error)     { return &httpFileInfo{".", 0, true}, nil }
func (d *httpDir) Readdir(count int) ([]os.FileInfo, error) {
	if count == 0 {
		return nil, nil
	}
	if d.off > 0 {
		return nil, io.EOF
	}
	d.off = 1
	return []os.FileInfo{&d.info}, nil
}

type httpFileInfo struct {
	name string
	size int64
	dir  bool
}

func (info *httpFileInfo) Name() string { return info.name }
func (info *httpFileInfo) Size() int64  { return info.size }
func (info *httpFileInfo) Mode() os.FileMode {
	if info.dir {
		return os.ModeDir | 0555
	}
	return 0444
}
func (info *httpFileInfo) ModTime() time.Time { return time.Time{} }
func (info *httpFileInfo) IsDir() bool        { return info.dir }
func (info *httpFileInfo) Sys() interface{}   { return nil }
