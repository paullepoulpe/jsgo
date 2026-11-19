// Clean-room specification based on java.util.LinkedHashMap & java.util.HashMap Javadoc.

package javautil

import (
	"iter"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

const orderedMapIndexThreshold = 8

// MapEntry pairs a key with a value in insertion order.
type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// LinkedHashMap is a deterministic, insertion-ordered map backed by a contiguous
// entry slice and an auxiliary index map populated above 8 elements.
// All iterations iterate over entries in insertion order, guaranteeing 100% determinism.
type LinkedHashMap[K comparable, V any] struct {
	entries []MapEntry[K, V]
	index   map[K]int
	holes   int
}

// HashMap is an alias to LinkedHashMap. In this port, all hash maps use insertion-order
// iteration to eliminate Go runtime map iteration randomization and guarantee byte-level parity.
type HashMap[K comparable, V any] = LinkedHashMap[K, V]

// NewLinkedHashMap creates an empty LinkedHashMap.
func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{}
}

// NewLinkedHashMapWithCapacity creates a LinkedHashMap with pre-allocated capacity.
func NewLinkedHashMapWithCapacity[K comparable, V any](capHint int) *LinkedHashMap[K, V] {
	if capHint <= 0 {
		return &LinkedHashMap[K, V]{}
	}
	m := &LinkedHashMap[K, V]{entries: make([]MapEntry[K, V], 0, capHint)}
	if capHint > orderedMapIndexThreshold {
		m.index = make(map[K]int, capHint)
	}
	return m
}

// NewHashMap creates an empty HashMap with deterministic insertion order.
func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return NewLinkedHashMap[K, V]()
}

// NewHashMapWithCapacity creates a HashMap with pre-allocated capacity.
func NewHashMapWithCapacity[K comparable, V any](capHint int) *HashMap[K, V] {
	return NewLinkedHashMapWithCapacity[K, V](capHint)
}

func (m *LinkedHashMap[K, V]) compact() {
	if m == nil || m.holes == 0 {
		return
	}
	w := 0
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; ok && idx == i {
			m.entries[w] = e
			m.index[e.Key] = w
			w++
		}
	}
	clear(m.entries[w:])
	m.entries = m.entries[:w]
	m.holes = 0
}

// Put inserts or updates k, preserving first-seen insertion order and returning previous value.
func (m *LinkedHashMap[K, V]) Put(k K, v V) V {
	var zero V
	if m.index != nil {
		if idx, ok := m.index[k]; ok {
			prev := m.entries[idx].Value
			m.entries[idx].Value = v
			return prev
		}
		m.index[k] = len(m.entries)
		m.entries = append(m.entries, MapEntry[K, V]{Key: k, Value: v})
		return zero
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			prev := m.entries[i].Value
			m.entries[i].Value = v
			return prev
		}
	}
	if len(m.entries) >= orderedMapIndexThreshold {
		m.index = make(map[K]int, len(m.entries)+1)
		for i, e := range m.entries {
			m.index[e.Key] = i
		}
		m.index[k] = len(m.entries)
		m.entries = append(m.entries, MapEntry[K, V]{Key: k, Value: v})
		return zero
	}
	m.entries = append(m.entries, MapEntry[K, V]{Key: k, Value: v})
	return zero
}

// PutIfAbsent inserts (k, v) only if k is not already present, returning true if inserted.
func (m *LinkedHashMap[K, V]) PutIfAbsent(k K, v V) bool {
	if m.ContainsKey(k) {
		return false
	}
	m.Put(k, v)
	return true
}

// Get looks up k, returning the value and whether it was present.
func (m *LinkedHashMap[K, V]) Get(k K) (V, bool) {
	var zero V
	if m == nil {
		return zero, false
	}
	if m.index != nil {
		idx, ok := m.index[k]
		if !ok {
			return zero, false
		}
		return m.entries[idx].Value, true
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			return m.entries[i].Value, true
		}
	}
	return zero, false
}

// GetVal returns the value for k, or the zero value if not present.
func (m *LinkedHashMap[K, V]) GetVal(k K) V {
	v, _ := m.Get(k)
	return v
}

// ContainsKey reports whether k is mapped in the map.
func (m *LinkedHashMap[K, V]) ContainsKey(k K) bool {
	if m == nil {
		return false
	}
	if m.index != nil {
		_, ok := m.index[k]
		return ok
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			return true
		}
	}
	return false
}

// ContainsValue reports whether any key maps to value v.
func (m *LinkedHashMap[K, V]) ContainsValue(v V) bool {
	if m == nil {
		return false
	}
	if m.holes == 0 {
		for i := range m.entries {
			if javalang.ObjectsEquals(m.entries[i].Value, v) {
				return true
			}
		}
		return false
	}
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; ok && idx == i {
			if javalang.ObjectsEquals(e.Value, v) {
				return true
			}
		}
	}
	return false
}

