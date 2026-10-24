// Port of com/google/javascript/jscomp/Xid.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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

package jscomp

import (
	"math"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Xid is a simple utility for shortening identifiers in a stable way.
// Generates short substitution strings deterministically, using a compact (1
// to 6 characters in length) repesentation of a 32-bit hash of the key. The
// string is suitable to be used as a JavaScript or CSS identifier.
// Collisions are possible but unlikely, depending on the underlying hash
// algorithm used.
//
// This substitution scheme uses case-sensitive names for maximum
// compression. Digits are also allowed in all but the first character of a
// class name. There are a few characters allowed by the CSS grammar that we
// choose not to use (e.g. the underscore and hyphen), to keep names simple.
//
// Xid should maintain as minimal dependencies as possible to ease its
// integration with other tools, such as server side HTML generators.
type Xid struct {
	hasher Xid_HashFunction
}

const (
	// Possible first chars in an identifier
	xid_START_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Possible non-first chars in an identifier
	xid_CHARS = xid_START_CHARS + "0123456789"

	xid_START_RADIX = len(xid_START_CHARS)

	xid_RADIX = len(xid_CHARS)
)

// Xid_HashFunction is the strategy for selecting the underlying hash code
// function to be used by Xid.
//
// Go-only: language-adapter — a functional interface, so a Go func type
// (Java: int hashCode(String value)).
type Xid_HashFunction func(value string) int32

// Java: Xid()
func NewXid() *Xid {
	return &Xid{hasher: jutil.StringHashCode}
}

// Java: Xid(HashFunction)
func NewXidWithHasher(hasher Xid_HashFunction) *Xid {
	return &Xid{hasher: hasher}
}

// Get gets the string that should be substituted for key. The same value
// will be consistently returned for any particular key.
//
// key is the text to be replaced (never null).
//
// Java: Xid.get(String)
func (x *Xid) Get(key string) string {
	return Xid_ToString(x.GetAsInt(key))
}

// GetAsInt returns the underlying integer representation of the given key.
//
// Java: Xid.getAsInt(String)
func (x *Xid) GetAsInt(key string) int32 {
	return x.hasher(key)
}

// Xid_ToString converts a 32-bit integer to a unique short string whose
// first character is in START_CHARS and whose subsequent characters, if
// any, are in CHARS. The result is 1-6 characters in length.
//
// Java: Xid.toString(int)
func Xid_ToString(i int32) string {
	buf := make([]byte, 6)
	length := 0
	l := int64(i) - int64(math.MinInt32)
	buf[length] = xid_START_CHARS[l%int64(xid_START_RADIX)]
	length++
	i = int32(l / int64(xid_START_RADIX))
	for i > 0 {
		buf[length] = xid_CHARS[i%int32(xid_RADIX)]
		length++
		i /= int32(xid_RADIX)
	}
	return string(buf[:length])
}
