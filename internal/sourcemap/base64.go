// Port of com/google/debugging/sourcemap/Base64.java @ closure-compiler v20261021

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Base64 is a utility class for working with Base64 values.
type Base64 struct{}

// Base64_BASE64_MAP is a map used to convert integer values in the range
// 0-63 to their base64 values.
const Base64_BASE64_MAP = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz" +
	"0123456789+/"

// Base64_BASE64_DECODE_MAP is a map used to convert base64 character into
// integer values.
var Base64_BASE64_DECODE_MAP = func() []int32 {
	m := make([]int32, 256)
	for i := range m {
		m[i] = -1
	}
	for i := 0; i < len(Base64_BASE64_MAP); i++ {
		m[Base64_BASE64_MAP[i]] = int32(i)
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
	result := Base64_BASE64_DECODE_MAP[c]
	return result
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
