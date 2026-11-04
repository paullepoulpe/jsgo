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
	"cmp"
	"math/bits"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Specialized non-generic HAMT maps for LinkedFlowScope:
//   - flowScopeMap: *TypedScope -> *LinkedFlowScope_OverlayScope
//   - flowSlotMap:  string      -> rhino.JSType
//
// Both preserve the exact 4-bit-per-level layout, root-minimum-hash invariant,
// pivot, deleteRoot, and Reconcile traversal order of rhino.HamtPMap while
// avoiding generic `any` boxing, 7-way nil switches, and repeated UTF-16
// string hashing on hot TypeInference paths.

const (
	flowHamtBits          = 4
	flowHamtBitsShift     = 32 - flowHamtBits
	flowSlotHashCacheBits = 12
)

var flowSlotHashCache [1 << flowSlotHashCacheBits]atomic.Uint64

func flowScopeHash(s *TypedScope) uint32 {
	if s == nil {
		panic(jutil.NewNullPointerException())
	}
	h := s.id() * 0x9e3779b1
	return bits.Reverse32(h ^ (h >> 16))
}

func flowSlotHash(s string) uint32 {
	return bits.Reverse32(uint32(jutil.StringHashCode(s)))
}

func flowSlotHashByID(id uint32, s string) uint32 {
	if id != 0 {
		idx := (id * 0x9e3779b1) >> (32 - flowSlotHashCacheBits)
		entry := flowSlotHashCache[idx].Load()
		if uint32(entry>>32) == id {
			return uint32(entry)
		}
		h := bits.Reverse32(uint32(jutil.StringHashCode(s)))
		flowSlotHashCache[idx].Store((uint64(id) << 32) | uint64(h))
		return h
	}
	return bits.Reverse32(uint32(jutil.StringHashCode(s)))
}

func flowSlotValueEquals(a, b rhino.JSType) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil || a == linkedFlowNullType || b == linkedFlowNullType {
		return false
	}
	ab := a.GetJSTypeBase()
	bb := b.GetJSTypeBase()
	if ab != nil && ab == bb {
		return true
	}
	return a.Equals(b)
}

func flowSameSlice[T any](a, b []T) bool {
	return len(a) == len(b) && (len(a) == 0 || &a[0] == &b[0])
}

// --- flowScopeMap (*TypedScope -> *LinkedFlowScope_OverlayScope) ---

type flowScopeMap struct {
	key      *TypedScope
	hash     uint32
	value    *LinkedFlowScope_OverlayScope
	mask     uint32
	children []*flowScopeMap
}

type flowScopeMapArena struct {
	nodes       jutil.ReuseSlab[flowScopeMap]
	childBlocks jutil.ReuseSlab[[32]*flowScopeMap]
	childSlice  []*flowScopeMap
}

func (a *flowScopeMapArena) Reset() {
	a.nodes.Reset(nil)
	a.childBlocks.Reset(nil)
	a.childSlice = nil
}

func (a *flowScopeMapArena) TrimLargeChunks() {
	a.Reset()
}

func allocFlowScopeChildren(a *flowScopeMapArena, n int) []*flowScopeMap {
	if n == 0 {
		return nil
	}
	if a == nil || n > 16 {
		return make([]*flowScopeMap, n)
	}
	if len(a.childSlice) < n {
		a.childSlice = a.childBlocks.Alloc()[:]
	}
	s := a.childSlice[:n:n]
	a.childSlice = a.childSlice[n:]
	return s
}

func allocFlowScopeMap(
	a *flowScopeMapArena,
	key *TypedScope,
	hash uint32,
	value *LinkedFlowScope_OverlayScope,
	mask uint32,
	children []*flowScopeMap,
) *flowScopeMap {
	var m *flowScopeMap
	if a != nil {
		m = a.nodes.Alloc()
	} else {
		m = &flowScopeMap{}
	}
	m.key = key
	m.hash = hash
	m.value = value
	m.mask = mask
	m.children = children
	return m
}

func (m *flowScopeMap) IsEmpty() bool { return m == nil }

func (m *flowScopeMap) index(bit uint32) int {
	return bits.OnesCount32(m.mask & (bit - 1))
}

func (m *flowScopeMap) getChild(bit uint32) *flowScopeMap {
	if (m.mask & bit) != 0 {
		return m.children[m.index(bit)]
	}
	return nil
}

func (m *flowScopeMap) Get(key *TypedScope) *LinkedFlowScope_OverlayScope {
	if m == nil {
		return nil
	}
	hash := flowScopeHash(key)
	for cur := m; cur != nil; {
		if hash == cur.hash && key == cur.key {
			return cur.value
		}
		bucketMask := uint32(1) << (hash >> flowHamtBitsShift)
		if (cur.mask & bucketMask) == 0 {
			return nil
		}
		cur = cur.children[cur.index(bucketMask)]
		hash <<= flowHamtBits
	}
	return nil
}

