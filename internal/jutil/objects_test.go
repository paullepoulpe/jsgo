package jutil

import "testing"

func TestObjectLinkedHashCollectionsEqualsAndHashZeroAllocs(t *testing.T) {
	m1 := NewObjectLinkedHashMap[*ptrHashableKey, int]()
	m2 := NewObjectLinkedHashMap[*ptrHashableKey, int]()
	s1 := NewObjectLinkedHashSet[*ptrHashableKey]()
	s2 := NewObjectLinkedHashSet[*ptrHashableKey]()
	ls1 := NewLinkedHashSet[string]()
	ls2 := NewLinkedHashSet[string]()

	for i := 0; i < 20; i++ {
		m1.Put(&ptrHashableKey{i}, i*10)
		m2.Put(&ptrHashableKey{19 - i}, (19-i)*10)
		s1.Add(&ptrHashableKey{i})
		s2.Add(&ptrHashableKey{19 - i})
	}
	ls1.AddAll("alpha", "beta", "gamma")
	ls2.AddAll("gamma", "alpha", "beta")

	if !m1.Equals(m2) || m1.HashCode() != m2.HashCode() {
		t.Fatal("ObjectLinkedHashMap Equals/HashCode mismatch")
	}
	if !s1.Equals(s2) || s1.HashCode() != s2.HashCode() {
		t.Fatal("ObjectLinkedHashSet Equals/HashCode mismatch")
	}
	if !ls1.Equals(ls2) || ls1.HashCode() != ls2.HashCode() {
		t.Fatal("LinkedHashSet Equals/HashCode mismatch")
	}
	if !ls1.Equals(NewObjectLinkedHashSet("alpha", "beta", "gamma")) {
		t.Fatal("cross-set Equals mismatch")
	}

	strSlice := []string{"alpha", "beta", "gamma"}
	if n := testing.AllocsPerRun(100, func() {
		if !m1.Equals(m2) || m1.HashCode() != m2.HashCode() ||
			!s1.Equals(s2) || s1.HashCode() != s2.HashCode() ||
			!ls1.Equals(ls2) || ls1.HashCode() != ls2.HashCode() ||
			SetHashCode(strSlice) != ls1.HashCode() {
			t.Fatal("unexpected mismatch")
		}
	}); n != 0 {
		t.Fatalf("expected 0 allocations, got %v", n)
	}
}
