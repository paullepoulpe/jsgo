package jutil

import (
	"iter"
	"strings"
)

// Go-only: language-adapter — Guava's HashBiMap and ImmutableBiMap. Both
// iterate in key insertion order (Guava 33's linked HashBiMap), and the
// inverse view iterates the same entries with key and value swapped. Keys
// and values use Go equality (like LinkedHashMap).

type biMapNode struct {
	kv         [2]any
	prev, next *biMapNode
}

type biMapList struct {
	head, tail *biMapNode
}

// BiMap is com.google.common.collect.BiMap as implemented by HashBiMap and
// ImmutableBiMap (immutability is not enforced).
type BiMap[K, V comparable] struct {
	byKey   map[K]*biMapNode
	byValue map[V]*biMapNode
	list    *biMapList
	side    int // index in biMapNode.kv of this view's key
	inverse *BiMap[V, K]
}

// HashBiMap_Create is HashBiMap.create().
func HashBiMap_Create[K, V comparable]() *BiMap[K, V] {
	m := &BiMap[K, V]{byKey: map[K]*biMapNode{}, byValue: map[V]*biMapNode{}, list: &biMapList{}}
	m.inverse = &BiMap[V, K]{byKey: m.byValue, byValue: m.byKey, list: m.list, side: 1, inverse: m}
	return m
}

// ImmutableBiMap_CopyOf is ImmutableBiMap.copyOf(Map): throws
// IllegalArgumentException if two keys map to the same value.
func ImmutableBiMap_CopyOf[K, V comparable](src *LinkedHashMap[K, V]) *BiMap[K, V] {
	m := HashBiMap_Create[K, V]()
	for k, v := range src.All() {
		if prev, ok := m.byValue[v]; ok {
			panic(NewIllegalArgumentExceptionMsg("Multiple entries with same value: " +
				ToJavaString(k) + "=" + ToJavaString(v) + " and " +
				ToJavaString(prev.kv[0]) + "=" + ToJavaString(prev.kv[1])))
		}
		m.Put(k, v)
	}
	return m
}

func (m *BiMap[K, V]) key(n *biMapNode) K   { return n.kv[m.side].(K) }
func (m *BiMap[K, V]) value(n *biMapNode) V { return n.kv[1-m.side].(V) }

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
	if n, ok := m.byKey[k]; ok {
		return m.value(n), true
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
	oldForKey, hasKey := m.byKey[k]
	if hasKey && m.value(oldForKey) == v {
		return v
	}
	if oldForValue, ok := m.byValue[v]; ok {
		if !force {
			panic(NewIllegalArgumentExceptionMsg("value already present: " + ToJavaString(v)))
		}
		m.delete(oldForValue)
	}
	n := &biMapNode{}
	n.kv[m.side] = k
	n.kv[1-m.side] = v
	if hasKey {
		old := m.value(oldForKey)
		delete(m.byValue, old)
		// Take over the position of the replaced entry.
		n.prev, n.next = oldForKey.prev, oldForKey.next
		if n.prev != nil {
			n.prev.next = n
		} else {
			m.list.head = n
		}
		if n.next != nil {
			n.next.prev = n
		} else {
			m.list.tail = n
		}
		oldForKey.prev, oldForKey.next = nil, nil
		m.byKey[k] = n
		m.byValue[v] = n
		return old
	}
	n.prev = m.list.tail
	if m.list.tail != nil {
		m.list.tail.next = n
	} else {
		m.list.head = n
	}
	m.list.tail = n
	m.byKey[k] = n
	m.byValue[v] = n
	var zero V
	return zero
}

func (m *BiMap[K, V]) delete(n *biMapNode) {
	delete(m.byKey, m.key(n))
	delete(m.byValue, m.value(n))
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		m.list.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		m.list.tail = n.prev
	}
	n.prev, n.next = nil, nil
}

// Remove is Map.remove; returns the removed value (zero if absent).
func (m *BiMap[K, V]) Remove(k K) V {
	n, ok := m.byKey[k]
	if !ok {
		var zero V
		return zero
	}
	v := m.value(n)
	m.delete(n)
	return v
}

// Clear is Map.clear().
func (m *BiMap[K, V]) Clear() {
	clear(m.byKey)
	clear(m.byValue)
	m.list.head, m.list.tail = nil, nil
}

// Inverse is BiMap.inverse(): a view sharing this map's entries.
func (m *BiMap[K, V]) Inverse() *BiMap[V, K] { return m.inverse }

// All iterates entries in key insertion order.
func (m *BiMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for n := m.list.head; n != nil; {
			next := n.next
			if !yield(m.key(n), m.value(n)) {
				return
			}
			n = next
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
