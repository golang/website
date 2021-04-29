// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package site

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/go.dev/cmd/internal/html/template"
	"golang.org/x/go.dev/cmd/internal/tmplfunc"
	"gopkg.in/yaml.v3"
)

func (site *Site) initTemplate() error {
	funcs := template.FuncMap{
		"absURL":      absURL,
		"default":     defaultFn,
		"dict":        dict,
		"fingerprint": fingerprint,
		"first":       first,
		"isset":       isset,
		"list":        list,
		"markdownify": markdownify,
		"path":        pathFn,
		"replace":     replace,
		"replaceRE":   replaceRE,
		"resources":   site.resources,
		"safeHTML":    safeHTML,
		"sort":        sortFn,
		"where":       where,
		"yaml":        yamlFn,
	}

	site.base = template.New("site").Funcs(funcs)
	if err := tmplfunc.ParseGlob(site.base, site.file("templates/*.tmpl")); err != nil && !strings.Contains(err.Error(), "pattern matches no files") {
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

func absURL(u string) string { return u }

func defaultFn(x, y string) string {
	if y != "" {
		return y
	}
	return x
}

type Fingerprint struct {
	r    *Resource
	Data struct {
		Integrity string
	}
	RelPermalink string
}

func fingerprint(r *Resource) *Fingerprint {
	f := &Fingerprint{r: r}
	sum := sha256.Sum256(r.data)
	ext := path.Ext(r.RelPermalink)
	f.RelPermalink = "/" + strings.TrimSuffix(r.RelPermalink, ext) + "." + hex.EncodeToString(sum[:]) + ext
	f.Data.Integrity = "sha256-" + base64.StdEncoding.EncodeToString(sum[:])
	return f
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

func isset(m map[string]interface{}, name string) bool {
	_, ok := m[name]
	return ok
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

// markdownify is the function provided to templates.
func markdownify(data interface{}) template.HTML {
	h := markdownToHTML(toString(data))
	s := strings.TrimSpace(string(h))
	if strings.HasPrefix(s, "<p>") && strings.HasSuffix(s, "</p>") && strings.Count(s, "<p>") == 1 {
		h = template.HTML(strings.TrimSpace(s[len("<p>") : len(s)-len("</p>")]))
	}
	return h
}

func pathFn() pathPkg { return pathPkg{} }

type pathPkg struct{}

func (pathPkg) Base(s interface{}) string { return path.Base(toString(s)) }
func (pathPkg) Dir(s interface{}) string  { return path.Dir(toString(s)) }
func (pathPkg) Join(args ...interface{}) string {
	var elem []string
	for _, a := range args {
		elem = append(elem, toString(a))
	}
	return path.Join(elem...)
}

func replace(input, x, y interface{}) string {
	return strings.ReplaceAll(toString(input), toString(x), toString(y))
}

func replaceRE(pattern, repl, input interface{}) string {
	re := regexp.MustCompile(toString(pattern))
	return re.ReplaceAllString(toString(input), toString(repl))
}

func safeHTML(s interface{}) template.HTML { return template.HTML(toString(s)) }

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
	return p.Site.pagesByID[p.section]
}

func (d *Page) HasMenuCurrent(x string, y *MenuItem) bool {
	return false
}

func (d *Page) IsMenuCurrent(x string, y *MenuItem) bool {
	return d.Permalink() == y.URL
}

func (p *Page) Param(key string) interface{} { return p.Params[key] }

func (p *Page) Parent() *Page {
	if p.IsHome {
		return nil
	}
	return p.Site.pagesByID[p.parent]
}

func (p *Page) Permalink() string {
	return strings.TrimRight(p.Site.BaseURL, "/") + p.RelPermalink()
}

func (p *Page) RelPermalink() string {
	if p.id == "" {
		return "/"
	}
	return "/" + p.id + "/"
}

func (p *Page) Resources() *PageResources {
	return &PageResources{p}
}

func (p *Page) Section() string {
	i := strings.Index(p.section, "/")
	if i < 0 {
		return p.section
	}
	return p.section[:i]
}

type PageResources struct{ p *Page }

func (r *PageResources) GetMatch(name string) (*Resource, error) {
	for _, rs := range r.p.TheResources {
		if name == rs.Name {
			if rs.data == nil {
				rs.RelPermalink = strings.TrimPrefix(filepath.ToSlash(filepath.Join(r.p.file, "../"+rs.Src)), "content")
				data, err := os.ReadFile(r.p.Site.file(r.p.file + "/../" + rs.Src))
				if err != nil {
					return nil, err
				}
				rs.data = data
			}
			return rs, nil
		}
	}
	return nil, nil
}

type Resource struct {
	data         []byte
	RelPermalink string
	Name         string
	Src          string
	Params       map[string]string
}

func (site *Site) resources() Resources { return Resources{site} }

type Resources struct{ site *Site }

func (r Resources) Get(name string) (*Resource, error) {
	data, err := os.ReadFile(r.site.file("assets/" + name))
	if err != nil {
		return nil, err
	}
	return &Resource{data: data, RelPermalink: name}, nil
}

func yamlFn(s string) (interface{}, error) {
	var d interface{}
	if err := yaml.Unmarshal([]byte(s), &d); err != nil {
		return nil, err
	}
	return d, nil
}
