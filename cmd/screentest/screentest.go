// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO(jba): sleep directive
// TODO(jba): specify percent of image that may differ
// TODO(jba): remove ints function in template (see cmd/golangorg/testdata/screentest/relnotes.txt)

package main

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"time"

	"cloud.google.com/go/storage"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/n7olkachev/imgdiff/pkg/imgdiff"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
)

// run compares testURL and wantURL using the test scripts in files and the options in opts.
func run(ctx context.Context, testURL, wantURL string, files []string, opts options) error {
	if opts.maxConcurrency < 1 {
		opts.maxConcurrency = 1
	}

	now := time.Now()

	if testURL == "" {
		return errors.New("missing URL or path to test")
	}
	if wantURL == "" {
		return errors.New("missing URL or path with expected results")
	}
	if len(files) == 0 {
		return errors.New("no files to run")
	}
	var cancel context.CancelFunc
	if opts.debuggerURL != "" {
		ctx, cancel = chromedp.NewRemoteAllocator(ctx, opts.debuggerURL)
	} else {
		ctx, cancel = chromedp.NewExecAllocator(ctx, append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.WindowSize(browserWidth, browserHeight),
		)...)
	}
	defer cancel()

	var buf bytes.Buffer
	for _, file := range files {
		tests, err := readTests(file, testURL, wantURL, opts)
		if err != nil {
			return fmt.Errorf("readTestdata(%q): %w", file, err)
		}
		if len(tests) == 0 && opts.run == "" {
			return fmt.Errorf("no tests found in %q", file)
		}
		if err := cleanOutput(ctx, tests); err != nil {
			return fmt.Errorf("cleanOutput(...): %w", err)
		}
		ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
		defer cancel()
		var hdr bool
		runConcurrently(len(tests), opts.maxConcurrency, func(i int) {
			tc := tests[i]
			if err := tc.run(ctx, opts.update); err != nil {
				if !hdr {
					fmt.Fprintf(&buf, "%s\n\n", file)
					hdr = true
				}
				fmt.Fprintf(&buf, "%v\n", err)
				fmt.Fprintf(&buf, "inspect diff at %s\n\n", tc.outDiff)
			}
			fmt.Println(tc.output.String())
		})
	}
	fmt.Printf("finished in %s\n\n", time.Since(now).Truncate(time.Millisecond))
	if buf.Len() > 0 {
		return errors.New(buf.String())
	}
	return nil
}

// cleanOutput clears the output locations of images not cached
// as part of a testcase, including diff output from previous test
// runs and obsolete screenshots. It ensures local directories exist
// for test output. GCS buckets must already exist prior to test run.
func cleanOutput(ctx context.Context, tests []*testcase) error {
	keepFiles := make(map[string]bool)
	bkts := make(map[string]bool)
	dirs := make(map[string]bool)
	// The extensions of files that are safe to delete
	safeExts := map[string]bool{
		"a.png":    true,
		"b.png":    true,
		"diff.png": true,
	}
	for _, t := range tests {
		if t.cacheA {
			keepFiles[t.outImgA] = true
		}
		if t.cacheB {
			keepFiles[t.outImgB] = true
		}
		if t.gcsBucket {
			bkt, _ := gcsParts(t.outDiff)
			bkts[bkt] = true
		} else {
			dirs[filepath.Dir(t.outDiff)] = true
		}
	}
	if err := cleanBkts(ctx, bkts, keepFiles, safeExts); err != nil {
		return fmt.Errorf("cleanBkts(...): %w", err)
	}
	if err := cleanDirs(dirs, keepFiles, safeExts); err != nil {
		return fmt.Errorf("cleanDirs(...): %w", err)
	}
	return nil
}

// cleanBkts clears all the GCS buckets in bkts of all objects not included
// in the set of keepFiles. Buckets that do not exist will cause an error.
func cleanBkts(ctx context.Context, bkts, keepFiles, safeExts map[string]bool) error {
	if len(bkts) == 0 {
		return nil
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient(ctx): %w", err)
	}
	defer client.Close()
	for bkt := range bkts {
		it := client.Bucket(bkt).Objects(ctx, nil)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("it.Next(): %w", err)
			}
			filename := "gs://" + attrs.Bucket + "/" + attrs.Name
			if !keepFiles[filename] && safeExts[ext(filename)] {
				if err := client.Bucket(attrs.Bucket).Object(attrs.Name).Delete(ctx); err != nil &&
					!errors.Is(err, storage.ErrObjectNotExist) {
					return fmt.Errorf("Object(%q).Delete: %v", attrs.Name, err)
				}
			}
		}
	}
	return client.Close()
}

