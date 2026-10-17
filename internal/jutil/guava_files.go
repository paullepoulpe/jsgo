package jutil

import "os"

// Files_CreateTempDir is Guava's (deprecated) Files.createTempDir(): a new
// directory under java.io.tmpdir, readable only by the user.
func Files_CreateTempDir() *File {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(NewIllegalStateExceptionMsg("Failed to create directory: " + err.Error()))
	}
	return NewFile(dir)
}
