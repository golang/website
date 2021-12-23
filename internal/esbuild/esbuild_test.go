package esbuild

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/website/internal/web"
)

func TestServeHTTP(t *testing.T) {
	exampleOut := `/**
 * @license
 * Copyright 2021 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
function sayHello(to) {
  console.log("Hello, " + to + "!");
}
const world = {
  name: "World",
  toString() {
    return this.name;
  }
};
sayHello(world);
`
	tests := []struct {
		name            string
		path            string
		wantCode        int
		wantBody        string
		wantCacheHeader bool
	}{
		{
			name:     "example code",
			path:     "/example.ts",
			wantCode: 200,
			wantBody: exampleOut,
		},
		{
			name:            "example code cached",
			path:            "/example.ts",
			wantCode:        200,
			wantBody:        exampleOut,
			wantCacheHeader: true,
		},
		{
			name:     "file not found",
			path:     "/notfound.ts",
			wantCode: 500,
			wantBody: fmt.Sprintf("\n\nopen testdata/notfound.ts: %s\n", syscall.ENOENT),
		},
		{
			name:     "syntax error",
			path:     "/error.ts",
			wantCode: 500,
			wantBody: "\n\nExpected identifier but found &#34;function&#34;\n\n",
		},
	}
	fsys := os.DirFS("testdata")
	server := NewServer(fsys, web.NewSite(fsys))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			got := httptest.NewRecorder()
			server.ServeHTTP(got, req)
			gotHeader := got.Header().Get(cacheHeader) == "true"
			if got.Code != tt.wantCode {
				t.Errorf("got status %d but wanted %d", got.Code, http.StatusOK)
			}
			if (tt.wantCacheHeader && !gotHeader) || (!tt.wantCacheHeader && gotHeader) {
				t.Errorf("got cache hit %v but wanted %v", gotHeader, tt.wantCacheHeader)
			}
			if diff := cmp.Diff(tt.wantBody, got.Body.String()); diff != "" {
				t.Errorf("ServeHTTP() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
