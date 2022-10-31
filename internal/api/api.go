// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file caches information about which standard library types, methods,
// and functions appeared in what version of Go

package api

import (
	"bufio"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// DB is a map of packages to information about those packages'
// symbols and when they were added to Go.
//
// Only things added after Go1 are tracked. Version strings are of the
// form "1.1", "1.2", etc.
type DB map[string]PkgDB // keyed by Go package ("net/http")

// PkgDB contains information about which version of Go added
// certain package symbols.
//
// Only things added after Go1 are tracked. Version strings are of the
// form "1.1", "1.2", etc.
type PkgDB struct {
	Type   map[string]string            // "Server" -> "1.7"
	Method map[string]map[string]string // "*Server" ->"Shutdown"->1.8
	Func   map[string]string            // "NewServer" -> "1.7"
	Field  map[string]map[string]string // "ClientTrace" -> "Got1xxResponse" -> "1.11"
}

// Func returns a string (such as "1.7") specifying which Go
// version introduced a symbol, unless it was introduced in Go1, in
// which case it returns the empty string.
//
// The kind is one of "type", "method", or "func".
//
// The receiver is only used for "methods" and specifies the receiver type,
// such as "*Server".
//
// The name is the symbol name ("Server") and the pkg is the package
// ("net/http").
func (v DB) Func(pkg, kind, receiver, name string) string {
	pv := v[pkg]
	switch kind {
	case "func":
		return pv.Func[name]
	case "type":
		return pv.Type[name]
	case "method":
		return pv.Method[receiver][name]
	}
	return ""
}

// Load loads a database from fsys's api/go*.txt files.
// Typically, fsys should be the root of a Go repository (a $GOROOT).
func Load(fsys fs.FS) (DB, error) {
	files, err := fs.Glob(fsys, "api/go*.txt")
	if err != nil {
		return nil, err
	}

	// Process files in go1.n, go1.n-1, ..., go1.2, go1.1, go1 order.
	//
	// It's rare, but the signature of an identifier may change
	// (for example, a function that accepts a type replaced with
	// an alias), and so an existing symbol may show up again in
	// a later api/go1.N.txt file. Parsing in reverse version
	// order means we end up with the earliest version of Go
	// when the symbol was added. See golang.org/issue/44081.
	//
	ver := func(name string) int {
		base := path.Base(name)
		ver := strings.TrimPrefix(strings.TrimSuffix(base, ".txt"), "go1.")
		if ver == "go1" {
			return 0
		}
		v, _ := strconv.Atoi(ver)
		return v
	}
	sort.Slice(files, func(i, j int) bool { return ver(files[i]) > ver(files[j]) })
	vp := new(parser)
	for _, f := range files {
		if err := vp.parseFile(fsys, f); err != nil {
			return nil, err
		}
	}
	return vp.res, nil
}

// parser parses $GOROOT/api/go*.txt files and stores them in its rows field.
type parser struct {
	res DB // initialized lazily
}

// parseFile parses the named $GOROOT/api/goVERSION.txt file.
//
// For each row, it updates the corresponding entry in
// vp.res to VERSION, overwriting any previous value.
// As a special case, if goVERSION is "go1", it deletes
// from the map instead.
func (vp *parser) parseFile(fsys fs.FS, name string) error {
	f, err := fsys.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	base := filepath.Base(name)
	ver := strings.TrimPrefix(strings.TrimSuffix(base, ".txt"), "go")

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		row, ok := parseRow(sc.Text())
		if !ok {
			continue
		}
		if vp.res == nil {
			vp.res = make(DB)
		}
		pkgi, ok := vp.res[row.pkg]
		if !ok {
			pkgi = PkgDB{
				Type:   make(map[string]string),
				Method: make(map[string]map[string]string),
				Func:   make(map[string]string),
				Field:  make(map[string]map[string]string),
			}
			vp.res[row.pkg] = pkgi
		}
		switch row.kind {
		case "func":
			if ver == "1" {
				delete(pkgi.Func, row.name)
				break
			}
			pkgi.Func[row.name] = ver
		case "type":
			if ver == "1" {
				delete(pkgi.Type, row.name)
				break
			}
			pkgi.Type[row.name] = ver
		case "method":
			if ver == "1" {
				delete(pkgi.Method[row.recv], row.name)
				break
			}
			if _, ok := pkgi.Method[row.recv]; !ok {
				pkgi.Method[row.recv] = make(map[string]string)
			}
			pkgi.Method[row.recv][row.name] = ver
		case "field":
			if ver == "1" {
				delete(pkgi.Field[row.structName], row.name)
				break
			}
			if _, ok := pkgi.Field[row.structName]; !ok {
				pkgi.Field[row.structName] = make(map[string]string)
			}
			pkgi.Field[row.structName][row.name] = ver
		}
	}
	return sc.Err()
}

// row represents an API feature, a parsed line of a
// $GOROOT/api/go.*txt file.
type row struct {
	pkg        string // "net/http"
	kind       string // "type", "func", "method", "field" TODO: "const", "var"
	recv       string // for methods, the receiver type ("Server", "*Server")
	name       string // name of type, (struct) field, func, method
	structName string // for struct fields, the outer struct name
}

func parseRow(s string) (vr row, ok bool) {
	if !strings.HasPrefix(s, "pkg ") {
		// Skip comments, blank lines, etc.
		return
	}
	rest := s[len("pkg "):]
	endPkg := strings.IndexFunc(rest, func(r rune) bool { return !(unicode.IsLetter(r) || r == '/' || unicode.IsDigit(r)) })
	if endPkg == -1 {
		return
	}
	vr.pkg, rest = rest[:endPkg], rest[endPkg:]
	if !strings.HasPrefix(rest, ", ") {
		// If the part after the pkg name isn't ", ", then it's a OS/ARCH-dependent line of the form:
		//   pkg syscall (darwin-amd64), const ImplementsGetwd = false
		// We skip those for now.
		return
	}
	rest = rest[len(", "):]

	switch {
	case strings.HasPrefix(rest, "type "):
		rest = rest[len("type "):]
		sp := strings.IndexByte(rest, ' ')
		if sp == -1 {
			return
		}
		vr.name, rest = rest[:sp], rest[sp+1:]
		if !strings.HasPrefix(rest, "struct, ") {
			vr.kind = "type"
			return vr, true
		}
		rest = rest[len("struct, "):]
		if i := strings.IndexByte(rest, ' '); i != -1 {
			vr.kind = "field"
			vr.structName = vr.name
			vr.name = rest[:i]
			return vr, true
		}
	case strings.HasPrefix(rest, "func "):
		vr.kind = "func"
		rest = rest[len("func "):]
		if i := strings.IndexByte(rest, '('); i != -1 {
			vr.name = rest[:i]
			return vr, true
		}
	case strings.HasPrefix(rest, "method "): // "method (*File) SetModTime(time.Time)"
		vr.kind = "method"
		rest = rest[len("method "):] // "(*File) SetModTime(time.Time)"
		sp := strings.IndexByte(rest, ' ')
		if sp == -1 {
			return
		}
		vr.recv = strings.Trim(rest[:sp], "()") // "*File"
		rest = rest[sp+1:]                      // SetMode(os.FileMode)
		paren := strings.IndexByte(rest, '(')
		if paren == -1 {
			return
		}
		vr.name = rest[:paren]
		return vr, true
	}
	return // TODO: handle more cases
}
