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

// InitWithBuffer makes l an empty list (new ArrayList<>(len(buf))) that
// stores its first len(buf) elements in buf; it grows into a new array, as
// any list does, past that.
//
// Go-only: performance — for a list held by value in another object with a
// small array next to it: the object, the list and its initial storage are
// then one allocation instead of three.
func (l *ArrayList[T]) InitWithBuffer(buf []T) {
	l.elems = buf[:0:len(buf)]
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

// lhmEntry is an entry of a LinkedHashMap. A removed entry keeps its next
// pointer, for the iterators that stand on it, and points prev at itself.
//
// Go-only: performance — prev == self marks removal instead of a bool
// field: with it, the entries of sets and of pointer-keyed maps took 48
// bytes instead of 32.
type lhmEntry[K comparable, V any] struct {
	key        K
	value      V
	prev, next *lhmEntry[K, V]
}

// removed reports whether e was removed from its map.
func (e *lhmEntry[K, V]) removed() bool { return e.prev == e }

// LinkedHashMap is java.util.LinkedHashMap (insertion order). It also models
// Guava ImmutableMap and any Java map whose iteration order is insertion
// order. Keys use Go equality, which matches Java equals() for strings,
// boxed numbers, enums and identity-compared objects (pointers).
//
// Go-only: performance — most maps of a compilation stay small (in the
// typed-app benchmark, half of them never get an entry and 98% never hold
// more than 8), and a Go map costs a header when it is made and a 200-byte
// group at its first insertion. So the Go map that indexes the entries is
// only built once the map holds more than lhmSmallSize entries; until then
// lookups walk the entry list, comparing keys with == as the Go map does.
type LinkedHashMap[K comparable, V any] struct {
	m          *lhmLarge[K, V] // nil until size exceeds lhmSmallSize
	head, tail *lhmEntry[K, V]
	size       int
}

// Go-only: performance — allocate entries of larger maps in batches.
// The many empty and tiny maps retain their original representation.
type lhmLarge[K comparable, V any] struct {
	index   map[K]*lhmEntry[K, V]
	entries Slab[lhmEntry[K, V]]
}

// lhmSmallSize is the most entries a LinkedHashMap or ObjectLinkedHashMap
// holds without a Go map.
const lhmSmallSize = 8

// NewLinkedHashMap returns an empty map.
func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{}
}

// find returns k's entry, or nil.
func (h *LinkedHashMap[K, V]) find(k K) *lhmEntry[K, V] {
	if h.m != nil {
		return h.m.index[k]
	}
	for e := h.head; e != nil; e = e.next {
		if e.key == k {
			return e
		}
	}
	return nil
}

func (h *LinkedHashMap[K, V]) Size() int {
	if h == nil {
		return 0
	}
	return h.size
}
func (h *LinkedHashMap[K, V]) IsEmpty() bool { return h.Size() == 0 }

// Get is Map.get: the zero value (Java null) when absent.
func (h *LinkedHashMap[K, V]) Get(k K) V {
	if h != nil {
		if e := h.find(k); e != nil {
			return e.value
		}
	}
	var zero V
	return zero
}

