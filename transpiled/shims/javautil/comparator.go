// Clean-room specification based on java.util.Comparator & java.util.Iterator Javadoc.

package javautil

import (
	"cmp"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Comparator represents an order relation matching java.util.Comparator<T>.
// Returns < 0 if a < b, 0 if a == b, > 0 if a > b.
type Comparator[T any] func(a, b T) int

// NaturalOrder returns a comparator ordering elements in their natural ascending order.
func NaturalOrder[T cmp.Ordered]() Comparator[T] {
	return func(a, b T) int {
		return cmp.Compare(a, b)
	}
}

// ReverseOrder returns a comparator ordering elements in descending order.
func ReverseOrder[T cmp.Ordered]() Comparator[T] {
	return func(a, b T) int {
		return cmp.Compare(b, a)
	}
}

// Comparing constructs a comparator extracting key U from T and comparing with keyCmp.
func Comparing[T, U any](keyExtractor func(T) U, keyCmp Comparator[U]) Comparator[T] {
	return func(a, b T) int {
		return keyCmp(keyExtractor(a), keyExtractor(b))
	}
}

// NullsFirst returns a null-friendly comparator that considers nil less than non-nil.
func NullsFirst[T any](comparator Comparator[T]) Comparator[T] {
	return func(a, b T) int {
		aNil := javalang.IsNil(a)
		bNil := javalang.IsNil(b)
		if aNil && bNil {
			return 0
		}
		if aNil {
			return -1
		}
		if bNil {
			return 1
		}
		return comparator(a, b)
	}
}

// NullsLast returns a null-friendly comparator that considers nil greater than non-nil.
func NullsLast[T any](comparator Comparator[T]) Comparator[T] {
	return func(a, b T) int {
		aNil := javalang.IsNil(a)
		bNil := javalang.IsNil(b)
		if aNil && bNil {
			return 0
		}
		if aNil {
			return 1
		}
		if bNil {
			return -1
		}
		return comparator(a, b)
	}
}

// Iterator matches java.util.Iterator<T>.
type Iterator[T any] interface {
	HasNext() bool
	Next() T
	Remove()
}

// Iterable matches java.lang.Iterable<T>.
type Iterable[T any] interface {
	Iterator() Iterator[T]
}
