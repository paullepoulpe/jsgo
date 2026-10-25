package jutil

import (
	"reflect"
	"slices"
)

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
	chunkBytes int // the last chunk's size class (see grow)
	fixedChunk int
	maxChunk   int
}

// NewFixedSlab returns a Slab that allocates chunks of fixedChunk elements.
func NewFixedSlab[T any](fixedChunk int) Slab[T] {
	return Slab[T]{fixedChunk: fixedChunk}
}

// NewMaxSlab returns a Slab whose chunks ramp up as usual but never hold
// more than maxChunk elements.
func NewMaxSlab[T any](maxChunk int) Slab[T] {
	return Slab[T]{maxChunk: maxChunk}
}

const (
	slabMinChunk      = 4
	slabMinChunkBytes = 256
	slabMaxChunkBytes = 8192
	slabHeaderBytes   = 8
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

// grow starts a new chunk, twice the size of the last one, so that a slab
// that serves few objects stays small (many serve only a handful: one
// traversal's scopes, say). Chunk sizes are the size classes 256, 512, ...,
// 8192 bytes (counting the 8-byte header Go puts in front of a pointerful
// allocation over 512 bytes), and the chunk gets every element its class
// holds. Those classes are the ones whose spans are a single 8 KB page:
// a chunk then never waits for a run of contiguous free pages, and any
// page the GC frees can serve one. Chunks of 10-60 KB (5-7-page spans)
// made the heap grow 14% (700 MB) on a 35 M-node input although the live
// heap was unchanged.
func (s *Slab[T]) grow() {
	n := s.fixedChunk
	if n == 0 {
		size := max(int(reflect.TypeFor[T]().Size()), 1)
		b := 2 * s.chunkBytes
		if b == 0 {
			b = slabMinChunkBytes
			for b < slabMinChunk*size+slabHeaderBytes {
				b *= 2
			}
		}
		b = min(b, slabMaxChunkBytes)
		s.chunkBytes = b
		n = max((b-slabHeaderBytes)/size, 1)
		if s.maxChunk > 0 {
			n = min(n, s.maxChunk)
		}
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
