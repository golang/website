// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/n7olkachev/imgdiff/pkg/imgdiff"
)

var bucketPath = flag.String("bucketpath", "", "bucket/prefix to test GCS I/O")

func TestReadTests(t *testing.T) {
	type args struct {
		testURL, wantURL string
		filename         string
	}
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
					testURL:        "https://go.dev/",
					wantURL:        "http://localhost:6060/",
					status:         200,
					testPath:       "readtests/go-dev-homepage.got.png",
					wantPath:       "readtests/go-dev-homepage.want.png",
					diffPath:       "readtests/go-dev-homepage.diff.png",
					viewportWidth:  1536,
					viewportHeight: 960,
					screenshotType: fullScreenshot,
					headers:        map[string]any{},
				},
				{
					name:           "go.dev homepage 540x1080",
					testURL:        "https://go.dev/",
					wantURL:        "http://localhost:6060/",
					status:         200,
					testPath:       "readtests/go-dev-homepage-540x1080.got.png",
					wantPath:       "readtests/go-dev-homepage-540x1080.want.png",
					diffPath:       "readtests/go-dev-homepage-540x1080.diff.png",
					viewportWidth:  540,
					viewportHeight: 1080,
					screenshotType: fullScreenshot,
					headers:        map[string]any{},
				},
				{
					name:           "about page",
					testURL:        "https://go.dev/about",
					wantURL:        "http://localhost:6060/about",
					status:         200,
					testPath:       "readtests/about-page.got.png",
					wantPath:       "readtests/about-page.want.png",
					diffPath:       "readtests/about-page.diff.png",
					screenshotType: fullScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
					headers:        map[string]any{},
				},
				{
					name:              "homepage element .go-Carousel",
					testURL:           "https://go.dev/",
					wantURL:           "http://localhost:6060/",
					status:            200,
					testPath:          "readtests/homepage-element--go-Carousel.got.png",
					wantPath:          "readtests/homepage-element--go-Carousel.want.png",
					diffPath:          "readtests/homepage-element--go-Carousel.diff.png",
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
					testURL:        "https://go.dev/net",
					wantURL:        "http://localhost:6060/net",
					status:         200,
					testPath:       "readtests/net-package-doc.got.png",
					wantPath:       "readtests/net-package-doc.want.png",
					diffPath:       "readtests/net-package-doc.diff.png",
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
					testURL:        "https://go.dev/net",
					wantURL:        "http://localhost:6060/net",
					status:         200,
					testPath:       "readtests/net-package-doc-540x1080.got.png",
					wantPath:       "readtests/net-package-doc-540x1080.want.png",
					diffPath:       "readtests/net-package-doc-540x1080.diff.png",
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
				testURL:  "some/directory",
				wantURL:  "http://localhost:8080",
				filename: "testdata/readtests2.txt",
			},
			opts: options{
				headers:   "Authorization:Bearer token",
				outputURL: "gs://bucket/prefix",
			},
			want: []*testcase{
				{
					name:            "about",
					wantURL:         "http://localhost:8080/about",
					headers:         map[string]any{"Authorization": "Bearer token"},
					status:          200,
					testPath:        "readtests2/about.got.png",
					wantPath:        "readtests2/about.want.png",
					diffPath:        "readtests2/about.diff.png",
					screenshotType:  viewportScreenshot,
					viewportWidth:   100,
					viewportHeight:  200,
					testImageReader: &dirImageReadWriter{dir: "some/directory"},
				},
				{
					name:           "eval",
					wantURL:        "http://localhost:8080/eval",
					headers:        map[string]interface{}{"Authorization": "Bearer token"},
					status:         200,
					testPath:       "readtests2/eval.got.png",
					wantPath:       "readtests2/eval.want.png",
					diffPath:       "readtests2/eval.diff.png",
					screenshotType: viewportScreenshot,
					viewportWidth:  100,
					viewportHeight: 200,
					tasks: chromedp.Tasks{
						chromedp.Evaluate("console.log('Hello, world!')", nil),
					},
					testImageReader: &dirImageReadWriter{dir: "some/directory"},
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
				cmpopts.IgnoreFields(testcase{}, "output", "tasks", "failImageWriter"),
				cmp.AllowUnexported(chromedp.Selector{}),
				cmpopts.IgnoreFields(chromedp.Selector{}, "by", "wait", "after"),
				cmp.AllowUnexported(dirImageReadWriter{}),
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
				output:  filepath.Join(cache, "fail"),
			},
			wantErr: true,
			wantFiles: []string{
				filepath.Join(cache, "fail", "homepage.diff.png"),
				filepath.Join(cache, "fail", "homepage.got.png"),
				filepath.Join(cache, "fail", "homepage.want.png"),
			},
		},
		{
			name: "want stored",
			args: args{
				testURL: "https://go.dev",
				wantURL: "testdata/screenshots",
				file:    "testdata/cached.txt",
				output:  "testdata/screenshots/cached",
			},
			opts: options{update: true},
			wantFiles: []string{
				filepath.Join("testdata", "screenshots", "cached", "homepage.want.png"),
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
		testURL:           "http://localhost:6061",
		wantURL:           "http://localhost:6061",
		headers:           map[string]interface{}{"Authorization": "Bearer token"},
		testPath:          filepath.Join("testdata", "screenshots", "headers", "headers-test.got.png"),
		wantPath:          filepath.Join("testdata", "screenshots", "headers", "headers-test.want.png"),
		diffPath:          filepath.Join("testdata", "screenshots", "headers", "headers-test.diff.png"),
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

func TestReadWriters(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 10, 10))
	ctx := context.Background()
	path := "sub/file.png"

	test := func(t *testing.T, rw imageReadWriter) {
		if err := rw.writeImage(ctx, path, img); err != nil {
			t.Fatal(err)
		}
		got, err := rw.readImage(ctx, path)
		if err != nil {
			t.Fatal(err)
		}
		result := imgdiff.Diff(img, got, &imgdiff.Options{})
		if !result.Equal {
			t.Error("images not equal")
		}
	}

	t.Run("dir", func(t *testing.T) {
		test(t, &dirImageReadWriter{t.TempDir()})
	})
	t.Run("gcs", func(t *testing.T) {
		if *bucketPath == "" {
			t.Skip("missing -bucketpath")
		}
		rw, err := newGCSImageReadWriter(ctx, "gs://"+*bucketPath)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%+v\n", rw)
		test(t, rw)
	})
}

