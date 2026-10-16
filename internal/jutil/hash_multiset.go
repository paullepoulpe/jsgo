package jutil

// HashMultiset is Guava's com.google.common.collect.HashMultiset, limited
// to the operations ported code uses: membership, counts and removal.
//
// Go-only: language-adapter — no iteration is offered: Java iterates a
// HashMultiset in hash order, which callers must not depend on.
type HashMultiset[T comparable] struct {
	counts *LinkedHashMap[T, int]
	size   int
}

// HashMultiset_Create is HashMultiset.create().
func HashMultiset_Create[T comparable]() *HashMultiset[T] {
	return &HashMultiset[T]{counts: NewLinkedHashMap[T, int]()}
}

// Add is Multiset.add(E): adds one occurrence; always returns true.
func (m *HashMultiset[T]) Add(e T) bool {
	m.counts.Put(e, m.counts.Get(e)+1)
	m.size++
	return true
}

// Count is Multiset.count(Object).
func (m *HashMultiset[T]) Count(e T) int { return m.counts.Get(e) }

// Contains is Multiset.contains(Object).
func (m *HashMultiset[T]) Contains(e T) bool { return m.counts.ContainsKey(e) }

// Size is Multiset.size(): the total number of occurrences.
func (m *HashMultiset[T]) Size() int { return m.size }

// IsEmpty is Collection.isEmpty().
func (m *HashMultiset[T]) IsEmpty() bool { return m.size == 0 }

// Remove is Multiset.remove(Object): removes one occurrence; reports
// whether one was present.
func (m *HashMultiset[T]) Remove(e T) bool {
	c, ok := m.counts.Lookup(e)
	if !ok {
		return false
	}
	if c == 1 {
		m.counts.Remove(e)
	} else {
		m.counts.Put(e, c-1)
	}
	m.size--
	return true
}

// RemoveAll is Multiset.removeAll(Collection): removes every occurrence of
// each element of c (however often it appears in c); reports whether the
// multiset changed.
func (m *HashMultiset[T]) RemoveAll(c []T) bool {
	changed := false
	for _, e := range c {
		if n, ok := m.counts.Lookup(e); ok {
			m.counts.Remove(e)
			m.size -= n
			changed = true
		}
	}
	return changed
}
