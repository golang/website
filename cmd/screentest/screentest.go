// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO(jba): remove ints function in template (see cmd/golangorg/testdata/screentest/relnotes.txt)

// TODO(jba): Provide a way to capture the results of an eval directive.
// If the second argument to chromedp.Evaluate is a *[]byte, the result will be written
// there. The problem is that we may screenshot twice, and currently the same list of
// tasks is used for both, and what's more the evaluations happen concurrently (see testcase.run).
// We need two locations for Evaluate, not one. Probably the simplest thing would be to build
// two copies of the tasks slice. But that is a lot of complexity for this one debugging feature.

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
	"strconv"
	"strings"
	"sync"
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
	start := time.Now()

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

	// Remove fail directory and all contents, to avoid confusion with previous runs.
	if err := c.failImageWriter.rmdir(ctx, "."); err != nil {
		return err
	}

	var (
		summary     bytes.Buffer
		nTests      int         // number of tests run
		failedTests []*testcase // tests that failed and wrote diffs
	)
	for _, file := range files {
		tests, err := readTests(file, testURL, wantURL, c)
		if err != nil {
			return err
		}

		if len(tests) == 0 {
			continue
		}
		nTests += len(tests)

		ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
		// TODO(jba): cancel after each iteration
		defer cancel()

		if opts.maxConcurrency < 1 {
			opts.maxConcurrency = 1
		}

		var (
			mu  sync.Mutex
			hdr bool
		)
		runConcurrently(len(tests), opts.maxConcurrency, func(i int) {
			tc := tests[i]
			if err := tc.run(ctx, opts.update); err != nil {
				mu.Lock()
				if !hdr {
					fmt.Fprintf(&summary, "%s\n", file)
					hdr = true
				}
				fmt.Fprintf(&summary, "%v\n", err)
				if tc.wroteDiff {
					failedTests = append(failedTests, tc)
				}
				mu.Unlock()
			}
			fmt.Println(tc.output.String())
		})
	}
	if nTests == 0 {
		log.Print("no tests to run")
		return nil
	}
	log.Printf("ran %d tests in %s\n", nTests, time.Since(start).Truncate(time.Millisecond))
	if summary.Len() > 0 {
		os.Stdout.Write(summary.Bytes())
		if len(failedTests) > 0 {
			data := failedImagesPage(failedTests)
			if err := c.failImageWriter.writeData(ctx, "index.html", data); err != nil {
				return err
			}
		}
		return fmt.Errorf("FAIL. Output at %s", c.failImageWriter.path())
	}
	if opts.update {
		log.Print("UPDATED")
	} else {
		log.Print("PASS")
	}
	return nil
}

// failedImagesPage builds a web page that displays the images for failed tests.
func failedImagesPage(failedTests []*testcase) []byte {
	var buf bytes.Buffer

	p := func(format string, args ...any) {
		fmt.Fprintf(&buf, format, args...)
	}

	p(`
    <html>
      <head>
        <style>
          table { border: 0 }
          /* image widths are 30% of the viewport width */
          img { width: 30vw }
          td { font-size: 1.2rem }
        </style>
      </head>
      <body>
        <h1>Screentest Failures</h1>
    `)
	for _, tc := range failedTests {
		p("<h2>%s</h2>\n", tc.name)
		p("<table>\n")
		p(`
        <tr>
          <td>Got: %s</td>
          <td>Want: %s</td>
          <td>Diff</td>
        </tr>
        `, tc.testOrigin(), tc.wantOrigin())
		p(`
        <tr valign=top>
          <td><img align='left' src='%s'/></td>
          <td><img align='left' src='%s'/></td>
          <td><img align='left' src='%s'/></td>
        </tr>
        `, tc.testPath, tc.wantPath, tc.diffPath)
		p("</table>\n")
	}
	p(`
      </body>
    </html>
    `)

	return buf.Bytes()
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
	retryPixels         int               // retry if difference <= this
}

