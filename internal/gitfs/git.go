// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gitfs presents a file tree downloaded from a remote Git repo as an in-memory fs.FS.
package gitfs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

// A Repo is a connection to a remote repository served over HTTP or HTTPS.
type Repo struct {
	url  string // trailing slash removed
	caps map[string]string
}

// NewRepo connects to a Git repository at the given http:// or https:// URL.
func NewRepo(url string) (*Repo, error) {
	r := &Repo{url: strings.TrimSuffix(url, "/")}
	if err := r.handshake(); err != nil {
		return nil, err
	}
	return r, nil
}

// handshake runs the initial Git opening handshake, learning the capabilities of the server.
// See https://git-scm.com/docs/protocol-v2#_initial_client_request.
func (r *Repo) handshake() error {
	req, _ := http.NewRequest("GET", r.url+"/info/refs?service=git-upload-pack", nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Git-Protocol", "version=2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("handshake: %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("handshake: %v\n%s", resp.Status, data)
	}
	if err != nil {
		return fmt.Errorf("handshake: reading body: %v", err)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/x-git-upload-pack-advertisement" {
		return fmt.Errorf("handshake: invalid response Content-Type: %v", ct)
	}

	pr := newPktLineReader(bytes.NewReader(data))
	lines, err := pr.Lines()
	if len(lines) == 1 && lines[0] == "# service=git-upload-pack" {
		lines, err = pr.Lines()
	}
	if err != nil {
		return fmt.Errorf("handshake: parsing response: %v", err)
	}
	caps := make(map[string]string)
	for _, line := range lines {
		verb, args, _ := strings.Cut(line, "=")
		caps[verb] = args
	}
	if _, ok := caps["version 2"]; !ok {
		return fmt.Errorf("handshake: not version 2: %q", lines)
	}
	r.caps = caps
	return nil
}

// Resolve looks up the given ref and returns the corresponding Hash.
func (r *Repo) Resolve(ref string) (Hash, error) {
	if h, err := parseHash(ref); err == nil {
		return h, nil
	}

	fail := func(err error) (Hash, error) {
		return Hash{}, fmt.Errorf("resolve %s: %v", ref, err)
	}
	refs, err := r.refs(ref)
	if err != nil {
		return fail(err)
	}
	for _, known := range refs {
		if known.name == ref {
			return known.hash, nil
		}
	}
	return fail(fmt.Errorf("unknown ref"))
}

// A ref is a single Git reference, like refs/heads/main, refs/tags/v1.0.0, or HEAD.
type ref struct {
	name string // "refs/heads/main", "refs/tags/v1.0.0", "HEAD"
	hash Hash   // hexadecimal hash
}

// refs executes an ls-refs command on the remote server
// to look up refs with the given prefixes.
// See https://git-scm.com/docs/protocol-v2#_ls_refs.
func (r *Repo) refs(prefixes ...string) ([]ref, error) {
	if _, ok := r.caps["ls-refs"]; !ok {
		return nil, fmt.Errorf("refs: server does not support ls-refs")
	}

	var buf bytes.Buffer
	pw := newPktLineWriter(&buf)
	pw.WriteString("command=ls-refs")
	pw.Delim()
	pw.WriteString("peel")
	pw.WriteString("symrefs")
	for _, prefix := range prefixes {
		pw.WriteString("ref-prefix " + prefix)
	}
	pw.Close()
	postbody := buf.Bytes()

	req, _ := http.NewRequest("POST", r.url+"/git-upload-pack", &buf)
	req.Header.Set("Content-Type", "application/x-git-upload-pack-request")
	req.Header.Set("Accept", "application/x-git-upload-pack-result")
	req.Header.Set("Git-Protocol", "version=2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refs: %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("refs: %v\n%s", resp.Status, data)
	}
	if err != nil {
		return nil, fmt.Errorf("refs: reading body: %v", err)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/x-git-upload-pack-result" {
		return nil, fmt.Errorf("refs: invalid response Content-Type: %v", ct)
	}

	var refs []ref
	lines, err := newPktLineReader(bytes.NewReader(data)).Lines()
	if err != nil {
		return nil, fmt.Errorf("refs: parsing response: %v %d\n%s\n%s", err, len(data), hex.Dump(postbody), hex.Dump(data))
	}
	for _, line := range lines {
		hash, rest, ok := strings.Cut(line, " ")
		if !ok {
			return nil, fmt.Errorf("refs: parsing response: invalid line: %q", line)
		}
		h, err := parseHash(hash)
		if err != nil {
			return nil, fmt.Errorf("refs: parsing response: invalid line: %q", line)
		}
		name, _, _ := strings.Cut(rest, " ")
		refs = append(refs, ref{hash: h, name: name})
	}
	return refs, nil
}

// Clone resolves the given ref to a hash and returns the corresponding fs.FS.
func (r *Repo) Clone(ref string) (Hash, fs.FS, error) {
	fail := func(err error) (Hash, fs.FS, error) {
		return Hash{}, nil, fmt.Errorf("clone %s: %v", ref, err)
	}
	h, err := r.Resolve(ref)
	if err != nil {
		return fail(err)
	}
	tfs, err := r.fetch(h)
	if err != nil {
		return fail(err)
	}
	return h, tfs, nil
}

// CloneHash returns the fs.FS for the given hash.
func (r *Repo) CloneHash(h Hash) (fs.FS, error) {
	tfs, err := r.fetch(h)
	if err != nil {
		return nil, fmt.Errorf("clone %s: %v", h, err)
	}
	return tfs, nil
}

// fetch returns the fs.FS for a given hash.
func (r *Repo) fetch(h Hash) (fs.FS, error) {
	// Fetch a shallow packfile from the remote server.
	// Shallow means it only contains the tree at that one commit,
	// not the entire history of the repo.
	// See https://git-scm.com/docs/protocol-v2#_fetch.
	opts, ok := r.caps["fetch"]
	if !ok {
		return nil, fmt.Errorf("fetch: server does not support fetch")
	}
	if !strings.Contains(" "+opts+" ", " shallow ") {
		return nil, fmt.Errorf("fetch: server does not support shallow fetch")
	}

	// Prepare and send request for pack file.
	var buf bytes.Buffer
	pw := newPktLineWriter(&buf)
	pw.WriteString("command=fetch")
	pw.Delim()
	pw.WriteString("deepen 1")
	pw.WriteString("want " + h.String())
	pw.WriteString("done")
	pw.Close()
	postbody := buf.Bytes()

	req, _ := http.NewRequest("POST", r.url+"/git-upload-pack", &buf)
	req.Header.Set("Content-Type", "application/x-git-upload-pack-request")
	req.Header.Set("Accept", "application/x-git-upload-pack-result")
	req.Header.Set("Git-Protocol", "version=2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fetch: %v\n%s\n%s", resp.Status, data, hex.Dump(postbody))
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/x-git-upload-pack-result" {
		return nil, fmt.Errorf("fetch: invalid response Content-Type: %v", ct)
	}

	// Response is sequence of pkt-line packets.
	// It is plain text output (printed by git) until we find "packfile".
	// Then it switches to packets with a single prefix byte saying
	// what kind of data is in that packet:
	// 1 for pack file data, 2 for text output, 3 for errors.
	var data []byte
	pr := newPktLineReader(resp.Body)
	sawPackfile := false
	for {
		line, err := pr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("fetch: parsing response: %v", err)
		}
		if line == nil { // ignore delimiter
			continue
		}
		if !sawPackfile {
			// Discard response lines until we get to packfile start.
			if strings.TrimSuffix(string(line), "\n") == "packfile" {
				sawPackfile = true
			}
			continue
		}
		if len(line) == 0 || line[0] == 0 || line[0] > 3 {
			return nil, fmt.Errorf("fetch: malformed response: invalid sideband: %q", line)
		}
		switch line[0] {
		case 1:
			data = append(data, line[1:]...)
		case 2:
			fmt.Printf("%s\n", line[1:])
		case 3:
			return nil, fmt.Errorf("fetch: server error: %s", line[1:])
		}
	}

	if !bytes.HasPrefix(data, []byte("PACK")) {
		return nil, fmt.Errorf("fetch: malformed response: not packfile")
	}

	// Unpack pack file and return fs.FS for the commit we downloaded.
	var s store
	if err := unpack(&s, data); err != nil {
		return nil, fmt.Errorf("fetch: %v", err)
	}
	tfs, err := s.commit(h)
	if err != nil {
		return nil, fmt.Errorf("fetch: %v", err)
	}
	return tfs, nil
}
