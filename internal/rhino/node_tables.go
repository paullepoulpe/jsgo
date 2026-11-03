package rhino

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/colors"
)

// Go-only: the side tables of a NodeStore.
//
// What most nodes do not have is kept beside the record instead of in it,
// so the record stays small and holds no pointer for the GC to trace:
//
//   - tc: a node's JSType or Color, as a 32-bit handle into the store's
//     types or colors table, in a column per block (see "tc" below).
//   - aux: a node's cold data (JSDoc and the rarely set fields, see
//     nodeAux), in a map of the store (see "aux" below).
//
// The handle tables:
//
//   - types: a JSType remembers its handle (JSTypeBase.nodeHandle); the
//     registry that owns the types clears their entries with
//     ReleaseNodeTypeHandles (called when the compiler drops its registry),
//     after which GetJSType on a node that still refers to one returns nil.
//   - colors: a Color remembers (store id, handle).
//
// Concurrency: handle and column lookups take no lock; a table only grows
// (under the store's mutex), its blocks never move, and a handle is
// published (stored in a column) only after its entry was written.

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

type handleBlock[P any] [tableBlockLen]atomic.Pointer[P]
type handleSubDir[P any] [tableSubLen]atomic.Pointer[handleBlock[P]]

// handleTable is a lazily allocated table of *P addressed by 32-bit handles
// (0 is never handed out). An empty handleTable adds zero heap allocations to
// NewNodeStore.
type handleTable[P any] struct {
	count atomic.Uint32
	blk0  atomic.Pointer[handleBlock[P]]
	dir   [tableDirLen]atomic.Pointer[handleSubDir[P]]
}

// load returns the entry of handle h.
func (t *handleTable[P]) load(h uint32) *P {
	if h < tableBlockLen {
		return t.blk0.Load()[h].Load()
	}
	k := h >> tableBlockShift
	return t.dir[k>>tableSubShift].Load()[k&tableSubMask].Load()[h&tableBlockMask].Load()
}

// peek returns the entry of h, or nil if h was not handed out by t.
func (t *handleTable[P]) peek(h uint32) *P {
	if h == 0 || h > t.count.Load() {
		return nil
	}
	if h < tableBlockLen {
		blk := t.blk0.Load()
		if blk == nil {
			return nil
		}
		return blk[h].Load()
	}
	k := h >> tableBlockShift
	topIdx := k >> tableSubShift
	if topIdx >= tableDirLen {
		return nil
	}
	sub := t.dir[topIdx].Load()
	if sub == nil {
		return nil
	}
	blk := sub[k&tableSubMask].Load()
	if blk == nil {
		return nil
	}
	return blk[h&tableBlockMask].Load()
}

// add stores v in a new entry and returns its handle; mu is the owning
// store's mutex (taken only to create a directory page or block).
func (t *handleTable[P]) add(mu *sync.Mutex, v *P) uint32 {
	h := t.count.Add(1)
	if h < tableBlockLen {
		blk := t.blk0.Load()
		if blk == nil {
			blk = t.addBlock(mu, 0, 0)
		}
		blk[h].Store(v)
		return h
	}
	k := h >> tableBlockShift
	topIdx := k >> tableSubShift
	if topIdx >= tableDirLen {
		panic("rhino.NodeStore: side table full")
	}
	subIdx := k & tableSubMask
	var blk *handleBlock[P]
	if sub := t.dir[topIdx].Load(); sub != nil {
		blk = sub[subIdx].Load()
	}
	if blk == nil {
		blk = t.addBlock(mu, topIdx, subIdx)
	}
	blk[h&tableBlockMask].Store(v)
	return h
}

//go:noinline
func (t *handleTable[P]) addBlock(mu *sync.Mutex, topIdx, subIdx uint32) *handleBlock[P] {
	mu.Lock()
	defer mu.Unlock()
	sub := t.dir[topIdx].Load()
	if sub == nil {
		sub = new(handleSubDir[P])
		blk := new(handleBlock[P])
		sub[subIdx].Store(blk)
		t.dir[topIdx].Store(sub)
		if topIdx == 0 && subIdx == 0 {
			t.blk0.Store(blk)
		}
		return blk
	}
	blk := sub[subIdx].Load()
	if blk == nil {
		blk = new(handleBlock[P])
		sub[subIdx].Store(blk)
		if topIdx == 0 && subIdx == 0 {
			t.blk0.Store(blk)
		}
	}
	return blk
}

