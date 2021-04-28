// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package site

import (
	"fmt"
	"strings"

	"github.com/russross/blackfriday"
	"golang.org/x/go.dev/cmd/internal/html/template"
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

// markdownWithShortCodesToHTML converts markdown to HTML,
// first expanding Hugo shortcodes in the markdown input.
// Shortcodes templates are given access to p as .Page.
func markdownWithShortCodesToHTML(markdown string, p *Page) (template.HTML, error) {
	// We replace each shortcode invocation in the markdown with
	// a keyword HUGOREPLACECODE0001 etc and then run the result
	// through markdown conversion, and then we substitute the actual
	// shortcode ouptuts for the keywords.
	var md string         // current markdown chunk
	var replaces []string // replacements to apply to all at end

	for i, elem := range p.parseCodes(markdown) {
		switch elem := elem.(type) {
		default:
			return "", fmt.Errorf("unexpected elem %T", elem)
		case string:
			md += elem

		case *ShortCode:
			code := elem
			html, err := code.run()
			if err != nil {
				return "", err
			}
			// Adjust shortcode output to match Hugo's line breaks.
			// This is weird but will go away when we retire shortcodes.
			if code.Inner != "" {
				html = "\n\n" + html
			} else if code.Kind == "%" {
				html = template.HTML(strings.TrimLeft(string(html), " \n"))
			}
			key := fmt.Sprintf("HUGOREPLACECODE%04d", i)
			md += key
			replaces = append(replaces, key, string(html), "<p>"+key+"</p>", string(html))
		}
	}
	html := markdownToHTML(md)
	return template.HTML(strings.NewReplacer(replaces...).Replace(string(html))), nil
}
