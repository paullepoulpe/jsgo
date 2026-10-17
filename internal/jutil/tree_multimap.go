package jutil

import (
	"iter"
	"strings"
)

// Go-only: language-adapter — Guava's TreeMultimap, as created by
// TreeMultimap.create(): a SetMultimap whose keys, and the values of each
// key, are kept sorted (here by the given comparisons). put(k, v) adds the
// pair unless it is already present; size() counts pairs; entries() and
// toString() ("{k=[v1, v2]}") follow key order, then value order.

// TreeMultimap is Guava's TreeMultimap.
type TreeMultimap[K any, V any] struct {
	compareValues func(a, b V) int
	m             *TreeMap[K, *TreeSet[V]]
	size          int
}

// TreeMultimap_Create is TreeMultimap.create(Comparator, Comparator) (pass
// the natural orders' comparisons for TreeMultimap.create()).
func TreeMultimap_Create[K any, V any](compareKeys func(a, b K) int, compareValues func(a, b V) int) *TreeMultimap[K, V] {
	return &TreeMultimap[K, V]{compareValues: compareValues, m: NewTreeMap[K, *TreeSet[V]](compareKeys)}
}

// Put is Multimap.put(K, V): it reports whether the pair was added.
func (t *TreeMultimap[K, V]) Put(k K, v V) bool {
	values, ok := t.m.Lookup(k)
	if !ok {
		values = NewTreeSet[V](t.compareValues)
		t.m.Put(k, values)
	}
	if !values.Add(v) {
		return false
	}
	t.size++
	return true
}

// Get is Multimap.get(K): the values of k in order (empty if none).
func (t *TreeMultimap[K, V]) Get(k K) []V {
	if values, ok := t.m.Lookup(k); ok {
		return values.ToSlice()
	}
	return nil
}

// ContainsKey is Multimap.containsKey(Object).
func (t *TreeMultimap[K, V]) ContainsKey(k K) bool { return t.m.ContainsKey(k) }

// ContainsEntry is Multimap.containsEntry(Object, Object).
func (t *TreeMultimap[K, V]) ContainsEntry(k K, v V) bool {
	values, ok := t.m.Lookup(k)
	return ok && values.Contains(v)
}

// Size is Multimap.size(): the number of key-value pairs.
func (t *TreeMultimap[K, V]) Size() int { return t.size }

// IsEmpty is Multimap.isEmpty().
func (t *TreeMultimap[K, V]) IsEmpty() bool { return t.size == 0 }

// KeySet is Multimap.keySet(): the distinct keys in order.
func (t *TreeMultimap[K, V]) KeySet() []K { return t.m.KeySlice() }

// Entries iterates Multimap.entries(): every pair, in key order then value
// order.
func (t *TreeMultimap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, values := range t.m.All() {
			for v := range values.All() {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// String is Multimap.toString(): asMap().toString().
func (t *TreeMultimap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, values := range t.m.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(ToJavaString(k) + "=" + values.String())
	}
	sb.WriteByte('}')
	return sb.String()
}
