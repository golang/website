// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"golang.org/x/website"
	"golang.org/x/website/internal/web"
	"golang.org/x/website/internal/webtest"
)

// Test that the release history page includes expected entries.
//
// At this time, the test is very strict and checks that all releases
// from Go 1 to Go 1.14.2 are included with exact HTML content.
// It can be relaxed whenever the presentation of the release history
// page needs to be changed.
func TestReleaseHistory(t *testing.T) {
	fsys := website.Content
	site, err := web.NewSite(fsys)
	if err != nil {
		t.Fatal(err)
	}
	mux := registerHandlers(fsys, site)

	webtest.TestHandler(t, "testdata/release.txt", mux)
}
