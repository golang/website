// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"
)

// A pageFile is a Page loaded from a file.
// It corresponds to some .md or .html file in the content tree.
type pageFile struct {
	file string      // .md file for page
	stat fs.FileInfo // stat for file when page was loaded
	url  string      // url excluding site.BaseURL; always begins with slash
	data []byte      // page data (markdown)
	page Page        // parameters passed to templates

	checked int64 // unix nano, atomically updated
}

// A Page is the data for a web page.
// See the package doc comment for details.
type Page map[string]interface{}

func (site *Site) openPage(file string) (*pageFile, error) {
	// Strip trailing .html or .md or /; it all names the same page.
	if strings.HasSuffix(file, "/index.md") {
		file = strings.TrimSuffix(file, "/index.md")
	} else if strings.HasSuffix(file, "/index.html") {
		file = strings.TrimSuffix(file, "/index.html")
	} else if file == "index.md" || file == "index.html" {
		file = "."
	} else if strings.HasSuffix(file, "/") {
		file = strings.TrimSuffix(file, "/")
	} else if strings.HasSuffix(file, ".html") {
		file = strings.TrimSuffix(file, ".html")
	} else {
		file = strings.TrimSuffix(file, ".md")
	}

	now := time.Now().UnixNano()
	if cp, ok := site.cache.Load(file); ok {
		// Have cache entry; only use if the underlying file hasn't changed.
		// To avoid continuous stats, only check it has been 3s since the last one.
		// TODO(rsc): Move caching into a more general layer and cache templates.
		p := cp.(*pageFile)
		if now-atomic.LoadInt64(&p.checked) >= 3e9 {
			info, err := fs.Stat(site.fs, p.file)
			if err == nil && info.ModTime().Equal(p.stat.ModTime()) && info.Size() == p.stat.Size() {
				atomic.StoreInt64(&p.checked, now)
				return p, nil
			}
		}
	}

	// Check md before html to work correctly when x/website is layered atop Go 1.15 goroot during Go 1.15 tests.
	// Want to find x/website's debugging_with_gdb.md not Go 1.15's debuging_with_gdb.html.
	files := []string{file + ".md", file + ".html", path.Join(file, "index.md"), path.Join(file, "index.html")}
	var filePath string
	var b []byte
	var err error
	var stat fs.FileInfo
	for _, filePath = range files {
		stat, err = fs.Stat(site.fs, filePath)
		if err == nil {
			b, err = site.readFile(".", filePath)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return nil, err
	}

	// If we read an index.md or index.html, the canonical relpath is without the index.md/index.html suffix.
	url := path.Join("/", file)
	if name := path.Base(filePath); name == "index.html" || name == "index.md" {
		url, _ = path.Split(path.Join("/", filePath))
	}

	params, body, err := parseMeta(b)
	if err != nil {
		return nil, err
	}

	p := &pageFile{
		file:    filePath,
		stat:    stat,
		url:     url,
		data:    body,
		page:    params,
		checked: now,
	}

	// File, FileData, URL
	p.page["File"] = filePath
	p.page["FileData"] = string(body)
	p.page["URL"] = p.url

	// User-specified redirect: overrides url but not URL.
	if redir, _ := p.page["redirect"].(string); redir != "" {
		p.url = redir
	}

	site.cache.Store(file, p)

	return p, nil
}

var (
	jsonStart = []byte("<!--{")
	jsonEnd   = []byte("}-->")

	yamlStart = []byte("---\n")
	yamlEnd   = []byte("\n---\n")
)

// parseMeta extracts top-of-file metadata from the file contents b.
// If there is no metadata, parseMeta returns Page{}, b, nil.
// Otherwise, the metdata is extracted, and parseMeta returns
// the metadata and the remainder of the file.
// The end of the metadata is overwritten in b to preserve
// the correct number of newlines so that the line numbers in tail
// match the line numbers in b.
//
// A JSON metadata object is bracketed by <!--{...}-->.
// A YAML metadata object is bracketed by "---\n" above and below the YAML.
//
// JSON is typically used in HTML; YAML is typically used in Markdown.
func parseMeta(b []byte) (meta Page, tail []byte, err error) {
	tail = b
	meta = make(Page)
	var end int
	if bytes.HasPrefix(b, jsonStart) {
		end = bytes.Index(b, jsonEnd)
		if end < 0 {
			return
		}
		b = b[len(jsonStart)-1 : end+1] // drop leading <!-- and include trailing }
		if err = json.Unmarshal(b, &meta); err != nil {
			return
		}
		end += len(jsonEnd)
		for k, v := range meta {
			delete(meta, k)
			meta[strings.ToLower(k)] = v
		}
	} else if bytes.HasPrefix(b, yamlStart) {
		end = bytes.Index(b, yamlEnd)
		if end < 0 {
			return
		}
		b = b[len(yamlStart) : end+1] // drop ---\n but include final \n
		if err = yaml.Unmarshal(b, &meta); err != nil {
			return
		}
		end += len(yamlEnd)
	}

	// Put the right number of \n at the start of tail to preserve line numbers.
	nl := bytes.Count(tail[:end], []byte("\n"))
	for i := 0; i < nl; i++ {
		end--
		tail[end] = '\n'
	}
	tail = tail[end:]
	return
}
