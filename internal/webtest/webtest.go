// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package webtest implements script-based testing for web servers.
//
// The scripts, described below, can be run against http.Handler
// implementations or against running servers. Testing against an
// http.Handler makes it easier to test handlers serving multiple sites
// as well as scheme-based features like redirecting to HTTPS.
// Testing against a running server provides a more complete end-to-end test.
//
// The test functions TestHandler and TestServer take a *testing.T
// and a glob pattern, which must match at least one file.
// They create a subtest of the top-level test for each script.
// Within each per-script subtest, they create a per-case subtest
// for each case in the script, making it easy to run selected cases.
//
// The functions CheckHandler and CheckServer are similar but do
// not require a *testing.T, making them suitable for use in other contexts.
// They run the entire script and return a multiline error summarizing
// any problems.
//
// # Scripts
//
// A script is a text file containing a sequence of cases, separated by blank lines.
// Lines beginning with # characters are ignored as comments.
// A case is a sequence of lines describing a request, along with checks to be
// applied to the response. For example, here is a trivial script:
//
//	GET /
//	body contains Go is an open source programming language
//
// This script has a single case. The first line describes the request.
// The second line describes a single check to be applied to the response.
// In this case, the request is a GET of the URL /, and the response body
// must contain the text “Go is an open source programming language”.
//
// # Requests
//
// Each case begins with a line starting with GET, HEAD, or POST.
// The argument (the remainder of the line) is the URL to be used in the request.
// Following this line, the request can be further customized using
// lines of the form
//
//	<verb> <text>
//
// where the verb is a single space-separated word and the text is arbitrary text
// to the end of the line, or multiline text (described below).
//
// The possible values for <verb> are as follows.
//
// The verb “hint” specifies text to be printed if the test case fails, as a
// hint about what might be wrong.
//
// The verbs “postbody”, “postquery”, and “posttype” customize a POST request.
//
// For example:
//
//	POST /api
//	posttype application/json
//	postbody {"go": true}
//
// This describes a POST request with a posted Content-Type of “application/json”
// and a body “{"go": true}”.
//
// The “postquery” verb specifies a post body in the form of a sequence of
// key-value pairs, query-encoded and concatenated automatically as a
// convenience. Using “postquery” also sets the default posted Content-Type
// to “application/x-www-form-urlencoded”.
//
// For example:
//
//	POST /api
//	postquery
//		x=hello world
//		y=Go & You
//
// This stanza sends a request with post body “x=hello+world&y=Go+%26+You”.
// (The multiline syntax is described in detail below.)
//
// # Checks
//
// By default, a stanza like the ones above checks only that the request
// succeeds in returning a response with HTTP status code 200 (OK).
// Additional checks are specified by more lines of the form
//
//	<value> [<key>] <op> <text>
//
// In the example above, <value> is “body”, there is no <key>,
// <op> is “contains”, and <text> is “Go is an open source programming language”.
// Whether there is a <key> depends on the <value>; “body” does not have one.
//
// The possible values for <value> are:
//
//	body - the full response body
//	code - the HTTP status code
//	header <key> - the value in the header line with the given key
//	redirect - the target of a redirect, as found in the Location header
//	trimbody - the response body, trimmed
//
// If a case contains no check of “code”, then it defaults to checking that
// the HTTP status code is 200, as described above, with one exception:
// if the case contains a check of “redirect”, then the code is required to
// be a 30x code.
//
// The “trimbody” value is the body with all runs of spaces and tabs
// reduced to single spaces, leading and trailing spaces removed on
// each line, and blank lines removed.
//
// The possible operators for <op> are:
//
//	== - the value must be equal to the text
//	!= - the value must not be equal to the text
//	~  - the value must match the text interpreted as a regular expression
//	!~ - the value must not match the text interpreted as a regular expression
//	contains  - the value must contain the text as a substring
//	!contains - the value must not contain the text as a substring
//
// For example:
//
//	GET /change/75944e2e3a63
//	hint no change redirect - hg to git mapping not registered?
//	code == 302
//	redirect contains bdb10cf
//	body contains bdb10cf
//	body !contains UA-
//
//	GET /pkg/net/http/httptrace/
//	body ~ Got1xxResponse.*// Go 1\.11
//	body ~ GotFirstResponseByte func\(\)\s*$
//
// # Multiline Texts
//
// The <text> in a request or check line can take a multiline form,
// by omitting it from the original line and then specifying the text
// as one or more following lines, each indented by a single tab.
// The text is taken to be the sequence of indented lines, including
// the final newline, but with the leading tab removed from each.
//
// The “postquery” example above showed the multiline syntax.
// Another common use is for multiline “body” checks. For example:
//
//	GET /hello
//	body ==
//		<!DOCTYPE html>
//		hello, world
package webtest

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"
)

