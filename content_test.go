// Copyright 2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package website_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	// Keep golang.org/x/tour in our go.mod require list for use during test.
	_ "golang.org/x/tour/wc"
)

// Test that all the .go files inside the _content/tour directory build
// and execute (without checking for output correctness).
// Files that contain the build tag "nobuild" are not built.
// Files that contain the build tag "norun" are not executed.
func TestTourContent(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skipf("skipping because 'go' executable not available: %v", err)
	}

	err := filepath.Walk(filepath.Join("_content", "tour"), func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			if err := testSnippet(filepath.ToSlash(path), t.TempDir()); err != nil {
				t.Errorf("%v: %v", path, err)
			}
		})
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func testSnippet(path, scratch string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	build := string(bytes.SplitN(b, []byte{'\n'}, 2)[0])
	if !strings.HasPrefix(build, "//go:build ") {
		return errors.New("first line is not a go:build comment")
	}
	if !strings.Contains(build, "OMIT") {
		return errors.New(`build comment does not contain "OMIT"`)
	}

	if strings.Contains(build, "nobuild") {
		return nil
	}
	bin := filepath.Join(scratch, filepath.Base(path)+".exe")
	out, err := exec.Command("go", "build", "-o", bin, path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("build error: %v\noutput:\n%s", err, out)
	}
	defer os.Remove(bin)

	if strings.Contains(build, "norun") {
		return nil
	}
	out, err = exec.Command(bin).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\nOutput:\n%s", err, out)
	}
	return nil
}
