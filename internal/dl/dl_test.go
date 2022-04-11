// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dl

import (
	"sort"
	"strings"
	"testing"
)

func TestParseVersion(t *testing.T) {
	for _, c := range []struct {
		in       string
		maj, min int
		tail     string
	}{
		{"go1.5", 5, 0, ""},
		{"go1.5beta1", 5, 0, "beta1"},
		{"go1.5.1", 5, 1, ""},
		{"go1.5.1rc1", 5, 1, "rc1"},
	} {
		maj, min, tail := parseVersion(c.in)
		if maj != c.maj || min != c.min || tail != c.tail {
			t.Errorf("parseVersion(%q) = %v, %v, %q; want %v, %v, %q",
				c.in, maj, min, tail, c.maj, c.min, c.tail)
		}
	}
}

func TestFileOrder(t *testing.T) {
	fs := []File{
		{Filename: "go1.16.src.tar.gz", Version: "go1.16", OS: "", Arch: "", Kind: "source"},
		{Filename: "go1.16.1.src.tar.gz", Version: "go1.16.1", OS: "", Arch: "", Kind: "source"},
		{Filename: "go1.16.linux-amd64.tar.gz", Version: "go1.16", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16.1.linux-amd64.tar.gz", Version: "go1.16.1", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16.darwin-amd64.tar.gz", Version: "go1.16", OS: "darwin", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16.darwin-amd64.pkg", Version: "go1.16", OS: "darwin", Arch: "amd64", Kind: "installer"},
		{Filename: "go1.16.darwin-arm64.tar.gz", Version: "go1.16", OS: "darwin", Arch: "arm64", Kind: "archive"},
		{Filename: "go1.16.darwin-arm64.pkg", Version: "go1.16", OS: "darwin", Arch: "arm64", Kind: "installer"},
		{Filename: "go1.16beta1.linux-amd64.tar.gz", Version: "go1.16beta1", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16beta2.linux-amd64.tar.gz", Version: "go1.16beta2", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16rc1.linux-amd64.tar.gz", Version: "go1.16rc1", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.15.linux-amd64.tar.gz", Version: "go1.15", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.15.2.linux-amd64.tar.gz", Version: "go1.15.2", OS: "linux", Arch: "amd64", Kind: "archive"},
	}
	sort.Sort(fileOrder(fs))
	var s []string
	for _, f := range fs {
		s = append(s, f.Filename)
	}
	got := strings.Join(s, "\n")
	want := strings.Join([]string{
		"go1.16.1.src.tar.gz",
		"go1.16.1.linux-amd64.tar.gz",
		"go1.16.src.tar.gz",
		"go1.16.darwin-amd64.tar.gz",
		"go1.16.darwin-amd64.pkg",
		"go1.16.darwin-arm64.tar.gz",
		"go1.16.darwin-arm64.pkg",
		"go1.16.linux-amd64.tar.gz",
		"go1.15.2.linux-amd64.tar.gz",
		"go1.15.linux-amd64.tar.gz",
		"go1.16rc1.linux-amd64.tar.gz",
		"go1.16beta2.linux-amd64.tar.gz",
		"go1.16beta1.linux-amd64.tar.gz",
	}, "\n")
	if got != want {
		t.Errorf("sort order is\n%s\nwant:\n%s", got, want)
	}
}

func TestFilesToReleases(t *testing.T) {
	fs := []File{
		{Version: "go1.7.4", OS: "darwin"},
		{Version: "go1.7.4", OS: "windows"},
		{Version: "go1.7", OS: "darwin"},
		{Version: "go1.7", OS: "windows"},
		{Version: "go1.6.2", OS: "darwin"},
		{Version: "go1.6.2", OS: "windows"},
		{Version: "go1.6", OS: "darwin"},
		{Version: "go1.6", OS: "windows"},
		{Version: "go1.5.2", OS: "darwin"},
		{Version: "go1.5.2", OS: "windows"},
		{Version: "go1.5", OS: "darwin"},
		{Version: "go1.5", OS: "windows"},
		{Version: "go1.5beta1", OS: "windows"},
	}
	stable, unstable, archive := filesToReleases(fs)
	if got, want := list(stable), "go1.7.4, go1.6.2"; got != want {
		t.Errorf("stable = %q; want %q", got, want)
	}
	if got, want := list(unstable), ""; got != want {
		t.Errorf("unstable = %q; want %q", got, want)
	}
	if got, want := list(archive), "go1.7, go1.6, go1.5.2, go1.5, go1.5beta1"; got != want {
		t.Errorf("archive = %q; want %q", got, want)
	}
}

func TestHighlightedFiles(t *testing.T) {
	fs := []File{
		{Filename: "go1.17.src.tar.gz", Version: "go1.17", OS: "", Arch: "", Kind: "source"},
		{Filename: "go1.17.linux-386.tar.gz", Version: "go1.17", OS: "linux", Arch: "386", Kind: "archive"},
		{Filename: "go1.17.linux-amd64.tar.gz", Version: "go1.17", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.17.darwin-amd64.tar.gz", Version: "go1.17", OS: "darwin", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.17.darwin-amd64.pkg", Version: "go1.17", OS: "darwin", Arch: "amd64", Kind: "installer"},
		{Filename: "go1.17.darwin-arm64.tar.gz", Version: "go1.17", OS: "darwin", Arch: "arm64", Kind: "archive"},
		{Filename: "go1.17.darwin-arm64.pkg", Version: "go1.17", OS: "darwin", Arch: "arm64", Kind: "installer"},
		{Filename: "go1.17.windows-386.zip", Version: "go1.17", OS: "windows", Arch: "386", Kind: "archive"},
		{Filename: "go1.17.windows-386.msi", Version: "go1.17", OS: "windows", Arch: "386", Kind: "installer"},
		{Filename: "go1.17.windows-amd64.zip", Version: "go1.17", OS: "windows", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.17.windows-amd64.msi", Version: "go1.17", OS: "windows", Arch: "amd64", Kind: "installer"},
		{Filename: "go1.17.windows-arm64.zip", Version: "go1.17", OS: "windows", Arch: "arm64", Kind: "archive"},
		{Filename: "go1.17.windows-arm64.msi", Version: "go1.17", OS: "windows", Arch: "arm64", Kind: "installer"},
	}
	sort.Sort(fileOrder(fs))
	var highlighted []string
	for _, f := range fs {
		if !f.Highlight() {
			continue
		}
		highlighted = append(highlighted, f.Filename)
	}
	got := strings.Join(highlighted, "\n")
	want := strings.Join([]string{
		"go1.17.src.tar.gz",
		"go1.17.darwin-amd64.pkg",
		"go1.17.darwin-arm64.pkg",
		"go1.17.linux-amd64.tar.gz",
		"go1.17.windows-amd64.msi",
		"go1.17.windows-arm64.msi",
	}, "\n")
	if got != want {
		t.Errorf("highlighted files:\n%s\nwant:\n%s", got, want)
	}
}

func TestOldUnstableNotShown(t *testing.T) {
	fs := []File{
		{Version: "go1.7.4"},
		{Version: "go1.7"},
		{Version: "go1.7beta1"},
	}
	_, unstable, archive := filesToReleases(fs)
	if len(unstable) != 0 {
		t.Errorf("got unstable, want none")
	}
	if got, want := list(archive), "go1.7, go1.7beta1"; got != want {
		t.Errorf("archive = %q; want %q", got, want)
	}
}

// A new beta should show up under unstable, but not show up under archive. See golang.org/issue/29669.
func TestNewUnstableShownOnce(t *testing.T) {
	fs := []File{
		{Version: "go1.12beta2"},
		{Version: "go1.11.4"},
		{Version: "go1.11"},
		{Version: "go1.10.7"},
		{Version: "go1.10"},
		{Version: "go1.9"},
	}
	stable, unstable, archive := filesToReleases(fs)
	if got, want := list(stable), "go1.11.4, go1.10.7"; got != want {
		t.Errorf("stable = %q; want %q", got, want)
	}
	if got, want := list(unstable), "go1.12beta2"; got != want {
		t.Errorf("unstable = %q; want %q", got, want)
	}
	if got, want := list(archive), "go1.11, go1.10, go1.9"; got != want {
		t.Errorf("archive = %q; want %q", got, want)
	}
}

func TestUnstableShown(t *testing.T) {
	fs := []File{
		{Version: "go1.8beta2"},
		{Version: "go1.8rc1"},
		{Version: "go1.7.4"},
		{Version: "go1.7"},
		{Version: "go1.7beta1"},
	}
	_, unstable, archive := filesToReleases(fs)
	// Show RCs ahead of betas.
	if got, want := list(unstable), "go1.8rc1"; got != want {
		t.Errorf("unstable = %q; want %q", got, want)
	}
	if got, want := list(archive), "go1.7, go1.8beta2, go1.7beta1"; got != want {
		t.Errorf("archive = %q; want %q", got, want)
	}
}

func TestFilesToFeatured(t *testing.T) {
	fs := []File{
		{Filename: "go1.16.3.src.tar.gz", Version: "go1.16.3", OS: "", Arch: "", Kind: "source"},
		{Filename: "go1.16.3.darwin-amd64.tar.gz", Version: "go1.16.3", OS: "darwin", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16.3.darwin-amd64.pkg", Version: "go1.16.3", OS: "darwin", Arch: "amd64", Kind: "installer"},
		{Filename: "go1.16.3.darwin-arm64.tar.gz", Version: "go1.16.3", OS: "darwin", Arch: "arm64", Kind: "archive"},
		{Filename: "go1.16.3.darwin-arm64.pkg", Version: "go1.16.3", OS: "darwin", Arch: "arm64", Kind: "installer"},
		{Filename: "go1.16.3.freebsd-386.tar.gz", Version: "go1.16.3", OS: "freebsd", Arch: "386", Kind: "archive"},
		{Filename: "go1.16.3.freebsd-amd64.tar.gz", Version: "go1.16.3", OS: "freebsd", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16.3.linux-386.tar.gz", Version: "go1.16.3", OS: "linux", Arch: "386", Kind: "archive"},
		{Filename: "go1.16.3.linux-amd64.tar.gz", Version: "go1.16.3", OS: "linux", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16.3.linux-arm64.tar.gz", Version: "go1.16.3", OS: "linux", Arch: "arm64", Kind: "archive"},
		{Filename: "go1.16.3.linux-armv6l.tar.gz", Version: "go1.16.3", OS: "linux", Arch: "armv6l", Kind: "archive"},
		{Filename: "go1.16.3.linux-ppc64le.tar.gz", Version: "go1.16.3", OS: "linux", Arch: "ppc64le", Kind: "archive"},
		{Filename: "go1.16.3.linux-s390x.tar.gz", Version: "go1.16.3", OS: "linux", Arch: "s390x", Kind: "archive"},
		{Filename: "go1.16.3.windows-386.zip", Version: "go1.16.3", OS: "windows", Arch: "386", Kind: "archive"},
		{Filename: "go1.16.3.windows-386.msi", Version: "go1.16.3", OS: "windows", Arch: "386", Kind: "installer"},
		{Filename: "go1.16.3.windows-amd64.zip", Version: "go1.16.3", OS: "windows", Arch: "amd64", Kind: "archive"},
		{Filename: "go1.16.3.windows-amd64.msi", Version: "go1.16.3", OS: "windows", Arch: "amd64", Kind: "installer"},
	}
	featured := filesToFeatured(fs)
	var s []string
	for _, f := range featured {
		s = append(s, f.Filename)
	}
	got := strings.Join(s, "\n")
	want := strings.Join([]string{
		"go1.16.3.windows-amd64.msi",
		"go1.16.3.darwin-arm64.pkg",
		"go1.16.3.darwin-amd64.pkg",
		"go1.16.3.linux-amd64.tar.gz",
		"go1.16.3.src.tar.gz",
	}, "\n")
	if got != want {
		t.Errorf("featured files:\n%s\nwant:\n%s", got, want)
	}
}

// list returns a version list string for the given releases.
func list(rs []Release) string {
	var s string
	for i, r := range rs {
		if i > 0 {
			s += ", "
		}
		s += r.Version
	}
	return s
}