func (m *flowScopeMap) PlusMapWithArena(
	arena *flowScopeMapArena, key *TypedScope, value *LinkedFlowScope_OverlayScope,
) *flowScopeMap {
	if value == nil {
		panic(jutil.NewNullPointerException())
	}
	hash := flowScopeHash(key)
	if m != nil {
		return m.plusWithHash(arena, key, hash, value)
	}
	return allocFlowScopeMap(arena, key, hash, value, 0, nil)
}

func (m *flowScopeMap) plusWithHash(
	arena *flowScopeMapArena, key *TypedScope, hash uint32, value *LinkedFlowScope_OverlayScope,
) *flowScopeMap {
	if hash == m.hash && key == m.key {
		if value == m.value {
			return m
		}
		return allocFlowScopeMap(arena, key, hash, value, m.mask, m.children)
	}
	if cmp.Compare(hash, m.hash) < 0 {
		return m.replaceRoot(arena, key, hash, value)
	}
	bucketMask := uint32(1) << (hash >> flowHamtBitsShift)
	hash <<= flowHamtBits
	idx := m.index(bucketMask)
	if (m.mask & bucketMask) != 0 {
		child := m.children[idx]
		newChild := child.plusWithHash(arena, key, hash, value)
		if child == newChild {
			return m
		}
		return m.withChildren(arena, m.mask, replaceFlowScopeChild(arena, m.children, idx, newChild))
	}
	newChild := allocFlowScopeMap(arena, key, hash, value, 0, nil)
	return m.withChildren(arena, m.mask|bucketMask, insertFlowScopeChild(arena, m.children, idx, newChild))
}

func (m *flowScopeMap) replaceRoot(
	arena *flowScopeMapArena, key *TypedScope, hash uint32, value *LinkedFlowScope_OverlayScope,
) *flowScopeMap {
	bucketMask := uint32(1) << (m.hash >> flowHamtBitsShift)
	leafHash := m.hash << flowHamtBits
	idx := m.index(bucketMask)
	var newChildren []*flowScopeMap
	if (m.mask & bucketMask) != 0 {
		newChildren = replaceFlowScopeChild(
			arena, m.children, idx, m.children[idx].plusWithHash(arena, m.key, leafHash, m.value))
	} else {
		newChild := allocFlowScopeMap(arena, m.key, leafHash, m.value, 0, nil)
		newChildren = insertFlowScopeChild(arena, m.children, idx, newChild)
	}
	return allocFlowScopeMap(arena, key, hash, value, m.mask|bucketMask, newChildren)
}

func (m *flowScopeMap) MinusMapWithArena(arena *flowScopeMapArena, key *TypedScope) *flowScopeMap {
	if m != nil {
		res, _, _ := m.minusWithHash(arena, key, flowScopeHash(key))
		return res
	}
	return nil
}

func (m *flowScopeMap) minusWithHash(
	arena *flowScopeMapArena, key *TypedScope, hash uint32,
) (*flowScopeMap, *LinkedFlowScope_OverlayScope, bool) {
	if hash == m.hash && key == m.key {
		return deleteFlowScopeRoot(arena, m.mask, m.children), m.value, true
	}
	bucketMask := uint32(1) << (hash >> flowHamtBitsShift)
	if (m.mask & bucketMask) == 0 {
		return m, nil, false
	}
	hash <<= flowHamtBits
	idx := m.index(bucketMask)
	child := m.children[idx]
	newChild, removed, found := child.minusWithHash(arena, key, hash)
	if !found {
		return m, nil, false
	}
	if newChild == nil {
		return m.withChildren(arena, m.mask&^bucketMask, deleteFlowScopeChild(arena, m.children, idx)), removed, true
	}
	return m.withChildren(arena, m.mask, replaceFlowScopeChild(arena, m.children, idx, newChild)), removed, true
}

func (m *flowScopeMap) ReconcileWithState(
	arena *flowScopeMapArena, that *flowScopeMap, st *linkedFlowJoinState,
) *flowScopeMap {
	if m == that {
		return m
	}
	return reconcileFlowScopeTrees(arena, m, that, st)
}

func reconcileSingleFlowScopeTree(
	arena *flowScopeMapArena, t *flowScopeMap, onLeft bool, st *linkedFlowJoinState,
) *flowScopeMap {
	var newValue *LinkedFlowScope_OverlayScope
	if onLeft {
		newValue = st.reconcileScope(t.key, t.value, nil)
	} else {
		newValue = st.reconcileScope(t.key, nil, t.value)
	}
	if newValue == nil {
		panic(jutil.NewNullPointerException())
	}
	var newChildrenBuf [16]*flowScopeMap
	n := len(t.children)
	sameChildren := true
	for i := 0; i < n; i++ {
		c := reconcileSingleFlowScopeTree(arena, t.children[i], onLeft, st)
		newChildrenBuf[i] = c
		if c != t.children[i] {
			sameChildren = false
		}
	}
	if sameChildren && t.value == newValue {
		return t
	}
	var newChildren []*flowScopeMap
	if sameChildren {
		newChildren = t.children
	} else {
		newChildren = allocFlowScopeChildren(arena, n)
		copy(newChildren, newChildrenBuf[:n])
	}
	return allocFlowScopeMap(arena, t.key, t.hash, newValue, t.mask, newChildren)
}

