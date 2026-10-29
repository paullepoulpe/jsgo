// Port of com/google/javascript/jscomp/base/JSCompDoubles.java @ closure-compiler v20261021

/*
 * Copyright 2020 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package base is the port of com.google.javascript.jscomp.base and
// com.google.javascript.jscomp.base.format.
package base

import (
	"math"
	"math/big"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSCompDoubles holds basic double functions used by JSComp.
type JSCompDoubles struct{}

// JSCompDoubles_IsExactInt32 reports whether a 32 bit int can exactly
// represent `x`.
//
// Many double values are not exact integers. Many that are integers are too
// large to fit into a Java int.
//
// This function does not guarantee that a value can be round-tripped from
// double to int to double and have an identical bit pattern. Notably, 0.0
// and -0.0 both represent exactly 0.
//
// Java: JSCompDoubles.isExactInt32(double)
func JSCompDoubles_IsExactInt32(x float64) bool {
	return !math.IsNaN(x) && float64(jutil.D2I(x)) == x
}

// JSCompDoubles_IsExactInt64 reports whether a 64 bit int can exactly
// represent `x`.
//
// Many double values are not exact integers. Many that are integers are too
// large to fit into a Java long.
//
// This function does not guarantee that a value can be round-tripped from
// double to long to double and have an identical bit pattern. Notably, 0.0
// and -0.0 both represent exactly 0.
//
// Java: JSCompDoubles.isExactInt64(double)
func JSCompDoubles_IsExactInt64(x float64) bool {
	return !math.IsNaN(x) && float64(jutil.D2L(x)) == x
}

// JSCompDoubles_IsMathematicalInteger reports whether `x` exactly represents
// a value with no fractional part.
//
// The value may be too large to fit in a primitive integral type, such as
// long.
//
// Returns false for NaN and Infinity.
//
// This should behave identically to Guava DoubleMath.isMathematicalInteger
// but is J2CL compatible.
//
// Java: JSCompDoubles.isMathematicalInteger(double)
func JSCompDoubles_IsMathematicalInteger(x float64) bool {
	return math.Mod(x, 1.0) == 0.0
}

// JSCompDoubles_IsAtLeastIntegerPrecision reports whether `x` has precision
// down to the "ones column".
//
// A double can hold exact integer values that are very large, but to do so
// it may loose precision at the scale of "ones". That is "largeDouble + 1.0
// == largeDouble" may be true.
//
// Returns false for NaN and Infinity.
//
// Java: JSCompDoubles.isAtLeastIntegerPrecision(double)
func JSCompDoubles_IsAtLeastIntegerPrecision(x float64) bool {
	// It's possible for a JVM to use floating-point values with more than 53
	// bits of precision, but we have to be conservative.
	return math.Abs(x) < JSCompDoubles_POW_2_53
}

// JSCompDoubles_IsNegative reports whether `x` carries a negative sign.
//
// Because -0.0 == 0.0, it is not enough to check `x < 0.0` to determine if a
// double is negative. This function identifies the -0.0 case.
//
// Throws if passed NaN since NaNs are neither positive nor negative.
//
// Java: JSCompDoubles.isNegative(double)
func JSCompDoubles_IsNegative(x float64) bool {
	jutil.CheckState(!math.IsNaN(x))
	return jutil.DoubleCompare(0.0, x) > 0
}

// JSCompDoubles_IsPositive reports whether `x` does not carry a negative
// sign.
//
// Because -0.0 == 0.0, it is not enough to check `x < 0.0` to determine if a
// double is negative. This function identifies the -0.0 case.
//
// Throws if passed NaN since NaNs are neither positive nor negative.
//
// Java: JSCompDoubles.isPositive(double)
func JSCompDoubles_IsPositive(x float64) bool {
	return !JSCompDoubles_IsNegative(x)
}

// JSCompDoubles_IsEitherZero reports whether `x` is positive or negative
// zero.
//
// Java: JSCompDoubles.isEitherZero(double)
func JSCompDoubles_IsEitherZero(x float64) bool {
	return x == 0.0
}

// JSCompDoubles_EcmascriptToInt32 is the ECMAScript ToInt32 abstract
// operation.
//
// See https://262.ecma-international.org/5.1/#sec-9.5
//
// Java: JSCompDoubles.ecmascriptToInt32(double)
func JSCompDoubles_EcmascriptToInt32(number float64) int32 {
	// Fast path for most common case. Also covers -0.0
	if JSCompDoubles_IsExactInt32(number) {
		return jutil.D2I(number)
	}

	// Step 2
	if math.IsNaN(number) || math.IsInf(number, 0) {
		return 0
	}

	// Step 3
	var posInt float64
	if number >= 0 {
		posInt = math.Floor(number)
	} else {
		posInt = math.Ceil(number)
	}

	// Step 4
	//
	// Afer modding, the value must fit into a long, but not necessarily an
	// int.
	//
	// If value produced by modding is negative, the spec says to add 2^32 so
	// that `int32bit` is positive. However, that addition doesn't affect the
	// trailing 32 bits of the `int32bit`, which are the only ones that matter.
	int32bit := jutil.D2L(math.Mod(posInt, JSCompDoubles_POW_2_32))

	// Step 5
	//
	// The instructions in the spec are equivalent to interpreting the last 32
	// bits of `int32bit` as a 2s-compliment integer, which is exactly the
	// representation Java uses for int.
	return int32(int32bit)
}

// JSCompDoubles_EcmascriptToUint32 is the ECMAScript ToUint32 abstract
// operation.
//
// Java has no uint types, so the caller must remember to treat the returned
// bits as a uint.
//
// See https://262.ecma-international.org/5.1/#sec-9.6
//
// Java: JSCompDoubles.ecmascriptToUint32(double)
func JSCompDoubles_EcmascriptToUint32(number float64) int32 {
	return JSCompDoubles_EcmascriptToInt32(number)
}

// Java: JSCompDoubles.POW_2_32
var JSCompDoubles_POW_2_32 = math.Pow(2, 32)

// Java: JSCompDoubles.POW_2_53
var JSCompDoubles_POW_2_53 = math.Pow(2, 53)

// NewJSCompDoubles is the private constructor, which always throws.
//
// Java: JSCompDoubles()
func NewJSCompDoubles() *JSCompDoubles {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Go-only: language-adapter — java.lang.Math methods whose Go counterparts
// differ from Java (special cases of max/min/pow, round's rounding, signum).

// MathMax is Math.max(double, double). Javadoc: the result is NaN when
// either argument is, and negative zero counts as less than positive zero
// (which Go's comparison operators treat as equal).
func MathMax(a, b float64) float64 {
	switch {
	case math.IsNaN(a):
		return a
	case math.IsNaN(b):
		return b
	case a == 0 && b == 0:
		// Two zeros: the one without the sign bit, if there is one.
		if math.Signbit(b) {
			return a
		}
		return b
	case a > b:
		return a
	}
	return b
}

// MathMin is Math.min(double, double). Javadoc: the result is NaN when
// either argument is, and negative zero counts as less than positive zero.
func MathMin(a, b float64) float64 {
	switch {
	case math.IsNaN(a):
		return a
	case math.IsNaN(b):
		return b
	case a == 0 && b == 0:
		// Two zeros: the one with the sign bit, if there is one.
		if math.Signbit(a) {
			return a
		}
		return b
	case a < b:
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

const (
	mathPowPrec    = 1280
	mathPowExpPrec = 320
)

func mathPowPositive(x, y float64, yIsInt bool) float64 {
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

func bigLn2Value() *big.Float {
	bigLn2Once.Do(func() {
		third := new(big.Float).SetPrec(mathPowExpPrec).Quo(
			new(big.Float).SetPrec(mathPowExpPrec).SetInt64(1),
			new(big.Float).SetPrec(mathPowExpPrec).SetInt64(3))
		bigLn2 = bigAtanhTimes2(third)
	})
	return bigLn2
}

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

func bigLog(x *big.Float) *big.Float {
	m := new(big.Float).SetPrec(mathPowExpPrec)
	e := x.MantExp(m)
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

func bigExp(t *big.Float) *big.Float {
	ln2 := bigLn2Value()
	kf := new(big.Float).SetPrec(mathPowExpPrec).Quo(t, ln2)
	kd, _ := kf.Float64()
	k := int64(math.RoundToEven(kd))
	r := new(big.Float).SetPrec(mathPowExpPrec).Mul(new(big.Float).SetPrec(mathPowExpPrec).SetInt64(k), ln2)
	r.Sub(t, r)
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
