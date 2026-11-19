package javanio

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/bits"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// CodingErrorAction models java.nio.charset.CodingErrorAction.
type CodingErrorAction int

const (
	CodingErrorAction_IGNORE CodingErrorAction = iota
	CodingErrorAction_REPLACE
	CodingErrorAction_REPORT
)

type CharacterCodingException struct {
	Msg string
}

func (e *CharacterCodingException) Error() string {
	return "java.nio.charset.CharacterCodingException: " + e.Msg
}

type MalformedInputException struct {
	InputLength int
}

func (e *MalformedInputException) Error() string {
	return fmt.Sprintf("java.nio.charset.MalformedInputException: Input length = %d", e.InputLength)
}

type UnmappableCharacterException struct {
	InputLength int
}

func (e *UnmappableCharacterException) Error() string {
	return fmt.Sprintf("java.nio.charset.UnmappableCharacterException: Input length = %d", e.InputLength)
}

type IllegalCharsetNameException struct {
	CharsetName string
}

func (e *IllegalCharsetNameException) Error() string {
	return "java.nio.charset.IllegalCharsetNameException: " + e.CharsetName
}

type UnsupportedCharsetException struct {
	CharsetName string
}

func (e *UnsupportedCharsetException) Error() string {
	return "java.nio.charset.UnsupportedCharsetException: " + e.CharsetName
}

type charsetOps interface {
	encode(s string) []byte
	decode(b []byte) (string, error)
	decodeReplace(b []byte) string
	canEncodeChar(ch uint16) bool
}

// Charset models java.nio.charset.Charset.
type Charset struct {
	name    string
	aliases []string
	build   func() charsetOps
	once    sync.Once
	ops     charsetOps
}

func (c *Charset) impl() charsetOps {
	c.once.Do(func() { c.ops = c.build() })
	return c.ops
}

func (c *Charset) Name() string   { return c.name }
func (c *Charset) String() string { return c.name }

func (c *Charset) Equals(other any) bool {
	o, ok := other.(*Charset)
	return ok && o != nil && strings.EqualFold(c.name, o.name)
}

func (c *Charset) HashCode() int32 {
	return javalang.StringHashCode(javalang.StringToLowerCaseASCII(c.name))
}

var emptyEncodedBytes = []byte{}

const (
	decodeDedupMaxLen        = 64
	decodeDedupSeenLen       = 16384
	decodeDedupSeenMask      = decodeDedupSeenLen - 1
	decodeDedupNumShards     = 64
	decodeDedupShardMask     = decodeDedupNumShards - 1
	decodeDedupSetsPerShard  = 256 // 64 * 256 * 2 = 32,768 in-place entries (0 bucket allocs)
	decodeDedupSetMask       = decodeDedupSetsPerShard - 1
)

type decodeDedupShard struct {
	mu   sync.RWMutex
	sets [decodeDedupSetsPerShard][2]string
}

var (
	decodeDedupSeen   [decodeDedupSeenLen]atomic.Uint32
	decodeDedupShards [decodeDedupNumShards]decodeDedupShard
)

func hashASCIIBytesUpTo64(b []byte) uint64 {
	n := len(b)
	if n >= 8 {
		mid := (n - 8) >> 1
		h := uint64(n)*0x9e3779b97f4a7c15 ^
			binary.LittleEndian.Uint64(b[:8]) ^
			(binary.LittleEndian.Uint64(b[mid:mid+8]) * 0xbf58476d1ce4e5b9) ^
			(binary.LittleEndian.Uint64(b[n-8:]) * 0x94d049bb133111eb)
		return h ^ (h >> 29)
	}
	if n >= 4 {
		h := uint64(n)*0x9e3779b97f4a7c15 ^
			(uint64(binary.LittleEndian.Uint32(b[:4])) * 0xbf58476d1ce4e5b9) ^
			(uint64(binary.LittleEndian.Uint32(b[n-4:])) * 0x94d049bb133111eb)
		return h ^ (h >> 29)
	}
	h := (uint64(b[0]) | uint64(b[n>>1])<<8 | uint64(b[n-1])<<16 | uint64(n)<<24) * 0x9e3779b97f4a7c15
	return h ^ (h >> 29)
}

func decodeShortASCIIDedup(b []byte) string {
	h := hashASCIIBytesUpTo64(b)
	seenIdx := h & decodeDedupSeenMask
	tag := uint32(h>>32) | 1
	if decodeDedupSeen[seenIdx].Load() != tag {
		decodeDedupSeen[seenIdx].Store(tag)
		return string(b)
	}
	sh := &decodeDedupShards[h&decodeDedupShardMask]
	slot := (h >> 6) & decodeDedupSetMask
	sh.mu.RLock()
	s0, s1 := sh.sets[slot][0], sh.sets[slot][1]
	sh.mu.RUnlock()
	if s0 == string(b) {
		return s0
	}
	if s1 == string(b) {
		return s1
	}
	out := string(b)
	if sh.mu.TryLock() {
		sh.sets[slot][1] = sh.sets[slot][0]
		sh.sets[slot][0] = out
		sh.mu.Unlock()
	}
	return out
}

func (c *Charset) Encode(s string) []byte {
	switch c {
	case StandardCharsets_UTF_8:
		if len(s) == 0 {
			return emptyEncodedBytes[:0:0]
		}
		if javalang.IsASCII(s) || utf8.ValidString(s) {
			return []byte(s)
		}
		return appendUTF8NonCanonical(make([]byte, 0, len(s)), s)
	case StandardCharsets_US_ASCII:
		if len(s) == 0 {
			return emptyEncodedBytes[:0:0]
		}
		if javalang.IsASCII(s) {
			return []byte(s)
		}
		return appendASCIINonPure(make([]byte, 0, len(s)), s)
	case StandardCharsets_ISO_8859_1:
		if len(s) == 0 {
			return emptyEncodedBytes[:0:0]
		}
		if javalang.IsASCII(s) {
			return []byte(s)
		}
		return appendISO88591NonPure(make([]byte, 0, len(s)), s)
	}
	return c.impl().encode(s)
}

