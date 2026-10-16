package jutil

// Go-only: language-adapter — java.lang.Math methods whose Go counterparts
// differ from Java (special cases of max/min/pow, round's rounding, signum),
// tested against JDK 21 by TestJavaMathMatchesJava.

import (
	"math"
	"math/big"
	"sync"
)

// MathMax is Math.max(double, double) (see its Javadoc): NaN when either
// argument is NaN, and +0.0 is taken to be greater than -0.0, which `>`
// alone does not distinguish.
//
// Observed behaviour (black-box, JDK 21): a single NaN argument is returned
// with its bits unchanged. With two NaNs the JDK returns the first in the
// interpreter and may return the second once the caller is JIT-compiled;
// this returns the first, as MathMin does.
func MathMax(a, b float64) float64 {
	if a != a { // a is NaN
		return a
	}
	if b != b {
		return b
	}
	if a == 0 && b == 0 {
		// Both zeros: the result is -0.0 only when both are -0.0.
		if math.Signbit(a) {
			return b
		}
		return a
	}
	if a >= b {
		return a
	}
	return b
}

// MathMin is Math.min(double, double): NaN if either argument is NaN (the
// first one if both are), and -0.0 below +0.0.
func MathMin(a, b float64) float64 {
	if a != a { // a is NaN
		return a
	}
	if a == 0 && b == 0 && math.Float64bits(b) == math.Float64bits(math.Copysign(0, -1)) {
		return b
	}
	if a <= b {
		return a
	}
	return b
}

// MathSignum is Math.signum(double): the argument itself for zeros and NaN,
// otherwise ±1.0.
func MathSignum(d float64) float64 {
	if d == 0 || d != d {
		return d
	}
	return math.Copysign(1.0, d)
}

// MathRound is Math.round(double) (see its Javadoc): the long nearest to a,
// halfway cases going up (towards positive infinity), so it is
// floor(a + 1/2) computed exactly; NaN gives 0, and values beyond the long
// range (infinities included) saturate to math.MinInt64 / math.MaxInt64.
func MathRound(a float64) int64 {
	if a != a {
		return 0
	}
	// The comparison of the fraction a - floor(a) with 0.5 is exact: the
	// subtraction is exact for a >= 0 and for a <= -1 (Sterbenz's lemma),
	// and for -1 < a < 0 it is 1+a, which rounds only when a > -0.5, where
	// it stays >= 0.5. Computing floor(a + 0.5) in floating point instead
	// would round the sum, e.g. for 0.49999999999999994 or 2^52+1.
	f := math.Floor(a)
	if a-f >= 0.5 {
		f++
	}
	switch {
	case f >= 1<<63: // 2^63 is the first double above math.MaxInt64
		return math.MaxInt64
	case f <= -(1 << 63):
		return math.MinInt64
	}
	return int64(f)
}

// MathPow is Math.pow(double, double).
//
// The special cases are exactly Java's. Otherwise the result is the
// correctly rounded value of x^y. Java only requires 1 ulp; HotSpot's x86_64
// intrinsic (Intel LIBM, using rcpps, whose output varies between CPU
// models) is correctly rounded on all the integer-power vectors and on all but
// ~0.05% of random arguments, where it is 1 ulp away. In particular, as Java
// requires, results that are exactly representable (integer powers of
// integers, powers of two, ...) are exact.
func MathPow(x, y float64) float64 {
	switch {
	case y == 0:
		return 1.0
	case y == 1:
		return x
	case y != y, x != x:
		return math.NaN()
	}
	ax := math.Abs(x)
	if math.IsInf(y, 0) {
		switch {
		case ax == 1:
			return math.NaN()
		case (ax > 1) == (y > 0):
			return math.Inf(1)
		default:
			return 0
		}
	}
	yIsInt := y == math.Trunc(y)
	yIsOdd := yIsInt && math.Abs(y) < 1<<53 && math.Mod(y, 2) != 0
	if x == 0 || math.IsInf(x, 0) {
		// ±0 and ±Inf: the magnitude is 0 or Inf, the sign comes from x only
		// for odd integer y.
		var r float64
		if (x == 0) == (y > 0) {
			r = 0
		} else {
			r = math.Inf(1)
		}
		if math.Signbit(x) && yIsOdd {
			r = -r
		}
		return r
	}
	negate := false
	if x < 0 {
		if !yIsInt {
			return math.NaN()
		}
		negate = yIsOdd
	}
	r := mathPowPositive(ax, y, yIsInt)
	if negate {
		return -r
	}
	return r
}

// mathPowPrec is the working precision (bits) of mathPowPositive's integer
// powers: enough for repeated squaring with exponents below 2^62 to be
// exact whenever the result is representable, and to be off by far less
// than half an ulp otherwise.
const mathPowPrec = 1280

// mathPowExpPrec is the working precision (bits) of the exp(y*ln(x)) path.
const mathPowExpPrec = 320

