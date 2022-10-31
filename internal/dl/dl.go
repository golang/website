// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dl implements a simple Go downloads frontend server.
//
// It accepts HTTP POST requests to create a new download metadata entity,
// and lists entities with sorting and filtering.
//
// The list of downloads, as well as individual files, are served at:
//
//	https://go.dev/dl/
//	https://go.dev/dl/{file}
//
// An optional query param, mode=json, serves the list of stable release
// downloads in JSON format:
//
//	https://go.dev/dl/?mode=json
//
// An additional query param, include=all, when used with the mode=json
// query param, will serve a full list of available downloads, including
// unstable, stable, and archived releases, in JSON format:
//
//	https://go.dev/dl/?mode=json&include=all
//
// Releases returned in JSON modes are sorted by version, newest to oldest.
package dl

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	cacheKey      = "download_list_5" // increment if listTemplateData changes
	cacheDuration = time.Hour
)

// File represents a file on the go.dev downloads page.
// It should be kept in sync with the upload code in x/build/internal/relui.
type File struct {
	Filename       string    `json:"filename"`
	OS             string    `json:"os"`
	Arch           string    `json:"arch"`
	Version        string    `json:"version"`
	Checksum       string    `json:"-" datastore:",noindex"` // SHA1; deprecated
	ChecksumSHA256 string    `json:"sha256" datastore:",noindex"`
	Size           int64     `json:"size" datastore:",noindex"`
	Kind           string    `json:"kind"` // "archive", "installer", "source"
	Uploaded       time.Time `json:"-"`
}

func (f File) ChecksumType() string {
	if f.ChecksumSHA256 != "" {
		return "SHA256"
	}
	return "SHA1"
}

func (f File) PrettyArch() string { return pretty(f.Arch) }
func (f File) PrettyKind() string { return pretty(f.Kind) }

func (f File) PrettyChecksum() string {
	if f.ChecksumSHA256 != "" {
		return f.ChecksumSHA256
	}
	return f.Checksum
}

func (f File) PrettyOS() string {
	if f.OS == "darwin" {
		// Some older releases, like Go 1.4,
		// still contain "osx" in the filename.
		switch {
		case strings.Contains(f.Filename, "osx10.8"):
			return "OS X 10.8+"
		case strings.Contains(f.Filename, "osx10.6"):
			return "OS X 10.6+"
		}
	}
	return pretty(f.OS)
}

func (f File) PrettySize() string {
	const mb = 1 << 20
	if f.Size == 0 {
		return ""
	}
	if f.Size < mb {
		// All Go releases are >1mb, but handle this case anyway.
		return fmt.Sprintf("%v bytes", f.Size)
	}
	return fmt.Sprintf("%.0fMB", float64(f.Size)/mb)
}

var primaryPorts = map[string]bool{
	"darwin/amd64":  true,
	"darwin/arm64":  true,
	"linux/386":     true,
	"linux/amd64":   true,
	"linux/armv6l":  true,
	"linux/arm64":   true,
	"windows/386":   true,
	"windows/amd64": true,
}

func (f File) PrimaryPort() bool {
	if f.Kind == "source" {
		return true
	}
	return primaryPorts[f.OS+"/"+f.Arch]
}

func (f File) Highlight() bool {
	switch {
	case f.Kind == "source":
		return true
	case f.OS == "linux" && f.Arch == "amd64":
		return true
	case f.OS == "windows" && f.Kind == "installer" && (f.Arch == "amd64" || f.Arch == "arm64"):
		return true
	case f.OS == "darwin" && f.Kind == "installer" && !strings.Contains(f.Filename, "osx10.6"):
		return true
	}
	return false
}

// URL returns the canonical URL of the file.
func (f File) URL() string {
	// The download URL of a Go release file is /dl/{name}. It is handled by getHandler.
	// Use a relative URL so it works for any host like go.dev and golang.google.cn.
	// Don't shortcut to the redirect target here, we want canonical URLs to be visible. See issue 38713.
	return "/dl/" + f.Filename
}

type Release struct {
	Version        string `json:"version"`
	Stable         bool   `json:"stable"`
	Files          []File `json:"files"`
	Visible        bool   `json:"-"` // show files on page load
	SplitPortTable bool   `json:"-"` // whether files should be split by primary/other ports.
}

type Feature struct {
	// The File field will be filled in by the first stable File
	// whose name matches the given fileRE.
	File
	fileRE *regexp.Regexp

	Platform     string // "Microsoft Windows", "Apple macOS", "Linux"
	Requirements string // "Windows XP and above, 64-bit Intel Processor"
}

// featuredFiles lists the platforms and files to be featured
// at the top of the downloads page.
var featuredFiles = []Feature{
	{
		Platform:     "Microsoft Windows",
		Requirements: "Windows 7 or later, Intel 64-bit processor",
		fileRE:       regexp.MustCompile(`\.windows-amd64\.msi$`),
	},
	{
		Platform:     "Apple macOS (ARM64)",
		Requirements: "macOS 11 or later, Apple 64-bit processor",
		fileRE:       regexp.MustCompile(`\.darwin-arm64\.pkg$`),
	},
	{
		Platform:     "Apple macOS (x86-64)",
		Requirements: "macOS 10.13 or later, Intel 64-bit processor",
		fileRE:       regexp.MustCompile(`\.darwin-amd64\.pkg$`),
	},
	{
		Platform:     "Linux",
		Requirements: "Linux 2.6.32 or later, Intel 64-bit processor",
		fileRE:       regexp.MustCompile(`\.linux-amd64\.tar\.gz$`),
	},
	{
		Platform: "Source",
		fileRE:   regexp.MustCompile(`\.src\.tar\.gz$`),
	},
}

