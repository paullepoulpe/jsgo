package jutil

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

var _ MapSource[string, int] = (*OrderedMap[string, int])(nil)

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

func NewOrderedMapFrom[K comparable, V any](src MapSource[K, V]) *OrderedMap[K, V] {
	if src == nil {
		return &OrderedMap[K, V]{}
	}
	if om, ok := src.(*OrderedMap[K, V]); ok {
		cloned := om.Clone()
		return &cloned
	}
	n := src.Size()
	if n <= 0 {
		return &OrderedMap[K, V]{}
	}
	out := NewOrderedMapWithCapacity[K, V](n)
	for k, v := range src.All() {
		out.Put(k, v)
	}
	return out
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

// PutNew inserts (k, v) when the caller has already verified that k is absent.
func (m *OrderedMap[K, V]) PutNew(k K, v V) {
	if m.index != nil {
		m.index[k] = len(m.entries)
		m.entries = append(m.entries, OrderedMapEntry[K, V]{Key: k, Value: v})
		return
	}
	m.appendUnique(k, v)
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

func (m *OrderedMap[K, V]) GetOrDefault(k K, def V) V {
	if v, ok := m.Lookup(k); ok {
		return v
	}
	return def
}

func (m *OrderedMap[K, V]) ContainsKey(k K) bool {
	_, ok := m.Lookup(k)
	return ok
}

// Clear removes all entries while retaining backing storage for reuse.
func (m *OrderedMap[K, V]) Clear() {
	if m == nil {
		return
	}
	clear(m.entries)
	m.entries = m.entries[:0]
	if m.index != nil {
		clear(m.index)
	}
	m.holes = 0
}

// Reset drops all entries and releases the backing slice and index.
func (m *OrderedMap[K, V]) Reset() {
	if m == nil {
		return
	}
	m.entries = nil
	m.index = nil
	m.holes = 0
}

// First returns the first key-value pair in insertion order, if any.
func (m *OrderedMap[K, V]) First() (K, V, bool) {
	var zeroK K
	var zeroV V
	if m == nil || len(m.entries) == m.holes {
		return zeroK, zeroV, false
	}
	if m.holes == 0 {
		e := m.entries[0]
		return e.Key, e.Value, true
	}
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; ok && idx == i {
			return e.Key, e.Value, true
		}
	}
	return zeroK, zeroV, false
}

// Index returns the 0-based insertion-order position of k and true if present,
// or (-1, false) if absent.
func (m *OrderedMap[K, V]) Index(k K) (int, bool) {
	if m == nil || len(m.entries) == m.holes {
		return -1, false
	}
	if m.holes == 0 {
		if m.index != nil {
			if idx, ok := m.index[k]; ok {
				return idx, true
			}
			return -1, false
		}
		for i := range m.entries {
			if m.entries[i].Key == k {
				return i, true
			}
		}
		return -1, false
	}
	rawIdx, ok := m.index[k]
	if !ok {
		return -1, false
	}
	pos := 0
	for i := 0; i < rawIdx; i++ {
		if idx, ok := m.index[m.entries[i].Key]; ok && idx == i {
			pos++
		}
	}
	return pos, true
}

// ValueAt returns the index-th value in insertion order without allocating a slice.
func (m *OrderedMap[K, V]) ValueAt(index int) (V, bool) {
	var zero V
	n := m.Size()
	if index < 0 || index >= n {
		return zero, false
	}
	if m.holes == 0 {
		return m.entries[index].Value, true
	}
	pos := 0
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; !ok || idx != i {
			continue
		}
		if pos == index {
			return e.Value, true
		}
		pos++
	}
	return zero, false
}

// Intern inserts k with the zero value of V if absent, and returns k's 0-based
// insertion-order index.
func (m *OrderedMap[K, V]) Intern(k K) int {
	m.compact()
	if m.index != nil {
		if idx, ok := m.index[k]; ok {
			return idx
		}
		idx := len(m.entries)
		m.index[k] = idx
		var zero V
		m.entries = append(m.entries, OrderedMapEntry[K, V]{Key: k, Value: zero})
		return idx
	}
	for i := range m.entries {
		if m.entries[i].Key == k {
			return i
		}
	}
	idx := len(m.entries)
	var zero V
	m.appendUnique(k, zero)
	return idx
}

func (m *OrderedMap[K, V]) IsNil() bool { return m == nil }

func (m *OrderedMap[K, V]) Size() int {
	if m == nil {
		return 0
	}
	return len(m.entries) - m.holes
}

