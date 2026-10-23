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
	slabMinChunk = 4
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
// slabMaxChunk (or s.maxChunk) elements, so that a slab that serves few
// objects stays small (many slabs serve only a handful: a scope creator's
// scopes for one function, say). The chunk gets every element its allocation's size
// class holds: Go rounds an allocation up to a size class (after putting an
// 8-byte header in front of a pointerful one over 512 bytes), so make([]T, n)
// would leave up to 15% of the chunk's memory unused.
func (s *Slab[T]) grow() {
	n := s.fixedChunk
	if n == 0 {
		maxC := slabMaxChunk
		if s.maxChunk > 0 {
			maxC = s.maxChunk
		}
		n = min(max(2*len(s.chunk), slabMinChunk), maxC)
	}
	s.chunk = SlabChunk[T](n)
	s.used = 0
}

// SlabChunk returns a zeroed []T of at least n elements, as many as the
// allocation's size class holds (see Slab.grow), for callers that carve
// small slices out of shared chunks.
func SlabChunk[T any](n int) []T {
	c := slices.Grow([]T(nil), n)
	return c[:cap(c)]
}
