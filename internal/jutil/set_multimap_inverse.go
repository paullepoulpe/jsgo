package jutil

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