func (m *OrderedMap[K, V]) Len() int { return m.Size() }

func (m *OrderedMap[K, V]) IsEmpty() bool { return m.Size() == 0 }

func (m *OrderedMap[K, V]) Entries() []OrderedMapEntry[K, V] {
	if m == nil || len(m.entries) == m.holes {
		return nil
	}
	if m.holes == 0 {
		return m.entries
	}
	out := make([]OrderedMapEntry[K, V], 0, len(m.entries)-m.holes)
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; !ok || idx != i {
			continue
		}
		out = append(out, e)
	}
	return out
}

func (m *OrderedMap[K, V]) KeySlice() []K {
	n := m.Size()
	if n == 0 {
		return nil
	}
	out := make([]K, 0, n)
	if m.holes == 0 {
		for _, e := range m.entries {
			out = append(out, e.Key)
		}
		return out
	}
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; !ok || idx != i {
			continue
		}
		out = append(out, e.Key)
	}
	return out
}

func (m *OrderedMap[K, V]) ValueSlice() []V {
	n := m.Size()
	if n == 0 {
		return nil
	}
	out := make([]V, 0, n)
	if m.holes == 0 {
		for _, e := range m.entries {
			out = append(out, e.Value)
		}
		return out
	}
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; !ok || idx != i {
			continue
		}
		out = append(out, e.Value)
	}
	return out
}

