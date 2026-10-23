package jutil

// Go-only: language-adapter — Guava ImmutableMap.Builder.

// ImmutableMapBuilder is Guava's ImmutableMap.Builder: entries keep their
// insertion order, and duplicate keys are rejected when the map is built.
type ImmutableMapBuilder[K comparable, V any] struct {
	slots []LhmSlot[K, V]
}

// NewImmutableMapBuilder is ImmutableMap.builder().
func NewImmutableMapBuilder[K comparable, V any]() *ImmutableMapBuilder[K, V] {
	return &ImmutableMapBuilder[K, V]{}
}

// Put is ImmutableMap.Builder.put(K, V).
func (b *ImmutableMapBuilder[K, V]) Put(k K, v V) *ImmutableMapBuilder[K, V] {
	// Java: CollectPreconditions.checkEntryNotNull(key, value)
	if IsNil(k) {
		panic(NewNullPointerExceptionMsg("null key in entry: null=" + ToJavaString(v)))
	} else if IsNil(v) {
		panic(NewNullPointerExceptionMsg("null value in entry: " + ToJavaString(k) + "=null"))
	}
	b.slots = append(b.slots, LhmSlot[K, V]{key: k, value: v})
	return b
}

// BuildOrThrow is ImmutableMap.Builder.buildOrThrow() (and build()).
// Guava's RegularImmutableMap.fromEntryArray inserts the entries from the
// last to the first, so the duplicate reported is the last entry (in
// insertion order) whose key is put again later, together with the nearest
// later entry with that key: IllegalArgumentException "Multiple entries
// with same key: k=later and k=earlier".
func (b *ImmutableMapBuilder[K, V]) BuildOrThrow() *LinkedHashMap[K, V] {
	n := len(b.slots)
	if n == 0 {
		return &LinkedHashMap[K, V]{}
	}
	slots := make([]LhmSlot[K, V], n)
	copy(slots, b.slots)
	if n <= lhmSmallSize {
		for i := n - 1; i >= 0; i-- {
			k := slots[i].key
			for j := i + 1; j < n; j++ {
				if slots[j].key == k {
					panic(NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
						ToJavaString(k) + "=" + ToJavaString(slots[j].value) + " and " +
						ToJavaString(k) + "=" + ToJavaString(slots[i].value)))
				}
			}
		}
		return &LinkedHashMap[K, V]{
			slots: slots,
			size:  int32(n),
		}
	}
	h := &LinkedHashMap[K, V]{
		m:     &lhmIndex{table: make([]uint32, lhmTableLen(n))},
		slots: slots,
		size:  int32(n),
	}
	for i := n - 1; i >= 0; i-- {
		k := slots[i].key
		if later := h.lookup(k); later >= 0 {
			panic(NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
				ToJavaString(k) + "=" + ToJavaString(slots[later].value) + " and " +
				ToJavaString(k) + "=" + ToJavaString(slots[i].value)))
		}
		h.m.insert(lhmHash(k), i)
	}
	return h
}

// PutAll is ImmutableMap.Builder.putAll(Map): puts each entry of m, in its
// iteration order.
func (b *ImmutableMapBuilder[K, V]) PutAll(m *LinkedHashMap[K, V]) *ImmutableMapBuilder[K, V] {
	for k, v := range m.All() {
		b.Put(k, v)
	}
	return b
}
