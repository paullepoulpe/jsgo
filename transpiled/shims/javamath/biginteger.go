// Clean-room specification based on java.math.BigInteger Javadoc.

package javamath

import (
	"math"
	"math/big"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

const bigIntegerMaxBitLength = math.MaxInt32

// BigInteger represents an immutable arbitrary-precision signed integer (java.math.BigInteger).
type BigInteger struct {
	v *big.Int
}

var (
	BigInteger_ZERO = &BigInteger{v: big.NewInt(0)}
	BigInteger_ONE  = &BigInteger{v: big.NewInt(1)}
	BigInteger_TWO  = &BigInteger{v: big.NewInt(2)}
	BigInteger_TEN  = &BigInteger{v: big.NewInt(10)}
)

func wrapBigInt(v *big.Int) *BigInteger {
	if v.BitLen() > bigIntegerMaxBitLength {
		panic(javalang.NewArithmeticException("BigInteger would overflow supported range"))
	}
	return &BigInteger{v: v}
}

// BigInteger_ValueOf implements BigInteger.valueOf(long).
func BigInteger_ValueOf(val int64) *BigInteger {
	switch val {
	case 0:
		return BigInteger_ZERO
	case 1:
		return BigInteger_ONE
	case 2:
		return BigInteger_TWO
	case 10:
		return BigInteger_TEN
	default:
		return &BigInteger{v: big.NewInt(val)}
	}
}

// NewBigInteger implements new BigInteger(String) (radix 10).
func NewBigInteger(s string) *BigInteger {
	return NewBigIntegerRadix(s, 10)
}

// NewBigIntegerRadix implements new BigInteger(String, int).
func NewBigIntegerRadix(s string, radix int) *BigInteger {
	if radix < javalang.CharacterMinRadix || radix > javalang.CharacterMaxRadix {
		panic(javalang.NewNumberFormatException("Radix out of range"))
	}
	units := javalang.StringToUTF16(s)
	start := 0
	neg := false
	if len(units) > 0 && (units[0] == '-' || units[0] == '+') {
		neg = units[0] == '-'
		start = 1
	}
	if start == len(units) {
		panic(javalang.NewNumberFormatException("Zero length BigInteger"))
	}
	for _, c := range units[start:] {
		if c == '-' || c == '+' {
			panic(javalang.NewNumberFormatException("Illegal embedded sign character"))
		}
	}
	for start < len(units) && javalang.CharacterDigit(units[start], radix) == 0 {
		start++
	}
	if start == len(units) {
		return BigInteger_ZERO
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
			d := javalang.CharacterDigit(digits[i], radix)
			if d < 0 {
				panic(javalang.BadNumberInput(javalang.StringFromUTF16(digits[from:to]), radix))
			}
			text[i] = "0123456789abcdefghijklmnopqrstuvwxyz"[d]
		}
	}
	x, _ := new(big.Int).SetString(string(text), radix)
	if neg {
		x.Neg(x)
	}
	return wrapBigInt(x)
}

func bigIntegerRadixGroup(radix int) int {
	k := 0
	for p := int64(radix); p <= math.MaxInt32; p *= int64(radix) {
		k++
	}
	return k
}

// NewBigIntegerBytes implements new BigInteger(byte[]) (big-endian two's-complement).
func NewBigIntegerBytes(b []byte) *BigInteger {
	if len(b) == 0 {
		panic(javalang.NewNumberFormatException("Zero length BigInteger"))
	}
	x := new(big.Int).SetBytes(b)
	if b[0]&0x80 != 0 {
		// Negative two's-complement value: subtract 2^(8*len(b)).
		twoPow := new(big.Int).Lsh(big.NewInt(1), uint(8*len(b)))
		x.Sub(x, twoPow)
	}
	return wrapBigInt(x)
}

// Raw returns a copy of the underlying *big.Int.
func (b *BigInteger) Raw() *big.Int {
	if b == nil || b.v == nil {
		return nil
	}
	return new(big.Int).Set(b.v)
}

func (b *BigInteger) Signum() int {
	return b.v.Sign()
}

func (b *BigInteger) BitLength() int {
	if b.v.Sign() >= 0 {
		return b.v.BitLen()
	}
	// Java BitLength for negative x is bitLength of (~x) = -x - 1.
	t := new(big.Int).Neg(b.v)
	t.Sub(t, big.NewInt(1))
	return t.BitLen()
}

func (b *BigInteger) Abs() *BigInteger {
	if b.v.Sign() >= 0 {
		return b
	}
	return wrapBigInt(new(big.Int).Abs(b.v))
}

