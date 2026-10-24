// Port of com/google/javascript/rhino/HamtPMap.java @ closure-compiler v20261021

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
	"iter"
	"math/bits"
	"reflect"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
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
// keys and values are not allowed. The implementation special-cases away the
// EMPTY map as soon as possible, using 'null' instead for all the logic
// (since EMPTY violates the invariant that key and value are non-null).
// Finally, we maintain an invariant that the entry with the smallest hash
// code is always at the root of the tree, which avoids almost all extra tree
// rebuilding during binary operations.
//
// TODO(sdh): Consider using tricks from
// https://bendyworks.com/blog/leveling-clojures-hash-maps
// We need a solid way to profile the results to see if it's actually worth
// the extra code.
//
// Go-only: language-adapter — keys are hashed and compared with Java
// hashCode/equals (jutil.ObjectHashCode / ObjectsEquals). A null key (only
// the EMPTY map and internal vacated roots have one) is recorded with
// keyIsNull, since K may be a Go value type such as string. Null values are
// the values for which jutil.IsNil is true.
type HamtPMap[K any, V any] struct {
	// Non-null key (exception: empty map has a null key).
	key       K
	keyIsNull bool

	// Hash of the key, right-shifted by BITS*depth.
	hash int32

	// Non-null value (exceptions: (1) empty map, (2) result of pivot, if not
	// found).
	value V

	// Bit mask indicating the children that are present (bitCount(mask) ==
	// children.length).
	mask int32

	// Non-null array of children. Elements are never reassigned.
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

func hamtIsNull[V any](v V) bool { return jutil.IsNil(any(v)) }

// hamtPMap_EMPTY holds the EMPTY singleton of each instantiation.
//
// Go-only: language-adapter — Java has one EMPTY instance cast to every
// <K, V>; Go generics need one per instantiation (identity is observable:
// minus() compares against EMPTY and empty() is always the same object).
var hamtPMap_EMPTY = map[reflect.Type]any{}

const debugHamtPMap = false

// HamtPMapArena provides bump-allocated storage for HamtPMap nodes and small
// child slices within a reusable analysis pass (e.g. TypeInference or
// MaybeReachingVariableUse).
type HamtPMapArena[K any, V any] struct {
	nodeBuf       [64]HamtPMap[K, V]
	nodeInlineLen int
	nodeChunks    []*[256]HamtPMap[K, V]
	nodeChunkIdx  int
	nodeOffset    int

	childBuf       [128]*HamtPMap[K, V]
	childInlineLen int
	childChunks    []*[1024]*HamtPMap[K, V]
	childChunkIdx  int
	childOffset    int
}

// Reset clears used arena storage and resets indices for reuse while retaining
// allocated chunk capacity.
func (a *HamtPMapArena[K, V]) Reset() {
	if a.nodeInlineLen > 0 {
		clear(a.nodeBuf[:a.nodeInlineLen])
		a.nodeInlineLen = 0
	}
	if a.nodeChunkIdx > 0 || a.nodeOffset > 0 {
		for i := 0; i < a.nodeChunkIdx; i++ {
			clear(a.nodeChunks[i][:])
		}
		if a.nodeOffset > 0 && a.nodeChunkIdx < len(a.nodeChunks) {
			clear(a.nodeChunks[a.nodeChunkIdx][:a.nodeOffset])
		}
		a.nodeChunkIdx = 0
		a.nodeOffset = 0
	}

	if a.childInlineLen > 0 {
		clear(a.childBuf[:a.childInlineLen])
		a.childInlineLen = 0
	}
	if a.childChunkIdx > 0 || a.childOffset > 0 {
		for i := 0; i < a.childChunkIdx; i++ {
			clear(a.childChunks[i][:])
		}
		if a.childOffset > 0 && a.childChunkIdx < len(a.childChunks) {
			clear(a.childChunks[a.childChunkIdx][:a.childOffset])
		}
		a.childChunkIdx = 0
		a.childOffset = 0
	}
}

// TrimLargeChunks releases excess arena chunks after a large global pass while
// retaining a small working set for function-level analyses.
func (a *HamtPMapArena[K, V]) TrimLargeChunks() {
	a.Reset()
	const maxRetained = 8
	if len(a.nodeChunks) > maxRetained {
		clear(a.nodeChunks[maxRetained:])
		a.nodeChunks = a.nodeChunks[:maxRetained]
	}
	if len(a.childChunks) > maxRetained {
		clear(a.childChunks[maxRetained:])
		a.childChunks = a.childChunks[:maxRetained]
	}
}

func (a *HamtPMapArena[K, V]) allocNode() *HamtPMap[K, V] {
	if a.nodeInlineLen < len(a.nodeBuf) {
		m := &a.nodeBuf[a.nodeInlineLen]
		a.nodeInlineLen++
		return m
	}
	return a.allocNodeChunk()
}

//go:noinline
func (a *HamtPMapArena[K, V]) allocNodeChunk() *HamtPMap[K, V] {
	if a.nodeOffset == 256 {
		a.nodeChunkIdx++
		a.nodeOffset = 0
	}
	if a.nodeChunkIdx == len(a.nodeChunks) {
		a.nodeChunks = append(a.nodeChunks, new([256]HamtPMap[K, V]))
	}
	m := &a.nodeChunks[a.nodeChunkIdx][a.nodeOffset]
	a.nodeOffset++
	return m
}

func allocHamtChildren[K any, V any](a *HamtPMapArena[K, V], n int) []*HamtPMap[K, V] {
	if n == 0 {
		return HamtPMap_EmptyChildren[K, V]()
	}
	if a == nil || n > 16 {
		return make([]*HamtPMap[K, V], n)
	}
	if a.childInlineLen+n <= len(a.childBuf) {
		off := a.childInlineLen
		a.childInlineLen = off + n
		return a.childBuf[off : off+n : off+n]
	}
	return a.allocChildrenChunk(n)
}

//go:noinline
func (a *HamtPMapArena[K, V]) allocChildrenChunk(n int) []*HamtPMap[K, V] {
	if a.childOffset+n > 1024 {
		if a.childOffset > 0 {
			a.childChunkIdx++
			a.childOffset = 0
		}
	}
	if a.childChunkIdx == len(a.childChunks) {
		a.childChunks = append(a.childChunks, new([1024]*HamtPMap[K, V]))
	}
	ch := a.childChunks[a.childChunkIdx]
	off := a.childOffset
	a.childOffset = off + n
	return ch[off : off+n : off+n]
}

func allocHamtPMap[K any, V any](a *HamtPMapArena[K, V], key K, keyIsNull bool, hash int32, value V, mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	var m *HamtPMap[K, V]
	if a != nil {
		m = a.allocNode()
	} else {
		m = &HamtPMap[K, V]{}
	}
	m.key = key
	m.keyIsNull = keyIsNull
	m.hash = hash
	m.value = value
	m.mask = mask
	m.children = children
	if debugHamtPMap {
		m.CheckInvariants()
	}
	return m
}

func allocHamtPMapNullKey[K any, V any](a *HamtPMapArena[K, V], value V, mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	var key K
	return allocHamtPMap(a, key, true, 0, value, mask, children)
}

// Java: HamtPMap(K, int, V, int, HamtPMap<K, V>[])
func newHamtPMap[K any, V any](key K, keyIsNull bool, hash int32, value V, mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	return allocHamtPMap(nil, key, keyIsNull, hash, value, mask, children)
}

func newHamtPMapNullKey[K any, V any](value V, mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	return allocHamtPMapNullKey(nil, value, mask, children)
}

// Java: HamtPMap.checkInvariants()
func (m *HamtPMap[K, V]) CheckInvariants() {
	jutil.CheckState(bits.OnesCount32(uint32(m.mask)) == len(m.children))
	for _, child := range m.children {
		jutil.CheckNotNull(child)
	}
}

// HamtPMap_Empty returns an empty map.
//
// Java: HamtPMap.empty()
func HamtPMap_Empty[K any, V any]() *HamtPMap[K, V] {
	t := reflect.TypeFor[HamtPMap[K, V]]()
	if e, ok := hamtPMap_EMPTY[t]; ok {
		return e.(*HamtPMap[K, V])
	}
	var nullValue V
	e := newHamtPMapNullKey[K, V](nullValue, 0, HamtPMap_EmptyChildren[K, V]())
	hamtPMap_EMPTY[t] = e
	return e
}

// HamtPMap_EmptyChildren returns an empty array of child maps.
//
// Java: HamtPMap.emptyChildren()
func HamtPMap_EmptyChildren[K any, V any]() []*HamtPMap[K, V] {
	return []*HamtPMap[K, V]{}
}

// Java: HamtPMap.toString()
func (m *HamtPMap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	if !m.IsEmpty() {
		m.AppendTo(&sb)
	}
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
	if m.keyIsNull {
		return "null"
	}
	return jutil.ToJavaString(m.key)
}

// IsEmpty returns whether this map is empty.
//
// Java: HamtPMap.isEmpty()
func (m *HamtPMap[K, V]) IsEmpty() bool {
	return m.keyIsNull
}

// Values returns an iterable for a (possibly null) tree.
//
// Java: HamtPMap.values()
func (m *HamtPMap[K, V]) Values() iter.Seq[V] {
	if m.IsEmpty() {
		return func(yield func(V) bool) {}
	}
	return func(yield func(V) bool) {
		it := NewHamtPMap_Iter(m, func(m *HamtPMap[K, V]) V { return m.value })
		for it.HasNext() {
			if !yield(it.Next()) {
				return
			}
		}
	}
}

// Keys returns an iterable for a (possibly null) tree.
//
// Java: HamtPMap.keys()
func (m *HamtPMap[K, V]) Keys() iter.Seq[K] {
	if m.IsEmpty() {
		return func(yield func(K) bool) {}
	}
	return func(yield func(K) bool) {
		it := NewHamtPMap_Iter(m, func(m *HamtPMap[K, V]) K { return m.key })
		for it.HasNext() {
			if !yield(it.Next()) {
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
	if !m.IsEmpty() {
		return m.GetWithHash(key, HamtPMap_Hash(key))
	}
	var null V
	return null
}

// keyEquals is Java's `key.equals(this.key)`.
func (m *HamtPMap[K, V]) keyEquals(key K) bool {
	// Go-only: performance — ObjectsEqualsOf compares String keys without
	// boxing them.
	return !m.keyIsNull && jutil.ObjectsEqualsOf(key, m.key)
}

// GetWithHash is the internal recursive implementation of get(K).
//
// Java: HamtPMap.get(K, int)
func (m *HamtPMap[K, V]) GetWithHash(key K, hash int32) V {
	if hash == m.hash && m.keyEquals(key) {
		return m.value
	}
	bucket := HamtPMap_Bucket(hash)
	bucketMask := jutil.Shl32(1, int64(bucket))
	if (m.mask & bucketMask) != 0 {
		return m.children[m.Index(bucketMask)].GetWithHash(key, HamtPMap_Shift(hash))
	}
	var null V
	return null
}

// Plus returns a new map with the given key-value pair added. If the value
// is already present, then this same map will be returned.
//
// Java: HamtPMap.plus(K, V)
func (m *HamtPMap[K, V]) Plus(key K, value V) PMap[K, V] {
	return m.PlusMap(key, value)
}

// PlusMap is plus(K, V) with its covariant HamtPMap return type.
//
// Go-only: language-adapter — Java's override narrows the return type to
// HamtPMap; a Go method satisfying PMap must return PMap.
func (m *HamtPMap[K, V]) PlusMap(key K, value V) *HamtPMap[K, V] {
	return m.PlusMapWithArena(nil, key, value)
}

// PlusMapWithArena is PlusMap using arena for node and child-slice allocations.
func (m *HamtPMap[K, V]) PlusMapWithArena(arena *HamtPMapArena[K, V], key K, value V) *HamtPMap[K, V] {
	jutil.CheckNotNull(any(value))

	if !m.IsEmpty() {
		return m.plusWithHash(arena, key, HamtPMap_Hash(key), value)
	}
	return allocHamtPMap(arena, key, false, HamtPMap_Hash(key), value, 0, HamtPMap_EmptyChildren[K, V]())
}

// PlusWithHash is the internal recursive implementation of plus(K, V).
//
// Java: HamtPMap.plus(K, int, V)
func (m *HamtPMap[K, V]) PlusWithHash(key K, hash int32, value V) *HamtPMap[K, V] {
	return m.plusWithHash(nil, key, hash, value)
}

func (m *HamtPMap[K, V]) plusWithHash(arena *HamtPMapArena[K, V], key K, hash int32, value V) *HamtPMap[K, V] {
	if hash == m.hash && m.keyEquals(key) {
		if jutil.ObjectsEqualsOf(value, m.value) {
			return m
		}
		return allocHamtPMap(arena, key, false, hash, value, m.mask, m.children)
	}
	if HamtPMap_CompareUnsigned(hash, m.hash) < 0 {
		return m.replaceRoot(arena, key, hash, value)
	}
	bucket := HamtPMap_Bucket(hash)
	hash = HamtPMap_Shift(hash)
	bucketMask := jutil.Shl32(1, int64(bucket))
	index := m.Index(bucketMask)
	if (m.mask & bucketMask) != 0 {
		// already a child, so overwrite
		child := m.children[index]
		newChild := child.plusWithHash(arena, key, hash, value)
		if child == newChild {
			return m
		}
		return m.withChildren(arena, m.mask, replaceChild(arena, m.children, index, newChild))
	} else {
		// insert at index
		newChild := allocHamtPMap(arena, key, false, hash, value, 0, HamtPMap_EmptyChildren[K, V]())
		return m.withChildren(arena, m.mask|bucketMask, insertChild(arena, m.children, index, newChild))
	}
}

// Java: HamtPMap.replaceRoot(K, int, V)
func (m *HamtPMap[K, V]) ReplaceRoot(key K, hash int32, value V) *HamtPMap[K, V] {
	return m.replaceRoot(nil, key, hash, value)
}

func (m *HamtPMap[K, V]) replaceRoot(arena *HamtPMapArena[K, V], key K, hash int32, value V) *HamtPMap[K, V] {
	bucket := HamtPMap_Bucket(m.hash)
	leafHash := HamtPMap_Shift(m.hash)
	bucketMask := jutil.Shl32(1, int64(bucket))
	index := m.Index(bucketMask)
	var newChildren []*HamtPMap[K, V]
	if (m.mask & bucketMask) != 0 {
		newChildren = replaceChild(arena, m.children, index, m.children[index].plusWithHash(arena, m.key, leafHash, m.value))
	} else {
		newChild := allocHamtPMap(arena, m.key, m.keyIsNull, leafHash, m.value, 0, HamtPMap_EmptyChildren[K, V]())
		newChildren = insertChild(arena, m.children, index, newChild)
	}
	return allocHamtPMap(arena, key, false, hash, value, m.mask|bucketMask, newChildren)
}

// Minus returns a new map with the given key removed. If the key was not
// present in the first place, then this same map will be returned.
//
// Java: HamtPMap.minus(K)
func (m *HamtPMap[K, V]) Minus(key K) PMap[K, V] {
	return m.MinusMap(key)
}

// MinusMap is minus(K) with its covariant HamtPMap return type.
//
// Go-only: language-adapter — see PlusMap.
func (m *HamtPMap[K, V]) MinusMap(key K) *HamtPMap[K, V] {
	return m.MinusMapWithArena(nil, key)
}

// MinusMapWithArena is MinusMap using arena for node and child-slice allocations.
func (m *HamtPMap[K, V]) MinusMapWithArena(arena *HamtPMapArena[K, V], key K) *HamtPMap[K, V] {
	if !m.IsEmpty() {
		res := m.minusWithHash(arena, key, HamtPMap_Hash(key), nil)
		if res != nil {
			return res
		}
		return HamtPMap_Empty[K, V]()
	}
	return m
}

// MinusWithHash is the internal recursive implementation of minus(K). The
// value of the removed node is returned via the 'value' array, if it is
// non-null.
//
// Java: HamtPMap.minus(K, int, V[])
func (m *HamtPMap[K, V]) MinusWithHash(key K, hash int32, value []V) *HamtPMap[K, V] {
	res := m.minusWithHash(nil, key, hash, value)
	if res != nil {
		return res
	}
	return HamtPMap_Empty[K, V]()
}

func (m *HamtPMap[K, V]) minusWithHash(arena *HamtPMapArena[K, V], key K, hash int32, value []V) *HamtPMap[K, V] {
	if hash == m.hash && m.keyEquals(key) {
		result := deleteRoot(arena, m.mask, m.children)
		if value != nil {
			value[0] = m.value
		}
		return result
	}
	bucket := HamtPMap_Bucket(hash)
	bucketMask := jutil.Shl32(1, int64(bucket))
	if (m.mask & bucketMask) == 0 {
		// not present, stop looking
		return m
	}
	hash = HamtPMap_Shift(hash)
	index := m.Index(bucketMask)
	child := m.children[index]
	newChild := child.minusWithHash(arena, key, hash, value)
	if newChild == child {
		return m
	} else if newChild == nil {
		return m.withChildren(arena, m.mask&^bucketMask, deleteChild(arena, m.children, index))
	} else {
		return m.withChildren(arena, m.mask, replaceChild(arena, m.children, index, newChild))
	}
}

// Java: HamtPMap.reconcile(PMap<K, V>, Reconciler<K, V>)
func (m *HamtPMap[K, V]) Reconcile(that PMap[K, V], joiner PMap_Reconciler[K, V]) PMap[K, V] {
	return m.ReconcileMap(that.(*HamtPMap[K, V]), joiner)
}

// ReconcileMap is reconcile(PMap<K, V>, Reconciler<K, V>) with its concrete
// *HamtPMap parameter and return type.
func (m *HamtPMap[K, V]) ReconcileMap(that *HamtPMap[K, V], joiner PMap_Reconciler[K, V]) *HamtPMap[K, V] {
	return m.ReconcileMapWithArena(nil, that, joiner)
}

// ReconcileMapWithArena is ReconcileMap using arena for node and child-slice allocations.
func (m *HamtPMap[K, V]) ReconcileMapWithArena(arena *HamtPMapArena[K, V], that *HamtPMap[K, V], joiner PMap_Reconciler[K, V]) *HamtPMap[K, V] {
	if m == that {
		return m
	}
	var t1, t2 *HamtPMap[K, V]
	if !m.IsEmpty() {
		t1 = m
	}
	if !that.IsEmpty() {
		t2 = that
	}
	result := hamtPMap_ReconcileTrees(arena, t1, t2, joiner, true)
	if result != nil {
		return result
	}
	if m.IsEmpty() {
		return m
	}
	if that.IsEmpty() {
		return that
	}
	return HamtPMap_Empty[K, V]()
}

// HamtPMap_ReconcileTrees is the internal recursive implementation of
// reconcile(HamtPMap, BiFunction), factoring out empies.
//
// Java: HamtPMap.reconcile(HamtPMap<K, V>, HamtPMap<K, V>, Reconciler<K, V>)
func HamtPMap_ReconcileTrees[K any, V any](t1 *HamtPMap[K, V], t2 *HamtPMap[K, V], joiner PMap_Reconciler[K, V]) *HamtPMap[K, V] {
	return hamtPMap_ReconcileTrees(nil, t1, t2, joiner, false)
}

func hamtPMap_ReconcileTrees[K any, V any](arena *HamtPMapArena[K, V], t1 *HamtPMap[K, V], t2 *HamtPMap[K, V], joiner PMap_Reconciler[K, V], checkNotNull bool) *HamtPMap[K, V] {
	var null V
	if t1 == t2 {
		return t1
	} else if t1 == nil {
		newValue := joiner(t2.key, null, t2.value)
		if checkNotNull {
			jutil.CheckNotNull(newValue)
		}
		var newChildrenBuf [16]*HamtPMap[K, V]
		n := len(t2.children)
		sameChildren := true
		for i := 0; i < n; i++ {
			c := hamtPMap_ReconcileTrees(arena, nil, t2.children[i], joiner, checkNotNull)
			newChildrenBuf[i] = c
			if c != t2.children[i] {
				sameChildren = false
			}
		}
		if sameChildren && jutil.ObjectsEqualsOf(t2.value, newValue) {
			return t2
		}
		var newChildren []*HamtPMap[K, V]
		if sameChildren {
			newChildren = t2.children
		} else {
			newChildren = allocHamtChildren(arena, n)
			copy(newChildren, newChildrenBuf[:n])
		}
		if checkNotNull || !hamtIsNull(newValue) {
			return allocHamtPMap(arena, t2.key, t2.keyIsNull, t2.hash, newValue, t2.mask, newChildren)
		}
		return deleteRoot(arena, t2.mask, newChildren)
	} else if t2 == nil {
		newValue := joiner(t1.key, t1.value, null)
		if checkNotNull {
			jutil.CheckNotNull(newValue)
		}
		var newChildrenBuf [16]*HamtPMap[K, V]
		n := len(t1.children)
		sameChildren := true
		for i := 0; i < n; i++ {
			c := hamtPMap_ReconcileTrees(arena, t1.children[i], nil, joiner, checkNotNull)
			newChildrenBuf[i] = c
			if c != t1.children[i] {
				sameChildren = false
			}
		}
		if sameChildren && jutil.ObjectsEqualsOf(t1.value, newValue) {
			return t1
		}
		var newChildren []*HamtPMap[K, V]
		if sameChildren {
			newChildren = t1.children
		} else {
			newChildren = allocHamtChildren(arena, n)
			copy(newChildren, newChildrenBuf[:n])
		}
		if checkNotNull || !hamtIsNull(newValue) {
			return allocHamtPMap(arena, t1.key, t1.keyIsNull, t1.hash, newValue, t1.mask, newChildren)
		}
		return deleteRoot(arena, t1.mask, newChildren)
	}

	// Try as hard as possible to return input trees exactly.
	sameChildrenAs1 := true
	sameChildrenAs2 := true

	// If the hashes are different, we need to keep the lower one at the top.
	hashCmp := HamtPMap_CompareUnsigned(t1.hash, t2.hash)
	key := t1.key
	keyIsNull := t1.keyIsNull
	hash := t1.hash
	var t1Val, t2Val V
	var vacatedBit1, vacatedBit2 int32
	if hashCmp < 0 {
		// t1.key is missing from t2; inline t2.VacateRoot() without allocating intermediate root/slice
		vacatedBit2 = jutil.Shl32(1, int64(HamtPMap_Bucket(t2.hash)))
		sameChildrenAs2 = false
		t1Val = t1.value
	} else if hashCmp > 0 {
		// t2.key is missing from t1; inline t1.VacateRoot() without allocating intermediate root/slice
		vacatedBit1 = jutil.Shl32(1, int64(HamtPMap_Bucket(t1.hash)))
		sameChildrenAs1 = false
		key = t2.key
		keyIsNull = t2.keyIsNull
		hash = t2.hash
		t2Val = t2.value
	} else {
		if !t2.keyEquals(t1.key) {
			// Hash collision: try to rearrange t2 to have the same root as t1
			t2 = t2.pivot(arena, t1.key, t1.hash)
			sameChildrenAs2 = false
		}
		t1Val = t1.value
		t2Val = t2.value
	}
	// Note: one or the other (but not both) tree may have a null key at root.

	var newValue V
	if vacatedBit1 == 0 && vacatedBit2 == 0 && jutil.ObjectsEqualsOf(t1Val, t2Val) {
		newValue = t1Val
	} else {
		newValue = joiner(key, t1Val, t2Val)
		if checkNotNull {
			jutil.CheckNotNull(newValue)
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
		childBit := mask & -mask // Integer.lowestOneBit
		mask &^= childBit
		child1 := t1.GetChild(childBit)
		child2 := t2.GetChild(childBit)
		var reconciled *HamtPMap[K, V]
		if childBit == vacatedBit2 {
			if (t2.mask & vacatedBit2) != 0 {
				vacatedChild2 := child2.plusWithHash(arena, t2.key, HamtPMap_Shift(t2.hash), t2.value)
				reconciled = hamtPMap_ReconcileTrees(arena, child1, vacatedChild2, joiner, checkNotNull)
			} else if child1 == nil {
				v := joiner(t2.key, null, t2.value)
				if checkNotNull {
					jutil.CheckNotNull(v)
				}
				if checkNotNull || !hamtIsNull(v) {
					reconciled = allocHamtPMap(arena, t2.key, t2.keyIsNull, HamtPMap_Shift(t2.hash), v, 0, HamtPMap_EmptyChildren[K, V]())
				}
			} else {
				vacatedChild2 := allocHamtPMap(arena, t2.key, t2.keyIsNull, HamtPMap_Shift(t2.hash), t2.value, 0, HamtPMap_EmptyChildren[K, V]())
				reconciled = hamtPMap_ReconcileTrees(arena, child1, vacatedChild2, joiner, checkNotNull)
			}
		} else if childBit == vacatedBit1 {
			if (t1.mask & vacatedBit1) != 0 {
				vacatedChild1 := child1.plusWithHash(arena, t1.key, HamtPMap_Shift(t1.hash), t1.value)
				reconciled = hamtPMap_ReconcileTrees(arena, vacatedChild1, child2, joiner, checkNotNull)
			} else if child2 == nil {
				v := joiner(t1.key, t1.value, null)
				if checkNotNull {
					jutil.CheckNotNull(v)
				}
				if checkNotNull || !hamtIsNull(v) {
					reconciled = allocHamtPMap(arena, t1.key, t1.keyIsNull, HamtPMap_Shift(t1.hash), v, 0, HamtPMap_EmptyChildren[K, V]())
				}
			} else {
				vacatedChild1 := allocHamtPMap(arena, t1.key, t1.keyIsNull, HamtPMap_Shift(t1.hash), t1.value, 0, HamtPMap_EmptyChildren[K, V]())
				reconciled = hamtPMap_ReconcileTrees(arena, vacatedChild1, child2, joiner, checkNotNull)
			}
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
	if sameChildrenAs1 && jutil.ObjectsEqualsOf(t1Val, newValue) {
		return t1
	} else if sameChildrenAs2 && jutil.ObjectsEqualsOf(t2Val, newValue) {
		return t2
	}
	var newChildren []*HamtPMap[K, V]
	if sameChildrenAs1 {
		newChildren = t1.children
	} else if sameChildrenAs2 {
		newChildren = t2.children
	} else if index == 0 {
		newChildren = HamtPMap_EmptyChildren[K, V]()
	} else {
		newChildren = allocHamtChildren(arena, index)
		copy(newChildren, newChildrenBuf[:index])
	}
	if !checkNotNull && hamtIsNull(newValue) {
		return deleteRoot(arena, newMask, newChildren)
	}
	return allocHamtPMap(arena, key, keyIsNull, hash, newValue, newMask, newChildren)
}

// Equivalent checks equality recursively based on the given equivalence.
// Short-circuits as soon as a 'false' result is found.
//
// Java: HamtPMap.equivalent(PMap<K, V>, BiPredicate<V, V>)
func (m *HamtPMap[K, V]) Equivalent(that PMap[K, V], equivalence func(V, V) bool) bool {
	return m.EquivalentMap(that.(*HamtPMap[K, V]), equivalence)
}

// EquivalentMap is equivalent(PMap<K, V>, BiPredicate<V, V>) with its concrete
// *HamtPMap parameter type.
func (m *HamtPMap[K, V]) EquivalentMap(that *HamtPMap[K, V], equivalence func(V, V) bool) bool {
	if m == that {
		return true
	}
	var t1, t2 *HamtPMap[K, V]
	if !m.IsEmpty() {
		t1 = m
	}
	if !that.IsEmpty() {
		t2 = that
	}
	return HamtPMap_EquivalentTrees(t1, t2, equivalence)
}

// HamtPMap_EquivalentTrees is the internal recursive implementation of
// equivalent(HamtPMap, BiPredicate).
//
// Java: HamtPMap.equivalent(HamtPMap<K, V>, HamtPMap<K, V>, BiPredicate<V, V>)
func HamtPMap_EquivalentTrees[K any, V any](t1 *HamtPMap[K, V], t2 *HamtPMap[K, V], equivalence func(V, V) bool) bool {
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
	} else if !t2.keyEquals(t1.key) {
		// Hash collision: try to rearrange t2 to have the same root as t1
		t2 = t2.Pivot(t1.key, t1.hash)
		if hamtIsNull(t2.value) {
			// pivot() returns a null 't2.value' when t1.key was not found in t2
			return false
		}
	}

	if !equivalence(t1.value, t2.value) {
		return false
	}
	mask := t1.mask | t2.mask
	for mask != 0 {
		childBit := mask & -mask // Integer.lowestOneBit
		mask &^= childBit
		if !HamtPMap_EquivalentTrees(t1.GetChild(childBit), t2.GetChild(childBit), equivalence) {
			return false
		}
	}
	return true
}

// Index returns the index into the 'children' array for the given bit,
// which must have exactly one bit set in its binary representation (i.e.
// must be a power of two).
//
// Java: HamtPMap.index(int)
func (m *HamtPMap[K, V]) Index(bit int32) int {
	return bits.OnesCount32(uint32(m.mask & (bit - 1)))
}

// GetChild returns the child for the given bit, which must have exactly one
// bit set. Returns null if there is no child for that bit.
//
// Java: HamtPMap.getChild(int)
func (m *HamtPMap[K, V]) GetChild(bit int32) *HamtPMap[K, V] {
	if (m.mask & bit) != 0 {
		return m.children[m.Index(bit)]
	}
	return nil
}

// HamtPMap_Hash performs the hash operation. This is done once per object as
// it enters the map, then never again since the result is passed around
// internally. We reverse the bits since certain types (such as small
// integers and short strings) have hash codes that only vary in the least
// significant bits, but we use the most significant bits for bucketing in
// order to maintain a canonical structure where the key with the smallest
// hash is always at the root of the tree.
//
// Java: HamtPMap.hash(Object)
func HamtPMap_Hash[K any](key K) int32 {
	// Go-only: performance — a String key (never null) is hashed without
	// being boxed; other keys as before.
	if s, ok := any(key).(string); ok {
		return int32(bits.Reverse32(uint32(jutil.StringHashCode(s))))
	}
	return int32(bits.Reverse32(uint32(jutil.ObjectHashCode(jutil.CheckNotNull(any(key))))))
}

// HamtPMap_Bucket returns the current bucket index from the hash.
//
// Java: HamtPMap.bucket(int)
func HamtPMap_Bucket(hash int32) int32 {
	return jutil.Ushr32(hash, HamtPMap_BITS_SHIFT)
}

// HamtPMap_Shift returns a new hash with the next bucket number shifted off.
//
// Java: HamtPMap.shift(int)
func HamtPMap_Shift(hash int32) int32 {
	return jutil.Shl32(hash, HamtPMap_BITS)
}

// HamtPMap_Unshift unshifts the bucket number back onto a hash.
//
// Java: HamtPMap.unshift(int, int)
func HamtPMap_Unshift(hash int32, bucket int32) int32 {
	return jutil.Ushr32(hash, HamtPMap_BITS) | jutil.Shl32(bucket, HamtPMap_BITS_SHIFT)
}

// HamtPMap_CompareUnsigned compares two unsigned integers.
//
// Java: HamtPMap.compareUnsigned(int, int)
func HamtPMap_CompareUnsigned(left int32, right int32) int32 {
	// NOTE: JDK 7 does not have a built-in operation for this, other than
	// casting to longs. In JDK 8 it's just Integer.compareUnsigned(left,
	// right). For now we emulate it by shifting the sign bit away, with a
	// fallback second compare only if needed.
	diff := jutil.Ushr32(left, 2) - jutil.Ushr32(right, 2)
	if diff != 0 {
		return diff
	}
	return (left & 3) - (right & 3)
}

// Pivot returns a new version of this map with the given key at the root,
// and the root element moved to some deeper node. If the key is not found in
// this map, then value will be null in the returned map.
//
// Java: HamtPMap.pivot(K, int)
func (m *HamtPMap[K, V]) Pivot(key K, hash int32) *HamtPMap[K, V] {
	return m.pivot(nil, key, hash)
}

func (m *HamtPMap[K, V]) pivot(arena *HamtPMapArena[K, V], key K, hash int32) *HamtPMap[K, V] {
	var result [1]V
	return m.pivotWithParent(arena, key, hash, nil, result[:])
}

// PivotWithParent is the internal recursive version of pivot.
//
// result is an array of size 1. The value, if found, is stored in the
// 'result' array, as a secondary return value. Otherwise the result array is
// unchanged. Returns, if parent is null, a map whose value is the value from
// the 'result' array (possibly null) and that is rooted at 'key'. Otherwise,
// a map rooted at 'parent.key'.
//
// Java: HamtPMap.pivot(K, int, HamtPMap<K, V>, V[])
func (m *HamtPMap[K, V]) PivotWithParent(key K, hash int32, parent *HamtPMap[K, V], result []V) *HamtPMap[K, V] {
	return m.pivotWithParent(nil, key, hash, parent, result)
}

func (m *HamtPMap[K, V]) pivotWithParent(arena *HamtPMapArena[K, V], key K, hash int32, parent *HamtPMap[K, V], result []V) *HamtPMap[K, V] {
	newMask := m.mask
	newChildren := m.children
	if hash == m.hash && m.keyEquals(key) {
		// Found the key: swap out this key/value with the parent and return
		// the result in the holder.
		result[0] = m.value
	} else {
		// Otherwise, see if the value might be present in a child.
		searchBucket := HamtPMap_Bucket(hash)
		replacementBucket := HamtPMap_Bucket(m.hash)

		if searchBucket == replacementBucket {
			bucketMask := jutil.Shl32(1, int64(searchBucket))
			if (m.mask & bucketMask) != 0 {
				// Found a child, call pivot recursively.
				index := m.Index(bucketMask)
				child := newChildren[index]
				newChild := child.pivotWithParent(arena, key, HamtPMap_Shift(hash), m, result)
				newChildren = replaceChild(arena, newChildren, index, newChild)
			}
		} else {
			searchMask := jutil.Shl32(1, int64(searchBucket))
			if (m.mask & searchMask) != 0 {
				index := m.Index(searchMask)
				child := newChildren[index]
				newChild := child.minusWithHash(arena, key, HamtPMap_Shift(hash), result)
				if newChild != nil {
					newChildren = replaceChild(arena, newChildren, index, newChild)
				} else {
					newChildren = deleteChild(arena, newChildren, index)
					newMask &^= searchMask
				}
			}
			replacementMask := jutil.Shl32(1, int64(replacementBucket))
			index := bits.OnesCount32(uint32(newMask & (replacementMask - 1)))
			if (m.mask & replacementMask) != 0 {
				child := newChildren[index]
				newChild := child.plusWithHash(arena, m.key, HamtPMap_Shift(m.hash), m.value)
				newChildren = replaceChild(arena, newChildren, index, newChild)
			} else {
				newChildren = insertChild(
					arena,
					newChildren,
					index,
					allocHamtPMap(arena, m.key, m.keyIsNull, HamtPMap_Shift(m.hash), m.value, 0, HamtPMap_EmptyChildren[K, V]()))
				newMask |= replacementMask
			}
		}
	}
	if parent != nil {
		return allocHamtPMap(arena, parent.key, parent.keyIsNull, HamtPMap_Shift(parent.hash), parent.value, newMask, newChildren)
	}
	return allocHamtPMap(arena, key, false, hash, result[0], newMask, newChildren)
}

// VacateRoot moves the root into the appropriate child.
//
// Java: HamtPMap.vacateRoot()
func (m *HamtPMap[K, V]) VacateRoot() *HamtPMap[K, V] {
	return m.vacateRoot(nil)
}

func (m *HamtPMap[K, V]) vacateRoot(arena *HamtPMapArena[K, V]) *HamtPMap[K, V] {
	var null V
	bucket := HamtPMap_Bucket(m.hash)
	bucketMask := jutil.Shl32(1, int64(bucket))
	index := m.Index(bucketMask)
	if (m.mask & bucketMask) != 0 {
		newChild := m.children[index].plusWithHash(arena, m.key, HamtPMap_Shift(m.hash), m.value)
		return allocHamtPMapNullKey(arena, null, m.mask, replaceChild(arena, m.children, index, newChild))
	}
	newChild := allocHamtPMap(arena, m.key, m.keyIsNull, HamtPMap_Shift(m.hash), m.value, 0, HamtPMap_EmptyChildren[K, V]())
	return allocHamtPMapNullKey(arena, null, m.mask|bucketMask, insertChild(arena, m.children, index, newChild))
}

// WithChildren returns a copy of this node with a different array of
// children.
//
// Java: HamtPMap.withChildren(int, HamtPMap<K, V>[])
func (m *HamtPMap[K, V]) WithChildren(mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	return m.withChildren(nil, mask, children)
}

func (m *HamtPMap[K, V]) withChildren(arena *HamtPMapArena[K, V], mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	if mask == m.mask && hamtSameArray(children, m.children) {
		return m
	}
	return allocHamtPMap(arena, m.key, m.keyIsNull, m.hash, m.value, mask, children)
}

// hamtSameArray is Java array identity (`a == b`) on Go slices.
func hamtSameArray[T any](a, b []T) bool {
	return len(a) == len(b) && (len(a) == 0 && cap(a) == cap(b) || len(a) > 0 && &a[0] == &b[0])
}

// HamtPMap_DeleteRoot returns a new map with the elements from children. One
// element is removed from one of the children and promoted to a root node.
// If there are no children, returns null.
//
// Java: HamtPMap.deleteRoot(int, HamtPMap<K, V>[])
func HamtPMap_DeleteRoot[K any, V any](mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	return deleteRoot(nil, mask, children)
}

func deleteRoot[K any, V any](arena *HamtPMapArena[K, V], mask int32, children []*HamtPMap[K, V]) *HamtPMap[K, V] {
	if mask == 0 {
		return nil
	}
	child := children[0]
	hashBits := int32(bits.TrailingZeros32(uint32(mask)))
	newHash := HamtPMap_Unshift(child.hash, hashBits)
	newChild := deleteRoot(arena, child.mask, child.children)
	if newChild == nil {
		newMask := mask &^ (mask & -mask) // Integer.lowestOneBit(mask)
		return allocHamtPMap(arena, child.key, child.keyIsNull, newHash, child.value, newMask, deleteChild(arena, children, 0))
	}
	return allocHamtPMap(arena, child.key, child.keyIsNull, newHash, child.value, mask, replaceChild(arena, children, 0, newChild))
}

// HamtPMap_InsertChild returns a new array of children with an additional
// child inserted at the given index.
//
// Java: HamtPMap.insertChild(HamtPMap<K, V>[], int, HamtPMap<K, V>)
func HamtPMap_InsertChild[K any, V any](children []*HamtPMap[K, V], index int, child *HamtPMap[K, V]) []*HamtPMap[K, V] {
	return insertChild(nil, children, index, child)
}

func insertChild[K any, V any](arena *HamtPMapArena[K, V], children []*HamtPMap[K, V], index int, child *HamtPMap[K, V]) []*HamtPMap[K, V] {
	newChildren := allocHamtChildren(arena, len(children)+1)
	newChildren[index] = child
	copy(newChildren[0:index], children[0:index])
	copy(newChildren[index+1:], children[index:])
	return newChildren
}

// HamtPMap_ReplaceChild returns a new array of children with the child at
// the given index replaced.
//
// Java: HamtPMap.replaceChild(HamtPMap<K, V>[], int, HamtPMap<K, V>)
func HamtPMap_ReplaceChild[K any, V any](children []*HamtPMap[K, V], index int, child *HamtPMap[K, V]) []*HamtPMap[K, V] {
	return replaceChild(nil, children, index, child)
}

func replaceChild[K any, V any](arena *HamtPMapArena[K, V], children []*HamtPMap[K, V], index int, child *HamtPMap[K, V]) []*HamtPMap[K, V] {
	newChildren := allocHamtChildren(arena, len(children))
	copy(newChildren, children)
	newChildren[index] = child
	return newChildren
}

// HamtPMap_DeleteChild returns a new array of children with the child at the
// given index deleted.
//
// Java: HamtPMap.deleteChild(HamtPMap<K, V>[], int)
func HamtPMap_DeleteChild[K any, V any](children []*HamtPMap[K, V], index int) []*HamtPMap[K, V] {
	return deleteChild(nil, children, index)
}

func deleteChild[K any, V any](arena *HamtPMapArena[K, V], children []*HamtPMap[K, V], index int) []*HamtPMap[K, V] {
	if len(children) == 1 {
		// Note: index should always be zero.
		return HamtPMap_EmptyChildren[K, V]()
	}
	newChildren := allocHamtChildren(arena, len(children)-1)
	copy(newChildren[0:index], children[0:index])
	copy(newChildren[index:], children[index+1:])
	return newChildren
}

// HamtPMap_Iter iterates sequentially over a tree.
type HamtPMap_Iter[K any, V any, O any] struct {
	queue       []*HamtPMap[K, V]
	transformer func(*HamtPMap[K, V]) O
}

// Java: HamtPMap.Iter(HamtPMap<K, V>, Function<HamtPMap<K, V>, O>)
func NewHamtPMap_Iter[K any, V any, O any](m *HamtPMap[K, V], transformer func(*HamtPMap[K, V]) O) *HamtPMap_Iter[K, V, O] {
	it := &HamtPMap_Iter[K, V, O]{transformer: transformer}
	if !m.IsEmpty() {
		it.queue = append(it.queue, m)
	}
	return it
}

// Java: HamtPMap.Iter.hasNext()
func (it *HamtPMap_Iter[K, V, O]) HasNext() bool {
	return len(it.queue) != 0
}

// Java: HamtPMap.Iter.next()
func (it *HamtPMap_Iter[K, V, O]) Next() O {
	if len(it.queue) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	top := it.queue[0]
	it.queue = it.queue[1:]
	for i := len(top.children) - 1; i >= 0; i-- {
		it.queue = append(it.queue, top.children[i])
	}
	return it.transformer(top)
}

// Java: HamtPMap.Iter.remove()
func (it *HamtPMap_Iter[K, V, O]) Remove() {
	panic(jutil.NewUnsupportedOperationException())
}

// AssertCorrectStructure throws an assertion error if the map invariant is
// violated.
//
// Java: HamtPMap.assertCorrectStructure()
func (m *HamtPMap[K, V]) AssertCorrectStructure() *HamtPMap[K, V] {
	if m.IsEmpty() {
		return m
	}
	hash := HamtPMap_Hash(m.key)
	for i := 0; i < len(m.children); i++ {
		childHash := HamtPMap_Hash(m.children[i].key)
		if HamtPMap_CompareUnsigned(childHash, hash) < 0 {
			panic(jutil.NewAssertionError("Invalid map has decreasing hash " + m.children[i].keyString() + "(" +
				jutil.IntegerToHexString(childHash) + ") beneath " + m.keyString() + "(" +
				jutil.IntegerToHexString(hash) + ": " + m.String()))
		}
		m.children[i].AssertCorrectStructure()
	}
	return m
}
