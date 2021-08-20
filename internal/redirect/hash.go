// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file provides a compact encoding of
// a map of Mercurial hashes to Git hashes.

package redirect

import (
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// hashMap is an encoded map of Mercurial hashes to Git hashes.
// It is a sequence of 8-byte entries, each of which is two little-endian uint32s
// giving a Mercurial, Git hash prefix pair.
// The map is sorted by Mercurial hash to allow binary search.
type hashMap []byte

// Lookup finds an hgHash in the map that matches the given prefix, and returns
// its corresponding gitHash. The prefix must be at least 8 characters long.
func (m hashMap) Lookup(s string) gitHash {
	if m == nil {
		return 0
	}
	hg, err := hgHashFromString(s)
	if err != nil {
		return 0
	}
	var git gitHash
	sort.Search(len(m)/8, func(i int) bool {
		entry := m[i*8 : (i+1)*8]
		v := hgHash(binary.LittleEndian.Uint32(entry[:4]))
		if v == hg {
			git = gitHash(binary.LittleEndian.Uint32(entry[4:]))
		}
		return v >= hg
	})
	return git
}

// hgHash represents the lower (leftmost) 32 bits of a Mercurial hash.
type hgHash uint32

func (h hgHash) String() string {
	return intToHash(int64(h))
}

func hgHashFromString(s string) (hgHash, error) {
	if len(s) < 8 {
		return 0, fmt.Errorf("string too small: len(s) = %d", len(s))
	}
	hash := s[:8]
	i, err := strconv.ParseInt(hash, 16, 64)
	if err != nil {
		return 0, err
	}
	return hgHash(i), nil
}

// gitHash represents the leftmost 28 bits of a Git hash in its upper 28 bits,
// and it encodes hash's repository in the lower 4  bits.
type gitHash uint32

func (h gitHash) Hash() string {
	return intToHash(int64(h))[:7]
}

func (h gitHash) Repo() string {
	return repo(h & 0xF).String()
}

func intToHash(i int64) string {
	s := strconv.FormatInt(i, 16)
	if len(s) < 8 {
		s = strings.Repeat("0", 8-len(s)) + s
	}
	return s
}

// repo represents a Go Git repository.
type repo byte

const (
	repoGo repo = iota
	repoBlog
	repoCrypto
	repoExp
	repoImage
	repoMobile
	repoNet
	repoSys
	repoTalks
	repoText
	repoTools
)

func (r repo) String() string {
	return map[repo]string{
		repoGo:     "go",
		repoBlog:   "blog",
		repoCrypto: "crypto",
		repoExp:    "exp",
		repoImage:  "image",
		repoMobile: "mobile",
		repoNet:    "net",
		repoSys:    "sys",
		repoTalks:  "talks",
		repoText:   "text",
		repoTools:  "tools",
	}[r]
}
