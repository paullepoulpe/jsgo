package rhino

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: copying compaction of the node store (see node_store.go).
//
// Links are ids, so the garbage collector never frees a dead node: a chunk
// lives as long as its generation. A Compactor copies the live tree in
// pre-order into a fresh generation (so the copy is also laid out in walk
// order), then every *Node held outside the tree is forwarded with Forward
// (a holder of a detached subtree gets that subtree copied as an extra
// root), and Finish releases the generations the live nodes came from.
//
// An old node remembers its copy: after copyNode, old.next is the id of the
// copy and old.previous is nodeForwardMarker; Forward reads that, so it is
// idempotent and holders can be forwarded in any order.
//
// Poison mode (JSGO_AST_POISON=1, for tests): the old chunks are kept alive
// forever, their table slots are cleared (any id into them panics in nodeAt
// with "stale node id") and every old node is overwritten with links and an
// id that can never resolve, so a holder that was not forwarded fails on
// its first use instead of reading a dead node.

// nodeForwardMarker is an id that is never allocated (the last chunk index
// is reserved, see addChunk).
const nodeForwardMarker nodeID = 0xFFFFFFFF

// NodePoisonMode reports whether compaction poisons the old nodes.
var NodePoisonMode = os.Getenv("JSGO_AST_POISON") == "1"

// Compactor copies a live AST into a fresh generation.
type Compactor struct {
	gen        *nodeGen
	old        map[*nodeGen]struct{}
	docs       map[*JSDocInfo]struct{}
	lastChunk  *nodeChunk
	Copied     int // nodes copied
	ExtraRoots int // detached subtrees copied for holders
	ExtraNodes int // nodes in those subtrees
	inExtra    bool
}

// NewCompactor starts a new generation (new nodes are allocated in it from
// now on) and returns the compactor that fills it.
func NewCompactor() *Compactor {
	return &Compactor{
		gen:  NewNodeGeneration().g,
		old:  make(map[*nodeGen]struct{}, 4),
		docs: make(map[*JSDocInfo]struct{}, 1024),
	}
}

func (k *Compactor) genOf(n *Node) *nodeGen {
	c := k.lastChunk
	if c == nil || c.hdr.index != uint32(n.self)>>nodeChunkShift {
		c = nodeChunkSlot(uint32(n.self) >> nodeChunkShift).Load()
		if c == nil {
			staleNodeID(n.self)
		}
		k.lastChunk = c
	}
	return c.hdr.owner
}

// copyNode copies n into the new generation (links cleared), forwards the
// side trees it owns, and marks n as forwarded.
func (k *Compactor) copyNode(old *Node) *Node {
	g := k.genOf(old)
	if _, seen := k.old[g]; !seen {
		if g == nodeInitialGen && NodePoisonMode {
			r := old
			for r.parent != 0 {
				r = nodeAt(r.parent)
			}
			panic(fmt.Sprintf("rhino.Compactor: a node of the initial generation (a package-level sentinel) is in the compacted tree: %s (root %s, source %s, extra root: %v)",
				old, r, r.GetSourceFileName(), k.inExtra))
		}
		k.old[g] = struct{}{}
	}
	nw := allocNode()
	self := nw.self
	*nw = *old
	nw.self = self
	nw.next, nw.first, nw.parent, nw.previous = 0, 0, 0, 0
	k.Copied++
	if (nw.propFlags & propFlagHasJSDoc) != 0 {
		if info := nw.GetJSDocInfo(); info != nil {
			if _, done := k.docs[info]; !done {
				k.docs[info] = struct{}{}
				info.forwardTypeNodes(k.Forward)
			}
		}
	}
	if (nw.kindFlags & kindFlagHasExtBox) != 0 {
		k.forwardNodeProps(nw)
	}
	old.next = nw.self
	old.previous = nodeForwardMarker
	old.first = 0
	old.parent = 0
	return nw
}

// forwardNodeProps forwards the subtrees hung off object props
// (DECLARED_TYPE_EXPR, CLOSURE_UNAWARE_SHADOW).
func (k *Compactor) forwardNodeProps(n *Node) {
	var props [4]Node_Prop
	var vals [4]*Node
	cnt := 0
	for x := n.coldPropListHead(); x != nil; x = x.next {
		if x.isObject {
			if v, ok := x.objectValue.(*Node); ok && v != nil {
				if cnt == len(props) {
					panic("rhino.Compactor: too many node-valued props")
				}
				props[cnt] = Node_Prop(x.propType)
				vals[cnt] = v
				cnt++
			}
		}
	}
	for i := 0; i < cnt; i++ {
		n.PutProp(props[i], k.Forward(vals[i]))
	}
}

