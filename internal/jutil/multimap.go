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

// GetOrEmpty returns the values for k, or empty when k is absent.
// Callers must not mutate the returned set.
func (mm *SetMultimap[K, V]) GetOrEmpty(k K, empty *ObjectLinkedHashSet[V]) *ObjectLinkedHashSet[V] {
	if vs := mm.m.Get(k); vs != nil {
		return vs
	}
	return empty
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

// Values lists the values of Entries(), in that order (Multimap.values()).
func (mm *SetMultimap[K, V]) Values() []V {
	out := make([]V, 0, mm.size)
	for _, v := range mm.Entries() {
		out = append(out, v)
	}
	return out
}

// String renders like Multimap.toString(): "{k1=[v1, v2], k2=[v3]}".
func (mm *SetMultimap[K, V]) String() string { return mm.m.String() }

// Inverse is ImmutableSetMultimap.inverse(): a multimap with every
// key-value pair swapped, built by putting the entries in entries() order
// (so the inverse's keys are the values in first-occurrence order, each with
// its keys in order).
//
// Go-only: language-adapter — Guava caches the inverse of an immutable
// multimap; here a new multimap is built on each call.
func (mm *SetMultimap[K, V]) Inverse() *SetMultimap[V, K] {
	out := NewSetMultimap[V, K]()
	for k, v := range mm.Entries() {
		out.Put(v, k)
	}
	return out
}

// JavaMultimap is implemented by the jutil types that model a Guava
// Multimap, so that `x instanceof Multimap` can be expressed as a Go type
// assertion.
type JavaMultimap interface {
	Size() int
	IsEmpty() bool
	// AsMapObject is Multimap.asMap() for callers that only know the value
	// is a Multimap (e.g. Gson): the map view as an Object.
	AsMapObject() any
}

// AsMap is Multimap.asMap(): a live view from each key to the set of its
// values, in key insertion order. Callers must not mutate it.
func (mm *SetMultimap[K, V]) AsMap() *ObjectLinkedHashMap[K, *ObjectLinkedHashSet[V]] {
	return mm.m
}

// AsMapObject is AsMap() as an Object.
func (mm *SetMultimap[K, V]) AsMapObject() any { return mm.AsMap() }
