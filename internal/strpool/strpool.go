// Port of the storage of com/google/javascript/rhino/RhinoStringPool.java @ closure-compiler v20261104

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

// EmptyID is the canonical ID of the empty string "".
const EmptyID = 1

const (
	numShards   = 64
	chunkShift  = 9
	chunkSize   = 1 << chunkShift // 512 strings per chunk (power of 2 for 1-cycle >> / & indexing)
	chunkMask   = chunkSize - 1
	maxChunks   = 1024 // up to 33.5M unique strings across 64 shards
	fnvOffset64 = 14695981039346656037
	fnvPrime64  = 1099511628211
)

// HashString returns the 64-bit FNV-1a hash of s.
func HashString(s string) uint64 {
	h := uint64(fnvOffset64)
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= fnvPrime64
	}
	return h
}

// HashBytes returns the 64-bit FNV-1a hash of b.
func HashBytes(b []byte) uint64 {
	h := uint64(fnvOffset64)
	for _, c := range b {
		h ^= uint64(c)
		h *= fnvPrime64
	}
	return h
}

type shard struct {
	mu           sync.RWMutex
	table        []uint64 // open-addressed hash table, at most 3/4 full: (hash32 << 32) | (localIdx + 1)
	count        uint32
	nextLocalIdx uint32
	chunks       [maxChunks]atomic.Pointer[[chunkSize]string]
}

var (
	shards      [numShards]shard
	seed        = maphash.MakeSeed()
	rewindEpoch atomic.Uint64
)

// BaselineSnapshot records the per-shard watermark of interned strings at a
// post-init baseline so RewindBaseline can reclaim per-request strings.
type BaselineSnapshot [numShards]uint32

// Epoch returns a monotonically increasing generation bumped whenever
// RewindBaseline rewinds any shard, allowing pooled L1 ID caches (e.g.
// IRFactory) to invalidate stale IDs across worker requests.
func Epoch() uint64 {
	return rewindEpoch.Load()
}

// CaptureBaseline snapshots the current per-shard string counts.
func CaptureBaseline() BaselineSnapshot {
	var snap BaselineSnapshot
	for i := range shards {
		s := &shards[i]
		s.mu.RLock()
		snap[i] = s.nextLocalIdx
		s.mu.RUnlock()
	}
	return snap
}

// RewindBaseline rewinds every shard back to snap, clearing post-baseline
// string references, releasing chunks allocated past snap, and shrinking
// bloated hash tables. Callers must invoke RewindBaseline only when no
// concurrent compilation is reading or interning post-baseline IDs.
func RewindBaseline(snap BaselineSnapshot) {
	rewound := false
	for i := range shards {
		if shards[i].rewind(snap[i]) {
			rewound = true
		}
	}
	if rewound {
		rewindEpoch.Add(1)
	}
}

func (s *shard) rewind(baseIdx uint32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	curIdx := s.nextLocalIdx
	if curIdx <= baseIdx {
		return false
	}

	firstDropChunk := (baseIdx + chunkMask) >> chunkShift
	lastUsedChunk := (curIdx - 1) >> chunkShift
	for c := firstDropChunk; c <= lastUsedChunk; c++ {
		s.chunks[c].Store(nil)
	}
	if baseIdx&chunkMask != 0 {
		tailChunkIdx := baseIdx >> chunkShift
		if chunk := s.chunks[tailChunkIdx].Load(); chunk != nil {
			endOff := min(uint32(chunkSize), curIdx-(tailChunkIdx<<chunkShift))
			clear(chunk[baseIdx&chunkMask : endOff])
		}
	}

	s.nextLocalIdx = baseIdx
	s.count = baseIdx

	if baseIdx == 0 {
		if len(s.table) > 128 {
			s.table = nil
		} else {
			clear(s.table)
		}
		return true
	}

	wantCap := 128
	for int(baseIdx)*4 > wantCap*3 {
		wantCap *= 2
	}
	if len(s.table) > wantCap*2 {
		newTable := make([]uint64, wantCap)
		mask := uint32(wantCap - 1)
		for _, entry := range s.table {
			if entry == 0 || uint32(entry)-1 >= baseIdx {
				continue
			}
			slot := uint32(entry>>32) & mask
			for newTable[slot] != 0 {
				slot = (slot + 1) & mask
			}
			newTable[slot] = entry
		}
		s.table = newTable
		return true
	}

	var stackBuf [256]uint64
	survivors := stackBuf[:0]
	if int(baseIdx) > len(stackBuf) {
		survivors = make([]uint64, 0, baseIdx)
	}
	for _, entry := range s.table {
		if entry != 0 && uint32(entry)-1 < baseIdx {
			survivors = append(survivors, entry)
		}
	}
	clear(s.table)
	mask := uint32(len(s.table) - 1)
	for _, entry := range survivors {
		slot := uint32(entry>>32) & mask
		for s.table[slot] != 0 {
			slot = (slot + 1) & mask
		}
		s.table[slot] = entry
	}
	return true
}

