// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package history

import (
	"sort"
	"testing"
	"time"
)

func dateTime(d Date) time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
}

func TestReleases(t *testing.T) {
	seen := make(map[Version]bool)
	lastMinor := make(map[Version]int)
	var last time.Time
	for _, r := range Releases {
		if seen[r.Version] {
			t.Errorf("version %v duplicated", r.Version)
			continue
		}
		seen[r.Version] = true
		rt := dateTime(r.Date)
		if !last.IsZero() && rt.After(last) {
			t.Errorf("version %v out of order: %v vs %v", r.Version, r.Date, last.Format("2006-01-02"))
		}
		last = rt
		major := r.Version
		major.Z = 0
		if z, ok := lastMinor[major]; ok && r.Version.Z != z-1 {
			old := major
			old.Z = z
			t.Errorf("jumped from %v to %v", old, r.Version)
		}
		lastMinor[major] = r.Version.Z
		if r.Security != nil {
			if p := r.Security.Packages; !sort.StringsAreSorted(p) {
				t.Errorf("version %v security fix packages out of order: %v vs %v", r.Version, p, sorted(p))
			}
		}
		if r.Bug != nil {
			if p := r.Bug.Packages; !sort.StringsAreSorted(p) {
				t.Errorf("version %v bug fix packages out of order: %v vs %v", r.Version, p, sorted(p))
			}
		}
	}
}

func sorted(s []string) []string {
	s = append([]string(nil), s...)
	sort.Strings(s)
	return s
}
