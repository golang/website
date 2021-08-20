// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitfs

import (
	"io/fs"
	"io/ioutil"
	"testing"
)

func TestGerrit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Gerrit network access in -short mode")
	}
	r, err := NewRepo("https://go.googlesource.com/scratch")
	if err != nil {
		t.Fatal(err)
	}
	_, fsys, err := r.Clone("HEAD")
	if err != nil {
		t.Fatal(err)
	}
	data, err := fs.ReadFile(fsys, "README.md")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestGitHub(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping GitHub network access in -short mode")
	}
	r, err := NewRepo("https://github.com/rsc/quote")
	if err != nil {
		t.Fatal(err)
	}
	_, fsys, err := r.Clone("HEAD")
	if err != nil {
		t.Fatal(err)
	}
	data, err := fs.ReadFile(fsys, "README.md")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestPack(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/scratch.pack")
	if err != nil {
		t.Fatal(err)
	}
	var s store
	err = unpack(&s, data)
	if err != nil {
		t.Fatal(err)
	}

	h := Hash{0xf6, 0xf7, 0x39, 0x2a, 0x99, 0x9b, 0x3d, 0x75, 0xe2, 0x1c, 0xae, 0xe3, 0x3a, 0xeb, 0x6d, 0x01, 0x92, 0xe8, 0xdc, 0x6b}
	tfs, err := s.commit(h)
	if err != nil {
		t.Fatal(err)
	}

	data, err = fs.ReadFile(tfs, "rsc/greeting.go")
	if err != nil {
		t.Fatal(err)
	}
	println(string(data))
}
