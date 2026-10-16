package jutil

import (
	"math"
	"strconv"
)

// Go-only: language-adapter — Double.parseDouble(String) and
// Double.valueOf(String). The accepted grammar is the one in the Javadoc of
// Double.valueOf(String); the conversion itself is delegated to
// strconv.ParseFloat, which rounds correctly (to nearest, ties to even) as
// the Javadoc requires, on a canonical form built here. The exception
// messages were determined black-box against JDK 21 (vectors from
// tools/oracle/java/JutilParsingVectors.java and JutilNumbersVectors.java).

// DoubleParseDouble is Double.parseDouble(String s).
//
// Per the Javadoc of Double.valueOf(String): leading and trailing chars
// <= U+0020 are ignored (as String.trim does); then an optional sign and
// either "NaN", "Infinity", a decimal floating-point literal or a
// hexadecimal one ("0x"/"0X", hex digits with an optional point, and a
// mandatory binary exponent 'p'/'P'), as in the JLS but with ASCII digits
// only and no underscores. A decimal literal needs at least one digit and
// may have an exponent 'e'/'E'; both kinds may end with one of 'f', 'F',
// 'd', 'D', which does not change the result (the value is always rounded
// to double).
//
// Observed behaviour (black-box): the messages are "empty String" when
// nothing is left after trimming; "multiple points" when the run of digits
// and '.' that follows the sign of a decimal literal holds two or more
// points; and otherwise `For input string: "<trimmed string>"`. NaN keeps
// no sign: "-NaN" is the canonical NaN.
func DoubleParseDouble(s string) float64 {
	t := StringTrim(s)
	if t == "" {
		panic(NewNumberFormatException("empty String"))
	}
	v, outcome := parseJavaDouble(t)
	switch outcome {
	case parseDoubleOK:
		return v
	case parseDoubleMultiplePoints:
		panic(NewNumberFormatException("multiple points"))
	}
	panic(NewNumberFormatException("For input string: \"" + t + "\""))
}

// Double_NaN is Double.NaN (the canonical NaN, raw bits 0x7ff8000000000000).
var Double_NaN = math.Float64frombits(0x7ff8000000000000)

const (
	parseDoubleOK = iota
	parseDoubleMultiplePoints
	parseDoubleInvalid
)

// parseJavaDouble parses a trimmed, non-empty string.
func parseJavaDouble(t string) (float64, int) {
	i := 0
	negative := false
	if t[0] == '-' || t[0] == '+' {
		negative = t[0] == '-'
		i = 1
	}
	body := t[i:]
	var v float64
	var outcome int
	switch {
	case body == "":
		return 0, parseDoubleInvalid
	case body[0] == 'N':
		if body != "NaN" {
			return 0, parseDoubleInvalid
		}
		return Double_NaN, parseDoubleOK
	case body[0] == 'I':
		if body != "Infinity" {
			return 0, parseDoubleInvalid
		}
		v, outcome = math.Inf(1), parseDoubleOK
	case len(body) >= 2 && body[0] == '0' && (body[1] == 'x' || body[1] == 'X'):
		v, outcome = parseJavaHexDouble(body[2:])
	default:
		v, outcome = parseJavaDecimalDouble(body)
	}
	if negative {
		v = -v
	}
	return v, outcome
}

func isASCIIDigit(c byte) bool { return c >= '0' && c <= '9' }

func isHexDigit(c byte) bool {
	return isASCIIDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// skipFloatSuffix steps over an optional float type suffix at s[i] and
// reports whether the literal ends there.
func skipFloatSuffix(s string, i int) bool {
	if i < len(s) {
		switch s[i] {
		case 'f', 'F', 'd', 'D':
			i++
		}
	}
	return i == len(s)
}

// scanSignedExponent reads an optional sign and one or more ASCII digits at
// s[i:], returning the value (its magnitude saturated at 1e15, far beyond
// any exponent that can still matter) and the index after the digits.
func scanSignedExponent(s string, i int) (value int64, end int, ok bool) {
	negative := false
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		negative = s[i] == '-'
		i++
	}
	start := i
	for i < len(s) && isASCIIDigit(s[i]) {
		if value < 1e15 {
			value = value*10 + int64(s[i]-'0')
		}
		i++
	}
	if i == start {
		return 0, i, false
	}
	if negative {
		value = -value
	}
	return value, i, true
}

