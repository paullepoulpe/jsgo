package jutil

import (
	"iter"
	"slices"
)

// Go-only: language-adapter — the read side of java.util.Set, and
// insertion-ordered Set / Object-keyed Map and Set implementations.

// Set is the read-only part of java.util.Set.
type Set[T any] interface {
	Size() int
	IsEmpty() bool
	Contains(e T) bool
	// All iterates in the set's iteration order.
	All() iter.Seq[T]
}

var (
	_ Set[string] = (*LinkedHashSet[string])(nil)
	_ Set[string] = (*ObjectLinkedHashSet[string])(nil)
)

// SetHashCode is Set.hashCode() for a slice standing in for a Java set: the
// sum of the element hash codes.
func SetHashCode[T any](elems []T) int32 {
	var h int32
	for _, e := range elems {
		h += objectHashCodeOf(e)
	}
	return h
}

// javaSet is implemented by set types so that Set.equals can compare across
// implementations.
type javaSet interface {
	Size() int
	containsAny(e any) bool
}

func setEntriesEquals[T any](size int, seq iter.Seq[T], o any) bool {
	other, ok := o.(javaSet)
	if !ok || IsNil(o) || other.Size() != size {
		return false
	}
	for e := range seq {
		if !other.containsAny(e) {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// LinkedHashSet

// LinkedHashSet is java.util.LinkedHashSet (insertion order); it also models
// Guava ImmutableSet.
type LinkedHashSet[T comparable] struct {
	// Go-only: performance — held by value: one allocation per set.
	m LinkedHashMap[T, struct{}]
}

func (s *LinkedHashSet[T]) table() *orderedTable[T, struct{}, struct{}] {
	if s == nil {
		return nil
	}
	return &s.m.orderedTable
}

// NewLinkedHashSet returns a set containing elems in order.
func NewLinkedHashSet[T comparable](elems ...T) *LinkedHashSet[T] {
	s := &LinkedHashSet[T]{}
	for _, e := range elems {
		s.Add(e)
	}
	return s
}

// NewLinkedHashSetWithCapacity returns an empty set pre-sized for cap elements.
func NewLinkedHashSetWithCapacity[T comparable](cap int) *LinkedHashSet[T] {
	if cap > 0 {
		return &LinkedHashSet[T]{
			m: LinkedHashMap[T, struct{}]{
				orderedTable: orderedTable[T, struct{}, struct{}]{
					slots: slices.Grow([]lhmSlot[T, struct{}](nil), cap),
				},
			},
		}
	}
	return &LinkedHashSet[T]{}
}

func (s *LinkedHashSet[T]) Size() int     { return s.table().sizeOrZero() }
func (s *LinkedHashSet[T]) IsEmpty() bool { return s.Size() == 0 }

// Add is Set.add; reports whether the set changed.
func (s *LinkedHashSet[T]) Add(e T) bool {
	if s.m.find(e) >= 0 {
		return false
	}
	s.m.add(e, struct{}{}, struct{}{}, lhmSlotHash[T, struct{}])
	return true
}

// AddNew inserts e when the caller has already verified that e is absent.
func (s *LinkedHashSet[T]) AddNew(e T) {
	s.m.add(e, struct{}{}, struct{}{}, lhmSlotHash[T, struct{}])
}

// AddAll adds every element; reports whether the set changed.
func (s *LinkedHashSet[T]) AddAll(es ...T) bool {
	changed := false
	for _, e := range es {
		if s.Add(e) {
			changed = true
		}
	}
	return changed
}

// Contains is Set.contains.
func (s *LinkedHashSet[T]) Contains(e T) bool { return s != nil && s.m.find(e) >= 0 }

// Remove is Set.remove; reports whether the element was present.
func (s *LinkedHashSet[T]) Remove(e T) bool {
	i := s.m.find(e)
	if i < 0 {
		return false
	}
	s.m.removeSlot(i, lhmSlotHash[T, struct{}])
	return true
}

// Clear is Set.clear.
func (s *LinkedHashSet[T]) Clear() { s.m.clear() }

// All iterates in insertion order.
func (s *LinkedHashSet[T]) All() iter.Seq[T] { return s.table().keys() }

// ToSlice returns the elements in insertion order.
func (s *LinkedHashSet[T]) ToSlice() []T { return s.table().keySlice() }

// First returns the first element in insertion order, if any, without closures.
func (s *LinkedHashSet[T]) First() (T, bool) { return s.table().firstKey() }

// String renders like AbstractCollection.toString(): "[a, b]".
func (s *LinkedHashSet[T]) String() string { return CollectionToString(s.ToSlice()) }

func (s *LinkedHashSet[T]) containsAny(e any) bool {
	x, ok := e.(T)
	return ok && s.Contains(x)
}

// Equals is Set.equals(Object): same size and every element contained.
func (s *LinkedHashSet[T]) Equals(o any) bool {
	if any(s) == o {
		return true
	}
	if other, ok := o.(*LinkedHashSet[T]); ok {
		if s == nil || other == nil || other.m.size != s.m.size {
			return false
		}
		if s.m.size == 0 {
			return true
		}
		removed := s.m.removedSlice()
		for i := int(s.m.first); i < len(s.m.slots); i++ {
			if i < len(removed) && removed[i] {
				continue
			}
			if other.m.find(s.m.slots[i].key) < 0 {
				return false
			}
		}
		return true
	}
	return setEntriesEquals(s.Size(), s.All(), o)
}

// HashCode is Set.hashCode(): the sum of the element hash codes.
func (s *LinkedHashSet[T]) HashCode() int32 {
	if s == nil || s.m.size == 0 {
		return 0
	}
	var h int32
	removed := s.m.removedSlice()
	for i := int(s.m.first); i < len(s.m.slots); i++ {
		if i < len(removed) && removed[i] {
			continue
		}
		h += objectHashCodeOf(s.m.slots[i].key)
	}
	return h
}

// ---------------------------------------------------------------------------
// ObjectLinkedHashMap / ObjectLinkedHashSet

// ObjectLinkedHashMap is java.util.LinkedHashMap (and Guava ImmutableMap)
// for keys whose classes override equals()/hashCode(): keys are compared
// with ObjectsEquals and hashed with ObjectHashCode instead of Go ==.
// Iteration is in insertion order.
//
// Go-only: performance — as for LinkedHashMap, the index is only built once
// the slot slice holds more than lhmSmallSize entries (in the typed-app
// benchmark, 65% of these maps stay empty); until then lookups walk the
// contiguous slots slice and call ObjectsEquals on the entries whose cached
// hash code matches.
type ObjectLinkedHashMap[K any, V any] struct {
	orderedTable[K, V, int32]
}

func olhmSlotHash[K any, V any](s *orderedSlot[K, V, int32]) uint {
	return lhmHash(s.meta)
}

func (h *ObjectLinkedHashMap[K, V]) table() *orderedTable[K, V, int32] {
	if h == nil {
		return nil
	}
	return &h.orderedTable
}

// NewObjectLinkedHashMap returns an empty map.
func NewObjectLinkedHashMap[K any, V any]() *ObjectLinkedHashMap[K, V] {
	return &ObjectLinkedHashMap[K, V]{}
}

// find returns k's slot index (or -1) and its hash code.
func (h *ObjectLinkedHashMap[K, V]) find(k K) (int, int32) {
	// Go-only: performance — the ...Of forms do not box String keys.
	hash := objectHashCodeOf(k)
	return h.findWithHash(k, hash), hash
}

func (h *ObjectLinkedHashMap[K, V]) findWithHash(k K, hash int32) int {
	if h == nil || h.size == 0 {
		return -1
	}
	if m := h.loadIndex(); m != nil {
		return h.lookupIn(m, k, hash)
	}
	if int(h.size) > lhmSmallSize || len(h.slots) > 2*lhmSmallSize {
		return h.lookupIn(h.buildIndexShared(olhmSlotHash[K, V]), k, hash)
	}
	for i := int(h.first); i < len(h.slots); i++ {
		if h.slots[i].meta == hash && ObjectsEqualsOf(h.slots[i].key, k) {
			return i
		}
	}
	return -1
}

func (h *ObjectLinkedHashMap[K, V]) lookupIn(m *lhmIndex, k K, hash int32) int {
	slots := h.slots
	return m.find(lhmHash(hash), func(slot int) bool {
		return slots[slot].meta == hash && ObjectsEqualsOf(slots[slot].key, k)
	})
}

func (h *ObjectLinkedHashMap[K, V]) Size() int     { return h.table().sizeOrZero() }
func (h *ObjectLinkedHashMap[K, V]) IsEmpty() bool { return h.Size() == 0 }

// Get is Map.get: the zero value (Java null) when absent.
func (h *ObjectLinkedHashMap[K, V]) Get(k K) V {
	v, _ := h.Lookup(k)
	return v
}

// Lookup returns the value and whether the key is present.
func (h *ObjectLinkedHashMap[K, V]) Lookup(k K) (V, bool) {
	if h != nil && h.size > 0 {
		if i, _ := h.find(k); i >= 0 {
			return h.slots[i].value, true
		}
	}
	var zero V
	return zero, false
}

// ContainsKey is Map.containsKey.
func (h *ObjectLinkedHashMap[K, V]) ContainsKey(k K) bool {
	if h == nil || h.size == 0 {
		return false
	}
	i, _ := h.find(k)
	return i >= 0
}

// Put is Map.put; re-putting an existing key keeps its position (and its
// original key object). Returns the previous value (zero if none).
func (h *ObjectLinkedHashMap[K, V]) Put(k K, v V) V {
	i, hash := h.find(k)
	if i >= 0 {
		old := h.slots[i].value
		h.slots[i].value = v
		return old
	}
	h.add(k, v, hash, olhmSlotHash[K, V])
	var zero V
	return zero
}

// Remove is Map.remove; returns the removed value (zero if absent).
func (h *ObjectLinkedHashMap[K, V]) Remove(k K) V {
	if h == nil || h.size == 0 {
		var zero V
		return zero
	}
	i, _ := h.find(k)
	if i < 0 {
		var zero V
		return zero
	}
	old := h.slots[i].value
	h.removeSlot(i, olhmSlotHash[K, V])
	return old
}

// All iterates entries in insertion order.
func (h *ObjectLinkedHashMap[K, V]) All() iter.Seq2[K, V] { return h.table().all() }

// Keys iterates keys in insertion order.
func (h *ObjectLinkedHashMap[K, V]) Keys() iter.Seq[K] { return h.table().keys() }

// Values iterates values in insertion order.
func (h *ObjectLinkedHashMap[K, V]) Values() iter.Seq[V] { return h.table().values() }

// KeySlice returns the keys in insertion order.
func (h *ObjectLinkedHashMap[K, V]) KeySlice() []K { return h.table().keySlice() }

func (h *ObjectLinkedHashMap[K, V]) lookupAny(k any) (any, bool) {
	key, ok := k.(K)
	if !ok {
		return nil, false
	}
	return h.Lookup(key)
}

// Equals is Map.equals(Object).
func (h *ObjectLinkedHashMap[K, V]) Equals(o any) bool {
	if any(h) == o {
		return true
	}
	if other, ok := o.(*ObjectLinkedHashMap[K, V]); ok {
		if h == nil || other == nil || other.size != h.size {
			return false
		}
		if h.size == 0 {
			return true
		}
		removed := h.removedSlice()
		for i := int(h.first); i < len(h.slots); i++ {
			if i < len(removed) && removed[i] {
				continue
			}
			oi := other.findWithHash(h.slots[i].key, h.slots[i].meta)
			if oi < 0 || !ObjectsEqualsOf(h.slots[i].value, other.slots[oi].value) {
				return false
			}
		}
		return true
	}
	return mapEntriesEquals(h.Size(), h.All(), o)
}

// HashCode is Map.hashCode().
func (h *ObjectLinkedHashMap[K, V]) HashCode() int32 {
	if h == nil || h.size == 0 {
		return 0
	}
	var sum int32
	removed := h.removedSlice()
	for i := int(h.first); i < len(h.slots); i++ {
		if i < len(removed) && removed[i] {
			continue
		}
		sum += h.slots[i].meta ^ objectHashCodeOf(h.slots[i].value)
	}
	return sum
}

// String renders like AbstractMap.toString().
func (h *ObjectLinkedHashMap[K, V]) String() string { return formatMapString(h.All()) }

// Clear is Map.clear() for ObjectLinkedHashMap.
func (h *ObjectLinkedHashMap[K, V]) Clear() { h.clear() }

// ObjectLinkedHashSet is java.util.LinkedHashSet (and Guava ImmutableSet)
// for elements whose classes override equals()/hashCode().
type ObjectLinkedHashSet[T any] struct {
	// Go-only: performance — held by value: one allocation per set.
	m ObjectLinkedHashMap[T, struct{}]
}

func (s *ObjectLinkedHashSet[T]) table() *orderedTable[T, struct{}, int32] {
	if s == nil {
		return nil
	}
	return &s.m.orderedTable
}

// NewObjectLinkedHashSet returns a set containing elems in order.
func NewObjectLinkedHashSet[T any](elems ...T) *ObjectLinkedHashSet[T] {
	s := &ObjectLinkedHashSet[T]{}
	for _, e := range elems {
		s.Add(e)
	}
	return s
}

func (s *ObjectLinkedHashSet[T]) Size() int     { return s.table().sizeOrZero() }
func (s *ObjectLinkedHashSet[T]) IsEmpty() bool { return s.Size() == 0 }

// Add is Set.add; reports whether the set changed.
func (s *ObjectLinkedHashSet[T]) Add(e T) bool {
	i, hash := s.m.find(e)
	if i >= 0 {
		return false
	}
	s.m.add(e, struct{}{}, hash, olhmSlotHash[T, struct{}])
	return true
}

// AddAll adds every element of the sequence; reports whether the set changed.
func (s *ObjectLinkedHashSet[T]) AddAll(es iter.Seq[T]) bool {
	changed := false
	for e := range es {
		if s.Add(e) {
			changed = true
		}
	}
	return changed
}

// Contains is Set.contains.
func (s *ObjectLinkedHashSet[T]) Contains(e T) bool { return s != nil && s.m.ContainsKey(e) }

// Remove is Set.remove.
func (s *ObjectLinkedHashSet[T]) Remove(e T) bool {
	if s == nil || s.m.size == 0 {
		return false
	}
	i, _ := s.m.find(e)
	if i < 0 {
		return false
	}
	s.m.removeSlot(i, olhmSlotHash[T, struct{}])
	return true
}

// All iterates in insertion order.
func (s *ObjectLinkedHashSet[T]) All() iter.Seq[T] { return s.table().keys() }

// ToSlice returns the elements in insertion order.
func (s *ObjectLinkedHashSet[T]) ToSlice() []T { return s.table().keySlice() }

func (s *ObjectLinkedHashSet[T]) containsAny(e any) bool {
	x, ok := e.(T)
	return ok && s.Contains(x)
}

// Equals is Set.equals(Object).
func (s *ObjectLinkedHashSet[T]) Equals(o any) bool {
	if any(s) == o {
		return true
	}
	if other, ok := o.(*ObjectLinkedHashSet[T]); ok {
		if s == nil || other == nil || other.m.size != s.m.size {
			return false
		}
		if s.m.size == 0 {
			return true
		}
		removed := s.m.removedSlice()
		for i := int(s.m.first); i < len(s.m.slots); i++ {
			if i < len(removed) && removed[i] {
				continue
			}
			if other.m.findWithHash(s.m.slots[i].key, s.m.slots[i].meta) < 0 {
				return false
			}
		}
		return true
	}
	return setEntriesEquals(s.Size(), s.All(), o)
}

// HashCode is Set.hashCode().
func (s *ObjectLinkedHashSet[T]) HashCode() int32 {
	if s == nil || s.m.size == 0 {
		return 0
	}
	var h int32
	removed := s.m.removedSlice()
	for i := int(s.m.first); i < len(s.m.slots); i++ {
		if i < len(removed) && removed[i] {
			continue
		}
		h += s.m.slots[i].meta
	}
	return h
}

// String renders like AbstractCollection.toString().
func (s *ObjectLinkedHashSet[T]) String() string { return CollectionToString(s.ToSlice()) }

// Clear is Set.clear() for ObjectLinkedHashSet.
func (s *ObjectLinkedHashSet[T]) Clear() { s.m.clear() }
