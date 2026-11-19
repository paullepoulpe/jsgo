// Clean-room specification based on Java Language Specification (JLS §5.1.3, §15.19).

package javalang

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// D2I performs Java's saturating narrowing cast (int) d (JLS §5.1.3):
//   - NaN -> 0
//   - >= math.MaxInt32 -> math.MaxInt32
//   - <= math.MinInt32 -> math.MinInt32
//   - otherwise truncates toward zero
func D2I(d float64) int32 {
	if math.IsNaN(d) {
		return 0
	}
	if d >= math.MaxInt32 {
		return math.MaxInt32
	}
	if d <= math.MinInt32 {
		return math.MinInt32
	}
	return int32(d)
}

// D2L performs Java's saturating narrowing cast (long) d (JLS §5.1.3):
//   - NaN -> 0
//   - >= math.MaxInt64 -> math.MaxInt64
//   - <= math.MinInt64 -> math.MinInt64
//   - otherwise truncates toward zero
func D2L(d float64) int64 {
	if math.IsNaN(d) {
		return 0
	}
	if d >= math.MaxInt64 {
		return math.MaxInt64
	}
	if d <= math.MinInt64 {
		return math.MinInt64
	}
	return int64(d)
}

// Shl32 is Java's 32-bit shift left (v << s), with shift distance masked to 5 bits (JLS §15.19).
func Shl32(v int32, s int) int32 {
	return v << (uint(s) & 31)
}

// Shr32 is Java's 32-bit signed shift right (v >> s), with shift distance masked to 5 bits (JLS §15.19).
func Shr32(v int32, s int) int32 {
	return v >> (uint(s) & 31)
}

// Ushr32 is Java's 32-bit unsigned shift right (v >>> s), with shift distance masked to 5 bits (JLS §15.19).
func Ushr32(v int32, s int) int32 {
	return int32(uint32(v) >> (uint(s) & 31))
}

// Shl64 is Java's 64-bit shift left (v << s), with shift distance masked to 6 bits (JLS §15.19).
func Shl64(v int64, s int) int64 {
	return v << (uint(s) & 63)
}

// Shr64 is Java's 64-bit signed shift right (v >> s), with shift distance masked to 6 bits (JLS §15.19).
func Shr64(v int64, s int) int64 {
	return v >> (uint(s) & 63)
}

// Ushr64 is Java's 64-bit unsigned shift right (v >>> s), with shift distance masked to 6 bits (JLS §15.19).
func Ushr64(v int64, s int) int64 {
	return int64(uint64(v) >> (uint(s) & 63))
}

// DRem is Java's double remainder `a % b` (IEEE 754 fmod).
func DRem(a, b float64) float64 {
	return math.Mod(a, b)
}

// DoubleToString formats d matching Java's Double.toString(d) (JDK 19+ Schubfach shortest round-trip).
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
	s := strconv.FormatFloat(a, 'e', -1, 64)
	digits, exp := splitSci(s)
	if len(digits) == 1 {
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
	if exp >= -3 && exp < 7 {
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

// FloatToString formats f matching Java's Float.toString(f) (JDK 19+ Schubfach shortest round-trip).
func FloatToString(f float32) string {
	d := float64(f)
	switch {
	case math.IsNaN(d):
		return "NaN"
	case math.IsInf(d, 1):
		return "Infinity"
	case math.IsInf(d, -1):
		return "-Infinity"
	case f == 0:
		if math.Signbit(d) {
			return "-0.0"
		}
		return "0.0"
	}
	neg := f < 0
	a := f
	if neg {
		a = -a
	}
	ad := float64(a)
	s := strconv.FormatFloat(ad, 'e', -1, 32)
	digits, exp := splitSci(s)
	if len(digits) == 1 {
		s2 := strconv.FormatFloat(ad, 'e', 1, 32)
		if v, err := strconv.ParseFloat(s2, 32); err == nil && float32(v) == a {
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
	if exp >= -3 && exp < 7 {
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

func splitSci(s string) (string, int) {
	e := strings.IndexByte(s, 'e')
	mant, expStr := s[:e], s[e+1:]
	exp, _ := strconv.Atoi(expStr)
	digits := strings.Replace(mant, ".", "", 1)
	return digits, exp
}

// IntegerToHexString returns the unsigned hexadecimal string representation of i.
func IntegerToHexString(i int32) string {
	return strconv.FormatUint(uint64(uint32(i)), 16)
}

// LongToHexString returns the unsigned hexadecimal string representation of i.
func LongToHexString(i int64) string {
	return strconv.FormatUint(uint64(i), 16)
}

// System_CurrentTimeMillis returns the current time in milliseconds.
func System_CurrentTimeMillis() int64 {
	return time.Now().UnixMilli()
}

// System_NanoTime returns the current value of the running Java Virtual Machine's high-resolution time source.
func System_NanoTime() int64 {
	return time.Now().UnixNano()
}

// System_Arraycopy copies an array from the specified source array to the destination array.
func System_Arraycopy[T any](src []T, srcPos int, dest []T, destPos int, length int) {
	if length < 0 || srcPos < 0 || destPos < 0 || srcPos+length > len(src) || destPos+length > len(dest) {
		panic(NewIndexOutOfBoundsException("Arraycopy index out of bounds"))
	}
	copy(dest[destPos:destPos+length], src[srcPos:srcPos+length])
}