// Lookup returns the value and whether the key is present.
func (h *LinkedHashMap[K, V]) Lookup(k K) (V, bool) {
	if h != nil {
		if e := h.find(k); e != nil {
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
	return h.find(k) != nil
}

// Put is Map.put; re-putting an existing key keeps its position. Returns the
// previous value (zero if none).
func (h *LinkedHashMap[K, V]) Put(k K, v V) V {
	if e := h.find(k); e != nil {
		old := e.value
		e.value = v
		return old
	}
	h.add(k, v)
	var zero V
	return zero
}

// add appends an entry for k, which must be absent.
func (h *LinkedHashMap[K, V]) add(k K, v V) {
	var e *lhmEntry[K, V]
	if h.m == nil {
		e = &lhmEntry[K, V]{}
	} else {
		e = h.m.entries.New()
	}
	*e = lhmEntry[K, V]{key: k, value: v, prev: h.tail}
	if h.tail != nil {
		h.tail.next = e
	} else {
		h.head = e
	}
	h.tail = e
	h.size++
	if h.m != nil {
		h.m.index[k] = e
	} else if h.size > lhmSmallSize {
		h.m = &lhmLarge[K, V]{index: make(map[K]*lhmEntry[K, V], 2*h.size)}
		for x := h.head; x != nil; x = x.next {
			h.m.index[x.key] = x
		}
	}
}

// PutIfAbsent is Map.putIfAbsent: if the key is absent or mapped to null
// (jutil.IsNil) it is (re)mapped to v and null is returned; otherwise the
// current value is returned.
func (h *LinkedHashMap[K, V]) PutIfAbsent(k K, v V) V {
	if e := h.find(k); e != nil && !IsNil(e.value) {
		return e.value
	}
	h.Put(k, v)
	var zero V
	return zero
}

// LoadOrStore returns the existing value and true, or inserts v and returns
// it with false. Unlike Java putIfAbsent it treats a present null as present.
// Go-only: performance — one lookup for analysis bookkeeping, including misses.
func (h *LinkedHashMap[K, V]) LoadOrStore(k K, v V) (V, bool) {
	if e := h.find(k); e != nil {
		return e.value, true
	}
	h.add(k, v)
	return v, false
}

// ComputeIfAbsent is Map.computeIfAbsent (the mapping function must not
// return Java null; Go has no way to express that distinction generically).
func (h *LinkedHashMap[K, V]) ComputeIfAbsent(k K, f func(K) V) V {
	if e := h.find(k); e != nil {
		return e.value
	}
	v := f(k)
	h.Put(k, v)
	return v
}

// Remove is Map.remove; returns the removed value (zero if absent).
func (h *LinkedHashMap[K, V]) Remove(k K) V {
	e := h.find(k)
	if e == nil {
		var zero V
		return zero
	}
	old := e.value
	h.removeEntry(e)
	return old
}

// removeEntry removes e, an entry of the map.
func (h *LinkedHashMap[K, V]) removeEntry(e *lhmEntry[K, V]) {
	if h.m != nil {
		delete(h.m.index, e.key)
	}
	h.size--
	h.unlink(e)
	var key K
	var value V
	e.key, e.value = key, value
}

func (h *LinkedHashMap[K, V]) unlink(e *lhmEntry[K, V]) {
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
	e.prev = e // removed
}

// Clear is Map.clear.
func (h *LinkedHashMap[K, V]) Clear() {
	for e := h.head; e != nil; e = e.next {
		e.prev = e // removed
	}
	if h.m != nil {
		clear(h.m.index)
		h.m.entries = Slab[lhmEntry[K, V]]{}
	}
	h.head, h.tail = nil, nil
	h.size = 0
}

// All iterates entries in insertion order. Entries removed during iteration
// are skipped; entries added during iteration are visited.
func (h *LinkedHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if h == nil {
			return
		}
		for e := h.head; e != nil; e = e.next {
			if e.removed() {
				continue
			}
			if !yield(e.key, e.value) {
				return
			}
		}
	}
}

// Keys iterates keys in insertion order (Java: map.keySet()).
//
// Go-only: performance — Keys, Values and the sets' All walk the entry
// list themselves, like All, instead of ranging over All: nested iterator
// closures escape to the heap (seven allocations per loop over a
// LinkedHashSet), flat ones are inlined into the caller's loop.
func (h *LinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if h == nil {
			return
		}
		for e := h.head; e != nil; e = e.next {
			if e.removed() {
				continue
			}
			if !yield(e.key) {
				return
			}
		}
	}
}

// Values iterates values in insertion order (Java: map.values()).
func (h *LinkedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if h == nil {
			return
		}
		for e := h.head; e != nil; e = e.next {
			if e.removed() {
				continue
			}
			if !yield(e.value) {
				return
			}
		}
	}
}

// KeySlice returns the keys in insertion order.
//
// Go-only: performance — this and the other KeySlice and ToSlice methods
// were slices.Collect over the iterator; they now allocate the result once,
// at the collection's size, and range over the iterator directly (passing
// it to slices.Collect makes its closures escape). An empty collection
// still gives nil.
func (h *LinkedHashMap[K, V]) KeySlice() []K {
	n := h.Size()
	if n == 0 {
		return nil
	}
	out := make([]K, 0, n)
	for k := range h.Keys() {
		out = append(out, k)
	}
	return out
}

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

