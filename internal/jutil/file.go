package jutil

import "strings"

// Go-only: language-adapter — java.io.File, as an abstract pathname (the
// methods Closure uses). Paths use the Unix separator.

// File is java.io.File.
type File struct {
	path string
}

// NewFile is new File(String pathname): the pathname is normalized like
// UnixFileSystem.normalize (duplicate separators and a trailing separator
// are removed).
func NewFile(pathname string) *File {
	return &File{path: fileNormalize(pathname)}
}

func fileNormalize(pathname string) string {
	var sb strings.Builder
	prev := byte(0)
	for i := 0; i < len(pathname); i++ {
		c := pathname[i]
		if c == '/' && prev == '/' {
			continue
		}
		sb.WriteByte(c)
		prev = c
	}
	s := sb.String()
	if len(s) > 1 && strings.HasSuffix(s, "/") {
		s = s[:len(s)-1]
	}
	return s
}

// GetPath is File.getPath().
func (f *File) GetPath() string { return f.path }

// String is File.toString() (the path).
func (f *File) String() string { return f.path }

// ToPath is File.toPath().
func (f *File) ToPath() *Path { return Path_Of(f.path) }
