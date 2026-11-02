package jscomp

import (
	"iter"
	"slices"
)

const orderedMapIndexThreshold = 8

// OrderedMapEntry pairs a key with a value in insertion order.
type OrderedMapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// OrderedMap is a compact, insertion-ordered map backed by a contiguous entry
// slice and a lazily populated hash index above orderedMapIndexThreshold keys.
// Its zero value is ready to use and allocates zero heap memory when empty.
type OrderedMap[K comparable, V any] struct {
	entries []OrderedMapEntry[K, V]
	index   map[K]int
	holes   int
}

type (
	OrderedStringEntry[V any] = OrderedMapEntry[string, V]
	OrderedStringMap[V any]   = OrderedMap[string, V]
)

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] { return &OrderedMap[K, V]{} }
func NewOrderedStringMap[V any]() *OrderedStringMap[V]      { return &OrderedStringMap[V]{} }

func NewOrderedMapWithCapacity[K comparable, V any](capHint int) *OrderedMap[K, V] {
	if capHint <= 0 {
		return &OrderedMap[K, V]{}
	}
	m := &OrderedMap[K, V]{entries: make([]OrderedMapEntry[K, V], 0, capHint)}
	if capHint > orderedMapIndexThreshold {
		m.index = make(map[K]int, capHint)
	}
	return m
}

func NewOrderedStringMapWithCapacity[V any](capHint int) *OrderedStringMap[V] {
	return NewOrderedMapWithCapacity[string, V](capHint)
}

func (m *OrderedMap[K, V]) compact() {
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

// Put inserts or updates k, preserving first-seen key order and returning the previous value.
func (m *OrderedMap[K, V]) Put(k K, v V) V {
	var zero V
	if m.index != nil {
		if idx, ok := m.index[k]; ok {
			prev := m.entries[idx].Value
			m.entries[idx].Value = v
			return prev
		}
		m.index[k] = len(m.entries)
		m.entries = append(m.entries, OrderedMapEntry[K, V]{Key: k, Value: v})
		return zero
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			prev := m.entries[i].Value
			m.entries[i].Value = v
			return prev
		}
	}
	m.appendUnique(k, v)
	return zero
}

// PutIfAbsent inserts (k, v) only if k is not already present, returning true if inserted.
func (m *OrderedMap[K, V]) PutIfAbsent(k K, v V) bool {
	if m.index != nil {
		if _, ok := m.index[k]; ok {
			return false
		}
		m.index[k] = len(m.entries)
		m.entries = append(m.entries, OrderedMapEntry[K, V]{Key: k, Value: v})
		return true
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			return false
		}
	}
	m.appendUnique(k, v)
	return true
}

func (m *OrderedMap[K, V]) appendUnique(k K, v V) {
	if len(m.entries) >= orderedMapIndexThreshold {
		m.index = make(map[K]int, len(m.entries)+4)
		for i, e := range m.entries {
			m.index[e.Key] = i
		}
		m.index[k] = len(m.entries)
	}
	m.entries = append(m.entries, OrderedMapEntry[K, V]{Key: k, Value: v})
}

func (m *OrderedMap[K, V]) ComputeIfAbsent(k K, mapping func(K) V) V {
	if v, ok := m.Lookup(k); ok {
		return v
	}
	v := mapping(k)
	m.Put(k, v)
	return v
}

func (m *OrderedMap[K, V]) Delete(k K) (V, bool) {
	var zero V
	if m == nil || len(m.entries) == m.holes {
		return zero, false
	}
	if m.index != nil {
		idx, ok := m.index[k]
		if !ok {
			return zero, false
		}
		prev := m.entries[idx].Value
		delete(m.index, k)
		m.entries[idx] = OrderedMapEntry[K, V]{}
		switch {
		case len(m.index) == 0:
			m.entries = m.entries[:0]
			m.holes = 0
		case idx == len(m.entries)-1:
			m.entries = m.entries[:idx]
		default:
			if m.holes++; m.holes*2 >= len(m.entries) {
				m.compact()
			}
		}
		return prev, true
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			prev := m.entries[i].Value
			m.entries = slices.Delete(m.entries, i, i+1)
			return prev, true
		}
	}
	return zero, false
}

func (m *OrderedMap[K, V]) Remove(k K) V {
	v, _ := m.Delete(k)
	return v
}

func (m *OrderedMap[K, V]) Lookup(k K) (V, bool) {
	var zero V
	if m == nil || len(m.entries) == m.holes {
		return zero, false
	}
	if m.index != nil {
		if idx, ok := m.index[k]; ok {
			return m.entries[idx].Value, true
		}
		return zero, false
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			return m.entries[i].Value, true
		}
	}
	return zero, false
}

func (m *OrderedMap[K, V]) Get(k K) V {
	v, _ := m.Lookup(k)
	return v
}

func (m *OrderedMap[K, V]) ContainsKey(k K) bool {
	_, ok := m.Lookup(k)
	return ok
}

