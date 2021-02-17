// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"errors"
	"io/fs"
	"sync"
	"time"
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

	// SummarizePackage optionally specifies a function to
	// summarize a package. It exists as an optimization to
	// avoid reading files to parse package comments.
	//
	// If SummarizePackage returns false for ok, the caller
	// ignores all return values and parses the files in the package
	// as if SummarizePackage were nil.
	//
	// If showList is false, the package is hidden from the
	// package listing.
	SummarizePackage func(pkg string) (summary string, showList, ok bool)

	// Send a value on this channel to trigger a metadata refresh.
	// It is buffered so that if a signal is not lost if sent
	// during a refresh.
	refreshMetadataSignal chan bool

	// file system information
	fsTree      rwValue // *Directory tree of packages, updated with each sync (but sync code is removed now)
	fsModified  rwValue // timestamp of last call to invalidateIndex
	docMetadata rwValue // mapping from paths to *Metadata

	// flag to check whether a corpus is initialized or not
	initMu   sync.RWMutex
	initDone bool

	// pkgAPIInfo contains the information about which package API
	// features were added in which version of Go.
	pkgAPIInfo apiVersions
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
	if err := c.initFSTree(); err != nil {
		return err
	}
	c.updateMetadata()
	go c.refreshMetadataLoop()

	c.initMu.Lock()
	c.initDone = true
	c.initMu.Unlock()
	return nil
}

func (c *Corpus) initFSTree() error {
	dir := c.newDirectory("/", -1)
	if dir == nil {
		return errors.New("godoc: corpus fstree is nil")
	}
	c.fsTree.Set(dir)
	return nil
}