// HandlerWithCheck returns an http.Handler that responds to each request
// by running the test script files mached by glob against the handler h.
// If the tests pass, the returned http.Handler responds with status code 200.
// If they fail, it prints the details and responds with status code 503
// (service unavailable).
func HandlerWithCheck(h http.Handler, path string, fsys fs.FS, glob string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path {
			err := CheckHandler(fsys, glob, h)
			if err != nil {
				http.Error(w, "webtest.CheckHandler failed:\n"+err.Error()+"\n", http.StatusInternalServerError)
			} else {
				fmt.Fprintf(w, "ok\n")
			}
			return
		}
		h.ServeHTTP(w, r)
	})
}

// CheckHandler runs the test script files in fsys matched by glob
// against the handler h. If any errors are encountered,
// CheckHandler returns an error listing the problems.
func CheckHandler(fsys fs.FS, glob string, h http.Handler) error {
	return check(fsys, glob, func(c *case_) error { return c.runHandler(h) })
}

func check(fsys fs.FS, glob string, do func(*case_) error) error {
	files, err := fs.Glob(fsys, glob)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no files match %#q", glob)
	}
	var buf bytes.Buffer
	for _, file := range files {
		data, err := fs.ReadFile(fsys, file)
		if err != nil {
			fmt.Fprintf(&buf, "# %s\n%v\n", file, err)
			continue
		}
		script, err := parseScript(file, string(data))
		if err != nil {
			fmt.Fprintf(&buf, "# %s\n%v\n", file, err)
			continue
		}
		hdr := false
		for _, c := range script.cases {
			if err := do(c); err != nil {
				if !hdr {
					fmt.Fprintf(&buf, "# %s\n", file)
					hdr = true
				}
				fmt.Fprintf(&buf, "## %s %s\n", c.method, c.url)
				fmt.Fprintf(&buf, "%v\n", err)
			}
		}
	}
	if buf.Len() > 0 {
		return errors.New(buf.String())
	}
	return nil
}

// TestHandler runs the test script files matched by glob
// against the handler h.
func TestHandler(t *testing.T, glob string, h http.Handler) {
	test(t, glob, func(c *case_) error { return c.runHandler(h) })
}

func test(t *testing.T, glob string, do func(*case_) error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatalf("no files match %#q", glob)
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			script, err := parseScript(file, string(data))
			if err != nil {
				t.Fatal(err)
			}
			for _, c := range script.cases {
				t.Run(c.method+"/"+strings.TrimPrefix(c.url, "/"), func(t *testing.T) {
					if err := do(c); err != nil {
						t.Fatal(err)
					}
				})
			}
		})
	}
}

// A script is a parsed test script.
type script struct {
	cases []*case_
}

// A case_ is a single test case (GET/HEAD/POST) in a script.
type case_ struct {
	file      string
	line      int
	method    string
	url       string
	postbody  string
	postquery string
	posttype  string
	hint      string
	checks    []*cmpCheck
}

// A cmp is a single comparison (check) made against a test case.
type cmpCheck struct {
	file    string
	line    int
	what    string
	whatArg string
	op      string
	want    string
	wantRE  *regexp.Regexp
}

// runHandler runs a test case against the handler h.
func (c *case_) runHandler(h http.Handler) error {
	w := httptest.NewRecorder()
	r, err := c.newRequest(c.url)
	if err != nil {
		return err
	}
	h.ServeHTTP(w, r)
	return c.check(w.Result(), w.Body.String())
}

// runServer runs a test case against the server at address addr.
func (c *case_) runServer(addr string) error {
	baseURL := ""
	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		// addr is a base URL
		if !strings.HasSuffix(addr, "/") {
			addr += "/"
		}
		baseURL = addr
	} else {
		// addr is an HTTP proxy
		baseURL = "http://" + addr + "/"
	}

	// Build full URL for request.
	u := c.url
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = strings.TrimSuffix(baseURL, "/")
		if !strings.HasPrefix(c.url, "/") {
			u += "/"
		}
		u += c.url
	}
	req, err := c.newRequest(u)

	if err != nil {
		return fmt.Errorf("%s:%d: %s %s: %s", c.file, c.line, c.method, c.url, err)
	}
	tr := &http.Transport{}
	if !strings.HasPrefix(u, baseURL) {
		// If u does not begin with baseURL, then we're in the proxy case
		// and we try to tunnel the network activity through the proxy's address.
		proxyURL, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("invalid addr: %v", err)
		}
		tr.Proxy = func(*http.Request) (*url.URL, error) { return proxyURL, nil }
	}
	resp, err := tr.RoundTrip(req)
	if err != nil {
		return fmt.Errorf("%s:%d: %s %s: %s", c.file, c.line, c.method, c.url, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("%s:%d: %s %s: reading body: %s", c.file, c.line, c.method, c.url, err)
	}
	return c.check(resp, string(body))
}

