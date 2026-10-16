package jutil

import "iter"

// Go-only: language-adapter — the read side of java.util.Set, for methods
// that return a set view (Collections.unmodifiableSet(map.keySet()), an
// AbstractSet over live state) rather than a copy. *LinkedHashSet and
// *ObjectLinkedHashSet implement it.

// Set is the read-only part of java.util.Set.
type Set[T any] interface {
	Size() int
	IsEmpty() bool
	Contains(e T) bool
	// All iterates in the set's iteration order.
	All() iter.Seq[T]
}

var (
	_ Set[string] = (*LinkedHashSet[string])(nil)
	_ Set[string] = (*ObjectLinkedHashSet[string])(nil)
)
