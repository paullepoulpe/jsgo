package jutil

// Go-only: language-adapter — Guava HashMultiset
// (com.google.common.collect.HashMultiset) for elements whose classes
// override equals()/hashCode(): elements are compared with ObjectsEquals and
// hashed with ObjectHashCode instead of Go == (HashMultiset is the variant
// for comparable elements). Java iterates a HashMultiset in hash order,
// which no port relies on, so no iteration is offered.

// ObjectHashMultiset is com.google.common.collect.HashMultiset.
type ObjectHashMultiset[T any] struct {
	counts *ObjectLinkedHashMap[T, int]
	size   int
}

// ObjectHashMultiset_Create is HashMultiset.create().
func ObjectHashMultiset_Create[T any]() *ObjectHashMultiset[T] {
	return &ObjectHashMultiset[T]{counts: NewObjectLinkedHashMap[T, int]()}
}

// Add is Multiset.add(E): adds one occurrence; always true.
func (m *ObjectHashMultiset[T]) Add(e T) bool {
	m.counts.Put(e, m.counts.Get(e)+1)
	m.size++
	return true
}

// AddOccurrences is Multiset.add(E, int): adds n occurrences (n >= 0) and
// returns the count before the operation.
func (m *ObjectHashMultiset[T]) AddOccurrences(e T, n int) int {
	CheckArgumentMsg1(n >= 0, "occurrences cannot be negative: %s", n)
	old := m.counts.Get(e)
	if n > 0 {
		m.counts.Put(e, old+n)
		m.size += n
	}
	return old
}

// Count is Multiset.count(Object): the number of occurrences of e.
func (m *ObjectHashMultiset[T]) Count(e T) int { return m.counts.Get(e) }

// Contains is Multiset.contains(Object).
func (m *ObjectHashMultiset[T]) Contains(e T) bool { return m.counts.ContainsKey(e) }

// Remove is Multiset.remove(Object): removes a single occurrence of e and
// reports whether there was one.
func (m *ObjectHashMultiset[T]) Remove(e T) bool {
	old, ok := m.counts.Lookup(e)
	if !ok {
		return false
	}
	if old == 1 {
		m.counts.Remove(e)
	} else {
		m.counts.Put(e, old-1)
	}
	m.size--
	return true
}

// Size is Multiset.size(): the total number of occurrences.
func (m *ObjectHashMultiset[T]) Size() int { return m.size }

// IsEmpty is Multiset.isEmpty().
func (m *ObjectHashMultiset[T]) IsEmpty() bool { return m.size == 0 }
