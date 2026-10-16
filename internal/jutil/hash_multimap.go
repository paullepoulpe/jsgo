package jutil

import "iter"

// Go-only: language-adapter — Guava's HashMultimap (and the SetMultimap that
// MultimapBuilder.hashKeys().hashSetValues().build() returns), as the Guava
// shaded into closure-compiler.jar implements them: an
// AbstractMapBasedMultimap over a JavaHashMap from key to a JavaHashSet of
// values, both sized as Guava sizes them, so keys, values and entries
// iterate in exactly Java's hash order. Keys and values are hashed with
// ObjectHashCode and compared with ObjectsEquals.

// HashMultimap is com.google.common.collect.HashMultimap.
type HashMultimap[K any, V any] struct {
	m                    *JavaHashMap[K, *JavaHashSet[V]]
	totalSize            int
	expectedValuesPerKey int
}

// hashMultimap_DEFAULT_VALUES_PER_KEY is HashMultimap.DEFAULT_VALUES_PER_KEY
// (and MultimapBuilder's DEFAULT_EXPECTED_VALUES_PER_KEY).
const hashMultimap_DEFAULT_VALUES_PER_KEY = 2

func newHashMultimap[K any, V any](expectedKeys, expectedValuesPerKey int) *HashMultimap[K, V] {
	return &HashMultimap[K, V]{
		m:                    Maps_NewHashMapWithExpectedSize[K, *JavaHashSet[V]](expectedKeys),
		expectedValuesPerKey: expectedValuesPerKey,
	}
}

// HashMultimap_Create is HashMultimap.create(): a backing map sized for 12
// keys, value sets for 2 values.
func HashMultimap_Create[K any, V any]() *HashMultimap[K, V] {
	return newHashMultimap[K, V](12, hashMultimap_DEFAULT_VALUES_PER_KEY)
}

// HashMultimap_CreateWithExpectedSizes is HashMultimap.create(int
// expectedKeys, int expectedValuesPerKey).
func HashMultimap_CreateWithExpectedSizes[K any, V any](expectedKeys, expectedValuesPerKey int) *HashMultimap[K, V] {
	mm := newHashMultimap[K, V](expectedKeys, hashMultimap_DEFAULT_VALUES_PER_KEY)
	CheckArgument(expectedValuesPerKey >= 0)
	mm.expectedValuesPerKey = expectedValuesPerKey
	return mm
}

// MultimapSource is the part of a Guava Multimap that HashMultimap.create(
// Multimap) reads: the number of distinct keys and the entries in iteration
// order.
type MultimapSource[K any, V any] interface {
	KeySetSize() int
	Entries() iter.Seq2[K, V]
}

// HashMultimap_CreateFrom is HashMultimap.create(Multimap): a backing map
// sized for the source's keys, then putAll(multimap) in its entries() order.
func HashMultimap_CreateFrom[K any, V any](src MultimapSource[K, V]) *HashMultimap[K, V] {
	mm := newHashMultimap[K, V](src.KeySetSize(), hashMultimap_DEFAULT_VALUES_PER_KEY)
	mm.PutAllMultimap(src)
	return mm
}

// createCollection is HashMultimap.createCollection():
// Sets.newHashSetWithExpectedSize(expectedValuesPerKey).
func (mm *HashMultimap[K, V]) createCollection() *JavaHashSet[V] {
	return Sets_NewHashSetWithExpectedSize[V](mm.expectedValuesPerKey)
}

// Size is Multimap.size(): the number of key-value pairs.
func (mm *HashMultimap[K, V]) Size() int { return mm.totalSize }

// IsEmpty is Multimap.isEmpty().
func (mm *HashMultimap[K, V]) IsEmpty() bool { return mm.totalSize == 0 }

// KeySetSize is keySet().size(): the number of distinct keys.
func (mm *HashMultimap[K, V]) KeySetSize() int { return mm.m.Size() }

// ContainsKey is Multimap.containsKey(Object).
func (mm *HashMultimap[K, V]) ContainsKey(key K) bool { return mm.m.ContainsKey(key) }

// ContainsValue is Multimap.containsValue(Object).
func (mm *HashMultimap[K, V]) ContainsValue(value V) bool {
	for vs := range mm.m.Values() {
		if vs.Contains(value) {
			return true
		}
	}
	return false
}

// ContainsEntry is Multimap.containsEntry(Object, Object).
func (mm *HashMultimap[K, V]) ContainsEntry(key K, value V) bool {
	vs := mm.m.Get(key)
	return vs != nil && vs.Contains(value)
}

