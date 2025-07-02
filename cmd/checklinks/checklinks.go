// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The checklinks command checks for broken links in the gopls
// documentation, recursively traversing from the main page.
//
// Example:
//
//	$ cd x/website
//	$ go run ./cmd/golangorg/ -gopls &
//	$ go run ./cmd/checklinks/
//
// Run golangorg with the GOLANGORG_LOCAL_X_TOOLS=dir environment variable
// to check the links of a locally edited x/tools repository in dir.
package main

import (
	"fmt"
	"iter"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// TODO(adonovan): check <img> elements too.

const verbose = false
const scope = "http://localhost:6060/go.dev/gopls/"

func main() {
	log.SetPrefix("checklinks: ")
	log.SetFlags(0)

	root, err := url.Parse(scope)
	if err != nil {
		log.Fatalf("invalid URL: %v", err)
	}
	if !root.IsAbs() {
		log.Fatalf("URL not absolute: %s", root)
	}
	check(root, root)
}

var seenURLs = make(map[string]bool)

// check loads the 'to' page and checks the validity of links within it.
// The 'from' parameter indicates where it was requested from.
func check(from, to *url.URL) {
	to.Fragment = ""
	str := to.String()

	if !strings.HasPrefix(str, scope) {
		if verbose {
			log.Printf("out of scope: %s", str)
		}
		return // out of scope
	}

	if seenURLs[str] {
		return
	}
	seenURLs[str] = true

	if verbose {
		log.Printf("%s -> %s", from, to)
	}

	doc, err := getHTML(to)
	if err != nil {
		log.Printf("from %s: %v", from, err)
		return
	}
	for href := range links(doc) {
		u, err := to.Parse(href)
		if err != nil {
			log.Printf("from %s: invalid URL: %v", to, err)
			continue
		}
		check(to, u)
	}
}

func getHTML(u *url.URL) (*html.Node, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %q failed: %v", u, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// links returns an iterator over the targets of <a href=...> elements.
func links(node *html.Node) iter.Seq[string] {
	return func(yield func(href string) bool) {
		_ = every(node, yield)
	}
}

func every(n *html.Node, f func(s string) bool) bool {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if !f(a.Val) {
					return false
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if !every(c, f) {
			return false
		}
	}
	return true
}
