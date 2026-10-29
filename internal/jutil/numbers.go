package jutil

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"unicode/utf8"
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

// DRem is Java's double remainder `a % b`: the IEEE 754 "fmod" result, with
// the sign of a and a magnitude below |b| (NaN if a is infinite, b is zero
// or either is NaN; a if b is infinite). It is exact, like Go's math.Mod.
func DRem(a, b float64) float64 { return math.Mod(a, b) }

// Shl32 is Java's int `x << n` (shift distance masked to 5 bits).
func Shl32(x int32, n int64) int32 { return x << (uint64(n) & 31) }

// Shr32 is Java's int `x >> n` (arithmetic).
func Shr32(x int32, n int64) int32 { return x >> (uint64(n) & 31) }

// Ushr32 is Java's int `x >>> n` (logical).
func Ushr32(x int32, n int64) int32 { return int32(uint32(x) >> (uint64(n) & 31)) }

// Shl64 is Java's long `x << n` (shift distance masked to 6 bits).
func Shl64(x int64, n int64) int64 { return x << (uint64(n) & 63) }

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
// A pointer to a string, bool or number is a nullable Java String or boxed
// primitive (Integer, Boolean, ...): non-nil, it prints its value.
func ToJavaString(v any) string {
	if IsNil(v) {
		return "null"
	}
	switch x := v.(type) {
	case string:
		return x
	case *string:
		// A nullable Java String (non-null here).
		return *x
	case *bool:
		return ToJavaString(*x)
	case *int:
		return ToJavaString(*x)
	case *int8:
		return ToJavaString(*x)
	case *int16:
		return ToJavaString(*x)
	case *int32:
		return ToJavaString(*x)
	case *int64:
		return ToJavaString(*x)
	case *float32:
		return ToJavaString(*x)
	case *float64:
		return ToJavaString(*x)
	case *uint16:
		return ToJavaString(*x)
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
	case []any:
		return CollectionToString(x)
	case []string:
		return CollectionToString(x)
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

// IntAdd is Java's a + b on ints (wrapping at 32 bits).
func IntAdd(a, b int) int { return int(int32(a) + int32(b)) }

// IntSub is Java's a - b on ints (wrapping at 32 bits).
func IntSub(a, b int) int { return int(int32(a) - int32(b)) }

// Integer is a boxed java.lang.Integer seen as a java.lang.Number (for
// example the value of JsonObject.addProperty(String, Number)). A Go int32
// stands for an Integer elsewhere; this named type is the Number.
//
// Written from the Javadoc of java.lang.Integer and java.lang.Number: the
// conversions are the JLS widening primitive conversions.
type Integer int32

// Java: Integer.intValue()
func (n Integer) IntValue() int32 { return int32(n) }

// Java: Integer.longValue()
func (n Integer) LongValue() int64 { return int64(n) }

// Java: Integer.doubleValue()
func (n Integer) DoubleValue() float64 { return float64(n) }

// Java: Integer.toString()
func (n Integer) String() string { return strconv.FormatInt(int64(n), 10) }

// Java: Integer.hashCode()
func (n Integer) HashCode() int32 { return int32(n) }

// Java: Integer.equals(Object)
func (n Integer) Equals(obj any) bool {
	other, ok := obj.(Integer)
	return ok && other == n
}

// longParseLong is Long.parseLong(String s, int radix).
func longParseLong(s string, radix int) int64 {
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
	if radix < character_MIN_RADIX || radix > character_MAX_RADIX {
		return 0, false
	}
	return parseRadixInteger(s, radix, math.MaxInt64)
}

// LongToString is Long.toString(long): the signed decimal form.
func LongToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

func checkParseRadix(radix int) {
	if radix < character_MIN_RADIX {
		panic(NewNumberFormatException("radix " + itoa(radix) + " less than Character.MIN_RADIX"))
	}
	if radix > character_MAX_RADIX {
		panic(NewNumberFormatException("radix " + itoa(radix) + " greater than Character.MAX_RADIX"))
	}
}

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

// bigDecimal is the part of java.math.BigDecimal that ported code uses: an
// immutable unscaled integer and a 32-bit scale, standing for
// unscaled × 10^-scale.
type bigDecimal struct {
	unscaled *big.Int
	scale    int32
}

const bigDecimalLongFormLength = 18

// newBigDecimalFromString is new BigDecimal(String).
func newBigDecimalFromString(s string) *bigDecimal {
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
	return &bigDecimal{unscaled: unscaled, scale: int32(scale)}
}

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
func (d *bigDecimal) Signum() int {
	return d.unscaled.Sign()
}

// Scale is BigDecimal.scale().
func (d *bigDecimal) Scale() int32 {
	return d.scale
}

// Precision is BigDecimal.precision(): the number of decimal digits of the
// unscaled value (1 for zero).
func (d *bigDecimal) Precision() int {
	if d.unscaled.Sign() == 0 {
		return 1
	}
	return len(new(big.Int).Abs(d.unscaled).String())
}

// ToBigInteger is BigDecimal.toBigInteger(): the value with any fractional
// part discarded (truncation toward zero).
func (d *bigDecimal) ToBigInteger() *big.Int {
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
		bigIntegerPowCheckEstimate(BigInteger_TEN, d.scale)
		return big.NewInt(0)
	}
	return new(big.Int).Quo(d.unscaled, BigIntegerPow(BigInteger_TEN, d.scale))
}

// LongValue is BigDecimal.longValue(): the low-order 64 bits of the value
// with its fractional part discarded.
func (d *bigDecimal) LongValue() int64 {
	return int64(d.truncatedLow64())
}

// IntValue is BigDecimal.intValue(): the low-order 32 bits of the value with
// its fractional part discarded.
func (d *bigDecimal) IntValue() int32 {
	return int32(uint32(d.truncatedLow64()))
}

func (d *bigDecimal) truncatedLow64() uint64 {
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