func (s *LinkedHashSet[T]) Size() int {
	if s == nil {
		return 0
	}
	return s.m.size
}
func (s *LinkedHashSet[T]) IsEmpty() bool { return s.Size() == 0 }

// Add is Set.add; reports whether the set changed.
func (s *LinkedHashSet[T]) Add(e T) bool {
	if s.m.find(e) != nil {
		return false
	}
	s.m.add(e, struct{}{})
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
func (s *LinkedHashSet[T]) Contains(e T) bool { return s != nil && s.m.find(e) != nil }

// Remove is Set.remove; reports whether the element was present.
func (s *LinkedHashSet[T]) Remove(e T) bool {
	x := s.m.find(e)
	if x == nil {
		return false
	}
	s.m.removeEntry(x)
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
		for e := h.head; e != nil; e = e.next {
			if e.removed() {
				continue
			}
			if !yield(e.key) {
				return
			}
		}
	}
}

// ToSlice returns the elements in insertion order.
func (s *LinkedHashSet[T]) ToSlice() []T {
	n := s.Size()
	if n == 0 {
		return nil
	}
	out := make([]T, 0, n)
	for e := range s.All() {
		out = append(out, e)
	}
	return out
}

// String renders like AbstractCollection.toString(): "[a, b]".
func (s *LinkedHashSet[T]) String() string { return CollectionToString(s.ToSlice()) }

// ---------------------------------------------------------------------------
// Collection / Map markers and value equality

// JavaCollection is implemented by the jutil types that model a
// java.util.Collection, so that `x instanceof Collection` can be expressed
// as a Go type assertion.
type JavaCollection interface {
	Size() int
	IsEmpty() bool
	javaCollection()
}

// JavaMap is implemented by the jutil types that model a java.util.Map.
type JavaMap interface {
	Size() int
	IsEmpty() bool
	javaMap()
}

func (*ArrayList[T]) javaCollection()           {}
func (*LinkedHashSet[T]) javaCollection()       {}
func (*LinkedHashMap[K, V]) javaMap()           {}
func (*ObjectLinkedHashSet[T]) javaCollection() {}
func (*ObjectLinkedHashMap[K, V]) javaMap()     {}
func (*TreeMap[K, V]) javaMap()                 {}

// javaList is implemented by list types so that List.equals can compare
// across implementations.
type javaList interface {
	javaListElems() []any
}

func (l *ArrayList[T]) javaListElems() []any {
	out := make([]any, len(l.elems))
	for i, e := range l.elems {
		out[i] = e
	}
	return out
}

// Equals is List.equals(Object): same size and pairwise Objects.equals.
func (l *ArrayList[T]) Equals(o any) bool {
	if any(l) == o {
		return true
	}
	other, ok := o.(javaList)
	if !ok || IsNil(o) {
		return false
	}
	oe := other.javaListElems()
	if len(oe) != len(l.elems) {
		return false
	}
	for i, e := range l.elems {
		if !ObjectsEquals(e, oe[i]) {
			return false
		}
	}
	return true
}

// HashCode is List.hashCode().
func (l *ArrayList[T]) HashCode() int32 {
	var h int32 = 1
	for _, e := range l.elems {
		h = 31*h + ObjectHashCode(e)
	}
	return h
}

// javaSet is implemented by set types so that Set.equals can compare across
// implementations.
type javaSet interface {
	Size() int
	containsAny(e any) bool
}

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
	var h int32
	for e := range s.All() {
		h += ObjectHashCode(e)
	}
	return h
}

// javaMapEntries is implemented by map types so that Map.equals can compare
// across implementations.
type javaMapEntries interface {
	Size() int
	lookupAny(k any) (any, bool)
}

func (h *LinkedHashMap[K, V]) lookupAny(k any) (any, bool) {
	key, ok := k.(K)
	if !ok {
		return nil, false
	}
	return h.Lookup(key)
}

