// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"sort"
	"strings"
	"time"

	"golang.org/x/website/internal/web"
)

// newest returns the pages sorted newest first,
// breaking ties by .linkTitle or else .title.
func newest(pages []web.Page) []web.Page {
	out := make([]web.Page, len(pages))
	copy(out, pages)

	sort.Slice(out, func(i, j int) bool {
		pi := out[i]
		pj := out[j]
		di, _ := pi["date"].(time.Time)
		dj, _ := pj["date"].(time.Time)
		if !di.Equal(dj) {
			return di.After(dj)
		}
		ti, _ := pi["linkTitle"].(string)
		if ti == "" {
			ti, _ = pi["title"].(string)
		}
		tj, _ := pj["linkTitle"].(string)
		if tj == "" {
			tj, _ = pj["title"].(string)
		}
		if ti != tj {
			return ti < tj
		}
		return false
	})
	return out
}

// section returns the site section for the given Page,
// defined as the first path element, or else an empty string.
// For example if p's URL is /x/y/z then section is "x".
func section(p web.Page) string {
	u, _ := p["URL"].(string)
	if !strings.HasPrefix(u, "/") {
		return ""
	}
	i := strings.Index(u[1:], "/")
	if i < 0 {
		return ""
	}
	return u[:1+i+1]
}
