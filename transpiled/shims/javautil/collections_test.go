package javautil

import (
	"testing"
)

func TestLinkedHashSetOrderAndAlgebra(t *testing.T) {
	s1 := NewLinkedHashSet("apple", "banana", "cherry")
	if s1.Size() != 3 {
		t.Errorf("s1.Size() got %d, want 3", s1.Size())
	}

	// Verify insertion order
	slice := s1.ToSlice()
	if slice[0] != "apple" || slice[1] != "banana" || slice[2] != "cherry" {
		t.Errorf("unexpected set order: %v", slice)
	}

	// Deduplication
	s1.Add("apple")
	if s1.Size() != 3 {
		t.Errorf("adding duplicate should not increase size")
	}

	s2 := NewLinkedHashSet("banana", "date", "elderberry")
	// Union: copy s1 and AddAll
	union := NewLinkedHashSet(s1.ToSlice()...)
	union.AddAllSet(s2)
	if union.Size() != 5 {
		t.Errorf("union size got %d, want 5", union.Size())
	}

	// Intersection
	inter := NewLinkedHashSet(s1.ToSlice()...)
	inter.RetainAll(s2)
	if inter.Size() != 1 || !inter.Contains("banana") {
		t.Errorf("intersection failed: %v", inter.ToSlice())
	}

	// Difference
	diff := NewLinkedHashSet(s1.ToSlice()...)
	diff.RemoveAll(s2)
	if diff.Size() != 2 || diff.Contains("banana") {
		t.Errorf("difference failed: %v", diff.ToSlice())
	}
}

func TestCollectionsUtilities(t *testing.T) {
	list := NewArrayListFrom(5, 3, 8, 1, 9)
	Collections_Sort(list, NaturalOrder[int]()) // [1, 3, 5, 8, 9]
	if list.Get(0) != 1 || list.Get(4) != 9 {
		t.Errorf("Collections_Sort failed: %v", list.ToSlice())
	}

	Collections_Reverse(list) // [9, 8, 5, 3, 1]
	if list.Get(0) != 9 || list.Get(4) != 1 {
		t.Errorf("Collections_Reverse failed: %v", list.ToSlice())
	}

	Collections_Reverse(list) // [1, 3, 5, 8, 9]
	idx := Collections_BinarySearch(list, 5, NaturalOrder[int]())
	if idx != 2 {
		t.Errorf("Collections_BinarySearch got %d, want 2", idx)
	}
	missingIdx := Collections_BinarySearch(list, 4, NaturalOrder[int]())
	if missingIdx >= 0 {
		t.Errorf("Collections_BinarySearch for missing key should be negative")
	}

	minVal := Collections_Min(list, NaturalOrder[int]())
	maxVal := Collections_Max(list, NaturalOrder[int]())
	if minVal != 1 || maxVal != 9 {
		t.Errorf("Min/Max failed: min=%d, max=%d", minVal, maxVal)
	}
}

func TestAuxStructures(t *testing.T) {
	// ArrayDeque LIFO/FIFO
	q := NewArrayDeque[int]()
	q.Push(1)
	q.Push(2)
	if q.Pop() != 2 || q.Pop() != 1 {
		t.Errorf("ArrayDeque stack pop failed")
	}

	// BitSet
	bs := NewBitSet()
	bs.Set(5)
	bs.Set(100)
	if !bs.Get(5) || !bs.Get(100) || bs.Get(6) {
		t.Errorf("BitSet Get failed")
	}
	if bs.Cardinality() != 2 {
		t.Errorf("BitSet Cardinality got %d, want 2", bs.Cardinality())
	}
	if bs.NextSetBit(0) != 5 || bs.NextSetBit(6) != 100 || bs.NextSetBit(101) != -1 {
		t.Errorf("BitSet NextSetBit failed")
	}

	// Optional
	opt1 := Optional_Of(42)
	if !opt1.IsPresent() || opt1.Get() != 42 {
		t.Errorf("Optional_Of failed")
	}
	opt2 := Optional_Empty[int]()
	if opt2.IsPresent() || opt2.OrElse(99) != 99 {
		t.Errorf("Optional_Empty failed")
	}
}
