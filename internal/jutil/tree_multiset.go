package jutil

import (
	"iter"
	"strconv"
	"strings"
)

// Go-only: language-adapter — Guava's TreeMultiset, as created by
// TreeMultiset.create() (natural order) or with a comparator: a multiset
// whose distinct elements are kept sorted. Iteration visits each distinct
// element in order, repeated as many times as it occurs; size() counts
// every occurrence; toString() prints "[a x 2, b]".

// TreeMultiset is Guava's TreeMultiset.
type TreeMultiset[E any] struct {
	counts *TreeMap[E, int]
	size   int
}

// TreeMultiset_Create is TreeMultiset.create(Comparator) (pass the natural
// order's comparison for TreeMultiset.create()).
func TreeMultiset_Create[E any](compare func(a, b E) int) *TreeMultiset[E] {
	return &TreeMultiset[E]{counts: NewTreeMap[E, int](compare)}
}

// Add is Multiset.add(E): it adds one occurrence and returns true.
func (m *TreeMultiset[E]) Add(e E) bool {
	m.counts.Put(e, m.counts.Get(e)+1)
	m.size++
	return true
}

// AddAll is Collection.addAll(Collection): it adds one occurrence of each
// element and reports whether the multiset changed.
func (m *TreeMultiset[E]) AddAll(es iter.Seq[E]) bool {
	changed := false
	for e := range es {
		m.Add(e)
		changed = true
	}
	return changed
}

// Count is Multiset.count(Object).
func (m *TreeMultiset[E]) Count(e E) int { return m.counts.Get(e) }

// Contains is Multiset.contains(Object).
func (m *TreeMultiset[E]) Contains(e E) bool { return m.counts.ContainsKey(e) }

// Size is Multiset.size(): the number of occurrences.
func (m *TreeMultiset[E]) Size() int { return m.size }

// IsEmpty is Multiset.isEmpty().
func (m *TreeMultiset[E]) IsEmpty() bool { return m.size == 0 }

// ElementSet is Multiset.elementSet(): the distinct elements in order.
func (m *TreeMultiset[E]) ElementSet() []E { return m.counts.KeySlice() }

// All iterates like Multiset.iterator(): each distinct element in order,
// once per occurrence.
func (m *TreeMultiset[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for e, n := range m.counts.All() {
			for i := 0; i < n; i++ {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// ToSlice lists the occurrences in iteration order.
func (m *TreeMultiset[E]) ToSlice() []E {
	var out []E
	for e := range m.All() {
		out = append(out, e)
	}
	return out
}

// String is Multiset.toString(): its entrySet, "[e]" or "[e x n]" per
// distinct element.
func (m *TreeMultiset[E]) String() string {
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for e, n := range m.counts.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(ToJavaString(e))
		if n != 1 {
			sb.WriteString(" x " + strconv.Itoa(n))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}
