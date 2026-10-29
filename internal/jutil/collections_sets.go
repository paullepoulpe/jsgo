package jutil

import (
	"iter"
	"slices"
	"strings"
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
		h += ObjectHashCode(e)
	}
	return h
}

// javaSet is implemented by set types so that Set.equals can compare across
// implementations.
type javaSet interface {
	Size() int
	containsAny(e any) bool
}

// ---------------------------------------------------------------------------
// LinkedHashSet

// LinkedHashSet is java.util.LinkedHashSet (insertion order); it also models
// Guava ImmutableSet.
type LinkedHashSet[T comparable] struct {
	// Go-only: performance — held by value: one allocation per set.
	m LinkedHashMap[T, struct{}]
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
				slots: slices.Grow([]lhmSlot[T, struct{}](nil), cap),
			},
		}
	}
	return &LinkedHashSet[T]{}
}

func (s *LinkedHashSet[T]) Size() int {
	if s == nil {
		return 0
	}
	return int(s.m.size)
}
func (s *LinkedHashSet[T]) IsEmpty() bool { return s.Size() == 0 }

// Add is Set.add; reports whether the set changed.
func (s *LinkedHashSet[T]) Add(e T) bool {
	if s.m.find(e) >= 0 {
		return false
	}
	s.m.add(e, struct{}{})
	return true
}

// AddNew inserts e when the caller has already verified that e is absent.
func (s *LinkedHashSet[T]) AddNew(e T) {
	s.m.add(e, struct{}{})
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
	s.m.removeSlot(i)
	return true
}

// Clear is Set.clear.
func (s *LinkedHashSet[T]) Clear() { s.m.Clear() }

// All iterates in insertion order (see LinkedHashMap.Keys).
func (s *LinkedHashSet[T]) All() iter.Seq[T] {
	var h *LinkedHashMap[T, struct{}]
	if s != nil {
		h = &s.m
	}
	return func(yield func(T) bool) {
		if h == nil {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if m := h.loadIndex(); i < int(h.first) || m != nil && m.isRemoved(i) {
				continue
			}
			if !yield(h.slots[i].key) {
				return
			}
		}
	}
}

// ToSlice returns the elements in insertion order.
func (s *LinkedHashSet[T]) ToSlice() []T {
	if s == nil {
		return nil
	}
	return s.m.KeySlice()
}

// First returns the first element in insertion order, if any, without closures.
func (s *LinkedHashSet[T]) First() (T, bool) {
	if s != nil && s.m.size > 0 {
		h := &s.m
		m := h.loadIndex()
		if m == nil || len(m.removed) == 0 {
			return h.slots[h.first].key, true
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if i >= len(m.removed) || !m.removed[i] {
				return h.slots[i].key, true
			}
		}
	}
	var zero T
	return zero, false
}

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
	other, ok := o.(javaSet)
	if !ok || IsNil(o) || other.Size() != s.Size() {
		return false
	}
	for e := range s.All() {
		if !other.containsAny(e) {
			return false
		}
	}
	return true
}

// HashCode is Set.hashCode(): the sum of the element hash codes.
func (s *LinkedHashSet[T]) HashCode() int32 {
	if s == nil || s.m.size == 0 {
		return 0
	}
	var h int32
	var removed []bool
	if m := s.m.loadIndex(); m != nil {
		removed = m.removed
	}
	for i := int(s.m.first); i < len(s.m.slots); i++ {
		if i < len(removed) && removed[i] {
			continue
		}
		h += ObjectHashCodeOf(s.m.slots[i].key)
	}
	return h
}

// ---------------------------------------------------------------------------
// ObjectLinkedHashMap / ObjectLinkedHashSet

type olhmSlot[K any, V any] struct {
	key     K
	value   V
	hash    int32 // ObjectHashCode(key)
	removed bool
}

// ObjectLinkedHashMap is java.util.LinkedHashMap (and Guava ImmutableMap)
// for keys whose classes override equals()/hashCode(): keys are compared
// with ObjectsEquals and hashed with ObjectHashCode instead of Go ==.
// Iteration is in insertion order.
//
// Go-only: performance — as for LinkedHashMap, the index is only built once
// the slot slice holds more than lhmSmallSize entries (in the typed-app
// benchmark, 65% of these maps stay empty); until then lookups walk the
// contiguous slots slice and call ObjectsEquals on the entries whose hash
// matches. The index hashes the hash codes stored in the slots.
type ObjectLinkedHashMap[K any, V any] struct {
	m     *lhmIndex // nil until len(slots) exceeds lhmSmallSize; m.removed is unused
	slots []olhmSlot[K, V]
	size  int32
	first int32
}

// NewObjectLinkedHashMap returns an empty map.
func NewObjectLinkedHashMap[K any, V any]() *ObjectLinkedHashMap[K, V] {
	return &ObjectLinkedHashMap[K, V]{}
}

