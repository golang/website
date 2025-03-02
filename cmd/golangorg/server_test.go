// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http/httptest"
	"net/url"
	"os"
	pathpkg "path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/website"
	"golang.org/x/website/internal/history"
	"golang.org/x/website/internal/webtest"
)

func TestWeb(t *testing.T) {
	h := NewHandler("../../_content", runtime.GOROOT())

	files, err := filepath.Glob("testdata/*.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if filepath.ToSlash(file) == "testdata/live.txt" {
			continue
		}
		webtest.TestHandler(t, file, h)
	}
}

var bads = []string{
	"&amp;lt;",
	"&amp;gt;",
	"&amp;amp;",
	" < ",
	"<-",
	"& ",
}

var ignoreBads = []string{
	// This JS appears on all the talks pages.
	`window["location"] && window["location"]["hostname"] == "go.dev/talks"`,
}

// findBad returns (only) the lines containing badly escaped HTML in body.
// If findBad returns the empty string, there is no badly escaped HTML.
func findBad(body string) string {
	lines := strings.SplitAfter(body, "\n")
	var out []string
Lines:
	for _, line := range lines {
		for _, ig := range ignoreBads {
			if strings.Contains(line, ig) {
				continue Lines
			}
		}
		for _, b := range bads {
			if strings.Contains(line, b) {
				out = append(out, line)
				break
			}
		}
	}
	return strings.Join(out, "")
}

func TestAll(t *testing.T) {
	h := NewHandler("../../_content", runtime.GOROOT())

	get := func(url string) (code int, body string, err error) {
		if url == "https://go.dev/rebuild" {
			// /rebuild reads from cloud storage so pretend it's fine.
			return 200, "", nil
		}
		rec := httptest.NewRecorder()
		rec.Body = new(bytes.Buffer)
		h.ServeHTTP(rec, httptest.NewRequest("GET", url, nil))
		if rec.Code != 200 && rec.Code/10 != 30 {
			return rec.Code, rec.Body.String(), fmt.Errorf("GET %s: %d, want 200 or 30x", url, rec.Code)
		}
		return rec.Code, rec.Body.String(), nil
	}

	// Assume any URL with these prefixes exists.
	skips := []string{
		"/issue/",
		"/pkg/",
		"/s/",
		"/wiki/",
		"/play/p/",
	}

	// Do not process these paths or path prefixes.
	ignores := []string{
		// Wiki is in a different repo; errors there should not block production push.
		"/wiki/",

		// Support files not meant to be served directly.
		"/doc/articles/wiki/",
		"/talks/2013/highperf/",
		"/talks/2016/refactor/",
		"/tour/static/partials/",
	}

	// Only check and report a URL the first time we see it.
	// Otherwise we recheck all the URLs in the page frames for every page.
	checked := make(map[string]bool)

	testTree := func(dir, prefix string) {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				t.Fatal(err)
			}
			path = filepath.ToSlash(path)
			siteURL := strings.TrimPrefix(path, dir)
			for _, ig := range ignores {
				if strings.HasPrefix(siteURL, ig) {
					return nil
				}
			}
			siteURL = prefix + siteURL // add https://go.dev/

			if strings.HasSuffix(path, ".md") ||
				strings.HasSuffix(path, ".html") ||
				strings.HasSuffix(path, ".article") ||
				strings.HasSuffix(path, ".slide") {
				if !strings.Contains(path, "/talks/") {
					siteURL = strings.TrimSuffix(siteURL, pathpkg.Ext(path))
				}
				if strings.HasSuffix(siteURL, "/index") {
					siteURL = strings.TrimSuffix(siteURL, "index")
				}

				// Check that page can be loaded.
				_, body, err := get(siteURL)
				if err != nil {
					t.Errorf("%v\n%s", err, body)
					return nil
				}

				// Check that page is valid HTML.
				// First check for over- or under-escaped HTML.
				bad := findBad(body)
				if bad != "" {
					t.Errorf("GET %s: contains improperly escaped HTML\n%s", siteURL, bad)
					return nil
				}

				// Now check all the links to other pages on this server.
				// (Pages on other servers are too expensive to check
				// and would cause test failures if servers went down
				// or moved their contents.)
				doc, err := html.Parse(strings.NewReader(body))
				if err != nil {
					t.Errorf("GET %s: parsing HTML: %v", siteURL, err)
					return nil
				}

				base, err := url.Parse(siteURL)
				if err != nil {
					t.Fatalf("cannot parse site URL: %v", err)
				}

				// Walk HTML looking for <a href=...>, <img src=...>, and <script src=...>.
				var checkLinks func(*html.Node)
				checkLinks = func(n *html.Node) {
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						checkLinks(c)
					}
					var targ string
					if n.Type == html.ElementNode {
						switch n.Data {
						case "a":
							targ = findAttr(n, "href")
						case "img", "script":
							targ = findAttr(n, "src")
						}
					}
					// Ignore no target or #fragment.
					if targ == "" || strings.HasPrefix(targ, "#") {
						return
					}

					// Parse target as URL.
					u, err := url.Parse(targ)
					if err != nil {
						t.Errorf("GET %s: found unparseable URL %s: %v", siteURL, targ, err)
						return
					}

					// Check whether URL is canonicalized properly.
					if fix := fixURL(u); fix != "" {
						t.Errorf("GET %s: found link to %s, should be %s", siteURL, targ, fix)
						return
					}

					// Skip checking URLs on other servers.
					if u.Scheme != "" || u.Host != "" {
						return
					}

					// Skip paths that we cannot really check in tests,
					// like the /s/ shortener or redirects to GitHub.
					for _, skip := range skips {
						if strings.HasPrefix(u.Path, skip) {
							return
						}
					}

					// Clear #fragment and build up fully qualified https://go.dev/ URL and check.
					// Only check each link one time during this test,
					// or else we re-check all the frame links on every page.
					u.Fragment = ""
					u.RawFragment = ""
					full := base.ResolveReference(u).String()
					if checked[full] {
						return
					}
					checked[full] = true
					if _, _, err := get(full); err != nil {
						t.Errorf("GET %s: found broken link to %s:\n%s", siteURL, targ, err)
					}
				}
				checkLinks(doc)
			}
			return nil
		})
	}

	testTree("../../_content", "https://go.dev")
}

