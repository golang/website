// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"sort"
	"strings"
)

// buildCSP builds the CSP header.
func buildCSP(kind string) string {
	var ks []string
	for k := range csp {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	var sb strings.Builder
	for _, k := range ks {
		sb.WriteString(k)
		sb.WriteString(" ")
		for _, v := range csp[k] {
			if (kind == "tour" || kind == "talks") && strings.HasPrefix(v, "'sha256-") {
				// Must drop sha256 entries to use unsafe-inline.
				continue
			}
			sb.WriteString(v)
			sb.WriteString(" ")
		}
		if kind == "tour" && k == "script-src" {
			sb.WriteString(" ")
			sb.WriteString(unsafeEval)
		}
		if (kind == "talks" || kind == "tour") && k == "script-src" {
			sb.WriteString(" ")
			sb.WriteString(unsafeInline)
		}
		sb.WriteString("; ")
	}
	return sb.String()
}

// addCSP returns a handler that adds the appropriate Content-Security-Policy header
// to the response and then invokes h.
func addCSP(h http.Handler) http.Handler {
	std := buildCSP("")
	tour := buildCSP("tour")
	talks := buildCSP("talks")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := std
		if strings.HasPrefix(r.URL.Path, "/tour/") {
			csp = tour
		}
		if strings.HasPrefix(r.URL.Path, "/talks/") {
			csp = talks
		}
		w.Header().Set("Content-Security-Policy", csp)
		h.ServeHTTP(w, r)
	})
}

const (
	self         = "'self'"
	none         = "'none'"
	unsafeInline = "'unsafe-inline'"
	unsafeEval   = "'unsafe-eval'"
)

var csp = map[string][]string{
	"connect-src": {
		"'self'",
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
		"player.vimeo.com",
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
		self,
	},
}
