// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package env provides environment information for the golangorg server
// running on golang.org.
package env

import (
	"log"
	"os"
	"strconv"
)

var requireDLSecretKey = boolEnv("GOLANGORG_REQUIRE_DL_SECRET_KEY")

// RequireDLSecretKey reports whether the download server secret key
// is expected to already exist, and the download server should panic
// on missing key instead of creating a new one.
func RequireDLSecretKey() bool {
	return requireDLSecretKey
}

func boolEnv(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		// TODO(dmitshur): In the future, consider detecting if running in App Engine,
		// and if so, making the environment variables mandatory rather than optional.
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatalf("environment variable %s (%q) must be a boolean", key, v)
	}
	return b
}
