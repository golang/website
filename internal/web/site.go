// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package web implements a basic web site serving framework.
// The two fundamental types in this package are Site and Page.
//
// # Sites
//
// A Site is an http.Handler that serves requests from a file system.
// Use NewSite(fsys) to create a new Site.
//
// The Site is defined primarily by the content of its file system fsys,
// which holds files to be served as well as templates for
// converting Markdown or HTML fragments into full HTML pages.
//
// # Pages
//
// A Page, which is a map[string]interface{}, is the raw data that a Site renders into a web page.
// Typically a Page is loaded from a *.html or *.md file in the file system fsys, although
// dynamic pages can be computed and passed to ServePage as well,
// as described in “Serving Dynamic Pages” below.
//
// For a Page loaded from the file system, the key-value pairs in the map
// are initialized from the YAML or JSON metadata block at the top of a Markdown or HTML file,
// which looks like (YAML):
//
//	---
//	key: value
//	...
//	---
//
// or (JSON):
//
//	<!--{
//		"Key": "value",
//		...
//	}-->
//
// By convention, key-value pairs loaded from a metadata block use lower-case keys.
// For historical reasons, keys in JSON metadata are converted to lower-case when read,
// so that the two headers above both refer to a key with a lower-case k.
//
// A few keys have special meanings:
//
// The key-value pair “status: n” sets the HTTP response status to the integer code n.
//
// The key-value pair “redirect: url” causes requests for this page redirect to the given
// relative or absolute URL.
//
// The key-value pair “layout: name” selects the page layout template with the given name.
// See the next section, “Page Rendering”, for details about layout and rendering.
//
// In addition to these explicit key-value pairs, pages loaded from the file system
// have a few implicit key-value pairs added by the page loading process:
//
//   - File: the path in fsys to the file containing the page
//   - FileData: the file body, with the key-value metadata stripped
//   - URL: this page's URL path (/x/y/z for x/y/z.md, /x/y/ for x/y/index.md)
//
// The key “Content” is added during the rendering process.
// See “Page Rendering” for details.
//
// # Page Rendering
//
// A Page's content is rendered in two steps: conversion to content, and framing of content.
//
// To convert a page to content, the page's file body (its FileData key, a []byte) is parsed
// and executed as an HTML template, with the page itself passed as the template input data.
// The template output is then interpreted as Markdown (perhaps with embedded HTML),
// and converted to HTML. The result is stored in the page under the key “Content”,
// with type template.HTML.
//
// A page's conversion to content can be skipped entirely in dynamically-generated pages
// by setting the “Content” key before passing the page to ServePage.
//
// The second step is framing the content in the overall site HTML, which is done by
// executing the site template, again using the Page itself as the template input data.
//
// The site template is constructed from two files in the file system.
// The first file is the fsys's “site.tmpl”, which provides the overall HTML frame for the site.
// The second file is a layout-specific template file, selected by the Page's
// “layout: name” key-value pair.
// The renderer searches for “name.tmpl” in the directory containing the page's file,
// then in the parent of that directory, and so on up to the root.
// If no such template is found, the rendering fails and reports that error.
// As a special case, “layout: none” skips the second file entirely.
//
// If there is no “layout: name” key-value pair, then the renderer tries using an
// implicit “layout: default”, but if no such “default.tmpl” template file can be found,
// the renderer uses an implicit “layout: none” instead.
//
// By convention, the site template and the layout-specific template are connected as follows.
// The site template, at the point where the content should be rendered, executes:
//
//	{{block "layout" .}}{{.Content}}{{end}}
//
// The layout-specific template overrides this block by defining its own template named “layout”.
// For example:
//
//	{{define "layout"}}
//	Here's some <blink>great</blink> content: {{.Content}}
//	{{end}}
//
// The use of the “block” template construct ensures that
// if there is no layout-specific template,
// the content will still be rendered.
//
// # Page Template Functions
//
// In this web server, templates can themselves be invoked as functions.
// See https://pkg.go.dev/rsc.io/tmplfunc for more details about that feature.
//
// During page rendering, both when rendering a page to content and when framing the content,
// the following template functions are available (in addition to those provided by the
// template package itself and the per-template functions just mentioned).
//
// In all functions taking a file path f, if the path begins with a slash,
// it is interpreted relative to the fsys root.
// Otherwise, it is interpreted relative to the directory of the current page's URL.
//
// The “{{add x y}}”, “{{sub x y}}”, “{{mul x y}}”, and “{{div x y}}” functions
// provide basic math on arguments of type int.
//
// The “{{code f [start [end]]}}” function returns a template.HTML of a formatted display
// of code lines from the file f.
// If both start and end are omitted, then the display shows the entire file.
// If only the start line is specified, then the display shows that single line.
// If both start and end are specified, then the display shows a range of lines
// starting at start up to and including end.
// The arguments start and end can take two forms: a number indicates a specific line number,
// and a string is taken to be a regular expression indicating the earliest matching line
// in the file (or, for end, the earliest matching line after the start line).
// Any lines ending in “OMIT” are elided from the display.
//
// For example:
//
//	{{code "hello.go" `^func main` `^}`}}
//
// The “{{data f}}” function reads the file f,
// decodes it as YAML, and then returns the resulting data,
// typically a map[string]interface{}.
// It is effectively shorthand for “{{yaml (file f)}}”.
//
// The “{{file f}}” function reads the file f and returns its content as a string.
//
// The “{{first n slice}}” function returns a slice of the first n elements of slice,
// or else slice itself when slice has fewer than n elements.
//
// The “{{markdown text}}” function interprets text (a string) as Markdown
// and returns the equivalent HTML as a template.HTML.
//
// The “{{page f}}” function returns the page data (a Page)
// for the static page contained in the file f.
// The lookup ignores trailing slashes in f as well as the presence or absence
// of extensions like .md, .html, /index.md, and /index.html,
// making it possible for f to be a relative or absolute URL path instead of a file path.
//
// The “{{pages glob}}” function returns a slice of page data (a []Page)
// for all pages loaded from files or directories
// in fsys matching the given glob (a string),
// according to the usual file path rules (if the glob starts with slash,
// it is interpreted relative to the fsys root, and otherwise
// relative to the directory of the page's URL).
// If the glob pattern matches a directory,
// the page for the directory's index.md or index.html is used.
//
// For example:
//
//	Here are all the articles:
//	{{range (pages "/articles/*")}}
//	- [{{.title}}]({{.URL}})
//	{{end}}
//
// The “{{raw s}}” function converts s (a string) to type template.HTML without any escaping,
// to allow using s as raw Markdown or HTML in the final output.
//
// The “{{yaml s}}” function decodes s (a string) as YAML and returns the resulting data.
// It is most useful for defining templates that accept YAML-structured data as a literal argument.
// For example:
//
//	{{define "quote info"}}
//	{{with (yaml .info)}}
//	.text
//	— .name{{if .title}}, .title{{end}}
//	{{end}}
//
//	{{quote `
//	  text: If a program is too slow, it must have a loop.
//	  name: Ken Thompson
//	`}}
//
// The “path” and “strings” functions return package objects with methods for every top-level
// function in these packages (except path.Split, which has more than one non-error result
// and would not be invokable). For example, “{{strings.ToUpper "abc"}}”.
//
// # Serving Requests
//
// A Site is an http.Handler that serves requests by consulting the underlying
// file system and constructing and rendering pages, as well as serving binary
// and text files.
//
// To serve a request for URL path /p, if fsys has a file
// p/index.md, p/index.html, p.md, or p.html
// (in that order of preference), then the Site opens that file,
// parses it into a Page, renders the page as described
// in the “Page Rendering” section above,
// and responds to the request with the generated HTML.
// If the request URL does not match the parsed page's URL,
// then the Site responds with a redirect to the canonical URL.
//
// Otherwise, if fsys has a directory p and the Site
// can find a template “dir.tmpl” in that directory or a parent,
// then the Site responds with the rendering of
//
//	Page{
//		"URL": "/p/",
//		"File": "p",
//		"layout": "dir",
//		"dir": []fs.FileInfo(dir),
//	}
//
// where dir is the directory contents.
//
// Otherwise, if fsys has a file p containing valid UTF-8 text
// (at least up to the first kilobyte of the file) and the Site
// can find a template “text.tmpl” in that file's directory or a parent,
// and the file is not named robots.txt,
// and the file does not have a .css, .js, .svg, or .ts extension,
// then the Site responds with the rendering of
//
//	Page{
//		"URL": "/p",
//		"File": "p",
//		"layout": "texthtml",
//		"texthtml": template.HTML(texthtml),
//	}
//
// where texthtml is the text file as rendered by the
// golang.org/x/website/internal/texthtml package.
// In the texthtml.Config, GoComments is set to true for
// file names ending in .go;
// the h URL query parameter, if present, is passed as Highlight,
// and the s URL query parameter, if set to lo:hi, is passed as a
// single-range Selection.
//
// If the request has the URL query parameter m=text,
// then the text file content is not rendered or framed and is instead
// served directly as a plain text response.
//
// If the request is for a file with a .ts extension the file contents
// are transformed from TypeScript to JavaScript and then served with
// a Content-Type=text/javascript header.
//
// Otherwise, if none of those cases apply but the request path p
// does exist in the file system, then the Site passes the
// request to an http.FileServer serving from fsys.
// This last case handles binary static content as well as
// textual static content excluded from the text file case above.
//
// Otherwise, the Site responds with the rendering of
//
//	Page{
//		"URL": r.URL.Path,
//		"status": 404,
//		"layout": "error",
//		"error": err,
//	}
//
// where err is the “not exist” error returned by fs.Stat(fsys, p).
// (See also the “Serving Errors” section below.)
//
// # Serving Dynamic Requests
//
// Of course, a web site may wish to serve more than static content.
// To allow dynamically generated web pages to make use of page
// rendering and site templates, the Site.ServePage method can be
// called with a dynamically generated Page value, which will then
// be rendered and served as the result of the request.
//
// # Serving Errors
//
// If an error occurs while serving a request r,
// the Site responds with the rendering of
//
//	Page{
//		"URL": r.URL.Path,
//		"status": 500,
//		"layout": "error",
//		"error": err,
//	}
//
// If that rendering itself fails, the Site responds with status 500
// and the cryptic page text “error rendering error”.
//
// The Site.ServeError and Site.ServeErrorStatus methods provide a way
// for dynamic servers to generate similar responses.
package web

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"golang.org/x/website/internal/backport/html/template"
	"golang.org/x/website/internal/spec"
	"golang.org/x/website/internal/texthtml"
)

