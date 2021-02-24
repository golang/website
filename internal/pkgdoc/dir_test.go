// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package pkgdoc

import (
	"go/token"
	"os"
	"runtime"
	"sort"
	"testing"
)

func TestNewDirTree(t *testing.T) {
	dir := newDir(os.DirFS(runtime.GOROOT()), token.NewFileSet(), "/src")
	processDir(t, dir)
}

func processDir(t *testing.T, dir *Dir) {
	var list []string
	for _, d := range dir.Dirs {
		list = append(list, d.Name())
		// recursively process the lower level
		processDir(t, d)
	}

	if sort.StringsAreSorted(list) == false {
		t.Errorf("list: %v is not sorted\n", list)
	}
}

func BenchmarkNewDirectory(b *testing.B) {
	if testing.Short() {
		b.Skip("not running tests requiring large file scan in short mode")
	}

	fs := os.DirFS(runtime.GOROOT())

	b.ResetTimer()
	b.ReportAllocs()
	for tries := 0; tries < b.N; tries++ {
		newDir(fs, token.NewFileSet(), "/src")
	}
}
