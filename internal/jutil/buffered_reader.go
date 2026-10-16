package jutil

import (
	"io"
	"strings"
)

// Go-only: language-adapter — java.io.BufferedReader over a StringReader or
// over an InputStreamReader(stream, UTF_8). The whole input is decoded on
// the first read (Java decodes incrementally; the characters are the same).

// BufferedReader is java.io.BufferedReader.
type BufferedReader struct {
	src    io.Reader
	text   string
	loaded bool
	pos    int
}

// NewBufferedReaderFromString is new BufferedReader(new StringReader(s)).
func NewBufferedReaderFromString(s string) *BufferedReader {
	return &BufferedReader{text: s, loaded: true}
}

// NewBufferedReaderUTF8 is new BufferedReader(new InputStreamReader(stream,
// UTF_8)): malformed input is replaced like new String(bytes, UTF_8).
func NewBufferedReaderUTF8(stream io.Reader) *BufferedReader {
	return &BufferedReader{src: stream}
}

func (r *BufferedReader) load() {
	if r.loaded {
		return
	}
	b, err := io.ReadAll(r.src)
	if err != nil {
		panic(NewIOException(err.Error()))
	}
	r.text = NewStringFromUTF8Bytes(b)
	r.loaded = true
}

// ReadLine is BufferedReader.readLine(): a line ends at '\n', '\r' or
// "\r\n" (not included); ok is false (Java null) at end of input.
func (r *BufferedReader) ReadLine() (line string, ok bool) {
	r.load()
	if r.pos >= len(r.text) {
		return "", false
	}
	rest := r.text[r.pos:]
	i := strings.IndexAny(rest, "\r\n")
	if i < 0 {
		r.pos = len(r.text)
		return rest, true
	}
	r.pos += i + 1
	if rest[i] == '\r' && i+1 < len(rest) && rest[i+1] == '\n' {
		r.pos++
	}
	return rest[:i], true
}

// Close is BufferedReader.close().
func (r *BufferedReader) Close() {}
