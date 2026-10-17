package jutil

import (
	"slices"
	"testing"
)

// Go-only: TreeMultiset and TreeMultimap against what Guava prints for the
// same calls (TreeMultiset.create() / TreeMultimap.create() over Strings).
func TestTreeMultisetTreeMultimap(t *testing.T) {
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

	mm := TreeMultimap_Create[string, string](StringCompareTo, StringCompareTo)
	if !mm.Put("lang", "es6") || !mm.Put("module", "goog") || !mm.Put("lang", "es5") || mm.Put("lang", "es6") {
		t.Errorf("TreeMultimap.put results are wrong")
	}
	if got, want := mm.String(), "{lang=[es5, es6], module=[goog]}"; got != want || mm.Size() != 3 {
		t.Errorf("TreeMultimap = %q (size %d), want %q", got, mm.Size(), want)
	}
	var pairs []string
	for k, v := range mm.Entries() {
		pairs = append(pairs, k+"="+v)
	}
	if !slices.Equal(pairs, []string{"lang=es5", "lang=es6", "module=goog"}) {
		t.Errorf("entries = %q", pairs)
	}
}
