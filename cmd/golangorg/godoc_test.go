// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
	"time"
)

// buildGodoc builds the godoc executable.
// It returns its path, and a cleanup function.
//
// TODO(adonovan): opt: do this at most once, and do the cleanup
// exactly once.  How though?  There's no atexit.
func buildGodoc(t *testing.T) (bin string, cleanup func()) {
	if runtime.GOARCH == "arm" {
		t.Skip("skipping test on arm platforms; too slow")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("skipping test because 'go' command unavailable: %v", err)
	}

	tmp, err := ioutil.TempDir("", "godoc-regtest-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if cleanup == nil { // probably, go build failed.
			os.RemoveAll(tmp)
		}
	}()

	bin = filepath.Join(tmp, "godoc")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Building godoc: %v", err)
	}

	return bin, func() { os.RemoveAll(tmp) }
}

func serverAddress(t *testing.T) string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		ln, err = net.Listen("tcp6", "[::1]:0")
	}
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().String()
}

func waitForServerReady(t *testing.T, addr string) {
	waitForServer(t,
		fmt.Sprintf("http://%v/", addr),
		"The Go Programming Language",
		15*time.Second,
		false)
}

func waitForSearchReady(t *testing.T, addr string) {
	waitForServer(t,
		fmt.Sprintf("http://%v/search?q=FALLTHROUGH", addr),
		"The list of tokens.",
		2*time.Minute,
		false)
}

func waitUntilScanComplete(t *testing.T, addr string) {
	waitForServer(t,
		fmt.Sprintf("http://%v/pkg", addr),
		"Scan is not yet complete",
		2*time.Minute,
		true,
	)
	// setting reverse as true, which means this waits
	// until the string is not returned in the response anymore
}

const pollInterval = 200 * time.Millisecond

func waitForServer(t *testing.T, url, match string, timeout time.Duration, reverse bool) {
	// "health check" duplicated from x/tools/cmd/tipgodoc/tip.go
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		time.Sleep(pollInterval)
		res, err := http.Get(url)
		if err != nil {
			continue
		}
		rbody, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err == nil && res.StatusCode == http.StatusOK {
			if bytes.Contains(rbody, []byte(match)) && !reverse {
				return
			}
			if !bytes.Contains(rbody, []byte(match)) && reverse {
				return
			}
		}
	}
	t.Fatalf("Server failed to respond in %v", timeout)
}

// hasTag checks whether a given release tag is contained in the current version
// of the go binary.
func hasTag(t string) bool {
	for _, v := range build.Default.ReleaseTags {
		if t == v {
			return true
		}
	}
	return false
}

func killAndWait(cmd *exec.Cmd) {
	cmd.Process.Kill()
	cmd.Wait()
}

// Basic integration test for godoc HTTP interface.
func TestWeb(t *testing.T) {
	testWeb(t, false)
}

// Basic integration test for godoc HTTP interface.
func TestWebIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in -short mode")
	}
	testWeb(t, true)
}

