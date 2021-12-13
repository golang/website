// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package screentest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
)

func TestReadTests(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				filename: "testdata/readtests.txt",
			},
			want: []*testcase{
				{
					name:           "go.dev homepage",
					originA:        "https://go.dev",
					originB:        "http://localhost:6060/go.dev",
					viewportWidth:  1536,
					viewportHeight: 960,
					screenshotType: fullScreenshot,
					pathame:        "/",
				},
				{
					name:           "go.dev homepage 540x1080",
					originA:        "https://go.dev",
					originB:        "http://localhost:6060/go.dev",
					viewportWidth:  540,
					viewportHeight: 1080,
					screenshotType: fullScreenshot,
					pathame:        "/",
				},
				{
					name:           "about page",
					originA:        "https://go.dev",
					originB:        "http://localhost:6060/go.dev",
					screenshotType: fullScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
					pathame:        "/about",
				},
				{
					name:              "pkg.go.dev homepage .go-Carousel",
					originA:           "https://pkg.go.dev",
					originB:           "https://beta.pkg.go.dev",
					screenshotType:    elementScreenshot,
					screenshotElement: ".go-Carousel",
					viewportWidth:     1536,
					viewportHeight:    960,
					pathame:           "/",
					tasks: chromedp.Tasks{
						chromedp.Click(".go-Carousel-dot"),
					},
				},
				{
					name:           "net package doc",
					originA:        "https://pkg.go.dev",
					originB:        "https://beta.pkg.go.dev",
					screenshotType: viewportScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
					pathame:        "/net",
					tasks: chromedp.Tasks{
						chromedp.WaitReady(`[role="treeitem"][aria-expanded="true"]`),
					},
				},
				{
					name:           "net package doc 540x1080",
					originA:        "https://pkg.go.dev",
					originB:        "https://beta.pkg.go.dev",
					screenshotType: viewportScreenshot,
					viewportWidth:  540,
					viewportHeight: 1080,
					pathame:        "/net",
					tasks: chromedp.Tasks{
						chromedp.WaitReady(`[role="treeitem"][aria-expanded="true"]`),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readTests(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("readTests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got,
				cmp.AllowUnexported(testcase{}),
				cmp.Comparer(func(a, b chromedp.ActionFunc) bool {
					return fmt.Sprint(a) == fmt.Sprint(b)
				}),
				cmp.Comparer(func(a, b chromedp.Selector) bool {
					return fmt.Sprint(a) == fmt.Sprint(b)
				}),
			); diff != "" {
				t.Errorf("readTests() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCheckHandler(t *testing.T) {
	// Skip this test if Google Chrome is not installed.
	_, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip()
	}
	type args struct {
		glob string
	}
	d, err := os.UserCacheDir()
	if err != nil {
		t.Errorf("os.UserCacheDir(): %v", err)
	}
	cache := filepath.Join(d, "screentest")
	var tests = []struct {
		name      string
		args      args
		wantErr   bool
		wantFiles []string
	}{
		{
			name: "pass",
			args: args{
				glob: "testdata/pass.txt",
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				glob: "testdata/fail.txt",
			},
			wantErr: true,
			wantFiles: []string{
				filepath.Join(cache, "fail-txt", "homepage.diff.png"),
				filepath.Join(cache, "fail-txt", "homepage.go-dev.png"),
				filepath.Join(cache, "fail-txt", "homepage.pkg-go-dev.png"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckHandler(tt.args.glob); (err != nil) != tt.wantErr {
				t.Fatalf("CheckHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				files, err := filepath.Glob(
					filepath.Join(cache, sanitized(filepath.Base(tt.args.glob)), "*.png"))
				if err != nil {
					t.Fatal("error reading diff output")
				}
				if diff := cmp.Diff(tt.wantFiles, files); diff != "" {
					t.Errorf("readTests() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
