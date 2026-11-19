// Clean-room specification based on java.io.Reader, Writer, BufferedReader, BufferedWriter, InputStreamReader, OutputStreamWriter Javadoc.

package javaio

import (
	"io"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Reader is the character stream reader interface.
type Reader interface {
	Read(cbuf []uint16, off, length int) int
	Close() error
}

// Writer is the character stream writer interface.
type Writer interface {
	Write(s string)
	AppendCharSequence(s string)
	AppendCharacter(c uint16)
	Flush()
	Close()
}

// StringReader reads characters from a string, deferring UTF-16 conversion
// until the first character read so unread DrainString/ReadAllString avoids round-tripping.
type StringReader struct {
	s      string
	units  []uint16
	ascii  bool
	inited bool
	offset int
}

// NewStringReader creates a character reader for s.
func NewStringReader(s string) *StringReader {
	return &StringReader{s: s}
}

func (r *StringReader) ensureInit() {
	if r.inited {
		return
	}
	r.inited = true
	if javalang.IsASCII(r.s) {
		r.ascii = true
		return
	}
	r.units = javalang.StringToUTF16(r.s)
	r.s = ""
}

// Read reads characters into cbuf starting at off with length length.
func (r *StringReader) Read(cbuf []uint16, off, length int) int {
	if length == 0 {
		return 0
	}
	r.ensureInit()
	if r.ascii {
		if r.offset >= len(r.s) {
			return -1
		}
		n := min(length, len(r.s)-r.offset)
		src := r.s[r.offset : r.offset+n]
		dst := cbuf[off : off+n]
		for i := 0; i < n; i++ {
			dst[i] = uint16(src[i])
		}
		r.offset += n
		return n
	}
	if r.offset >= len(r.units) {
		return -1
	}
	n := min(length, len(r.units)-r.offset)
	copy(cbuf[off:off+n], r.units[r.offset:r.offset+n])
	r.offset += n
	return n
}

// ReadChar reads a single character, returning false on EOF.
func (r *StringReader) ReadChar() (uint16, bool) {
	r.ensureInit()
	if r.ascii {
		if r.offset >= len(r.s) {
			return 0, false
		}
		ch := uint16(r.s[r.offset])
		r.offset++
		return ch, true
	}
	if r.offset >= len(r.units) {
		return 0, false
	}
	ch := r.units[r.offset]
	r.offset++
	return ch, true
}

// DrainString drains all remaining characters directly without UTF-16 chunk copying.
func (r *StringReader) DrainString() (string, bool) {
	if !r.inited {
		s := r.s
		r.s = ""
		r.inited = true
		r.ascii = true
		r.offset = 0
		if utf8.ValidString(s) {
			return s, true
		}
		return javalang.StringFromUTF16(javalang.StringToUTF16(s)), true
	}
	if r.ascii {
		if r.offset >= len(r.s) {
			return "", true
		}
		out := r.s[r.offset:]
		r.offset = len(r.s)
		return out, true
	}
	if r.offset >= len(r.units) {
		return "", true
	}
	out := javalang.StringFromUTF16(r.units[r.offset:])
	r.offset = len(r.units)
	return out, true
}

// ReadAllString reads all remaining characters as a string.
func (r *StringReader) ReadAllString() string {
	s, _ := r.DrainString()
	return s
}

// Close closes the stream.
func (r *StringReader) Close() error {
	r.inited = true
	r.ascii = true
	r.s = ""
	r.units = nil
	r.offset = 0
	return nil
}

const (
	isrByteBufSize        = 8192
	bufferedReaderBufSize = 4096
)

// InputStreamReader decodes bytes from an io.Reader into UTF-16 code units.
type InputStreamReader struct {
	r      io.Reader
	latin1 bool
	rawBuf [isrByteBufSize]byte
	rawPos int
	rawEnd int
	pendLo int32
	eof    bool
}

// NewInputStreamReader creates a UTF-8 InputStreamReader over r.
func NewInputStreamReader(r io.Reader) *InputStreamReader {
	return &InputStreamReader{r: r, pendLo: -1}
}

// NewInputStreamReaderCharset creates an InputStreamReader with the named charset ("UTF-8", "ISO-8859-1", "US-ASCII").
func NewInputStreamReaderCharset(r io.Reader, charset string) *InputStreamReader {
	latin1 := strings.EqualFold(charset, "ISO-8859-1") ||
		strings.EqualFold(charset, "ISO_8859_1") ||
		strings.EqualFold(charset, "LATIN1")
	return &InputStreamReader{r: r, latin1: latin1, pendLo: -1}
}

func (r *InputStreamReader) fillBytes() bool {
	if r.eof || r.r == nil {
		return false
	}
	avail := r.rawEnd - r.rawPos
	if avail > 0 && r.rawPos > 0 {
		copy(r.rawBuf[:avail], r.rawBuf[r.rawPos:r.rawEnd])
	}
	r.rawPos = 0
	r.rawEnd = avail
	for {
		n, err := r.r.Read(r.rawBuf[r.rawEnd:])
		if n > 0 {
			r.rawEnd += n
			if err != nil {
				r.eof = true
			}
			return true
		}
		if err != nil {
			r.eof = true
			return false
		}
	}
}

// Read reads characters into cbuf[off:off+length].
func (r *InputStreamReader) Read(cbuf []uint16, off, length int) int {
	if length == 0 {
		return 0
	}
	written := 0
	if r.pendLo >= 0 {
		cbuf[off] = uint16(r.pendLo)
		r.pendLo = -1
		written = 1
		if length == 1 {
			return 1
		}
	}
	for written < length {
		avail := r.rawEnd - r.rawPos
		if r.latin1 {
			if avail == 0 {
				if written > 0 || !r.fillBytes() {
					break
				}
				avail = r.rawEnd - r.rawPos
			}
			n := min(length-written, avail)
			src := r.rawBuf[r.rawPos : r.rawPos+n]
			dst := cbuf[off+written : off+written+n]
			for i := 0; i < n; i++ {
				dst[i] = uint16(src[i])
			}
			r.rawPos += n
			written += n
			continue
		}

		if avail == 0 {
			if written > 0 || !r.fillBytes() {
				break
			}
			avail = r.rawEnd - r.rawPos
		}
		// Fast path for ASCII runs.
		for r.rawPos < r.rawEnd && written < length {
			b := r.rawBuf[r.rawPos]
			if b >= 0x80 {
				break
			}
			cbuf[off+written] = uint16(b)
			r.rawPos++
			written++
		}
		if r.rawPos >= r.rawEnd || written >= length {
			continue
		}
		rem := r.rawBuf[r.rawPos:r.rawEnd]
		if len(rem) < utf8.UTFMax && !r.eof && !utf8.FullRune(rem) {
			if written > 0 {
				break
			}
			if !r.fillBytes() {
				rem = r.rawBuf[r.rawPos:r.rawEnd]
				if len(rem) == 0 {
					break
				}
			} else {
				continue
			}
		}
		ru, sz := utf8.DecodeRune(rem)
		r.rawPos += sz
		if ru <= 0xFFFF {
			cbuf[off+written] = uint16(ru)
			written++
		} else {
			hi, lo := utf16.EncodeRune(ru)
			cbuf[off+written] = uint16(hi)
			written++
			if written < length {
				cbuf[off+written] = uint16(lo)
				written++
			} else {
				r.pendLo = int32(uint16(lo))
			}
		}
	}
	if written == 0 {
		return -1
	}
	return written
}

// DrainString drains an unread InputStreamReader directly into a string without UTF-16 chunking.
func (r *InputStreamReader) DrainString() (string, bool) {
	if r.pendLo >= 0 || r.rawPos != r.rawEnd {
		return "", false
	}
	if r.eof || r.r == nil {
		return "", true
	}
	data, _ := io.ReadAll(r.r)
	r.eof = true
	if len(data) == 0 {
		return "", true
	}
	if r.latin1 {
		ascii := true
		for _, b := range data {
			if b >= 0x80 {
				ascii = false
				break
			}
		}
		if ascii {
			return string(data), true
		}
		var sb strings.Builder
		sb.Grow(len(data))
		for _, b := range data {
			sb.WriteRune(rune(b))
		}
		return sb.String(), true
	}
	if utf8.Valid(data) {
		return string(data), true
	}
	var sb strings.Builder
	sb.Grow(len(data))
	for len(data) > 0 {
		ru, sz := utf8.DecodeRune(data)
		data = data[sz:]
		sb.WriteRune(ru)
	}
	return sb.String(), true
}

// ReadAllString reads all remaining characters as a WTF-8/UTF-8 string.
func (r *InputStreamReader) ReadAllString() string {
	if s, ok := r.DrainString(); ok {
		return s
	}
	var units []uint16
	var buf [bufferedReaderBufSize]uint16
	for {
		n := r.Read(buf[:], 0, len(buf))
		if n <= 0 {
			break
		}
		units = append(units, buf[:n]...)
	}
	return javalang.StringFromUTF16(units)
}

// Close closes the underlying reader if it implements io.Closer.
func (r *InputStreamReader) Close() error {
	if c, ok := r.r.(io.Closer); ok {
		return ignoreAlreadyClosed(c.Close())
	}
	return nil
}

// BufferedReader buffers characters in a reusable chunk buffer for zero-alloc per-char reading and fast ReadLine.
type BufferedReader struct {
	r       Reader
	cbuf    [bufferedReaderBufSize]uint16
	pos     int
	end     int
	skipLF  bool
	lineBuf []uint16
}

// NewBufferedReader creates a buffering character reader.
func NewBufferedReader(r Reader) *BufferedReader {
	return &BufferedReader{r: r}
}

func (b *BufferedReader) fill() bool {
	for b.pos >= b.end {
		if b.r == nil {
			b.skipLF = false
			return false
		}
		n := b.r.Read(b.cbuf[:], 0, len(b.cbuf))
		if n <= 0 {
			b.pos = 0
			b.end = 0
			b.skipLF = false
			return false
		}
		b.pos = 0
		b.end = n
		if b.skipLF {
			b.skipLF = false
			if b.cbuf[0] == '\n' {
				b.pos = 1
			}
		}
	}
	return true
}

// Read reads characters into cbuf[off:off+length].
func (b *BufferedReader) Read(cbuf []uint16, off, length int) int {
	if length == 0 {
		return 0
	}
	copied := 0
	for copied < length {
		if b.pos < b.end {
			n := min(length-copied, b.end-b.pos)
			copy(cbuf[off+copied:off+copied+n], b.cbuf[b.pos:b.pos+n])
			b.pos += n
			copied += n
			if copied == length {
				return copied
			}
		}
		if !b.skipLF && length-copied >= len(b.cbuf) {
			if b.r == nil {
				break
			}
			n := b.r.Read(cbuf, off+copied, length-copied)
			if n <= 0 {
				break
			}
			copied += n
			continue
		}
		if !b.fill() {
			break
		}
	}
	if copied == 0 {
		return -1
	}
	return copied
}

// ReadLine reads a line of text terminated by '\n', '\r', "\r\n", or EOF.
func (b *BufferedReader) ReadLine() (string, error) {
	if b.pos >= b.end && !b.fill() {
		return "", io.EOF
	}
	start := b.pos
	buf := b.cbuf[start:b.end]
	for i, ch := range buf {
		if ch == '\n' || ch == '\r' {
			b.pos = start + i + 1
			if ch == '\r' {
				if b.pos < b.end {
					if b.cbuf[b.pos] == '\n' {
						b.pos++
					}
				} else {
					b.skipLF = true
				}
			}
			return javalang.StringFromUTF16(buf[:i]), nil
		}
	}
	line := append(b.lineBuf[:0], buf...)
	b.pos = b.end
	for b.fill() {
		start = b.pos
		buf = b.cbuf[start:b.end]
		for i, ch := range buf {
			if ch == '\n' || ch == '\r' {
				b.pos = start + i + 1
				if ch == '\r' {
					if b.pos < b.end {
						if b.cbuf[b.pos] == '\n' {
							b.pos++
						}
					} else {
						b.skipLF = true
					}
				}
				line = append(line, buf[:i]...)
				s := javalang.StringFromUTF16(line)
				if cap(line) <= 65536 {
					b.lineBuf = line[:0]
				} else {
					b.lineBuf = nil
				}
				return s, nil
			}
		}
		line = append(line, buf...)
		b.pos = b.end
	}
	s := javalang.StringFromUTF16(line)
	if cap(line) <= 65536 {
		b.lineBuf = line[:0]
	} else {
		b.lineBuf = nil
	}
	return s, nil
}

// DrainString drains the underlying reader directly when BufferedReader has no buffered characters.
func (b *BufferedReader) DrainString() (string, bool) {
	if b.pos == 0 && b.end == 0 && !b.skipLF && b.r != nil {
		if ds, ok := b.r.(interface{ DrainString() (string, bool) }); ok {
			return ds.DrainString()
		}
	}
	return "", false
}

// Close closes the underlying reader.
func (b *BufferedReader) Close() error {
	if b.r != nil {
		return b.r.Close()
	}
	return nil
}

// StringWriter collects characters in a string builder.
type StringWriter struct {
	sb strings.Builder
}

// NewStringWriter creates an empty StringWriter.
func NewStringWriter() *StringWriter {
	return &StringWriter{}
}

func (w *StringWriter) Write(s string) {
	w.sb.WriteString(s)
}

func (w *StringWriter) Flush() {}

func (w *StringWriter) Close() {}

func (w *StringWriter) String() string {
	return w.sb.String()
}

func (w *StringWriter) AppendCharSequence(s string) {
	w.sb.WriteString(s)
}

func (w *StringWriter) AppendCharacter(c uint16) {
	w.sb.WriteString(javalang.StringFromChar(c))
}

// OutputStreamWriter encodes characters to an underlying io.Writer.
type OutputStreamWriter struct {
	w       io.Writer
	highSur uint16
	hasHigh bool
}

// NewOutputStreamWriter creates a UTF-8 OutputStreamWriter wrapping w.
func NewOutputStreamWriter(w io.Writer) *OutputStreamWriter {
	return &OutputStreamWriter{w: w}
}

func (w *OutputStreamWriter) flushDanglingHigh() {
	if w.hasHigh {
		w.hasHigh = false
		var buf [3]byte
		n := utf8.EncodeRune(buf[:], utf8.RuneError)
		_, _ = w.w.Write(buf[:n])
	}
}

func (w *OutputStreamWriter) Write(s string) {
	if len(s) == 0 {
		return
	}
	if !w.hasHigh && utf8.ValidString(s) {
		_, _ = io.WriteString(w.w, s)
		return
	}
	units := javalang.StringToUTF16(s)
	for _, u := range units {
		w.AppendCharacter(u)
	}
}

func (w *OutputStreamWriter) AppendCharSequence(s string) {
	w.Write(s)
}

func (w *OutputStreamWriter) AppendCharacter(c uint16) {
	if w.hasHigh {
		hi := w.highSur
		w.hasHigh = false
		if c >= 0xDC00 && c <= 0xDFFF {
			cp := utf16.DecodeRune(rune(hi), rune(c))
			var buf [4]byte
			n := utf8.EncodeRune(buf[:], cp)
			_, _ = w.w.Write(buf[:n])
			return
		}
		var errBuf [3]byte
		n := utf8.EncodeRune(errBuf[:], utf8.RuneError)
		_, _ = w.w.Write(errBuf[:n])
	}
	if c >= 0xD800 && c <= 0xDBFF {
		w.highSur = c
		w.hasHigh = true
		return
	}
	if c >= 0xDC00 && c <= 0xDFFF {
		var errBuf [3]byte
		n := utf8.EncodeRune(errBuf[:], utf8.RuneError)
		_, _ = w.w.Write(errBuf[:n])
		return
	}
	var buf [3]byte
	n := utf8.EncodeRune(buf[:], rune(c))
	_, _ = w.w.Write(buf[:n])
}

func (w *OutputStreamWriter) Flush() {
	if f, ok := w.w.(interface{ Flush() error }); ok {
		_ = f.Flush()
	}
}

func (w *OutputStreamWriter) Close() {
	w.flushDanglingHigh()
	w.Flush()
	if c, ok := w.w.(io.Closer); ok {
		_ = c.Close()
	}
}

// BufferedWriter wraps a Writer with newline helper support.
type BufferedWriter struct {
	w      Writer
	closed bool
}

// NewBufferedWriter creates a BufferedWriter wrapping w.
func NewBufferedWriter(w Writer) *BufferedWriter {
	return &BufferedWriter{w: w}
}

func (b *BufferedWriter) Write(s string) {
	b.w.Write(s)
}

func (b *BufferedWriter) NewLine() {
	b.w.Write("\n")
}

func (b *BufferedWriter) AppendCharSequence(s string) {
	b.w.AppendCharSequence(s)
}

func (b *BufferedWriter) AppendCharacter(c uint16) {
	b.w.AppendCharacter(c)
}

func (b *BufferedWriter) Flush() {
	b.w.Flush()
}

func (b *BufferedWriter) Close() {
	if b.closed {
		return
	}
	b.closed = true
	b.w.Flush()
	b.w.Close()
}
