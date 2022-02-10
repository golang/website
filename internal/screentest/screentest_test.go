// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package screentest

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
		filename string
	}
	d, err := os.UserCacheDir()
	if err != nil {
		t.Errorf("os.UserCacheDir(): %v", err)
	}
	cache := filepath.Join(d, "screentest")
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
					urlA:           "https://go.dev/",
					urlB:           "http://localhost:6060/go.dev/",
					status:         200,
					outImgA:        filepath.Join(cache, "readtests-txt", "go-dev-homepage.a.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "go-dev-homepage.b.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "go-dev-homepage.diff.png"),
					viewportWidth:  1536,
					viewportHeight: 960,
					screenshotType: fullScreenshot,
				},
				{
					name:           "go.dev homepage 540x1080",
					urlA:           "https://go.dev/",
					urlB:           "http://localhost:6060/go.dev/",
					status:         200,
					outImgA:        filepath.Join(cache, "readtests-txt", "go-dev-homepage-540x1080.a.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "go-dev-homepage-540x1080.b.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "go-dev-homepage-540x1080.diff.png"),
					viewportWidth:  540,
					viewportHeight: 1080,
					screenshotType: fullScreenshot,
				},
				{
					name:           "about page",
					urlA:           "https://go.dev/about",
					urlB:           "http://localhost:6060/go.dev/about",
					status:         200,
					outImgA:        filepath.Join(cache, "readtests-txt", "about-page.a.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "about-page.b.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "about-page.diff.png"),
					screenshotType: fullScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
				},
				{
					name:              "pkg.go.dev homepage .go-Carousel",
					urlA:              "https://pkg.go.dev/",
					urlB:              "https://beta.pkg.go.dev/",
					status:            200,
					outImgA:           filepath.Join(cache, "readtests-txt", "pkg-go-dev-homepage--go-Carousel.a.png"),
					outImgB:           filepath.Join(cache, "readtests-txt", "pkg-go-dev-homepage--go-Carousel.b.png"),
					outDiff:           filepath.Join(cache, "readtests-txt", "pkg-go-dev-homepage--go-Carousel.diff.png"),
					screenshotType:    elementScreenshot,
					screenshotElement: ".go-Carousel",
					viewportWidth:     1536,
					viewportHeight:    960,
					tasks: chromedp.Tasks{
						chromedp.Click(".go-Carousel-dot"),
					},
				},
				{
					name:           "net package doc",
					urlA:           "https://pkg.go.dev/net",
					urlB:           "https://beta.pkg.go.dev/net",
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
				},
				{
					name:           "net package doc 540x1080",
					urlA:           "https://pkg.go.dev/net",
					urlB:           "https://beta.pkg.go.dev/net",
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
				},
				{
					name:           "about",
					urlA:           "https://pkg.go.dev/about",
					cacheA:         true,
					urlB:           "http://localhost:8080/about",
					headers:        map[string]interface{}{"Authorization": "Bearer token"},
					status:         200,
					outImgA:        filepath.Join(cache, "readtests-txt", "about.a.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "about.b.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "about.diff.png"),
					screenshotType: viewportScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
				},
				{
					name:           "eval",
					urlA:           "https://pkg.go.dev/eval",
					cacheA:         true,
					urlB:           "http://localhost:8080/eval",
					headers:        map[string]interface{}{"Authorization": "Bearer token"},
					status:         200,
					outImgA:        filepath.Join(cache, "readtests-txt", "eval.a.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "eval.b.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "eval.diff.png"),
					screenshotType: viewportScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
					tasks: chromedp.Tasks{
						chromedp.Evaluate("console.log('Hello, world!')", nil),
					},
				},
				{
					name:           "gcs-output",
					urlA:           "https://pkg.go.dev/gcs-output",
					cacheA:         true,
					urlB:           "http://localhost:8080/gcs-output",
					gcsBucket:      true,
					headers:        map[string]interface{}{"Authorization": "Bearer token"},
					status:         200,
					outImgA:        "gs://bucket-name/gcs-output.a.png",
					outImgB:        "gs://bucket-name/gcs-output.b.png",
					outDiff:        "gs://bucket-name/gcs-output.diff.png",
					screenshotType: viewportScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readTests(tt.args.filename, map[string]string{"Authorization": "Bearer token"})
			if (err != nil) != tt.wantErr {
				t.Errorf("readTests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got,
				cmp.AllowUnexported(testcase{}),
				cmpopts.IgnoreFields(testcase{}, "output"),
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
		glob   string
		output string
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
				output: filepath.Join(cache, "fail-txt"),
				glob:   "testdata/fail.txt",
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
				output: "testdata/screenshots/cached",
				glob:   "testdata/cached.txt",
			},
			wantFiles: []string{
				filepath.Join("testdata", "screenshots", "cached", "homepage.a.png"),
				filepath.Join("testdata", "screenshots", "cached", "homepage.b.png"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckHandler(tt.args.glob, CheckOptions{}); (err != nil) != tt.wantErr {
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

func TestTestHandler(t *testing.T) {
	// Skip this test if Google Chrome is not installed.
	_, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip()
	}
	TestHandler(t, "testdata/pass.txt", TestOpts{})
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
