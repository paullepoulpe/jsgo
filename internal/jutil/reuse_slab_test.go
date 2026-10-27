package jutil

import "testing"

func TestReuseSlabUniquePointersAndReset(t *testing.T) {
	var s ReuseSlab[int]
	seen := make(map[*int]int, 3500)
	for i := range 3500 {
		p := s.Alloc()
		if *p != 0 {
			t.Fatalf("alloc %d not zeroed: %d", i, *p)
		}
		if prev, dup := seen[p]; dup {
			t.Fatalf("alloc %d reused pointer from %d before Reset", i, prev)
		}
		*p = i + 1
		seen[p] = i
	}
	s.Reset(nil)
	seenAfter := make(map[*int]int, 1000)
	for i := range 1000 {
		p := s.Alloc()
		if *p != 0 {
			t.Fatalf("post-Reset alloc %d not zeroed: %d", i, *p)
		}
		if prev, dup := seenAfter[p]; dup {
			t.Fatalf("post-Reset alloc %d reused pointer from %d", i, prev)
		}
		*p = i + 1
		seenAfter[p] = i
	}
}
