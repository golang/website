// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestReadTests(t *testing.T) {
	type args struct {
		testURL, wantURL string
		filename         string
	}
	d, err := os.UserCacheDir()
	if err != nil {
		t.Errorf("os.UserCacheDir(): %v", err)
	}
	cache := filepath.Join(d, "screentest")
	tests := []struct {
		name    string
		args    args
		opts    options
		want    any
		wantErr bool
	}{
		{
			name: "readtests",
			args: args{
				testURL:  "https://go.dev",
				wantURL:  "http://localhost:6060",
				filename: "testdata/readtests.txt",
			},
			opts: options{
				vars: "Authorization:Bearer token",
			},
			want: []*testcase{
				{
					name:           "go.dev homepage",
					urlA:           "https://go.dev/",
					urlB:           "http://localhost:6060/",
					status:         200,
					outImgA:        filepath.Join(cache, "readtests-txt", "go-dev-homepage.a.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "go-dev-homepage.b.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "go-dev-homepage.diff.png"),
					viewportWidth:  1536,
					viewportHeight: 960,
					screenshotType: fullScreenshot,
					headers:        map[string]any{},
				},
				{
					name:           "go.dev homepage 540x1080",
					urlA:           "https://go.dev/",
					urlB:           "http://localhost:6060/",
					status:         200,
					outImgA:        filepath.Join(cache, "readtests-txt", "go-dev-homepage-540x1080.a.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "go-dev-homepage-540x1080.b.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "go-dev-homepage-540x1080.diff.png"),
					viewportWidth:  540,
					viewportHeight: 1080,
					screenshotType: fullScreenshot,
					headers:        map[string]any{},
				},
				{
					name:           "about page",
					urlA:           "https://go.dev/about",
					urlB:           "http://localhost:6060/about",
					status:         200,
					outImgA:        filepath.Join(cache, "readtests-txt", "about-page.a.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "about-page.b.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "about-page.diff.png"),
					screenshotType: fullScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
					headers:        map[string]any{},
				},
				{
					name:              "homepage element .go-Carousel",
					urlA:              "https://go.dev/",
					urlB:              "http://localhost:6060/",
					status:            200,
					outImgA:           filepath.Join(cache, "readtests-txt", "homepage-element--go-Carousel.a.png"),
					outImgB:           filepath.Join(cache, "readtests-txt", "homepage-element--go-Carousel.b.png"),
					outDiff:           filepath.Join(cache, "readtests-txt", "homepage-element--go-Carousel.diff.png"),
					screenshotType:    elementScreenshot,
					screenshotElement: ".go-Carousel",
					viewportWidth:     1536,
					viewportHeight:    960,
					tasks: chromedp.Tasks{
						chromedp.Click(".go-Carousel-dot"),
					},
					headers: map[string]any{},
				},
				{
					name:           "net package doc",
					urlA:           "https://go.dev/net",
					urlB:           "http://localhost:6060/net",
					status:         200,
					outImgA:        filepath.Join(cache, "readtests-txt", "net-package-doc.a.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "net-package-doc.b.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "net-package-doc.diff.png"),
					screenshotType: viewportScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
					tasks: chromedp.Tasks{
						chromedp.WaitReady(`[role="treeitem"][aria-expanded="true"]`),
					},
					headers: map[string]any{},
				},
				{
					name:           "net package doc 540x1080",
					urlA:           "https://go.dev/net",
					urlB:           "http://localhost:6060/net",
					status:         200,
					outImgA:        filepath.Join(cache, "readtests-txt", "net-package-doc-540x1080.a.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "net-package-doc-540x1080.b.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "net-package-doc-540x1080.diff.png"),
					screenshotType: viewportScreenshot,
					viewportWidth:  540,
					viewportHeight: 1080,
					tasks: chromedp.Tasks{
						chromedp.WaitReady(`[role="treeitem"][aria-expanded="true"]`),
					},
					headers: map[string]any{},
				},
			},
			wantErr: false,
		},
		{
			name: "readtests2",
			args: args{
				testURL:  "https://pkg.go.dev::cache",
				wantURL:  "http://localhost:8080",
				filename: "testdata/readtests2.txt",
			},
			opts: options{
				headers:   "Authorization:Bearer token",
				outputURL: "gs://bucket/prefix",
			},
			want: []*testcase{
				{
					name:           "about",
					urlA:           "https://pkg.go.dev/about",
					cacheA:         true,
					urlB:           "http://localhost:8080/about",
					headers:        map[string]any{"Authorization": "Bearer token"},
					status:         200,
					gcsBucket:      true,
					outImgA:        "gs://bucket/prefix/readtests2-txt/about.a.png",
					outImgB:        "gs://bucket/prefix/readtests2-txt/about.b.png",
					outDiff:        "gs://bucket/prefix/readtests2-txt/about.diff.png",
					screenshotType: viewportScreenshot,
					viewportWidth:  100,
					viewportHeight: 200,
				},
				{
					name:           "eval",
					urlA:           "https://pkg.go.dev/eval",
					cacheA:         true,
					urlB:           "http://localhost:8080/eval",
					headers:        map[string]interface{}{"Authorization": "Bearer token"},
					status:         200,
					gcsBucket:      true,
					outImgA:        "gs://bucket/prefix/readtests2-txt/eval.a.png",
					outImgB:        "gs://bucket/prefix/readtests2-txt/eval.b.png",
					outDiff:        "gs://bucket/prefix/readtests2-txt/eval.diff.png",
					screenshotType: viewportScreenshot,
					viewportWidth:  100,
					viewportHeight: 200,
					tasks: chromedp.Tasks{
						chromedp.Evaluate("console.log('Hello, world!')", nil),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readTests(tt.args.filename, tt.args.testURL, tt.args.wantURL, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("readTests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got,
				cmp.AllowUnexported(testcase{}),
				cmpopts.IgnoreFields(testcase{}, "output", "tasks"),
				cmp.AllowUnexported(chromedp.Selector{}),
				cmpopts.IgnoreFields(chromedp.Selector{}, "by", "wait", "after"),
			); diff != "" {
				t.Errorf("readTests() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRun(t *testing.T) {
	// Skip this test if Google Chrome is not installed.
	_, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip("google-chrome not installed")
	}
	type args struct {
		testURL, wantURL string
		output           string
		file             string
	}
	d, err := os.UserCacheDir()
	if err != nil {
		t.Errorf("os.UserCacheDir(): %v", err)
	}
	cache := filepath.Join(d, "screentest")
	var tests = []struct {
		name      string
		args      args
		opts      options
		wantErr   bool
		wantFiles []string
	}{
		{
			name: "pass",
			args: args{
				testURL: "https://go.dev",
				wantURL: "https://go.dev",
				file:    "testdata/pass.txt",
			},
			opts:    options{},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				testURL: "https://go.dev",
				wantURL: "https://pkg.go.dev",
				file:    "testdata/fail.txt",
				output:  filepath.Join(cache, "fail-txt"),
			},
			wantErr: true,
			wantFiles: []string{
				filepath.Join(cache, "fail-txt", "homepage.a.png"),
				filepath.Join(cache, "fail-txt", "homepage.b.png"),
				filepath.Join(cache, "fail-txt", "homepage.diff.png"),
			},
		},
		{
			name: "cached",
			args: args{
				testURL: "https://go.dev::cache",
				wantURL: "https://go.dev::cache",
				file:    "testdata/cached.txt",
				output:  "testdata/screenshots/cached",
			},
			opts: options{
				outputURL: "testdata/screenshots/cached",
			},
			wantFiles: []string{
				filepath.Join("testdata", "screenshots", "cached", "homepage.a.png"),
				filepath.Join("testdata", "screenshots", "cached", "homepage.b.png"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(context.Background(), tt.args.testURL, tt.args.wantURL, []string{tt.args.file}, tt.opts); (err != nil) != tt.wantErr {
				t.Fatalf("CheckHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.wantFiles) != 0 {
				files, err := filepath.Glob(
					filepath.Join(tt.args.output, "*.png"))
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

func TestHeaders(t *testing.T) {
	// Skip this test if Google Chrome is not installed.
	_, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip()
	}
	go headerServer()
	tc := &testcase{
		name:              "go.dev homepage",
		urlA:              "http://localhost:6061",
		cacheA:            true,
		urlB:              "http://localhost:6061",
		headers:           map[string]interface{}{"Authorization": "Bearer token"},
		outImgA:           filepath.Join("testdata", "screenshots", "headers", "headers-test.a.png"),
		outImgB:           filepath.Join("testdata", "screenshots", "headers", "headers-test.b.png"),
		outDiff:           filepath.Join("testdata", "screenshots", "headers", "headers-test.diff.png"),
		viewportWidth:     1536,
		viewportHeight:    960,
		screenshotType:    elementScreenshot,
		screenshotElement: "#result",
	}
	if err := tc.run(context.Background(), false); err != nil {
		t.Fatal(err)
	}
}

func headerServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, `<!doctype html>
		<html>
		<body>
		  <span id="result">%s</span>
		</body>
		</html>`, req.Header.Get("Authorization"))
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", 6061), mux)
}

func Test_gcsParts(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name       string
		args       args
		wantBucket string
		wantObject string
	}{
		{
			args: args{
				filename: "gs://bucket-name/object-name",
			},
			wantBucket: "bucket-name",
			wantObject: "object-name",
		},
		{
			args: args{
				filename: "gs://bucket-name/subdir/object-name",
			},
			wantBucket: "bucket-name",
			wantObject: "subdir/object-name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBucket, gotObject := gcsParts(tt.args.filename)
			if gotBucket != tt.wantBucket {
				t.Errorf("gcsParts() gotBucket = %v, want %v", gotBucket, tt.wantBucket)
			}
			if gotObject != tt.wantObject {
				t.Errorf("gcsParts() gotObject = %v, want %v", gotObject, tt.wantObject)
			}
		})
	}
}

func Test_cleanDirs(t *testing.T) {
	f, err := os.Create("testdata/screenshots/cached/should-delete.a.png")
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	type args struct {
		dirs      map[string]bool
		keepFiles map[string]bool
		safeExts  map[string]bool
	}
	tests := []struct {
		name      string
		args      args
		wantFiles map[string]bool
	}{
		{
			name: "keeps files in keepFiles",
			args: args{
				dirs: map[string]bool{
					"testdata/screenshots/cached":  true,
					"testdata/screenshots/headers": true,
					"testdata":                     true,
				},
				keepFiles: map[string]bool{
					"testdata/screenshots/cached/homepage.a.png":      true,
					"testdata/screenshots/cached/homepage.b.png":      true,
					"testdata/screenshots/headers/headers-test.a.png": true,
				},
				safeExts: map[string]bool{
					"a.png": true,
					"b.png": true,
				},
			},
			wantFiles: map[string]bool{
				"testdata/screenshots/cached/homepage.a.png":      true,
				"testdata/screenshots/headers/headers-test.a.png": true,
			},
		},
		{
			name: "keeps files without matching extension",
			args: args{
				dirs: map[string]bool{
					"testdata": true,
				},
				safeExts: map[string]bool{
					"a.png": true,
				},
			},
			wantFiles: map[string]bool{
				"testdata/cached.txt":    true,
				"testdata/fail.txt":      true,
				"testdata/pass.txt":      true,
				"testdata/readtests.txt": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cleanDirs(tt.args.dirs, tt.args.keepFiles, tt.args.safeExts); err != nil {
				t.Fatal(err)
			}
			for file := range tt.wantFiles {
				if _, err := os.Stat(file); err != nil {
					t.Errorf("cleanDirs() error = %v, wantErr %v", err, nil)
				}
			}
		})
	}
}

func TestSplitDimensions(t *testing.T) {
	for _, tc := range []struct {
		in   string
		w, h int
	}{
		{"1x2", 1, 2},
		{"23x40", 23, 40},
	} {
		gw, gh, err := splitDimensions(tc.in)
		if err != nil {
			t.Errorf("%q: %v", tc.in, err)
		} else if gw != tc.w || gh != tc.h {
			t.Errorf("%s: got (%d, %d), want (%d, %d)", tc.in, gw, gh, tc.w, tc.h)
		}
	}

	// Expect errors.
	for _, in := range []string{
		"", "1", "1x2a", "  1x2", "1 x 2", "3x-4",
	} {
		if _, _, err := splitDimensions(in); err == nil {
			t.Errorf("%q: got nil, want error", in)
		}
	}
}
