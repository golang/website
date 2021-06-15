// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"golang.org/x/website/internal/webtest"
)

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

func init() {
	// TestWeb reinvokes the test binary (us) with -be-main
	// to simulate running the actual golangorg binary.
	if len(os.Args) >= 2 && os.Args[1] == "-be-main" {
		os.Args = os.Args[1:]
		os.Args[0] = "(golangorg)"
		main()
		os.Exit(0)
	}
}

func TestWeb(t *testing.T) {
	h := NewHandler("../../_content", runtime.GOROOT())
	files, err := filepath.Glob("testdata/*.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if filepath.ToSlash(file) != "testdata/live.txt" {
			webtest.TestHandler(t, file, h)
		}
	}
}

// Regression tests to run against a production instance of golangorg.

var host = flag.String("regtest.host", "", "host to run regression test against")

func TestLiveServer(t *testing.T) {
	*host = strings.TrimSuffix(*host, "/")
	if *host == "" {
		t.Skip("regtest.host flag missing.")
	}

	webtest.TestServer(t, "testdata/*.txt", *host)
}
