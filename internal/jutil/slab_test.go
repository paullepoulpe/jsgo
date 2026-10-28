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

// Chunk requested sizes ramp 16 -> 32 -> 64 -> 128, each filled to what its
// allocation size class holds (and capped by maxChunk when set).
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
	c16 := len(SlabChunk[node](16))
	c32 := len(SlabChunk[node](32))
	c64 := len(SlabChunk[node](64))
	c128 := len(SlabChunk[node](128))
	want := []int{c16, c32, c64, c128}
	for n := c16 + c32 + c64 + c128; n < 3000; n += c128 {
		want = append(want, c128)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("chunk sizes %v, want %v", got, want)
	}
	m := NewMaxSlab[node](40)
	var maxGot []int
	for i := 0; i < 150; i++ {
		m.New()
		if m.used == 1 {
			maxGot = append(maxGot, len(m.chunk))
		}
	}
	for _, sz := range maxGot {
		if sz > 40 {
			t.Fatalf("NewMaxSlab(40) produced chunk of %d elements (> 40)", sz)
		}
	}
	f := NewFixedSlab[node](100)
	f.New()
	if len(f.chunk) != len(SlabChunk[node](100)) {
		t.Fatalf("fixed chunk of %d elements, want %d", len(f.chunk), len(SlabChunk[node](100)))
	}
}
