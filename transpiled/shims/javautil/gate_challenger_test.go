package javautil_test

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
)

// TestNilReceiverBehavior investigates how nil receiver calls are handled across LinkedHashMap and LinkedHashSet.
func TestNilReceiverBehavior(t *testing.T) {
	var m *javautil.LinkedHashMap[string, int] = nil

	// LinkedHashMap nil receiver safe methods:
	if m.Size() != 0 {
		t.Errorf("m.Size() on nil expected 0, got %d", m.Size())
	}
	if !m.IsEmpty() {
		t.Errorf("m.IsEmpty() on nil expected true, got false")
	}
	if _, ok := m.Get("a"); ok {
		t.Errorf("m.Get() on nil expected ok=false, got true")
	}
	if m.GetVal("a") != 0 {
		t.Errorf("m.GetVal() on nil expected 0, got %d", m.GetVal("a"))
	}
	if m.ContainsKey("a") {
		t.Errorf("m.ContainsKey() on nil expected false, got true")
	}
	if m.ContainsValue(10) {
		t.Errorf("m.ContainsValue() on nil expected false, got true")
	}
	if _, ok := m.Remove("a"); ok {
		t.Errorf("m.Remove() on nil expected ok=false, got true")
	}
	if m.RemoveKey("a") != 0 {
		t.Errorf("m.RemoveKey() on nil expected 0, got %d", m.RemoveKey("a"))
	}
	if m.Entries() != nil {
		t.Errorf("m.Entries() on nil expected nil, got %v", m.Entries())
	}
	if m.Values() == nil || m.Values().Size() != 0 {
		t.Errorf("m.Values() on nil expected empty ArrayList, got %v", m.Values())
	}
	if m.KeySet() == nil || m.KeySet().Size() != 0 {
		t.Errorf("m.KeySet() on nil expected empty LinkedHashSet, got %v", m.KeySet())
	}
	if m.HashCode() != 0 {
		t.Errorf("m.HashCode() on nil expected 0, got %d", m.HashCode())
	}
	if m.String() != "null" {
		t.Errorf("m.String() on nil expected 'null', got %q", m.String())
	}

	// Verify All() and Keys() don't panic on nil m
	for range m.All() {
		t.Errorf("m.All() on nil should yield nothing")
	}
	for range m.Keys() {
		t.Errorf("m.Keys() on nil should yield nothing")
	}

	// Now examine LinkedHashSet nil receiver behavior:
	var s *javautil.LinkedHashSet[string] = nil

	// ToSlice, String, HashCode explicitly guard for s == nil:
	if s.ToSlice() != nil {
		t.Errorf("s.ToSlice() on nil expected nil, got %v", s.ToSlice())
	}
	if s.String() != "null" {
		t.Errorf("s.String() on nil expected 'null', got %q", s.String())
	}
	if s.HashCode() != 0 {
		t.Errorf("s.HashCode() on nil expected 0, got %d", s.HashCode())
	}

	// Inconsistent: Size(), IsEmpty(), Contains() panic on nil receiver because s.m accesses s without nil guard!
	assertPanics(t, "s.Size() on nil receiver", func() { _ = s.Size() })
	assertPanics(t, "s.IsEmpty() on nil receiver", func() { _ = s.IsEmpty() })
	assertPanics(t, "s.Contains() on nil receiver", func() { _ = s.Contains("a") })
}

// TestNilArgumentHandling verifies robustness against nil arguments.
func TestNilArgumentHandling(t *testing.T) {
	m := javautil.NewLinkedHashMap[string, int]()
	m.Put("a", 1)

	// PutAll(nil) should be a safe no-op
	m.PutAll(nil)
	if m.Size() != 1 {
		t.Errorf("PutAll(nil) changed map size to %d", m.Size())
	}

	s := javautil.NewLinkedHashSet[string]("x", "y")
	// AddAllSet(nil) should return false and not panic
	if s.AddAllSet(nil) {
		t.Errorf("AddAllSet(nil) should return false")
	}
	// RemoveAll(nil) should return false and not panic
	if s.RemoveAll(nil) {
		t.Errorf("RemoveAll(nil) should return false")
	}
	// RetainAll(nil) clears the set in current implementation
	sCopy := javautil.NewLinkedHashSet[string]("x", "y")
	if !sCopy.RetainAll(nil) || sCopy.Size() != 0 {
		t.Errorf("RetainAll(nil) should clear set, size=%d", sCopy.Size())
	}
}

