// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strings"
)

// googleCN reports whether request r is considered to be arriving from China.
// Typically that means the request is for host golang.google.cn,
// but we also report true for requests that set googlecn=1 as a query parameter.
func googleCN(r *http.Request) bool {
	return r.FormValue("googlecn") != "" || strings.HasSuffix(r.Host, ".cn")
}
