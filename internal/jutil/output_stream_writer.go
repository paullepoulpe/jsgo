package jutil

import (
	"io"
	"unicode/utf8"
)

// Go-only: language-adapter — java.io.OutputStreamWriter with the UTF-8
// charset (or a UTF-8 CharsetEncoder) and CodingErrorAction.REPLACE, over a
// Go io.WriteCloser. Observed behaviour (writers_test.go, against JDK 21's
// public API): a high surrogate that ends one write pairs with a low
// surrogate that starts the next, so a pair split between writes is encoded
// whole, and an unpaired surrogate is encoded as the replacement byte '?'.
// Bytes are written through immediately (Java buffers them, which only
// changes when they reach the stream).

// OutputStreamWriter is java.io.OutputStreamWriter.
type OutputStreamWriter struct {
	out                     io.WriteCloser
	pendingHighSurrogate    uint16
	hasPendingHighSurrogate bool
	closed                  bool
}

// NewOutputStreamWriter is new OutputStreamWriter(out, UTF_8).
func NewOutputStreamWriter(out io.WriteCloser) *OutputStreamWriter {
	if out == nil {
		panic(NewNullPointerException())
	}
	return &OutputStreamWriter{out: out}
}

func (w *OutputStreamWriter) checkOpen() {
	if w.closed {
		panic(NewIOException("Stream closed"))
	}
}

// Write is OutputStreamWriter.write(String).
func (w *OutputStreamWriter) Write(s string) {
	w.checkOpen()
	units := StringToUTF16(s)
	if w.hasPendingHighSurrogate {
		units = append([]uint16{w.pendingHighSurrogate}, units...)
		w.hasPendingHighSurrogate = false
	}
	var b []byte
	for i := 0; i < len(units); i++ {
		c := units[i]
		switch {
		case c >= 0xD800 && c <= 0xDBFF:
			if i+1 == len(units) {
				w.pendingHighSurrogate, w.hasPendingHighSurrogate = c, true
				continue
			}
			if lo := units[i+1]; lo >= 0xDC00 && lo <= 0xDFFF {
				b = utf8.AppendRune(b, (rune(c)-0xD800)<<10+(rune(lo)-0xDC00)+0x10000)
				i++
				continue
			}
			b = append(b, '?')
		case c >= 0xDC00 && c <= 0xDFFF:
			b = append(b, '?')
		default:
			b = utf8.AppendRune(b, rune(c))
		}
	}
	w.writeBytes(b)
}

func (w *OutputStreamWriter) writeBytes(b []byte) {
	if len(b) == 0 {
		return
	}
	if _, err := w.out.Write(b); err != nil {
		panic(NewIOException(err.Error()))
	}
}

// Flush is OutputStreamWriter.flush(); the bytes are already in the stream
// (see above).
func (w *OutputStreamWriter) Flush() {
	w.checkOpen()
}

// Close is OutputStreamWriter.close(): an unpaired high surrogate left
// over from the last write is encoded as '?', then the stream is closed.
// Closing a closed writer has no effect.
func (w *OutputStreamWriter) Close() {
	if w.closed {
		return
	}
	w.closed = true
	if w.hasPendingHighSurrogate {
		w.hasPendingHighSurrogate = false
		w.writeBytes([]byte{'?'})
	}
	if err := w.out.Close(); err != nil {
		panic(NewIOException(err.Error()))
	}
}
