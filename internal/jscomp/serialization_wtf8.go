// Port of com/google/javascript/jscomp/serialization/Wtf8.java @ closure-compiler v20260415

package jscomp

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Wtf8 is an encoder/decoder into https://simonsapin.github.io/wtf-8/
type Wtf8 struct{}

const wtf8_CONTINUATION_MASK = 0x3F

// Wtf8_EncodeToWtf8 encodes s (a Java String: UTF-16 code units, possibly
// with lone surrogates) as WTF-8. ByteString is []byte.
//
// Java: Wtf8.encodeToWtf8(String)
func Wtf8_EncodeToWtf8(s string) []byte {
	units := jutil.StringToUTF16(s)
	length := len(units)
	output := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		codepoint := wtf8_codePointAt(units, i)

		if codepoint < 0x80 {
			output = append(output, byte(codepoint))
		} else if codepoint < 0x800 {
			output = append(output, byte(0xC0|(0x1f&(codepoint>>6))))
			output = append(output, byte(0x80|(wtf8_CONTINUATION_MASK&codepoint)))
		} else if codepoint < 0x10000 {
			output = append(output, byte(0xE0|(0xf&(codepoint>>12))))
			output = append(output, byte(0x80|(wtf8_CONTINUATION_MASK&(codepoint>>6))))
			output = append(output, byte(0x80|(wtf8_CONTINUATION_MASK&codepoint)))
		} else {
			// This codepoints takes two UTF-16 code units, so we need an extra
			// increment.
			i++
			output = append(output, byte(0xF0|(0x7&(codepoint>>18))))
			output = append(output, byte(0x80|(wtf8_CONTINUATION_MASK&(codepoint>>12))))
			output = append(output, byte(0x80|(wtf8_CONTINUATION_MASK&(codepoint>>6))))
			output = append(output, byte(0x80|(wtf8_CONTINUATION_MASK&codepoint)))
		}
	}
	return output
}

// wtf8_codePointAt is String.codePointAt(int) on UTF-16 units.
func wtf8_codePointAt(units []uint16, i int) int {
	c := units[i]
	if c >= 0xD800 && c <= 0xDBFF && i+1 < len(units) {
		if lo := units[i+1]; lo >= 0xDC00 && lo <= 0xDFFF {
			return 0x10000 + (int(c)-0xD800)<<10 + (int(lo) - 0xDC00)
		}
	}
	return int(c)
}

// Wtf8_Decoder_ is Wtf8.decoder(int).
//
// Go-only: language-adapter — the trailing underscore keeps the mechanical
// name Wtf8_Decoder for the nested Decoder type.
func Wtf8_Decoder_(maxLength int) *Wtf8_Decoder {
	return newWtf8_Decoder(maxLength)
}

// Wtf8_Decoder decodes strings from WTF8 bytes.
//
// This class is not threadsafe! In order to minimize allocations, and
// maximize performance, it reuses buffers between decodings.
type Wtf8_Decoder struct {
	codepointBuffer []int32
}

func newWtf8_Decoder(maxLength int) *Wtf8_Decoder {
	if maxLength < 0 {
		e := &jutil.RuntimeException{}
		msg := strconv.Itoa(maxLength)
		jutil.InitThrowable(e, "java.lang.NegativeArraySizeException", &msg, nil)
		panic(e)
	}
	return &Wtf8_Decoder{codepointBuffer: make([]int32, maxLength)}
}

// wtf8_byteAt is ByteString.byteAt(int) (signed).
func wtf8_byteAt(encoded []byte, i int) int32 {
	if i < 0 || i >= len(encoded) {
		panic(jutil.NewArrayIndexOutOfBoundsExceptionForLength(i, len(encoded)))
	}
	return int32(int8(encoded[i]))
}

// Java: Wtf8.Decoder.decode(ByteString)
func (d *Wtf8_Decoder) Decode(encoded []byte) string {
	encodedByteCount := len(encoded)
	codepointCount := 0

	for i := 0; i < encodedByteCount; i++ {
		b := wtf8_byteAt(encoded, i)
		var codepoint int32

		if (b & 0x80) == 0 {
			// 0xxx xxxx: 1 byte
			codepoint = b
		} else if (b & 0xE0) == 0xC0 {
			// 110x xxxx: 2 bytes
			firstByte := 0x1F & b
			i++
			secondByte := wtf8_byteAt(encoded, i) & wtf8_CONTINUATION_MASK
			codepoint = (firstByte << 6) | secondByte
		} else if (b & 0xF0) == 0xE0 {
			// 1110 xxxx: 3 bytes
			firstByte := 0xF & b
			i++
			secondByte := wtf8_byteAt(encoded, i) & wtf8_CONTINUATION_MASK
			i++
			thirdByte := wtf8_byteAt(encoded, i) & wtf8_CONTINUATION_MASK
			codepoint = (firstByte << 12) | (secondByte << 6) | thirdByte
		} else if (b & 0xF8) == 0xF0 {
			// 1111 0xxx: 4 bytes
			firstByte := 0x7 & b
			i++
			secondByte := wtf8_byteAt(encoded, i) & wtf8_CONTINUATION_MASK
			i++
			thirdByte := wtf8_byteAt(encoded, i) & wtf8_CONTINUATION_MASK
			i++
			fourthByte := wtf8_byteAt(encoded, i) & wtf8_CONTINUATION_MASK
			codepoint = (firstByte << 18) | (secondByte << 12) | (thirdByte << 6) | fourthByte
		} else {
			panic(jutil.NewAssertionErrorNoMessage())
		}

		if codepointCount >= len(d.codepointBuffer) {
			panic(jutil.NewArrayIndexOutOfBoundsExceptionForLength(codepointCount, len(d.codepointBuffer)))
		}
		d.codepointBuffer[codepointCount] = codepoint
		codepointCount++
	}

	return wtf8_newStringFromCodePoints(d.codepointBuffer[:codepointCount])
}

// wtf8_newStringFromCodePoints is new String(int[] codePoints, 0, n).
func wtf8_newStringFromCodePoints(codePoints []int32) string {
	units := make([]uint16, 0, len(codePoints))
	for _, cp := range codePoints {
		switch {
		case cp >= 0 && cp < 0x10000:
			units = append(units, uint16(cp))
		case cp >= 0x10000 && cp <= 0x10FFFF:
			cp -= 0x10000
			units = append(units, uint16(0xD800+(cp>>10)), uint16(0xDC00+(cp&0x3FF)))
		default:
			panic(jutil.NewIllegalArgumentExceptionMsg("Not a valid Unicode code point: 0x" +
				strconv.FormatUint(uint64(uint32(cp)), 16)))
		}
	}
	return jutil.StringFromUTF16(units)
}

// Java: Wtf8()
func newWtf8() *Wtf8 { return &Wtf8{} }
