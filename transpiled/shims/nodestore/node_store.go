// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261104

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
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

const (
	staticBit = 1 << 31 // IDs of the static store

	blockShift = 14 // 16,384 records (512 KiB) per block
	blockLen   = 1 << blockShift
	blockMask  = blockLen - 1
	numBlocks  = 1 << (32 - blockShift)

	// Threshold to wrap block ID allocation when persistent CLI workers run many compilations.
	nextBlockWrapThreshold = (staticBit >> blockShift) / 2
)

// BlockShift is the log2 number of NodeData records per NodeBlock (14 -> 16,384).
const BlockShift = blockShift

// NodeBlock is a 16,384-record (512 KiB) noscan block of NodeData records.
type NodeBlock [blockLen]NodeData

type nodeBlock = NodeBlock

// Global 2-level directory mapping block index to block pointer.
// Access compiles to 2 pointer dereferences with 0 mutex contention.
var nodeBlocks [numBlocks]atomic.Pointer[nodeBlock]

// NodeStore manages memory blocks and handle assignment for a compilation unit.
type NodeStore struct {
	id           uint32
	static       bool
	base         uint32
	next         atomic.Uint32
	blocks       []*nodeBlock
	mu           sync.Mutex
	blockAllocMu sync.Mutex

	tcs              []*NodeColumn
	origNames        []*NodeColumn
	len16Pages       []len16Record
	lengthOverflows  []lenOverflowRecord
	negLenNodes      atomic.Uint32
	sfOverflows      []sfOverflowRecord
	castBits         []*castBitBlock
	nums             numTable
	numMu            sync.Mutex
	hasExtChangeTime atomic.Bool
	hasExtCast       atomic.Bool
	hasExtExtra      atomic.Bool
	propShards       [propShards]propShard
	jsdocShards      [jsdocShards]jsdocShard
	aux              [auxShards]auxShard
	jsdocValMu       sync.Mutex
	jsdocs           AnyTable
	castValMu        sync.Mutex
	casts            AnyTable
}

const fastAllocBatch = 1024

var (
	nodeStoreLastID atomic.Uint32
	nextBlock       atomic.Uint32
	staticStore     = newNodeStore(true)

	// current is the store new nodes go to; written between passes only.
	current           atomic.Pointer[NodeStore]
	currentMu         sync.Mutex
	standaloneStoreMu sync.Mutex

	storeResetHooksMu sync.Mutex
	storeResetHooks   []func()

	// Fast single-owner batch cursor for AllocNodeInit / AllocNodeToken / AllocNodeRaw.
	// Invariant: fastCur, fastLast, fastBlk, and fastStore are ONLY mutated by
	// the single-owner goroutine while fastOwnerAllowed.Load() && fastSuspendCount.Load() == 0.
	// Whenever concurrency begins (SuspendSingleOwnerAlloc transitioning 0 -> 1) or ends
	// (ExitSingleOwnerAlloc), FlushFastAlloc resets fastCur = 0 and fastLast = 0 before
	// any happens-before edge to other goroutines, so concurrent callers only read 0 < 0.
	fastCur   Node
	fastLast  Node
	fastBlk   *nodeBlock
	fastSlice []NodeData
	fastStore *NodeStore

	fastOwnerAllowed atomic.Bool
	fastSuspendCount atomic.Int32
)

func init() {
	initLen16Top()
	nextBlock.Store(1) // block 0 is never used: ID 0 is NilNode
	current.Store(staticStore)
}

// LockStandaloneStore acquires the process-wide standalone compilation mutex.
func LockStandaloneStore() {
	standaloneStoreMu.Lock()
}

// UnlockStandaloneStore releases the process-wide standalone compilation mutex.
func UnlockStandaloneStore() {
	standaloneStoreMu.Unlock()
}

// NewNodeStore returns an empty NodeStore. Call SetCurrentNodeStore to make it active.
func NewNodeStore() *NodeStore {
	return newNodeStore(false)
}

func newNodeStore(static bool) *NodeStore {
	s := &NodeStore{id: nodeStoreLastID.Add(1), static: static}
	if static {
		s.setBase(staticBit)
	}
	return s
}

func (s *NodeStore) setBase(base uint32) {
	s.base = base
	s.next.Store(base)
}

func (s *NodeStore) register() {
	if neg := s.negLenNodes.Load(); neg != 0 {
		negLenCount.Add(neg)
	}
	first := s.base >> blockShift
	tcField.register(first, s.tcs)
	for _, r := range s.len16Pages {
		storeLen16Page(r.idx, r.page)
	}
	for _, r := range s.lengthOverflows {
		storeLenOverflowPage(r.idx, r.page)
	}
	for _, r := range s.sfOverflows {
		storeSFOverflowPage(r.idx, r.page)
	}
	origNameField.register(first, s.origNames)
	castBitField.register(first, s.castBits)
	for k, b := range s.blocks {
		nodeBlocks[first+uint32(k)].Store(b)
	}
}

