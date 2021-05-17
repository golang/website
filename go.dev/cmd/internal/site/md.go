// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package site

import (
	"bytes"
	"strings"

	"github.com/russross/blackfriday"
	"golang.org/x/go.dev/cmd/internal/html/template"
	"golang.org/x/go.dev/cmd/internal/tmplfunc"
)

// markdownToHTML converts markdown to HTML using the renderer and settings that Hugo uses.
func markdownToHTML(markdown string) template.HTML {
	markdown = strings.TrimLeft(markdown, "\n")
	renderer := blackfriday.HtmlRenderer(blackfriday.HTML_USE_XHTML|
		blackfriday.HTML_USE_SMARTYPANTS|
		blackfriday.HTML_SMARTYPANTS_FRACTIONS|
		blackfriday.HTML_SMARTYPANTS_DASHES|
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES|
		blackfriday.HTML_NOREFERRER_LINKS|
		blackfriday.HTML_HREF_TARGET_BLANK,
		"", "")
	options := blackfriday.Options{
		Extensions: blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
			blackfriday.EXTENSION_TABLES |
			blackfriday.EXTENSION_FENCED_CODE |
			blackfriday.EXTENSION_AUTOLINK |
			blackfriday.EXTENSION_STRIKETHROUGH |
			blackfriday.EXTENSION_SPACE_HEADERS |
			blackfriday.EXTENSION_HEADER_IDS |
			blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
			blackfriday.EXTENSION_DEFINITION_LISTS |
			blackfriday.EXTENSION_AUTO_HEADER_IDS,
	}
	return template.HTML(blackfriday.MarkdownOptions([]byte(markdown), renderer, options))
}

// markdownTemplateToHTML converts a markdown template to HTML,
// first applying the template execution engine and then interpreting
// the result as markdown to be converted to HTML.
// This is the same logic used by the Go web site.
func (site *Site) markdownTemplateToHTML(markdown string, p *page) (template.HTML, error) {
	t := site.clone().New(p.file)
	if err := tmplfunc.Parse(t, string(p.data)); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, p.params); err != nil {
		return "", err
	}
	return markdownToHTML(buf.String()), nil
}
