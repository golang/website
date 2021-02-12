// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.16

// Package static exports the static content as an embed.FS.
package static

import "embed"

// FS is the static content as a file system.
//go:embed *
var FS embed.FS
