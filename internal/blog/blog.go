// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blog

import (
	"encoding/json"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/website/internal/blog/atom"
	"golang.org/x/website/internal/web"
)

const maxFeed = 10

// atomFeed returns the Atom feed for the go.dev blog, given the go.dev site.
func atomFeed(site *web.Site) ([]byte, error) {
	pages, err := feedPages(site)
	if err != nil {
		return nil, err
	}

	var updated time.Time
	if len(pages) > 0 {
		updated, _ = pages[0]["date"].(time.Time)
	}

	baseURL := "https://go.dev"
	feed := &atom.Feed{
		Title:   "The Go Blog",
		ID:      "tag:blog.golang.org,2013:blog.golang.org", // keep original blog ID
		Updated: atom.Time(updated),
		Link: []atom.Link{{
			Rel:  "self",
			Href: baseURL + "/blog/feed.atom",
		}},
	}

	for _, p := range pages {
		title, _ := p["title"].(string)
		url, _ := p["URL"].(string)
		date, _ := p["date"].(time.Time)
		summary, _ := p["summary"].(string)
		by, _ := p["by"].([]string)
		content, err := site.RenderContent(p, "blogfeed.tmpl")
		if err != nil {
			return nil, err
		}

		e := &atom.Entry{
			Title: title,
			ID:    feed.ID + strings.TrimPrefix(url, "/blog"),
			Link: []atom.Link{{
				Rel:  "alternate",
				Href: baseURL + url,
			}},
			Published: atom.Time(date),
			Updated:   atom.Time(date),
			Summary: &atom.Text{
				Type: "html",
				Body: html.EscapeString(summary),
			},
			Content: &atom.Text{
				Type: "html",
				Body: string(content),
			},
			Author: &atom.Person{
				Name: authors(by),
			},
		}
		feed.Entry = append(feed.Entry, e)
	}

	return xml.Marshal(feed)
}

type jsonItem struct {
	Title   string
	Link    string
	Time    time.Time
	Summary string
	Content string
	Author  string
}

// jsonFeed returns the JSON feed for the go.dev blog, given the go.dev site.
func jsonFeed(site *web.Site) ([]byte, error) {
	pages, err := feedPages(site)
	if err != nil {
		return nil, err
	}

	baseURL := "https://go.dev"
	var feed []jsonItem
	for _, p := range pages {
		title, _ := p["title"].(string)
		url, _ := p["URL"].(string)
		date, _ := p["date"].(time.Time)
		summary, _ := p["summary"].(string)
		by, _ := p["by"].([]string)
		content, err := site.RenderContent(p, "blogfeed.tmpl")
		if err != nil {
			return nil, err
		}
		item := jsonItem{
			Title:   title,
			Link:    baseURL + url,
			Time:    date,
			Summary: summary,
			Content: string(content),
			Author:  authors(by),
		}
		feed = append(feed, item)
	}

	return json.Marshal(feed)
}

func feedPages(site *web.Site) ([]web.Page, error) {
	pages, err := site.Pages("/blog/*")
	if err != nil {
		return nil, err
	}
	sort.Slice(pages, func(i, j int) bool {
		ti, _ := pages[i]["date"].(time.Time)
		tj, _ := pages[j]["date"].(time.Time)
		return ti.After(tj)
	})
	if len(pages) > maxFeed {
		pages = pages[:maxFeed]
	}
	for ; len(pages) > 0; pages = pages[:len(pages)-1] {
		last := pages[len(pages)-1]
		t, _ := last["date"].(time.Time)
		if !t.IsZero() {
			break
		}
	}
	return pages, nil
}

func authors(by []string) string {
	switch len(by) {
	case 0:
		return ""
	case 1:
		return by[0]
	case 2:
		return by[0] + " and " + by[1]
	default:
		return strings.Join(by[:len(by)-1], ", ") + ", and " + by[len(by)-1]
	}
}

var validJSONPFunc = regexp.MustCompile(`(?i)^[a-z_][a-z0-9_.]*$`)

// RegisterFeeds registers the blog Atom and JSON feeds for site on mux,
// using host as a host prefix on the registered paths.
func RegisterFeeds(mux *http.ServeMux, host string, site *web.Site) error {
	atom, err := atomFeed(site)
	if err != nil {
		return err
	}
	atomHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/atom+xml; charset=utf-8")
		w.Write(atom)
	}
	mux.HandleFunc("/blog/feed.atom", atomHandler)
	mux.HandleFunc("/blog/feeds/posts/default", atomHandler)

	json, err := jsonFeed(site)
	if err != nil {
		return err
	}
	jsonHandler := func(w http.ResponseWriter, r *http.Request) {
		if p := r.FormValue("jsonp"); validJSONPFunc.MatchString(p) {
			w.Header().Set("Content-type", "application/javascript; charset=utf-8")
			io.WriteString(w, p+"(")
			defer io.WriteString(w, ")")
		} else {
			w.Header().Set("Content-type", "application/json; charset=utf-8")
		}
		w.Write(json)
	}
	mux.HandleFunc("/blog/.json", jsonHandler)
	return nil
}