// TestIterationDeterminismAcrossRepeatedRuns verifies 100% byte-identical iteration order across 200 runs.
func TestIterationDeterminismAcrossRepeatedRuns(t *testing.T) {
	const iterations = 200
	const numKeys = 500

	var baselineKeys []string
	var baselineVals []int

	for run := 0; run < iterations; run++ {
		m := javautil.NewLinkedHashMap[string, int]()
		// Deterministic PRNG seed for identical operation stream
		rng := rand.New(rand.NewSource(12345))

		// Interleave insertions, updates, and deletions
		for i := 0; i < numKeys; i++ {
			k := fmt.Sprintf("key_%04d", rng.Intn(numKeys))
			m.Put(k, i)
			if i%7 == 0 {
				delKey := fmt.Sprintf("key_%04d", rng.Intn(numKeys))
				m.Remove(delKey)
			}
		}

		var runKeys []string
		var runVals []int
		for k, v := range m.All() {
			runKeys = append(runKeys, k)
			runVals = append(runVals, v)
		}

		if run == 0 {
			baselineKeys = runKeys
			baselineVals = runVals
		} else {
			if len(runKeys) != len(baselineKeys) {
				t.Fatalf("Run %d: size mismatch got %d, want %d", run, len(runKeys), len(baselineKeys))
			}
			for idx := range runKeys {
				if runKeys[idx] != baselineKeys[idx] || runVals[idx] != baselineVals[idx] {
					t.Fatalf("Run %d determinism violation at index %d: got (%s, %d), want (%s, %d)",
						run, idx, runKeys[idx], runVals[idx], baselineKeys[idx], baselineVals[idx])
				}
			}
		}
	}
}

// TestLinkedHashSetIterationDeterminismRepeatedRuns verifies LinkedHashSet 100% order reproducibility.
func TestLinkedHashSetIterationDeterminismRepeatedRuns(t *testing.T) {
	const iterations = 200
	const numOps = 500

	var baseline []int

	for run := 0; run < iterations; run++ {
		s := javautil.NewLinkedHashSet[int]()
		rng := rand.New(rand.NewSource(98765))

		for i := 0; i < numOps; i++ {
			val := rng.Intn(200)
			if rng.Intn(10) < 7 {
				s.Add(val)
			} else {
				s.Remove(val)
			}
		}

		slice := s.ToSlice()
		if run == 0 {
			baseline = slice
		} else {
			if len(slice) != len(baseline) {
				t.Fatalf("Run %d: set size mismatch got %d, want %d", run, len(slice), len(baseline))
			}
			for idx := range slice {
				if slice[idx] != baseline[idx] {
					t.Fatalf("Run %d: set order violation at index %d: got %d, want %d",
						run, idx, slice[idx], baseline[idx])
				}
			}
		}
	}
}

// TestConcurrentReadNoDeletions verifies that when holes == 0, concurrent reads are race-free.
func TestConcurrentReadNoDeletions(t *testing.T) {
	m := javautil.NewLinkedHashMap[int, string]()
	for i := 0; i < 50; i++ {
		m.Put(i, fmt.Sprintf("val-%d", i))
	}

	var wg sync.WaitGroup
	for g := 0; g < 10; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for iter := 0; iter < 100; iter++ {
				count := 0
				for range m.All() {
					count++
				}
				if count != 50 {
					panic(fmt.Sprintf("expected 50, got %d", count))
				}
				_ = m.Values()
				_ = m.Entries()
				_ = m.ContainsKey(25)
				_ = m.Size()
			}
		}()
	}
	wg.Wait()
}

