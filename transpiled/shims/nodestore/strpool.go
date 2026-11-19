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

package nodestore

import (
	"hash/maphash"
	"strings"
	"sync"
	"sync/atomic"
)

// StringID is a 32-bit handle to an interned string in the concurrent string pool:
//   - 0: null / unset
//   - 1: "" (empty string)
//   - >= 2: ((localIdx << 6) | shardIdx) + 2
type StringID uint32

// EmptyStringID is the canonical ID of the empty string "".
const EmptyStringID StringID = 1

const (
	strpoolNumShards   = 64
	strpoolChunkShift  = 9
	strpoolChunkSize   = 1 << strpoolChunkShift // 512 strings per chunk
	strpoolChunkMask   = strpoolChunkSize - 1
	strpoolMaxChunks   = 1024 // up to 33.5M unique strings across 64 shards
	strpoolIdxBits     = 20
	strpoolIdxMask     = (1 << strpoolIdxBits) - 1 // 0x000FFFFF: holds localIdx+1 (1..524,288)
	strpoolFPMask      = ^uint32(strpoolIdxMask)   // 0xFFF00000: 12-bit upper hash fingerprint
	strpoolFnvOffset64 = 14695981039346656037
	strpoolFnvPrime64  = 1099511628211
)

// HashString returns the 64-bit FNV-1a hash of s.
func HashString(s string) uint64 {
	h := uint64(strpoolFnvOffset64)
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= strpoolFnvPrime64
	}
	return h
}

// HashBytes returns the 64-bit FNV-1a hash of b.
func HashBytes(b []byte) uint64 {
	h := uint64(strpoolFnvOffset64)
	for _, c := range b {
		h ^= uint64(c)
		h *= strpoolFnvPrime64
	}
	return h
}

// poolChunk packs 512 interned strings into a single noscan backing string
// once filled (whenever total byte length <= 65,535, which covers >99.99% of chunks),
// replacing 512 individual heap string objects and 512 GC-scanned string pointers
// (8 KiB) with 1 string pointer + 1,026 bytes of noscan uint16 cumulative offsets
// (~1.1 KiB total, ptrdata = 24).
type poolChunk struct {
	slab   string
	raw    *[strpoolChunkSize]string
	offs16 [strpoolChunkSize + 1]uint16
}

func (c *poolChunk) at(slot uint32) string {
	if c.raw == nil {
		return c.slab[c.offs16[slot]:c.offs16[slot+1]]
	}
	return c.raw[slot]
}

type strpoolShard struct {
	mu           sync.RWMutex
	table        []uint32 // open-addressed hash table, at most 3/4 full: (hash32 & 0xFFF00000) | (localIdx + 1)
	count        uint32
	nextLocalIdx uint32
}

var (
	poolShards     [strpoolNumShards]strpoolShard
	poolChunks     [strpoolMaxChunks * strpoolNumShards]atomic.Pointer[poolChunk]
	emptyPoolChunk = &poolChunk{}
	poolSeed       = maphash.MakeSeed()
	rewindEpoch    atomic.Uint64
)

func init() {
	for i := range poolChunks {
		poolChunks[i].Store(emptyPoolChunk)
	}
	// Reserve localIdx 0 and 1 in shard 0 as "" so StringID 0 (unset) and
	// EmptyStringID (1) resolve to "" via poolChunks[0].
	c0 := &poolChunk{raw: new([strpoolChunkSize]string)}
	poolChunks[0].Store(c0)
	poolShards[0].nextLocalIdx = 2
	poolShards[0].count = 2
}

// BaselineSnapshot records the per-shard watermark of interned strings at a
// post-init baseline so RewindBaseline can reclaim per-request strings.
type BaselineSnapshot [strpoolNumShards]uint32

// Epoch returns a monotonically increasing generation bumped whenever
// RewindBaseline rewinds any shard, allowing pooled L1 ID caches to invalidate
// stale IDs across worker requests.
func Epoch() uint64 {
	return rewindEpoch.Load()
}