func reconcileVacatedFlowScopeChild(
	arena *flowScopeMapArena,
	vacated, ownChild, otherChild *flowScopeMap,
	vacatedOnLeft bool,
	st *linkedFlowJoinState,
) *flowScopeMap {
	shiftedHash := vacated.hash << flowHamtBits
	if ownChild == nil && otherChild == nil {
		var v *LinkedFlowScope_OverlayScope
		if vacatedOnLeft {
			v = st.reconcileScope(vacated.key, vacated.value, nil)
		} else {
			v = st.reconcileScope(vacated.key, nil, vacated.value)
		}
		if v == nil {
			panic(jutil.NewNullPointerException())
		}
		return allocFlowScopeMap(arena, vacated.key, shiftedHash, v, 0, nil)
	}
	var vacatedChild *flowScopeMap
	if ownChild != nil {
		vacatedChild = ownChild.plusWithHash(arena, vacated.key, shiftedHash, vacated.value)
	} else {
		vacatedChild = allocFlowScopeMap(arena, vacated.key, shiftedHash, vacated.value, 0, nil)
	}
	if vacatedOnLeft {
		return reconcileFlowScopeTrees(arena, vacatedChild, otherChild, st)
	}
	return reconcileFlowScopeTrees(arena, otherChild, vacatedChild, st)
}

func reconcileFlowScopeTrees(
	arena *flowScopeMapArena, t1, t2 *flowScopeMap, st *linkedFlowJoinState,
) *flowScopeMap {
	if t1 == t2 {
		return t1
	} else if t1 == nil {
		return reconcileSingleFlowScopeTree(arena, t2, false, st)
	} else if t2 == nil {
		return reconcileSingleFlowScopeTree(arena, t1, true, st)
	}

	sameChildrenAs1 := true
	sameChildrenAs2 := true

	hashCmp := cmp.Compare(t1.hash, t2.hash)
	key := t1.key
	hash := t1.hash
	var t1Val, t2Val *LinkedFlowScope_OverlayScope
	var vacatedBit1, vacatedBit2 uint32
	bothPresent := false
	if hashCmp < 0 {
		vacatedBit2 = uint32(1) << (t2.hash >> flowHamtBitsShift)
		sameChildrenAs2 = false
		t1Val = t1.value
	} else if hashCmp > 0 {
		vacatedBit1 = uint32(1) << (t1.hash >> flowHamtBitsShift)
		sameChildrenAs1 = false
		key = t2.key
		hash = t2.hash
		t2Val = t2.value
	} else {
		if t1.key != t2.key {
			t2, bothPresent = t2.pivot(arena, t1.key, t1.hash)
			sameChildrenAs2 = false
		} else {
			bothPresent = true
		}
		t1Val = t1.value
		t2Val = t2.value
	}

	var newValue *LinkedFlowScope_OverlayScope
	if bothPresent && t1Val == t2Val {
		newValue = t1Val
	} else {
		newValue = st.reconcileScope(key, t1Val, t2Val)
		if newValue == nil {
			panic(jutil.NewNullPointerException())
		}
	}

	t1Mask := t1.mask | vacatedBit1
	t2Mask := t2.mask | vacatedBit2
	newMask := t1Mask | t2Mask
	sameChildrenAs1 = sameChildrenAs1 && (newMask == t1Mask)
	sameChildrenAs2 = sameChildrenAs2 && (newMask == t2Mask)

	var newChildrenBuf [16]*flowScopeMap
	mask := newMask
	idx := 0
	for mask != 0 {
		childBit := mask & -mask
		mask &^= childBit
		child1 := t1.getChild(childBit)
		child2 := t2.getChild(childBit)
		var reconciled *flowScopeMap
		if childBit == vacatedBit2 {
			reconciled = reconcileVacatedFlowScopeChild(arena, t2, child2, child1, false, st)
		} else if childBit == vacatedBit1 {
			reconciled = reconcileVacatedFlowScopeChild(arena, t1, child1, child2, true, st)
		} else {
			reconciled = reconcileFlowScopeTrees(arena, child1, child2, st)
		}
		newChildrenBuf[idx] = reconciled
		sameChildrenAs1 = sameChildrenAs1 && (reconciled == child1)
		sameChildrenAs2 = sameChildrenAs2 && (reconciled == child2)
		if reconciled != nil {
			idx++
		} else {
			newMask &^= childBit
		}
	}
	if sameChildrenAs1 && t1Val == newValue {
		return t1
	} else if sameChildrenAs2 && t2Val == newValue {
		return t2
	}
	var newChildren []*flowScopeMap
	if sameChildrenAs1 {
		newChildren = t1.children
	} else if sameChildrenAs2 {
		newChildren = t2.children
	} else if idx > 0 {
		newChildren = allocFlowScopeChildren(arena, idx)
		copy(newChildren, newChildrenBuf[:idx])
	}
	return allocFlowScopeMap(arena, key, hash, newValue, newMask, newChildren)
}

