package jutil

// ContainsEntry is Multimap.containsEntry(Object, Object): whether k is
// mapped to v.
func (mm *SetMultimap[K, V]) ContainsEntry(k K, v V) bool {
	vs := mm.m.Get(k)
	return vs != nil && vs.Contains(v)
}
