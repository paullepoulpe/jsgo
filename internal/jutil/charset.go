package jutil

import (
	"strings"
	"sync"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

// Go-only: language-adapter — java.nio.charset.Charset for supportedCharsets,
// and String.getBytes(UTF_8) and new String(byte[], UTF_8). Each charset's
// conversions (charsetCoder) are built on first use.

// Charset is java.nio.charset.Charset.
type Charset struct {
	name    string
	aliases []string
	build   func() charsetCoder
	once    sync.Once
	coder   charsetCoder
}

// impl returns the charset's conversions, building them on first use.
func (c *Charset) impl() charsetCoder {
	c.once.Do(func() { c.coder = c.build() })
	return c.coder
}

// UnsupportedCharsetException is java.nio.charset.UnsupportedCharsetException.
type UnsupportedCharsetException struct{ IllegalArgumentException }

// NewUnsupportedCharsetException is new UnsupportedCharsetException(name).
func NewUnsupportedCharsetException(charsetName string) *UnsupportedCharsetException {
	e := &UnsupportedCharsetException{}
	InitThrowable(e, "java.nio.charset.UnsupportedCharsetException", msgPtr(charsetName), nil)
	return e
}

// charsetLookup finds a supported charset by canonical name or alias,
// ignoring ASCII case (see the Javadoc of java.nio.charset.Charset), or
// returns nil.
func charsetLookup(charsetName string) *Charset {
	charsetIndexOnce.Do(buildCharsetIndex)
	return charsetIndex[StringToLowerCaseASCII(charsetName)]
}

var (
	charsetIndexOnce sync.Once
	charsetIndex     map[string]*Charset
)

func buildCharsetIndex() {
	charsetIndex = map[string]*Charset{}
	for _, c := range supportedCharsets {
		charsetIndex[StringToLowerCaseASCII(c.name)] = c
		for _, a := range c.aliases {
			charsetIndex[StringToLowerCaseASCII(a)] = c
		}
	}
}

// Charset_ForName is Charset.forName(String): it panics with
// IllegalCharsetNameException for an illegal name and with
// UnsupportedCharsetException (message: the name) for a charset that is not
// in charset_registry.go.
func Charset_ForName(charsetName string) *Charset {
	requireLegalCharsetName(charsetName)
	if c := charsetLookup(charsetName); c != nil {
		return c
	}
	panic(NewUnsupportedCharsetException(charsetName))
}

// Name is Charset.name().
func (c *Charset) Name() string { return c.name }

// String is Charset.toString().
func (c *Charset) String() string { return c.name }

// CanEncodeChar is charset.newEncoder().canEncode(char). A lone surrogate
// can never be encoded on its own.
func (c *Charset) CanEncodeChar(ch uint16) bool {
	if isSurrogate(ch) {
		return false
	}
	return c.impl().canEncode(ch)
}

// StringGetBytesUTF8 is String.getBytes(StandardCharsets.UTF_8): unpaired
// surrogates become '?'.
func StringGetBytesUTF8(s string) []byte {
	if utf8.ValidString(s) {
		return []byte(s)
	}
	return AppendStringBytesUTF8(make([]byte, 0, len(s)), s)
}

// AppendStringBytesUTF8 appends the bytes of String.getBytes(StandardCharsets.UTF_8)
// to dst without allocating an intermediate []byte when s is valid UTF-8.
func AppendStringBytesUTF8(dst []byte, s string) []byte {
	if utf8.ValidString(s) {
		return append(dst, s...)
	}
	units := StringToUTF16(s)
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
		case isSurrogate(c):
			dst = append(dst, '?')
		default:
			dst = append(dst, byte(0xE0|c>>12), byte(0x80|(c>>6)&0x3F), byte(0x80|c&0x3F))
		}
	}
	return dst
}

