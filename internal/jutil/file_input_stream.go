package jutil

import (
	"io"
	"os"
)

// Go-only: language-adapter — java.io.FileInputStream as a Go
// io.ReadCloser over an *os.File.

// NewFileInputStream is new FileInputStream(String name): it throws
// FileNotFoundException("<name> (<reason>)") when the file does not exist,
// is a directory, or cannot be opened for reading.
func NewFileInputStream(name string) io.ReadCloser {
	f, err := os.Open(name)
	if err == nil {
		if st, serr := f.Stat(); serr == nil && st.IsDir() {
			f.Close()
			panic(NewFileNotFoundException(name + " (Is a directory)"))
		}
		return f
	}
	panic(FileNotFoundExceptionFromError(err, name))
}