// A Site is an http.Handler that serves requests from a file system.
// See the package doc comment for details.
type Site struct {
	fs         fs.FS            // from NewSite
	fileServer http.Handler     // http.FileServer(http.FS(fs))
	funcs      template.FuncMap // accumulated from s.Funcs
	cache      sync.Map         // canonical file path -> *pageFile, for site.openPage
}

// NewSite returns a new Site for serving pages from the file system fsys.
func NewSite(fsys fs.FS) *Site {
	return &Site{
		fs:         fsys,
		fileServer: http.FileServer(http.FS(fsys)),
	}
}

// Funcs adds the functions in m to the set of functions available to templates.
// Funcs must not be called concurrently with any page rendering.
func (s *Site) Funcs(m template.FuncMap) {
	if s.funcs == nil {
		s.funcs = make(template.FuncMap)
	}
	for k, v := range m {
		s.funcs[k] = v
	}
}

// readFile returns the content of the named file in the site's file system.
// If file begins with a slash, it is interpreted relative to the root of the file system.
// Otherwise, it is interpreted relative to dir.
func (site *Site) readFile(dir, file string) ([]byte, error) {
	if strings.HasPrefix(file, "/") {
		file = path.Clean(file)
	} else {
		file = path.Join(dir, file)
	}
	file = strings.Trim(file, "/")
	if file == "" {
		file = "."
	}
	return fs.ReadFile(site.fs, file)
}

