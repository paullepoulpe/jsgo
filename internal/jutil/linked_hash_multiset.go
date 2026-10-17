package jutil

import (
	"iter"
	"strconv"
	"strings"
)

// Go-only: language-adapter — Guava com.google.common.collect.
// LinkedHashMultiset: a multiset whose distinct elements iterate in
// first-insertion order (an element whose last occurrence is removed and
// that is added again moves to the end), each element repeated count times
// in a row. Limited to the operations ported code uses.

// LinkedHashMultiset is com.google.common.collect.LinkedHashMultiset.
type LinkedHashMultiset[T comparable] struct {
	counts *LinkedHashMap[T, int]
	size   int
}

// LinkedHashMultiset_Create is LinkedHashMultiset.create().
func LinkedHashMultiset_Create[T comparable]() *LinkedHashMultiset[T] {
	return &LinkedHashMultiset[T]{counts: NewLinkedHashMap[T, int]()}
}

// Add is Multiset.add(E): adds one occurrence; always returns true.
func (m *LinkedHashMultiset[T]) Add(e T) bool {
	m.counts.Put(e, m.counts.Get(e)+1)
	m.size++
	return true
}

// Remove is Multiset.remove(Object): removes one occurrence; reports
// whether one was present.
func (m *LinkedHashMultiset[T]) Remove(e T) bool {
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

// Count is Multiset.count(Object).
func (m *LinkedHashMultiset[T]) Count(e T) int { return m.counts.Get(e) }

// Contains is Multiset.contains(Object).
func (m *LinkedHashMultiset[T]) Contains(e T) bool { return m.counts.ContainsKey(e) }

// Size is Multiset.size(): the total number of occurrences.
func (m *LinkedHashMultiset[T]) Size() int { return m.size }

// IsEmpty is Collection.isEmpty().
func (m *LinkedHashMultiset[T]) IsEmpty() bool { return m.size == 0 }

// All is Multiset.iterator(): each distinct element, in first-insertion
// order, count times in a row.
func (m *LinkedHashMultiset[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e, c := range m.counts.All() {
			for i := 0; i < c; i++ {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// EntrySet iterates the (element, count) pairs of Multiset.entrySet().
func (m *LinkedHashMultiset[T]) EntrySet() iter.Seq2[T, int] { return m.counts.All() }

// String is Multiset.toString(): "[a x 2, b]".
func (m *LinkedHashMultiset[T]) String() string {
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for e, c := range m.counts.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(ToJavaString(e))
		if c != 1 {
			sb.WriteString(" x ")
			sb.WriteString(strconv.Itoa(c))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}
