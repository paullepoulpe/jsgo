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

// BuildOrThrow is ImmutableMap.Builder.buildOrThrow() (and build()): the
// first entry whose key was already put throws IllegalArgumentException
// "Multiple entries with same key: k=new and k=old".
func (b *ImmutableMapBuilder[K, V]) BuildOrThrow() *LinkedHashMap[K, V] {
	m := NewLinkedHashMap[K, V]()
	for i, k := range b.keys {
		if old, ok := m.Lookup(k); ok {
			panic(NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
				ToJavaString(k) + "=" + ToJavaString(b.values[i]) + " and " +
				ToJavaString(k) + "=" + ToJavaString(old)))
		}
		m.Put(k, b.values[i])
	}
	return m
}
