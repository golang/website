// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"io/fs"
	"sync"
	"time"

	"golang.org/x/website/internal/api"
)

// A Corpus holds all the state related to serving and indexing a
// collection of Go code.
//
// Construct a new Corpus with NewCorpus, then modify options,
// then call its Init method.
type Corpus struct {
	fs fs.FS

	// Verbose logging.
	Verbose bool

	// Send a value on this channel to trigger a metadata refresh.
	// It is buffered so that if a signal is not lost if sent
	// during a refresh.
	refreshMetadataSignal chan bool

	// file system information
	fsModified  rwValue // timestamp of last call to invalidateIndex
	docMetadata rwValue // mapping from paths to *Metadata

	// flag to check whether a corpus is initialized or not
	initMu   sync.RWMutex
	initDone bool

	// pkgAPIInfo contains the information about which package API
	// features were added in which version of Go.
	pkgAPIInfo api.DB
}

// NewCorpus returns a new Corpus from a filesystem.
// The returned corpus has all indexing enabled and MaxResults set to 1000.
// Change or set any options on Corpus before calling the Corpus.Init method.
func NewCorpus(fsys fs.FS) *Corpus {
	c := &Corpus{
		fs:                    fsys,
		refreshMetadataSignal: make(chan bool, 1),
	}
	return c
}

func (c *Corpus) FSModifiedTime() time.Time {
	_, ts := c.fsModified.Get()
	return ts
}

// Init initializes Corpus, once options on Corpus are set.
// It must be called before any subsequent method calls.
func (c *Corpus) Init() error {
	c.updateMetadata()
	go c.refreshMetadataLoop()

	c.initMu.Lock()
	c.initDone = true
	c.initMu.Unlock()
	return nil
}