func (s *NodeStore) unregister() {
	first := s.base >> blockShift
	for k, b := range s.blocks {
		nodeBlocks[first+uint32(k)].CompareAndSwap(b, nil)
	}
	tcField.unregister(first, s.tcs)
	origNameField.unregister(first, s.origNames)
	for _, r := range s.len16Pages {
		casLen16Page(r.idx, r.page)
	}
	for _, r := range s.lengthOverflows {
		casLenOverflowPage(r.idx, r.page)
	}
	for _, r := range s.sfOverflows {
		casSFOverflowPage(r.idx, r.page)
	}
	castBitField.unregister(first, s.castBits)
	if neg := s.negLenNodes.Load(); neg != 0 {
		negLenCount.Add(-neg)
	}
}

func (s *NodeStore) blockEnd() uint32 {
	return s.base>>blockShift + uint32(len(s.blocks))
}

// CurrentNodeStore returns the currently active NodeStore for node allocations.
func CurrentNodeStore() *NodeStore {
	return current.Load()
}

// FlushFastAlloc returns any unconsumed node IDs from the single-owner batch
// reservation back to the NodeStore and clears the fast-path cursor.
// Unconsumed slots in (fastCur, fastLast] were never written since new(nodeBlock),
// so rewinding s.next leaves them zero-initialized.
func FlushFastAlloc() {
	if fastStore != nil {
		if fastLast > fastCur {
			fastStore.next.CompareAndSwap(uint32(fastLast)+1, uint32(fastCur)+1)
		}
		fastCur = 0
		fastLast = 0
		fastBlk = nil
		fastSlice = nil
		fastStore = nil
	}
	SealActivePoolChunks()
}

// CanUseSingleOwnerAlloc reports whether no concurrent Java threads or pool tasks
// have suspended single-owner fast allocation.
func CanUseSingleOwnerAlloc() bool {
	return fastSuspendCount.Load() == 0
}

// IsSingleOwnerAllocActive reports whether a single-owner fast allocation session
// is currently entered.
func IsSingleOwnerAllocActive() bool {
	return fastOwnerAllowed.Load()
}

// TryEnterSingleOwnerAlloc enables fast 1024-ID batch allocation for the calling
// goroutine if no single-owner session is already active and no concurrent tasks
// have suspended it. If it returns true, the caller must call ExitSingleOwnerAlloc
// on the same goroutine before handing control back to another goroutine.
func TryEnterSingleOwnerAlloc() bool {
	if fastSuspendCount.Load() != 0 {
		return false
	}
	if !fastOwnerAllowed.CompareAndSwap(false, true) {
		return false
	}
	if fastSuspendCount.Load() != 0 {
		fastOwnerAllowed.Store(false)
		return false
	}
	return true
}

// EnterSingleOwnerAlloc enables fast 1024-ID batch allocation for the calling
// single-owner goroutine when no concurrent tasks are active.
func EnterSingleOwnerAlloc() {
	_ = TryEnterSingleOwnerAlloc()
}

// ExitSingleOwnerAlloc flushes any remaining batch reservation and disables
// single-owner fast allocation. Must be called on the single-owner goroutine
// before signaling completion (e.g. closing Future.done) to another goroutine.
func ExitSingleOwnerAlloc() {
	if !fastOwnerAllowed.Load() {
		return
	}
	if fastSuspendCount.Load() == 0 {
		FlushFastAlloc()
	}
	fastOwnerAllowed.Store(false)
}

// SuspendSingleOwnerAlloc flushes any active single-owner batch on the calling
// goroutine and suspends single-owner fast allocation while a concurrent thread
// or pool task is in flight. Must be called on the parent/submitting goroutine
// BEFORE spawning the goroutine or sending the task over a channel.
func SuspendSingleOwnerAlloc() {
	if fastOwnerAllowed.Load() && fastSuspendCount.Load() == 0 {
		FlushFastAlloc()
	}
	fastSuspendCount.Add(1)
}

// ResumeSingleOwnerAlloc decrements the concurrent suspension count. Must be
// called by the finishing worker/thread BEFORE signaling completion (e.g. closing
// Future.done) to the waiting parent goroutine.
func ResumeSingleOwnerAlloc() {
	rem := fastSuspendCount.Add(-1)
	if rem < 0 {
		panic("nodestore: negative fastSuspendCount")
	}
	if rem == 0 {
		SealActivePoolChunks()
	}
}

// RegisterStoreResetHook registers fn to be invoked whenever an active non-static
// compilation NodeStore is unregistered in SetCurrentNodeStore (including ResetNodeStore).
func RegisterStoreResetHook(fn func()) {
	if fn == nil {
		return
	}
	storeResetHooksMu.Lock()
	storeResetHooks = append(storeResetHooks, fn)
	storeResetHooksMu.Unlock()
}

