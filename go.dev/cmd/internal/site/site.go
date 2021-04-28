// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package site implements generation of content for serving from go.dev.
// It is meant to support a transition from being a Hugo-based web site
// to being a site compatible with x/website.
package site

import (
	"bytes"
	"crypto/sha256"
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

	"github.com/BurntSushi/toml"
	"golang.org/x/go.dev/cmd/internal/html/template"
	"gopkg.in/yaml.v3"
)

// A Site holds metadata about the entire site.
type Site struct {
	BaseURL      string
	LanguageCode string
	Title        string
	Menus        map[string][]*MenuItem `toml:"menu"`
	IsServer     bool
	Data         map[string]interface{}
	pages        []*Page
	pagesByID    map[string]*Page
	dir          string
	redirects    map[string]string
	base         *template.Template
}

// A MenuItem is a single entry in a menu.
type MenuItem struct {
	Identifier string
	Name       string
	Title      string
	URL        string
	Parent     string
	Weight     int
	Children   []*MenuItem
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

	// Read site config from config.toml.
	if _, err := toml.DecodeFile(site.file("config.toml"), &site); err != nil {
		return nil, fmt.Errorf("parsing site config.toml: %v", err)
	}

	// Group and sort menus.
	for name, list := range site.Menus {
		// Collect top-level items and assign children.
		topsByID := make(map[string]*MenuItem)
		var tops []*MenuItem
		for _, item := range list {
			if p := topsByID[item.Parent]; p != nil {
				p.Children = append(p.Children, item)
				continue
			}
			tops = append(tops, item)
			if item.Identifier != "" {
				topsByID[item.Identifier] = item
			}
		}
		// Sort each top-level item's child list.
		for _, item := range tops {
			c := item.Children
			sort.Slice(c, func(i, j int) bool { return c[i].Weight < c[j].Weight })
		}
		site.Menus[name] = tops
	}

	// Load site data files.
	// site.Data is a directory tree in which each key points at
	// either another directory tree (a subdirectory)
	// or a parsed yaml file.
	site.Data = make(map[string]interface{})
	root := site.file("data")
	err = filepath.Walk(root, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if name == root {
			name = "."
		} else {
			name = name[len(root)+1:]
		}
		if info.IsDir() {
			site.Data[name] = make(map[string]interface{})
			return nil
		}
		if strings.HasSuffix(name, ".yaml") {
			data, err := ioutil.ReadFile(filepath.Join(root, name))
			if err != nil {
				return err
			}
			var d interface{}
			if err := yaml.Unmarshal(data, &d); err != nil {
				return fmt.Errorf("unmarshaling %v: %v", name, err)
			}

			elems := strings.Split(name, "/")
			m := site.Data
			for _, elem := range elems[:len(elems)-1] {
				m = m[elem].(map[string]interface{})
			}
			m[strings.TrimSuffix(elems[len(elems)-1], ".yaml")] = d
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("loading data: %v", err)
	}

	// Implicit home page.
	home := site.newPage("")
	home.Params["Series"] = ""
	home.IsHome = true
	home.Title = site.Title

	// Load site pages from md files.
	err = filepath.Walk(site.file("content"), func(name string, info os.FileInfo, err error) error {
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
			if pi.Weight != pj.Weight {
				return pi.Weight > pj.Weight
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
		Site:   site,
		id:     short,
		Params: make(map[string]interface{}),
	}
	site.pages = append(site.pages, p)
	site.pagesByID[p.id] = p
	return p
}

// Open returns the content to serve at the given path.
// This function makes Site an http.FileServer, for easy HTTP serving.
func (site *Site) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(name, "/")
	switch ext := path.Ext(name); ext {
	case ".css", ".jpeg", ".jpg", ".js", ".png", ".svg", ".txt":
		if f, err := os.Open(site.file("content/" + name)); err == nil {
			return f, nil
		}
		if f, err := os.Open(site.file("static/" + name)); err == nil {
			return f, nil
		}

		// Maybe it is name.hash.ext. Check hash.
		// We will stop generating these eventually,
		// so it doesn't matter that this is slow.
		prefix := name[:len(name)-len(ext)]
		hash := path.Ext(prefix)
		prefix = prefix[:len(prefix)-len(hash)]
		if len(hash) == 1+64 {
			file := site.file("assets/" + prefix + ext)
			if data, err := ioutil.ReadFile(file); err == nil && fmt.Sprintf(".%x", sha256.Sum256(data)) == hash {
				if f, err := os.Open(file); err == nil {
					return f, nil
				}
			}
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
