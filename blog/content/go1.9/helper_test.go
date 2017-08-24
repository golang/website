// Copyright 2017 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build OMIT

package p

import "testing"

func failure(t *testing.T) {
	t.Helper() // This call silences this function in error reports.
	t.Fatal("failure")
}

func Test(t *testing.T) {
	failure(t)
}