// parseJavaDecimalDouble parses an unsigned decimal literal: digits with at
// most one point (at least one digit), an optional exponent and an
// optional type suffix.
func parseJavaDecimalDouble(s string) (float64, int) {
	i, points, digitCount := 0, 0, 0
	for i < len(s) && (isASCIIDigit(s[i]) || s[i] == '.') {
		if s[i] == '.' {
			points++
		} else {
			digitCount++
		}
		i++
	}
	if points > 1 {
		return 0, parseDoubleMultiplePoints
	}
	if digitCount == 0 {
		return 0, parseDoubleInvalid
	}
	significandEnd := i
	exponent := int64(0)
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		var ok bool
		exponent, i, ok = scanSignedExponent(s, i+1)
		if !ok {
			return 0, parseDoubleInvalid
		}
	}
	literalEnd := i
	if !skipFloatSuffix(s, i) {
		return 0, parseDoubleInvalid
	}
	if exponent > -10000 && exponent < 10000 {
		// Already in a form strconv.ParseFloat reads exactly (small exponents
		// are never saturated by it).
		if v, err := strconv.ParseFloat(s[:literalEnd], 64); err == nil || isRangeError(err) {
			return v, parseDoubleOK
		}
	}
	// Canonical form 0.<significant digits>e<exponent>, with the exponent
	// adjusted here so that strconv never sees a huge one.
	sig := make([]byte, 0, digitCount)
	intDigits := int64(0) // significant digits before the point
	seenPoint := false
	for j := 0; j < significandEnd; j++ {
		switch c := s[j]; {
		case c == '.':
			seenPoint = true
		case c == '0' && len(sig) == 0:
			// A leading zero only matters after the point, where it moves
			// the significant digits one place to the right.
			if seenPoint {
				intDigits--
			}
		default:
			sig = append(sig, c)
			if !seenPoint {
				intDigits++
			}
		}
	}
	if len(sig) == 0 {
		return 0, parseDoubleOK
	}
	// The value is 0.<sig> × 10^scale10, with 0.1 <= 0.<sig> < 1.
	scale10 := intDigits + exponent
	switch {
	case scale10 > 310:
		return math.Inf(1), parseDoubleOK
	case scale10 < -400:
		return 0, parseDoubleOK
	}
	v, _ := strconv.ParseFloat("0."+string(sig)+"e"+strconv.FormatInt(scale10, 10), 64)
	return v, parseDoubleOK
}

// parseJavaHexDouble parses what follows "0x": hex digits with an optional
// point (at least one digit), the binary exponent and an optional type
// suffix.
func parseJavaHexDouble(s string) (float64, int) {
	i, digitCount := 0, 0
	sig := make([]byte, 0, len(s))
	intDigits := int64(0) // significant hex digits before the point
	seenPoint := false
	for i < len(s) {
		c := s[i]
		switch {
		case isHexDigit(c):
			digitCount++
			if c != '0' || len(sig) > 0 {
				sig = append(sig, c)
				if !seenPoint {
					intDigits++
				}
			} else if seenPoint {
				intDigits-- // a leading zero after the point
			}
		case c == '.' && !seenPoint:
			seenPoint = true
		default:
			goto digitsDone
		}
		i++
	}
digitsDone:
	if digitCount == 0 || i == len(s) || (s[i] != 'p' && s[i] != 'P') {
		return 0, parseDoubleInvalid
	}
	exponent, end, ok := scanSignedExponent(s, i+1)
	if !ok || !skipFloatSuffix(s, end) {
		return 0, parseDoubleInvalid
	}
	if len(sig) == 0 {
		return 0, parseDoubleOK
	}
	// The value is 0x0.<sig> × 2^scale2, with 1/16 <= 0x0.<sig> < 1.
	scale2 := 4*intDigits + exponent
	switch {
	case scale2 > 1100:
		return math.Inf(1), parseDoubleOK
	case scale2 < -1200:
		return 0, parseDoubleOK
	}
	v, _ := strconv.ParseFloat("0x0."+string(sig)+"p"+strconv.FormatInt(scale2, 10), 64)
	return v, parseDoubleOK
}

func isRangeError(err error) bool {
	ne, ok := err.(*strconv.NumError)
	return ok && ne.Err == strconv.ErrRange
}