// data to send to the template; increment cacheKey if you change this.
type listTemplateData struct {
	Featured                  []Feature
	Stable, Unstable, Archive []Release
}

func filesToFeatured(fs []File) (featured []Feature) {
	for _, feature := range featuredFiles {
		for _, file := range fs {
			if feature.fileRE.MatchString(file.Filename) {
				feature.File = file
				featured = append(featured, feature)
				break
			}
		}
	}
	return
}

func filesToReleases(fs []File) (stable, unstable, archive []Release) {
	sort.Sort(fileOrder(fs))

	var r *Release
	var stableMaj, stableMin int
	add := func() {
		if r == nil {
			return
		}
		if !r.Stable {
			if len(unstable) != 0 {
				// Only show one (latest) unstable version,
				// consider the older ones to be archived.
				archive = append(archive, *r)
				return
			}
			maj, min, _ := parseVersion(r.Version)
			if maj < stableMaj || maj == stableMaj && min <= stableMin {
				// Display unstable version only if newer than the
				// latest stable release, otherwise consider it archived.
				archive = append(archive, *r)
				return
			}
			unstable = append(unstable, *r)
			return
		}

		// Reports whether the release is the most recent minor version of the
		// two most recent major versions.
		shouldAddStable := func() bool {
			if len(stable) >= 2 {
				// Show up to two stable versions.
				return false
			}
			if len(stable) == 0 {
				// Most recent stable version.
				stableMaj, stableMin, _ = parseVersion(r.Version)
				return true
			}
			if maj, _, _ := parseVersion(r.Version); maj == stableMaj {
				// Older minor version of most recent major version.
				return false
			}
			// Second most recent stable version.
			return true
		}
		if !shouldAddStable() {
			archive = append(archive, *r)
			return
		}

		// Split the file list into primary/other ports for the stable releases.
		// NOTE(cbro): This is only done for stable releases because maintaining the historical
		// nature of primary/other ports for older versions is infeasible.
		// If freebsd is considered primary some time in the future, we'd not want to
		// mark all of the older freebsd binaries as "primary".
		// It might be better if we set that as a flag when uploading.
		r.SplitPortTable = true
		r.Visible = true // Toggle open all stable releases.
		stable = append(stable, *r)
	}
	for _, f := range fs {
		if r == nil || f.Version != r.Version {
			add()
			r = &Release{
				Version: f.Version,
				Stable:  isStable(f.Version),
			}
		}
		r.Files = append(r.Files, f)
	}
	add()
	return
}

// isStable reports whether the version string v is a stable version.
func isStable(v string) bool {
	return !strings.Contains(v, "beta") && !strings.Contains(v, "rc")
}

type fileOrder []File

func (s fileOrder) Len() int      { return len(s) }
func (s fileOrder) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s fileOrder) Less(i, j int) bool {
	a, b := s[i], s[j]
	if av, bv := a.Version, b.Version; av != bv {
		// Put stable releases first.
		if isStable(av) != isStable(bv) {
			return isStable(av)
		}
		return versionLess(av, bv)
	}
	if a.OS != b.OS {
		return a.OS < b.OS
	}
	if a.Arch != b.Arch {
		return a.Arch < b.Arch
	}
	if a.Kind != b.Kind {
		return a.Kind < b.Kind
	}
	return a.Filename < b.Filename
}

func versionLess(a, b string) bool {
	maja, mina, ta := parseVersion(a)
	majb, minb, tb := parseVersion(b)
	if maja == majb {
		if mina == minb {
			if ta == "" {
				return true
			} else if tb == "" {
				return false
			}
			return ta >= tb
		}
		return mina >= minb
	}
	return maja >= majb
}

func parseVersion(v string) (maj, min int, tail string) {
	if i := strings.Index(v, "beta"); i > 0 {
		tail = v[i:]
		v = v[:i]
	}
	if i := strings.Index(v, "rc"); i > 0 {
		tail = v[i:]
		v = v[:i]
	}
	p := strings.Split(strings.TrimPrefix(v, "go1."), ".")
	maj, _ = strconv.Atoi(p[0])
	if len(p) < 2 {
		return
	}
	min, _ = strconv.Atoi(p[1])
	return
}

// validUser controls whether the named gomote user is allowed to upload
// Go release binaries via the /dl/upload endpoint.
func validUser(user string) bool {
	switch user {
	case "amedee", "cherryyz", "dmitshur", "drchase", "heschi", "mknyszek", "rakoczy", "thanm":
		return true
	case "relui":
		return true
	}
	return false
}

var (
	fileRe  = regexp.MustCompile(`^go[0-9a-z.]+\.[0-9a-z.-]+\.(tar\.gz|tar\.gz\.asc|pkg|msi|zip)$`)
	goGetRe = regexp.MustCompile(`^go[0-9a-z.]+\.[0-9a-z.-]+$`)
)

// pretty returns a human-readable version of the given OS, Arch, or Kind.
func pretty(s string) string {
	t, ok := prettyStrings[s]
	if !ok {
		return s
	}
	return t
}

var prettyStrings = map[string]string{
	"darwin":  "macOS",
	"freebsd": "FreeBSD",
	"linux":   "Linux",
	"windows": "Windows",

	"386":    "x86",
	"amd64":  "x86-64",
	"armv6l": "ARMv6",
	"arm64":  "ARM64",

	"archive":   "Archive",
	"installer": "Installer",
	"source":    "Source",
}
