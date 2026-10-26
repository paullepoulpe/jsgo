package jutil

import (
	"iter"
	"strconv"
	"strings"
)

// Go-only: language-adapter — Guava ImmutableMultiset (and its Builder /
// toImmutableMultiset collector): distinct elements in first-occurrence
// order, each iterated count times in a row; size() is the total count.

// ImmutableMultiset is com.google.common.collect.ImmutableMultiset.
type ImmutableMultiset[T comparable] struct {
	counts *LinkedHashMap[T, int]
	size   int
}

var emptyStringImmutableMultiset ImmutableMultiset[string]

// ImmutableMultiset_Of is ImmutableMultiset.of(E...).
func ImmutableMultiset_Of[T comparable](elems ...T) *ImmutableMultiset[T] {
	if len(elems) == 0 {
		if empty, ok := any(&emptyStringImmutableMultiset).(*ImmutableMultiset[T]); ok {
			return empty
		}
		return &ImmutableMultiset[T]{}
	}
	var b ImmutableMultiset_Builder[T]
	for _, e := range elems {
		b.Add(e)
	}
	return b.Build()
}

// ImmutableMultiset_Builder is ImmutableMultiset.Builder.
type ImmutableMultiset_Builder[T comparable] struct {
	counts *LinkedHashMap[T, int]
	size   int
	built  bool
}

// NewImmutableMultiset_Builder is ImmutableMultiset.builder().
func NewImmutableMultiset_Builder[T comparable]() *ImmutableMultiset_Builder[T] {
	return &ImmutableMultiset_Builder[T]{}
}

// Add is ImmutableMultiset.Builder.add(E).
func (b *ImmutableMultiset_Builder[T]) Add(e T) *ImmutableMultiset_Builder[T] {
	if IsNil(any(e)) {
		panic(NewNullPointerException())
	}
	if b.counts == nil {
		b.counts = NewLinkedHashMap[T, int]()
	} else if b.built {
		b.counts = NewLinkedHashMapFrom(b.counts)
		b.built = false
	}
	b.counts.Put(e, b.counts.Get(e)+1)
	b.size++
	return b
}

// AddAll is ImmutableMultiset.Builder.addAll(Iterable).
func (b *ImmutableMultiset_Builder[T]) AddAll(es iter.Seq[T]) *ImmutableMultiset_Builder[T] {
	for e := range es {
		b.Add(e)
	}
	return b
}

// Contains reports whether e has been added to the builder.
func (b *ImmutableMultiset_Builder[T]) Contains(e T) bool {
	return b.counts != nil && b.counts.ContainsKey(e)
}

// Build is ImmutableMultiset.Builder.build().
func (b *ImmutableMultiset_Builder[T]) Build() *ImmutableMultiset[T] {
	if b.size == 0 {
		return ImmutableMultiset_Of[T]()
	}
	b.built = true
	return &ImmutableMultiset[T]{counts: b.counts, size: b.size}
}

// Size is Multiset.size(): the total number of occurrences.
func (m *ImmutableMultiset[T]) Size() int { return m.size }

// IsEmpty is Collection.isEmpty().
func (m *ImmutableMultiset[T]) IsEmpty() bool { return m.size == 0 }

// Count is Multiset.count(Object).
func (m *ImmutableMultiset[T]) Count(e T) int { return m.counts.Get(e) }

// Contains is Collection.contains(Object).
func (m *ImmutableMultiset[T]) Contains(e T) bool { return m.counts.ContainsKey(e) }

// ElementSet is Multiset.elementSet(): the distinct elements, in order.
func (m *ImmutableMultiset[T]) ElementSet() *LinkedHashSet[T] {
	s := NewLinkedHashSet[T]()
	for e := range m.counts.Keys() {
		s.Add(e)
	}
	return s
}

// All iterates the elements, each repeated count times.
func (m *ImmutableMultiset[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e, c := range m.counts.All() {
			for i := 0; i < c; i++ {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// EntrySet iterates (element, count) pairs in order.
func (m *ImmutableMultiset[T]) EntrySet() iter.Seq2[T, int] { return m.counts.All() }

// Equals is Multiset.equals(Object): same elements with the same counts.
func (m *ImmutableMultiset[T]) Equals(o any) bool {
	that, ok := o.(*ImmutableMultiset[T])
	if !ok || that == nil {
		return false
	}
	if m == that {
		return true
	}
	if m.size != that.size || m.counts.Size() != that.counts.Size() {
		return false
	}
	for e, c := range m.counts.All() {
		if that.counts.Get(e) != c {
			return false
		}
	}
	return true
}

// HashCode is Multiset.hashCode(): entrySet().hashCode().
func (m *ImmutableMultiset[T]) HashCode() int32 {
	var h int32
	for e, c := range m.counts.All() {
		h += ObjectHashCode(e) ^ int32(c)
	}
	return h
}

// String is Multiset.toString(): "[a x 2, b]".
func (m *ImmutableMultiset[T]) String() string {
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for e, c := range m.counts.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(ToJavaString(e))
		if c != 1 {
			sb.WriteString(" x ")
			sb.WriteString(strconv.Itoa(c))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}
