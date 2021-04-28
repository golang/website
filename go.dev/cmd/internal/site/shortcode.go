// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package site

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/go.dev/cmd/internal/html/template"
)

// A ShortCode is a parsed Hugo shortcode like {{% foo %}} or {{< foo >}}.
// It is Hugo's wrapping of a template call and will be replaced by actual template calls.
type ShortCode struct {
	Kind  string
	Name  string
	Args  []string
	Keys  map[string]string
	Inner template.HTML
	Page  *Page
}

func (c *ShortCode) run() (template.HTML, error) {
	return c.Page.Site.runTemplate("layouts/shortcodes/"+c.Name+".html", c)
}

func (c *ShortCode) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("<code %s %v %v %q>", c.Name, c.Args, c.Keys, c.Inner)
}

func (c *ShortCode) Get(x interface{}) string {
	switch x := x.(type) {
	case int:
		if 0 <= x && x < len(c.Args) {
			return c.Args[x]
		}
		return ""
	case string:
		return c.Keys[x]
	}
	panic(fmt.Sprintf("bad Get %v", x))
}

func (c *ShortCode) IsNamedParams() bool {
	return len(c.Keys) != 0
}

// parseCodes parses the shortcode invocations in the Hugo markdown file.
// It returns a slice containing only two types of elements: string and *ShortCode.
// The even indexes are strings and the odd indexes are *ShortCode.
// There are an odd number of elements (the slice begins and ends with a string).
func (p *Page) parseCodes(markdown string) []interface{} {
	t1, c1, t2, kind1 := findCode(markdown)
	t2, c2, t3, kind2 := findCode(t2)
	var ret []interface{}
	for c1 != nil {
		c1.Kind = kind1
		c1.Page = p
		if c2 != nil && c2.Name == "/"+c1.Name {
			c1.Inner = template.HTML(t2)
			t2, c2, t3, kind2 = findCode(t3)
			continue
		}
		ret = append(ret, t1)
		ret = append(ret, c1)
		t1, c1, kind1 = t2, c2, kind2
		t2, c2, t3, kind2 = findCode(t3)
	}
	ret = append(ret, t1)
	return ret
}

func findCode(text string) (before string, code *ShortCode, after string, kind string) {
	end := "%}}"
	kind = "%"
	i := strings.Index(text, "{{%")
	j := strings.Index(text, "{{<")
	if i < 0 || j >= 0 && j < i {
		i = j
		kind = "<"
		end = ">}}"
	}
	if i < 0 {
		return text, nil, "", ""
	}
	j = strings.Index(text[i+3:], end)
	if j < 0 {
		return text, nil, "", ""
	}
	before, codeText, after := text[:i], text[i+3:i+3+j], text[i+3+j+3:]
	codeText = strings.TrimSpace(codeText)
	name, args, _ := cutAny(codeText, " \t\r\n")
	if name == "" {
		log.Fatalf("empty code")
	}
	args = strings.TrimSpace(args)
	code = &ShortCode{Name: name, Keys: make(map[string]string)}
	for args != "" {
		k, v := "", args
		if strings.HasPrefix(args, `"`) {
			goto Value
		}
		{
			i := strings.Index(args, "=")
			if i < 0 {
				goto Value
			}
			for j := 0; j < i; j++ {
				if args[j] == ' ' || args[j] == '\t' {
					goto Value
				}
			}
			k, v = args[:i], args[i+1:]
		}
	Value:
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, `"`) {
			j := 1
			for ; ; j++ {
				if j >= len(v) {
					log.Fatalf("unterminated quoted string: %s", args)
				}
				if v[j] == '"' {
					v, args = v[:j+1], v[j+1:]
					break
				}
				if v[j] == '\\' {
					j++
				}
			}
			u, err := strconv.Unquote(v)
			if err != nil {
				log.Fatalf("malformed k=v: %s=%s", k, v)
			}
			v = u
		} else {
			v, args, _ = cutAny(v, " \t\r\n")
		}
		if k == "" {
			code.Args = append(code.Args, v)
		} else {
			code.Keys[k] = v
		}
		args = strings.TrimSpace(args)
	}
	return
}

func cutAny(s, any string) (before, after string, ok bool) {
	if i := strings.IndexAny(s, any); i >= 0 {
		_, size := utf8.DecodeRuneInString(s[i:])
		return s[:i], s[i+size:], true
	}
	return s, "", false
}