func (b *BigInteger) Negate() *BigInteger {
	if b.v.Sign() == 0 {
		return BigInteger_ZERO
	}
	return &BigInteger{v: new(big.Int).Neg(b.v)}
}

func (b *BigInteger) Add(other *BigInteger) *BigInteger {
	return wrapBigInt(new(big.Int).Add(b.v, other.v))
}

func (b *BigInteger) Subtract(other *BigInteger) *BigInteger {
	return wrapBigInt(new(big.Int).Sub(b.v, other.v))
}

func (b *BigInteger) Multiply(other *BigInteger) *BigInteger {
	if b.v.Sign() == 0 || other.v.Sign() == 0 {
		return BigInteger_ZERO
	}
	if int64(b.v.BitLen())+int64(other.v.BitLen())-1 > bigIntegerMaxBitLength {
		panic(javalang.NewArithmeticException("BigInteger would overflow supported range"))
	}
	return wrapBigInt(new(big.Int).Mul(b.v, other.v))
}

// Divide implements BigInteger.divide(BigInteger) (truncates toward zero).
func (b *BigInteger) Divide(other *BigInteger) *BigInteger {
	if other.v.Sign() == 0 {
		panic(javalang.NewArithmeticException("BigInteger divide by zero"))
	}
	return &BigInteger{v: new(big.Int).Quo(b.v, other.v)}
}

// Remainder implements BigInteger.remainder(BigInteger) (sign matches dividend).
func (b *BigInteger) Remainder(other *BigInteger) *BigInteger {
	if other.v.Sign() == 0 {
		panic(javalang.NewArithmeticException("BigInteger divide by zero"))
	}
	return &BigInteger{v: new(big.Int).Rem(b.v, other.v)}
}

// DivideAndRemainder implements BigInteger.divideAndRemainder(BigInteger).
func (b *BigInteger) DivideAndRemainder(other *BigInteger) []*BigInteger {
	if other.v.Sign() == 0 {
		panic(javalang.NewArithmeticException("BigInteger divide by zero"))
	}
	q, r := new(big.Int).QuoRem(b.v, other.v, new(big.Int))
	return []*BigInteger{{v: q}, {v: r}}
}

// Mod implements BigInteger.mod(BigInteger) (non-negative remainder for m > 0).
func (b *BigInteger) Mod(m *BigInteger) *BigInteger {
	if m.v.Sign() <= 0 {
		panic(javalang.NewArithmeticException("BigInteger: modulus not positive"))
	}
	return &BigInteger{v: new(big.Int).Mod(b.v, m.v)}
}

// Pow implements BigInteger.pow(int).
func (b *BigInteger) Pow(exponent int) *BigInteger {
	if exponent < 0 {
		panic(javalang.NewArithmeticException("Negative exponent"))
	}
	if exponent == 0 {
		return BigInteger_ONE
	}
	if b.v.Sign() == 0 {
		return BigInteger_ZERO
	}
	neg := b.v.Sign() < 0 && exponent&1 == 1
	mag := new(big.Int).Abs(b.v)
	var res *big.Int
	if twos := mag.TrailingZeroBits(); uint(mag.BitLen()) == twos+1 {
		nbits := int64(twos) * int64(exponent)
		if nbits+1 > bigIntegerMaxBitLength {
			panic(javalang.NewArithmeticException("BigInteger would overflow supported range"))
		}
		res = new(big.Int).Lsh(big.NewInt(1), uint(nbits))
	} else {
		if int64(mag.BitLen())*int64(exponent)/32 > 1<<26 {
			panic(javalang.NewArithmeticException("BigInteger would overflow supported range"))
		}
		res = new(big.Int).Exp(mag, big.NewInt(int64(exponent)), nil)
	}
	if neg {
		res.Neg(res)
	}
	return wrapBigInt(res)
}

func (b *BigInteger) And(other *BigInteger) *BigInteger {
	return &BigInteger{v: new(big.Int).And(b.v, other.v)}
}

func (b *BigInteger) Or(other *BigInteger) *BigInteger {
	return &BigInteger{v: new(big.Int).Or(b.v, other.v)}
}

func (b *BigInteger) Xor(other *BigInteger) *BigInteger {
	return &BigInteger{v: new(big.Int).Xor(b.v, other.v)}
}

func (b *BigInteger) Not() *BigInteger {
	return &BigInteger{v: new(big.Int).Not(b.v)}
}