// AppendEncoded appends the charset-encoded bytes of s directly to dst without
// allocating an intermediate []byte on the common UTF-8 / ASCII / Latin-1 paths.
func (c *Charset) AppendEncoded(dst []byte, s string) []byte {
	if len(s) == 0 {
		return dst
	}
	switch c {
	case StandardCharsets_UTF_8:
		if javalang.IsASCII(s) || utf8.ValidString(s) {
			return append(dst, s...)
		}
		return appendUTF8NonCanonical(dst, s)
	case StandardCharsets_US_ASCII:
		if javalang.IsASCII(s) {
			return append(dst, s...)
		}
		return appendASCIINonPure(dst, s)
	case StandardCharsets_ISO_8859_1:
		if javalang.IsASCII(s) {
			return append(dst, s...)
		}
		return appendISO88591NonPure(dst, s)
	}
	return append(dst, c.impl().encode(s)...)
}

//go:noinline
func (c *Charset) decodeReplaceSlow(b []byte) string {
	switch c {
	case StandardCharsets_UTF_8:
		return utf8Ops{}.decodeReplace(b)
	case StandardCharsets_US_ASCII:
		return asciiOps{}.decodeReplace(b)
	case StandardCharsets_ISO_8859_1:
		return decodeISO88591NonASCII(b)
	case StandardCharsets_UTF_16:
		return utf16Ops{bigEndian: true, marked: true}.decodeReplace(b)
	case StandardCharsets_UTF_16BE:
		return utf16Ops{bigEndian: true, marked: false}.decodeReplace(b)
	case StandardCharsets_UTF_16LE:
		return utf16Ops{bigEndian: false, marked: false}.decodeReplace(b)
	}
	cp := make([]byte, len(b))
	copy(cp, b)
	return c.impl().decodeReplace(cp)
}

//go:noinline
func (c *Charset) decodeStrictSlow(b []byte) (string, error) {
	switch c {
	case StandardCharsets_UTF_16:
		return utf16Ops{bigEndian: true, marked: true}.decode(b)
	case StandardCharsets_UTF_16BE:
		return utf16Ops{bigEndian: true, marked: false}.decode(b)
	case StandardCharsets_UTF_16LE:
		return utf16Ops{bigEndian: false, marked: false}.decode(b)
	}
	cp := make([]byte, len(b))
	copy(cp, b)
	return c.impl().decode(cp)
}

func (c *Charset) Decode(b []byte) string {
	n := len(b)
	switch c {
	case StandardCharsets_UTF_8:
		if n == 0 {
			return ""
		}
		if n == 1 && b[0] < 0x80 {
			return javalang.StringFromChar(uint16(b[0]))
		}
		if n <= decodeDedupMaxLen && isBytesASCIIUpTo64(b) {
			return decodeShortASCIIDedup(b)
		}
		if _, err := validateUTF8Chunk(b, true); err == nil {
			return string(b)
		}
	case StandardCharsets_US_ASCII, StandardCharsets_ISO_8859_1:
		if n == 0 {
			return ""
		}
		if n == 1 && b[0] < 0x80 {
			return javalang.StringFromChar(uint16(b[0]))
		}
		if n <= decodeDedupMaxLen && isBytesASCIIUpTo64(b) {
			return decodeShortASCIIDedup(b)
		}
		if isBytesASCII(b) {
			return string(b)
		}
	}
	return c.decodeReplaceSlow(b)
}

func (c *Charset) DecodeStrict(b []byte) (string, error) {
	n := len(b)
	switch c {
	case StandardCharsets_UTF_8:
		if n == 0 {
			return "", nil
		}
		if n == 1 && b[0] < 0x80 {
			return javalang.StringFromChar(uint16(b[0])), nil
		}
		if n <= decodeDedupMaxLen && isBytesASCIIUpTo64(b) {
			return decodeShortASCIIDedup(b), nil
		}
		if _, err := validateUTF8Chunk(b, true); err != nil {
			return "", err
		}
		return string(b), nil
	case StandardCharsets_US_ASCII:
		if n == 0 {
			return "", nil
		}
		if n == 1 && b[0] < 0x80 {
			return javalang.StringFromChar(uint16(b[0])), nil
		}
		if n <= decodeDedupMaxLen && isBytesASCIIUpTo64(b) {
			return decodeShortASCIIDedup(b), nil
		}
		if isBytesASCII(b) {
			return string(b), nil
		}
		return "", &MalformedInputException{InputLength: 1}
	case StandardCharsets_ISO_8859_1:
		if n == 0 {
			return "", nil
		}
		if n == 1 && b[0] < 0x80 {
			return javalang.StringFromChar(uint16(b[0])), nil
		}
		if n <= decodeDedupMaxLen && isBytesASCIIUpTo64(b) {
			return decodeShortASCIIDedup(b), nil
		}
		if isBytesASCII(b) {
			return string(b), nil
		}
		return decodeISO88591NonASCII(b), nil
	}
	return c.decodeStrictSlow(b)
}

func (c *Charset) CanEncodeChar(ch uint16) bool {
	if ch >= 0xD800 && ch <= 0xDFFF {
		return false
	}
	return c.impl().canEncodeChar(ch)
}

func (c *Charset) CanEncode(s string) bool {
	units := javalang.StringToUTF16(s)
	for i := 0; i < len(units); i++ {
		u := units[i]
		if u >= 0xD800 && u <= 0xDBFF {
			if i+1 >= len(units) || units[i+1] < 0xDC00 || units[i+1] > 0xDFFF {
				return false
			}
			// Supplementary code point: encodable in Unicode charsets, check via encode round-trip otherwise.
			if c != StandardCharsets_UTF_8 && c != StandardCharsets_UTF_16 &&
				c != StandardCharsets_UTF_16BE && c != StandardCharsets_UTF_16LE &&
				!strings.HasPrefix(c.name, "UTF-32") && !strings.HasPrefix(c.name, "X-UTF-32") && c.name != "GB18030" {
				return false
			}
			i++
			continue
		}
		if u >= 0xDC00 && u <= 0xDFFF {
			return false
		}
		if !c.impl().canEncodeChar(u) {
			return false
		}
	}
	return true
}

