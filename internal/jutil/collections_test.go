package jutil

import (
	"fmt"
	"slices"
	"sync"
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

func TestOrderedMap_ConcurrentReadsWithHoles(t *testing.T) {
	for _, size := range []int{6, 24} {
		m := NewOrderedMap[string, int]()
		s := NewOrderedSet[string]()
		for i := 0; i < size; i++ {
			k := fmt.Sprintf("k%d", i)
			m.Put(k, i)
			s.Add(k)
		}
		// Leave tombstones (holes > 0) before concurrent read-only access.
		m.Remove("k1")
		m.Remove("k3")
		s.Remove("k1")
		s.Remove("k3")

		var wg sync.WaitGroup
		for g := 0; g < 8; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for iter := 0; iter < 50; iter++ {
					if m.Size() != size-2 || s.Size() != size-2 {
						t.Errorf("unexpected size: m=%d s=%d", m.Size(), s.Size())
					}
					_ = m.KeySlice()
					_ = m.ValueSlice()
					_ = m.Entries()
					_, _, _ = m.First()
					_, _ = m.ValueAt(1)
					_, _ = m.Index("k2")
					_ = m.Clone()
					_ = m.ToMap()
					_ = m.String()
					_ = m.HashCode()
					_ = m.Equals(m)
					_ = s.ToSlice()
					_, _ = s.First()
					_ = s.Clone()
					_ = s.String()
					_ = s.HashCode()
					_ = s.Equals(s)
				}
			}()
		}
		wg.Wait()
	}
}
