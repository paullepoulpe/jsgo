package jutil

import "slices"

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