// Equals is Map.equals(Object): same size and equal values for every key.
func (h *LinkedHashMap[K, V]) Equals(o any) bool {
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

// HashCode is Map.hashCode(): the sum of key^value hash codes.
func (h *LinkedHashMap[K, V]) HashCode() int32 {
	var sum int32
	for k, v := range h.All() {
		sum += ObjectHashCode(k) ^ ObjectHashCode(v)
	}
	return sum
}

// NewLinkedHashMapFrom is new LinkedHashMap<>(map) / ImmutableMap.copyOf.
func NewLinkedHashMapFrom[K comparable, V any](src *LinkedHashMap[K, V]) *LinkedHashMap[K, V] {
	out := NewLinkedHashMap[K, V]()
	for k, v := range src.All() {
		out.Put(k, v)
	}
	return out
}

// ---------------------------------------------------------------------------
// ObjectLinkedHashMap / ObjectLinkedHashSet

type olhmEntry[K any, V any] struct {
	key        K
	value      V
	prev, next *olhmEntry[K, V]
	// hnext is the next entry with the same hash, in insertion order, once
	// the map has its buckets.
	hnext   *olhmEntry[K, V]
	hash    int32
	removed bool
}

// ObjectLinkedHashMap is java.util.LinkedHashMap (and Guava ImmutableMap)
// for keys whose classes override equals()/hashCode(): keys are compared
// with ObjectsEquals and hashed with ObjectHashCode instead of Go ==.
// Iteration is in insertion order.
//
// Go-only: performance — as for LinkedHashMap, the buckets are only built
// once the map holds more than lhmSmallSize entries (in the typed-app
// benchmark, 65% of these maps stay empty); until then lookups walk the
// entry list and call ObjectsEquals on the entries whose hash matches, the
// same entries in the same order as the bucket. A bucket chains its
// entries through hnext instead of a slice per hash.
type ObjectLinkedHashMap[K any, V any] struct {
	buckets    *olhmLarge[K, V] // nil until size exceeds lhmSmallSize
	size       int
	head, tail *olhmEntry[K, V]
}

// Go-only: performance — the object-equality map uses the same batched
// allocation after its small-map threshold.
type olhmLarge[K any, V any] struct {
	index   map[int32]*olhmEntry[K, V]
	entries Slab[olhmEntry[K, V]]
}

// NewObjectLinkedHashMap returns an empty map.
func NewObjectLinkedHashMap[K any, V any]() *ObjectLinkedHashMap[K, V] {
	return &ObjectLinkedHashMap[K, V]{}
}

func (h *ObjectLinkedHashMap[K, V]) find(k K) (*olhmEntry[K, V], int32) {
	// Go-only: performance — the ...Of forms do not box String keys.
	hash := ObjectHashCodeOf(k)
	if h == nil {
		return nil, hash
	}
	if h.buckets != nil {
		for e := h.buckets.index[hash]; e != nil; e = e.hnext {
			if ObjectsEqualsOf(e.key, k) {
				return e, hash
			}
		}
		return nil, hash
	}
	for e := h.head; e != nil; e = e.next {
		if e.hash == hash && ObjectsEqualsOf(e.key, k) {
			return e, hash
		}
	}
	return nil, hash
}

// chain appends e to the bucket of its hash.
func (h *ObjectLinkedHashMap[K, V]) chain(e *olhmEntry[K, V]) {
	last := h.buckets.index[e.hash]
	if last == nil {
		h.buckets.index[e.hash] = e
		return
	}
	for last.hnext != nil {
		last = last.hnext
	}
	last.hnext = e
}

func (h *ObjectLinkedHashMap[K, V]) Size() int {
	if h == nil {
		return 0
	}
	return h.size
}
func (h *ObjectLinkedHashMap[K, V]) IsEmpty() bool { return h.Size() == 0 }

// Get is Map.get: the zero value (Java null) when absent.
func (h *ObjectLinkedHashMap[K, V]) Get(k K) V {
	v, _ := h.Lookup(k)
	return v
}

// Lookup returns the value and whether the key is present.
func (h *ObjectLinkedHashMap[K, V]) Lookup(k K) (V, bool) {
	if e, _ := h.find(k); e != nil {
		return e.value, true
	}
	var zero V
	return zero, false
}

// ContainsKey is Map.containsKey.
func (h *ObjectLinkedHashMap[K, V]) ContainsKey(k K) bool {
	e, _ := h.find(k)
	return e != nil
}

// Put is Map.put; re-putting an existing key keeps its position (and its
// original key object). Returns the previous value (zero if none).
func (h *ObjectLinkedHashMap[K, V]) Put(k K, v V) V {
	e, hash := h.find(k)
	if e != nil {
		old := e.value
		e.value = v
		return old
	}
	if h.buckets == nil {
		e = &olhmEntry[K, V]{}
	} else {
		e = h.buckets.entries.New()
	}
	*e = olhmEntry[K, V]{key: k, hash: hash, value: v, prev: h.tail}
	if h.tail != nil {
		h.tail.next = e
	} else {
		h.head = e
	}
	h.tail = e
	h.size++
	if h.buckets != nil {
		h.chain(e)
	} else if h.size > lhmSmallSize {
		h.buckets = &olhmLarge[K, V]{index: make(map[int32]*olhmEntry[K, V], 2*h.size)}
		for x := h.head; x != nil; x = x.next {
			h.chain(x)
		}
	}
	var zero V
	return zero
}

// Remove is Map.remove; returns the removed value (zero if absent).
func (h *ObjectLinkedHashMap[K, V]) Remove(k K) V {
	e, hash := h.find(k)
	if e == nil {
		var zero V
		return zero
	}
	if h.buckets != nil {
		if first := h.buckets.index[hash]; first == e {
			if e.hnext != nil {
				h.buckets.index[hash] = e.hnext
			} else {
				delete(h.buckets.index, hash)
			}
		} else {
			p := first
			for p.hnext != e {
				p = p.hnext
			}
			p.hnext = e.hnext
		}
		e.hnext = nil
	}
	old := e.value
	h.size--
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
	var key K
	var value V
	e.key, e.value = key, value
	return old
}

// All iterates entries in insertion order.
func (h *ObjectLinkedHashMap[K, V]) All() iter.Seq2[K, V] {
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

// Keys iterates keys in insertion order (see LinkedHashMap.Keys).
func (h *ObjectLinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if h == nil {
			return
		}
		for e := h.head; e != nil; e = e.next {
			if e.removed {
				continue
			}
			if !yield(e.key) {
				return
			}
		}
	}
}

// Values iterates values in insertion order.
func (h *ObjectLinkedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if h == nil {
			return
		}
		for e := h.head; e != nil; e = e.next {
			if e.removed {
				continue
			}
			if !yield(e.value) {
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
	for k := range h.Keys() {
		out = append(out, k)
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
	var sum int32
	for k, v := range h.All() {
		sum += ObjectHashCode(k) ^ ObjectHashCode(v)
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
	if s.m.ContainsKey(e) {
		return false
	}
	s.m.Put(e, struct{}{})
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
	if !s.m.ContainsKey(e) {
		return false
	}
	s.m.Remove(e)
	return true
}

// All iterates in insertion order (see LinkedHashMap.Keys).
func (s *ObjectLinkedHashSet[T]) All() iter.Seq[T] {
	var h *ObjectLinkedHashMap[T, struct{}]
	if s != nil {
		h = &s.m
	}
	return func(yield func(T) bool) {
		if h == nil {
			return
		}
		for e := h.head; e != nil; e = e.next {
			if e.removed {
				continue
			}
			if !yield(e.key) {
				return
			}
		}
	}
}

// ToSlice returns the elements in insertion order.
func (s *ObjectLinkedHashSet[T]) ToSlice() []T {
	n := s.Size()
	if n == 0 {
		return nil
	}
	out := make([]T, 0, n)
	for e := range s.All() {
		out = append(out, e)
	}
	return out
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
	var h int32
	for e := range s.All() {
		h += ObjectHashCode(e)
	}
	return h
}

// String renders like AbstractCollection.toString().
func (s *ObjectLinkedHashSet[T]) String() string { return CollectionToString(s.ToSlice()) }
