package jutil

import (
	"math"
	"math/big"
)

// Go-only: language-adapter — java.math.BigInteger on math/big.Int. A Java
// BigInteger is an immutable value; ported code treats *big.Int values the
// same way (never mutates one it did not just allocate).

// NewBigIntegerFromBytes is new BigInteger(byte[]): a big-endian two's
// complement representation.
func NewBigIntegerFromBytes(b []byte) *big.Int {
	x := new(big.Int).SetBytes(b)
	if len(b) > 0 && b[0]&0x80 != 0 {
		x.Sub(x, new(big.Int).Lsh(big.NewInt(1), uint(8*len(b))))
	}
	return x
}

// BigIntegerToByteArray is BigInteger.toByteArray(): the minimal big-endian
// two's complement representation (always at least one byte, with a sign
// bit).
func BigIntegerToByteArray(x *big.Int) []byte {
	n := x.BitLen()/8 + 1 // room for the sign bit
	out := make([]byte, n)
	if x.Sign() >= 0 {
		x.FillBytes(out)
		return out
	}
	// Two's complement of a negative number: 2^(8n) + x.
	t := new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), uint(8*n)), x)
	t.FillBytes(out)
	// Trim redundant leading 0xff bytes, keeping the sign bit set.
	i := 0
	for i < len(out)-1 && out[i] == 0xff && out[i+1]&0x80 != 0 {
		i++
	}
	return out[i:]
}

// NewBigIntegerFromString is new BigInteger(String) (radix 10).
func NewBigIntegerFromString(s string) *big.Int {
	if s == "" {
		panic(NewNumberFormatException("Zero length BigInteger"))
	}
	body := s
	if body[0] == '-' || body[0] == '+' {
		body = body[1:]
	}
	if body == "" {
		panic(NewNumberFormatException("Zero length BigInteger"))
	}
	for i := 0; i < len(body); i++ {
		if body[i] < '0' || body[i] > '9' {
			panic(NewNumberFormatException("For input string: \"" + s + "\""))
		}
	}
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic(NewNumberFormatException("For input string: \"" + s + "\""))
	}
	return x
}

// BigIntegerIntValue is BigInteger.intValue(): the low-order 32 bits.
func BigIntegerIntValue(x *big.Int) int32 {
	m := new(big.Int).And(x, big.NewInt(0xffffffff))
	return int32(uint32(m.Uint64()))
}

// BigIntegerEquals is BigInteger.equals(Object) for two non-null values.
func BigIntegerEquals(a, b *big.Int) bool { return a.Cmp(b) == 0 }

// BigInteger_ZERO, BigInteger_ONE and BigInteger_TEN are the BigInteger
// constants. Callers must not mutate them.
var (
	BigInteger_ZERO = big.NewInt(0)
	BigInteger_ONE  = big.NewInt(1)
	BigInteger_TEN  = big.NewInt(10)
)

// bigIntegerMaxBitLength is the largest supported magnitude, in bits. The
// BigInteger class Javadoc (implementation note) gives the supported range
// as -2^Integer.MAX_VALUE to +2^Integer.MAX_VALUE, both exclusive, and
// operations whose result leaves it throw ArithmeticException.
const bigIntegerMaxBitLength = math.MaxInt32

// panicBigIntegerRange raises the ArithmeticException that BigInteger
// operations throw for a result they cannot represent. Its message,
// "BigInteger would overflow supported range", was observed black-box.
func panicBigIntegerRange() {
	panic(NewArithmeticException("BigInteger would overflow supported range"))
}

// bigIntegerWithinRange passes x through when its magnitude has at most
// bigIntegerMaxBitLength bits, and calls panicBigIntegerRange otherwise.
func bigIntegerWithinRange(x *big.Int) *big.Int {
	if x.BitLen() > bigIntegerMaxBitLength {
		panicBigIntegerRange()
	}
	return x
}

