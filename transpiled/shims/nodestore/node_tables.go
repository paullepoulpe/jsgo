// Port of the side tables of com/google/javascript/rhino/Node.java @ closure-compiler v20261104

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
	"math"
	"reflect"
	"sync"
	"sync/atomic"
)

// Side tables for sparse and cold node attributes, keeping NodeData exactly
// 32 bytes with zero Go pointers for noscan allocation.

const (
	tableBlockShift = 12
	tableBlockLen   = 1 << tableBlockShift
	tableBlockMask  = tableBlockLen - 1

	tableSubShift  = 8
	tableSubLen    = 1 << tableSubShift
	tableSubMask   = tableSubLen - 1
	tableDirLen    = 64
	tableMaxBlocks = tableDirLen * tableSubLen // 2^26 handles per table
)

type anyBlock [tableBlockLen]any

// AnyTable is a lazily allocated table of contiguous [4096]any blocks
// addressed by 1-based 32-bit handles (up to 2^26 = 67,108,864 handles).
// Zero is never handed out. Each slot h is written once in Add before h is
// published atomically to readers, eliminating per-entry *any heap boxes
// while remaining 100% race-free and branch/call-free on Load.
type AnyTable struct {
	count  atomic.Uint32
	blk0   atomic.Pointer[anyBlock]
	blocks [tableMaxBlocks]atomic.Pointer[anyBlock]
}

// Load returns the entry of handle h, or nil if unmapped.
func (t *AnyTable) Load(h uint32) any {
	if h-1 < tableBlockLen-1 {
		if blk := t.blk0.Load(); blk != nil {
			return blk[h]
		}
		return nil
	}
	if h-1 < tableMaxBlocks*tableBlockLen-1 {
		if blk := t.blocks[h>>tableBlockShift].Load(); blk != nil {
			return blk[h&tableBlockMask]
		}
	}
	return nil
}

// Peek returns the entry of h or nil if h was not allocated yet.
func (t *AnyTable) Peek(h uint32) any {
	if h == 0 || h > t.count.Load() {
		return nil
	}
	return t.Load(h)
}

// Add stores v in a new 1-based entry and returns its 32-bit handle.
func (t *AnyTable) Add(mu *sync.Mutex, v any) uint32 {
	h := t.count.Add(1)
	if h < tableBlockLen {
		blk := t.blk0.Load()
		if blk == nil {
			blk = t.addBlock(mu, 0)
		}
		blk[h] = v
		return h
	}
	idx := h >> tableBlockShift
	if idx >= tableMaxBlocks {
		panic("nodestore: AnyTable capacity exceeded")
	}
	blk := t.blocks[idx].Load()
	if blk == nil {
		blk = t.addBlock(mu, idx)
	}
	blk[h&tableBlockMask] = v
	return h
}

//go:noinline
func (t *AnyTable) addBlock(mu *sync.Mutex, idx uint32) *anyBlock {
	mu.Lock()
	defer mu.Unlock()
	blk := t.blocks[idx].Load()
	if blk == nil {
		blk = new(anyBlock)
		t.blocks[idx].Store(blk)
		if idx == 0 {
			t.blk0.Store(blk)
		}
	}
	return blk
}

// set overwrites an already-allocated 1-based handle h with v.
// Caller must hold exclusive synchronization for handle h.
func (t *AnyTable) set(h uint32, v any) {
	if h < tableBlockLen {
		t.blk0.Load()[h] = v
		return
	}
	t.blocks[h>>tableBlockShift].Load()[h&tableBlockMask] = v
}

// clear zeroes an already-allocated 1-based handle h so GC can reclaim its value.
// Caller must hold exclusive synchronization for handle h.
func (t *AnyTable) clear(h uint32) {
	if h < tableBlockLen {
		t.blk0.Load()[h] = nil
		return
	}
	t.blocks[h>>tableBlockShift].Load()[h&tableBlockMask] = nil
}

// --- Noscan Number Table ---

type numBlock [tableBlockLen]atomic.Uint64
type numSubDir [tableSubLen]atomic.Pointer[numBlock]

// numTable is a lazily allocated 2-level table of noscan 64-bit IEEE-754 float64 payloads
// indexed by 1-based uint32 handles stored in NodeData.Payload (0 represents +0.0).
type numTable struct {
	count atomic.Uint32
	blk0  atomic.Pointer[numBlock]
	dir   [tableDirLen]atomic.Pointer[numSubDir]
}

func (t *numTable) load(h uint32) uint64 {
	if h < tableBlockLen {
		blk := t.blk0.Load()
		if blk == nil {
			return 0
		}
		return blk[h].Load()
	}
	k := h >> tableBlockShift
	topIdx := k >> tableSubShift
	if topIdx >= tableDirLen {
		return 0
	}
	sub := t.dir[topIdx].Load()
	if sub == nil {
		return 0
	}
	blk := sub[k&tableSubMask].Load()
	if blk == nil {
		return 0
	}
	return blk[h&tableBlockMask].Load()
}

func (t *numTable) add(mu *sync.Mutex, bits uint64) uint32 {
	h := t.count.Add(1)
	if h < tableBlockLen {
		blk := t.blk0.Load()
		if blk == nil {
			blk = t.addBlock(mu, 0, 0)
		}
		blk[h].Store(bits)
		return h
	}
	k := h >> tableBlockShift
	topIdx := k >> tableSubShift
	if topIdx >= tableDirLen {
		panic("nodestore: numTable capacity exceeded")
	}
	subIdx := k & tableSubMask
	var blk *numBlock
	if sub := t.dir[topIdx].Load(); sub != nil {
		blk = sub[subIdx].Load()
	}
	if blk == nil {
		blk = t.addBlock(mu, topIdx, subIdx)
	}
	blk[h&tableBlockMask].Store(bits)
	return h
}

func (t *numTable) addBlock(mu *sync.Mutex, topIdx, subIdx uint32) *numBlock {
	mu.Lock()
	defer mu.Unlock()
	sub := t.dir[topIdx].Load()
	if sub == nil {
		sub = new(numSubDir)
		blk := new(numBlock)
		sub[subIdx].Store(blk)
		t.dir[topIdx].Store(sub)
		if topIdx == 0 && subIdx == 0 {
			t.blk0.Store(blk)
		}
		return blk
	}
	blk := sub[subIdx].Load()
	if blk == nil {
		blk = new(numBlock)
		sub[subIdx].Store(blk)
		if topIdx == 0 && subIdx == 0 {
			t.blk0.Store(blk)
		}
	}
	return blk
}

// --- Columns & Handle Holders ---

// ColorMarker is implemented by generated Color structs so type/color storage can
// record kindFlagTypeIsColor directly on NodeData.KindFlags without interface probes on read.
type ColorMarker interface {
	Implements_colors_Color()
}

// TCHandleHolder is implemented by generated JSType and Color structs so each
// distinct type/color gets a single 32-bit handle cached atomically on the struct.
type TCHandleHolder interface {
	IsJavaNil_() bool
	TCHandlePtr_() *atomic.Uint32
}

// SourceFileHandleHolder is implemented by generated StaticSourceFile structs so each
// distinct source file gets a single 32-bit handle cached atomically on the struct.
type SourceFileHandleHolder interface {
	IsJavaNil_() bool
	SourceFileHandlePtr_() *atomic.Uint32
}

// JSDocHandleHolder is implemented by generated JSDocInfo structs so each
// distinct JSDocInfo gets a single 32-bit handle cached atomically on the struct.
type JSDocHandleHolder interface {
	IsJavaNil_() bool
	JSDocHandlePtr_() *atomic.Uint32
}

type javaNil interface {
	IsJavaNil_() bool
}

var (
	tcHandles           AnyTable
	tcHandlesMu         sync.Mutex
	sourceFileHandles   AnyTable
	sourceFileHandlesMu sync.Mutex
	sourceFileFallback  sync.Map
)

// NodeColumn stores 16,384 uint32 values corresponding to one nodeBlock.
// Block publication uses atomic.Pointer[NodeColumn]; disjoint Node IDs address
// independent uint32 array elements (matching nodeBlock's [16384]NodeData model),
// avoiding redundant per-element acquire barriers on ARM64.
type NodeColumn [blockLen]uint32

type columnField [numBlocks]atomic.Pointer[NodeColumn]

var (
	tcField       columnField
	origNameField columnField
)

func (f *columnField) Get(n Node) uint32 {
	if c := f[n>>blockShift].Load(); c != nil {
		return c[n&blockMask]
	}
	return 0
}

func (f *columnField) Set(n Node, v uint32) {
	if c := f[n>>blockShift].Load(); c != nil {
		if v != 0 || c[n&blockMask] != 0 {
			c[n&blockMask] = v
		}
		return
	}
	if v != 0 {
		f.setSlow(n, v)
	}
}

//go:noinline
func (f *columnField) setSlow(n Node, v uint32) {
	if n == 0 {
		_ = nodeBlocks[0].Load()[0]
	}
	c := f[n>>blockShift].Load()
	if c == nil {
		s := storeForNode(n)
		switch f {
		case &tcField:
			c = s.addColumn(f, &s.tcs, uint32(n))
		case &origNameField:
			c = s.addColumn(f, &s.origNames, uint32(n))
		default:
			panic("nodestore: unknown columnField")
		}
	}
	c[n&blockMask] = v
}

func (s *NodeStore) addColumn(f *columnField, cols *[]*NodeColumn, id uint32) *NodeColumn {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := int((id - s.base) >> blockShift)
	for len(*cols) <= k {
		*cols = append(*cols, nil)
	}
	if (*cols)[k] == nil {
		col := new(NodeColumn)
		(*cols)[k] = col
		f[id>>blockShift].Store(col)
	}
	return (*cols)[k]
}

