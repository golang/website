// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package site

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/go.dev/cmd/internal/html/template"
	"golang.org/x/go.dev/cmd/internal/tmplfunc"
	"gopkg.in/yaml.v3"
)

func (site *Site) initTemplate() error {
	funcs := template.FuncMap{
		"data":     site.data,
		"dict":     dict,
		"first":    first,
		"list":     list,
		"markdown": markdown,
		"replace":  replace,
		"rawhtml":  rawhtml,
		"sort":     sortFn,
		"where":    where,
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
	out := reflect.Zero(list.Type())

	for i := 0; i < list.Len() && i < n; i++ {
		out = reflect.Append(out, list.Index(i))
	}
	return out
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
func markdown(data interface{}) template.HTML {
	h := markdownToHTML(toString(data))
	s := strings.TrimSpace(string(h))
	if strings.HasPrefix(s, "<p>") && strings.HasSuffix(s, "</p>") && strings.Count(s, "<p>") == 1 {
		h = template.HTML(strings.TrimSpace(s[len("<p>") : len(s)-len("</p>")]))
	}
	return h
}

func replace(input, x, y interface{}) string {
	return strings.ReplaceAll(toString(input), toString(x), toString(y))
}

func rawhtml(s interface{}) template.HTML {
	return template.HTML(toString(s))
}

func sortFn(list reflect.Value, key, asc string) (reflect.Value, error) {
	out := reflect.Zero(list.Type())
	var keys []string
	var perm []int
	for i := 0; i < list.Len(); i++ {
		elem := list.Index(i)
		v, ok := eval(elem, key)
		if !ok {
			return reflect.Value{}, fmt.Errorf("no key %s", key)
		}
		keys = append(keys, strings.ToLower(v))
		perm = append(perm, i)
	}
	sort.Slice(perm, func(i, j int) bool {
		return keys[perm[i]] < keys[perm[j]]
	})
	for _, i := range perm {
		out = reflect.Append(out, list.Index(i))
	}
	return out, nil
}

func where(list reflect.Value, key, val string) reflect.Value {
	out := reflect.Zero(list.Type())
	for i := 0; i < list.Len(); i++ {
		elem := list.Index(i)
		v, ok := eval(elem, key)
		if ok && v == val {
			out = reflect.Append(out, elem)
		}
	}
	return out
}

func eval(elem reflect.Value, key string) (string, bool) {
	for _, k := range strings.Split(key, ".") {
		if !elem.IsValid() {
			return "", false
		}
		m := elem.MethodByName(k)
		if m.IsValid() {
			elem = m.Call(nil)[0]
			continue
		}
		if elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr {
			if elem.IsNil() {
				return "", false
			}
			elem = elem.Elem()
		}
		switch elem.Kind() {
		case reflect.Struct:
			elem = elem.FieldByName(k)
			continue
		case reflect.Map:
			elem = elem.MapIndex(reflect.ValueOf(k))
			continue
		}
		return "", false
	}
	if !elem.IsValid() {
		return "", false
	}
	if elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr {
		if elem.IsNil() {
			return "", false
		}
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.String {
		return "", false
	}
	return elem.String(), true
}

func (p *Page) CurrentSection() *Page {
	return p.site.pagesByID[p.Section()]
}

func (p *Page) IsHome() bool { return p.id == "" }

func (p *Page) Parent() *Page {
	if p.IsHome() {
		return nil
	}
	return p.site.pagesByID[p.parent]
}

func (p *Page) URL() string {
	return strings.TrimRight(p.site.URL, "/") + p.Path()
}

func (p *Page) Path() string {
	if p.id == "" {
		return "/"
	}
	if strings.HasSuffix(p.file, "/index.md") {
		return "/" + p.id + "/"
	}
	return "/" + p.id
}

func (p *Page) Section() string {
	i := strings.Index(p.section, "/")
	if i < 0 {
		return p.section
	}
	return p.section[:i]
}

func yamlFn(s string) (interface{}, error) {
	var d interface{}
	if err := yaml.Unmarshal([]byte(s), &d); err != nil {
		return nil, err
	}
	return d, nil
}

func (p *Page) Dir() string {
	return strings.TrimPrefix(filepath.ToSlash(filepath.Dir(p.file)), "_content")
}
