// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package markdown provides a wrapper for rendering Markdown. It is intended
// to be used on the golang.org website.
//
// This package is not intended for general use, and its API is not guaranteed
// to be stable.
package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

// Render converts a limited and opinionated flavor of Markdown (compliant with
// CommonMark 0.29) to HTML for the purposes of golang.org websites. This should
// not be adjusted except for the needs of *.golang.org.
//
// The Markdown source may contain raw HTML and Go templates. Sanitization of
// untrusted content is not performed: the caller is responsible for ensuring
// that only trusted content is provided.
func Render(src []byte) ([]byte, error) {
	// html.WithUnsafe allows use of raw HTML, which we need for tables.
	md := goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
