// Port of com/google/javascript/rhino/HamtPMap.java @ closure-compiler v20261104

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"cmp"
	"iter"
	"math/bits"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// HamtPMap is an immutable sorted map with efficient (persistent) updates.
//
// Uses a hash array mapped trie:
// http://en.wikipedia.org/wiki/Hash_array_mapped_trie.
//
// This implementation stores the bare minimum in each node: a key (with its
// hash), value, mask, and children. It is also optimized to take maximum
// advantage of binary operations on shared trees. Specifically, reconcile
// avoids recursing into entire subtrees if they are identical objects. Null
// keys and values are not allowed. A nil *HamtPMap represents the empty map.
// Finally, we maintain an invariant that the entry with the smallest hash
// code is always at the root of the tree, which avoids almost all extra tree
// rebuilding during binary operations.
//
// TODO(sdh): Consider using tricks from
// https://bendyworks.com/blog/leveling-clojures-hash-maps
// We need a solid way to profile the results to see if it's actually worth
// the extra code.
//
// Go-only: language-adapter — keys are compared with == and hashed with
// unboxed fast paths for string, strpool.ID, Node, int, and HashCode() int32.
// A nil *HamtPMap[K, V] is the empty map.
type HamtPMap[K, V comparable] struct {
	// Non-null key.
	key K

	// Hash of the key, left-shifted by BITS*depth.
	hash uint32

	// Non-null value (exception: result of pivot, if not found).
	value V

	// Bit mask indicating the children that are present (bitCount(mask) ==
	// children.length).
	mask uint32

	// Array of children (nil when empty). Elements are never reassigned.
	children []*HamtPMap[K, V]
}

// Number of bits of fan-out at each level. May be anywhere from 1 (a binary
// tree) to 5 (for a fan-out of 32). Anything larger requires using a long
// for the mask (which is ugly for the JavaScript-compiled version), and zero
// is essentially a linked list. The relevant trade-off is that more fan-out
// means shallower trees, which should lead to quicker look-up times and more
// efficient updates, but less fan-out makes reconcile() and equivalent() more
// efficient for maps that are mostly the same, since bigger common trees can
// be pruned from the operation. Preliminary profiling suggests that for
// type-checking purposes, all nonzero values are roughly comparable.
const HamtPMap_BITS = 4

// Number of bits to shift off to get the most significant BITS number of
// bits.
const HamtPMap_BITS_SHIFT = 32 - HamtPMap_BITS

func hamtIsNull[V comparable](v V) bool {
	switch any(v).(type) {
	case nil:
		return true
	case string, int, int32, int64, int16, int8, uint, uint64, uint32, uint16, uint8, bool, float64, float32:
		return false
	}
	if n, ok := any(v).(Node); ok {
		return n.IsNil()
	}
	if _, ok := any(v).(interface{ IsNil() bool }); ok {
		return any(v).(interface{ IsNil() bool }).IsNil()
	}
	var zero V
	if v != zero {
		return false
	}
	return jutil.IsNil(any(v))
}

const debugHamtPMap = false

// HamtPMapArena provides bump-allocated storage for HamtPMap nodes and small
// child slices within a reusable analysis pass (e.g. TypeInference).
type HamtPMapArena[K, V comparable] struct {
	nodes       jutil.ReuseSlab[HamtPMap[K, V]]
	childBlocks jutil.ReuseSlab[[32]*HamtPMap[K, V]]
	childSlice  []*HamtPMap[K, V]
}

// Reset clears used arena storage and resets indices for reuse while retaining
// a bounded working set of allocated chunks (via jutil.ReuseSlab.Reset).
func (a *HamtPMapArena[K, V]) Reset() {
	a.nodes.Reset(nil)
	a.childBlocks.Reset(nil)
	a.childSlice = nil
}

// TrimLargeChunks resets the arena and drops chunks beyond jutil.ReuseSlab's
// bounded working set after a large global pass, retaining the initial working
// set for subsequent function-level analyses.
func (a *HamtPMapArena[K, V]) TrimLargeChunks() {
	a.Reset()
}

