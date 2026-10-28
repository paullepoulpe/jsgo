package jutil

import (
	"math"
	"math/big"
	"strconv"
	"unicode/utf8"
)

// Go-only: language-adapter — the integer parsers of java.lang.Long and
// java.lang.Integer, and java.math.BigDecimal's String constructor with the
// accessors Gson needs. Written from the Javadoc of Long.parseLong(String,
// int), Integer.parseInt(String, int) and BigDecimal(String); the exception
// messages, which the Javadoc does not give, were determined black-box
// against JDK 21 (tools/oracle/java/JutilNumbersVectors.java).

// LongParseLong is Long.parseLong(String s, int radix).
//
// Per the Javadoc: an optional leading '-' or '+' (only when more characters
// follow), then one or more characters that Character.digit accepts in
// radix; the value must fit in a long.
func LongParseLong(s string, radix int) int64 {
	checkParseRadix(radix)
	v, ok := parseRadixInteger(s, radix, math.MaxInt64)
	if !ok {
		panic(badNumberInput(s, radix))
	}
	return v
}

// TryParseLong parses s in radix using Long.parseLong's grammar and range,
// returning (0, false) where Long.parseLong would throw NumberFormatException.
func TryParseLong(s string, radix int) (int64, bool) {
	if radix < Character_MIN_RADIX || radix > Character_MAX_RADIX {
		return 0, false
	}
	return parseRadixInteger(s, radix, math.MaxInt64)
}

// LongToString is Long.toString(long): the signed decimal form.
func LongToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

// checkParseRadix rejects a radix outside [Character.MIN_RADIX,
// Character.MAX_RADIX], as Integer.parseInt and Long.parseLong do before
// looking at the string.
//
// Observed behaviour (black-box): the radix is checked first, even for an
// empty string, with these two messages.
func checkParseRadix(radix int) {
	if radix < Character_MIN_RADIX {
		panic(NewNumberFormatException("radix " + itoa(radix) + " less than Character.MIN_RADIX"))
	}
	if radix > Character_MAX_RADIX {
		panic(NewNumberFormatException("radix " + itoa(radix) + " greater than Character.MAX_RADIX"))
	}
}

// parseRadixInteger reads s as Integer.parseInt/Long.parseLong do: an
// optional sign (not on its own), then digits in radix (any char that
// Character.digit maps to a value below radix, so non-ASCII decimal digits
// count). A positive result may be at most maxPositive and a negative one
// at least -maxPositive-1. ok is false when s is not such a number.
//
// The digits are Java chars. A WTF-8 sequence of a supplementary code point
// stands for two surrogate chars, and surrogates are never digits, so such
// a sequence (or any undecodable byte) makes s invalid.
func parseRadixInteger(s string, radix int, maxPositive uint64) (int64, bool) {
	i := 0
	negative := false
	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		negative = s[0] == '-'
		i = 1
	}
	if i == len(s) {
		return 0, false // empty, or a sign alone
	}
	limit := maxPositive
	if negative {
		limit++
	}
	base := uint64(radix)
	var magnitude uint64
	for i < len(s) {
		var digit int
		if c := s[i]; c < utf8.RuneSelf {
			digit = CharacterDigit(uint16(c), radix)
			i++
		} else {
			r, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 || r > 0xFFFF {
				return 0, false
			}
			digit = CharacterDigit(uint16(r), radix)
			i += size
		}
		if digit < 0 {
			return 0, false
		}
		if magnitude > (limit-uint64(digit))/base {
			return 0, false // the value would leave the type's range
		}
		magnitude = magnitude*base + uint64(digit)
	}
	if negative {
		return -int64(magnitude), true // -2^63 wraps to itself, as wanted
	}
	return int64(magnitude), true
}

// BigDecimal is the part of java.math.BigDecimal that ported code uses: an
// immutable unscaled integer and a 32-bit scale, standing for
// unscaled × 10^-scale.
type BigDecimal struct {
	unscaled *big.Int
	scale    int32
}

