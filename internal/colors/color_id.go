// Port of com/google/javascript/jscomp/colors/ColorId.java @ closure-compiler v20261021

/*
 * Copyright 2020 The Closure Compiler Authors.
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

// Package colors is the port of com.google.javascript.jscomp.colors.
package colors

import (
	"encoding/binary"
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ColorId is a probablistically unique ID for a Color.
//
// This is a semantic wrapper around UUID bytes to provide better typing when
// passing those bytes around. It also hosts some utility methods to ensure
// consistent byte handling and help with testing/debugging.
//
// If two Colors have the same ID, they should be considered to represent the
// same "logical" Color, perhaps inferred in different libraries. Ideally,
// those instances would be reconciled before either Color is instantiated.
//
// IDs have an intentionally unspecified width so that it may expand if
// needed in the future. For the purposes of comparison, serialization, and
// display, leading 0s are ignored.
//
// Go-only: language-adapter — ColorId is an immutable value class whose
// equals/hashCode compare the wrapped long, so it is a Go value type: Go ==
// is Java equals, and maps can be keyed by it directly. Nullable ColorId
// references are *ColorId. ColorId implements LogsGson.Able (ToLogsGson)
// structurally.
type ColorId struct {
	// Assume for now that we have at most 8 bytes.
	rightAligned int64
}

// Java: ColorId.fromUnsigned(long)
func ColorId_FromUnsigned(x int64) ColorId {
	return NewColorId(x)
}

// Java: ColorId.fromUnsigned(int)
func ColorId_FromUnsignedInt(x int32) ColorId {
	return NewColorId(int64(x) & 0xFFFFFFFF)
}

// ColorId_FromBytes builds an ID from a protobuf ByteString (Go []byte).
//
// Java: ColorId.fromBytes(ByteString)
func ColorId_FromBytes(bytes []byte) ColorId {
	if len(bytes) > 8 {
		panic(jutil.NewIllegalStateExceptionMsg(strconv.Itoa(len(bytes))))
	}
	var rightAligned int64
	for _, b := range bytes {
		rightAligned = (rightAligned << 8) | int64(b)
	}
	return NewColorId(rightAligned)
}

// Java: ColorId.fromAscii(String)
func ColorId_FromAscii(str string) ColorId {
	return ColorId_FromBytesAt(str, len(str), ColorId_ASCII_AT)
}

// Java: ColorId.fromBytesAt(T, int, ByteAt<T>)
func ColorId_FromBytesAt[T any](source T, length int, byteAt ColorId_ByteAt[T]) ColorId {
	if length > 8 {
		panic(jutil.NewIllegalStateExceptionMsg(strconv.Itoa(length)))
	}

	var rightAligned int64 = 0
	for i := 0; i < length; i++ {
		rightAligned <<= 8
		rightAligned |= int64(byteAt(source, i)) & 0xFF
	}

	return NewColorId(rightAligned)
}

// ColorId_Union defines how IDs for unions are combined.
//
// Ideally, this function would preserve the algebraic properties of the
// union operation. However, doing so is difficult (impossible?) while also
// providing good hashing. Those properties are summarized below.
//
// This method must never be passed IDs of existing unions. Lack of
// associativity means the result will not be the same as unioning all the
// IDs simultaneously.
//
//	| name          | has   | example                       | note                        |
//	|---------------|-------|-------------------------------|-----------------------------|
//	| associativity | false | f(x, f(y, z)) = f(f(x, y), z) |                             |
//	| commutativity | true  |       f(x, y) = f(y, x)       |                             |
//	| idempotence   | true  |    f(x, x, y) = f(x, y)       | `ids` is a Set              |
//	| identity      | true  |          f(x) = x             | singleton unions are banned |
//
// Java: ColorId.union(Set<ColorId>)
func ColorId_Union(ids *jutil.LinkedHashSet[ColorId]) ColorId {
	return ColorId_UnionSlice(ids.ToSlice())
}

// ColorId_UnionSlice combines a slice of ColorIds into a union ColorId.
func ColorId_UnionSlice(ids []ColorId) ColorId {
	if len(ids) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	if len(ids) == 1 {
		return ids[0]
	}
	var stackBuf [32]int64
	var sorted []int64
	if len(ids) <= len(stackBuf) {
		sorted = stackBuf[:len(ids)]
	} else {
		sorted = make([]int64, len(ids))
	}
	for i, c := range ids {
		sorted[i] = c.rightAligned
	}
	return unionRightAlignedIds(sorted)
}

func unionRightAlignedIds(ids []int64) ColorId {
	slices.Sort(ids)
	ids = slices.Compact(ids)
	n := len(ids)
	if n == 1 {
		return NewColorId(ids[0])
	}
	if n == 2 {
		var buf [16]byte
		binary.LittleEndian.PutUint64(buf[0:8], uint64(ids[0]))
		binary.LittleEndian.PutUint64(buf[8:16], uint64(ids[1]))
		return NewColorId(int64(jutil.FarmHashFingerprint64(buf[:])))
	}
	var stackBytes [256]byte
	var raw []byte
	if n*8 <= len(stackBytes) {
		raw = stackBytes[:n*8]
	} else {
		raw = make([]byte, n*8)
	}
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint64(raw[i*8:], uint64(ids[i]))
	}
	return NewColorId(int64(jutil.FarmHashFingerprint64(raw)))
}

// Java: ColorId(long)
func NewColorId(rightAligned int64) ColorId {
	return ColorId{rightAligned: rightAligned}
}

// Java: ColorId.equals(Object)
func (c ColorId) Equals(x any) bool {
	that, ok := x.(ColorId)
	if !ok {
		if p, isPtr := x.(*ColorId); isPtr && p != nil {
			that, ok = *p, true
		}
	}
	if !ok {
		return false
	}

	return c.rightAligned == that.rightAligned
}

// Java: ColorId.hashCode()
func (c ColorId) HashCode() int32 {
	return int32(c.rightAligned)
}

// AsByteString returns the eight big-endian bytes of the ID (protobuf
// ByteString as Go []byte).
//
// Java: ColorId.asByteString()
func (c ColorId) AsByteString() []byte {
	out := make([]byte, 8)
	binary.BigEndian.PutUint64(out, uint64(c.rightAligned))
	return out
}

// PutBigEndian writes the eight big-endian bytes of the ID into dst[:8].
func (c ColorId) PutBigEndian(dst []byte) {
	binary.BigEndian.PutUint64(dst, uint64(c.rightAligned))
}

// Java: ColorId.toString()
func (c ColorId) String() string {
	// Use hex encoding so there are no ambiguous trailing chars like Base64
	// would have.
	return jutil.LongToHexString(c.rightAligned)
}

// Java: ColorId.toLogsGson()
func (c ColorId) ToLogsGson() string {
	return c.String()
}

// ColorId_ByteAt reads byte `index` of `source`.
//
// Java: ColorId.ByteAt<T>
type ColorId_ByteAt[T any] func(source T, index int) int8

// Java: ColorId.ASCII_AT
var ColorId_ASCII_AT ColorId_ByteAt[string] = func(s string, pos int) int8 {
	if pos < 0 || pos >= len(s) {
		panic(jutil.NewStringIndexOutOfBoundsException(pos, len(s)))
	}
	c := int(s[pos])
	jutil.CheckStateMsg1(c < 128, "%s", c)
	return int8(c)
}