func (b *BigInteger) ShiftLeft(n int) *BigInteger {
	if n == 0 || b.v.Sign() == 0 {
		return b
	}
	if n < 0 {
		return &BigInteger{v: new(big.Int).Rsh(b.v, uint(-int64(n)))}
	}
	if int64(b.v.BitLen())+int64(n) > bigIntegerMaxBitLength {
		panic(javalang.NewArithmeticException("BigInteger would overflow supported range"))
	}
	return wrapBigInt(new(big.Int).Lsh(b.v, uint(n)))
}

func (b *BigInteger) ShiftRight(n int) *BigInteger {
	if n == 0 || b.v.Sign() == 0 {
		return b
	}
	if n < 0 {
		return b.ShiftLeft(-n)
	}
	return &BigInteger{v: new(big.Int).Rsh(b.v, uint(n))}
}

// IntValue returns the low-order 32 bits as a signed int32.
func (b *BigInteger) IntValue() int32 {
	words := b.v.Bits()
	if len(words) == 0 {
		return 0
	}
	low := uint32(words[0])
	if b.v.Sign() < 0 {
		return int32(-low)
	}
	return int32(low)
}

// IntValueExact converts b to int32, panicking with ArithmeticException if out of range.
func (b *BigInteger) IntValueExact() int32 {
	if !b.v.IsInt64() {
		panic(javalang.NewArithmeticException("BigInteger out of int range"))
	}
	i64 := b.v.Int64()
	if i64 < math.MinInt32 || i64 > math.MaxInt32 {
		panic(javalang.NewArithmeticException("BigInteger out of int range"))
	}
	return int32(i64)
}

// LongValue returns the low-order 64 bits as a signed int64.
func (b *BigInteger) LongValue() int64 {
	words := b.v.Bits()
	if len(words) == 0 {
		return 0
	}
	var low uint64
	if math.MaxUint == math.MaxUint64 {
		low = uint64(words[0])
	} else {
		low = uint64(words[0])
		if len(words) > 1 {
			low |= uint64(words[1]) << 32
		}
	}
	if b.v.Sign() < 0 {
		return int64(-low)
	}
	return int64(low)
}

// LongValueExact converts b to int64, panicking with ArithmeticException if out of range.
func (b *BigInteger) LongValueExact() int64 {
	if !b.v.IsInt64() {
		panic(javalang.NewArithmeticException("BigInteger out of long range"))
	}
	return b.v.Int64()
}

// DoubleValue converts b to an IEEE-754 float64 using round-to-nearest-even.
func (b *BigInteger) DoubleValue() float64 {
	f, _ := new(big.Float).SetInt(b.v).Float64()
	return f
}

// ToByteArray returns the minimal big-endian two's-complement byte slice.
func (b *BigInteger) ToByteArray() []byte {
	n := b.v.BitLen()/8 + 1
	out := make([]byte, n)
	if b.v.Sign() >= 0 {
		b.v.FillBytes(out)
		return out
	}
	t := new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), uint(8*n)), b.v)
	t.FillBytes(out)
	i := 0
	for i < len(out)-1 && out[i] == 0xff && out[i+1]&0x80 != 0 {
		i++
	}
	return out[i:]
}

func (b *BigInteger) CompareTo(other *BigInteger) int {
	return b.v.Cmp(other.v)
}

func (b *BigInteger) Equals(other any) bool {
	o, ok := other.(*BigInteger)
	if !ok || o == nil || b == nil {
		return false
	}
	return b.v.Cmp(o.v) == 0
}

// HashCode implements Java BigInteger.hashCode() (32-bit unsigned big-endian limbs with sign).
func (b *BigInteger) HashCode() int32 {
	if b == nil || b.v.Sign() == 0 {
		return 0
	}
	mag := new(big.Int).Abs(b.v).Bytes()
	// Split big-endian bytes into 32-bit big-endian ints.
	rem := len(mag) % 4
	var hashCode int32
	idx := 0
	if rem != 0 {
		var limb uint32
		for idx < rem {
			limb = (limb << 8) | uint32(mag[idx])
			idx++
		}
		hashCode = int32(limb)
	}
	for idx < len(mag) {
		limb := uint32(mag[idx])<<24 | uint32(mag[idx+1])<<16 | uint32(mag[idx+2])<<8 | uint32(mag[idx+3])
		hashCode = 31*hashCode + int32(limb)
		idx += 4
	}
	if b.v.Sign() < 0 {
		return -hashCode
	}
	return hashCode
}

func (b *BigInteger) String() string {
	if b == nil || b.v == nil {
		return "null"
	}
	return b.v.Text(10)
}

func (b *BigInteger) ToStringRadix(radix int) string {
	if radix < javalang.CharacterMinRadix || radix > javalang.CharacterMaxRadix {
		radix = 10
	}
	return b.v.Text(radix)
}
