// Port of com/google/javascript/rhino/dtoa/DToA.java @ closure-compiler v20261104

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Waldemar Horwat
 *   Roger Lawrence
 *   Attila Szegedi
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

/****************************************************************
 *
 * The author of this software is David M. Gay.
 *
 * Copyright (c) 1991, 2000, 2001 by Lucent Technologies.
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose without fee is hereby granted, provided that this entire notice
 * is included in all copies of any software which is or includes a copy
 * or modification of this software and in all copies of the supporting
 * documentation for such software.
 *
 * THIS SOFTWARE IS BEING PROVIDED "AS IS", WITHOUT ANY EXPRESS OR IMPLIED
 * WARRANTY.  IN PARTICULAR, NEITHER THE AUTHOR NOR LUCENT MAKES ANY
 * REPRESENTATION OR WARRANTY OF ANY KIND CONCERNING THE MERCHANTABILITY
 * OF THIS SOFTWARE OR ITS FITNESS FOR ANY PARTICULAR PURPOSE.
 *
 ***************************************************************/

// Package dtoa is the port of com.google.javascript.rhino.dtoa: JavaScript
// number-to-string conversion.
package dtoa

import (
	"math"
	"math/big"
	"math/bits"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

var (
	bigFive = big.NewInt(5)
	bigTen  = big.NewInt(10)
)

// DToA_NumberToString is JavaScript's Number.prototype.toString().
//
// Java: DToA.numberToString(double)
func DToA_NumberToString(value float64) string {
	if value == 0.0 {
		return "0"
	}
	if value >= -9007199254740991 && value <= 9007199254740991 {
		if iv := int64(value); float64(iv) == value {
			return strconv.FormatInt(iv, 10)
		}
	}
	if math.IsNaN(value) {
		return "NaN"
	}
	if math.IsInf(value, 1) {
		return "Infinity"
	}
	if math.IsInf(value, -1) {
		return "-Infinity"
	}

	absVal := math.Abs(value)
	// Preserve byte-for-byte parity with Rhino DToA.java on subnormals whose
	// shifted significand fits entirely in the low 32-bit word (where DToA.d2b
	// computes bbits = 0).
	if absBits := math.Float64bits(absVal); absBits < 1<<52 {
		if y := uint32(absBits); y != 0 && (uint32(absBits>>32)&0xfffff)>>bits.TrailingZeros32(y) == 0 {
			return rhinoSubnormalToString(value < 0, absBits)
		}
	}

	// Format shortest IEEE-754 round-trip round-to-even representation ("d.dddde±dd" or "de±dd") via Ryū.
	var raw [32]byte
	b := strconv.AppendFloat(raw[:0], absVal, 'e', -1, 64)

	var digits [20]byte
	digits[0] = b[0]
	k := 1
	ePos := 1
	if b[1] == '.' {
		for ePos = 2; b[ePos] != 'e'; ePos++ {
			digits[k] = b[ePos]
			k++
		}
	}

	exp := 0
	for i := ePos + 2; i < len(b); i++ {
		exp = exp*10 + int(b[i]-'0')
	}
	if b[ePos+1] == '-' {
		exp = -exp
	}

	var out [32]byte
	return string(appendECMADigits(out[:0], value < 0, digits[:k], exp))
}

// appendECMADigits formats significand digits and base-10 exponent according to
// ECMAScript Number::toString(x, 10), where n = exp + 1 is the decimal point offset.
func appendECMADigits(out []byte, neg bool, digits []byte, exp int) []byte {
	if neg {
		out = append(out, '-')
	}
	k := len(digits)
	if n := exp + 1; -5 <= n && n <= 21 {
		switch {
		case k <= n:
			out = append(out, digits...)
			for ; k < n; k++ {
				out = append(out, '0')
			}
		case n > 0:
			out = append(out, digits[:n]...)
			out = append(out, '.')
			out = append(out, digits[n:]...)
		default:
			out = append(out, '0', '.')
			for range -n {
				out = append(out, '0')
			}
			out = append(out, digits...)
		}
		return out
	}
	out = append(out, digits[0])
	if k > 1 {
		out = append(out, '.')
		out = append(out, digits[1:]...)
	}
	out = append(out, 'e')
	if exp >= 0 {
		out = append(out, '+')
	}
	return strconv.AppendInt(out, int64(exp), 10)
}

// rhinoSubnormalToString reproduces Rhino DToA.JS_dtoa (mode 0) + JS_dtostr
// (DTOSTR_STANDARD) on the small subnormals where DToA.d2b computes bbits = 0.
func rhinoSubnormalToString(neg bool, absBits uint64) string {
	d1 := int32(absBits)
	tz := int32(bits.TrailingZeros32(uint32(d1)))
	b := new(big.Int).SetUint64((absBits & 0x000fffffffffffff) >> uint(tz))

	xd := float64(int64(d1) << uint(32-tz))
	xdBits := math.Float64bits(xd)
	d2 := math.Float64frombits((uint64(uint32(int32(xdBits>>32)-31*0x100000)) << 32) | (xdBits & 0xffffffff))
	i := tz - 1075

	ds := (d2-1.5)*0.289529654602168 + 0.1760912590558 + float64(i)*0.301029995663981
	k := jutil.D2I(ds)
	if ds < 0.0 && ds != float64(k) {
		k--
	}
	kCheck := true
	if k >= 0 && k <= 22 {
		// Positive subnormal < 1.0 <= 10^k whenever 0 <= k <= 22.
		k--
		kCheck = false
	}

	var b2, b5, s5 int32
	s2 := 1074 - tz
	if k >= 0 {
		s5 = k
		s2 += k
	} else {
		b2 -= k
		b5 = -k
	}

	m2, m5 := b2, b5
	step := tz + 1
	b2 += step
	s2 += step
	if m2 > 0 && s2 > 0 {
		c := min(m2, s2)
		b2 -= c
		m2 -= c
		s2 -= c
	}

	mhi := big.NewInt(1)
	if b5 > 0 {
		mhi = jutil.BigIntegerPow(bigFive, m5)
		b = jutil.BigIntegerMultiply(mhi, b)
	}
	S := big.NewInt(1)
	if s5 > 0 {
		S = jutil.BigIntegerPow(bigFive, s5)
	}

	sBytes := jutil.BigIntegerToByteArray(S)
	var sHiWord int32
	for idx := range 4 {
		sHiWord <<= 8
		if idx < len(sBytes) {
			sHiWord |= int32(sBytes[idx])
		}
	}
	hi := int32(1)
	if s5 != 0 {
		hi = 32 - int32(bits.LeadingZeros32(uint32(sHiWord)))
	}
	if pad := (hi + s2) & 0x1f; pad != 0 {
		i = 32 - pad
	} else {
		i = 0
	}
	if i > 4 {
		i -= 4
		b2 += i
		m2 += i
		s2 += i
	} else if i < 4 {
		i += 28
		b2 += i
		m2 += i
		s2 += i
	}
	if b2 > 0 {
		b = jutil.BigIntegerShiftLeft(b, b2)
	}
	if s2 > 0 {
		S = jutil.BigIntegerShiftLeft(S, s2)
	}
	if kCheck && b.Cmp(S) < 0 {
		k--
		b.Mul(b, bigTen)
		mhi.Mul(mhi, bigTen)
	}
	if m2 > 0 {
		mhi = jutil.BigIntegerShiftLeft(mhi, m2)
	}

	var digits []byte
	var q, rem, delta big.Int
	even := (d1 & 1) == 0
	for {
		q.QuoRem(b, S, &rem)
		b.Set(&rem)
		dig := byte(q.Int64() + '0')
		j := b.Cmp(mhi)
		delta.Sub(S, mhi)
		j1 := 1
		if delta.Sign() > 0 {
			j1 = b.Cmp(&delta)
		}
		if j1 == 0 && even {
			if j > 0 || dig == '9' {
				digits = dtoaRoundUpDigit(digits, dig, &k)
			} else {
				digits = append(digits, dig)
			}
			break
		}
		if j < 0 || (j == 0 && even) {
			if j1 > 0 {
				b.Lsh(b, 1)
				if c := b.Cmp(S); c > 0 || (c == 0 && (dig&1) == 1) {
					digits = dtoaRoundUpDigit(digits, dig, &k)
					break
				}
			}
			digits = append(digits, dig)
			break
		}
		if j1 > 0 {
			digits = dtoaRoundUpDigit(digits, dig, &k)
			break
		}
		digits = append(digits, dig)
		b.Mul(b, bigTen)
		mhi.Mul(mhi, bigTen)
	}
	return string(appendECMADigits(nil, neg, digits, int(k)))
}

func dtoaRoundUpDigit(buf []byte, dig byte, k *int32) []byte {
	if dig != '9' {
		return append(buf, dig+1)
	}
	for i := len(buf) - 1; i >= 0; i-- {
		if buf[i] != '9' {
			buf[i]++
			return buf[:i+1]
		}
	}
	*k++
	return append(buf[:0], '1')
}
