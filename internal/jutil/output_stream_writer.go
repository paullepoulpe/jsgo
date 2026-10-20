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
	byteBuf                 []byte
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
	b := w.byteBuf[:0]
	if cap(b) < len(s) {
		b = make([]byte, 0, max(len(s), 1024))
	}
	for i := 0; i < len(s); {
		if !w.hasPendingHighSurrogate && s[i] < 0x80 {
			j := i + 1
			for j < len(s) && s[j] < 0x80 {
				j++
			}
			b = append(b, s[i:j]...)
			i = j
			continue
		}
		r, sz := decodeWTF8(s[i:])
		i += sz
		if r >= 0x10000 {
			if w.hasPendingHighSurrogate {
				b = append(b, '?')
				w.hasPendingHighSurrogate = false
			}
			b = utf8.AppendRune(b, r)
			continue
		}
		c := uint16(r)
		switch {
		case c >= 0xD800 && c <= 0xDBFF:
			if w.hasPendingHighSurrogate {
				b = append(b, '?')
			}
			w.pendingHighSurrogate, w.hasPendingHighSurrogate = c, true
		case c >= 0xDC00 && c <= 0xDFFF:
			if w.hasPendingHighSurrogate {
				b = utf8.AppendRune(b, (rune(w.pendingHighSurrogate)-0xD800)<<10+(rune(c)-0xDC00)+0x10000)
				w.hasPendingHighSurrogate = false
			} else {
				b = append(b, '?')
			}
		default:
			if w.hasPendingHighSurrogate {
				b = append(b, '?')
				w.hasPendingHighSurrogate = false
			}
			b = utf8.AppendRune(b, rune(c))
		}
	}
	w.byteBuf = b[:0]
	w.writeBytes(b)
}

func (w *OutputStreamWriter) writeUTF16(units []uint16) {
	w.checkOpen()
	b := w.byteBuf[:0]
	if cap(b) < len(units) {
		b = make([]byte, 0, max(len(units), 1024))
	}
	i := 0
	if w.hasPendingHighSurrogate && len(units) > 0 {
		hi := w.pendingHighSurrogate
		w.hasPendingHighSurrogate = false
		if lo := units[0]; lo >= 0xDC00 && lo <= 0xDFFF {
			b = utf8.AppendRune(b, (rune(hi)-0xD800)<<10+(rune(lo)-0xDC00)+0x10000)
			i = 1
		} else {
			b = append(b, '?')
		}
	}
	for ; i < len(units); i++ {
		c := units[i]
		if c < 0x80 {
			b = append(b, byte(c))
			continue
		}
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
	w.byteBuf = b[:0]
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
		b := append(w.byteBuf[:0], '?')
		w.byteBuf = b[:0]
		w.writeBytes(b)
	}
	w.byteBuf = nil
	if err := w.out.Close(); err != nil {
		panic(NewIOException(err.Error()))
	}
}
