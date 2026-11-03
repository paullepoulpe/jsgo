package jutil

import (
	"iter"
	"slices"
)

// Go-only: language-adapter — the read side of java.util.Set, and
// insertion-ordered LinkedHashSet.

// Set is the read-only part of java.util.Set.
type Set[T any] interface {
	Size() int
	IsEmpty() bool
	Contains(e T) bool
	// All iterates in the set's iteration order.
	All() iter.Seq[T]
}

var _ Set[string] = (*LinkedHashSet[string])(nil)

// SetHashCode is Set.hashCode() for a slice standing in for a Java set: the
// sum of the element hash codes.
func SetHashCode[T any](elems []T) int32 {
	var h int32
	for _, e := range elems {
		h += ObjectHashCode(any(e))
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

func (s *LinkedHashSet[T]) table() *orderedTable[T, struct{}] {
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
				orderedTable: orderedTable[T, struct{}]{
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
		h += ObjectHashCode(any(s.m.slots[i].key))
	}
	return h
}