// testcase is a test case.
type testcase struct {
	common
	name               string // name of the test (arg to 'test' directive)
	tasks              chromedp.Tasks
	path               string // path of URL to visit
	testURL, wantURL   string // URL to visit if the command-line arg is http/https
	testPath, wantPath string // slash-separated path to use if the command-line arg is file, gs or a path
	diffPath           string // output path for failed tests
	wroteDiff          bool   // test failed and diffPath was written
	status             int
	viewportWidth      int
	viewportHeight     int
	screenshotType     screenshotType
	screenshotElement  string
	blockedURLs        []string
	output             bytes.Buffer
}

// testOrigin returns the origin of the test image: either an http(s) URL or a
// storage path.
func (tc *testcase) testOrigin() string { return cmp.Or(tc.testURL, tc.testPath) }

// wantOrigin returns the origin of the want image: either an http(s) URL or a
// storage path.
func (tc *testcase) wantOrigin() string { return cmp.Or(tc.wantURL, tc.wantPath) }

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
	if opts.update && c.wantImageReadWriter == nil {
		return common{}, fmt.Errorf("cannot update a non-storage wantURL: %s", wantURL)
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
	if opts.filterRegexp != "" {
		re, err := regexp.Compile(opts.filterRegexp)
		if err != nil {
			return common{}, err
		}
		c.filter = re.MatchString
	}

	c.vars, err = splitList(opts.vars)
	if err != nil {
		return common{}, err
	}
	c.retryPixels = opts.retryPixels
	return c, nil
}

