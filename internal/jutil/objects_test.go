package jutil

import "testing"

func TestLinkedHashCollectionsEqualsAndHash(t *testing.T) {
	m1 := NewLinkedHashMap[int, *ptrHashableKey]()
	m2 := NewLinkedHashMap[int, *ptrHashableKey]()
	ls1 := NewLinkedHashSet[string]()
	ls2 := NewLinkedHashSet[string]()

	for i := 0; i < 20; i++ {
		m1.Put(i, &ptrHashableKey{i})
		m2.Put(19-i, &ptrHashableKey{19 - i})
	}
	ls1.AddAll("alpha", "beta", "gamma")
	ls2.AddAll("gamma", "alpha", "beta")

	if !m1.Equals(m2) || m1.HashCode() != m2.HashCode() {
		t.Fatal("LinkedHashMap Equals/HashCode mismatch")
	}
	if !ls1.Equals(ls2) || ls1.HashCode() != ls2.HashCode() {
		t.Fatal("LinkedHashSet Equals/HashCode mismatch")
	}
	strSlice := []string{"alpha", "beta", "gamma"}
	if SetHashCode(strSlice) != ls1.HashCode() {
		t.Fatal("SetHashCode mismatch")
	}

	if n := testing.AllocsPerRun(100, func() {
		if !ls1.Equals(ls2) {
			t.Fatal("unexpected mismatch")
		}
	}); n != 0 {
		t.Fatalf("expected 0 allocations for LinkedHashSet.Equals, got %v", n)
	}
}