// CaptureBaseline snapshots the current per-shard string counts.
func CaptureBaseline() BaselineSnapshot {
	var snap BaselineSnapshot
	for i := range poolShards {
		s := &poolShards[i]
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
	for i := range poolShards {
		if poolShards[i].rewind(uint32(i), snap[i]) {
			rewound = true
		}
	}
	if rewound {
		rewindEpoch.Add(1)
	}
}

func poolChunkIndex(shardIdx, chunkIdx uint32) uint32 {
	return (chunkIdx << 6) | shardIdx
}

func (s *strpoolShard) rewind(shardIdx, baseIdx uint32) bool {
	if shardIdx == 0 && baseIdx < 2 {
		baseIdx = 2
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	curIdx := s.nextLocalIdx
	if curIdx <= baseIdx {
		return false
	}

	if baseIdx == 0 || (shardIdx == 0 && baseIdx == 2) {
		if len(s.table) > 128 {
			s.table = nil
		} else {
			clear(s.table)
		}
	} else {
		wantCap := 128
		for int(baseIdx)*4 > wantCap*3 {
			wantCap *= 2
		}
		if len(s.table) > wantCap*2 {
			newTable := make([]uint32, wantCap)
			mask := uint32(wantCap - 1)
			for _, entry := range s.table {
				if entry == 0 {
					continue
				}
				localIdx := (entry & strpoolIdxMask) - 1
				if localIdx >= baseIdx {
					continue
				}
				h32 := uint32(poolHash(stringAtShard(shardIdx, localIdx)) >> 32)
				slot := h32 & mask
				for newTable[slot] != 0 {
					slot = (slot + 1) & mask
				}
				newTable[slot] = entry
			}
			s.table = newTable
		} else {
			var stackBuf [256]uint32
			survivors := stackBuf[:0]
			if int(baseIdx) > len(stackBuf) {
				survivors = make([]uint32, 0, baseIdx)
			}
			for _, entry := range s.table {
				if entry != 0 && (entry&strpoolIdxMask)-1 < baseIdx {
					survivors = append(survivors, entry)
				}
			}
			clear(s.table)
			mask := uint32(len(s.table) - 1)
			for _, entry := range survivors {
				localIdx := (entry & strpoolIdxMask) - 1
				h32 := uint32(poolHash(stringAtShard(shardIdx, localIdx)) >> 32)
				slot := h32 & mask
				for s.table[slot] != 0 {
					slot = (slot + 1) & mask
				}
				s.table[slot] = entry
			}
		}
	}

	firstDropChunk := (baseIdx + strpoolChunkMask) >> strpoolChunkShift
	lastUsedChunk := (curIdx - 1) >> strpoolChunkShift
	for c := firstDropChunk; c <= lastUsedChunk; c++ {
		poolChunks[poolChunkIndex(shardIdx, c)].Store(emptyPoolChunk)
	}
	if keep := baseIdx & strpoolChunkMask; keep != 0 {
		tailChunkIdx := baseIdx >> strpoolChunkShift
		flatIdx := poolChunkIndex(shardIdx, tailChunkIdx)
		if chunk := poolChunks[flatIdx].Load(); chunk != emptyPoolChunk {
			if chunk.raw != nil {
				endOff := min(uint32(strpoolChunkSize), curIdx-(tailChunkIdx<<strpoolChunkShift))
				clear(chunk.raw[keep:endOff])
			} else {
				raw := new([strpoolChunkSize]string)
				for i := uint32(0); i < keep; i++ {
					raw[i] = strings.Clone(chunk.slab[chunk.offs16[i]:chunk.offs16[i+1]])
				}
				poolChunks[flatIdx].Store(&poolChunk{raw: raw})
			}
		}
	}

	s.nextLocalIdx = baseIdx
	s.count = baseIdx
	return true
}

func poolHash(s string) uint64 {
	return maphash.String(poolSeed, s)
}

func stringAtShard(shardIdx, localIdx uint32) string {
	return poolChunks[((localIdx>>strpoolChunkShift)<<6)|shardIdx].Load().at(localIdx & strpoolChunkMask)
}

func sealPoolChunk(chunk *poolChunk, lastStr string, owned bool) *poolChunk {
	raw := chunk.raw
	total := len(lastStr)
	for i := 0; i < strpoolChunkSize-1; i++ {
		total += len(raw[i])
	}
	if total > 0xFFFF {
		if !owned {
			lastStr = strings.Clone(lastStr)
		}
		raw[strpoolChunkSize-1] = lastStr
		return chunk
	}
	var b strings.Builder
	b.Grow(total)
	sealed := new(poolChunk)
	pos := uint16(0)
	for i := 0; i < strpoolChunkSize-1; i++ {
		si := raw[i]
		b.WriteString(si)
		pos += uint16(len(si))
		sealed.offs16[i+1] = pos
	}
	b.WriteString(lastStr)
	pos += uint16(len(lastStr))
	sealed.offs16[strpoolChunkSize] = pos
	sealed.slab = b.String()
	return sealed
}

func (s *strpoolShard) allocString(shardIdx uint32, str string, owned bool) uint32 {
	localIdx := s.nextLocalIdx
	chunkIdx := localIdx >> strpoolChunkShift
	if chunkIdx >= strpoolMaxChunks {
		panic("strpool: shard capacity exceeded")
	}
	s.nextLocalIdx = localIdx + 1
	flatIdx := poolChunkIndex(shardIdx, chunkIdx)
	slot := localIdx & strpoolChunkMask
	chunk := poolChunks[flatIdx].Load()
	if slot == strpoolChunkMask {
		if sealed := sealPoolChunk(chunk, str, owned); sealed != chunk {
			poolChunks[flatIdx].Store(sealed)
		}
		return localIdx
	}
	if !owned {
		// Clone on miss so substring slices from parser inputs never pin SourceFile backing buffers.
		str = strings.Clone(str)
	}
	if chunk == emptyPoolChunk {
		raw := new([strpoolChunkSize]string)
		raw[slot] = str
		poolChunks[flatIdx].Store(&poolChunk{raw: raw})
		return localIdx
	}
	chunk.raw[slot] = str
	return localIdx
}

func (s *strpoolShard) growTable(shardIdx uint32) {
	newCap := max(len(s.table)*2, 128)
	newTable := make([]uint32, newCap)
	mask := uint32(newCap - 1)
	for _, entry := range s.table {
		if entry == 0 {
			continue
		}
		localIdx := (entry & strpoolIdxMask) - 1
		h32 := uint32(poolHash(stringAtShard(shardIdx, localIdx)) >> 32)
		slot := h32 & mask
		for newTable[slot] != 0 {
			slot = (slot + 1) & mask
		}
		newTable[slot] = entry
	}
	s.table = newTable
}

func poolProbeLocked[S ~string | ~[]byte](s *strpoolShard, shardIdx, hash32 uint32, key S) (uint32, bool) {
	tbl := s.table
	if len(tbl) == 0 {
		return 0, false
	}
	fp := hash32 & strpoolFPMask
	mask := uint32(len(tbl) - 1)
	slot := hash32 & mask
	for {
		entry := tbl[slot]
		if entry == 0 {
			return 0, false
		}
		if entry&strpoolFPMask == fp {
			localIdx := (entry & strpoolIdxMask) - 1
			if stringAtShard(shardIdx, localIdx) == string(key) {
				return localIdx, true
			}
		}
		slot = (slot + 1) & mask
	}
}

func makeStringID(localIdx, shardIdx uint32) StringID {
	return StringID((((localIdx >> strpoolChunkShift) << (6 + strpoolChunkShift)) | (shardIdx << strpoolChunkShift) | (localIdx & strpoolChunkMask)))
}

func (s *strpoolShard) internSlow(shardIdx, hash32 uint32, str string, owned bool) StringID {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.table) == 0 {
		s.table = make([]uint32, 128)
	} else if localIdx, ok := poolProbeLocked(s, shardIdx, hash32, str); ok {
		return makeStringID(localIdx, shardIdx)
	}
	if int(s.count+1)*4 > len(s.table)*3 {
		s.growTable(shardIdx)
	}
	localIdx := s.allocString(shardIdx, str, owned)
	mask := uint32(len(s.table) - 1)
	slot := hash32 & mask
	for s.table[slot] != 0 {
		slot = (slot + 1) & mask
	}
	s.table[slot] = (hash32 & strpoolFPMask) | (localIdx + 1)
	s.count++
	return makeStringID(localIdx, shardIdx)
}

// Intern interns str in the process-wide pool and returns its 32-bit handle.
func Intern(str string) StringID {
	if len(str) == 0 {
		return EmptyStringID
	}
	h := poolHash(str)
	shardIdx := uint32(h) & (strpoolNumShards - 1)
	hash32 := uint32(h >> 32)
	s := &poolShards[shardIdx]

	s.mu.RLock()
	if localIdx, ok := poolProbeLocked(s, shardIdx, hash32, str); ok {
		s.mu.RUnlock()
		return makeStringID(localIdx, shardIdx)
	}
	s.mu.RUnlock()
	return s.internSlow(shardIdx, hash32, str, false)
}

// String returns the string value for id, or "" if id is 0 or 1.
func String(id StringID) string {
	c := poolChunks[uint16(id>>strpoolChunkShift)].Load()
	s := id & strpoolChunkMask
	if c.raw == nil {
		return c.slab[c.offs16[s]:c.offs16[s+1]]
	}
	return c.raw[s]
}

// Lookup returns the ID of str if already interned, or 0 if not present.
func Lookup(str string) StringID {
	if len(str) == 0 {
		return EmptyStringID
	}
	h := poolHash(str)
	shardIdx := uint32(h) & (strpoolNumShards - 1)
	hash32 := uint32(h >> 32)
	s := &poolShards[shardIdx]

	s.mu.RLock()
	localIdx, ok := poolProbeLocked(s, shardIdx, hash32, str)
	s.mu.RUnlock()
	if !ok {
		return 0
	}
	return makeStringID(localIdx, shardIdx)
}

// InternBytes interns raw bytes without heap allocating a string if already interned.
func InternBytes(b []byte) StringID {
	if len(b) == 0 {
		return EmptyStringID
	}
	h := maphash.Bytes(poolSeed, b)
	shardIdx := uint32(h) & (strpoolNumShards - 1)
	hash32 := uint32(h >> 32)
	s := &poolShards[shardIdx]

	s.mu.RLock()
	if localIdx, ok := poolProbeLocked(s, shardIdx, hash32, b); ok {
		s.mu.RUnlock()
		return makeStringID(localIdx, shardIdx)
	}
	s.mu.RUnlock()
	return s.internSlow(shardIdx, hash32, string(b), true)
}

const (
	joinDot2CacheBits  = 14
	joinDot2CacheLen   = 1 << joinDot2CacheBits // 16384 buckets = 32,768 2-way entries (128 KiB BSS)
	joinDot3CacheBits  = 13
	joinDot3CacheLen   = 1 << joinDot3CacheBits // 8192 buckets = 16,384 2-way entries (64 KiB BSS)
	joinDot4CacheBits  = 13
	joinDot4CacheLen   = 1 << joinDot4CacheBits // 8192 buckets = 16,384 2-way entries (64 KiB BSS)
	joinDotMaxTotalLen = 512
	joinDotStackBufLen = 512
	joinDotFastBufLen  = 128
)

// joinDot2Bucket is an immutable 48-byte 2-way set-associative cache bucket
// with only 2 Go pointers (ptrdata = 24). Neither caller slice a nor b is stored,
// preventing any pinning of large parser input buffers while remaining 100% lock-free.
type joinDot2Bucket struct {
	out0  string
	out1  string
	h0    uint32
	h1    uint32
	alen0 uint16
	alen1 uint16
}

// joinDot3Bucket is an immutable 48-byte 2-way set-associative cache bucket
// with only 2 Go pointers (ptrdata = 24).
type joinDot3Bucket struct {
	out0  string
	out1  string
	h0    uint32
	h1    uint32
	alen0 uint16
	blen0 uint16
	alen1 uint16
	blen1 uint16
}

// joinDot4Bucket is an immutable 64-byte 2-way set-associative cache bucket
// with only 2 Go pointers (ptrdata = 24).
type joinDot4Bucket struct {
	out0  string
	out1  string
	h0    uint32
	h1    uint32
	alen0 uint16
	blen0 uint16
	clen0 uint16
	alen1 uint16
	blen1 uint16
	clen1 uint16
}

var (
	joinDot2Cache [joinDot2CacheLen]atomic.Pointer[joinDot2Bucket]
	joinDot3Cache [joinDot3CacheLen]atomic.Pointer[joinDot3Bucket]
	joinDot4Cache [joinDot4CacheLen]atomic.Pointer[joinDot4Bucket]
)

func formatJoin2(a, b string) string {
	total := len(a) + len(b) + 1
	if total <= joinDotFastBufLen {
		var buf [joinDotFastBufLen]byte
		i := copy(buf[:], a)
		buf[i] = '.'
		copy(buf[i+1:], b)
		return string(buf[:total])
	}
	if total <= joinDotStackBufLen {
		var buf [joinDotStackBufLen]byte
		i := copy(buf[:], a)
		buf[i] = '.'
		copy(buf[i+1:], b)
		return string(buf[:total])
	}
	return a + "." + b
}

func formatJoin3(a, b, c string) string {
	total := len(a) + len(b) + len(c) + 2
	if total <= joinDotFastBufLen {
		var buf [joinDotFastBufLen]byte
		i := copy(buf[:], a)
		buf[i] = '.'
		i += 1 + copy(buf[i+1:], b)
		buf[i] = '.'
		copy(buf[i+1:], c)
		return string(buf[:total])
	}
	if total <= joinDotStackBufLen {
		var buf [joinDotStackBufLen]byte
		i := copy(buf[:], a)
		buf[i] = '.'
		i += 1 + copy(buf[i+1:], b)
		buf[i] = '.'
		copy(buf[i+1:], c)
		return string(buf[:total])
	}
	return a + "." + b + "." + c
}

func formatJoin4(a, b, c, d string) string {
	total := len(a) + len(b) + len(c) + len(d) + 3
	if total <= joinDotFastBufLen {
		var buf [joinDotFastBufLen]byte
		i := copy(buf[:], a)
		buf[i] = '.'
		i += 1 + copy(buf[i+1:], b)
		buf[i] = '.'
		i += 1 + copy(buf[i+1:], c)
		buf[i] = '.'
		copy(buf[i+1:], d)
		return string(buf[:total])
	}
	if total <= joinDotStackBufLen {
		var buf [joinDotStackBufLen]byte
		i := copy(buf[:], a)
		buf[i] = '.'
		i += 1 + copy(buf[i+1:], b)
		buf[i] = '.'
		i += 1 + copy(buf[i+1:], c)
		buf[i] = '.'
		copy(buf[i+1:], d)
		return string(buf[:total])
	}
	return a + "." + b + "." + c + "." + d
}

// JoinDot returns a + "." + b, deduplicating repeated qualified-name joins via
// a bounded lock-free 2-way set-associative cache without polluting the
// permanent AST string pool.
func JoinDot(a, b string) string {
	total := len(a) + len(b) + 1
	if total > joinDotMaxTotalLen {
		return a + "." + b
	}
	h := maphash.String(poolSeed, a)*0x9E3779B97F4A7C15 ^ maphash.String(poolSeed, b)
	idx := uint32(h >> (64 - joinDot2CacheBits))
	h32 := uint32(h)
	bkt := joinDot2Cache[idx].Load()
	if bkt != nil {
		if bkt.h0 == h32 && len(bkt.out0) == total && int(bkt.alen0) == len(a) {
			if s := bkt.out0; s[:len(a)] == a && s[len(a)+1:] == b {
				return s
			}
		}
		if bkt.h1 == h32 && len(bkt.out1) == total && int(bkt.alen1) == len(a) {
			if s := bkt.out1; s[:len(a)] == a && s[len(a)+1:] == b {
				return s
			}
		}
	}
	return joinDot2Miss(a, b, idx, h32, bkt)
}

//go:noinline
func joinDot2Miss(a, b string, idx, h32 uint32, old *joinDot2Bucket) string {
	out := formatJoin2(a, b)
	nb := &joinDot2Bucket{
		out0:  out,
		h0:    h32,
		alen0: uint16(len(a)),
	}
	if old != nil {
		nb.out1 = old.out0
		nb.h1 = old.h0
		nb.alen1 = old.alen0
	}
	joinDot2Cache[idx].Store(nb)
	return out
}

// JoinDot3 returns a + "." + b + "." + c, deduplicating repeated 3-part
// qualified-name joins via a bounded lock-free 2-way set-associative cache
// without polluting the permanent AST string pool.
func JoinDot3(a, b, c string) string {
	total := len(a) + len(b) + len(c) + 2
	if total > joinDotMaxTotalLen {
		return a + "." + b + "." + c
	}
	h := maphash.String(poolSeed, a)*0x9E3779B97F4A7C15 ^ maphash.String(poolSeed, b)*0xBF58476D1CE4E5B9 ^ maphash.String(poolSeed, c)
	idx := uint32(h >> (64 - joinDot3CacheBits))
	h32 := uint32(h)
	bkt := joinDot3Cache[idx].Load()
	if bkt != nil {
		if bkt.h0 == h32 && len(bkt.out0) == total && int(bkt.alen0) == len(a) && int(bkt.blen0) == len(b) {
			if s := bkt.out0; s[:len(a)] == a && s[len(a)+1:len(a)+1+len(b)] == b && s[len(a)+len(b)+2:] == c {
				return s
			}
		}
		if bkt.h1 == h32 && len(bkt.out1) == total && int(bkt.alen1) == len(a) && int(bkt.blen1) == len(b) {
			if s := bkt.out1; s[:len(a)] == a && s[len(a)+1:len(a)+1+len(b)] == b && s[len(a)+len(b)+2:] == c {
				return s
			}
		}
	}
	return joinDot3Miss(a, b, c, idx, h32, bkt)
}

//go:noinline
func joinDot3Miss(a, b, c string, idx, h32 uint32, old *joinDot3Bucket) string {
	out := formatJoin3(a, b, c)
	nb := &joinDot3Bucket{
		out0:  out,
		h0:    h32,
		alen0: uint16(len(a)),
		blen0: uint16(len(b)),
	}
	if old != nil {
		nb.out1 = old.out0
		nb.h1 = old.h0
		nb.alen1 = old.alen0
		nb.blen1 = old.blen0
	}
	joinDot3Cache[idx].Store(nb)
	return out
}

// JoinDot4 returns a + "." + b + "." + c + "." + d, deduplicating repeated 4-part
// qualified-name joins via a bounded lock-free 2-way set-associative cache
// without polluting the permanent AST string pool.
func JoinDot4(a, b, c, d string) string {
	total := len(a) + len(b) + len(c) + len(d) + 3
	if total > joinDotMaxTotalLen {
		return a + "." + b + "." + c + "." + d
	}
	h := maphash.String(poolSeed, a)*0x9E3779B97F4A7C15 ^ maphash.String(poolSeed, b)*0xBF58476D1CE4E5B9 ^ maphash.String(poolSeed, c)*0x94D049BB133111EB ^ maphash.String(poolSeed, d)
	idx := uint32(h >> (64 - joinDot4CacheBits))
	h32 := uint32(h)
	bkt := joinDot4Cache[idx].Load()
	if bkt != nil {
		if bkt.h0 == h32 && len(bkt.out0) == total && int(bkt.alen0) == len(a) && int(bkt.blen0) == len(b) && int(bkt.clen0) == len(c) {
			abEnd := len(a) + 1 + len(b)
			abcEnd := abEnd + 1 + len(c)
			if s := bkt.out0; s[:len(a)] == a && s[len(a)+1:abEnd] == b && s[abEnd+1:abcEnd] == c && s[abcEnd+1:] == d {
				return s
			}
		}
		if bkt.h1 == h32 && len(bkt.out1) == total && int(bkt.alen1) == len(a) && int(bkt.blen1) == len(b) && int(bkt.clen1) == len(c) {
			abEnd := len(a) + 1 + len(b)
			abcEnd := abEnd + 1 + len(c)
			if s := bkt.out1; s[:len(a)] == a && s[len(a)+1:abEnd] == b && s[abEnd+1:abcEnd] == c && s[abcEnd+1:] == d {
				return s
			}
		}
	}
	return joinDot4Miss(a, b, c, d, idx, h32, bkt)
}

//go:noinline
func joinDot4Miss(a, b, c, d string, idx, h32 uint32, old *joinDot4Bucket) string {
	out := formatJoin4(a, b, c, d)
	nb := &joinDot4Bucket{
		out0:  out,
		h0:    h32,
		alen0: uint16(len(a)),
		blen0: uint16(len(b)),
		clen0: uint16(len(c)),
	}
	if old != nil {
		nb.out1 = old.out0
		nb.h1 = old.h0
		nb.alen1 = old.alen0
		nb.blen1 = old.blen0
		nb.clen1 = old.clen0
	}
	joinDot4Cache[idx].Store(nb)
	return out
}