func runStoreResetHooks() {
	storeResetHooksMu.Lock()
	hooks := storeResetHooks
	storeResetHooksMu.Unlock()
	for _, fn := range hooks {
		fn()
	}
}

// SetCurrentNodeStore sets s as the active NodeStore and returns the previous active store.
func SetCurrentNodeStore(s *NodeStore) *NodeStore {
	if s == staticStore {
		s = nil
	}
	currentMu.Lock()
	if fastOwnerAllowed.Load() && fastSuspendCount.Load() == 0 {
		FlushFastAlloc()
	} else {
		SealActivePoolChunks()
	}
	prev := current.Load()
	if prev == staticStore {
		prev = nil
	}
	runHooks := false
	if prev != s {
		if prev != nil {
			raiseNextBlock(prev.blockEnd())
			prev.unregister()
			if s == nil {
				releaseJSTypesLocked(prev, false)
			}
			runHooks = true
		}
		if s != nil {
			if s.base == 0 {
				s.setBase(nextBlock.Load() << blockShift)
			}
			s.register()
			current.Store(s)
		} else {
			current.Store(staticStore)
		}
	}
	if s == nil && nextBlock.Load() >= nextBlockWrapThreshold {
		nextBlock.Store(1)
	}
	currentMu.Unlock()
	if runHooks {
		runStoreResetHooks()
	}
	return prev
}

func raiseNextBlock(v uint32) {
	for {
		old := nextBlock.Load()
		if v <= old || nextBlock.CompareAndSwap(old, v) {
			return
		}
	}
}

// ResetNodeStore deactivates the current NodeStore.
func ResetNodeStore() {
	SetCurrentNodeStore(nil)
}

// EnsureCompilationStore activates a fresh NodeStore if none is active.
func EnsureCompilationStore() *NodeStore {
	currentMu.Lock()
	defer currentMu.Unlock()
	s := current.Load()
	if s == nil || s == staticStore {
		if fastOwnerAllowed.Load() && fastSuspendCount.Load() == 0 {
			FlushFastAlloc()
		}
		s = NewNodeStore()
		if s.base == 0 {
			s.setBase(nextBlock.Load() << blockShift)
		}
		s.register()
		current.Store(s)
	}
	return s
}

// IsStatic reports whether s is the static template store.
func (s *NodeStore) IsStatic() bool {
	return s.static
}

// NodeCount returns the number of nodes allocated in s.
func (s *NodeStore) NodeCount() int {
	if s.base == 0 {
		return 0
	}
	n := int(s.next.Load() - s.base)
	if fastOwnerAllowed.Load() && fastSuspendCount.Load() == 0 && fastStore == s && fastLast > fastCur {
		n -= int(fastLast - fastCur)
	}
	return n
}

// Alloc allocates a new zeroed node handle in store s.
func (s *NodeStore) Alloc() Node {
	id := s.next.Add(1) - 1
	if nodeBlocks[id>>blockShift].Load() == nil {
		s.addBlock(id)
	}
	return Node(id)
}

// AllocNode allocates a new zeroed node handle in the current store.
func AllocNode() Node {
	return CurrentNodeStore().Alloc()
}

type allocKind uint8

// AllocNodeFast pops a fresh zero-initialized node handle from the single-owner batch
// reservation (or allocates one via allocNodeFastSlow).
func AllocNodeFast() Node {
	if fastCur < fastLast {
		fastCur++
		return fastCur
	}
	return allocNodeFastSlow()
}

// AllocNodeFastRec pops a fresh zero-initialized node handle and its NodeData record pointer
// from the single-owner batch reservation without a second nodeBlocks lookup.
func AllocNodeFastRec() (id Node, rec *NodeData) {
	if fastCur < fastLast {
		fastCur++
		id, rec = fastCur, &fastSlice[fastCur&blockMask]
	} else {
		id, rec = allocNodeRawSlow()
	}
	return
}

//go:noinline
func allocNodeFastSlow() Node {
	id, _ := allocNodeRawSlow()
	return id
}

//go:noinline
func allocNodeRawSlow() (Node, *NodeData) {
	s := current.Load()
	if s != nil && !s.static && fastOwnerAllowed.Load() && fastSuspendCount.Load() == 0 {
		for {
			cur := s.next.Load()
			remInBlk := uint32(blockLen) - (cur & blockMask)
			nReserve := uint32(fastAllocBatch)
			if remInBlk < nReserve {
				nReserve = remInBlk
			}
			if s.next.CompareAndSwap(cur, cur+nReserve) {
				blk := nodeBlocks[cur>>blockShift].Load()
				if blk == nil {
					blk = s.addBlock(cur)
				}
				fastStore = s
				fastBlk = blk
				fastSlice = blk[:]
				fastCur = Node(cur)
				fastLast = Node(cur + nReserve - 1)
				return Node(cur), &blk[cur&blockMask]
			}
		}
	}
	id := s.next.Add(1) - 1
	blk := nodeBlocks[id>>blockShift].Load()
	if blk == nil {
		blk = s.addBlock(id)
	}
	return Node(id), &blk[id&blockMask]
}

