package jutil

// ArrayListMultimap_Create is ArrayListMultimap.create() and
// MultimapBuilder.hashKeys().arrayListValues().build(): a mutable
// ListMultimap whose values for a key keep insertion order and duplicates.
//
// Go-only: language-adapter — keys iterate in insertion order where Java
// iterates a hash-keyed multimap in hash order; callers must not depend on
// key order.
func ArrayListMultimap_Create[K any, V any]() *ListMultimap[K, V] {
	return &ListMultimap[K, V]{m: NewObjectLinkedHashMap[K, []V]()}
}

// Put is Multimap.put(K, V) on a list multimap: appends v to k's values;
// always returns true.
func (mm *ListMultimap[K, V]) Put(k K, v V) bool {
	mm.m.Put(k, append(mm.m.Get(k), v))
	mm.size++
	return true
}

// RemoveAll is Multimap.removeAll(Object): removes k and returns its
// values (empty when absent).
func (mm *ListMultimap[K, V]) RemoveAll(k K) []V {
	vs, ok := mm.m.Lookup(k)
	if !ok {
		return nil
	}
	mm.m.Remove(k)
	mm.size -= len(vs)
	return vs
}
