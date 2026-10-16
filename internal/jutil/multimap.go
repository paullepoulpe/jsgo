package jutil

import "iter"

// Go-only: language-adapter — Guava SetMultimap built with
// MultimapBuilder.linkedHashKeys().linkedHashSetValues(), and
// ImmutableSetMultimap: keys and values in insertion order, compared with
// Java equals()/hashCode() (ObjectsEquals / ObjectHashCode).
type SetMultimap[K any, V any] struct {
	m    *ObjectLinkedHashMap[K, *ObjectLinkedHashSet[V]]
	size int
}

// NewSetMultimap returns an empty multimap.
func NewSetMultimap[K any, V any]() *SetMultimap[K, V] {
	return &SetMultimap[K, V]{m: NewObjectLinkedHashMap[K, *ObjectLinkedHashSet[V]]()}
}

// SetMultimapCopyOf is ImmutableSetMultimap.copyOf(multimap).
func SetMultimapCopyOf[K any, V any](src *SetMultimap[K, V]) *SetMultimap[K, V] {
	out := NewSetMultimap[K, V]()
	for k, vs := range src.m.All() {
		for v := range vs.All() {
			out.Put(k, v)
		}
	}
	return out
}

// Put is Multimap.put; reports whether the multimap changed.
func (mm *SetMultimap[K, V]) Put(k K, v V) bool {
	vs := mm.m.Get(k)
	if vs == nil {
		vs = NewObjectLinkedHashSet[V]()
		mm.m.Put(k, vs)
	}
	if vs.Add(v) {
		mm.size++
		return true
	}
	return false
}

// Get is SetMultimap.get: the values for k (an empty set when absent).
// Callers must not mutate the returned set.
func (mm *SetMultimap[K, V]) Get(k K) *ObjectLinkedHashSet[V] {
	if vs := mm.m.Get(k); vs != nil {
		return vs
	}
	return NewObjectLinkedHashSet[V]()
}

// ContainsKey is Multimap.containsKey.
func (mm *SetMultimap[K, V]) ContainsKey(k K) bool { return mm.m.ContainsKey(k) }

// Size is Multimap.size(): the number of key-value pairs.
func (mm *SetMultimap[K, V]) Size() int { return mm.size }

// IsEmpty is Multimap.isEmpty().
func (mm *SetMultimap[K, V]) IsEmpty() bool { return mm.size == 0 }

// KeySet iterates the distinct keys in insertion order (Multimap.keySet()).
func (mm *SetMultimap[K, V]) KeySet() iter.Seq[K] { return mm.m.Keys() }

// Entries iterates key-value pairs grouped by key (Multimap.entries()).
func (mm *SetMultimap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, vs := range mm.m.All() {
			for v := range vs.All() {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// String renders like Multimap.toString(): "{k1=[v1, v2], k2=[v3]}".
func (mm *SetMultimap[K, V]) String() string { return mm.m.String() }
