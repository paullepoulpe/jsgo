package jutil

import "io"

// Go-only: language-adapter — java.io.OutputStreamWriter(OutputStream,
// Charset), encoding with CodingErrorAction.REPLACE as the Javadoc
// specifies for this constructor: an unmappable char (a surrogate pair
// counts as one) or an unpaired surrogate becomes the encoder's
// replacement (charsetEncoder). UTF-8 is jutil.OutputStreamWriter.
// Observed behaviour (black-box, cli_vectors_test.go and
// charset_vectors_test.go): a high surrogate that ends one write() pairs with
// a low surrogate that starts the next; close() encodes one left over from
// the last write as unpaired, then ends the output (ISO-2022-JP returns to
// ASCII there); a charset that writes a byte-order mark writes it before
// its first char, so an empty output stays empty.

// NewOutputStreamWriterWithCharset is new OutputStreamWriter(out, cs).
func NewOutputStreamWriterWithCharset(out io.WriteCloser, cs *Charset) Writer {
	if out == nil || cs == nil {
		panic(NewNullPointerException())
	}
	if cs == StandardCharsets_UTF_8 {
		return NewOutputStreamWriter(out)
	}
	return &charsetOutputStreamWriter{out: out, enc: cs.impl().newEncoder()}
}

type charsetOutputStreamWriter struct {
	out                     io.WriteCloser
	enc                     charsetEncoder
	pendingHighSurrogate    uint16
	hasPendingHighSurrogate bool
	closed                  bool
}

func (w *charsetOutputStreamWriter) checkOpen() {
	if w.closed {
		panic(NewIOException("Stream closed"))
	}
}

// Write is OutputStreamWriter.write(String).
func (w *charsetOutputStreamWriter) Write(s string) {
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
		case isHighSurrogateUnit(c):
			if i+1 == len(units) {
				w.pendingHighSurrogate, w.hasPendingHighSurrogate = c, true
				continue
			}
			if lo := units[i+1]; isLowSurrogateUnit(lo) {
				b = w.enc.encodeRune(b, 0x10000+(rune(c)-0xD800)<<10+(rune(lo)-0xDC00))
				i++
				continue
			}
			b = w.enc.encodeMalformed(b)
		case isLowSurrogateUnit(c):
			b = w.enc.encodeMalformed(b)
		default:
			b = w.enc.encodeRune(b, rune(c))
		}
	}
	w.emit(b)
}

func (w *charsetOutputStreamWriter) emit(b []byte) {
	if len(b) == 0 {
		return
	}
	if _, err := w.out.Write(b); err != nil {
		panic(NewIOException(err.Error()))
	}
}

// Flush is OutputStreamWriter.flush(); bytes are written through
// immediately.
func (w *charsetOutputStreamWriter) Flush() {
	w.checkOpen()
}

// Close is OutputStreamWriter.close(): a high surrogate left over from the
// last write is encoded as unpaired, the output is ended, then the stream
// is closed. Closing a closed writer has no effect.
func (w *charsetOutputStreamWriter) Close() {
	if w.closed {
		return
	}
	w.closed = true
	var b []byte
	if w.hasPendingHighSurrogate {
		w.hasPendingHighSurrogate = false
		b = w.enc.encodeMalformed(b)
	}
	w.emit(w.enc.finish(b))
	if err := w.out.Close(); err != nil {
		panic(NewIOException(err.Error()))
	}
}
