package jutil

import "testing"

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
