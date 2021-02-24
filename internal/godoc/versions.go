// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

// This file caches information about which standard library types, methods,
// and functions appeared in what version of Go

package godoc

import (
	"log"

	"golang.org/x/website/internal/api"
)

// InitVersionInfo parses the $GOROOT/api/go*.txt API definition files to discover
// which API features were added in which Go releases.
func (c *Corpus) InitVersionInfo() {
	var err error
	c.pkgAPIInfo, err = api.Load()
	if err != nil {
		// TODO: consider making this fatal, after the Go 1.11 cycle.
		log.Printf("godoc: error parsing API version files: %v", err)
	}
}
