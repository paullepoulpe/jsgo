package jutil

import (
	"bufio"
	"io"
	"unicode/utf8"
)

// Go-only: language-adapter — java.io.OutputStreamWriter with the UTF-8
// charset (or any supported Charset) and CodingErrorAction.REPLACE, over a
// Go io.WriteCloser. Observed behaviour (writers_test.go, cli_vectors_test.go,
// charset_vectors_test.go, against JDK 21's public API): a high surrogate
// that ends one write pairs with a low surrogate that starts the next, so a
// pair split between writes is encoded whole, and an unpaired surrogate is
// encoded as the encoder's replacement ('?'). Close encodes a trailing
// unpaired high surrogate and ends the output (ISO-2022-JP returns to ASCII
// there); a charset that writes a byte-order mark writes it before its first
// char, so an empty output stays empty.

// StreamEncoder encodes strings incrementally in a Charset with
// CodingErrorAction.REPLACE, preserving split surrogate pairs across writes.
type StreamEncoder struct {
	enc                     charsetEncoder
	pendingHighSurrogate    uint16
	hasPendingHighSurrogate bool
	isUTF8                  bool
	isASCII                 bool
}

// NewStreamEncoder creates a streaming encoder for c.
func (c *Charset) NewStreamEncoder() StreamEncoder {
	if c == StandardCharsets_UTF_8 {
		return StreamEncoder{isUTF8: true}
	}
	return StreamEncoder{
		enc:     c.impl().newEncoder(),
		isASCII: c == StandardCharsets_US_ASCII || c == StandardCharsets_ISO_8859_1,
	}
}

// Append appends the encoded bytes of s to dst.
func (se *StreamEncoder) Append(dst []byte, s string) []byte {
	if se.isUTF8 {
		if !se.hasPendingHighSurrogate && utf8.ValidString(s) {
			return append(dst, s...)
		}
		for i := 0; i < len(s); {
			if !se.hasPendingHighSurrogate && s[i] < 0x80 {
				j := i + 1
				for j < len(s) && s[j] < 0x80 {
					j++
				}
				dst = append(dst, s[i:j]...)
				i = j
				continue
			}
			r, sz := decodeWTF8(s[i:])
			i += sz
			if r >= 0x10000 {
				if se.hasPendingHighSurrogate {
					dst = append(dst, '?')
					se.hasPendingHighSurrogate = false
				}
				dst = utf8.AppendRune(dst, r)
				continue
			}
			c := uint16(r)
			switch {
			case c >= 0xD800 && c <= 0xDBFF:
				if se.hasPendingHighSurrogate {
					dst = append(dst, '?')
				}
				se.pendingHighSurrogate, se.hasPendingHighSurrogate = c, true
			case c >= 0xDC00 && c <= 0xDFFF:
				if se.hasPendingHighSurrogate {
					dst = utf8.AppendRune(dst, (rune(se.pendingHighSurrogate)-0xD800)<<10+(rune(c)-0xDC00)+0x10000)
					se.hasPendingHighSurrogate = false
				} else {
					dst = append(dst, '?')
				}
			default:
				if se.hasPendingHighSurrogate {
					dst = append(dst, '?')
					se.hasPendingHighSurrogate = false
				}
				dst = utf8.AppendRune(dst, rune(c))
			}
		}
		return dst
	}
	if se.isASCII && !se.hasPendingHighSurrogate {
		asciiEnd := 0
		for asciiEnd < len(s) && s[asciiEnd] < 0x80 {
			asciiEnd++
		}
		if asciiEnd > 0 {
			dst = append(dst, s[:asciiEnd]...)
			if asciiEnd == len(s) {
				return dst
			}
			s = s[asciiEnd:]
		}
	}
	return se.appendUnits(dst, StringToUTF16(s))
}

func (se *StreamEncoder) appendUnits(dst []byte, units []uint16) []byte {
	i := 0
	if se.hasPendingHighSurrogate && len(units) > 0 {
		hi := se.pendingHighSurrogate
		se.hasPendingHighSurrogate = false
		if lo := units[0]; isLowSurrogateUnit(lo) {
			dst = se.enc.encodeRune(dst, 0x10000+(rune(hi)-0xD800)<<10+(rune(lo)-0xDC00))
			i = 1
		} else {
			dst = se.enc.encodeMalformed(dst)
		}
	}
	if se.isASCII {
		for ; i < len(units); i++ {
			c := units[i]
			if c < 0x80 {
				dst = append(dst, byte(c))
				continue
			}
			switch {
			case isHighSurrogateUnit(c):
				if i+1 == len(units) {
					se.pendingHighSurrogate, se.hasPendingHighSurrogate = c, true
					continue
				}
				if lo := units[i+1]; isLowSurrogateUnit(lo) {
					dst = se.enc.encodeRune(dst, 0x10000+(rune(c)-0xD800)<<10+(rune(lo)-0xDC00))
					i++
					continue
				}
				dst = se.enc.encodeMalformed(dst)
			case isLowSurrogateUnit(c):
				dst = se.enc.encodeMalformed(dst)
			default:
				dst = se.enc.encodeRune(dst, rune(c))
			}
		}
	} else {
		for ; i < len(units); i++ {
			c := units[i]
			switch {
			case isHighSurrogateUnit(c):
				if i+1 == len(units) {
					se.pendingHighSurrogate, se.hasPendingHighSurrogate = c, true
					continue
				}
				if lo := units[i+1]; isLowSurrogateUnit(lo) {
					dst = se.enc.encodeRune(dst, 0x10000+(rune(c)-0xD800)<<10+(rune(lo)-0xDC00))
					i++
					continue
				}
				dst = se.enc.encodeMalformed(dst)
			case isLowSurrogateUnit(c):
				dst = se.enc.encodeMalformed(dst)
			default:
				dst = se.enc.encodeRune(dst, rune(c))
			}
		}
	}
	return dst
}