func (m *flowScopeMap) Equivalent(that *flowScopeMap) bool {
	if m == that {
		return true
	} else if m == nil || that == nil {
		return false
	}
	if m.hash != that.hash {
		return false
	} else if m.key != that.key {
		var ok bool
		that, ok = that.pivot(nil, m.key, m.hash)
		if !ok {
			return false
		}
	}
	if !LinkedFlowScope_EqualScopes(m.value, that.value) {
		return false
	}
	mask := m.mask | that.mask
	for mask != 0 {
		childBit := mask & -mask
		mask &^= childBit
		if !m.getChild(childBit).Equivalent(that.getChild(childBit)) {
			return false
		}
	}
	return true
}

func (m *flowScopeMap) pivot(arena *flowScopeMapArena, key *TypedScope, hash uint32) (*flowScopeMap, bool) {
	res, _, ok := m.pivotWithParent(arena, key, hash, nil)
	return res, ok
}

func (m *flowScopeMap) pivotWithParent(
	arena *flowScopeMapArena, key *TypedScope, hash uint32, parent *flowScopeMap,
) (*flowScopeMap, *LinkedFlowScope_OverlayScope, bool) {
	var found *LinkedFlowScope_OverlayScope
	var ok bool
	newMask := m.mask
	newChildren := m.children
	if hash == m.hash && key == m.key {
		found = m.value
		ok = true
	} else {
		searchBucket := hash >> flowHamtBitsShift
		replacementBucket := m.hash >> flowHamtBitsShift
		if searchBucket == replacementBucket {
			bucketMask := uint32(1) << searchBucket
			if (m.mask & bucketMask) != 0 {
				idx := m.index(bucketMask)
				child := newChildren[idx]
				var newChild *flowScopeMap
				newChild, found, ok = child.pivotWithParent(arena, key, hash<<flowHamtBits, m)
				newChildren = replaceFlowScopeChild(arena, newChildren, idx, newChild)
			}
		} else {
			searchMask := uint32(1) << searchBucket
			if (m.mask & searchMask) != 0 {
				idx := m.index(searchMask)
				child := newChildren[idx]
				var newChild *flowScopeMap
				newChild, found, ok = child.minusWithHash(arena, key, hash<<flowHamtBits)
				if newChild != nil {
					newChildren = replaceFlowScopeChild(arena, newChildren, idx, newChild)
				} else {
					newChildren = deleteFlowScopeChild(arena, newChildren, idx)
					newMask &^= searchMask
				}
			}
			replacementMask := uint32(1) << replacementBucket
			idx := bits.OnesCount32(newMask & (replacementMask - 1))
			if (m.mask & replacementMask) != 0 {
				child := newChildren[idx]
				newChild := child.plusWithHash(arena, m.key, m.hash<<flowHamtBits, m.value)
				newChildren = replaceFlowScopeChild(arena, newChildren, idx, newChild)
			} else {
				newChildren = insertFlowScopeChild(
					arena, newChildren, idx,
					allocFlowScopeMap(arena, m.key, m.hash<<flowHamtBits, m.value, 0, nil))
				newMask |= replacementMask
			}
		}
	}
	if parent != nil {
		return allocFlowScopeMap(arena, parent.key, parent.hash<<flowHamtBits, parent.value, newMask, newChildren), found, ok
	}
	return allocFlowScopeMap(arena, key, hash, found, newMask, newChildren), found, ok
}

func (m *flowScopeMap) withChildren(arena *flowScopeMapArena, mask uint32, children []*flowScopeMap) *flowScopeMap {
	if mask == m.mask && flowSameSlice(children, m.children) {
		return m
	}
	return allocFlowScopeMap(arena, m.key, m.hash, m.value, mask, children)
}

func deleteFlowScopeRoot(arena *flowScopeMapArena, mask uint32, children []*flowScopeMap) *flowScopeMap {
	if mask == 0 {
		return nil
	}
	child := children[0]
	hashBits := uint32(bits.TrailingZeros32(mask))
	newHash := (child.hash >> flowHamtBits) | (hashBits << flowHamtBitsShift)
	newChild := deleteFlowScopeRoot(arena, child.mask, child.children)
	if newChild == nil {
		newMask := mask &^ (mask & -mask)
		return allocFlowScopeMap(arena, child.key, newHash, child.value, newMask, deleteFlowScopeChild(arena, children, 0))
	}
	return allocFlowScopeMap(arena, child.key, newHash, child.value, mask, replaceFlowScopeChild(arena, children, 0, newChild))
}