func allocHamtChildren[K, V comparable](a *HamtPMapArena[K, V], n int) []*HamtPMap[K, V] {
	if n == 0 {
		return nil
	}
	if a == nil || n > 16 {
		return make([]*HamtPMap[K, V], n)
	}
	if len(a.childSlice) < n {
		a.childSlice = a.childBlocks.Alloc()[:]
	}
	s := a.childSlice[:n:n]
	a.childSlice = a.childSlice[n:]
	return s
}

func allocHamtPMap[K, V comparable](a *HamtPMapArena[K, V], key K, hash uint32, value V, mask uint32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	var m *HamtPMap[K, V]
	if a != nil {
		m = a.nodes.Alloc()
	} else {
		m = &HamtPMap[K, V]{}
	}
	m.key = key
	m.hash = hash
	m.value = value
	m.mask = mask
	m.children = children
	if debugHamtPMap {
		m.CheckInvariants()
	}
	return m
}

// Java: HamtPMap.checkInvariants()
func (m *HamtPMap[K, V]) CheckInvariants() {
	jutil.CheckState(bits.OnesCount32(m.mask) == len(m.children))
	for _, child := range m.children {
		jutil.CheckNotNull(child)
	}
}

// HamtPMap_Empty returns an empty map (represented by nil).
//
// Java: HamtPMap.empty()
func HamtPMap_Empty[K, V comparable]() *HamtPMap[K, V] { return nil }

// Java: HamtPMap.toString()
func (m *HamtPMap[K, V]) String() string {
	if m == nil {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteString("{")
	m.AppendTo(&sb)
	sb.WriteString("}")
	return sb.String()
}

// AppendTo appends this map's contents to a string builder.
//
// Java: HamtPMap.appendTo(StringBuilder)
func (m *HamtPMap[K, V]) AppendTo(sb *strings.Builder) {
	if sb.Len() > 1 {
		sb.WriteString(", ")
	}
	sb.WriteString(m.keyString())
	sb.WriteString(": ")
	sb.WriteString(jutil.ToJavaString(m.value))
	for _, child := range m.children {
		child.AppendTo(sb)
	}
}

func (m *HamtPMap[K, V]) keyString() string {
	return jutil.ToJavaString(m.key)
}

// IsEmpty returns whether this map is empty.
//
// Java: HamtPMap.isEmpty()
func (m *HamtPMap[K, V]) IsEmpty() bool {
	return m == nil
}

// All returns an iterator over all key-value pairs in a (possibly nil) tree.
func (m *HamtPMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m == nil {
			return
		}
		var buf [16]*HamtPMap[K, V]
		queue := append(buf[:0], m)
		head := 0
		for head < len(queue) {
			top := queue[head]
			head++
			if !yield(top.key, top.value) {
				return
			}
			n := len(top.children)
			if n == 0 {
				if head == len(queue) {
					head = 0
					queue = queue[:0]
				}
				continue
			}
			if head > 0 && len(queue)+n > cap(queue) {
				rem := copy(queue[:cap(queue)], queue[head:])
				queue = queue[:rem]
				head = 0
			}
			for i := n - 1; i >= 0; i-- {
				queue = append(queue, top.children[i])
			}
		}
	}
}

// Values returns an iterable for a (possibly nil) tree.
//
// Java: HamtPMap.values()
func (m *HamtPMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}

