package rhino

import (
	"fmt"
	"math/bits"
	"sort"
	"sync"
	"sync/atomic"

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
	liveCount uint32
	ranked    bool
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
	k.ranked = false
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

// Total returns the total number of node slots in the store snapshot.
func (k *Compactor) Total() int {
	return int(k.end - k.base)
}

func (k *Compactor) ensureRanks() uint32 {
	if !k.ranked {
		live := uint32(0)
		if cap(k.ranks) >= len(k.marks) {
			k.ranks = k.ranks[:len(k.marks)]
		} else {
			k.ranks = make([]uint32, len(k.marks))
		}
		for w, m := range k.marks {
			k.ranks[w] = live
			live += uint32(bits.OnesCount64(m))
		}
		k.liveCount = live
		k.ranked = true
	}
	return k.liveCount
}

// CountMarked returns the number of marked live nodes.
func (k *Compactor) CountMarked() int {
	return int(k.ensureRanks())
}

// locateRank finds the marks word index, masked remaining bits starting at
// targetPos, and source index for the targetPos-th live node (0 <= targetPos < live).
func (k *Compactor) locateRank(targetPos uint32) (w int, m uint64, srcIdx uint32) {
	// Largest w with k.ranks[w] <= targetPos; guaranteed to have k.marks[w] != 0.
	w = sort.Search(len(k.ranks), func(i int) bool {
		return k.ranks[i] > targetPos
	}) - 1
	m = k.marks[w]
	skip := int(targetPos - k.ranks[w])
	for range skip {
		m &= m - 1
	}
	srcIdx = uint32(w)<<6 | uint32(bits.TrailingZeros64(m))
	return w, m, srcIdx
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
	k.SlideWithThreads(1)
}

// SlideWithThreads moves the live nodes to the front of the store using up to
// threads parallel workers across destination block stripes.
func (k *Compactor) SlideWithThreads(threads int) {
	if k.slid {
		panic("rhino.Compactor: Slide twice")
	}
	s := k.s
	live := k.ensureRanks()
	raiseNextBlock(s.blockEnd())
	k.newBase = nextBlock.Load() << blockShift

	// Unregister old blocks and columns from the global directories before any
	// block or column pointers in s are overwritten or trimmed.
	s.unregister()

	keep := int((live + blockLen - 1) >> blockShift)
	workers := min(max(threads, 1), keep)

	if workers <= 1 {
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
	} else {
		hasTCs := len(s.tcs) > 0
		if hasTCs {
			for len(s.tcs) < keep {
				s.tcs = append(s.tcs, nil)
			}
		}
		hasOrigNames := len(s.origNames) > 0
		if hasOrigNames {
			for len(s.origNames) < keep {
				s.origNames = append(s.origNames, nil)
			}
		}

		type stripe struct {
			startBlock int
			endBlock   int
			posStart   uint32
			posEnd     uint32
			startW     int
			startMask  uint64
			tailSrcBlk int
		}
		stripes := make([]stripe, workers)
		for u := range workers {
			sb := u * keep / workers
			eb := (u + 1) * keep / workers
			ps := uint32(sb) << blockShift
			pe := min(uint32(eb)<<blockShift, live)
			sw, sm, _ := k.locateRank(ps)
			_, _, lastSrc := k.locateRank(pe - 1)
			stripes[u] = stripe{
				startBlock: sb,
				endBlock:   eb,
				posStart:   ps,
				posEnd:     pe,
				startW:     sw,
				startMask:  sm,
				tailSrcBlk: int(lastSrc >> blockShift),
			}
		}

		// Destination blocks b <= tailSrcBlk[u-1] are still read by earlier
		// stripes; write those into fresh blocks so s.blocks[b] stays read-only
		// until all workers finish. Non-overlapping blocks write in place.
		outBlocks := make([]*nodeBlock, keep)
		var outTCs []*nodeColumn
		if hasTCs {
			outTCs = make([]*nodeColumn, keep)
		}
		var outOrigs []*pagedColumn
		if hasOrigNames {
			outOrigs = make([]*pagedColumn, keep)
		}
		for u := 1; u < workers; u++ {
			overlapEnd := min(stripes[u].endBlock, stripes[u-1].tailSrcBlk+1)
			for b := stripes[u].startBlock; b < overlapEnd; b++ {
				outBlocks[b] = new(nodeBlock)
			}
		}

		var wg sync.WaitGroup
		var firstPanic atomic.Pointer[any]
		runStripe := func(st stripe) {
			defer func() {
				if r := recover(); r != nil {
					firstPanic.CompareAndSwap(nil, &r)
				}
			}()
			pos := st.posStart
			limit := st.posEnd
			w := st.startW
			m := st.startMask
			for pos < limit {
				for m == 0 {
					w++
					m = k.marks[w]
				}
				i := uint32(w)<<6 | uint32(bits.TrailingZeros64(m))
				m &= m - 1

				r := s.blocks[i>>blockShift][i&blockMask]
				r.next, r.first = k.link(r.next), k.link(r.first)
				r.parent, r.previous = k.link(r.parent), k.link(r.previous)

				dstB := int(pos >> blockShift)
				dstOff := pos & blockMask
				if ob := outBlocks[dstB]; ob != nil {
					ob[dstOff] = r
					if hasTCs {
						setColumnInBlock(&outTCs[dstB], dstOff, columnAt(s.tcs, i))
					}
					if hasOrigNames {
						setPagedColumnInBlock(&outOrigs[dstB], dstOff, pagedColumnAt(s.origNames, i))
					}
				} else {
					s.blocks[dstB][dstOff] = r
					if hasTCs {
						setColumnInBlock(&s.tcs[dstB], dstOff, columnAt(s.tcs, i))
					}
					if hasOrigNames {
						setPagedColumnInBlock(&s.origNames[dstB], dstOff, pagedColumnAt(s.origNames, i))
					}
				}
				pos++
			}
		}
		for u := 1; u < workers; u++ {
			st := stripes[u]
			wg.Add(1)
			go func() {
				defer wg.Done()
				runStripe(st)
			}()
		}
		runStripe(stripes[0])
		wg.Wait()
		if p := firstPanic.Load(); p != nil {
			panic(*p)
		}

		for b := range keep {
			if outBlocks[b] != nil {
				s.blocks[b] = outBlocks[b]
				if hasTCs {
					s.tcs[b] = outTCs[b]
				}
				if hasOrigNames {
					s.origNames[b] = outOrigs[b]
				}
			}
		}
	}

	// Free the blocks left empty, and clear the rest of the last one: new
	// nodes expect zeroed records.
	for p := live; p < uint32(keep)<<blockShift; p++ {
		s.blocks[p>>blockShift][p&blockMask] = nodeData{}
		setColumnAt(&s.tcs, p, 0)
		setPagedColumnAt(&s.origNames, p, 0)
	}
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
	var extPropNodes []Node
	for _, m := range oldAux {
		// order-independent: moves each entry to its new id; extPropNodes is only used to forward each node's own props
		for id, a := range m {
			if nn, ok := newNode(id); ok {
				nn.auxShard().put(uint32(nn), a)
				if a.ext != nil && a.ext.propListHead != nil {
					extPropNodes = append(extPropNodes, nn)
				}
			}
		}
	}
	for _, m := range oldCold {
		// order-independent: moves each entry to its new id
		for id, c := range m {
			if nn, ok := newNode(id); ok {
				nn.auxShard().putCold(uint32(nn), c)
			}
		}
	}
	k.slid = true
	k.Live = int(live)

	// Node ids held by the JSDoc of live nodes and by node-valued props.
	// order-independent: forwards each JSDocInfo's own type nodes; Forward is idempotent
	for info := range k.docs {
		info.forwardTypeNodes(k.Forward)
	}
	for _, n := range extPropNodes {
		for _, v := range nodeProps(n) {
			n.PutProp(v.prop, k.Forward(v.node))
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

// PruneDeadAux drops auxiliary entries for unmarked nodes without changing any
// live node IDs, when compaction aborts before Slide.
func (k *Compactor) PruneDeadAux() {
	s := k.s
	for sh := range s.aux {
		// order-independent: deletes dead entries
		for id := range s.aux[sh].m {
			if i := id - k.base; id >= k.base && id < k.end && !k.isMarked(i) {
				delete(s.aux[sh].m, id)
			}
		}
		// order-independent: deletes dead entries
		for id := range s.aux[sh].cold {
			if i := id - k.base; id >= k.base && id < k.end && !k.isMarked(i) {
				delete(s.aux[sh].cold, id)
			}
		}
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