func (f *columnField) register(first uint32, cols []*NodeColumn) {
	for k, c := range cols {
		if c != nil {
			f[first+uint32(k)].Store(c)
		}
	}
}

func (f *columnField) unregister(first uint32, cols []*NodeColumn) {
	for k, c := range cols {
		if c != nil {
			f[first+uint32(k)].CompareAndSwap(c, nil)
		}
	}
}

// --- Lazy 2-Level uint16 Length Side Table (8 KiB noscan per 4096-node page) ---

const (
	len16PageShift = 12
	len16PageLen   = 1 << len16PageShift // 4096 uint16 entries (8 KiB)
	len16PageMask  = len16PageLen - 1
	len16SubShift  = 8
	len16SubLen    = 1 << len16SubShift // 256 sub slots (2 KiB)
	len16SubMask   = len16SubLen - 1
	len16TopShift  = len16PageShift + len16SubShift // 20 -> 4096 top slots (32 KiB BSS)
	len16TopLen    = 1 << (32 - len16TopShift)
)

type len16Page [len16PageLen]uint16
type len16Sub [len16SubLen]atomic.Pointer[len16Page]

type len16Record struct {
	idx  uint32
	page *len16Page
}

var (
	len16Mu               sync.Mutex
	emptyLen16Page        len16Page
	unallocWriteLen16Page len16Page
	emptyLen16Sub         len16Sub
	emptyLen16WriteSub    len16Sub
	len16Top              [len16TopLen]atomic.Pointer[len16Sub]
	len16WriteTop         [len16TopLen]atomic.Pointer[len16Sub]
)

func initLen16Top() {
	for i := range unallocWriteLen16Page {
		unallocWriteLen16Page[i] = 0xFFFF
	}
	for i := range emptyLen16Sub {
		emptyLen16Sub[i].Store(&emptyLen16Page)
		emptyLen16WriteSub[i].Store(&unallocWriteLen16Page)
	}
	for i := range len16Top {
		len16Top[i].Store(&emptyLen16Sub)
		len16WriteTop[i].Store(&emptyLen16WriteSub)
	}
}

func ensureLen16Subs(topIdx uint32) (*len16Sub, *len16Sub) {
	rSub := len16Top[topIdx].Load()
	wSub := len16WriteTop[topIdx].Load()
	if rSub == &emptyLen16Sub || wSub == &emptyLen16WriteSub {
		len16Mu.Lock()
		rSub = len16Top[topIdx].Load()
		wSub = len16WriteTop[topIdx].Load()
		if rSub == &emptyLen16Sub {
			freshR := new(len16Sub)
			freshW := new(len16Sub)
			for i := range freshR {
				freshR[i].Store(&emptyLen16Page)
				freshW[i].Store(&unallocWriteLen16Page)
			}
			len16Top[topIdx].Store(freshR)
			len16WriteTop[topIdx].Store(freshW)
			rSub = freshR
			wSub = freshW
		}
		len16Mu.Unlock()
	}
	return rSub, wSub
}

func loadLen16Page(idx uint32) *len16Page {
	return len16Top[idx>>len16SubShift].Load()[idx&len16SubMask].Load()
}

func storeLen16Page(idx uint32, p *len16Page) {
	rSub, wSub := ensureLen16Subs(idx >> len16SubShift)
	subIdx := idx & len16SubMask
	rSub[subIdx].Store(p)
	wSub[subIdx].Store(p)
}

func casLen16Page(idx uint32, old *len16Page) {
	topIdx := idx >> len16SubShift
	subIdx := idx & len16SubMask
	if wSub := len16WriteTop[topIdx].Load(); wSub != &emptyLen16WriteSub {
		wSub[subIdx].CompareAndSwap(old, &unallocWriteLen16Page)
	}
	if rSub := len16Top[topIdx].Load(); rSub != &emptyLen16Sub {
		rSub[subIdx].CompareAndSwap(old, &emptyLen16Page)
	}
}

func ensureLen16Page(n Node, idx uint32) *len16Page {
	if p := loadLen16Page(idx); p != &emptyLen16Page {
		return p
	}
	s := storeForNode(n)
	s.mu.Lock()
	p := loadLen16Page(idx)
	if p == &emptyLen16Page {
		for _, r := range s.len16Pages {
			if r.idx == idx {
				p = r.page
				break
			}
		}
		if p == &emptyLen16Page || p == nil {
			p = new(len16Page)
			s.len16Pages = append(s.len16Pages, len16Record{idx: idx, page: p})
		}
		storeLen16Page(idx, p)
	}
	s.mu.Unlock()
	return p
}

// --- Sparse Length Overflow Pages (1 KiB noscan per 256-node page with <0 or >=65535 length) ---

const (
	lenOverflowShift    = 8
	lenOverflowLen      = 1 << lenOverflowShift
	lenOverflowMask     = lenOverflowLen - 1
	lenOverflowSubShift = 11
	lenOverflowSubLen   = 1 << lenOverflowSubShift
	lenOverflowSubMask  = lenOverflowSubLen - 1
	lenOverflowTopShift = lenOverflowShift + lenOverflowSubShift // 19 -> 8192 top slots (64 KiB BSS)
	lenOverflowTopLen   = 1 << (32 - lenOverflowTopShift)
)

type lenOverflowPage [lenOverflowLen]uint32
type lenOverflowSub [lenOverflowSubLen]atomic.Pointer[lenOverflowPage]

type lenOverflowRecord struct {
	idx  uint32
	page *lenOverflowPage
}

var (
	lenOverflowMu  sync.Mutex
	lenOverflowTop [lenOverflowTopLen]atomic.Pointer[lenOverflowSub]
	negLenCount    atomic.Uint32
)

func ensureLenOverflowSub(topIdx uint32) *lenOverflowSub {
	sub := lenOverflowTop[topIdx].Load()
	if sub == nil {
		lenOverflowMu.Lock()
		sub = lenOverflowTop[topIdx].Load()
		if sub == nil {
			sub = new(lenOverflowSub)
			lenOverflowTop[topIdx].Store(sub)
		}
		lenOverflowMu.Unlock()
	}
	return sub
}

func loadLenOverflowPage(idx uint32) *lenOverflowPage {
	if sub := lenOverflowTop[idx>>lenOverflowSubShift].Load(); sub != nil {
		return sub[idx&lenOverflowSubMask].Load()
	}
	return nil
}

func storeLenOverflowPage(idx uint32, p *lenOverflowPage) {
	ensureLenOverflowSub(idx >> lenOverflowSubShift)[idx&lenOverflowSubMask].Store(p)
}

func casLenOverflowPage(idx uint32, old *lenOverflowPage) {
	if sub := lenOverflowTop[idx>>lenOverflowSubShift].Load(); sub != nil {
		sub[idx&lenOverflowSubMask].CompareAndSwap(old, nil)
	}
}

func setLenOverflow(n Node, v uint32) {
	idx := uint32(n) >> lenOverflowShift
	if p := loadLenOverflowPage(idx); p != nil {
		if v != 0 || p[uint32(n)&lenOverflowMask] != 0 {
			p[uint32(n)&lenOverflowMask] = v
		}
		return
	}
	if v == 0 {
		return
	}
	setLenOverflowSlow(n, idx, v)
}

//go:noinline
func setLenOverflowSlow(n Node, idx, v uint32) {
	s := storeForNode(n)
	s.mu.Lock()
	p := loadLenOverflowPage(idx)
	if p == nil {
		for _, r := range s.lengthOverflows {
			if r.idx == idx {
				p = r.page
				break
			}
		}
		if p == nil {
			p = new(lenOverflowPage)
			s.lengthOverflows = append(s.lengthOverflows, lenOverflowRecord{idx: idx, page: p})
		}
		storeLenOverflowPage(idx, p)
	}
	s.mu.Unlock()
	p[uint32(n)&lenOverflowMask] = v
}

// --- Sparse SourceFile Overflow Pages (8 KiB noscan per 2048-node page with >=65535 handle) ---

const (
	sfOverflowShift    = 11
	sfOverflowLen      = 1 << sfOverflowShift
	sfOverflowMask     = sfOverflowLen - 1
	sfOverflowSubShift = 10
	sfOverflowSubLen   = 1 << sfOverflowSubShift
	sfOverflowSubMask  = sfOverflowSubLen - 1
	sfOverflowTopShift = sfOverflowShift + sfOverflowSubShift // 21 -> 2048 top slots (16 KiB BSS)
	sfOverflowTopLen   = 1 << (32 - sfOverflowTopShift)
)

type sfOverflowPage [sfOverflowLen]uint32
type sfOverflowSub [sfOverflowSubLen]atomic.Pointer[sfOverflowPage]

type sfOverflowRecord struct {
	idx  uint32
	page *sfOverflowPage
}

var (
	sfOverflowMu  sync.Mutex
	sfOverflowTop [sfOverflowTopLen]atomic.Pointer[sfOverflowSub]
)

func ensureSFOverflowSub(topIdx uint32) *sfOverflowSub {
	sub := sfOverflowTop[topIdx].Load()
	if sub == nil {
		sfOverflowMu.Lock()
		sub = sfOverflowTop[topIdx].Load()
		if sub == nil {
			sub = new(sfOverflowSub)
			sfOverflowTop[topIdx].Store(sub)
		}
		sfOverflowMu.Unlock()
	}
	return sub
}

func loadSFOverflowPage(idx uint32) *sfOverflowPage {
	if sub := sfOverflowTop[idx>>sfOverflowSubShift].Load(); sub != nil {
		return sub[idx&sfOverflowSubMask].Load()
	}
	return nil
}