// each calls f with every slot handed out so far.
func (t *handleTable[P]) each(f func(*atomic.Pointer[P])) {
	n := t.count.Load()
	if n == 0 {
		return
	}
	lastBlk := min(n>>tableBlockShift, uint32(tableMaxBlocks-1))
	for k := uint32(0); k <= lastBlk; k++ {
		sub := t.dir[k>>tableSubShift].Load()
		if sub == nil {
			k |= tableSubMask
			continue
		}
		blk := sub[k&tableSubMask].Load()
		if blk == nil {
			continue
		}
		start := uint32(0)
		if k == 0 {
			start = 1
		}
		end := uint32(tableBlockLen)
		if k == lastBlk && (n>>tableBlockShift) == lastBlk {
			end = (n & tableBlockMask) + 1
		}
		for off := start; off < end; off++ {
			f(&blk[off])
		}
	}
}

// nodeAux is the cold data of a node: its ext box when the record has
// kindFlagHasExtBox, else its JSDocInfo when it has propFlagHasJSDoc.
type nodeAux struct {
	ext   *nodeExtBox
	jsDoc *JSDocInfo
}

// --- types ---

// typeHandle returns the handle of t in s, registering it on first use.
func (s *NodeStore) typeHandle(t *JSTypeBase) uint32 {
	if h := t.nodeHandle.Load(); h != 0 && s.types.peek(h) == t {
		return h
	}
	if reg := t.GetRegistry(); reg != nil {
		reg.noteNodeStore(s)
	}
	h := s.types.add(&s.mu, t)
	t.nodeHandle.Store(h)
	return h
}

func (s *NodeStore) typeAt(h uint32) *JSTypeBase {
	if h == 0 {
		return nil
	}
	return s.types.load(h)
}

// noteNodeStore records that r's types are registered in s, so that
// ReleaseNodeTypeHandles clears them there.
func (r *JSTypeRegistry) noteNodeStore(s *NodeStore) {
	r.nodeStoresMu.Lock()
	defer r.nodeStoresMu.Unlock()
	for _, x := range r.nodeStores {
		if x == s {
			return
		}
	}
	r.nodeStores = append(r.nodeStores, s)
}

// ReleaseNodeTypeHandles drops every type of r from the node stores' type
// tables: nodes that still refer to one read a nil JSType from now on, and
// the types are no longer kept alive by the AST.
func (r *JSTypeRegistry) ReleaseNodeTypeHandles() {
	r.nodeStoresMu.Lock()
	stores := r.nodeStores
	r.nodeStores = nil
	r.nodeStoresMu.Unlock()
	for _, s := range stores {
		s.types.each(func(e *atomic.Pointer[JSTypeBase]) {
			if v := e.Load(); v != nil && v.GetRegistry() == r {
				e.CompareAndSwap(v, nil)
			}
		})
	}
}

// --- colors ---

func (s *NodeStore) colorHandle(c *colors.Color) uint32 {
	if storeID, h := c.LoadNodeHandle(); storeID == s.id && s.colors.peek(h) == c {
		return h
	}
	h := s.colors.add(&s.mu, c)
	c.StoreNodeHandle(s.id, h)
	return h
}

func (s *NodeStore) colorAt(h uint32) *colors.Color {
	if h == 0 {
		return nil
	}
	return s.colors.load(h)
}

// --- aux ---

// The cold data of the few nodes that have some (about 1 %) is kept in a
// map per store instead of a field in every record. The map is split in
// shards, each with its own lock, because parallel parsing attaches JSDoc
// from several goroutines; readers only look when the record's flags say
// the node has cold data (propWordJSDoc, kindFlagHasExtBox). Entries of
// dead nodes stay until compaction rebuilds the map for the live ones.

const auxShards = 16

type auxShard struct {
	mu sync.RWMutex
	m  map[uint32]nodeAux
}

// hasAux reports whether the record's flags say n may have cold data.
func (n Node) hasAux() bool {
	r := n.rec()
	return r.propFlags&propWordJSDoc != 0 || r.kindFlags&kindFlagHasExtBox != 0
}

// auxOf returns the node's cold-data entry; ok is false if it has none.
func (n Node) auxOf() (a nodeAux, ok bool) {
	if !n.hasAux() {
		return nodeAux{}, false
	}
	return n.store().auxGet(uint32(n))
}

func (s *NodeStore) auxGet(id uint32) (nodeAux, bool) {
	sh := &s.aux[id%auxShards]
	sh.mu.RLock()
	a, ok := sh.m[id]
	sh.mu.RUnlock()
	return a, ok
}

// setAux stores a as the node's cold data; the zero nodeAux clears it.
func (n Node) setAux(a nodeAux) {
	sh := n.auxShard()
	sh.mu.Lock()
	sh.put(uint32(n), a)
	sh.mu.Unlock()
}

func (n Node) auxShard() *auxShard {
	return &n.store().aux[uint32(n)%auxShards]
}

