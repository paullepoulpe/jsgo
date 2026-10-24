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
	"hash/maphash"
	"iter"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// RhinoStringPool is an interning pool for strings used by the Rhino
// package.
//
// As of 2021-03-16 and again as of 2024-12-09, this custom pool is
// measurably more performant than `String::intern`.
type RhinoStringPool struct{}

// rhinoStringID is a 32-bit pointer-free handle to an interned string:
//   - 0: null / unset
//   - 1: "" (empty string)
//   - >= 2: ((localIdx << 6) | shardIdx) + 2
type rhinoStringID uint32

const (
	rhinoStringPool_SHARDS         = 64
	rhinoStringPool_MAX_SHARD_SIZE = 56000
	rhinoStringPageSize            = 32760 // <= maxSmallSize - gc.MallocHeaderSize
	rhinoSpanChunkSize             = 2048  // 16384B <= 32760B (power of 2)
)

type rhinoStringPoolShard struct {
	mu           sync.RWMutex
	table        []uint64 // noscan open-addressed hash table, at most 3/4 full: (hash32 << 32) | (localIdx + 1)
	count        uint32
	nextLocalIdx uint32
	curPageIdx   uint32
	curPageOff   uint32
	hasCurPage   bool
	spanChunks   [2048]*[rhinoSpanChunkSize]uint64
	pages        [4096]*[rhinoStringPageSize]byte
	largeStrings []string
	ptrMap       map[uint32]*string
	ptrSlab      []string
}

var (
	rhinoStringPool_shards [rhinoStringPool_SHARDS]rhinoStringPoolShard
	rhinoStringPoolCleared bool
	rhinoEmptyString       = ""
	rhinoStringPoolSeed    = maphash.MakeSeed()
)

func rhinoStringPoolHash(s string) uint64 {
	return maphash.String(rhinoStringPoolSeed, s)
}

func (s *rhinoStringPoolShard) stringAtLocked(localIdx uint32) string {
	desc := s.spanChunks[localIdx>>11][localIdx&2047]
	if int64(desc) < 0 {
		return s.largeStrings[uint32(desc)]
	}
	pageOff := uint32(desc & 0xFFFFF)
	length := uint32((desc >> 20) & 0xFFFFF)
	pageIdx := uint32((desc >> 40) & 0xFFFFF)
	return unsafe.String(&s.pages[pageIdx][pageOff], length)
}

func (s *rhinoStringPoolShard) allocSpan(str string) uint32 {
	localIdx := s.nextLocalIdx
	s.nextLocalIdx = localIdx + 1
	chunkIdx := localIdx >> 11
	if s.spanChunks[chunkIdx] == nil {
		s.spanChunks[chunkIdx] = new([rhinoSpanChunkSize]uint64)
	}
	n := len(str)
	if n > rhinoStringPageSize || (s.hasCurPage && int(s.curPageOff)+n > rhinoStringPageSize && int(s.curPageIdx)+1 >= len(s.pages)) {
		largeIdx := uint32(len(s.largeStrings))
		// Copy string bytes so substring slices don't pin large source buffers.
		s.largeStrings = append(s.largeStrings, strings.Clone(str))
		s.spanChunks[chunkIdx][localIdx&2047] = (uint64(1) << 63) | uint64(largeIdx)
		return localIdx
	}
	if !s.hasCurPage || int(s.curPageOff)+n > rhinoStringPageSize {
		if s.hasCurPage {
			s.curPageIdx++
		} else {
			s.hasCurPage = true
		}
		s.curPageOff = 0
		s.pages[s.curPageIdx] = new([rhinoStringPageSize]byte)
	}
	off := s.curPageOff
	s.curPageOff = off + uint32(n)
	copy(s.pages[s.curPageIdx][off:off+uint32(n)], str)
	desc := uint64(off) | (uint64(n) << 20) | (uint64(s.curPageIdx) << 40)
	s.spanChunks[chunkIdx][localIdx&2047] = desc
	return localIdx
}

func (s *rhinoStringPoolShard) growTable() {
	newCap := len(s.table) * 2
	if newCap < 128 {
		newCap = 128
	}
	newTable := make([]uint64, newCap)
	mask := uint32(newCap - 1)
	for _, entry := range s.table {
		if entry == 0 {
			continue
		}
		h32 := uint32(entry >> 32)
		slot := h32 & mask
		for newTable[slot] != 0 {
			slot = (slot + 1) & mask
		}
		newTable[slot] = entry
	}
	s.table = newTable
}

