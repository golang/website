// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkgdoc

import (
	"testing"
	"testing/fstest"

	"golang.org/x/website/internal/web"
)

// TestIgnoredGoFiles tests the scenario where a folder has no .go or .c files,
// but has an ignored go file.
func TestIgnoredGoFiles(t *testing.T) {
	packagePath := "github.com/package"
	packageComment := "main is documented in an ignored .go file"

	fs := fstest.MapFS{
		"lib/godoc/x.html": {},
		"src/" + packagePath + "/ignored.go": {Data: []byte(`// +build ignore

// ` + packageComment + `
package main`)},
	}
	site := web.NewSite(fs)
	h, err := NewServer(fs, site, nil)
	if err != nil {
		t.Fatal(err)
	}
	d := h.(*docs)
	pInfo := d.open("src/"+packagePath, modeAll, "linux", "amd64")

	if pInfo.PDoc == nil {
		t.Error("pInfo.PDoc = nil; want non-nil.")
	} else {
		if got, want := pInfo.PDoc.Doc, packageComment+"\n"; got != want {
			t.Errorf("pInfo.PDoc.Doc = %q; want %q.", got, want)
		}
		if got, want := pInfo.PDoc.Name, "main"; got != want {
			t.Errorf("pInfo.PDoc.Name = %q; want %q.", got, want)
		}
		if got, want := pInfo.PDoc.ImportPath, packagePath; got != want {
			t.Errorf("pInfo.PDoc.ImportPath = %q; want %q.", got, want)
		}
	}
	if pInfo.fset == nil {
		t.Error("pInfo.fset = nil; want non-nil.")
	}
}

func TestIssue5247(t *testing.T) {
	const packagePath = "example.com/p"
	fs := fstest.MapFS{
		"lib/godoc/x.html": {},
		"src/" + packagePath + "/p.go": {Data: []byte(`package p

//line notgen.go:3
// F doc //line 1 should appear
// line 2 should appear
func F()
//line foo.go:100`)}, // No newline at end to check corner cases.
	}

	site := web.NewSite(fs)
	h, err := NewServer(fs, site, nil)
	if err != nil {
		t.Fatal(err)
	}
	d := h.(*docs)

	pInfo := d.open("src/"+packagePath, 0, "linux", "amd64")
	if got, want := pInfo.PDoc.Funcs[0].Doc, "F doc //line 1 should appear\nline 2 should appear\n"; got != want {
		t.Errorf("pInfo.PDoc.Funcs[0].Doc = %q; want %q", got, want)
	}
}
