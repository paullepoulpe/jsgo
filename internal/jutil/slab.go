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
	nextChunk  int
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
	minChunk        = 16
	defaultMaxChunk = 128
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

// grow starts a new chunk, ramping 16 -> 32 -> 64 -> 128 elements (capped by
// maxChunk when set) so that a slab that serves few objects stays small while
// larger slabs amortize allocations without spanning many 8 KB pages.
func (s *Slab[T]) grow() {
	n := s.fixedChunk
	if n == 0 {
		limit := defaultMaxChunk
		if s.maxChunk > 0 {
			limit = s.maxChunk
		}
		if s.nextChunk == 0 {
			n = min(minChunk, limit)
		} else {
			n = min(s.nextChunk, limit)
		}
		s.nextChunk = min(2*n, limit)
	}
	chunk := SlabChunk[T](n)
	if s.maxChunk > 0 && len(chunk) > s.maxChunk {
		chunk = chunk[:s.maxChunk]
	}
	s.chunk = chunk
	s.used = 0
}

// SlabChunk returns a zeroed []T of at least n elements, as many as the
// allocation's size class holds, for callers that carve small slices out of
// shared chunks.
func SlabChunk[T any](n int) []T {
	c := slices.Grow([]T(nil), n)
	return c[:cap(c)]
}

// ReuseSlab hands out pointers to Ts that stay valid until Reset: zeroed
// ones, or ones that Reset recycled. Chunks double in size from
// reuseSlabMinChunk elements up to reuseSlabMaxChunk, each rounded up to its
// allocation size class.
//
// Go-only: performance — objects that live only as long as one analysis
// (graph nodes and edges, flow scopes) are allocated here, one allocation
// per chunk instead of one per object; when the analysis is redone for the
// next function (see ReusableCfg), its slabs are rewound instead of
// allocated again.
type ReuseSlab[T any] struct {
	chunks [][]T // chunks[:cur] are full; chunks[cur][:used] is in use
	cur    int
	used   int
}

const (
	reuseSlabMinChunk = 8
	reuseSlabMaxChunk = 256
	// A reset slab keeps at most this many chunks (about 2,500 elements).
	reuseSlabKeepChunks = 15
)

func (s *ReuseSlab[T]) Alloc() *T {
	if s.cur < len(s.chunks) && s.used == len(s.chunks[s.cur]) {
		s.cur++
		s.used = 0
	}
	if s.cur == len(s.chunks) {
		n := reuseSlabMinChunk
		if s.cur > 0 {
			n = min(2*len(s.chunks[s.cur-1]), reuseSlabMaxChunk)
		}
		chunk := slices.Grow([]T(nil), n)
		s.chunks = append(s.chunks, chunk[:cap(chunk)])
	}
	p := &s.chunks[s.cur][s.used]
	s.used++
	return p
}

// Reset makes s hand out its elements again from the start, after zeroing
// those that Alloc handed out (so that they do not keep the old analysis
// alive), or after calling recycle on them when it is non-nil. It drops all
// but the first reuseSlabKeepChunks chunks.
func (s *ReuseSlab[T]) Reset(recycle func(*T)) {
	for i := 0; i <= s.cur && i < len(s.chunks); i++ {
		used := s.chunks[i]
		if i == s.cur {
			used = used[:s.used]
		}
		if recycle == nil {
			clear(used)
			continue
		}
		for j := range used {
			recycle(&used[j])
		}
	}
	s.Rewind()
}

// Rewind is Reset for Ts whose users re-initialize every field of what Alloc
// returns: it hands the elements out again as they are, so that the storage
// they own can be reused too.
func (s *ReuseSlab[T]) Rewind() {
	if len(s.chunks) > reuseSlabKeepChunks {
		clear(s.chunks[reuseSlabKeepChunks:])
		s.chunks = s.chunks[:reuseSlabKeepChunks]
	}
	s.cur = 0
	s.used = 0
}