// NewStringFromUTF8Bytes is new String(byte[], StandardCharsets.UTF_8):
// each malformed subsequence becomes U+FFFD, as Java's decoder replaces it.
func NewStringFromUTF8Bytes(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	var sb strings.Builder
	sb.Grow(len(b))
	n := len(b)
	i := 0
	cont := func(j int, lo, hi byte) bool { return j < n && b[j] >= lo && b[j] <= hi }
	for i < n {
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
			// Java accepts ED A0..BF as a well-formed prefix: a complete
			// encoded surrogate is one malformed char (one U+FFFD).
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

// charsetCoder holds the conversions of one charset.
type charsetCoder interface {
	decode(b []byte, report bool) (string, codingError)
	canEncode(c uint16) bool
	newEncoder() charsetEncoder
}

// codingError describes a malformed or unmappable input sequence, as a
// CoderResult does (length 0: no error).
type codingError struct {
	length     int
	unmappable bool
}

// charsetStepper decodes one character (or skips a sequence that only
// changes state) at a time.
type charsetStepper interface {
	step(out []uint16, b []byte, i int) ([]uint16, int, codingError)
}

// charsetEncoder encodes code points in REPLACE mode, keeping any shift
// state between calls.
type charsetEncoder interface {
	encodeRune(dst []byte, r rune) []byte
	encodeMalformed(dst []byte) []byte
	finish(dst []byte) []byte
}

func decodeSteps(d charsetStepper, b []byte, report bool) (string, codingError) {
	out := make([]uint16, 0, len(b))
	for i := 0; i < len(b); {
		var n int
		var e codingError
		out, n, e = d.step(out, b, i)
		if e.length > 0 {
			if report {
				return "", e
			}
			out = append(out, 0xFFFD)
			i += e.length
			continue
		}
		i += n
	}
	return StringFromUTF16(out), codingError{}
}

func malformed(n int) codingError  { return codingError{length: n} }
func unmappable(n int) codingError { return codingError{length: n, unmappable: true} }

// charTable maps BMP chars to uint32 values (0: absent), in pages of 256.
type charTable struct{ pages [256]*[256]uint32 }

func (t *charTable) get(c uint16) uint32 {
	if p := t.pages[c>>8]; p != nil {
		return p[c&0xFF]
	}
	return 0
}

func (t *charTable) set(c uint16, v uint32) {
	p := t.pages[c>>8]
	if p == nil {
		p = new([256]uint32)
		t.pages[c>>8] = p
	}
	p[c&0xFF] = v
}

func (t *charTable) setIfAbsent(c uint16, v uint32) {
	if t.get(c) == 0 {
		t.set(c, v)
	}
}

func packBytes(b ...byte) uint32 {
	var v uint32
	for _, x := range b {
		v = v<<8 | uint32(x)
	}
	return uint32(len(b))<<24 | v
}

func appendPacked(dst []byte, v uint32) []byte {
	switch v >> 24 {
	case 1:
		return append(dst, byte(v))
	case 2:
		return append(dst, byte(v>>8), byte(v))
	default:
		return append(dst, byte(v>>16), byte(v>>8), byte(v))
	}
}

type tableEncoder struct{ table *charTable }

func (e tableEncoder) encodeRune(dst []byte, r rune) []byte {
	if r < 0x10000 {
		if v := e.table.get(uint16(r)); v != 0 {
			return appendPacked(dst, v)
		}
	}
	return append(dst, '?')
}

func (tableEncoder) encodeMalformed(dst []byte) []byte { return append(dst, '?') }
func (tableEncoder) finish(dst []byte) []byte          { return dst }

// Decode is new String(bytes, charset): each malformed or unmappable input
// sequence is replaced with U+FFFD.
func (c *Charset) Decode(b []byte) string {
	s, _ := c.impl().decode(b, false)
	return s
}

// DecodeStrict decodes b and panics at the first malformed input sequence
// with MalformedInputException, or at the first unmappable one with
// UnmappableCharacterException (CodingErrorAction.REPORT).
func (c *Charset) DecodeStrict(b []byte) string {
	s, err := c.TryDecodeStrict(b)
	if err != nil {
		panic(err)
	}
	return s
}

// TryDecodeStrict decodes b and returns MalformedInputException or
// UnmappableCharacterException instead of panicking on invalid input.
func (c *Charset) TryDecodeStrict(b []byte) (string, IIOException) {
	s, e := c.impl().decode(b, true)
	switch {
	case e.length == 0:
		return s, nil
	case e.unmappable:
		return "", NewUnmappableCharacterException(e.length)
	default:
		return "", NewMalformedInputException(e.length)
	}
}

type utf8Coder struct{}

func (utf8Coder) decode(b []byte, report bool) (string, codingError) {
	if utf8.Valid(b) {
		return string(b), codingError{}
	}
	if report {
		return "", malformed(charsetUTF8MalformedLength(b))
	}
	return NewStringFromUTF8Bytes(b), codingError{}
}

func (utf8Coder) canEncode(uint16) bool      { return true }
func (utf8Coder) newEncoder() charsetEncoder { return utf8Encoder{} }

type utf8Encoder struct{}

func (utf8Encoder) encodeRune(dst []byte, r rune) []byte { return utf8.AppendRune(dst, r) }
func (utf8Encoder) encodeMalformed(dst []byte) []byte    { return append(dst, '?') }
func (utf8Encoder) finish(dst []byte) []byte             { return dst }

func charsetUTF8MalformedLength(b []byte) int {
	for i := 0; i < len(b); {
		r, size := utf8.DecodeRune(b[i:])
		if r != utf8.RuneError || size > 1 {
			i += size
			continue
		}
		b1 := b[i]
		n := len(b) - i
		cont := func(j int, lo, hi byte) bool { return j < n && b[i+j] >= lo && b[i+j] <= hi }
		switch {
		case b1 >= 0xE0 && b1 <= 0xEF:
			lo := byte(0x80)
			if b1 == 0xE0 {
				lo = 0xA0
			}
			if !cont(1, lo, 0xBF) {
				return 1
			}
			if b1 == 0xED && b[i+1] >= 0xA0 && cont(2, 0x80, 0xBF) {
				return 3
			}
			return 2
		case b1 >= 0xF0 && b1 <= 0xF4:
			lo, hi := byte(0x80), byte(0xBF)
			if b1 == 0xF0 {
				lo = 0x90
			} else if b1 == 0xF4 {
				hi = 0x8F
			}
			if !cont(1, lo, hi) {
				return 1
			}
			if !cont(2, 0x80, 0xBF) {
				return 2
			}
			return 3
		}
		return 1
	}
	return 1
}

type utf16Coder struct {
	bigEndian bool
	marked    bool
}

func (c utf16Coder) canEncode(uint16) bool      { return true }
func (c utf16Coder) newEncoder() charsetEncoder { return &utf16Encoder{coder: c} }

func (c utf16Coder) decode(b []byte, report bool) (string, codingError) {
	bigEndian := c.bigEndian
	if c.marked && len(b) >= 2 {
		switch {
		case b[0] == 0xFE && b[1] == 0xFF:
			b = b[2:]
		case b[0] == 0xFF && b[1] == 0xFE:
			bigEndian = false
			b = b[2:]
		}
	}
	unit := func(i int) uint16 {
		if bigEndian {
			return uint16(b[i])<<8 | uint16(b[i+1])
		}
		return uint16(b[i+1])<<8 | uint16(b[i])
	}
	units := make([]uint16, 0, len(b)/2+1)
	for i := 0; i < len(b); {
		bad := 0
		switch rest := len(b) - i; {
		case rest < 2:
			bad = rest
		default:
			u := unit(i)
			switch {
			case isLowSurrogateUnit(u):
				bad = 2
			case !isHighSurrogateUnit(u):
				units = append(units, u)
				i += 2
				continue
			case rest < 4:
				bad = rest
			default:
				if v := unit(i + 2); isLowSurrogateUnit(v) {
					units = append(units, u, v)
					i += 4
					continue
				}
				bad = 4
			}
		}
		if report {
			return "", malformed(bad)
		}
		units = append(units, 0xFFFD)
		i += bad
	}
	return StringFromUTF16(units), codingError{}
}

type utf16Encoder struct {
	coder   utf16Coder
	started bool
}

func (e *utf16Encoder) unit(dst []byte, u uint16) []byte {
	if e.coder.marked && !e.started {
		dst = append(dst, 0xFE, 0xFF)
	}
	e.started = true
	if e.coder.bigEndian {
		return append(dst, byte(u>>8), byte(u))
	}
	return append(dst, byte(u), byte(u>>8))
}

func (e *utf16Encoder) encodeRune(dst []byte, r rune) []byte {
	if r >= 0x10000 {
		r -= 0x10000
		dst = e.unit(dst, uint16(0xD800+r>>10))
		return e.unit(dst, uint16(0xDC00+r&0x3FF))
	}
	return e.unit(dst, uint16(r))
}

func (e *utf16Encoder) encodeMalformed(dst []byte) []byte { return e.unit(dst, 0xFFFD) }
func (e *utf16Encoder) finish(dst []byte) []byte          { return dst }

type utf32Coder struct {
	bigEndian bool
	detect    bool
	writeMark bool
}

func (c utf32Coder) decode(b []byte, report bool) (string, codingError) {
	return decodeSteps(&utf32Decoder{bigEndian: c.bigEndian, detect: c.detect}, b, report)
}

func (utf32Coder) canEncode(uint16) bool        { return true }
func (c utf32Coder) newEncoder() charsetEncoder { return &utf32Encoder{coder: c} }

type utf32Decoder struct {
	bigEndian bool
	detect    bool
}

func (d *utf32Decoder) step(out []uint16, b []byte, i int) ([]uint16, int, codingError) {
	if rest := len(b) - i; rest < 4 {
		return out, 0, malformed(rest)
	}
	be := uint32(b[i])<<24 | uint32(b[i+1])<<16 | uint32(b[i+2])<<8 | uint32(b[i+3])
	le := uint32(b[i+3])<<24 | uint32(b[i+2])<<16 | uint32(b[i+1])<<8 | uint32(b[i])
	if i == 0 {
		switch {
		case d.detect && be == 0xFEFF:
			return out, 4, codingError{}
		case d.detect && le == 0xFEFF:
			d.bigEndian = false
			return out, 4, codingError{}
		case !d.detect && d.bigEndian && be == 0xFEFF, !d.detect && !d.bigEndian && le == 0xFEFF:
			return out, 4, codingError{}
		}
	}
	v := le
	if d.bigEndian {
		v = be
	}
	switch {
	case v < 0x10000:
		return append(out, uint16(v)), 4, codingError{}
	case v <= 0x10FFFF:
		v -= 0x10000
		return append(out, uint16(0xD800+v>>10), uint16(0xDC00+v&0x3FF)), 4, codingError{}
	}
	return out, 0, malformed(4)
}

type utf32Encoder struct {
	coder   utf32Coder
	started bool
}

func (e *utf32Encoder) unit(dst []byte, v rune) []byte {
	if e.coder.writeMark && !e.started {
		dst = e.raw(dst, 0xFEFF)
	}
	e.started = true
	return e.raw(dst, v)
}

func (e *utf32Encoder) raw(dst []byte, v rune) []byte {
	if e.coder.bigEndian {
		return append(dst, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	return append(dst, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

func (e *utf32Encoder) encodeRune(dst []byte, r rune) []byte { return e.unit(dst, r) }
func (e *utf32Encoder) encodeMalformed(dst []byte) []byte    { return e.unit(dst, 0xFFFD) }
func (e *utf32Encoder) finish(dst []byte) []byte             { return dst }

// IllegalCharsetNameException is
// java.nio.charset.IllegalCharsetNameException: thrown when a string that
// is not a legal charset name is used as such. Its message is the name.
type IllegalCharsetNameException struct {
	IllegalArgumentException
	charsetName string
}

func (*IllegalCharsetNameException) illegalCharsetNameException() {}

// IIllegalCharsetNameException is the catch marker for
// IllegalCharsetNameException.
type IIllegalCharsetNameException interface {
	IIllegalArgumentException
	illegalCharsetNameException()
}

// NewIllegalCharsetNameException is new IllegalCharsetNameException(String).
func NewIllegalCharsetNameException(charsetName string) *IllegalCharsetNameException {
	e := &IllegalCharsetNameException{charsetName: charsetName}
	InitThrowable(e, "java.nio.charset.IllegalCharsetNameException", &charsetName, nil)
	return e
}

// GetCharsetName is IllegalCharsetNameException.getCharsetName().
func (e *IllegalCharsetNameException) GetCharsetName() string { return e.charsetName }

// requireLegalCharsetName applies the Javadoc's rules for legal charset names:
// non-empty, made of the letters A-Z and a-z, the digits 0-9, '-', '+',
// '.', ':' and '_', and beginning with a letter or a digit.
func requireLegalCharsetName(name string) {
	if name == "" {
		panic(NewIllegalCharsetNameException(name))
	}
	for i := 0; i < len(name); i++ {
		c := name[i]
		letterOrDigit := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
		if letterOrDigit {
			continue
		}
		if i > 0 && (c == '-' || c == '+' || c == '.' || c == ':' || c == '_') {
			continue
		}
		panic(NewIllegalCharsetNameException(name))
	}
}

// Charset_IsSupported is Charset.isSupported(String): whether the named
// charset is supported (supportedCharsets), throwing
// IllegalCharsetNameException for an illegal name.
func Charset_IsSupported(charsetName string) bool {
	requireLegalCharsetName(charsetName)
	return charsetLookup(charsetName) != nil
}

// Go-only: language-adapter — the charsets jutil supports, with the
// canonical names and aliases the JDK reports for them (Charset.name() and
// Charset.aliases(), observed black-box on JDK 21; see
// tools/oracle/java/JutilCharsetVectors.java). Charset.forName finds a
// charset by any of these names, ignoring case. The JDK supports more
// charsets (EBCDIC, the IBM and Mac code pages not listed here, ...); for
// those, Charset.isSupported returns false and Charset.forName throws
// UnsupportedCharsetException, as a JDK without them would.

var (
	StandardCharsets_US_ASCII   = newCharset("US-ASCII", newUSASCII, "646", "ANSI_X3.4-1968", "ANSI_X3.4-1986", "ASCII", "IBM367", "ISO646-US", "ISO_646.irv:1991", "ascii7", "cp367", "csASCII", "iso-ir-6", "iso_646.irv:1983", "us")
	StandardCharsets_ISO_8859_1 = newCharset("ISO-8859-1", newSBCS(charmap.ISO8859_1, false), "819", "8859_1", "IBM-819", "IBM819", "ISO8859-1", "ISO8859_1", "ISO_8859-1", "ISO_8859-1:1987", "ISO_8859_1", "cp819", "csISOLatin1", "iso-ir-100", "l1", "latin1")
	StandardCharsets_UTF_8      = newCharset("UTF-8", func() charsetCoder { return utf8Coder{} }, "UTF8", "unicode-1-1-utf-8")
	StandardCharsets_UTF_16     = newCharset("UTF-16", func() charsetCoder { return utf16Coder{bigEndian: true, marked: true} }, "UTF_16", "UnicodeBig", "unicode", "utf16")
	StandardCharsets_UTF_16BE   = newCharset("UTF-16BE", func() charsetCoder { return utf16Coder{bigEndian: true} }, "ISO-10646-UCS-2", "UTF_16BE", "UnicodeBigUnmarked", "X-UTF-16BE")
	StandardCharsets_UTF_16LE   = newCharset("UTF-16LE", func() charsetCoder { return utf16Coder{} }, "UTF_16LE", "UnicodeLittleUnmarked", "X-UTF-16LE")
)

func newCharset(name string, build func() charsetCoder, aliases ...string) *Charset {
	return &Charset{name: name, aliases: aliases, build: build}
}

// supportedCharsets are the charsets Charset.forName finds.
var supportedCharsets = []*Charset{
	StandardCharsets_US_ASCII, StandardCharsets_ISO_8859_1, StandardCharsets_UTF_8,
	StandardCharsets_UTF_16, StandardCharsets_UTF_16BE, StandardCharsets_UTF_16LE,

	newCharset("UTF-32", func() charsetCoder { return utf32Coder{bigEndian: true, detect: true} }, "UTF32", "UTF_32"),
	newCharset("UTF-32BE", func() charsetCoder { return utf32Coder{bigEndian: true} }, "UTF_32BE", "X-UTF-32BE"),
	newCharset("UTF-32LE", func() charsetCoder { return utf32Coder{} }, "UTF_32LE", "X-UTF-32LE"),
	newCharset("X-UTF-32BE-BOM", func() charsetCoder { return utf32Coder{bigEndian: true, writeMark: true} }, "UTF-32BE-BOM", "UTF_32BE_BOM"),
	newCharset("X-UTF-32LE-BOM", func() charsetCoder { return utf32Coder{writeMark: true} }, "UTF-32LE-BOM", "UTF_32LE_BOM"),

	newCharset("windows-1250", newSBCS(charmap.Windows1250, false), "cp1250", "cp5346"),
	newCharset("windows-1251", newSBCS(charmap.Windows1251, false), "ansi-1251", "cp1251", "cp5347"),
	newCharset("windows-1252", newSBCS(charmap.Windows1252, false), "cp1252", "cp5348", "ibm-1252", "ibm1252"),
	newCharset("windows-1253", newSBCS(charmap.Windows1253, false), "cp1253", "cp5349"),
	newCharset("windows-1254", newSBCS(charmap.Windows1254, false), "cp1254", "cp5350"),
	newCharset("windows-1255", sbcsWindows1255, "cp1255"),
	newCharset("windows-1256", newSBCS(charmap.Windows1256, false), "cp1256"),
	newCharset("windows-1257", newSBCS(charmap.Windows1257, false), "cp1257", "cp5353"),
	newCharset("windows-1258", newSBCS(charmap.Windows1258, false), "cp1258"),
	newCharset("x-windows-874", newSBCS(charmap.Windows874, false), "ms-874", "ms874", "windows-874"),

	newCharset("ISO-8859-2", newSBCS(charmap.ISO8859_2, true), "8859_2", "912", "ISO8859-2", "ISO_8859-2", "ISO_8859-2:1987", "cp912", "csISOLatin2", "ibm-912", "ibm912", "iso-ir-101", "iso8859_2", "l2", "latin2"),
	newCharset("ISO-8859-3", newSBCS(charmap.ISO8859_3, true), "8859_3", "913", "ISO8859-3", "ISO_8859-3", "ISO_8859-3:1988", "cp913", "csISOLatin3", "ibm-913", "ibm913", "iso-ir-109", "iso8859_3", "l3", "latin3"),
	newCharset("ISO-8859-4", newSBCS(charmap.ISO8859_4, true), "8859_4", "914", "ISO_8859-4", "ISO_8859-4:1988", "cp914", "csISOLatin4", "ibm-914", "ibm914", "iso-ir-110", "iso8859-4", "iso8859_4", "l4", "latin4"),
	newCharset("ISO-8859-5", newSBCS(charmap.ISO8859_5, true), "8859_5", "915", "ISO8859-5", "ISO_8859-5", "ISO_8859-5:1988", "cp915", "csISOLatinCyrillic", "cyrillic", "ibm-915", "ibm915", "iso-ir-144", "iso8859_5"),
	newCharset("ISO-8859-6", newSBCS(charmap.ISO8859_6, true), "1089", "8859_6", "ASMO-708", "ECMA-114", "ISO8859-6", "ISO_8859-6", "ISO_8859-6:1987", "arabic", "cp1089", "csISOLatinArabic", "ibm-1089", "ibm1089", "iso-ir-127", "iso8859_6"),
	newCharset("ISO-8859-7", newSBCS(charmap.ISO8859_7, true), "813", "8859_7", "ECMA-118", "ELOT_928", "ISO_8859-7", "ISO_8859-7:1987", "cp813", "csISOLatinGreek", "greek", "greek8", "ibm-813", "ibm813", "iso-ir-126", "iso8859-7", "iso8859_7", "sun_eu_greek"),
	newCharset("ISO-8859-8", newSBCS(charmap.ISO8859_8, true), "8859_8", "916", "ISO8859-8", "ISO_8859-8", "ISO_8859-8:1988", "cp916", "csISOLatinHebrew", "hebrew", "ibm-916", "ibm916", "iso-ir-138", "iso8859_8"),
	newCharset("ISO-8859-9", newSBCS(charmap.ISO8859_9, true), "8859_9", "920", "ISO8859-9", "ISO_8859-9", "ISO_8859-9:1989", "cp920", "csISOLatin5", "ibm-920", "ibm920", "iso-ir-148", "iso8859_9", "l5", "latin5"),
	// x-iso-8859-11 is windows-874 with the C1 controls at 0x80 to 0x9F.
	newCharset("x-iso-8859-11", newSBCS(charmap.Windows874, true), "iso-8859-11", "iso8859_11"),
	newCharset("ISO-8859-13", newSBCS(charmap.ISO8859_13, true), "8859_13", "ISO8859-13", "iso8859_13", "iso_8859-13"),
	newCharset("ISO-8859-15", newSBCS(charmap.ISO8859_15, true), "8859_15", "923", "IBM-923", "IBM923", "ISO-8859-15", "ISO8859-15", "ISO8859_15", "ISO8859_15_FDIS", "ISO_8859-15", "L9", "LATIN0", "LATIN9", "Latin-9", "cp923", "csISO885915", "csISOlatin0", "csISOlatin9"),
	newCharset("ISO-8859-16", newSBCS(charmap.ISO8859_16, true), "ISO8859_16", "ISO_8859-16", "ISO_8859-16:2001", "csISO885916", "iso-ir-226", "l10", "latin10"),
	newCharset("TIS-620", sbcsTIS620, "tis620", "tis620.2533"),

	newCharset("KOI8-R", newSBCS(charmap.KOI8R, false), "cskoi8r", "koi8", "koi8_r"),
	newCharset("KOI8-U", sbcsKOI8U, "koi8_u"),

	newCharset("IBM437", newSBCS(charmap.CodePage437, false), "437", "cp437", "cspc8codepage437", "ibm-437", "ibm437", "windows-437"),
	newCharset("IBM850", newSBCS(charmap.CodePage850, false), "850", "cp850", "cspc850multilingual", "ibm-850", "ibm850"),
	newCharset("IBM852", newSBCS(charmap.CodePage852, false), "852", "cp852", "csPCp852", "ibm-852", "ibm852"),
	newCharset("IBM855", newSBCS(charmap.CodePage855, false), "855", "cp855", "cspcp855", "ibm-855", "ibm855"),
	newCharset("IBM00858", newSBCS(charmap.CodePage858, false), "858", "PC-Multilingual-850+euro", "ccsid00858", "cp00858", "cp858", "ibm-858", "ibm858"),
	newCharset("IBM860", newSBCS(charmap.CodePage860, false), "860", "cp860", "csIBM860", "ibm-860", "ibm860"),
	newCharset("IBM862", newSBCS(charmap.CodePage862, false), "862", "cp862", "csIBM862", "cspc862latinhebrew", "ibm-862", "ibm862"),
	newCharset("IBM863", newSBCS(charmap.CodePage863, false), "863", "cp863", "csIBM863", "ibm-863", "ibm863"),
	newCharset("IBM865", newSBCS(charmap.CodePage865, false), "865", "cp865", "csIBM865", "ibm-865", "ibm865"),
	newCharset("IBM866", newSBCS(charmap.CodePage866, false), "866", "cp866", "csIBM866", "ibm-866", "ibm866"),

	newCharset("x-MacRoman", newSBCS(charmap.Macintosh, false), "MacRoman"),
	newCharset("x-MacCyrillic", sbcsMacCyrillic, "MacCyrillic"),

	newCharset("Shift_JIS", newShiftJIS, "csShiftJIS", "ms_kanji", "shift-jis", "shift_jis", "sjis", "x-sjis"),
	newCharset("windows-31j", newWindows31J, "MS932", "csWindows31J", "windows-932"),
	newCharset("EUC-JP", newEUCJP, "Extended_UNIX_Code_Packed_Format_for_Japanese", "csEUCPkdFmtjapanese", "euc_jp", "eucjis", "eucjp", "x-euc-jp", "x-eucjp"),
	newCharset("ISO-2022-JP", newISO2022JP, "csISO2022JP", "csjisencoding", "iso2022jp", "jis", "jis_encoding"),
	newCharset("GBK", newGBK, "CP936", "windows-936"),
	newCharset("GB2312", newGB2312, "EUC_CN", "euc-cn", "euccn", "gb2312", "gb2312-1980", "gb2312-80", "x-EUC-CN"),
	newCharset("GB18030", newGB18030, "gb18030-2022"),
	newCharset("Big5", newBig5, "csBig5"),
	newCharset("EUC-KR", newEUCKR, "5601", "csEUCKR", "euc_kr", "euckr", "ks_c_5601-1987", "ksc5601", "ksc5601-1987", "ksc5601_1987", "ksc_5601"),
	newCharset("x-windows-949", newWindows949, "ms949", "ms_949", "windows-949", "windows949"),
}
