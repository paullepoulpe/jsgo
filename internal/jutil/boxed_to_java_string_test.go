package jutil

import "testing"

// Go-only: test — ToJavaString of nullable boxed primitives (pointers).
func TestToJavaStringBoxed(t *testing.T) {
	i, b, d, c := int32(3), true, 0.5, uint16('x')
	var nilInt *int32
	for _, tc := range []struct {
		v    any
		want string
	}{
		{&i, "3"}, {&b, "true"}, {&d, "0.5"}, {&c, "x"}, {nilInt, "null"},
	} {
		if got := ToJavaString(tc.v); got != tc.want {
			t.Errorf("ToJavaString(%T) = %q, want %q", tc.v, got, tc.want)
		}
	}
	if got := LenientFormat("Invalid Polymer version: (%s)", &i); got != "Invalid Polymer version: (3)" {
		t.Errorf("LenientFormat = %q", got)
	}
}
