// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webtest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

func TestWebtestHandler(t *testing.T) {
	h := http.FileServer(http.Dir("testdata"))
	testWebtest(t, "testdata/fs*.txt", func(c *case_) error { return c.runHandler(h) })
}

func echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%v %v\n", r.Method, r.RequestURI)
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "parsing form: %v\n", err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "%q: %q\n", k, v)
	}
	if len(r.Form) == 0 {
		fmt.Fprintf(w, "no query\n")
	}
}

func TestEchoHandler(t *testing.T) {
	TestHandler(t, "testdata/echo.txt", http.HandlerFunc(echo))
}

func testWebtest(t *testing.T, glob string, do func(*case_) error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			script, err := parseScript(file, string(data))
			if err != nil {
				t.Fatal(err)
			}
			for _, c := range script.cases {
				t.Run(c.method+"/"+strings.TrimPrefix(c.url, "/"), func(t *testing.T) {
					hint := c.hint
					c.hint = ""
					if err := do(c); err != nil {
						if hint == "" {
							t.Fatal(err)
						}
						if !strings.Contains(err.Error(), hint) {
							t.Fatalf("unexpected error %v (want %q)", err, hint)
						}
						return
					}
					if hint != "" {
						t.Fatalf("unexpected success (want %q)", hint)
					}
				})
			}
		})
	}
}
