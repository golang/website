// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/tools/godoc"
	"golang.org/x/website/internal/env"
)

// This file holds common code from the x/tools/godoc serving engine.
// It's being used during the transition. See golang.org/issue/29206.

// extractMetadata extracts the godoc.Metadata from a byte slice.
// It returns the godoc.Metadata value and the remaining data.
// If no metadata is present the original byte slice is returned.
//
func extractMetadata(b []byte) (meta godoc.Metadata, tail []byte, _ error) {
	tail = b
	if !bytes.HasPrefix(b, jsonStart) {
		return godoc.Metadata{}, tail, nil
	}
	end := bytes.Index(b, jsonEnd)
	if end < 0 {
		return godoc.Metadata{}, tail, nil
	}
	b = b[len(jsonStart)-1 : end+1] // drop leading <!-- and include trailing }
	if err := json.Unmarshal(b, &meta); err != nil {
		return godoc.Metadata{}, nil, err
	}
	tail = tail[end+len(jsonEnd):]
	return meta, tail, nil
}

var (
	jsonStart = []byte("<!--{")
	jsonEnd   = []byte("}-->")
)

// googleCN reports whether request r is considered
// to be served from golang.google.cn.
// TODO: This is duplicated within internal/proxy. Move to a common location.
func googleCN(r *http.Request) bool {
	if r.FormValue("googlecn") != "" {
		return true
	}
	if strings.HasSuffix(r.Host, ".cn") {
		return true
	}
	if !env.CheckCountry() {
		return false
	}
	switch r.Header.Get("X-Appengine-Country") {
	case "", "ZZ", "CN":
		return true
	}
	return false
}
