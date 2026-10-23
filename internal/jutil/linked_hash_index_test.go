package jutil

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
	"unsafe"
)

// Maps far past lhmSmallSize: the index grows several times, tombstones
// pile up and are dropped, and the slots are compacted.
func TestLinkedHashMapIndexedLarge(t *testing.T) {
	for seed := int64(1); seed <= 2; seed++ {
		for _, domain := range []int{200, 1500} {
			checkLinkedHashMapAgainstModel(t, seed, domain, func(i int) string { return fmt.Sprint("k", i) })
			checkLinkedHashMapAgainstModel(t, seed, domain, func(i int) int32 { return int32(i * 7919) })
		}
	}
}

// checkLinkedHashMapKeys checks h's iteration order and that every key of
// want is found with its position as value and absent is not found.
func checkLinkedHashMapKeys(t *testing.T, what string, h *LinkedHashMap[string, int], want []string, absent []string) {
	t.Helper()
	if got := h.KeySlice(); !slices.Equal(got, want) {
		t.Fatalf("%s: keys %v, want %v", what, got, want)
	}
	for _, k := range want {
		if !h.ContainsKey(k) {
			t.Fatalf("%s: %s not found", what, k)
		}
	}
	for _, k := range absent {
		if h.ContainsKey(k) {
			t.Fatalf("%s: removed %s found", what, k)
		}
	}
	if h.Size() != len(want) {
		t.Fatalf("%s: size %d, want %d", what, h.Size(), len(want))
	}
}

// Removing most entries from the middle, then adding: the slots are
// compacted and the index is rebuilt for the moved entries.
func TestLinkedHashMapCompactAfterRemovals(t *testing.T) {
	h := NewLinkedHashMap[string, int]()
	var keys []string
	for i := 0; i < 500; i++ {
		k := fmt.Sprint("k", i)
		keys = append(keys, k)
		h.Put(k, i)
	}
	var kept, removed []string
	for i, k := range keys {
		if i%10 == 3 {
			kept = append(kept, k)
		} else {
			h.Remove(k)
			removed = append(removed, k)
		}
	}
	checkLinkedHashMapKeys(t, "after removals", h, kept, removed)
	for i := 500; i < 520; i++ {
		k := fmt.Sprint("k", i)
		kept = append(kept, k)
		h.Put(k, i)
	}
	if len(h.slots) != len(kept) {
		t.Fatalf("slots not compacted: %d slots for %d entries", len(h.slots), len(kept))
	}
	checkLinkedHashMapKeys(t, "after compaction", h, kept, removed)
	// Re-inserting removed keys appends them at the end.
	for _, k := range removed[:5] {
		h.Put(k, 0)
		kept = append(kept, k)
	}
	checkLinkedHashMapKeys(t, "after re-insertion", h, kept, removed[5:])
}

// A map used as a queue or a working set of constant size keeps a bounded
// index: tombstones are dropped when the table is rebuilt.
func TestLinkedHashMapTombstoneChurn(t *testing.T) {
	for _, fromFront := range []bool{true, false} {
		h := NewLinkedHashMap[string, int]()
		var live, gone []string
		for i := 0; i < 20000; i++ {
			k := fmt.Sprint("k", i)
			h.Put(k, i)
			live = append(live, k)
			if len(live) > 40 {
				j := 0
				if !fromFront {
					j = len(live) / 2
				}
				h.Remove(live[j])
				gone = append(gone, live[j])
				live = slices.Delete(live, j, j+1)
			}
			if h.m != nil && len(h.m.table) > 128 {
				t.Fatalf("fromFront=%v step %d: table of %d buckets for %d entries", fromFront, i, len(h.m.table), h.Size())
			}
		}
		checkLinkedHashMapKeys(t, fmt.Sprint("fromFront=", fromFront), h, live, gone)
	}
}

// ObjectLinkedHashMap far past lhmSmallSize, with hash codes that collide
// in groups (id % 50) and through the table mask, against the model.
func TestObjectLinkedHashMapIndexedLarge(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	h := NewObjectLinkedHashMap[olhmKey, int]()
	var m orderedModel[olhmKey]
	for step := 0; step < 30000; step++ {
		id := r.Intn(600)
		k := olhmKey{id, id % 50 * 1024}
		switch op := r.Intn(100); {
		case op < 50:
			v := r.Intn(1000) + 1
			if got, want := h.Put(k, v), func() int { o, _ := m.put(k, v); return o }(); got != want {
				t.Fatalf("step %d: Put(%v) = %d, want %d", step, k, got, want)
			}
		case op < 80:
			got := h.Remove(k)
			if want, _ := m.remove(k); got != want {
				t.Fatalf("step %d: Remove(%v) = %d, want %d", step, k, got, want)
			}
		case op < 81 && step%7 == 0:
			h.Clear()
			m = orderedModel[olhmKey]{}
		default:
			v, ok := h.Lookup(k)
			if i := m.index(k); ok != (i >= 0) || ok && v != m.vals[i] {
				t.Fatalf("step %d: Lookup(%v) = %d, %v; model index %d", step, k, v, ok, i)
			}
		}
		if h.Size() != len(m.keys) {
			t.Fatalf("step %d: size %d, want %d", step, h.Size(), len(m.keys))
		}
		if step%100 == 0 {
			if got := h.KeySlice(); !slices.Equal(got, m.keys) && !(len(got) == 0 && len(m.keys) == 0) {
				t.Fatalf("step %d: keys %v, want %v", step, got, m.keys)
			}
		}
	}
}