// AllocNodeInit allocates a new node handle in the current store and initializes its
// default Token (0xFF), LinenoCharno (-1), and KindFlags (kind) in a single lookup.
func AllocNodeInit(kind uint8) Node {
	return allocNodeInit(kind, allocNodeInitSlow)
}

func allocNodeInit(kind uint8, slow func(uint8) Node) Node {
	if fastCur < fastLast {
		fastCur++
		rec := &fastBlk[fastCur&blockMask]
		rec.Token, rec.KindFlags, rec.LinenoCharno = 0xFF, kind&kindFlagsMask, -1
		return fastCur
	}
	return slow(kind)
}

//go:noinline
func allocNodeInitSlow(kind uint8) Node {
	id, rec := allocNodeRawSlow()
	rec.Token = 0xFF
	rec.KindFlags = kind & kindFlagsMask
	rec.LinenoCharno = -1
	return id
}

// AllocNodeToken allocates a new node handle in the current store and initializes its
// Token, LinenoCharno (-1), and KindFlags (kind) in a single lookup.
func AllocNodeToken(kind, token uint8) Node {
	return allocNodeToken(kind, token, allocNodeTokenSlow)
}

func allocNodeToken(kind, token uint8, slow func(uint8, uint8) Node) Node {
	if fastCur < fastLast {
		fastCur++
		rec := &fastBlk[fastCur&blockMask]
		rec.Token, rec.KindFlags, rec.LinenoCharno = token, kind&kindFlagsMask, -1
		return fastCur
	}
	return slow(kind, token)
}

//go:noinline
func allocNodeTokenSlow(kind, token uint8) Node {
	id, rec := allocNodeRawSlow()
	rec.Token = token
	rec.KindFlags = kind & kindFlagsMask
	rec.LinenoCharno = -1
	return id
}

// NewTemplateNode creates an immutable template node in the static store.
func NewTemplateNode(token uint8) Node {
	n := staticStore.Alloc()
	n.SetToken(token)
	return n
}

func (s *NodeStore) addBlock(id uint32) *nodeBlock {
	s.blockAllocMu.Lock()
	blkIdx := int((id - s.base) >> blockShift)
	s.mu.Lock()
	if blkIdx < len(s.blocks) {
		if blk := s.blocks[blkIdx]; blk != nil {
			s.mu.Unlock()
			s.blockAllocMu.Unlock()
			return blk
		}
	}
	s.mu.Unlock()

	active := !s.static && current.Load() == s
	if active && blkIdx >= 2 {
		chunk := new([4]nodeBlock)
		first := s.base >> blockShift
		s.mu.Lock()
		for len(s.blocks) < blkIdx+4 {
			s.blocks = append(s.blocks, nil)
		}
		for i := 0; i < 4; i++ {
			idx := blkIdx + i
			if s.blocks[idx] == nil {
				b := &chunk[i]
				s.blocks[idx] = b
				if current.Load() == s {
					nodeBlocks[first+uint32(idx)].Store(b)
				}
			}
		}
		blk := s.blocks[blkIdx]
		s.mu.Unlock()
		s.blockAllocMu.Unlock()
		return blk
	}

	fresh := new(nodeBlock)
	s.mu.Lock()
	for len(s.blocks) <= blkIdx {
		s.blocks = append(s.blocks, nil)
	}
	blk := s.blocks[blkIdx]
	if blk == nil {
		blk = fresh
		s.blocks[blkIdx] = blk
		if s.static || current.Load() == s {
			nodeBlocks[id>>blockShift].Store(blk)
		}
	}
	s.mu.Unlock()
	s.blockAllocMu.Unlock()
	return blk
}

// isSingleOwnerFastPath reports whether single-owner mode is active with zero
// concurrent worker goroutines in flight.
func isSingleOwnerFastPath() bool {
	return fastOwnerAllowed.Load() && fastSuspendCount.Load() == 0
}

// --- NodeAllocator: Thread-Local Batch Allocator ---

type internCacheTable [256]StringID

var internCachePool = sync.Pool{
	New: func() any {
		return new(internCacheTable)
	},
}

// NodeAllocator reserves contiguous single-block batches of node handles (256 or 1024)
// to eliminate atomic contention and L1/L2 cache-line false sharing across parse workers.
type NodeAllocator struct {
	store         *NodeStore
	blk           *nodeBlock
	internCache   *internCacheTable
	next, end     Node
	batchEnd      Node
	restEnd       uint32
	rangeSize     uint32
	tmplNode      Node
	tmplPropFlags uint32
	tmplSF16      uint16
	tmplCold      bool
	internDirty   bool
}

const (
	allocatorDefaultRange = 1024
	allocatorLocalRange   = 256
)

