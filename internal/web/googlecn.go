// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"
	"strings"
)

// GoogleCN reports whether request r is considered to be arriving from China.
// Typically that means the request is for host golang.google.cn,
// but we also report true for requests that set googlecn=1 as a query parameter
// and requests that App Engine geolocates in China or in “unknown country.”
func GoogleCN(r *http.Request) bool {
	if r.FormValue("googlecn") != "" {
		return true
	}
	if strings.HasSuffix(r.Host, ".cn") {
		return true
	}
	switch r.Header.Get("X-Appengine-Country") {
	case "ZZ", "CN":
		return true
	}
	return false
}
