// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO(jba): sleep directive
// TODO(jba): specify percent of image that may differ
// TODO(jba): remove ints function in template (see cmd/golangorg/testdata/screentest/relnotes.txt)
// TODO(jba): write index.html to outdir with a nice view of all the failures
// TODO(jba): debug -run regexp matching

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
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
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
	now := time.Now()

	if testURL == "" {
		return errors.New("missing URL or path to test")
	}
	if wantURL == "" {
		return errors.New("missing URL or path with expected results")
	}
	if _, err := url.Parse(testURL); err != nil {
		return err
	}
	if _, err := url.Parse(wantURL); err != nil {
		return err
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

	c, err := commonValues(ctx, testURL, wantURL, opts)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for _, file := range files {
		tests, err := readTests(file, testURL, wantURL, c)
		if err != nil {
			return err
		}
		if len(tests) == 0 && opts.run == "" {
			return fmt.Errorf("no tests found in %q", file)
		}

		// Remove fail directory and all contents, to avoid confusion with previous runs.
		if err := c.failImageWriter.rmdir(ctx, fileDir(file)); err != nil {
			return err
		}

		ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
		defer cancel()

		if opts.maxConcurrency < 1 {
			opts.maxConcurrency = 1
		}
		var hdr bool
		runConcurrently(len(tests), opts.maxConcurrency, func(i int) {
			tc := tests[i]
			if err := tc.run(ctx, opts.update); err != nil {
				if !hdr {
					fmt.Fprintf(&buf, "%s\n\n", file)
					hdr = true
				}
				fmt.Fprintf(&buf, "%v\n", err)
				fmt.Fprintf(&buf, "inspect diff at %s\n\n", path.Join(tc.failImageWriter.path(), tc.diffPath))
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

const (
	browserWidth  = 1536
	browserHeight = 960
)

type screenshotType int

const (
	fullScreenshot screenshotType = iota
	viewportScreenshot
	elementScreenshot
)

// common contains values common to all test files.
type common struct {
	testImageReader     imageReader       // read images for comparison or update
	wantImageReadWriter imageReadWriter   // read images for comparison, write them for update
	failImageWriter     imageWriter       // write images for failed tests
	headers             map[string]any    // any to match chromedp arg
	filter              func(string) bool // filter out tests by name, from -run flag
	vars                map[string]string // variables for template execution
}

// testcase is a test case.
type testcase struct {
	common
	name               string // name of the test (arg to 'test' directive)
	tasks              chromedp.Tasks
	testURL, wantURL   string // URL to visit if the command-line arg is http/https
	testPath, wantPath string // slash-separated path to use if the command-line arg is file, gs or a path
	diffPath           string // output path for failed tests
	status             int
	viewportWidth      int
	viewportHeight     int
	screenshotType     screenshotType
	screenshotElement  string
	blockedURLs        []string
	output             bytes.Buffer
}

func (t *testcase) String() string {
	return t.name
}

// commonValues returns values common to all test files.
func commonValues(ctx context.Context, testURL, wantURL string, opts options) (c common, err error) {
	// The test/want image readers/writers are relative to the test/want URLs, so
	// they are common to all files. See test/wantPath for the file- and test-relative components.
	// They may be nil if a URL has an http or https scheme.
	c.testImageReader, err = newImageReadWriter(ctx, testURL)
	if err != nil {
		return common{}, err
	}
	c.wantImageReadWriter, err = newImageReadWriter(ctx, wantURL)
	if err != nil {
		return common{}, err
	}

	outDirPath := opts.outputDirURL
	if outDirPath == "" {
		cache, err := os.UserCacheDir()
		if err != nil {
			return common{}, fmt.Errorf("os.UserCacheDir(): %w", err)
		}
		outDirPath = path.Join(filepath.ToSlash(cache), "screentest")
	}
	c.failImageWriter, err = newImageReadWriter(ctx, outDirPath)
	if err != nil {
		return common{}, err
	}
	if c.failImageWriter == nil {
		return common{}, fmt.Errorf("cannot write images to %q", outDirPath)
	}

	hs, err := splitList(opts.headers)
	if err != nil {
		return common{}, err
	}
	if len(hs) > 0 {
		c.headers = map[string]any{}
		for k, v := range hs {
			c.headers[k] = v
		}
	}

	c.filter = func(string) bool { return true }
	if opts.run != "" {
		re, err := regexp.Compile(opts.run)
		if err != nil {
			return common{}, err
		}
		c.filter = re.MatchString
	}

	c.vars, err = splitList(opts.vars)
	if err != nil {
		return common{}, err
	}

	return c, nil
}

// readTests parses the testcases from a text file.
func readTests(file, testURL, wantURL string, common common) (_ []*testcase, err error) {
	// Test files are templates, so first execute them.
	data, err := executeFileTemplate(file, common.vars)
	if err != nil {
		return nil, err
	}
	var tests []*testcase
	testNames := map[string]bool{}
	var (
		testName, pathname string
		tasks              chromedp.Tasks
		status             int = http.StatusOK
		width, height      int
		lineNo             int
		blockedURLs        []string
		// URLs for HTTP(s) requests
		testRequestURL string
		wantRequestURL string
	)

	defer wrapf(&err, "%s:%d", file, lineNo)

	scan := bufio.NewScanner(bytes.NewReader(data))
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
			return nil, fmt.Errorf("the %q directive should only occur in a test", strings.ToLower(field))
		}
		switch field {
		case "":
			// We've reached an empty line, reset properties scoped to a single test.
			testName, pathname = "", ""
			tasks = nil
			status = http.StatusOK
			testRequestURL, wantRequestURL = "", ""

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
			if testNames[testName] {
				return nil, fmt.Errorf("duplicate test name %q", testName)
			}
			testNames[testName] = true

		case "PATHNAME":
			pathname = args
			// If there is no imageReader, then assume the URL is an http(s) URL.
			if common.testImageReader == nil {
				u, err := url.JoinPath(testURL, pathname)
				if err != nil {
					return nil, err
				}
				testRequestURL = u
			}
			if common.wantImageReadWriter == nil {
				u, err := url.JoinPath(wantURL, pathname)
				if err != nil {
					return nil, err
				}
				wantRequestURL = u
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
			if !common.filter(testName) {
				continue
			}
			test := &testcase{
				common:         common,
				name:           testName,
				tasks:          tasks,
				status:         status,
				testURL:        testRequestURL,
				wantURL:        wantRequestURL,
				blockedURLs:    blockedURLs,
				screenshotType: viewportScreenshot, // default to viewportScreenshot
				viewportWidth:  width,
				viewportHeight: height,
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
			filePath := filepath.ToSlash(fileDir(file))
			fnPath := path.Join(filePath, sanitize(test.name))
			test.testPath = fnPath + ".got.png"
			test.wantPath = fnPath + ".want.png"
			test.diffPath = fnPath + ".diff.png"

		default:
			return nil, fmt.Errorf("unknown directive %q", field)
		}
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

// executeFileTemplate reads file and executes it with text/template, passing vars as the argument.
func executeFileTemplate(file string, vars map[string]string) ([]byte, error) {
	tmpl := template.New(filepath.Base(file)).Funcs(template.FuncMap{
		"ints": func(start, end int) []int {
			var out []int
			for i := start; i < end; i++ {
				out = append(out, i)
			}
			return out
		},
	})

	if _, err := tmpl.ParseFiles(file); err != nil {
		return nil, fmt.Errorf("%s: could not parse template: %w", file, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// fileDir returns the output directory for a test filename, which is the filename
// base without the extension.
func fileDir(filename string) string {
	base := filepath.Base(filename)
	return strings.TrimSuffix(base, filepath.Ext(base))
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
	var testScreen, wantScreen image.Image
	g, ctx := errgroup.WithContext(ctx)
	// If the hosts are the same, chrome (or chromedp) does not handle concurrent requests well.
	// This wouldn't make sense in an actual test, but it does happen in this package's tests.
	urla, erra := url.Parse(tc.testURL)
	urlb, errb := url.Parse(tc.wantURL)
	if err := cmp.Or(erra, errb); err != nil {
		return err
	}
	if urla.Host == urlb.Host {
		g.SetLimit(1)
	}

	g.Go(func() error {
		testScreen, err = tc.screenshot(ctx, tc.testURL, tc.testPath, tc.testImageReader)
		return err
	})
	if !update {
		g.Go(func() error {
			wantScreen, err = tc.screenshot(ctx, tc.wantURL, tc.wantPath, tc.wantImageReadWriter)
			return err
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Fprint(&tc.output, "\n", err)
		return err
	}

	// Update means overwrite the golden with the test result.
	if update {
		fmt.Fprintf(&tc.output, "updating %s\n", tc.wantURL)
		return tc.wantImageReadWriter.writeImage(ctx, tc.wantPath, testScreen)
	}

	result := imgdiff.Diff(testScreen, wantScreen, &imgdiff.Options{
		Threshold: 0.1,
		DiffImage: true,
	})
	since := time.Since(now).Truncate(time.Millisecond)
	if result.Equal {
		fmt.Fprintf(&tc.output, "(%s)\n", since)
		return nil
	}
	fmt.Fprintf(&tc.output, "(%s)\nFAIL %s != %s (%d pixels differ)\n", since, tc.testURL, tc.wantURL, result.DiffPixelsCount)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error { return tc.failImageWriter.writeImage(gctx, tc.testPath, testScreen) })
	g.Go(func() error { return tc.failImageWriter.writeImage(gctx, tc.wantPath, wantScreen) })
	g.Go(func() error { return tc.failImageWriter.writeImage(gctx, tc.diffPath, result.Image) })
	if err := g.Wait(); err != nil {
		return err
	}
	fmt.Fprintf(&tc.output, "wrote diff to %s\n", path.Join(tc.failImageWriter.path(), tc.diffPath))
	return fmt.Errorf("%s != %s", tc.testURL, tc.wantURL)
}

// screenshot gets a screenshot for a testcase url. If reader is non-nil
// it reads the pathname from reader. Otherwise it captures a new screenshot from url.
func (tc *testcase) screenshot(ctx context.Context, url, pathname string, reader imageReader) (image.Image, error) {
	if reader != nil {
		return reader.readImage(ctx, pathname)
	} else {
		data, err := tc.captureScreenshot(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("captureScreenshot(ctx, %q, %q): %w", url, tc, err)
		}
		img, _, err := image.Decode(bytes.NewReader(data))
		return img, err
	}
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

var sanitizeRegexp = regexp.MustCompile("[.*<>?`'|/\\: ]")

// sanitize transforms text into a string suitable for use in a
// filename part.
func sanitize(text string) string {
	return sanitizeRegexp.ReplaceAllString(text, "-")
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

// An imageReader reads images from slash-separated paths.
type imageReader interface {
	readImage(ctx context.Context, path string) (image.Image, error) // get an image with the given name
}

// An imageWriter writes images to slash-separated paths.
type imageWriter interface {
	writeImage(ctx context.Context, path string, img image.Image) error
	rmdir(ctx context.Context, path string) error
	path() string // return the slash-separated path that this was created with
}

type imageReadWriter interface {
	imageReader
	imageWriter
}

var validSchemes = []string{"file", "gs", "http", "https"}

// newImageReadWriter returns an imageReadWriter for loc.
// loc can be a URL with a scheme or a slash-separated file path.
func newImageReadWriter(ctx context.Context, loc string) (imageReadWriter, error) {
	scheme, _, _ := strings.Cut(loc, ":")
	scheme = strings.ToLower(scheme)
	switch scheme {
	case "http", "https":
		return nil, nil
	case "file", "gs":
		u, err := url.Parse(loc)
		if err != nil {
			return nil, err
		}
		if scheme == "file" {
			return &dirImageReadWriter{dir: path.Clean(u.Path)}, nil
		}
		return newGCSImageReadWriter(ctx, loc)
	default:
		// Assume a file path; Windows paths can start with a drive letter.
		return &dirImageReadWriter{dir: path.Clean(loc)}, nil
	}
}

// A dirImageReadWriter reads and writes images to a filesystem directory.
// dir should be slash-separated.
type dirImageReadWriter struct {
	dir string
}

func (rw *dirImageReadWriter) readImage(_ context.Context, path string) (_ image.Image, err error) {
	path = rw.nativePathname(path)
	defer wrapf(&err, "reading image from %s", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func (rw *dirImageReadWriter) writeImage(_ context.Context, path string, img image.Image) (err error) {
	path = rw.nativePathname(path)
	defer wrapf(&err, "writing %s", path)

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()
	return png.Encode(f, img)
}

func (rw *dirImageReadWriter) rmdir(_ context.Context, path string) error {
	return os.RemoveAll(rw.nativePathname(path))
}

func (rw *dirImageReadWriter) nativePathname(pth string) string {
	spath := path.Join(rw.dir, pth)
	return filepath.FromSlash(spath)
}

func (rw *dirImageReadWriter) path() string {
	return rw.dir
}

type gcsImageReadWriter struct {
	url    string // URL with scheme "gs" referring to bucket and prefix
	bucket *storage.BucketHandle
	prefix string // initial path of objects; effectively a directory
}

func newGCSImageReadWriter(ctx context.Context, urlstr string) (*gcsImageReadWriter, error) {
	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(urlstr)
	if err != nil {
		return nil, err
	}
	return &gcsImageReadWriter{
		url:    urlstr,
		bucket: c.Bucket(u.Host),
		prefix: u.Path[1:], //remove initial slash
	}, nil
}

func (rw *gcsImageReadWriter) readImage(ctx context.Context, pth string) (_ image.Image, err error) {
	defer wrapf(&err, "reading %s", path.Join(rw.url, pth))

	r, err := rw.bucket.Object(rw.objectName(pth)).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	img, _, err := image.Decode(r)
	return img, err
}

func (rw *gcsImageReadWriter) writeImage(ctx context.Context, pth string, img image.Image) (err error) {
	defer wrapf(&err, "writing %s", path.Join(rw.url, pth))

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	w := rw.bucket.Object(rw.objectName(pth)).NewWriter(cctx)
	if err := png.Encode(w, img); err != nil {
		cancel()
		_ = w.Close()
		return err
	}
	return w.Close()
}

func (rw *gcsImageReadWriter) rmdir(ctx context.Context, pth string) (err error) {
	defer wrapf(&err, "rmdir %s", path.Join(rw.url, pth))

	prefix := path.Join(rw.prefix, pth)
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	it := rw.bucket.Objects(ctx, &storage.Query{
		Prefix:                   prefix,
		Projection:               storage.ProjectionNoACL,
		IncludeTrailingDelimiter: true,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("iterating over %s: %w", rw.url, err)
		}
		err = rw.bucket.Object(attrs.Name).Delete(ctx)
		if err != nil && !errors.Is(err, storage.ErrObjectNotExist) {
			return fmt.Errorf("deleting %q from bucket %q: %w", attrs.Name, attrs.Bucket, err)
		}
	}
	return nil
}

func (rw *gcsImageReadWriter) path() string { return rw.url }

func (rw *gcsImageReadWriter) objectName(pth string) string {
	// The relevant parts of pth have already been sanitized.
	return path.Join(rw.prefix, pth)
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

// wrapf prepends a non-nil *errp with the given message, formatted by fmt.Sprintf.
func wrapf(errp *error, format string, args ...any) {
	if *errp != nil {
		*errp = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *errp)
	}
}
