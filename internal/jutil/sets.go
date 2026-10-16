package jutil

import "iter"

// Go-only: language-adapter — the hash-set factories of Guava's
// com.google.common.collect.Sets, as the Guava shaded into
// closure-compiler.jar implements them (checked with javap).

// Sets_NewHashSet is Sets.newHashSet(): new HashSet<>().
func Sets_NewHashSet[T any]() *JavaHashSet[T] { return NewJavaHashSet[T]() }

// Sets_NewHashSetOf is Sets.newHashSet(E... elements):
// newHashSetWithExpectedSize(elements.length), then the elements in order.
func Sets_NewHashSetOf[T any](elements ...T) *JavaHashSet[T] {
	s := Sets_NewHashSetWithExpectedSize[T](len(elements))
	for _, e := range elements {
		s.Add(e)
	}
	return s
}

// Sets_NewHashSetFromCollection is Sets.newHashSet(Iterable) for an
// Iterable that is a Collection: new HashSet<>(collection).
func Sets_NewHashSetFromCollection[T any](elements Collection[T]) *JavaHashSet[T] {
	return NewJavaHashSetFrom[T](elements)
}

// Sets_NewHashSetFromIterator is Sets.newHashSet(Iterator), and
// Sets.newHashSet(Iterable) for an Iterable that is not a Collection: new
// HashSet<>() filled in iteration order.
func Sets_NewHashSetFromIterator[T any](elements iter.Seq[T]) *JavaHashSet[T] {
	s := NewJavaHashSet[T]()
	s.AddAll(elements)
	return s
}

// Sets_NewHashSetWithExpectedSize is Sets.newHashSetWithExpectedSize(int):
// new HashSet<>(Maps.capacity(expectedSize)).
func Sets_NewHashSetWithExpectedSize[T any](expectedSize int) *JavaHashSet[T] {
	return NewJavaHashSetWithCapacity[T](Maps_Capacity(expectedSize))
}
