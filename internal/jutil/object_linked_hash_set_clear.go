package jutil

// Clear is Set.clear() for ObjectLinkedHashSet.
func (s *ObjectLinkedHashSet[T]) Clear() {
	s.m.Clear()
}

// Clear is Map.clear() for ObjectLinkedHashMap.
func (h *ObjectLinkedHashMap[K, V]) Clear() {
	clear(h.slots)
	h.slots = h.slots[:0]
	if h.buckets != nil {
		clear(h.buckets)
	}
	h.size = 0
	h.first = 0
}