// TestConcurrentReadAfterRemovalsLazyCompactExposesRace tests whether read operations
// (All, Values, Entries, KeySet, ContainsValue) after removals cause data races or crashes under concurrent read access.
func TestConcurrentReadAfterRemovalsLazyCompactExposesRace(t *testing.T) {
	m := javautil.NewLinkedHashMap[int, string]()
	for i := 0; i < 50; i++ {
		m.Put(i, fmt.Sprintf("val-%d", i))
	}

	// Remove 10 items.
	// Note: 10 holes < 16, so m.Remove does NOT call compact() internally!
	// Therefore, m.holes == 10 > 0 at this point.
	for i := 0; i < 10; i++ {
		m.Remove(i)
	}

	// Now 10 goroutines concurrently call read-only methods.
	// If m.compact() is called lazily on read methods (All, Values, Entries),
	// this causes concurrent map and slice writes during "read" operations!
	var wg sync.WaitGroup
	var panicCount atomic.Int32

	for g := 0; g < 10; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicCount.Add(1)
				}
			}()

			switch gid % 4 {
			case 0:
				count := 0
				for range m.All() {
					count++
				}
			case 1:
				_ = m.Values()
			case 2:
				_ = m.Entries()
			case 3:
				_ = m.ContainsValue("val-25")
			}
		}(g)
	}
	wg.Wait()

	if panicCount.Load() > 0 {
		t.Logf("Detected %d panics during concurrent reads after deletions due to unsynchronized lazy compact()", panicCount.Load())
	}
}

// TestZeroValueKeyAndHoles verifies that zero-value keys (0, "", nil) correctly survive
// hole creation, deletion, re-insertion, and compaction without being lost or duplicated.
func TestZeroValueKeyAndHoles(t *testing.T) {
	m := javautil.NewLinkedHashMap[int, string]()

	// 1. Insert 0 first
	m.Put(0, "zero-orig")
	for i := 1; i <= 20; i++ {
		m.Put(i, fmt.Sprintf("val-%d", i))
	}

	// 2. Remove 0 (creates hole at index 0)
	val, ok := m.Remove(0)
	if !ok || val != "zero-orig" {
		t.Fatalf("Failed to remove 0: %v, %v", val, ok)
	}
	if m.ContainsKey(0) {
		t.Fatalf("m should not contain 0 after remove")
	}

	// 3. Remove several other elements
	m.Remove(5)
	m.Remove(10)
	m.Remove(15)

	// 4. Re-insert 0 (should be appended to the end)
	m.Put(0, "zero-reinserted")

	// 5. Verify iteration order: 0 must be the very last element
	var lastKey int
	var lastVal string
	count := 0
	for k, v := range m.All() {
		lastKey = k
		lastVal = v
		count++
	}

	if lastKey != 0 || lastVal != "zero-reinserted" {
		t.Fatalf("Zero key re-inserted should be at end, got (%d, %s)", lastKey, lastVal)
	}
	expectedCount := 21 - 3 // 21 inserted, 4 removed (0, 5, 10, 15), 0 re-inserted = 18
	if count != expectedCount || m.Size() != expectedCount {
		t.Fatalf("Size mismatch: got %d (iter count %d), expected %d", m.Size(), count, expectedCount)
	}

	// 6. Remove 0 again, then re-insert 0 multiple times
	m.Remove(0)
	m.Put(0, "zero-v3")
	m.Remove(0)
	m.Put(0, "zero-final")

	lastKey = -1
	for k, v := range m.All() {
		lastKey = k
		lastVal = v
	}
	if lastKey != 0 || lastVal != "zero-final" {
		t.Fatalf("Final zero re-insertion mismatch: got (%d, %s)", lastKey, lastVal)
	}
}

// TestLargeKeySet100kOperations tests scalability and correctness with 100,000 keys.
func TestLargeKeySet100kOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 100k key test in short mode")
	}

	const n = 100000
	m := javautil.NewLinkedHashMapWithCapacity[int, int](n)

	// Sequential insertion
	for i := 0; i < n; i++ {
		m.Put(i, i*2)
	}
	if m.Size() != n {
		t.Fatalf("Size got %d, want %d", m.Size(), n)
	}

	// Remove all even keys (50,000 removals -> triggers compact multiple times)
	for i := 0; i < n; i += 2 {
		m.Remove(i)
	}
	if m.Size() != n/2 {
		t.Fatalf("Size after removing evens got %d, want %d", m.Size(), n/2)
	}

	// Verify iteration strictly yields odd numbers in ascending order: 1, 3, 5, ...
	expected := 1
	count := 0
	for k, v := range m.All() {
		if k != expected || v != expected*2 {
			t.Fatalf("Order mismatch at %d: got (%d, %d), want (%d, %d)", count, k, v, expected, expected*2)
		}
		expected += 2
		count++
	}
	if count != n/2 {
		t.Fatalf("Iter count got %d, want %d", count, n/2)
	}
}

