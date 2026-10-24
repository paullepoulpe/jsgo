package rhino

import (
	"sort"
	"unsafe"
)

// Go-only: SPIKE, not used by the compiler. Measures what flat slice storage
// (Node links as indexes into a []Node) would cost against the chunk store,
// on trees copied out of the chunk store. See
// docs/memory-research/round3-flat-slices.md.

// FlatStore holds nodes in one slice; a link is an index into Nodes, 0 is nil.
type FlatStore struct {
	Nodes []Node
}

// FlatOrder picks the order nodes are laid out in when copied.
type FlatOrder int

const (
	FlatAllocOrder FlatOrder = iota // the chunk store's allocation order (what parsing produces)
	FlatPreOrder                    // tree pre-order (what a rebuild after a pass would produce)
)

// FlattenInto appends the tree under root to s (links become indexes into
// s.Nodes) and returns the root's index.
func FlattenInto(s *FlatStore, root *Node, order FlatOrder) uint32 {
	if len(s.Nodes) == 0 {
		s.Nodes = append(s.Nodes, Node{}) // index 0 = nil
	}
	var all []*Node
	var collect func(n *Node)
	collect = func(n *Node) {
		all = append(all, n)
		for c := n.firstNode(); c != nil; c = c.nextNode() {
			collect(c)
		}
	}
	collect(root)
	if order == FlatAllocOrder {
		sort.Slice(all, func(i, j int) bool { return all[i].self < all[j].self })
	}
	base := uint32(len(s.Nodes))
	remap := make(map[nodeID]uint32, len(all))
	for i, n := range all {
		remap[n.self] = base + uint32(i)
	}
	m := func(id nodeID) nodeID {
		if id == 0 {
			return 0
		}
		return nodeID(remap[id]) // 0 when outside the tree (the root's parent)
	}
	for _, n := range all {
		c := *n
		c.next, c.first, c.parent, c.previous = m(n.next), m(n.first), m(n.parent), m(n.previous)
		c.self = nodeID(remap[n.self])
		s.Nodes = append(s.Nodes, c)
	}
	return remap[root.self]
}

// --- workloads over a FlatStore (ids are indexes) ---

// FlatWalk visits every node under id recursively (first child, next
// sibling), summing tokens, as NodeTraversal does.
func (s *FlatStore) FlatWalk(id uint32) uint64 {
	nodes := s.Nodes
	n := &nodes[id]
	sum := uint64(n.token)
	for c := uint32(n.first); c != 0; c = uint32(nodes[c].next) {
		sum += s.FlatWalk(c)
	}
	return sum
}

// FlatDepthSum walks from every node in [lo, hi) up to its root.
func (s *FlatStore) FlatDepthSum(lo, hi uint32) uint64 {
	nodes := s.Nodes
	var sum uint64
	for i := lo; i < hi; i++ {
		for p := uint32(nodes[i].parent); p != 0; p = uint32(nodes[p].parent) {
			sum++
		}
	}
	return sum
}

// --- two-level lookup: per-file slices addressed by a global id ---

// FileStores addresses per-file stores with id = file<<FileShift | local:
// what a pass needs when it follows a link into another file without the
// unified slice.
type FileStores struct {
	Files [][]Node
}

const FileShift = 20

func (f *FileStores) at(id uint32) *Node {
	return &f.Files[id>>FileShift][id&(1<<FileShift-1)]
}

// TwoLevelWalk is FlatWalk through the two-level lookup; links inside a file
// are local, so the walk carries the file's high bits.
func (f *FileStores) TwoLevelWalk(id uint32) uint64 {
	n := f.at(id)
	sum := uint64(n.token)
	hi := id &^ (1<<FileShift - 1)
	for c := uint32(n.first); c != 0; {
		g := hi | c
		sum += f.TwoLevelWalk(g)
		c = uint32(f.at(g).next)
	}
	return sum
}

// --- the chunk store, for comparison ---

// ChunkWalk is FlatWalk through the public *Node API (what passes use today).
func ChunkWalk(n *Node) uint64 {
	sum := uint64(n.token)
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		sum += ChunkWalk(c)
	}
	return sum
}

// ChunkDepthSum walks every node in nodes up to its root through GetParent.
func ChunkDepthSum(nodes []*Node) uint64 {
	var sum uint64
	for _, n := range nodes {
		for p := n.GetParent(); p != nil; p = p.GetParent() {
			sum++
		}
	}
	return sum
}

// nodeAtInl is nodeAt without the atomic, the stale check and the nil id
// check, small enough to inline: the chunk table at its cheapest.
func nodeAtInl(id nodeID) *Node {
	c := *(*unsafe.Pointer)(unsafe.Pointer(&nodeChunkTable[uint32(id)>>nodeChunkShift]))
	return (*Node)(unsafe.Add(c, uintptr(id&nodeChunkMask)*nodeSize))
}

// ChunkWalkInl is ChunkWalk with an inlined table lookup.
func ChunkWalkInl(n *Node) uint64 {
	sum := uint64(n.token)
	for c := n.first; c != 0; {
		cn := nodeAtInl(c)
		sum += ChunkWalkInl(cn)
		c = cn.next
	}
	return sum
}

// ChunkDepthSumInl is ChunkDepthSum with an inlined table lookup.
func ChunkDepthSumInl(nodes []*Node) uint64 {
	var sum uint64
	for _, n := range nodes {
		for p := n.parent; p != 0; p = nodeAtInl(p).parent {
			sum++
		}
	}
	return sum
}

// CollectNodes returns every node under root, in allocation order.
func CollectNodes(root *Node) []*Node {
	var all []*Node
	var collect func(n *Node)
	collect = func(n *Node) {
		all = append(all, n)
		for c := n.firstNode(); c != nil; c = c.nextNode() {
			collect(c)
		}
	}
	collect(root)
	sort.Slice(all, func(i, j int) bool { return all[i].self < all[j].self })
	return all
}

// MergeStores concatenates per-file stores into one, rebasing every link:
// the switch from per-file to unified storage. roots are the per-file root
// indexes; the result's roots are returned.
func MergeStores(stores []*FlatStore, roots []uint32) (*FlatStore, []uint32) {
	total := 1
	for _, s := range stores {
		total += len(s.Nodes) - 1
	}
	u := &FlatStore{Nodes: make([]Node, 1, total)}
	out := make([]uint32, len(stores))
	for i, s := range stores {
		base := nodeID(len(u.Nodes) - 1) // local index 1 lands at base+1
		rb := func(id nodeID) nodeID {
			if id == 0 {
				return 0
			}
			return id + base
		}
		for _, n := range s.Nodes[1:] {
			n.next, n.first, n.parent, n.previous, n.self = rb(n.next), rb(n.first), rb(n.parent), rb(n.previous), rb(n.self)
			u.Nodes = append(u.Nodes, n)
		}
		out[i] = roots[i] + uint32(base)
	}
	return u, out
}
