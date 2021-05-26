// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Regression tests to run against a production instance of golangorg.

package main_test

import (
	"flag"
	"strings"
	"testing"

	"golang.org/x/website/internal/webtest"
)

var host = flag.String("regtest.host", "", "host to run regression test against")

func TestLiveServer(t *testing.T) {
	*host = strings.TrimSuffix(*host, "/")
	if *host == "" {
		t.Skip("regtest.host flag missing.")
	}

	webtest.TestServer(t, "testdata/live.txt", *host)
}
