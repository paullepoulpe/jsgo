package main

import (
	"bytes"
	"compress/gzip"
	"io"
)

// comparableOutput returns what is compared of an output file. A gzip file
// (--typed_ast_output_file, a saved compiler state, any .gz output) is
// compared by its 10-byte member header and its decompressed content, not
// by its compressed bytes: jsgo compresses with Go's compress/flate and
// Java with zlib, so the same content can compress differently. This is an
// intentional divergence (tools/oracle/README.md, "Compressed output").
// Anything else, and a gzip file that does not decompress, is compared as
// is.
func comparableOutput(b []byte) []byte {
	if len(b) < 10 || b[0] != 0x1f || b[1] != 0x8b {
		return b
	}
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return b
	}
	content, err := io.ReadAll(r)
	if err != nil {
		return b
	}
	return append(append([]byte{}, b[:10]...), content...)
}

// sameOutput reports whether two output files are the same, gzip files
// compared decompressed (comparableOutput).
func sameOutput(a, b []byte) bool {
	return bytes.Equal(comparableOutput(a), comparableOutput(b))
}
