package jutil

import (
	"iter"
	"slices"
	"strings"
)

// Go-only: language-adapter — java.util / Guava collections with Java
// iteration-order semantics. Ported code uses these (never a ranged Go map)
// wherever Java iteration order is observable.

// ---------------------------------------------------------------------------
// ArrayList

// ArrayList is java.util.ArrayList: a growable list with reference
// semantics (mutations through any alias are visible to all holders, unlike
// a bare Go slice). Use plain []T for Java ImmutableList / List.of values
// that are never mutated after construction.
type ArrayList[T any] struct {
	elems []T
}

// NewArrayList returns an empty list.
func NewArrayList[T any]() *ArrayList[T] { return &ArrayList[T]{} }

// NewArrayListWithCapacity is new ArrayList<>(capacity).
func NewArrayListWithCapacity[T any](capacity int) *ArrayList[T] {
	return &ArrayList[T]{elems: make([]T, 0, capacity)}
}

// NewArrayListFrom is new ArrayList<>(collection) / Lists.newArrayList(...).
func NewArrayListFrom[T any](elems []T) *ArrayList[T] {
	return &ArrayList[T]{elems: slices.Clone(elems)}
}

func (l *ArrayList[T]) Size() int     { return len(l.elems) }
func (l *ArrayList[T]) IsEmpty() bool { return len(l.elems) == 0 }

// Add appends e; always returns true like Collection.add.
func (l *ArrayList[T]) Add(e T) bool {
	l.elems = append(l.elems, e)
	return true
}

// AddAt is List.add(int, E).
func (l *ArrayList[T]) AddAt(index int, e T) {
	if index < 0 || index > len(l.elems) {
		panic(NewIndexOutOfBoundsException(LenientFormat("Index: %s, Size: %s", index, len(l.elems))))
	}
	l.elems = slices.Insert(l.elems, index, e)
}

// AddAll appends all elements.
func (l *ArrayList[T]) AddAll(es ...T) bool {
	l.elems = append(l.elems, es...)
	return len(es) > 0
}

func (l *ArrayList[T]) checkIndex(index int) {
	if uint(index) >= uint(len(l.elems)) {
		l.indexOutOfBounds(index)
	}
}

// indexOutOfBounds is kept out of line so that the checks above inline.
//
//go:noinline
func (l *ArrayList[T]) indexOutOfBounds(index int) {
	panic(NewIndexOutOfBoundsException(LenientFormat("Index %s out of bounds for length %s", index, len(l.elems))))
}

// Get is List.get(int).
func (l *ArrayList[T]) Get(index int) T {
	if uint(index) >= uint(len(l.elems)) {
		l.indexOutOfBounds(index)
	}
	return l.elems[index]
}

// Set is List.set(int, E); returns the previous element.
func (l *ArrayList[T]) Set(index int, e T) T {
	l.checkIndex(index)
	old := l.elems[index]
	l.elems[index] = e
	return old
}

// RemoveAt is List.remove(int); returns the removed element.
func (l *ArrayList[T]) RemoveAt(index int) T {
	l.checkIndex(index)
	old := l.elems[index]
	l.elems = slices.Delete(l.elems, index, index+1)
	return old
}

// Clear is List.clear().
func (l *ArrayList[T]) Clear() { clear(l.elems); l.elems = l.elems[:0] }

// IndexOf is List.indexOf(Object) using Go == on the dynamic values
// (identity for pointers, value equality for strings and numbers).
func (l *ArrayList[T]) IndexOf(e T) int {
	for i, x := range l.elems {
		if any(x) == any(e) {
			return i
		}
	}
	return -1
}

// Contains is List.contains(Object) with the same equality as IndexOf.
func (l *ArrayList[T]) Contains(e T) bool { return l.IndexOf(e) >= 0 }

// Remove is List.remove(Object): removes the first equal element.
func (l *ArrayList[T]) Remove(e T) bool {
	if i := l.IndexOf(e); i >= 0 {
		l.elems = slices.Delete(l.elems, i, i+1)
		return true
	}
	return false
}

