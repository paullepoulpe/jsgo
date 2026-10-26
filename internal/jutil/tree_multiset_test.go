package jutil

import (
	"slices"
	"testing"
)

// Go-only: TreeMultiset against what Guava prints for the same calls
// (TreeMultiset.create() over Strings).
func TestTreeMultiset(t *testing.T) {
	ms := TreeMultiset_Create[string](StringCompareTo)
	for _, s := range []string{"b", "a", "b", "é", "B"} {
		ms.Add(s)
	}
	if got, want := ms.String(), "[B, a, b x 2, é]"; got != want {
		t.Errorf("TreeMultiset = %q, want %q", got, want)
	}
	if got := ms.ToSlice(); !slices.Equal(got, []string{"B", "a", "b", "b", "é"}) || ms.Size() != 5 || ms.Count("b") != 2 {
		t.Errorf("TreeMultiset iteration = %q, size %d", got, ms.Size())
	}
	if got := ms.ElementSet(); !slices.Equal(got, []string{"B", "a", "b", "é"}) {
		t.Errorf("elementSet = %q", got)
	}
}