// bigDecimalLongFormLength is the string length (after an optional leading
// sign) above which BigDecimal(String) reports any unexpected character with
// the "missing e notation" message instead of naming the character.
//
// Observed behaviour (black-box): with at most 18 chars after the sign,
// "1x" fails with `Character x is neither a decimal digit number, decimal
// point, nor "e" notation exponential mark.`; with 19 or more, every such
// failure says `Character array is missing "e" notation exponential mark.`
// All other messages are the same for both lengths.
const bigDecimalLongFormLength = 18

// NewBigDecimalFromString is new BigDecimal(String).
//
// The grammar, from the Javadoc: an optional sign, then a significand (digits
// with at most one decimal point, at least one digit), then optionally 'e'
// or 'E' and a signed integer exponent. Digits are any chars that
// Character.isDigit accepts. The unscaled value is the significand's digits
// read as an integer; the scale is the number of fraction digits minus the
// exponent, and must fit in an int.
//
// Observed behaviour (black-box): the string is read left to right and the
// first problem found decides the message:
//   - "" → NumberFormatException with no message;
//   - a second '.' → "Character array contains more than one decimal point.";
//   - another non-digit before the exponent → see bigDecimalLongFormLength;
//   - nothing after 'e', or after the exponent's sign → no message;
//   - more than 10 exponent chars left after the leading zeros (of any
//     script) → "Too many nonzero exponent digits." (checked before the
//     chars are looked at);
//   - a non-digit in the exponent → "Not a digit.";
//   - after the whole string: no significand digit → "No digits found.";
//     then a scale outside the int range → "Exponent overflow.".
func NewBigDecimalFromString(s string) *BigDecimal {
	units := StringToUTF16(s)
	n := len(units)
	if n == 0 {
		panic(newNumberFormatExceptionWithoutMessage())
	}
	pos := 0
	negative := false
	if units[0] == '-' || units[0] == '+' {
		negative = units[0] == '-'
		pos = 1
	}
	longForm := n-pos > bigDecimalLongFormLength
	digits := make([]byte, 0, n)
	fractionDigits := int64(0)
	seenPoint := false
	exponent := int64(0)
scan:
	for ; pos < n; pos++ {
		c := units[pos]
		if d := CharacterDigit(c, 10); d >= 0 {
			digits = append(digits, byte('0'+d))
			if seenPoint {
				fractionDigits++
			}
			continue
		}
		switch c {
		case '.':
			if seenPoint {
				panic(NewNumberFormatException("Character array contains more than one decimal point."))
			}
			seenPoint = true
		case 'e', 'E':
			exponent = parseBigDecimalExponent(units[pos+1:])
			break scan
		default:
			if longForm {
				panic(NewNumberFormatException("Character array is missing \"e\" notation exponential mark."))
			}
			panic(NewNumberFormatException("Character " + StringFromUTF16(units[pos:pos+1]) +
				" is neither a decimal digit number, decimal point, nor \"e\" notation exponential mark."))
		}
	}
	if len(digits) == 0 {
		panic(NewNumberFormatException("No digits found."))
	}
	scale := fractionDigits - exponent
	if scale < math.MinInt32 || scale > math.MaxInt32 {
		panic(NewNumberFormatException("Exponent overflow."))
	}
	unscaled, _ := new(big.Int).SetString(string(digits), 10)
	if negative {
		unscaled.Neg(unscaled)
	}
	return &BigDecimal{unscaled: unscaled, scale: int32(scale)}
}

// parseBigDecimalExponent reads the chars after 'e' in BigDecimal(String):
// an optional sign and decimal digits (see NewBigDecimalFromString for the
// failure messages). At most 10 significant digits are read, so the value
// always fits in an int64.
func parseBigDecimalExponent(units []uint16) int64 {
	if len(units) == 0 {
		panic(newNumberFormatExceptionWithoutMessage())
	}
	negative := false
	if units[0] == '-' || units[0] == '+' {
		negative = units[0] == '-'
		units = units[1:]
		if len(units) == 0 {
			panic(newNumberFormatExceptionWithoutMessage())
		}
	}
	for len(units) > 0 && CharacterDigit(units[0], 10) == 0 {
		units = units[1:] // any zero digit, e.g. U+0660, counts as leading
	}
	if len(units) > 10 {
		panic(NewNumberFormatException("Too many nonzero exponent digits."))
	}
	var v int64
	for _, c := range units {
		d := CharacterDigit(c, 10)
		if d < 0 {
			panic(NewNumberFormatException("Not a digit."))
		}
		v = v*10 + int64(d)
	}
	if negative {
		return -v
	}
	return v
}