func TestNewImageReadWriter(t *testing.T) {
	for _, tc := range []struct {
		in         string
		wantDir    string // implies dirImageReadWriter
		wantPrefix string // implies gcsImageReadWriter
	}{
		{
			in:      "unix/path",
			wantDir: "unix/path",
		},
		{
			in:      "unix/path/../dir",
			wantDir: "unix/dir",
		},
		{
			in:      "c:/windows/path",
			wantDir: "c:/windows/path",
		},
		{
			in:      "c:/windows/../dir",
			wantDir: "c:/dir",
		},
		{
			in:      "file:///file/path",
			wantDir: "/file/path",
		},
		{
			in:         "gs://bucket/prefix",
			wantPrefix: "prefix",
		},
		{
			in: "http://example.com",
		},
	} {
		got, err := newImageReadWriter(context.Background(), tc.in)
		if err != nil {
			t.Fatal(err)
		}
		if tc.wantDir != "" {
			d, ok := got.(*dirImageReadWriter)
			if !ok || d.dir != tc.wantDir {
				t.Errorf("%s: got %+v, want dirImageReadWriter{dir: %q}", tc.in, got, tc.wantDir)
			}
		} else if tc.wantPrefix != "" {
			g, ok := got.(*gcsImageReadWriter)
			if !ok || g.prefix != tc.wantPrefix {
				t.Errorf("%s: got %+v, want gcsImageReadWriter{prefix: %q}", tc.in, got, tc.wantPrefix)
			}
		} else if got != nil {
			t.Errorf("%s: got %+v, want nil", tc.in, got)
		}
	}
}
