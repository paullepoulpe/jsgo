package jscomp

// Go-only: test-helper — small helpers shared by the w5-es-modules test ports.

// esModulesTest_mapOf is ImmutableMap.of(k1, v1, ...) for String keys and
// values.
func esModulesTest_mapOf(kv ...string) *OrderedStringMap[string] {
	m := NewOrderedStringMapWithCapacity[string](len(kv) / 2)
	for i := 0; i < len(kv); i += 2 {
		m.Put(kv[i], kv[i+1])
	}
	return m
}
