package jutil

import (
	"strings"
	"unicode/utf8"
)

// Go-only: language-adapter — java.nio.charset.Charset (the standard
// charsets only), String.getBytes(UTF_8) and new String(byte[], UTF_8).

// Charset is java.nio.charset.Charset, limited to the charsets every Java
// platform must support (StandardCharsets).
type Charset struct {
	name    string
	aliases []string
	// maxChar is the largest char the charset encodes (0xFFFF for the
	// Unicode encodings, which still reject lone surrogates).
	maxChar uint16
}

var (
	StandardCharsets_US_ASCII   = &Charset{name: "US-ASCII", aliases: []string{"646", "ANSI_X3.4-1968", "ANSI_X3.4-1986", "ASCII", "IBM367", "ISO646-US", "ISO_646.irv:1991", "ascii7", "cp367", "csASCII", "iso-ir-6", "iso_646.irv:1983", "us"}, maxChar: 0x7F}
	StandardCharsets_ISO_8859_1 = &Charset{name: "ISO-8859-1", aliases: []string{"819", "8859_1", "IBM-819", "IBM819", "ISO8859-1", "ISO8859_1", "ISO_8859-1", "ISO_8859-1:1987", "ISO_8859_1", "cp819", "csISOLatin1", "iso-ir-100", "l1", "latin1"}, maxChar: 0xFF}
	StandardCharsets_UTF_8      = &Charset{name: "UTF-8", aliases: []string{"UTF8", "unicode-1-1-utf-8"}, maxChar: 0xFFFF}
	StandardCharsets_UTF_16     = &Charset{name: "UTF-16", aliases: []string{"UTF_16", "UnicodeBig", "unicode", "utf16"}, maxChar: 0xFFFF}
	StandardCharsets_UTF_16BE   = &Charset{name: "UTF-16BE", aliases: []string{"ISO-10646-UCS-2", "UTF_16BE", "UnicodeBigUnmarked", "X-UTF-16BE"}, maxChar: 0xFFFF}
	StandardCharsets_UTF_16LE   = &Charset{name: "UTF-16LE", aliases: []string{"UTF_16LE", "UnicodeLittleUnmarked", "X-UTF-16LE"}, maxChar: 0xFFFF}
)

var standardCharsets = []*Charset{
	StandardCharsets_US_ASCII, StandardCharsets_ISO_8859_1, StandardCharsets_UTF_8,
	StandardCharsets_UTF_16, StandardCharsets_UTF_16BE, StandardCharsets_UTF_16LE,
}

// UnsupportedCharsetException is java.nio.charset.UnsupportedCharsetException.
type UnsupportedCharsetException struct{ IllegalArgumentException }

// NewUnsupportedCharsetException is new UnsupportedCharsetException(name).
func NewUnsupportedCharsetException(charsetName string) *UnsupportedCharsetException {
	e := &UnsupportedCharsetException{}
	InitThrowable(e, "java.nio.charset.UnsupportedCharsetException", msgPtr(charsetName), nil)
	return e
}

// Charset_ForName is Charset.forName(String): canonical names and aliases,
// case-insensitively. Charsets other than the standard ones are reported
// as unsupported.
func Charset_ForName(charsetName string) *Charset {
	lower := StringToLowerCaseASCII(charsetName)
	for _, c := range standardCharsets {
		if lower == StringToLowerCaseASCII(c.name) {
			return c
		}
		for _, a := range c.aliases {
			if lower == StringToLowerCaseASCII(a) {
				return c
			}
		}
	}
	panic(NewUnsupportedCharsetException(charsetName))
}

// Name is Charset.name().
func (c *Charset) Name() string { return c.name }

// String is Charset.toString().
func (c *Charset) String() string { return c.name }

// CanEncodeChar is charset.newEncoder().canEncode(char): a lone surrogate
// can never be encoded on its own.
func (c *Charset) CanEncodeChar(ch uint16) bool {
	if isSurrogate(ch) {
		return false
	}
	return ch <= c.maxChar
}

// StringGetBytesUTF8 is String.getBytes(StandardCharsets.UTF_8): unpaired
// surrogates become '?'.
func StringGetBytesUTF8(s string) []byte {
	if utf8.ValidString(s) {
		return []byte(s)
	}
	units := StringToUTF16(s)
	out := make([]byte, 0, len(s))
	for i := 0; i < len(units); i++ {
		c := units[i]
		switch {
		case c < 0x80:
			out = append(out, byte(c))
		case c < 0x800:
			out = append(out, byte(0xC0|c>>6), byte(0x80|c&0x3F))
		case c >= 0xD800 && c <= 0xDBFF && i+1 < len(units) && units[i+1] >= 0xDC00 && units[i+1] <= 0xDFFF:
			r := 0x10000 + (rune(c)-0xD800)<<10 + (rune(units[i+1]) - 0xDC00)
			out = utf8.AppendRune(out, r)
			i++
		case isSurrogate(c):
			out = append(out, '?')
		default:
			out = append(out, byte(0xE0|c>>12), byte(0x80|(c>>6)&0x3F), byte(0x80|c&0x3F))
		}
	}
	return out
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
