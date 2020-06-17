// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package short

import (
	"net/http/httptest"
	"testing"
)

func TestExtractKey(t *testing.T) {
	testCases := []struct {
		in                     string
		wantKey, wantRemaining string
		wantErr                bool
	}{
		{in: "/s/foo", wantKey: "foo", wantRemaining: ""},
		{in: "/s/foo/", wantKey: "foo", wantRemaining: "/"},
		{in: "/s/foo/bar/", wantKey: "foo", wantRemaining: "/bar/"},
		{in: "/s/foo.bar/baz", wantKey: "foo.bar", wantRemaining: "/baz"},
		{in: "/s/s/s/s", wantKey: "s", wantRemaining: "/s/s"},
		{in: "/", wantErr: true},
		{in: "/s/", wantErr: true},
		{in: "/s", wantErr: true},
		{in: "/t/foo", wantErr: true},
		{in: "/s/foo*", wantErr: true},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("GET", tc.in, nil)
		gotKey, gotRemaining, gotErr := extractKey(req)
		if gotKey != tc.wantKey || gotRemaining != tc.wantRemaining || (gotErr != nil) != tc.wantErr {
			t.Errorf("extractKey(%q) = (%q, %q, %v), want (%q, %q, err=%v)", tc.in, gotKey, gotRemaining, gotErr, tc.wantKey, tc.wantRemaining, tc.wantErr)
		}
	}
}