// CharsetEncoder models java.nio.charset.CharsetEncoder.
type CharsetEncoder struct {
	cs         *Charset
	malformed  CodingErrorAction
	unmappable CodingErrorAction
}

func (c *Charset) NewEncoder() *CharsetEncoder {
	return &CharsetEncoder{
		cs:         c,
		malformed:  CodingErrorAction_REPORT,
		unmappable: CodingErrorAction_REPORT,
	}
}

func (e *CharsetEncoder) Charset() *Charset { return e.cs }

func (e *CharsetEncoder) OnMalformedInput(action CodingErrorAction) *CharsetEncoder {
	e.malformed = action
	return e
}

func (e *CharsetEncoder) OnUnmappableCharacter(action CodingErrorAction) *CharsetEncoder {
	e.unmappable = action
	return e
}

func (e *CharsetEncoder) CanEncodeChar(ch uint16) bool { return e.cs.CanEncodeChar(ch) }
func (e *CharsetEncoder) CanEncode(s string) bool      { return e.cs.CanEncode(s) }

func StringGetBytes(s string, cs *Charset) []byte {
	if cs == nil {
		cs = StandardCharsets_UTF_8
	}
	return cs.Encode(s)
}

func NewStringFromBytes(b []byte, cs *Charset) string {
	if cs == nil {
		cs = StandardCharsets_UTF_8
	}
	return cs.Decode(b)
}

// --- Reader / Writer Charset Bridges ---

const charsetWriterBufSize = 8192

var (
	charsetWriterBufChan = make(chan *[charsetWriterBufSize]byte, 16)
	charsetWriterBufPool = sync.Pool{
		New: func() any {
			return new([charsetWriterBufSize]byte)
		},
	}
)

type charsetWriter struct {
	out    io.Writer
	closer io.Closer
	cs     *Charset
	buf    *[charsetWriterBufSize]byte
	n      int
	closed bool
}

func NewOutputStreamWriterWithCharset(out io.Writer, cs *Charset) javaio.Writer {
	if cs == nil {
		cs = StandardCharsets_UTF_8
	}
	c, _ := out.(io.Closer)
	return &charsetWriter{out: out, closer: c, cs: cs}
}

func NewOutputStreamWriterWithEncoder(out io.Writer, enc *CharsetEncoder) javaio.Writer {
	cs := StandardCharsets_UTF_8
	if enc != nil && enc.cs != nil {
		cs = enc.cs
	}
	return NewOutputStreamWriterWithCharset(out, cs)
}

func (w *charsetWriter) ensureBuf() *[charsetWriterBufSize]byte {
	if w.buf == nil {
		select {
		case b := <-charsetWriterBufChan:
			w.buf = b
		default:
			w.buf = charsetWriterBufPool.Get().(*[charsetWriterBufSize]byte)
		}
	}
	return w.buf
}

func (w *charsetWriter) flushBuf() {
	if w.n > 0 {
		n := w.n
		w.n = 0
		if _, err := w.out.Write(w.buf[:n]); err != nil {
			panic(javaio.NewIOException(err.Error()))
		}
	}
}

func (w *charsetWriter) Write(s string) {
	if w.closed {
		panic(javaio.NewIOException("Stream closed"))
	}
	if len(s) == 0 {
		return
	}
	if (w.cs == StandardCharsets_UTF_8 && (javalang.IsASCII(s) || utf8.ValidString(s))) ||
		((w.cs == StandardCharsets_US_ASCII || w.cs == StandardCharsets_ISO_8859_1) && javalang.IsASCII(s)) {
		if len(s) >= charsetWriterBufSize {
			w.flushBuf()
			if sw, ok := w.out.(io.StringWriter); ok {
				if _, err := sw.WriteString(s); err != nil {
					panic(javaio.NewIOException(err.Error()))
				}
				return
			}
			buf := w.ensureBuf()
			for len(s) >= charsetWriterBufSize {
				copy(buf[:], s[:charsetWriterBufSize])
				if _, err := w.out.Write(buf[:]); err != nil {
					panic(javaio.NewIOException(err.Error()))
				}
				s = s[charsetWriterBufSize:]
			}
			if len(s) > 0 {
				copy(buf[:], s)
				w.n = len(s)
			}
			return
		}
		buf := w.ensureBuf()
		if w.n+len(s) > charsetWriterBufSize {
			w.flushBuf()
		}
		copy(buf[w.n:], s)
		w.n += len(s)
		return
	}
	if (w.cs == StandardCharsets_UTF_8 || w.cs == StandardCharsets_US_ASCII || w.cs == StandardCharsets_ISO_8859_1) &&
		len(s) <= charsetWriterBufSize {
		buf := w.ensureBuf()
		if w.n+len(s) > charsetWriterBufSize {
			w.flushBuf()
		}
		out := w.cs.AppendEncoded(buf[w.n:w.n], s)
		w.n += len(out)
		return
	}
	enc := w.cs.Encode(s)
	if len(enc) >= charsetWriterBufSize {
		w.flushBuf()
		if _, err := w.out.Write(enc); err != nil {
			panic(javaio.NewIOException(err.Error()))
		}
		return
	}
	buf := w.ensureBuf()
	if w.n+len(enc) > charsetWriterBufSize {
		w.flushBuf()
	}
	copy(buf[w.n:], enc)
	w.n += len(enc)
}

func (w *charsetWriter) AppendCharSequence(s string) {
	w.Write(s)
}

func (w *charsetWriter) AppendCharacter(c uint16) {
	if w.closed {
		panic(javaio.NewIOException("Stream closed"))
	}
	if c < 0x80 && (w.cs == StandardCharsets_UTF_8 || w.cs == StandardCharsets_US_ASCII || w.cs == StandardCharsets_ISO_8859_1) {
		buf := w.ensureBuf()
		if w.n == charsetWriterBufSize {
			w.flushBuf()
		}
		buf[w.n] = byte(c)
		w.n++
		return
	}
	w.Write(javalang.StringFromChar(c))
}

func (w *charsetWriter) WriteChar(c uint16) {
	w.AppendCharacter(c)
}

