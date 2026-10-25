package ast

import "iter"

// ReadAhead iterates over a list from its first element, as returned by a
// First<X> accessor:
//
//	for member := range ast.ReadAhead(obj.FirstMember()) {
//		member.Detach() // fine: the next member was already read
//	}
//
// It reads the next element before the loop body runs, so the body may
// detach, move or replace the current element. The list iterators
// (obj.Members(), block.Statements(), ...) read it after the body, like a
// for loop over GetNext; use them when the body changes the elements after
// the current one instead.
func ReadAhead[T interface{ Next() (T, bool) }](first T, ok bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		x, more := first, ok
		for more {
			next, nextOk := x.Next()
			if !yield(x) {
				return
			}
			x, more = next, nextOk
		}
	}
}
