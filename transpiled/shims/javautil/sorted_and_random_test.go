package javautil_test

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
)

func TestRandomLCGAndShuffleParity(t *testing.T) {
	// Standard java.util.Random(12345L) sequence per Javadoc 48-bit LCG:
	r := javautil.NewRandomSeed(12345)
	if got := r.NextInt(); got != 1553932502 {
		t.Fatalf("Random(12345).NextInt() = %d, want 1553932502", got)
	}
	if got := r.NextIntBound(100); got != 80 {
		t.Fatalf("Random(12345) 2nd NextInt(100) = %d, want 80", got)
	}

	list := javautil.NewArrayListFrom("a", "b", "c", "d", "e")
	javautil.Collections_Shuffle(list, javautil.NewRandomSeed(42))
	list2 := javautil.NewArrayListFrom("a", "b", "c", "d", "e")
	javautil.Collections_Shuffle(list2, javautil.NewRandomSeed(42))
	if !slices.Equal(list.ToSlice(), list2.ToSlice()) {
		t.Fatalf("Collections_Shuffle non-deterministic: %v vs %v", list.ToSlice(), list2.ToSlice())
	}
}

func TestTreeMapTreeSetAndPriorityQueue(t *testing.T) {
	tm := javautil.NewTreeMap[string, int]()
	tm.Put("cherry", 3)
	tm.Put("apple", 1)
	tm.Put("banana", 2)
	if tm.FirstKey() != "apple" || tm.LastKey() != "cherry" {
		t.Fatalf("TreeMap FirstKey/LastKey mismatch: %s / %s", tm.FirstKey(), tm.LastKey())
	}
	if fk, ok := tm.FloorKey("blueberry"); !ok || fk != "banana" {
		t.Fatalf("TreeMap FloorKey got %q, %v", fk, ok)
	}
	if ck, ok := tm.CeilingKey("blueberry"); !ok || ck != "cherry" {
		t.Fatalf("TreeMap CeilingKey got %q, %v", ck, ok)
	}

	ts := javautil.NewTreeSet[int]()
	for _, x := range []int{5, 1, 3, 1, 4, 2} {
		ts.Add(x)
	}
	if !slices.Equal(ts.ToSlice(), []int{1, 2, 3, 4, 5}) {
		t.Fatalf("TreeSet order mismatch: %v", ts.ToSlice())
	}

	pq := javautil.NewPriorityQueue[int]()
	for _, x := range []int{9, 4, 7, 1, 3} {
		pq.Offer(x)
	}
	var popped []int
	for !pq.IsEmpty() {
		v, _ := pq.Poll()
		popped = append(popped, v)
	}
	if !slices.Equal(popped, []int{1, 3, 4, 7, 9}) {
		t.Fatalf("PriorityQueue order mismatch: %v", popped)
	}
}
