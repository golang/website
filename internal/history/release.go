// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package history stores historical data for the Go project.
package history

import (
	"fmt"
	"html/template"
	"time"
)

// Release contains release metadata and a summary of release content.
type Release struct {
	// Release metadata.
	Date     Date // Date of the release.
	Security bool // Security release.
	Future   bool // Future is whether the release hasn't happened yet.

	// Release content summary.
	Quantifier    string          // Optional quantifier. Empty string for unspecified amount of fixes (typical), "a" for a single fix, "two", "three" for multiple fixes, etc.
	Components    []template.HTML // Components involved. For example, "cgo", "the <code>go</code> command", "the runtime", etc.
	Packages      []string        // Packages involved. For example, "net/http", "crypto/x509", etc.
	More          template.HTML   // Additional release content.
	CustomSummary template.HTML   // CustomSummary, if non-empty, replaces the entire release content summary with custom HTML.
}

// Releases summarizes the changes between official stable releases of Go.
//
// It contains entries for releases of Go 1.9 and newer.
// Older releases are listed in doc/devel/release.html.
var Releases = map[GoVer]Release{
	{1, 14, 7}: {
		Date:     Date{2020, 8, 6},
		Security: true,
		Packages: []string{"encoding/binary"},
	},
	{1, 13, 15}: {
		Date:     Date{2020, 8, 6},
		Security: true,
		Packages: []string{"encoding/binary"},
	},

	{1, 14, 6}: {
		Date: Date{2020, 7, 16},

		Components: []template.HTML{"the <code>go</code> command", "the compiler", "the linker", "vet"},
		Packages:   []string{"database/sql", "encoding/json", "net/http", "reflect", "testing"},
	},
	{1, 13, 14}: {
		Date: Date{2020, 7, 16},

		Components: []template.HTML{"the compiler", "vet"},
		Packages:   []string{"database/sql", "net/http", "reflect"},
	},

	{1, 14, 5}: {
		Date:     Date{2020, 7, 14},
		Security: true,

		Packages: []string{"crypto/x509", "net/http"},
	},
	{1, 13, 13}: {
		Date:     Date{2020, 7, 14},
		Security: true,

		Packages: []string{"crypto/x509", "net/http"},
	},

	{1, 14, 4}: {
		Date: Date{2020, 6, 1},

		Components: []template.HTML{"the <code>go</code> <code>doc</code> command", "the runtime"},
		Packages:   []string{"encoding/json", "os"},
	},
	{1, 13, 12}: {
		Date: Date{2020, 6, 1},

		Components: []template.HTML{"the runtime"},
		Packages:   []string{"go/types", "math/big"},
	},

	{1, 14, 3}: {
		Date: Date{2020, 5, 14},

		Components: []template.HTML{"cgo", "the compiler", "the runtime"},
		Packages:   []string{"go/doc", "math/big"},
	},
	{1, 13, 11}: {
		Date: Date{2020, 5, 14},

		Components: []template.HTML{"the compiler"},
	},

	{1, 14, 2}: {
		Date: Date{2020, 4, 8},

		Components: []template.HTML{"cgo", "the go command", "the runtime"},
		Packages:   []string{"os/exec", "testing"},
	},
	{1, 13, 10}: {
		Date: Date{2020, 4, 8},

		Components: []template.HTML{"the go command", "the runtime"},
		Packages:   []string{"os/exec", "time"},
	},

	{1, 14, 1}: {
		Date: Date{2020, 3, 19},

		Components: []template.HTML{"the go command", "tools", "the runtime"},
	},
	{1, 13, 9}: {
		Date: Date{2020, 3, 19},

		Components: []template.HTML{"the go command", "tools", "the runtime", "the toolchain"},
		Packages:   []string{"crypto/cypher"},
	},

	{1, 14, 0}: {
		Date: Date{2020, 2, 25},
	},

	{1, 13, 8}: {
		Date: Date{2020, 2, 12},

		Components: []template.HTML{"the runtime"},
		Packages:   []string{"crypto/x509", "net/http"},
	},
	{1, 12, 17}: {
		Date: Date{2020, 2, 12},

		Quantifier: "a",
		Components: []template.HTML{"the runtime"},
	},

	{1, 13, 7}: {
		Date:     Date{2020, 1, 28},
		Security: true,

		Quantifier: "two",
		Packages:   []string{"crypto/x509"},
	},
	{1, 12, 16}: {
		Date:     Date{2020, 1, 28},
		Security: true,

		Quantifier: "two",
		Packages:   []string{"crypto/x509"},
	},

	{1, 13, 6}: {
		Date: Date{2020, 1, 9},

		Components: []template.HTML{"the runtime"},
		Packages:   []string{"net/http"},
	},
	{1, 12, 15}: {
		Date: Date{2020, 1, 9},

		Components: []template.HTML{"the runtime"},
		Packages:   []string{"net/http"},
	},

	{1, 13, 5}: {
		Date: Date{2019, 12, 4},

		Components: []template.HTML{"the go command", "the runtime", "the linker"},
		Packages:   []string{"net/http"},
	},
	{1, 12, 14}: {
		Date: Date{2019, 12, 4},

		Quantifier: "a",
		Components: []template.HTML{"the runtime"},
	},

	{1, 13, 4}: {
		Date: Date{2019, 10, 31},

		Packages: []string{"net/http", "syscall"},
		More: `It also fixes an issue on macOS 10.15 Catalina
where the non-notarized installer and binaries were being
<a href="https://golang.org/issue/34986">rejected by Gatekeeper</a>.`,
	},
	{1, 12, 13}: {
		Date: Date{2019, 10, 31},

		CustomSummary: `fixes an issue on macOS 10.15 Catalina
where the non-notarized installer and binaries were being
<a href="https://golang.org/issue/34986">rejected by Gatekeeper</a>.
Only macOS users who hit this issue need to update.`,
	},

	{1, 13, 3}: {
		Date: Date{2019, 10, 17},

		Components: []template.HTML{"the go command", "the toolchain", "the runtime"},
		Packages:   []string{"syscall", "net", "net/http", "crypto/ecdsa"},
	},
	{1, 12, 12}: {
		Date: Date{2019, 10, 17},

		Components: []template.HTML{"the go command", "runtime"},
		Packages:   []string{"syscall", "net"},
	},

	{1, 13, 2}: {
		Date:     Date{2019, 10, 17},
		Security: true,

		Components: []template.HTML{"the compiler"},
		Packages:   []string{"crypto/dsa"},
	},
	{1, 12, 11}: {
		Date:     Date{2019, 10, 17},
		Security: true,

		Packages: []string{"crypto/dsa"},
	},

	{1, 13, 1}: {
		Date:     Date{2019, 9, 25},
		Security: true,

		Packages: []string{"net/http", "net/textproto"},
	},
	{1, 12, 10}: {
		Date:     Date{2019, 9, 25},
		Security: true,

		Packages: []string{"net/http", "net/textproto"},
	},

	{1, 13, 0}: {
		Date: Date{2019, 9, 3},
	},

	{1, 12, 9}: {
		Date: Date{2019, 8, 15},

		Components: []template.HTML{"the linker"},
		Packages:   []string{"os", "math/big"},
	},

	{1, 12, 8}: {
		Date:     Date{2019, 8, 13},
		Security: true,

		Packages: []string{"net/http", "net/url"},
	},
	{1, 11, 13}: {
		Date:     Date{2019, 8, 13},
		Security: true,

		Packages: []string{"net/http", "net/url"},
	},

	{1, 12, 7}: {
		Date: Date{2019, 7, 8},

		Components: []template.HTML{"cgo", "the compiler", "the linker"},
	},
	{1, 11, 12}: {
		Date: Date{2019, 7, 8},

		Components: []template.HTML{"the compiler", "the linker"},
	},

	{1, 12, 6}: {
		Date: Date{2019, 6, 11},

		Components: []template.HTML{"the compiler", "the linker", "the go command"},
		Packages:   []string{"crypto/x509", "net/http", "os"},
	},
	{1, 11, 11}: {
		Date: Date{2019, 6, 11},

		Quantifier: "a",
		Packages:   []string{"crypto/x509"},
	},

	{1, 12, 5}: {
		Date: Date{2019, 5, 6},

		Components: []template.HTML{"the compiler", "the linker", "the go command", "the runtime"},
		Packages:   []string{"os"},
	},
	{1, 11, 10}: {
		Date: Date{2019, 5, 6},

		Components: []template.HTML{"the runtime", "the linker"},
	},

	{1, 12, 4}: {
		Date: Date{2019, 4, 11},

		CustomSummary: `fixes an issue where using the prebuilt binary
releases on older versions of GNU/Linux
<a href="https://golang.org/issues/31293">led to failures</a>
when linking programs that used cgo.
Only Linux users who hit this issue need to update.`,
	},
	{1, 11, 9}: {
		Date: Date{2019, 4, 11},

		CustomSummary: `fixes an issue where using the prebuilt binary
releases on older versions of GNU/Linux
<a href="https://golang.org/issues/31293">led to failures</a>
when linking programs that used cgo.
Only Linux users who hit this issue need to update.`,
	},

	{1, 12, 3}: {
		Date: Date{2019, 4, 8},

		CustomSummary: `was accidentally released without its
intended fix. It is identical to go1.12.2, except for its version
number. The intended fix is in go1.12.4.`,
	},
	{1, 11, 8}: {
		Date: Date{2019, 4, 8},

		CustomSummary: `was accidentally released without its
intended fix. It is identical to go1.11.7, except for its version
number. The intended fix is in go1.11.9.`,
	},

	{1, 12, 2}: {
		Date: Date{2019, 4, 5},

		Components: []template.HTML{"the compiler", "the go command", "the runtime"},
		Packages:   []string{"doc", "net", "net/http/httputil", "os"},
	},
	{1, 11, 7}: {
		Date: Date{2019, 4, 5},

		Components: []template.HTML{"the runtime"},
		Packages:   []string{"net"},
	},

	{1, 12, 1}: {
		Date: Date{2019, 3, 14},

		Components: []template.HTML{"cgo", "the compiler", "the go command"},
		Packages:   []string{"fmt", "net/smtp", "os", "path/filepath", "sync", "text/template"},
	},
	{1, 11, 6}: {
		Date: Date{2019, 3, 14},

		Components: []template.HTML{"cgo", "the compiler", "linker", "runtime", "go command"},
		Packages:   []string{"crypto/x509", "encoding/json", "net", "net/url"},
	},

	{1, 12, 0}: {
		Date: Date{2019, 2, 25},
	},

	{1, 11, 5}: {
		Date:     Date{2019, 1, 23},
		Security: true,

		Quantifier: "a",
		Packages:   []string{"crypto/elliptic"},
	},
	{1, 10, 8}: {
		Date:     Date{2019, 1, 23},
		Security: true,

		Quantifier: "a",
		Packages:   []string{"crypto/elliptic"},
	},

	{1, 11, 4}: {
		Date: Date{2018, 12, 14},

		Components: []template.HTML{"cgo", "the compiler", "linker", "runtime", "documentation", "go command"},
		Packages:   []string{"net/http", "go/types"},
		More: `It includes a fix to a bug introduced in Go 1.11.3 that broke <code>go</code>
<code>get</code> for import path patterns containing "<code>...</code>".`,
	},
	{1, 10, 7}: {
		Date: Date{2018, 12, 14},

		// TODO: Modify to follow usual pattern, say it includes a fix to the go command.
		CustomSummary: `includes a fix to a bug introduced in Go 1.10.6
that broke <code>go</code> <code>get</code> for import path patterns containing
"<code>...</code>".
See the <a href="https://github.com/golang/go/issues?q=milestone%3AGo1.10.7+label%3ACherryPickApproved">
Go 1.10.7 milestone</a> on our issue tracker for details.`,
	},

	{1, 11, 3}: {
		Date:     Date{2018, 12, 12},
		Security: true,

		Quantifier: "three",
		Components: []template.HTML{`"go get"`},
		Packages:   []string{"crypto/x509"},
	},
	{1, 10, 6}: {
		Date:     Date{2018, 12, 12},
		Security: true,

		Quantifier: "three",
		Components: []template.HTML{`"go get"`},
		Packages:   []string{"crypto/x509"},
		More:       "It contains the same fixes as Go 1.11.3 and was released at the same time.",
	},

	{1, 11, 2}: {
		Date: Date{2018, 11, 2},

		Components: []template.HTML{"the compiler", "linker", "documentation", "go command"},
		Packages:   []string{"database/sql", "go/types"},
	},
	{1, 10, 5}: {
		Date: Date{2018, 11, 2},

		Components: []template.HTML{"the go command", "linker", "runtime"},
		Packages:   []string{"database/sql"},
	},

	{1, 11, 1}: {
		Date: Date{2018, 10, 1},

		Components: []template.HTML{"the compiler", "documentation", "go command", "runtime"},
		Packages:   []string{"crypto/x509", "encoding/json", "go/types", "net", "net/http", "reflect"},
	},

	{1, 11, 0}: {
		Date: Date{2018, 8, 24},
	},
	{1, 10, 4}: {
		Date: Date{2018, 8, 24},

		Components: []template.HTML{"the go command", "linker"},
		Packages:   []string{"net/http", "mime/multipart", "ld/macho", "bytes", "strings"},
	},

	{1, 10, 3}: {
		Date: Date{2018, 6, 5},

		Components: []template.HTML{"the go command"},
		Packages:   []string{"crypto/tls", "crypto/x509", "strings"},
		More: `In particular, it adds <a href="https://go.googlesource.com/go/+/d4e21288e444d3ffd30d1a0737f15ea3fc3b8ad9">
minimal support to the go command for the vgo transition</a>.`,
	},
	{1, 9, 7}: {
		Date: Date{2018, 6, 5},

		Components: []template.HTML{"the go command"},
		Packages:   []string{"crypto/x509", "strings"},
		More: `In particular, it adds <a href="https://go.googlesource.com/go/+/d4e21288e444d3ffd30d1a0737f15ea3fc3b8ad9">
minimal support to the go command for the vgo transition</a>.`,
	},

	{1, 10, 2}: {
		Date: Date{2018, 5, 1},

		Components: []template.HTML{"the compiler", "linker", "go command"},
	},
	{1, 9, 6}: {
		Date: Date{2018, 5, 1},

		Components: []template.HTML{"the compiler", "go command"},
	},

	{1, 10, 1}: {
		Date: Date{2018, 3, 28},

		Components: []template.HTML{"the compiler", "runtime"},
		Packages:   []string{"archive/zip", "crypto/tls", "crypto/x509", "encoding/json", "net", "net/http", "net/http/pprof"},
	},
	{1, 9, 5}: {
		Date: Date{2018, 3, 28},

		Components: []template.HTML{"the compiler", "go command"},
		Packages:   []string{"net/http/pprof"},
	},

	{1, 10, 0}: {
		Date: Date{2018, 2, 16},
	},

	{1, 9, 4}: {
		Date:     Date{2018, 2, 7},
		Security: true,

		Quantifier: "a",
		Components: []template.HTML{`"go get"`},
	},

	{1, 9, 3}: {
		Date: Date{2018, 1, 22},

		Components: []template.HTML{"the compiler", "runtime"},
		Packages:   []string{"database/sql", "math/big", "net/http", "net/url"},
	},

	{1, 9, 2}: {
		Date: Date{2017, 10, 25},

		Components: []template.HTML{"the compiler", "linker", "runtime", "documentation", "<code>go</code> command"},
		Packages:   []string{"crypto/x509", "database/sql", "log", "net/smtp"},
		More: `It includes a fix to a bug introduced in Go 1.9.1 that broke <code>go</code> <code>get</code>
of non-Git repositories under certain conditions.`,
	},

	{1, 9, 1}: {
		Date:     Date{2017, 10, 4},
		Security: true,

		Quantifier: "two",
	},

	{1, 9, 0}: {
		Date: Date{2017, 8, 24},
	},
}

