package jutil

import (
	"encoding/base64"
	"strconv"
	"unicode/utf8"
)

// Go-only: language-adapter — Guava BaseEncoding.base64() (RFC 4648
// standard alphabet, '=' padding): encode, and Guava's lenient-padding,
// strict-alphabet decode (Base64Encoding.decodeTo).

// BaseEncoding is com.google.common.io.BaseEncoding.
type BaseEncoding struct{}

var baseEncodingBase64 = &BaseEncoding{}

// BaseEncoding_Base64 is BaseEncoding.base64().
func BaseEncoding_Base64() *BaseEncoding { return baseEncodingBase64 }

// Encode is BaseEncoding.encode(byte[]).
func (b *BaseEncoding) Encode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// BaseEncoding_DecodingException is BaseEncoding.DecodingException.
type BaseEncoding_DecodingException struct{ IOException }

func newBaseEncodingDecodingException(message string) *BaseEncoding_DecodingException {
	e := &BaseEncoding_DecodingException{}
	InitThrowable(e, "com.google.javascript.jscomp.jarjar.com.google.common.io.BaseEncoding$DecodingException",
		msgPtr(message), nil)
	return e
}

const baseEncodingBase64Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var baseEncodingBase64DecTable = func() [256]int8 {
	var t [256]int8
	for i := range t {
		t[i] = -1
	}
	for i := 0; i < len(baseEncodingBase64Alphabet); i++ {
		t[baseEncodingBase64Alphabet[i]] = int8(i)
	}
	return t
}()

func baseEncodingBase64Decode(ch uint16) int {
	if ch > 0x7F {
		panic(newBaseEncodingDecodingException("Unrecognized character: 0x" + strconv.FormatInt(int64(ch), 16)))
	}
	if v := baseEncodingBase64DecTable[ch]; v >= 0 {
		return int(v)
	}
	if ch <= 0x20 || ch == 0x7F {
		panic(newBaseEncodingDecodingException("Unrecognized character: 0x" + strconv.FormatInt(int64(ch), 16)))
	}
	panic(newBaseEncodingDecodingException("Unrecognized character: " + StringFromChar(ch)))
}

// Decode is BaseEncoding.decode(CharSequence): it panics with
// IllegalArgumentException (caused by a DecodingException) on bad input.
func (b *BaseEncoding) Decode(chars string) []byte {
	l := len(chars)
	for l > 0 && chars[l-1] == '=' {
		l--
	}
	if l%4 != 1 {
		target := make([]byte, 0, l*3/4)
		i := 0
		ok := true
		for i+4 <= l {
			c0 := baseEncodingBase64DecTable[chars[i]]
			c1 := baseEncodingBase64DecTable[chars[i+1]]
			c2 := baseEncodingBase64DecTable[chars[i+2]]
			c3 := baseEncodingBase64DecTable[chars[i+3]]
			if (c0 | c1 | c2 | c3) < 0 {
				ok = false
				break
			}
			chunk := uint32(c0)<<18 | uint32(c1)<<12 | uint32(c2)<<6 | uint32(c3)
			target = append(target, byte(chunk>>16), byte(chunk>>8), byte(chunk))
			i += 4
		}
		if ok {
			rem := l - i
			if rem == 2 {
				c0 := baseEncodingBase64DecTable[chars[i]]
				c1 := baseEncodingBase64DecTable[chars[i+1]]
				if (c0 | c1) >= 0 {
					chunk := uint32(c0)<<18 | uint32(c1)<<12
					return append(target, byte(chunk>>16))
				}
			} else if rem == 3 {
				c0 := baseEncodingBase64DecTable[chars[i]]
				c1 := baseEncodingBase64DecTable[chars[i+1]]
				c2 := baseEncodingBase64DecTable[chars[i+2]]
				if (c0 | c1 | c2) >= 0 {
					chunk := uint32(c0)<<18 | uint32(c1)<<12 | uint32(c2)<<6
					return append(target, byte(chunk>>16), byte(chunk>>8))
				}
			} else {
				return target
			}
		}
	}
	return b.decodeSlow(chars)
}

func (b *BaseEncoding) decodeSlow(chars string) []byte {
	var out []byte
	Try(func() {
		out = b.decodeChecked(chars)
	}, func(badInput *BaseEncoding_DecodingException) {
		e := &IllegalArgumentException{}
		InitThrowable(e, "java.lang.IllegalArgumentException", nil, badInput)
		panic(e)
	})
	return out
}