// NewNodeAllocator creates a batch allocator for store s (or CurrentNodeStore() if nil).
func NewNodeAllocator(s *NodeStore) *NodeAllocator {
	if s == nil {
		s = CurrentNodeStore()
	}
	return &NodeAllocator{store: s, rangeSize: allocatorDefaultRange}
}

// InitLocal configures a for single-file local allocation against the current NodeStore
// with a 256-node batch size; call Flush when done to return unused trailing IDs.
func (a *NodeAllocator) InitLocal() {
	a.store = CurrentNodeStore()
	a.blk = nil
	a.next = 0
	a.end = 0
	a.batchEnd = 0
	a.restEnd = 0
	a.rangeSize = allocatorLocalRange
	a.tmplNode = 0
	a.tmplPropFlags = 0
	a.tmplSF16 = 0
	a.tmplCold = false
	if a.internDirty {
		if a.internCache != nil {
			clear(a.internCache[:])
		}
		a.internDirty = false
	}
}

// Flush returns any unused trailing IDs in [a.next, a.restEnd) back to a.store when no
// concurrent reservation has advanced past restEnd, and clears a.
func (a *NodeAllocator) Flush() {
	if a.restEnd > uint32(a.next) && a.store != nil {
		a.store.next.CompareAndSwap(a.restEnd, uint32(a.next))
	}
	c := a.internCache
	if c != nil {
		if a.internDirty {
			clear(c[:])
		}
		internCachePool.Put(c)
	}
	*a = NodeAllocator{}
}

//go:noinline
func (a *NodeAllocator) initInternCache() *internCacheTable {
	c := internCachePool.Get().(*internCacheTable)
	a.internCache = c
	return c
}

// Intern returns the interned StringID for s using a goroutine-local 256-slot
// direct-mapped L1 cache to bypass global strpool locks on repeated identifiers.
func (a *NodeAllocator) Intern(s string) StringID {
	if len(s) == 0 {
		return EmptyStringID
	}
	c := a.internCache
	if c == nil {
		c = a.initInternCache()
	}
	idx := (uint32(len(s))*2654435761 ^ uint32(s[0]) ^ (uint32(s[len(s)-1]) << 8) ^ (uint32(s[len(s)>>1]) << 16)) & 255
	if id := c[idx]; id != 0 && String(id) == s {
		return id
	}
	id := Intern(s)
	c[idx] = id
	a.internDirty = true
	return id
}

// SetTemplate caches tmpl's SourceFile16 and PropFlags so AllocNode1..4, AllocStringNode,
// and AllocNumberNode stamp template properties inline without calling ClonePropsFrom.
func (a *NodeAllocator) SetTemplate(tmpl Node) {
	if tmpl == 0 {
		a.tmplNode = 0
		a.tmplPropFlags = 0
		a.tmplSF16 = 0
		a.tmplCold = false
		a.end = a.batchEnd
		return
	}
	rec := tmpl.Rec()
	sf16 := rec.SourceFile16
	cold := rec.KindFlags&kindFlagClonedMask != 0 || sf16 == 0xFFFF
	a.tmplNode = tmpl
	a.tmplPropFlags = rec.PropFlags
	a.tmplSF16 = sf16
	a.tmplCold = cold
	if cold {
		a.end = 0
	} else {
		a.end = a.batchEnd
	}
}

// Alloc allocates a zeroed node handle from the allocator's reservation batch.
func (a *NodeAllocator) Alloc() Node {
	if a.next < a.batchEnd {
		id := a.next
		a.next++
		return id
	}
	return a.allocSlow()
}

//go:noinline
func (a *NodeAllocator) allocSlow() Node {
	id, _ := a.allocRecSlow()
	return id
}

func (a *NodeAllocator) allocRec() (id Node, rec *NodeData) {
	if a.next < a.batchEnd {
		n := a.next
		a.next = n + 1
		id, rec = n, &a.blk[n&blockMask]
	} else {
		id, rec = a.allocRecSlow()
	}
	return
}

//go:noinline
func (a *NodeAllocator) allocRecSlow() (Node, *NodeData) {
	s := a.store
	if s == nil {
		s = CurrentNodeStore()
		a.store = s
	}
	if s.static {
		id := s.Alloc()
		return id, id.Rec()
	}
	rangeSize := a.rangeSize
	if rangeSize == 0 {
		rangeSize = allocatorDefaultRange
	}
	for {
		cur := s.next.Load()
		remInBlk := uint32(blockLen) - (cur & blockMask)
		nReserve := rangeSize
		if remInBlk < nReserve {
			nReserve = remInBlk
		}
		if s.next.CompareAndSwap(cur, cur+nReserve) {
			blk := nodeBlocks[cur>>blockShift].Load()
			if blk == nil {
				blk = s.addBlock(cur)
			}
			a.blk = blk
			a.next = Node(cur + 1)
			end := cur + nReserve
			a.batchEnd = Node(end)
			a.restEnd = end
			if a.tmplCold {
				a.end = 0
			} else {
				a.end = Node(end)
			}
			return Node(cur), &blk[cur&blockMask]
		}
	}
}