// ServeError is ServeErrorStatus with HTTP status code 500 (internal server error).
func (s *Site) ServeError(w http.ResponseWriter, r *http.Request, err error) {
	s.ServeErrorStatus(w, r, err, http.StatusInternalServerError)
}

// ServeErrorStatus responds to the request
// with the given error and HTTP status.
// It is equivalent to calling ServePage(w, r, p) where p is:
//
//	Page{
//		"URL": r.URL.Path,
//		"status": status,
//		"layout": error,
//		"error": err,
//	}
func (s *Site) ServeErrorStatus(w http.ResponseWriter, r *http.Request, err error, status int) {
	s.serveErrorStatus(w, r, err, status, false)
}

func (s *Site) serveErrorStatus(w http.ResponseWriter, r *http.Request, err error, status int, renderingError bool) {

	if renderingError {
		log.Printf("error rendering error: %v", err)
		w.WriteHeader(status)
		w.Write([]byte("error rendering error"))
		return
	}

	p := Page{
		"URL":    r.URL.Path,
		"status": status,
		"layout": "error",
		"error":  err,
	}
	s.servePage(w, r, p, true)
}

// ServePage renders the page p to HTML and writes that HTML to w.
// See the package doc comment for details about page rendering.
//
// So that all templates can assume the presence of p["URL"],
// if p["URL"] is unset or does not have type string, then ServePage
// sets p["URL"] to r.URL.Path in a clone of p before rendering the page.
func (s *Site) ServePage(w http.ResponseWriter, r *http.Request, p Page) {
	s.servePage(w, r, p, false)
}

