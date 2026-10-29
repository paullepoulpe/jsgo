// Port of com/google/debugging/sourcemap/Base64.java @ closure-compiler v20261021
// Port of com/google/debugging/sourcemap/Base64VLQ.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sourcemap

import (
	"io"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Base64_BASE64_MAP is a map used to convert integer values in the range
// 0-63 to their base64 values.
const Base64_BASE64_MAP = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz" +
	"0123456789+/"

// Base64_BASE64_DECODE_MAP is a map used to convert base64 character into
// integer values.
var Base64_BASE64_DECODE_MAP = func() [256]int8 {
	var m [256]int8
	for i := range m {
		m[i] = -1
	}
	for i := 0; i < len(Base64_BASE64_MAP); i++ {
		m[Base64_BASE64_MAP[i]] = int8(i)
	}
	return m
}()

// Base64_ToBase64 returns a base64 digit for a value in the range of 0-63.
//
// Java: Base64.toBase64(int)
func Base64_ToBase64(value int32) uint16 {
	if uint32(value) >= uint32(len(Base64_BASE64_MAP)) {
		panic(jutil.NewStringIndexOutOfBoundsExceptionForIndex(int(value), len(Base64_BASE64_MAP)))
	}
	return uint16(Base64_BASE64_MAP[value])
}

// Base64_FromBase64 returns a value in the range of 0-63 for a base64
// digit.
//
// Java: Base64.fromBase64(char)
func Base64_FromBase64(c uint16) int32 {
	return int32(Base64_BASE64_DECODE_MAP[c])
}

// Base64_Base64EncodeInt returns the six digit long base64 encoded value of
// the integer.
//
// Java: Base64.base64EncodeInt(int)
func Base64_Base64EncodeInt(value int32) string {
	var c [6]byte
	for i := 0; i < 5; i++ {
		c[i] = byte(Base64_ToBase64(jutil.Shr32(value, int64(26-i*6)) & 0x3f))
	}
	c[5] = byte(Base64_ToBase64(jutil.Shl32(value, 4) & 0x3f))
	return string(c[:])
}

// Base64VLQ: we encode our variable length numbers as base64 encoded
// strings with the least significant digit coming first. Each base64 digit
// encodes a 5-bit value (0-31) and a continuation bit. Signed values can be
// represented by using the least significant bit of the value as the sign
// bit.

// A Base64 VLQ digit can represent 5 bits, so it is base-32.
const (
	Base64VLQ_VLQ_BASE_SHIFT = 5
	Base64VLQ_VLQ_BASE       = 1 << Base64VLQ_VLQ_BASE_SHIFT

	// A mask of bits for a VLQ digit (11111), 31 decimal.
	Base64VLQ_VLQ_BASE_MASK = Base64VLQ_VLQ_BASE - 1

	// The continuation bit is the 6th bit.
	Base64VLQ_VLQ_CONTINUATION_BIT = Base64VLQ_VLQ_BASE
)

// Base64VLQ_ToVLQSigned converts from a two-complement value to a value
// where the sign bit is is placed in the least significant bit. For
// example, as decimals:
//
//	1 becomes 2 (10 binary), -1 becomes 3 (11 binary)
//	2 becomes 4 (100 binary), -2 becomes 5 (101 binary)
//
// Java: Base64VLQ.toVLQSigned(int)
func Base64VLQ_ToVLQSigned(value int32) int32 {
	if value < 0 {
		return ((-value) << 1) + 1
	} else {
		return (value << 1) + 0
	}
}

// Base64VLQ_FromVLQSigned converts to a two-complement value from a value
// where the sign bit is is placed in the least significant bit. For
// example, as decimals:
//
//	2 (10 binary) becomes 1, 3 (11 binary) becomes -1
//	4 (100 binary) becomes 2, 5 (101 binary) becomes -2
//
// Java: Base64VLQ.fromVLQSigned(int)
func Base64VLQ_FromVLQSigned(value int32) int32 {
	negate := (value & 1) == 1
	value = jutil.Ushr32(value, 1)
	if !negate {
		return value
	}
	// We need to OR 0x80000000 here to ensure the 32nd bit (the sign bit)
	// is always set for negative numbers. If `value` were 1, (meaning
	// `negate` is true and all other bits were zeros), `value` would now be
	// 0. -0 is just 0, and doesn't flip the 32nd bit as intended. All
	// positive numbers will successfully flip the 32nd bit without issue, so
	// it's a noop for them.
	return -value | -0x80000000
}

// Base64VLQ_Encode writes a VLQ encoded value to the provide appendable.
//
// Java: Base64VLQ.encode(Appendable, int)
func Base64VLQ_Encode(out io.StringWriter, value int32) {
	var buf [8]byte
	n := base64VLQ_encodeToBytes(&buf, value)
	appendBytes(out, buf[:n])
}

func base64VLQ_encodeToBytes(buf *[8]byte, value int32) int {
	value = Base64VLQ_ToVLQSigned(value)
	n := 0
	for {
		digit := value & Base64VLQ_VLQ_BASE_MASK
		value = jutil.Ushr32(value, Base64VLQ_VLQ_BASE_SHIFT)
		if value > 0 {
			digit |= Base64VLQ_VLQ_CONTINUATION_BIT
		}
		buf[n] = byte(Base64_ToBase64(digit))
		n++
		if !(value > 0) {
			break
		}
	}
	return n
}

// Base64VLQ_CharIterator is a simple interface for advancing through a
// sequence of characters, that communicates that advance back to the
// source.
type Base64VLQ_CharIterator interface {
	// Java: Base64VLQ.CharIterator.hasNext()
	HasNext() bool

	// Java: Base64VLQ.CharIterator.next()
	Next() uint16
}

// Base64VLQ_Decode decodes the next VLQValue from the provided
// CharIterator.
//
// Java: Base64VLQ.decode(CharIterator)
func Base64VLQ_Decode(in Base64VLQ_CharIterator) int32 {
	var result int32
	var continuation bool
	var shift int32
	for {
		c := in.Next()
		digit := Base64_FromBase64(c)
		continuation = (digit & Base64VLQ_VLQ_CONTINUATION_BIT) != 0
		digit &= Base64VLQ_VLQ_BASE_MASK
		result = result + jutil.Shl32(digit, int64(shift))
		shift = shift + Base64VLQ_VLQ_BASE_SHIFT
		if !continuation {
			break
		}
	}

	return Base64VLQ_FromVLQSigned(result)
}

// Go-only: language-adapter — Java's Appendable.append(CharSequence) and
// append(char) on an io.StringWriter; a write error is the IOException
// Appendable.append throws.
func appendString(out io.StringWriter, s string) {
	if _, err := out.WriteString(s); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
}

func appendBytes(out io.StringWriter, b []byte) {
	if w, ok := out.(io.Writer); ok {
		if _, err := w.Write(b); err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
		return
	}
	appendString(out, string(b))
}

func appendByte(out io.StringWriter, b byte) {
	if bw, ok := out.(io.ByteWriter); ok {
		if err := bw.WriteByte(b); err != nil {
			panic(jutil.NewIOException(err.Error()))
		}
		return
	}
	var tmp [1]byte
	tmp[0] = b
	appendBytes(out, tmp[:])
}

func appendChar(out io.StringWriter, c uint16) {
	appendString(out, jutil.StringFromChar(c))
}
