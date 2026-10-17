package jutil

import "testing"

// A nullable Java String is a *string: Objects.equals and Objects.hashCode
// compare and hash the string, not the pointer.
func TestObjectsNullableString(t *testing.T) {
	a, b, c := "x", "x", "y"
	if !ObjectsEquals(&a, &b) || ObjectsEquals(&a, &c) || ObjectsEquals(&a, "x") {
		t.Errorf("ObjectsEquals(*string)")
	}
	var n *string
	if !ObjectsEquals(n, (*string)(nil)) || ObjectsEquals(&a, n) || ObjectsEquals(n, &a) {
		t.Errorf("ObjectsEquals(nil *string)")
	}
	if ObjectHashCode(&a) != StringHashCode("x") || ObjectHashCode(n) != 0 {
		t.Errorf("ObjectHashCode(*string)")
	}
}
