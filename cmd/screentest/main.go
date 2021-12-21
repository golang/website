// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command screentest runs the visual diff check for the set of scripts
// provided by the flag -testdata.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/website/internal/screentest"
)

var (
	update  = flag.Bool("update", false, "update cached screenshots")
	headers = flag.String("H", "", "set request headers")
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: screentest [OPTIONS] glob\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	hdr := make(map[string]interface{})
	if *headers != "" {
		for _, h := range strings.Split(*headers, ",") {
			parts := strings.Split(h, ":")
			if len(parts) != 2 {
				log.Fatalf("invalid header %s", h)
			}
			hdr[parts[0]] = parts[1]
		}
	}
	if err := screentest.CheckHandler(args[0], *update, hdr); err != nil {
		log.Fatal(err)
	}
}