// Remove removes k from the map, returning the previous value and whether it was present.
func (m *LinkedHashMap[K, V]) Remove(k K) (V, bool) {
	var zero V
	if m == nil {
		return zero, false
	}
	if m.index != nil {
		idx, ok := m.index[k]
		if !ok {
			return zero, false
		}
		delete(m.index, k)
		prev := m.entries[idx].Value
		m.entries[idx] = MapEntry[K, V]{}
		m.holes++
		if m.holes*2 > len(m.entries) && m.holes > 16 {
			m.compact()
		}
		return prev, true
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			prev := m.entries[i].Value
			copy(m.entries[i:], m.entries[i+1:])
			m.entries[len(m.entries)-1] = MapEntry[K, V]{}
			m.entries = m.entries[:len(m.entries)-1]
			return prev, true
		}
	}
	return zero, false
}

// RemoveKey removes k and returns its previous value (or zero).
func (m *LinkedHashMap[K, V]) RemoveKey(k K) V {
	v, _ := m.Remove(k)
	return v
}

// Size returns the number of active entries in the map.
func (m *LinkedHashMap[K, V]) Size() int {
	if m == nil {
		return 0
	}
	if m.index != nil {
		return len(m.index)
	}
	return len(m.entries)
}

// IsEmpty reports whether the map contains no entries.
func (m *LinkedHashMap[K, V]) IsEmpty() bool {
	return m.Size() == 0
}

// Clear removes all entries from the map.
func (m *LinkedHashMap[K, V]) Clear() {
	if m == nil {
		return
	}
	clear(m.entries)
	m.entries = m.entries[:0]
	clear(m.index)
	m.index = nil
	m.holes = 0
}

// ComputeIfAbsent computes and inserts mapping(k) if k is not mapped, returning the value.
func (m *LinkedHashMap[K, V]) ComputeIfAbsent(k K, mapping func(K) V) V {
	if v, ok := m.Get(k); ok {
		return v
	}
	v := mapping(k)
	m.Put(k, v)
	return v
}

// PutAll copies all mappings from src to this map.
func (m *LinkedHashMap[K, V]) PutAll(src *LinkedHashMap[K, V]) {
	if src == nil {
		return
	}
	for k, v := range src.All() {
		m.Put(k, v)
	}
}

// KeySet returns a LinkedHashSet containing all keys in insertion order.
func (m *LinkedHashMap[K, V]) KeySet() *LinkedHashSet[K] {
	s := NewLinkedHashSetWithCapacity[K](m.Size())
	for k := range m.Keys() {
		s.Add(k)
	}
	return s
}

// Values returns an ArrayList containing all values in insertion order.
func (m *LinkedHashMap[K, V]) Values() *ArrayList[V] {
	list := NewArrayListWithCapacity[V](m.Size())
	if m == nil {
		return list
	}
	if m.holes == 0 {
		for i := range m.entries {
			list.Add(m.entries[i].Value)
		}
		return list
	}
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; ok && idx == i {
			list.Add(e.Value)
		}
	}
	return list
}

// Entries returns a copy of all entries in insertion order.
func (m *LinkedHashMap[K, V]) Entries() []MapEntry[K, V] {
	if m == nil {
		return nil
	}
	if m.holes == 0 {
		cp := make([]MapEntry[K, V], len(m.entries))
		copy(cp, m.entries)
		return cp
	}
	cp := make([]MapEntry[K, V], 0, m.Size())
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; ok && idx == i {
			cp = append(cp, e)
		}
	}
	return cp
}

// All returns an iter.Seq2[K, V] yielding (key, value) in deterministic insertion order.
func (m *LinkedHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m == nil {
			return
		}
		if m.holes == 0 {
			for i := range m.entries {
				if !yield(m.entries[i].Key, m.entries[i].Value) {
					return
				}
			}
			return
		}
		for i, e := range m.entries {
			if idx, ok := m.index[e.Key]; ok && idx == i {
				if !yield(e.Key, e.Value) {
					return
				}
			}
		}
	}
}

// Keys returns an iter.Seq[K] yielding keys in deterministic insertion order.
func (m *LinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if m == nil {
			return
		}
		if m.holes == 0 {
			for i := range m.entries {
				if !yield(m.entries[i].Key) {
					return
				}
			}
			return
		}
		for i, e := range m.entries {
			if idx, ok := m.index[e.Key]; ok && idx == i {
				if !yield(e.Key) {
					return
				}
			}
		}
	}
}

// Equals reports whether other is a map with identical key-value mappings.
func (m *LinkedHashMap[K, V]) Equals(other any) bool {
	o, ok := other.(*LinkedHashMap[K, V])
	if !ok {
		return false
	}
	if m == o {
		return true
	}
	if m.Size() != o.Size() {
		return false
	}
	for k, v := range m.All() {
		ov, ok := o.Get(k)
		if !ok || !javalang.ObjectsEquals(v, ov) {
			return false
		}
	}
	return true
}

// HashCode returns the Java map hash code: sum of (k.hashCode() ^ v.hashCode()).
func (m *LinkedHashMap[K, V]) HashCode() int32 {
	if m == nil {
		return 0
	}
	var h int32
	for k, v := range m.All() {
		h += javalang.ObjectHashCode(k) ^ javalang.ObjectHashCode(v)
	}
	return h
}

// String formats the map as "{k0=v0, k1=v1, ...}".
func (m *LinkedHashMap[K, V]) String() string {
	if m == nil {
		return "null"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, v := range m.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(javalang.ToJavaString(k))
		sb.WriteByte('=')
		sb.WriteString(javalang.ToJavaString(v))
	}
	sb.WriteByte('}')
	return sb.String()
}
