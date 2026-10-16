package jutil

import "os"

// Go-only: language-adapter — the reading methods of java.nio.file.Files
// (see files.go).

// Files_ReadString is Files.readString(Path, Charset): it panics with
// MalformedInputException on malformed input, and with an IOException
// (NoSuchFileException, ...) when the file cannot be read.
func Files_ReadString(path *Path, cs *Charset) string {
	b, err := os.ReadFile(path.String())
	if err != nil {
		panic(IOExceptionFromError(err, path.String()))
	}
	return cs.DecodeStrict(b)
}

// Files_NewBufferedReader is Files.newBufferedReader(Path, Charset): the
// file is opened now, and malformed input is reported while reading.
func Files_NewBufferedReader(path *Path, cs *Charset) JavaReader {
	f, err := os.Open(path.String())
	if err != nil {
		panic(IOExceptionFromError(err, path.String()))
	}
	r := NewInputStreamReader(f, cs)
	r.strict = true
	return r
}
