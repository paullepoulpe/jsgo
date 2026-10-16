package jutil

import "iter"

// RemoveAll is Multimap.removeAll(Object) on a set multimap: removes k and
// returns its values (an empty set when absent).
func (mm *SetMultimap[K, V]) RemoveAll(k K) *ObjectLinkedHashSet[V] {
	vs, ok := mm.m.Lookup(k)
	if !ok {
		return NewObjectLinkedHashSet[V]()
	}
	mm.m.Remove(k)
	mm.size -= vs.Size()
	return vs
}

// PutAll is Multimap.putAll(K, Iterable): puts each value under k; reports
// whether the multimap changed.
func (mm *SetMultimap[K, V]) PutAll(k K, vs iter.Seq[V]) bool {
	changed := false
	for v := range vs {
		if mm.Put(k, v) {
			changed = true
		}
	}
	return changed
}
