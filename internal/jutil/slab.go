package jutil

import "unsafe"

// Slab allocates values of T in chunks: New returns a pointer to a new zero
// T, and only one call in a few hundred allocates. It suits many small
// objects that are created together and become garbage together, such as
// the tokens and source positions of one parse: a pointer to any element
// keeps its whole chunk alive, so it must not be used for objects that may
// outlive their neighbours. The zero Slab is ready to use.
//
// Go-only: performance — the JVM allocates small objects by bumping a
// pointer; Go's malloc costs far more per object, and in a line-for-line
// port one allocation per token or position dominates the parser's profile.
type Slab[T any] struct {
	chunk      []T
	used       int
	fixedChunk int
	maxChunk   int
}

// NewFixedSlab returns a Slab that allocates chunks of fixedChunk elements.
func NewFixedSlab[T any](fixedChunk int) Slab[T] {
	return Slab[T]{fixedChunk: fixedChunk}
}

// NewMaxSlab returns a Slab that ramps exponentially from 16 up to maxChunk elements.
func NewMaxSlab[T any](maxChunk int) Slab[T] {
	return Slab[T]{maxChunk: maxChunk}
}

const (
	slabMinChunk       = 16
	slabMaxChunk       = 128
	maxSmallAllocBytes = 32760 // Go runtime maxSmallSize (32768) - gc.MallocHeaderSize (8)
)

// New returns a pointer to a new zero T, distinct from every other pointer
// it has returned.
func (s *Slab[T]) New() *T {
	if s.used == len(s.chunk) {
		s.grow()
	}
	p := &s.chunk[s.used]
	s.used++
	return p
}

// grow starts a new chunk, twice the size of the last one up to
// slabMaxChunk (or s.maxChunk) elements, clamped to maxSmallAllocBytes so
// chunks stay in P-local mcache small spans instead of mallocgcLarge.
func (s *Slab[T]) grow() {
	elemSize := unsafe.Sizeof(*new(T))
	maxSmallElems := 0
	if elemSize > 0 {
		maxSmallElems = max(int(maxSmallAllocBytes/elemSize), 1)
	}
	if s.fixedChunk > 0 {
		fc := s.fixedChunk
		if maxSmallElems > 0 && fc > maxSmallElems {
			fc = maxSmallElems
		}
		s.chunk = make([]T, fc)
		s.used = 0
		return
	}
	maxC := slabMaxChunk
	if s.maxChunk > 0 {
		maxC = s.maxChunk
	}
	if maxSmallElems > 0 && maxC > maxSmallElems {
		maxC = maxSmallElems
	}
	n := 2 * cap(s.chunk)
	if n < slabMinChunk {
		n = slabMinChunk
	}
	if n > maxC {
		n = maxC
	}
	s.chunk = make([]T, n)
	s.used = 0
}
