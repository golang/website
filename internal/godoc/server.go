// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16
// +build go1.16

package godoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	htmlpkg "html"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/website/internal/pkgdoc"
	"golang.org/x/website/internal/spec"
	"golang.org/x/website/internal/texthtml"
)

// toFS returns the io/fs name for path (no leading slash).
func toFS(name string) string {
	if name == "/" {
		return "."
	}
	return path.Clean(strings.TrimPrefix(name, "/"))
}

// docServer serves a package doc tree (/cmd or /pkg).
type docServer struct {
	p *Presentation
	d *pkgdoc.Docs
}

func (h *docServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if redirect(w, r) {
		return
	}

	// TODO(rsc): URL should be clean already.
	relpath := path.Clean(strings.TrimPrefix(r.URL.Path, "/pkg/"))

	abspath := path.Join("/src", relpath)
	mode := pkgdoc.ParseMode(r.FormValue("m"))
	if relpath == "builtin" {
		// The fake built-in package contains unexported identifiers,
		// but we want to show them. Also, disable type association,
		// since it's not helpful for this fake package (see issue 6645).
		mode |= pkgdoc.ModeAll | pkgdoc.ModeBuiltin
	}
	info := pkgdoc.Doc(h.d, abspath, relpath, mode, r.FormValue("GOOS"), r.FormValue("GOARCH"))
	if info.Err != nil {
		log.Print(info.Err)
		h.p.ServeError(w, r, relpath, info.Err)
		return
	}

	var tabtitle, title, subtitle string
	switch {
	case info.PAst != nil:
		for _, ast := range info.PAst {
			tabtitle = ast.Name.Name
			break
		}
	case info.PDoc != nil:
		tabtitle = info.PDoc.Name
	default:
		tabtitle = info.Dirname
		title = "Directory "
	}
	if title == "" {
		if info.IsMain {
			// assume that the directory name is the command name
			_, tabtitle = path.Split(relpath)
			title = "Command "
		} else {
			title = "Package "
		}
	}
	title += tabtitle

	// special cases for top-level package/command directories
	switch tabtitle {
	case "/src":
		title = "Packages"
		tabtitle = "Packages"
	case "/src/cmd":
		title = "Commands"
		tabtitle = "Commands"
	}

	info.GoogleCN = h.p.googleCN(r)
	var body []byte
	if info.Dirname == "/src" {
		body = applyTemplate(h.p.PackageRootHTML, "packageRootHTML", info)
	} else {
		body = applyTemplate(h.p.PackageHTML, "packageHTML", info)
	}
	h.p.ServePage(w, Page{
		Title:    title,
		Tabtitle: tabtitle,
		Subtitle: subtitle,
		Body:     body,
		GoogleCN: info.GoogleCN,
	})
}

func modeQueryString(m pkgdoc.Mode) string {
	s := m.String()
	if s == "" {
		return ""
	}
	return "?m=" + s
}

func applyTemplate(t *template.Template, name string, data interface{}) []byte {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		log.Printf("%s.Execute: %s", name, err)
	}
	return buf.Bytes()
}

type writerCapturesErr struct {
	w   io.Writer
	err error
}

