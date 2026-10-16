package jutil

import "iter"

// TreeSet is java.util.TreeSet<E> with a comparator (for String elements,
// StringCompareTo is the natural ordering): a sorted set, backed by a TreeMap
// as in the JDK.
type TreeSet[T any] struct {
	m *TreeMap[T, struct{}]
}

// NewTreeSet is new TreeSet<>(comparator).
func NewTreeSet[T any](compare func(a, b T) int) *TreeSet[T] {
	return &TreeSet[T]{m: NewTreeMap[T, struct{}](compare)}
}

// Add is TreeSet.add(E): reports whether the set did not already contain e.
func (s *TreeSet[T]) Add(e T) bool {
	if s.m.ContainsKey(e) {
		return false
	}
	s.m.Put(e, struct{}{})
	return true
}

// AddAll is TreeSet.addAll(Collection): reports whether the set changed.
func (s *TreeSet[T]) AddAll(es iter.Seq[T]) bool {
	changed := false
	for e := range es {
		if s.Add(e) {
			changed = true
		}
	}
	return changed
}

// Contains is TreeSet.contains(Object).
func (s *TreeSet[T]) Contains(e T) bool { return s.m.ContainsKey(e) }

// Remove is TreeSet.remove(Object): reports whether e was present.
func (s *TreeSet[T]) Remove(e T) bool {
	if !s.m.ContainsKey(e) {
		return false
	}
	s.m.Remove(e)
	return true
}

// Size is TreeSet.size().
func (s *TreeSet[T]) Size() int { return s.m.Size() }

// IsEmpty is TreeSet.isEmpty().
func (s *TreeSet[T]) IsEmpty() bool { return s.m.IsEmpty() }

// All iterates in ascending order.
func (s *TreeSet[T]) All() iter.Seq[T] { return s.m.Keys() }

// String is AbstractCollection.toString().
func (s *TreeSet[T]) String() string {
	var elems []T
	for e := range s.All() {
		elems = append(elems, e)
	}
	return CollectionToString(elems)
}
