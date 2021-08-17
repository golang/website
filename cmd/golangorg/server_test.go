// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func TestAll(t *testing.T) {
	h := NewHandler("../../_content", runtime.GOROOT())

	testTree := func(dir, prefix string) {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				t.Fatal(err)
			}
			path = filepath.ToSlash(path)
			if strings.HasSuffix(path, ".md") {
				rec := httptest.NewRecorder()
				rec.Body = new(bytes.Buffer)
				url := prefix + strings.TrimSuffix(strings.TrimPrefix(path, dir), ".md")
				if strings.HasSuffix(url, "/index") {
					url = strings.TrimSuffix(url, "index")
				}
				h.ServeHTTP(rec, httptest.NewRequest("GET", url, nil))
				if rec.Code != 200 && rec.Code != 301 {
					t.Errorf("GET %s: %d, want 200\n%s", url, rec.Code, rec.Body.String())
				}
			}
			return nil
		})
	}

	testTree("../../_content", "https://golang.org")
	testTree("../../go.dev/_content", "https://go.dev")
}
