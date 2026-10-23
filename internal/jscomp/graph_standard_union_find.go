// Port of com/google/javascript/jscomp/graph/StandardUnionFind.java @ closure-compiler v20261015

/*
 * Copyright 2008 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// StandardUnionFind is a Union-Find implementation.
//
// This class implements Union-Find algorithm with rank and path
// compression.
//
// See https://algorithmist.com/wiki/Union_find for more detail.
//
// E is the element type. Go-only: language-adapter — elements are keyed
// with Go == (Java equals() for strings, numbers and identity-compared
// objects).
type StandardUnionFind[E comparable] struct {
	// All values with the same root node are in the same equivalence set.
	elmap    *jutil.LinkedHashMap[E, *standardUnionFind_Node[E]]
	nodeSlab jutil.Slab[standardUnionFind_Node[E]]
}

// NewStandardUnionFind creates an empty UnionFind structure.
//
// Java: StandardUnionFind()
func NewStandardUnionFind[E comparable]() *StandardUnionFind[E] {
	return &StandardUnionFind[E]{
		elmap: jutil.NewLinkedHashMap[E, *standardUnionFind_Node[E]](),
	}
}

// NewStandardUnionFindCopyOf creates an UnionFind structure being a copy of
// other structure. The created structure is optimal in a sense that the
// paths to the root from any element will have a length of at most 1.
//
// Java: StandardUnionFind(UnionFind<E>)
func NewStandardUnionFindCopyOf[E comparable](other UnionFind[E]) *StandardUnionFind[E] {
	u := NewStandardUnionFind[E]()
	for elem := range other.Elements().All() {
		u.Union(other.Find(elem), elem)
	}
	return u
}

// Java: StandardUnionFind.add(E)
func (u *StandardUnionFind[E]) Add(e E) {
	u.Union(e, e)
}

// Java: StandardUnionFind.addAll(Iterable<E>)
func (u *StandardUnionFind[E]) AddAll(es iter.Seq[E]) {
	for e := range es {
		u.Add(e)
	}
}

// Java: StandardUnionFind.union(E, E)
func (u *StandardUnionFind[E]) Union(a E, b E) E {
	nodeA := u.findRootOrCreateNode(a)
	nodeB := u.findRootOrCreateNode(b)

	if nodeA == nodeB {
		return nodeA.element
	}
	// If possible, prefer nodeA over nodeB, to preserve insertion order.
	if nodeA.rank >= nodeB.rank {
		nodeB.parent = nodeA
		nodeA.size += nodeB.size
		if nodeA.rank == nodeB.rank {
			nodeA.rank++
		}
		return nodeA.element
	}
	nodeA.parent = nodeB
	nodeB.size += nodeA.size
	temp := nodeB.element
	nodeB.element = nodeA.element
	nodeA.element = temp
	return nodeB.element
}

// Java: StandardUnionFind.find(E)
func (u *StandardUnionFind[E]) Find(e E) E {
	jutil.CheckArgumentMsg1(u.elmap.ContainsKey(e), "Element does not exist: %s", e)
	return u.findRoot(u.elmap.Get(e)).element
}

// Java: StandardUnionFind.areEquivalent(E, E)
func (u *StandardUnionFind[E]) AreEquivalent(a E, b E) bool {
	aRep := u.Find(a)
	bRep := u.Find(b)
	return aRep == bRep
}

// Java: StandardUnionFind.elements()
func (u *StandardUnionFind[E]) Elements() jutil.Set[E] {
	return &standardUnionFind_keySet[E]{m: u.elmap}
}

// standardUnionFind_keySet is Collections.unmodifiableSet(elmap.keySet()).
type standardUnionFind_keySet[E comparable] struct {
	m *jutil.LinkedHashMap[E, *standardUnionFind_Node[E]]
}

func (s *standardUnionFind_keySet[E]) Size() int         { return s.m.Size() }
func (s *standardUnionFind_keySet[E]) IsEmpty() bool     { return s.m.IsEmpty() }
func (s *standardUnionFind_keySet[E]) Contains(e E) bool { return s.m.ContainsKey(e) }
func (s *standardUnionFind_keySet[E]) All() iter.Seq[E]  { return s.m.Keys() }
func (s *standardUnionFind_keySet[E]) String() string {
	return jutil.CollectionToString(s.m.KeySlice())
}

// Java: StandardUnionFind.allEquivalenceClasses()
func (u *StandardUnionFind[E]) AllEquivalenceClasses() []*jutil.LinkedHashSet[E] {
	// MultimapBuilder.linkedHashKeys().linkedHashSetValues(); the keys are
	// Nodes (identity).
	groupsTmp := jutil.NewLinkedHashMap[*standardUnionFind_Node[E], *jutil.LinkedHashSet[E]]()
	for elem := range u.elmap.Values() {
		root := u.findRoot(elem)
		group := groupsTmp.Get(root)
		if group == nil {
			group = jutil.NewLinkedHashSet[E]()
			groupsTmp.Put(root, group)
		}
		group.Add(elem.element)
	}
	result := []*jutil.LinkedHashSet[E]{}
	for group := range groupsTmp.Values() {
		result = append(result, jutil.NewLinkedHashSet(group.ToSlice()...))
	}
	return result
}

// AllRepresentatives returns the reprsentative elements of all the
// equivalence classes.
//
// This is a "snapshot" view of the representatives at the time the method
// was called.
//
// Java: StandardUnionFind.allRepresentatives()
func (u *StandardUnionFind[E]) AllRepresentatives() *jutil.LinkedHashSet[E] {
	result := jutil.NewLinkedHashSet[E]()
	for n := range u.elmap.Values() {
		if n == n.parent {
			result.Add(n.element)
		}
	}
	return result
}

// findRootOrCreateNode: if e is already in a non-trivial equivalence class,
// that is, a class with more than two elements, then return the Node
// corresponding to the representative element. Otherwise, if e sits in an
// equivalence class by itself, then create a Node, put it into elmap and
// return it.
//
// Java: StandardUnionFind.findRootOrCreateNode(E)
func (u *StandardUnionFind[E]) findRootOrCreateNode(e E) *standardUnionFind_Node[E] {
	node := u.elmap.Get(e)
	if node != nil {
		return u.findRoot(node)
	}
	node = u.nodeSlab.New()
	node.element = e
	node.rank = 0
	node.size = 1
	node.parent = node
	u.elmap.Put(e, node)
	return node
}

// findRoot: given a Node, walk the parent field as far as possible, until
// reaching the root, which is the Node for the current representative of
// this equivalence class. To achieve low runtime complexity, also compress
// the path, by making each node a direct child of the root.
//
// Java: StandardUnionFind.findRoot(Node<E>)
func (u *StandardUnionFind[E]) findRoot(node *standardUnionFind_Node[E]) *standardUnionFind_Node[E] {
	if node.parent != node {
		node.parent = u.findRoot(node.parent)
	}
	return node.parent
}

// Java: StandardUnionFind.findAll(E)
func (u *StandardUnionFind[E]) FindAll(value E) jutil.Set[E] {
	jutil.CheckArgumentMsg1(u.elmap.ContainsKey(value), "Element does not exist: %s", value)

	isSameRoot := &standardUnionFind_isSameRoot[E]{u: u, value: value,
		// some node that's close to the root, or null
		nodeForValue: u.elmap.Get(value)}

	return &standardUnionFind_findAllSet[E]{u: u, value: value, isSameRoot: isSameRoot}
}

// standardUnionFind_isSameRoot is the anonymous Predicate<Object> of
// findAll.
type standardUnionFind_isSameRoot[E comparable] struct {
	u     *StandardUnionFind[E]
	value E

	// some node that's close to the root, or null
	nodeForValue *standardUnionFind_Node[E]
}

// Apply is Predicate.apply(Object).
func (p *standardUnionFind_isSameRoot[E]) Apply(b E) bool {
	if jutil.ObjectsEquals(p.value, b) {
		return true
	}
	nodeForB := p.u.elmap.Get(b)
	if nodeForB == nil {
		return false
	}
	p.nodeForValue = p.u.findRoot(p.nodeForValue)
	return p.u.findRoot(nodeForB) == p.nodeForValue
}

// standardUnionFind_findAllSet is the AbstractSet view returned by
// findAll.
type standardUnionFind_findAllSet[E comparable] struct {
	u          *StandardUnionFind[E]
	value      E
	isSameRoot *standardUnionFind_isSameRoot[E]
}

// Contains is AbstractSet.contains(Object).
func (s *standardUnionFind_findAllSet[E]) Contains(o E) bool {
	return s.isSameRoot.Apply(o)
}

// All is the iterator(): filter(elmap.keySet().iterator(), isSameRoot).
func (s *standardUnionFind_findAllSet[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for k := range s.u.elmap.Keys() {
			if s.isSameRoot.Apply(k) {
				if !yield(k) {
					return
				}
			}
		}
	}
}

// Size is AbstractSet.size().
func (s *standardUnionFind_findAllSet[E]) Size() int {
	return s.u.findRoot(s.u.elmap.Get(s.value)).size
}

// IsEmpty is AbstractCollection.isEmpty().
func (s *standardUnionFind_findAllSet[E]) IsEmpty() bool { return s.Size() == 0 }

// String is AbstractCollection.toString().
func (s *standardUnionFind_findAllSet[E]) String() string {
	var elems []E
	for e := range s.All() {
		elems = append(elems, e)
	}
	return jutil.CollectionToString(elems)
}

// standardUnionFind_Node is the internal node representation.
type standardUnionFind_Node[E comparable] struct {
	// The parent node of this element.
	parent *standardUnionFind_Node[E]

	// The element represented by this node.
	element E

	// A bound on the depth of the subtree rooted to this node.
	rank int

	// If this node is the root of a tree, this is the number of elements in
	// the tree. Otherwise, it's undefined.
	size int
}

// Java: StandardUnionFind.Node(E)
func newStandardUnionFind_Node[E comparable](element E) *standardUnionFind_Node[E] {
	n := &standardUnionFind_Node[E]{element: element, rank: 0, size: 1}
	n.parent = n
	return n
}
