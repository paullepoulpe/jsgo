package jutil

import "testing"

// ArrayListMultimap.remove(key, value): first equal value only; the key goes
// when its last value does (Guava AbstractMapBasedMultimap semantics).
func TestListMultimapRemove(t *testing.T) {
	mm := ArrayListMultimap_Create[string, string]()
	mm.Put("a", "x")
	mm.Put("a", "y")
	mm.Put("a", "x")
	mm.Put("b", "z")
	if mm.Remove("a", "q") || mm.Remove("c", "x") {
		t.Fatal("removed an absent entry")
	}
	if !mm.Remove("a", "x") {
		t.Fatal("did not remove a:x")
	}
	if got := mm.String(); got != "{a=[y, x], b=[z]}" {
		t.Fatalf("got %s", got)
	}
	mm.Put("a", "x")
	if got := mm.String(); got != "{a=[y, x, x], b=[z]}" {
		t.Fatalf("got %s", got)
	}
	if !mm.Remove("b", "z") || mm.ContainsKey("b") || mm.Size() != 3 {
		t.Fatalf("got %s size %d", mm.String(), mm.Size())
	}
	mm.Put("b", "w")
	if got := mm.String(); got != "{a=[y, x, x], b=[w]}" {
		t.Fatalf("got %s", got)
	}
}
