package jutil

import (
	"iter"
	"strings"
)

// Go-only: language-adapter — Guava's HashBiMap and ImmutableBiMap. Both
// iterate in key insertion order (Guava 33's linked HashBiMap), and the
// inverse view iterates the same entries with key and value swapped. Keys
// and values use Go equality (like LinkedHashMap).

type biMapLink struct {
	prev, next int32
}

type biMapLinks struct {
	nodes              []biMapLink
	head, tail, freeHd int32
}

// BiMap is com.google.common.collect.BiMap as implemented by HashBiMap and
// ImmutableBiMap (immutability is not enforced).
type BiMap[K, V comparable] struct {
	byKey   map[K]int32 // 1-based index into *keys / *values (0 means absent)
	byValue map[V]int32 // 1-based index into *keys / *values (0 means absent)
	keys    *[]K
	values  *[]V
	links   *biMapLinks
	inverse *BiMap[V, K]
}

// HashBiMap_Create is HashBiMap.create().
func HashBiMap_Create[K, V comparable]() *BiMap[K, V] {
	return newBiMapWithCapacity[K, V](0)
}

func newBiMapWithCapacity[K, V comparable](capHint int) *BiMap[K, V] {
	var (
		ks      []K
		vs      []V
		byKey   map[K]int32
		byValue map[V]int32
		nodes   []biMapLink
	)
	if capHint > 0 {
		ks = make([]K, 0, capHint)
		vs = make([]V, 0, capHint)
		byKey = make(map[K]int32, capHint)
		byValue = make(map[V]int32, capHint)
		nodes = make([]biMapLink, 0, capHint)
	} else {
		byKey = map[K]int32{}
		byValue = map[V]int32{}
	}
	links := &biMapLinks{nodes: nodes, head: -1, tail: -1, freeHd: -1}
	m := &BiMap[K, V]{byKey: byKey, byValue: byValue, keys: &ks, values: &vs, links: links}
	m.inverse = &BiMap[V, K]{byKey: byValue, byValue: byKey, keys: &vs, values: &ks, links: links, inverse: m}
	return m
}

// ImmutableBiMap_CopyOf is ImmutableBiMap.copyOf(Map): throws
// IllegalArgumentException if two keys map to the same value.
func ImmutableBiMap_CopyOf[K, V comparable](src *LinkedHashMap[K, V]) *BiMap[K, V] {
	m := newBiMapWithCapacity[K, V](src.Size())
	for k, v := range src.All() {
		if prevSlot, ok := m.byValue[v]; ok {
			prevIdx := prevSlot - 1
			panic(NewIllegalArgumentExceptionMsg("Multiple entries with same value: " +
				ToJavaString(k) + "=" + ToJavaString(v) + " and " +
				ToJavaString((*m.keys)[prevIdx]) + "=" + ToJavaString((*m.values)[prevIdx])))
		}
		m.Put(k, v)
	}
	return m
}

// Size is Map.size().
func (m *BiMap[K, V]) Size() int { return len(m.byKey) }

// IsEmpty is Map.isEmpty().
func (m *BiMap[K, V]) IsEmpty() bool { return len(m.byKey) == 0 }

// Get is Map.get: the zero value (Java null) when absent.
func (m *BiMap[K, V]) Get(k K) V {
	v, _ := m.Lookup(k)
	return v
}

// Lookup returns the value and whether the key is present.
func (m *BiMap[K, V]) Lookup(k K) (V, bool) {
	if slot, ok := m.byKey[k]; ok {
		return (*m.values)[slot-1], true
	}
	var zero V
	return zero, false
}

// ContainsKey is Map.containsKey.
func (m *BiMap[K, V]) ContainsKey(k K) bool {
	_, ok := m.byKey[k]
	return ok
}

// ContainsValue is Map.containsValue.
func (m *BiMap[K, V]) ContainsValue(v V) bool {
	_, ok := m.byValue[v]
	return ok
}

// Put is BiMap.put: throws IllegalArgumentException if v is already bound
// to a different key. Replacing the value of an existing key keeps the
// key's position. Returns the previous value (zero if none).
func (m *BiMap[K, V]) Put(k K, v V) V { return m.put(k, v, false) }

// ForcePut is BiMap.forcePut: an existing entry with value v is removed
// first.
func (m *BiMap[K, V]) ForcePut(k K, v V) V { return m.put(k, v, true) }

