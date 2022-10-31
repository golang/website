// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package history holds the Go project release history.
package history

import (
	"fmt"
	"html"
	"sort"
	"strings"
	"time"

	"golang.org/x/website/internal/backport/html/template"
)

// A Release describes a single Go release.
type Release struct {
	Version Version
	Date    Date
	Future  bool // if true, the release hasn't happened yet

	// Release content summary.
	Security      *FixSummary   // Security fixes, if any.
	Bug           *FixSummary   // Bug fixes, if any.
	More          template.HTML // Additional release content.
	CustomSummary template.HTML // CustomSummary, if non-empty, replaces the entire release content summary with custom HTML.
}

// FixSummary summarizes fixes in a Go release, listing components and packages involved.
type FixSummary struct {
	Quantifier string          // Optional quantifier. Empty string for unspecified amount of fixes (typical), "a" for a single fix, "two", "three" for multiple fixes, etc.
	Components []template.HTML // Components involved. For example, "cgo", "the compiler", "runtime", "the <code>go</code> command", etc.
	Packages   []string        // Packages involved. For example, "crypto/x509", "net/http", etc.
}

// A Version is a Go release version.
//
// In contrast to Semantic Versioning 2.0.0,
// trailing zero components are omitted,
// a version like Go 1.14 is considered a major Go release,
// a version like Go 1.14.1 is considered a minor Go release.
//
// See proposal golang.org/issue/32450 for background, details,
// and a discussion of the costs involved in making a change.
type Version struct {
	X int // X is the 1st component of a Go X.Y.Z version. It must be 1 or higher.
	Y int // Y is the 2nd component of a Go X.Y.Z version. It must be 0 or higher.
	Z int // Z is the 3rd component of a Go X.Y.Z version. It must be 0 or higher.
}

// String returns the Go release version string,
// like "1.14", "1.14.1", "1.14.2", and so on.
func (v Version) String() string {
	switch {
	case v.Z != 0:
		return fmt.Sprintf("%d.%d.%d", v.X, v.Y, v.Z)
	case v.Y != 0:
		return fmt.Sprintf("%d.%d", v.X, v.Y)
	default:
		return fmt.Sprintf("%d", v.X)
	}
}

// Before reports whether version v comes before version u.
func (v Version) Before(u Version) bool {
	if v.X != u.X {
		return v.X < u.X
	}
	if v.Y != u.Y {
		return v.Y < u.Y
	}
	return v.Z < u.Z
}

// IsMajor reports whether version v is considered to be a major Go release.
// For example, Go 1.14 and 1.13 are major Go releases.
func (v Version) IsMajor() bool { return v.Z == 0 }

// IsMinor reports whether version v is considered to be a minor Go release.
// For example, Go 1.14.1 and 1.13.9 are minor Go releases.
func (v Version) IsMinor() bool { return v.Z != 0 }

// A Date represents the date (year, month, day) of a Go release.
//
// This type does not include location information, and
// therefore does not describe a unique 24-hour timespan.
type Date struct {
	Year  int        // Year (e.g., 2009).
	Month time.Month // Month of the year (January = 1, ...).
	Day   int        // Day of the month, starting at 1.
}

func (d Date) String() string {
	return d.Format("2006-01-02")
}

func (d Date) Format(format string) string {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC).Format(format)
}

// A Major describes a major Go release and its minor revisions.
type Major struct {
	*Release
	Minor []*Release // oldest first
}

var Majors []*Major = majors() // major versions, newest first

// majors returns a list of major versions, sorted newest first.
func majors() []*Major {
	byVersion := make(map[Version]*Major)
	for _, r := range Releases {
		v := r.Version
		v.Z = 0 // make major version
		m := byVersion[v]
		if m == nil {
			m = new(Major)
			byVersion[v] = m
		}
		if r.Version.Z == 0 {
			m.Release = r
		} else {
			m.Minor = append(m.Minor, r)
		}
	}

	var majors []*Major
	for _, m := range byVersion {
		majors = append(majors, m)
		// minors oldest first
		sort.Slice(m.Minor, func(i, j int) bool {
			return m.Minor[i].Version.Before(m.Minor[j].Version)
		})
	}
	// majors newest first
	sort.Slice(majors, func(i, j int) bool {
		return !majors[i].Version.Before(majors[j].Version)
	})
	return majors
}

// ComponentsAndPackages joins components and packages involved
// in a Go release for the purposes of being displayed on the
// release history page, keeping English grammar rules in mind.
//
// The different special cases are:
//
//	c1
//	c1 and c2
//	c1, c2, and c3
//
//	the p1 package
//	the p1 and p2 packages
//	the p1, p2, and p3 packages
//
//	c1 and [1 package]
//	c1, and [2 or more packages]
//	c1, c2, and [1 or more packages]
func (f *FixSummary) ComponentsAndPackages() template.HTML {
	var buf strings.Builder

	// List components, if any.
	for i, comp := range f.Components {
		if len(f.Packages) == 0 {
			// No packages, so components are joined with more rules.
			switch {
			case i != 0 && len(f.Components) == 2:
				buf.WriteString(" and ")
			case i != 0 && len(f.Components) >= 3 && i != len(f.Components)-1:
				buf.WriteString(", ")
			case i != 0 && len(f.Components) >= 3 && i == len(f.Components)-1:
				buf.WriteString(", and ")
			}
		} else {
			// When there are packages, all components are comma-separated.
			if i != 0 {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(string(comp))
	}

	// Join components and packages using a comma and/or "and" as needed.
	if len(f.Components) > 0 && len(f.Packages) > 0 {
		if len(f.Components)+len(f.Packages) >= 3 {
			buf.WriteString(",")
		}
		buf.WriteString(" and ")
	}

	// List packages, if any.
	if len(f.Packages) > 0 {
		buf.WriteString("the ")
	}
	for i, pkg := range f.Packages {
		switch {
		case i != 0 && len(f.Packages) == 2:
			buf.WriteString(" and ")
		case i != 0 && len(f.Packages) >= 3 && i != len(f.Packages)-1:
			buf.WriteString(", ")
		case i != 0 && len(f.Packages) >= 3 && i == len(f.Packages)-1:
			buf.WriteString(", and ")
		}
		buf.WriteString("<code>" + html.EscapeString(pkg) + "</code>")
	}
	switch {
	case len(f.Packages) == 1:
		buf.WriteString(" package")
	case len(f.Packages) >= 2:
		buf.WriteString(" packages")
	}

	return template.HTML(buf.String())
}
