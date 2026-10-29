package jutil

import "testing"

func TestIsNilTypedPointersAndSlices(t *testing.T) {
	var (
		nilInt     *int
		nilInt32   *int32
		nilInt64   *int64
		nilBool    *bool
		nilFloat64 *float64
		nilSlice   []int
		nilMap     map[string]int
		nilFunc    func()
	)
	vInt := 1
	vInt32 := int32(2)
	vInt64 := int64(3)
	vBool := true
	vFloat64 := 4.5
	nonNilSlice := []int{}
	nonNilMap := map[string]int{}
	nonNilFunc := func() {}

	cases := []struct {
		v    any
		want bool
	}{
		{nilInt, true},
		{&vInt, false},
		{nilInt32, true},
		{&vInt32, false},
		{nilInt64, true},
		{&vInt64, false},
		{nilBool, true},
		{&vBool, false},
		{nilFloat64, true},
		{&vFloat64, false},
		{nilSlice, true},
		{nonNilSlice, false},
		{nilMap, true},
		{nonNilMap, false},
		{nilFunc, true},
		{nonNilFunc, false},
	}
	for _, tc := range cases {
		if got := IsNil(tc.v); got != tc.want {
			t.Fatalf("IsNil(%#v) = %v, want %v", tc.v, got, tc.want)
		}
	}
	if !IsNilOf(nilSlice) || IsNilOf(nonNilSlice) || !IsNilOf(nilMap) || IsNilOf(nonNilMap) || !IsNilOf(nilFunc) || IsNilOf(nonNilFunc) {
		t.Fatal("IsNilOf mismatch for slice/map/func")
	}
}
