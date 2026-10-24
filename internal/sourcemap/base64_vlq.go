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
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Base64VLQ: we encode our variable length numbers as base64 encoded
// strings with the least significant digit coming first. Each base64 digit
// encodes a 5-bit value (0-31) and a continuation bit. Signed values can be
// represented by using the least significant bit of the value as the sign
// bit.
type Base64VLQ struct{}

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
	appendString(out, unsafe.String(&buf[0], n))
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

func appendChar(out io.StringWriter, c uint16) {
	appendString(out, jutil.StringFromChar(c))
}