func storeSFOverflowPage(idx uint32, p *sfOverflowPage) {
	ensureSFOverflowSub(idx >> sfOverflowSubShift)[idx&sfOverflowSubMask].Store(p)
}

func casSFOverflowPage(idx uint32, old *sfOverflowPage) {
	if sub := sfOverflowTop[idx>>sfOverflowSubShift].Load(); sub != nil {
		sub[idx&sfOverflowSubMask].CompareAndSwap(old, nil)
	}
}

func setSFOverflow(n Node, v uint32) {
	idx := uint32(n) >> sfOverflowShift
	if p := loadSFOverflowPage(idx); p != nil {
		if v != 0 || p[uint32(n)&sfOverflowMask] != 0 {
			p[uint32(n)&sfOverflowMask] = v
		}
		return
	}
	if v == 0 {
		return
	}
	setSFOverflowSlow(n, idx, v)
}

//go:noinline
func setSFOverflowSlow(n Node, idx, v uint32) {
	s := storeForNode(n)
	s.mu.Lock()
	p := loadSFOverflowPage(idx)
	if p == nil {
		for _, r := range s.sfOverflows {
			if r.idx == idx {
				p = r.page
				break
			}
		}
		if p == nil {
			p = new(sfOverflowPage)
			s.sfOverflows = append(s.sfOverflows, sfOverflowRecord{idx: idx, page: p})
		}
		storeSFOverflowPage(idx, p)
	}
	s.mu.Unlock()
	p[uint32(n)&sfOverflowMask] = v
}

// --- Per-Block Exact Cast Presence Bitmap (2 KiB noscan per block with casts) ---

type castBitBlock [blockLen >> 6]uint64

type castBitColumn [numBlocks]atomic.Pointer[castBitBlock]

var castBitField castBitColumn

func (f *castBitColumn) register(first uint32, cols []*castBitBlock) {
	for k, b := range cols {
		if b != nil {
			f[first+uint32(k)].Store(b)
		}
	}
}

func (f *castBitColumn) unregister(first uint32, cols []*castBitBlock) {
	for k, b := range cols {
		if b != nil {
			f[first+uint32(k)].CompareAndSwap(b, nil)
		}
	}
}

func (s *NodeStore) addCastBitBlock(id uint32) *castBitBlock {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := int((id - s.base) >> blockShift)
	for len(s.castBits) <= k {
		s.castBits = append(s.castBits, nil)
	}
	b := s.castBits[k]
	if b == nil {
		b = new(castBitBlock)
		s.castBits[k] = b
		castBitField[id>>blockShift].Store(b)
	}
	return b
}

func (s *NodeStore) setCastBit(id uint32) {
	b := castBitField[id>>blockShift].Load()
	if b == nil {
		b = s.addCastBitBlock(id)
	}
	w := &b[(id>>6)&((blockLen>>6)-1)]
	bit := uint64(1) << (id & 63)
	for {
		old := atomic.LoadUint64(w)
		if old&bit != 0 || atomic.CompareAndSwapUint64(w, old, old|bit) {
			return
		}
	}
}

func (s *NodeStore) clearCastBit(id uint32) {
	b := castBitField[id>>blockShift].Load()
	if b == nil {
		s.mu.Lock()
		k := int((id - s.base) >> blockShift)
		if k >= 0 && k < len(s.castBits) {
			b = s.castBits[k]
		}
		s.mu.Unlock()
		if b == nil {
			return
		}
	}
	w := &b[(id>>6)&((blockLen>>6)-1)]
	bit := uint64(1) << (id & 63)
	for {
		old := atomic.LoadUint64(w)
		if old&bit == 0 || atomic.CompareAndSwapUint64(w, old, old&^bit) {
			return
		}
	}
}

func storeForNode(n Node) *NodeStore {
	if uint32(n)&staticBit != 0 {
		return staticStore
	}
	return CurrentNodeStore()
}

// --- Cold PropList Map (64-Way Sharded, Gated by kindFlagHasColdProps) ---

const propShards = 64

type propShard struct {
	mu sync.RWMutex
	m  map[uint32]any
}

// --- Dedicated Noscan-Indexed JSDoc Table (64-Way Sharded) ---

const (
	jsdocShards  = 64
	jsdocInitCap = 16
)

// jsdocSlot is 8 bytes (id + handle) with zero Go pointers (100% noscan):
//   - id == 0: empty slot
//   - id != 0 && h != 0: live entry mapping node id -> 1-based s.jsdocs handle h
//   - id != 0 && h == 0: deleted tombstone
type jsdocSlot struct {
	id uint32
	h  uint32
}

type jsdocShard struct {
	mu    sync.RWMutex
	count uint32
	used  uint32
	slots []jsdocSlot
}

func jsdocShardIndex(id uint32) uint32 {
	return (id * 0x9e3779b9) >> 26
}

func jsdocSlotHash(id uint32) uint32 {
	h := id * 0x85ebca6b
	return h ^ (h >> 16)
}

func (sh *jsdocShard) rehash() {
	newCap := jsdocInitCap
	need := int(sh.count) + 1
	for need*10 >= newCap*7 {
		newCap <<= 1
	}
	if need*2 > newCap {
		newCap <<= 1
	}
	newSlots := make([]jsdocSlot, newCap)
	if sh.count != 0 {
		mask := uint32(newCap - 1)
		for _, sl := range sh.slots {
			if sl.h == 0 {
				continue
			}
			pos := jsdocSlotHash(sl.id) & mask
			for newSlots[pos].id != 0 {
				pos = (pos + 1) & mask
			}
			newSlots[pos] = sl
		}
	}
	sh.slots = newSlots
	sh.used = sh.count
}

func (s *NodeStore) jsdocLen() int {
	total := 0
	for i := range s.jsdocShards {
		sh := &s.jsdocShards[i]
		sh.mu.RLock()
		total += int(sh.count)
		sh.mu.RUnlock()
	}
	return total
}

func (s *NodeStore) getJSDoc(id uint32) any {
	if id == 0 {
		return nil
	}
	sh := &s.jsdocShards[jsdocShardIndex(id)]
	if !s.static && isSingleOwnerFastPath() {
		if sh.count == 0 {
			return nil
		}
		slots := sh.slots
		mask := uint32(len(slots) - 1)
		pos := jsdocSlotHash(id) & mask
		for {
			sl := slots[pos]
			if sl.id == id {
				if sl.h == 0 {
					return nil
				}
				return s.jsdocs.Load(sl.h)
			}
			if sl.id == 0 {
				return nil
			}
			pos = (pos + 1) & mask
		}
	}
	sh.mu.RLock()
	if sh.count == 0 {
		sh.mu.RUnlock()
		return nil
	}
	slots := sh.slots
	mask := uint32(len(slots) - 1)
	pos := jsdocSlotHash(id) & mask
	for {
		sl := slots[pos]
		if sl.id == id {
			var v any
			if sl.h != 0 {
				v = s.jsdocs.Load(sl.h)
			}
			sh.mu.RUnlock()
			return v
		}
		if sl.id == 0 {
			sh.mu.RUnlock()
			return nil
		}
		pos = (pos + 1) & mask
	}
}

func (s *NodeStore) hasJSDoc(id uint32) bool {
	if id == 0 {
		return false
	}
	sh := &s.jsdocShards[jsdocShardIndex(id)]
	if !s.static && isSingleOwnerFastPath() {
		if sh.count == 0 {
			return false
		}
		slots := sh.slots
		mask := uint32(len(slots) - 1)
		pos := jsdocSlotHash(id) & mask
		for {
			sl := slots[pos]
			if sl.id == id {
				return sl.h != 0
			}
			if sl.id == 0 {
				return false
			}
			pos = (pos + 1) & mask
		}
	}
	sh.mu.RLock()
	if sh.count == 0 {
		sh.mu.RUnlock()
		return false
	}
	slots := sh.slots
	mask := uint32(len(slots) - 1)
	pos := jsdocSlotHash(id) & mask
	for {
		sl := slots[pos]
		if sl.id == id {
			ok := sl.h != 0
			sh.mu.RUnlock()
			return ok
		}
		if sl.id == 0 {
			sh.mu.RUnlock()
			return false
		}
		pos = (pos + 1) & mask
	}
}

func (s *NodeStore) setJSDoc(id uint32, v any) {
	if id == 0 {
		return
	}
	sh := &s.jsdocShards[jsdocShardIndex(id)]
	if s.static || !isSingleOwnerFastPath() {
		sh.mu.Lock()
		defer sh.mu.Unlock()
	}
	if len(sh.slots) == 0 || int(sh.used+1)*10 >= len(sh.slots)*7 {
		sh.rehash()
	}
	slots := sh.slots
	mask := uint32(len(slots) - 1)
	pos := jsdocSlotHash(id) & mask
	tomb := -1
	for {
		sl := slots[pos]
		if sl.id == id {
			if sl.h != 0 {
				s.jsdocs.set(sl.h, v)
				return
			}
			slots[pos].h = s.jsdocs.Add(&s.jsdocValMu, v)
			sh.count++
			return
		}
		if sl.id == 0 {
			h := s.jsdocs.Add(&s.jsdocValMu, v)
			if tomb >= 0 {
				slots[tomb] = jsdocSlot{id: id, h: h}
			} else {
				slots[pos] = jsdocSlot{id: id, h: h}
				sh.used++
			}
			sh.count++
			return
		}
		if sl.h == 0 && tomb < 0 {
			tomb = int(pos)
		}
		pos = (pos + 1) & mask
	}
}

