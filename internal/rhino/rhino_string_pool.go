// Port of com/google/javascript/rhino/RhinoStringPool.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"iter"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// RhinoStringPool is an interning pool for strings used by the Rhino
// package.
//
// As of 2021-03-16 and again as of 2024-12-09, this custom pool is
// measurably more performant than `String::intern`.
type RhinoStringPool struct{}

// rhinoStringID is a 32-bit pointer-free handle to an interned string (see
// package strpool): 0 is null / unset, 1 is "".
type rhinoStringID = strpool.ID

func rhinoStringPoolIntern(s string) rhinoStringID  { return strpool.Intern(s) }
func rhinoStringFromID(id rhinoStringID) string     { return strpool.String(id) }
func rhinoStringPtrFromID(id rhinoStringID) *string { return strpool.StringPtr(id) }
func rhinoStringPoolLookup(s string) rhinoStringID  { return strpool.Lookup(s) }
func rhinoStringPoolClearForTest()                  { strpool.ClearForTest() }
func rhinoStringPoolWasCleared() bool               { return strpool.WasClearedForTest() }

// RhinoStringPool_UncheckedEquals checks if two strings are the same
// according to interning.
//
// The caller is responsible for ensuring that strings passed to this method
// are actually interned. This method mainly exists to highlight where
// equality checks depend on interning for correctness.
//
// Ideally we could use something like a branded-type, for interned strings,
// to verify correct usage. However, Java doesn't support type-brands, and
// using like wrapper objects would undermine performance. This also needs to
// be ergonomic enough that callers don't resort to using `==` directly.
//
// Go-only: language-adapter — Java compares the references; Go compares the
// string values, which is the same for interned strings.
//
// Java: RhinoStringPool.uncheckedEquals(String, String)
func RhinoStringPool_UncheckedEquals(a string, b string) bool {
	return base.JSCompObjects_Identical(a, b)
}

// RhinoStringPool_UncheckedEqualsNullable is uncheckedEquals for Java
// Strings that may be null (*string in the port): null is identical only to
// null.
//
// Go-only: language-adapter — nullable String overload of
// RhinoStringPool.uncheckedEquals(String, String).
func RhinoStringPool_UncheckedEqualsNullable(a *string, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return RhinoStringPool_UncheckedEquals(*a, *b)
}

// Java: RhinoStringPool.addOrGet(String)
func RhinoStringPool_AddOrGet(s string) string {
	return rhinoStringFromID(rhinoStringPoolIntern(s))
}

// RhinoStringPool_EmptyID is the canonical 32-bit ID of the empty string "".
const RhinoStringPool_EmptyID uint32 = 1

// RhinoStringPool_InternID interns s and returns its 32-bit pointer-free ID.
func RhinoStringPool_InternID(s string) uint32 {
	return uint32(rhinoStringPoolIntern(s))
}

// RhinoStringPool_InternBytesID interns b and returns its 32-bit pointer-free ID
// without allocating a string when b is already interned.
func RhinoStringPool_InternBytesID(b []byte) uint32 {
	return uint32(strpool.InternBytes(b))
}

// RhinoStringPool_LookupID returns the 32-bit ID of s if already interned, or 0 if not.
func RhinoStringPool_LookupID(s string) uint32 {
	return uint32(rhinoStringPoolLookup(s))
}

// RhinoStringPool_StringFromID returns the interned string for a 32-bit ID.
func RhinoStringPool_StringFromID(id uint32) string {
	return rhinoStringFromID(rhinoStringID(id))
}

// RhinoStringPool_WasClearedForTest reports whether rhinoStringPoolClearForTest has been called.
func RhinoStringPool_WasClearedForTest() bool {
	return rhinoStringPoolWasCleared()
}

// Java: RhinoStringPool()
func NewRhinoStringPool() *RhinoStringPool { return &RhinoStringPool{} }

// RhinoStringPool_LazyInternedStringList is a list of strings that is lazily
// interned into a RhinoStringPool as they are accessed.
//
// The only reason to use this class is for performance optimizations - the
// LazyInternedStringList#newStringNode method skips re-interning the string.
// This is a performance win if creating a lot of string nodes with the same
// string.
type RhinoStringPool_LazyInternedStringList struct {
	pool       *jutil.ArrayList[string]
	ids        []uint32
	isInterned *RhinoStringPool_WriteOnlyBitset
}

// Java: RhinoStringPool.LazyInternedStringList(List<String>)
func NewRhinoStringPool_LazyInternedStringList(pool []string) *RhinoStringPool_LazyInternedStringList {
	return &RhinoStringPool_LazyInternedStringList{
		pool:       jutil.NewArrayListFrom(pool),
		ids:        make([]uint32, len(pool)),
		isInterned: NewRhinoStringPool_WriteOnlyBitset(len(pool)),
	}
}

// NewRhinoStringPool_LazyInternedStringListFromIDs creates a LazyInternedStringList
// backed directly by pre-interned 32-bit rhinoStringIDs (noscan, zero pointers).
func NewRhinoStringPool_LazyInternedStringListFromIDs(ids []uint32) *RhinoStringPool_LazyInternedStringList {
	if len(ids) == 0 {
		_ = NewRhinoStringPool_WriteOnlyBitset(0)
	}
	return &RhinoStringPool_LazyInternedStringList{
		ids: ids,
	}
}

