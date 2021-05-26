// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"golang.org/x/website/internal/webtest"
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
		15*time.Second)
}

const pollInterval = 200 * time.Millisecond

func waitForServer(t *testing.T, url, match string, timeout time.Duration) {
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
			if bytes.Contains(rbody, []byte(match)) {
				return
			}
		}
	}
	t.Fatalf("Server failed to respond in %v", timeout)
}

func killAndWait(cmd *exec.Cmd) {
	cmd.Process.Kill()
	cmd.Wait()
}

// Basic integration test for godoc HTTP interface.
func TestWeb(t *testing.T) {
	testWeb(t)
}

// Basic integration test for godoc HTTP interface.
func testWeb(t *testing.T) {
	if runtime.GOOS == "plan9" {
		t.Skip("skipping on plan9; fails to start up quickly enough")
	}
	bin, cleanup := buildGodoc(t)
	defer cleanup()
	addr := serverAddress(t)
	args := []string{fmt.Sprintf("-http=%s", addr)}
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Args[0] = "godoc"

	// Set GOPATH variable to non-existing path
	// and GOPROXY=off to disable module fetches.
	// We cannot just unset GOPATH variable because godoc would default it to ~/go.
	// (We don't want the server looking at the local workspace during tests.)
	cmd.Env = append(os.Environ(),
		"GOPATH=does_not_exist",
		"GOPROXY=off",
		"GO111MODULE=off")

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start godoc: %s", err)
	}
	defer killAndWait(cmd)

	waitForServerReady(t, addr)

	webtest.TestServer(t, "testdata/web.txt", addr)
}
