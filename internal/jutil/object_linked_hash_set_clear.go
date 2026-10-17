package jutil

// Clear is Set.clear() for ObjectLinkedHashSet.
func (s *ObjectLinkedHashSet[T]) Clear() {
	s.m = ObjectLinkedHashMap[T, struct{}]{}
}

// Clear is Map.clear() for ObjectLinkedHashMap.
func (h *ObjectLinkedHashMap[K, V]) Clear() {
	*h = *NewObjectLinkedHashMap[K, V]()
}
