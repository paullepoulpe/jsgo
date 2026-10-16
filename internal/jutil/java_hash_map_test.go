package jutil

import (
	"slices"
	"testing"
)

// Semantics of JavaHashMap / JavaHashSet / HashMultimap that the JDK vectors
// (java_hash_vectors_test.go) do not exercise.

func TestJavaHashMapFailFast(t *testing.T) {
	m := NewJavaHashMap[string, int]()
	for _, k := range []string{"a", "b", "c"} {
		m.Put(k, 1)
	}
	expectCME := func(name string, f func()) {
		t.Helper()
		thrown := false
		Try(f, func(e IRuntimeException) {
			_, thrown = e.(IConcurrentModificationException)
		})
		if !thrown {
			t.Errorf("%s: no ConcurrentModificationException", name)
		}
	}
	expectCME("All", func() {
		for k := range m.All() {
			if k == "a" {
				m.Put("d", 2)
			}
		}
	})
	// Replacing a value is not a structural change.
	for k := range m.Keys() {
		m.Put(k, 3)
	}
	visited := 0
	expectCME("ForEach", func() {
		m.ForEach(func(k string, v int) {
			visited++
			if k == "b" {
				m.Remove("zz")
				m.Put("e", 1)
			}
		})
	})
	// Java's forEach walks the whole table (including "e", added to a later
	// bin) and only then throws.
	if visited != 5 {
		t.Errorf("ForEach visited %d entries before throwing, want 5", visited)
	}
	expectCME("ComputeIfAbsent", func() {
		m.ComputeIfAbsent("f", func(string) int { m.Put("g", 1); return 1 })
	})
	it := m.Iterator()
	Try(func() { it.Remove(); t.Error("Remove before Next did not throw") },
		func(e IRuntimeException) {
			if _, ok := e.(IIllegalStateException); !ok {
				panic(e)
			}
		})
}

func TestJavaHashCollectionsEquality(t *testing.T) {
	h := NewJavaHashMap[string, int]()
	l := NewLinkedHashMap[string, int]()
	for i, k := range []string{"x", "y", "z", "Aa", "BB"} {
		h.Put(k, i)
		l.Put(k, i)
	}
	if !h.Equals(l) || !l.Equals(h) || h.HashCode() != l.HashCode() {
		t.Error("JavaHashMap and LinkedHashMap with the same entries differ")
	}
	c := NewJavaHashMapFrom[string, int](l)
	c.Put("x", 9)
	if h.Equals(c) {
		t.Error("maps with different values are equal")
	}
	s := NewJavaHashSetFromSlice([]string{"b", "a", "c"})
	ls := NewLinkedHashSet("a", "b", "c")
	if !s.Equals(ls) || !ls.Equals(s) || s.HashCode() != ls.HashCode() {
		t.Error("JavaHashSet and LinkedHashSet with the same elements differ")
	}
	if !h.KeySet().Equals(NewLinkedHashSet("x", "y", "z", "Aa", "BB")) {
		t.Error("keySet view differs from an equal set")
	}
	mm := HashMultimap_Create[string, string]()
	mm2 := MultimapBuilder_HashKeys_HashSetValues[string, string]()
	for _, kv := range [][2]string{{"a", "1"}, {"b", "2"}, {"a", "3"}} {
		mm.Put(kv[0], kv[1])
		mm2.Put(kv[0], kv[1])
	}
	if !mm.Equals(mm2) || mm.HashCode() != mm2.HashCode() {
		t.Error("equal multimaps differ")
	}
	if mm.Equals(h) {
		t.Error("multimap equals a map")
	}
	if got := mm.String(); got != "{a=[1, 3], b=[2]}" {
		t.Errorf("multimap toString %s", got)
	}
}

func TestJavaHashMapToString(t *testing.T) {
	m := NewJavaHashMap[any, any]()
	if m.String() != "{}" {
		t.Errorf("empty map %s", m.String())
	}
	m.Put("k", m)
	m.Put(int32(1), nil)
	if got := m.String(); got != "{1=null, k=(this Map)}" {
		t.Errorf("toString %s", got)
	}
	if got := NewJavaHashSetFromSlice([]int32{3, 1, 2}).String(); got != "[1, 2, 3]" {
		t.Errorf("set toString %s", got)
	}
}

func TestGuavaHashCapacities(t *testing.T) {
	// Maps.capacity in the shaded Guava (javap): n+1 below 3, ceil(n/0.75)
	// below 2^30.
	for n, want := range map[int]int{0: 1, 1: 2, 2: 3, 3: 4, 4: 6, 6: 8, 12: 16, 13: 18, 96: 128, 1 << 30: 1<<31 - 1} {
		if got := Maps_Capacity(n); got != want {
			t.Errorf("Maps.capacity(%d) = %d, want %d", n, got, want)
		}
	}
	for n, want := range map[int]int32{0: 1, 1: 1, 2: 2, 3: 4, 16: 16, 17: 32, 1 << 30: 1 << 30} {
		if got := javaHashMapTableSizeFor(int32(n)); got != want {
			t.Errorf("tableSizeFor(%d) = %d, want %d", n, got, want)
		}
	}
	Try(func() { Maps_Capacity(-1); t.Error("no exception") }, func(e IRuntimeException) {
		if e.GetMessage() != "expectedSize cannot be negative but was: -1" {
			t.Errorf("message %v", e.GetMessage())
		}
	})
	// HashMultimap.removeAll copies into a new set sized for
	// expectedValuesPerKey, which can iterate differently from the old one:
	// the old set grew to 32 bins (1 before 4), the copy has 4 (4 before 1).
	mm := HashMultimap_CreateWithExpectedSizes[int32, int32](1, 0)
	for v := int32(1); v <= 20; v++ {
		mm.Put(5, v)
	}
	for v := int32(1); v <= 20; v++ {
		if v != 1 && v != 4 {
			mm.Remove(5, v)
		}
	}
	if got := slices.Collect(mm.Get(5).All()); !slices.Equal(got, []int32{1, 4}) {
		t.Errorf("get %v", got)
	}
	if got := slices.Collect(mm.RemoveAll(5).All()); !slices.Equal(got, []int32{4, 1}) {
		t.Errorf("removeAll %v", got)
	}
	if !mm.IsEmpty() || mm.ContainsKey(5) {
		t.Error("removeAll left the key")
	}
}

func BenchmarkJavaHashMapStringPutGet(b *testing.B) {
	keys := make([]string, 1000)
	for i := range keys {
		keys[i] = "key" + ToJavaString(i)
	}
	b.ReportAllocs()
	for b.Loop() {
		m := NewJavaHashMap[string, int]()
		for i, k := range keys {
			m.Put(k, i)
		}
		for _, k := range keys {
			_ = m.Get(k)
		}
	}
}
