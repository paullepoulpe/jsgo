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

// Chunk byte sizes double from 512 (the smallest holding slabMinChunk
// 64-byte elements) to 8192, each filled to what its size class holds after
// the 8-byte header of a pointerful allocation.
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
	want := []int{7, 15, 31, 63, 127}
	for n := 7 + 15 + 31 + 63 + 127; n < 3000; n += 127 {
		want = append(want, 127)
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