// mathPowPositive returns the correctly rounded value of x^y for finite
// x > 0 and finite y not in {0, 1}.
func mathPowPositive(x, y float64, yIsInt bool) float64 {
	// Results whose binary exponent is far outside the double range
	// overflow or underflow; this also keeps big.Float's exponent small.
	est := y * math.Log2(x)
	if est > 2000 {
		return math.Inf(1)
	}
	if est < -2000 {
		return 0
	}
	bx := new(big.Float).SetPrec(mathPowPrec).SetFloat64(x)
	var res *big.Float
	if yIsInt && math.Abs(y) < 1<<62 {
		n := int64(math.Abs(y))
		res = new(big.Float).SetPrec(mathPowPrec).SetInt64(1)
		base := new(big.Float).SetPrec(mathPowPrec).Set(bx)
		for n > 0 {
			if n&1 == 1 {
				res.Mul(res, base)
			}
			n >>= 1
			if n > 0 {
				base.Mul(base, base)
			}
		}
		if y < 0 {
			res.Quo(new(big.Float).SetPrec(mathPowPrec).SetInt64(1), res)
		}
	} else {
		t := bigLog(new(big.Float).SetPrec(mathPowExpPrec).SetFloat64(x))
		t.Mul(t, new(big.Float).SetPrec(mathPowExpPrec).SetFloat64(y))
		res = bigExp(t)
	}
	f, _ := res.Float64()
	return f
}

var (
	bigLn2Once sync.Once
	bigLn2     *big.Float
)

// bigLn2Value is ln(2) to mathPowExpPrec bits.
func bigLn2Value() *big.Float {
	bigLn2Once.Do(func() {
		third := new(big.Float).SetPrec(mathPowExpPrec).Quo(
			new(big.Float).SetPrec(mathPowExpPrec).SetInt64(1),
			new(big.Float).SetPrec(mathPowExpPrec).SetInt64(3))
		bigLn2 = bigAtanhTimes2(third)
	})
	return bigLn2
}

// bigAtanhTimes2 returns 2*atanh(z) = ln((1+z)/(1-z)) for |z| <= 1/3, by
// its power series.
func bigAtanhTimes2(z *big.Float) *big.Float {
	z2 := new(big.Float).SetPrec(mathPowExpPrec).Mul(z, z)
	sum := new(big.Float).SetPrec(mathPowExpPrec).Set(z)
	term := new(big.Float).SetPrec(mathPowExpPrec).Set(z)
	q := new(big.Float).SetPrec(mathPowExpPrec)
	for k := int64(1); ; k++ {
		term.Mul(term, z2)
		q.Quo(term, new(big.Float).SetPrec(mathPowExpPrec).SetInt64(2*k+1))
		if q.Sign() == 0 || q.MantExp(nil)-sum.MantExp(nil) < -mathPowExpPrec-8 {
			break
		}
		sum.Add(sum, q)
	}
	return sum.Mul(sum, new(big.Float).SetPrec(mathPowExpPrec).SetInt64(2))
}

// bigLog returns ln(x) for x > 0.
func bigLog(x *big.Float) *big.Float {
	m := new(big.Float).SetPrec(mathPowExpPrec)
	e := x.MantExp(m) // x = m * 2^e, m in [0.5, 1)
	// Move m into [sqrt(1/2), sqrt(2)) for faster convergence.
	if m.Cmp(big.NewFloat(math.Sqrt2/2)) < 0 {
		m.SetMantExp(m, 1)
		e--
	}
	one := new(big.Float).SetPrec(mathPowExpPrec).SetInt64(1)
	num := new(big.Float).SetPrec(mathPowExpPrec).Sub(m, one)
	den := new(big.Float).SetPrec(mathPowExpPrec).Add(m, one)
	z := num.Quo(num, den)
	res := bigAtanhTimes2(z)
	if e != 0 {
		le := new(big.Float).SetPrec(mathPowExpPrec).SetInt64(int64(e))
		le.Mul(le, bigLn2Value())
		res.Add(res, le)
	}
	return res
}

// bigExp returns e^t for |t| < 2000*ln(2).
func bigExp(t *big.Float) *big.Float {
	ln2 := bigLn2Value()
	// t = k*ln2 + r, |r| <= ln2/2.
	kf := new(big.Float).SetPrec(mathPowExpPrec).Quo(t, ln2)
	kd, _ := kf.Float64()
	k := int64(math.RoundToEven(kd))
	r := new(big.Float).SetPrec(mathPowExpPrec).Mul(new(big.Float).SetPrec(mathPowExpPrec).SetInt64(k), ln2)
	r.Sub(t, r)
	// e^r = (e^(r/2^s))^(2^s).
	const s = 16
	r.SetMantExp(r, -s)
	sum := new(big.Float).SetPrec(mathPowExpPrec).SetInt64(1)
	term := new(big.Float).SetPrec(mathPowExpPrec).SetInt64(1)
	for i := int64(1); ; i++ {
		term.Mul(term, r)
		term.Quo(term, new(big.Float).SetPrec(mathPowExpPrec).SetInt64(i))
		if term.Sign() == 0 || term.MantExp(nil) < -mathPowExpPrec-8 {
			break
		}
		sum.Add(sum, term)
	}
	for i := 0; i < s; i++ {
		sum.Mul(sum, sum)
	}
	return sum.SetMantExp(sum, int(k))
}
