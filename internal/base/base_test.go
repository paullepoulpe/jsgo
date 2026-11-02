package base

import (
	"math"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestTri(t *testing.T) {
	vals := [3]Tri{Tri_FALSE, Tri_UNKNOWN, Tri_TRUE}
	wantOr := [3][3]Tri{
		{Tri_FALSE, Tri_UNKNOWN, Tri_TRUE},
		{Tri_UNKNOWN, Tri_UNKNOWN, Tri_TRUE},
		{Tri_TRUE, Tri_TRUE, Tri_TRUE},
	}
	wantAnd := [3][3]Tri{
		{Tri_FALSE, Tri_FALSE, Tri_FALSE},
		{Tri_FALSE, Tri_UNKNOWN, Tri_UNKNOWN},
		{Tri_FALSE, Tri_UNKNOWN, Tri_TRUE},
	}
	wantXor := [3][3]Tri{
		{Tri_FALSE, Tri_UNKNOWN, Tri_TRUE},
		{Tri_UNKNOWN, Tri_UNKNOWN, Tri_UNKNOWN},
		{Tri_TRUE, Tri_UNKNOWN, Tri_FALSE},
	}
	wantNot := [3]Tri{Tri_TRUE, Tri_UNKNOWN, Tri_FALSE}
	for i, a := range vals {
		if got := a.Not(); got != wantNot[i] {
			t.Fatalf("%v.Not() = %v, want %v", a, got, wantNot[i])
		}
		for j, b := range vals {
			if got := a.Or(b); got != wantOr[i][j] {
				t.Fatalf("%v.Or(%v) = %v, want %v", a, b, got, wantOr[i][j])
			}
			if got := a.And(b); got != wantAnd[i][j] {
				t.Fatalf("%v.And(%v) = %v, want %v", a, b, got, wantAnd[i][j])
			}
			if got := a.Xor(b); got != wantXor[i][j] {
				t.Fatalf("%v.Xor(%v) = %v, want %v", a, b, got, wantXor[i][j])
			}
		}
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
	posZero := 0.0
	if !math.IsNaN(MathMax(math.NaN(), 1)) || !math.IsNaN(MathMax(1, math.NaN())) || !math.IsNaN(MathMax(math.NaN(), math.NaN())) {
		t.Fatal("MathMax NaN")
	}
	if math.Signbit(MathMax(negZero, posZero)) || math.Signbit(MathMax(posZero, negZero)) ||
		!math.Signbit(MathMax(negZero, negZero)) || math.Signbit(MathMax(posZero, posZero)) ||
		MathMax(-2, 3) != 3 || MathMax(3, -2) != 3 {
		t.Fatal("MathMax zeros/values")
	}
	if !math.IsNaN(MathMin(math.NaN(), 1)) || !math.IsNaN(MathMin(1, math.NaN())) || !math.IsNaN(MathMin(math.NaN(), math.NaN())) {
		t.Fatal("MathMin NaN")
	}
	if !math.Signbit(MathMin(negZero, posZero)) || !math.Signbit(MathMin(posZero, negZero)) ||
		!math.Signbit(MathMin(negZero, negZero)) || math.Signbit(MathMin(posZero, posZero)) ||
		MathMin(-2, 3) != -2 || MathMin(3, -2) != -2 {
		t.Fatal("MathMin zeros/values")
	}
	if !math.Signbit(MathSignum(negZero)) || MathSignum(-3.5) != -1 || !math.IsNaN(MathSignum(math.NaN())) {
		t.Fatal("MathSignum")
	}
	if MathRound(0.5) != 1 || MathRound(-0.5) != 0 || MathRound(math.Float64frombits(0x3fdfffffffffffff)) != 0 || MathRound(math.NaN()) != 0 {
		t.Fatal("MathRound")
	}
	if MathPow(10, 2) != 100 || !math.Signbit(MathPow(negZero, 3)) || !math.IsNaN(MathPow(1, math.Inf(1))) || MathPow(2, 0.5) != math.Sqrt2 {
		t.Fatal("MathPow")
	}
}
