package jutil

// KeySlice returns the keys in ascending order (a snapshot of keySet()).
func (m *TreeMap[K, V]) KeySlice() []K {
	out := make([]K, len(m.entries))
	for i, e := range m.entries {
		out[i] = e.key
	}
	return out
}

// Values iterates values in ascending key order (TreeMap.values()).
func (m *TreeMap[K, V]) Values() func(yield func(V) bool) {
	return func(yield func(V) bool) {
		for i := 0; i < len(m.entries); i++ {
			if !yield(m.entries[i].value) {
				return
			}
		}
	}
}
