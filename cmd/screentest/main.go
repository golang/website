// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command screentest runs the visual diff check for the set of scripts
// provided by the flag -testdata.
package main

import (
	"flag"
	"log"

	"golang.org/x/website/internal/screentest"
)

var (
	testdata = flag.String("testdata", "cmd/screentest/testdata/*.txt", "directory to look for testdata")
	update   = flag.Bool("update", false, "use this flag to update cached screenshots")
)

func main() {
	flag.Parse()
	if err := screentest.CheckHandler(*testdata, *update); err != nil {
		log.Fatal(err)
	}
}
