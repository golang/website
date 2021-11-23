// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tour

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/tools/present"
	"golang.org/x/website"
)

var (
	uiContent      []byte
	lessons        = make(map[string][]byte)
	lessonNotFound = fmt.Errorf("lesson not found")
)

var contentTour = website.TourOnly()

// initTour loads tour.article and the relevant HTML templates from root.
func initTour(mux *http.ServeMux, transport string) error {
	// Make sure playground is enabled before rendering.
	present.PlayEnabled = true

	// Set up templates.
	tmpl, err := present.Template().ParseFS(contentTour, "tour/template/action.tmpl")
	if err != nil {
		return fmt.Errorf("parse templates: %v", err)
	}

	// Init lessons.
	if err := initLessons(tmpl); err != nil {
		return fmt.Errorf("init lessons: %v", err)
	}

	// Init UI.
	ui, err := template.ParseFS(contentTour, "tour/template/index.tmpl")
	if err != nil {
		return fmt.Errorf("parse index.tmpl: %v", err)
	}
	buf := new(bytes.Buffer)

	data := struct {
		AnalyticsHTML template.HTML
	}{analyticsHTML}

	if err := ui.Execute(buf, data); err != nil {
		return fmt.Errorf("render UI: %v", err)
	}
	uiContent = buf.Bytes()

	mux.HandleFunc("/tour/", rootHandler)
	mux.HandleFunc("/tour/lesson/", lessonHandler)
	mux.Handle("/tour/static/", http.FileServer(http.FS(contentTour)))

	return initScript(mux, socketAddr(), transport)
}

// initLessonss finds all the lessons in the content directory, renders them,
// using the given template and saves the content in the lessons map.
func initLessons(tmpl *template.Template) error {
	files, err := fs.ReadDir(contentTour, "tour")
	if err != nil {
		return err
	}
	for _, f := range files {
		if path.Ext(f.Name()) != ".article" {
			continue
		}
		content, err := parseLesson(f.Name(), tmpl)
		if err != nil {
			return fmt.Errorf("parsing %v: %v", f.Name(), err)
		}
		name := strings.TrimSuffix(f.Name(), ".article")
		lessons[name] = content
	}
	return nil
}

// file defines the JSON form of a code file in a page.
type file struct {
	Name    string
	Content string
	Hash    string
}

// page defines the JSON form of a tour lesson page.
type page struct {
	Title   string
	Content string
	Files   []file
}

// lesson defines the JSON form of a tour lesson.
type lesson struct {
	Title       string
	Description string
	Pages       []page
}

// parseLesson parses and returns a lesson content given its path
// relative to root ('/'-separated) and the template to render it.
func parseLesson(path string, tmpl *template.Template) ([]byte, error) {
	f, err := contentTour.Open("tour/" + path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ctx := &present.Context{
		ReadFile: func(filename string) ([]byte, error) {
			return fs.ReadFile(contentTour, "tour/"+filepath.ToSlash(filename))
		},
	}
	doc, err := ctx.Parse(prepContent(f), path, 0)
	if err != nil {
		return nil, err
	}

	lesson := lesson{
		doc.Title,
		doc.Subtitle,
		make([]page, len(doc.Sections)),
	}

	for i, sec := range doc.Sections {
		p := &lesson.Pages[i]
		w := new(bytes.Buffer)
		if err := sec.Render(w, tmpl); err != nil {
			return nil, fmt.Errorf("render section: %v", err)
		}
		p.Title = sec.Title
		p.Content = w.String()
		codes := findPlayCode(sec)
		p.Files = make([]file, len(codes))
		for i, c := range codes {
			f := &p.Files[i]
			f.Name = c.FileName
			f.Content = string(c.Raw)
			hash := sha1.Sum(c.Raw)
			f.Hash = base64.StdEncoding.EncodeToString(hash[:])
		}
	}

	w := new(bytes.Buffer)
	if err := json.NewEncoder(w).Encode(lesson); err != nil {
		return nil, fmt.Errorf("encode lesson: %v", err)
	}
	return w.Bytes(), nil
}

// findPlayCode returns a slide with all the Code elements in the given
// Elem with Play set to true.
func findPlayCode(e present.Elem) []*present.Code {
	var r []*present.Code
	switch v := e.(type) {
	case present.Code:
		if v.Play {
			r = append(r, &v)
		}
	case present.Section:
		for _, s := range v.Elem {
			r = append(r, findPlayCode(s)...)
		}
	}
	return r
}

// writeLesson writes the tour content to the provided Writer.
func writeLesson(name string, w io.Writer) error {
	if uiContent == nil {
		panic("writeLesson called before successful initTour")
	}
	if len(name) == 0 {
		return writeAllLessons(w)
	}
	l, ok := lessons[name]
	if !ok {
		return lessonNotFound
	}
	_, err := w.Write(l)
	return err
}

func writeAllLessons(w io.Writer) error {
	if _, err := fmt.Fprint(w, "{"); err != nil {
		return err
	}
	nLessons := len(lessons)
	for k, v := range lessons {
		if _, err := fmt.Fprintf(w, "%q:%s", k, v); err != nil {
			return err
		}
		nLessons--
		if nLessons != 0 {
			if _, err := fmt.Fprint(w, ","); err != nil {
				return err
			}
		}
	}
	_, err := fmt.Fprint(w, "}")
	return err
}

// renderUI writes the tour UI to the provided Writer.
func renderUI(w io.Writer) error {
	if uiContent == nil {
		panic("renderUI called before successful initTour")
	}
	_, err := w.Write(uiContent)
	return err
}

// initScript concatenates all the javascript files needed to render
// the tour UI and serves the result on /script.js.
func initScript(mux *http.ServeMux, socketAddr, transport string) error {
	modTime := time.Now()
	b := new(bytes.Buffer)

	// Keep this list in dependency order
	files := []string{
		"../js/playground.js",
		"static/lib/jquery.min.js",
		"static/lib/jquery-ui.min.js",
		"static/lib/angular.min.js",
		"static/lib/codemirror/lib/codemirror.js",
		"static/lib/codemirror/mode/go/go.js",
		"static/lib/angular-ui.min.js",
		"static/js/app.js",
		"static/js/controllers.js",
		"static/js/directives.js",
		"static/js/services.js",
		"static/js/values.js",
	}

	for _, file := range files {
		f, err := fs.ReadFile(contentTour, path.Clean("tour/"+file))
		if err != nil {
			return err
		}
		b.Write(f)
	}

	f, err := fs.ReadFile(contentTour, "tour/static/js/page.js")
	if err != nil {
		return err
	}
	s := string(f)
	s = strings.ReplaceAll(s, "{{.SocketAddr}}", socketAddr)
	s = strings.ReplaceAll(s, "{{.Transport}}", transport)
	b.WriteString(s)

	mux.HandleFunc("/tour/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/javascript")
		// Set expiration time in one week.
		w.Header().Set("Cache-control", "max-age=604800")
		http.ServeContent(w, r, "", modTime, bytes.NewReader(b.Bytes()))
	})

	return nil
}