// GetFirst / GetLast are List.getFirst()/getLast() (Java 21).
func (l *ArrayList[T]) GetFirst() T { l.checkIndex(0); return l.elems[0] }
func (l *ArrayList[T]) GetLast() T  { l.checkIndex(len(l.elems) - 1); return l.elems[len(l.elems)-1] }

// RemoveLast is List.removeLast() (Java 21).
func (l *ArrayList[T]) RemoveLast() T { return l.RemoveAt(len(l.elems) - 1) }

// Slice returns the backing slice. Callers must not retain it across
// mutations of the list.
func (l *ArrayList[T]) Slice() []T {
	if l == nil {
		return nil
	}
	return l.elems
}

// ToSlice returns a copy (Java: ImmutableList.copyOf(list)).
func (l *ArrayList[T]) ToSlice() []T { return slices.Clone(l.elems) }

// All iterates in index order (Java: for (T e : list)).
func (l *ArrayList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, e := range l.elems {
			if !yield(e) {
				return
			}
		}
	}
}

// String renders like AbstractCollection.toString(): "[a, b]".
func (l *ArrayList[T]) String() string { return CollectionToString(l.elems) }

// CollectionToString renders elements like Java's AbstractCollection.toString.
func CollectionToString[T any](elems []T) string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, e := range elems {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ToJavaString(e))
	}
	sb.WriteByte(']')
	return sb.String()
}

// ---------------------------------------------------------------------------
// LinkedHashMap

type lhmEntry[K comparable, V any] struct {
	key        K
	value      V
	prev, next *lhmEntry[K, V]
	removed    bool
}

// LinkedHashMap is java.util.LinkedHashMap (insertion order). It also models
// Guava ImmutableMap and any Java map whose iteration order is insertion
// order. Keys use Go equality, which matches Java equals() for strings,
// boxed numbers, enums and identity-compared objects (pointers).
type LinkedHashMap[K comparable, V any] struct {
	m          map[K]*lhmEntry[K, V]
	head, tail *lhmEntry[K, V]
}

// NewLinkedHashMap returns an empty map.
func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{m: make(map[K]*lhmEntry[K, V])}
}

func (h *LinkedHashMap[K, V]) Size() int {
	if h == nil {
		return 0
	}
	return len(h.m)
}
func (h *LinkedHashMap[K, V]) IsEmpty() bool { return h.Size() == 0 }

// Get is Map.get: the zero value (Java null) when absent.
func (h *LinkedHashMap[K, V]) Get(k K) V {
	if h != nil {
		if e, ok := h.m[k]; ok {
			return e.value
		}
	}
	var zero V
	return zero
}

// Lookup returns the value and whether the key is present.
func (h *LinkedHashMap[K, V]) Lookup(k K) (V, bool) {
	if h != nil {
		if e, ok := h.m[k]; ok {
			return e.value, true
		}
	}
	var zero V
	return zero, false
}

// GetOrDefault is Map.getOrDefault.
func (h *LinkedHashMap[K, V]) GetOrDefault(k K, def V) V {
	if v, ok := h.Lookup(k); ok {
		return v
	}
	return def
}

// ContainsKey is Map.containsKey.
func (h *LinkedHashMap[K, V]) ContainsKey(k K) bool {
	if h == nil {
		return false
	}
	_, ok := h.m[k]
	return ok
}

// Put is Map.put; re-putting an existing key keeps its position. Returns the
// previous value (zero if none).
func (h *LinkedHashMap[K, V]) Put(k K, v V) V {
	if e, ok := h.m[k]; ok {
		old := e.value
		e.value = v
		return old
	}
	e := &lhmEntry[K, V]{key: k, value: v, prev: h.tail}
	if h.tail != nil {
		h.tail.next = e
	} else {
		h.head = e
	}
	h.tail = e
	h.m[k] = e
	var zero V
	return zero
}

// PutIfAbsent is Map.putIfAbsent; returns the existing value if present.
func (h *LinkedHashMap[K, V]) PutIfAbsent(k K, v V) V {
	if e, ok := h.m[k]; ok {
		return e.value
	}
	h.Put(k, v)
	var zero V
	return zero
}