func (w *charsetWriter) WriteChars(cbuf []uint16, off, length int) {
	if w.closed {
		panic(javaio.NewIOException("Stream closed"))
	}
	if length <= 0 {
		return
	}
	sub := cbuf[off : off+length]
	if w.cs == StandardCharsets_UTF_8 || w.cs == StandardCharsets_US_ASCII || w.cs == StandardCharsets_ISO_8859_1 {
		buf := w.ensureBuf()
		for i, u := range sub {
			if u >= 0x80 {
				w.Write(javalang.StringFromUTF16(sub[i:]))
				return
			}
			if w.n == charsetWriterBufSize {
				w.flushBuf()
			}
			buf[w.n] = byte(u)
			w.n++
		}
		return
	}
	w.Write(javalang.StringFromUTF16(sub))
}

func (w *charsetWriter) Flush() {
	if w.closed {
		panic(javaio.NewIOException("Stream closed"))
	}
	w.flushBuf()
	if f, ok := w.out.(interface{ Flush() error }); ok {
		if err := f.Flush(); err != nil {
			panic(javaio.NewIOException(err.Error()))
		}
	} else if f2, ok := w.out.(interface{ Flush() }); ok {
		f2.Flush()
	}
}

func (w *charsetWriter) Close() {
	if w.closed {
		return
	}
	w.closed = true
	var flushErr error
	if w.n > 0 {
		n := w.n
		w.n = 0
		_, flushErr = w.out.Write(w.buf[:n])
	}
	if w.buf != nil {
		select {
		case charsetWriterBufChan <- w.buf:
		default:
			charsetWriterBufPool.Put(w.buf)
		}
		w.buf = nil
	}
	var closeErr error
	if w.closer != nil {
		closeErr = w.closer.Close()
	}
	if flushErr != nil {
		panic(javaio.NewIOException(flushErr.Error()))
	}
	if closeErr != nil {
		panic(javaio.NewIOException(closeErr.Error()))
	}
}

func NewInputStreamReaderWithCharset(in io.Reader, cs *Charset) javaio.Reader {
	if cs == nil || cs.Equals(StandardCharsets_UTF_8) {
		return javaio.NewInputStreamReader(in)
	}
	data, err := io.ReadAll(in)
	if c, ok := in.(io.Closer); ok {
		_ = c.Close()
	}
	if err != nil {
		return &errReader{err: err}
	}
	return javaio.NewStringReader(cs.Decode(data))
}

type errReader struct{ err error }

func (r *errReader) Read(_ []uint16, _, _ int) int {
	panic(javaio.NewIOException(r.err.Error()))
}
func (r *errReader) Close() error { return nil }

// --- Built-in Charset Implementations ---

type utf8Ops struct{}

func appendUTF8NonCanonical(dst []byte, s string) []byte {
	units := javalang.StringToUTF16(s)
	for i := 0; i < len(units); i++ {
		c := units[i]
		switch {
		case c < 0x80:
			dst = append(dst, byte(c))
		case c < 0x800:
			dst = append(dst, byte(0xC0|c>>6), byte(0x80|c&0x3F))
		case c >= 0xD800 && c <= 0xDBFF && i+1 < len(units) && units[i+1] >= 0xDC00 && units[i+1] <= 0xDFFF:
			r := 0x10000 + (rune(c)-0xD800)<<10 + (rune(units[i+1]) - 0xDC00)
			dst = utf8.AppendRune(dst, r)
			i++
		case c >= 0xD800 && c <= 0xDFFF:
			dst = append(dst, '?')
		default:
			dst = append(dst, byte(0xE0|c>>12), byte(0x80|(c>>6)&0x3F), byte(0x80|c&0x3F))
		}
	}
	return dst
}

func (utf8Ops) encode(s string) []byte {
	if len(s) == 0 {
		return emptyEncodedBytes[:0:0]
	}
	if javalang.IsASCII(s) || utf8.ValidString(s) {
		return []byte(s)
	}
	return appendUTF8NonCanonical(make([]byte, 0, len(s)), s)
}

func isBytesASCIIUpTo64(b []byte) bool {
	n := len(b)
	if n > 32 {
		wHead := binary.LittleEndian.Uint64(b[:8]) |
			binary.LittleEndian.Uint64(b[8:16]) |
			binary.LittleEndian.Uint64(b[16:24]) |
			binary.LittleEndian.Uint64(b[24:32])
		wTail := binary.LittleEndian.Uint64(b[n-32:n-24]) |
			binary.LittleEndian.Uint64(b[n-24:n-16]) |
			binary.LittleEndian.Uint64(b[n-16:n-8]) |
			binary.LittleEndian.Uint64(b[n-8:])
		return ((wHead | wTail) & 0x8080808080808080) == 0
	}
	if n >= 16 {
		w := binary.LittleEndian.Uint64(b[:8]) |
			binary.LittleEndian.Uint64(b[8:16]) |
			binary.LittleEndian.Uint64(b[n-16:n-8]) |
			binary.LittleEndian.Uint64(b[n-8:])
		return (w & 0x8080808080808080) == 0
	}
	if n >= 8 {
		w := binary.LittleEndian.Uint64(b[:8]) | binary.LittleEndian.Uint64(b[n-8:])
		return (w & 0x8080808080808080) == 0
	}
	if n >= 4 {
		w := binary.LittleEndian.Uint32(b[:4]) | binary.LittleEndian.Uint32(b[n-4:])
		return (w & 0x80808080) == 0
	}
	if n > 0 {
		return (b[0] | b[n>>1] | b[n-1]) < 0x80
	}
	return true
}

func isBytesASCII(b []byte) bool {
	if len(b) <= 64 {
		return isBytesASCIIUpTo64(b)
	}
	for len(b) >= 32 {
		w0 := binary.LittleEndian.Uint64(b)
		w1 := binary.LittleEndian.Uint64(b[8:])
		w2 := binary.LittleEndian.Uint64(b[16:])
		w3 := binary.LittleEndian.Uint64(b[24:])
		if (w0|w1|w2|w3)&0x8080808080808080 != 0 {
			return false
		}
		b = b[32:]
	}
	for len(b) >= 8 {
		if binary.LittleEndian.Uint64(b)&0x8080808080808080 != 0 {
			return false
		}
		b = b[8:]
	}
	for _, x := range b {
		if x >= 0x80 {
			return false
		}
	}
	return true
}

