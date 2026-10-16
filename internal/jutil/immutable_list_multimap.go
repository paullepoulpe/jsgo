package jutil

import (
	"iter"
	"slices"
	"strings"
)

// Go-only: language-adapter — Guava ImmutableListMultimap (what
// ImmutableMultimap.builder() builds) and its Builder, including
// orderKeysBy: keys in insertion order (or sorted by the key comparator),
// each key's values in insertion order, duplicates kept. Keys are compared
// with Java equals()/hashCode() (ObjectsEquals / ObjectHashCode).
type ListMultimap[K any, V any] struct {
	m    *ObjectLinkedHashMap[K, []V]
	size int
}

// ImmutableListMultimapBuilder is ImmutableMultimap.Builder /
// ImmutableListMultimap.Builder.
type ImmutableListMultimapBuilder[K any, V any] struct {
	m             *ObjectLinkedHashMap[K, []V]
	size          int
	keyComparator func(a, b K) int
}

// NewImmutableListMultimapBuilder is ImmutableMultimap.builder().
func NewImmutableListMultimapBuilder[K any, V any]() *ImmutableListMultimapBuilder[K, V] {
	return &ImmutableListMultimapBuilder[K, V]{m: NewObjectLinkedHashMap[K, []V]()}
}

// OrderKeysBy is Builder.orderKeysBy(Comparator): the built multimap's keys
// are sorted (stably) by cmp.
func (b *ImmutableListMultimapBuilder[K, V]) OrderKeysBy(cmp func(a, b K) int) *ImmutableListMultimapBuilder[K, V] {
	b.keyComparator = cmp
	return b
}

// Put is Builder.put(K, V).
func (b *ImmutableListMultimapBuilder[K, V]) Put(k K, v V) *ImmutableListMultimapBuilder[K, V] {
	// Java: CollectPreconditions.checkEntryNotNull(key, value)
	if IsNil(k) {
		panic(NewNullPointerExceptionMsg("null key in entry: null=" + ToJavaString(v)))
	} else if IsNil(v) {
		panic(NewNullPointerExceptionMsg("null value in entry: " + ToJavaString(k) + "=null"))
	}
	b.m.Put(k, append(b.m.Get(k), v))
	b.size++
	return b
}

// Build is Builder.build().
func (b *ImmutableListMultimapBuilder[K, V]) Build() *ListMultimap[K, V] {
	keys := b.m.KeySlice()
	if b.keyComparator != nil {
		slices.SortStableFunc(keys, b.keyComparator)
	}
	out := &ListMultimap[K, V]{m: NewObjectLinkedHashMap[K, []V](), size: b.size}
	for _, k := range keys {
		out.m.Put(k, slices.Clone(b.m.Get(k)))
	}
	return out
}

// Get is ListMultimap.get: the values for k (empty when absent). Callers
// must not mutate the returned slice.
func (mm *ListMultimap[K, V]) Get(k K) []V { return mm.m.Get(k) }

// ContainsKey is Multimap.containsKey.
func (mm *ListMultimap[K, V]) ContainsKey(k K) bool { return mm.m.ContainsKey(k) }

// Size is Multimap.size(): the number of key-value pairs.
func (mm *ListMultimap[K, V]) Size() int { return mm.size }

// IsEmpty is Multimap.isEmpty().
func (mm *ListMultimap[K, V]) IsEmpty() bool { return mm.size == 0 }

// KeySet iterates the distinct keys in order (Multimap.keySet()).
func (mm *ListMultimap[K, V]) KeySet() iter.Seq[K] { return mm.m.Keys() }

// Entries iterates key-value pairs grouped by key (Multimap.entries()).
func (mm *ListMultimap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, vs := range mm.m.All() {
			for _, v := range vs {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// AsMap is Multimap.asMap(): each key to its values, in key order. Callers
// must not mutate it.
func (mm *ListMultimap[K, V]) AsMap() *ObjectLinkedHashMap[K, []V] { return mm.m }

// AsMapObject is AsMap() as an Object.
func (mm *ListMultimap[K, V]) AsMapObject() any { return mm.AsMap() }

// String renders like Multimap.toString(): "{k1=[v1, v2], k2=[v3]}".
func (mm *ListMultimap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, vs := range mm.m.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(ToJavaString(k))
		sb.WriteByte('=')
		sb.WriteString(CollectionToString(vs))
	}
	sb.WriteByte('}')
	return sb.String()
}
