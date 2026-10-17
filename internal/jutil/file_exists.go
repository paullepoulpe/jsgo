package jutil

import "os"

// Go-only: language-adapter — java.io.File.exists(), written from its
// Javadoc: whether the file or directory that the abstract pathname denotes
// exists (a relative pathname is resolved against the current directory). A
// pathname the system cannot examine, such as the empty pathname, does not
// exist.

// Exists is File.exists().
func (f *File) Exists() bool {
	if f.path == "" {
		return false
	}
	_, err := os.Stat(f.path)
	return err == nil
}
