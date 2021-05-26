// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"testing"

	"golang.org/x/website/internal/webtest"
)

func TestXHandler(t *testing.T) {
	webtest.TestHandler(t, "testdata/x.txt", http.HandlerFunc(xHandler))
}
