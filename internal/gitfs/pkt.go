// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitfs

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// A pktLineReader reads Git pkt-line-formatted packets.
//
// Each n-byte packet is preceded by a 4-digit hexadecimal length
// encoding n+4 (the length counts its own bytes), like "0006a\n" for "a\n".
//
// A packet starting with 0000 is a so-called flush packet.
// A packet starting with 0001 is a delimiting marker,
// which usually marks the end of a sequence in the stream.
//
// See https://git-scm.com/docs/protocol-common#_pkt_line_format
// for the official documentation, although it fails to mention the 0001 packets.
type pktLineReader struct {
	b    *bufio.Reader
	size [4]byte
}

// newPktLineReader returns a new pktLineReader reading from r.
func newPktLineReader(r io.Reader) *pktLineReader {
	return &pktLineReader{b: bufio.NewReader(r)}
}

// Next returns the payload of the next packet from the stream.
// If the next packet is a flush packet (length 0000), Next returns nil, io.EOF.
// If the next packet is a delimiter packet (length 0001), Next returns nil, nil.
// If the data stream has ended, Next returns nil, io.ErrUnexpectedEOF.
func (r *pktLineReader) Next() ([]byte, error) {
	_, err := io.ReadFull(r.b, r.size[:])
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	n, err := strconv.ParseUint(string(r.size[:]), 16, 0)
	if err != nil || n == 2 || n == 3 {
		return nil, fmt.Errorf("malformed pkt-line")
	}
	if n == 1 {
		return nil, nil // delimiter
	}
	if n == 0 {
		return nil, io.EOF
	}
	buf := make([]byte, n-4)
	_, err = io.ReadFull(r.b, buf)
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	return buf, nil
}

// Lines reads packets from r until a flush packet.
// It returns a string for each packet, with any trailing newline trimmed.
func (r *pktLineReader) Lines() ([]string, error) {
	var lines []string
	for {
		line, err := r.Next()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return lines, err
		}
		lines = append(lines, strings.TrimSuffix(string(line), "\n"))
	}
}

// A pktLineWriter writes Git pkt-line-formatted packets.
// See pktLineReader for a description of the packet format.
type pktLineWriter struct {
	b    *bufio.Writer
	size [4]byte
}

// newPktLineWriter returns a new pktLineWriter writing to w.
func newPktLineWriter(w io.Writer) *pktLineWriter {
	return &pktLineWriter{b: bufio.NewWriter(w)}
}

// writeSize writes a four-digit hexadecimal length packet for n.
// Typically n is len(data)+4.
func (w *pktLineWriter) writeSize(n int) {
	hex := "0123456789abcdef"
	w.size[0] = hex[n>>12]
	w.size[1] = hex[(n>>8)&0xf]
	w.size[2] = hex[(n>>4)&0xf]
	w.size[3] = hex[(n>>0)&0xf]
	w.b.Write(w.size[:])
}

// Write writes b as a single packet.
func (w *pktLineWriter) Write(b []byte) (int, error) {
	n := len(b)
	if n+4 > 0xffff {
		return 0, fmt.Errorf("write too large")
	}
	w.writeSize(n + 4)
	w.b.Write(b)
	return n, nil
}

// WriteString writes s as a single packet.
func (w *pktLineWriter) WriteString(s string) (int, error) {
	n := len(s)
	if n+4 > 0xffff {
		return 0, fmt.Errorf("write too large")
	}
	w.writeSize(n + 4)
	w.b.WriteString(s)
	return n, nil
}

// Close writes a terminating flush packet
// and flushes buffered data to the underlying writer.
func (w *pktLineWriter) Close() error {
	w.b.WriteString("0000")
	w.b.Flush()
	return nil
}

// Delim writes a delimiter packet.
func (w *pktLineWriter) Delim() {
	w.b.WriteString("0001")
}
