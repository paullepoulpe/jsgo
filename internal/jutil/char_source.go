package jutil

import "os"

// Go-only: language-adapter — Guava com.google.common.io.CharSource, as
// returned by Files.asCharSource(File, Charset).

// CharSource is Guava's CharSource over a file.
type CharSource struct {
	file    *File
	charset *Charset
}

// Files_AsCharSource is Guava Files.asCharSource(File, Charset).
func Files_AsCharSource(file *File, charset *Charset) *CharSource {
	return &CharSource{file: CheckNotNull(file), charset: CheckNotNull(charset)}
}

// OpenStream is CharSource.openStream(): a new InputStreamReader over a new
// FileInputStream for the file. It panics with FileNotFoundException when
// the file cannot be opened.
func (c *CharSource) OpenStream() JavaReader {
	f, err := os.Open(c.file.GetPath())
	if err == nil {
		var st os.FileInfo
		if st, err = f.Stat(); err == nil && st.IsDir() {
			f.Close()
			panic(NewFileNotFoundException(c.file.GetPath() + " (Is a directory)"))
		}
	}
	if err != nil {
		panic(FileNotFoundExceptionFromError(err, c.file.GetPath()))
	}
	return NewInputStreamReader(f, c.charset)
}