func (a *NodeAllocator) applyTemplateSlow(id Node, rec *NodeData, tmpl Node) {
	if a.tmplNode == tmpl && !a.tmplCold {
		rec.SourceFile16 = a.tmplSF16
		rec.PropFlags = a.tmplPropFlags
		return
	}
	if tmpl != 0 {
		id.ClonePropsFrom(tmpl)
	}
}

func checkDetachedAllocatorRec(c Node) *NodeData {
	if c == 0 {
		panic(javalang.NewNullPointerException())
	}
	rec := c.Rec()
	if (rec.Parent | rec.Next | rec.Previous) != 0 {
		checkDetachedAllocatorFail(c, rec)
	}
	return rec
}

//go:noinline
func checkDetachedAllocatorFail(c Node, rec *NodeData) {
	if rec.Parent != 0 {
		panic(javalang.NewIllegalStateException(fmt.Sprintf("Has parent: %d", c)))
	}
	if rec.Next != 0 {
		panic(javalang.NewIllegalStateException(fmt.Sprintf("Has next: %d", c)))
	}
	if rec.Previous != 0 {
		panic(javalang.NewIllegalStateException(fmt.Sprintf("Has previous: %d", c)))
	}
}

// AllocNode1Cached allocates a leaf node with token tok using the cached SetTemplate state.
func (a *NodeAllocator) AllocNode1Cached(tok uint8) Node {
	return allocNode1Cached(a, tok, (*NodeAllocator).allocNode1CachedSlow)
}

func allocNode1Cached(a *NodeAllocator, tok uint8, slow func(*NodeAllocator, uint8) Node) Node {
	if id := a.next; id < a.end {
		a.next++
		r := &a.blk[id&blockMask]
		r.Token, r.SourceFile16, r.PropFlags, r.LinenoCharno = tok, a.tmplSF16, a.tmplPropFlags, -1
		return id
	}
	return slow(a, tok)
}

//go:noinline
func (a *NodeAllocator) allocNode1CachedSlow(tok uint8) Node {
	return a.allocNode1Slow(tok, a.tmplNode)
}

// AllocNode1 allocates a leaf node with token tok and template properties from tmpl.
func (a *NodeAllocator) AllocNode1(tok uint8, tmpl Node) Node {
	return allocNode1(a, tok, tmpl, (*NodeAllocator).allocNode1Slow)
}

func allocNode1(a *NodeAllocator, tok uint8, tmpl Node, slow func(*NodeAllocator, uint8, Node) Node) Node {
	if id := a.next; id < a.end && a.tmplNode == tmpl {
		a.next++
		r := &a.blk[id&blockMask]
		r.Token, r.SourceFile16, r.PropFlags, r.LinenoCharno = tok, a.tmplSF16, a.tmplPropFlags, -1
		return id
	}
	return slow(a, tok, tmpl)
}

//go:noinline
func (a *NodeAllocator) allocNode1Slow(tok uint8, tmpl Node) Node {
	id, rec := a.allocRec()
	rec.Token = tok
	rec.LinenoCharno = -1
	a.applyTemplateSlow(id, rec, tmpl)
	return id
}

// AllocNode2 allocates a 1-child parent node with token tok, detached child c1, and template tmpl.
func (a *NodeAllocator) AllocNode2(tok uint8, c1, tmpl Node) Node {
	if id := a.next; id < a.end && a.tmplNode == tmpl && c1^id < blockLen {
		if c := &a.blk[c1&blockMask]; (c.Parent | c.Previous | c.Next) == 0 {
			a.next = id + 1
			r := &a.blk[id&blockMask]
			r.First = c1
			r.Token = tok
			r.SourceFile16 = a.tmplSF16
			r.PropFlags = a.tmplPropFlags
			r.LinenoCharno = -1
			c.Parent = id
			c.Previous = c1
			return id
		}
	}
	return a.allocNode2Slow(tok, c1, tmpl)
}

//go:noinline
func (a *NodeAllocator) allocNode2Slow(tok uint8, c1, tmpl Node) Node {
	id, rec := a.allocRec()
	rec.First = c1
	rec.Token = tok
	rec.LinenoCharno = -1
	c1Rec := checkDetachedAllocatorRec(c1)
	c1Rec.Parent = id
	c1Rec.Previous = c1
	a.applyTemplateSlow(id, rec, tmpl)
	return id
}

