package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"path"
	"strings"
	"unicode/utf8"
)

// fileText returns the text that the content checks (GPL mentions,
// excluded corpus paths) read in the tracked file name with contents b:
//   - a gzip file (by its magic number; a ".gz" name without it is an
//     error): its decompressed contents, every member. Matching the raw
//     bytes would find chance matches in the compressed data and miss the
//     text;
//   - a file in another compressed or archive format (by magic number or
//     name: zip and jar, bzip2, xz, zstd, 7z): an error, since the checks
//     cannot read it. None is tracked; decompress it, or teach fileText;
//   - any other binary file: skip is true, and the checks skip it. A file
//     is binary when it has a NUL byte in its first 8000 bytes (git's
//     rule) and is not valid UTF-8 (some JavaScript test inputs hold a
//     NUL character; they are text);
//   - anything else: b.
func fileText(name string, b []byte) (text []byte, skip bool, err error) {
	if bytes.HasPrefix(b, gzipMagic) {
		zr, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, false, fmt.Errorf("%s: %v", name, err)
		}
		text, err := io.ReadAll(zr)
		if err != nil {
			return nil, false, fmt.Errorf("%s: %v", name, err)
		}
		return text, false, nil
	}
	ext := strings.ToLower(path.Ext(name))
	if ext == ".gz" || ext == ".tgz" {
		return nil, false, fmt.Errorf("%s: named like a gzip file but is not one", name)
	}
	for _, f := range compressedFormats {
		if f.is(b) || ext == f.ext {
			return nil, false, fmt.Errorf("%s: %s file: the checks cannot read it (decompress it, or add the format to tools/licenses/content.go)", name, f.name)
		}
	}
	if bytes.IndexByte(b[:min(len(b), 8000)], 0) >= 0 && !utf8.Valid(b) {
		return nil, true, nil
	}
	return b, false, nil
}

// mentionsGPL reports whether text mentions the GPL.
func mentionsGPL(text []byte) bool {
	// Most files mention neither word; a case-folded search for the
	// literal part of each alternative of gplPattern rules them out fast.
	lower := bytes.ToLower(text)
	if !bytes.Contains(lower, []byte("g"+"pl")) && !bytes.Contains(lower, []byte("public licen")) {
		return false
	}
	return gplPattern.Match(text)
}

var gzipMagic = []byte{0x1f, 0x8b}

// compressedFormats are the formats fileText recognises but cannot read:
// their name extension, and a test of their first bytes.
var compressedFormats = []struct {
	name, ext string
	is        func(b []byte) bool
}{
	{"zip", ".zip", prefix("PK\x03\x04")},
	{"jar", ".jar", prefix("PK\x03\x04")},
	{"bzip2", ".bz2", func(b []byte) bool { return len(b) > 3 && prefix("BZh")(b) && '1' <= b[3] && b[3] <= '9' }},
	{"xz", ".xz", prefix("\xfd7zXZ\x00")},
	{"zstd", ".zst", prefix("\x28\xb5\x2f\xfd")},
	{"7z", ".7z", prefix("7z\xbc\xaf\x27\x1c")},
}

func prefix(magic string) func([]byte) bool {
	return func(b []byte) bool { return bytes.HasPrefix(b, []byte(magic)) }
}