// find returns k's slot index (or -1) and its hash code.
func (h *ObjectLinkedHashMap[K, V]) find(k K) (int, int32) {
	// Go-only: performance — the ...Of forms do not box String keys.
	hash := ObjectHashCodeOf(k)
	if h == nil || h.size == 0 {
		return -1, hash
	}
	slots := h.slots
	if h.m != nil {
		t := h.m.table
		mask := uint(len(t) - 1)
		for b := lhmHash(hash) & mask; ; b = (b + 1) & mask {
			switch e := t[b]; e {
			case 0:
				return -1, hash
			case lhmTombstone:
			default:
				if slots[e-1].hash == hash && ObjectsEqualsOf(slots[e-1].key, k) {
					return int(e - 1), hash
				}
			}
		}
	}
	for i := int(h.first); i < len(slots); i++ {
		if slots[i].hash == hash && !slots[i].removed && ObjectsEqualsOf(slots[i].key, k) {
			return i, hash
		}
	}
	return -1, hash
}

func (h *ObjectLinkedHashMap[K, V]) Size() int {
	if h == nil {
		return 0
	}
	return int(h.size)
}
func (h *ObjectLinkedHashMap[K, V]) IsEmpty() bool { return h.Size() == 0 }

// Get is Map.get: the zero value (Java null) when absent.
func (h *ObjectLinkedHashMap[K, V]) Get(k K) V {
	v, _ := h.Lookup(k)
	return v
}

// Lookup returns the value and whether the key is present.
func (h *ObjectLinkedHashMap[K, V]) Lookup(k K) (V, bool) {
	if h == nil || h.size == 0 {
		var zero V
		return zero, false
	}
	if i, _ := h.find(k); i >= 0 {
		return h.slots[i].value, true
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
	h.addWithHash(k, v, hash)
	var zero V
	return zero
}

// addWithHash appends an entry for k, which must be absent and have hash
// code hash.
func (h *ObjectLinkedHashMap[K, V]) addWithHash(k K, v V, hash int32) {
	if cap(h.slots) == 0 {
		h.slots = make([]olhmSlot[K, V], 0, 2)
		h.first = 0
	} else if len(h.slots) >= 64 && int(h.size)*4 < len(h.slots) {
		h.compact()
	}
	idx := len(h.slots)
	h.slots = append(h.slots, olhmSlot[K, V]{key: k, value: v, hash: hash})
	h.size++
	switch {
	case h.m != nil && !h.m.needsGrow():
		h.m.insert(lhmHash(hash), idx)
	case h.m != nil || int(h.size) > lhmSmallSize || len(h.slots) > 2*lhmSmallSize:
		h.reindex()
	}
}

// reindex (re)builds the index from the live slots, sized for them.
func (h *ObjectLinkedHashMap[K, V]) reindex() {
	if h.m == nil {
		h.m = &lhmIndex{}
	}
	h.m.reset(lhmTableLen(int(h.size)))
	for i := int(h.first); i < len(h.slots); i++ {
		if !h.slots[i].removed {
			h.m.insert(lhmHash(h.slots[i].hash), i)
		}
	}
}

// compact drops the removed slots, keeping the order of the others.
//
//go:noinline
func (h *ObjectLinkedHashMap[K, V]) compact() {
	w := 0
	for i := int(h.first); i < len(h.slots); i++ {
		if !h.slots[i].removed {
			h.slots[w] = h.slots[i]
			w++
		}
	}
	clear(h.slots[w:])
	h.slots = h.slots[:w]
	h.first = 0
	if h.m != nil {
		h.reindex()
	}
}

// Remove is Map.remove; returns the removed value (zero if absent).
func (h *ObjectLinkedHashMap[K, V]) Remove(k K) V {
	if h == nil || h.size == 0 {
		var zero V
		return zero
	}
	i, hash := h.find(k)
	if i < 0 {
		var zero V
		return zero
	}
	old := h.slots[i].value
	h.removeSlot(i, hash)
	return old
}

// removeSlot removes the entry at slot index i, whose hash code is hash.
func (h *ObjectLinkedHashMap[K, V]) removeSlot(i int, hash int32) {
	if h.m != nil {
		h.m.erase(lhmHash(hash), i)
	}
	h.size--
	h.slots[i] = olhmSlot[K, V]{removed: true}
	if i == int(h.first) {
		first := i + 1
		for first < len(h.slots) && h.slots[first].removed {
			first++
		}
		h.first = int32(first)
	}
}

// All iterates entries in insertion order.
func (h *ObjectLinkedHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if h == nil || h.size == 0 {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if h.slots[i].removed {
				continue
			}
			if !yield(h.slots[i].key, h.slots[i].value) {
				return
			}
		}
	}
}

// Keys iterates keys in insertion order (see LinkedHashMap.Keys).
func (h *ObjectLinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if h == nil || h.size == 0 {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if h.slots[i].removed {
				continue
			}
			if !yield(h.slots[i].key) {
				return
			}
		}
	}
}

// Values iterates values in insertion order.
func (h *ObjectLinkedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if h == nil || h.size == 0 {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if h.slots[i].removed {
				continue
			}
			if !yield(h.slots[i].value) {
				return
			}
		}
	}
}

