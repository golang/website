// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testHosts = map[string]string{
	"":                     "foo.example.test",
	"dev.example.test":     "foo-dev.example.test",
	"staging.example.test": "foo-staging.example.test",
}

func TestRedirect(t *testing.T) {
	tests := []struct {
		desc     string
		target   string
		hosts    map[string]string
		want     string
		wantCode int
	}{
		{
			desc:     "basic redirect",
			target:   "/",
			hosts:    testHosts,
			want:     "https://foo.example.test/",
			wantCode: http.StatusFound,
		},
		{
			desc:     "redirect keeps query and path",
			target:   "/github.com/golang/glog?tab=overview",
			hosts:    testHosts,
			want:     "https://foo.example.test/github.com/golang/glog?tab=overview",
			wantCode: http.StatusFound,
		},
		{
			desc:     "redirects to the correct host",
			target:   "https://dev.example.test/",
			hosts:    testHosts,
			want:     "https://foo-dev.example.test/",
			wantCode: http.StatusFound,
		},
		{
			desc:     "renders 404 if hosts are missing",
			target:   "https://dev.example.test/",
			hosts:    nil,
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.target, nil)
			w := httptest.NewRecorder()
			redirectHosts(tt.hosts).ServeHTTP(w, req)
			resp := w.Result()
			if resp.StatusCode != tt.wantCode {
				t.Errorf("w.Result().StatusCode = %v, wanted %v", resp.StatusCode, tt.wantCode)
			}
			l, err := resp.Location()
			if resp.StatusCode == http.StatusFound && (l == nil || l.String() != tt.want || err != nil) {
				t.Errorf("resp.Location() = %v, %v, wanted %v, no error", l, err, tt.want)
			}
		})
	}
}

var siteTests = []struct {
	target string
	want   []string
}{
	{"/", []string{"Go is an open source programming language supported by Google"}},
	{"/solutions/", []string{"Using Go at Google"}},
	{"/solutions/dropbox", []string{"About Dropbox"}},
}

func TestSite(t *testing.T) {
	h, err := NewHandler("../../_content")
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
