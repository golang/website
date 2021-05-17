// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package site

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"golang.org/x/go.dev/cmd/internal/tmplfunc"
	"gopkg.in/yaml.v3"
)

// A page is a single web page.
// It corresponds to some .md file in the content tree.
type page struct {
	id     string // page ID (url path excluding site.BaseURL and trailing slash)
	file   string // .md file for page
	data   []byte // page data (markdown)
	html   []byte // rendered page (HTML)
	params tPage  // parameters passed to templates
}

// A tPage is the template form of the page, the data passed to rendering templates.
type tPage map[string]interface{}

// loadPage loads the site's page from the given file.
// It returns the page but also adds the page to site.pages and site.pagesByID.
func (site *Site) loadPage(file string) (*page, error) {
	id := strings.TrimPrefix(file, "_content/")
	if id == "index.md" {
		id = ""
	} else if strings.HasSuffix(id, "/index.md") {
		id = strings.TrimSuffix(id, "/index.md")
	} else {
		id = strings.TrimSuffix(id, ".md")
	}

	p := site.newPage(id)
	p.file = file

	// Load content, including leading yaml.
	data, err := ioutil.ReadFile(site.file(file))
	if err != nil {
		return nil, err
	}
	if bytes.HasPrefix(data, []byte("---\n")) {
		i := bytes.Index(data, []byte("\n---\n"))
		if i < 0 {
			if bytes.HasSuffix(data, []byte("\n---")) {
				i = len(data) - 4
			}
		}
		if i >= 0 {
			meta := data[4 : i+1]
			err := yaml.Unmarshal(meta, p.params)
			if err != nil {
				return nil, fmt.Errorf("load %s: %v", file, err)
			}

			// Drop YAML but insert the right number of newlines to keep line numbers correct in template errors.
			nl := 0
			for _, c := range data[:i+4] {
				if c == '\n' {
					nl++
				}
			}
			i += 4
			for ; nl > 0; nl-- {
				i--
				data[i] = '\n'
			}
			data = data[i:]
		}
	}
	p.data = data

	// Default linkTitle to title
	if _, ok := p.params["linkTitle"]; !ok {
		p.params["linkTitle"] = p.params["title"]
	}

	// Parse date to Date.
	// Note that YAML parser may have done it for us (!)
	p.params["Date"] = time.Time{}
	if d, ok := p.params["date"].(string); ok {
		t, err := parseDate(d)
		if err != nil {
			return nil, err
		}
		p.params["Date"] = t
	} else if d, ok := p.params["date"].(time.Time); ok {
		p.params["Date"] = d
	}

	// Path, Dir, URL
	urlPath := "/" + p.id
	if strings.HasSuffix(p.file, "/index.md") && p.id != "" {
		urlPath += "/"
	}
	p.params["Path"] = urlPath
	p.params["Dir"] = path.Dir(urlPath)
	p.params["URL"] = strings.TrimRight(site.URL, "/") + urlPath

	// Parent
	if p.id != "" {
		parent := path.Dir("/" + p.id)
		if parent != "/" {
			parent += "/"
		}
		p.params["Parent"] = parent
	}

	// Section
	section := "/"
	if i := strings.Index(p.id, "/"); i >= 0 {
		section = "/" + p.id[:i+1]
	} else if strings.HasSuffix(p.file, "/index.md") {
		section = "/" + p.id + "/"
	}
	p.params["Section"] = section

	// Register aliases. Needs URL.
	aliases, _ := p.params["aliases"].([]interface{})
	for _, alias := range aliases {
		if a, ok := alias.(string); ok {
			site.redirects[strings.Trim(a, "/")] = p.params["URL"].(string)
		}
	}
	return p, nil
}

// renderHTML renders the HTML for the page, leaving it in p.html.
func (site *Site) renderHTML(p *page) error {
	content, err := site.markdownTemplateToHTML(string(p.data), p)
	if err != nil {
		return err
	}
	p.params["Content"] = content

	// Load base template.
	base, err := ioutil.ReadFile(site.file("_templates/layouts/site.tmpl"))
	if err != nil {
		return err
	}
	t := site.clone().New("_templates/layouts/site.tmpl")
	if err := tmplfunc.Parse(t, string(base)); err != nil {
		return err
	}

	// Load page-specific layout template.
	layout, _ := p.params["layout"].(string)
	if layout == "" {
		layout = "default"
	}
	data, err := ioutil.ReadFile(site.file("_templates/layouts/" + layout + ".tmpl"))
	if err != nil {
		return err
	}
	if err := tmplfunc.Parse(t.New(layout), string(data)); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, p.params); err != nil {
		return err
	}
	p.html = buf.Bytes()
	return nil
}

var dateFormats = []string{
	"2006-01-02",
	time.RFC3339,
}

func parseDate(d string) (time.Time, error) {
	for _, f := range dateFormats {
		if tt, err := time.Parse(f, d); err == nil {
			return tt, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date: %s", d)
}
