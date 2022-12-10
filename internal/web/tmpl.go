// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/present"
	"gopkg.in/yaml.v3"
)

// A siteDir is a site extended with a known directory for interpreting relative paths.
type siteDir struct {
	*Site
	dir string
}

func toString(x interface{}) string {
	switch x := x.(type) {
	case string:
		return x
	case template.HTML:
		return string(x)
	case nil:
		return ""
	default:
		panic(fmt.Sprintf("cannot toString %T", x))
	}
}

// data parses the named yaml file (relative to dir) and returns its structured data.
func (site *siteDir) data(name string) (interface{}, error) {
	data, err := site.readFile(site.dir, name)
	if err != nil {
		return nil, err
	}
	var d interface{}
	if err := yaml.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	return d, nil
}

func first(n int, list reflect.Value) reflect.Value {
	if !list.IsValid() {
		return list
	}
	if list.Kind() == reflect.Interface {
		if list.IsNil() {
			return list
		}
		list = list.Elem()
	}

	if list.Len() < n {
		return list
	}
	return list.Slice(0, n)
}

// markdown is the function provided to templates.
func markdown(data interface{}) (template.HTML, error) {
	h, err := markdownToHTML(toString(data))
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(string(h))
	if strings.HasPrefix(s, "<p>") && strings.HasSuffix(s, "</p>") && strings.Count(s, "<p>") == 1 {
		h = template.HTML(strings.TrimSpace(s[len("<p>") : len(s)-len("</p>")]))
	}
	return h, nil
}

func (site *siteDir) readfile(name string) (string, error) {
	data, err := site.readFile(site.dir, name)
	return string(data), err
}

// page returns the page params for the page with a given url u.
// The url may or may not have its leading slash.
func (site *siteDir) page(u string) (Page, error) {
	if !path.IsAbs(u) {
		u = path.Join(site.dir, u)
	}
	p, err := site.openPage(strings.Trim(u, "/"))
	if err != nil {
		return nil, err
	}
	return p.page, nil
}

// Pages returns the pages found in files matching glob.
func (site *Site) Pages(glob string) ([]Page, error) {
	return (&siteDir{site, "."}).pages(glob)
}

// pages returns the page params for pages with urls matching glob.
func (site *siteDir) pages(glob string) ([]Page, error) {
	if !path.IsAbs(glob) {
		glob = path.Join(site.dir, glob)
	}
	// TODO(rsc): Add a cache?
	_, err := path.Match(glob, "")
	if err != nil {
		return nil, err
	}
	glob = strings.Trim(glob, "/")
	if glob == "" {
		glob = "."
	}
	matches, err := fs.Glob(site.fs, glob)
	if err != nil {
		return nil, err
	}
	var out []Page
	for _, file := range matches {
		if !strings.HasSuffix(file, ".md") && !strings.HasSuffix(file, ".html") {
			f := path.Join(file, "index.md")
			if _, err := fs.Stat(site.fs, f); err != nil {
				f = path.Join(file, "index.html")
				if _, err = fs.Stat(site.fs, f); err != nil {
					continue
				}
			}
			file = f
		}
		p, err := site.openPage(file)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", file, err)
		}
		out = append(out, p.page)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i]["URL"].(string) < out[j]["URL"].(string)
	})
	return out, nil
}

// file parses the named file (relative to dir) and returns its content as a string.
func (site *siteDir) file(name string) (string, error) {
	data, err := site.readFile(site.dir, name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func raw(s interface{}) template.HTML {
	return template.HTML(toString(s))
}

func yamlFn(s string) (interface{}, error) {
	var d interface{}
	if err := yaml.Unmarshal([]byte(s), &d); err != nil {
		return nil, err
	}
	return d, nil
}

func presentStyle(s string) template.HTML {
	return template.HTML(present.Style(s))
}
