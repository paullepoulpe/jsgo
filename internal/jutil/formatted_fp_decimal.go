package jutil

import (
	"math"
	"strings"
)

// Go-only: language-adapter — the 'f' conversion of java.util.Formatter
// with an explicit precision, for a double (String.format("%.<n>f", d)).

// Formatter_FormatDecimalFloat is String.format("%.<precision>f", value)
// in a locale with ASCII digits, '.' as decimal separator and no grouping.
//
// The Formatter Javadoc ('f' conversion) specifies the digits: those of
// Double.toString(value), rounded with the round-half-up algorithm when the
// precision keeps fewer fraction digits, and padded with zeros otherwise.
// NaN and the infinities print as "NaN", "Infinity" and "-Infinity". The
// rounding works on the decimal digit string, so it is exact.
//
// Observed behaviour (black-box): the result equals
// new BigDecimal(Double.toString(value)).setScale(precision, HALF_UP)
// .toPlainString() for every value tried (400,000 random doubles with
// precisions 0..24, and tools/oracle/java/JutilNumbersVectors.java), except
// that a negative value (including -0.0) keeps its '-' even when it rounds
// to zero ("%.1f" of -0.04 is "-0.0"). So the rounding is not that of the
// double's exact binary value: 0.15 prints as "0.2" with "%.1f" although
// the double is slightly below 0.15, and 1e23 prints as
// "100000000000000000000000".
func Formatter_FormatDecimalFloat(value float64, precision int) string {
	switch {
	case math.IsNaN(value):
		return "NaN"
	case math.IsInf(value, 1):
		return "Infinity"
	case math.IsInf(value, -1):
		return "-Infinity"
	}
	digits, pointPos := javaShortestDigits(math.Abs(value))
	// scaled holds the digits of round(|value| × 10^precision).
	var scaled string
	keep := pointPos + precision // digits that stay before the rounding point
	switch {
	case keep < 0:
		scaled = "0"
	case keep >= len(digits):
		scaled = digits + strings.Repeat("0", keep-len(digits))
	default:
		scaled = digits[:keep]
		if digits[keep] >= '5' {
			scaled = incrementDecimalString(scaled)
		}
		if scaled == "" {
			scaled = "0"
		}
	}
	if len(scaled) <= precision {
		scaled = strings.Repeat("0", precision+1-len(scaled)) + scaled
	}
	var sb strings.Builder
	if math.Signbit(value) {
		sb.WriteByte('-')
	}
	intLen := len(scaled) - precision
	sb.WriteString(scaled[:intLen])
	if precision > 0 {
		sb.WriteByte('.')
		sb.WriteString(scaled[intLen:])
	}
	return sb.String()
}

// javaShortestDigits returns the decimal digits of Double.toString(a) for a
// finite a >= 0, without leading zeros ("0" for zero), and the position of
// the decimal point: a = 0.<digits> × 10^pointPos.
func javaShortestDigits(a float64) (string, int) {
	s := DoubleToString(a) // "123.45", "0.00123" or "1.2345E-5"
	exp := 0
	if e := strings.IndexByte(s, 'E'); e >= 0 {
		exp = int(IntegerParseInt(s[e+1:], 10))
		s = s[:e]
	}
	point := strings.IndexByte(s, '.')
	digits := s[:point] + s[point+1:]
	pointPos := point + exp
	for len(digits) > 1 && digits[0] == '0' {
		digits = digits[1:]
		pointPos--
	}
	digits = strings.TrimRight(digits, "0")
	if digits == "" {
		return "0", 1
	}
	return digits, pointPos
}

// incrementDecimalString adds one to a string of decimal digits ("" counts
// as 0), carrying into a new leading digit when needed.
func incrementDecimalString(s string) string {
	b := []byte(s)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] < '9' {
			b[i]++
			return string(b)
		}
		b[i] = '0'
	}
	return "1" + string(b)
}
