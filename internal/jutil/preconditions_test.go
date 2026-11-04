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
}

func assertNullPointerPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("%s: expected NullPointerException panic, got none", name)
		}
		if _, ok := r.(*NullPointerException); !ok {
			t.Fatalf("%s: expected *NullPointerException, got %T (%v)", name, r, r)
		}
	}()
	fn()
}

func TestCheckNotNullAndPtrBehavior(t *testing.T) {
	type sample struct{ x int }
	s := &sample{x: 7}
	if got := CheckNotNull(s); got != s {
		t.Fatalf("CheckNotNull(sample) = %v, want %v", got, s)
	}
	if got := CheckNotNullPtr(s); got != s {
		t.Fatalf("CheckNotNullPtr(sample) = %v, want %v", got, s)
	}
	if got := CheckNotNull([]string{"a"}); len(got) != 1 {
		t.Fatalf("CheckNotNull([]string) = %v", got)
	}
	if got := CheckNotNull(map[string]string{"k": "v"}); len(got) != 1 {
		t.Fatalf("CheckNotNull(map) = %v", got)
	}
	if got := CheckNotNull(testNodeID(9)); got != 9 {
		t.Fatalf("CheckNotNull(testNodeID(9)) = %v", got)
	}

	var nilSample *sample
	var nilStrings []string
	var nilMap map[string]string
	var nilAny any
	var nilFn func()

	assertNullPointerPanic(t, "typed nil struct ptr", func() { CheckNotNull(nilSample) })
	assertNullPointerPanic(t, "CheckNotNullPtr nil", func() { CheckNotNullPtr(nilSample) })
	assertNullPointerPanic(t, "nil []string", func() { CheckNotNull(nilStrings) })
	assertNullPointerPanic(t, "nil map[string]string", func() { CheckNotNull(nilMap) })
	assertNullPointerPanic(t, "nil interface", func() { CheckNotNull(nilAny) })
	assertNullPointerPanic(t, "typed nil in any", func() { CheckNotNull(any(nilSample)) })
	assertNullPointerPanic(t, "nil func", func() { CheckNotNull(nilFn) })
	assertNullPointerPanic(t, "nil custom handle", func() { CheckNotNull(testNodeID(0)) })
}
