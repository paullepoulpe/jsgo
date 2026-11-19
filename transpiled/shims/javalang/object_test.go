package javalang

import (
	"math"
	"testing"
)

type customVal struct {
	x int
}

func (c *customVal) Equals(other any) bool {
	o, ok := other.(*customVal)
	return ok && o != nil && c.x == o.x
}

func (c *customVal) HashCode() int32 {
	return int32(c.x * 31)
}

func TestObjectsEqualsAndHash(t *testing.T) {
	if !ObjectsEquals(nil, nil) {
		t.Errorf("ObjectsEquals(nil, nil) should be true")
	}
	var nilPtr *string
	if !ObjectsEquals(nilPtr, nil) {
		t.Errorf("ObjectsEquals(typed nil, nil) should be true")
	}
	if !ObjectsEquals("a", "a") {
		t.Errorf("ObjectsEquals(\"a\", \"a\") should be true")
	}
	if ObjectsEquals("a", "b") {
		t.Errorf("ObjectsEquals(\"a\", \"b\") should be false")
	}

	// IEEE 754 bit-pattern equality (Double.equals in Java: NaN == NaN, +0.0 != -0.0)
	nan1 := math.NaN()
	nan2 := math.NaN()
	if !ObjectsEquals(nan1, nan2) {
		t.Errorf("Double.equals(NaN, NaN) should be true")
	}
	negZero := math.Copysign(0.0, -1.0)
	if ObjectsEquals(0.0, negZero) {
		t.Errorf("Double.equals(+0.0, -0.0) should be false")
	}

	// Custom Equaler
	v1 := &customVal{x: 10}
	v2 := &customVal{x: 10}
	v3 := &customVal{x: 20}
	if !ObjectsEquals(v1, v2) {
		t.Errorf("v1 and v2 should be equal via custom Equals")
	}
	if ObjectsEquals(v1, v3) {
		t.Errorf("v1 and v3 should not be equal")
	}
	if ObjectHashCode(v1) != 310 {
		t.Errorf("ObjectHashCode(v1) got %d, want 310", ObjectHashCode(v1))
	}
}

func TestIsNil(t *testing.T) {
	if !IsNil(nil) {
		t.Errorf("IsNil(nil) should be true")
	}
	var ptr *customVal
	if !IsNil(ptr) {
		t.Errorf("IsNil(typed nil pointer) should be true")
	}
	var slice []int
	if !IsNil(slice) {
		t.Errorf("IsNil(nil slice) should be true")
	}
	if IsNil("hello") {
		t.Errorf("IsNil(\"hello\") should be false")
	}
	if IsNil(0) {
		t.Errorf("IsNil(0) should be false")
	}
}

type mockJavaObj struct {
	id int
}

func (m *mockJavaObj) IsJavaNil_() bool { return m == nil }

type uncomparableFieldObj struct {
	items []int
	extra any
}

func (u *uncomparableFieldObj) IsJavaNil_() bool { return u == nil }

func TestJavaNilCheckerAndFastPaths(t *testing.T) {
	var nilObj *mockJavaObj
	nonNilObj := &mockJavaObj{id: 1}
	otherObj := &mockJavaObj{id: 1}

	if !IsNil(nilObj) {
		t.Errorf("IsNil(nil JavaNilChecker) should be true")
	}
	if IsNil(nonNilObj) {
		t.Errorf("IsNil(non-nil JavaNilChecker) should be false")
	}

	// Default Java Object.equals on non-Equaler structs is pointer identity without heap allocs.
	if !ObjectsEquals(nonNilObj, nonNilObj) {
		t.Errorf("ObjectsEquals(same pointer) should be true")
	}
	if ObjectsEquals(nonNilObj, otherObj) {
		t.Errorf("ObjectsEquals(distinct pointers without Equals) should be false")
	}
	var nilCustom *customVal
	if !ObjectsEquals(nilObj, nilCustom) {
		t.Errorf("ObjectsEquals(different typed nils) should be true")
	}
	if ObjectsEquals(nilCustom, &customVal{x: 1}) {
		t.Errorf("ObjectsEquals(typed nil Equaler, non-nil) should be false")
	}

	u1 := &uncomparableFieldObj{items: []int{1, 2}, extra: []string{"a"}}
	u2 := &uncomparableFieldObj{items: []int{1, 2}, extra: []string{"a"}}
	if !ObjectsEquals(u1, u1) {
		t.Errorf("ObjectsEquals(u1, u1) should be true")
	}
	if ObjectsEquals(u1, u2) {
		t.Errorf("ObjectsEquals(u1, u2) should be false")
	}
	// Value struct containing an interface holding a slice must not panic in ObjectsEquals.
	vs1 := struct{ x any }{x: []int{1, 2}}
	vs2 := struct{ x any }{x: []int{1, 2}}
	if !ObjectsEquals(vs1, vs2) {
		t.Errorf("ObjectsEquals(value structs with slice in any) should fall back to DeepEqual")
	}

	allocs := testing.AllocsPerRun(100, func() {
		_ = IsNil(any(nonNilObj))
		_ = IsNil(any(nilObj))
		_ = ObjectsEquals(any(u1), any(u2))
		_ = ObjectsEquals(any(u1), any(u1))
	})
	if allocs != 0 {
		t.Errorf("expected 0 allocs/op for IsNil and ObjectsEquals on JavaNilChecker pointers, got %v", allocs)
	}
}

