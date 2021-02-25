// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"io/fs"

	"golang.org/x/website/internal/api"
)

// A Corpus holds all the state related to serving and indexing a
// collection of Go code.
//
// Construct a new Corpus with NewCorpus, then modify options,
// then call its Init method.
type Corpus struct {
	fs fs.FS

	// pkgAPIInfo contains the information about which package API
	// features were added in which version of Go.
	pkgAPIInfo api.DB
}

// NewCorpus returns a new Corpus from a filesystem.
// The returned corpus has all indexing enabled and MaxResults set to 1000.
// Change or set any options on Corpus before calling the Corpus.Init method.
func NewCorpus(fsys fs.FS) *Corpus {
	c := &Corpus{
		fs: fsys,
	}
	return c
}