func insertFlowScopeChild(arena *flowScopeMapArena, children []*flowScopeMap, index int, child *flowScopeMap) []*flowScopeMap {
	newChildren := allocFlowScopeChildren(arena, len(children)+1)
	newChildren[index] = child
	copy(newChildren[:index], children[:index])
	copy(newChildren[index+1:], children[index:])
	return newChildren
}

func replaceFlowScopeChild(arena *flowScopeMapArena, children []*flowScopeMap, index int, child *flowScopeMap) []*flowScopeMap {
	newChildren := allocFlowScopeChildren(arena, len(children))
	copy(newChildren, children)
	newChildren[index] = child
	return newChildren
}

func deleteFlowScopeChild(arena *flowScopeMapArena, children []*flowScopeMap, index int) []*flowScopeMap {
	if len(children) == 1 {
		return nil
	}
	newChildren := allocFlowScopeChildren(arena, len(children)-1)
	copy(newChildren[:index], children[:index])
	copy(newChildren[index:], children[index+1:])
	return newChildren
}

// --- flowSlotMap (string -> rhino.JSType) ---

type flowSlotMap struct {
	key      string
	hash     uint32
	value    rhino.JSType
	mask     uint32
	children []*flowSlotMap
}

type flowSlotMapArena struct {
	nodes       jutil.ReuseSlab[flowSlotMap]
	childBlocks jutil.ReuseSlab[[32]*flowSlotMap]
	childSlice  []*flowSlotMap
}

func (a *flowSlotMapArena) Reset() {
	a.nodes.Reset(nil)
	a.childBlocks.Reset(nil)
	a.childSlice = nil
}

func (a *flowSlotMapArena) TrimLargeChunks() {
	a.Reset()
}

func allocFlowSlotChildren(a *flowSlotMapArena, n int) []*flowSlotMap {
	if n == 0 {
		return nil
	}
	if a == nil || n > 16 {
		return make([]*flowSlotMap, n)
	}
	if len(a.childSlice) < n {
		a.childSlice = a.childBlocks.Alloc()[:]
	}
	s := a.childSlice[:n:n]
	a.childSlice = a.childSlice[n:]
	return s
}

func allocFlowSlotMap(
	a *flowSlotMapArena,
	key string,
	hash uint32,
	value rhino.JSType,
	mask uint32,
	children []*flowSlotMap,
) *flowSlotMap {
	var m *flowSlotMap
	if a != nil {
		m = a.nodes.Alloc()
	} else {
		m = &flowSlotMap{}
	}
	m.key = key
	m.hash = hash
	m.value = value
	m.mask = mask
	m.children = children
	return m
}

func (m *flowSlotMap) IsEmpty() bool { return m == nil }

func (m *flowSlotMap) index(bit uint32) int {
	return bits.OnesCount32(m.mask & (bit - 1))
}

func (m *flowSlotMap) getChild(bit uint32) *flowSlotMap {
	if (m.mask & bit) != 0 {
		return m.children[m.index(bit)]
	}
	return nil
}

func (m *flowSlotMap) Get(key string) rhino.JSType {
	if m == nil {
		return nil
	}
	return m.getWithHash(key, flowSlotHash(key))
}

func (m *flowSlotMap) GetByID(id uint32, key string) rhino.JSType {
	if m == nil {
		return nil
	}
	return m.getWithHash(key, flowSlotHashByID(id, key))
}

func (m *flowSlotMap) getWithHash(key string, hash uint32) rhino.JSType {
	for cur := m; cur != nil; {
		if hash == cur.hash && key == cur.key {
			return cur.value
		}
		bucketMask := uint32(1) << (hash >> flowHamtBitsShift)
		if (cur.mask & bucketMask) == 0 {
			return nil
		}
		cur = cur.children[cur.index(bucketMask)]
		hash <<= flowHamtBits
	}
	return nil
}

func (m *flowSlotMap) PlusMapWithArena(
	arena *flowSlotMapArena, key string, value rhino.JSType,
) *flowSlotMap {
	return m.PlusMapWithHash(arena, key, flowSlotHash(key), value)
}

func (m *flowSlotMap) PlusMapWithHash(
	arena *flowSlotMapArena, key string, hash uint32, value rhino.JSType,
) *flowSlotMap {
	if value == nil {
		panic(jutil.NewNullPointerException())
	}
	if m != nil {
		return m.plusWithHash(arena, key, hash, value)
	}
	return allocFlowSlotMap(arena, key, hash, value, 0, nil)
}

