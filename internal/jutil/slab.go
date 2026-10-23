package jutil

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
}

// NewFixedSlab returns a Slab that allocates chunks of fixedChunk elements.
func NewFixedSlab[T any](fixedChunk int) Slab[T] {
	return Slab[T]{fixedChunk: fixedChunk}
}

const (
	slabMinChunk = 16
	slabMaxChunk = 128
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
// slabMaxChunk elements, so that a slab that serves few objects stays small.
func (s *Slab[T]) grow() {
	if s.fixedChunk > 0 {
		s.chunk = make([]T, s.fixedChunk)
		s.used = 0
		return
	}
	n := 2 * len(s.chunk)
	if n < slabMinChunk {
		n = slabMinChunk
	} else if n > slabMaxChunk {
		n = slabMaxChunk
	}
	s.chunk = make([]T, n)
	s.used = 0
}