func (w *writerCapturesErr) Write(p []byte) (int, error) {
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
	w := &writerCapturesErr{w: rw}
	err := t.Execute(w, data)
	// There are some cases where template.Execute does not return an error when
	// rw returns an error, and some where it does.  So check w.err first.
	if w.err == nil && err != nil {
		// Log template errors.
		log.Printf("%s.Execute: %s", t.Name(), err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) (redirected bool) {
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

func redirectFile(w http.ResponseWriter, r *http.Request) (redirected bool) {
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

func (p *Presentation) serveTextFile(w http.ResponseWriter, r *http.Request, abspath, relpath string) {
	src, err := fs.ReadFile(p.fs, toFS(abspath))
	if err != nil {
		log.Printf("ReadFile: %s", err)
		p.ServeError(w, r, relpath, err)
		return
	}

	if r.FormValue("m") == "text" {
		p.ServeText(w, src)
		return
	}

	cfg := texthtml.Config{
		GoComments: path.Ext(abspath) == ".go",
		Highlight:  r.FormValue("h"),
		Selection:  rangeSelection(r.FormValue("s")),
		Line:       1,
	}

	var buf bytes.Buffer
	buf.WriteString("<pre>")
	buf.Write(texthtml.Format(src, cfg))
	buf.WriteString("</pre>")

	fmt.Fprintf(&buf, `<p><a href="/%s?m=text">View as plain text</a></p>`, htmlpkg.EscapeString(relpath))

	title := "Text file"
	if strings.HasSuffix(relpath, ".go") {
		title = "Source file"
	}
	p.ServePage(w, Page{
		Title:    title,
		SrcPath:  relpath,
		Tabtitle: relpath,
		Body:     buf.Bytes(),
		GoogleCN: p.googleCN(r),
	})
}

func (p *Presentation) serveDirectory(w http.ResponseWriter, r *http.Request, abspath, relpath string) {
	if redirect(w, r) {
		return
	}

	list, err := fs.ReadDir(p.fs, toFS(abspath))
	if err != nil {
		p.ServeError(w, r, relpath, err)
		return
	}

	var info []fs.FileInfo
	for _, d := range list {
		i, err := d.Info()
		if err == nil {
			info = append(info, i)
		}
	}

	p.ServePage(w, Page{
		Title:    "Directory",
		SrcPath:  relpath,
		Tabtitle: relpath,
		Body:     applyTemplate(p.DirlistHTML, "dirlistHTML", info),
		GoogleCN: p.googleCN(r),
	})
}

func (p *Presentation) serveHTML(w http.ResponseWriter, r *http.Request, f *file) {
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
		GoogleCN: p.googleCN(r),
	}

	// evaluate as template if indicated
	if f.Template {
		tmpl, err := template.New("main").Funcs(p.DocFuncs).Parse(string(src))
		if err != nil {
			log.Printf("parsing template %s: %v", f.Path, err)
			p.ServeError(w, r, f.Path, err)
			return
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, page); err != nil {
			log.Printf("executing template %s: %v", f.Path, err)
			p.ServeError(w, r, f.Path, err)
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
			p.ServeError(w, r, f.Path, err)
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

	page.Body = src
	p.ServePage(w, page)
}

func (p *Presentation) ServeFile(w http.ResponseWriter, r *http.Request) {
	p.serveFile(w, r)
}

func (p *Presentation) serveFile(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/index.html") {
		// We'll show index.html for the directory.
		// Use the dir/ version as canonical instead of dir/index.html.
		http.Redirect(w, r, r.URL.Path[0:len(r.URL.Path)-len("index.html")], http.StatusMovedPermanently)
		return
	}

	// Check to see if we need to redirect or serve another file.
	abspath := r.URL.Path
	if f := open(p.fs, abspath); f != nil {
		if f.Path != abspath {
			// Redirect to canonical path.
			http.Redirect(w, r, f.Path, http.StatusMovedPermanently)
			return
		}
		// Serve from the actual filesystem path.
		p.serveHTML(w, r, f)
		return
	}

	relpath := abspath[1:] // strip leading slash

	dir, err := fs.Stat(p.fs, toFS(abspath))
	if err != nil {
		// Check for spurious trailing slash.
		if strings.HasSuffix(abspath, "/") {
			trimmed := abspath[:len(abspath)-1]
			if _, err := fs.Stat(p.fs, toFS(trimmed)); err == nil ||
				open(p.fs, trimmed) != nil {
				http.Redirect(w, r, trimmed, http.StatusMovedPermanently)
				return
			}
		}
		p.ServeError(w, r, relpath, err)
		return
	}

	fsPath := toFS(abspath)
	if dir != nil && dir.IsDir() {
		if redirect(w, r) {
			return
		}
		p.serveDirectory(w, r, abspath, relpath)
		return
	}

	if isTextFile(p.fs, fsPath) {
		if redirectFile(w, r) {
			return
		}
		p.serveTextFile(w, r, abspath, relpath)
		return
	}

	p.fileServer.ServeHTTP(w, r)
}

func (p *Presentation) ServeText(w http.ResponseWriter, text []byte) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(text)
}

func marshalJSON(x interface{}) []byte {
	var data []byte
	var err error
	const indentJSON = false // for easier debugging
	if indentJSON {
		data, err = json.MarshalIndent(x, "", "    ")
	} else {
		data, err = json.Marshal(x)
	}
	if err != nil {
		panic(fmt.Sprintf("json.Marshal failed: %s", err))
	}
	return data
}
