package rhino

import (
	"fmt"
	"os"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: copying compaction of the AST (see node_store.go).
//
// Links are ids, so the garbage collector never frees a dead node: a store
// keeps every node allocated in it. A Compactor copies the live tree in
// pre-order into a fresh store (so the copy is also laid out in walk order)
// and makes that store current; then every Node held outside the tree is
// forwarded with Forward (a holder of a detached subtree gets that subtree
// copied as an extra root). Once nothing refers to the old store any more,
// the GC frees it with all its dead nodes.
//
// An old node remembers its copy: copyNode sets kindFlagForwarded and puts
// the copy's id in old.next (the copy is in old.store.forwardedTo); Forward
// reads that, so it is idempotent and holders can be forwarded in any
// order.
//
// Poison mode (JSGO_AST_POISON=1, for tests): Finish overwrites every copied
// node with links and handles that can never resolve, so a holder that was
// not forwarded fails on its first use instead of reading a dead node.

// NodePoisonMode reports whether compaction poisons the old nodes.
var NodePoisonMode = os.Getenv("JSGO_AST_POISON") == "1"

// Compactor copies a live AST into a fresh store.
type Compactor struct {
	dst        *NodeStore
	srcs       map[*NodeStore]struct{}
	docs       map[*JSDocInfo]struct{}
	Copied     int    // nodes copied
	ExtraRoots int    // detached subtrees copied for holders
	ExtraNodes int    // nodes in those subtrees
	inExtra    bool
}

// NewCompactor creates a new store, makes it current (new nodes are
// allocated in it from now on) and returns the compactor that fills it.
func NewCompactor() *Compactor {
	dst := NewNodeStore()
	SetCurrentNodeStore(dst)
	return &Compactor{
		dst:  dst,
		srcs: make(map[*NodeStore]struct{}, 2),
		docs: make(map[*JSDocInfo]struct{}, 1024),
	}
}

// Store returns the store the compactor copies into.
func (k *Compactor) Store() *NodeStore {
	return k.dst
}

// copyNode copies old into the new store (links cleared), forwards the side
// trees it owns, and marks old as forwarded.
func (k *Compactor) copyNode(old Node) Node {
	src := old.store
	if _, seen := k.srcs[src]; !seen {
		if src.static {
			r := old
			for r.rec.parent != 0 {
				r = r.parentNode()
			}
			name := ""
			if p := r.GetSourceFileName(); p != nil {
				name = *p
			}
			panic(fmt.Sprintf("rhino.Compactor: a package-level node or template is in the compacted tree: %s (root %s, source %s, extra root: %v)",
				old, r, name, k.inExtra))
		}
		if src.forwardedTo != nil && src.forwardedTo != k.dst {
			panic("rhino.Compactor: a node of an already compacted store was not forwarded by that compaction")
		}
		src.forwardedTo = k.dst
		k.srcs[src] = struct{}{}
	}
	nw := k.dst.alloc()
	self := nw.rec.self
	*nw.rec = *old.rec
	nw.rec.self = self
	nw.rec.next, nw.rec.first, nw.rec.parent, nw.rec.previous = 0, 0, 0, 0
	nw.rec.tc = k.dst.translateTC(src, old.rec.tc, old.rec.kindFlags&kindFlagTypeIsColor != 0)
	nw.rec.aux = 0
	if a := old.auxValue(); a != (nodeAux{}) {
		if a.ext != nil && a.ext.typeBeforeCast != 0 {
			a.ext.typeBeforeCast = k.dst.translateTypeHandle(src, a.ext.typeBeforeCast)
		}
		nw.setAux(a)
	}
	k.Copied++
	if (nw.rec.propFlags & propFlagHasJSDoc) != 0 {
		if info := nw.GetJSDocInfo(); info != nil {
			if _, done := k.docs[info]; !done {
				k.docs[info] = struct{}{}
				info.forwardTypeNodes(k.Forward)
			}
		}
	}
	if (nw.rec.kindFlags & kindFlagHasExtBox) != 0 {
		k.forwardNodeProps(nw)
	}
	old.rec.kindFlags |= kindFlagForwarded
	old.rec.next = nw.rec.self
	old.rec.first = 0
	old.rec.parent = 0
	old.rec.previous = 0
	return nw
}

// forwardNodeProps forwards the subtrees hung off object props
// (DECLARED_TYPE_EXPR, CLOSURE_UNAWARE_SHADOW).
func (k *Compactor) forwardNodeProps(n Node) {
	var props [4]Node_Prop
	var vals [4]Node
	cnt := 0
	for x := n.coldPropListHead(); x != nil; x = x.next {
		if x.isObject {
			if v, ok := x.objectValue.(Node); ok && !v.IsNil() {
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
func (k *Compactor) copyTree(old Node) Node {
	src := old.store
	child := old.rec.first
	nw := k.copyNode(old)
	var last Node
	for child != 0 {
		c := src.node(child)
		child = c.rec.next
		cc := k.copyTree(c)
		cc.rec.parent = nw.rec.self
		if last.IsNil() {
			nw.rec.first = cc.rec.self
		} else {
			last.rec.next = cc.rec.self
			cc.rec.previous = last.rec.self
		}
		last = cc
	}
	if !last.IsNil() {
		nw.firstNode().rec.previous = last.rec.self
	}
	return nw
}

// forwarded returns the copy of a forwarded node (following the copies of
// earlier compactions too).
func forwarded(n Node) Node {
	for n.rec.kindFlags&kindFlagForwarded != 0 {
		n = n.store.forwardedTo.node(n.rec.next)
	}
	return n
}

// CopyTree copies the tree rooted at root (which must be detached) and
// returns the copy.
func (k *Compactor) CopyTree(root Node) Node {
	if root.IsNil() {
		return Node{}
	}
	if root.rec.kindFlags&kindFlagForwarded != 0 {
		return forwarded(root)
	}
	if root.store == k.dst || root.store.static {
		return root
	}
	if root.rec.parent != 0 {
		panic("rhino.Compactor: CopyTree of an attached node")
	}
	return k.copyTree(root)
}

// Forward returns the copy of n: n itself when it is already in the new
// store (or a package-level node), its copy when it was copied, and
// otherwise copies the detached subtree n belongs to as an extra root.
func (k *Compactor) Forward(n Node) Node {
	if n.IsNil() {
		return n
	}
	if n.rec.kindFlags&kindFlagForwarded != 0 {
		return forwarded(n)
	}
	if n.store == k.dst || n.store.static {
		return n
	}
	r := n
	for r.rec.parent != 0 {
		r = r.parentNode()
		if r.rec.kindFlags&kindFlagForwarded != 0 {
			panic("rhino.Compactor: an ancestor was forwarded but the node was not copied")
		}
	}
	before := k.Copied
	k.inExtra = true
	k.copyTree(r)
	k.inExtra = false
	k.ExtraRoots++
	k.ExtraNodes += k.Copied - before
	return forwarded(n)
}

// ForwardSlice forwards every element of s in place.
func (k *Compactor) ForwardSlice(s []Node) {
	for i, n := range s {
		s[i] = k.Forward(n)
	}
}

// poisonID is an id no store ever hands out (beyond maxSegments).
const poisonID nodeID = 0xFFFFFFF0

// Finish ends the compaction: the old stores release their segments and
// tables, so a holder that was not forwarded pins at most the segment its
// node is in, and following one of its links panics. In poison mode every
// copied node is also overwritten so that such a holder fails on its first
// use.
func (k *Compactor) Finish() {
	if NodePoisonMode {
		for s := range k.srcs {
			for i := range s.segs {
				seg := s.segs[i].Load()
				if seg == nil {
					continue
				}
				for j := range seg {
					if r := &seg[j]; r.kindFlags&kindFlagForwarded != 0 {
						r.next, r.first, r.parent, r.previous, r.self = poisonID, poisonID, poisonID, poisonID, poisonID
						r.tc, r.aux = 0xFFFFFFF0, 0xFFFFFFF0
						r.token = Token(0xFF)
						r.kindFlags &^= kindFlagForwarded
					}
				}
			}
		}
	}
	for s := range k.srcs {
		s.release()
	}
	k.srcs = nil
	k.docs = nil
}

// String describes the compaction (for logs).
func (k *Compactor) String() string {
	return fmt.Sprintf("copied %d nodes (%d extra roots, %d nodes)", k.Copied, k.ExtraRoots, k.ExtraNodes)
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