func (s *NodeStore) deleteJSDoc(id uint32) bool {
	if id == 0 {
		return false
	}
	sh := &s.jsdocShards[jsdocShardIndex(id)]
	if s.static || !isSingleOwnerFastPath() {
		sh.mu.Lock()
		defer sh.mu.Unlock()
	}
	if sh.count == 0 {
		return false
	}
	slots := sh.slots
	mask := uint32(len(slots) - 1)
	pos := jsdocSlotHash(id) & mask
	for {
		sl := slots[pos]
		if sl.id == id {
			if sl.h != 0 {
				s.jsdocs.clear(sl.h)
				slots[pos].h = 0
				sh.count--
				return true
			}
			return false
		}
		if sl.id == 0 {
			return false
		}
		pos = (pos + 1) & mask
	}
}

// --- Cold Aux Table (64-Way Sharded, Noscan Open-Addressed + Per-Shard Cast Bloom) ---

const (
	auxShards  = 64
	auxInitCap = 16
)

func auxShardIndex(id uint32) uint32 {
	return (id * 0x9e3779b9) >> 26
}

func auxSlotHash(id uint32) uint32 {
	h := id * 0x85ebca6b
	return h ^ (h >> 16)
}

func auxCastBloomMask(id uint32) uint64 {
	h := id * 0xcc9e2d51
	h ^= h >> 15
	return (uint64(1) << (h & 63)) | (uint64(1) << ((h >> 6) & 63))
}

// NodeAux holds exported cold / rarely accessed metadata for a Node.
type NodeAux struct {
	JSDoc any
	Extra any
}

// nodeAuxEntry holds ultra-rare (< 0.3% of nodes) cold attributes in 24 bytes:
// typeBeforeCast (16B) + changeTime (4B) + length (4B).
type nodeAuxEntry struct {
	typeBeforeCast any
	changeTime     int32
	length         int32
}

type nodeAux = nodeAuxEntry

// auxSlot is 16 bytes with zero Go pointers (100% noscan):
//   - id == 0: empty slot
//   - id != 0 && (castH != 0 || changeTime != 0 || length != 0): live entry
//   - id != 0 && castH == 0 && changeTime == 0 && length == 0: deleted tombstone
type auxSlot struct {
	id         uint32
	castH      uint32
	changeTime int32
	length     int32
}

type auxShard struct {
	castCount atomic.Uint32
	castBloom atomic.Uint64
	mu        sync.RWMutex
	count     uint32
	used      uint32
	slots     []auxSlot
	freeCastH []uint32
	extra     map[uint32]any
}

func (sh *auxShard) rehash() {
	newCap := auxInitCap
	need := int(sh.count) + 1
	for need*10 >= newCap*7 {
		newCap <<= 1
	}
	if need*2 > newCap {
		newCap <<= 1
	}
	newSlots := make([]auxSlot, newCap)
	if sh.count != 0 {
		mask := uint32(newCap - 1)
		for _, sl := range sh.slots {
			if sl.castH == 0 && (sl.changeTime|sl.length) == 0 {
				continue
			}
			pos := auxSlotHash(sl.id) & mask
			for newSlots[pos].id != 0 {
				pos = (pos + 1) & mask
			}
			newSlots[pos] = sl
		}
	}
	sh.slots = newSlots
	sh.used = sh.count
}

func (s *NodeStore) allocCastHandle(sh *auxShard, v any) uint32 {
	if n := len(sh.freeCastH); n > 0 {
		h := sh.freeCastH[n-1]
		sh.freeCastH = sh.freeCastH[:n-1]
		s.casts.set(h, v)
		return h
	}
	return s.casts.Add(&s.castValMu, v)
}

func (s *NodeStore) auxLen() int {
	total := 0
	for i := range s.aux {
		sh := &s.aux[i]
		sh.mu.RLock()
		total += int(sh.count) + len(sh.extra)
		sh.mu.RUnlock()
	}
	return total
}

func (s *NodeStore) getAuxTypeBeforeCast(id uint32) any {
	if !s.hasExtCast.Load() {
		return nil
	}
	sh := &s.aux[auxShardIndex(id)]
	if sh.castCount.Load() == 0 {
		return nil
	}
	mask := auxCastBloomMask(id)
	if sh.castBloom.Load()&mask != mask {
		return nil
	}
	if !s.static && isSingleOwnerFastPath() {
		slots := sh.slots
		if len(slots) == 0 {
			return nil
		}
		slotMask := uint32(len(slots) - 1)
		pos := auxSlotHash(id) & slotMask
		for {
			sl := slots[pos]
			if sl.id == id {
				if sl.castH == 0 {
					return nil
				}
				return s.casts.Load(sl.castH)
			}
			if sl.id == 0 {
				return nil
			}
			pos = (pos + 1) & slotMask
		}
	}
	sh.mu.RLock()
	if sh.castCount.Load() == 0 || len(sh.slots) == 0 {
		sh.mu.RUnlock()
		return nil
	}
	slots := sh.slots
	slotMask := uint32(len(slots) - 1)
	pos := auxSlotHash(id) & slotMask
	for {
		sl := slots[pos]
		if sl.id == id {
			var v any
			if sl.castH != 0 {
				v = s.casts.Load(sl.castH)
			}
			sh.mu.RUnlock()
			return v
		}
		if sl.id == 0 {
			sh.mu.RUnlock()
			return nil
		}
		pos = (pos + 1) & slotMask
	}
}

func (s *NodeStore) hasAuxTypeBeforeCast(id uint32) bool {
	if !s.hasExtCast.Load() {
		return false
	}
	return s.getAuxTypeBeforeCast(id) != nil
}

func (s *NodeStore) clearAuxTypeBeforeCastShard(sh *auxShard, id uint32) bool {
	if s.static || !isSingleOwnerFastPath() {
		sh.mu.Lock()
		defer sh.mu.Unlock()
	}
	if sh.castCount.Load() == 0 || len(sh.slots) == 0 {
		return false
	}
	slots := sh.slots
	slotMask := uint32(len(slots) - 1)
	pos := auxSlotHash(id) & slotMask
	for {
		sl := slots[pos]
		if sl.id == id {
			if sl.castH == 0 {
				return false
			}
			s.casts.clear(sl.castH)
			sh.freeCastH = append(sh.freeCastH, sl.castH)
			slots[pos].castH = 0
			if (sl.changeTime | sl.length) == 0 {
				sh.count--
			}
			rem := sh.castCount.Load() - 1
			sh.castCount.Store(rem)
			if rem == 0 {
				sh.castBloom.Store(0)
			}
			s.clearCastBit(id)
			return true
		}
		if sl.id == 0 {
			return false
		}
		pos = (pos + 1) & slotMask
	}
}

func (s *NodeStore) setAuxTypeBeforeCastNonNil(id uint32, v any) {
	sh := &s.aux[auxShardIndex(id)]
	if s.static || !isSingleOwnerFastPath() {
		sh.mu.Lock()
		defer sh.mu.Unlock()
	}
	if len(sh.slots) == 0 || int(sh.used+1)*10 >= len(sh.slots)*7 {
		sh.rehash()
	}
	slots := sh.slots
	slotMask := uint32(len(slots) - 1)
	pos := auxSlotHash(id) & slotMask
	tomb := -1
	for {
		sl := slots[pos]
		if sl.id == id {
			if sl.castH != 0 {
				s.casts.set(sl.castH, v)
				return
			}
			slots[pos].castH = s.allocCastHandle(sh, v)
			if (sl.changeTime | sl.length) == 0 {
				sh.count++
			}
			sh.castBloom.Store(sh.castBloom.Load() | auxCastBloomMask(id))
			sh.castCount.Store(sh.castCount.Load() + 1)
			if !s.hasExtCast.Load() {
				s.hasExtCast.Store(true)
			}
			s.setCastBit(id)
			return
		}
		if sl.id == 0 {
			h := s.allocCastHandle(sh, v)
			if tomb >= 0 {
				slots[tomb] = auxSlot{id: id, castH: h}
			} else {
				slots[pos] = auxSlot{id: id, castH: h}
				sh.used++
			}
			sh.count++
			sh.castBloom.Store(sh.castBloom.Load() | auxCastBloomMask(id))
			sh.castCount.Store(sh.castCount.Load() + 1)
			if !s.hasExtCast.Load() {
				s.hasExtCast.Store(true)
			}
			s.setCastBit(id)
			return
		}
		if sl.castH == 0 && (sl.changeTime|sl.length) == 0 && tomb < 0 {
			tomb = int(pos)
		}
		pos = (pos + 1) & slotMask
	}
}

func (s *NodeStore) setAuxTypeBeforeCast(id uint32, v any) {
	if v == nil {
		if !s.hasExtCast.Load() {
			return
		}
		sh := &s.aux[auxShardIndex(id)]
		if sh.castCount.Load() == 0 {
			return
		}
		mask := auxCastBloomMask(id)
		if sh.castBloom.Load()&mask != mask {
			return
		}
		_ = s.clearAuxTypeBeforeCastShard(sh, id)
		return
	}
	s.setAuxTypeBeforeCastNonNil(id, v)
}

func (s *NodeStore) getAuxChangeTime(id uint32) int32 {
	if !s.hasExtChangeTime.Load() {
		return 0
	}
	sh := &s.aux[auxShardIndex(id)]
	if !s.static && isSingleOwnerFastPath() {
		if sh.count == 0 {
			return 0
		}
		slots := sh.slots
		slotMask := uint32(len(slots) - 1)
		pos := auxSlotHash(id) & slotMask
		for {
			sl := slots[pos]
			if sl.id == id {
				return sl.changeTime
			}
			if sl.id == 0 {
				return 0
			}
			pos = (pos + 1) & slotMask
		}
	}
	sh.mu.RLock()
	if sh.count == 0 {
		sh.mu.RUnlock()
		return 0
	}
	slots := sh.slots
	slotMask := uint32(len(slots) - 1)
	pos := auxSlotHash(id) & slotMask
	for {
		sl := slots[pos]
		if sl.id == id {
			ct := sl.changeTime
			sh.mu.RUnlock()
			return ct
		}
		if sl.id == 0 {
			sh.mu.RUnlock()
			return 0
		}
		pos = (pos + 1) & slotMask
	}
}

