package jutil

import (
	"slices"
	"testing"
)

func TestOrderedSetMultimap_BasicAndOrder(t *testing.T) {
	var nilMap *OrderedSetMultimap[string, int]
	if !nilMap.IsEmpty() || nilMap.Size() != 0 || nilMap.String() != "{}" {
		t.Fatalf("nil multimap mismatch: empty=%v size=%d str=%q", nilMap.IsEmpty(), nilMap.Size(), nilMap.String())
	}
	if !nilMap.Get("x").IsEmpty() || len(nilMap.Get("x").ToSlice()) != 0 {
		t.Fatalf("nil Get mismatch")
	}

	m := NewOrderedSetMultimap[string, int]()
	if !m.IsEmpty() || m.String() != "{}" {
		t.Fatalf("empty multimap mismatch: %v %q", m.IsEmpty(), m.String())
	}

	if !m.Put("b", 2) || !m.Put("a", 1) || !m.Put("b", 3) {
		t.Fatalf("expected initial puts to succeed")
	}
	if m.Put("b", 2) || m.Put("a", 1) {
		t.Fatalf("expected duplicate puts to return false")
	}

	if m.Size() != 3 || m.IsEmpty() {
		t.Fatalf("got size=%d empty=%v, want 3 false", m.Size(), m.IsEmpty())
	}
	if got := m.String(); got != "{b=[2, 3], a=[1]}" {
		t.Fatalf("got %q, want %q", got, "{b=[2, 3], a=[1]}")
	}

	bSet := m.Get("b")
	if bSet.Size() != 2 || !bSet.Contains(2) || !bSet.Contains(3) || bSet.Contains(9) {
		t.Fatalf("Get(b) mismatch: %v", bSet)
	}
	if got := slices.Collect(bSet.All()); !slices.Equal(got, []int{2, 3}) {
		t.Fatalf("Get(b).All() = %v, want [2 3]", got)
	}

	var keys []string
	var vals []int
	for k, v := range m.Entries() {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	if !slices.Equal(keys, []string{"b", "b", "a"}) || !slices.Equal(vals, []int{2, 3, 1}) {
		t.Fatalf("Entries() = (%v, %v)", keys, vals)
	}
}

func TestOrderedSetMultimap_CloneAndInverse(t *testing.T) {
	m := NewOrderedSetMultimap[string, string]()
	m.Put("c1", "loc0")
	m.Put("c2", "loc0")
	m.Put("c1", "loc1")

	cloned := m.Clone()
	cloned.Put("c3", "loc2")
	cloned.Put("c1", "loc9")
	if m.Size() != 3 || cloned.Size() != 5 {
		t.Fatalf("clone mutated original: m=%d cloned=%d", m.Size(), cloned.Size())
	}
	if !slices.Equal(m.Get("c1").ToSlice(), []string{"loc0", "loc1"}) {
		t.Fatalf("original c1 mutated: %v", m.Get("c1").ToSlice())
	}

	inv := m.Inverse()
	if inv.Size() != 3 {
		t.Fatalf("inverse size=%d, want 3", inv.Size())
	}
	if got := inv.String(); got != "{loc0=[c1, c2], loc1=[c1]}" {
		t.Fatalf("inverse = %q, want %q", got, "{loc0=[c1, c2], loc1=[c1]}")
	}

	var mapKeys []string
	for k := range inv.AsMap() {
		mapKeys = append(mapKeys, k)
	}
	if !slices.Equal(mapKeys, []string{"loc0", "loc1"}) {
		t.Fatalf("AsMap keys = %v", mapKeys)
	}
}
