package rhino

import (
	"fmt"
	"math/bits"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: compaction of the AST (see node_store.go).
//
// A store keeps every node allocated in it, dead or alive. Compaction keeps
// only the live ones, in place:
//
//  1. Mark: the compiler calls Mark for every node it holds; Mark marks the
//     whole tree that node belongs to (with the type expressions of its
//     JSDoc and the subtrees hung off its properties).
//  2. Slide: the live records move down to the front of the store's own
//     blocks, in id order, and the blocks left empty at the end are freed.
//     No second copy of the tree is ever allocated.
//  3. Forward: the compiler calls Forward for every node it holds, again,
//     and stores the result: the live nodes have new ids, in a range no store
//     has used before, so an id kept from before compaction (one that was
//     not forwarded) fails like a nil pointer on its first use.
//
// A mark bit per old id says which nodes are live, and a node's new id is
// the number of live ids before it: no table of new ids is needed.

// Compactor compacts the current node store.
type Compactor struct {
	s         *NodeStore
	base, end uint32   // the store's ids before compaction
	marks     []uint64 // bit i: id base+i is live
	ranks     []uint32 // ranks[w]: the number of live ids before word w
	docs      map[*JSDocInfo]struct{}
	newBase   uint32
	slid      bool

	Live       int // live nodes (after Slide)
	ExtraRoots int // detached trees marked because a holder refers to them
}

// NewCompactor starts the compaction of the current store.
func NewCompactor() *Compactor {
	s := current
	if s == nil || s.base == 0 {
		panic("rhino.Compactor: no current node store")
	}
	end := s.next.Load()
	return &Compactor{
		s:     s,
		base:  s.base,
		end:   end,
		marks: make([]uint64, (end-s.base+63)/64),
		docs:  make(map[*JSDocInfo]struct{}, 1024),
	}
}

// Store returns the store being compacted.
func (k *Compactor) Store() *NodeStore {
	return k.s
}

func (k *Compactor) isMarked(i uint32) bool { return k.marks[i>>6]&(1<<(i&63)) != 0 }

// inStore returns n's index in the store before compaction; ok is false
// for the nil node and the static store's nodes, which compaction leaves
// alone. It panics for a node of an older store.
func (k *Compactor) inStore(n Node) (i uint32, ok bool) {
	id := uint32(n)
	if id == 0 || id&staticBit != 0 {
		return 0, false
	}
	if id < k.base || id >= k.end {
		panic(fmt.Sprintf("rhino.Compactor: node id %d is not in the store being compacted", id))
	}
	return id - k.base, true
}

// Mark marks n and the whole tree it belongs to as live, and returns n.
func (k *Compactor) Mark(n Node) Node {
	if k.slid {
		panic("rhino.Compactor: Mark after Slide")
	}
	i, ok := k.inStore(n)
	if !ok || k.isMarked(i) {
		return n
	}
	r := n
	for r.rec().parent != 0 {
		r = r.parentNode()
	}
	if ri, _ := k.inStore(r); !k.isMarked(ri) {
		k.ExtraRoots++
		k.markTree(r)
	}
	return n
}

// markTree marks the subtree rooted at n.
func (k *Compactor) markTree(n Node) {
	i, ok := k.inStore(n)
	if !ok {
		panic(fmt.Sprintf("rhino.Compactor: a package-level node or template is in a compacted tree: %s", n))
	}
	if k.isMarked(i) {
		return
	}
	k.marks[i>>6] |= 1 << (i & 63)
	r := n.rec()
	if r.propFlags&propWordJSDoc != 0 {
		if info := n.GetJSDocInfo(); info != nil {
			if _, done := k.docs[info]; !done {
				k.docs[info] = struct{}{}
				info.forwardTypeNodes(k.Mark)
			}
		}
	}
	if r.kindFlags&kindFlagHasExtBox != 0 {
		for _, v := range nodeProps(n) {
			k.Mark(v.node)
		}
	}
	for c := n.firstNode(); c != 0; c = c.nextNode() {
		k.markTree(c)
	}
}

type nodeProp struct {
	prop Node_Prop
	node Node
}

// nodeProps returns the node-valued props of n (DECLARED_TYPE_EXPR,
// CLOSURE_UNAWARE_SHADOW).
func nodeProps(n Node) []nodeProp {
	var out []nodeProp
	for x := n.coldPropListHead(); x != nil; x = x.next {
		if x.isObject {
			if v, ok := x.objectValue.(Node); ok && !v.IsNil() {
				out = append(out, nodeProp{Node_Prop(x.propType), v})
			}
		}
	}
	return out
}

// newID returns the new id of the live node with index i.
func (k *Compactor) newID(i uint32) uint32 {
	w := i >> 6
	return k.newBase + k.ranks[w] + uint32(bits.OnesCount64(k.marks[w]&(1<<(i&63)-1)))
}

// link returns the new id for a link to the old id id (0 stays 0).
func (k *Compactor) link(id uint32) uint32 {
	if id == 0 {
		return 0
	}
	i := id - k.base
	if !k.isMarked(i) {
		panic(fmt.Sprintf("rhino.Compactor: a live node links to node %d, which was not marked", id))
	}
	return k.newID(i)
}

// Slide moves the live nodes to the front of the store, under their new
// ids, and frees the blocks left empty.
func (k *Compactor) Slide() {
	if k.slid {
		panic("rhino.Compactor: Slide twice")
	}
	s := k.s
	live := uint32(0)
	k.ranks = make([]uint32, len(k.marks))
	for w, m := range k.marks {
		k.ranks[w] = live
		live += uint32(bits.OnesCount64(m))
	}
	raiseNextBlock(s.blockEnd())
	k.newBase = nextBlock.Load() << blockShift

	// Move the records, and their column values, down. A node's new
	// position is never after its old one, so going up in id order never
	// overwrites a record before it has moved.
	pos := uint32(0)
	for w, m := range k.marks {
		for ; m != 0; m &= m - 1 {
			i := uint32(w)<<6 | uint32(bits.TrailingZeros64(m))
			r := s.blocks[i>>blockShift][i&blockMask]
			r.next, r.first = k.link(r.next), k.link(r.first)
			r.parent, r.previous = k.link(r.parent), k.link(r.previous)
			s.blocks[pos>>blockShift][pos&blockMask] = r
			setColumnAt(&s.tcs, pos, columnAt(s.tcs, i))
			setPagedColumnAt(&s.origNames, pos, pagedColumnAt(s.origNames, i))
			pos++
		}
	}

	// Free the blocks left empty, and clear the rest of the last one: new
	// nodes expect zeroed records.
	keep := int((live + blockLen - 1) >> blockShift)
	for p := live; p < uint32(keep)<<blockShift; p++ {
		s.blocks[p>>blockShift][p&blockMask] = nodeData{}
		setColumnAt(&s.tcs, p, 0)
		setPagedColumnAt(&s.origNames, p, 0)
	}
	s.unregister()
	clear(s.blocks[keep:])
	s.blocks = s.blocks[:keep]
	if len(s.tcs) > keep {
		clear(s.tcs[keep:])
		s.tcs = s.tcs[:keep]
	}
	if len(s.origNames) > keep {
		clear(s.origNames[keep:])
		s.origNames = s.origNames[:keep]
	}
	s.base = k.newBase
	s.next.Store(k.newBase + live)
	s.register()

	// The cold data is keyed by id (no other goroutine runs during
	// compaction, so the shards are not locked).
	var oldAux [auxShards]map[uint32]nodeAux
	var oldCold [auxShards]map[uint32]uint32
	for sh := range s.aux {
		oldAux[sh], s.aux[sh].m = s.aux[sh].m, nil
		oldCold[sh], s.aux[sh].cold = s.aux[sh].cold, nil
	}
	newNode := func(id uint32) (Node, bool) {
		if i := id - k.base; id >= k.base && id < k.end && k.isMarked(i) {
			return Node(k.newID(i)), true
		}
		return NilNode, false
	}
	for _, m := range oldAux {
		for id, a := range m {
			if nn, ok := newNode(id); ok {
				nn.auxShard().put(uint32(nn), a)
			}
		}
	}
	for _, m := range oldCold {
		for id, c := range m {
			if nn, ok := newNode(id); ok {
				nn.auxShard().putCold(uint32(nn), c)
			}
		}
	}
	k.slid = true
	k.Live = int(live)

	// Node ids held by the JSDoc of live nodes and by node-valued props.
	for info := range k.docs {
		info.forwardTypeNodes(k.Forward)
	}
	for id := k.newBase; id < k.newBase+live; id++ {
		n := Node(id)
		if n.rec().kindFlags&kindFlagHasExtBox != 0 {
			for _, v := range nodeProps(n) {
				n.PutProp(v.prop, k.Forward(v.node))
			}
		}
	}
}

// Forward returns the new id of a node marked before Slide (or n itself for
// the nil node, a static node, or a node already forwarded).
func (k *Compactor) Forward(n Node) Node {
	if !k.slid {
		panic("rhino.Compactor: Forward before Slide")
	}
	id := uint32(n)
	if id == 0 || id&staticBit != 0 || (id >= k.newBase && id < k.newBase+uint32(k.Live)) {
		return n
	}
	i, _ := k.inStore(n)
	if !k.isMarked(i) {
		panic(fmt.Sprintf("rhino.Compactor: node id %d is forwarded but was not marked", id))
	}
	return Node(k.newID(i))
}

// ForwardSlice forwards every element of s in place.
func (k *Compactor) ForwardSlice(s []Node) {
	for i, n := range s {
		s[i] = k.Forward(n)
	}
}

// Finish ends the compaction.
func (k *Compactor) Finish() {
	k.marks, k.ranks, k.docs = nil, nil, nil
}

// String describes the compaction (for logs).
func (k *Compactor) String() string {
	return fmt.Sprintf("kept %d of %d nodes (%d extra roots)", k.Live, k.end-k.base, k.ExtraRoots)
}

// forwardTypeNodes forwards the roots of the type-expression trees and type
// transformations the JSDocInfo holds.
func (i *JSDocInfo) forwardTypeNodes(fwd func(Node) Node) {
	for it := i.TypeExprIter(); it.Next(); {
		e := it.Expr()
		if !e.root.IsNil() {
			e.root = fwd(e.root)
		}
	}
	if m := jsdocGet[*jutil.LinkedHashMap[string, Node]](JSDocInfo_TYPE_TRANSFORMATIONS.Get(i)); m != nil {
		for _, key := range m.KeySlice() {
			m.Put(key, fwd(m.Get(key)))
		}
	}
}
