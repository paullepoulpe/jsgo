package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/tools/oracle/corpus"
)

// fileText returns the text that the content checks (GPL mentions,
// excluded corpus paths) read in the tracked file name with contents b:
//   - a gzip file (by its magic number; a ".gz" name without it is an
//     error): its decompressed contents, every member. Matching the raw
//     bytes would find chance matches in the compressed data and miss the
//     text;
//   - a zip or jar archive (by magic number or name): the text of each of
//     its files, read the same way (so a gzip member is decompressed and a
//     binary member skipped), one after the other. A file that is not a
//     valid archive is an error;
//   - a file in another compressed or archive format (by magic number or
//     name: bzip2, xz, zstd, 7z): an error, since the checks cannot read
//     it. None is tracked; decompress it, or teach fileText;
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
	for _, f := range archiveFormats {
		if f.is(b) || ext == f.ext {
			text, err := archiveText(name, f.name, b)
			return text, false, err
		}
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

// archiveText is fileText for a zip archive (kind "zip" or "jar") named
// name: the texts of its files, each followed by a newline.
func archiveText(name, kind string, b []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, fmt.Errorf("%s: %s file: %v", name, kind, err)
	}
	var out []byte
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("%s: %s: %v", name, f.Name, err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("%s: %s: %v", name, f.Name, err)
		}
		text, skip, err := fileText(name+"!/"+f.Name, data)
		if err != nil {
			return nil, err
		}
		if skip {
			continue
		}
		out = append(out, text...)
		out = append(out, '\n')
	}
	return out, nil
}

// mentionsGPL reports whether text mentions the GPL. The name of the list
// of excluded corpus files (tools/oracle/corpus) does not count.
func mentionsGPL(text []byte) bool {
	// Most files mention neither word; a case-folded search for the
	// literal part of each alternative of gplPattern rules them out fast.
	lower := bytes.ToLower(text)
	if !bytes.Contains(lower, []byte("g"+"pl")) && !bytes.Contains(lower, []byte("public licen")) {
		return false
	}
	if name := []byte(corpus.ListName); bytes.Contains(text, name) {
		text = bytes.ReplaceAll(text, name, []byte(" "))
	}
	return gplPattern.Match(text)
}

// excludedIndex finds the excluded corpus files that a text mentions: it
// holds the text that makes a tracked file mention each one, its key
// (corpus.Key), grouped by the key's directory. The excluded files sit in
// few directories, so one search per directory rules out most texts.
type excludedIndex struct {
	dirs  [][]byte   // distinct key directories, with a trailing "/"
	files [][]string // per directory, the excluded files
	keys  map[string][]byte
}

func newExcludedIndex(excluded []string) *excludedIndex {
	x := &excludedIndex{keys: map[string][]byte{}}
	at := map[string]int{}
	for _, e := range excluded {
		k := corpus.Key(e)
		x.keys[e] = []byte(k)
		dir := k[:strings.LastIndex(k, "/")+1]
		if dir == "" {
			dir = k
		}
		i, ok := at[dir]
		if !ok {
			i = len(x.dirs)
			at[dir] = i
			x.dirs = append(x.dirs, []byte(dir))
			x.files = append(x.files, nil)
		}
		x.files[i] = append(x.files[i], e)
	}
	return x
}

// mentioned returns the excluded corpus files that text mentions.
func (x *excludedIndex) mentioned(text []byte) []string {
	var found []string
	for i, dir := range x.dirs {
		if !bytes.Contains(text, dir) {
			continue
		}
		for _, e := range x.files[i] {
			if bytes.Contains(text, x.keys[e]) {
				found = append(found, e)
			}
		}
	}
	return found
}

var gzipMagic = []byte{0x1f, 0x8b}

// fileFormat is a format fileText recognises: its name extension, and a
// test of its first bytes.
type fileFormat struct {
	name, ext string
	is        func(b []byte) bool
}

// archiveFormats are the zip-based archives fileText reads (archiveText).
var archiveFormats = []fileFormat{
	{"zip", ".zip", prefix("PK\x03\x04")},
	{"jar", ".jar", prefix("PK\x03\x04")},
}

// compressedFormats are the formats fileText recognises but cannot read.
var compressedFormats = []fileFormat{
	{"bzip2", ".bz2", func(b []byte) bool { return len(b) > 3 && prefix("BZh")(b) && '1' <= b[3] && b[3] <= '9' }},
	{"xz", ".xz", prefix("\xfd7zXZ\x00")},
	{"zstd", ".zst", prefix("\x28\xb5\x2f\xfd")},
	{"7z", ".7z", prefix("7z\xbc\xaf\x27\x1c")},
}

func prefix(magic string) func([]byte) bool {
	return func(b []byte) bool { return bytes.HasPrefix(b, []byte(magic)) }
}