// Finish flushes any trailing unpaired high surrogate and encoder trailer to dst.
func (se *StreamEncoder) Finish(dst []byte) []byte {
	if se.isUTF8 {
		if se.hasPendingHighSurrogate {
			se.hasPendingHighSurrogate = false
			dst = append(dst, '?')
		}
		return dst
	}
	if se.hasPendingHighSurrogate {
		se.hasPendingHighSurrogate = false
		dst = se.enc.encodeMalformed(dst)
	}
	return se.enc.finish(dst)
}

// WriteTo encodes s directly into bw.
func (se *StreamEncoder) WriteTo(bw *bufio.Writer, s string) error {
	if len(s) == 0 {
		return nil
	}
	if se.isUTF8 && !se.hasPendingHighSurrogate && utf8.ValidString(s) {
		_, err := bw.WriteString(s)
		return err
	}
	if se.isASCII && !se.hasPendingHighSurrogate {
		asciiEnd := 0
		for asciiEnd < len(s) && s[asciiEnd] < 0x80 {
			asciiEnd++
		}
		if asciiEnd == len(s) {
			_, err := bw.WriteString(s)
			return err
		}
		if asciiEnd > 0 {
			if _, err := bw.WriteString(s[:asciiEnd]); err != nil {
				return err
			}
			s = s[asciiEnd:]
		}
	}
	b := se.Append(bw.AvailableBuffer(), s)
	if len(b) == 0 {
		return nil
	}
	_, err := bw.Write(b)
	return err
}

// WriteCharTo encodes a single UTF-16 code unit directly into bw.
func (se *StreamEncoder) WriteCharTo(bw *bufio.Writer, c uint16) error {
	if c < 0x80 && (se.isUTF8 || se.isASCII) && !se.hasPendingHighSurrogate {
		return bw.WriteByte(byte(c))
	}
	return se.WriteTo(bw, StringFromChar(c))
}

// FinishTo flushes any trailing unpaired high surrogate and encoder trailer into bw.
func (se *StreamEncoder) FinishTo(bw *bufio.Writer) error {
	b := se.Finish(bw.AvailableBuffer())
	if len(b) == 0 {
		return nil
	}
	_, err := bw.Write(b)
	return err
}

// OutputStreamWriter is java.io.OutputStreamWriter.
type OutputStreamWriter struct {
	out     io.WriteCloser
	byteBuf []byte
	enc     StreamEncoder
	closed  bool
}

// NewOutputStreamWriter is new OutputStreamWriter(out, UTF_8).
func NewOutputStreamWriter(out io.WriteCloser) *OutputStreamWriter {
	if out == nil {
		panic(NewNullPointerException())
	}
	return &OutputStreamWriter{out: out, enc: StreamEncoder{isUTF8: true}}
}

// NewOutputStreamWriterWithCharset is new OutputStreamWriter(out, cs).
func NewOutputStreamWriterWithCharset(out io.WriteCloser, cs *Charset) Writer {
	if out == nil || cs == nil {
		panic(NewNullPointerException())
	}
	return &OutputStreamWriter{out: out, enc: cs.NewStreamEncoder()}
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
	b = w.enc.Append(b, s)
	w.byteBuf = b[:0]
	return w.tryWriteBytes(b)
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

// Flush is OutputStreamWriter.flush(); the bytes are already in the stream.
func (w *OutputStreamWriter) Flush() {
	if ex := w.TryFlush(); ex != nil {
		panic(ex)
	}
}

// TryFlush checks that the stream is open and returns an IIOException if closed.
func (w *OutputStreamWriter) TryFlush() IIOException {
	if w.closed {
		return NewIOException("Stream closed")
	}
	return nil
}

// Close is OutputStreamWriter.close(): an unpaired high surrogate left
// over from the last write is encoded as unpaired, the output is ended, then
// the stream is closed. Closing a closed writer has no effect.
func (w *OutputStreamWriter) Close() {
	if w.closed {
		return
	}
	w.closed = true
	b := w.enc.Finish(w.byteBuf[:0])
	w.byteBuf = nil
	if ex := w.tryWriteBytes(b); ex != nil {
		_ = w.out.Close()
		panic(ex)
	}
	if err := w.out.Close(); err != nil {
		panic(NewIOException(err.Error()))
	}
}