// ComputeIfAbsent is Map.computeIfAbsent (the mapping function must not
// return Java null; Go has no way to express that distinction generically).
func (h *LinkedHashMap[K, V]) ComputeIfAbsent(k K, f func(K) V) V {
	if e, ok := h.m[k]; ok {
		return e.value
	}
	v := f(k)
	h.Put(k, v)
	return v
}

// Remove is Map.remove; returns the removed value (zero if absent).
func (h *LinkedHashMap[K, V]) Remove(k K) V {
	e, ok := h.m[k]
	if !ok {
		var zero V
		return zero
	}
	delete(h.m, k)
	h.unlink(e)
	return e.value
}

func (h *LinkedHashMap[K, V]) unlink(e *lhmEntry[K, V]) {
	e.removed = true
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		h.head = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		h.tail = e.prev
	}
}

// Clear is Map.clear.
func (h *LinkedHashMap[K, V]) Clear() {
	for e := h.head; e != nil; e = e.next {
		e.removed = true
	}
	clear(h.m)
	h.head, h.tail = nil, nil
}

// All iterates entries in insertion order. Entries removed during iteration
// are skipped; entries added during iteration are visited.
func (h *LinkedHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if h == nil {
			return
		}
		for e := h.head; e != nil; e = e.next {
			if e.removed {
				continue
			}
			if !yield(e.key, e.value) {
				return
			}
		}
	}
}

// Keys iterates keys in insertion order (Java: map.keySet()).
func (h *LinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range h.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values iterates values in insertion order (Java: map.values()).
func (h *LinkedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range h.All() {
			if !yield(v) {
				return
			}
		}
	}
}

// KeySlice returns the keys in insertion order.
func (h *LinkedHashMap[K, V]) KeySlice() []K { return slices.Collect(h.Keys()) }

// String renders like AbstractMap.toString(): "{k1=v1, k2=v2}".
func (h *LinkedHashMap[K, V]) String() string {
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

// ---------------------------------------------------------------------------
// LinkedHashSet

// LinkedHashSet is java.util.LinkedHashSet (insertion order); it also models
// Guava ImmutableSet.
type LinkedHashSet[T comparable] struct {
	m *LinkedHashMap[T, struct{}]
}

// NewLinkedHashSet returns a set containing elems in order.
func NewLinkedHashSet[T comparable](elems ...T) *LinkedHashSet[T] {
	s := &LinkedHashSet[T]{m: NewLinkedHashMap[T, struct{}]()}
	for _, e := range elems {
		s.Add(e)
	}
	return s
}

func (s *LinkedHashSet[T]) Size() int {
	if s == nil {
		return 0
	}
	return s.m.Size()
}
func (s *LinkedHashSet[T]) IsEmpty() bool { return s.Size() == 0 }

// Add is Set.add; reports whether the set changed.
func (s *LinkedHashSet[T]) Add(e T) bool {
	if s.m.ContainsKey(e) {
		return false
	}
	s.m.Put(e, struct{}{})
	return true
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
func (s *LinkedHashSet[T]) Contains(e T) bool { return s != nil && s.m.ContainsKey(e) }

// Remove is Set.remove; reports whether the element was present.
func (s *LinkedHashSet[T]) Remove(e T) bool {
	if !s.m.ContainsKey(e) {
		return false
	}
	s.m.Remove(e)
	return true
}

// Clear is Set.clear.
func (s *LinkedHashSet[T]) Clear() { s.m.Clear() }

// All iterates in insertion order.
func (s *LinkedHashSet[T]) All() iter.Seq[T] {
	if s == nil {
		return func(func(T) bool) {}
	}
	return s.m.Keys()
}

// ToSlice returns the elements in insertion order.
func (s *LinkedHashSet[T]) ToSlice() []T { return slices.Collect(s.All()) }

// String renders like AbstractCollection.toString(): "[a, b]".
func (s *LinkedHashSet[T]) String() string { return CollectionToString(s.ToSlice()) }
