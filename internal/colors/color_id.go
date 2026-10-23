// Port of com/google/javascript/jscomp/colors/ColorId.java @ closure-compiler v20261015

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
	"slices"

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

// Java: ColorId.fromUnsigned(byte)
func ColorId_FromUnsignedByte(x int8) ColorId {
	return NewColorId(int64(x) & 0xFF)
}

// ColorId_FromBytes builds an ID from a protobuf ByteString (Go []byte).
//
// Java: ColorId.fromBytes(ByteString)
func ColorId_FromBytes(bytes []byte) ColorId {
	return ColorId_FromBytesAt(bytes, len(bytes), func(s []byte, i int) int8 { return int8(s[i]) })
}

// Java: ColorId.fromBytes(byte[])
func ColorId_FromByteArray(bytes []byte) ColorId {
	return ColorId_FromBytesAt(bytes, len(bytes), func(s []byte, i int) int8 { return int8(s[i]) })
}

// Java: ColorId.fromAscii(String)
func ColorId_FromAscii(str string) ColorId {
	return ColorId_FromBytesAt(str, jutil.StringLength(str), ColorId_ASCII_AT)
}

// Java: ColorId.fromBytesAt(T, int, ByteAt<T>)
func ColorId_FromBytesAt[T any](source T, length int, byteAt ColorId_ByteAt[T]) ColorId {
	jutil.CheckStateMsg1(length <= 8, "%s", length)

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
	n := ids.Size()
	if n <= 1 {
		return iterablesGetOnlyElement(ids)
	}
	if n == 2 {
		var a, b int64
		first := true
		for id := range ids.All() {
			if first {
				a = id.rightAligned
				first = false
			} else {
				b = id.rightAligned
			}
		}
		if a > b {
			a, b = b, a
		}
		var hasher jutil.FarmHashFingerprint64Hasher
		hasher.Reset()
		hasher.PutLong(a)
		hasher.PutLong(b)
		return NewColorId(hasher.HashAsLong())
	}

	var stackBuf [32]int64
	var sorted []int64
	if n <= len(stackBuf) {
		sorted = stackBuf[:n]
	} else {
		sorted = make([]int64, n)
	}
	index := 0
	for id := range ids.All() {
		sorted[index] = id.rightAligned
		index++
	}
	slices.Sort(sorted)

	var hasher jutil.FarmHashFingerprint64Hasher
	hasher.Reset()
	for i := 0; i < len(sorted); i++ {
		hasher.PutLong(sorted[i])
	}

	return NewColorId(hasher.HashAsLong())
}

// iterablesGetOnlyElement is Guava's Iterables.getOnlyElement(Iterable).
func iterablesGetOnlyElement[T comparable](s *jutil.LinkedHashSet[T]) T {
	if s.Size() == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	if s.Size() > 1 {
		panic(jutil.NewIllegalArgumentExceptionMsg("expected one element but was: " + s.String()))
	}
	for e := range s.All() {
		return e
	}
	panic(jutil.NewNoSuchElementException())
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
	cp := c.rightAligned
	out := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		out[i] = byte(cp)
		cp = jutil.Ushr64(cp, 8)
	}
	return out
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

// ColorId_FARM_64 is Hashing.farmHashFingerprint64(), represented by its
// newHasher() method.
//
// Java: ColorId.FARM_64
var ColorId_FARM_64 = jutil.Hashing_FarmHashFingerprint64NewHasher

// ColorId_ByteAt reads byte `index` of `source`.
//
// Java: ColorId.ByteAt<T>
type ColorId_ByteAt[T any] func(source T, index int) int8

// Java: ColorId.ASCII_AT
var ColorId_ASCII_AT ColorId_ByteAt[string] = func(s string, i int) int8 {
	c := jutil.StringCodePointAt(s, i)
	jutil.CheckStateMsg1(0 <= c && c < 128, "%s", c)
	return int8(c)
}
