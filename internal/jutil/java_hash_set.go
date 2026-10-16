package jutil

import (
	"iter"
	"slices"
	"strconv"
)

// Go-only: language-adapter — java.util.HashSet as JDK 21 behaves: a
// JavaHashMap from the elements to a placeholder value (Javadoc of
// java.util.HashSet: "backed by a hash table (actually a HashMap
// instance)"), so its iteration order is the map's key order. The bulk
// operations are those of AbstractCollection/AbstractSet, as their Javadoc
// describes them.

// JavaHashSet is java.util.HashSet.
type JavaHashSet[T any] struct{ m *JavaHashMap[T, struct{}] }

// NewJavaHashSet is new HashSet<>().
func NewJavaHashSet[T any]() *JavaHashSet[T] {
	return &JavaHashSet[T]{NewJavaHashMap[T, struct{}]()}
}

// NewJavaHashSetWithCapacity is new HashSet<>(initialCapacity).
func NewJavaHashSetWithCapacity[T any](initialCapacity int) *JavaHashSet[T] {
	return &JavaHashSet[T]{NewJavaHashMapWithCapacity[T, struct{}](initialCapacity)}
}

// NewJavaHashSetWithCapacityAndLoadFactor is new HashSet<>(initialCapacity,
// loadFactor).
func NewJavaHashSetWithCapacityAndLoadFactor[T any](initialCapacity int, loadFactor float32) *JavaHashSet[T] {
	return &JavaHashSet[T]{NewJavaHashMapWithCapacityAndLoadFactor[T, struct{}](initialCapacity, loadFactor)}
}

// NewJavaHashSetFrom is new HashSet<>(c): an initial capacity of
// ceil(c.size() / 0.75), but at least 16 (observed), then c's elements
// added in its iteration order.
func NewJavaHashSetFrom[T any](c Collection[T]) *JavaHashSet[T] {
	s := NewJavaHashSetWithCapacity[T](max(capacityForMappings(c.Size()), javaHashMapDefaultCapacity))
	s.AddAll(c.All())
	return s
}

// NewJavaHashSetFromSlice is new HashSet<>(Arrays.asList(elems)).
func NewJavaHashSetFromSlice[T any](elems []T) *JavaHashSet[T] {
	s := NewJavaHashSetWithCapacity[T](max(capacityForMappings(len(elems)), javaHashMapDefaultCapacity))
	s.AddAll(slices.Values(elems))
	return s
}

// JavaHashSet_NewHashSet is HashSet.newHashSet(numElements).
func JavaHashSet_NewHashSet[T any](numElements int) *JavaHashSet[T] {
	if numElements < 0 {
		panic(NewIllegalArgumentExceptionMsg("Negative number of elements: " + strconv.Itoa(numElements)))
	}
	return NewJavaHashSetWithCapacity[T](capacityForMappings(numElements))
}

// Size is Set.size().
func (s *JavaHashSet[T]) Size() int {
	if s == nil {
		return 0
	}
	return s.m.Size()
}

// IsEmpty is Set.isEmpty().
func (s *JavaHashSet[T]) IsEmpty() bool { return s.Size() == 0 }

// Contains is Set.contains(Object).
func (s *JavaHashSet[T]) Contains(e T) bool {
	return s != nil && s.m.ContainsKey(e)
}

// Add is Set.add: reports whether e was absent (a present element keeps its
// place and its original object).
func (s *JavaHashSet[T]) Add(e T) bool {
	_, present := s.m.insertKey(hashMapHashOf(e), e, struct{}{}, false)
	return !present
}

// Remove is Set.remove(Object): reports whether e was present.
func (s *JavaHashSet[T]) Remove(e T) bool {
	return s.m.deleteKey(hashMapHashOf(e), e, struct{}{}, false, false) != nil
}

// Clear is Set.clear(); the table keeps its size.
func (s *JavaHashSet[T]) Clear() { s.m.Clear() }

// AddAll is Collection.addAll: adds each element in order; reports whether
// the set changed.
func (s *JavaHashSet[T]) AddAll(es iter.Seq[T]) bool {
	changed := false
	for e := range es {
		if s.Add(e) {
			changed = true
		}
	}
	return changed
}