// Removing every entry, from the front and from the middle, then clearing
// and refilling the map.
func TestLinkedHashMapIndexedRemoveAllAndClear(t *testing.T) {
	h := NewLinkedHashMap[int, int]()
	for i := 0; i < 100; i++ {
		h.Put(i, i)
	}
	for i := 50; i < 100; i++ {
		h.Remove(i)
	}
	for i := 0; i < 50; i++ {
		h.Remove(i)
	}
	if h.Size() != 0 || len(h.KeySlice()) != 0 || h.ContainsKey(0) || h.ContainsKey(99) {
		t.Fatalf("not empty: %v", h.KeySlice())
	}
	h.Put(7, 7)
	h.Clear()
	for i := 0; i < 30; i++ {
		h.Put(i*3, i)
	}
	for i := 0; i < 90; i++ {
		v, ok := h.Lookup(i)
		if ok != (i%3 == 0) || ok && v != i/3 {
			t.Fatalf("Lookup(%d) = %d, %v", i, v, ok)
		}
	}
	if k, _ := h.KeyAt(29); k != 87 {
		t.Fatalf("KeyAt(29) = %d", k)
	}
}

// Zero-value keys ("" and nil) behave as ordinary keys in both modes, also
// after slots holding zeroed (removed) entries exist.
func TestLinkedHashMapZeroKeys(t *testing.T) {
	for _, n := range []int{2, 3 * lhmSmallSize} {
		h := NewLinkedHashMap[string, int]()
		for i := 0; i < n; i++ {
			h.Put(fmt.Sprint(i), i)
		}
		h.Remove("1")
		if h.ContainsKey("") {
			t.Fatalf("n=%d: empty key found before insertion", n)
		}
		h.Put("", -1)
		if v, ok := h.Lookup(""); !ok || v != -1 {
			t.Fatalf("n=%d: Lookup(\"\") = %d, %v", n, v, ok)
		}
		h.Remove("")
		if h.ContainsKey("") || h.Size() != n-1 {
			t.Fatalf("n=%d: empty key still found", n)
		}
	}
}

// LinkedHashSet slots carry no padding for their struct{} values.
func TestLinkedHashSetSlotSize(t *testing.T) {
	var p LhmSlot[*int, struct{}]
	var s LhmSlot[string, struct{}]
	if got := unsafe.Sizeof(p); got != 8 {
		t.Errorf("LhmSlot[*int, struct{}] is %d bytes", got)
	}
	if got := unsafe.Sizeof(s); got != 16 {
		t.Errorf("LhmSlot[string, struct{}] is %d bytes", got)
	}
}

type lhmGenericNode[T any] interface{ get() T }

type lhmGenericNodeImpl[T any] struct{ v T }

func (n *lhmGenericNodeImpl[T]) get() T { return n.v }

// fillGenericInterfaceSet builds, in generic code, a set whose element type
// is a generic interface type (the case lhmHash hashes as any).
func fillGenericInterfaceSet[T any](vals []T) ([]lhmGenericNode[T], *LinkedHashSet[lhmGenericNode[T]]) {
	var nodes []lhmGenericNode[T]
	s := NewLinkedHashSet[lhmGenericNode[T]]()
	for _, v := range vals {
		n := &lhmGenericNodeImpl[T]{v}
		nodes = append(nodes, n)
		s.Add(n)
		s.Add(n)
	}
	return nodes, s
}

func TestLinkedHashSetGenericInterfaceKeys(t *testing.T) {
	nodes, s := fillGenericInterfaceSet([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	if !slices.Equal(s.ToSlice(), nodes) {
		t.Fatalf("elements %v", s.ToSlice())
	}
	for _, n := range nodes {
		if !s.Contains(n) {
			t.Fatalf("%v not found", n)
		}
	}
	if s.Contains(&lhmGenericNodeImpl[int]{1}) {
		t.Fatal("found a new node")
	}
}

func TestLinkedHashMapConcurrentReadsLazyIndex(t *testing.T) {
	for trial := 0; trial < 32; trial++ {
		h := NewLinkedHashMap[string, int]()
		s := NewLinkedHashSet[string]()
		const n = 24
		keys := make([]string, n)
		for i := 0; i < n; i++ {
			k := fmt.Sprintf("key_%d_%d", trial, i)
			keys[i] = k
			h.PutNew(k, i)
			s.AddNew(k)
		}
		if h.m != nil || s.m.m != nil {
			t.Fatal("expected nil index before first lookup")
		}

		const workers = 8
		start := make(chan struct{})
		done := make(chan struct{}, workers)
		for w := 0; w < workers; w++ {
			w := w
			go func() {
				defer func() { done <- struct{}{} }()
				<-start
				for i, k := range keys {
					idx := (i + w) % len(keys)
					qk := keys[idx]
					if got, ok := h.Lookup(qk); !ok || got != idx {
						t.Errorf("Lookup(%q) = (%d, %v), want (%d, true)", qk, got, ok, idx)
						return
					}
					if !s.Contains(k) {
						t.Errorf("Contains(%q) = false", k)
						return
					}
				}
				if h.ContainsKey("missing") || s.Contains("missing") {
					t.Errorf("unexpected hit for missing key")
					return
				}
				cnt := 0
				for range h.All() {
					cnt++
				}
				if cnt != n {
					t.Errorf("All count = %d, want %d", cnt, n)
				}
				_ = h.KeySlice()
				_ = h.HashCode()
				_ = s.ToSlice()
				_ = s.HashCode()
				_, _ = s.First()
				_, _, _ = s.FirstTwo()
			}()
		}
		close(start)
		for w := 0; w < workers; w++ {
			<-done
		}
	}
}
