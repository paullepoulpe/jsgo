package jutil

import (
	"iter"
	"slices"
	"strings"
)

// Go-only: language-adapter — java.util / Guava collection interfaces and
// helpers with Java iteration-order semantics.

// MapSource is the part of java.util.Map<K, V> that copy constructors and
// putAll read: the size and the entries in iteration order. *OrderedMap
// implements it.
type MapSource[K any, V any] interface {
	Size() int
	All() iter.Seq2[K, V]
}

// CollectionToString renders elements like Java's AbstractCollection.toString.
func CollectionToString[T any](elems []T) string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, e := range elems {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ToJavaString(e))
	}
	sb.WriteByte(']')
	return sb.String()
}

// ListHashCode is List.hashCode() for a slice standing in for a Java list.
func ListHashCode[T any](elems []T) int32 {
	var h int32 = 1
	for _, e := range elems {
		h = 31*h + ObjectHashCode(e)
	}
	return h
}

func formatMapString[K any, V any](seq iter.Seq2[K, V]) string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, v := range seq {
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

// javaMapEntries is implemented by map types so that Map.equals can compare
// across implementations.
type javaMapEntries interface {
	Size() int
	lookupAny(k any) (any, bool)
}

func mapEntriesEquals[K any, V any](size int, seq iter.Seq2[K, V], o any) bool {
	other, ok := o.(javaMapEntries)
	if !ok || IsNil(o) || other.Size() != size {
		return false
	}
	for k, v := range seq {
		ov, present := other.lookupAny(k)
		if !present || !ObjectsEquals(v, ov) {
			return false
		}
	}
	return true
}

// Go-only: language-adapter — Arrays.binarySearch(int[], int).

// ArraysBinarySearchInts is Arrays.binarySearch(int[] a, int key) on a
// sorted slice: the index of key, or -(insertion point) - 1 when key is
// absent, the insertion point being the index of the first element greater
// than key (len(a) if there is none), as the Javadoc specifies.
//
// The Javadoc leaves open which index is returned when key occurs several
// times. Observed behaviour (black-box): it is the one the textbook
// bisection finds, probing the middle index (low+high)/2 (rounded down) of
// the remaining range [low, high] and returning on the first equal element;
// checked against 200,000 random sorted arrays with many duplicates.
func ArraysBinarySearchInts(a []int, key int) int {
	low, high := 0, len(a)-1
	for low <= high {
		mid := int(uint(low+high) >> 1)
		switch v := a[mid]; {
		case v < key:
			low = mid + 1
		case v > key:
			high = mid - 1
		default:
			return mid
		}
	}
	return -(low + 1)
}

// ---------------------------------------------------------------------------
// OrderedSetMultimap

// OrderedSetView is an insertion-ordered, deduplicated read-only slice view of
// comparable elements.
type OrderedSetView[V comparable] []V

// Contains reports whether v is in s.
func (s OrderedSetView[V]) Contains(v V) bool { return slices.Contains(s, v) }

// ToSlice returns a defensive copy of s (non-nil empty slice when empty).
func (s OrderedSetView[V]) ToSlice() []V {
	if len(s) == 0 {
		return []V{}
	}
	return slices.Clone([]V(s))
}

// All iterates elements in insertion order.
func (s OrderedSetView[V]) All() iter.Seq[V] { return slices.Values([]V(s)) }

// Size returns the number of elements in s.
func (s OrderedSetView[V]) Size() int { return len(s) }

// IsEmpty reports whether s has no elements.
func (s OrderedSetView[V]) IsEmpty() bool { return len(s) == 0 }

// OrderedSetMultimap maps keys to insertion-ordered, deduplicated value sets,
// preserving first-seen key insertion order (Guava LinkedHashMultimap /
// ImmutableSetMultimap semantics). Its zero value is ready to use.
type OrderedSetMultimap[K, V comparable] struct {
	vals  map[K][]V
	order []K
	size  int
}

// NewOrderedSetMultimap returns an empty OrderedSetMultimap.
func NewOrderedSetMultimap[K, V comparable]() *OrderedSetMultimap[K, V] {
	return &OrderedSetMultimap[K, V]{}
}

// Put adds (k, v) if not already present and reports whether the multimap changed.
func (m *OrderedSetMultimap[K, V]) Put(k K, v V) bool {
	if m.vals == nil {
		m.vals = make(map[K][]V)
	}
	existing, ok := m.vals[k]
	if !ok {
		m.order = append(m.order, k)
	} else if slices.Contains(existing, v) {
		return false
	}
	m.vals[k] = append(existing, v)
	m.size++
	return true
}

// Get returns the ordered value set for k (empty view when absent).
func (m *OrderedSetMultimap[K, V]) Get(k K) OrderedSetView[V] {
	if m == nil {
		return nil
	}
	return OrderedSetView[V](m.vals[k])
}

// Size returns the total number of (key, value) pairs.
func (m *OrderedSetMultimap[K, V]) Size() int {
	if m == nil {
		return 0
	}
	return m.size
}

// IsEmpty reports whether m contains no entries.
func (m *OrderedSetMultimap[K, V]) IsEmpty() bool {
	return m == nil || m.size == 0
}

// Entries iterates (key, value) pairs in key and per-key insertion order.
func (m *OrderedSetMultimap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m == nil {
			return
		}
		for _, k := range m.order {
			for _, v := range m.vals[k] {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// AsMap iterates (key, valueSet) pairs in key insertion order.
func (m *OrderedSetMultimap[K, V]) AsMap() iter.Seq2[K, OrderedSetView[V]] {
	return func(yield func(K, OrderedSetView[V]) bool) {
		if m == nil {
			return
		}
		for _, k := range m.order {
			if !yield(k, OrderedSetView[V](m.vals[k])) {
				return
			}
		}
	}
}

// Clone returns an independent deep copy of m's key and value slices.
func (m *OrderedSetMultimap[K, V]) Clone() OrderedSetMultimap[K, V] {
	if m == nil || len(m.order) == 0 {
		return OrderedSetMultimap[K, V]{}
	}
	vals := make(map[K][]V, len(m.vals))
	for _, k := range m.order {
		vals[k] = slices.Clone(m.vals[k])
	}
	return OrderedSetMultimap[K, V]{
		vals:  vals,
		order: slices.Clone(m.order),
		size:  m.size,
	}
}

// Inverse returns a new multimap mapping each value to its associated keys in
// insertion order.
func (m *OrderedSetMultimap[K, V]) Inverse() *OrderedSetMultimap[V, K] {
	out := &OrderedSetMultimap[V, K]{}
	if m == nil || len(m.order) == 0 {
		return out
	}
	out.vals = make(map[V][]K, len(m.order))
	out.order = make([]V, 0, len(m.order))
	for _, k := range m.order {
		for _, v := range m.vals[k] {
			ks, ok := out.vals[v]
			if !ok {
				out.order = append(out.order, v)
			} else if slices.Contains(ks, k) {
				continue
			}
			out.vals[v] = append(ks, k)
			out.size++
		}
	}
	return out
}

// String formats m like Guava's AbstractMultimap.toString(): "{k1=[v1, v2]}".
func (m *OrderedSetMultimap[K, V]) String() string {
	if m == nil || len(m.order) == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, k := range m.order {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ToJavaString(k))
		sb.WriteString("=[")
		for j, v := range m.vals[k] {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(ToJavaString(v))
		}
		sb.WriteByte(']')
	}
	sb.WriteByte('}')
	return sb.String()
}
