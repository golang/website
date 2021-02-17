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

	"golang.org/x/website/internal/godoc/golangorgenv"
)

var (
	requireDLSecretKey = boolEnv("GOLANGORG_REQUIRE_DL_SECRET_KEY")
)

// RequireDLSecretKey reports whether the download server secret key
// is expected to already exist, and the download server should panic
// on missing key instead of creating a new one.
func RequireDLSecretKey() bool {
	return requireDLSecretKey
}

// Use the golangorgenv package for common configuration, instead
// of duplicating it. This reduces the risk of divergence between
// the environment variables that this env package uses, and ones
// that golangorgenv uses.
//
// TODO(dmitshur): When the golang.org/x/tools/playground package becomes unused,
// and golang.org/x/tools/godoc is modified to accept configuration explicitly,
// the golang.org/x/tools/godoc/golangorgenv package can be deleted.
// At that time, its implementation can be inlined into this package, as needed.

// CheckCountry reports whether country restrictions should be enforced.
func CheckCountry() bool {
	return golangorgenv.CheckCountry()
}

// EnforceHosts reports whether host filtering should be enforced.
func EnforceHosts() bool {
	return golangorgenv.EnforceHosts()
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
