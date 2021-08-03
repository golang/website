// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"rsc.io/rf/diff"
)

func TestGolden(t *testing.T) {
	start := time.Now()
	h, err := godevHandler("../../go.dev/_content")
	if err != nil {
		t.Fatal(err)
	}
	total := time.Since(start)
	t.Logf("Load %v\n", total)

	root := "../../go.dev/testdata/golden"
	err = filepath.Walk(root, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		name = filepath.ToSlash(name[len(root)+1:])
		switch name {
		case "index.xml",
			"categories/index.html",
			"categories/index.xml",
			"learn/index.xml",
			"series/index.html",
			"series/index.xml",
			"series/case-studies/index.html",
			"series/case-studies/index.xml",
			"series/use-cases/index.html",
			"series/use-cases/index.xml",
			"sitemap.xml",
			"solutions/google/index.xml",
			"solutions/index.xml",
			"tags/index.html",
			"tags/index.xml":
			t.Logf("%s <- SKIP\n", name)
			return nil
		}

		want, err := ioutil.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}

		start := time.Now()
		r := httptest.NewRequest("GET", "/"+name, nil)
		resp := httptest.NewRecorder()
		resp.Body = new(bytes.Buffer)
		h.ServeHTTP(resp, r)
		for nredir := 0; resp.Code/10 == 30; nredir++ {
			if nredir > 10 {
				t.Fatalf("%s <- redirect loop!", name)
			}
			r.URL.Path = resp.Result().Header.Get("Location")
			resp = httptest.NewRecorder()
			resp.Body = new(bytes.Buffer)
			h.ServeHTTP(resp, r)
		}
		if resp.Code != 200 {
			t.Fatalf("GET %s <- %d, want 200", r.URL, resp.Code)
		}
		have := resp.Body.Bytes()
		total += time.Since(start)

		if path.Ext(name) == ".html" {
			have = canonicalize(have)
			want = canonicalize(want)
			if !bytes.Equal(have, want) {
				d, err := diff.Diff("hugo", want, "newgo", have)
				if err != nil {
					panic(err)
				}
				t.Fatalf("%s: diff:\n%s", name, d)
			}
			t.Logf("%s <- OK!\n", name)
			return nil
		}

		if !bytes.Equal(have, want) {
			t.Fatalf("%s: wrong bytes", name)
		}
		return nil
	})
	t.Logf("total %v", total)

	if err != nil {
		t.Fatal(err)
	}
}

// canonicalize trims trailing spaces and tabs at the ends of lines,
// removes blank lines, and removes leading spaces before HTML tags.
// This gives us a little more leeway in cases where it is difficult
// to match Hugo's whitespace heuristics exactly or where we are
// refactoring templates a little which changes spacing in inconsequential ways.
func canonicalize(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte("<li>"), []byte("<li>\n"))
	data = bytes.ReplaceAll(data, []byte("</p>"), []byte("</p>\n"))
	data = bytes.ReplaceAll(data, []byte("</ul>"), []byte("</ul>\n"))
	data = regexp.MustCompile(`(<(img|hr)([^<>]*[^ <>])?) */>`).ReplaceAll(data, []byte("$1>")) // <img/> to <img>

	lines := bytes.Split(data, []byte("\n"))
	for i, line := range lines {
		lines[i] = bytes.Trim(line, " \t")
	}
	var out [][]byte
	for _, line := range lines {
		if len(line) > 0 {
			out = append(out, line)
		}
	}
	return bytes.Join(out, []byte("\n"))
}