// newRequest creates a new request for the case c,
// using the URL u.
func (c *case_) newRequest(u string) (*http.Request, error) {
	body := c.requestBody()
	r, err := http.NewRequest(c.method, u, body)
	if err != nil {
		return nil, err
	}
	typ := c.posttype
	if body != nil && typ == "" {
		typ = "application/x-www-form-urlencoded"
	}
	if typ != "" {
		r.Header.Set("Content-Type", typ)
	}
	return r, nil
}

// requestBody returns the body for the case's request.
func (c *case_) requestBody() io.Reader {
	if c.postbody == "" {
		return nil
	}
	return strings.NewReader(c.postbody)
}

// check checks the response against the comparisons for the case.
func (c *case_) check(resp *http.Response, body string) error {
	var msg bytes.Buffer
	for _, chk := range c.checks {
		what := chk.what
		if chk.whatArg != "" {
			what += " " + chk.whatArg
		}
		var value string
		switch chk.what {
		default:
			value = "unknown what: " + chk.what
		case "body":
			value = body
		case "trimbody":
			value = trim(body)
		case "code":
			value = fmt.Sprint(resp.StatusCode)
		case "header":
			value = resp.Header.Get(chk.whatArg)
		case "redirect":
			if resp.StatusCode/10 == 30 {
				value = resp.Header.Get("Location")
			}
		}

		switch chk.op {
		default:
			fmt.Fprintf(&msg, "%s:%d: unknown operator %s\n", chk.file, chk.line, chk.op)
		case "==":
			if value != chk.want {
				fmt.Fprintf(&msg, "%s:%d: %s = %q, want %q\n", chk.file, chk.line, what, value, chk.want)
			}
		case "!=":
			if value == chk.want {
				fmt.Fprintf(&msg, "%s:%d: %s == %q (but want !=)\n", chk.file, chk.line, what, value)
			}
		case "~":
			if !chk.wantRE.MatchString(value) {
				fmt.Fprintf(&msg, "%s:%d: %s does not match %#q (but should)\n\t%s\n", chk.file, chk.line, what, chk.want, indent(value))
			}
		case "!~":
			if chk.wantRE.MatchString(value) {
				fmt.Fprintf(&msg, "%s:%d: %s matches %#q (but should not)\n\t%s\n", chk.file, chk.line, what, chk.want, indent(value))
			}
		case "contains":
			if !strings.Contains(value, chk.want) {
				fmt.Fprintf(&msg, "%s:%d: %s does not contain %#q (but should)\n\t%s\n", chk.file, chk.line, what, chk.want, indent(value))
			}
		case "!contains":
			if strings.Contains(value, chk.want) {
				fmt.Fprintf(&msg, "%s:%d: %s contains %#q (but should not)\n\t%s\n", chk.file, chk.line, what, chk.want, indent(value))
			}
		}
	}
	if msg.Len() > 0 && c.hint != "" {
		fmt.Fprintf(&msg, "hint: %s\n", indent(c.hint))
	}

	if msg.Len() > 0 {
		return fmt.Errorf("%s:%d: %s %s\n%s", c.file, c.line, c.method, c.url, msg.String())
	}
	return nil
}

// trim returns a trimming of s, in which all runs of spaces and tabs have
// been collapsed to a single space, leading and trailing spaces have been
// removed from each line, and blank lines are removed entirely.
func trim(s string) string {
	s = regexp.MustCompile(`[ \t]+`).ReplaceAllString(s, " ")
	s = regexp.MustCompile(`(?m)(^ | $)`).ReplaceAllString(s, "")
	s = strings.TrimLeft(s, "\n")
	s = regexp.MustCompile(`\n\n+`).ReplaceAllString(s, "\n")
	return s
}

// indent indents text for formatting in a message.
func indent(text string) string {
	if text == "" {
		return "(empty)"
	}
	if text == "\n" {
		return "(blank line)"
	}
	text = strings.TrimRight(text, "\n")
	if text == "" {
		return "(blank lines)"
	}
	text = strings.ReplaceAll(text, "\n", "\n\t")
	return text
}

