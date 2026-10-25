// Port of the storage of com/google/javascript/rhino/RhinoStringPool.java @ closure-compiler v20261021

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

// Package strpool is the process-wide interning pool of strings behind
// rhino.RhinoStringPool: an interned string is a 32-bit ID.
//
// Go-only: it is its own package, below rhino, so that the parser can intern
// identifier names too (parser.IdentifierToken keeps an ID, not a string).
package strpool

import (
	"hash/maphash"
	"strings"
	"sync"
	"unsafe"
)

// ID is a 32-bit pointer-free handle to an interned string:
//   - 0: null / unset
//   - 1: "" (empty string)
//   - >= 2: ((localIdx << 6) | shardIdx) + 2
type ID uint32

const (
	numShards     = 64
	maxShardSize  = 56000
	pageSize      = 32760 // <= maxSmallSize - gc.MallocHeaderSize
	spanChunkSize = 2048  // 16384B <= 32760B (power of 2)
)

type shard struct {
	mu           sync.RWMutex
	table        []uint64 // noscan open-addressed hash table, at most 3/4 full: (hash32 << 32) | (localIdx + 1)
	count        uint32
	nextLocalIdx uint32
	curPageIdx   uint32
	curPageOff   uint32
	hasCurPage   bool
	spanChunks   [2048]*[spanChunkSize]uint64
	pages        [4096]*[pageSize]byte
	largeStrings []string
	ptrMap       map[uint32]*string
	ptrSlab      []string
}

var (
	shards      [numShards]shard
	cleared     bool
	emptyString = ""
	seed        = maphash.MakeSeed()
)

func hash(s string) uint64 {
	return maphash.String(seed, s)
}

func (s *shard) stringAtLocked(localIdx uint32) string {
	desc := s.spanChunks[localIdx>>11][localIdx&2047]
	if int64(desc) < 0 {
		return s.largeStrings[uint32(desc)]
	}
	pageOff := uint32(desc & 0xFFFFF)
	length := uint32((desc >> 20) & 0xFFFFF)
	pageIdx := uint32((desc >> 40) & 0xFFFFF)
	return unsafe.String(&s.pages[pageIdx][pageOff], length)
}

func (s *shard) allocSpan(str string) uint32 {
	localIdx := s.nextLocalIdx
	s.nextLocalIdx = localIdx + 1
	chunkIdx := localIdx >> 11
	if s.spanChunks[chunkIdx] == nil {
		s.spanChunks[chunkIdx] = new([spanChunkSize]uint64)
	}
	n := len(str)
	if n > pageSize || (s.hasCurPage && int(s.curPageOff)+n > pageSize && int(s.curPageIdx)+1 >= len(s.pages)) {
		largeIdx := uint32(len(s.largeStrings))
		// Copy string bytes so substring slices don't pin large source buffers.
		s.largeStrings = append(s.largeStrings, strings.Clone(str))
		s.spanChunks[chunkIdx][localIdx&2047] = (uint64(1) << 63) | uint64(largeIdx)
		return localIdx
	}
	if !s.hasCurPage || int(s.curPageOff)+n > pageSize {
		if s.hasCurPage {
			s.curPageIdx++
		} else {
			s.hasCurPage = true
		}
		s.curPageOff = 0
		s.pages[s.curPageIdx] = new([pageSize]byte)
	}
	off := s.curPageOff
	s.curPageOff = off + uint32(n)
	copy(s.pages[s.curPageIdx][off:off+uint32(n)], str)
	desc := uint64(off) | (uint64(n) << 20) | (uint64(s.curPageIdx) << 40)
	s.spanChunks[chunkIdx][localIdx&2047] = desc
	return localIdx
}

func (s *shard) growTable() {
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

func Intern(str string) ID {
	if len(str) == 0 {
		return 1
	}
	h := hash(str)
	shardIdx := uint32(h) & (numShards - 1)
	hash32 := uint32(h >> 32)
	shard := &shards[shardIdx]

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
					return ID(((localIdx << 6) | shardIdx) + 2)
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
					return ID(((localIdx << 6) | shardIdx) + 2)
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
	return ID(((localIdx << 6) | shardIdx) + 2)
}

//go:noinline
func (s *shard) largeStringAt(largeIdx uint32) string {
	s.mu.RLock()
	v := s.largeStrings[largeIdx]
	s.mu.RUnlock()
	return v
}

func String(id ID) string {
	if id <= 1 {
		return ""
	}
	idx := uint32(id - 2)
	s := &shards[idx&(numShards-1)]
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

func StringPtr(id ID) *string {
	if id == 0 {
		return nil
	}
	if id == 1 {
		return &emptyString
	}
	idx := uint32(id - 2)
	s := &shards[idx&(numShards-1)]
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

func ClearForTest() {
	cleared = true
	for i := range shards {
		s := &shards[i]
		s.mu.Lock()
		clear(s.table)
		s.count = 0
		s.mu.Unlock()
	}
}

func Lookup(str string) ID {
	if len(str) == 0 {
		return 1
	}
	h := hash(str)
	shardIdx := uint32(h) & (numShards - 1)
	hash32 := uint32(h >> 32)
	shard := &shards[shardIdx]

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
					return ID(((localIdx << 6) | shardIdx) + 2)
				}
			}
			slot = (slot + 1) & mask
		}
	}
	shard.mu.RUnlock()
	return 0
}

// InternBytes interns b and returns its ID, without allocating a string
// when b is already interned.
func InternBytes(b []byte) ID {
	if len(b) == 0 {
		return 1
	}
	h := maphash.Bytes(seed, b)
	shardIdx := uint32(h) & (numShards - 1)
	hash32 := uint32(h >> 32)
	shard := &shards[shardIdx]

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
				if shard.stringAtLocked(localIdx) == string(b) {
					shard.mu.RUnlock()
					return ID(((localIdx << 6) | shardIdx) + 2)
				}
			}
			slot = (slot + 1) & mask
		}
	}
	shard.mu.RUnlock()
	return Intern(string(b))
}

// WasClearedForTest reports whether ClearForTest has been called.
func WasClearedForTest() bool { return cleared }