func (s *Site) servePage(w http.ResponseWriter, r *http.Request, p Page, renderingError bool) {
	html, err := s.renderHTML(p, "site.tmpl", r)
	if err != nil {
		s.serveErrorStatus(w, r, fmt.Errorf("template execution: %v", err), http.StatusInternalServerError, renderingError)
		return
	}
	if code, ok := p["status"].(int); ok {
		w.WriteHeader(code)
	}
	w.Write(html)
}

// ServeHTTP implements http.Handler, serving from a file in the site.
// See the Site type documentation for details about how requests are handled.
func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	abspath := r.URL.Path
	relpath := path.Clean(strings.TrimPrefix(abspath, "/"))

	// Is it a TypeScript file?
	if strings.HasSuffix(relpath, ".ts") {
		s.serveTypeScript(w, r)
		return
	}

	// Is it a page we can generate?
	if p, err := s.openPage(relpath); err == nil {
		if p.url != abspath {
			// Redirect to canonical path.
			status := http.StatusMovedPermanently
			if i, ok := p.page["status"].(int); ok {
				status = i
			}
			http.Redirect(w, r, p.url, status)
			return
		}
		// Serve from the actual filesystem path.
		s.serveHTML(w, r, p)
		return
	}

	// Is it a directory or file we can serve?
	info, err := fs.Stat(s.fs, relpath)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, fs.ErrNotExist) {
			status = http.StatusNotFound
		}
		s.ServeErrorStatus(w, r, err, status)
		return
	}

	// Serve directory.
	if info != nil && info.IsDir() {
		if _, ok := s.findLayout(relpath, "dir"); ok {
			if !maybeRedirect(w, r) {
				s.serveDir(w, r, relpath)
			}
			return
		}
	}

	// Serve text file.
	if isTextFile(s.fs, relpath) {
		if _, ok := s.findLayout(path.Dir(relpath), "texthtml"); ok {
			if !maybeRedirectFile(w, r) {
				s.serveText(w, r, relpath)
			}
			return
		}
	}

	// Serve raw bytes.
	s.fileServer.ServeHTTP(w, r)
}

func maybeRedirect(w http.ResponseWriter, r *http.Request) (redirected bool) {
	canonical := path.Clean(r.URL.Path)
	if !strings.HasSuffix(canonical, "/") {
		canonical += "/"
	}
	if r.URL.Path != canonical {
		url := *r.URL
		url.Path = canonical
		http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
		redirected = true
	}
	return
}

func maybeRedirectFile(w http.ResponseWriter, r *http.Request) (redirected bool) {
	c := path.Clean(r.URL.Path)
	c = strings.TrimRight(c, "/")
	if r.URL.Path != c {
		url := *r.URL
		url.Path = c
		http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
		redirected = true
	}
	return
}

func (s *Site) serveHTML(w http.ResponseWriter, r *http.Request, p *pageFile) {
	src, _ := p.page["FileData"].(string)
	filePath, _ := p.page["File"].(string)
	isMarkdown := strings.HasSuffix(filePath, ".md")

	// if it begins with "<!DOCTYPE " assume it is standalone
	// html that doesn't need the template wrapping.
	if strings.HasPrefix(src, "<!DOCTYPE ") {
		w.Write([]byte(src))
		return
	}

	// if it's the language spec, add tags to EBNF productions
	if strings.HasSuffix(filePath, "ref/spec.html") {
		var buf bytes.Buffer
		spec.Linkify(&buf, []byte(src))
		src = buf.String()
	}

	// Template is enabled always in Markdown.
	// It can only be disabled for HTML files.
	isTemplate, _ := p.page["template"].(bool)
	if !isTemplate && !isMarkdown {
		p.page["Content"] = template.HTML(src)
	}
	s.ServePage(w, r, p.page)
}