func newNumberFormatExceptionWithoutMessage() *NumberFormatException {
	e := &NumberFormatException{}
	InitThrowable(e, "java.lang.NumberFormatException", nil, nil)
	return e
}

// Signum is BigDecimal.signum(): -1, 0 or 1.
func (d *BigDecimal) Signum() int {
	return d.unscaled.Sign()
}

// Scale is BigDecimal.scale().
func (d *BigDecimal) Scale() int32 {
	return d.scale
}

// Precision is BigDecimal.precision(): the number of decimal digits of the
// unscaled value (1 for zero).
func (d *BigDecimal) Precision() int {
	if d.unscaled.Sign() == 0 {
		return 1
	}
	return len(new(big.Int).Abs(d.unscaled).String())
}

// ToBigInteger is BigDecimal.toBigInteger(): the value with any fractional
// part discarded (truncation toward zero).
//
// Observed behaviour (black-box): zero gives 0 whatever the scale. For any
// other value the power of ten is subject to BigInteger.pow's range check
// (see BigIntegerPow), so "1e-2147483647" throws "BigInteger would overflow
// supported range" although the result would be 0, and a scale of
// Integer.MIN_VALUE throws ArithmeticException "Underflow".
func (d *BigDecimal) ToBigInteger() *big.Int {
	switch {
	case d.unscaled.Sign() == 0:
		return big.NewInt(0)
	case d.scale == 0:
		return d.unscaled
	case d.scale == math.MinInt32:
		panic(NewArithmeticException("Underflow"))
	case d.scale < 0:
		return BigIntegerMultiply(d.unscaled, BigIntegerPow(BigInteger_TEN, -d.scale))
	}
	if int64(d.scale) >= int64(d.Precision()) {
		// |value| < 1. Only the overflow check of the divisor is observable.
		bigIntegerPowCheckEstimate(BigInteger_TEN, d.scale)
		return big.NewInt(0)
	}
	return new(big.Int).Quo(d.unscaled, BigIntegerPow(BigInteger_TEN, d.scale))
}

// LongValue is BigDecimal.longValue(): the low-order 64 bits of the value
// with its fractional part discarded.
func (d *BigDecimal) LongValue() int64 {
	return int64(d.truncatedLow64())
}

// IntValue is BigDecimal.intValue(): the low-order 32 bits of the value with
// its fractional part discarded.
func (d *BigDecimal) IntValue() int32 {
	return int32(uint32(d.truncatedLow64()))
}

// truncatedLow64 returns the integer part of d modulo 2^64. It never builds
// a huge power of ten: 10^k is a multiple of 2^64 for k >= 64, and a scale
// of at least the precision means |d| < 1.
func (d *BigDecimal) truncatedLow64() uint64 {
	if d.unscaled.Sign() == 0 {
		return 0
	}
	var whole *big.Int
	switch {
	case d.scale <= 0:
		if d.scale <= -64 {
			return 0
		}
		whole = new(big.Int).Mul(d.unscaled, new(big.Int).Exp(BigInteger_TEN, big.NewInt(-int64(d.scale)), nil))
	case int64(d.scale) >= int64(d.Precision()):
		return 0
	default:
		whole = new(big.Int).Quo(d.unscaled, new(big.Int).Exp(BigInteger_TEN, big.NewInt(int64(d.scale)), nil))
	}
	low := new(big.Int).And(new(big.Int).Abs(whole), new(big.Int).SetUint64(math.MaxUint64)).Uint64()
	if whole.Sign() < 0 {
		return -low // two's complement of the magnitude's low bits
	}
	return low
}