// ValidateAndCheckCanonicalUTF8 validates chars according to Guava's
// BaseEncoding.base64().decode(CharSequence) without allocating on the heap
// (panicking with the identical IllegalArgumentException on malformed input),
// and reports whether chars is canonical RFC 4648 base64 whose decoded bytes
// are valid UTF-8 without a leading UTF-8 BOM (so re-encoding GetCode() yields
// chars verbatim).
func (b *BaseEncoding) ValidateAndCheckCanonicalUTF8(chars string) bool {
	l := len(chars)
	for l > 0 && chars[l-1] == '=' {
		l--
	}
	if l%4 == 1 {
		b.decodeSlow(chars)
		return false
	}
	pad := len(chars) - l
	// "77u/" decodes to the 3-byte UTF-8 BOM (0xEF, 0xBB, 0xBF), which
	// SourceFile.SetCodeAndDoBookkeeping strips on load.
	canonical := len(chars)%4 == 0 && pad <= 2 && !(len(chars) >= 4 && chars[:4] == "77u/")

	var buf [772]byte
	carry := 0
	flushUTF8 := func(n int, final bool) {
		if !canonical {
			return
		}
		if final || n == 0 || buf[n-1] < 0x80 {
			if !utf8.Valid(buf[:n]) {
				canonical = false
			}
			carry = 0
			return
		}
		split := n - 1
		for split > 0 && n-split < 4 && (buf[split]&0xC0) == 0x80 {
			split--
		}
		if !utf8.FullRune(buf[split:n]) {
			if !utf8.Valid(buf[:split]) {
				canonical = false
				carry = 0
				return
			}
			carry = copy(buf[:], buf[split:n])
			return
		}
		if !utf8.Valid(buf[:n]) {
			canonical = false
		}
		carry = 0
	}

	pos := carry
	i := 0
	for i+4 <= l {
		c0 := baseEncodingBase64DecTable[chars[i]]
		c1 := baseEncodingBase64DecTable[chars[i+1]]
		c2 := baseEncodingBase64DecTable[chars[i+2]]
		c3 := baseEncodingBase64DecTable[chars[i+3]]
		if (c0 | c1 | c2 | c3) < 0 {
			b.decodeSlow(chars)
			return false
		}
		if canonical {
			chunk := uint32(c0)<<18 | uint32(c1)<<12 | uint32(c2)<<6 | uint32(c3)
			buf[pos] = byte(chunk >> 16)
			buf[pos+1] = byte(chunk >> 8)
			buf[pos+2] = byte(chunk)
			pos += 3
			if pos+3 > len(buf) {
				flushUTF8(pos, false)
				pos = carry
			}
		}
		i += 4
	}
	rem := l - i
	if rem == 2 {
		c0 := baseEncodingBase64DecTable[chars[i]]
		c1 := baseEncodingBase64DecTable[chars[i+1]]
		if (c0 | c1) < 0 {
			b.decodeSlow(chars)
			return false
		}
		if canonical {
			if (c1 & 0x0F) != 0 {
				canonical = false
			} else {
				chunk := uint32(c0)<<18 | uint32(c1)<<12
				buf[pos] = byte(chunk >> 16)
				pos++
			}
		}
	} else if rem == 3 {
		c0 := baseEncodingBase64DecTable[chars[i]]
		c1 := baseEncodingBase64DecTable[chars[i+1]]
		c2 := baseEncodingBase64DecTable[chars[i+2]]
		if (c0 | c1 | c2) < 0 {
			b.decodeSlow(chars)
			return false
		}
		if canonical {
			if (c2 & 0x03) != 0 {
				canonical = false
			} else {
				chunk := uint32(c0)<<18 | uint32(c1)<<12 | uint32(c2)<<6
				buf[pos] = byte(chunk >> 16)
				buf[pos+1] = byte(chunk >> 8)
				pos += 2
			}
		}
	}
	if canonical && pos > 0 {
		flushUTF8(pos, true)
	}
	return canonical
}

func (b *BaseEncoding) decodeChecked(s string) []byte {
	chars := StringToUTF16(s)
	// trimTrailingPadding
	l := len(chars)
	for l > 0 && chars[l-1] == '=' {
		l--
	}
	chars = chars[:l]
	switch len(chars) % 4 {
	case 0, 2, 3:
	default:
		panic(newBaseEncodingDecodingException("Invalid input length " + strconv.Itoa(len(chars))))
	}
	target := make([]byte, 0, len(chars)*3/4+3)
	for i := 0; i < len(chars); {
		chunk := baseEncodingBase64Decode(chars[i]) << 18
		i++
		chunk |= baseEncodingBase64Decode(chars[i]) << 12
		i++
		target = append(target, byte(chunk>>16))
		if i < len(chars) {
			chunk |= baseEncodingBase64Decode(chars[i]) << 6
			i++
			target = append(target, byte((chunk>>8)&0xFF))
			if i < len(chars) {
				chunk |= baseEncodingBase64Decode(chars[i])
				i++
				target = append(target, byte(chunk&0xFF))
			}
		}
	}
	return target
}
