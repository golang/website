// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package pkgdoc

import (
	"bytes"
	"go/ast"
	"go/build"
	"go/doc"
	"go/token"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Docs struct {
	fs   fs.FS
	root *Dir
}

func NewDocs(fsys fs.FS) *Docs {
	src := newDir(fsys, token.NewFileSet(), "/src")
	root := &Dir{
		Path: "/",
		Dirs: []*Dir{src},
	}
	return &Docs{
		fs:   fsys,
		root: root,
	}
}

type Page struct {
	Dirname string // directory containing the package
	Err     error  // error or nil

	Mode Mode // display metadata from query string

	// package info
	FSet       *token.FileSet // nil if no package documentation
	PDoc       *doc.Package   // nil if no package documentation
	Examples   []*doc.Example // nil if no example code
	Bugs       []*doc.Note    // nil if no BUG comments
	IsMain     bool           // true for package main
	IsFiltered bool           // true if results were filtered

	// directory info
	Dirs    *DirList // nil if no directory information
	DirFlat bool     // if set, show directory in a flat (non-indented) manner
}

func (info *Page) IsEmpty() bool {
	return info.Err != nil || info.PDoc == nil && info.Dirs == nil
}

type Mode uint

const (
	ModeAll     Mode = 1 << iota // do not filter exports
	ModeFlat                     // show directory in a flat (non-indented) manner
	ModeMethods                  // show all embedded methods
	ModeBuiltin                  // don't associate consts, vars, and factory functions with types (not exposed via ?m= query parameter, used for package builtin, see issue 6645)
)

// modeNames defines names for each PageInfoMode flag.
// The order here must match the order of the constants above.
var modeNames = []string{
	"all",
	"flat",
	"methods",
}

// generate a query string for persisting PageInfoMode between pages.
func (m Mode) String() string {
	s := ""
	for i, name := range modeNames {
		if m&(1<<i) != 0 && name != "" {
			if s != "" {
				s += ","
			}
			s += name
		}
	}
	return s
}

// ParseMode computes the PageInfoMode flags by analyzing the request
// URL form value "m". It is value is a comma-separated list of mode names (for example, "all,flat").
func ParseMode(text string) Mode {
	var mode Mode
	for _, k := range strings.Split(text, ",") {
		k = strings.TrimSpace(k)
		for i, name := range modeNames {
			if name == k {
				mode |= 1 << i
			}
		}
	}
	return mode
}

// Doc returns the Page for a package directory abspath.
// Package documentation (Page.PDoc) is extracted from the AST.
// If there is no corresponding package in the
// directory, Page.PDoc is nil. If there are no sub-
// directories, Page.Dirs is nil. If an error occurred, PageInfo.Err is
// set to the respective error but the error is not logged.
func Doc(d *Docs, abspath, relpath string, mode Mode, goos, goarch string) *Page {
	info := &Page{Dirname: abspath, Mode: mode}

	// Restrict to the package files that would be used when building
	// the package on this system.  This makes sure that if there are
	// separate implementations for, say, Windows vs Unix, we don't
	// jumble them all together.
	// Note: If goos/goarch aren't set, the current binary's GOOS/GOARCH
	// are used.
	ctxt := build.Default
	ctxt.IsAbsPath = path.IsAbs
	ctxt.IsDir = func(path string) bool {
		fi, err := fs.Stat(d.fs, toFS(filepath.ToSlash(path)))
		return err == nil && fi.IsDir()
	}
	ctxt.ReadDir = func(dir string) ([]os.FileInfo, error) {
		f, err := fs.ReadDir(d.fs, toFS(filepath.ToSlash(dir)))
		filtered := make([]os.FileInfo, 0, len(f))
		for _, i := range f {
			if mode&ModeAll != 0 || i.Name() != "internal" {
				info, err := i.Info()
				if err == nil {
					filtered = append(filtered, info)
				}
			}
		}
		return filtered, err
	}
	ctxt.OpenFile = func(name string) (r io.ReadCloser, err error) {
		data, err := fs.ReadFile(d.fs, toFS(filepath.ToSlash(name)))
		if err != nil {
			return nil, err
		}
		return ioutil.NopCloser(bytes.NewReader(data)), nil
	}

	// Make the syscall/js package always visible by default.
	// It defaults to the host's GOOS/GOARCH, and golang.org's
	// linux/amd64 means the wasm syscall/js package was blank.
	// And you can't run godoc on js/wasm anyway, so host defaults
	// don't make sense here.
	if goos == "" && goarch == "" && relpath == "syscall/js" {
		goos, goarch = "js", "wasm"
	}
	if goos != "" {
		ctxt.GOOS = goos
	}
	if goarch != "" {
		ctxt.GOARCH = goarch
	}

	pkginfo, err := ctxt.ImportDir(abspath, 0)
	// continue if there are no Go source files; we still want the directory info
	if _, nogo := err.(*build.NoGoError); err != nil && !nogo {
		info.Err = err
		return info
	}

	// collect package files
	pkgname := pkginfo.Name
	pkgfiles := append(pkginfo.GoFiles, pkginfo.CgoFiles...)
	if len(pkgfiles) == 0 {
		// Commands written in C have no .go files in the build.
		// Instead, documentation may be found in an ignored file.
		// The file may be ignored via an explicit +build ignore
		// constraint (recommended), or by defining the package
		// documentation (historic).
		pkgname = "main" // assume package main since pkginfo.Name == ""
		pkgfiles = pkginfo.IgnoredGoFiles
	}

	// get package information, if any
	if len(pkgfiles) > 0 {
		// build package AST
		fset := token.NewFileSet()
		files, err := parseFiles(d.fs, fset, relpath, abspath, pkgfiles)
		if err != nil {
			info.Err = err
			return info
		}

		// ignore any errors - they are due to unresolved identifiers
		pkg, _ := ast.NewPackage(fset, files, simpleImporter, nil)

		// extract package documentation
		info.FSet = fset
		info.IsMain = pkgname == "main"
		// show extracted documentation
		var m doc.Mode
		if mode&ModeAll != 0 {
			m |= doc.AllDecls
		}
		if mode&ModeMethods != 0 {
			m |= doc.AllMethods
		}
		info.PDoc = doc.New(pkg, path.Clean(relpath), m) // no trailing '/' in importpath
		if mode&ModeBuiltin != 0 {
			for _, t := range info.PDoc.Types {
				info.PDoc.Consts = append(info.PDoc.Consts, t.Consts...)
				info.PDoc.Vars = append(info.PDoc.Vars, t.Vars...)
				info.PDoc.Funcs = append(info.PDoc.Funcs, t.Funcs...)
				t.Consts = nil
				t.Vars = nil
				t.Funcs = nil
			}
			// for now we cannot easily sort consts and vars since
			// go/doc.Value doesn't export the order information
			sort.Sort(funcsByName(info.PDoc.Funcs))
		}

		// collect examples
		testfiles := append(pkginfo.TestGoFiles, pkginfo.XTestGoFiles...)
		files, err = parseFiles(d.fs, fset, relpath, abspath, testfiles)
		if err != nil {
			log.Println("parsing examples:", err)
		}
		info.Examples = collectExamples(pkg, files)
		info.Bugs = info.PDoc.Notes["BUG"]
	}

	info.Dirs = d.root.Lookup(abspath).List(func(path string) bool { return d.includePath(path, mode) })
	info.DirFlat = mode&ModeFlat != 0

	return info
}

func (d *Docs) includePath(path string, mode Mode) (r bool) {
	// if the path includes 'internal', don't list unless we are in the NoFiltering mode.
	if mode&ModeAll != 0 {
		return true
	}
	if strings.Contains(path, "internal") || strings.Contains(path, "vendor") {
		for _, c := range strings.Split(filepath.Clean(path), string(os.PathSeparator)) {
			if c == "internal" || c == "vendor" {
				return false
			}
		}
	}
	return true
}

// simpleImporter returns a (dummy) package object named
// by the last path component of the provided package path
// (as is the convention for packages). This is sufficient
// to resolve package identifiers without doing an actual
// import. It never returns an error.
//
func simpleImporter(imports map[string]*ast.Object, path string) (*ast.Object, error) {
	pkg := imports[path]
	if pkg == nil {
		// note that strings.LastIndex returns -1 if there is no "/"
		pkg = ast.NewObj(ast.Pkg, path[strings.LastIndex(path, "/")+1:])
		pkg.Data = ast.NewScope(nil) // required by ast.NewPackage for dot-import
		imports[path] = pkg
	}
	return pkg, nil
}

// packageExports is a local implementation of ast.PackageExports
// which correctly updates each package file's comment list.
// (The ast.PackageExports signature is frozen, hence the local
// implementation).
//
func packageExports(fset *token.FileSet, pkg *ast.Package) {
	for _, src := range pkg.Files {
		cmap := ast.NewCommentMap(fset, src, src.Comments)
		ast.FileExports(src)
		src.Comments = cmap.Filter(src).Comments()
	}
}

type funcsByName []*doc.Func

func (s funcsByName) Len() int { return len(s) }

func (s funcsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s funcsByName) Less(i, j int) bool { return s[i].Name < s[j].Name }

// collectExamples collects examples for pkg from testfiles.
func collectExamples(pkg *ast.Package, testfiles map[string]*ast.File) []*doc.Example {
	var files []*ast.File
	for _, f := range testfiles {
		files = append(files, f)
	}

	var examples []*doc.Example
	globals := globalNames(pkg)
	for _, e := range doc.Examples(files...) {
		name := TrimExampleSuffix(e.Name)
		if name == "" || globals[name] {
			examples = append(examples, e)
		}
	}

	return examples
}

// globalNames returns a set of the names declared by all package-level
// declarations. Method names are returned in the form Receiver_Method.
func globalNames(pkg *ast.Package) map[string]bool {
	names := make(map[string]bool)
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			addNames(names, decl)
		}
	}
	return names
}