// parseScript parses the test script in text.
// Errors are reported as being from file, but file is not directly read.
func parseScript(file, text string) (*script, error) {
	var current struct {
		Case      *case_
		Multiline *string
	}
	script := new(script)
	lastLineWasBlank := true
	lineno := 0
	line := ""
	errorf := func(format string, args ...interface{}) error {
		if line != "" {
			line = "\n" + line
		}
		return fmt.Errorf("%s:%d: %v%s", file, lineno, fmt.Sprintf(format, args...), line)
	}
	for text != "" {
		lineno++
		prevLine := line
		line, text, _ = cut(text, "\n")
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimRight(line, " \t")
		if line == "" {
			lastLineWasBlank = true
			continue
		}
		what, args := splitOneField(line)

		// Add indented line to current multiline check, or else it ends.
		if what == "" {
			// Line is indented.
			if current.Multiline != nil {
				lastLineWasBlank = false
				*current.Multiline += args + "\n"
				continue
			}
			return nil, errorf("unexpected indented line")
		}

		// Multiline text is over; must be present.
		if current.Multiline != nil && *current.Multiline == "" {
			lineno--
			line = prevLine
			return nil, errorf("missing multiline text")
		}
		current.Multiline = nil

		// Look for start of new check.
		switch what {
		case "GET", "HEAD", "POST":
			if !lastLineWasBlank {
				return nil, errorf("missing blank line before start of case")
			}
			if args == "" {
				return nil, errorf("missing %s URL", what)
			}
			cas := &case_{method: what, url: args, file: file, line: lineno}
			script.cases = append(script.cases, cas)
			current.Case = cas
			lastLineWasBlank = false
			continue
		}

		if lastLineWasBlank || current.Case == nil {
			return nil, errorf("missing GET/HEAD/POST at start of check")
		}

		// Look for case metadata.
		var targ *string
		switch what {
		case "postbody":
			targ = &current.Case.postbody
		case "postquery":
			targ = &current.Case.postquery
		case "posttype":
			targ = &current.Case.posttype
		case "hint":
			targ = &current.Case.hint
		}
		if targ != nil {
			if strings.HasPrefix(what, "post") && current.Case.method != "POST" {
				return nil, errorf("need POST (not %v) for %v", current.Case.method, what)
			}
			if args != "" {
				*targ = args
			} else {
				current.Multiline = targ
			}
			continue
		}

		// Start a comparison check.
		chk := &cmpCheck{file: file, line: lineno, what: what}
		current.Case.checks = append(current.Case.checks, chk)
		switch what {
		case "body", "code", "redirect":
			// no WhatArg
		case "header":
			chk.whatArg, args = splitOneField(args)
			if chk.whatArg == "" {
				return nil, errorf("missing header name")
			}
		}

		// Opcode, with optional leading "not"
		chk.op, args = splitOneField(args)
		switch chk.op {
		case "==", "!=", "~", "!~", "contains", "!contains":
			// ok
		default:
			return nil, errorf("unknown check operator %q", chk.op)
		}

		if args != "" {
			chk.want = args
		} else {
			current.Multiline = &chk.want
		}
	}

	// Finish each case.
	// Compute POST body from POST query.
	// Check that each regexp compiles, and insert "code equals 200"
	// in each case that doesn't already have a code check.
	for _, cas := range script.cases {
		if cas.postquery != "" {
			if cas.postbody != "" {
				line = ""
				lineno = cas.line
				return nil, errorf("case has postbody and postquery")
			}
			for _, kv := range strings.Split(cas.postquery, "\n") {
				kv = strings.TrimSpace(kv)
				if kv == "" {
					continue
				}
				k, v, ok := cut(kv, "=")
				if !ok {
					lineno = cas.line // close enough
					line = kv
					return nil, errorf("postquery has non key=value line")
				}
				if cas.postbody != "" {
					cas.postbody += "&"
				}
				cas.postbody += url.QueryEscape(k) + "=" + url.QueryEscape(v)
			}
		}
		sawCode := false
		for _, chk := range cas.checks {
			if chk.what == "code" || chk.what == "redirect" {
				sawCode = true
			}
			if chk.op == "~" || chk.op == "!~" {
				re, err := regexp.Compile(`(?m)` + chk.want)
				if err != nil {
					lineno = chk.line
					line = chk.want
					return nil, errorf("invalid regexp: %s", err)
				}
				chk.wantRE = re
			}
		}
		if !sawCode {
			line := cas.line
			if len(cas.checks) > 0 {
				line = cas.checks[0].line
			}
			chk := &cmpCheck{file: cas.file, line: line, what: "code", op: "==", want: "200"}
			cas.checks = append(cas.checks, chk)
		}
	}
	return script, nil
}

// cut returns the result of cutting s around the first instance of sep.
func cut(s, sep string) (before, after string, ok bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

// cutAny returns the result of cutting s around the first instance of
// any code point from any.
func cutAny(s, any string) (before, after string, ok bool) {
	if i := strings.IndexAny(s, any); i >= 0 {
		_, size := utf8.DecodeRuneInString(s[i:])
		return s[:i], s[i+size:], true
	}
	return s, "", false
}

// splitOneField splits text at the first space or tab
// and returns that first field and the remaining text.
func splitOneField(text string) (field, rest string) {
	i := strings.IndexAny(text, " \t")
	if i < 0 {
		return text, ""
	}
	return text[:i], strings.TrimLeft(text[i:], " \t")
}
