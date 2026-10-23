package rhino

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/colors"
)

// Node side tables: a Node refers to its JSType or Color (Node.tc) and to its
// cold data (Node.aux: a *JSDocInfo or a *nodeExtBox) through 32-bit handles
// instead of pointers, so the node carries no pointer the GC has to trace
// for them and so the type column can be dropped wholesale at
// types-to-colors (Java's generational collector leaves the old AST alone;
// Go re-marks every node pointer every cycle, see
// docs/memory-research/round2-design.md).
//
// Go-only: performance.
//
// Layout. A handle is blockIndex<<handleBlockBits | slot. Each of the three
// tables (types, colors, aux) is a process-wide directory of blocks of 4096
// entries. The directory is append-only: block indexes are never reused, so
// a stale handle (one whose block was released) reads nil instead of another
// owner's entry; with 2^20 blocks of 4096 entries a process can hand out
// 4 G handles per table before it would wrap, which it never does (jsgo
// panics instead).
//
// Ownership and release. Entries are written by an allocator (handleAlloc),
// which owns whole blocks and releases them all at once by storing nil in
// their directory slots; the GC then frees the block and everything only it
// kept alive.
//   - JSType handles are allocated by the type's own JSTypeRegistry (lazily,
//     at SetJSType / SetJSTypeBeforeCast) and released by
//     JSTypeRegistry.ReleaseNodeTypeHandles, which the Compiler calls from
//     ClearJSTypeRegistry: every type a node referred to is freed there, even
//     the ones held by dead nodes that are still reachable (a detached
//     subtree in a side map), and GetJSType on such a node returns nil.
//   - Color and aux handles are allocated from the NodeTables that are
//     current when the setter runs (NodeTables.Install, called by the
//     Compiler constructor; outside a compiler, a default NodeTables that is
//     never released). The Compiler releases its NodeTables (and the type
//     handles of the registries created while they were current) when it is
//     garbage collected (runtime.AddCleanup) — nothing in the tables may
//     reach the Compiler (its registry reports through a weak reference, an
//     InputId is never a field of its CompilerInput), and the lifetime test
//     in jscomp guards that.
//     "Current" is the newest tables, so a node of an OLDER compiler that is
//     still being worked on (a test harness keeps the compiler under test
//     while it creates another one to parse the expected output) allocates
//     into the newer tables. The blocks of a tables are therefore freed only
//     once it and every older tables are dead (Release marks; reap frees in
//     creation order), which is exactly when no live compiler can hold a
//     handle into them.
//
// Concurrency. Readers never lock: a lookup is two atomic loads (directory,
// block) and an index; a handle published through a node is published
// after its entry was written, so the reader sees the entry. Allocation
// takes the allocator's mutex (parallel parsing creates JSDoc on different
// files, parallel per-function passes set colors and properties
// concurrently); directory growth takes the table's mutex and publishes the
// new directory atomically. A Color or JSType remembers its handle; two
// goroutines racing to register the same object both succeed with distinct
// handles and the last store wins, which is harmless (a handle only needs
// to resolve back to the object).
//
// What is NOT reclaimed: an entry whose node died keeps its JSDoc or ext box
// alive until the owner is released (the Compiler dies, or the registry is
// cleared). Compaction (round-2 step 4) rebuilds the aux table from the
// live tree; until then this pins the cold data of dead nodes.

const (
	handleBlockBits = 12
	handleBlockSize = 1 << handleBlockBits
	handleSlotMask  = handleBlockSize - 1
	handleMaxBlocks = 1 << (32 - handleBlockBits)
)

type handleBlock [handleBlockSize]unsafe.Pointer

// handleTable is one process-wide directory of blocks.
type handleTable struct {
	mu        sync.Mutex // growth of dir and nextBlock
	dir       atomic.Pointer[[]atomic.Pointer[handleBlock]]
	nextBlock uint32 // next block index to hand out; 0 is never used
}

// get resolves a handle; nil for 0, for a released block and for a handle
// from a block that does not exist (never allocated).
func (t *handleTable) get(h uint32) unsafe.Pointer {
	d := t.dir.Load()
	b := h >> handleBlockBits
	if d == nil || b >= uint32(len(*d)) {
		return nil
	}
	blk := (*d)[b].Load()
	if blk == nil {
		return nil
	}
	return blk[h&handleSlotMask]
}

