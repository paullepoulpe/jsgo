package javalang

import (
	"math"
	"testing"
)

func TestClassDescriptor(t *testing.T) {
	c1 := NewClass("java.lang.String", "String")
	c2 := NewClass("java.lang.String", "String")
	c3 := NewClass("java.lang.Integer", "Integer")

	if c1.GetName() != "java.lang.String" {
		t.Errorf("GetName got %q, want java.lang.String", c1.GetName())
	}
	if c1.GetSimpleName() != "String" {
		t.Errorf("GetSimpleName got %q, want String", c1.GetSimpleName())
	}
	if c1.String() != "class java.lang.String" {
		t.Errorf("String got %q, want class java.lang.String", c1.String())
	}
	if !c1.Equals(c2) {
		t.Errorf("c1.Equals(c2) should be true")
	}
	if c1.Equals(c3) {
		t.Errorf("c1.Equals(c3) should be false")
	}
	if c1.HashCode() != StringHashCode("java.lang.String") {
		t.Errorf("c1.HashCode() mismatch")
	}
}

func TestThrowableHierarchy(t *testing.T) {
	npe := NewNullPointerExceptionMsg("foo was null")
	if npe.GetMessage() != "foo was null" {
		t.Errorf("npe.GetMessage() got %q", npe.GetMessage())
	}
	if npe.Error() != "java.lang.NullPointerException: foo was null" {
		t.Errorf("npe.Error() got %q", npe.Error())
	}

	// Test marker interface promotion
	var _ IRuntimeException = npe
	var _ IException = npe
	var _ Throwable = npe

	// Test panic & recovery
	var caught Throwable
	func() {
		defer func() {
			caught = RecoverException(recover())
		}()
		panic(NewIllegalArgumentException("bad arg"))
	}()
	if caught == nil {
		t.Fatalf("failed to recover exception")
	}
	if caught.GetMessage() != "bad arg" {
		t.Errorf("recovered message got %q, want \"bad arg\"", caught.GetMessage())
	}
	if _, ok := caught.(IIllegalArgumentException); !ok {
		t.Errorf("recovered exception should implement IIllegalArgumentException")
	}
}

func TestPrimitivesConversions(t *testing.T) {
	// Narrowing casts
	if D2I(123.45) != 123 {
		t.Errorf("D2I(123.45) got %d, want 123", D2I(123.45))
	}
	if D2I(1e12) != 2147483647 {
		t.Errorf("D2I(1e12) saturation failed: got %d", D2I(1e12))
	}
	if D2I(-1e12) != -2147483648 {
		t.Errorf("D2I(-1e12) saturation failed: got %d", D2I(-1e12))
	}

	// Shifts with masking
	if Shl32(1, 33) != 2 {
		t.Errorf("Shl32(1, 33) got %d, want 2", Shl32(1, 33))
	}
	if Ushr32(-1, 1) != 2147483647 {
		t.Errorf("Ushr32(-1, 1) got %d, want 2147483647", Ushr32(-1, 1))
	}

	// DoubleToString
	if DoubleToString(0.0) != "0.0" {
		t.Errorf("DoubleToString(0.0) got %q, want \"0.0\"", DoubleToString(0.0))
	}
	if DoubleToString(12345.67) != "12345.67" {
		t.Errorf("DoubleToString(12345.67) got %q, want \"12345.67\"", DoubleToString(12345.67))
	}

	// System arraycopy
	src := []int{1, 2, 3, 4, 5}
	dest := make([]int, 5)
	System_Arraycopy(src, 1, dest, 2, 3)
	if dest[2] != 2 || dest[3] != 3 || dest[4] != 4 {
		t.Errorf("System_Arraycopy failed: got %v", dest)
	}
}

func TestFloatToString(t *testing.T) {
	cases := []struct {
		in   float32
		want string
	}{
		{0.0, "0.0"},
		{float32(math.Copysign(0, -1)), "-0.0"},
		{float32(math.NaN()), "NaN"},
		{float32(math.Inf(1)), "Infinity"},
		{float32(math.Inf(-1)), "-Infinity"},
		{1e-4, "1.0E-4"},
		{-1e-4, "-1.0E-4"},
		{1e-3, "0.001"},
		{-1e-3, "-0.001"},
		{0.1, "0.1"},
		{1.0, "1.0"},
		{123.0, "123.0"},
		{-123.45, "-123.45"},
		{9999999.0, "9999999.0"},
		{1e7, "1.0E7"},
		{-1e7, "-1.0E7"},
		{1e20, "1.0E20"},
		{math.SmallestNonzeroFloat32, "1.4E-45"},
		{math.MaxFloat32, "3.4028235E38"},
	}
	for _, tc := range cases {
		got := FloatToString(tc.in)
		if got != tc.want {
			t.Errorf("FloatToString(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
