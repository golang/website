// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitfs

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

// unpack parses data, which is a Git pack-formatted archive,
// writing every object it contains to the store s.
//
// See https://git-scm.com/docs/pack-format for format documentation.
func unpack(s *store, data []byte) error {
	// If the store is empty, pre-allocate the length of data.
	// This should be about the right order of magnitude for the eventual data,
	// avoiding many growing steps during append.
	if len(s.data) == 0 {
		s.data = make([]byte, 0, len(data))
	}

	// Pack data starts with 12-byte header: "PACK" version[4] nobj[4].
	if len(data) < 12+20 {
		return fmt.Errorf("malformed git pack: too short")
	}
	hdr := data[:12]
	vers := binary.BigEndian.Uint32(hdr[4:8])
	nobj := binary.BigEndian.Uint32(hdr[8:12])
	if string(hdr[:4]) != "PACK" || vers != 2 && vers != 3 || len(data) < 12+20 || int64(nobj) >= int64(len(data)) {
		return fmt.Errorf("malformed git pack")
	}
	if vers == 3 {
		return fmt.Errorf("cannot read git pack v3")
	}

	// Pack data ends with SHA1 of the entire pack.
	sum := sha1.Sum(data[:len(data)-20])
	if !bytes.Equal(sum[:], data[len(data)-20:]) {
		return fmt.Errorf("malformed git pack: bad checksum")
	}

	// Object data is everything between hdr and ending SHA1.
	// Unpack every object into the store.
	objs := data[12 : len(data)-20]
	off := 0
	for i := 0; i < int(nobj); i++ {
		_, _, _, encSize, err := unpackObject(s, objs, off)
		if err != nil {
			return fmt.Errorf("unpack: malformed git pack: %v", err)
		}
		off += encSize
	}
	if off != len(objs) {
		return fmt.Errorf("malformed git pack: junk after objects")
	}
	return nil
}

