package guava

import (
	"testing"
)

func TestImmutableCollections(t *testing.T) {
	// ImmutableList
	list := ImmutableList_Of("alpha", "beta", "gamma")
	if list.Size() != 3 {
		t.Errorf("list.Size() got %d, want 3", list.Size())
	}
	if list.Get(1) != "beta" {
		t.Errorf("list.Get(1) got %q, want \"beta\"", list.Get(1))
	}
	if !list.Contains("gamma") {
		t.Errorf("list should contain \"gamma\"")
	}

	// ImmutableSet
	set := ImmutableSet_Of("k1", "k2", "k3", "k1")
	if set.Size() != 3 {
		t.Errorf("set deduplication failed: got %d, want 3", set.Size())
	}
	if !set.Contains("k2") {
		t.Errorf("set should contain \"k2\"")
	}

	// ImmutableMap
	m := ImmutableMap_Of2("key1", 100, "key2", 200)
	if m.Size() != 2 {
		t.Errorf("m.Size() got %d, want 2", m.Size())
	}
	if m.GetVal("key1") != 100 || m.GetVal("key2") != 200 {
		t.Errorf("m.GetVal values incorrect")
	}
}

func TestListsMapsSetsHelpers(t *testing.T) {
	// Lists_Transform
	orig := Lists_NewArrayListFrom(1, 2, 3, 4)
	transformed := Lists_Transform(orig, func(x int) int { return x * 10 })
	if transformed.Get(0) != 10 || transformed.Get(3) != 40 {
		t.Errorf("Lists_Transform failed: %v", transformed.ToSlice())
	}

	// Maps_UniqueIndex
	type Item struct {
		id   string
		name string
	}
	items := []Item{
		{"id1", "Item One"},
		{"id2", "Item Two"},
	}
	idxMap := Maps_UniqueIndex(items, func(item Item) string { return item.id })
	if idxMap.Size() != 2 || idxMap.GetVal("id1").name != "Item One" {
		t.Errorf("Maps_UniqueIndex failed")
	}

	// Sets helpers
	s1 := Sets_NewHashSet("A", "B", "C")
	s2 := Sets_NewHashSet("B", "C", "D")

	diff := Sets_Difference(s1, s2)
	if diff.Size() != 1 || !diff.Contains("A") {
		t.Errorf("Sets_Difference failed: %v", diff.ToSlice())
	}

	inter := Sets_Intersection(s1, s2)
	if inter.Size() != 2 || !inter.Contains("B") || !inter.Contains("C") {
		t.Errorf("Sets_Intersection failed: %v", inter.ToSlice())
	}

	union := Sets_Union(s1, s2)
	if union.Size() != 4 {
		t.Errorf("Sets_Union failed: %v", union.ToSlice())
	}
}
