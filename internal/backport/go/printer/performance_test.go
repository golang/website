// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements a simple printer performance benchmark:
// go test -bench=BenchmarkPrint

package printer

import (
	"bytes"
	"golang.org/x/website/internal/backport/go/ast"
	"golang.org/x/website/internal/backport/go/parser"
	"io"
	"log"
	"os"
	"testing"
)

var (
	testfile *ast.File
	testsize int64
)

func testprint(out io.Writer, file *ast.File) {
	if err := (&Config{TabIndent | UseSpaces | normalizeNumbers, 8, 0}).Fprint(out, fset, file); err != nil {
		log.Fatalf("print error: %s", err)
	}
}

// cannot initialize in init because (printer) Fprint launches goroutines.
func initialize() {
	const filename = "testdata/parser.go"

	src, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("%s", err)
	}

	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		log.Fatalf("%s", err)
	}

	var buf bytes.Buffer
	testprint(&buf, file)
	if !bytes.Equal(buf.Bytes(), src) {
		log.Fatalf("print error: %s not idempotent", filename)
	}

	testfile = file
	testsize = int64(len(src))
}

func BenchmarkPrint(b *testing.B) {
	if testfile == nil {
		initialize()
	}
	b.ReportAllocs()
	b.SetBytes(testsize)
	for i := 0; i < b.N; i++ {
		testprint(io.Discard, testfile)
	}
}