func (m *flowSlotMap) plusWithHash(
	arena *flowSlotMapArena, key string, hash uint32, value rhino.JSType,
) *flowSlotMap {
	if hash == m.hash && key == m.key {
		if flowSlotValueEquals(value, m.value) {
			return m
		}
		return allocFlowSlotMap(arena, key, hash, value, m.mask, m.children)
	}
	if cmp.Compare(hash, m.hash) < 0 {
		return m.replaceRoot(arena, key, hash, value)
	}
	bucketMask := uint32(1) << (hash >> flowHamtBitsShift)
	hash <<= flowHamtBits
	idx := m.index(bucketMask)
	if (m.mask & bucketMask) != 0 {
		child := m.children[idx]
		newChild := child.plusWithHash(arena, key, hash, value)
		if child == newChild {
			return m
		}
		return m.withChildren(arena, m.mask, replaceFlowSlotChild(arena, m.children, idx, newChild))
	}
	newChild := allocFlowSlotMap(arena, key, hash, value, 0, nil)
	return m.withChildren(arena, m.mask|bucketMask, insertFlowSlotChild(arena, m.children, idx, newChild))
}

func (m *flowSlotMap) replaceRoot(
	arena *flowSlotMapArena, key string, hash uint32, value rhino.JSType,
) *flowSlotMap {
	bucketMask := uint32(1) << (m.hash >> flowHamtBitsShift)
	leafHash := m.hash << flowHamtBits
	idx := m.index(bucketMask)
	var newChildren []*flowSlotMap
	if (m.mask & bucketMask) != 0 {
		newChildren = replaceFlowSlotChild(
			arena, m.children, idx, m.children[idx].plusWithHash(arena, m.key, leafHash, m.value))
	} else {
		newChild := allocFlowSlotMap(arena, m.key, leafHash, m.value, 0, nil)
		newChildren = insertFlowSlotChild(arena, m.children, idx, newChild)
	}
	return allocFlowSlotMap(arena, key, hash, value, m.mask|bucketMask, newChildren)
}

func (m *flowSlotMap) minusWithHash(
	arena *flowSlotMapArena, key string, hash uint32,
) (*flowSlotMap, rhino.JSType, bool) {
	if hash == m.hash && key == m.key {
		return deleteFlowSlotRoot(arena, m.mask, m.children), m.value, true
	}
	bucketMask := uint32(1) << (hash >> flowHamtBitsShift)
	if (m.mask & bucketMask) == 0 {
		return m, nil, false
	}
	hash <<= flowHamtBits
	idx := m.index(bucketMask)
	child := m.children[idx]
	newChild, removed, found := child.minusWithHash(arena, key, hash)
	if !found {
		return m, nil, false
	}
	if newChild == nil {
		return m.withChildren(arena, m.mask&^bucketMask, deleteFlowSlotChild(arena, m.children, idx)), removed, true
	}
	return m.withChildren(arena, m.mask, replaceFlowSlotChild(arena, m.children, idx, newChild)), removed, true
}

func (m *flowSlotMap) ReconcileWithState(
	arena *flowSlotMapArena, that *flowSlotMap, st *linkedFlowJoinState,
) *flowSlotMap {
	if m == that {
		return m
	}
	return reconcileFlowSlotTrees(arena, m, that, st)
}

func reconcileSingleFlowSlotTree(
	arena *flowSlotMapArena, t *flowSlotMap, onLeft bool, st *linkedFlowJoinState,
) *flowSlotMap {
	var newValue rhino.JSType
	if onLeft {
		newValue = st.reconcileSlotEntry(t.key, t.value, nil)
	} else {
		newValue = st.reconcileSlotEntry(t.key, nil, t.value)
	}
	if newValue == nil {
		panic(jutil.NewNullPointerException())
	}
	var newChildrenBuf [16]*flowSlotMap
	n := len(t.children)
	sameChildren := true
	for i := 0; i < n; i++ {
		c := reconcileSingleFlowSlotTree(arena, t.children[i], onLeft, st)
		newChildrenBuf[i] = c
		if c != t.children[i] {
			sameChildren = false
		}
	}
	if sameChildren && flowSlotValueEquals(t.value, newValue) {
		return t
	}
	var newChildren []*flowSlotMap
	if sameChildren {
		newChildren = t.children
	} else {
		newChildren = allocFlowSlotChildren(arena, n)
		copy(newChildren, newChildrenBuf[:n])
	}
	return allocFlowSlotMap(arena, t.key, t.hash, newValue, t.mask, newChildren)
}

func reconcileVacatedFlowSlotChild(
	arena *flowSlotMapArena,
	vacated, ownChild, otherChild *flowSlotMap,
	vacatedOnLeft bool,
	st *linkedFlowJoinState,
) *flowSlotMap {
	shiftedHash := vacated.hash << flowHamtBits
	if ownChild == nil && otherChild == nil {
		var v rhino.JSType
		if vacatedOnLeft {
			v = st.reconcileSlotEntry(vacated.key, vacated.value, nil)
		} else {
			v = st.reconcileSlotEntry(vacated.key, nil, vacated.value)
		}
		if v == nil {
			panic(jutil.NewNullPointerException())
		}
		return allocFlowSlotMap(arena, vacated.key, shiftedHash, v, 0, nil)
	}
	var vacatedChild *flowSlotMap
	if ownChild != nil {
		vacatedChild = ownChild.plusWithHash(arena, vacated.key, shiftedHash, vacated.value)
	} else {
		vacatedChild = allocFlowSlotMap(arena, vacated.key, shiftedHash, vacated.value, 0, nil)
	}
	if vacatedOnLeft {
		return reconcileFlowSlotTrees(arena, vacatedChild, otherChild, st)
	}
	return reconcileFlowSlotTrees(arena, otherChild, vacatedChild, st)
}

