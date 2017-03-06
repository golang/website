// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"strings"
)

var strip = flag.Bool("strip", false, "strip included files")

func main() {
	flag.Parse()
	for _, arg := range flag.Args() {
		do(arg)
	}
}

func do(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	var out []byte
	skip := false
	for _, line := range bytes.SplitAfter(data, []byte("\n")) {
		if skip && bytes.HasPrefix(line, []byte("<!--end")) {
			skip = false
		}
		if skip {
			continue
		}
		out = append(out, line...)
		if bytes.HasPrefix(line, []byte("<!--include")) {
			if !*strip {
				more, err := ioutil.ReadFile(strings.Fields(string(line))[1])
				if err != nil {
					log.Fatal(err)
				}
				if bytes.HasPrefix(more, xmlHeader) {
					more = more[len(xmlHeader):]
				}
				if len(more) > 0 && more[len(more)-1] != '\n' {
					more = append(more, '\n')
				}
				out = append(out, more...)
			}
			skip = true
		}
	}

	if err := ioutil.WriteFile(file, out, 0666); err != nil {
		log.Fatal(err)
	}
}

var xmlHeader = []byte(`<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" 
  "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
`)