// Put is Multimap.put(K, V); reports whether the multimap changed. A new key
// gets a new value set, which is added to the backing map after the value.
func (mm *HashMultimap[K, V]) Put(key K, value V) bool {
	vs := mm.m.Get(key)
	if vs == nil {
		vs = mm.createCollection()
		vs.Add(value)
		mm.totalSize++
		mm.m.Put(key, vs)
		return true
	}
	if vs.Add(value) {
		mm.totalSize++
		return true
	}
	return false
}

// PutAll is Multimap.putAll(K, Iterable): adds each value under key;
// reports whether the multimap changed.
func (mm *HashMultimap[K, V]) PutAll(key K, values iter.Seq[V]) bool {
	vs := mm.m.Get(key)
	isNew := vs == nil
	if isNew {
		vs = mm.createCollection()
	}
	oldSize := vs.Size()
	changed := vs.AddAll(values)
	if changed {
		mm.totalSize += vs.Size() - oldSize
		if isNew {
			mm.m.Put(key, vs)
		}
	}
	return changed
}

// PutAllMultimap is Multimap.putAll(Multimap): put for each entry in order.
func (mm *HashMultimap[K, V]) PutAllMultimap(src MultimapSource[K, V]) bool {
	changed := false
	for k, v := range src.Entries() {
		if mm.Put(k, v) {
			changed = true
		}
	}
	return changed
}

// Remove is Multimap.remove(Object, Object); a key whose last value is
// removed is removed from the backing map.
func (mm *HashMultimap[K, V]) Remove(key K, value V) bool {
	vs := mm.m.Get(key)
	if vs == nil || !vs.Remove(value) {
		return false
	}
	mm.totalSize--
	if vs.IsEmpty() {
		mm.m.Remove(key)
	}
	return true
}

// RemoveAll is Multimap.removeAll(Object): removes the key and returns its
// values in a new set (sized for expectedValuesPerKey and filled in the old
// set's order, so its order can differ from the old set's). Callers must not
// mutate the result (Java returns it unmodifiable).
func (mm *HashMultimap[K, V]) RemoveAll(key K) *JavaHashSet[V] {
	vs := mm.m.Remove(key)
	if vs == nil {
		return NewJavaHashSet[V]()
	}
	out := mm.createCollection()
	out.AddAll(vs.All())
	mm.totalSize -= vs.Size()
	vs.Clear()
	return out
}

// Clear is Multimap.clear(). The backing map keeps its capacity.
func (mm *HashMultimap[K, V]) Clear() {
	for vs := range mm.m.Values() {
		vs.Clear()
	}
	mm.m.Clear()
	mm.totalSize = 0
}

// Get is SetMultimap.get(K): the values for key (a new empty set when
// absent). Callers must not mutate the returned set (Java returns a view
// that writes through to the multimap).
func (mm *HashMultimap[K, V]) Get(key K) *JavaHashSet[V] {
	if vs := mm.m.Get(key); vs != nil {
		return vs
	}
	return mm.createCollection()
}

// KeySet iterates the distinct keys (Multimap.keySet()).
func (mm *HashMultimap[K, V]) KeySet() iter.Seq[K] { return mm.m.Keys() }

// Keys iterates the keys once per value (Multimap.keys()), grouped by key.
func (mm *HashMultimap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range mm.Entries() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values iterates every value (Multimap.values()), grouped by key.
func (mm *HashMultimap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range mm.Entries() {
			if !yield(v) {
				return
			}
		}
	}
}

// Entries iterates the key-value pairs (Multimap.entries()): keys in the
// backing map's order, each key's values in its set's order.
func (mm *HashMultimap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, vs := range mm.m.All() {
			for v := range vs.All() {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// AsMap is Multimap.asMap(): a live view from each key to its values.
// Callers must not mutate it.
func (mm *HashMultimap[K, V]) AsMap() *JavaHashMap[K, *JavaHashSet[V]] { return mm.m }

// AsMapObject is AsMap() as an Object.
func (mm *HashMultimap[K, V]) AsMapObject() any { return mm.m }

// Equals is Multimap.equals(Object): asMap().equals(other.asMap()).
func (mm *HashMultimap[K, V]) Equals(o any) bool {
	if any(mm) == o {
		return true
	}
	other, ok := o.(JavaMultimap)
	if !ok || IsNil(o) {
		return false
	}
	return mm.m.Equals(other.AsMapObject())
}

// HashCode is Multimap.hashCode(): asMap().hashCode().
func (mm *HashMultimap[K, V]) HashCode() int32 { return mm.m.HashCode() }

// String is Multimap.toString(): asMap().toString(), "{k1=[v1, v2], k2=[v3]}".
func (mm *HashMultimap[K, V]) String() string { return mm.m.String() }

var (
	_ JavaMultimap                   = (*HashMultimap[string, string])(nil)
	_ MultimapSource[string, string] = (*HashMultimap[string, string])(nil)
	_ Hashable                       = (*HashMultimap[string, string])(nil)
)