// KeySlice returns the keys in insertion order.
func (h *ObjectLinkedHashMap[K, V]) KeySlice() []K {
	n := h.Size()
	if n == 0 {
		return nil
	}
	out := make([]K, 0, n)
	slots := h.slots
	if n == len(slots) {
		for i := range slots {
			out = append(out, slots[i].key)
		}
		return out
	}
	for i := int(h.first); i < len(slots); i++ {
		if !slots[i].removed {
			out = append(out, slots[i].key)
		}
	}
	return out
}

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
	other, ok := o.(javaMapEntries)
	if !ok || IsNil(o) || other.Size() != h.Size() {
		return false
	}
	for k, v := range h.All() {
		ov, present := other.lookupAny(k)
		if !present || !ObjectsEquals(v, ov) {
			return false
		}
	}
	return true
}

// HashCode is Map.hashCode().
func (h *ObjectLinkedHashMap[K, V]) HashCode() int32 {
	if h == nil || h.size == 0 {
		return 0
	}
	var sum int32
	for i := int(h.first); i < len(h.slots); i++ {
		if h.slots[i].removed {
			continue
		}
		sum += ObjectHashCodeOf(h.slots[i].key) ^ ObjectHashCodeOf(h.slots[i].value)
	}
	return sum
}

// String renders like AbstractMap.toString().
func (h *ObjectLinkedHashMap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, v := range h.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(ToJavaString(k))
		sb.WriteByte('=')
		sb.WriteString(ToJavaString(v))
	}
	sb.WriteByte('}')
	return sb.String()
}

// Clear is Map.clear() for ObjectLinkedHashMap.
func (h *ObjectLinkedHashMap[K, V]) Clear() {
	clear(h.slots)
	h.slots = h.slots[:0]
	if h.m != nil {
		h.m.reset(len(h.m.table))
	}
	h.size = 0
	h.first = 0
}

// ObjectLinkedHashSet is java.util.LinkedHashSet (and Guava ImmutableSet)
// for elements whose classes override equals()/hashCode().
type ObjectLinkedHashSet[T any] struct {
	// Go-only: performance — held by value: one allocation per set.
	m ObjectLinkedHashMap[T, struct{}]
}

// NewObjectLinkedHashSet returns a set containing elems in order.
func NewObjectLinkedHashSet[T any](elems ...T) *ObjectLinkedHashSet[T] {
	s := &ObjectLinkedHashSet[T]{}
	for _, e := range elems {
		s.Add(e)
	}
	return s
}

func (s *ObjectLinkedHashSet[T]) Size() int {
	if s == nil {
		return 0
	}
	return s.m.Size()
}
func (s *ObjectLinkedHashSet[T]) IsEmpty() bool { return s.Size() == 0 }

// Add is Set.add; reports whether the set changed.
func (s *ObjectLinkedHashSet[T]) Add(e T) bool {
	i, hash := s.m.find(e)
	if i >= 0 {
		return false
	}
	s.m.addWithHash(e, struct{}{}, hash)
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
	i, hash := s.m.find(e)
	if i < 0 {
		return false
	}
	s.m.removeSlot(i, hash)
	return true
}

// All iterates in insertion order (see LinkedHashMap.Keys).
func (s *ObjectLinkedHashSet[T]) All() iter.Seq[T] {
	var h *ObjectLinkedHashMap[T, struct{}]
	if s != nil {
		h = &s.m
	}
	return func(yield func(T) bool) {
		if h == nil || h.size == 0 {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if h.slots[i].removed {
				continue
			}
			if !yield(h.slots[i].key) {
				return
			}
		}
	}
}

// ToSlice returns the elements in insertion order.
func (s *ObjectLinkedHashSet[T]) ToSlice() []T {
	if s == nil {
		return nil
	}
	return s.m.KeySlice()
}

func (s *ObjectLinkedHashSet[T]) containsAny(e any) bool {
	x, ok := e.(T)
	return ok && s.Contains(x)
}

// Equals is Set.equals(Object).
func (s *ObjectLinkedHashSet[T]) Equals(o any) bool {
	if any(s) == o {
		return true
	}
	other, ok := o.(javaSet)
	if !ok || IsNil(o) || other.Size() != s.Size() {
		return false
	}
	for e := range s.All() {
		if !other.containsAny(e) {
			return false
		}
	}
	return true
}

// HashCode is Set.hashCode().
func (s *ObjectLinkedHashSet[T]) HashCode() int32 {
	if s == nil || s.m.size == 0 {
		return 0
	}
	var h int32
	for i := int(s.m.first); i < len(s.m.slots); i++ {
		if s.m.slots[i].removed {
			continue
		}
		h += ObjectHashCodeOf(s.m.slots[i].key)
	}
	return h
}

// String renders like AbstractCollection.toString().
func (s *ObjectLinkedHashSet[T]) String() string { return CollectionToString(s.ToSlice()) }

// Clear is Set.clear() for ObjectLinkedHashSet.
func (s *ObjectLinkedHashSet[T]) Clear() {
	s.m.Clear()
}