// block returns the block of handle h, or nil if it was released or never
// allocated.
func (t *handleTable) block(h uint32) *handleBlock {
	d := t.dir.Load()
	b := h >> handleBlockBits
	if d == nil || b >= uint32(len(*d)) {
		return nil
	}
	return (*d)[b].Load()
}

// newBlock appends a fresh block and returns its index (>= 1).
func (t *handleTable) newBlock() (uint32, *handleBlock) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.nextBlock == 0 {
		t.nextBlock = 1
	}
	b := t.nextBlock
	if b >= handleMaxBlocks {
		panic("rhino: node handle table exhausted")
	}
	t.nextBlock = b + 1
	d := t.dir.Load()
	if d == nil || b >= uint32(len(*d)) {
		n := 64
		if d != nil {
			n = len(*d) * 2
		}
		for uint32(n) <= b {
			n *= 2
		}
		nd := make([]atomic.Pointer[handleBlock], n)
		if d != nil {
			for i := range *d {
				nd[i].Store((*d)[i].Load())
			}
		}
		t.dir.Store(&nd)
		d = &nd
	}
	blk := new(handleBlock)
	(*d)[b].Store(blk)
	return b, blk
}

// release drops block b: its entries become unreachable through the table.
func (t *handleTable) release(b uint32) {
	d := t.dir.Load()
	if d != nil && b < uint32(len(*d)) {
		(*d)[b].Store(nil)
	}
}

var (
	nodeTypeTable  handleTable // *JSTypeBase
	nodeColorTable handleTable // *colors.Color
	nodeAuxTable   handleTable // *JSDocInfo or *nodeExtBox (Node.kindFlagHasExtBox says which)
)

// handleAlloc hands out handles of one table from blocks it owns.
type handleAlloc struct {
	tab    *handleTable
	mu     sync.Mutex
	cur    uint32 // next free handle in blk; valid iff cur < limit
	limit  uint32
	blk    *handleBlock
	blocks []uint32
}

// alloc stores v in a fresh slot and returns its handle (never 0).
func (a *handleAlloc) alloc(v unsafe.Pointer) uint32 {
	a.mu.Lock()
	if a.cur >= a.limit {
		b, blk := a.tab.newBlock()
		a.blocks = append(a.blocks, b)
		a.blk = blk
		a.cur = b << handleBlockBits
		a.limit = a.cur + handleBlockSize
	}
	h := a.cur
	a.cur++
	a.blk[h&handleSlotMask] = v
	a.mu.Unlock()
	return h
}

// release drops every block this allocator handed out. Later allocations
// start fresh blocks (which nobody will release: the owner is gone).
func (a *handleAlloc) release() {
	a.mu.Lock()
	for _, b := range a.blocks {
		a.tab.release(b)
	}
	a.blocks = nil
	a.blk = nil
	a.cur, a.limit = 0, 0
	a.mu.Unlock()
}

// auxSet overwrites the entry of an aux handle. The handle's block may have
// been released (the node outlived its owner); the write is then dropped and
// the node reads nil, like every other stale handle.
func auxSet(h uint32, v unsafe.Pointer) {
	if blk := nodeAuxTable.block(h); blk != nil {
		blk[h&handleSlotMask] = v
	}
}

// NodeTables owns the color and aux handles of one compilation, plus the
// node type handles of the registries created while it was current.
type NodeTables struct {
	id         uint32 // creation order; distinguishes the NodeTables a Color registered with
	colors     handleAlloc
	aux        handleAlloc
	mu         sync.Mutex
	registries []*handleAlloc
	dead       bool // Release was called; freed once every older tables is dead too
}

var (
	nodeTablesLastID  atomic.Uint32
	currentNodeTables atomic.Pointer[NodeTables]
	defaultNodeTables = &NodeTables{colors: handleAlloc{tab: &nodeColorTable}, aux: handleAlloc{tab: &nodeAuxTable}}

	// nodeTablesMu guards nodeTablesLive: every installed NodeTables not yet
	// freed, in creation order.
	nodeTablesMu   sync.Mutex
	nodeTablesLive []*NodeTables
)

