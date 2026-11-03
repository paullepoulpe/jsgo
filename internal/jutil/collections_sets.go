package jutil

import "iter"

// Go-only: language-adapter — the read side of java.util.Set.

// Set is the read-only part of java.util.Set.
type Set[T any] interface {
	Size() int
	IsEmpty() bool
	Contains(e T) bool
	// All iterates in the set's iteration order.
	All() iter.Seq[T]
}

var _ Set[string] = (*OrderedSet[string])(nil)

// SetHashCode is Set.hashCode() for a slice standing in for a Java set: the
// sum of the element hash codes.
func SetHashCode[T any](elems []T) int32 {
	var h int32
	for _, e := range elems {
		h += ObjectHashCode(any(e))
	}
	return h
}

// javaSet is implemented by set types so that Set.equals can compare across
// implementations.
type javaSet interface {
	Size() int
	containsAny(e any) bool
}

func setEntriesEquals[T any](size int, seq iter.Seq[T], o any) bool {
	other, ok := o.(javaSet)
	if !ok || IsNil(o) || other.Size() != size {
		return false
	}
	for e := range seq {
		if !other.containsAny(e) {
			return false
		}
	}
	return true
}
