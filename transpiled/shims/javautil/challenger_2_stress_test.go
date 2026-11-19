package javautil_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
)

// Reference LinkedHashMap model
type refLinkedHashMap[K comparable, V comparable] struct {
	order []K
	m     map[K]V
}

func newRefLinkedHashMap[K comparable, V comparable]() *refLinkedHashMap[K, V] {
	return &refLinkedHashMap[K, V]{
		order: make([]K, 0),
		m:     make(map[K]V),
	}
}

func (r *refLinkedHashMap[K, V]) put(k K, v V) {
	if _, exists := r.m[k]; !exists {
		r.order = append(r.order, k)
	}
	r.m[k] = v
}

func (r *refLinkedHashMap[K, V]) remove(k K) {
	if _, exists := r.m[k]; exists {
		delete(r.m, k)
		newOrder := make([]K, 0, len(r.order)-1)
		for _, x := range r.order {
			if x != k {
				newOrder = append(newOrder, x)
			}
		}
		r.order = newOrder
	}
}

func TestLinkedHashMap10kRandomizedOperations(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	m := javautil.NewLinkedHashMap[int, string]()
	ref := newRefLinkedHashMap[int, string]()

	keyDomain := 200 // ensures collisions, overwrites, removals, and re-insertions

	for op := 0; op < 10000; op++ {
		action := rng.Intn(100)
		key := rng.Intn(keyDomain)

		switch {
		case action < 55: // 55% Put
			val := fmt.Sprintf("val-%d-%d", key, op)
			m.Put(key, val)
			ref.put(key, val)
		case action < 80: // 25% Remove
			m.Remove(key)
			ref.remove(key)
		case action < 90: // 10% Get / ContainsKey check
			v, ok := m.Get(key)
			refV, refOk := ref.m[key]
			if ok != refOk {
				t.Fatalf("op %d: Get(%d) ok=%v, expected %v", op, key, ok, refOk)
			}
			if ok && v != refV {
				t.Fatalf("op %d: Get(%d) val=%s, expected %s", op, key, v, refV)
			}
		default: // 10% Full iteration and order check
			if m.Size() != len(ref.order) {
				t.Fatalf("op %d: Size() = %d, expected %d", op, m.Size(), len(ref.order))
			}
			idx := 0
			for k, v := range m.All() {
				if idx >= len(ref.order) {
					t.Fatalf("op %d: Extra element at index %d: key=%d", op, idx, k)
				}
				expectedKey := ref.order[idx]
				expectedVal := ref.m[expectedKey]
				if k != expectedKey || v != expectedVal {
					t.Fatalf("op %d: Order mismatch at index %d: got (%d, %s), expected (%d, %s)",
						op, idx, k, v, expectedKey, expectedVal)
				}
				idx++
			}
			if idx != len(ref.order) {
				t.Fatalf("op %d: Iteration yielded %d elements, expected %d", op, idx, len(ref.order))
			}
		}
	}

	// Final verification of 100% order
	if m.Size() != len(ref.order) {
		t.Fatalf("Final size mismatch: %d vs %d", m.Size(), len(ref.order))
	}
	idx := 0
	for k, v := range m.All() {
		expectedKey := ref.order[idx]
		expectedVal := ref.m[expectedKey]
		if k != expectedKey || v != expectedVal {
			t.Fatalf("Final order mismatch at %d: got (%d, %s), expected (%d, %s)",
				idx, k, v, expectedKey, expectedVal)
		}
		idx++
	}
}

func TestLinkedHashMapZeroKeyEdgeCases(t *testing.T) {
	m := javautil.NewLinkedHashMap[string, int]()
	// Insert zero key ""
	m.Put("", 0)
	for i := 1; i <= 15; i++ {
		m.Put(fmt.Sprintf("k%d", i), i)
	}
	// Remove zero key
	val, ok := m.Remove("")
	if !ok || val != 0 {
		t.Fatalf("Remove zero key got (%d, %v), expected (0, true)", val, ok)
	}
	if m.ContainsKey("") {
		t.Fatalf("ContainsKey('') is true after removal")
	}
	// Re-insert zero key
	m.Put("", 999)
	if m.GetVal("") != 999 {
		t.Fatalf("GetVal('') = %d, expected 999", m.GetVal(""))
	}
	// Iteration check - "" must be at the very end
	keys := make([]string, 0, m.Size())
	for k := range m.Keys() {
		keys = append(keys, k)
	}
	if len(keys) != 16 || keys[len(keys)-1] != "" {
		t.Fatalf("Expected '' at end of keys, got %v", keys)
	}

	// Clear and re-use
	m.Clear()
	if m.Size() != 0 || !m.IsEmpty() {
		t.Fatalf("Size after Clear() = %d, expected 0", m.Size())
	}
	m.Put("fresh", 42)
	if m.Size() != 1 || m.GetVal("fresh") != 42 {
		t.Fatalf("Re-use after Clear() failed: %v", m)
	}
}