//go:noinline
func (l *RhinoStringPool_LazyInternedStringList) getIDSlow(offset int) rhinoStringID {
	if l.pool == nil {
		wordIdx := RhinoStringPool_WriteOnlyBitset_WordIndex(offset)
		wordLen := RhinoStringPool_WriteOnlyBitset_WordIndex(len(l.ids)-1) + 1
		panic(jutil.NewArrayIndexOutOfBoundsExceptionForLength(wordIdx, wordLen))
	}
	if uint(offset) < uint(len(l.ids)) {
		if id := atomic.LoadUint32(&l.ids[offset]); id != 0 {
			return rhinoStringID(id)
		}
	} else {
		_ = l.isInterned.Get(offset)
	}
	uninterned := l.pool.Get(offset)
	id := rhinoStringPoolIntern(uninterned)
	atomic.StoreUint32(&l.ids[offset], uint32(id))
	l.isInterned.Set(offset)
	return id
}

// GetID returns the interned rhinoStringID at offset, interning on first access.
func (l *RhinoStringPool_LazyInternedStringList) GetID(offset int) rhinoStringID {
	if l.pool == nil && uint(offset) < uint(len(l.ids)) {
		return rhinoStringID(l.ids[offset])
	}
	return l.getIDSlow(offset)
}

// Get returns the interned string at offset.
//
// if Java supported type-brands we would want to use one here to indicate
// that the returned string is always interned but that's not possible
// without an expensive wrapper type - see also the comment on
// uncheckedEquals().
//
// Java: RhinoStringPool.LazyInternedStringList.get(int)
func (l *RhinoStringPool_LazyInternedStringList) Get(offset int) string {
	return rhinoStringFromID(l.GetID(offset))
}

// Java: RhinoStringPool.LazyInternedStringList.stream()
func (l *RhinoStringPool_LazyInternedStringList) Stream() iter.Seq[string] {
	if l.pool != nil {
		return l.pool.All()
	}
	return func(yield func(string) bool) {
		for _, id := range l.ids {
			if !yield(rhinoStringFromID(rhinoStringID(id))) {
				return
			}
		}
	}
}

// Slice returns the underlying string slice.
func (l *RhinoStringPool_LazyInternedStringList) Slice() []string {
	if l.pool != nil {
		return l.pool.Slice()
	}
	out := make([]string, len(l.ids))
	for i, id := range l.ids {
		out[i] = rhinoStringFromID(rhinoStringID(id))
	}
	return out
}

// Len returns the number of strings in the list.
func (l *RhinoStringPool_LazyInternedStringList) Len() int {
	return len(l.ids)
}

// RhinoStringPool_WriteOnlyBitset is a bitset implementation that only
// supports setting bits from false => true, not vice versa, and is not
// resizable once created.
//
// This class was originally copied from the JDK BitSet implementation, but
// has been heavily modified to be thread-safe and remove unnecessary
// functionality.
//
// Go-only: language-adapter — the AtomicIntegerArray is a plain []int32
// (the port is single-threaded; see PORTING.md §5); out-of-range accesses
// throw ArrayIndexOutOfBoundsException like the JDK's.
type RhinoStringPool_WriteOnlyBitset struct {
	bits []int32
}

// BitSets are packed into arrays of "words." Currently a word is an int,
// which consists of 32 bits, requiring 5 address bits. The choice of word
// size comes from the convenience of using AtomicIntegerArray.
const RhinoStringPool_WriteOnlyBitset_ADDRESS_BITS_PER_WORD = 5

// Java: RhinoStringPool.WriteOnlyBitset(int)
func NewRhinoStringPool_WriteOnlyBitset(size int) *RhinoStringPool_WriteOnlyBitset {
	n := RhinoStringPool_WriteOnlyBitset_WordIndex(size-1) + 1
	if n < 0 {
		panic(jutil.NewRuntimeException("java.lang.NegativeArraySizeException: " + jutil.ToJavaString(n)))
	}
	return &RhinoStringPool_WriteOnlyBitset{bits: make([]int32, n)}
}

func (b *RhinoStringPool_WriteOnlyBitset) checkIndex(wordIndex int) {
	if wordIndex < 0 || wordIndex >= len(b.bits) {
		panic(jutil.NewArrayIndexOutOfBoundsExceptionForLength(wordIndex, len(b.bits)))
	}
}

// Java: RhinoStringPool.WriteOnlyBitset.set(int)
func (b *RhinoStringPool_WriteOnlyBitset) Set(offset int) {
	wordIndex := RhinoStringPool_WriteOnlyBitset_WordIndex(offset)
	b.checkIndex(wordIndex)
	atomic.OrInt32(&b.bits[wordIndex], jutil.Shl32(1, int64(offset)))
}

// Java: RhinoStringPool.WriteOnlyBitset.get(int)
func (b *RhinoStringPool_WriteOnlyBitset) Get(offset int) bool {
	wordIndex := RhinoStringPool_WriteOnlyBitset_WordIndex(offset)
	b.checkIndex(wordIndex)
	return (atomic.LoadInt32(&b.bits[wordIndex]) & jutil.Shl32(1, int64(offset))) != 0
}

// RhinoStringPool_WriteOnlyBitset_WordIndex returns, given a bit index, the
// word index containing it.
//
// Java: RhinoStringPool.WriteOnlyBitset.wordIndex(int)
func RhinoStringPool_WriteOnlyBitset_WordIndex(bitIndex int) int {
	return int(int32(bitIndex) >> RhinoStringPool_WriteOnlyBitset_ADDRESS_BITS_PER_WORD)
}
