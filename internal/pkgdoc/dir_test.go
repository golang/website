// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkgdoc

import (
	"go/token"
	"os"
	"runtime"
	"sort"
	"testing"
	"testing/fstest"
)

func TestNewDirTree(t *testing.T) {
	d := newDir(os.DirFS(runtime.GOROOT()), token.NewFileSet(), "src")
	processDir(t, d)
}

func processDir(t *testing.T, d *Dir) {
	var list []string
	for _, child := range d.Dirs {
		list = append(list, child.Name())
		// recursively process the lower level
		processDir(t, child)
	}

	if sort.StringsAreSorted(list) == false {
		t.Errorf("list: %v is not sorted\n", list)
	}
}

func TestIssue45614(t *testing.T) {
	fs := fstest.MapFS{
		"src/index/suffixarray/gen.go": {
			Data: []byte(`// P1: directory contains a main package
package main
`)},
		"src/index/suffixarray/suffixarray.go": {
			Data: []byte(`// P0: directory name matches package name
package suffixarray
`)},
	}

	d := newDir(fs, token.NewFileSet(), "src/index/suffixarray")
	if got, want := d.Synopsis, "P0: directory name matches package name"; got != want {
		t.Errorf("d.Synopsis = %q; want %q", got, want)
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
		newDir(fs, token.NewFileSet(), "src")
	}
}