// AllocNode3 allocates a 2-child parent node with token tok, detached children c1, c2, and template tmpl.
func (a *NodeAllocator) AllocNode3(tok uint8, c1, c2, tmpl Node) Node {
	if id := a.next; id < a.end && a.tmplNode == tmpl &&
		((c1^id)|(c2^id))>>blockShift == 0 && c1 != c2 {
		c1Rec, c2Rec := &a.blk[c1&blockMask], &a.blk[c2&blockMask]
		if (c1Rec.Parent | c1Rec.Previous | c1Rec.Next | c2Rec.Parent | c2Rec.Previous | c2Rec.Next) == 0 {
			a.next++
			rec := &a.blk[id&blockMask]
			rec.First = c1
			rec.Token = tok
			rec.SourceFile16 = a.tmplSF16
			rec.PropFlags = a.tmplPropFlags
			rec.LinenoCharno = -1
			c1Rec.Next = c2
			c1Rec.Parent = id
			c1Rec.Previous = c2
			c2Rec.Parent = id
			c2Rec.Previous = c1
			return id
		}
	}
	return a.allocNode3Slow(tok, c1, c2, tmpl)
}

//go:noinline
func (a *NodeAllocator) allocNode3Slow(tok uint8, c1, c2, tmpl Node) Node {
	id, rec := a.allocRec()
	rec.First = c1
	rec.Token = tok
	rec.LinenoCharno = -1
	c1Rec := checkDetachedAllocatorRec(c1)
	c1Rec.Next = c2
	c1Rec.Parent = id
	c1Rec.Previous = c2
	c2Rec := checkDetachedAllocatorRec(c2)
	c2Rec.Parent = id
	c2Rec.Previous = c1
	a.applyTemplateSlow(id, rec, tmpl)
	return id
}

// AllocNode4 allocates a 3-child parent node with token tok, detached children c1, c2, c3, and template tmpl.
func (a *NodeAllocator) AllocNode4(tok uint8, c1, c2, c3, tmpl Node) Node {
	if id := a.next; id < a.end && a.tmplNode == tmpl &&
		((c1^id)|(c2^id)|(c3^id))>>blockShift == 0 &&
		c1 != c2 && c1 != c3 && c2 != c3 {
		c1Rec := &a.blk[c1&blockMask]
		c2Rec := &a.blk[c2&blockMask]
		c3Rec := &a.blk[c3&blockMask]
		if (c1Rec.Parent | c1Rec.Previous | c1Rec.Next |
			c2Rec.Parent | c2Rec.Previous | c2Rec.Next |
			c3Rec.Parent | c3Rec.Previous | c3Rec.Next) == 0 {
			a.next++
			rec := &a.blk[id&blockMask]
			rec.First = c1
			rec.Token = tok
			rec.SourceFile16 = a.tmplSF16
			rec.PropFlags = a.tmplPropFlags
			rec.LinenoCharno = -1
			c1Rec.Next = c2
			c1Rec.Parent = id
			c1Rec.Previous = c3
			c2Rec.Next = c3
			c2Rec.Parent = id
			c2Rec.Previous = c1
			c3Rec.Parent = id
			c3Rec.Previous = c2
			return id
		}
	}
	return a.allocNode4Slow(tok, c1, c2, c3, tmpl)
}

//go:noinline
func (a *NodeAllocator) allocNode4Slow(tok uint8, c1, c2, c3, tmpl Node) Node {
	id, rec := a.allocRec()
	rec.First = c1
	rec.Token = tok
	rec.LinenoCharno = -1
	c1Rec := checkDetachedAllocatorRec(c1)
	c1Rec.Next = c2
	c1Rec.Parent = id
	c1Rec.Previous = c3
	c2Rec := checkDetachedAllocatorRec(c2)
	c2Rec.Next = c3
	c2Rec.Parent = id
	c2Rec.Previous = c1
	c3Rec := checkDetachedAllocatorRec(c3)
	c3Rec.Parent = id
	c3Rec.Previous = c2
	a.applyTemplateSlow(id, rec, tmpl)
	return id
}

// AllocStringNodeCached allocates a StringNode (KindFlags = 3) with token tok and interned strID
// using the cached SetTemplate state.
func (a *NodeAllocator) AllocStringNodeCached(tok uint8, strID StringID) Node {
	return allocStringNodeCached(a, tok, strID, (*NodeAllocator).allocStringNodeCachedSlow)
}

func allocStringNodeCached(a *NodeAllocator, tok uint8, strID StringID, slow func(*NodeAllocator, uint8, StringID) Node) Node {
	if id := a.next; id < a.end {
		a.next++
		r := &a.blk[id&blockMask]
		r.Token, r.KindFlags, r.SourceFile16, r.Payload, r.PropFlags, r.LinenoCharno = tok, 3, a.tmplSF16, uint32(strID), a.tmplPropFlags, -1
		return id
	}
	return slow(a, tok, strID)
}

//go:noinline
func (a *NodeAllocator) allocStringNodeCachedSlow(tok uint8, strID StringID) Node {
	return a.allocStringNodeSlow(tok, strID, a.tmplNode)
}

