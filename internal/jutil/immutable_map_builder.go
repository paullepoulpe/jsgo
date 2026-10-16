package jutil

// Go-only: language-adapter — Guava ImmutableMap.Builder.

// ImmutableMapBuilder is Guava's ImmutableMap.Builder: entries keep their
// insertion order, and duplicate keys are rejected when the map is built.
type ImmutableMapBuilder[K comparable, V any] struct {
	keys   []K
	values []V
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
	b.keys = append(b.keys, k)
	b.values = append(b.values, v)
	return b
}

// BuildOrThrow is ImmutableMap.Builder.buildOrThrow() (and build()).
// Guava's RegularImmutableMap.fromEntryArray inserts the entries from the
// last to the first, so the duplicate reported is the last entry (in
// insertion order) whose key is put again later, together with the nearest
// later entry with that key: IllegalArgumentException "Multiple entries
// with same key: k=later and k=earlier".
func (b *ImmutableMapBuilder[K, V]) BuildOrThrow() *LinkedHashMap[K, V] {
	later := make(map[K]V, len(b.keys))
	for i := len(b.keys) - 1; i >= 0; i-- {
		k := b.keys[i]
		if v, ok := later[k]; ok {
			panic(NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
				ToJavaString(k) + "=" + ToJavaString(v) + " and " +
				ToJavaString(k) + "=" + ToJavaString(b.values[i])))
		}
		later[k] = b.values[i]
	}
	m := NewLinkedHashMap[K, V]()
	for i, k := range b.keys {
		m.Put(k, b.values[i])
	}
	return m
}
