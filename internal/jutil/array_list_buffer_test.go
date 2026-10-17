package jutil

import (
	"slices"
	"testing"
)

// A list started with InitWithBuffer behaves as one made with
// NewArrayListWithCapacity: it fills the buffer in place, then grows away
// from it.
func TestArrayListInitWithBuffer(t *testing.T) {
	var holder struct {
		list ArrayList[int]
		buf  [2]int
	}
	l := &holder.list
	l.InitWithBuffer(holder.buf[:])
	if !l.IsEmpty() || l.Size() != 0 {
		t.Fatalf("new list not empty: %v", l.Slice())
	}
	l.Add(1)
	l.Add(2)
	if holder.buf != [2]int{1, 2} {
		t.Fatalf("first elements not in the buffer: %v", holder.buf)
	}
	l.Add(3)
	l.Add(4)
	if got := l.Slice(); !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Fatalf("after growing: %v", got)
	}
	l.Remove(1)
	if got := l.Slice(); !slices.Equal(got, []int{2, 3, 4}) {
		t.Fatalf("after Remove: %v", got)
	}
	if holder.buf != [2]int{1, 2} {
		t.Fatalf("buffer written after the list grew: %v", holder.buf)
	}
	// Holder, list and buffer are one allocation.
	type pointers struct {
		list ArrayList[*int]
		buf  [1]*int
	}
	var sink *pointers
	if n := testing.AllocsPerRun(100, func() {
		h := new(pointers)
		h.list.InitWithBuffer(h.buf[:])
		h.list.Add(nil)
		sink = h
	}); n != 1 {
		t.Fatalf("a holder with a filled buffer took %v allocations", n)
	}
	if sink.list.Size() != 1 {
		t.Fatal("size")
	}
}
