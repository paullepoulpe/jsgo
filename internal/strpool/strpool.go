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
	"sync/atomic"
)

// ID is a 32-bit pointer-free handle to an interned string:
//   - 0: null / unset
//   - 1: "" (empty string)
//   - >= 2: ((localIdx << 6) | shardIdx) + 2
type ID uint32

const (
	numShards  = 64
	chunkShift = 9
	chunkSize  = 1 << chunkShift // 512 strings per chunk (power of 2 for 1-cycle >> / & indexing)
	chunkMask  = chunkSize - 1
	maxChunks  = 1024 // up to 33.5M unique strings across 64 shards
)

type shard struct {
	mu           sync.RWMutex
	table        []uint64 // open-addressed hash table, at most 3/4 full: (hash32 << 32) | (localIdx + 1)
	count        uint32
	nextLocalIdx uint32
	chunks       [maxChunks]atomic.Pointer[[chunkSize]string]
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

func (s *shard) stringAt(localIdx uint32) string {
	return s.chunks[localIdx>>chunkShift].Load()[localIdx&chunkMask]
}

func (s *shard) allocString(str string) uint32 {
	localIdx := s.nextLocalIdx
	s.nextLocalIdx = localIdx + 1
	chunkIdx := localIdx >> chunkShift
	chunk := s.chunks[chunkIdx].Load()
	if chunk == nil {
		chunk = new([chunkSize]string)
		s.chunks[chunkIdx].Store(chunk)
	}
	// Clone so substring slices from parser inputs never pin full source files.
	chunk[localIdx&chunkMask] = strings.Clone(str)
	return localIdx
}

func (s *shard) growTable() {
	newCap := max(len(s.table)*2, 128)
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

func (s *shard) probeLocked(hash32 uint32, str string) (uint32, bool) {
	tbl := s.table
	if len(tbl) == 0 {
		return 0, false
	}
	mask := uint32(len(tbl) - 1)
	slot := hash32 & mask
	for {
		entry := tbl[slot]
		if entry == 0 {
			return 0, false
		}
		if uint32(entry>>32) == hash32 {
			localIdx := uint32(entry) - 1
			if s.stringAt(localIdx) == str {
				return localIdx, true
			}
		}
		slot = (slot + 1) & mask
	}
}

func makeID(localIdx, shardIdx uint32) ID {
	return ID(((localIdx << 6) | shardIdx) + 2)
}

func Intern(str string) ID {
	if len(str) == 0 {
		return 1
	}
	h := hash(str)
	shardIdx := uint32(h) & (numShards - 1)
	hash32 := uint32(h >> 32)
	s := &shards[shardIdx]

	s.mu.RLock()
	if localIdx, ok := s.probeLocked(hash32, str); ok {
		s.mu.RUnlock()
		return makeID(localIdx, shardIdx)
	}
	s.mu.RUnlock()

	s.mu.Lock()
	if len(s.table) == 0 {
		s.table = make([]uint64, 128)
	} else if localIdx, ok := s.probeLocked(hash32, str); ok {
		s.mu.Unlock()
		return makeID(localIdx, shardIdx)
	}
	if int(s.count+1)*4 > len(s.table)*3 {
		s.growTable()
	}
	localIdx := s.allocString(str)
	mask := uint32(len(s.table) - 1)
	slot := hash32 & mask
	for s.table[slot] != 0 {
		slot = (slot + 1) & mask
	}
	s.table[slot] = (uint64(hash32) << 32) | uint64(localIdx+1)
	s.count++
	s.mu.Unlock()
	return makeID(localIdx, shardIdx)
}

func String(id ID) string {
	if id <= 1 {
		return ""
	}
	idx := uint32(id - 2)
	s := &shards[idx&(numShards-1)]
	localIdx := idx >> 6
	return s.chunks[localIdx>>chunkShift].Load()[localIdx&chunkMask]
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
	return &s.chunks[localIdx>>chunkShift].Load()[localIdx&chunkMask]
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
	s := &shards[shardIdx]

	s.mu.RLock()
	localIdx, ok := s.probeLocked(hash32, str)
	s.mu.RUnlock()
	if !ok {
		return 0
	}
	return makeID(localIdx, shardIdx)
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
	s := &shards[shardIdx]

	s.mu.RLock()
	if tbl := s.table; len(tbl) > 0 {
		mask := uint32(len(tbl) - 1)
		slot := hash32 & mask
		for {
			entry := tbl[slot]
			if entry == 0 {
				break
			}
			if uint32(entry>>32) == hash32 {
				localIdx := uint32(entry) - 1
				if s.stringAt(localIdx) == string(b) {
					s.mu.RUnlock()
					return makeID(localIdx, shardIdx)
				}
			}
			slot = (slot + 1) & mask
		}
	}
	s.mu.RUnlock()
	return Intern(string(b))
}

// WasClearedForTest reports whether ClearForTest has been called.
func WasClearedForTest() bool { return cleared }