// GoVer represents a Go release version.
//
// In contrast to Semantic Versioning 2.0.0,
// trailing zero components are omitted,
// a version like Go 1.14 is considered a major Go release,
// a version like Go 1.14.1 is considered a minor Go release.
//
// See proposal golang.org/issue/32450 for background, details,
// and a discussion of the costs involved in making a change.
type GoVer struct {
	X int // X is the 1st component of a Go X.Y.Z version. It must be 1 or higher.
	Y int // Y is the 2nd component of a Go X.Y.Z version. It must be 0 or higher.
	Z int // Z is the 3rd component of a Go X.Y.Z version. It must be 0 or higher.
}

// String returns the Go release version string,
// like "1.14", "1.14.1", "1.14.2", and so on.
func (v GoVer) String() string {
	switch {
	case v.Z != 0:
		return fmt.Sprintf("%d.%d.%d", v.X, v.Y, v.Z)
	case v.Y != 0:
		return fmt.Sprintf("%d.%d", v.X, v.Y)
	default:
		return fmt.Sprintf("%d", v.X)
	}
}

// IsMajor reports whether version v is considered to be a major Go release.
// For example, Go 1.14 and 1.13 are major Go releases.
func (v GoVer) IsMajor() bool { return v.Z == 0 }

// IsMinor reports whether version v is considered to be a minor Go release.
// For example, Go 1.14.1 and 1.13.9 are minor Go releases.
func (v GoVer) IsMinor() bool { return v.Z != 0 }

// Date represents the date (year, month, day) of a Go release.
//
// This type does not include location information, and
// therefore does not describe a unique 24-hour timespan.
type Date struct {
	Year  int        // Year (e.g., 2009).
	Month time.Month // Month of the year (January = 1, ...).
	Day   int        // Day of the month, starting at 1.
}
