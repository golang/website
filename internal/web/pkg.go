// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"path"
	"strings"
	"unicode"
)

type pkgPath struct{}

func (pkgPath) Base(a string) string            { return path.Base(a) }
func (pkgPath) Clean(a string) string           { return path.Clean(a) }
func (pkgPath) Dir(a string) string             { return path.Dir(a) }
func (pkgPath) Ext(a string) string             { return path.Ext(a) }
func (pkgPath) IsAbs(a string) bool             { return path.IsAbs(a) }
func (pkgPath) Join(a ...string) string         { return path.Join(a...) }
func (pkgPath) Match(a, b string) (bool, error) { return path.Match(a, b) }
func (pkgPath) Split(a string) (string, string) { return path.Split(a) }

type pkgStrings struct{}

func (pkgStrings) Compare(a, b string) int                         { return strings.Compare(a, b) }
func (pkgStrings) Contains(a, b string) bool                       { return strings.Contains(a, b) }
func (pkgStrings) ContainsAny(a, b string) bool                    { return strings.ContainsAny(a, b) }
func (pkgStrings) ContainsRune(a string, b rune) bool              { return strings.ContainsRune(a, b) }
func (pkgStrings) Count(a, b string) int                           { return strings.Count(a, b) }
func (pkgStrings) EqualFold(a, b string) bool                      { return strings.EqualFold(a, b) }
func (pkgStrings) Fields(a string) []string                        { return strings.Fields(a) }
func (pkgStrings) FieldsFunc(a string, b func(rune) bool) []string { return strings.FieldsFunc(a, b) }
func (pkgStrings) HasPrefix(a, b string) bool                      { return strings.HasPrefix(a, b) }
func (pkgStrings) HasSuffix(a, b string) bool                      { return strings.HasSuffix(a, b) }
func (pkgStrings) Index(a, b string) int                           { return strings.Index(a, b) }
func (pkgStrings) IndexAny(a, b string) int                        { return strings.IndexAny(a, b) }
func (pkgStrings) IndexByte(a string, b byte) int                  { return strings.IndexByte(a, b) }
func (pkgStrings) IndexFunc(a string, b func(rune) bool) int       { return strings.IndexFunc(a, b) }
func (pkgStrings) IndexRune(a string, b rune) int                  { return strings.IndexRune(a, b) }
func (pkgStrings) Join(a []string, b string) string                { return strings.Join(a, b) }
func (pkgStrings) LastIndex(a, b string) int                       { return strings.LastIndex(a, b) }
func (pkgStrings) LastIndexAny(a, b string) int                    { return strings.LastIndexAny(a, b) }
func (pkgStrings) LastIndexByte(a string, b byte) int              { return strings.LastIndexByte(a, b) }
func (pkgStrings) LastIndexFunc(a string, b func(rune) bool) int {
	return strings.LastIndexFunc(a, b)
}
func (pkgStrings) Map(a func(rune) rune, b string) string    { return strings.Map(a, b) }
func (pkgStrings) NewReader(a string) *strings.Reader        { return strings.NewReader(a) }
func (pkgStrings) NewReplacer(a ...string) *strings.Replacer { return strings.NewReplacer(a...) }
func (pkgStrings) Repeat(a string, b int) string             { return strings.Repeat(a, b) }
func (pkgStrings) Replace(a, b, c string, d int) string      { return strings.Replace(a, b, c, d) }
func (pkgStrings) ReplaceAll(a, b, c string) string          { return strings.ReplaceAll(a, b, c) }
func (pkgStrings) Split(a, b string) []string                { return strings.Split(a, b) }
func (pkgStrings) SplitAfter(a, b string) []string           { return strings.SplitAfter(a, b) }
func (pkgStrings) SplitAfterN(a, b string, c int) []string   { return strings.SplitAfterN(a, b, c) }
func (pkgStrings) SplitN(a, b string, c int) []string        { return strings.SplitN(a, b, c) }
func (pkgStrings) Title(a string) string                     { return strings.Title(a) }
func (pkgStrings) ToLower(a string) string                   { return strings.ToLower(a) }
func (pkgStrings) ToLowerSpecial(a unicode.SpecialCase, b string) string {
	return strings.ToLowerSpecial(a, b)
}
func (pkgStrings) ToTitle(a string) string { return strings.ToTitle(a) }
func (pkgStrings) ToTitleSpecial(a unicode.SpecialCase, b string) string {
	return strings.ToTitleSpecial(a, b)
}
func (pkgStrings) ToUpper(a string) string { return strings.ToUpper(a) }
func (pkgStrings) ToUpperSpecial(a unicode.SpecialCase, b string) string {
	return strings.ToUpperSpecial(a, b)
}
func (pkgStrings) ToValidUTF8(a, b string) string                  { return strings.ToValidUTF8(a, b) }
func (pkgStrings) Trim(a, b string) string                         { return strings.Trim(a, b) }
func (pkgStrings) TrimFunc(a string, b func(rune) bool) string     { return strings.TrimFunc(a, b) }
func (pkgStrings) TrimLeft(a, b string) string                     { return strings.TrimLeft(a, b) }
func (pkgStrings) TrimLeftFunc(a string, b func(rune) bool) string { return strings.TrimLeftFunc(a, b) }
func (pkgStrings) TrimPrefix(a, b string) string                   { return strings.TrimPrefix(a, b) }
func (pkgStrings) TrimRight(a, b string) string                    { return strings.TrimRight(a, b) }
func (pkgStrings) TrimRightFunc(a string, b func(rune) bool) string {
	return strings.TrimRightFunc(a, b)
}
func (pkgStrings) TrimSpace(a string) string     { return strings.TrimSpace(a) }
func (pkgStrings) TrimSuffix(a, b string) string { return strings.TrimSuffix(a, b) }
