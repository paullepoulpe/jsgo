package jutil

import (
	"strings"
	"unicode/utf8"
)

// Go-only: language-adapter — decoding with the standard Charsets
// (CharsetDecoder in REPLACE mode, as new String(bytes, cs) and
// InputStreamReader do, and in REPORT mode, as Files.readString does).

// Decode is new String(bytes, charset): each malformed input sequence is
// replaced with U+FFFD.
func (c *Charset) Decode(b []byte) string {
	s, _ := c.decode(b, false)
	return s
}

// DecodeStrict decodes b and panics with MalformedInputException at the
// first malformed input sequence (CodingErrorAction.REPORT).
func (c *Charset) DecodeStrict(b []byte) string {
	s, malformed := c.decode(b, true)
	if malformed > 0 {
		panic(NewMalformedInputException(malformed))
	}
	return s
}

// decode returns the decoded string; when report is set it stops at the
// first malformed sequence and returns its length (> 0).
func (c *Charset) decode(b []byte, report bool) (string, int) {
	switch c {
	case StandardCharsets_UTF_8:
		if utf8.Valid(b) {
			return string(b), 0
		}
		if report {
			return "", charsetUTF8MalformedLength(b)
		}
		return NewStringFromUTF8Bytes(b), 0
	case StandardCharsets_ISO_8859_1:
		var sb strings.Builder
		for _, x := range b {
			sb.WriteRune(rune(x))
		}
		return sb.String(), 0
	case StandardCharsets_US_ASCII:
		var sb strings.Builder
		for _, x := range b {
			if x >= 0x80 {
				if report {
					return "", 1
				}
				sb.WriteRune(utf8.RuneError)
			} else {
				sb.WriteByte(x)
			}
		}
		return sb.String(), 0
	case StandardCharsets_UTF_16, StandardCharsets_UTF_16BE, StandardCharsets_UTF_16LE:
		return c.decodeUTF16(b, report)
	}
	panic(NewUnsupportedCharsetException(c.name))
}

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

// decodeUTF16 decodes UTF-16, UTF-16BE or UTF-16LE (see the Javadoc of
// java.nio.charset.StandardCharsets and Charset: UTF-16 takes its byte order
// from an optional byte-order mark and defaults to big-endian; the BE and LE
// forms have a fixed order and keep U+FEFF as a character).
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
func (c *Charset) decodeUTF16(b []byte, report bool) (string, int) {
	bigEndian := c != StandardCharsets_UTF_16LE
	if c == StandardCharsets_UTF_16 && len(b) >= 2 {
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
		malformed := 0
		switch rest := len(b) - i; {
		case rest < 2:
			malformed = rest
		default:
			u := unit(i)
			switch {
			case isLowSurrogateUnit(u):
				malformed = 2
			case !isHighSurrogateUnit(u):
				units = append(units, u)
				i += 2
				continue
			case rest < 4:
				malformed = rest
			default:
				if v := unit(i + 2); isLowSurrogateUnit(v) {
					units = append(units, u, v)
					i += 4
					continue
				}
				malformed = 4
			}
		}
		if report {
			return "", malformed
		}
		units = append(units, 0xFFFD)
		i += malformed
	}
	return StringFromUTF16(units), 0
}
