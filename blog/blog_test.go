// Copyright 2018 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"golang.org/x/website/internal/webtest"
)

func TestServer(t *testing.T) {
	h, err := blogHandler()
	if err != nil {
		t.Fatal(err)
	}
	webtest.TestHandler(t, "testdata/blog.txt", h)
}
