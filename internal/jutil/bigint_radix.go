package jutil

import (
	"math"
	"math/big"
)

// Go-only: language-adapter — new BigInteger(String val, int radix), written
// from its Javadoc: an optional leading '-' or '+', then one or more digits
// in radix (any char Character.digit accepts, so non-ASCII decimal digits
// count); NumberFormatException otherwise or for a radix outside
// [Character.MIN_RADIX, Character.MAX_RADIX]. The messages were determined
// black-box against JDK 21 (vectors from tools/oracle/java/
// JutilParsingVectors.java and JutilNumbersVectors.java).

// NewBigIntegerFromStringRadix is new BigInteger(String val, int radix).
//
// Observed behaviour (black-box), in this order:
//   - a bad radix → "Radix out of range";
//   - "" or a sign alone → "Zero length BigInteger";
//   - a '-' or '+' anywhere after the first char → "Illegal embedded sign
//     character";
//   - a char that is not a digit → the message of Integer.parseInt (see
//     numberFormatForInputString) for the group of digits that contains it.
//     After the leading zeros, the digits are split into groups of
//     bigIntegerRadixGroup(radix) chars, the first group taking the
//     remainder (or a whole group), and the first bad group is named.
func NewBigIntegerFromStringRadix(s string, radix int) *big.Int {
	if radix < Character_MIN_RADIX || radix > Character_MAX_RADIX {
		panic(NewNumberFormatException("Radix out of range"))
	}
	units := StringToUTF16(s)
	start := 0
	negative := false
	if len(units) > 0 && (units[0] == '-' || units[0] == '+') {
		negative = units[0] == '-'
		start = 1
	}
	if start == len(units) {
		panic(NewNumberFormatException("Zero length BigInteger"))
	}
	for _, c := range units[start:] {
		if c == '-' || c == '+' {
			panic(NewNumberFormatException("Illegal embedded sign character"))
		}
	}
	for start < len(units) && CharacterDigit(units[start], radix) == 0 {
		start++
	}
	if start == len(units) {
		return big.NewInt(0)
	}
	digits := units[start:]
	group := bigIntegerRadixGroup(radix)
	first := len(digits) % group
	if first == 0 {
		first = group
	}
	text := make([]byte, len(digits))
	for from, to := 0, first; from < len(digits); from, to = to, to+group {
		for i := from; i < to; i++ {
			d := CharacterDigit(digits[i], radix)
			if d < 0 {
				panic(numberFormatForInputString(StringFromUTF16(digits[from:to]), radix))
			}
			text[i] = "0123456789abcdefghijklmnopqrstuvwxyz"[d]
		}
	}
	x, _ := new(big.Int).SetString(string(text), radix)
	if negative {
		x.Neg(x)
	}
	return bigIntegerCheckRange(x)
}

// bigIntegerRadixGroup is the size of the digit groups in which
// NewBigIntegerFromStringRadix reports a bad digit: the largest k such that
// radix^k does not exceed Integer.MAX_VALUE, so any k digits fit in an int.
//
// Observed behaviour (black-box): measured for every radix from 2 to 36 (30
// for radix 2, 9 for 10, 7 for 16, 5 for 36); each matches this rule.
func bigIntegerRadixGroup(radix int) int {
	k := 0
	for p := int64(radix); p <= math.MaxInt32; p *= int64(radix) {
		k++
	}
	return k
}