// readTests parses the testcases from a text file.
func readTests(file, testURL, wantURL string, common common) (_ []*testcase, err error) {
	// Test files are templates, so first execute them.
	data, err := executeFileTemplate(file, common.vars)
	if err != nil {
		return nil, err
	}
	var (
		tests         []*testcase
		test          *testcase // test currently being constructed
		width, height int       // from windowsize directive
		blockedURLs   []string  // from block directive
		lineNo        int
		lastDirective string
	)

	testNames := map[string]bool{} // to detect duplicates

	defer func() {
		wraperr(&err, "%s:%d", file, lineNo)
	}()

	scan := bufio.NewScanner(bytes.NewReader(data))
	for scan.Scan() {
		lineNo++
		line := strings.TrimSpace(scan.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimRight(line, " \t")
		origDirective, args := splitOneField(line)
		directive := strings.ToUpper(origDirective)
		switch directive {
		case "":
			// An empty line means the end of a test (if one is active).
			// A test must end with a capture.
			if test != nil && lastDirective != "CAPTURE" {
				return nil, errors.New("test does not end with capture")
			}
			test = nil

		case "WINDOWSIZE":
			width, height, err = splitDimensions(args)
			if err != nil {
				return nil, err
			}

		case "BLOCK":
			urls := strings.Fields(args)
			if test != nil {
				test.blockedURLs = append(test.blockedURLs, urls...)
			} else {
				blockedURLs = append(blockedURLs, urls...)
			}

		case "TEST":
			if test != nil {
				return nil, errors.New("no blank lines between tests")
			}
			test = &testcase{
				common:      common,
				name:        args,
				status:      http.StatusOK,
				blockedURLs: blockedURLs,
			}
			if testNames[test.name] {
				return nil, fmt.Errorf("duplicate test name %q", test.name)
			}
			testNames[test.name] = true

		case "STATUS":
			if test == nil {
				return nil, errors.New("directive must be in a test")
			}
			test.status, err = strconv.Atoi(args)
			if err != nil {
				return nil, fmt.Errorf("strconv.Atoi(%q): %w", args, err)
			}

		case "PATH":
			if test == nil {
				return nil, errors.New("directive must be in a test")
			}
			test.path = args
			// If there is no imageReader, then assume the URL is an http(s) URL.
			if common.testImageReader == nil {
				test.testURL = joinURL(testURL, test.path)
			}
			// Ditto for want.
			if common.wantImageReadWriter == nil {
				test.wantURL = joinURL(wantURL, test.path)
			}

		case "CLICK":
			if test == nil {
				return nil, errors.New("directive must be in a test")
			}
			test.tasks = append(test.tasks, chromedp.Click(args))

		case "WAIT":
			if test == nil {
				return nil, errors.New("directive must be in a test")
			}
			test.tasks = append(test.tasks, chromedp.WaitReady(args))

		case "EVAL":
			if test == nil {
				return nil, errors.New("directive must be in a test")
			}
			// Warn about a quoted argument to eval.
			// The quotes are not stripped, so JS sees a string, not an interesting
			// expression.
			// It's only a warning, not an error, because without more sophisticated
			// parsing we can't distinguish 'ab' from 'a' + 'b'.
			if len(args) >= 2 {
				s := args[0]
				e := args[len(args)-1]
				if (s == '\'' && e == '\'') || (s == '"' && e == '"') {
					fmt.Printf("WARNING: quoted argument %s to eval will evaluate to itself\n", args)
				}
			}
			test.tasks = append(test.tasks, chromedp.Evaluate(args, nil))

		case "SLEEP":
			if test == nil {
				return nil, errors.New("directive must be in a test")
			}
			dur, err := time.ParseDuration(args)
			if err != nil {
				return nil, err
			}
			test.tasks = append(test.tasks, chromedp.Sleep(dur))

		case "CAPTURE":
			if test == nil {
				return nil, errors.New("directive must be in a test")
			}
			if test.path == "" {
				return nil, errors.New("missing path")
			}
			test.screenshotType = viewportScreenshot // default to viewportScreenshot
			test.viewportWidth = width
			test.viewportHeight = height
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
			if common.filter(test.name) {
				tests = append(tests, test)
			}
			// Copy the test in case there's another capture directive.
			// This is safe because all non-shallow fields are append-only.
			clone := *test
			test = &clone

		default:
			return nil, fmt.Errorf("unknown directive %q", origDirective)
		}
		if directive != "" {
			lastDirective = directive
		}
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	if lastDirective != "CAPTURE" {
		return nil, errors.New("test file does not end with capture")
	}
	return tests, nil
}

// joinURL joins the left and right parts of a URL
// with a slash.
// Unlike url.JoinPath, it does not escape the URL path.
// Unlike path.Join, it does not remove doubled slashes.
func joinURL(left, right string) string {
	return strings.TrimSuffix(left, "/") + "/" + strings.TrimPrefix(right, "/")
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

// Maximum number of retries if diff <= retryPixels.
const maxRetries = 3

// run generates screenshots for a given test case and a diff if the
// screenshots do not match.
func (tc *testcase) run(ctx context.Context, update bool) (err error) {
	defer wraperr(&err, "test %s", tc.name)
	now := time.Now()
	var (
		since                  time.Duration
		result                 *imgdiff.Result
		testScreen, wantScreen image.Image
	)
	fmt.Fprintf(&tc.output, "test %s ", tc.name)
	var failReason string
	for try := 0; try < maxRetries; try++ {
		g, gctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			testScreen, err = tc.screenshot(gctx, tc.testURL, tc.testPath, tc.testImageReader)
			return err
		})
		if !update {
			g.Go(func() error {
				wantScreen, err = tc.screenshot(gctx, tc.wantURL, tc.wantPath, tc.wantImageReadWriter)
				return err
			})
		}
		if err := g.Wait(); err != nil {
			fmt.Fprint(&tc.output, "\n", err)
			return err
		}

		// Update means overwrite the golden with the test result.
		if update {
			fmt.Fprintf(&tc.output, "- updating %s", tc.wantURL)
			return tc.wantImageReadWriter.writeImage(ctx, tc.wantPath, testScreen)
		}

		// Expect the images to start at (0, 0).
		if p := testScreen.Bounds().Min; p != image.ZP {
			return fmt.Errorf("test image starts at %s, not (0, 0)", p)
		}
		if p := wantScreen.Bounds().Min; p != image.ZP {
			return fmt.Errorf("want image starts at %s, not (0, 0)", p)
		}
		// If the images are different sizes, don't even compare them. imgdiff does
		// not handle differently sized images properly.
		if tmax, wmax := testScreen.Bounds().Max, wantScreen.Bounds().Max; tmax != wmax {
			failReason = fmt.Sprintf("test image is %s but want image is %s", tmax, wmax)
			break
		}
		result = imgdiff.Diff(testScreen, wantScreen, &imgdiff.Options{
			Threshold: 0.1,
			DiffImage: true,
		})
		since = time.Since(now).Truncate(time.Millisecond)
		if result.Equal {
			fmt.Fprintf(&tc.output, "(%s)", since)
			return nil
		}
		failReason = fmt.Sprintf("%d pixels differ", result.DiffPixelsCount)
		if result.DiffPixelsCount > uint64(tc.retryPixels) {
			break
		}
		fmt.Fprintf(&tc.output, "difference is <= %d pixels\n", tc.retryPixels)
	}
	fmt.Fprintf(&tc.output, "(%s)\n    FAIL %s != %s: %s\n",
		since, tc.testOrigin(), tc.wantOrigin(), failReason)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error { return tc.failImageWriter.writeImage(gctx, tc.testPath, testScreen) })
	g.Go(func() error { return tc.failImageWriter.writeImage(gctx, tc.wantPath, wantScreen) })
	if result != nil && result.Image != nil {
		g.Go(func() error { return tc.failImageWriter.writeImage(gctx, tc.diffPath, result.Image) })
	}
	if err := g.Wait(); err != nil {
		return err
	}
	fmt.Fprintf(&tc.output, "    wrote diff to %s", path.Join(tc.failImageWriter.path(), tc.diffPath))
	tc.wroteDiff = true
	return fmt.Errorf("%s != %s", tc.testOrigin(), tc.wantOrigin())
}