func (s *NodeStore) setAuxChangeTime(id uint32, v int32) (hadExt bool) {
	if v == 0 && !s.hasExtChangeTime.Load() {
		return false
	}
	sh := &s.aux[auxShardIndex(id)]
	if s.static || !isSingleOwnerFastPath() {
		sh.mu.Lock()
		defer sh.mu.Unlock()
	}
	if v == 0 {
		if sh.count == 0 {
			return false
		}
		slots := sh.slots
		slotMask := uint32(len(slots) - 1)
		pos := auxSlotHash(id) & slotMask
		for {
			sl := slots[pos]
			if sl.id == id {
				if sl.changeTime == 0 {
					return false
				}
				slots[pos].changeTime = 0
				if sl.castH == 0 && sl.length == 0 {
					sh.count--
				}
				return true
			}
			if sl.id == 0 {
				return false
			}
			pos = (pos + 1) & slotMask
		}
	}
	if !s.hasExtChangeTime.Load() {
		s.hasExtChangeTime.Store(true)
	}
	if len(sh.slots) == 0 || int(sh.used+1)*10 >= len(sh.slots)*7 {
		sh.rehash()
	}
	slots := sh.slots
	slotMask := uint32(len(slots) - 1)
	pos := auxSlotHash(id) & slotMask
	tomb := -1
	for {
		sl := slots[pos]
		if sl.id == id {
			if sl.castH == 0 && (sl.changeTime|sl.length) == 0 {
				sh.count++
			}
			slots[pos].changeTime = v
			return true
		}
		if sl.id == 0 {
			if tomb >= 0 {
				slots[tomb] = auxSlot{id: id, changeTime: v}
			} else {
				slots[pos] = auxSlot{id: id, changeTime: v}
				sh.used++
			}
			sh.count++
			return true
		}
		if sl.castH == 0 && (sl.changeTime|sl.length) == 0 && tomb < 0 {
			tomb = int(pos)
		}
		pos = (pos + 1) & slotMask
	}
}

func (s *NodeStore) getAuxLengthOverflow(id uint32) int32 {
	sh := &s.aux[auxShardIndex(id)]
	if !s.static && isSingleOwnerFastPath() {
		if sh.count == 0 {
			return 0
		}
		slots := sh.slots
		slotMask := uint32(len(slots) - 1)
		pos := auxSlotHash(id) & slotMask
		for {
			sl := slots[pos]
			if sl.id == id {
				return sl.length
			}
			if sl.id == 0 {
				return 0
			}
			pos = (pos + 1) & slotMask
		}
	}
	sh.mu.RLock()
	if sh.count == 0 {
		sh.mu.RUnlock()
		return 0
	}
	slots := sh.slots
	slotMask := uint32(len(slots) - 1)
	pos := auxSlotHash(id) & slotMask
	for {
		sl := slots[pos]
		if sl.id == id {
			lo := sl.length
			sh.mu.RUnlock()
			return lo
		}
		if sl.id == 0 {
			sh.mu.RUnlock()
			return 0
		}
		pos = (pos + 1) & slotMask
	}
}

func (s *NodeStore) setAuxLengthOverflow(id uint32, v int32) {
	sh := &s.aux[auxShardIndex(id)]
	if s.static || !isSingleOwnerFastPath() {
		sh.mu.Lock()
		defer sh.mu.Unlock()
	}
	if v == 0 {
		if sh.count == 0 {
			return
		}
		slots := sh.slots
		slotMask := uint32(len(slots) - 1)
		pos := auxSlotHash(id) & slotMask
		for {
			sl := slots[pos]
			if sl.id == id {
				if sl.length == 0 {
					return
				}
				slots[pos].length = 0
				if sl.castH == 0 && sl.changeTime == 0 {
					sh.count--
				}
				return
			}
			if sl.id == 0 {
				return
			}
			pos = (pos + 1) & slotMask
		}
	}
	if len(sh.slots) == 0 || int(sh.used+1)*10 >= len(sh.slots)*7 {
		sh.rehash()
	}
	slots := sh.slots
	slotMask := uint32(len(slots) - 1)
	pos := auxSlotHash(id) & slotMask
	tomb := -1
	for {
		sl := slots[pos]
		if sl.id == id {
			if sl.castH == 0 && (sl.changeTime|sl.length) == 0 {
				sh.count++
			}
			slots[pos].length = v
			return
		}
		if sl.id == 0 {
			if tomb >= 0 {
				slots[tomb] = auxSlot{id: id, length: v}
			} else {
				slots[pos] = auxSlot{id: id, length: v}
				sh.used++
			}
			sh.count++
			return
		}
		if sl.castH == 0 && (sl.changeTime|sl.length) == 0 && tomb < 0 {
			tomb = int(pos)
		}
		pos = (pos + 1) & slotMask
	}
}

func (s *NodeStore) getAuxExtra(id uint32) any {
	if !s.hasExtExtra.Load() {
		return nil
	}
	sh := &s.aux[auxShardIndex(id)]
	if !s.static && isSingleOwnerFastPath() {
		if len(sh.extra) == 0 {
			return nil
		}
		return sh.extra[id]
	}
	sh.mu.RLock()
	var v any
	if len(sh.extra) != 0 {
		v = sh.extra[id]
	}
	sh.mu.RUnlock()
	return v
}

func (s *NodeStore) setAuxExtra(id uint32, v any) {
	if v == nil && !s.hasExtExtra.Load() {
		return
	}
	sh := &s.aux[auxShardIndex(id)]
	singleOwner := !s.static && isSingleOwnerFastPath()
	if !singleOwner {
		sh.mu.Lock()
		defer sh.mu.Unlock()
	}
	if v == nil {
		delete(sh.extra, id)
		return
	}
	if !s.hasExtExtra.Load() {
		s.hasExtExtra.Store(true)
	}
	if sh.extra == nil {
		sh.extra = make(map[uint32]any)
	}
	sh.extra[id] = v
}

const tcColorBit uint32 = 1 << 31

// TC returns the Type/Color handle from the TC column for node n.
func (n Node) TC() uint32 {
	return tcField.Get(n) &^ tcColorBit
}

// SetTC stores the Type/Color handle in the TC column for node n.
func (n Node) SetTC(h uint32) {
	tcField.Set(n, h)
}

// OrigNameID returns the original name interned StringID for node n.
func (n Node) OrigNameID() StringID {
	if blk := nodeBlocks[n>>blockShift].Load(); blk != nil {
		if r := &blk[n&blockMask]; r.KindFlags&kindFlagAnyOrigName != 0 {
			if r.KindFlags >= kindFlagInlineOrigName {
				return StringID(r.Payload)
			}
			return StringID(origNameField[n>>blockShift].Load()[n&blockMask])
		}
	}
	return 0
}

// SetOrigNameID stores the original name interned StringID for node n.
func (n Node) SetOrigNameID(v StringID) {
	rec := n.Rec()
	flags := rec.KindFlags
	if v == 0 {
		if flags&kindFlagAnyOrigName == 0 {
			return
		}
		if flags&kindFlagInlineOrigName != 0 {
			if flags&kindFlagsMask == 0 {
				rec.Payload = 0
			}
			rec.KindFlags = flags &^ kindFlagAnyOrigName
			return
		}
		rec.KindFlags = flags &^ kindFlagHasExtOrigName
		if c := origNameField[n>>blockShift].Load(); c != nil {
			c[n&blockMask] = 0
		}
		return
	}
	if flags&kindFlagsMask == 3 && uint32(v) == rec.Payload {
		if flags&kindFlagInlineOrigName == 0 || flags&kindFlagHasExtOrigName != 0 {
			if flags&kindFlagHasExtOrigName != 0 {
				if c := origNameField[n>>blockShift].Load(); c != nil {
					c[n&blockMask] = 0
				}
			}
			rec.KindFlags = (flags &^ kindFlagHasExtOrigName) | kindFlagInlineOrigName
		}
		return
	}
	if flags&(kindFlagsMask|kindFlagHasChangeTime) == 0 && (rec.Payload == 0 || flags&kindFlagInlineOrigName != 0) {
		if flags&kindFlagHasExtOrigName != 0 {
			if c := origNameField[n>>blockShift].Load(); c != nil {
				c[n&blockMask] = 0
			}
		}
		rec.KindFlags = (flags &^ kindFlagHasExtOrigName) | kindFlagInlineOrigName
		rec.Payload = uint32(v)
		return
	}
	if flags&kindFlagInlineOrigName != 0 {
		if flags&kindFlagsMask == 0 {
			rec.Payload = 0
		}
		flags &^= kindFlagInlineOrigName
	}
	origNameField.Set(n, uint32(v))
	rec.KindFlags = flags | kindFlagHasExtOrigName
}

func tcHandleFor(v any) (uint32, bool) {
	if holder, ok := v.(TCHandleHolder); ok {
		if holder.IsJavaNil_() {
			return 0, false
		}
		hp := holder.TCHandlePtr_()
		if h := hp.Load(); h != 0 {
			return h, true
		}
		return tcHandleRegisterHolder(v, hp), true
	}
	return tcHandleForSlow(v)
}