func TestToJavaStringPrimitives(t *testing.T) {
	cases := []struct {
		in   any
		want string
	}{
		{nil, "null"},
		{(*mockJavaObj)(nil), "null"},
		{"abc", "abc"},
		{float64(1e100), "1.0E100"},
		{float64(0.0), "0.0"},
		{float64(math.Copysign(0, -1)), "-0.0"},
		{float32(1e20), "1.0E20"},
		{float32(0.1), "0.1"},
		{int32(-42), "-42"},
		{int64(9007199254740993), "9007199254740993"},
		{int(123), "123"},
		{int8(-8), "-8"},
		{int16(-16), "-16"},
		{uint8(255), "255"},
		{uint16('Z'), "Z"},
		{uint32(4000000000), "4000000000"},
		{uint64(18446744073709551615), "18446744073709551615"},
		{true, "true"},
		{false, "false"},
	}
	for _, tc := range cases {
		if got := ToJavaString(tc.in); got != tc.want {
			t.Errorf("ToJavaString(%T(%v)) = %q, want %q", tc.in, tc.in, got, tc.want)
		}
	}
}

type mockIdentityObj struct {
	hash int32
}

func (m *mockIdentityObj) IsJavaNil_() bool    { return m == nil }
func (m *mockIdentityObj) HasIdentityEquals_() {}
func (m *mockIdentityObj) IdentityHash_() int32 {
	return m.hash
}

type mockEqualerObj struct {
	val    int32
	idHash int32
}

func (m *mockEqualerObj) IsJavaNil_() bool     { return m == nil }
func (m *mockEqualerObj) IdentityHash_() int32 { return m.idHash }
func (m *mockEqualerObj) HashCode() int32      { return m.val * 17 }
func (m *mockEqualerObj) Equals(other any) bool {
	o, ok := other.(*mockEqualerObj)
	return ok && o != nil && o.val == m.val
}

type mockDualPredicate uint32

func (mockDualPredicate) IsJavaNil_() bool    { return false }
func (mockDualPredicate) HasIdentityEquals_() {}
func (p mockDualPredicate) Equals(other any) bool {
	switch o := other.(type) {
	case mockDualPredicate:
		return o == p
	case *mockEqualerObj:
		return o != nil && uint32(o.val) == uint32(p)
	default:
		return false
	}
}

