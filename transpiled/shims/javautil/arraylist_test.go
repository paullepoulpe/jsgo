package javautil

import (
	"testing"
)

func TestArrayListBasics(t *testing.T) {
	list := NewArrayList[string]()
	if !list.IsEmpty() {
		t.Errorf("expected empty list")
	}
	list.Add("a")
	list.Add("b")
	list.Add("c")

	if list.Size() != 3 {
		t.Errorf("list.Size() got %d, want 3", list.Size())
	}
	if list.Get(1) != "b" {
		t.Errorf("list.Get(1) got %q, want \"b\"", list.Get(1))
	}
	if list.Set(1, "B") != "b" {
		t.Errorf("list.Set(1, \"B\") should return old value \"b\"")
	}
	if list.Get(1) != "B" {
		t.Errorf("list.Get(1) got %q, want \"B\"", list.Get(1))
	}

	list.AddAt(1, "x") // ["a", "x", "B", "c"]
	if list.Size() != 4 || list.Get(1) != "x" || list.Get(2) != "B" {
		t.Errorf("AddAt failed: %v", list.ToSlice())
	}

	rem := list.Remove(1) // ["a", "B", "c"]
	if rem != "x" || list.Size() != 3 {
		t.Errorf("Remove failed: rem=%q, size=%d", rem, list.Size())
	}

	if !list.Contains("B") {
		t.Errorf("list should contain \"B\"")
	}
	if list.IndexOf("B") != 1 {
		t.Errorf("IndexOf(\"B\") got %d, want 1", list.IndexOf("B"))
	}

	// Sublist
	sub := list.SubList(0, 2)
	if sub.Size() != 2 || sub.Get(0) != "a" || sub.Get(1) != "B" {
		t.Errorf("SubList failed: %v", sub.ToSlice())
	}

	// Sort
	list.Sort(NaturalOrder[string]()) // ["B", "a", "c"]
	if list.Get(0) != "B" || list.Get(1) != "a" || list.Get(2) != "c" {
		t.Errorf("Sort failed: %v", list.ToSlice())
	}
}

func TestArrayListIteratorAndSequence(t *testing.T) {
	list := NewArrayListFrom("x", "y", "z")
	it := list.Iterator()
	var elems []string
	for it.HasNext() {
		elems = append(elems, it.Next())
	}
	if len(elems) != 3 || elems[0] != "x" || elems[1] != "y" || elems[2] != "z" {
		t.Errorf("Iterator traversal failed: %v", elems)
	}

	var seqElems []string
	for e := range list.All() {
		seqElems = append(seqElems, e)
	}
	if len(seqElems) != 3 || seqElems[0] != "x" || seqElems[1] != "y" || seqElems[2] != "z" {
		t.Errorf("All() traversal failed: %v", seqElems)
	}
}