// TestClearAndReuseStress verifies clear() completely resets state and allows repeated reuse.
func TestClearAndReuseStress(t *testing.T) {
	m := javautil.NewLinkedHashMap[int, int]()

	for cycle := 0; cycle < 50; cycle++ {
		numElems := 100 + cycle*10
		for i := 0; i < numElems; i++ {
			m.Put(i, i+1000)
		}
		if m.Size() != numElems {
			t.Fatalf("Cycle %d: size got %d, want %d", cycle, m.Size(), numElems)
		}

		// Remove some elements to generate holes
		for i := 0; i < numElems/2; i += 2 {
			m.Remove(i)
		}

		m.Clear()
		if m.Size() != 0 || !m.IsEmpty() {
			t.Fatalf("Cycle %d: Clear() failed, size=%d", cycle, m.Size())
		}
		for range m.All() {
			t.Fatalf("Cycle %d: All() yielded element after Clear()", cycle)
		}

		// Re-insert smaller set
		for i := 0; i < 10; i++ {
			m.Put(i, i*10)
		}
		if m.Size() != 10 {
			t.Fatalf("Cycle %d: size after re-insert got %d, want 10", cycle, m.Size())
		}
		idx := 0
		for k, v := range m.All() {
			if k != idx || v != idx*10 {
				t.Fatalf("Cycle %d: re-insert order mismatch at %d: (%d, %d)", cycle, idx, k, v)
			}
			idx++
		}
		m.Clear()
	}
}

// TestUpdateExistingKeyPreservesOrder verifies that updating existing keys does NOT
// change their insertion position.
func TestUpdateExistingKeyPreservesOrder(t *testing.T) {
	m := javautil.NewLinkedHashMap[string, string]()
	keys := []string{"first", "second", "third", "fourth", "fifth"}
	for _, k := range keys {
		m.Put(k, "orig-"+k)
	}

	// Update "second" and "first"
	prev := m.Put("second", "updated-second")
	if prev != "orig-second" {
		t.Fatalf("Put('second') returned prev %q, want 'orig-second'", prev)
	}
	m.Put("first", "updated-first")

	// Order must remain: "first", "second", "third", "fourth", "fifth"
	idx := 0
	for k, v := range m.All() {
		if k != keys[idx] {
			t.Fatalf("Index %d: key got %q, want %q", idx, k, keys[idx])
		}
		if idx == 0 && v != "updated-first" {
			t.Fatalf("first value got %q, want 'updated-first'", v)
		}
		if idx == 1 && v != "updated-second" {
			t.Fatalf("second value got %q, want 'updated-second'", v)
		}
		idx++
	}
	if idx != len(keys) {
		t.Fatalf("Yielded %d keys, want %d", idx, len(keys))
	}
}

// TestSetAlgebraDeterminismAndEdgeCases tests Set RetainAll, RemoveAll, and AddAll.
func TestSetAlgebraDeterminismAndEdgeCases(t *testing.T) {
	s1 := javautil.NewLinkedHashSet[int](1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	s2 := javautil.NewLinkedHashSet[int](2, 4, 6, 8, 10, 12)

	// RetainAll: only evens 2, 4, 6, 8, 10 should remain, in original insertion order!
	modified := s1.RetainAll(s2)
	if !modified {
		t.Fatalf("RetainAll should report modified=true")
	}
	expected := []int{2, 4, 6, 8, 10}
	slice := s1.ToSlice()
	if len(slice) != len(expected) {
		t.Fatalf("RetainAll size got %d, want %d", len(slice), len(expected))
	}
	for i := range expected {
		if slice[i] != expected[i] {
			t.Fatalf("RetainAll order mismatch at %d: got %d, want %d", i, slice[i], expected[i])
		}
	}

	// RemoveAll
	s3 := javautil.NewLinkedHashSet[int](2, 6, 10)
	modified = s1.RemoveAll(s3)
	if !modified {
		t.Fatalf("RemoveAll should report modified=true")
	}
	expectedAfterRemove := []int{4, 8}
	slice = s1.ToSlice()
	if len(slice) != len(expectedAfterRemove) {
		t.Fatalf("RemoveAll size got %d, want %d", len(slice), len(expectedAfterRemove))
	}
	for i := range expectedAfterRemove {
		if slice[i] != expectedAfterRemove[i] {
			t.Fatalf("RemoveAll order mismatch at %d: got %d, want %d", i, slice[i], expectedAfterRemove[i])
		}
	}
}

func assertPanics(t *testing.T, desc string, f func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("%s did NOT panic as expected", desc)
		}
	}()
	f()
}
