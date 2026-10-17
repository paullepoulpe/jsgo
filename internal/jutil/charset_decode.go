package jutil

import (
	"unicode/utf8"
)

// Go-only: language-adapter — decoding with a Charset (CharsetDecoder in
// REPLACE mode, as new String(bytes, cs) and InputStreamReader do, and in
// REPORT mode, as Files.readString does), and the UTF-8 and UTF-16 coders.

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
	s, e := c.impl().decode(b, true)
	switch {
	case e.length == 0:
		return s
	case e.unmappable:
		panic(NewUnmappableCharacterException(e.length))
	default:
		panic(NewMalformedInputException(e.length))
	}
}

// utf8Coder is UTF-8.
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

func (utf8Coder) canEncode(uint16) bool { return true }

func (utf8Coder) newEncoder() charsetEncoder { return utf8Encoder{} }

// utf8Encoder replaces an unpaired surrogate with '?' (see
// OutputStreamWriter).
type utf8Encoder struct{}

func (utf8Encoder) encodeRune(dst []byte, r rune) []byte { return utf8.AppendRune(dst, r) }
func (utf8Encoder) encodeMalformed(dst []byte) []byte    { return append(dst, '?') }
func (utf8Encoder) finish(dst []byte) []byte             { return dst }

// charsetUTF8MalformedLength is the length Java's UTF-8 decoder reports for
// the first malformed sequence of b.
func charsetUTF8MalformedLength(b []byte) int {
	for i := 0; i < len(b); {
		r, size := utf8.DecodeRune(b[i:])
		if r != utf8.RuneError || size > 1 {
			i += size
			continue
		}
		// The malformed sequence starts at i: Java reports the length of the
		// well-formed prefix of a multi-byte sequence (at least 1).
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
			// Observed behaviour (black-box): ED A0..BF, the prefix of an
			// encoded surrogate, is malformed with length 3 when complete
			// and 2 when cut short, as in REPLACE mode.
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

// utf16Coder is UTF-16 (byte order from an optional byte-order mark,
// big-endian by default), UTF-16BE or UTF-16LE (fixed order); see the
// Javadoc of java.nio.charset.StandardCharsets and Charset.
type utf16Coder struct {
	bigEndian bool
	// marked: the decoder honours (and drops) a leading byte-order mark
	// and the encoder writes one.
	marked bool
}

func (c utf16Coder) canEncode(uint16) bool { return true }

func (c utf16Coder) newEncoder() charsetEncoder { return &utf16Encoder{coder: c} }

// decode decodes UTF-16.
//
// Observed behaviour (black-box, CharsetDecoder in REPLACE and REPORT mode):
//   - UTF-16 looks for a mark only in the first two bytes: FE FF selects
//     big-endian and FF FE little-endian, and the mark is dropped. Any later
//     U+FEFF or U+FFFE is an ordinary character. The mark does not count in
//     the reported length of a later malformed sequence;
//   - a trailing odd byte is malformed, length 1;
//   - a lone low surrogate is malformed, length 2;
//   - a high surrogate followed by a unit that is not a low surrogate is
//     malformed with length 4: both units are consumed and replaced by a
//     single U+FFFD (so the second unit is lost even when it is valid);
//   - a high surrogate followed by fewer than two bytes is malformed with
//     length 2 plus the leftover bytes (2 or 3).
//
// In REPLACE mode each malformed sequence becomes one U+FFFD.
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

// utf16Encoder replaces an unpaired surrogate with U+FFFD. The marked form
// writes the big-endian mark U+FEFF before its first char.
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