func rhinoStringPoolIntern(str string) rhinoStringID {
	if len(str) == 0 {
		return 1
	}
	h := rhinoStringPoolHash(str)
	shardIdx := uint32(h) & (rhinoStringPool_SHARDS - 1)
	hash32 := uint32(h >> 32)
	shard := &rhinoStringPool_shards[shardIdx]

	shard.mu.RLock()
	if tbl := shard.table; len(tbl) > 0 {
		mask := uint32(len(tbl) - 1)
		slot := hash32 & mask
		for {
			entry := tbl[slot]
			if entry == 0 {
				break
			}
			if uint32(entry>>32) == hash32 {
				localIdx := uint32(entry) - 1
				if shard.stringAtLocked(localIdx) == str {
					shard.mu.RUnlock()
					return rhinoStringID(((localIdx << 6) | shardIdx) + 2)
				}
			}
			slot = (slot + 1) & mask
		}
	}
	shard.mu.RUnlock()

	shard.mu.Lock()
	if len(shard.table) == 0 {
		shard.table = make([]uint64, 128)
	} else {
		mask := uint32(len(shard.table) - 1)
		slot := hash32 & mask
		for {
			entry := shard.table[slot]
			if entry == 0 {
				break
			}
			if uint32(entry>>32) == hash32 {
				localIdx := uint32(entry) - 1
				if shard.stringAtLocked(localIdx) == str {
					shard.mu.Unlock()
					return rhinoStringID(((localIdx << 6) | shardIdx) + 2)
				}
			}
			slot = (slot + 1) & mask
		}
	}
	if int(shard.count+1)*4 > len(shard.table)*3 {
		shard.growTable()
	}
	localIdx := shard.allocSpan(str)
	mask := uint32(len(shard.table) - 1)
	slot := hash32 & mask
	for shard.table[slot] != 0 {
		slot = (slot + 1) & mask
	}
	shard.table[slot] = (uint64(hash32) << 32) | uint64(localIdx+1)
	shard.count++
	shard.mu.Unlock()
	return rhinoStringID(((localIdx << 6) | shardIdx) + 2)
}

//go:noinline
func (s *rhinoStringPoolShard) largeStringAt(largeIdx uint32) string {
	s.mu.RLock()
	v := s.largeStrings[largeIdx]
	s.mu.RUnlock()
	return v
}

func rhinoStringFromID(id rhinoStringID) string {
	if id <= 1 {
		return ""
	}
	idx := uint32(id - 2)
	s := &rhinoStringPool_shards[idx&(rhinoStringPool_SHARDS-1)]
	localIdx := idx >> 6
	desc := s.spanChunks[localIdx>>11][localIdx&2047]
	if int64(desc) < 0 {
		return s.largeStringAt(uint32(desc))
	}
	pageOff := uint32(desc & 0xFFFFF)
	length := uint32((desc >> 20) & 0xFFFFF)
	pageIdx := uint32((desc >> 40) & 0xFFFFF)
	return unsafe.String(&s.pages[pageIdx][pageOff], length)
}

func rhinoStringPtrFromID(id rhinoStringID) *string {
	if id == 0 {
		return nil
	}
	if id == 1 {
		return &rhinoEmptyString
	}
	idx := uint32(id - 2)
	s := &rhinoStringPool_shards[idx&(rhinoStringPool_SHARDS-1)]
	localIdx := idx >> 6

	s.mu.RLock()
	p := s.ptrMap[localIdx]
	s.mu.RUnlock()
	if p != nil {
		return p
	}

	s.mu.Lock()
	if s.ptrMap == nil {
		s.ptrMap = make(map[uint32]*string, 64)
	} else if p = s.ptrMap[localIdx]; p != nil {
		s.mu.Unlock()
		return p
	}
	if len(s.ptrSlab) == 0 {
		// 1023 strings and the 8-byte malloc header fill the 16 KiB size class.
		s.ptrSlab = make([]string, 1023)
	}
	p = &s.ptrSlab[0]
	s.ptrSlab = s.ptrSlab[1:]
	*p = s.stringAtLocked(localIdx)
	s.ptrMap[localIdx] = p
	s.mu.Unlock()
	return p
}

func rhinoStringPoolClearForTest() {
	rhinoStringPoolCleared = true
	for i := range rhinoStringPool_shards {
		s := &rhinoStringPool_shards[i]
		s.mu.Lock()
		clear(s.table)
		s.count = 0
		s.mu.Unlock()
	}
}

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

func rhinoStringPoolLookup(str string) rhinoStringID {
	if len(str) == 0 {
		return 1
	}
	h := rhinoStringPoolHash(str)
	shardIdx := uint32(h) & (rhinoStringPool_SHARDS - 1)
	hash32 := uint32(h >> 32)
	shard := &rhinoStringPool_shards[shardIdx]

	shard.mu.RLock()
	if tbl := shard.table; len(tbl) > 0 {
		mask := uint32(len(tbl) - 1)
		slot := hash32 & mask
		for {
			entry := tbl[slot]
			if entry == 0 {
				break
			}
			if uint32(entry>>32) == hash32 {
				localIdx := uint32(entry) - 1
				if shard.stringAtLocked(localIdx) == str {
					shard.mu.RUnlock()
					return rhinoStringID(((localIdx << 6) | shardIdx) + 2)
				}
			}
			slot = (slot + 1) & mask
		}
	}
	shard.mu.RUnlock()
	return 0
}

// RhinoStringPool_EmptyID is the canonical 32-bit ID of the empty string "".
const RhinoStringPool_EmptyID uint32 = 1

// RhinoStringPool_InternID interns s and returns its 32-bit pointer-free ID.
func RhinoStringPool_InternID(s string) uint32 {
	return uint32(rhinoStringPoolIntern(s))
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
	return rhinoStringPoolCleared
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