func reconcileFlowSlotTrees(
	arena *flowSlotMapArena, t1, t2 *flowSlotMap, st *linkedFlowJoinState,
) *flowSlotMap {
	if t1 == t2 {
		return t1
	} else if t1 == nil {
		return reconcileSingleFlowSlotTree(arena, t2, false, st)
	} else if t2 == nil {
		return reconcileSingleFlowSlotTree(arena, t1, true, st)
	}

	sameChildrenAs1 := true
	sameChildrenAs2 := true

	hashCmp := cmp.Compare(t1.hash, t2.hash)
	key := t1.key
	hash := t1.hash
	var t1Val, t2Val rhino.JSType
	var vacatedBit1, vacatedBit2 uint32
	bothPresent := false
	if hashCmp < 0 {
		vacatedBit2 = uint32(1) << (t2.hash >> flowHamtBitsShift)
		sameChildrenAs2 = false
		t1Val = t1.value
	} else if hashCmp > 0 {
		vacatedBit1 = uint32(1) << (t1.hash >> flowHamtBitsShift)
		sameChildrenAs1 = false
		key = t2.key
		hash = t2.hash
		t2Val = t2.value
	} else {
		if t1.key != t2.key {
			t2, bothPresent = t2.pivot(arena, t1.key, t1.hash)
			sameChildrenAs2 = false
		} else {
			bothPresent = true
		}
		t1Val = t1.value
		t2Val = t2.value
	}

	var newValue rhino.JSType
	if bothPresent && flowSlotValueEquals(t1Val, t2Val) {
		newValue = t1Val
	} else {
		newValue = st.reconcileSlotEntry(key, t1Val, t2Val)
		if newValue == nil {
			panic(jutil.NewNullPointerException())
		}
	}

	t1Mask := t1.mask | vacatedBit1
	t2Mask := t2.mask | vacatedBit2
	newMask := t1Mask | t2Mask
	sameChildrenAs1 = sameChildrenAs1 && (newMask == t1Mask)
	sameChildrenAs2 = sameChildrenAs2 && (newMask == t2Mask)

	var newChildrenBuf [16]*flowSlotMap
	mask := newMask
	idx := 0
	for mask != 0 {
		childBit := mask & -mask
		mask &^= childBit
		child1 := t1.getChild(childBit)
		child2 := t2.getChild(childBit)
		var reconciled *flowSlotMap
		if childBit == vacatedBit2 {
			reconciled = reconcileVacatedFlowSlotChild(arena, t2, child2, child1, false, st)
		} else if childBit == vacatedBit1 {
			reconciled = reconcileVacatedFlowSlotChild(arena, t1, child1, child2, true, st)
		} else {
			reconciled = reconcileFlowSlotTrees(arena, child1, child2, st)
		}
		newChildrenBuf[idx] = reconciled
		sameChildrenAs1 = sameChildrenAs1 && (reconciled == child1)
		sameChildrenAs2 = sameChildrenAs2 && (reconciled == child2)
		if reconciled != nil {
			idx++
		} else {
			newMask &^= childBit
		}
	}
	if sameChildrenAs1 && flowSlotValueEquals(t1Val, newValue) {
		return t1
	} else if sameChildrenAs2 && flowSlotValueEquals(t2Val, newValue) {
		return t2
	}
	var newChildren []*flowSlotMap
	if sameChildrenAs1 {
		newChildren = t1.children
	} else if sameChildrenAs2 {
		newChildren = t2.children
	} else if idx > 0 {
		newChildren = allocFlowSlotChildren(arena, idx)
		copy(newChildren, newChildrenBuf[:idx])
	}
	return allocFlowSlotMap(arena, key, hash, newValue, newMask, newChildren)
}

func (m *flowSlotMap) Equivalent(that *flowSlotMap) bool {
	if m == that {
		return true
	} else if m == nil || that == nil {
		return false
	}
	if m.hash != that.hash {
		return false
	} else if m.key != that.key {
		var ok bool
		that, ok = that.pivot(nil, m.key, m.hash)
		if !ok {
			return false
		}
	}
	if !linkedFlowScope_EqualTypes(m.value, that.value) {
		return false
	}
	mask := m.mask | that.mask
	for mask != 0 {
		childBit := mask & -mask
		mask &^= childBit
		if !m.getChild(childBit).Equivalent(that.getChild(childBit)) {
			return false
		}
	}
	return true
}

