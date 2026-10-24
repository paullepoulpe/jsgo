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
	tableMaxBlocks  = 1 << 14 // 2^26 handles per table
)

// handleTable is a growable table of T addressed by 32-bit handles (0 is
// never handed out).
type handleTable[T any] struct {
	next   atomic.Uint32
	blocks []atomic.Pointer[[tableBlockLen]T]
}

func (t *handleTable[T]) init() {
	t.blocks = make([]atomic.Pointer[[tableBlockLen]T], tableMaxBlocks)
	t.next.Store(1)
}

// at returns the entry of handle h.
func (t *handleTable[T]) at(h uint32) *T {
	return &t.blocks[h>>tableBlockShift].Load()[h&tableBlockMask]
}

// peek returns the entry of h, or the zero T if h was not handed out by t.
func (t *handleTable[T]) peek(h uint32) T {
	var zero T
	if h == 0 || h >= t.next.Load() {
		return zero
	}
	blk := t.blocks[h>>tableBlockShift].Load()
	if blk == nil {
		return zero
	}
	return blk[h&tableBlockMask]
}

// add stores v in a new entry and returns its handle; mu is the owning
// store's mutex (taken only to create a block).
func (t *handleTable[T]) add(mu *sync.Mutex, v T) uint32 {
	h := t.next.Add(1) - 1
	k := h >> tableBlockShift
	if k >= tableMaxBlocks {
		panic("rhino.NodeStore: side table full")
	}
	blk := t.blocks[k].Load()
	if blk == nil {
		mu.Lock()
		if blk = t.blocks[k].Load(); blk == nil {
			blk = new([tableBlockLen]T)
			t.blocks[k].Store(blk)
		}
		mu.Unlock()
	}
	blk[h&tableBlockMask] = v
	return h
}

// each calls f with every entry handed out so far.
func (t *handleTable[T]) each(f func(*T)) {
	n := t.next.Load()
	for h := uint32(1); h < n; h++ {
		if blk := t.blocks[h>>tableBlockShift].Load(); blk != nil {
			f(&blk[h&tableBlockMask])
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
	if h := t.nodeHandle; h != 0 && s.types.peek(h) == t {
		return h
	}
	if reg := t.GetRegistry(); reg != nil {
		reg.noteNodeStore(s)
	}
	h := s.types.add(&s.mu, t)
	t.nodeHandle = h
	return h
}

func (s *NodeStore) typeAt(h uint32) *JSTypeBase {
	if h == 0 {
		return nil
	}
	return *s.types.at(h)
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
		s.types.each(func(e **JSTypeBase) {
			if *e != nil && (*e).GetRegistry() == r {
				*e = nil
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
	return *s.colors.at(h)
}

// --- aux ---

// The cold data of the few nodes that have some (about 1 %) is kept in a
// map per store instead of a field in every record. The map is split in
// shards, each with its own lock, because parallel parsing attaches JSDoc
// from several goroutines; readers only look when the record's flags say
// the node has cold data (propFlagHasJSDoc, kindFlagHasExtBox). Entries of
// dead nodes stay until compaction copies the live nodes into a new store.

const auxShards = 16

type auxShard struct {
	mu sync.RWMutex
	m  map[uint32]nodeAux
}

// hasAux reports whether the record's flags say n may have cold data.
func (n Node) hasAux() bool {
	r := n.rec()
	return r.propFlags&propFlagHasJSDoc != 0 || r.kindFlags&kindFlagHasExtBox != 0
}

// auxOf returns a copy of the node's cold-data entry, or nil if it has none.
func (n Node) auxOf() *nodeAux {
	if !n.hasAux() {
		return nil
	}
	a, ok := n.store().auxGet(uint32(n))
	if !ok {
		return nil
	}
	return &a
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
	s := n.store()
	id := uint32(n)
	sh := &s.aux[id%auxShards]
	sh.mu.Lock()
	if a == (nodeAux{}) {
		delete(sh.m, id)
	} else {
		if sh.m == nil {
			sh.m = make(map[uint32]nodeAux)
		}
		sh.m[id] = a
	}
	sh.mu.Unlock()
}

// auxValue returns a copy of the node's cold-data entry (zero if none).
func (n Node) auxValue() nodeAux {
	if a := n.auxOf(); a != nil {
		return *a
	}
	return nodeAux{}
}

// --- tc ---

// A node's JSType or Color handle lives in a column beside its block (a
// uint32 per record), allocated for a block when the first node in it gets
// a type or color: untyped compilations never allocate one.

type tcColumn [blockLen]uint32

// nodeTCs is the column directory, parallel to nodeBlocks.
var nodeTCs [numBlocks]atomic.Pointer[tcColumn]

// tc returns n's JSType or Color handle (0: none).
func (n Node) tc() uint32 {
	if c := nodeTCs[uint32(n)>>blockShift].Load(); c != nil {
		return c[uint32(n)&blockMask]
	}
	return 0
}

// setTC sets n's JSType or Color handle.
func (n Node) setTC(h uint32) {
	c := nodeTCs[uint32(n)>>blockShift].Load()
	if c == nil {
		if h == 0 {
			return
		}
		c = n.store().addTCColumn(uint32(n))
	}
	c[uint32(n)&blockMask] = h
}

// addTCColumn creates the tc column of id's block (unless another goroutine
// just did).
//
//go:noinline
func (s *NodeStore) addTCColumn(id uint32) *tcColumn {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := int((id - s.base) >> blockShift)
	for len(s.tcs) <= k {
		s.tcs = append(s.tcs, nil)
	}
	if s.tcs[k] == nil {
		s.tcs[k] = new(tcColumn)
		nodeTCs[id>>blockShift].Store(s.tcs[k])
	}
	return s.tcs[k]
}

// tcAt returns the tc handle of the node at index i of s.
func (s *NodeStore) tcAt(i uint32) uint32 {
	if k := int(i >> blockShift); k < len(s.tcs) && s.tcs[k] != nil {
		return s.tcs[k][i&blockMask]
	}
	return 0
}

// setTCAt sets the tc handle of the node at index i of s (for compaction:
// the column is not registered in the directory yet).
func (s *NodeStore) setTCAt(i uint32, h uint32) {
	k := int(i >> blockShift)
	if k >= len(s.tcs) || s.tcs[k] == nil {
		if h == 0 {
			return
		}
		for len(s.tcs) <= k {
			s.tcs = append(s.tcs, nil)
		}
		s.tcs[k] = new(tcColumn)
	}
	s.tcs[k][i&blockMask] = h
}

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
