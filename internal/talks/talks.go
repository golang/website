// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package talks

import (
	"io/fs"
	"log"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/present"
	"golang.org/x/website/internal/web"
)

func RegisterHandlers(mux *http.ServeMux, site *web.Site, content fs.FS) error {
	h := &handler{content: content, site: site}
	mux.Handle("/talks/", h)
	return nil
}

type handler struct {
	content fs.FS
	site    *web.Site
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	if h.isDoc(name) {
		err := h.renderDoc(w, r, strings.Trim(name, "/"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	handled, err := h.dirList(w, r)
	if err != nil {
		addr, _, e := net.SplitHostPort(r.RemoteAddr)
		if e != nil {
			addr = r.RemoteAddr
		}
		log.Printf("request from %s: %s", addr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !handled {
		http.FileServer(http.FS(h.content)).ServeHTTP(w, r)
	}
}

func (h *handler) isDoc(name string) bool {
	switch path.Ext(name) {
	case ".slide", ".article":
		return true
	}
	return false
}

func playable(c present.Code) bool {
	// Restrict playable files to only Go source files when using play.golang.org,
	// since there is no method to execute shell scripts there.
	return c.Ext == ".go"
}

// renderDoc reads the present file, gets its template representation,
// and executes the template, sending output to w.
func (h *handler) renderDoc(w http.ResponseWriter, r *http.Request, docFile string) error {
	// Read the input and build the doc structure.
	doc, err := h.parse(docFile, 0)
	if err != nil {
		println("PARSE", err.Error())
		return err
	}

	ext := strings.TrimPrefix(path.Ext(r.URL.Path), ".")
	h.site.ServePage(w, r, web.Page{
		"layout": "/talks/" + ext,
		"doc":    doc,
		"title":  doc.Title,
	})
	return nil
}

func (h *handler) readFile(name string) ([]byte, error) {
	return fs.ReadFile(h.content, filepath.ToSlash(name))
}

func (h *handler) parse(name string, mode present.ParseMode) (*present.Doc, error) {
	ctx := &present.Context{ReadFile: h.readFile}
	f, err := h.content.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ctx.Parse(f, name, mode)
}

// dirList scans the given path and writes a directory listing to w.
// It parses the first part of each .slide file it encounters to display the
// presentation title in the listing.
// If the given path is not a directory, it returns (handled == false, err == nil)
// and writes nothing to w.
func (h *handler) dirList(w http.ResponseWriter, r *http.Request) (handled bool, err error) {
	name := strings.Trim(r.URL.Path, "/")
	info, err := fs.Stat(h.content, name)
	if err != nil || !info.IsDir() {
		return false, err
	}
	if !strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusFound)
		return
	}
	files, err := fs.ReadDir(h.content, name)
	if err != nil {
		return false, err
	}
	d := &dirListData{Path: name}
	for _, fi := range files {
		// skip the golang.org directory
		if name == "." && fi.Name() == "golang.org" {
			continue
		}
		e := dirEntry{
			Name: fi.Name(),
			Path: path.Join(name, fi.Name()),
		}
		if fi.IsDir() && h.showDir(e.Name) {
			d.Dirs = append(d.Dirs, e)
			continue
		}
		if h.isDoc(e.Name) {
			fn := path.Join(name, fi.Name())
			if p, err := h.parse(fn, present.TitlesOnly); err != nil {
				log.Printf("parse(%q, present.TitlesOnly): %v", fn, err)
			} else {
				e.Title = p.Title
			}
			switch path.Ext(e.Path) {
			case ".article":
				d.Articles = append(d.Articles, e)
			case ".slide":
				d.Slides = append(d.Slides, e)
			}
		} else if h.showFile(e.Name) {
			d.Other = append(d.Other, e)
		}
	}
	if d.Path == "." {
		d.Path = ""
	}
	sort.Sort(d.Dirs)
	sort.Sort(d.Slides)
	sort.Sort(d.Articles)
	sort.Sort(d.Other)

	h.site.ServePage(w, r, web.Page{
		"layout": "/talks/dir",
		"title":  d.Path,
		"dir":    d,
	})
	return true, nil
}

// showFile reports whether the given file should be displayed in the list.
func (h *handler) showFile(n string) bool {
	switch path.Ext(n) {
	case ".pdf":
	case ".html":
	case ".go":
	default:
		return h.isDoc(n)
	}
	return true
}

// showDir reports whether the given directory should be displayed in the list.
func (h *handler) showDir(n string) bool {
	if len(n) > 0 && (n[0] == '.' || n[0] == '_') || n == "present" {
		return false
	}
	return true
}

type dirListData struct {
	Path                          string
	Dirs, Slides, Articles, Other dirEntrySlice
}

type dirEntry struct {
	Name, Path, Title string
}

type dirEntrySlice []dirEntry

func (s dirEntrySlice) Len() int           { return len(s) }
func (s dirEntrySlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s dirEntrySlice) Less(i, j int) bool { return s[i].Name < s[j].Name }