// fixURL returns the corrected URL for u,
// or the empty string if u is fine.
func fixURL(u *url.URL) string {
	switch u.Host {
	case "golang.org":
		if strings.HasPrefix(u.Path, "/x/") {
			return ""
		}
		fallthrough
	case "go.dev":
		u.Host = ""
		u.Scheme = ""
		if u.Path == "" {
			u.Path = "/"
		}
		return u.String()
	case "blog.golang.org",
		"blog.go.dev",
		"learn.golang.org",
		"learn.go.dev",
		"play.golang.org",
		"play.go.dev",
		"tour.golang.org",
		"tour.go.dev",
		"talks.golang.org",
		"talks.go.dev":
		name, _, _ := strings.Cut(u.Host, ".")
		u.Host = ""
		u.Scheme = ""
		u.Path = "/" + name + u.Path
		return u.String()
	case "github.com":
		if strings.HasPrefix(u.Path, "/golang/go/issues/") {
			u.Host = ""
			u.Scheme = ""
			u.Path = "/issue/" + strings.TrimPrefix(u.Path, "/golang/go/issues/")
			if u.Path == "/issue/new/choose" {
				// A special case to deal with the '/choose' suffix.
				// See comment in internal/redirect.newIssueHandler.
				u.Path = "/issue/new"
			}
			return u.String()
		}
		if strings.HasPrefix(u.Path, "/golang/go/wiki/") {
			u.Host = ""
			u.Scheme = ""
			u.Path = "/wiki/" + strings.TrimPrefix(u.Path, "/golang/go/wiki/")
			return u.String()
		}
	}
	return ""
}

// findAttr returns the value for n's attribute with the given name.
func findAttr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

// TestReleaseNotesHaveDate tests that release notes
// include the date of the corresponding major release.
// See go.dev/issue/54170.
func TestReleaseNotesHaveDate(t *testing.T) {
	for _, r := range history.Majors {
		if r.Version.Before(history.Version{X: 1, Y: 24}) {
			// No dates in release notes before Go 1.24.
			break
		}
		maj := r.Version.MajorPrefix()
		t.Run(maj, func(t *testing.T) {
			name := fmt.Sprintf("doc/go%s.md", maj)
			have, err := fs.ReadFile(website.Content(), name)
			if err != nil {
				t.Fatalf("Go %s release notes (_content/%s) can't be read: %v", maj, name, err)
			}
			want := fmt.Sprintf("[%s %d](/doc/devel/release#go%s)", r.Date.Month, r.Date.Year, r.Version)
			if r.Future {
				want = fmt.Sprintf("%s %d", r.Date.Month, r.Date.Year)
			}
			if !strings.Contains(string(have), want) {
				t.Errorf("Go %s release notes (_content/%s) doesn't contain the release date and link to release history page %q", maj, name, want)
			}
		})
	}
}