// screenshot gets a screenshot for a testcase url. If reader is non-nil
// it reads the pathname from reader. Otherwise it captures a new screenshot from url.
func (tc *testcase) screenshot(ctx context.Context, url, pathname string, reader imageReader) (image.Image, error) {
	if reader != nil {
		return reader.readImage(ctx, pathname)
	} else {
		data, err := tc.captureScreenshot(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("captureScreenshot(ctx, %q, %q): %w", url, tc.name, err)
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
		closed := false
		cctx, cancel := context.WithCancel(ctx)
		defer cancel()
		chromedp.ListenTarget(cctx, func(ev any) {

			switch e := ev.(type) {
			case *page.EventLifecycleEvent:
				if e.Name == eventName {
					if !closed {
						close(ch)
						closed = true
					}
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
	writeData(ct context.Context, path string, data []byte) error
	rmdir(ctx context.Context, path string) error
	path() string // return the slash-separated path that this was created with
}

type imageReadWriter interface {
	imageReader
	imageWriter
}

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
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decoding image from %s: %w", path, err)
	}
	return img, nil
}

func (rw *dirImageReadWriter) writeImage(ctx context.Context, path string, img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}
	return rw.writeData(ctx, path, buf.Bytes())
}

func (rw *dirImageReadWriter) writeData(_ context.Context, path string, data []byte) (err error) {
	path = rw.nativePathname(path)
	defer wraperr(&err, "writing %s", path)

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
	_, err = f.Write(data)
	return err
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
	defer wraperr(&err, "reading %s", path.Join(rw.url, pth))

	r, err := rw.bucket.Object(rw.objectName(pth)).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	img, _, err := image.Decode(r)
	return img, err
}

func (rw *gcsImageReadWriter) writeImage(ctx context.Context, path string, img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}
	return rw.writeData(ctx, path, buf.Bytes())
}

func (rw *gcsImageReadWriter) writeData(ctx context.Context, pth string, data []byte) (err error) {
	defer wraperr(&err, "writing %s", path.Join(rw.url, pth))

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	w := rw.bucket.Object(rw.objectName(pth)).NewWriter(cctx)
	_, err = w.Write(data)
	return errors.Join(err, w.Close())
}

func (rw *gcsImageReadWriter) rmdir(ctx context.Context, pth string) (err error) {
	defer wraperr(&err, "rmdir %s", path.Join(rw.url, pth))

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

// wraperr prepends a non-nil *errp with the given message, formatted by fmt.Sprintf.
func wraperr(errp *error, format string, args ...any) {
	if *errp != nil {
		*errp = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *errp)
	}
}
