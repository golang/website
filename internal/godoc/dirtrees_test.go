// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"os"
	"runtime"
	"sort"
	"testing"
)

func TestNewDirTree(t *testing.T) {
	c := NewCorpus(os.DirFS(runtime.GOROOT()))
	// 3 levels deep is enough for testing
	dir := c.newDirectory("/", 3)

	processDir(t, dir)
}

func processDir(t *testing.T, dir *Directory) {
	var list []string
	for _, d := range dir.Dirs {
		list = append(list, d.Name)
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
		corpus := NewCorpus(fs)
		corpus.newDirectory("/", -1)
	}
}

func TestIsGOROOT(t *testing.T) {
	tests := []struct {
		path     string
		isGOROOT bool
	}{
		{runtime.GOROOT(), true},
		{"/tmp/", false},
	}

	for _, item := range tests {
		fs := os.DirFS(item.path)
		if isGOROOT(fs) != item.isGOROOT {
			t.Errorf("%s: isGOROOT = %v, want %v", item.path, !item.isGOROOT, item.isGOROOT)
		}
	}
}
