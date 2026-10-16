package jutil

// Clear is Multimap.clear().
func (mm *SetMultimap[K, V]) Clear() {
	mm.m = NewObjectLinkedHashMap[K, *ObjectLinkedHashSet[V]]()
	mm.size = 0
}
