package jutil

// ContainsValue is Map.containsValue(Object): whether some entry's value
// equals v (ObjectsEquals), scanning the entries in order.
func (h *LinkedHashMap[K, V]) ContainsValue(v V) bool {
	for e := range h.Values() {
		if ObjectsEquals(any(e), any(v)) {
			return true
		}
	}
	return false
}
