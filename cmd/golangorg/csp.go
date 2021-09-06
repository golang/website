// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"sort"
	"strings"
)

// addCSP returns a handler that adds the appropriate Content-Security-Policy header
// to the response and then invokes h.
func addCSP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ks []string
		for k := range csp {
			ks = append(ks, k)
		}
		sort.Strings(ks)

		var sb strings.Builder
		for _, k := range ks {
			sb.WriteString(k)
			sb.WriteString(" ")
			sb.WriteString(strings.Join(csp[k], " "))
			sb.WriteString("; ")
		}

		w.Header().Set("Content-Security-Policy", sb.String())

		h.ServeHTTP(w, r)
	})
}

const (
	self         = "'self'"
	none         = "'none'"
	unsafeInline = "'unsafe-inline'"
)

var csp = map[string][]string{
	"connect-src": {
		"'self'",
		"https://golang.org",
		"https://play.golang.org", // For running playground snippets on the blog.
		"www.google-analytics.com",
		"stats.g.doubleclick.net",
	},
	"default-src": {
		self,
	},
	"font-src": {
		self,
		"fonts.googleapis.com",
		"fonts.gstatic.com",
		"data:",
	},
	"style-src": {
		self,
		unsafeInline,
		"fonts.googleapis.com",
		"feedback.googleusercontent.com",
		"www.gstatic.com",
		"gstatic.com",
		"tagmanager.google.com",
	},
	"frame-src": {
		self,
		"www.google.com",
		"feedback.googleusercontent.com",
		"www.googletagmanager.com",
		"scone-pa.clients6.google.com",
		"www.youtube.com",
	},
	"img-src": {
		self,
		"www.google.com",
		"www.google-analytics.com",
		"ssl.gstatic.com",
		"www.gstatic.com",
		"gstatic.com",
		"data: *",
	},
	"object-src": {
		none,
	},
	"script-src": {
		self,
		"'sha256-n6OdwTrm52KqKm6aHYgD0TFUdMgww4a0GQlIAVrMzck='", // Google Tag Manager main snippet
		"'sha256-4ryYrf7Y5daLOBv0CpYtyBIcJPZkRD2eBPdfqsN3r1M='", // Google Tag Manager Preview mode
		"'sha256-sVKX08+SqOmnWhiySYk3xC7RDUgKyAkmbXV2GWts4fo='", // Google Tag Manager Preview mode
		"www.google.com",
		"apis.google.com",
		"www.gstatic.com",
		"gstatic.com",
		"support.google.com",
		"www.googletagmanager.com",
		"www.google-analytics.com",
		"ssl.google-analytics.com",
		"tagmanager.google.com",
	},
	"frame-ancestors": {
		none,
	},
}
