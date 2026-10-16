package jutil

import "testing"

func TestObjectsHashAndCollectionHashCodes(t *testing.T) {
	// Objects.hash("a", 1) == 31*(31*1+97)+1 == 3969
	if h := ObjectsHash("a", 1); h != 3969 {
		t.Fatalf("Objects.hash(\"a\", 1) = %d", h)
	}
	if h := ObjectsHash(nil); h != 31 {
		t.Fatalf("Objects.hash(null) = %d", h)
	}
	if h := ListHashCode([]string{"a", "b"}); h != 31*(31+97)+98 {
		t.Fatalf("List.hashCode = %d", h)
	}
	if h := SetHashCode([]string{"a", "b"}); h != 97+98 {
		t.Fatalf("Set.hashCode = %d", h)
	}
}

func TestBitSetWithSize(t *testing.T) {
	// Java: new BitSet().size() == 64, new BitSet(0).size() == 0,
	// new BitSet(65).size() == 128.
	if s := NewBitSetWithSize(64).Size(); s != 64 {
		t.Fatalf("size = %d", s)
	}
	if s := NewBitSetWithSize(0).Size(); s != 0 {
		t.Fatalf("size = %d", s)
	}
	b := NewBitSetWithSize(65)
	if b.Size() != 128 || b.Length() != 0 {
		t.Fatalf("size = %d", b.Size())
	}
	b.Set(70)
	if b.Size() != 128 || b.Length() != 71 || !b.Get(70) {
		t.Fatalf("size = %d length = %d", b.Size(), b.Length())
	}
}

func TestTreeMapKeySliceAndValues(t *testing.T) {
	m := NewTreeMap[string, int](StringCompareTo)
	m.Put("b", 2)
	m.Put("a", 1)
	if k := m.KeySlice(); len(k) != 2 || k[0] != "a" || k[1] != "b" {
		t.Fatalf("keys = %v", k)
	}
	var vs []int
	for v := range m.Values() {
		vs = append(vs, v)
	}
	if len(vs) != 2 || vs[0] != 1 || vs[1] != 2 {
		t.Fatalf("values = %v", vs)
	}
}
