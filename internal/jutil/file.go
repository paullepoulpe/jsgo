package jutil

import "strings"

// Go-only: language-adapter — java.io.File, as an abstract pathname (the
// methods Closure uses). Paths use the Unix separator.

// File is java.io.File.
type File struct {
	path string
}

// NewFile is new File(String pathname). The pathname is normalized as
// observed through File.getPath() (file_test.go): runs of '/' collapse to
// one and a trailing '/' is dropped, except in the root path "/".
func NewFile(pathname string) *File {
	return &File{path: fileNormalize(pathname)}
}

// fileNormalize drops the empty names of a '/'-separated pathname, keeping
// a leading '/'.
func fileNormalize(pathname string) string {
	if !strings.Contains(pathname, "//") && (len(pathname) <= 1 || !strings.HasSuffix(pathname, "/")) {
		return pathname
	}
	var names []string
	for _, name := range strings.Split(pathname, "/") {
		if name != "" {
			names = append(names, name)
		}
	}
	s := strings.Join(names, "/")
	if strings.HasPrefix(pathname, "/") {
		s = "/" + s
	}
	return s
}

// GetPath is File.getPath().
func (f *File) GetPath() string { return f.path }

// String is File.toString() (the path).
func (f *File) String() string { return f.path }

// ToPath is File.toPath().
func (f *File) ToPath() *Path { return Path_Of(f.path) }
