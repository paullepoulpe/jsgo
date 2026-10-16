// Port of com/google/javascript/jscomp/base/LinkedIdentityHashMap.java @ closure-compiler v20260415

package base

import "github.com/paullepoulpe/jsgo/internal/jutil"

// LinkedIdentityHashMap is a deterministic drop-in replacement for
// IdentityHashMap.
//
// This class is actually "map-like" in that it doesn't fully implement the
// Map interface. Only the methods needed to eliminate IdentityHashMap were
// defined. The entire interface could be satisfied in the future if
// necessary.
//
// Go-only: language-adapter — the IdentityHashMap is a Go map keyed by K
// with Go ==, which is identity for the pointer types used as keys.
type LinkedIdentityHashMap[K comparable, V any] struct {
	innerMap       map[K]*LinkedIdentityHashMap_Entry[K, V]
	iterationOrder *jutil.ArrayList[*LinkedIdentityHashMap_Entry[K, V]]
}

// NewLinkedIdentityHashMap is the implicit default constructor.
//
// Java: LinkedIdentityHashMap()
func NewLinkedIdentityHashMap[K comparable, V any]() *LinkedIdentityHashMap[K, V] {
	return &LinkedIdentityHashMap[K, V]{
		innerMap:       make(map[K]*LinkedIdentityHashMap_Entry[K, V]),
		iterationOrder: jutil.NewArrayList[*LinkedIdentityHashMap_Entry[K, V]](),
	}
}

// Get returns the value for key, or the zero value (Java null).
//
// Java: LinkedIdentityHashMap.get(K)
func (m *LinkedIdentityHashMap[K, V]) Get(key K) V {
	entry := m.innerMap[key]
	if entry == nil {
		var zero V
		return zero
	}
	return entry.value
}

// Put associates value with key and returns the previous value (zero if
// none).
//
// Java: LinkedIdentityHashMap.put(K, V)
func (m *LinkedIdentityHashMap[K, V]) Put(key K, value V) V {
	entry := m.innerMap[key]
	if entry == nil {
		entry = NewLinkedIdentityHashMap_Entry(key, value)
		m.innerMap[key] = entry
		m.iterationOrder.Add(entry)
		var zero V
		return zero
	}
	oldValue := entry.value
	entry.value = value
	return oldValue
}

// ComputeIfAbsent returns the existing non-null value, or computes, stores
// and returns a new one.
//
// Java: LinkedIdentityHashMap.computeIfAbsent(K, Function<? super K, ? extends V>)
func (m *LinkedIdentityHashMap[K, V]) ComputeIfAbsent(key K, fn func(K) V) V {
	oldValue := m.Get(key)
	if !jutil.IsNil(oldValue) {
		return oldValue
	}

	newValue := fn(key)
	m.Put(key, newValue)
	return newValue
}

// Merge stores value if the key has no (non-null) value, else the result of
// fn(oldValue, value).
//
// Java: LinkedIdentityHashMap.merge(K, V, BiFunction<? super V, ? super V, ? extends V>)
func (m *LinkedIdentityHashMap[K, V]) Merge(key K, value V, fn func(V, V) V) V {
	oldValue := m.Get(key)
	if jutil.IsNil(oldValue) {
		m.Put(key, value)
		return value
	}

	mergeValue := fn(oldValue, value)
	m.Put(key, mergeValue)
	return mergeValue
}

// ForEach calls fn for each entry in insertion order.
//
// Java: LinkedIdentityHashMap.forEach(BiConsumer<? super K, ? super V>)
func (m *LinkedIdentityHashMap[K, V]) ForEach(fn func(K, V)) {
	for entry := range m.iterationOrder.All() {
		fn(entry.key, entry.value)
	}
}

// LinkedIdentityHashMap_Entry is a key/value pair.
type LinkedIdentityHashMap_Entry[K any, V any] struct {
	key   K
	value V
}

// Java: LinkedIdentityHashMap.Entry(K, V)
func NewLinkedIdentityHashMap_Entry[K any, V any](key K, value V) *LinkedIdentityHashMap_Entry[K, V] {
	return &LinkedIdentityHashMap_Entry[K, V]{key: key, value: value}
}
