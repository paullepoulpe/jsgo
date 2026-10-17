// Package corpus keeps the corpus files listed in
// tools/oracle/gpl-excluded-corpus.txt out of the Go side of the oracles.
//
// Those files carry licence notices that PORTING.md does not allow in this
// repository, so no golden, dump or list may be made from them. The Java
// dumps skip them (tools/oracle/java/CorpusExclusions.java); the Go tests
// that read a corpus file list or walk a corpus directory pass what they
// read through ReadList or Filter; and "go run ./tools/licenses -check"
// fails if a tracked file mentions one of them (see tools/oracle/README.md).
//
// A path names an excluded file when, with forward slashes, it equals the
// file's key or ends with "/" and the key. The key is the listed path
// inside its npm package: the part after "<name>@<version>/package/" (see
// Key). So the match holds whatever directory the corpus root is.
package corpus

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ListPath is the excluded-files list, relative to the repository root.
const ListPath = "tools/oracle/" + ListName

// ListName is the base name of the list.
const ListName = "gpl-excluded-corpus.txt"

var (
	loadOnce sync.Once
	loaded   []string
	loadErr  error
)

// Excluded returns the paths of the list, in its order. It finds the list
// by walking up from the current directory (a test runs in its package's
// directory, a tool at the repository root).
func Excluded() ([]string, error) {
	loadOnce.Do(func() {
		var b []byte
		b, loadErr = readList()
		if loadErr == nil {
			loaded = ParseList(b)
		}
	})
	return loaded, loadErr
}

func readList() ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for {
		b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(ListPath)))
		if err == nil {
			return b, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, fmt.Errorf("%s not found above the current directory", ListPath)
		}
		dir = parent
	}
}

// ParseList returns the paths of a list in the format of
// gpl-excluded-corpus.txt: one per line; blank lines and lines starting
// with "#" are ignored.
func ParseList(b []byte) []string {
	var paths []string
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		l := strings.TrimSpace(sc.Text())
		if l != "" && !strings.HasPrefix(l, "#") {
			paths = append(paths, l)
		}
	}
	return paths
}

// Key returns the part of an excluded path that names the file whatever
// the corpus root: the path inside its npm package, or the whole path when
// it has no "<name>@<version>/package/" prefix.
func Key(excluded string) string {
	parts := strings.Split(excluded, "/")
	for i := 0; i+1 < len(parts); i++ {
		if strings.Contains(parts[i], "@") && parts[i+1] == "package" && i+2 < len(parts) {
			return strings.Join(parts[i+2:], "/")
		}
	}
	return excluded
}

// Match reports whether path names one of the excluded files (see the
// package comment).
func Match(excluded []string, path string) bool {
	p := filepath.ToSlash(path)
	for _, e := range excluded {
		k := Key(e)
		if p == k || strings.HasSuffix(p, "/"+k) {
			return true
		}
	}
	return false
}

// IsExcluded reports whether path names a file of the list. It panics if
// the list cannot be read: without it the guard cannot work, and an
// excluded file must never get through.
func IsExcluded(path string) bool {
	return Match(mustExcluded(), path)
}

func mustExcluded() []string {
	ex, err := Excluded()
	if err != nil {
		panic("tools/oracle/corpus: " + err.Error())
	}
	return ex
}

// Filter returns paths without the excluded files, in order. It panics if
// the list cannot be read (see IsExcluded).
func Filter(paths []string) []string {
	ex := mustExcluded()
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if !Match(ex, p) {
			out = append(out, p)
		}
	}
	return out
}

// ReadList reads a corpus file list (one path per line, as the oracle
// dumps' @<listfile> and --list=<listfile>) and returns its non-empty
// lines without the excluded files.
func ReadList(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var files []string
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(nil, 1<<20)
	for sc.Scan() {
		if l := sc.Text(); l != "" {
			files = append(files, l)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("%s: %v", path, err)
	}
	ex, err := Excluded()
	if err != nil {
		return nil, err
	}
	out := files[:0]
	for _, f := range files {
		if !Match(ex, f) {
			out = append(out, f)
		}
	}
	return out, nil
}