func (m *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m == nil {
			return
		}
		for i := 0; i < len(m.entries); i++ {
			e := m.entries[i]
			if m.holes > 0 {
				if idx, ok := m.index[e.Key]; !ok || idx != i {
					continue
				}
			}
			if !yield(e.Key, e.Value) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if m == nil {
			return
		}
		for i := 0; i < len(m.entries); i++ {
			e := m.entries[i]
			if m.holes > 0 {
				if idx, ok := m.index[e.Key]; !ok || idx != i {
					continue
				}
			}
			if !yield(e.Key) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if m == nil {
			return
		}
		for i := 0; i < len(m.entries); i++ {
			e := m.entries[i]
			if m.holes > 0 {
				if idx, ok := m.index[e.Key]; !ok || idx != i {
					continue
				}
			}
			if !yield(e.Value) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) Clone() OrderedMap[K, V] {
	n := m.Size()
	if n == 0 {
		return OrderedMap[K, V]{}
	}
	var cloned []OrderedMapEntry[K, V]
	if m.holes == 0 {
		cloned = slices.Clone(m.entries)
	} else {
		cloned = make([]OrderedMapEntry[K, V], 0, n)
		for i, e := range m.entries {
			if idx, ok := m.index[e.Key]; !ok || idx != i {
				continue
			}
			cloned = append(cloned, e)
		}
	}
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
	n := m.Size()
	out := make(map[K]V, n)
	if n == 0 {
		return out
	}
	if m.holes == 0 {
		for _, e := range m.entries {
			out[e.Key] = e.Value
		}
		return out
	}
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; !ok || idx != i {
			continue
		}
		out[e.Key] = e.Value
	}
	return out
}

func (m *OrderedMap[K, V]) EqualFunc(other *OrderedMap[K, V], eq func(V, V) bool) bool {
	if m.Size() != other.Size() {
		return false
	}
	if m == nil {
		return true
	}
	if m.holes == 0 {
		for _, e := range m.entries {
			ov, ok := other.Lookup(e.Key)
			if !ok || !eq(e.Value, ov) {
				return false
			}
		}
		return true
	}
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; !ok || idx != i {
			continue
		}
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

// String renders like AbstractMap.toString(): "{k1=v1, k2=v2}".
func (m *OrderedMap[K, V]) String() string { return formatMapString(m.All()) }

func (m *OrderedMap[K, V]) lookupAny(k any) (any, bool) {
	key, ok := k.(K)
	if !ok {
		return nil, false
	}
	return m.Lookup(key)
}

// Equals is Map.equals(Object): same size and equal values for every key.
func (m *OrderedMap[K, V]) Equals(o any) bool {
	if any(m) == o {
		return true
	}
	return mapEntriesEquals(m.Size(), m.All(), o)
}

// HashCode is Map.hashCode(): the sum of key.hashCode() ^ value.hashCode() across all entries.
func (m *OrderedMap[K, V]) HashCode() int32 {
	if m == nil || len(m.entries) == m.holes {
		return 0
	}
	var h int32
	if m.holes == 0 {
		for _, e := range m.entries {
			h += ObjectHashCode(any(e.Key)) ^ ObjectHashCode(any(e.Value))
		}
		return h
	}
	for i, e := range m.entries {
		if idx, ok := m.index[e.Key]; !ok || idx != i {
			continue
		}
		h += ObjectHashCode(any(e.Key)) ^ ObjectHashCode(any(e.Value))
	}
	return h
}

// OrderedSet is a compact, insertion-ordered set backed by OrderedMap[K, struct{}].
// Its zero value is ready to use and allocates zero heap memory when empty.
type OrderedSet[K comparable] struct {
	m OrderedMap[K, struct{}]
}

var _ Set[string] = (*OrderedSet[string])(nil)

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

func (s *OrderedSet[K]) mapOrNil() *OrderedMap[K, struct{}] {
	if s == nil {
		return nil
	}
	return &s.m
}

func (s *OrderedSet[K]) Add(k K) bool { return s.m.PutIfAbsent(k, struct{}{}) }

func (s *OrderedSet[K]) AddNew(k K) { s.m.PutNew(k, struct{}{}) }

func (s *OrderedSet[K]) AddAll(elems ...K) bool {
	changed := false
	for _, k := range elems {
		if s.m.PutIfAbsent(k, struct{}{}) {
			changed = true
		}
	}
	return changed
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

// Clear removes all elements while retaining backing storage for reuse.
func (s *OrderedSet[K]) Clear() {
	if s == nil {
		return
	}
	s.m.Clear()
}

// First returns the first element in insertion order, if any.
func (s *OrderedSet[K]) First() (K, bool) {
	k, _, ok := s.mapOrNil().First()
	return k, ok
}

// Index returns the 0-based insertion-order position of k and true if present,
// or (-1, false) if absent.
func (s *OrderedSet[K]) Index(k K) (int, bool) {
	return s.mapOrNil().Index(k)
}

// Intern inserts k if absent and returns its 0-based insertion-order index.
func (s *OrderedSet[K]) Intern(k K) int { return s.m.Intern(k) }

func (s *OrderedSet[K]) IsNil() bool { return s == nil }
func (s *OrderedSet[K]) Len() int    { return s.Size() }
func (s *OrderedSet[K]) Size() int   { return s.mapOrNil().Size() }
func (s *OrderedSet[K]) IsEmpty() bool {
	return s.Size() == 0
}

func (s *OrderedSet[K]) All() iter.Seq[K] {
	return s.mapOrNil().Keys()
}

func (s *OrderedSet[K]) ToSlice() []K {
	return s.mapOrNil().KeySlice()
}

func (s *OrderedSet[K]) Entries() []OrderedMapEntry[K, struct{}] {
	return s.mapOrNil().Entries()
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
	if s == nil {
		return true
	}
	if s.m.holes == 0 {
		for _, e := range s.m.entries {
			if !other.Has(e.Key) {
				return false
			}
		}
		return true
	}
	for i, e := range s.m.entries {
		if idx, ok := s.m.index[e.Key]; !ok || idx != i {
			continue
		}
		if !other.Has(e.Key) {
			return false
		}
	}
	return true
}

// String renders like AbstractCollection.toString(): "[a, b]".
func (s *OrderedSet[K]) String() string { return CollectionToString(s.ToSlice()) }

func (s *OrderedSet[K]) containsAny(e any) bool {
	x, ok := e.(K)
	return ok && s.Contains(x)
}

// Equals is Set.equals(Object): same size and every element contained.
func (s *OrderedSet[K]) Equals(o any) bool {
	if any(s) == o {
		return true
	}
	if other, ok := o.(*OrderedSet[K]); ok {
		if s == nil || other == nil {
			return false
		}
		return s.Equal(other)
	}
	return setEntriesEquals(s.Size(), s.All(), o)
}

// HashCode is Set.hashCode(): the sum of elem.hashCode() across all elements.
func (s *OrderedSet[K]) HashCode() int32 {
	if s == nil || len(s.m.entries) == s.m.holes {
		return 0
	}
	var h int32
	if s.m.holes == 0 {
		for _, e := range s.m.entries {
			h += ObjectHashCode(any(e.Key))
		}
		return h
	}
	for i, e := range s.m.entries {
		if idx, ok := s.m.index[e.Key]; !ok || idx != i {
			continue
		}
		h += ObjectHashCode(any(e.Key))
	}
	return h
}