func (s *Site) serveDir(w http.ResponseWriter, r *http.Request, relpath string) {
	if maybeRedirect(w, r) {
		return
	}

	list, err := fs.ReadDir(s.fs, relpath)
	if err != nil {
		s.ServeError(w, r, err)
		return
	}

	var info []fs.FileInfo
	for _, d := range list {
		i, err := d.Info()
		if err == nil {
			info = append(info, i)
		}
	}

	s.ServePage(w, r, Page{
		"URL":    r.URL.Path,
		"File":   relpath,
		"layout": "dir",
		"dir":    info,
	})
}

func (s *Site) serveText(w http.ResponseWriter, r *http.Request, relpath string) {
	src, err := fs.ReadFile(s.fs, relpath)
	if err != nil {
		log.Printf("ReadFile: %s", err)
		s.ServeError(w, r, err)
		return
	}

	if r.FormValue("m") == "text" {
		s.serveRawText(w, src)
		return
	}

	cfg := texthtml.Config{
		GoComments: path.Ext(relpath) == ".go",
		Highlight:  r.FormValue("h"),
		Selection:  rangeSelection(r.FormValue("s")),
		Line:       1,
	}

	var buf bytes.Buffer
	buf.WriteString("<pre>")
	buf.Write(texthtml.Format(src, cfg))
	buf.WriteString("</pre>")

	fmt.Fprintf(&buf, `<p><a href="/%s?m=text">View as plain text</a></p>`, html.EscapeString(relpath))

	s.ServePage(w, r, Page{
		"URL":      r.URL.Path,
		"File":     relpath,
		"layout":   "texthtml",
		"texthtml": template.HTML(buf.String()),
	})
}

var selRx = regexp.MustCompile(`^([0-9]+):([0-9]+)`)

// rangeSelection computes the Selection for a text range described
// by the argument str, of the form Start:End, where Start and End
// are decimal byte offsets.
func rangeSelection(str string) texthtml.Selection {
	m := selRx.FindStringSubmatch(str)
	if len(m) >= 2 {
		from, _ := strconv.Atoi(m[1])
		to, _ := strconv.Atoi(m[2])
		if from < to {
			return texthtml.Spans(texthtml.Span{Start: from, End: to})
		}
	}
	return nil
}

func (s *Site) serveRawText(w http.ResponseWriter, text []byte) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(text)
}

const cacheHeader = "X-Go-Dev-Cache-Hit"

type jsout struct {
	output []byte
	stat   fs.FileInfo // stat for file when page was loaded
}

func (s *Site) serveTypeScript(w http.ResponseWriter, r *http.Request) {
	filename := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	if cjs, ok := s.cache.Load(filename); ok {
		js := cjs.(*jsout)
		info, err := fs.Stat(s.fs, filename)
		if err == nil && info.ModTime().Equal(js.stat.ModTime()) {
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			w.Header().Set(cacheHeader, "true")
			http.ServeContent(w, r, filename, info.ModTime(), bytes.NewReader(js.output))
			return
		}
	}
	file, err := s.fs.Open(filename)
	if err != nil {
		s.ServeError(w, r, err)
		return
	}
	var contents bytes.Buffer
	_, err = io.Copy(&contents, file)
	if err != nil {
		s.ServeError(w, r, err)
		return
	}
	result := api.Transform(contents.String(), api.TransformOptions{
		Loader: api.LoaderTS,
		Target: api.ES2018,
	})
	var buf bytes.Buffer
	for _, v := range result.Errors {
		fmt.Fprintln(&buf, v.Text)
	}
	if buf.Len() > 0 {
		s.ServeError(w, r, errors.New(buf.String()))
		return
	}
	info, err := file.Stat()
	if err != nil {
		s.ServeError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	http.ServeContent(w, r, filename, info.ModTime(), bytes.NewReader(result.Code))
	s.cache.Store(filename, &jsout{
		output: result.Code,
		stat:   info,
	})
}
