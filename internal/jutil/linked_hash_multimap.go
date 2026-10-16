package jutil

import (
	"iter"
	"strings"
)

// Go-only: language-adapter — Guava LinkedHashMultimap
// (com.google.common.collect.LinkedHashMultimap). Keys and values are
// compared with Java equals()/hashCode() (ObjectsEquals / ObjectHashCode); a
// key-value pair is stored at most once. Get and RemoveAll return the values
// of a key in the order they were added; Entries and Values iterate across
// all key-value pairs in the order they were added to the multimap (not
// grouped by key); KeySet iterates the keys in the order they were first
// added (a key whose values were all removed and that is added again moves
// to the end).

type linkedHashMultimapEntry[K any, V any] struct {
	key        K
	value      V
	prev, next *linkedHashMultimapEntry[K, V]
	removed    bool
}

// LinkedHashMultimap is com.google.common.collect.LinkedHashMultimap.
type LinkedHashMultimap[K any, V any] struct {
	m          *ObjectLinkedHashMap[K, *ObjectLinkedHashMap[V, *linkedHashMultimapEntry[K, V]]]
	head, tail *linkedHashMultimapEntry[K, V]
	size       int
}

// NewLinkedHashMultimap is LinkedHashMultimap.create().
func NewLinkedHashMultimap[K any, V any]() *LinkedHashMultimap[K, V] {
	return &LinkedHashMultimap[K, V]{
		m: NewObjectLinkedHashMap[K, *ObjectLinkedHashMap[V, *linkedHashMultimapEntry[K, V]]](),
	}
}

// Put is Multimap.put; reports whether the multimap changed (false when the
// key-value pair was already present).
func (mm *LinkedHashMultimap[K, V]) Put(k K, v V) bool {
	vs := mm.m.Get(k)
	if vs == nil {
		vs = NewObjectLinkedHashMap[V, *linkedHashMultimapEntry[K, V]]()
		mm.m.Put(k, vs)
	} else if vs.ContainsKey(v) {
		return false
	}
	e := &linkedHashMultimapEntry[K, V]{key: k, value: v, prev: mm.tail}
	if mm.tail != nil {
		mm.tail.next = e
	} else {
		mm.head = e
	}
	mm.tail = e
	vs.Put(v, e)
	mm.size++
	return true
}

func (mm *LinkedHashMultimap[K, V]) unlink(e *linkedHashMultimapEntry[K, V]) {
	e.removed = true
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		mm.head = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		mm.tail = e.prev
	}
	mm.size--
}

// Get is Multimap.get: a snapshot of the values for k, in the order they
// were added (empty when absent).
func (mm *LinkedHashMultimap[K, V]) Get(k K) []V {
	vs := mm.m.Get(k)
	if vs == nil {
		return nil
	}
	return vs.KeySlice()
}

// ContainsKey is Multimap.containsKey.
func (mm *LinkedHashMultimap[K, V]) ContainsKey(k K) bool { return mm.m.ContainsKey(k) }

// ContainsEntry is Multimap.containsEntry.
func (mm *LinkedHashMultimap[K, V]) ContainsEntry(k K, v V) bool {
	vs := mm.m.Get(k)
	return vs != nil && vs.ContainsKey(v)
}

// Remove is Multimap.remove(key, value); reports whether the multimap
// changed.
func (mm *LinkedHashMultimap[K, V]) Remove(k K, v V) bool {
	vs := mm.m.Get(k)
	if vs == nil {
		return false
	}
	e, ok := vs.Lookup(v)
	if !ok {
		return false
	}
	vs.Remove(v)
	mm.unlink(e)
	if vs.IsEmpty() {
		mm.m.Remove(k)
	}
	return true
}

// RemoveAll is Multimap.removeAll: removes every value of k and returns
// them in the order they were added (empty when absent).
func (mm *LinkedHashMultimap[K, V]) RemoveAll(k K) []V {
	vs := mm.m.Get(k)
	if vs == nil {
		return nil
	}
	mm.m.Remove(k)
	out := make([]V, 0, vs.Size())
	for v, e := range vs.All() {
		out = append(out, v)
		mm.unlink(e)
	}
	return out
}

// Size is Multimap.size(): the number of key-value pairs.
func (mm *LinkedHashMultimap[K, V]) Size() int { return mm.size }

// IsEmpty is Multimap.isEmpty().
func (mm *LinkedHashMultimap[K, V]) IsEmpty() bool { return mm.size == 0 }

// KeySet iterates the distinct keys (Multimap.keySet()).
func (mm *LinkedHashMultimap[K, V]) KeySet() iter.Seq[K] { return mm.m.Keys() }

// Entries iterates the key-value pairs in the order they were added
// (Multimap.entries()).
func (mm *LinkedHashMultimap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := mm.head; e != nil; e = e.next {
			if e.removed {
				continue
			}
			if !yield(e.key, e.value) {
				return
			}
		}
	}
}

// Values iterates the values of all key-value pairs in the order they were
// added (Multimap.values()).
func (mm *LinkedHashMultimap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range mm.Entries() {
			if !yield(v) {
				return
			}
		}
	}
}

// String renders like Multimap.toString(): "{k1=[v1, v2], k2=[v3]}".
func (mm *LinkedHashMultimap[K, V]) String() string {
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
		sb.WriteString(CollectionToString(vs.KeySlice()))
	}
	sb.WriteByte('}')
	return sb.String()
}
