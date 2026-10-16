package jutil

import "testing"

func TestLinkedHashMapContainsValue(t *testing.T) {
	m := NewLinkedHashMap[string, string]()
	if m.ContainsValue("a") {
		t.Fatal("empty map contains a value")
	}
	m.Put("1", "a")
	m.Put("2", "b")
	if !m.ContainsValue("a") || !m.ContainsValue("b") || m.ContainsValue("c") || m.ContainsValue("1") {
		t.Fatal("containsValue mismatch")
	}
	m.Remove("1")
	if m.ContainsValue("a") {
		t.Fatal("removed value still contained")
	}
	p := NewLinkedHashMap[string, *string]()
	p.Put("x", nil)
	if !p.ContainsValue(nil) {
		t.Fatal("null value not found")
	}
}