func TestWave20AObjectsEqualsAndIdentityHashFastPaths(t *testing.T) {
	id1 := &mockIdentityObj{hash: 0x12345678}
	id2 := &mockIdentityObj{hash: 0x12345678}
	var nilID *mockIdentityObj

	eq1 := &mockEqualerObj{val: 42, idHash: 0x1111}
	eq2 := &mockEqualerObj{val: 42, idHash: 0x2222}
	eq3 := &mockEqualerObj{val: 99, idHash: 0x3333}
	var nilEq *mockEqualerObj
	var nilCustom *customVal

	// IdentityHasher & HashCoder fast paths (including typed nils).
	if got := IdentityHashCode(id1); got != 0x12345678 {
		t.Fatalf("IdentityHashCode(id1) = 0x%x, want 0x12345678", got)
	}
	if got := ObjectHashCode(id1); got != 0x12345678 {
		t.Fatalf("ObjectHashCode(id1) = 0x%x, want 0x12345678", got)
	}
	if IdentityHashCode(nilID) != 0 || ObjectHashCode(nilID) != 0 {
		t.Fatalf("expected 0 hash for typed-nil IdentityHasher")
	}
	// System.identityHashCode on Equaler/HashCoder uses IdentityHash_, whereas ObjectHashCode uses HashCode().
	if got := IdentityHashCode(eq1); got != 0x1111 {
		t.Fatalf("IdentityHashCode(eq1) = 0x%x, want 0x1111", got)
	}
	if got := ObjectHashCode(eq1); got != 42*17 {
		t.Fatalf("ObjectHashCode(eq1) = %d, want %d", got, 42*17)
	}
	if IdentityHashCode(nilEq) != 0 || ObjectHashCode(nilEq) != 0 || ObjectHashCode(nilCustom) != 0 {
		t.Fatalf("expected 0 hash for typed-nil Equaler/HashCoder")
	}

	// IdentityEqualer pointer identity & cross-type/nil equality.
	if !ObjectsEquals(id1, id1) || ObjectsEquals(id1, id2) {
		t.Fatalf("IdentityEqualer pointer identity failed")
	}
	if ObjectsEquals(id1, nil) || ObjectsEquals(nil, id1) || ObjectsEquals(id1, nilID) || ObjectsEquals(nilID, id1) {
		t.Fatalf("non-nil IdentityEqualer matched nil")
	}
	if !ObjectsEquals(nilID, nil) || !ObjectsEquals(nil, nilID) || !ObjectsEquals(nilID, nilEq) || !ObjectsEquals(nilEq, nilID) {
		t.Fatalf("typed-nil IdentityEqualer/Equaler should equal nil")
	}

	// Equaler fast path & priority over HasIdentityEquals_ on dual types (e.g. predEqualNode_).
	if !ObjectsEquals(eq1, eq2) || ObjectsEquals(eq1, eq3) {
		t.Fatalf("Equaler equality mismatch")
	}
	if !ObjectsEquals(nilCustom, nilID) || ObjectsEquals(nilCustom, id1) {
		t.Fatalf("non-JavaNilChecker typed-nil Equaler mismatch")
	}
	dual := mockDualPredicate(42)
	if !ObjectsEquals(dual, eq1) || ObjectsEquals(dual, eq3) {
		t.Fatalf("Equaler priority over HasIdentityEquals_ failed")
	}
	// IdentityEqualer compared against uncomparable slice in any must return false without panicking.
	if ObjectsEquals(id1, []int{1, 2, 3}) {
		t.Fatalf("IdentityEqualer should not equal []int")
	}

	// NullString symmetry across ObjectsEquals, IsNil, ObjectHashCode, IdentityHashCode, ToJavaString.
	if !IsNil(NullString) || ToJavaString(NullString) != "null" {
		t.Fatalf("NullString IsNil/ToJavaString mismatch")
	}
	if ObjectHashCode(NullString) != 0 || IdentityHashCode(NullString) != 0 {
		t.Fatalf("NullString hash should be 0")
	}
	nullCases := []any{nil, NullString, nilID, nilEq, nilCustom, (*string)(nil), ([]int)(nil)}
	for i, a := range nullCases {
		for j, b := range nullCases {
			if !ObjectsEquals(a, b) {
				t.Fatalf("ObjectsEquals(nullCases[%d]=%T, nullCases[%d]=%T) = false, want true", i, a, j, b)
			}
		}
		if ObjectsEquals(a, "") || ObjectsEquals("", a) || ObjectsEquals(a, eq1) || ObjectsEquals(eq1, a) || ObjectsEquals(a, id1) || ObjectsEquals(id1, a) {
			t.Fatalf("nullCases[%d]=%T matched non-null value", i, a)
		}
	}

	// Zero-alloc verification on hot fast paths.
	allocs := testing.AllocsPerRun(100, func() {
		_ = ObjectsEquals(any(id1), any(id2))
		_ = ObjectsEquals(any(eq1), any(eq2))
		_ = ObjectsEquals(any("abc"), any("abd"))
		_ = IdentityHashCode(any(id1))
		_ = ObjectHashCode(any(id1))
		_ = ObjectHashCode(any(eq1))
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs on ObjectsEquals/ObjectHashCode/IdentityHashCode fast paths, got %v", allocs)
	}
}