// cleanDirs ensures the set of directories in dirs exists and
// clears dirs of all files not included in the set of keepFiles.
func cleanDirs(dirs, keepFiles, safeExts map[string]bool) error {
	for dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("os.MkdirAll(%q): %w", dir, err)
		}
	}
	for dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("os.ReadDir(%q): %w", dir, err)
		}
		for _, f := range files {
			filename := dir + "/" + f.Name()
			if !keepFiles[filename] && safeExts[ext(filename)] {
				if err := os.Remove(filename); err != nil && !errors.Is(err, fs.ErrNotExist) {
					return fmt.Errorf("os.Remove(%q): %w", filename, err)
				}
			}
		}
	}
	return nil
}

func ext(filename string) string {
	// If the filename has multiple dots use the first one as
	// the split for the extension.
	if strings.Count(filename, ".") > 1 {
		base := filepath.Base(filename)
		parts := strings.SplitN(base, ".", 2)
		return parts[1]
	}
	return filepath.Ext(filename)
}

const (
	browserWidth  = 1536
	browserHeight = 960
	cacheSuffix   = "::cache"
	gcsScheme     = "gs://"
)

type screenshotType int

const (
	fullScreenshot screenshotType = iota
	viewportScreenshot
	elementScreenshot
)

type testcase struct {
	name                      string
	tasks                     chromedp.Tasks
	urlA, urlB                string
	headers                   map[string]any // to match chromedp arg
	status                    int
	cacheA, cacheB            bool
	gcsBucket                 bool
	outImgA, outImgB, outDiff string
	viewportWidth             int
	viewportHeight            int
	screenshotType            screenshotType
	screenshotElement         string
	blockedURLs               []string
	output                    bytes.Buffer
}

func (t *testcase) String() string {
	return t.name
}

