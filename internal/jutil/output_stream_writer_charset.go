package jutil

import "io"

// Go-only: language-adapter — java.io.OutputStreamWriter(OutputStream,
// Charset) for the standard charsets (charset.go), encoding with
// CodingErrorAction.REPLACE as the Javadoc specifies for this constructor.
// Verified black-box against the JDK (output_stream_writer_charset_test.go):
//   - US-ASCII and ISO-8859-1 replace each unmappable character (a
//     surrogate pair counts as one) and each unpaired surrogate with '?';
//   - UTF-16BE and UTF-16LE replace an unpaired surrogate with U+FFFD;
//   - UTF-16 is big-endian and writes the byte-order mark U+FEFF before
//     its first bytes (an empty output stays empty);
//   - UTF-8 is jutil.OutputStreamWriter.
// Observed behaviour (same tests): a high surrogate that ends one write()
// pairs with a low surrogate that starts the next, and close() encodes one
// left over from the last write as unpaired.

// NewOutputStreamWriterWithCharset is new OutputStreamWriter(out, cs).
func NewOutputStreamWriterWithCharset(out io.WriteCloser, cs *Charset) Writer {
	if out == nil || cs == nil {
		panic(NewNullPointerException())
	}
	if cs == StandardCharsets_UTF_8 {
		return NewOutputStreamWriter(out)
	}
	return &charsetOutputStreamWriter{out: out, cs: cs}
}

type charsetOutputStreamWriter struct {
	out                     io.WriteCloser
	cs                      *Charset
	pendingHighSurrogate    uint16
	hasPendingHighSurrogate bool
	wroteBOM                bool
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
		if c >= 0xD800 && c <= 0xDBFF {
			if i+1 == len(units) {
				w.pendingHighSurrogate, w.hasPendingHighSurrogate = c, true
				continue
			}
			if lo := units[i+1]; lo >= 0xDC00 && lo <= 0xDFFF {
				b = w.encodePair(b, c, lo)
				i++
				continue
			}
			b = w.encodeMalformed(b)
			continue
		}
		if c >= 0xDC00 && c <= 0xDFFF {
			b = w.encodeMalformed(b)
			continue
		}
		b = w.encodeChar(b, c)
	}
	w.emit(b)
}

func (w *charsetOutputStreamWriter) encodeChar(b []byte, c uint16) []byte {
	switch w.cs {
	case StandardCharsets_US_ASCII, StandardCharsets_ISO_8859_1:
		if c <= w.cs.maxChar {
			return append(b, byte(c))
		}
		return append(b, '?')
	case StandardCharsets_UTF_16LE:
		return append(b, byte(c), byte(c>>8))
	default: // UTF-16, UTF-16BE
		return append(b, byte(c>>8), byte(c))
	}
}

func (w *charsetOutputStreamWriter) encodePair(b []byte, hi, lo uint16) []byte {
	switch w.cs {
	case StandardCharsets_US_ASCII, StandardCharsets_ISO_8859_1:
		return append(b, '?')
	default:
		return w.encodeChar(w.encodeChar(b, hi), lo)
	}
}

func (w *charsetOutputStreamWriter) encodeMalformed(b []byte) []byte {
	switch w.cs {
	case StandardCharsets_US_ASCII, StandardCharsets_ISO_8859_1:
		return append(b, '?')
	default:
		return w.encodeChar(b, 0xFFFD)
	}
}

func (w *charsetOutputStreamWriter) emit(b []byte) {
	if len(b) == 0 {
		return
	}
	if w.cs == StandardCharsets_UTF_16 && !w.wroteBOM {
		w.wroteBOM = true
		b = append([]byte{0xFE, 0xFF}, b...)
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
// last write is encoded as unpaired, then the stream is closed. Closing a
// closed writer has no effect.
func (w *charsetOutputStreamWriter) Close() {
	if w.closed {
		return
	}
	w.closed = true
	if w.hasPendingHighSurrogate {
		w.hasPendingHighSurrogate = false
		w.emit(w.encodeMalformed(nil))
	}
	if err := w.out.Close(); err != nil {
		panic(NewIOException(err.Error()))
	}
}
