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

// PutIfAbsent is Map.putIfAbsent: if the key is absent or mapped to null
// (jutil.IsNil) it is (re)mapped to v and null is returned; otherwise the
// current value is returned.
func (h *LinkedHashMap[K, V]) PutIfAbsent(k K, v V) V {
	if e, ok := h.m[k]; ok && !IsNil(e.value) {
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
	hash       int32
	value      V
	prev, next *olhmEntry[K, V]
	removed    bool
}

// ObjectLinkedHashMap is java.util.LinkedHashMap (and Guava ImmutableMap)
// for keys whose classes override equals()/hashCode(): keys are compared
// with ObjectsEquals and hashed with ObjectHashCode instead of Go ==.
// Iteration is in insertion order.
type ObjectLinkedHashMap[K any, V any] struct {
	buckets    map[int32][]*olhmEntry[K, V]
	size       int
	head, tail *olhmEntry[K, V]
}

// NewObjectLinkedHashMap returns an empty map.
func NewObjectLinkedHashMap[K any, V any]() *ObjectLinkedHashMap[K, V] {
	return &ObjectLinkedHashMap[K, V]{buckets: make(map[int32][]*olhmEntry[K, V])}
}

func (h *ObjectLinkedHashMap[K, V]) find(k K) (*olhmEntry[K, V], int32) {
	hash := ObjectHashCode(k)
	if h == nil {
		return nil, hash
	}
	for _, e := range h.buckets[hash] {
		if ObjectsEquals(e.key, k) {
			return e, hash
		}
	}
	return nil, hash
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
	e = &olhmEntry[K, V]{key: k, hash: hash, value: v, prev: h.tail}
	if h.tail != nil {
		h.tail.next = e
	} else {
		h.head = e
	}
	h.tail = e
	h.buckets[hash] = append(h.buckets[hash], e)
	h.size++
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
	b := h.buckets[hash]
	for i, x := range b {
		if x == e {
			h.buckets[hash] = append(b[:i:i], b[i+1:]...)
			break
		}
	}
	if len(h.buckets[hash]) == 0 {
		delete(h.buckets, hash)
	}
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
	return e.value
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

// Keys iterates keys in insertion order.
func (h *ObjectLinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range h.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values iterates values in insertion order.
func (h *ObjectLinkedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range h.All() {
			if !yield(v) {
				return
			}
		}
	}
}

// KeySlice returns the keys in insertion order.
func (h *ObjectLinkedHashMap[K, V]) KeySlice() []K { return slices.Collect(h.Keys()) }

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
	m *ObjectLinkedHashMap[T, struct{}]
}

// NewObjectLinkedHashSet returns a set containing elems in order.
func NewObjectLinkedHashSet[T any](elems ...T) *ObjectLinkedHashSet[T] {
	s := &ObjectLinkedHashSet[T]{m: NewObjectLinkedHashMap[T, struct{}]()}
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

// All iterates in insertion order.
func (s *ObjectLinkedHashSet[T]) All() iter.Seq[T] {
	if s == nil {
		return func(func(T) bool) {}
	}
	return s.m.Keys()
}

// ToSlice returns the elements in insertion order.
func (s *ObjectLinkedHashSet[T]) ToSlice() []T { return slices.Collect(s.All()) }

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
