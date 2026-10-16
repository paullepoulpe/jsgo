package jutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Go-only: language-adapter — Java numeric semantics.

// D2I is the Java narrowing cast (int) d: NaN → 0, saturating at the int32
// range, otherwise truncation toward zero.
func D2I(d float64) int32 {
	switch {
	case math.IsNaN(d):
		return 0
	case d >= math.MaxInt32:
		return math.MaxInt32
	case d <= math.MinInt32:
		return math.MinInt32
	}
	return int32(d)
}

// D2L is the Java narrowing cast (long) d.
func D2L(d float64) int64 {
	switch {
	case math.IsNaN(d):
		return 0
	case d >= math.MaxInt64:
		return math.MaxInt64
	case d <= math.MinInt64:
		return math.MinInt64
	}
	return int64(d)
}

// Shl32 is Java's int `x << n` (shift distance masked to 5 bits).
func Shl32(x int32, n int64) int32 { return x << (uint64(n) & 31) }

// Shr32 is Java's int `x >> n` (arithmetic).
func Shr32(x int32, n int64) int32 { return x >> (uint64(n) & 31) }

// Ushr32 is Java's int `x >>> n` (logical).
func Ushr32(x int32, n int64) int32 { return int32(uint32(x) >> (uint64(n) & 31)) }

// Shl64 is Java's long `x << n` (shift distance masked to 6 bits).
func Shl64(x int64, n int64) int64 { return x << (uint64(n) & 63) }

// Shr64 is Java's long `x >> n`.
func Shr64(x int64, n int64) int64 { return x >> (uint64(n) & 63) }

// Ushr64 is Java's long `x >>> n`.
func Ushr64(x int64, n int64) int64 { return int64(uint64(x) >> (uint64(n) & 63)) }

// DoubleToString is Java's Double.toString(double) (JDK 19+ algorithm): the
// shortest decimal that rounds to d (two digits considered when the shortest
// has one digit, picking the closest), rendered as plain decimal when
// 1e-3 <= |d| < 1e7 and as computerized scientific notation otherwise.
func DoubleToString(d float64) string {
	switch {
	case math.IsNaN(d):
		return "NaN"
	case math.IsInf(d, 1):
		return "Infinity"
	case math.IsInf(d, -1):
		return "-Infinity"
	case d == 0:
		if math.Signbit(d) {
			return "-0.0"
		}
		return "0.0"
	}
	neg := d < 0
	a := math.Abs(d)
	// Shortest round-tripping digits.
	s := strconv.FormatFloat(a, 'e', -1, 64) // e.g. "1.2345e+06"
	digits, exp := splitSci(s)
	if len(digits) == 1 {
		// Java considers 2-digit decimals when the shortest has 1 digit and
		// picks the one closest to the exact value (e.g. 4.9E-324).
		s2 := strconv.FormatFloat(a, 'e', 1, 64)
		if v, err := strconv.ParseFloat(s2, 64); err == nil && v == a {
			d2, e2 := splitSci(s2)
			d2 = strings.TrimRight(d2, "0")
			if d2 == "" {
				d2 = "0"
			}
			digits, exp = d2, e2
		}
	}
	var sb strings.Builder
	if neg {
		sb.WriteByte('-')
	}
	// digits represent 0.d1d2... * 10^(exp+1), i.e. d1.d2... * 10^exp.
	if a >= 1e-3 && a < 1e7 {
		if exp >= 0 {
			intLen := exp + 1
			if len(digits) <= intLen {
				sb.WriteString(digits)
				sb.WriteString(strings.Repeat("0", intLen-len(digits)))
				sb.WriteString(".0")
			} else {
				sb.WriteString(digits[:intLen])
				sb.WriteByte('.')
				sb.WriteString(digits[intLen:])
			}
		} else {
			sb.WriteString("0.")
			sb.WriteString(strings.Repeat("0", -exp-1))
			sb.WriteString(digits)
		}
		return sb.String()
	}
	sb.WriteByte(digits[0])
	sb.WriteByte('.')
	if len(digits) > 1 {
		sb.WriteString(digits[1:])
	} else {
		sb.WriteByte('0')
	}
	sb.WriteByte('E')
	sb.WriteString(strconv.Itoa(exp))
	return sb.String()
}

// splitSci splits Go's "d.ddde±XX" into the digit string and exponent.
func splitSci(s string) (string, int) {
	e := strings.IndexByte(s, 'e')
	mant, expStr := s[:e], s[e+1:]
	exp, _ := strconv.Atoi(expStr)
	digits := strings.Replace(mant, ".", "", 1)
	return digits, exp
}

// ToJavaString is String.valueOf(Object) for values flowing through Java
// string concatenation or formatting: nil → "null", float64 →
// Double.toString, uint16 (Java char) → the character, Stringer → String().
func ToJavaString(v any) string {
	if IsNil(v) {
		return "null"
	}
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return DoubleToString(x)
	case float32:
		return DoubleToString(float64(x))
	case bool:
		if x {
			return "true"
		}
		return "false"
	case uint16:
		return StringFromChar(x)
	case int:
		return strconv.Itoa(x)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case fmt.Stringer:
		return x.String()
	case error:
		return x.Error()
	}
	return fmt.Sprint(v)
}

// DoubleCompare is Double.compare(double, double): numeric order, with -0.0
// below 0.0 and NaN above everything (and equal to itself).
func DoubleCompare(a, b float64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	ab, bb := int64(doubleBits(a)), int64(doubleBits(b))
	switch {
	case ab == bb:
		return 0
	case ab < bb:
		return -1
	}
	return 1
}

// DoubleToLongBits is Double.doubleToLongBits(double).
func DoubleToLongBits(d float64) int64 { return int64(doubleBits(d)) }

// DoubleToRawLongBits is Double.doubleToRawLongBits(double).
func DoubleToRawLongBits(d float64) int64 { return int64(math.Float64bits(d)) }

// LongBitsToDouble is Double.longBitsToDouble(long).
func LongBitsToDouble(bits int64) float64 { return math.Float64frombits(uint64(bits)) }