// readTests parses the testcases from a text file.
func readTests(file, testURL, wantURL string, opts options) ([]*testcase, error) {
	tmpl := template.New(filepath.Base(file)).Funcs(template.FuncMap{
		"ints": func(start, end int) []int {
			var out []int
			for i := start; i < end; i++ {
				out = append(out, i)
			}
			return out
		},
	})

	_, err := tmpl.ParseFiles(file)
	if err != nil {
		return nil, fmt.Errorf("template.ParseFiles(%q): %w", file, err)
	}

	parsedVars, err := splitList(opts.vars)
	if err != nil {
		return nil, err
	}

	var tmplout bytes.Buffer
	if err := tmpl.Execute(&tmplout, parsedVars); err != nil {
		return nil, fmt.Errorf("tmpl.Execute(...): %w", err)
	}
	var tests []*testcase
	var (
		testName, pathname string
		tasks              chromedp.Tasks
		originA, originB   string
		status             int = http.StatusOK
		cacheA, cacheB     bool
		gcsBucket          bool
		width, height      int
		lineNo             int
		blockedURLs        []string
	)
	cache, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("os.UserCacheDir(): %w", err)
	}
	originA = testURL
	if strings.HasSuffix(originA, cacheSuffix) {
		originA = strings.TrimSuffix(originA, cacheSuffix)
		cacheA = true
	}
	if _, err := url.Parse(originA); err != nil {
		return nil, fmt.Errorf("url.Parse(%q): %w", originA, err)
	}
	originB = wantURL
	if strings.HasSuffix(originB, cacheSuffix) {
		originB = strings.TrimSuffix(originB, cacheSuffix)
		cacheB = true
	}
	if _, err := url.Parse(originB); err != nil {
		return nil, fmt.Errorf("url.Parse(%q): %w", originB, err)
	}

	headers := map[string]any{}
	hs, err := splitList(opts.headers)
	if err != nil {
		return nil, err
	}
	for k, v := range hs {
		headers[k] = v
	}
	dir := cmp.Or(opts.outputURL, filepath.Join(cache, "screentest"))
	out, err := outDir(dir, file)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(out, gcsScheme) {
		gcsBucket = true
	}

	filter := func(string) bool { return true }
	if opts.run != "" {
		re, err := regexp.Compile(opts.run)
		if err != nil {
			return nil, err
		}
		filter = re.MatchString
	}

	scan := bufio.NewScanner(&tmplout)
	for scan.Scan() {
		lineNo++
		line := strings.TrimSpace(scan.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimRight(line, " \t")
		field, args := splitOneField(line)
		field = strings.ToUpper(field)
		if testName == "" && !slices.Contains([]string{"", "TEST", "BLOCK", "WINDOWSIZE"}, field) {
			return nil, fmt.Errorf("%s:%d: the %q directive should only occur in a test", file, lineNo, strings.ToLower(field))
		}
		switch field {
		case "":
			// We've reached an empty line, reset properties scoped to a single test.
			testName, pathname = "", ""
			tasks = nil
			status = http.StatusOK
		case "STATUS":
			status, err = strconv.Atoi(args)
			if err != nil {
				return nil, fmt.Errorf("strconv.Atoi(%q): %w", args, err)
			}
		case "WINDOWSIZE":
			width, height, err = splitDimensions(args)
			if err != nil {
				return nil, err
			}
		case "TEST":
			testName = args
			for _, t := range tests {
				if t.name == testName {
					return nil, fmt.Errorf("%s:%d: duplicate test name %q", file, lineNo, testName)
				}
			}
		case "PATHNAME":
			if _, err := url.Parse(originA + args); err != nil {
				return nil, fmt.Errorf("url.Parse(%q): %w", originA+args, err)
			}
			if _, err := url.Parse(originB + args); err != nil {
				return nil, fmt.Errorf("url.Parse(%q): %w", originB+args, err)
			}
			pathname = args
			if testName == "" {
				testName = pathname[1:]
			}
			for _, t := range tests {
				if t.name == testName {
					return nil, fmt.Errorf(
						"duplicate test with pathname %q on line %d", pathname, lineNo)
				}
			}
		case "CLICK":
			tasks = append(tasks, chromedp.Click(args))
		case "WAIT":
			tasks = append(tasks, chromedp.WaitReady(args))
		case "EVAL":
			tasks = append(tasks, chromedp.Evaluate(args, nil))
		case "BLOCK":
			blockedURLs = append(blockedURLs, strings.Fields(args)...)
		case "CAPTURE":
			if pathname == "" {
				return nil, fmt.Errorf("missing pathname for capture on line %d", lineNo)
			}
			urlA, err := url.Parse(originA + pathname)
			if err != nil {
				return nil, fmt.Errorf("url.Parse(%q): %w", originA+pathname, err)
			}
			urlB, err := url.Parse(originB + pathname)
			if err != nil {
				return nil, fmt.Errorf("url.Parse(%q): %w", originB+pathname, err)
			}
			if !filter(testName) {
				continue
			}
			test := &testcase{
				name:        testName,
				tasks:       tasks,
				urlA:        urlA.String(),
				urlB:        urlB.String(),
				headers:     headers,
				status:      status,
				blockedURLs: blockedURLs,
				// Default to viewportScreenshot
				screenshotType: viewportScreenshot,
				viewportWidth:  width,
				viewportHeight: height,
				cacheA:         cacheA,
				cacheB:         cacheB,
				gcsBucket:      gcsBucket,
			}
			tests = append(tests, test)
			field, args := splitOneField(args)
			field = strings.ToUpper(field)
			switch field {
			case "FULLSCREEN", "VIEWPORT":
				if field == "FULLSCREEN" {
					test.screenshotType = fullScreenshot
				}
				if args != "" {
					w, h, err := splitDimensions(args)
					if err != nil {
						return nil, err
					}
					test.name += fmt.Sprintf(" %dx%d", w, h)
					test.viewportWidth = w
					test.viewportHeight = h
				}
			case "ELEMENT":
				test.name += fmt.Sprintf(" %s", args)
				test.screenshotType = elementScreenshot
				test.screenshotElement = args
			case "":
				// nothing to do
			default:
				return nil, fmt.Errorf("first argument to capture must be 'fullscreen', 'viewport' or 'element'")
			}
			outfile := filepath.Join(out, sanitize(test.name))
			if gcsBucket {
				outfile, err = url.JoinPath(out, sanitize(test.name))
			}
			test.outImgA = outfile + ".a.png"
			test.outImgB = outfile + ".b.png"
			test.outDiff = outfile + ".diff.png"
		default:
			// We should never reach this error.
			return nil, fmt.Errorf("invalid syntax on line %d: %q", lineNo, line)
		}
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("scan.Err(): %v", err)
	}
	return tests, nil
}

