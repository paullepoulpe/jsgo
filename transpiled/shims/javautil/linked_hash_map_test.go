package javautil

import (
	"fmt"
	"testing"
)

func TestLinkedHashMapDeterminism(t *testing.T) {
	m := NewLinkedHashMap[string, int]()
	// Insert > 8 elements to trigger index map creation
	const n = 20
	for i := 0; i < n; i++ {
		m.Put(fmt.Sprintf("key_%02d", i), i*10)
	}

	if m.Size() != n {
		t.Errorf("m.Size() got %d, want %d", m.Size(), n)
	}

	// Verify insertion order is strictly preserved
	idx := 0
	for k, v := range m.All() {
		expectedKey := fmt.Sprintf("key_%02d", idx)
		expectedVal := idx * 10
		if k != expectedKey || v != expectedVal {
			t.Fatalf("entry %d: got (%q, %d), want (%q, %d)", idx, k, v, expectedKey, expectedVal)
		}
		idx++
	}

	// Remove an entry and verify holes and order
	remVal, ok := m.Remove("key_05")
	if !ok || remVal != 50 || m.Size() != n-1 {
		t.Fatalf("failed to remove key_05")
	}
	if m.ContainsKey("key_05") {
		t.Errorf("key_05 should no longer be present")
	}

	// Verify remaining order without key_05
	idx = 0
	for k, v := range m.All() {
		if idx == 5 {
			idx++ // skip 5
		}
		expectedKey := fmt.Sprintf("key_%02d", idx)
		expectedVal := idx * 10
		if k != expectedKey || v != expectedVal {
			t.Fatalf("entry %d after remove: got (%q, %d), want (%q, %d)", idx, k, v, expectedKey, expectedVal)
		}
		idx++
	}

	// Re-insert key_05 -> should append at the end (insertion order!)
	m.Put("key_05", 555)
	var lastKey string
	var lastVal int
	for k, v := range m.All() {
		lastKey = k
		lastVal = v
	}
	if lastKey != "key_05" || lastVal != 555 {
		t.Errorf("re-inserted key should be at the end, got (%q, %d)", lastKey, lastVal)
	}
}