// put stores a as the entry of id (deleting it if a is zero); the caller
// holds sh.mu.
func (sh *auxShard) put(id uint32, a nodeAux) {
	if a == (nodeAux{}) {
		delete(sh.m, id)
		return
	}
	if sh.m == nil {
		sh.m = make(map[uint32]nodeAux)
	}
	sh.m[id] = a
}

// auxValue returns a copy of the node's cold-data entry (zero if none).
func (n Node) auxValue() nodeAux {
	a, _ := n.auxOf()
	return a
}

// --- columns ---

// A column holds a uint32 per record of a block, beside the block, for a
// field that many nodes do not have: a block's column is allocated when the
// first node in it gets a non-zero value. Two fields live in columns:
//
//   - tc: a node's JSType or Color handle (untyped compilations never
//     allocate a tc column);
//   - origName: a node's original name (the name before renaming, for
//     source maps), when it is not simply the node's own string (see
//     kindFlagOrigNameIsString); in a paged column (see below).
//
// A process-wide directory per field maps a block number to its column, as
// nodeBlocks does for the records; the store owns the columns of its blocks.

type nodeColumn [blockLen]uint32

// columnField is one column-stored field: the directory of its columns.
type columnField struct {
	dir [numBlocks]atomic.Pointer[nodeColumn]
}

var tcField columnField

// get returns n's value (0: none).
func (f *columnField) get(n Node) uint32 {
	if c := f.dir[uint32(n)>>blockShift].Load(); c != nil {
		return c[uint32(n)&blockMask]
	}
	return 0
}

// set sets n's value; cols is the store's slice of columns for f.
func (f *columnField) set(n Node, cols *[]*nodeColumn, v uint32) {
	c := f.dir[uint32(n)>>blockShift].Load()
	if c == nil {
		if v == 0 {
			return
		}
		c = n.store().addColumn(f, cols, uint32(n))
	}
	c[uint32(n)&blockMask] = v
}

// addColumn creates the column of f for id's block (unless another
// goroutine just did).
//
//go:noinline
func (s *NodeStore) addColumn(f *columnField, cols *[]*nodeColumn, id uint32) *nodeColumn {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := int((id - s.base) >> blockShift)
	for len(*cols) <= k {
		*cols = append(*cols, nil)
	}
	if (*cols)[k] == nil {
		(*cols)[k] = new(nodeColumn)
		f.dir[id>>blockShift].Store((*cols)[k])
	}
	return (*cols)[k]
}

// registerColumns puts the columns cols of the blocks starting at block
// number first in f's directory.
func (f *columnField) register(first uint32, cols []*nodeColumn) {
	for k, c := range cols {
		if c != nil {
			f.dir[first+uint32(k)].Store(c)
		}
	}
}

// unregister takes the columns cols of the blocks starting at block number
// first out of f's directory.
func (f *columnField) unregister(first uint32, cols []*nodeColumn) {
	for k, c := range cols {
		if c != nil {
			f.dir[first+uint32(k)].CompareAndSwap(c, nil)
		}
	}
}

// columnAt returns the value at index i of a store's columns (for
// compaction: the columns are not in the directory then).
func columnAt(cols []*nodeColumn, i uint32) uint32 {
	if k := int(i >> blockShift); k < len(cols) && cols[k] != nil {
		return cols[k][i&blockMask]
	}
	return 0
}

// setColumnInBlock sets the value at offset off within a single block's column.
func setColumnInBlock(col **nodeColumn, off uint32, v uint32) {
	c := *col
	if c == nil {
		if v == 0 {
			return
		}
		c = new(nodeColumn)
		*col = c
	}
	c[off] = v
}

// setColumnAt sets the value at index i of a store's columns (see columnAt).
func setColumnAt(cols *[]*nodeColumn, i uint32, v uint32) {
	k := int(i >> blockShift)
	if k >= len(*cols) {
		if v == 0 {
			return
		}
		for len(*cols) <= k {
			*cols = append(*cols, nil)
		}
	}
	setColumnInBlock(&(*cols)[k], i&blockMask, v)
}

// --- paged columns ---

// A paged column is a column split into pages of 64 values, each allocated
// on the first non-zero write: for a field that most nodes of a block do not
// have. The original name is one: when the type checker peaks, about 3 % of
// the nodes have one that is not their own string (most pages stay empty),
// while after renaming most names do.

const (
	columnPageShift = 6
	columnPageLen   = 1 << columnPageShift
	columnPageMask  = columnPageLen - 1
	pagesPerBlock   = blockLen / columnPageLen
)

type columnPage [columnPageLen]uint32

// pagedColumn is the pages of one block (2 KB, where a column is 64 KB).
type pagedColumn [pagesPerBlock]atomic.Pointer[columnPage]