// splitList splits a list of key:value pairs separated by commas.
// Whitespace is trimmed around comma-separated elements, keys, and values.
// Empty names are an error; empty values are OK.
func splitList(s string) (map[string]string, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil, nil
	}
	m := map[string]string{}
	for _, h := range strings.Split(s, ",") {
		name, value, ok := strings.Cut(h, ":")
		if !ok || name == "" {
			return nil, fmt.Errorf("invalid name:value pair: %q", h)
		}
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		m[name] = value
	}
	return m, nil
}

// outDir gets a diff output directory for a given testfile.
// If dir points to a GCS bucket or testfile is empty it just
// returns dir.
func outDir(dir, testfile string) (string, error) {
	tf := sanitize(filepath.Base(testfile))
	if strings.HasPrefix(dir, gcsScheme) {
		return url.JoinPath(dir, tf)
	}
	return filepath.Clean(filepath.Join(dir, tf)), nil
}

// splitOneField splits text at the first space or tab
// and returns that first field and the remaining text.
func splitOneField(text string) (field, rest string) {
	i := strings.IndexAny(text, " \t")
	if i < 0 {
		return text, ""
	}
	return text[:i], strings.TrimSpace(text[i:])
}

// splitDimensions parses a window dimension string into int values
// for width and height.
func splitDimensions(text string) (width, height int, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("splitDimensions(%q): %w", text, err)
		}
	}()

	windowsize := strings.Split(text, "x")
	if len(windowsize) != 2 {
		return 0, 0, errors.New("syntax error")
	}
	width, err = strconv.Atoi(windowsize[0])
	if err != nil {
		return 0, 0, err
	}
	height, err = strconv.Atoi(windowsize[1])
	if err != nil {
		return 0, 0, err
	}
	if width < 0 || height < 0 {
		return 0, 0, errors.New("negative dimension")
	}
	return width, height, nil
}