//go:noinline
func tcHandleRegisterHolder(v any, hp *atomic.Uint32) uint32 {
	newH := tcHandles.Add(&tcHandlesMu, v)
	if hp.CompareAndSwap(0, newH) {
		return newH
	}
	return hp.Load()
}

//go:noinline
func tcHandleForSlow(v any) (uint32, bool) {
	if jn, ok := v.(javaNil); ok && jn.IsJavaNil_() {
		return 0, false
	}
	return tcHandles.Add(&tcHandlesMu, v), true
}

// TypeOrColor returns the type or color value stored via the noscan TC handle column for node n.
func (n Node) TypeOrColor() any {
	if c := tcField[n>>blockShift].Load(); c != nil {
		if h := c[n&blockMask] &^ tcColorBit; h-1 < tableMaxBlocks*tableBlockLen-1 {
			if ab := tcHandles.blocks[h>>tableBlockShift].Load(); ab != nil {
				return ab[h&tableBlockMask]
			}
		}
	}
	return nil
}

// SetTypeOrColor stores the type or color value via a 32-bit handle in the noscan TC column for node n.
func (n Node) SetTypeOrColor(v any) {
	if n != 0 {
		if c := tcField[n>>blockShift].Load(); c != nil {
			if v == nil {
				if c[n&blockMask] != 0 {
					c[n&blockMask] = 0
				}
				return
			}
			if holder, ok := v.(TCHandleHolder); ok && !holder.IsJavaNil_() {
				if h := holder.TCHandlePtr_().Load(); h != 0 {
					if _, isColor := v.(ColorMarker); isColor {
						h |= tcColorBit
					}
					c[n&blockMask] = h
					return
				}
			}
		} else if v == nil {
			return
		}
	}
	setTypeOrColorSlow(n, v)
}

//go:noinline
func setTypeOrColorSlow(n Node, v any) {
	if n == 0 {
		_ = (*NodeData)(nil).Token
	}
	if v == nil {
		tcField.Set(n, 0)
		return
	}
	h, ok := tcHandleFor(v)
	if !ok {
		tcField.Set(n, 0)
		return
	}
	if _, isColor := v.(ColorMarker); isColor {
		h |= tcColorBit
	}
	tcField.Set(n, h)
}

// JSTypeObj returns the JSType stored on n, or nil if unset or if n holds a Color.
func (n Node) JSTypeObj() any {
	if c := tcField[n>>blockShift].Load(); c != nil {
		if h := c[n&blockMask]; h-1 < tableMaxBlocks*tableBlockLen-1 {
			if ab := tcHandles.blocks[h>>tableBlockShift].Load(); ab != nil {
				return ab[h&tableBlockMask]
			}
		}
	} else if n == 0 {
		_ = (*NodeData)(nil).Token
	}
	return nil
}

// SetJSTypeObj stores or clears the JSType on n, returning false without modifying n if n holds a Color.
func (n Node) SetJSTypeObj(v any) bool {
	if n != 0 {
		if c := tcField[n>>blockShift].Load(); c != nil {
			raw := c[n&blockMask]
			if raw&tcColorBit != 0 {
				return false
			}
			if v == nil {
				if raw != 0 {
					c[n&blockMask] = 0
				}
				return true
			}
			if holder, ok := v.(TCHandleHolder); ok && !holder.IsJavaNil_() {
				if h := holder.TCHandlePtr_().Load(); h != 0 {
					c[n&blockMask] = h
					return true
				}
			}
		} else if v == nil {
			return true
		}
	}
	return setJSTypeObjSlow(n, v)
}

//go:noinline
func setJSTypeObjSlow(n Node, v any) bool {
	if n == 0 {
		_ = (*NodeData)(nil).Token
	}
	c := tcField[n>>blockShift].Load()
	if c != nil && c[n&blockMask]&tcColorBit != 0 {
		return false
	}
	if v == nil {
		if c != nil && c[n&blockMask] != 0 {
			c[n&blockMask] = 0
		}
		return true
	}
	h, ok := tcHandleFor(v)
	if !ok {
		if c != nil && c[n&blockMask] != 0 {
			c[n&blockMask] = 0
		}
		return true
	}
	if c != nil {
		c[n&blockMask] = h
		return true
	}
	tcField.setSlow(n, h)
	return true
}

// ColorObj returns the Color stored on n, or nil if unset or if n holds a JSType.
func (n Node) ColorObj() any {
	if c := tcField[n>>blockShift].Load(); c != nil {
		if h := c[n&blockMask] ^ tcColorBit; h-1 < tableMaxBlocks*tableBlockLen-1 {
			if ab := tcHandles.blocks[h>>tableBlockShift].Load(); ab != nil {
				return ab[h&tableBlockMask]
			}
		}
	} else if n == 0 {
		_ = (*NodeData)(nil).Token
	}
	return nil
}

// SetColorObj stores or clears the Color on n, returning false without modifying n
// if n currently holds a non-Color value (e.g. JSType).
func (n Node) SetColorObj(v any) bool {
	if n != 0 {
		if c := tcField[n>>blockShift].Load(); c != nil {
			raw := c[n&blockMask]
			if raw != 0 && raw&tcColorBit == 0 {
				return false
			}
			if v == nil {
				if raw != 0 {
					c[n&blockMask] = 0
				}
				return true
			}
			if holder, ok := v.(TCHandleHolder); ok && !holder.IsJavaNil_() {
				if h := holder.TCHandlePtr_().Load(); h != 0 {
					c[n&blockMask] = h | tcColorBit
					return true
				}
			}
		} else if v == nil {
			return true
		}
	}
	return setColorObjSlow(n, v)
}

//go:noinline
func setColorObjSlow(n Node, v any) bool {
	if n == 0 {
		_ = (*NodeData)(nil).Token
	}
	c := tcField[n>>blockShift].Load()
	if c != nil {
		if raw := c[n&blockMask]; raw != 0 && raw&tcColorBit == 0 {
			return false
		}
	}
	if v == nil {
		if c != nil && c[n&blockMask] != 0 {
			c[n&blockMask] = 0
		}
		return true
	}
	h, ok := tcHandleFor(v)
	if !ok {
		if c != nil && c[n&blockMask] != 0 {
			c[n&blockMask] = 0
		}
		return true
	}
	if c != nil {
		c[n&blockMask] = h | tcColorBit
		return true
	}
	tcField.setSlow(n, h|tcColorBit)
	return true
}

// ReleaseCurrentStoreJSTypes releases all non-Color JSType entries from tcHandles
// (resetting their cached TCHandlePtr_ handles to 0 and freeing empty blocks) and
// clears non-Color TC/typeBeforeCast references on the active compilation NodeStore
// so GC can immediately reclaim the JSType graph (e.g. after ConvertTypesToColors).
// Color handles (implementing ColorMarker) and handles referenced by staticStore are preserved.
func ReleaseCurrentStoreJSTypes() {
	currentMu.Lock()
	defer currentMu.Unlock()
	releaseJSTypesLocked(current.Load(), true)
}

func releaseJSTypesLocked(s *NodeStore, clearStoreAndNonHolders bool) {
	tcHandlesMu.Lock()
	defer tcHandlesMu.Unlock()

	var staticKeep map[uint32]struct{}
	staticStore.mu.Lock()
	for _, col := range staticStore.tcs {
		if col == nil {
			continue
		}
		for _, raw := range col {
			if h := raw &^ tcColorBit; h != 0 {
				if staticKeep == nil {
					staticKeep = make(map[uint32]struct{})
				}
				staticKeep[h] = struct{}{}
			}
		}
	}
	staticStore.mu.Unlock()

	if n := tcHandles.count.Load(); n > 0 {
		curBlk := n >> tableBlockShift
		lastBlk := min(curBlk, uint32(tableMaxBlocks-1))
		for idx := uint32(0); idx <= lastBlk; idx++ {
			blk := tcHandles.blocks[idx].Load()
			if blk == nil {
				continue
			}
			start := uint32(0)
			if idx == 0 {
				start = 1
			}
			end := uint32(tableBlockLen)
			if idx == lastBlk && curBlk == lastBlk {
				end = (n & tableBlockMask) + 1
			}
			hasLive := false
			for off := start; off < end; off++ {
				v := blk[off]
				if v == nil {
					continue
				}
				if _, isColor := v.(ColorMarker); isColor {
					hasLive = true
					continue
				}
				h := (idx << tableBlockShift) | off
				if staticKeep != nil {
					if _, keep := staticKeep[h]; keep {
						hasLive = true
						continue
					}
				}
				holder, isHolder := v.(TCHandleHolder)
				if !clearStoreAndNonHolders && !isHolder {
					hasLive = true
					continue
				}
				if isHolder {
					holder.TCHandlePtr_().Store(0)
				}
				blk[off] = nil
			}
			if !hasLive && idx < curBlk {
				tcHandles.blocks[idx].Store(nil)
				if idx == 0 {
					tcHandles.blk0.Store(nil)
				}
			}
		}
	}

	if s != nil && !s.static && clearStoreAndNonHolders {
		s.mu.Lock()
		cols := s.tcs
		s.mu.Unlock()
		for _, col := range cols {
			if col == nil {
				continue
			}
			for i, raw := range col {
				if raw != 0 && raw&tcColorBit == 0 {
					col[i] = 0
				}
			}
		}
		for i := range s.aux {
			sh := &s.aux[i]
			sh.mu.Lock()
			if sh.castCount.Load() != 0 {
				var newBloom uint64
				var newCastCount uint32
				for pos := range sh.slots {
					sl := &sh.slots[pos]
					if sl.id == 0 || sl.castH == 0 {
						continue
					}
					if v := s.casts.Load(sl.castH); v != nil {
						if _, isColor := v.(ColorMarker); isColor {
							newBloom |= auxCastBloomMask(sl.id)
							newCastCount++
							continue
						}
					}
					s.casts.clear(sl.castH)
					sh.freeCastH = append(sh.freeCastH, sl.castH)
					sl.castH = 0
					s.clearCastBit(sl.id)
					if (sl.changeTime | sl.length) == 0 {
						sh.count--
					}
					if !s.hasJSDoc(sl.id) {
						if blk := nodeBlocks[sl.id>>blockShift].Load(); blk != nil {
							blk[sl.id&blockMask].KindFlags &^= kindFlagHasJSDocOrCast
						}
					}
				}
				sh.castBloom.Store(newBloom)
				sh.castCount.Store(newCastCount)
			}
			sh.mu.Unlock()
		}
	}
}

