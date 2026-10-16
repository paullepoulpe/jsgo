package jutil

import (
	"encoding/base64"
	"strconv"
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

func baseEncodingBase64Decode(ch uint16) int {
	if ch > 0x7F {
		panic(newBaseEncodingDecodingException("Unrecognized character: 0x" + strconv.FormatInt(int64(ch), 16)))
	}
	for i := 0; i < len(baseEncodingBase64Alphabet); i++ {
		if uint16(baseEncodingBase64Alphabet[i]) == ch {
			return i
		}
	}
	if ch <= 0x20 || ch == 0x7F {
		panic(newBaseEncodingDecodingException("Unrecognized character: 0x" + strconv.FormatInt(int64(ch), 16)))
	}
	panic(newBaseEncodingDecodingException("Unrecognized character: " + StringFromChar(ch)))
}

// Decode is BaseEncoding.decode(CharSequence): it panics with
// IllegalArgumentException (caused by a DecodingException) on bad input.
func (b *BaseEncoding) Decode(chars string) []byte {
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
