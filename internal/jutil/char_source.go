package jutil

import "os"

// Go-only: language-adapter — Guava com.google.common.io.CharSource, as
// returned by Files.asCharSource(File, Charset) and CharSource.wrap(CharSequence)
// (char_source_wrap.go).

// CharSource is Guava's CharSource over a file, or over a string.
type CharSource struct {
	file    *File
	charset *Charset
	// wrapped is the string of a CharSource.wrap source (file is then nil).
	wrapped *string
}

// Files_AsCharSource is Guava Files.asCharSource(File, Charset).
func Files_AsCharSource(file *File, charset *Charset) *CharSource {
	return &CharSource{file: CheckNotNull(file), charset: CheckNotNull(charset)}
}

// OpenStream is CharSource.openStream(): a new InputStreamReader over a new
// FileInputStream for the file. It panics with FileNotFoundException when
// the file cannot be opened.
func (c *CharSource) OpenStream() JavaReader {
	if c.wrapped != nil {
		return NewStringReader(*c.wrapped)
	}
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

// Read is CharSource.read(): the whole content of the source as a string
// (for a file, decoded with the source's charset; it panics with an
// IOException (FileNotFoundException) when the file cannot be read).
func (c *CharSource) Read() string {
	r := c.OpenStream()
	defer r.Close()
	var sb StringBuilder
	buf := make([]uint16, 8192)
	for {
		n := r.Read(buf, 0, len(buf))
		if n < 0 {
			break
		}
		for _, u := range buf[:n] {
			sb.AppendChar(u)
		}
	}
	return sb.String()
}
