// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Wrap wraps long lines intelligently.
//
// Usage:
//
//	go run wrap.go [-w] [file ...]
//
// By default, wrap prints the line-wrapped version of the input files to standard output.
// If no input file is listed, standard input is used.
//
// The -w flag causes wrap to update the files in place, overwriting each with its
// line-wrapped equivalent. Wrap leaves the file unchanged if only a few lines
// need wrapping.
//
// Examples
//
//	go run wrap.go -w your.article
//	go run wrap.go -w *.article
//
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: wrap [-w] [file ...]\n")
	os.Exit(2)
}

var writeBack = flag.Bool("w", false, "write conversions back to original files")
var exitStatus = 0

func main() {
	log.SetPrefix("wrap: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		if *writeBack {
			log.Fatalf("cannot use -w with standard input")
		}
		convert(os.Stdin, "")
		return
	}

	for _, arg := range args {
		f, err := os.Open(arg)
		if err != nil {
			log.Print(err)
			exitStatus = 1
			continue
		}
		target := ""
		if *writeBack {
			target = arg
		}
		err = convert(f, target)
		f.Close()
		if err != nil {
			log.Print(err)
			exitStatus = 1
		}
	}
	os.Exit(exitStatus)
}

// convert reads content from r and line-wraps it.
// If target is the empty string, convert writes the wrapped result to standard output.
// If target is non-empty and there were more than a few line wraps needed,
// convert writes the result to the file named by target.
func convert(r io.Reader, target string) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	changes := 0
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if len(line) < 80 {
			continue
		}
		switch line[0] {
		case '.', '#', ':', '\t', ' ':
			continue
		}
		w := wrap(line)
		if line != w {
			if strings.HasPrefix(line, "- ") {
				w = strings.Replace(w, "\n", "\n  ", -1)
			}
			changes += strings.Count(w, "\n")
			lines[i] = w
		}
	}
	wrapped := []byte(strings.Join(lines, "\n"))

	if target == "" {
		os.Stdout.Write(wrapped)
		return nil
	}

	// Don't rewrite if not much changed.
	if changes < 5 {
		return nil
	}
	return ioutil.WriteFile(target, wrapped, 0666)
}

func wrap(text string) string {
	if len(text) < 120 {
		return text
	}
	// Wrap after sentence boundaries when possible.
	// See https://rhodesmill.org/brandon/2012/one-sentence-per-line/.
	// Note: wrapAt guarantees line[i] == ' '.
	text = wrapAt(text, func(line string, i int) bool {
		return i > 40 && i+20 < len(line) &&
			(line[i-1] == '.' || line[i-1] == '!') && (line[i+1] == ' ' || !('a' <= line[i+1] && line[i+1] <= 'z'))
	})

	// Wrap after phrase boundaries next.
	text = wrapAt(text, func(line string, i int) bool {
		return i > 40 && i+20 < len(line) && (line[i-1] == ';' || line[i-1] == ':')
	})
	text = wrapAt(text, func(line string, i int) bool {
		return i > 40 && i+20 < len(line) && line[i-1] == ','
	})

	// Wrap long lines that are left at spaces.
	text = wrapAt(text, func(line string, i int) bool {
		return i > 70 && i+20 < len(line)
	})

	return text
}

// wrapAt wraps long lines in text, returning the result.
// It calls canWrapAt(line, start, i) to ask whether the given line
// should be wrapped just before offset i; line[i] is known to be a space.
func wrapAt(text string, canWrapAt func(line string, i int) bool) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		var out string
		start := 0
		brackets := 0
		for i := 0; i < len(line); {
			if line[i] == '[' {
				brackets++
			}
			if line[i] == ']' {
				brackets--
			}
			if brackets == 0 && line[i] == ' ' && i+1 < len(line) && line[i+1] != '-' && canWrapAt(line[start:], i-start) {
				out += trimRight(line[start:i]) + "\n"
				i++
				for i < len(line) && line[i] == ' ' {
					i++
				}
				start = i
				continue
			}
			i++
		}
		out += trimRight(line[start:])
		lines[i] = out
	}
	return strings.Join(lines, "\n")
}

func trimRight(s string) string {
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
