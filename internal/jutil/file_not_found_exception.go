package jutil

import (
	"errors"
	"io/fs"
	"strings"
	"syscall"
)

// FileNotFoundException is java.io.FileNotFoundException.
type FileNotFoundException struct{ IOException }

func (*FileNotFoundException) fileNotFoundException() {}

// IFileNotFoundException is the catch marker for FileNotFoundException.
type IFileNotFoundException interface {
	IIOException
	fileNotFoundException()
}

// NewFileNotFoundException is new FileNotFoundException(String).
func NewFileNotFoundException(message string) *FileNotFoundException {
	e := &FileNotFoundException{}
	InitThrowable(e, "java.io.FileNotFoundException", &message, nil)
	return e
}

// FileNotFoundExceptionFromError is the exception java.io's
// FileInputStream/FileOutputStream constructors throw when opening path
// fails: FileNotFoundException("<path> (<reason>)"), the reason being the
// C library's strerror text (e.g. "No such file or directory", "Is a
// directory").
func FileNotFoundExceptionFromError(err error, path string) *FileNotFoundException {
	reason := err.Error()
	var errno syscall.Errno
	if errors.As(err, &errno) {
		reason = errno.Error()
	} else if pe := (*fs.PathError)(nil); errors.As(err, &pe) {
		reason = pe.Err.Error()
	}
	// Go's errno strings are strerror's, lowercased in the first letter.
	if reason != "" {
		reason = strings.ToUpper(reason[:1]) + reason[1:]
	}
	return NewFileNotFoundException(path + " (" + reason + ")")
}