func isStringASCII(s string) bool {
	return javalang.IsASCII(s)
}

// validateUTF8Chunk validates b in-place using 32-byte uint64 ASCII word skips
// and utf8.DecodeRune on non-ASCII bytes. When !atEOF and b ends with a valid
// incomplete 1..3-byte prefix of a multi-byte UTF-8 rune, it returns the prefix
// length before that trailing sequence so streaming callers can carry those
// 1..3 bytes into the next read chunk without extra allocations.
func validateUTF8Chunk(b []byte, atEOF bool) (int, error) {
	i := 0
	n := len(b)
	for i < n {
		for i+32 <= n {
			w0 := binary.LittleEndian.Uint64(b[i:])
			w1 := binary.LittleEndian.Uint64(b[i+8:])
			w2 := binary.LittleEndian.Uint64(b[i+16:])
			w3 := binary.LittleEndian.Uint64(b[i+24:])
			if (w0|w1|w2|w3)&0x8080808080808080 != 0 {
				if m := w0 & 0x8080808080808080; m != 0 {
					i += bits.TrailingZeros64(m) >> 3
					goto checkNonASCII
				}
				if m := w1 & 0x8080808080808080; m != 0 {
					i += 8 + (bits.TrailingZeros64(m) >> 3)
					goto checkNonASCII
				}
				if m := w2 & 0x8080808080808080; m != 0 {
					i += 16 + (bits.TrailingZeros64(m) >> 3)
					goto checkNonASCII
				}
				i += 24 + (bits.TrailingZeros64(w3&0x8080808080808080) >> 3)
				goto checkNonASCII
			}
			i += 32
		}
		for i+8 <= n {
			if m := binary.LittleEndian.Uint64(b[i:]) & 0x8080808080808080; m != 0 {
				i += bits.TrailingZeros64(m) >> 3
				goto checkNonASCII
			}
			i += 8
		}
		for i < n && b[i] < 0x80 {
			i++
		}
		if i >= n {
			break
		}
	checkNonASCII:
		if !atEOF && n-i < utf8.UTFMax && !utf8.FullRune(b[i:]) {
			return i, nil
		}
		r, width := utf8.DecodeRune(b[i:])
		if r == utf8.RuneError && width == 1 {
			return 0, &MalformedInputException{InputLength: utf8MalformedLength(b[i:])}
		}
		i += width
	}
	return n, nil
}

func utf8MalformedLength(b []byte) int {
	n := len(b)
	cont := func(j int, lo, hi byte) bool { return j < n && b[j] >= lo && b[j] <= hi }
	for i := 0; i < n; {
		b1 := b[i]
		switch {
		case b1 < 0x80:
			i++
		case b1 >= 0xC2 && b1 <= 0xDF:
			if cont(i+1, 0x80, 0xBF) {
				i += 2
			} else {
				return 1
			}
		case b1 >= 0xE0 && b1 <= 0xEF:
			lo, hi := byte(0x80), byte(0xBF)
			if b1 == 0xE0 {
				lo = 0xA0
			} else if b1 == 0xED {
				hi = 0x9F
			}
			if !cont(i+1, lo, hi) {
				return 1
			}
			if !cont(i+2, 0x80, 0xBF) {
				return 2
			}
			i += 3
		case b1 >= 0xF0 && b1 <= 0xF4:
			lo, hi := byte(0x80), byte(0xBF)
			if b1 == 0xF0 {
				lo = 0x90
			} else if b1 == 0xF4 {
				hi = 0x8F
			}
			if !cont(i+1, lo, hi) {
				return 1
			}
			if !cont(i+2, 0x80, 0xBF) {
				return 2
			}
			if !cont(i+3, 0x80, 0xBF) {
				return 3
			}
			i += 4
		default:
			return 1
		}
	}
	return 1
}

func (u utf8Ops) decode(b []byte) (string, error) {
	if _, err := validateUTF8Chunk(b, true); err != nil {
		return "", err
	}
	return string(b), nil
}

func (utf8Ops) decodeReplace(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	var sb strings.Builder
	sb.Grow(len(b))
	n := len(b)
	cont := func(j int, lo, hi byte) bool { return j < n && b[j] >= lo && b[j] <= hi }
	for i := 0; i < n; {
		b1 := b[i]
		switch {
		case b1 < 0x80:
			sb.WriteByte(b1)
			i++
		case b1 >= 0xC2 && b1 <= 0xDF:
			if cont(i+1, 0x80, 0xBF) {
				sb.Write(b[i : i+2])
				i += 2
			} else {
				sb.WriteRune(utf8.RuneError)
				i++
			}
		case b1 >= 0xE0 && b1 <= 0xEF:
			lo := byte(0x80)
			if b1 == 0xE0 {
				lo = 0xA0
			}
			if !cont(i+1, lo, 0xBF) {
				sb.WriteRune(utf8.RuneError)
				i++
			} else if !cont(i+2, 0x80, 0xBF) {
				sb.WriteRune(utf8.RuneError)
				i += 2
			} else if b1 == 0xED && b[i+1] >= 0xA0 {
				sb.WriteRune(utf8.RuneError)
				i += 3
			} else {
				sb.Write(b[i : i+3])
				i += 3
			}
		case b1 >= 0xF0 && b1 <= 0xF4:
			lo, hi := byte(0x80), byte(0xBF)
			if b1 == 0xF0 {
				lo = 0x90
			} else if b1 == 0xF4 {
				hi = 0x8F
			}
			if !cont(i+1, lo, hi) {
				sb.WriteRune(utf8.RuneError)
				i++
			} else if !cont(i+2, 0x80, 0xBF) {
				sb.WriteRune(utf8.RuneError)
				i += 2
			} else if !cont(i+3, 0x80, 0xBF) {
				sb.WriteRune(utf8.RuneError)
				i += 3
			} else {
				sb.Write(b[i : i+4])
				i += 4
			}
		default:
			sb.WriteRune(utf8.RuneError)
			i++
		}
	}
	return sb.String()
}

