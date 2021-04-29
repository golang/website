// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package site

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/go.dev/cmd/internal/html/template"
	"golang.org/x/go.dev/cmd/internal/tmplfunc"
	"gopkg.in/yaml.v3"
)

// A Page is a single web page.
// It corresponds to some .md file in the content tree.
// Although page is not exported for use by other Go code,
// its exported fields and methods are available to templates.
type Page struct {
	id      string // page ID (url path excluding site.BaseURL and trailing slash)
	file    string // .md file for page
	section string // page section ID
	parent  string // parent page ID
	data    []byte // page data (markdown)
	html    []byte // rendered page (HTML)

	// yaml metadata and data available to templates
	Aliases      []string
	Content      template.HTML
	Date         anyTime
	Description  string `yaml:"description"`
	IsHome       bool
	LinkTitle    string `yaml:"linkTitle"`
	Pages        []*Page
	Params       map[string]interface{}
	Site         *Site
	TheResources []*Resource `yaml:"resources"`
	Title        string
	Weight       int
}

// loadPage loads the site's page from the given file.
// It returns the page but also adds the page to site.pages and site.pagesByID.
func (site *Site) loadPage(file string) (*Page, error) {
	id := strings.TrimPrefix(file, "content/")
	if strings.HasSuffix(id, "/_index.md") {
		id = strings.TrimSuffix(id, "/_index.md")
	} else if strings.HasSuffix(id, "/index.md") {
		id = strings.TrimSuffix(id, "/index.md")
	} else {
		id = strings.TrimSuffix(id, ".md")
	}

	p := site.newPage(id)
	p.file = file
	p.Params["Series"] = ""
	p.Params["series"] = ""

	// Determine section.
	for dir := path.Dir(file); dir != "."; dir = path.Dir(dir) {
		if _, err := os.Stat(site.file(dir + "/_index.md")); err == nil {
			p.section = strings.TrimPrefix(dir, "content/")
			break
		}
	}

	// Determine parent.
	p.parent = p.section
	if p.parent == p.id {
		p.parent = ""
		for dir := path.Dir("content/" + p.id); dir != "."; dir = path.Dir(dir) {
			if _, err := os.Stat(site.file(dir + "/_index.md")); err == nil {
				p.parent = strings.TrimPrefix(dir, "content/")
				break
			}
		}
	}

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
			err := yaml.Unmarshal(meta, p.Params)
			if err != nil {
				return nil, fmt.Errorf("load %s: %v", file, err)
			}
			err = yaml.Unmarshal(meta, p)
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

	// Set a few defaults.
	p.Params["Series"] = p.Params["series"]
	if p.LinkTitle == "" {
		p.LinkTitle = p.Title
	}

	// Register aliases.
	for _, alias := range p.Aliases {
		site.redirects[strings.Trim(alias, "/")] = p.Permalink()
	}

	return p, nil
}

// renderHTML renders the HTML for the page, leaving it in p.html.
func (p *Page) renderHTML() error {
	var err error
	p.Content, err = markdownTemplateToHTML(string(p.data), p)
	if err != nil {
		return err
	}

	// Load base template.
	base, err := ioutil.ReadFile(p.Site.file("layouts/_default/baseof.html"))
	if err != nil {
		return err
	}
	t := p.Site.clone().New("layouts/_default/baseof.html")
	if err := tmplfunc.Parse(t, string(base)); err != nil {
		return err
	}

	// Load page-specific layout template.
	// There are general rules in Hugo, but we don't need to reproduce them here
	// since this will go away.
	layout := "layouts/_default/single.html"
	switch p.id {
	case "":
		layout = "layouts/index.html"
	case "learn":
		layout = "layouts/learn/section.html"
	case "solutions":
		layout = "layouts/solutions/section.html"
	}
	if strings.HasPrefix(p.id, "solutions/") {
		layout = "layouts/solutions/single.html"
	}
	data, err := ioutil.ReadFile(p.Site.file(layout))
	if err != nil {
		return err
	}
	if err := tmplfunc.Parse(t.New(layout), string(data)); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, p); err != nil {
		return err
	}
	p.html = buf.Bytes()
	return nil
}

// An anyTime is a time.Time that accepts any of the anyTimeFormats when unmarshaling.
type anyTime struct {
	time.Time
}

var anyTimeFormats = []string{
	"2006-01-02",
	time.RFC3339,
}

func (t *anyTime) UnmarshalText(data []byte) error {
	for _, f := range anyTimeFormats {
		if tt, err := time.Parse(f, string(data)); err == nil {
			t.Time = tt
			return nil
		}
	}
	return fmt.Errorf("invalid time: %s", data)
}