// AllocStringNode allocates a StringNode (KindFlags = 3) with token tok, interned strID, and template tmpl.
func (a *NodeAllocator) AllocStringNode(tok uint8, strID StringID, tmpl Node) Node {
	return allocStringNode(a, tok, strID, tmpl, (*NodeAllocator).allocStringNodeSlow)
}

func allocStringNode(a *NodeAllocator, tok uint8, strID StringID, tmpl Node, slow func(*NodeAllocator, uint8, StringID, Node) Node) Node {
	if id := a.next; id < a.end && a.tmplNode == tmpl {
		a.next++
		a.blk[id&blockMask] = NodeData{
			Token:        tok,
			KindFlags:    3,
			SourceFile16: a.tmplSF16,
			Payload:      uint32(strID),
			PropFlags:    a.tmplPropFlags,
			LinenoCharno: -1,
		}
		return id
	}
	return slow(a, tok, strID, tmpl)
}

//go:noinline
func (a *NodeAllocator) allocStringNodeSlow(tok uint8, strID StringID, tmpl Node) Node {
	id, rec := a.allocRec()
	rec.Token = tok
	rec.KindFlags = 3
	rec.Payload = uint32(strID)
	rec.LinenoCharno = -1
	a.applyTemplateSlow(id, rec, tmpl)
	return id
}

// AllocNumberNode allocates a NumberNode (KindFlags = 1) with token tok, IEEE-754 float64 val, and template tmpl.
// +0.0 stores Payload = 0 without touching numTable; -0.0 (bits == 0x8000000000000000) allocates a non-zero handle.
func (a *NodeAllocator) AllocNumberNode(tok uint8, val float64, tmpl Node) Node {
	if id := a.next; id < a.end && a.tmplNode == tmpl {
		a.next++
		rec := &a.blk[id&blockMask]
		rec.Token = tok
		rec.KindFlags = 1
		if bits := math.Float64bits(val); bits != 0 {
			rec.Payload = a.allocNumBits(bits)
		}
		rec.SourceFile16 = a.tmplSF16
		rec.PropFlags = a.tmplPropFlags
		rec.LinenoCharno = -1
		return id
	}
	return a.allocNumberNodeSlow(tok, val, tmpl)
}

//go:noinline
func (a *NodeAllocator) allocNumBits(bits uint64) uint32 {
	return a.store.nums.add(&a.store.numMu, bits)
}

//go:noinline
func (a *NodeAllocator) allocNumberNodeSlow(tok uint8, val float64, tmpl Node) Node {
	id, rec := a.allocRec()
	rec.Token = tok
	rec.KindFlags = 1
	rec.LinenoCharno = -1
	if bits := math.Float64bits(val); bits != 0 {
		s := a.store
		if s == nil || (uint32(id)&staticBit != 0) != s.static {
			s = storeForNode(id)
		}
		rec.Payload = s.nums.add(&s.numMu, bits)
	}
	a.applyTemplateSlow(id, rec, tmpl)
	return id
}

// AllocWithToken allocates a node handle and initializes its token.
func (a *NodeAllocator) AllocWithToken(token uint8) Node {
	n := a.Alloc()
	n.SetToken(token)
	return n
}

// AllocCachedTemplate creates a single leaf node with a token.
func (a *NodeAllocator) AllocCachedTemplate(token uint8) Node {
	return a.AllocWithToken(token)
}

// AllocChildCachedTemplate creates a parent node with one child.
func (a *NodeAllocator) AllocChildCachedTemplate(token uint8, child Node) Node {
	parent := a.AllocWithToken(token)
	parent.SetFirstChild(child)
	if child != 0 {
		child.SetParent(parent)
		child.SetPrevious(child)
	}
	return parent
}

// Alloc2ChildrenCachedTemplate creates a parent node with two children.
func (a *NodeAllocator) Alloc2ChildrenCachedTemplate(token uint8, c1, c2 Node) Node {
	parent := a.AllocWithToken(token)
	parent.SetFirstChild(c1)
	if c1 != 0 {
		c1.SetParent(parent)
		c1.SetNext(c2)
		if c2 != 0 {
			c2.SetParent(parent)
			c2.SetPrevious(c1)
			c1.SetPrevious(c2)
		} else {
			c1.SetPrevious(c1)
		}
	}
	return parent
}

// Alloc3ChildrenCachedTemplate creates a parent node with three children.
func (a *NodeAllocator) Alloc3ChildrenCachedTemplate(token uint8, c1, c2, c3 Node) Node {
	parent := a.AllocWithToken(token)
	parent.SetFirstChild(c1)
	if c1 != 0 {
		c1.SetParent(parent)
		c1.SetNext(c2)
		if c2 != 0 {
			c2.SetParent(parent)
			c2.SetNext(c3)
			if c3 != 0 {
				c3.SetParent(parent)
				c3.SetPrevious(c2)
				c1.SetPrevious(c3)
			} else {
				c1.SetPrevious(c2)
			}
		} else {
			c1.SetPrevious(c1)
		}
	}
	return parent
}
