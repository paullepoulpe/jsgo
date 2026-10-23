package jutil

import (
	"slices"
	"testing"
)

func TestSlabNewReturnsDistinctZeroValues(t *testing.T) {
	type pos struct {
		src    *string
		offset int
	}
	var s Slab[pos]
	src := "f"
	seen := make(map[*pos]bool)
	var all []*pos
	for i := 0; i < 5000; i++ {
		p := s.New()
		if p.src != nil || p.offset != 0 {
			t.Fatalf("New() #%d is not zero: %+v", i, *p)
		}
		if seen[p] {
			t.Fatalf("New() #%d returned a pointer it had returned before", i)
		}
		seen[p] = true
		p.src, p.offset = &src, i
		all = append(all, p)
	}
	for i, p := range all {
		if p.src != &src || p.offset != i {
			t.Fatalf("value #%d was overwritten: %+v", i, *p)
		}
	}
}

// Chunks double from slabMinChunk to slabMaxChunk elements, each rounded up
// to what its size class holds: 64-byte pointerful elements fill 9472-byte
// chunks with 147 elements, not 128.
func TestSlabChunkSizes(t *testing.T) {
	type node struct {
		p [4]*int
		x [4]int
	}
	var s Slab[node]
	var got []int
	for i := 0; i < 3000; i++ {
		s.New()
		if s.used == 1 {
			got = append(got, len(s.chunk))
		}
	}
	want := []int{4, 8, 17, 35, 75, 147}
	for n := 4 + 8 + 17 + 35 + 75 + 147; n < 3000; n += 147 {
		want = append(want, 147)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("chunk sizes %v, want %v", got, want)
	}
	f := NewFixedSlab[node](100)
	f.New()
	if len(f.chunk) != 101 {
		t.Fatalf("fixed chunk of %d elements, want 101", len(f.chunk))
	}
}
