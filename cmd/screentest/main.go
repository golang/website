// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Screentest compares images of rendered web pages.
It compares images obtained from two sources, one to test and one for the expected result.
The comparisons are driven by a script file in a format described below.

# Usage

	screentest [flags] [glob]

The flags are:

	-test URL
		  URL or path being tested. Required.
	-want URL
		  URL or path for expected results. Required.
	-c
		  Number of test cases to run concurrently.
	-d
		  URL of a Chrome websocket debugger. If omitted, screentest tries to find the
		  Chrome executable on the system and starts a new instance.
	-headers
		  HTTP(S) headers to send with each request, as a comma-separated list of name:value.
	-run REGEXP
		  Run only tests matching regexp.
	-o
		  URL or path for output files. If omitted, files are written to a subdirectory of the
		  user's cache directory.
	-u
		  Update cached screenshots.
	-v
		  Variables provided to script templates as comma separated KEY:VALUE pairs.

# Scripts

A script file contains one or more test cases described as a sequence of lines. The
file is first processed as Go template using the text/template package, with a map
of the variables given by the -v flag provided as `.`.

The script format is line-oriented.
Lines beginning with # characters are ignored as comments.
Each non-blank, non-comment line is a directive, listed below.

Each test case begins with the 'test' directive and ends with a blank line.
A test case describes actions to take on a page, along
with the dimensions of the screenshots to be compared. For example, here is
a trivial script:

	test about
	pathname /about
	capture fullscreen

This script has a single test case. The first line names the test.
The second line sets the page to visit at each origin. The last line
captures full-page screenshots of the pages and generates a diff image if they
do not match.

# Directives

Use windowsize WIDTHxHEIGHT to set the default window size for all test cases
that follow.

	windowsize 540x1080

Use block URL ... to set URL patterns to block. Wildcards ('*') are allowed.

	block https://codecov.io/* https://travis-ci.com/*

The directives above apply to all test cases that follow.
The ones below must appear inside a test case and apply only to that case.

Use test NAME to create a name for the test case.

	test about page

Use pathname PATH to set the page to visit at each origin.

	pathname /about

Use status CODE to set an expected HTTP status code. The default is 200.

	status 404

Use click SELECTOR to add a click an element on the page.

	click button.submit

Use wait SELECTOR to wait for an element to appear.

	wait [role="treeitem"][aria-expanded="true"]

Use capture [SIZE] [ARG] to create a test case with the properties
defined in the test case. If present, the first argument to capture should be one of
'fullscreen', 'viewport' or 'element'.

	capture fullscreen 540x1080

When taking an element screenshot provide a selector.

	capture element header

Use eval JS to evaluate JavaScript snippets to hide elements or prepare the page in
some other way.

	eval 'document.querySelector(".selector").remove();'
	eval 'window.scrollTo({top: 0});'

Each capture command to creates a new test case for a single page.

	windowsize 1536x960

	test homepage
	pathname /
	capture viewport
	capture viewport 540x1080
	capture viewport 400x1000

	test about page
	pathname /about
	capture viewport
	capture viewport 540x1080
	capture viewport 400x1000
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var flags options

func init() {
	flag.StringVar(&flags.testURL, "test", "", "URL or file path to test")
	flag.StringVar(&flags.wantURL, "want", "", "URL or file path with expected results")
	flag.BoolVar(&flags.update, "u", false, "update cached screenshots")
	flag.StringVar(&flags.vars, "v", "", "variables provided to script templates as comma separated KEY:VALUE pairs")
	flag.IntVar(&flags.maxConcurrency, "c", (runtime.NumCPU()+1)/2, "number of test cases to run concurrently")
	flag.StringVar(&flags.debuggerURL, "d", "", "chrome debugger URL")
	flag.StringVar(&flags.outputURL, "o", "", "path for output: file path or URL with 'file' or 'gs' scheme")
	flag.StringVar(&flags.headers, "headers", "", "HTTP headers: comma-separated list of name:value")
	flag.StringVar(&flags.run, "run", "", "regexp to match test")
}

// options are the options for the program.
// See the top command and the flag.XXXVar calls above for documentation.
type options struct {
	testURL        string
	wantURL        string
	update         bool
	vars           string
	maxConcurrency int
	debuggerURL    string
	run            string
	outputURL      string
	headers        string
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: screentest [flags] [glob]\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	// Require testdata glob when invoked as an installed command.
	if len(args) != 1 && os.Args[0] == "screentest" {
		flag.Usage()
		os.Exit(1)
	}
	glob := filepath.Join("cmd", "screentest", "testdata", "*")
	if len(args) == 1 {
		glob = args[0]
	}

	if err := run(context.Background(), glob, flags); err != nil {
		log.Fatal(err)
	}
}