// Basic integration test for godoc HTTP interface.
func testWeb(t *testing.T, withIndex bool) {
	if runtime.GOOS == "plan9" {
		t.Skip("skipping on plan9; fails to start up quickly enough")
	}
	bin, cleanup := buildGodoc(t)
	defer cleanup()
	addr := serverAddress(t)
	args := []string{fmt.Sprintf("-http=%s", addr)}
	if withIndex {
		args = append(args, "-index", "-index_interval=-1s")
	}
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Args[0] = "godoc"

	// Set GOPATH variable to non-existing path
	// and GOPROXY=off to disable module fetches.
	// We cannot just unset GOPATH variable because godoc would default it to ~/go.
	// (We don't want the indexer looking at the local workspace during tests.)
	cmd.Env = append(os.Environ(),
		"GOPATH=does_not_exist",
		"GOPROXY=off",
		"GO111MODULE=off")

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start godoc: %s", err)
	}
	defer killAndWait(cmd)

	if withIndex {
		waitForSearchReady(t, addr)
	} else {
		waitForServerReady(t, addr)
		waitUntilScanComplete(t, addr)
	}

	tests := []struct {
		path        string
		contains    []string // substring
		match       []string // regexp
		notContains []string
		needIndex   bool
		releaseTag  string // optional release tag that must be in go/build.ReleaseTags
	}{
		{
			path:     "/",
			contains: []string{"Go is an open source programming language"},
		},
		{
			path:     "/pkg/fmt/",
			contains: []string{"Package fmt implements formatted I/O"},
		},
		{
			path:     "/src/fmt/",
			contains: []string{"scan_test.go"},
		},
		{
			path:     "/src/fmt/print.go",
			contains: []string{"// Println formats using"},
		},
		{
			path: "/pkg",
			contains: []string{
				"Standard library",
				"Package fmt implements formatted I/O",
			},
			notContains: []string{
				"internal/syscall",
				"cmd/gc",
			},
		},
		{
			path: "/pkg/?m=all",
			contains: []string{
				"Standard library",
				"Package fmt implements formatted I/O",
				"internal/syscall/?m=all",
			},
			notContains: []string{
				"cmd/gc",
			},
		},
		{
			path: "/search?q=ListenAndServe",
			contains: []string{
				"/src",
			},
			notContains: []string{
				"/pkg/bootstrap",
			},
			needIndex: true,
		},
		{
			path: "/pkg/strings/",
			contains: []string{
				`href="/src/strings/strings.go"`,
			},
		},
		{
			path: "/cmd/compile/internal/amd64/",
			contains: []string{
				`href="/src/cmd/compile/internal/amd64/ssa.go"`,
			},
		},
		{
			path: "/pkg/math/bits/",
			contains: []string{
				`Added in Go 1.9`,
			},
		},
		{
			path: "/pkg/net/",
			contains: []string{
				`// IPv6 scoped addressing zone; added in Go 1.1`,
			},
		},
		{
			path: "/pkg/net/http/httptrace/",
			match: []string{
				`Got1xxResponse.*// Go 1\.11`,
			},
			releaseTag: "go1.11",
		},
		// Verify we don't add version info to a struct field added the same time
		// as the struct itself:
		{
			path: "/pkg/net/http/httptrace/",
			match: []string{
				`(?m)GotFirstResponseByte func\(\)\s*$`,
			},
		},
		// Remove trailing periods before adding semicolons:
		{
			path: "/pkg/database/sql/",
			contains: []string{
				"The number of connections currently in use; added in Go 1.11",
				"The number of idle connections; added in Go 1.11",
			},
			releaseTag: "go1.11",
		},
	}
	for _, test := range tests {
		if test.needIndex && !withIndex {
			continue
		}
		url := fmt.Sprintf("http://%s%s", addr, test.path)
		resp, err := http.Get(url)
		if err != nil {
			t.Errorf("GET %s failed: %s", url, err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		strBody := string(body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("GET %s: failed to read body: %s (response: %v)", url, err, resp)
		}
		isErr := false
		for _, substr := range test.contains {
			if test.releaseTag != "" && !hasTag(test.releaseTag) {
				continue
			}
			if !bytes.Contains(body, []byte(substr)) {
				t.Errorf("GET %s: wanted substring %q in body", url, substr)
				isErr = true
			}
		}
		for _, re := range test.match {
			if test.releaseTag != "" && !hasTag(test.releaseTag) {
				continue
			}
			if ok, err := regexp.MatchString(re, strBody); !ok || err != nil {
				if err != nil {
					t.Fatalf("Bad regexp %q: %v", re, err)
				}
				t.Errorf("GET %s: wanted to match %s in body", url, re)
				isErr = true
			}
		}
		for _, substr := range test.notContains {
			if bytes.Contains(body, []byte(substr)) {
				t.Errorf("GET %s: didn't want substring %q in body", url, substr)
				isErr = true
			}
		}
		if isErr {
			t.Errorf("GET %s: got:\n%s", url, body)
		}
	}
}