// addNames adds the names declared by decl to the names set.
// Method names are added in the form ReceiverTypeName_Method.
func addNames(names map[string]bool, decl ast.Decl) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		name := d.Name.Name
		if d.Recv != nil {
			var typeName string
			switch r := d.Recv.List[0].Type.(type) {
			case *ast.StarExpr:
				typeName = r.X.(*ast.Ident).Name
			case *ast.Ident:
				typeName = r.Name
			}
			name = typeName + "_" + name
		}
		names[name] = true
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				names[s.Name.Name] = true
			case *ast.ValueSpec:
				for _, id := range s.Names {
					names[id.Name] = true
				}
			}
		}
	}
}

func SplitExampleName(s string) (name, suffix string) {
	i := strings.LastIndex(s, "_")
	if 0 <= i && i < len(s)-1 && !startsWithUppercase(s[i+1:]) {
		name = s[:i]
		suffix = " (" + strings.Title(s[i+1:]) + ")"
		return
	}
	name = s
	return
}

// TrimExampleSuffix strips lowercase braz in Foo_braz or Foo_Bar_braz from name
// while keeping uppercase Braz in Foo_Braz.
func TrimExampleSuffix(name string) string {
	if i := strings.LastIndex(name, "_"); i != -1 {
		if i < len(name)-1 && !startsWithUppercase(name[i+1:]) {
			name = name[:i]
		}
	}
	return name
}

func startsWithUppercase(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}
