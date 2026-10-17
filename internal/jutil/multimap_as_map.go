package jutil

// Go-only: language-adapter — Guava's Multimap as a type (`x instanceof
// Multimap`) and Multimap.asMap().

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
