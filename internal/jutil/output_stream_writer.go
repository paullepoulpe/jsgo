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
	if err := w.TryWrite(s); err != nil {
		panic(err)
	}
}

// TryWrite writes s and returns an IIOException instead of panicking on I/O failure.
func (w *OutputStreamWriter) TryWrite(s string) IIOException {
	if w.closed {
		return NewIOException("Stream closed")
	}
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
	return w.tryWriteBytes(b)
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

func (w *OutputStreamWriter) tryWriteBytes(b []byte) IIOException {
	if len(b) == 0 {
		return nil
	}
	if _, err := w.out.Write(b); err != nil {
		return NewIOException(err.Error())
	}
	return nil
}

func (w *OutputStreamWriter) writeBytes(b []byte) {
	if err := w.tryWriteBytes(b); err != nil {
		panic(err)
	}
}

// Flush is OutputStreamWriter.flush(); the bytes are already in the stream
// (see above).
func (w *OutputStreamWriter) Flush() {
	w.checkOpen()
}

// TryFlush checks that the stream is open and returns an IIOException if closed.
func (w *OutputStreamWriter) TryFlush() IIOException {
	if w.closed {
		return NewIOException("Stream closed")
	}
	return nil
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
	return &charsetOutputStreamWriter{
		out:     out,
		enc:     cs.impl().newEncoder(),
		byteBuf: make([]byte, 0, 8192),
		isASCII: cs == StandardCharsets_US_ASCII || cs == StandardCharsets_ISO_8859_1,
	}
}

type charsetOutputStreamWriter struct {
	out                     io.WriteCloser
	enc                     charsetEncoder
	byteBuf                 []byte
	pendingHighSurrogate    uint16
	hasPendingHighSurrogate bool
	isASCII                 bool
	closed                  bool
}

func (w *charsetOutputStreamWriter) checkOpen() {
	if w.closed {
		panic(NewIOException("Stream closed"))
	}
}

// Write is OutputStreamWriter.write(String).
func (w *charsetOutputStreamWriter) Write(s string) {
	if err := w.TryWrite(s); err != nil {
		panic(err)
	}
}

// TryWrite writes s and returns an IIOException instead of panicking on I/O failure.
func (w *charsetOutputStreamWriter) TryWrite(s string) IIOException {
	if w.closed {
		return NewIOException("Stream closed")
	}
	b := w.byteBuf[:0]
	if cap(b) < len(s) {
		b = make([]byte, 0, max(len(s), 8192))
	}
	if w.isASCII && !w.hasPendingHighSurrogate {
		asciiEnd := 0
		for asciiEnd < len(s) && s[asciiEnd] < 0x80 {
			asciiEnd++
		}
		if asciiEnd > 0 {
			b = append(b, s[:asciiEnd]...)
			if asciiEnd == len(s) {
				w.byteBuf = b[:0]
				return w.tryEmit(b)
			}
			s = s[asciiEnd:]
		}
	}
	units := StringToUTF16(s)
	return w.tryEncodeUnits(b, units)
}

func (w *charsetOutputStreamWriter) writeUTF16(units []uint16) {
	w.checkOpen()
	b := w.byteBuf[:0]
	if cap(b) < len(units) {
		b = make([]byte, 0, max(len(units), 8192))
	}
	if err := w.tryEncodeUnits(b, units); err != nil {
		panic(err)
	}
}

func (w *charsetOutputStreamWriter) tryEncodeUnits(b []byte, units []uint16) IIOException {
	i := 0
	if w.hasPendingHighSurrogate && len(units) > 0 {
		hi := w.pendingHighSurrogate
		w.hasPendingHighSurrogate = false
		if lo := units[0]; isLowSurrogateUnit(lo) {
			b = w.enc.encodeRune(b, 0x10000+(rune(hi)-0xD800)<<10+(rune(lo)-0xDC00))
			i = 1
		} else {
			b = w.enc.encodeMalformed(b)
		}
	}
	if w.isASCII {
		for ; i < len(units); i++ {
			c := units[i]
			if c < 0x80 {
				b = append(b, byte(c))
				continue
			}
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
	} else {
		for ; i < len(units); i++ {
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
	}
	w.byteBuf = b[:0]
	return w.tryEmit(b)
}

func (w *charsetOutputStreamWriter) tryEmit(b []byte) IIOException {
	if len(b) == 0 {
		return nil
	}
	if _, err := w.out.Write(b); err != nil {
		return NewIOException(err.Error())
	}
	return nil
}

func (w *charsetOutputStreamWriter) emit(b []byte) {
	if err := w.tryEmit(b); err != nil {
		panic(err)
	}
}

// Flush is OutputStreamWriter.flush(); bytes are written through
// immediately.
func (w *charsetOutputStreamWriter) Flush() {
	w.checkOpen()
}

// TryFlush checks that the stream is open and returns an IIOException if closed.
func (w *charsetOutputStreamWriter) TryFlush() IIOException {
	if w.closed {
		return NewIOException("Stream closed")
	}
	return nil
}

// Close is OutputStreamWriter.close(): a high surrogate left over from the
// last write is encoded as unpaired, the output is ended, then the stream
// is closed. Closing a closed writer has no effect.
func (w *charsetOutputStreamWriter) Close() {
	if w.closed {
		return
	}
	w.closed = true
	b := w.byteBuf[:0]
	if w.hasPendingHighSurrogate {
		w.hasPendingHighSurrogate = false
		b = w.enc.encodeMalformed(b)
	}
	b = w.enc.finish(b)
	w.byteBuf = nil
	w.emit(b)
	if err := w.out.Close(); err != nil {
		panic(NewIOException(err.Error()))
	}
}
