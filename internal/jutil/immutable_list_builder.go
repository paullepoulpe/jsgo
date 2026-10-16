package jutil

// Go-only: language-adapter — Guava ImmutableList.Builder, building a Go
// slice (ImmutableList is a slice in the port). Null elements are rejected
// with NullPointerException, as in Guava.

// ImmutableList_Builder is ImmutableList.Builder.
type ImmutableList_Builder[T any] struct {
	elems []T
}

// ImmutableList_Builder_ is ImmutableList.builder().
func ImmutableList_Builder_[T any]() *ImmutableList_Builder[T] {
	return &ImmutableList_Builder[T]{}
}

// Add is ImmutableList.Builder.add(E).
func (b *ImmutableList_Builder[T]) Add(e T) *ImmutableList_Builder[T] {
	CheckNotNull(e)
	b.elems = append(b.elems, e)
	return b
}

// AddAll is ImmutableList.Builder.addAll(Iterable).
func (b *ImmutableList_Builder[T]) AddAll(es []T) *ImmutableList_Builder[T] {
	for _, e := range es {
		b.Add(e)
	}
	return b
}

// Build is ImmutableList.Builder.build(): a copy of the elements added so
// far (never nil).
func (b *ImmutableList_Builder[T]) Build() []T {
	return append(make([]T, 0, len(b.elems)), b.elems...)
}
