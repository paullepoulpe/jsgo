package jutil

import (
	"os"
	"unicode/utf8"
	"unsafe"
)

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
	if (cs == nil || cs == StandardCharsets_UTF_8) && utf8.Valid(b) {
		return unsafe.String(unsafe.SliceData(b), len(b))
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

// Files_ReadStringUTF8 is java.nio.file.Files.readString(Path), which
// decodes with UTF-8.
func Files_ReadStringUTF8(path *Path) string {
	return Files_ReadString(path, StandardCharsets_UTF_8)
}