// NewNodeTables creates tables that nothing uses until Install.
func NewNodeTables() *NodeTables {
	return &NodeTables{
		id:     nodeTablesLastID.Add(1),
		colors: handleAlloc{tab: &nodeColorTable},
		aux:    handleAlloc{tab: &nodeAuxTable},
	}
}

// Install makes t the tables new color and aux handles are allocated from.
func (t *NodeTables) Install() {
	nodeTablesMu.Lock()
	nodeTablesLive = append(nodeTablesLive, t)
	nodeTablesMu.Unlock()
	currentNodeTables.Store(t)
}

// Release marks t dead: its owner no longer uses any handle. The blocks of
// t (and of its registries) are freed once every older tables is dead as
// well, because an older compiler may have allocated into t while t was
// current; nodes that still hold such handles read nil (no JSDoc, no type,
// no color) from then on. If t is current, the default tables become
// current.
func (t *NodeTables) Release() {
	currentNodeTables.CompareAndSwap(t, nil)
	nodeTablesMu.Lock()
	t.dead = true
	var free []*NodeTables
	for len(nodeTablesLive) > 0 && nodeTablesLive[0].dead {
		free = append(free, nodeTablesLive[0])
		nodeTablesLive[0] = nil
		nodeTablesLive = nodeTablesLive[1:]
	}
	if len(nodeTablesLive) == 0 {
		nodeTablesLive = nil
	}
	nodeTablesMu.Unlock()
	for _, f := range free {
		f.free()
	}
}

// free drops every handle block t and its registries own.
func (t *NodeTables) free() {
	t.mu.Lock()
	regs := t.registries
	t.registries = nil
	t.mu.Unlock()
	t.colors.release()
	t.aux.release()
	for _, r := range regs {
		r.release()
	}
}

// adoptRegistry makes t release r's node type handles with its own.
func (t *NodeTables) adoptRegistry(r *JSTypeRegistry) {
	t.mu.Lock()
	t.registries = append(t.registries, r.nodeTypes)
	t.mu.Unlock()
}

func currentTables() *NodeTables {
	if t := currentNodeTables.Load(); t != nil {
		return t
	}
	return defaultNodeTables
}

// --- types ---

// nodeTypeHandle returns the handle of t, registering it with its registry
// on first use (or after the registry released its handles).
func nodeTypeHandle(t *JSTypeBase) uint32 {
	if h := t.nodeHandle; h != 0 && nodeTypeTable.get(h) == unsafe.Pointer(t) {
		return h
	}
	reg := t.GetRegistry()
	var a *handleAlloc
	if reg != nil {
		a = reg.nodeTypes
	} else {
		a = &registrylessTypes
	}
	h := a.alloc(unsafe.Pointer(t))
	t.nodeHandle = h
	return h
}

// registrylessTypes serves types without a registry (test fixtures); it is
// never released.
var registrylessTypes = handleAlloc{tab: &nodeTypeTable}

func nodeTypeAt(h uint32) *JSTypeBase {
	return (*JSTypeBase)(nodeTypeTable.get(h))
}

// ReleaseNodeTypeHandles frees the handles of every type that was attached
// to a Node through this registry. Nodes that still hold one read a nil
// JSType afterwards.
func (r *JSTypeRegistry) ReleaseNodeTypeHandles() {
	r.nodeTypes.release()
}

// --- colors ---

func nodeColorHandle(c *colors.Color) uint32 {
	t := currentTables()
	if ownerID, h := c.LoadNodeHandle(); ownerID == t.id && nodeColorTable.get(h) == unsafe.Pointer(c) {
		return h
	}
	h := t.colors.alloc(unsafe.Pointer(c))
	c.StoreNodeHandle(t.id, h)
	return h
}

func nodeColorAt(h uint32) *colors.Color {
	return (*colors.Color)(nodeColorTable.get(h))
}

// --- aux ---

func auxAlloc(v unsafe.Pointer) uint32 {
	return currentTables().aux.alloc(v)
}
