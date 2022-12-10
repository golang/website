// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the code dealing with package directory trees.

package pkgdoc

import (
	"bytes"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path"
	"strings"
)

type Dir struct {
	Path     string // directory path
	HasPkg   bool   // true if the directory contains at least one package
	Synopsis string // package documentation, if any
	Dirs     []*Dir // subdirectories
}

func (d *Dir) Name() string {
	return path.Base(d.Path)
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

// lookup looks for the *Dir for a given named path, relative to d.
func (d *Dir) lookup(name string) *Dir {
	name = path.Join(d.Path, name)
	if name == d.Path {
		return d
	}
	if d.Path != "." {
		if !strings.HasPrefix(name, d.Path) || name[len(d.Path)] != '/' {
			return nil
		}
		name = name[len(d.Path)+1:]
	}
Walk:
	for i := 0; i <= len(name); i++ {
		if i == len(name) || name[i] == '/' {
			// Find next child along path.
			for _, sub := range d.Dirs {
				if sub.Path == name[:i] {
					d = sub
					continue Walk
				}
			}
			return nil
		}
	}
	return d
}

// list creates a (linear) directory list from a directory tree.
// If filter is set, only the directory entries whose paths match the filter
// are included.
func (d *Dir) list(filter func(string) bool) []DirEntry {
	if d == nil {
		return nil
	}

	root := d
	var list []DirEntry
	root.walk(func(d *Dir, depth int) {
		if depth == 0 || filter != nil && !filter(d.Path) {
			return
		}
		// the path is relative to root.Path - remove the root.Path
		// prefix (the prefix should always be present but avoid
		// crashes and check)
		path := strings.TrimPrefix(d.Path, root.Path)
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
	return list
}

// newDir returns a Dir describing dirpath in fsys.
// When Go files need to be parsed, newDir uses fset.
// If there are no package files and no subdirectories containing packages,
// newDir returns nil.
func newDir(fsys fs.FS, fset *token.FileSet, dirpath string) *Dir {
	var synopses [3]string // prioritized package documentation (0 == highest priority)

	hasPkgFiles := false
	haveSummary := false

	list, err := fs.ReadDir(fsys, dirpath)
	if err != nil {
		// TODO: propagate more. See golang.org/issue/14252.
		log.Printf("newDirTree reading %s: %v", dirpath, err)
	}

	// determine number of subdirectories and if there are package files
	var dirs []*Dir

	for _, de := range list {
		filename := path.Join(dirpath, de.Name())
		switch {
		case isPkgDir(de):
			d := newDir(fsys, fset, filename)
			if d != nil {
				dirs = append(dirs, d)
			}

		case !haveSummary && isPkgFile(de):
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
				case path.Base(dirpath):
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

	return &Dir{
		Path:     dirpath,
		HasPkg:   hasPkgFiles,
		Synopsis: synopsis,
		Dirs:     dirs,
	}
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

// walk calls f for each directory in the tree rooted at d, including d itself.
// The depth argument specifies the depth of the directory in the tree.
// The depth of d itself is 0.
func (d *Dir) walk(f func(d *Dir, depth int)) {
	walkDirs(f, d, 0)
}

func walkDirs(f func(d *Dir, depth int), d *Dir, depth int) {
	f(d, depth)
	for _, sub := range d.Dirs {
		walkDirs(f, sub, depth+1)
	}
}

func parseFile(fsys fs.FS, fset *token.FileSet, filename string, mode parser.Mode) (*ast.File, error) {
	src, err := fs.ReadFile(fsys, filename)
	if err != nil {
		return nil, err
	}

	// Temporary ad-hoc fix for issue 5247.
	// TODO(gri,dmitshur) Remove this in favor of a better fix, eventually (see issue 32092).
	replaceLinePrefixCommentsWithBlankLine(src)

	return parser.ParseFile(fset, filename, src, mode)
}

func parseFiles(fsys fs.FS, fset *token.FileSet, dirname string, localnames []string) (map[string]*ast.File, error) {
	files := make(map[string]*ast.File)
	for _, f := range localnames {
		filename := path.Join(dirname, f)
		file, err := parseFile(fsys, fset, filename, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files[filename] = file
	}

	return files, nil
}

var linePrefix = []byte("//line ")

// This function replaces source lines starting with "//line " with a blank line.
// It does this irrespective of whether the line is truly a line comment or not;
// e.g., the line may be inside a string, or a /*-style comment; however that is
// rather unlikely (proper testing would require a full Go scan which we want to
// avoid for performance).
func replaceLinePrefixCommentsWithBlankLine(src []byte) {
	for {
		i := bytes.Index(src, linePrefix)
		if i < 0 {
			break // we're done
		}
		// 0 <= i && i+len(linePrefix) <= len(src)
		if i == 0 || src[i-1] == '\n' {
			// at beginning of line: blank out line
			for i < len(src) && src[i] != '\n' {
				src[i] = ' '
				i++
			}
		} else {
			// not at beginning of line: skip over prefix
			i += len(linePrefix)
		}
		// i <= len(src)
		src = src[i:]
	}
}
