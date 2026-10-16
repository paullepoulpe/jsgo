package jutil

import "testing"

// Expected values are those of java.util.BitSet (JDK 21).
func TestBitSetMatchesJava(t *testing.T) {
	b := NewBitSet()
	if b.Size() != 64 || b.Length() != 0 {
		t.Fatalf("new BitSet(): size=%d length=%d", b.Size(), b.Length())
	}
	z := NewBitSetWithSize(0)
	if z.Size() != 0 {
		t.Fatalf("new BitSet(0).size() = %d, want 0", z.Size())
	}
	if s := NewBitSetWithSize(65).Size(); s != 128 {
		t.Fatalf("new BitSet(65).size() = %d, want 128", s)
	}
	z.Set(3)
	// ensureCapacity(1) from 0 words: max(0, 1) = 1 word.
	if z.Size() != 64 || z.Length() != 4 || !z.Get(3) || z.Get(2) {
		t.Fatalf("after set(3): size=%d length=%d", z.Size(), z.Length())
	}
	b.Set(130)
	// ensureCapacity(3) from 1 word: max(2, 3) = 3 words.
	if b.Size() != 192 || b.Length() != 131 {
		t.Fatalf("after set(130): size=%d length=%d", b.Size(), b.Length())
	}
	b.SetValue(130, false)
	if b.Length() != 0 || b.Get(130) {
		t.Fatalf("after clear(130): length=%d", b.Length())
	}
	o := NewBitSetWithSize(10)
	o.Set(1)
	o.Set(70)
	z.Or(o)
	if !z.Get(1) || !z.Get(3) || !z.Get(70) || z.Length() != 71 {
		t.Fatalf("or: length=%d", z.Length())
	}
}

func TestObjectsHash(t *testing.T) {
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