// BigIntegerPow is BigInteger.pow(int exponent): x^exponent, with 0^0 = 1,
// throwing ArithmeticException for a negative exponent (Javadoc).
//
// Observed behaviour (black-box): unless |x| is a power of two, pow throws
// "BigInteger would overflow supported range" straight away when
// bitLength(|x|) * exponent / 32 (integer division) exceeds 2^26, even when
// the exact result would be in range (5^800000000 has about 1.86e9 bits,
// fewer than 2^31-1, and throws). Probing bases 3, 5, -5, 10, 12, 255 and
// 2^40+1 at the boundary: 3^1073741839 starts computing, 3^1073741840
// throws; 10^536870919 computes, 10^536870920 throws. Powers of two are not
// estimated: they throw only when the exact result is out of range
// (2^2147483646 succeeds, 2^1073741840 succeeds).
func BigIntegerPow(x *big.Int, exponent int32) *big.Int {
	if exponent < 0 {
		panic(NewArithmeticException("Negative exponent"))
	}
	if exponent == 0 {
		return big.NewInt(1)
	}
	if x.Sign() == 0 {
		return big.NewInt(0)
	}
	negative := x.Sign() < 0 && exponent&1 == 1
	magnitude := new(big.Int).Abs(x)
	var result *big.Int
	if twos := magnitude.TrailingZeroBits(); uint(magnitude.BitLen()) == twos+1 {
		// |x| = 2^twos, so the result is 2^(twos*exponent).
		bits := int64(twos) * int64(exponent)
		if bits+1 > bigIntegerMaxBitLength {
			panicBigIntegerRange()
		}
		result = new(big.Int).Lsh(big.NewInt(1), uint(bits))
	} else {
		bigIntegerPowCheckEstimate(magnitude, exponent)
		result = bigIntegerWithinRange(new(big.Int).Exp(magnitude, big.NewInt(int64(exponent)), nil))
	}
	if negative {
		result.Neg(result)
	}
	return result
}

// bigIntegerPowCheckEstimate throws when BigInteger.pow rejects
// magnitude^exponent up front (see BigIntegerPow) for a magnitude that is
// not a power of two.
func bigIntegerPowCheckEstimate(magnitude *big.Int, exponent int32) {
	if int64(magnitude.BitLen())*int64(exponent)/32 > 1<<26 {
		panicBigIntegerRange()
	}
}

// BigIntegerMultiply is BigInteger.multiply(BigInteger), including the
// overflow check for a result outside the supported range.
func BigIntegerMultiply(a, b *big.Int) *big.Int {
	if a.Sign() == 0 || b.Sign() == 0 {
		return big.NewInt(0)
	}
	// The product of an m-bit and an n-bit magnitude has at least m+n-1 bits.
	if int64(a.BitLen())+int64(b.BitLen())-1 > bigIntegerMaxBitLength {
		panicBigIntegerRange()
	}
	return bigIntegerWithinRange(new(big.Int).Mul(a, b))
}

// BigIntegerShiftLeft is BigInteger.shiftLeft(int) for n >= 0 (returning x
// itself for n == 0, like the JDK).
func BigIntegerShiftLeft(x *big.Int, n int32) *big.Int {
	if n == 0 || x.Sign() == 0 {
		return x
	}
	if n < 0 {
		return new(big.Int).Rsh(x, uint(-int64(n)))
	}
	if int64(x.BitLen())+int64(n) > bigIntegerMaxBitLength {
		panicBigIntegerRange()
	}
	return new(big.Int).Lsh(x, uint(n))
}

// bigIntegerBitLength is BigInteger.bitLength(): the number of bits of the
// minimal two's-complement representation, excluding the sign bit.
func bigIntegerBitLength(x *big.Int) int {
	if x.Sign() >= 0 {
		return x.BitLen()
	}
	return new(big.Int).Sub(new(big.Int).Neg(x), BigInteger_ONE).BitLen()
}

