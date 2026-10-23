package jutil

import "slices"

// Slab allocates values of T in chunks: New returns a pointer to a new zero
// T, and only one call in a few hundred allocates. The zero Slab is ready to
// use.
//
// Use it only for objects that are created together and become garbage
// together, such as one pass's Vars or union-find nodes: a pointer to any
// element keeps its whole chunk alive (and everything the chunk's elements
// point to). Never use it for long-lived objects with individual lifetimes
// such as AST nodes or types: on typed-full, slab-allocated AST nodes kept
// dead neighbours (and through them the whole type registry) alive, and
// one allocation per node cut the peak live heap by 7 %.
//
// Go-only: performance — the JVM allocates small objects by bumping a
// pointer; Go's malloc costs far more per object.
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
// The chunk gets every element its allocation's size class holds: Go rounds
// an allocation up to a size class (after putting an 8-byte header in front
// of a pointerful one over 512 bytes), so make([]T, n) would leave up to 15%
// of the chunk's memory unused.
func (s *Slab[T]) grow() {
	n := s.fixedChunk
	if n == 0 {
		n = min(max(2*len(s.chunk), slabMinChunk), slabMaxChunk)
	}
	s.chunk = slices.Grow([]T(nil), n)
	s.chunk = s.chunk[:cap(s.chunk)]
	s.used = 0
}
