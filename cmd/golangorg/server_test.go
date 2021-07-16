// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/website/internal/webtest"
)

func init() {
	isTestBinary = true
}

func TestWeb(t *testing.T) {
	h := NewHandler("../../_content", runtime.GOROOT())
	files, err := filepath.Glob("testdata/*.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if filepath.ToSlash(file) != "testdata/live.txt" {
			webtest.TestHandler(t, file, h)
		}
	}
}
