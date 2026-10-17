package base

import (
	"math"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestTri(t *testing.T) {
	if Tri_TRUE.Or(Tri_FALSE) != Tri_TRUE || Tri_TRUE.And(Tri_UNKNOWN) != Tri_UNKNOWN ||
		Tri_TRUE.Xor(Tri_TRUE) != Tri_FALSE || Tri_UNKNOWN.Xor(Tri_TRUE) != Tri_UNKNOWN ||
		Tri_FALSE.Not() != Tri_TRUE || Tri_UNKNOWN.Not() != Tri_UNKNOWN {
		t.Fatal("ternary logic")
	}
	if Tri_UNKNOWN.String() != "unknown" || Tri_UNKNOWN.Name() != "UNKNOWN" {
		t.Fatal("toString")
	}
	if !Tri_UNKNOWN.ToBoolean(true) || Tri_UNKNOWN.ToBoolean(false) {
		t.Fatal("toBoolean")
	}
}

func TestJSCompDoubles(t *testing.T) {
	negZero := math.Copysign(0, -1)
	if !JSCompDoubles_IsNegative(negZero) || JSCompDoubles_IsNegative(0) || !JSCompDoubles_IsPositive(0) {
		t.Fatal("isNegative")
	}
	threw := false
	jutil.Try(func() { JSCompDoubles_IsNegative(math.NaN()) }, func(e jutil.IIllegalStateException) { threw = true })
	if !threw {
		t.Fatal("isNegative(NaN) must throw")
	}
	cases := map[float64]int32{
		0: 0, 1.5: 1, -1.5: -1, 4294967296: 0, 4294967297: 1, 2147483648: -2147483648,
		-2147483649: 2147483647, math.Inf(1): 0, math.NaN(): 0, 1e20: 1661992960,
	}
	for in, want := range cases {
		if got := JSCompDoubles_EcmascriptToInt32(in); got != want {
			t.Errorf("ToInt32(%v) = %d, want %d", in, got, want)
		}
	}
	if !JSCompDoubles_IsExactInt32(2147483647) || JSCompDoubles_IsExactInt32(2147483648) ||
		JSCompDoubles_IsExactInt32(math.NaN()) || !JSCompDoubles_IsMathematicalInteger(1e300) ||
		JSCompDoubles_IsMathematicalInteger(math.Inf(1)) || !JSCompDoubles_IsAtLeastIntegerPrecision(1e15) {
		t.Fatal("exactness predicates")
	}
}

func TestLinkedIdentityHashSet(t *testing.T) {
	type obj struct{ _ byte }
	a, b := &obj{}, &obj{}
	s := NewLinkedIdentityHashSet[*obj]()
	if !s.Add(a) || s.Add(a) || !s.Add(b) || !s.Contains(a) || !s.Remove(a) || s.Contains(a) || !s.Add(a) {
		t.Fatal("set semantics")
	}
	m := NewLinkedIdentityHashMap[*obj, string]()
	m.Put(b, "b")
	m.Put(a, "a")
	m.Merge(a, "x", func(o, n string) string { return o + n })
	var order string
	m.ForEach(func(_ *obj, v string) { order += v + "," })
	if order != "b,ax," {
		t.Fatalf("order %q", order)
	}
}

// A non-nillable V (Java Integer): a missing key is null, a present 0 is not.
func TestLinkedIdentityHashMapIntValues(t *testing.T) {
	type obj struct{ _ byte }
	a, b := &obj{}, &obj{}
	m := NewLinkedIdentityHashMap[*obj, int]()
	if got := m.ComputeIfAbsent(a, func(*obj) int { return 0 }); got != 0 {
		t.Fatalf("computeIfAbsent(a) = %d", got)
	}
	if got := m.ComputeIfAbsent(a, func(*obj) int { return 7 }); got != 0 {
		t.Fatalf("computeIfAbsent(a) again = %d, want the stored 0", got)
	}
	if got := m.Merge(b, 5, func(o, n int) int { return o + n }); got != 5 {
		t.Fatalf("merge(b) on a missing key = %d", got)
	}
	if v, ok := m.Lookup(a); !ok || v != 0 {
		t.Fatalf("lookup(a) = %d, %v", v, ok)
	}
	if _, ok := m.Lookup(&obj{}); ok {
		t.Fatal("lookup of a missing key reports present")
	}
	var order []int
	m.ForEach(func(_ *obj, v int) { order = append(order, v) })
	if len(order) != 2 || order[0] != 0 || order[1] != 5 {
		t.Fatalf("order %v", order)
	}
}

func TestFormatExceptions(t *testing.T) {
	if e := NewIllegalFormatConversionException('d', "class java.lang.String"); e.GetMessage() != "ConversionType = d class = class java.lang.String" {
		t.Fatal(e.GetMessage())
	}
	if e := NewIllegalFormatWidthException(3); e.Error() != "com.google.javascript.jscomp.base.format.IllegalFormatWidthException: Width:3" {
		t.Fatal(e.Error())
	}
}
