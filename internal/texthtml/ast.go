// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package texthtml

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// A goLink describes the (HTML) link information for a Go identifier.
// The zero value of a link represents "no link".
type goLink struct {
	path, name string // package path, identifier name
	isVal      bool   // identifier is defined in a const or var declaration
	oldDocs    bool   // link to ?m=old docs
}

func (l *goLink) tags() (start, end string) {
	switch {
	case l.path != "" && l.name == "":
		// package path
		return `<a href="/pkg/` + l.path + `/` + l.docSuffix() + `">`, `</a>`
	case l.path != "" && l.name != "":
		// qualified identifier
		return `<a href="/pkg/` + l.path + `/` + l.docSuffix() + `#` + l.name + `">`, `</a>`
	case l.path == "" && l.name != "":
		// local identifier
		if l.isVal {
			return `<span id="` + l.name + `">`, `</span>`
		}
		if ast.IsExported(l.name) {
			return `<a href="#` + l.name + `">`, `</a>`
		}
	}
	return "", ""
}

func (l *goLink) docSuffix() string {
	if l.oldDocs {
		return "?m=old"
	}
	return ""
}

// goLinksFor returns the list of links for the identifiers used
// by node in the same order as they appear in the source.
func goLinksFor(node ast.Node) (links []goLink) {
	// linkMap tracks link information for each ast.Ident node. Entries may
	// be created out of source order (for example, when we visit a parent
	// definition node). These links are appended to the returned slice when
	// their ast.Ident nodes are visited.
	linkMap := make(map[*ast.Ident]goLink)

	ast.Inspect(node, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.Field:
			for _, n := range n.Names {
				linkMap[n] = goLink{}
			}
		case *ast.ImportSpec:
			if name := n.Name; name != nil {
				linkMap[name] = goLink{}
			}
		case *ast.ValueSpec:
			for _, n := range n.Names {
				linkMap[n] = goLink{name: n.Name, isVal: true}
			}
		case *ast.FuncDecl:
			linkMap[n.Name] = goLink{}
		case *ast.TypeSpec:
			linkMap[n.Name] = goLink{}
		case *ast.AssignStmt:
			// Short variable declarations only show up if we apply
			// this code to all source code (as opposed to exported
			// declarations only).
			if n.Tok == token.DEFINE {
				// Some of the lhs variables may be re-declared,
				// so technically they are not defs. We don't
				// care for now.
				for _, x := range n.Lhs {
					// Each lhs expression should be an
					// ident, but we are conservative and check.
					if n, _ := x.(*ast.Ident); n != nil {
						linkMap[n] = goLink{isVal: true}
					}
				}
			}
		case *ast.SelectorExpr:
			// Detect qualified identifiers of the form pkg.ident.
			// If anything fails we return true and collect individual
			// identifiers instead.
			if x, _ := n.X.(*ast.Ident); x != nil {
				// Create links only if x is a qualified identifier.
				if obj := x.Obj; obj != nil && obj.Kind == ast.Pkg {
					if spec, _ := obj.Decl.(*ast.ImportSpec); spec != nil {
						// spec.Path.Value is the import path
						if path, err := strconv.Unquote(spec.Path.Value); err == nil {
							// Register two links, one for the package
							// and one for the qualified identifier.
							linkMap[x] = goLink{path: path}
							linkMap[n.Sel] = goLink{path: path, name: n.Sel.Name}
						}
					}
				}
			}
		case *ast.CompositeLit:
			// Detect field names within composite literals. These links should
			// be prefixed by the type name.
			fieldPath := ""
			prefix := ""
			switch typ := n.Type.(type) {
			case *ast.Ident:
				prefix = typ.Name + "."
			case *ast.SelectorExpr:
				if x, _ := typ.X.(*ast.Ident); x != nil {
					// Create links only if x is a qualified identifier.
					if obj := x.Obj; obj != nil && obj.Kind == ast.Pkg {
						if spec, _ := obj.Decl.(*ast.ImportSpec); spec != nil {
							// spec.Path.Value is the import path
							if path, err := strconv.Unquote(spec.Path.Value); err == nil {
								// Register two links, one for the package
								// and one for the qualified identifier.
								linkMap[x] = goLink{path: path}
								linkMap[typ.Sel] = goLink{path: path, name: typ.Sel.Name}
								fieldPath = path
								prefix = typ.Sel.Name + "."
							}
						}
					}
				}
			}
			for _, e := range n.Elts {
				if kv, ok := e.(*ast.KeyValueExpr); ok {
					if k, ok := kv.Key.(*ast.Ident); ok {
						// Note: there is some syntactic ambiguity here. We cannot determine
						// if this is a struct literal or a map literal without type
						// information. We assume struct literal.
						name := prefix + k.Name
						linkMap[k] = goLink{path: fieldPath, name: name}
					}
				}
			}
		case *ast.Ident:
			if l, ok := linkMap[n]; ok {
				links = append(links, l)
			} else {
				l := goLink{name: n.Name}
				if n.Obj == nil && doc.IsPredeclared(n.Name) {
					l.path = "builtin"
				}
				links = append(links, l)
			}
		}
		return true
	})
	return
}