func (utf8Ops) canEncodeChar(uint16) bool { return true }

type asciiOps struct{}

func appendASCIINonPure(out []byte, s string) []byte {
	units := javalang.StringToUTF16(s)
	for i := 0; i < len(units); i++ {
		c := units[i]
		if c < 0x80 {
			out = append(out, byte(c))
		} else {
			if c >= 0xD800 && c <= 0xDBFF && i+1 < len(units) && units[i+1] >= 0xDC00 && units[i+1] <= 0xDFFF {
				i++
			}
			out = append(out, '?')
		}
	}
	return out
}

func (asciiOps) encode(s string) []byte {
	if javalang.IsASCII(s) {
		return []byte(s)
	}
	return appendASCIINonPure(make([]byte, 0, len(s)), s)
}

func (asciiOps) decode(b []byte) (string, error) {
	if isBytesASCII(b) {
		return string(b), nil
	}
	return "", &MalformedInputException{InputLength: 1}
}

func (asciiOps) decodeReplace(b []byte) string {
	if isBytesASCII(b) {
		return string(b)
	}
	var sb strings.Builder
	sb.Grow(len(b))
	for _, x := range b {
		if x < 0x80 {
			sb.WriteByte(x)
		} else {
			sb.WriteRune(utf8.RuneError)
		}
	}
	return sb.String()
}

func (asciiOps) canEncodeChar(ch uint16) bool { return ch < 0x80 }

type iso88591Ops struct{}

func appendISO88591NonPure(out []byte, s string) []byte {
	units := javalang.StringToUTF16(s)
	for i := 0; i < len(units); i++ {
		c := units[i]
		if c <= 0xFF {
			out = append(out, byte(c))
		} else {
			if c >= 0xD800 && c <= 0xDBFF && i+1 < len(units) && units[i+1] >= 0xDC00 && units[i+1] <= 0xDFFF {
				i++
			}
			out = append(out, '?')
		}
	}
	return out
}

func (iso88591Ops) encode(s string) []byte {
	if javalang.IsASCII(s) {
		return []byte(s)
	}
	return appendISO88591NonPure(make([]byte, 0, len(s)), s)
}

func (iso88591Ops) decode(b []byte) (string, error) {
	return iso88591Ops{}.decodeReplace(b), nil
}

func decodeISO88591NonASCII(b []byte) string {
	nonASCII := 0
	for _, x := range b {
		if x >= 0x80 {
			nonASCII++
		}
	}
	var sb strings.Builder
	sb.Grow(len(b) + nonASCII)
	start := 0
	for i, x := range b {
		if x >= 0x80 {
			if i > start {
				sb.Write(b[start:i])
			}
			sb.WriteByte(0xC0 | (x >> 6))
			sb.WriteByte(0x80 | (x & 0x3F))
			start = i + 1
		}
	}
	if start < len(b) {
		sb.Write(b[start:])
	}
	return sb.String()
}

func (iso88591Ops) decodeReplace(b []byte) string {
	if isBytesASCII(b) {
		return string(b)
	}
	return decodeISO88591NonASCII(b)
}

func (iso88591Ops) canEncodeChar(ch uint16) bool { return ch <= 0xFF }

type utf16Ops struct {
	bigEndian bool
	marked    bool
}

func (o utf16Ops) encode(s string) []byte {
	units := javalang.StringToUTF16(s)
	if len(units) == 0 {
		return []byte{}
	}
	out := make([]byte, 0, len(units)*2+2)
	if o.marked {
		out = append(out, 0xFE, 0xFF)
	}
	put := func(u uint16) {
		if o.bigEndian {
			out = append(out, byte(u>>8), byte(u))
		} else {
			out = append(out, byte(u), byte(u>>8))
		}
	}
	for i := 0; i < len(units); i++ {
		u := units[i]
		if u >= 0xD800 && u <= 0xDBFF {
			if i+1 < len(units) && units[i+1] >= 0xDC00 && units[i+1] <= 0xDFFF {
				put(u)
				put(units[i+1])
				i++
			} else {
				put(0xFFFD)
			}
		} else if u >= 0xDC00 && u <= 0xDFFF {
			put(0xFFFD)
		} else {
			put(u)
		}
	}
	return out
}

func (o utf16Ops) decodeImpl(b []byte, strict bool) (string, error) {
	be := o.bigEndian
	if o.marked && len(b) >= 2 {
		if b[0] == 0xFE && b[1] == 0xFF {
			b = b[2:]
		} else if b[0] == 0xFF && b[1] == 0xFE {
			be = false
			b = b[2:]
		}
	}
	unit := func(i int) uint16 {
		if be {
			return uint16(b[i])<<8 | uint16(b[i+1])
		}
		return uint16(b[i+1])<<8 | uint16(b[i])
	}
	units := make([]uint16, 0, len(b)/2+1)
	for i := 0; i < len(b); {
		rest := len(b) - i
		if rest < 2 {
			if strict {
				return "", &MalformedInputException{InputLength: rest}
			}
			units = append(units, 0xFFFD)
			break
		}
		u := unit(i)
		if u >= 0xDC00 && u <= 0xDFFF {
			if strict {
				return "", &MalformedInputException{InputLength: 2}
			}
			units = append(units, 0xFFFD)
			i += 2
			continue
		}
		if u >= 0xD800 && u <= 0xDBFF {
			if rest >= 4 {
				v := unit(i + 2)
				if v >= 0xDC00 && v <= 0xDFFF {
					units = append(units, u, v)
					i += 4
					continue
				}
				if strict {
					return "", &MalformedInputException{InputLength: 4}
				}
				units = append(units, 0xFFFD)
				i += 4
				continue
			}
			if strict {
				return "", &MalformedInputException{InputLength: rest}
			}
			units = append(units, 0xFFFD)
			break
		}
		units = append(units, u)
		i += 2
	}
	return javalang.StringFromUTF16(units), nil
}