// run generates screenshots for a given test case and a diff if the
// screenshots do not match.
func (tc *testcase) run(ctx context.Context, update bool) (err error) {
	now := time.Now()
	fmt.Fprintf(&tc.output, "test %s ", tc.name)
	var screenA, screenB *image.Image
	g, ctx := errgroup.WithContext(ctx)
	// If the hosts are the same, chrome (or chromedp) does not handle concurrent requests well.
	// This wouldn't make sense in an actual test, but it does happen in this package's tests.
	urla, erra := url.Parse(tc.urlA)
	urlb, errb := url.Parse(tc.urlA)
	if err := cmp.Or(erra, errb); err != nil {
		return err
	}
	if urla.Host == urlb.Host {
		g.SetLimit(1)
	}
	g.Go(func() error {
		screenA, err = tc.screenshot(ctx, tc.urlA, tc.outImgA, tc.cacheA, update)
		if err != nil {
			return fmt.Errorf("screenshot(ctx, %q, %q, %q, %v): %w", tc, tc.urlA, tc.outImgA, tc.cacheA, err)
		}
		return nil
	})
	g.Go(func() error {
		screenB, err = tc.screenshot(ctx, tc.urlB, tc.outImgB, tc.cacheB, update)
		if err != nil {
			return fmt.Errorf("screenshot(ctx, %q, %q, %q, %v): %w", tc, tc.urlB, tc.outImgB, tc.cacheB, err)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		fmt.Fprint(&tc.output, "\n", err)
		return err
	}
	result := imgdiff.Diff(*screenA, *screenB, &imgdiff.Options{
		Threshold: 0.1,
		DiffImage: true,
	})
	since := time.Since(now).Truncate(time.Millisecond)
	if result.Equal {
		fmt.Fprintf(&tc.output, "(%s)\n", since)
		return nil
	}
	fmt.Fprintf(&tc.output, "(%s)\nFAIL %s != %s (%d pixels differ)\n", since, tc.urlA, tc.urlB, result.DiffPixelsCount)
	g = &errgroup.Group{}
	g.Go(func() error {
		return writePNG(&result.Image, tc.outDiff)
	})
	// Only write screenshots if they haven't already been written to the cache.
	if !tc.cacheA {
		g.Go(func() error {
			return writePNG(screenA, tc.outImgA)
		})
	}
	if !tc.cacheB {
		g.Go(func() error {
			return writePNG(screenB, tc.outImgB)
		})
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf("writePNG(...): %w", err)
	}
	fmt.Fprintf(&tc.output, "wrote diff to %s\n", tc.outDiff)
	return fmt.Errorf("%s != %s", tc.urlA, tc.urlB)
}

// screenshot gets a screenshot for a testcase url. When cache is true it will
// attempt to read the screenshot from a cache or capture a new screenshot
// and write it to the cache if it does not exist.
func (tc *testcase) screenshot(ctx context.Context, url, file string,
	cache, update bool,
) (_ *image.Image, err error) {
	var data []byte
	// If cache is enabled, try to read the file from the cache.
	if cache && tc.gcsBucket {
		client, err := storage.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("storage.NewClient(err): %w", err)
		}
		defer client.Close()
		bkt, obj := gcsParts(file)
		r, err := client.Bucket(bkt).Object(obj).NewReader(ctx)
		if err != nil && !errors.Is(err, storage.ErrObjectNotExist) {
			return nil, fmt.Errorf("object.NewReader(ctx): %w", err)
		} else if err == nil {
			defer r.Close()
			data, err = io.ReadAll(r)
			if err != nil {
				return nil, fmt.Errorf("io.ReadAll(...): %w", err)
			}
		}
	} else if cache {
		data, err = os.ReadFile(file)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("os.ReadFile(...): %w", err)
		}
	}
	// If cache is false, an update is requested, or this is the first test run
	// we capture a new screenshot from a live URL.
	if !cache || update || data == nil {
		update = true
		data, err = tc.captureScreenshot(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("captureScreenshot(ctx, %q, %q): %w", url, tc, err)
		}
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("image.Decode(...): %w", err)
	}
	// Write to the cache.
	if cache && update {
		if err := writePNG(&img, file); err != nil {
			return nil, fmt.Errorf("os.WriteFile(...): %w", err)
		}
		fmt.Fprintf(&tc.output, "updated %s\n", file)
	}
	return &img, nil
}

type response struct {
	Status int
}

// captureScreenshot runs a series of browser actions, including navigating to url,
// and takes a screenshot of the resulting webpage in an instance of headless chrome.
func (tc *testcase) captureScreenshot(ctx context.Context, url string) ([]byte, error) {
	var buf []byte
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, time.Minute)
	defer cancel()
	var tasks chromedp.Tasks
	if len(tc.headers) > 0 {
		tasks = append(tasks, network.SetExtraHTTPHeaders(tc.headers))
	}
	if tc.blockedURLs != nil {
		tasks = append(tasks, network.SetBlockedURLS(tc.blockedURLs))
	}
	var res response
	tasks = append(tasks,
		getResponse(url, &res),
		chromedp.EmulateViewport(int64(tc.viewportWidth), int64(tc.viewportHeight)),
		chromedp.Navigate(url),
		waitForEvent("networkIdle"),
		reduceMotion(),
		checkResponse(tc, &res),
		tc.tasks,
	)
	switch tc.screenshotType {
	case fullScreenshot:
		tasks = append(tasks, chromedp.FullScreenshot(&buf, 100))
	case viewportScreenshot:
		tasks = append(tasks, chromedp.CaptureScreenshot(&buf))
	case elementScreenshot:
		tasks = append(tasks, chromedp.Screenshot(tc.screenshotElement, &buf))
	}
	if err := chromedp.Run(ctx, tasks); err != nil {
		return nil, fmt.Errorf("chromedp.Run(...): %w", err)
	}
	return buf, nil
}