// Keys returns an iterable for a (possibly nil) tree.
//
// Java: HamtPMap.keys()
func (m *HamtPMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Get retrieves the value associated with the given key from the map, or
// returns null if it is not present.
//
// Java: HamtPMap.get(K)
func (m *HamtPMap[K, V]) Get(key K) V {
	if m != nil {
		return m.getWithHash(key, hamtHash(key))
	}
	var null V
	return null
}

// getWithHash is the internal recursive implementation of get(K).
//
// Java: HamtPMap.get(K, int)
func (m *HamtPMap[K, V]) getWithHash(key K, hash uint32) V {
	if hash == m.hash && key == m.key {
		return m.value
	}
	bucketMask := uint32(1) << (hash >> HamtPMap_BITS_SHIFT)
	if (m.mask & bucketMask) != 0 {
		return m.children[m.index(bucketMask)].getWithHash(key, hash<<HamtPMap_BITS)
	}
	var null V
	return null
}

// HamtPMap_Reconciler merges two values; see HamtPMap.Reconcile.
type HamtPMap_Reconciler[K, V comparable] func(key K, thisVal V, thatVal V) V

// Plus returns a new map with the given key-value pair added. If the value
// is already present, then this same map will be returned.
//
// Java: HamtPMap.plus(K, V)
func (m *HamtPMap[K, V]) Plus(key K, value V) *HamtPMap[K, V] {
	return m.PlusMapWithArena(nil, key, value)
}

// PlusMapWithArena is Plus using arena for node and child-slice allocations.
func (m *HamtPMap[K, V]) PlusMapWithArena(arena *HamtPMapArena[K, V], key K, value V) *HamtPMap[K, V] {
	hamtCheckNotNull(value)

	if m != nil {
		return m.plusWithHash(arena, key, hamtHash(key), value)
	}
	return allocHamtPMap(arena, key, hamtHash(key), value, 0, nil)
}

// plusWithHash is the internal recursive implementation of plus(K, V).
//
// Java: HamtPMap.plus(K, int, V)
func (m *HamtPMap[K, V]) plusWithHash(arena *HamtPMapArena[K, V], key K, hash uint32, value V) *HamtPMap[K, V] {
	if hash == m.hash && key == m.key {
		if hamtValueEquals(value, m.value) {
			return m
		}
		return allocHamtPMap(arena, key, hash, value, m.mask, m.children)
	}
	if cmp.Compare(hash, m.hash) < 0 {
		return m.replaceRoot(arena, key, hash, value)
	}
	bucketMask := uint32(1) << (hash >> HamtPMap_BITS_SHIFT)
	hash <<= HamtPMap_BITS
	index := m.index(bucketMask)
	if (m.mask & bucketMask) != 0 {
		// already a child, so overwrite
		child := m.children[index]
		newChild := child.plusWithHash(arena, key, hash, value)
		if child == newChild {
			return m
		}
		return m.withChildren(arena, m.mask, replaceChild(arena, m.children, index, newChild))
	}
	// insert at index
	newChild := allocHamtPMap(arena, key, hash, value, 0, nil)
	return m.withChildren(arena, m.mask|bucketMask, insertChild(arena, m.children, index, newChild))
}

// Java: HamtPMap.replaceRoot(K, int, V)
func (m *HamtPMap[K, V]) replaceRoot(arena *HamtPMapArena[K, V], key K, hash uint32, value V) *HamtPMap[K, V] {
	bucketMask := uint32(1) << (m.hash >> HamtPMap_BITS_SHIFT)
	leafHash := m.hash << HamtPMap_BITS
	index := m.index(bucketMask)
	var newChildren []*HamtPMap[K, V]
	if (m.mask & bucketMask) != 0 {
		newChildren = replaceChild(arena, m.children, index, m.children[index].plusWithHash(arena, m.key, leafHash, m.value))
	} else {
		newChild := allocHamtPMap(arena, m.key, leafHash, m.value, 0, nil)
		newChildren = insertChild(arena, m.children, index, newChild)
	}
	return allocHamtPMap(arena, key, hash, value, m.mask|bucketMask, newChildren)
}

// Minus returns a new map with the given key removed. If the key was not
// present in the first place, then this same map will be returned.
//
// Java: HamtPMap.minus(K)
func (m *HamtPMap[K, V]) Minus(key K) *HamtPMap[K, V] {
	return m.MinusMapWithArena(nil, key)
}

// MinusMapWithArena is Minus using arena for node and child-slice allocations.
func (m *HamtPMap[K, V]) MinusMapWithArena(arena *HamtPMapArena[K, V], key K) *HamtPMap[K, V] {
	if m != nil {
		res, _, _ := m.minusWithHash(arena, key, hamtHash(key))
		return res
	}
	return nil
}

// minusWithHash is the internal recursive implementation of minus(K), returning
// the updated subtree, the removed value, and whether key was present.
//
// Java: HamtPMap.minus(K, int, V[])
func (m *HamtPMap[K, V]) minusWithHash(arena *HamtPMapArena[K, V], key K, hash uint32) (*HamtPMap[K, V], V, bool) {
	var zero V
	if hash == m.hash && key == m.key {
		return deleteRoot(arena, m.mask, m.children), m.value, true
	}
	bucketMask := uint32(1) << (hash >> HamtPMap_BITS_SHIFT)
	if (m.mask & bucketMask) == 0 {
		// not present, stop looking
		return m, zero, false
	}
	hash <<= HamtPMap_BITS
	index := m.index(bucketMask)
	child := m.children[index]
	newChild, removed, found := child.minusWithHash(arena, key, hash)
	if !found {
		return m, zero, false
	} else if newChild == nil {
		return m.withChildren(arena, m.mask&^bucketMask, deleteChild(arena, m.children, index)), removed, true
	} else {
		return m.withChildren(arena, m.mask, replaceChild(arena, m.children, index, newChild)), removed, true
	}
}

// Reconcile performs a reconcile operation to merge `m` and `that`.
//
// Java: HamtPMap.reconcile(PMap<K, V>, Reconciler<K, V>)
func (m *HamtPMap[K, V]) Reconcile(that *HamtPMap[K, V], joiner HamtPMap_Reconciler[K, V]) *HamtPMap[K, V] {
	return m.ReconcileMapWithArena(nil, that, joiner)
}

// ReconcileMapWithArena is Reconcile using arena for node and child-slice allocations.
func (m *HamtPMap[K, V]) ReconcileMapWithArena(arena *HamtPMapArena[K, V], that *HamtPMap[K, V], joiner HamtPMap_Reconciler[K, V]) *HamtPMap[K, V] {
	if m == that {
		return m
	}
	return hamtPMap_ReconcileTrees(arena, m, that, joiner, true)
}

func reconcileSingleTree[K, V comparable](arena *HamtPMapArena[K, V], t *HamtPMap[K, V], onLeft bool, joiner HamtPMap_Reconciler[K, V], checkNotNull bool) *HamtPMap[K, V] {
	var null, newValue V
	if onLeft {
		newValue = joiner(t.key, t.value, null)
	} else {
		newValue = joiner(t.key, null, t.value)
	}
	if checkNotNull {
		hamtCheckNotNull(newValue)
	}
	var newChildrenBuf [16]*HamtPMap[K, V]
	n := len(t.children)
	sameChildren := true
	for i := 0; i < n; i++ {
		c := reconcileSingleTree(arena, t.children[i], onLeft, joiner, checkNotNull)
		newChildrenBuf[i] = c
		if c != t.children[i] {
			sameChildren = false
		}
	}
	if sameChildren && hamtValueEquals(t.value, newValue) {
		return t
	}
	var newChildren []*HamtPMap[K, V]
	if sameChildren {
		newChildren = t.children
	} else {
		newChildren = allocHamtChildren(arena, n)
		copy(newChildren, newChildrenBuf[:n])
	}
	if checkNotNull || !hamtIsNull(newValue) {
		return allocHamtPMap(arena, t.key, t.hash, newValue, t.mask, newChildren)
	}
	return deleteRoot(arena, t.mask, newChildren)
}

func reconcileVacatedChild[K, V comparable](arena *HamtPMapArena[K, V], vacated, ownChild, otherChild *HamtPMap[K, V], vacatedOnLeft bool, joiner HamtPMap_Reconciler[K, V], checkNotNull bool) *HamtPMap[K, V] {
	shiftedHash := vacated.hash << HamtPMap_BITS
	if ownChild == nil && otherChild == nil {
		var null, v V
		if vacatedOnLeft {
			v = joiner(vacated.key, vacated.value, null)
		} else {
			v = joiner(vacated.key, null, vacated.value)
		}
		if checkNotNull {
			hamtCheckNotNull(v)
		}
		if checkNotNull || !hamtIsNull(v) {
			return allocHamtPMap(arena, vacated.key, shiftedHash, v, 0, nil)
		}
		return nil
	}
	var vacatedChild *HamtPMap[K, V]
	if ownChild != nil {
		vacatedChild = ownChild.plusWithHash(arena, vacated.key, shiftedHash, vacated.value)
	} else {
		vacatedChild = allocHamtPMap(arena, vacated.key, shiftedHash, vacated.value, 0, nil)
	}
	if vacatedOnLeft {
		return hamtPMap_ReconcileTrees(arena, vacatedChild, otherChild, joiner, checkNotNull)
	}
	return hamtPMap_ReconcileTrees(arena, otherChild, vacatedChild, joiner, checkNotNull)
}

// hamtPMap_ReconcileTrees is the internal recursive implementation of
// reconcile(HamtPMap, BiFunction), factoring out empties.
//
// Java: HamtPMap.reconcile(HamtPMap<K, V>, HamtPMap<K, V>, Reconciler<K, V>)
func hamtPMap_ReconcileTrees[K, V comparable](arena *HamtPMapArena[K, V], t1 *HamtPMap[K, V], t2 *HamtPMap[K, V], joiner HamtPMap_Reconciler[K, V], checkNotNull bool) *HamtPMap[K, V] {
	if t1 == t2 {
		return t1
	} else if t1 == nil {
		return reconcileSingleTree(arena, t2, false, joiner, checkNotNull)
	} else if t2 == nil {
		return reconcileSingleTree(arena, t1, true, joiner, checkNotNull)
	}

	// Try as hard as possible to return input trees exactly.
	sameChildrenAs1 := true
	sameChildrenAs2 := true

	// If the hashes are different, we need to keep the lower one at the top.
	hashCmp := cmp.Compare(t1.hash, t2.hash)
	key := t1.key
	hash := t1.hash
	var t1Val, t2Val V
	var vacatedBit1, vacatedBit2 uint32
	bothPresent := false
	if hashCmp < 0 {
		// t1.key is missing from t2; inline t2.VacateRoot() without allocating intermediate root/slice
		vacatedBit2 = uint32(1) << (t2.hash >> HamtPMap_BITS_SHIFT)
		sameChildrenAs2 = false
		t1Val = t1.value
	} else if hashCmp > 0 {
		// t2.key is missing from t1; inline t1.VacateRoot() without allocating intermediate root/slice
		vacatedBit1 = uint32(1) << (t1.hash >> HamtPMap_BITS_SHIFT)
		sameChildrenAs1 = false
		key = t2.key
		hash = t2.hash
		t2Val = t2.value
	} else {
		if t1.key != t2.key {
			// Hash collision: try to rearrange t2 to have the same root as t1
			t2, bothPresent = t2.pivot(arena, t1.key, t1.hash)
			sameChildrenAs2 = false
		} else {
			bothPresent = true
		}
		t1Val = t1.value
		t2Val = t2.value
	}

	var newValue V
	if bothPresent && hamtValueEquals(t1Val, t2Val) {
		newValue = t1Val
	} else {
		newValue = joiner(key, t1Val, t2Val)
		if checkNotNull {
			hamtCheckNotNull(newValue)
		}
	}
	t1Mask := t1.mask | vacatedBit1
	t2Mask := t2.mask | vacatedBit2
	newMask := t1Mask | t2Mask
	sameChildrenAs1 = sameChildrenAs1 && (newMask == t1Mask)
	sameChildrenAs2 = sameChildrenAs2 && (newMask == t2Mask)

	var newChildrenBuf [16]*HamtPMap[K, V]
	mask := newMask
	index := 0
	for mask != 0 {
		childBit := mask & -mask // lowest set bit
		mask &^= childBit
		child1 := t1.getChild(childBit)
		child2 := t2.getChild(childBit)
		var reconciled *HamtPMap[K, V]
		if childBit == vacatedBit2 {
			reconciled = reconcileVacatedChild(arena, t2, child2, child1, false, joiner, checkNotNull)
		} else if childBit == vacatedBit1 {
			reconciled = reconcileVacatedChild(arena, t1, child1, child2, true, joiner, checkNotNull)
		} else {
			reconciled = hamtPMap_ReconcileTrees(arena, child1, child2, joiner, checkNotNull)
		}
		newChildrenBuf[index] = reconciled
		sameChildrenAs1 = sameChildrenAs1 && (reconciled == child1)
		sameChildrenAs2 = sameChildrenAs2 && (reconciled == child2)
		if reconciled != nil {
			index++
		} else {
			newMask &^= childBit
		}
	}
	if sameChildrenAs1 && hamtValueEquals(t1Val, newValue) {
		return t1
	} else if sameChildrenAs2 && hamtValueEquals(t2Val, newValue) {
		return t2
	}
	var newChildren []*HamtPMap[K, V]
	if sameChildrenAs1 {
		newChildren = t1.children
	} else if sameChildrenAs2 {
		newChildren = t2.children
	} else if index > 0 {
		newChildren = allocHamtChildren(arena, index)
		copy(newChildren, newChildrenBuf[:index])
	}
	if !checkNotNull && hamtIsNull(newValue) {
		return deleteRoot(arena, newMask, newChildren)
	}
	return allocHamtPMap(arena, key, hash, newValue, newMask, newChildren)
}

// Equivalent checks equality recursively based on the given equivalence.
// Short-circuits as soon as a 'false' result is found.
//
// Java: HamtPMap.equivalent(PMap<K, V>, BiPredicate<V, V>)
func (m *HamtPMap[K, V]) Equivalent(that *HamtPMap[K, V], equivalence func(V, V) bool) bool {
	return HamtPMap_EquivalentTrees(m, that, equivalence)
}

// HamtPMap_EquivalentTrees is the internal recursive implementation of
// equivalent(HamtPMap, BiPredicate).
//
// Java: HamtPMap.equivalent(HamtPMap<K, V>, HamtPMap<K, V>, BiPredicate<V, V>)
func HamtPMap_EquivalentTrees[K, V comparable](t1 *HamtPMap[K, V], t2 *HamtPMap[K, V], equivalence func(V, V) bool) bool {
	if t1 == t2 {
		return true
	} else if t1 == nil || t2 == nil {
		return false
	}

	if t1.hash != t2.hash {
		// Due to the invariant that keys that are .equals must have the same
		// hash code, we can safely conclude that there's a discrepancy in the
		// keys without any extra work.
		return false
	} else if t1.key != t2.key {
		// Hash collision: try to rearrange t2 to have the same root as t1
		var ok bool
		t2, ok = t2.pivot(nil, t1.key, t1.hash)
		if !ok {
			return false
		}
	}

	if !equivalence(t1.value, t2.value) {
		return false
	}
	mask := t1.mask | t2.mask
	for mask != 0 {
		childBit := mask & -mask // lowest set bit
		mask &^= childBit
		if !HamtPMap_EquivalentTrees(t1.getChild(childBit), t2.getChild(childBit), equivalence) {
			return false
		}
	}
	return true
}

// index returns the index into the 'children' array for the given bit,
// which must have exactly one bit set in its binary representation (i.e.
// must be a power of two).
//
// Java: HamtPMap.index(int)
func (m *HamtPMap[K, V]) index(bit uint32) int {
	return bits.OnesCount32(m.mask & (bit - 1))
}

// getChild returns the child for the given bit, which must have exactly one
// bit set. Returns nil if there is no child for that bit.
//
// Java: HamtPMap.getChild(int)
func (m *HamtPMap[K, V]) getChild(bit uint32) *HamtPMap[K, V] {
	if (m.mask & bit) != 0 {
		return m.children[m.index(bit)]
	}
	return nil
}

// hamtHash performs the hash operation. This is done once per object as
// it enters the map, then never again since the result is passed around
// internally. We reverse the bits since certain types (such as small
// integers and short strings) have hash codes that only vary in the least
// significant bits, but we use the most significant bits for bucketing in
// order to maintain a canonical structure where the key with the smallest
// hash is always at the root of the tree.
//
// Java: HamtPMap.hash(Object)
func hamtHash[K comparable](key K) uint32 {
	if any(key) == nil {
		failJSTypeNull()
	}
	// Go-only: performance — common scalar/handle keys are hashed via direct
	// non-escaping type assertions without boxing onto the heap.
	if id, ok := any(key).(strpool.ID); ok {
		if id == 0 {
			failJSTypeNull()
		}
		h := uint32(id) * 0x9e3779b1
		return bits.Reverse32(h ^ (h >> 16))
	}
	if s, ok := any(key).(string); ok {
		return bits.Reverse32(uint32(jutil.StringHashCode(s)))
	}
	if n, ok := any(key).(Node); ok {
		return bits.Reverse32(uint32(CheckNotNull(n).IdentityHashCode()))
	}
	if x, ok := any(key).(int); ok {
		return bits.Reverse32(uint32(int32(x)))
	}
	if hamtIsNull(key) {
		failJSTypeNull()
	}
	if h, ok := any(key).(interface{ HashCode() int32 }); ok {
		return bits.Reverse32(uint32(h.HashCode()))
	}
	panic(jutil.NewUnsupportedOperationException())
}

func hamtValueEquals[V comparable](a, b V) bool {
	if a == b {
		return true
	}
	if hamtIsNull(a) || hamtIsNull(b) {
		return false
	}
	if ja, ok := any(a).(JSType); ok {
		jb, ok := any(b).(JSType)
		return ok && jsTypeEquals(ja, jb)
	}
	if ea, ok := any(a).(interface{ Equals(any) bool }); ok {
		return ea.Equals(any(b))
	}
	return false
}

func hamtCheckNotNull[V comparable](v V) {
	if hamtIsNull(v) {
		failJSTypeNull()
	}
}

// pivot returns a new version of this map with the given key at the root,
// and the root element moved to some deeper node, along with whether key was
// found in this map.
//
// Java: HamtPMap.pivot(K, int)
func (m *HamtPMap[K, V]) pivot(arena *HamtPMapArena[K, V], key K, hash uint32) (*HamtPMap[K, V], bool) {
	res, _, ok := m.pivotWithParent(arena, key, hash, nil)
	return res, ok
}

// pivotWithParent is the internal recursive version of pivot, returning the
// pivoted subtree, the found value (zero V if not present), and whether key
// was present. If parent is nil, returns a map rooted at key with the found
// value; otherwise returns a map rooted at parent.key.
//
// Java: HamtPMap.pivot(K, int, HamtPMap<K, V>, V[])
func (m *HamtPMap[K, V]) pivotWithParent(arena *HamtPMapArena[K, V], key K, hash uint32, parent *HamtPMap[K, V]) (*HamtPMap[K, V], V, bool) {
	var found V
	var ok bool
	newMask := m.mask
	newChildren := m.children
	if hash == m.hash && key == m.key {
		found = m.value
		ok = true
	} else {
		// Otherwise, see if the value might be present in a child.
		searchBucket := hash >> HamtPMap_BITS_SHIFT
		replacementBucket := m.hash >> HamtPMap_BITS_SHIFT

		if searchBucket == replacementBucket {
			bucketMask := uint32(1) << searchBucket
			if (m.mask & bucketMask) != 0 {
				// Found a child, call pivot recursively.
				index := m.index(bucketMask)
				child := newChildren[index]
				var newChild *HamtPMap[K, V]
				newChild, found, ok = child.pivotWithParent(arena, key, hash<<HamtPMap_BITS, m)
				newChildren = replaceChild(arena, newChildren, index, newChild)
			}
		} else {
			searchMask := uint32(1) << searchBucket
			if (m.mask & searchMask) != 0 {
				index := m.index(searchMask)
				child := newChildren[index]
				var newChild *HamtPMap[K, V]
				newChild, found, ok = child.minusWithHash(arena, key, hash<<HamtPMap_BITS)
				if newChild != nil {
					newChildren = replaceChild(arena, newChildren, index, newChild)
				} else {
					newChildren = deleteChild(arena, newChildren, index)
					newMask &^= searchMask
				}
			}
			replacementMask := uint32(1) << replacementBucket
			index := bits.OnesCount32(newMask & (replacementMask - 1))
			if (m.mask & replacementMask) != 0 {
				child := newChildren[index]
				newChild := child.plusWithHash(arena, m.key, m.hash<<HamtPMap_BITS, m.value)
				newChildren = replaceChild(arena, newChildren, index, newChild)
			} else {
				newChildren = insertChild(
					arena,
					newChildren,
					index,
					allocHamtPMap(arena, m.key, m.hash<<HamtPMap_BITS, m.value, 0, nil))
				newMask |= replacementMask
			}
		}
	}
	if parent != nil {
		return allocHamtPMap(arena, parent.key, parent.hash<<HamtPMap_BITS, parent.value, newMask, newChildren), found, ok
	}
	return allocHamtPMap(arena, key, hash, found, newMask, newChildren), found, ok
}

// withChildren returns a copy of this node with a different array of
// children.
//
// Java: HamtPMap.withChildren(int, HamtPMap<K, V>[])
func (m *HamtPMap[K, V]) withChildren(arena *HamtPMapArena[K, V], mask uint32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	if mask == m.mask && hamtSameArray(children, m.children) {
		return m
	}
	return allocHamtPMap(arena, m.key, m.hash, m.value, mask, children)
}

// hamtSameArray is Java array identity (`a == b`) on Go slices.
func hamtSameArray[T any](a, b []T) bool {
	return len(a) == len(b) && (len(a) == 0 || &a[0] == &b[0])
}

// deleteRoot returns a new map with the elements from children. One
// element is removed from one of the children and promoted to a root node.
// If there are no children, returns nil.
//
// Java: HamtPMap.deleteRoot(int, HamtPMap<K, V>[])
func deleteRoot[K, V comparable](arena *HamtPMapArena[K, V], mask uint32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	if mask == 0 {
		return nil
	}
	child := children[0]
	hashBits := uint32(bits.TrailingZeros32(mask))
	newHash := (child.hash >> HamtPMap_BITS) | (hashBits << HamtPMap_BITS_SHIFT)
	newChild := deleteRoot(arena, child.mask, child.children)
	if newChild == nil {
		newMask := mask &^ (mask & -mask) // lowest set bit
		return allocHamtPMap(arena, child.key, newHash, child.value, newMask, deleteChild(arena, children, 0))
	}
	return allocHamtPMap(arena, child.key, newHash, child.value, mask, replaceChild(arena, children, 0, newChild))
}

// insertChild returns a new array of children with an additional
// child inserted at the given index.
//
// Java: HamtPMap.insertChild(HamtPMap<K, V>[], int, HamtPMap<K, V>)
func insertChild[K, V comparable](arena *HamtPMapArena[K, V], children []*HamtPMap[K, V], index int, child *HamtPMap[K, V]) []*HamtPMap[K, V] {
	newChildren := allocHamtChildren(arena, len(children)+1)
	newChildren[index] = child
	copy(newChildren[0:index], children[0:index])
	copy(newChildren[index+1:], children[index:])
	return newChildren
}

// replaceChild returns a new array of children with the child at
// the given index replaced.
//
// Java: HamtPMap.replaceChild(HamtPMap<K, V>[], int, HamtPMap<K, V>)
func replaceChild[K, V comparable](arena *HamtPMapArena[K, V], children []*HamtPMap[K, V], index int, child *HamtPMap[K, V]) []*HamtPMap[K, V] {
	newChildren := allocHamtChildren(arena, len(children))
	copy(newChildren, children)
	newChildren[index] = child
	return newChildren
}

// deleteChild returns a new array of children with the child at the
// given index deleted.
//
// Java: HamtPMap.deleteChild(HamtPMap<K, V>[], int)
func deleteChild[K, V comparable](arena *HamtPMapArena[K, V], children []*HamtPMap[K, V], index int) []*HamtPMap[K, V] {
	if len(children) == 1 {
		// Note: index should always be zero.
		return nil
	}
	newChildren := allocHamtChildren(arena, len(children)-1)
	copy(newChildren[0:index], children[0:index])
	copy(newChildren[index:], children[index+1:])
	return newChildren
}

// AssertCorrectStructure throws an assertion error if the map invariant is
// violated.
//
// Java: HamtPMap.assertCorrectStructure()
func (m *HamtPMap[K, V]) AssertCorrectStructure() *HamtPMap[K, V] {
	if m == nil {
		return nil
	}
	hash := hamtHash(m.key)
	for i := 0; i < len(m.children); i++ {
		childHash := hamtHash(m.children[i].key)
		if cmp.Compare(childHash, hash) < 0 {
			panic(jutil.NewAssertionError("Invalid map has decreasing hash " + m.children[i].keyString() + "(" +
				jutil.IntegerToHexString(int32(childHash)) + ") beneath " + m.keyString() + "(" +
				jutil.IntegerToHexString(int32(hash)) + ": " + m.String()))
		}
		m.children[i].AssertCorrectStructure()
	}
	return m
}
