package jutil

import (
	"io"
	"os"
)

// Go-only: language-adapter — java.io.FileOutputStream as a Go
// io.WriteCloser over an *os.File.

// NewFileOutputStream is new FileOutputStream(String name): it creates or
// truncates the file, and throws FileNotFoundException("<name> (<reason>)")
// when the file cannot be opened for writing. A write error is thrown as
// IOException by the writers that wrap it.
func NewFileOutputStream(name string) io.WriteCloser {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		panic(FileNotFoundExceptionFromError(err, name))
	}
	return f
}
