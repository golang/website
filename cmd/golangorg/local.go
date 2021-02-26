// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16 && !prod
// +build go1.16,!prod

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	// This package registers "/compile" and "/share" handlers
	// that redirect to the golang.org playground.
	_ "golang.org/x/tools/playground"
)

func earlySetup() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintln(os.Stderr, "runtime.Caller failed: cannot find templates for -a mode.")
		os.Exit(2)
	}
	dir := filepath.Join(file, "../../../_content")
	if _, err := os.Stat(filepath.Join(dir, "lib/godoc/site.html")); err != nil {
		log.Printf("warning: cannot find template dir; using embedded copy")
		return
	}
	*templateDir = dir
}

func lateSetup(mux *http.ServeMux) {
	// Register a redirect handler for /dl/ to the golang.org download page.
	mux.Handle("/dl/", http.RedirectHandler("https://golang.org/dl/", http.StatusFound))
}