func TestLinkedHashSet10kRandomizedOperations(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	s := javautil.NewLinkedHashSet[int]()
	var refOrder []int
	refSet := make(map[int]bool)

	keyDomain := 200

	for op := 0; op < 10000; op++ {
		action := rng.Intn(100)
		elem := rng.Intn(keyDomain)

		switch {
		case action < 55: // Add
			added := s.Add(elem)
			wasPresent := refSet[elem]
			if added == wasPresent {
				t.Fatalf("op %d: Add(%d) returned %v, wasPresent=%v", op, elem, added, wasPresent)
			}
			if !wasPresent {
				refSet[elem] = true
				refOrder = append(refOrder, elem)
			}
		case action < 80: // Remove
			removed := s.Remove(elem)
			wasPresent := refSet[elem]
			if removed != wasPresent {
				t.Fatalf("op %d: Remove(%d) returned %v, wasPresent=%v", op, elem, removed, wasPresent)
			}
			if wasPresent {
				delete(refSet, elem)
				newOrder := make([]int, 0, len(refOrder)-1)
				for _, x := range refOrder {
					if x != elem {
						newOrder = append(newOrder, x)
					}
				}
				refOrder = newOrder
			}
		default: // Iteration verification
			if s.Size() != len(refOrder) {
				t.Fatalf("op %d: Size() = %d, expected %d", op, s.Size(), len(refOrder))
			}
			idx := 0
			for val := range s.All() {
				if idx >= len(refOrder) {
					t.Fatalf("op %d: Extra element in set at %d: %d", op, idx, val)
				}
				if val != refOrder[idx] {
					t.Fatalf("op %d: Order mismatch in set at %d: got %d, expected %d", op, idx, val, refOrder[idx])
				}
				idx++
			}
			if idx != len(refOrder) {
				t.Fatalf("op %d: Set iteration yielded %d, expected %d", op, idx, len(refOrder))
			}
		}
	}
}

func TestArrayListEdgeCases(t *testing.T) {
	list := javautil.NewArrayList[int]()
	list.Add(10)
	list.Add(20)
	list.Add(30)

	// Negative index panics
	assertPanicsIndexOutOfBounds(t, "Get(-1)", func() { list.Get(-1) })
	assertPanicsIndexOutOfBounds(t, "Set(-1)", func() { list.Set(-1, 99) })
	assertPanicsIndexOutOfBounds(t, "Remove(-1)", func() { list.Remove(-1) })
	assertPanicsIndexOutOfBounds(t, "AddAt(-1)", func() { list.AddAt(-1, 99) })

	// Out-of-bounds panics
	assertPanicsIndexOutOfBounds(t, "Get(3)", func() { list.Get(3) })
	assertPanicsIndexOutOfBounds(t, "Set(3)", func() { list.Set(3, 99) })
	assertPanicsIndexOutOfBounds(t, "Remove(3)", func() { list.Remove(3) })
	assertPanicsIndexOutOfBounds(t, "AddAt(4)", func() { list.AddAt(4, 99) })

	// AddAt(3) at end is allowed
	list.AddAt(3, 40)
	if list.Size() != 4 || list.Get(3) != 40 {
		t.Fatalf("AddAt(3) failed: %v", list)
	}

	// SubList out-of-bounds
	assertPanicsIndexOutOfBounds(t, "SubList(-1, 2)", func() { list.SubList(-1, 2) })
	assertPanicsIndexOutOfBounds(t, "SubList(2, 5)", func() { list.SubList(2, 5) })
	assertPanicsIndexOutOfBounds(t, "SubList(3, 2)", func() { list.SubList(3, 2) })

	// SubList valid operations
	sub := list.SubList(1, 3) // [20, 30]
	if sub.Size() != 2 || sub.Get(0) != 20 || sub.Get(1) != 30 {
		t.Fatalf("SubList(1, 3) got %v", sub)
	}
	subEmpty := list.SubList(2, 2)
	if subEmpty.Size() != 0 {
		t.Fatalf("SubList(2, 2) expected empty, got size %d", subEmpty.Size())
	}

	// Concurrent read-only access
	var wg sync.WaitGroup
	for g := 0; g < 10; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for iter := 0; iter < 1000; iter++ {
				_ = list.Size()
				_ = list.Get(0)
				_ = list.Get(1)
				_ = list.Contains(20)
				_ = list.IndexOf(30)
				_ = list.ToArray()
				sum := 0
				for v := range list.All() {
					sum += v
				}
				if sum != 100 {
					panic(fmt.Sprintf("Concurrent read got unexpected sum %d", sum))
				}
			}
		}()
	}
	wg.Wait()
}

func assertPanicsIndexOutOfBounds(t *testing.T, desc string, f func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("%s did not panic", desc)
		}
	}()
	f()
}
