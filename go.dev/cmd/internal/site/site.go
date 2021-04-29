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

	"golang.org/x/go.dev/cmd/internal/html/template"
	"gopkg.in/yaml.v3"
)

// A Site holds metadata about the entire site.
type Site struct {
	URL   string
	Title string

	pages     []*Page
	pagesByID map[string]*Page
	dir       string
	redirects map[string]string
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
		redirects: make(map[string]string),
		pagesByID: make(map[string]*Page),
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

	// Assign pages to sections and sort section lists.
	for _, p := range site.pages {
		p.Pages = append(p.Pages, p)
	}
	for _, p := range site.pages {
		if parent := site.pagesByID[p.parent]; parent != nil {
			parent.Pages = append(parent.Pages, p)
		}
	}
	for _, p := range site.pages {
		pages := p.Pages[1:]
		sort.Slice(pages, func(i, j int) bool {
			pi := pages[i]
			pj := pages[j]
			if !pi.Date.Equal(pj.Date.Time) {
				return pi.Date.After(pj.Date.Time)
			}
			ti := pi.LinkTitle
			tj := pj.LinkTitle
			if ti != tj {
				return ti < tj
			}
			return false
		})
	}

	// Now that all pages are loaded and set up, can render all.
	// (Pages can refer to other pages.)
	for _, p := range site.pages {
		if err := p.renderHTML(); err != nil {
			return nil, err
		}
	}

	return site, nil
}

// file returns the full path to the named file within the site.
func (site *Site) file(name string) string { return filepath.Join(site.dir, name) }

// newPage returns a new page belonging to site.
func (site *Site) newPage(short string) *Page {
	p := &Page{
		site:   site,
		id:     short,
		Params: make(map[string]interface{}),
	}
	site.pages = append(site.pages, p)
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
		if target := site.redirects[id]; target != "" {
			s := fmt.Sprintf(redirectFmt, target)
			return &httpFile{strings.NewReader(s), int64(len(s))}, nil
		}
		if p := site.pagesByID[id]; p != nil {
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
