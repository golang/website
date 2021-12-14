// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dl

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"sort"
	"testing"
)

func TestServeJSON(t *testing.T) {
	data := listTemplateData{
		Stable:   []Release{{Version: "Stable"}},
		Unstable: []Release{{Version: "Unstable"}},
		Archive:  []Release{{Version: "Archived"}},
	}
	testCases := []struct {
		desc     string
		method   string
		target   string
		status   int
		versions []string
	}{
		{
			desc:     "basic",
			method:   "GET",
			target:   "/",
			status:   200,
			versions: []string{"Stable"},
		},
		{
			desc:     "include all versions",
			method:   "GET",
			target:   "/?include=all",
			status:   200,
			versions: []string{"Stable", "Unstable", "Archived"},
		},
		{
			desc:   "CORS preflight request",
			method: "OPTIONS",
			target: "/",
			status: 204,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.target, nil)
			w := httptest.NewRecorder()
			serveJSON(w, r, &data)

			resp := w.Result()
			defer resp.Body.Close()
			if got, want := resp.StatusCode, tc.status; got != want {
				t.Errorf("Response status code = %d; want %d", got, want)
			}
			for k, v := range map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
			} {
				if got, want := resp.Header.Get(k), v; got != want {
					t.Errorf("%s = %q; want %q", k, got, want)
				}
			}
			if tc.versions == nil {
				return
			}

			if got, want := resp.Header.Get("Content-Type"), "application/json"; got != want {
				t.Errorf("Content-Type = %q; want %q", got, want)
			}
			var rs []Release
			if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
				t.Fatalf("json.Decode: got unexpected error: %v", err)
			}
			sort.Slice(rs, func(i, j int) bool {
				return rs[i].Version < rs[j].Version
			})
			sort.Strings(tc.versions)
			if got, want := len(rs), len(tc.versions); got != want {
				t.Fatalf("Number of releases = %d; want %d", got, want)
			}
			for i := range rs {
				if got, want := rs[i].Version, tc.versions[i]; got != want {
					t.Errorf("Got version %q; want %q", got, want)
				}
			}
		})
	}
}

func TestSortedJSON(t *testing.T) {
	// sortedVersion is populated with versions sorted from newest to oldest.
	var sortedVersions []File
	for maj := 30; maj >= 0; maj-- {
		for min := 30; min >= 0; min-- {
			minStr := ""
			if min > 0 {
				minStr = fmt.Sprintf(".%d", min)
			}
			for _, tail := range []string{"", "rc2", "rc1", "beta2", "beta1"} {
				version := fmt.Sprintf("go1.%d%s%s", maj, minStr, tail)
				sortedVersions = append(sortedVersions, File{Version: version})
			}
		}
	}
	shuffledVersions := append([]File{}, sortedVersions...)
	rand.Shuffle(len(shuffledVersions), func(i, j int) {
		shuffledVersions[i], shuffledVersions[j] = shuffledVersions[j], shuffledVersions[i]
	})

	d := listTemplateData{}
	d.Stable, d.Unstable, d.Archive = filesToReleases(shuffledVersions)
	r := httptest.NewRequest("GET", "/?mode=json&include=all", nil)
	w := httptest.NewRecorder()
	serveJSON(w, r, &d)
	rsp := w.Result()
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		t.Errorf("response status code = %d; want 200", rsp.StatusCode)
	}
	var rs []Release
	if err := json.NewDecoder(rsp.Body).Decode(&rs); err != nil {
		t.Fatalf("json.Decode: got unexpected error: %v", err)
	}
	if len(rs) != len(sortedVersions) {
		t.Fatalf("got %d versions; want %d versions", len(rs), len(sortedVersions))
	}
	for i := range rs {
		if rs[i].Version != sortedVersions[i].Version {
			t.Errorf("unexpected order at index %d of %s; want %s", i, rs[i].Version, sortedVersions[i].Version)
		}
	}
}