// pagedColumnField is one paged-column field: the directory of its columns.
type pagedColumnField struct {
	dir [numBlocks]atomic.Pointer[pagedColumn]
}

var origNameField pagedColumnField

// get returns n's value (0: none).
func (f *pagedColumnField) get(n Node) uint32 {
	c := f.dir[uint32(n)>>blockShift].Load()
	if c == nil {
		return 0
	}
	if p := c[(uint32(n)&blockMask)>>columnPageShift].Load(); p != nil {
		return p[uint32(n)&columnPageMask]
	}
	return 0
}

// set sets n's value; cols is the store's slice of columns for f.
func (f *pagedColumnField) set(n Node, cols *[]*pagedColumn, v uint32) {
	c := f.dir[uint32(n)>>blockShift].Load()
	if c == nil {
		if v == 0 {
			return
		}
		c = n.store().addPagedColumn(f, cols, uint32(n))
	}
	slot := &c[(uint32(n)&blockMask)>>columnPageShift]
	p := slot.Load()
	if p == nil {
		if v == 0 {
			return
		}
		slot.CompareAndSwap(nil, new(columnPage)) // fails if another goroutine just made it
		p = slot.Load()
	}
	p[uint32(n)&columnPageMask] = v
}

// addPagedColumn creates the paged column of f for id's block (unless
// another goroutine just did).
//
//go:noinline
func (s *NodeStore) addPagedColumn(f *pagedColumnField, cols *[]*pagedColumn, id uint32) *pagedColumn {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := int((id - s.base) >> blockShift)
	for len(*cols) <= k {
		*cols = append(*cols, nil)
	}
	if (*cols)[k] == nil {
		(*cols)[k] = new(pagedColumn)
		f.dir[id>>blockShift].Store((*cols)[k])
	}
	return (*cols)[k]
}

func (f *pagedColumnField) register(first uint32, cols []*pagedColumn) {
	for k, c := range cols {
		if c != nil {
			f.dir[first+uint32(k)].Store(c)
		}
	}
}

func (f *pagedColumnField) unregister(first uint32, cols []*pagedColumn) {
	for k, c := range cols {
		if c != nil {
			f.dir[first+uint32(k)].CompareAndSwap(c, nil)
		}
	}
}

// pagedColumnAt returns the value at index i of a store's paged columns
// (for compaction: the columns are not in the directory then).
func pagedColumnAt(cols []*pagedColumn, i uint32) uint32 {
	if k := int(i >> blockShift); k < len(cols) && cols[k] != nil {
		if p := cols[k][(i&blockMask)>>columnPageShift].Load(); p != nil {
			return p[i&columnPageMask]
		}
	}
	return 0
}

// setPagedColumnInBlock sets the value at offset off within a single block's paged column.
func setPagedColumnInBlock(col **pagedColumn, off uint32, v uint32) {
	c := *col
	if c == nil {
		if v == 0 {
			return
		}
		c = new(pagedColumn)
		*col = c
	}
	slot := &c[off>>columnPageShift]
	p := slot.Load()
	if p == nil {
		if v == 0 {
			return
		}
		p = new(columnPage)
		slot.Store(p)
	}
	p[off&columnPageMask] = v
}

// setPagedColumnAt sets the value at index i of a store's paged columns
// (see pagedColumnAt; no other goroutine runs during compaction).
func setPagedColumnAt(cols *[]*pagedColumn, i uint32, v uint32) {
	k := int(i >> blockShift)
	if k >= len(*cols) {
		if v == 0 {
			return
		}
		for len(*cols) <= k {
			*cols = append(*cols, nil)
		}
	}
	setPagedColumnInBlock(&(*cols)[k], i&blockMask, v)
}

// tc returns n's JSType or Color handle (0: none).
func (n Node) tc() uint32 { return tcField.get(n) }

// setTC sets n's JSType or Color handle.
func (n Node) setTC(h uint32) { tcField.set(n, &n.store().tcs, h) }

// translateTC returns the tc handle in s for the JSType or Color that tc
// names in from (the same handle when from is s).
func (s *NodeStore) translateTC(from *NodeStore, tc uint32, isColor bool) uint32 {
	if tc == 0 || from == s {
		return tc
	}
	if isColor {
		if c := from.colorAt(tc); c != nil {
			return s.colorHandle(c)
		}
		return 0
	}
	return s.translateTypeHandle(from, tc)
}

// translateTypeHandle returns the type handle in s for the JSType that h
// names in from.
func (s *NodeStore) translateTypeHandle(from *NodeStore, h uint32) uint32 {
	if h == 0 || from == s {
		return h
	}
	if t := from.typeAt(h); t != nil {
		return s.typeHandle(t)
	}
	return 0
}
