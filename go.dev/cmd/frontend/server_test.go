// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"golang.org/x/website/internal/webtest"
)

func TestWeb(t *testing.T) {
	h, err := NewHandler("../../_content")
	if err != nil {
		t.Fatal(err)
	}
	webtest.TestHandler(t, "testdata/*.txt", h)
}