// SourceFileByHandle returns the StaticSourceFile object for 1-based handle h, or nil if h == 0.
func SourceFileByHandle(h uint32) any {
	return sourceFileHandles.Load(h)
}

// SourceFileObj returns the StaticSourceFile object stored via NodeData.SourceFile16 for node n.
func (n Node) SourceFileObj() any {
	return sourceFileObj(n, sourceFileObjOverflowSlow)
}

func sourceFileObj(n Node, slow func(Node) any) any {
	if blk := nodeBlocks[n>>blockShift].Load(); blk != nil {
		if sf16 := blk[n&blockMask].SourceFile16; sf16 != 0 {
			if sf16 != 0xFFFF {
				if ab := sourceFileHandles.blocks[sf16>>tableBlockShift].Load(); ab != nil {
					return ab[sf16&tableBlockMask]
				}
				return nil
			}
			return slow(n)
		}
	}
	return nil
}

//go:noinline
func sourceFileObjOverflowSlow(n Node) any {
	return sourceFileHandles.Load(getSourceFileOverflowSlow(n))
}

// SetSourceFileObj stores the StaticSourceFile object via a 32-bit handle for node n.
func (n Node) SetSourceFileObj(v any) {
	if v == nil {
		n.SetSourceFile(0)
		return
	}
	if holder, ok := v.(SourceFileHandleHolder); ok {
		if holder.IsJavaNil_() {
			n.SetSourceFile(0)
			return
		}
		hp := holder.SourceFileHandlePtr_()
		h := hp.Load()
		if h == 0 {
			newH := sourceFileHandles.Add(&sourceFileHandlesMu, v)
			if hp.CompareAndSwap(0, newH) {
				h = newH
			} else {
				h = hp.Load()
			}
		}
		n.SetSourceFile(h)
		return
	}
	if jn, ok := v.(javaNil); ok && jn.IsJavaNil_() {
		n.SetSourceFile(0)
		return
	}
	if reflect.TypeOf(v).Comparable() {
		if existing, ok := sourceFileFallback.Load(v); ok {
			n.SetSourceFile(existing.(uint32))
			return
		}
		newH := sourceFileHandles.Add(&sourceFileHandlesMu, v)
		actual, _ := sourceFileFallback.LoadOrStore(v, newH)
		n.SetSourceFile(actual.(uint32))
		return
	}
	n.SetSourceFile(sourceFileHandles.Add(&sourceFileHandlesMu, v))
}

// JSDocObj returns the JSDocInfo object stored in NodeStore for node n, or nil if none is set.
func (n Node) JSDocObj() any {
	return jsdocObj(n, jsdocSlow)
}

func jsdocObj(n Node, slow func(Node) any) any {
	if blk := nodeBlocks[n>>blockShift].Load(); blk != nil && blk[n&blockMask].KindFlags&kindFlagHasJSDoc != 0 {
		return slow(n)
	}
	return nil
}

//go:noinline
func jsdocSlow(n Node) any {
	return storeForNode(n).getJSDoc(uint32(n))
}

// SetJSDocObj stores or clears the JSDocInfo object in NodeStore for node n.
func (n Node) SetJSDocObj(v any) {
	setJSDocObj(n, v, setJSDocSlow)
}

// ClearJSDoc clears any JSDocInfo on n without interface boxing.
func (n Node) ClearJSDoc() {
	clearJSDoc(n, clearJSDocSlow)
}

func clearJSDoc(n Node, slow func(Node, *NodeData)) {
	rec := n.Rec()
	if rec.KindFlags&kindFlagHasJSDoc == 0 {
		return
	}
	slow(n, rec)
}

//go:noinline
func clearJSDocSlow(n Node, rec *NodeData) {
	s := storeForNode(n)
	id := uint32(n)
	if s.deleteJSDoc(id) && !s.hasAuxTypeBeforeCast(id) {
		rec.KindFlags &^= kindFlagHasJSDoc
	}
}

func setJSDocObj(n Node, v any, slow func(Node, *NodeData, any)) {
	rec := n.Rec()
	if v == nil && rec.KindFlags&kindFlagHasJSDoc == 0 {
		return
	}
	slow(n, rec, v)
}

//go:noinline
func setJSDocSlow(n Node, rec *NodeData, v any) {
	flags := rec.KindFlags
	if v == nil {
		if flags&kindFlagHasJSDoc != 0 {
			clearJSDocSlow(n, rec)
		}
		return
	}
	if jn, ok := v.(javaNil); ok && jn.IsJavaNil_() {
		if flags&kindFlagHasJSDoc != 0 {
			clearJSDocSlow(n, rec)
		}
		return
	}
	storeForNode(n).setJSDoc(uint32(n), v)
	if flags&kindFlagHasJSDoc == 0 {
		rec.KindFlags = flags | kindFlagHasJSDoc
	}
}

// JSTypeBeforeCastObj returns the JSType before cast stored in NodeStore.aux for node n, or nil if none is set.
func (n Node) JSTypeBeforeCastObj() any {
	return jsTypeBeforeCastObj(n, typeBeforeCastSlow)
}

func jsTypeBeforeCastObj(n Node, slow func(Node) any) any {
	if b := castBitField[n>>blockShift].Load(); b != nil && atomic.LoadUint64(&b[(n>>6)&((blockLen>>6)-1)])&(1<<(n&63)) != 0 {
		return slow(n)
	}
	return nil
}

//go:noinline
func typeBeforeCastSlow(n Node) any {
	return storeForNode(n).getAuxTypeBeforeCast(uint32(n))
}

// SetJSTypeBeforeCastObj stores or clears the JSType before cast directly in NodeStore.aux for node n.
func (n Node) SetJSTypeBeforeCastObj(v any) {
	setJSTypeBeforeCastObj(n, v, setJSTypeBeforeCastSlow)
}

// ClearJSTypeBeforeCast clears any JSType before cast on n without interface boxing.
func (n Node) ClearJSTypeBeforeCast() {
	clearJSTypeBeforeCast(n, clearJSTypeBeforeCastSlow)
}

func clearJSTypeBeforeCast(n Node, slow func(Node)) {
	if n != 0 {
		if b := castBitField[n>>blockShift].Load(); b == nil || atomic.LoadUint64(&b[(n>>6)&((blockLen>>6)-1)])&(1<<(n&63)) == 0 {
			return
		}
	}
	slow(n)
}

//go:noinline
func clearJSTypeBeforeCastSlow(n Node) {
	rec := n.Rec()
	if rec.KindFlags&kindFlagHasTypeBeforeCast != 0 {
		clearJSTypeBeforeCastRecSlow(n, rec)
	}
}

func clearJSTypeBeforeCastRecSlow(n Node, rec *NodeData) {
	s := current.Load()
	if uint32(n)&staticBit != 0 {
		s = staticStore
	}
	if !s.hasExtCast.Load() {
		return
	}
	id := uint32(n)
	sh := &s.aux[auxShardIndex(id)]
	if sh.castCount.Load() == 0 {
		return
	}
	mask := auxCastBloomMask(id)
	if sh.castBloom.Load()&mask != mask {
		return
	}
	if s.clearAuxTypeBeforeCastShard(sh, id) && !s.hasJSDoc(id) {
		rec.KindFlags &^= kindFlagHasTypeBeforeCast
	}
}

func setJSTypeBeforeCastObj(n Node, v any, slow func(Node, any)) {
	if v == nil && n != 0 {
		if b := castBitField[n>>blockShift].Load(); b == nil || atomic.LoadUint64(&b[(n>>6)&((blockLen>>6)-1)])&(1<<(n&63)) == 0 {
			return
		}
	}
	slow(n, v)
}

//go:noinline
func setJSTypeBeforeCastSlow(n Node, v any) {
	if n == 0 {
		_ = (*NodeData)(nil).Token
	}
	if v == nil {
		if b := castBitField[n>>blockShift].Load(); b != nil && atomic.LoadUint64(&b[(n>>6)&((blockLen>>6)-1)])&(1<<(n&63)) != 0 {
			clearJSTypeBeforeCastSlow(n)
		}
		return
	}
	if jn, ok := v.(javaNil); ok && jn.IsJavaNil_() {
		if b := castBitField[n>>blockShift].Load(); b != nil && atomic.LoadUint64(&b[(n>>6)&((blockLen>>6)-1)])&(1<<(n&63)) != 0 {
			clearJSTypeBeforeCastSlow(n)
		}
		return
	}
	rec := n.Rec()
	s := current.Load()
	if uint32(n)&staticBit != 0 {
		s = staticStore
	}
	s.setAuxTypeBeforeCastNonNil(uint32(n), v)
	if rec.KindFlags&kindFlagHasTypeBeforeCast == 0 {
		rec.KindFlags |= kindFlagHasTypeBeforeCast
	}
}

