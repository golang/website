// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/website/internal/backport/html/template"
	"golang.org/x/website/internal/backport/httpfs"
	"golang.org/x/website/internal/backport/io/fs"
	"golang.org/x/website/internal/spec"
	"golang.org/x/website/internal/texthtml"
)

// Site is a website served from a file system.
type Site struct {
	fs fs.FS

	mux        *http.ServeMux
	fileServer http.Handler

	Templates *template.Template

	// GoogleAnalytics optionally adds Google Analytics via the provided
	// tracking ID to each page.
	GoogleAnalytics string

	docFuncs template.FuncMap
}

var siteFuncs = template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
	"mul": func(a, b int) int { return a * b },
	"div": func(a, b int) int { return a / b },

	"basename": path.Base,

	"split":      strings.Split,
	"join":       strings.Join,
	"hasPrefix":  strings.HasPrefix,
	"hasSuffix":  strings.HasSuffix,
	"trimPrefix": strings.TrimPrefix,
	"trimSuffix": strings.TrimSuffix,
}

// NewSite returns a new Presentation from a file system.
func NewSite(fsys fs.FS) (*Site, error) {
	p := &Site{
		fs:         fsys,
		mux:        http.NewServeMux(),
		fileServer: http.FileServer(httpfs.FS(fsys)),
	}
	p.mux.HandleFunc("/", p.serveFile)
	p.initDocFuncs()

	t, err := template.New("").Funcs(siteFuncs).ParseFS(fsys, "lib/godoc/*.html")
	if err != nil {
		return nil, err
	}
	p.Templates = t

	return p, nil
}

// ServeError responds to the request with the given error.
func (s *Site) ServeError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusNotFound)
	s.ServePage(w, r, Page{
		Title:    r.URL.Path,
		Template: "error.html",
		Data:     err,
	})
}

// ServeHTTP implements http.Handler, dispatching the request appropriately.
func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// ServePage responds to the request with the content described by page.
func (s *Site) ServePage(w http.ResponseWriter, r *http.Request, page Page) {
	page = s.fullPage(r, page)
	if d, ok := page.Data.(interface{ SetWebPage(*Page) }); ok {
		d.SetWebPage(&page)
	}

	if page.Template != "" {
		t := s.Templates.Lookup(page.Template)
		var buf bytes.Buffer
		if err := t.Execute(&buf, &page); err != nil {
			log.Printf("%s.Execute: %s", t.Name(), err)
		}
		page.HTML = template.HTML(buf.String())
	} else {
		page.HTML = page.Data.(template.HTML)
	}

	applyTemplateToResponseWriter(w, s.Templates.Lookup("site.html"), &page)
}

// A Page describes the contents of a webpage to be served.
//
// A Page's Methods are for use by the templates rendering the page.
type Page struct {
	Title    string // <h1>
	TabTitle string // prefix in <title>; defaults to Title
	Subtitle string // subtitle (date for spec, memory model)
	SrcPath  string // path to file in /src for text view

	// Template and Data describe the data to be
	// rendered into the overall site frame template.
	// If Template is empty, then Data should be a template.HTML
	// holding raw HTML to render into the site frame.
	// Otherwise, Template should be the name of a template file
	// in _content/lib/godoc (for example, "package.html"),
	// and that template will be executed
	// (with the *Page as its data argument) to produce HTML.
	//
	// The overall site template site.html is also invoked with
	// the *Page as its data argument. It is what arranges to call Template.
	Template string      // template to apply to data (empty string when Data is raw template.HTML)
	Data     interface{} // data to be rendered into page frame

	HTML template.HTML

	// Filled in automatically by ServePage
	GoogleCN        bool   // served on golang.google.cn
	GoogleAnalytics string // Google Analytics tag
	Version         string
	Site            *Site
}

// fullPage returns a copy of page with the “automatic” fields filled in.
func (s *Site) fullPage(r *http.Request, page Page) Page {
	if page.TabTitle == "" {
		page.TabTitle = page.Title
	}
	page.Version = runtime.Version()
	page.GoogleCN = GoogleCN(r)
	page.GoogleAnalytics = s.GoogleAnalytics
	page.Site = s
	return page
}

type writeErrorSaver struct {
	w   io.Writer
	err error
}

func (w *writeErrorSaver) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	if err != nil {
		w.err = err
	}
	return n, err
}

