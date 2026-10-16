package jutil

import (
	"iter"
	"slices"
)

// Go-only: language-adapter — java.util.TreeSet with a comparator. The
// Javadoc of java.util.TreeMap (which backs TreeSet) specifies a red-black
// tree using the algorithms of Cormen, Leiserson and Rivest's "Introduction
// to Algorithms"; this is that textbook tree. Iteration is in comparator
// order. The tree's shape is observable only through a comparator that is
// inconsistent over time (testdata/jscompa_treeset.tsv has such vectors), so
// the descent compares compare(probe, nodeKey) at each node, insertion uses
// the textbook fix-up, and deletion of a node with two children puts its
// in-order successor in its place.
//
// Observed behaviour (black-box): adding to an empty set calls
// compare(e, e) once (TreeMap's type check), lookups in an empty set call
// no comparator.

type treeSetNode[E any] struct {
	key                 E
	parent, left, right *treeSetNode[E]
	red                 bool
}

// TreeSet is java.util.TreeSet<E> with a comparator.
type TreeSet[E any] struct {
	compare func(a, b E) int
	size    int
	root    *treeSetNode[E]
}

// NewTreeSet is new TreeSet<>(comparator).
func NewTreeSet[E any](compare func(a, b E) int) *TreeSet[E] {
	return &TreeSet[E]{compare: compare}
}

// Size is Set.size().
func (s *TreeSet[E]) Size() int {
	if s == nil {
		return 0
	}
	return s.size
}

// IsEmpty is Set.isEmpty().
func (s *TreeSet[E]) IsEmpty() bool { return s.Size() == 0 }

// Clear is Set.clear().
func (s *TreeSet[E]) Clear() {
	s.root = nil
	s.size = 0
}

// find is the node whose key compares equal to e, or nil.
func (s *TreeSet[E]) find(e E) *treeSetNode[E] {
	if s == nil {
		return nil
	}
	for p := s.root; p != nil; {
		switch c := s.compare(e, p.key); {
		case c < 0:
			p = p.left
		case c > 0:
			p = p.right
		default:
			return p
		}
	}
	return nil
}

// Contains is Set.contains(Object): whether some element compares equal.
func (s *TreeSet[E]) Contains(e E) bool { return s.find(e) != nil }

// Add is Set.add: inserts key unless an element compares equal to it (which
// is then kept); reports whether the set changed.
func (s *TreeSet[E]) Add(key E) bool {
	if s.root == nil {
		s.compare(key, key)
		s.root = &treeSetNode[E]{key: key}
		s.size = 1
		return true
	}
	p := s.root
	for {
		c := s.compare(key, p.key)
		var next **treeSetNode[E]
		switch {
		case c < 0:
			next = &p.left
		case c > 0:
			next = &p.right
		default:
			return false
		}
		if *next == nil {
			x := &treeSetNode[E]{key: key, parent: p, red: true}
			*next = x
			s.insertFixup(x)
			s.size++
			return true
		}
		p = *next
	}
}

// Remove is Set.remove(Object): removes the element comparing equal to key.
func (s *TreeSet[E]) Remove(key E) bool {
	z := s.find(key)
	if z == nil {
		return false
	}
	s.delete(z)
	s.size--
	return true
}

// First is SortedSet.first(): the lowest element; NoSuchElementException
// when empty.
func (s *TreeSet[E]) First() E {
	if s.IsEmpty() {
		panic(NewNoSuchElementException())
	}
	return treeSetMin(s.root).key
}

// Last is SortedSet.last(): the highest element; NoSuchElementException
// when empty.
func (s *TreeSet[E]) Last() E {
	if s.IsEmpty() {
		panic(NewNoSuchElementException())
	}
	return treeSetMax(s.root).key
}

// All iterates in ascending order.
func (s *TreeSet[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		if s == nil || s.root == nil {
			return
		}
		for x := treeSetMin(s.root); x != nil; x = treeSetSuccessor(x) {
			if !yield(x.key) {
				return
			}
		}
	}
}

// Descending iterates in descending order (descendingIterator()).
func (s *TreeSet[E]) Descending() iter.Seq[E] {
	return func(yield func(E) bool) {
		if s == nil || s.root == nil {
			return
		}
		for x := treeSetMax(s.root); x != nil; x = treeSetPredecessor(x) {
			if !yield(x.key) {
				return
			}
		}
	}
}

// ToSlice returns the elements in ascending order.
func (s *TreeSet[E]) ToSlice() []E { return slices.Collect(s.All()) }

// String is AbstractCollection.toString(): "[a, b]".
func (s *TreeSet[E]) String() string { return javaCollectionString(any(s), s.All()) }

// AddAll is Collection.addAll: adds each element in order; reports whether
// the set changed.
// Go-only: language-adapter — Java's TreeSet.addAll bulk-builds an empty set
// from a SortedSet with the same comparator; a sequence carries no
// comparator, so the elements are always added one by one (same contents
// and order; only the shape of the tree can differ).
func (s *TreeSet[E]) AddAll(es iter.Seq[E]) bool {
	changed := false
	for e := range es {
		if s.Add(e) {
			changed = true
		}
	}
	return changed
}