func (o utf16Ops) decode(b []byte) (string, error) { return o.decodeImpl(b, true) }
func (o utf16Ops) decodeReplace(b []byte) string {
	s, _ := o.decodeImpl(b, false)
	return s
}
func (utf16Ops) canEncodeChar(uint16) bool { return true }

type xtextOps struct {
	enc         encoding.Encoding
	isoControls bool
}

func newXTextOps(enc encoding.Encoding, isoControls bool) func() charsetOps {
	return func() charsetOps { return xtextOps{enc: enc, isoControls: isoControls} }
}

func (o xtextOps) encode(s string) []byte {
	valid := string(utf8Ops{}.encode(s))
	enc := encoding.ReplaceUnsupported(o.enc.NewEncoder())
	out, err := enc.Bytes([]byte(valid))
	if err != nil {
		return []byte(strings.Repeat("?", utf8.RuneCountInString(valid)))
	}
	return out
}

func (o xtextOps) decode(b []byte) (string, error) {
	if cm, ok := o.enc.(*charmap.Charmap); ok && o.isoControls {
		units := make([]uint16, len(b))
		for i, x := range b {
			if x >= 0x80 && x < 0xA0 {
				units[i] = uint16(x)
				continue
			}
			r := cm.DecodeByte(x)
			if r == utf8.RuneError {
				return "", &UnmappableCharacterException{InputLength: 1}
			}
			units[i] = uint16(r)
		}
		return javalang.StringFromUTF16(units), nil
	}
	out, err := o.enc.NewDecoder().Bytes(b)
	if err != nil {
		return "", &MalformedInputException{InputLength: 1}
	}
	if !utf8.Valid(out) || strings.ContainsRune(string(out), utf8.RuneError) {
		// Check whether input literally had U+FFFD or failed mapping.
		return string(out), nil
	}
	return string(out), nil
}

func (o xtextOps) decodeReplace(b []byte) string {
	if cm, ok := o.enc.(*charmap.Charmap); ok {
		units := make([]uint16, len(b))
		for i, x := range b {
			if o.isoControls && x >= 0x80 && x < 0xA0 {
				units[i] = uint16(x)
				continue
			}
			r := cm.DecodeByte(x)
			if r == utf8.RuneError {
				units[i] = 0xFFFD
			} else {
				units[i] = uint16(r)
			}
		}
		return javalang.StringFromUTF16(units)
	}
	out, _ := o.enc.NewDecoder().Bytes(b)
	return utf8Ops{}.decodeReplace(out)
}

func (o xtextOps) canEncodeChar(ch uint16) bool {
	if cm, ok := o.enc.(*charmap.Charmap); ok {
		if o.isoControls && ch >= 0x80 && ch < 0xA0 {
			return true
		}
		_, ok := cm.EncodeRune(rune(ch))
		return ok
	}
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], rune(ch))
	_, err := o.enc.NewEncoder().Bytes(buf[:n])
	return err == nil
}

// --- StandardCharsets & Registry ---

func newCharset(name string, build func() charsetOps, aliases ...string) *Charset {
	return &Charset{name: name, aliases: aliases, build: build}
}

var (
	StandardCharsets_US_ASCII   = newCharset("US-ASCII", func() charsetOps { return asciiOps{} }, "646", "ANSI_X3.4-1968", "ANSI_X3.4-1986", "ASCII", "IBM367", "ISO646-US", "ISO_646.irv:1991", "ascii7", "cp367", "csASCII", "iso-ir-6", "iso_646.irv:1983", "us")
	StandardCharsets_ISO_8859_1 = newCharset("ISO-8859-1", func() charsetOps { return iso88591Ops{} }, "819", "8859_1", "IBM-819", "IBM819", "ISO8859-1", "ISO8859_1", "ISO_8859-1", "ISO_8859-1:1987", "ISO_8859_1", "cp819", "csISOLatin1", "iso-ir-100", "l1", "latin1")
	StandardCharsets_UTF_8      = newCharset("UTF-8", func() charsetOps { return utf8Ops{} }, "UTF8", "unicode-1-1-utf-8")
	StandardCharsets_UTF_16     = newCharset("UTF-16", func() charsetOps { return utf16Ops{bigEndian: true, marked: true} }, "UTF_16", "UnicodeBig", "unicode", "utf16")
	StandardCharsets_UTF_16BE   = newCharset("UTF-16BE", func() charsetOps { return utf16Ops{bigEndian: true} }, "ISO-10646-UCS-2", "UTF_16BE", "UnicodeBigUnmarked", "X-UTF-16BE")
	StandardCharsets_UTF_16LE   = newCharset("UTF-16LE", func() charsetOps { return utf16Ops{bigEndian: false} }, "UTF_16LE", "UnicodeLittleUnmarked", "X-UTF-16LE")
)

