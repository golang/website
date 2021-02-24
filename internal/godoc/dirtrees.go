// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

// This file contains the code dealing with package directory trees.

package godoc

import (
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path"
	"sort"
	"strings"
)

type Directory struct {
	Path     string       // directory path
	HasPkg   bool         // true if the directory contains at least one package
	Synopsis string       // package documentation, if any
	Dirs     []*Directory // subdirectories
}

func (d *Directory) Name() string {
	return path.Base(d.Path)
}

func isPkgFile(fi fs.DirEntry) bool {
	name := fi.Name()
	return !fi.IsDir() &&
		path.Ext(name) == ".go" &&
		!strings.HasSuffix(fi.Name(), "_test.go") // ignore test files
}

func isPkgDir(fi fs.DirEntry) bool {
	name := fi.Name()
	return fi.IsDir() &&
		name != "testdata" &&
		len(name) > 0 && name[0] != '_' && name[0] != '.' // ignore _files and .files
}

func newDirTree(fsys fs.FS, fset *token.FileSet, abspath string) *Directory {
	var synopses [3]string // prioritized package documentation (0 == highest priority)

	hasPkgFiles := false
	haveSummary := false

	list, err := fs.ReadDir(fsys, toFS(abspath))
	if err != nil {
		// TODO: propagate more. See golang.org/issue/14252.
		log.Printf("newDirTree reading %s: %v", abspath, err)
	}

	// determine number of subdirectories and if there are package files
	var dirchs []chan *Directory
	var dirs []*Directory

	for _, d := range list {
		name := d.Name()
		filename := path.Join(abspath, name)
		switch {
		case isPkgDir(d):
			dir := newDirTree(fsys, fset, filename)
			if dir != nil {
				dirs = append(dirs, dir)
			}

		case !haveSummary && isPkgFile(d):
			// looks like a package file, but may just be a file ending in ".go";
			// don't just count it yet (otherwise we may end up with hasPkgFiles even
			// though the directory doesn't contain any real package files - was bug)
			// no "optimal" package synopsis yet; continue to collect synopses
			const flags = parser.ParseComments | parser.PackageClauseOnly
			file, err := parseFile(fsys, fset, filename, flags)
			if err != nil {
				log.Printf("parsing %v: %v", filename, err)
				break
			}

			hasPkgFiles = true
			if file.Doc != nil {
				// prioritize documentation
				i := -1
				switch file.Name.Name {
				case name:
					i = 0 // normal case: directory name matches package name
				case "main":
					i = 1 // directory contains a main package
				default:
					i = 2 // none of the above
				}
				if 0 <= i && i < len(synopses) && synopses[i] == "" {
					synopses[i] = doc.Synopsis(file.Doc.Text())
				}
			}
			haveSummary = synopses[0] != ""
		}
	}

	// create subdirectory tree
	for _, ch := range dirchs {
		if d := <-ch; d != nil {
			dirs = append(dirs, d)
		}
	}

	// We need to sort the dirs slice because
	// it is appended again after reading from dirchs.
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Path < dirs[j].Path
	})

	// if there are no package files and no subdirectories
	// containing package files, ignore the directory
	if !hasPkgFiles && len(dirs) == 0 {
		return nil
	}

	// select the highest-priority synopsis for the directory entry, if any
	synopsis := ""
	for _, synopsis = range synopses {
		if synopsis != "" {
			break
		}
	}

	return &Directory{
		Path:     abspath,
		HasPkg:   hasPkgFiles,
		Synopsis: synopsis,
		Dirs:     dirs,
	}
}

// toFS returns the io/fs name for path (no leading slash).
func toFS(name string) string {
	if name == "/" {
		return "."
	}
	return path.Clean(strings.TrimPrefix(name, "/"))
}

// walk calls f(d, depth) for each directory d in the tree rooted at dir, including dir itself.
// The depth argument specifies the depth of d in the tree.
// The depth of dir itself is 0.
func (dir *Directory) walk(f func(d *Directory, depth int)) {
	walkDirs(f, dir, 0)
}

func walkDirs(f func(d *Directory, depth int), d *Directory, depth int) {
	f(d, depth)
	for _, sub := range d.Dirs {
		walkDirs(f, sub, depth+1)
	}
}

// lookup looks for the *Directory for a given named path, relative to dir.
func (dir *Directory) lookup(name string) *Directory {
	name = path.Join(dir.Path, name)
	if name == dir.Path {
		return dir
	}
	dirPathLen := len(dir.Path)
	if dir.Path == "/" {
		dirPathLen = 0 // so path[dirPathLen] is a slash
	}
	if !strings.HasPrefix(name, dir.Path) || name[dirPathLen] != '/' {
		println("NO", name, dir.Path)
		return nil
	}
	d := dir
Walk:
	for i := dirPathLen + 1; i <= len(name); i++ {
		if i == len(name) || name[i] == '/' {
			// Find next child along path.
			for _, sub := range d.Dirs {
				if sub.Path == name[:i] {
					d = sub
					continue Walk
				}
			}
			println("LOST", name[:i])
			return nil
		}
	}
	return d
}

// DirEntry describes a directory entry.
// The Depth gives the directory depth relative to the overall list,
// for use in presenting a hierarchical directory entry.
type DirEntry struct {
	Depth    int    // >= 0
	Path     string // relative path to directory from listing start
	HasPkg   bool   // true if the directory contains at least one package
	Synopsis string // package documentation, if any
}

func (d *DirEntry) Name() string {
	return path.Base(d.Path)
}

type DirList struct {
	List []DirEntry
}

// listing creates a (linear) directory listing from a directory tree.
// If skipRoot is set, the root directory itself is excluded from the list.
// If filter is set, only the directory entries whose paths match the filter
// are included.
//
func (dir *Directory) listing(filter func(string) bool) *DirList {
	if dir == nil {
		return nil
	}

	var list []DirEntry
	dir.walk(func(d *Directory, depth int) {
		if depth == 0 || filter != nil && !filter(d.Path) {
			return
		}
		// the path is relative to root.Path - remove the root.Path
		// prefix (the prefix should always be present but avoid
		// crashes and check)
		path := strings.TrimPrefix(d.Path, dir.Path)
		// remove leading separator if any - path must be relative
		path = strings.TrimPrefix(path, "/")
		list = append(list, DirEntry{
			Depth:    depth,
			Path:     path,
			HasPkg:   d.HasPkg,
			Synopsis: d.Synopsis,
		})
	})

	if len(list) == 0 {
		return nil
	}
	return &DirList{list}
}
