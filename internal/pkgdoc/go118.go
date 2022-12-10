// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.19
// +build !go1.19

package pkgdoc

import (
	"bytes"
	"go/doc"
)

func docPackageHTML(_ *doc.Package, text string) []byte {
	var buf bytes.Buffer
	doc.ToHTML(&buf, text, nil)
	return buf.Bytes()
}