var supportedCharsets = []*Charset{
	StandardCharsets_US_ASCII,
	StandardCharsets_ISO_8859_1,
	StandardCharsets_UTF_8,
	StandardCharsets_UTF_16,
	StandardCharsets_UTF_16BE,
	StandardCharsets_UTF_16LE,

	newCharset("windows-1250", newXTextOps(charmap.Windows1250, false), "cp1250", "cp5346"),
	newCharset("windows-1251", newXTextOps(charmap.Windows1251, false), "ansi-1251", "cp1251", "cp5347"),
	newCharset("windows-1252", newXTextOps(charmap.Windows1252, false), "cp1252", "cp5348", "ibm-1252", "ibm1252"),
	newCharset("windows-1253", newXTextOps(charmap.Windows1253, false), "cp1253", "cp5349"),
	newCharset("windows-1254", newXTextOps(charmap.Windows1254, false), "cp1254", "cp5350"),
	newCharset("windows-1255", newXTextOps(charmap.Windows1255, false), "cp1255"),
	newCharset("windows-1256", newXTextOps(charmap.Windows1256, false), "cp1256"),
	newCharset("windows-1257", newXTextOps(charmap.Windows1257, false), "cp1257", "cp5353"),
	newCharset("windows-1258", newXTextOps(charmap.Windows1258, false), "cp1258"),
	newCharset("x-windows-874", newXTextOps(charmap.Windows874, false), "ms-874", "ms874", "windows-874"),

	newCharset("ISO-8859-2", newXTextOps(charmap.ISO8859_2, true), "8859_2", "ISO8859-2", "ISO_8859-2", "iso8859_2", "l2", "latin2"),
	newCharset("ISO-8859-3", newXTextOps(charmap.ISO8859_3, true), "8859_3", "ISO8859-3", "ISO_8859-3", "iso8859_3", "l3", "latin3"),
	newCharset("ISO-8859-4", newXTextOps(charmap.ISO8859_4, true), "8859_4", "ISO8859-4", "ISO_8859-4", "iso8859_4", "l4", "latin4"),
	newCharset("ISO-8859-5", newXTextOps(charmap.ISO8859_5, true), "8859_5", "ISO8859-5", "ISO_8859-5", "iso8859_5", "cyrillic"),
	newCharset("ISO-8859-6", newXTextOps(charmap.ISO8859_6, true), "8859_6", "ISO8859-6", "ISO_8859-6", "iso8859_6", "arabic"),
	newCharset("ISO-8859-7", newXTextOps(charmap.ISO8859_7, true), "8859_7", "ISO8859-7", "ISO_8859-7", "iso8859_7", "greek"),
	newCharset("ISO-8859-8", newXTextOps(charmap.ISO8859_8, true), "8859_8", "ISO8859-8", "ISO_8859-8", "iso8859_8", "hebrew"),
	newCharset("ISO-8859-9", newXTextOps(charmap.ISO8859_9, true), "8859_9", "ISO8859-9", "ISO_8859-9", "iso8859_9", "l5", "latin5"),
	newCharset("ISO-8859-13", newXTextOps(charmap.ISO8859_13, true), "8859_13", "ISO8859-13", "iso8859_13"),
	newCharset("ISO-8859-15", newXTextOps(charmap.ISO8859_15, true), "8859_15", "ISO-8859-15", "ISO8859-15", "ISO8859_15", "L9", "LATIN9", "Latin-9"),
	newCharset("ISO-8859-16", newXTextOps(charmap.ISO8859_16, true), "ISO8859_16", "ISO_8859-16", "l10", "latin10"),

	newCharset("KOI8-R", newXTextOps(charmap.KOI8R, false), "cskoi8r", "koi8", "koi8_r"),
	newCharset("KOI8-U", newXTextOps(charmap.KOI8U, false), "koi8_u"),
	newCharset("IBM437", newXTextOps(charmap.CodePage437, false), "437", "cp437", "ibm437"),
	newCharset("IBM850", newXTextOps(charmap.CodePage850, false), "850", "cp850", "ibm850"),
	newCharset("IBM866", newXTextOps(charmap.CodePage866, false), "866", "cp866", "ibm866"),
	newCharset("x-MacRoman", newXTextOps(charmap.Macintosh, false), "MacRoman"),

	newCharset("Shift_JIS", newXTextOps(japanese.ShiftJIS, false), "csShiftJIS", "ms_kanji", "shift-jis", "shift_jis", "sjis", "x-sjis"),
	newCharset("EUC-JP", newXTextOps(japanese.EUCJP, false), "euc_jp", "eucjis", "eucjp", "x-euc-jp"),
	newCharset("ISO-2022-JP", newXTextOps(japanese.ISO2022JP, false), "csISO2022JP", "iso2022jp", "jis"),
	newCharset("GBK", newXTextOps(simplifiedchinese.GBK, false), "CP936", "windows-936"),
	newCharset("GB2312", newXTextOps(simplifiedchinese.HZGB2312, false), "EUC_CN", "euc-cn", "euccn", "gb2312"),
	newCharset("GB18030", newXTextOps(simplifiedchinese.GB18030, false), "gb18030-2022"),
	newCharset("Big5", newXTextOps(traditionalchinese.Big5, false), "csBig5"),
	newCharset("EUC-KR", newXTextOps(korean.EUCKR, false), "5601", "csEUCKR", "euc_kr", "euckr", "ksc5601"),
}

var (
	charsetIndexOnce sync.Once
	charsetIndex     map[string]*Charset
)

func lookupCharset(name string) *Charset {
	charsetIndexOnce.Do(func() {
		charsetIndex = make(map[string]*Charset, len(supportedCharsets)*4)
		for _, cs := range supportedCharsets {
			charsetIndex[javalang.StringToLowerCaseASCII(cs.name)] = cs
			for _, alias := range cs.aliases {
				charsetIndex[javalang.StringToLowerCaseASCII(alias)] = cs
			}
		}
	})
	return charsetIndex[javalang.StringToLowerCaseASCII(name)]
}

func requireLegalCharsetName(name string) error {
	if name == "" {
		return &IllegalCharsetNameException{CharsetName: name}
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			continue
		}
		if i > 0 && (c == '-' || c == '+' || c == '.' || c == ':' || c == '_') {
			continue
		}
		return &IllegalCharsetNameException{CharsetName: name}
	}
	return nil
}

func Charset_DefaultCharset() *Charset {
	return StandardCharsets_UTF_8
}

func fastStandardCharset(charsetName string) *Charset {
	switch charsetName {
	case "UTF-8", "utf-8", "UTF8", "utf8":
		return StandardCharsets_UTF_8
	case "US-ASCII", "us-ascii", "ASCII", "ascii":
		return StandardCharsets_US_ASCII
	case "ISO-8859-1", "iso-8859-1", "latin1", "ISO8859_1":
		return StandardCharsets_ISO_8859_1
	default:
		return nil
	}
}

func Charset_IsSupported(charsetName string) bool {
	if fastStandardCharset(charsetName) != nil {
		return true
	}
	if err := requireLegalCharsetName(charsetName); err != nil {
		panic(err)
	}
	return lookupCharset(charsetName) != nil
}

func Charset_ForName(charsetName string) *Charset {
	if cs := fastStandardCharset(charsetName); cs != nil {
		return cs
	}
	if err := requireLegalCharsetName(charsetName); err != nil {
		panic(err)
	}
	if cs := lookupCharset(charsetName); cs != nil {
		return cs
	}
	panic(&UnsupportedCharsetException{CharsetName: charsetName})
}
