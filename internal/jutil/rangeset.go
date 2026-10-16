package jutil

import "sort"

// Go-only: language-adapter — Guava's Range and TreeRangeSet, for closed
// ranges (Range.closed), which is the only kind the port creates. A Java
// Comparable element type is given as a comparator.

// Range is a Guava Range<C> created by Range.closed(lower, upper): every
// value v with lower <= v <= upper.
type Range[T any] struct {
	lower, upper T
}

// Range_Closed is Range.closed(lower, upper). It panics with an
// IllegalArgumentException if lower > upper, as Guava does.
func Range_Closed[T any](lower, upper T, cmp func(a, b T) int) Range[T] {
	if cmp(lower, upper) > 0 {
		panic(NewIllegalArgumentExceptionMsg(
			"Invalid range: [" + ToJavaString(lower) + ".." + ToJavaString(upper) + "]"))
	}
	return Range[T]{lower: lower, upper: upper}
}

// LowerEndpoint is Range.lowerEndpoint().
func (r Range[T]) LowerEndpoint() T { return r.lower }

// UpperEndpoint is Range.upperEndpoint().
func (r Range[T]) UpperEndpoint() T { return r.upper }

// TreeRangeSet is Guava's TreeRangeSet<C>: a set of disjoint ranges, kept
// sorted, where adding a range coalesces it with every range it is connected
// to.
type TreeRangeSet[T any] struct {
	cmp    func(a, b T) int
	ranges []Range[T] // sorted by lower endpoint, pairwise disconnected
}

// NewTreeRangeSet is TreeRangeSet.create() for elements ordered by cmp.
func NewTreeRangeSet[T any](cmp func(a, b T) int) *TreeRangeSet[T] {
	return &TreeRangeSet[T]{cmp: cmp}
}

// IsEmpty is RangeSet.isEmpty().
func (s *TreeRangeSet[T]) IsEmpty() bool { return len(s.ranges) == 0 }

// Add is RangeSet.add(Range): the set afterwards contains every value of r,
// and r is merged with the ranges it overlaps or touches.
func (s *TreeRangeSet[T]) Add(r Range[T]) {
	// First range whose upper endpoint is >= r.lower: every range before it is
	// strictly below r and stays as it is.
	i := sort.Search(len(s.ranges), func(k int) bool { return s.cmp(s.ranges[k].upper, r.lower) >= 0 })
	j := i
	lower, upper := r.lower, r.upper
	for j < len(s.ranges) && s.cmp(s.ranges[j].lower, r.upper) <= 0 {
		if s.cmp(s.ranges[j].lower, lower) < 0 {
			lower = s.ranges[j].lower
		}
		if s.cmp(s.ranges[j].upper, upper) > 0 {
			upper = s.ranges[j].upper
		}
		j++
	}
	merged := Range[T]{lower: lower, upper: upper}
	if i == j {
		s.ranges = append(s.ranges, Range[T]{})
		copy(s.ranges[i+1:], s.ranges[i:])
		s.ranges[i] = merged
		return
	}
	s.ranges[i] = merged
	s.ranges = append(s.ranges[:i+1], s.ranges[j:]...)
}

// Contains is RangeSet.contains(C): whether some range of the set contains v.
func (s *TreeRangeSet[T]) Contains(v T) bool {
	i := sort.Search(len(s.ranges), func(k int) bool { return s.cmp(s.ranges[k].upper, v) >= 0 })
	return i < len(s.ranges) && s.cmp(s.ranges[i].lower, v) <= 0
}

// AsRanges is RangeSet.asRanges(): the disconnected ranges, in ascending
// order.
func (s *TreeRangeSet[T]) AsRanges() []Range[T] {
	return append([]Range[T](nil), s.ranges...)
}
