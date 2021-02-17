// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

Godoc extracts and generates documentation for Go programs.

It runs as a web server and presents the documentation as a
web page.

	godoc -http=:6060

Usage:

	godoc [flag]

The flags are:

	-v
		verbose mode
	-timestamps=true
		show timestamps with directory listings
	-play=false
		enable playground
	-links=true
		link identifiers to their declarations
	-notes="BUG"
		regular expression matching note markers to show
		(e.g., "BUG|TODO", ".*")
	-goroot=$GOROOT
		Go root directory
	-http=addr
		HTTP service address (e.g., '127.0.0.1:6060' or just ':6060')
	-templates=""
		directory containing alternate template files; if set,
		the directory may provide alternative template files
		for the files in _content/

By default, golangorg looks at the packages it finds via $GOROOT (if set).
This behavior can be altered by providing an alternative $GOROOT with the -goroot
flag.

By default, godoc uses the system's GOOS/GOARCH. You can provide the URL parameters
"GOOS" and "GOARCH" to set the output on the web page for the target system.

The presentation mode of web pages served by godoc can be controlled with the
"m" URL parameter; it accepts a comma-separated list of flag names as value:

	all	show documentation for all declarations, not just the exported ones
	methods	show all embedded methods, not just those of unexported anonymous fields
	src	show the original source code rather than the extracted documentation
	flat	present flat (not indented) directory listings using full paths

For instance, https://golang.org/pkg/math/big/?m=all shows the documentation
for all (not just the exported) declarations of package big.

Godoc serves files from the file system of the underlying OS.

Godoc documentation is converted to HTML or to text using the go/doc package;
see https://golang.org/pkg/go/doc/#ToHTML for the exact rules.
Godoc also shows example code that is runnable by the testing package;
see https://golang.org/pkg/testing/#hdr-Examples for the conventions.
See "Godoc: documenting Go code" for how to write good comments for godoc:
https://golang.org/doc/articles/godoc_documenting_go_code.html

*/
package main // import "golang.org/x/website/cmd/golangorg"