// unpackObject unpacks the object at objs[off:] and writes it to the store s.
// It returns the type, hash, and content of the object, as well as the encoded size,
// meaning the number of bytes at the start of objs[off:] that this record occupies.
func unpackObject(s *store, objs []byte, off int) (typ objType, h Hash, content []byte, encSize int, err error) {
	fail := func(err error) (objType, Hash, []byte, int, error) {
		return 0, Hash{}, nil, 0, err
	}
	if off < 0 || off >= len(objs) {
		return fail(fmt.Errorf("invalid object offset"))
	}

	// Object starts with varint-encoded type and length n.
	// (The length n is the length of the compressed data that follows,
	// not the length of the actual object.)
	u, size := binary.Uvarint(objs[off:])
	if size <= 0 {
		return fail(fmt.Errorf("invalid object: bad varint header"))
	}
	typ = objType((u >> 4) & 7)
	n := int(u&15 | u>>7<<4)

	// Git often stores objects that differ very little (different revs of a file).
	// It can save space by encoding one as "start with this other object and apply these diffs".
	// There are two ways to specify "this other object": an object ref (20-byte SHA1)
	// or as a relative offset to an earlier position in the objs slice.
	// For either of these, we need to fetch the other object's type and data (deltaTyp and deltaBase).
	// The Git docs call this the "deltified representation".
	var deltaTyp objType
	var deltaBase []byte
	switch typ {
	case objRefDelta:
		if len(objs)-(off+size) < 20 {
			return fail(fmt.Errorf("invalid object: bad delta ref"))
		}
		// Base block identified by SHA1 of an already unpacked hash.
		var h Hash
		copy(h[:], objs[off+size:])
		size += 20
		deltaTyp, deltaBase = s.object(h)
		if deltaTyp == 0 {
			return fail(fmt.Errorf("invalid object: unknown delta ref %v", h))
		}

	case objOfsDelta:
		i := off + size
		if len(objs)-i < 20 {
			return fail(fmt.Errorf("invalid object: too short"))
		}
		// Base block identified by relative offset to earlier position in objs,
		// using a varint-like but not-quite-varint encoding.
		// Look for "offset encoding:" in https://git-scm.com/docs/pack-format.
		d := int64(objs[i] & 0x7f)
		for objs[i]&0x80 != 0 {
			i++
			if i-(off+size) > 10 {
				return fail(fmt.Errorf("invalid object: malformed delta offset"))
			}
			d = d<<7 | int64(objs[i]&0x7f)
			d += 1 << 7
		}
		i++
		size = i - off

		// Re-unpack the object at the earlier offset to find its type and content.
		if d == 0 || d > int64(off) {
			return fail(fmt.Errorf("invalid object: bad delta offset"))
		}
		var err error
		deltaTyp, _, deltaBase, _, err = unpackObject(s, objs, off-int(d))
		if err != nil {
			return fail(fmt.Errorf("invalid object: bad delta offset"))
		}
	}

	// The main encoded data is a zlib-compressed stream.
	br := bytes.NewReader(objs[off+size:])
	zr, err := zlib.NewReader(br)
	if err != nil {
		return fail(fmt.Errorf("invalid object deflate: %v", err))
	}
	data, err := ioutil.ReadAll(zr)
	if err != nil {
		return fail(fmt.Errorf("invalid object: bad deflate: %v", err))
	}
	if len(data) != n {
		return fail(fmt.Errorf("invalid object: deflate size %d != %d", len(data), n))
	}
	encSize = len(objs[off:]) - br.Len()

	// If we fetched a base object above, the stream is an encoded delta.
	// Otherwise it is the raw data.
	switch typ {
	default:
		return fail(fmt.Errorf("invalid object: unknown object type"))
	case objCommit, objTree, objBlob, objTag:
		// ok
	case objRefDelta, objOfsDelta:
		// Actual object type is the type of the base object.
		typ = deltaTyp

		// Delta encoding starts with size of base object and size of new object.
		baseSize, s := binary.Uvarint(data)
		data = data[s:]
		if baseSize != uint64(len(deltaBase)) {
			return fail(fmt.Errorf("invalid object: mismatched delta src size"))
		}
		targSize, s := binary.Uvarint(data)
		data = data[s:]

		// Apply delta to base object, producing new object.
		targ := make([]byte, targSize)
		if err := applyDelta(targ, deltaBase, data); err != nil {
			return fail(fmt.Errorf("invalid object: %v", err))
		}
		data = targ
	}

	h, data = s.add(typ, data)
	return typ, h, data, encSize, nil
}

// applyDelta applies the delta encoding to src, producing dst,
// which has already been allocated to the expected final size.
// See https://git-scm.com/docs/pack-format#_deltified_representation for docs.
func applyDelta(dst, src, delta []byte) error {
	for len(delta) > 0 {
		// Command byte says what comes next.
		cmd := delta[0]
		delta = delta[1:]
		switch {
		case cmd == 0:
			// cmd == 0 is reserved.
			return fmt.Errorf("invalid delta cmd")

		case cmd&0x80 != 0:
			// Copy from base object, 4-byte offset, 3-byte size.
			// But any zero byte in the offset or size can be omitted.
			// The bottom 7 bits of cmd say which offset/size bytes are present.
			var off, size int64
			for i := uint(0); i < 4; i++ {
				if cmd&(1<<i) != 0 {
					off |= int64(delta[0]) << (8 * i)
					delta = delta[1:]
				}
			}
			for i := uint(0); i < 3; i++ {
				if cmd&(0x10<<i) != 0 {
					size |= int64(delta[0]) << (8 * i)
					delta = delta[1:]
				}
			}
			// Size 0 means size 0x10000 for some reason. (!)
			if size == 0 {
				size = 0x10000
			}
			copy(dst[:size], src[off:off+size])
			dst = dst[size:]

		default:
			// Up to 0x7F bytes of literal data, length in bottom 7 bits of cmd.
			n := int(cmd)
			copy(dst[:n], delta[:n])
			dst = dst[n:]
			delta = delta[n:]
		}
	}
	if len(dst) != 0 {
		return fmt.Errorf("delta encoding too short")
	}
	return nil
}