// postFormatAST makes any appropriate changes to the formatting of node in buf.
// Specifically, it adds span links to each struct field, so they can be linked properly.
// TODO(rsc): Why not do this as part of the linking above?
func postFormatAST(buf *bytes.Buffer, node ast.Node) {
	if st, name := isStructTypeDecl(node); st != nil {
		addStructFieldIDAttributes(buf, name, st)
	}
}

// isStructTypeDecl checks whether n is a struct declaration.
// It either returns a non-nil StructType and its name, or zero values.
func isStructTypeDecl(n ast.Node) (st *ast.StructType, name string) {
	gd, ok := n.(*ast.GenDecl)
	if !ok || gd.Tok != token.TYPE {
		return nil, ""
	}
	if gd.Lparen > 0 {
		// Parenthesized type. Who does that, anyway?
		// TODO: Reportedly gri does. Fix this to handle that too.
		return nil, ""
	}
	if len(gd.Specs) != 1 {
		return nil, ""
	}
	ts, ok := gd.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil, ""
	}
	st, ok = ts.Type.(*ast.StructType)
	if !ok {
		return nil, ""
	}
	return st, ts.Name.Name
}

// addStructFieldIDAttributes modifies the contents of buf such that
// all struct fields of the named struct have <span id='name.Field'>
// in them, so people can link to /#Struct.Field.
func addStructFieldIDAttributes(buf *bytes.Buffer, name string, st *ast.StructType) {
	if st.Fields == nil {
		return
	}
	// needsLink is a set of identifiers that still need to be
	// linked, where value == key, to avoid an allocation in func
	// linkedField.
	needsLink := make(map[string]string)

	for _, f := range st.Fields.List {
		if len(f.Names) == 0 {
			continue
		}
		fieldName := f.Names[0].Name
		needsLink[fieldName] = fieldName
	}
	var newBuf bytes.Buffer
	foreachLine(buf.Bytes(), func(line []byte) {
		if fieldName := linkedField(line, needsLink); fieldName != "" {
			fmt.Fprintf(&newBuf, `<span id="%s.%s"></span>`, name, fieldName)
			delete(needsLink, fieldName)
		}
		newBuf.Write(line)
	})
	buf.Reset()
	buf.Write(newBuf.Bytes())
}

// foreachLine calls fn for each line of in, where a line includes
// the trailing "\n", except on the last line, if it doesn't exist.
func foreachLine(in []byte, fn func(line []byte)) {
	for len(in) > 0 {
		nl := bytes.IndexByte(in, '\n')
		if nl == -1 {
			fn(in)
			return
		}
		fn(in[:nl+1])
		in = in[nl+1:]
	}
}

// commentPrefix is the line prefix for comments after they've been HTMLified.
var commentPrefix = []byte(`<span class="comment">// `)

// linkedField determines whether the given line starts with an
// identifier in the provided ids map (mapping from identifier to the
// same identifier). The line can start with either an identifier or
// an identifier in a comment. If one matches, it returns the
// identifier that matched. Otherwise it returns the empty string.
func linkedField(line []byte, ids map[string]string) string {
	line = bytes.TrimSpace(line)

	// For fields with a doc string of the
	// conventional form, we put the new span into
	// the comment instead of the field.
	// The "conventional" form is a complete sentence
	// per https://golang.org/s/style#comment-sentences like:
	//
	//    // Foo is an optional Fooer to foo the foos.
	//    Foo Fooer
	//
	// In this case, we want the #StructName.Foo
	// link to make the browser go to the comment
	// line "Foo is an optional Fooer" instead of
	// the "Foo Fooer" line, which could otherwise
	// obscure the docs above the browser's "fold".
	//
	// TODO: do this better, so it works for all
	// comments, including unconventional ones.
	line = bytes.TrimPrefix(line, commentPrefix)
	id := scanIdentifier(line)
	if len(id) == 0 {
		// No leading identifier. Avoid map lookup for
		// somewhat common case.
		return ""
	}
	return ids[string(id)]
}

// scanIdentifier scans a valid Go identifier off the front of v and
// either returns a subslice of v if there's a valid identifier, or
// returns a zero-length slice.
func scanIdentifier(v []byte) []byte {
	var n int // number of leading bytes of v belonging to an identifier
	for {
		r, width := utf8.DecodeRune(v[n:])
		if !(isLetter(r) || n > 0 && isDigit(r)) {
			break
		}
		n += width
	}
	return v[:n]
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}