// ChangeTime returns the change-time stamp for node n (inline in Payload on plain nodes, or in NodeStore.aux), or 0 if unset.
func (n Node) ChangeTime() int32 {
	return nodeChangeTime(n, changeTimeSlow)
}

func nodeChangeTime(n Node, slow func(Node) int32) int32 {
	if n == 0 {
		_ = (*NodeData)(nil).Token
	}
	if blk := nodeBlocks[n>>blockShift].Load(); blk != nil && blk[n&blockMask].KindFlags&kindFlagHasChangeTime != 0 {
		return slow(n)
	}
	return 0
}

//go:noinline
func changeTimeSlow(n Node) int32 {
	rec := n.Rec()
	s := storeForNode(n)
	if rec.KindFlags&(kindFlagsMask|kindFlagInlineOrigName) == 0 && rec.Payload != 0 {
		if s.hasExtChangeTime.Load() {
			if ct := s.getAuxChangeTime(uint32(n)); ct != 0 {
				return ct
			}
		}
		return int32(rec.Payload)
	}
	return s.getAuxChangeTime(uint32(n))
}

// SetChangeTime stores or clears the change-time stamp for node n (inline in Payload on plain nodes, or in NodeStore.aux).
func (n Node) SetChangeTime(v int32) {
	setNodeChangeTime(n, v, setChangeTimeSlow)
}

func setNodeChangeTime(n Node, v int32, slow func(Node, int32)) {
	if n.Rec().KindFlags&kindFlagHasChangeTime == 0 && v == 0 {
		return
	}
	slow(n, v)
}

//go:noinline
func setChangeTimeSlow(n Node, v int32) {
	rec := n.Rec()
	flags := rec.KindFlags
	s := storeForNode(n)
	if v == 0 {
		if flags&kindFlagHasChangeTime == 0 {
			return
		}
		rec.KindFlags = flags &^ kindFlagHasChangeTime
		if s.setAuxChangeTime(uint32(n), 0) {
			return
		}
		if flags&(kindFlagsMask|kindFlagInlineOrigName) == 0 {
			rec.Payload = 0
		}
		return
	}
	if flags&kindFlagsMask == 0 {
		if flags&kindFlagInlineOrigName != 0 {
			if rec.Payload != 0 {
				origNameField.Set(n, rec.Payload)
				flags = (flags &^ kindFlagInlineOrigName) | kindFlagHasExtOrigName
			} else {
				flags &^= kindFlagInlineOrigName
			}
			rec.Payload = uint32(v)
			rec.KindFlags = flags | kindFlagHasChangeTime
			s.setAuxChangeTime(uint32(n), 0)
			return
		}
		if rec.Payload == 0 || (flags&kindFlagHasChangeTime != 0 && s.getAuxChangeTime(uint32(n)) == 0) {
			rec.Payload = uint32(v)
			rec.KindFlags = flags | kindFlagHasChangeTime
			s.setAuxChangeTime(uint32(n), 0)
			return
		}
	}
	s.setAuxChangeTime(uint32(n), v)
	rec.KindFlags = flags | kindFlagHasChangeTime
}

// ClonePropsFrom copies inline SourceFile16 (plus overflow SourceFile), PropFlags, JSDoc, TypeBeforeCast, ChangeTime, and cold PropList from other to n.
// It returns false without modifying n if n already has any property set.
func (n Node) ClonePropsFrom(other Node) bool {
	thisRec := n.Rec()
	if thisRec.SourceFile16 != 0 || thisRec.PropFlags != 0 || thisRec.KindFlags&kindFlagClonedMask != 0 {
		return false
	}
	otherRec := other.Rec()
	sf16 := otherRec.SourceFile16
	thisRec.SourceFile16 = sf16
	thisRec.PropFlags = otherRec.PropFlags
	otherFlags := otherRec.KindFlags
	if otherFlags&kindFlagClonedMask != 0 || sf16 == 0xFFFF {
		n.cloneColdPropsFrom(thisRec, other, otherFlags)
	}
	return true
}

//go:noinline
func (n Node) cloneColdPropsFrom(thisRec *NodeData, other Node, otherFlags uint8) {
	if thisRec.SourceFile16 == 0xFFFF {
		setSFOverflow(n, getSourceFileOverflowSlow(other))
	}
	if otherFlags&kindFlagHasJSDocOrCast != 0 {
		otherStore := storeForNode(other)
		otherID := uint32(other)
		jsDoc := otherStore.getJSDoc(otherID)
		var cast any
		if b := castBitField[otherID>>blockShift].Load(); b != nil && atomic.LoadUint64(&b[(otherID>>6)&((blockLen>>6)-1)])&(1<<(otherID&63)) != 0 {
			cast = otherStore.getAuxTypeBeforeCast(otherID)
		}
		if jsDoc != nil || cast != nil {
			thisStore := storeForNode(n)
			thisID := uint32(n)
			if jsDoc != nil {
				thisStore.setJSDoc(thisID, jsDoc)
			}
			if cast != nil {
				thisStore.setAuxTypeBeforeCast(thisID, cast)
			}
			thisRec.KindFlags |= kindFlagHasJSDocOrCast
		}
	}
	if otherFlags&kindFlagHasChangeTime != 0 {
		n.SetChangeTime(other.ChangeTime())
	}
	if otherFlags&kindFlagHasColdProps != 0 {
		n.SetPropList(other.PropList())
	}
}

// NumberVal returns the float64 value stored for a NumberNode via NodeData.Payload.
func (n Node) NumberVal() float64 {
	h := n.Rec().Payload
	if h == 0 {
		return 0
	}
	s := storeForNode(n)
	return math.Float64frombits(s.nums.load(h))
}

// SetNumberVal stores a float64 value for a NumberNode via NodeData.Payload, keeping
// NodeData.PropFlags free for inline node properties.
func (n Node) SetNumberVal(v float64) {
	rec := n.Rec()
	if rec.KindFlags&kindFlagInlineOrigName != 0 {
		n.materializeOrigName(rec, rec.Payload)
	} else if rec.KindFlags&(kindFlagsMask|kindFlagHasChangeTime) == kindFlagHasChangeTime {
		n.migrateInlineChangeTime(rec)
	}
	bits := math.Float64bits(v)
	if bits == 0 {
		rec.Payload = 0
		return
	}
	s := storeForNode(n)
	if oldH := rec.Payload; oldH != 0 && s.nums.load(oldH) == bits {
		return
	}
	rec.Payload = s.nums.add(&s.numMu, bits)
}

// PropList returns the cold property-list head for node n, or nil if none is set.
func (n Node) PropList() any {
	return nodePropList(n, propListSlow)
}

func nodePropList(n Node, slow func(Node) any) any {
	if blk := nodeBlocks[n>>blockShift].Load(); blk != nil && blk[n&blockMask].KindFlags&kindFlagHasColdProps != 0 {
		return slow(n)
	}
	return nil
}

//go:noinline
func propListSlow(n Node) any {
	id := uint32(n)
	s := storeForNode(n)
	sh := &s.propShards[(id^(id>>6))&(propShards-1)]
	if !s.static && isSingleOwnerFastPath() {
		return sh.m[id]
	}
	sh.mu.RLock()
	v := sh.m[id]
	sh.mu.RUnlock()
	return v
}

// SetPropList stores or deletes the cold property-list head for node n in its NodeStore.
func (n Node) SetPropList(v any) {
	rec := n.Rec()
	if v == nil {
		if rec.KindFlags&kindFlagHasColdProps == 0 {
			return
		}
		rec.KindFlags &^= kindFlagHasColdProps
		n.deletePropListSlow()
		return
	}
	if jn, ok := v.(javaNil); ok && jn.IsJavaNil_() {
		if rec.KindFlags&kindFlagHasColdProps == 0 {
			return
		}
		rec.KindFlags &^= kindFlagHasColdProps
		n.deletePropListSlow()
		return
	}
	id := uint32(n)
	s := storeForNode(n)
	sh := &s.propShards[(id^(id>>6))&(propShards-1)]
	if !s.static && isSingleOwnerFastPath() {
		if sh.m == nil {
			sh.m = make(map[uint32]any)
		}
		sh.m[id] = v
	} else {
		sh.mu.Lock()
		if sh.m == nil {
			sh.m = make(map[uint32]any)
		}
		sh.m[id] = v
		sh.mu.Unlock()
	}
	if rec.KindFlags&kindFlagHasColdProps == 0 {
		rec.KindFlags |= kindFlagHasColdProps
	}
}

//go:noinline
func (n Node) deletePropListSlow() {
	id := uint32(n)
	s := storeForNode(n)
	sh := &s.propShards[(id^(id>>6))&(propShards-1)]
	if !s.static && isSingleOwnerFastPath() {
		delete(sh.m, id)
		return
	}
	sh.mu.Lock()
	delete(sh.m, id)
	sh.mu.Unlock()
}

// GetAux retrieves the cold metadata for node n from the active store.
func (n Node) GetAux() (NodeAux, bool) {
	id := uint32(n)
	if id == 0 {
		return NodeAux{}, false
	}
	s := storeForNode(n)
	jsDoc := s.getJSDoc(id)
	extra := s.getAuxExtra(id)
	if jsDoc == nil && extra == nil {
		return NodeAux{}, false
	}
	return NodeAux{JSDoc: jsDoc, Extra: extra}, true
}

// SetAux stores cold metadata for node n in the active store.
func (n Node) SetAux(a NodeAux) {
	id := uint32(n)
	if id == 0 {
		return
	}
	s := storeForNode(n)
	if a.JSDoc != nil {
		s.setJSDoc(id, a.JSDoc)
	} else {
		s.deleteJSDoc(id)
	}
	s.setAuxExtra(id, a.Extra)
}