// ContainsAll is Collection.containsAll.
func (s *JavaHashSet[T]) ContainsAll(c Collection[T]) bool {
	for e := range c.All() {
		if !s.Contains(e) {
			return false
		}
	}
	return true
}

// RemoveAll is AbstractSet.removeAll: when this set is larger than c, each
// element of c is removed with remove(); otherwise this set is iterated and
// the elements c contains are removed through the iterator (Javadoc of
// AbstractSet.removeAll; equal sizes iterate this set, observed).
func (s *JavaHashSet[T]) RemoveAll(c Collection[T]) bool {
	changed := false
	if s.Size() > c.Size() {
		for e := range c.All() {
			if s.Remove(e) {
				changed = true
			}
		}
		return changed
	}
	return s.RemoveIf(c.Contains)
}

// RetainAll is AbstractCollection.retainAll: removes, through the iterator,
// the elements c does not contain.
func (s *JavaHashSet[T]) RetainAll(c Collection[T]) bool {
	return s.RemoveIf(func(e T) bool { return !c.Contains(e) })
}

// RemoveIf is Collection.removeIf: removes the matching elements through the
// iterator; reports whether any was removed.
func (s *JavaHashSet[T]) RemoveIf(pred func(T) bool) bool {
	if pred == nil {
		panic(NewNullPointerException())
	}
	changed := false
	it := s.Iterator()
	for it.HasNext() {
		if pred(it.Next()) {
			it.Remove()
			changed = true
		}
	}
	return changed
}

// JavaHashSetIterator is the set's iterator.
type JavaHashSetIterator[T any] struct {
	it *JavaHashMapIterator[T, struct{}]
}

// Iterator is Set.iterator().
func (s *JavaHashSet[T]) Iterator() JavaHashSetIterator[T] {
	return JavaHashSetIterator[T]{s.m.Iterator()}
}

// HasNext is Iterator.hasNext().
func (it JavaHashSetIterator[T]) HasNext() bool { return it.it.HasNext() }

// Next is Iterator.next().
func (it JavaHashSetIterator[T]) Next() T {
	e, _ := it.it.Next()
	return e
}

// Remove is Iterator.remove().
func (it JavaHashSetIterator[T]) Remove() { it.it.Remove() }

// All iterates the elements in Java's order.
func (s *JavaHashSet[T]) All() iter.Seq[T] {
	if s == nil {
		return func(func(T) bool) {}
	}
	return s.m.Keys()
}

// ToSlice returns the elements in iteration order.
func (s *JavaHashSet[T]) ToSlice() []T { return slices.Collect(s.All()) }

// Clone is HashSet.clone(): a shallow copy backed by the map's clone.
func (s *JavaHashSet[T]) Clone() *JavaHashSet[T] { return &JavaHashSet[T]{s.m.Clone()} }

// String is AbstractCollection.toString(): "[a, b]", with "(this
// Collection)" for the set itself.
func (s *JavaHashSet[T]) String() string { return javaCollectionString(any(s), s.All()) }

func (*JavaHashSet[T]) javaCollection() {}

func (s *JavaHashSet[T]) containsAny(e any) bool {
	x, ok := e.(T)
	if !ok {
		// Only a null e can match an element that is not a T: the zero T,
		// when T is nullable.
		var zero T
		if !IsNil(e) || !IsNil(any(zero)) {
			return false
		}
		x = zero
	}
	return s.Contains(x)
}

// Equals is Set.equals(Object): same size and every element contained.
func (s *JavaHashSet[T]) Equals(o any) bool {
	if any(s) == o {
		return true
	}
	return javaSetEquals(s.Size(), s.All(), o)
}

// HashCode is Set.hashCode(): the sum of the elements' hash codes.
func (s *JavaHashSet[T]) HashCode() int32 {
	var sum int32
	for e := range s.All() {
		sum += ObjectHashCode(any(e))
	}
	return sum
}

var (
	_ Set[string]        = (*JavaHashSet[string])(nil)
	_ Collection[string] = (*JavaHashSet[string])(nil)
	_ JavaCollection     = (*JavaHashSet[string])(nil)
	_ Hashable           = (*JavaHashSet[string])(nil)
	_ javaSet            = (*JavaHashSet[string])(nil)
)
