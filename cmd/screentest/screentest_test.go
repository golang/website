// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	stdcmp "cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"image"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/n7olkachev/imgdiff/pkg/imgdiff"
)

var bucket = flag.String("bucket", "", "bucket to test GCS I/O")

func TestReadTests(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name             string
		testURL, wantURL string
		opts             options
		want             []*testcase
		wantErr          bool
	}{
		{
			name:    "readtests",
			testURL: "https://go.dev",
			wantURL: "http://localhost:6060",
			opts: options{
				vars: "Authorization:Bearer token",
			},
			want: []*testcase{
				{
					common: common{
						vars: map[string]string{"Authorization": "Bearer token"},
					},
					name:           "go.dev homepage",
					path:           "/",
					testURL:        "https://go.dev/",
					wantURL:        "http://localhost:6060/",
					status:         200,
					testPath:       "readtests/go-dev-homepage.got.png",
					wantPath:       "readtests/go-dev-homepage.want.png",
					diffPath:       "readtests/go-dev-homepage.diff.png",
					viewportWidth:  1536,
					viewportHeight: 960,
					screenshotType: fullScreenshot,
				},
				{
					common: common{
						vars: map[string]string{"Authorization": "Bearer token"},
					},
					name:           "go.dev homepage 540x1080",
					path:           "/",
					testURL:        "https://go.dev/",
					wantURL:        "http://localhost:6060/",
					status:         200,
					testPath:       "readtests/go-dev-homepage-540x1080.got.png",
					wantPath:       "readtests/go-dev-homepage-540x1080.want.png",
					diffPath:       "readtests/go-dev-homepage-540x1080.diff.png",
					viewportWidth:  540,
					viewportHeight: 1080,
					screenshotType: fullScreenshot,
				},
				{
					common: common{
						vars: map[string]string{"Authorization": "Bearer token"},
					},
					name:           "about page",
					path:           "/about",
					testURL:        "https://go.dev/about",
					wantURL:        "http://localhost:6060/about",
					status:         200,
					testPath:       "readtests/about-page.got.png",
					wantPath:       "readtests/about-page.want.png",
					diffPath:       "readtests/about-page.diff.png",
					screenshotType: fullScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
				},
				{
					common: common{
						vars: map[string]string{"Authorization": "Bearer token"},
					},
					name:              "homepage element .go-Carousel",
					path:              "/",
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
				},
				{
					common: common{
						vars: map[string]string{"Authorization": "Bearer token"},
					},
					name:           "net package doc",
					path:           "/net",
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
				},
				{
					common: common{
						vars: map[string]string{"Authorization": "Bearer token"},
					},
					name:           "net package doc 540x1080",
					path:           "/net",
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
				},
			},
		},
		{
			name:    "readtests2",
			testURL: "some/directory",
			wantURL: "http://localhost:8080",
			opts: options{
				headers:      "Authorization:Bearer token",
				outputDirURL: "gs://bucket/prefix",
			},
			want: []*testcase{
				{
					common: common{
						testImageReader: &dirImageReadWriter{dir: "some/directory"},
						headers:         map[string]any{"Authorization": "Bearer token"},
					},
					name:           "about",
					path:           "/about",
					wantURL:        "http://localhost:8080/about",
					status:         200,
					testPath:       "readtests2/about.got.png",
					wantPath:       "readtests2/about.want.png",
					diffPath:       "readtests2/about.diff.png",
					screenshotType: viewportScreenshot,
					viewportWidth:  100,
					viewportHeight: 200,
				},
				{
					common: common{
						testImageReader: &dirImageReadWriter{dir: "some/directory"},
						headers:         map[string]interface{}{"Authorization": "Bearer token"},
					},
					name:           "eval",
					path:           "/eval",
					wantURL:        "http://localhost:8080/eval",
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
				},
			},
		},
		{
			name:    "readtests-no-capture",
			wantErr: true,
		},
		{
			name:    "readtests-ends-with-capture",
			wantErr: true,
		},
		{
			name: "readtests-filter",
			opts: options{
				filterRegexp: `foo \d`,
			},
			// There are three tests, "foo", "foo 20x30" and "bar". The filter
			// should select only the second.
			want: []*testcase{
				{
					name: "foo 20x30",
					common: common{
						testImageReader:     &dirImageReadWriter{dir: "."},
						wantImageReadWriter: &dirImageReadWriter{dir: "."},
					},
					path:           "p", // all the tests specify the path "p"
					status:         200,
					screenshotType: viewportScreenshot,
					viewportWidth:  20,
					viewportHeight: 30,
					testPath:       "readtests-filter/foo-20x30.got.png",
					wantPath:       "readtests-filter/foo-20x30.want.png",
					diffPath:       "readtests-filter/foo-20x30.diff.png",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join("testdata", tt.name+".txt")
			comm, err := commonValues(ctx, tt.testURL, tt.wantURL, tt.opts)
			if err != nil {
				t.Fatal(err)
			}
			got, err := readTests(filename, tt.testURL, tt.wantURL, comm)
			if (err != nil) != tt.wantErr {
				t.Fatalf("readTests() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.want, got,
				cmp.AllowUnexported(testcase{}),
				cmpopts.IgnoreFields(testcase{}, "output", "tasks"),
				cmp.AllowUnexported(common{}),
				cmpopts.IgnoreFields(common{}, "failImageWriter", "filter"),
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
		common: common{
			headers: map[string]interface{}{"Authorization": "Bearer token"},
		},
		name:              "go.dev homepage",
		testURL:           "http://localhost:6061",
		wantURL:           "http://localhost:6061",
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
	dir := "filedir"
	fpath := dir + "/file.png"

	test := func(t *testing.T, rw imageReadWriter) {
		if err := rw.writeImage(ctx, fpath, img); err != nil {
			t.Fatal(err)
		}
		got, err := rw.readImage(ctx, fpath)
		if err != nil {
			t.Fatal(err)
		}
		result := imgdiff.Diff(img, got, &imgdiff.Options{})
		if !result.Equal {
			t.Error("images not equal")
		}

		t.Logf("removing %q", dir)
		if err := rw.rmdir(ctx, dir); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("dir", func(t *testing.T) {
		tmp := t.TempDir()
		deldir := filepath.Join(tmp, dir)
		test(t, &dirImageReadWriter{tmp})
		if _, err := os.Stat(deldir); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("directory %s still exists", deldir)
		}
	})
	t.Run("gcs", func(t *testing.T) {
		if *bucket == "" {
			t.Skip("missing -bucket")
		}
		// Construct a prefix (directory) that does not already exist.
		prefix := fmt.Sprintf("screentest-test-%s-%s",
			stdcmp.Or(os.Getenv("USER"), "unknown"), time.Now().Format("2006-01-02-03-04-05"))

		rw, err := newGCSImageReadWriter(ctx, fmt.Sprintf("gs://%s/%s", *bucket, prefix))
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("testing with %s", rw.url)
		test(t, rw)

		client, err := storage.NewClient(ctx)
		if err != nil {
			t.Fatal(err)
		}
		deldir := path.Join(rw.prefix, dir)
		if _, err := client.Bucket(*bucket).Object(deldir).Attrs(ctx); !errors.Is(err, storage.ErrObjectNotExist) {
			t.Errorf("object %s still exists", deldir)
		}
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
