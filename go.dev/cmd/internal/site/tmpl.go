// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package site

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/go.dev/cmd/internal/html/template"
	"golang.org/x/go.dev/cmd/internal/tmplfunc"
	"gopkg.in/yaml.v3"
)

func (site *Site) initTemplate() error {
	funcs := template.FuncMap{
		"add":      func(i, j int) int { return i + j },
		"data":     site.data,
		"dict":     dict,
		"first":    first,
		"markdown": markdown,
		"newest":   newest,
		"page":     site.pageByPath,
		"pages":    site.pagesGlob,
		"replace":  replace,
		"rawhtml":  rawhtml,
		"trim":     strings.Trim,
		"yaml":     yamlFn,
	}

	site.base = template.New("site").Funcs(funcs)
	if err := tmplfunc.ParseGlob(site.base, site.file("_templates/*.tmpl")); err != nil && !strings.Contains(err.Error(), "pattern matches no files") {
		return err
	}
	return nil
}

func (site *Site) clone() *template.Template {
	t := template.Must(site.base.Clone())
	if err := tmplfunc.Funcs(t); err != nil {
		panic(err)
	}
	return t
}

func toString(x interface{}) string {
	switch x := x.(type) {
	case string:
		return x
	case template.HTML:
		return string(x)
	case nil:
		return ""
	default:
		panic(fmt.Sprintf("cannot toString %T", x))
	}
}

func first(n int, list reflect.Value) reflect.Value {
	if !list.IsValid() {
		return list
	}
	if list.Kind() == reflect.Interface {
		if list.IsNil() {
			return list
		}
		list = list.Elem()
	}

	if list.Len() < n {
		return list
	}
	return list.Slice(0, n)
}

func dict(args ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		m[args[i].(string)] = args[i+1]
	}
	m["Identifier"] = "IDENT"
	return m
}

func list(args ...interface{}) []interface{} {
	return args
}

// markdown is the function provided to templates.
func markdown(data interface{}) (template.HTML, error) {
	h, err := markdownToHTML(toString(data))
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(string(h))
	if strings.HasPrefix(s, "<p>") && strings.HasSuffix(s, "</p>") && strings.Count(s, "<p>") == 1 {
		h = template.HTML(strings.TrimSpace(s[len("<p>") : len(s)-len("</p>")]))
	}
	return h, nil
}

func replace(input, x, y interface{}) string {
	return strings.ReplaceAll(toString(input), toString(x), toString(y))
}

func rawhtml(s interface{}) template.HTML {
	return template.HTML(toString(s))
}

func yamlFn(s string) (interface{}, error) {
	var d interface{}
	if err := yaml.Unmarshal([]byte(s), &d); err != nil {
		return nil, err
	}
	return d, nil
}