func hash(s string) uint64 {
	return maphash.String(seed, s)
}

func (s *shard) stringAt(localIdx uint32) string {
	return s.chunks[localIdx>>chunkShift].Load()[localIdx&chunkMask]
}

func (s *shard) allocString(str string) uint32 {
	localIdx := s.nextLocalIdx
	chunkIdx := localIdx >> chunkShift
	if int(chunkIdx) >= len(s.chunks) {
		panic("strpool: shard capacity exceeded")
	}
	s.nextLocalIdx = localIdx + 1
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

func probeLocked[S ~string | ~[]byte](s *shard, hash32 uint32, key S) (uint32, bool) {
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
			if s.stringAt(localIdx) == string(key) {
				return localIdx, true
			}
		}
		slot = (slot + 1) & mask
	}
}

func makeID(localIdx, shardIdx uint32) ID {
	return ID(((localIdx << 6) | shardIdx) + 2)
}

func (s *shard) internSlow(shardIdx, hash32 uint32, str string) ID {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.table) == 0 {
		s.table = make([]uint64, 128)
	} else if localIdx, ok := probeLocked(s, hash32, str); ok {
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
	return makeID(localIdx, shardIdx)
}

func Intern(str string) ID {
	if len(str) == 0 {
		return EmptyID
	}
	h := hash(str)
	shardIdx := uint32(h) & (numShards - 1)
	hash32 := uint32(h >> 32)
	s := &shards[shardIdx]

	s.mu.RLock()
	if localIdx, ok := probeLocked(s, hash32, str); ok {
		s.mu.RUnlock()
		return makeID(localIdx, shardIdx)
	}
	s.mu.RUnlock()
	return s.internSlow(shardIdx, hash32, str)
}

func String(id ID) string {
	if id <= EmptyID {
		return ""
	}
	idx := uint32(id - 2)
	s := &shards[idx&(numShards-1)]
	localIdx := idx >> 6
	chunk := s.chunks[localIdx>>chunkShift].Load()
	if chunk == nil {
		return ""
	}
	return chunk[localIdx&chunkMask]
}

func Lookup(str string) ID {
	if len(str) == 0 {
		return EmptyID
	}
	h := hash(str)
	shardIdx := uint32(h) & (numShards - 1)
	hash32 := uint32(h >> 32)
	s := &shards[shardIdx]

	s.mu.RLock()
	localIdx, ok := probeLocked(s, hash32, str)
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
		return EmptyID
	}
	h := maphash.Bytes(seed, b)
	shardIdx := uint32(h) & (numShards - 1)
	hash32 := uint32(h >> 32)
	s := &shards[shardIdx]

	s.mu.RLock()
	if localIdx, ok := probeLocked(s, hash32, b); ok {
		s.mu.RUnlock()
		return makeID(localIdx, shardIdx)
	}
	s.mu.RUnlock()
	return s.internSlow(shardIdx, hash32, string(b))
}