// BigIntegerIntValueExact is BigInteger.intValueExact(): the value as an
// int, or an ArithmeticException if it does not fit.
func BigIntegerIntValueExact(x *big.Int) int32 {
	if bigIntegerBitLength(x) <= 31 {
		return int32(x.Int64())
	}
	panic(NewArithmeticException("BigInteger out of int range"))
}

// BigIntegerMod is BigInteger.mod(BigInteger): the non-negative remainder
// for a positive modulus, and an ArithmeticException otherwise.
func BigIntegerMod(x, m *big.Int) *big.Int {
	if m.Sign() <= 0 {
		panic(NewArithmeticException("BigInteger: modulus not positive"))
	}
	return new(big.Int).Mod(x, m)
}

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
//     badNumberInput) for the group of digits that contains it.
//     After the leading zeros, the digits are split into groups of
//     bigIntegerRadixGroup(radix) chars, the first group taking the
//     remainder (or a whole group), and the first bad group is named.
func NewBigIntegerFromStringRadix(s string, radix int) *big.Int {
	if radix < character_MIN_RADIX || radix > character_MAX_RADIX {
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
				panic(badNumberInput(StringFromUTF16(digits[from:to]), radix))
			}
			text[i] = "0123456789abcdefghijklmnopqrstuvwxyz"[d]
		}
	}
	x, _ := new(big.Int).SetString(string(text), radix)
	if negative {
		x.Neg(x)
	}
	return bigIntegerWithinRange(x)
}

// TryParseBigInt parses s in radix using Java's BigInteger(String, int)
// grammar, returning (nil, false) where NewBigIntegerFromStringRadix would
// throw NumberFormatException or ArithmeticException.
func TryParseBigInt(s string, radix int) (*big.Int, bool) {
	if radix < character_MIN_RADIX || radix > character_MAX_RADIX {
		return nil, false
	}
	if IsASCII(s) {
		start := 0
		negative := false
		if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
			negative = s[0] == '-'
			start = 1
		}
		if start == len(s) {
			return nil, false
		}
		for i := start; i < len(s); i++ {
			c := s[i]
			var d int
			switch {
			case c >= '0' && c <= '9':
				d = int(c - '0')
			case c >= 'a' && c <= 'z':
				d = int(c-'a') + 10
			case c >= 'A' && c <= 'Z':
				d = int(c-'A') + 10
			default:
				return nil, false
			}
			if d >= radix {
				return nil, false
			}
		}
		for start < len(s) && s[start] == '0' {
			start++
		}
		if start == len(s) {
			return big.NewInt(0), true
		}
		x, ok := new(big.Int).SetString(s[start:], radix)
		if !ok || x.BitLen() > bigIntegerMaxBitLength {
			return nil, false
		}
		if negative {
			x.Neg(x)
		}
		return x, true
	}
	units := StringToUTF16(s)
	start := 0
	negative := false
	if len(units) > 0 && (units[0] == '-' || units[0] == '+') {
		negative = units[0] == '-'
		start = 1
	}
	if start == len(units) {
		return nil, false
	}
	for _, c := range units[start:] {
		if c == '-' || c == '+' {
			return nil, false
		}
	}
	for start < len(units) && CharacterDigit(units[start], radix) == 0 {
		start++
	}
	if start == len(units) {
		return big.NewInt(0), true
	}
	digits := units[start:]
	text := make([]byte, len(digits))
	for i, c := range digits {
		d := CharacterDigit(c, radix)
		if d < 0 {
			return nil, false
		}
		text[i] = "0123456789abcdefghijklmnopqrstuvwxyz"[d]
	}
	x, ok := new(big.Int).SetString(string(text), radix)
	if !ok || x.BitLen() > bigIntegerMaxBitLength {
		return nil, false
	}
	if negative {
		x.Neg(x)
	}
	return x, true
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