func (m *flowSlotMap) pivot(arena *flowSlotMapArena, key string, hash uint32) (*flowSlotMap, bool) {
	res, _, ok := m.pivotWithParent(arena, key, hash, nil)
	return res, ok
}

func (m *flowSlotMap) pivotWithParent(
	arena *flowSlotMapArena, key string, hash uint32, parent *flowSlotMap,
) (*flowSlotMap, rhino.JSType, bool) {
	var found rhino.JSType
	var ok bool
	newMask := m.mask
	newChildren := m.children
	if hash == m.hash && key == m.key {
		found = m.value
		ok = true
	} else {
		searchBucket := hash >> flowHamtBitsShift
		replacementBucket := m.hash >> flowHamtBitsShift
		if searchBucket == replacementBucket {
			bucketMask := uint32(1) << searchBucket
			if (m.mask & bucketMask) != 0 {
				idx := m.index(bucketMask)
				child := newChildren[idx]
				var newChild *flowSlotMap
				newChild, found, ok = child.pivotWithParent(arena, key, hash<<flowHamtBits, m)
				newChildren = replaceFlowSlotChild(arena, newChildren, idx, newChild)
			}
		} else {
			searchMask := uint32(1) << searchBucket
			if (m.mask & searchMask) != 0 {
				idx := m.index(searchMask)
				child := newChildren[idx]
				var newChild *flowSlotMap
				newChild, found, ok = child.minusWithHash(arena, key, hash<<flowHamtBits)
				if newChild != nil {
					newChildren = replaceFlowSlotChild(arena, newChildren, idx, newChild)
				} else {
					newChildren = deleteFlowSlotChild(arena, newChildren, idx)
					newMask &^= searchMask
				}
			}
			replacementMask := uint32(1) << replacementBucket
			idx := bits.OnesCount32(newMask & (replacementMask - 1))
			if (m.mask & replacementMask) != 0 {
				child := newChildren[idx]
				newChild := child.plusWithHash(arena, m.key, m.hash<<flowHamtBits, m.value)
				newChildren = replaceFlowSlotChild(arena, newChildren, idx, newChild)
			} else {
				newChildren = insertFlowSlotChild(
					arena, newChildren, idx,
					allocFlowSlotMap(arena, m.key, m.hash<<flowHamtBits, m.value, 0, nil))
				newMask |= replacementMask
			}
		}
	}
	if parent != nil {
		return allocFlowSlotMap(arena, parent.key, parent.hash<<flowHamtBits, parent.value, newMask, newChildren), found, ok
	}
	return allocFlowSlotMap(arena, key, hash, found, newMask, newChildren), found, ok
}

func (m *flowSlotMap) withChildren(arena *flowSlotMapArena, mask uint32, children []*flowSlotMap) *flowSlotMap {
	if mask == m.mask && flowSameSlice(children, m.children) {
		return m
	}
	return allocFlowSlotMap(arena, m.key, m.hash, m.value, mask, children)
}

func deleteFlowSlotRoot(arena *flowSlotMapArena, mask uint32, children []*flowSlotMap) *flowSlotMap {
	if mask == 0 {
		return nil
	}
	child := children[0]
	hashBits := uint32(bits.TrailingZeros32(mask))
	newHash := (child.hash >> flowHamtBits) | (hashBits << flowHamtBitsShift)
	newChild := deleteFlowSlotRoot(arena, child.mask, child.children)
	if newChild == nil {
		newMask := mask &^ (mask & -mask)
		return allocFlowSlotMap(arena, child.key, newHash, child.value, newMask, deleteFlowSlotChild(arena, children, 0))
	}
	return allocFlowSlotMap(arena, child.key, newHash, child.value, mask, replaceFlowSlotChild(arena, children, 0, newChild))
}

func insertFlowSlotChild(arena *flowSlotMapArena, children []*flowSlotMap, index int, child *flowSlotMap) []*flowSlotMap {
	newChildren := allocFlowSlotChildren(arena, len(children)+1)
	newChildren[index] = child
	copy(newChildren[:index], children[:index])
	copy(newChildren[index+1:], children[index:])
	return newChildren
}

func replaceFlowSlotChild(arena *flowSlotMapArena, children []*flowSlotMap, index int, child *flowSlotMap) []*flowSlotMap {
	newChildren := allocFlowSlotChildren(arena, len(children))
	copy(newChildren, children)
	newChildren[index] = child
	return newChildren
}

func deleteFlowSlotChild(arena *flowSlotMapArena, children []*flowSlotMap, index int) []*flowSlotMap {
	if len(children) == 1 {
		return nil
	}
	newChildren := allocFlowSlotChildren(arena, len(children)-1)
	copy(newChildren[:index], children[:index])
	copy(newChildren[index:], children[index+1:])
	return newChildren
}