func (m *OrderedMap[K, V]) Size() int {
	if m == nil {
		return 0
	}
	return len(m.entries) - m.holes
}

func (m *OrderedMap[K, V]) Len() int { return m.Size() }

func (m *OrderedMap[K, V]) IsEmpty() bool { return m.Size() == 0 }

func (m *OrderedMap[K, V]) Entries() []OrderedMapEntry[K, V] {
	if m == nil {
		return nil
	}
	m.compact()
	return m.entries
}

func (m *OrderedMap[K, V]) KeySlice() []K {
	entries := m.Entries()
	if len(entries) == 0 {
		return nil
	}
	out := make([]K, len(entries))
	for i, e := range entries {
		out[i] = e.Key
	}
	return out
}

func (m *OrderedMap[K, V]) ValueSlice() []V {
	entries := m.Entries()
	if len(entries) == 0 {
		return nil
	}
	out := make([]V, len(entries))
	for i, e := range entries {
		out[i] = e.Value
	}
	return out
}

func (m *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, e := range m.Entries() {
			if !yield(e.Key, e.Value) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, e := range m.Entries() {
			if !yield(e.Key) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, e := range m.Entries() {
			if !yield(e.Value) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) Clone() OrderedMap[K, V] {
	entries := m.Entries()
	if len(entries) == 0 {
		return OrderedMap[K, V]{}
	}
	cloned := slices.Clone(entries)
	var idx map[K]int
	if len(cloned) > orderedMapIndexThreshold {
		idx = make(map[K]int, len(cloned))
		for i, e := range cloned {
			idx[e.Key] = i
		}
	}
	return OrderedMap[K, V]{entries: cloned, index: idx}
}

func (m *OrderedMap[K, V]) ToMap() map[K]V {
	entries := m.Entries()
	out := make(map[K]V, len(entries))
	for _, e := range entries {
		out[e.Key] = e.Value
	}
	return out
}

func (m *OrderedMap[K, V]) EqualFunc(other *OrderedMap[K, V], eq func(V, V) bool) bool {
	if m.Size() != other.Size() {
		return false
	}
	for _, e := range m.Entries() {
		ov, ok := other.Lookup(e.Key)
		if !ok || !eq(e.Value, ov) {
			return false
		}
	}
	return true
}

func OrderedMapEqual[K comparable, V comparable](a, b *OrderedMap[K, V]) bool {
	return a.EqualFunc(b, func(x, y V) bool { return x == y })
}

// OrderedSet is a compact, insertion-ordered set backed by OrderedMap[K, struct{}].
// Its zero value is ready to use and allocates zero heap memory when empty.
type OrderedSet[K comparable] struct {
	m OrderedMap[K, struct{}]
}

func NewOrderedSet[K comparable](elems ...K) *OrderedSet[K] {
	s := NewOrderedSetWithCapacity[K](len(elems))
	s.AddAll(elems...)
	return s
}

func NewOrderedSetWithCapacity[K comparable](capHint int) *OrderedSet[K] {
	if capHint <= 0 {
		return &OrderedSet[K]{}
	}
	return &OrderedSet[K]{m: *NewOrderedMapWithCapacity[K, struct{}](capHint)}
}

func (s *OrderedSet[K]) Add(k K) bool { return s.m.PutIfAbsent(k, struct{}{}) }

func (s *OrderedSet[K]) AddAll(elems ...K) {
	for _, k := range elems {
		s.m.PutIfAbsent(k, struct{}{})
	}
}

func (s *OrderedSet[K]) Has(k K) bool      { return s != nil && s.m.ContainsKey(k) }
func (s *OrderedSet[K]) Contains(k K) bool { return s != nil && s.m.ContainsKey(k) }

func (s *OrderedSet[K]) Remove(k K) bool {
	if s == nil {
		return false
	}
	_, ok := s.m.Delete(k)
	return ok
}

func (s *OrderedSet[K]) Len() int { return s.Size() }
func (s *OrderedSet[K]) Size() int {
	if s == nil {
		return 0
	}
	return s.m.Size()
}
func (s *OrderedSet[K]) IsEmpty() bool { return s.Size() == 0 }

func (s *OrderedSet[K]) All() iter.Seq[K] {
	if s == nil {
		return func(func(K) bool) {}
	}
	return s.m.Keys()
}

func (s *OrderedSet[K]) ToSlice() []K {
	if s == nil {
		return nil
	}
	return s.m.KeySlice()
}

func (s *OrderedSet[K]) Clone() OrderedSet[K] {
	if s == nil || s.IsEmpty() {
		return OrderedSet[K]{}
	}
	return OrderedSet[K]{m: s.m.Clone()}
}

func (s *OrderedSet[K]) Equal(other *OrderedSet[K]) bool {
	if s.Size() != other.Size() {
		return false
	}
	for _, e := range s.m.Entries() {
		if !other.Has(e.Key) {
			return false
		}
	}
	return true
}