// reduceMotion returns a chromedp action that will minimize motion during a screen capture.
func reduceMotion() chromedp.Action {
	css := `*, ::before, ::after {
		animation-delay: -1ms !important;
		animation-duration: 1ms !important;
		animation-iteration-count: 1 !important;
		background-attachment: initial !important;
		caret-color: transparent;
		scroll-behavior: auto !important;
		transition-duration: 0s !important;
		transition-delay: 0s !important;
	}`
	script := `
	(() => {
		const style = document.createElement('style');
		style.type = 'text/css';
		style.appendChild(document.createTextNode(` + "`" + css + "`" + `));
		document.head.appendChild(style);
	})()
	`
	return chromedp.Evaluate(script, nil)
}

// writePNG writes image data to a png file.
func writePNG(i *image.Image, filename string) error {
	var f io.WriteCloser
	if strings.HasPrefix(filename, gcsScheme) {
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("storage.NewClient(ctx): %w", err)
		}
		defer client.Close()
		bkt, obj := gcsParts(filename)
		f = client.Bucket(bkt).Object(obj).NewWriter(ctx)
	} else {
		var err error
		f, err = os.Create(filename)
		if err != nil {
			return err
		}
	}
	if err := png.Encode(f, *i); err != nil {
		// Ignore f.Close() error, since png.Encode returned an error.
		_ = f.Close()
		return fmt.Errorf("png.Encode(...): %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("f.Close(): %w", err)
	}
	return nil
}

var sanitizeRegexp = regexp.MustCompile("[.*<>?`'|/\\: ]")

// sanitize transforms text into a string suitable for use in a
// filename part.
func sanitize(text string) string {
	return sanitizeRegexp.ReplaceAllString(text, "-")
}

// gcsParts splits a Cloud Storage filename into bucket name and
// object name parts.
func gcsParts(filename string) (bucket, object string) {
	filename = strings.TrimPrefix(filename, gcsScheme)
	n := strings.Index(filename, "/")
	bucket = filename[:n]
	object = filename[n+1:]
	return bucket, object
}

// waitForEvent waits for browser lifecycle events. This is useful for
// ensuring the page is fully loaded before capturing screenshots.
func waitForEvent(eventName string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		ch := make(chan struct{})
		cctx, cancel := context.WithCancel(ctx)
		chromedp.ListenTarget(cctx, func(ev any) {
			switch e := ev.(type) {
			case *page.EventLifecycleEvent:
				if e.Name == eventName {
					cancel()
					close(ch)
				}
			}
		})
		select {
		case <-ch:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func getResponse(u string, res *response) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		chromedp.ListenTarget(ctx, func(ev any) {
			// URL fragments are dropped in request targets so we must strip the fragment
			// from the URL to make a comparison.
			_u, _ := url.Parse(u)
			_u.Fragment = ""
			switch e := ev.(type) {
			case *network.EventResponseReceived:
				if e.Response.URL == _u.String() {
					res.Status = int(e.Response.Status)
				}
			// Capture the status from a redirected response.
			case *network.EventRequestWillBeSent:
				if e.RedirectResponse != nil && e.RedirectResponse.URL == _u.String() {
					res.Status = int(e.RedirectResponse.Status)
				}
			}
		})
		return nil
	}
}

func checkResponse(tc *testcase, res *response) chromedp.ActionFunc {
	return func(context.Context) error {
		if res.Status != tc.status {
			fmt.Fprintf(&tc.output, "\nFAIL http status mismatch: got %d; want %d", res.Status, tc.status)
			return fmt.Errorf("bad status: %d", res.Status)
		}
		return nil
	}
}

// runConcurrently calls f on each integer from 0 to n-1,
// with at most max invocations active at once.
// It waits for all invocations to complete.
func runConcurrently(n, max int, f func(int)) {
	tokens := make(chan struct{}, max)
	for i := 0; i < n; i++ {
		i := i
		tokens <- struct{}{} // wait until the number of goroutines is below the limit
		go func() {
			f(i)
			<-tokens // let another goroutine run
		}()
	}
	// Wait for all goroutines to finish.
	for i := 0; i < cap(tokens); i++ {
		tokens <- struct{}{}
	}
}
