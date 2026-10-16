// Port of com/google/javascript/jscomp/base/JSCompDoubles.java @ closure-compiler v20260415

// Package base is the port of com.google.javascript.jscomp.base and
// com.google.javascript.jscomp.base.format.
package base

import (
	"math"

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
