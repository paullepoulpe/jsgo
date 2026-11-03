package jutil

import "testing"

// A nullable Java String is a *string: Objects.equals and Objects.hashCode
// compare and hash the string, not the pointer.
func TestObjectsNullableString(t *testing.T) {
	a, b, c := "x", "x", "y"
	if !ObjectsEquals(&a, &b) || ObjectsEquals(&a, &c) || !ObjectsEquals(&a, "x") || !ObjectsEquals("x", &a) || ObjectsEquals(&a, "y") || ObjectsEquals("y", &a) {
		t.Errorf("ObjectsEquals(*string, string)")
	}
	if !ObjectsEquals("abc", StringPtr("abc")) || !ObjectsEquals(StringPtr("abc"), "abc") {
		t.Errorf("ObjectsEquals(StringPtr)")
	}
	var n *string
	if !ObjectsEquals(n, (*string)(nil)) || ObjectsEquals(&a, n) || ObjectsEquals(n, &a) || ObjectsEquals("", n) || ObjectsEquals(n, "") {
		t.Errorf("ObjectsEquals(nil *string)")
	}
	if ObjectHashCode(&a) != StringHashCode("x") || ObjectHashCode(&a) != ObjectHashCode("x") || ObjectHashCode(n) != 0 {
		t.Errorf("ObjectHashCode(*string)")
	}
}
