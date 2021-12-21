// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command screentest runs the screentest check for a set of scripts.
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
	update = flag.Bool("update", false, "update cached screenshots")
	vars   = flag.String("vars", "", "provide variables to the script template as comma separated KEY:VALUE pairs")
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
	parsedVars := make(map[string]string)
	if *vars != "" {
		for _, pair := range strings.Split(*vars, ",") {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) != 2 {
				log.Fatal(fmt.Errorf("invalid key value pair, %q", pair))
			}
			parsedVars[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	if err := screentest.CheckHandler(args[0], *update, parsedVars); err != nil {
		log.Fatal(err)
	}
}