// applyTemplateToResponseWriter uses an http.ResponseWriter as the io.Writer
// for the call to template.Execute.  It uses an io.Writer wrapper to capture
// errors from the underlying http.ResponseWriter.  Errors are logged only when
// they come from the template processing and not the Writer; this avoid
// polluting log files with error messages due to networking issues, such as
// client disconnects and http HEAD protocol violations.
func applyTemplateToResponseWriter(rw http.ResponseWriter, t *template.Template, data interface{}) {
	w := &writeErrorSaver{w: rw}
	err := t.Execute(w, data)
	// There are some cases where template.Execute does not return an error when
	// rw returns an error, and some where it does.  So check w.err first.
	if w.err == nil && err != nil {
		// Log template errors.
		log.Printf("%s.Execute: %s", t.Name(), err)
	}
}

func (s *Site) serveFile(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/index.html") {
		// We'll show index.html for the directory.
		// Use the dir/ version as canonical instead of dir/index.html.
		http.Redirect(w, r, r.URL.Path[0:len(r.URL.Path)-len("index.html")], http.StatusMovedPermanently)
		return
	}

	// Check to see if we need to redirect or serve another file.
	abspath := r.URL.Path
	relpath := path.Clean(strings.TrimPrefix(abspath, "/"))
	if f := open(s.fs, relpath); f != nil {
		if f.Path != abspath {
			// Redirect to canonical path.
			http.Redirect(w, r, f.Path, http.StatusMovedPermanently)
			return
		}
		// Serve from the actual filesystem path.
		s.serveHTML(w, r, f)
		return
	}

	dir, err := fs.Stat(s.fs, relpath)
	if err != nil {
		// Check for spurious trailing slash.
		if strings.HasSuffix(abspath, "/") {
			trimmed := relpath[:len(relpath)-1]
			if _, err := fs.Stat(s.fs, trimmed); err == nil ||
				open(s.fs, trimmed) != nil {
				http.Redirect(w, r, "/"+trimmed, http.StatusMovedPermanently)
				return
			}
		}
		s.ServeError(w, r, err)
		return
	}

	if dir != nil && dir.IsDir() {
		if maybeRedirect(w, r) {
			return
		}
		s.serveDir(w, r, relpath)
		return
	}

	if isTextFile(s.fs, relpath) {
		if maybeRedirectFile(w, r) {
			return
		}
		s.serveText(w, r, relpath)
		return
	}

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

var doctype = []byte("<!DOCTYPE ")

func (s *Site) serveHTML(w http.ResponseWriter, r *http.Request, f *file) {
	src := f.Body
	isMarkdown := strings.HasSuffix(f.FilePath, ".md")

	// if it begins with "<!DOCTYPE " assume it is standalone
	// html that doesn't need the template wrapping.
	if bytes.HasPrefix(src, doctype) {
		w.Write(src)
		return
	}

	page := Page{
		Title:    f.Title,
		Subtitle: f.Subtitle,
	}

	// evaluate as template if indicated
	if f.Template {
		page = s.fullPage(r, page)
		tmpl, err := template.New("main").Funcs(s.docFuncs).Parse(string(src))
		if err != nil {
			log.Printf("parsing template %s: %v", f.Path, err)
			s.ServeError(w, r, err)
			return
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, page); err != nil {
			log.Printf("executing template %s: %v", f.Path, err)
			s.ServeError(w, r, err)
			return
		}
		src = buf.Bytes()
	}

	// Apply markdown as indicated.
	// (Note template applies before Markdown.)
	if isMarkdown {
		html, err := renderMarkdown(src)
		if err != nil {
			log.Printf("executing markdown %s: %v", f.Path, err)
			s.ServeError(w, r, err)
			return
		}
		src = html
	}

	// if it's the language spec, add tags to EBNF productions
	if strings.HasSuffix(f.FilePath, "go_spec.html") {
		var buf bytes.Buffer
		spec.Linkify(&buf, src)
		src = buf.Bytes()
	}

	page.Data = template.HTML(src)
	s.ServePage(w, r, page)
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

	dirpath := strings.TrimSuffix(relpath, "/") + "/"
	s.ServePage(w, r, Page{
		Title:    "Directory",
		SrcPath:  dirpath,
		TabTitle: dirpath,
		Template: "dirlist.html",
		Data:     info,
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

	title := "Text file"
	if strings.HasSuffix(relpath, ".go") {
		title = "Source file"
	}
	s.ServePage(w, r, Page{
		Title:    title,
		SrcPath:  relpath,
		TabTitle: relpath,
		Data:     template.HTML(buf.String()),
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
