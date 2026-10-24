package jutil

import (
	"os"
	"unicode/utf8"
)

// Go-only: language-adapter — the reading methods of java.nio.file.Files
// (see files.go).

// Files_ReadString is Files.readString(Path, Charset): it panics with
// MalformedInputException on malformed input, and with an IOException
// (NoSuchFileException, ...) when the file cannot be read.
func Files_ReadString(path *Path, cs *Charset) string {
	return Files_ReadStringRaw(path.String(), cs)
}

// Files_ReadStringRaw is Files_ReadString for a pre-normalized path string.
func Files_ReadStringRaw(pathStr string, cs *Charset) string {
	b, err := os.ReadFile(pathStr)
	if err != nil {
		panic(IOExceptionFromError(err, pathStr))
	}
	if (cs == nil || cs == StandardCharsets_UTF_8) && utf8.Valid(b) {
		return string(b)
	}
	return cs.DecodeStrict(b)
}

// Files_NewBufferedReader is Files.newBufferedReader(Path, Charset): the
// file is opened now, and malformed input is reported while reading.
func Files_NewBufferedReader(path *Path, cs *Charset) JavaReader {
	return Files_NewBufferedReaderRaw(path.String(), cs)
}

// Files_NewBufferedReaderRaw is Files_NewBufferedReader for a pre-normalized path string.
func Files_NewBufferedReaderRaw(pathStr string, cs *Charset) JavaReader {
	f, err := os.Open(pathStr)
	if err != nil {
		panic(IOExceptionFromError(err, pathStr))
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
