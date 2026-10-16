package jutil

// Go-only: language-adapter — the hash-keyed builds of Guava's
// com.google.common.collect.MultimapBuilder that ported code uses.
// MultimapBuilder.hashKeys().arrayListValues() and
// .hashKeys().linkedHashSetValues() are not here: Closure only looks those
// up by key, so their key order is never observable and the insertion-ordered
// jutil multimaps stand for them (see PORTING.md).

// MultimapBuilder_HashKeys_HashSetValues is
// MultimapBuilder.hashKeys().hashSetValues().build(): a backing map sized for
// 8 keys, value sets for 2 values.
func MultimapBuilder_HashKeys_HashSetValues[K any, V any]() *HashMultimap[K, V] {
	return newHashMultimap[K, V](8, hashMultimap_DEFAULT_VALUES_PER_KEY)
}

// MultimapBuilder_HashKeysN_HashSetValuesN is
// MultimapBuilder.hashKeys(expectedKeys).hashSetValues(expectedValuesPerKey)
// .build().
func MultimapBuilder_HashKeysN_HashSetValuesN[K any, V any](expectedKeys, expectedValuesPerKey int) *HashMultimap[K, V] {
	collectPreconditionsCheckNonnegative(expectedKeys, "expectedKeys")
	collectPreconditionsCheckNonnegative(expectedValuesPerKey, "expectedValuesPerKey")
	return newHashMultimap[K, V](expectedKeys, expectedValuesPerKey)
}
