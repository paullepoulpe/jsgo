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

// NextLineFromString scans the next line starting at pos in text, ending at
// '\n', '\r' or "\r\n" (not included). Returns ok == false when pos >= len(text).
func NextLineFromString(text string, pos int) (line string, nextPos int, ok bool) {
	if pos >= len(text) {
		return "", pos, false
	}
	rest := text[pos:]
	nl := strings.IndexByte(rest, '\n')
	if nl >= 0 {
		cr := strings.IndexByte(rest[:nl], '\r')
		if cr < 0 {
			return rest[:nl], pos + nl + 1, true
		}
		if cr == nl-1 {
			return rest[:cr], pos + nl + 1, true
		}
		return rest[:cr], pos + cr + 1, true
	}
	cr := strings.IndexByte(rest, '\r')
	if cr < 0 {
		return rest, len(text), true
	}
	return rest[:cr], pos + cr + 1, true
}

// ReadLine is BufferedReader.readLine(): a line ends at '\n', '\r' or
// "\r\n" (not included); ok is false (Java null) at end of input.
func (r *BufferedReader) ReadLine() (line string, ok bool) {
	r.load()
	line, r.pos, ok = NextLineFromString(r.text, r.pos)
	return line, ok
}

// Close is BufferedReader.close().
func (r *BufferedReader) Close() {}
