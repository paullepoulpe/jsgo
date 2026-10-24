package rhino

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: copying compaction of the AST (see node_store.go).
//
// A store keeps every node allocated in it, dead or alive. A Compactor
// copies the live tree in pre-order into a fresh store (so the copy is also
// laid out in walk order) and makes that store current; then every Node
// held outside the tree is forwarded with Forward (a holder of a detached
// subtree gets that subtree copied as an extra root). After Finish nothing
// refers to the old store and the GC frees it.
//
// The new store's ids start after the old store's, so Forward knows which
// ids are already copies, and an id that was not forwarded is outside the
// current store and panics on its first use.
//
// An old node remembers its copy: copyNode sets kindFlagForwarded and puts
// the copy's id in old.next; Forward reads that, so it is idempotent and
// holders can be forwarded in any order.

// Compactor copies a live AST into a fresh store.
type Compactor struct {
	src, dst   *NodeStore
	docs       map[*JSDocInfo]struct{}
	Copied     int // nodes copied
	ExtraRoots int // detached subtrees copied for holders
	ExtraNodes int // nodes in those subtrees
	inExtra    bool
}

// NewCompactor creates a new store after the current one, makes it current
// (new nodes are allocated in it from now on) and returns the compactor
// that fills it. The old store stays readable until Finish.
func NewCompactor() *Compactor {
	src := current
	if src == nil {
		panic("rhino.Compactor: no current node store")
	}
	dst := NewNodeStore()
	compactSource = src
	SetCurrentNodeStore(dst)
	return &Compactor{
		src:  src,
		dst:  dst,
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
	if uint32(old)&staticBit != 0 {
		r := old
		for r.rec().parent != 0 {
			r = r.parentNode()
		}
		name := ""
		if p := r.GetSourceFileName(); p != nil {
			name = *p
		}
		panic(fmt.Sprintf("rhino.Compactor: a package-level node or template is in the compacted tree: %s (root %s, source %s, extra root: %v)",
			old, r, name, k.inExtra))
	}
	or := old.rec()
	nw := k.dst.alloc()
	nr := nw.rec()
	*nr = *or
	nr.next, nr.first, nr.parent, nr.previous = 0, 0, 0, 0
	nr.tc = k.dst.translateTC(k.src, or.tc, or.kindFlags&kindFlagTypeIsColor != 0)
	nr.aux = 0
	if a := old.auxValue(); a != (nodeAux{}) {
		if a.ext != nil && a.ext.typeBeforeCast != 0 {
			a.ext.typeBeforeCast = k.dst.translateTypeHandle(k.src, a.ext.typeBeforeCast)
		}
		nw.setAux(a)
	}
	k.Copied++
	if (nr.propFlags & propFlagHasJSDoc) != 0 {
		if info := nw.GetJSDocInfo(); info != nil {
			if _, done := k.docs[info]; !done {
				k.docs[info] = struct{}{}
				info.forwardTypeNodes(k.Forward)
			}
		}
	}
	if (nr.kindFlags & kindFlagHasExtBox) != 0 {
		k.forwardNodeProps(nw)
	}
	or.kindFlags |= kindFlagForwarded
	or.next = uint32(nw)
	or.first, or.parent, or.previous = 0, 0, 0
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
	child := Node(old.rec().first)
	nw := k.copyNode(old)
	var last Node
	for child != 0 {
		c := child
		child = Node(c.rec().next) // read before the copy overwrites it
		cc := k.copyTree(c)
		cc.rec().parent = uint32(nw)
		if last == 0 {
			nw.rec().first = uint32(cc)
		} else {
			last.rec().next = uint32(cc)
			cc.rec().previous = uint32(last)
		}
		last = cc
	}
	if last != 0 {
		nw.firstNode().rec().previous = uint32(last)
	}
	return nw
}

// copied reports whether n needs no copy (nil, already in the new store,
// or a package-level node) and returns its copy if it was copied.
func (k *Compactor) copied(n Node) (Node, bool) {
	id := uint32(n)
	if n == 0 || id&staticBit != 0 || k.dst.contains(id) {
		return n, true
	}
	if !k.src.contains(id) {
		panic(fmt.Sprintf("rhino.Compactor: node id %d is from an older store", id))
	}
	if r := n.rec(); r.kindFlags&kindFlagForwarded != 0 {
		return Node(r.next), true
	}
	return 0, false
}

// CopyTree copies the tree rooted at root (which must be detached) and
// returns the copy.
func (k *Compactor) CopyTree(root Node) Node {
	if c, ok := k.copied(root); ok {
		return c
	}
	if root.rec().parent != 0 {
		panic("rhino.Compactor: CopyTree of an attached node")
	}
	return k.copyTree(root)
}

// Forward returns the copy of n: n itself when it is already in the new
// store (or a package-level node), its copy when it was copied, and
// otherwise copies the detached subtree n belongs to as an extra root.
func (k *Compactor) Forward(n Node) Node {
	if c, ok := k.copied(n); ok {
		return c
	}
	r := n
	for r.rec().parent != 0 {
		r = r.parentNode()
		if r.rec().kindFlags&kindFlagForwarded != 0 {
			panic("rhino.Compactor: an ancestor was forwarded but the node was not copied")
		}
	}
	before := k.Copied
	k.inExtra = true
	k.copyTree(r)
	k.inExtra = false
	k.ExtraRoots++
	k.ExtraNodes += k.Copied - before
	c, _ := k.copied(n)
	return c
}

// ForwardSlice forwards every element of s in place.
func (k *Compactor) ForwardSlice(s []Node) {
	for i, n := range s {
		s[i] = k.Forward(n)
	}
}

// Finish ends the compaction: the old store leaves the block directory (its
// ids fail like nil pointers from now on) and is freed once the GC sees
// nothing refers to it.
func (k *Compactor) Finish() {
	k.src.unregister()
	compactSource = nil
	k.src = nil
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
