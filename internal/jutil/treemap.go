package jutil

import (
	"iter"
	"sort"
	"strings"
)

// Go-only: language-adapter — java.util.TreeMap.

type treeEntry[K any, V any] struct {
	key   K
	value V
}

// TreeMap is java.util.TreeMap: a map whose iteration order is the key order
// given by a comparator (Java's natural ordering is passed explicitly as the
// key type's compareTo). It is backed by a sorted slice, which has the same
// observable behaviour.
type TreeMap[K any, V any] struct {
	compare func(a, b K) int
	entries []treeEntry[K, V]
}

// NewTreeMap is new TreeMap<>(comparator).
func NewTreeMap[K any, V any](compare func(a, b K) int) *TreeMap[K, V] {
	return &TreeMap[K, V]{compare: compare}
}

func (m *TreeMap[K, V]) search(k K) (int, bool) {
	i := sort.Search(len(m.entries), func(i int) bool { return m.compare(m.entries[i].key, k) >= 0 })
	return i, i < len(m.entries) && m.compare(m.entries[i].key, k) == 0
}

func (m *TreeMap[K, V]) Size() int     { return len(m.entries) }
func (m *TreeMap[K, V]) IsEmpty() bool { return len(m.entries) == 0 }

// Get is Map.get: the zero value (Java null) when absent.
func (m *TreeMap[K, V]) Get(k K) V {
	if i, ok := m.search(k); ok {
		return m.entries[i].value
	}
	var zero V
	return zero
}

// Lookup returns the value and whether the key is present.
func (m *TreeMap[K, V]) Lookup(k K) (V, bool) {
	if i, ok := m.search(k); ok {
		return m.entries[i].value, true
	}
	var zero V
	return zero, false
}

// ContainsKey is Map.containsKey.
func (m *TreeMap[K, V]) ContainsKey(k K) bool {
	_, ok := m.search(k)
	return ok
}

// Put is Map.put; returns the previous value (zero if none).
func (m *TreeMap[K, V]) Put(k K, v V) V {
	i, ok := m.search(k)
	if ok {
		old := m.entries[i].value
		m.entries[i].value = v
		return old
	}
	m.entries = append(m.entries, treeEntry[K, V]{})
	copy(m.entries[i+1:], m.entries[i:])
	m.entries[i] = treeEntry[K, V]{key: k, value: v}
	var zero V
	return zero
}

// PutIfAbsent is Map.putIfAbsent: if the key is absent or mapped to null it
// is (re)mapped to v and null is returned; otherwise the current value is
// returned.
func (m *TreeMap[K, V]) PutIfAbsent(k K, v V) V {
	i, ok := m.search(k)
	if ok && !IsNil(m.entries[i].value) {
		return m.entries[i].value
	}
	m.Put(k, v)
	var zero V
	return zero
}

// Remove is Map.remove; returns the removed value (zero if absent).
func (m *TreeMap[K, V]) Remove(k K) V {
	i, ok := m.search(k)
	if !ok {
		var zero V
		return zero
	}
	old := m.entries[i].value
	m.entries = append(m.entries[:i], m.entries[i+1:]...)
	return old
}

// Clear is Map.clear.
func (m *TreeMap[K, V]) Clear() { m.entries = nil }

// All iterates entries in ascending key order. Values replaced with Put
// during iteration are seen (Java: Map.Entry.setValue); structural changes
// during iteration are not supported (Java throws
// ConcurrentModificationException).
func (m *TreeMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := 0; i < len(m.entries); i++ {
			if !yield(m.entries[i].key, m.entries[i].value) {
				return
			}
		}
	}
}

// Keys iterates keys in ascending order.
//
// Go-only: performance — walks the entries itself (see
// LinkedHashMap.Keys).
func (m *TreeMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for i := 0; i < len(m.entries); i++ {
			if !yield(m.entries[i].key) {
				return
			}
		}
	}
}

// Copy is new TreeMap<>(sortedMap).
func (m *TreeMap[K, V]) Copy() *TreeMap[K, V] {
	return &TreeMap[K, V]{compare: m.compare, entries: append([]treeEntry[K, V](nil), m.entries...)}
}

// String renders like AbstractMap.toString().
func (m *TreeMap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	for i, e := range m.entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ToJavaString(e.key))
		sb.WriteByte('=')
		sb.WriteString(ToJavaString(e.value))
	}
	sb.WriteByte('}')
	return sb.String()
}
