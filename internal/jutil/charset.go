package jutil

import (
	"strings"
	"sync"
	"unicode/utf8"
)

// Go-only: language-adapter — java.nio.charset.Charset for the charsets of
// charset_registry.go, and String.getBytes(UTF_8) and new String(byte[],
// UTF_8). Each charset's conversions (charsetCoder) are built on first use.

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