func (m *BiMap[K, V]) put(k K, v V, force bool) V {
	oldKeySlot, hasKey := m.byKey[k]
	if hasKey && (*m.values)[oldKeySlot-1] == v {
		return v
	}
	if oldValSlot, ok := m.byValue[v]; ok {
		if !force {
			panic(NewIllegalArgumentExceptionMsg("value already present: " + ToJavaString(v)))
		}
		m.deleteSlot(oldValSlot - 1)
	}
	if hasKey {
		idx := oldKeySlot - 1
		old := (*m.values)[idx]
		delete(m.byValue, old)
		(*m.values)[idx] = v
		m.byValue[v] = oldKeySlot
		return old
	}
	var idx int32
	if m.links.freeHd >= 0 {
		idx = m.links.freeHd
		m.links.freeHd = m.links.nodes[idx].next
		(*m.keys)[idx] = k
		(*m.values)[idx] = v
		m.links.nodes[idx] = biMapLink{prev: m.links.tail, next: -1}
	} else {
		idx = int32(len(*m.keys))
		*m.keys = append(*m.keys, k)
		*m.values = append(*m.values, v)
		m.links.nodes = append(m.links.nodes, biMapLink{prev: m.links.tail, next: -1})
	}
	if m.links.tail >= 0 {
		m.links.nodes[m.links.tail].next = idx
	} else {
		m.links.head = idx
	}
	m.links.tail = idx
	m.byKey[k] = idx + 1
	m.byValue[v] = idx + 1
	var zero V
	return zero
}

func (m *BiMap[K, V]) deleteSlot(idx int32) {
	delete(m.byKey, (*m.keys)[idx])
	delete(m.byValue, (*m.values)[idx])
	var zeroK K
	var zeroV V
	(*m.keys)[idx] = zeroK
	(*m.values)[idx] = zeroV
	nd := m.links.nodes[idx]
	if nd.prev >= 0 {
		m.links.nodes[nd.prev].next = nd.next
	} else {
		m.links.head = nd.next
	}
	if nd.next >= 0 {
		m.links.nodes[nd.next].prev = nd.prev
	} else {
		m.links.tail = nd.prev
	}
	m.links.nodes[idx] = biMapLink{prev: -2, next: m.links.freeHd}
	m.links.freeHd = idx
}

// Remove is Map.remove; returns the removed value (zero if absent).
func (m *BiMap[K, V]) Remove(k K) V {
	slot, ok := m.byKey[k]
	if !ok {
		var zero V
		return zero
	}
	idx := slot - 1
	v := (*m.values)[idx]
	m.deleteSlot(idx)
	return v
}

// Clear is Map.clear().
func (m *BiMap[K, V]) Clear() {
	clear(m.byKey)
	clear(m.byValue)
	clear(*m.keys)
	clear(*m.values)
	*m.keys = (*m.keys)[:0]
	*m.values = (*m.values)[:0]
	m.links.nodes = m.links.nodes[:0]
	m.links.head, m.links.tail, m.links.freeHd = -1, -1, -1
}

// Inverse is BiMap.inverse(): a view sharing this map's entries.
func (m *BiMap[K, V]) Inverse() *BiMap[V, K] { return m.inverse }

// All iterates entries in key insertion order.
func (m *BiMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for idx := m.links.head; idx >= 0; {
			next := m.links.nodes[idx].next
			if !yield((*m.keys)[idx], (*m.values)[idx]) {
				return
			}
			idx = next
		}
	}
}

// Keys iterates the keys in insertion order.
func (m *BiMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values iterates the values in key insertion order.
func (m *BiMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}

// ToLinkedHashMap copies the entries, in order, into a LinkedHashMap.
func (m *BiMap[K, V]) ToLinkedHashMap() *LinkedHashMap[K, V] {
	out := NewLinkedHashMap[K, V]()
	for k, v := range m.All() {
		out.Put(k, v)
	}
	return out
}

// String renders like AbstractMap.toString(): "{k1=v1, k2=v2}".
func (m *BiMap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, v := range m.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(ToJavaString(k))
		sb.WriteByte('=')
		sb.WriteString(ToJavaString(v))
	}
	sb.WriteByte('}')
	return sb.String()
}

func (*BiMap[K, V]) javaMap() {}