// copyTree copies the subtree rooted at old in pre-order.
func (k *Compactor) copyTree(old *Node) *Node {
	child := old.first
	nw := k.copyNode(old)
	var last *Node
	for child != 0 {
		c := nodeAt(child)
		child = c.next
		cc := k.copyTree(c)
		cc.parent = nw.self
		if last == nil {
			nw.first = cc.self
		} else {
			last.next = cc.self
			cc.previous = last.self
		}
		last = cc
	}
	if last != nil {
		nodeAt(nw.first).previous = last.self
	}
	return nw
}

// CopyTree copies the tree rooted at root (which must be detached) and
// returns the copy.
func (k *Compactor) CopyTree(root *Node) *Node {
	if root == nil {
		return nil
	}
	if root.previous == nodeForwardMarker {
		return nodeAt(root.next)
	}
	if root.self == 0 || k.genOf(root) == k.gen {
		return root
	}
	if root.parent != 0 {
		panic("rhino.Compactor: CopyTree of an attached node")
	}
	return k.copyTree(root)
}

// Forward returns the copy of n: n itself when it is already in the new
// generation (or not in any chunk), its copy when it was copied, and
// otherwise copies the detached subtree n belongs to as an extra root.
func (k *Compactor) Forward(n *Node) *Node {
	if n == nil || n.self == 0 {
		return n
	}
	if n.previous == nodeForwardMarker {
		return nodeAt(n.next)
	}
	if g := k.genOf(n); g == k.gen || g == nodeInitialGen {
		// Already copied, or a package-level sentinel (the initial
		// generation is never released, so holders may keep it).
		return n
	}
	r := n
	for r.parent != 0 {
		r = nodeAt(r.parent)
		if r.previous == nodeForwardMarker {
			panic("rhino.Compactor: an ancestor was forwarded but the node was not copied")
		}
	}
	before := k.Copied
	k.inExtra = true
	k.copyTree(r)
	k.inExtra = false
	k.ExtraRoots++
	k.ExtraNodes += k.Copied - before
	return nodeAt(n.next)
}

// ForwardSlice forwards every element of s in place.
func (k *Compactor) ForwardSlice(s []*Node) {
	for i, n := range s {
		s[i] = k.Forward(n)
	}
}

var (
	nodePoisonMu     sync.Mutex
	nodePoisonChunks []*nodeChunk // kept alive forever in poison mode
)

// Finish releases the generations the copied nodes came from and returns
// the number of chunks released.
func (k *Compactor) Finish() int {
	released := 0
	for g := range k.old {
		if g == k.gen || g == nodeInitialGen {
			continue
		}
		g.mu.Lock()
		chunks := g.chunks
		g.chunks = nil
		g.mu.Unlock()
		for _, c := range chunks {
			nodeChunkSlot(c.hdr.index).Store(nil)
			released++
			if NodePoisonMode {
				poisonChunk(c)
			}
		}
		if NodePoisonMode {
			nodePoisonMu.Lock()
			nodePoisonChunks = append(nodePoisonChunks, chunks...)
			nodePoisonMu.Unlock()
		}
	}
	k.old = nil
	k.docs = nil
	k.lastChunk = nil
	return released
}

func poisonChunk(c *nodeChunk) {
	for i := 1; i < nodeChunkSlots; i++ {
		n := (*Node)(unsafe.Add(unsafe.Pointer(c), uintptr(i)*nodeSize))
		n.next, n.first, n.parent, n.previous, n.self = nodeForwardMarker, nodeForwardMarker, nodeForwardMarker, nodeForwardMarker, nodeForwardMarker
		n.token = Token(0xFF)
	}
}

// String describes the compaction (for logs).
func (k *Compactor) String() string {
	return fmt.Sprintf("copied %d nodes (%d extra roots, %d nodes)", k.Copied, k.ExtraRoots, k.ExtraNodes)
}

// forwardTypeNodes forwards the roots of the type-expression trees and type
// transformations the JSDocInfo holds.
func (i *JSDocInfo) forwardTypeNodes(fwd func(*Node) *Node) {
	for it := i.TypeExprIter(); it.Next(); {
		e := it.Expr()
		if e.root != nil {
			e.root = fwd(e.root)
		}
	}
	if m := jsdocGet[*jutil.LinkedHashMap[string, *Node]](JSDocInfo_TYPE_TRANSFORMATIONS.Get(i)); m != nil {
		for _, key := range m.KeySlice() {
			m.Put(key, fwd(m.Get(key)))
		}
	}
}
