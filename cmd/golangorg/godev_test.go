// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"
)

var siteTests = []struct {
	target string
	want   []string
}{
	{"/", []string{"Go is an open source programming language supported by Google"}},
	{"/solutions/", []string{"Using Go at Google"}},
	{"/solutions/dropbox", []string{"About Dropbox"}},
}

func TestSite(t *testing.T) {
	h, err := godevHandler("../../go.dev/_content")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range siteTests {
		t.Run(tt.target, func(t *testing.T) {
			r := httptest.NewRequest("GET", tt.target, nil)
			resp := httptest.NewRecorder()
			resp.Body = new(bytes.Buffer)
			h.ServeHTTP(resp, r)
			if resp.Code != 200 {
				t.Fatalf("Code = %d, want 200", resp.Code)
			}
			body := resp.Body.String()
			for _, str := range tt.want {
				if !strings.Contains(body, str) {
					t.Fatalf("Body does not contain %q:\n%s", str, body)
				}
			}
		})
	}
}