func treeSetMin[E any](x *treeSetNode[E]) *treeSetNode[E] {
	for x.left != nil {
		x = x.left
	}
	return x
}

func treeSetMax[E any](x *treeSetNode[E]) *treeSetNode[E] {
	for x.right != nil {
		x = x.right
	}
	return x
}

func treeSetSuccessor[E any](x *treeSetNode[E]) *treeSetNode[E] {
	if x.right != nil {
		return treeSetMin(x.right)
	}
	p := x.parent
	for p != nil && x == p.right {
		x, p = p, p.parent
	}
	return p
}

func treeSetPredecessor[E any](x *treeSetNode[E]) *treeSetNode[E] {
	if x.left != nil {
		return treeSetMax(x.left)
	}
	p := x.parent
	for p != nil && x == p.left {
		x, p = p, p.parent
	}
	return p
}

func treeSetIsRed[E any](x *treeSetNode[E]) bool { return x != nil && x.red }

func (s *TreeSet[E]) rotateLeft(x *treeSetNode[E]) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	s.replaceChild(x, y)
	y.left = x
	x.parent = y
}

func (s *TreeSet[E]) rotateRight(x *treeSetNode[E]) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	s.replaceChild(x, y)
	y.right = x
	x.parent = y
}

// replaceChild links v (possibly nil) where u hangs from u's parent.
func (s *TreeSet[E]) replaceChild(u, v *treeSetNode[E]) {
	switch up := u.parent; {
	case up == nil:
		s.root = v
	case up.left == u:
		up.left = v
	default:
		up.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

// insertFixup restores the red-black properties after red leaf z was
// attached (textbook cases: red uncle recolours; otherwise one or two
// rotations).
func (s *TreeSet[E]) insertFixup(z *treeSetNode[E]) {
	for treeSetIsRed(z.parent) {
		zp := z.parent
		zpp := zp.parent
		if zp == zpp.left {
			if y := zpp.right; treeSetIsRed(y) {
				zp.red, y.red, zpp.red = false, false, true
				z = zpp
				continue
			}
			if z == zp.right {
				z = zp
				s.rotateLeft(z)
				zp = z.parent
			}
			zp.red, zpp.red = false, true
			s.rotateRight(zpp)
		} else {
			if y := zpp.left; treeSetIsRed(y) {
				zp.red, y.red, zpp.red = false, false, true
				z = zpp
				continue
			}
			if z == zp.left {
				z = zp
				s.rotateRight(z)
				zp = z.parent
			}
			zp.red, zpp.red = false, true
			s.rotateLeft(zpp)
		}
	}
	s.root.red = false
}

// delete unlinks node z; a node with two children is replaced by its
// in-order successor.
func (s *TreeSet[E]) delete(z *treeSetNode[E]) {
	removedRed := z.red
	var x, xp *treeSetNode[E]
	switch {
	case z.left == nil:
		x, xp = z.right, z.parent
		s.replaceChild(z, x)
	case z.right == nil:
		x, xp = z.left, z.parent
		s.replaceChild(z, x)
	default:
		y := treeSetMin(z.right)
		removedRed = y.red
		x = y.right
		if y.parent == z {
			xp = y
		} else {
			xp = y.parent
			s.replaceChild(y, x)
			y.right = z.right
			y.right.parent = y
		}
		s.replaceChild(z, y)
		y.left = z.left
		y.left.parent = y
		y.red = z.red
	}
	if !removedRed {
		s.deleteFixup(x, xp)
	}
}

// deleteFixup restores the red-black properties after a black node was
// removed; x (possibly nil, with parent xp) carries the extra black.
func (s *TreeSet[E]) deleteFixup(x, xp *treeSetNode[E]) {
	for x != s.root && !treeSetIsRed(x) {
		if x == xp.left {
			w := xp.right
			if treeSetIsRed(w) {
				w.red, xp.red = false, true
				s.rotateLeft(xp)
				w = xp.right
			}
			if !treeSetIsRed(w.left) && !treeSetIsRed(w.right) {
				w.red = true
				x, xp = xp, xp.parent
				continue
			}
			if !treeSetIsRed(w.right) {
				w.left.red, w.red = false, true
				s.rotateRight(w)
				w = xp.right
			}
			w.red, xp.red = xp.red, false
			w.right.red = false
			s.rotateLeft(xp)
			x = s.root
		} else {
			w := xp.left
			if treeSetIsRed(w) {
				w.red, xp.red = false, true
				s.rotateRight(xp)
				w = xp.left
			}
			if !treeSetIsRed(w.left) && !treeSetIsRed(w.right) {
				w.red = true
				x, xp = xp, xp.parent
				continue
			}
			if !treeSetIsRed(w.left) {
				w.right.red, w.red = false, true
				s.rotateLeft(w)
				w = xp.left
			}
			w.red, xp.red = xp.red, false
			w.left.red = false
			s.rotateRight(xp)
			x = s.root
		}
	}
	if x != nil {
		x.red = false
	}
}
