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

// bigIntegerReportOverflow throws the ArithmeticException of a result
// outside the supported range.
//
// Observed behaviour (black-box): the message is "BigInteger would overflow
// supported range".
func bigIntegerReportOverflow() {
	panic(NewArithmeticException("BigInteger would overflow supported range"))
}

// bigIntegerCheckRange returns x, or throws if x is outside the supported
// range (see bigIntegerMaxBitLength).
func bigIntegerCheckRange(x *big.Int) *big.Int {
	if x.BitLen() > bigIntegerMaxBitLength {
		bigIntegerReportOverflow()
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
			bigIntegerReportOverflow()
		}
		result = new(big.Int).Lsh(big.NewInt(1), uint(bits))
	} else {
		bigIntegerPowCheckEstimate(magnitude, exponent)
		result = bigIntegerCheckRange(new(big.Int).Exp(magnitude, big.NewInt(int64(exponent)), nil))
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
		bigIntegerReportOverflow()
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
		bigIntegerReportOverflow()
	}
	return bigIntegerCheckRange(new(big.Int).Mul(a, b))
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
		bigIntegerReportOverflow()
	}
	return new(big.Int).Lsh(x, uint(n))
}
