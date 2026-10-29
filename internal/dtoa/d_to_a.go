// Port of com/google/javascript/rhino/dtoa/DToA.java @ closure-compiler v20261021

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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DToA converts doubles to strings the way JavaScript does.
type DToA struct{}

const (
	DToA_DTOSTR_STANDARD             = 0 /* Either fixed or exponential format; round-trip */
	DToA_DTOSTR_STANDARD_EXPONENTIAL = 1 /* Always exponential format; round-trip */
	DToA_DTOSTR_FIXED                = 2 /* Round to <precision> digits after the decimal point; exponential if number is large */
	DToA_DTOSTR_EXPONENTIAL          = 3 /* Always exponential format; <precision> significant digits */
	DToA_DTOSTR_PRECISION            = 4 /* Either fixed or exponential format; <precision> significant digits */
)

const (
	DToA_Frac_mask int32 = 0xfffff
	DToA_Exp_shift int32 = 20
	DToA_Exp_msk1  int32 = 0x100000

	DToA_Bias int32 = 1023
	DToA_P    int32 = 53

	DToA_Exp_shift1 int32 = 20
	DToA_Exp_mask   int32 = 0x7ff00000
	DToA_Bndry_mask int32 = 0xfffff
	DToA_Log2P      int32 = 1

	DToA_Sign_bit   int32 = -0x80000000 // 0x80000000
	DToA_Exp_11     int32 = 0x3ff00000
	DToA_Ten_pmax   int32 = 22
	DToA_Quick_max  int32 = 14
	DToA_Bletch     int32 = 0x10
	DToA_Frac_mask1 int32 = 0xfffff
	DToA_Int_max    int32 = 14
	DToA_n_bigtens  int32 = 5
)

// Java: DToA.tens
var DToA_tens = []float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16,
	1e17, 1e18, 1e19, 1e20, 1e21, 1e22,
}

// Java: DToA.bigtens
var DToA_bigtens = []float64{1e16, 1e32, 1e64, 1e128, 1e256}

// Java: DToA.lo0bits(int)
func DToA_Lo0bits(y int32) int32 {
	var k int32
	x := y

	if (x & 7) != 0 {
		if (x & 1) != 0 {
			return 0
		}
		if (x & 2) != 0 {
			return 1
		}
		return 2
	}
	k = 0
	if (x & 0xffff) == 0 {
		k = 16
		x = jutil.Ushr32(x, 16)
	}
	if (x & 0xff) == 0 {
		k += 8
		x = jutil.Ushr32(x, 8)
	}
	if (x & 0xf) == 0 {
		k += 4
		x = jutil.Ushr32(x, 4)
	}
	if (x & 0x3) == 0 {
		k += 2
		x = jutil.Ushr32(x, 2)
	}
	if (x & 1) == 0 {
		k++
		x = jutil.Ushr32(x, 1)
		if (x & 1) == 0 {
			return 32
		}
	}
	return k
}

// DToA_Hi0bits returns the number (0 through 32) of most significant zero
// bits in x.
//
// Java: DToA.hi0bits(int)
func DToA_Hi0bits(x int32) int32 {
	var k int32 = 0

	if uint32(x)&0xffff0000 == 0 {
		k = 16
		x <<= 16
	}
	if uint32(x)&0xff000000 == 0 {
		k += 8
		x <<= 8
	}
	if uint32(x)&0xf0000000 == 0 {
		k += 4
		x <<= 4
	}
	if uint32(x)&0xc0000000 == 0 {
		k += 2
		x <<= 2
	}
	if uint32(x)&0x80000000 == 0 {
		k++
		if (x & 0x40000000) == 0 {
			return 32
		}
	}
	return k
}

// Java: DToA.stuffBits(byte[], int, int)
func DToA_StuffBits(bits []byte, offset int, val int32) {
	bits[offset] = byte(val >> 24)
	bits[offset+1] = byte(val >> 16)
	bits[offset+2] = byte(val >> 8)
	bits[offset+3] = byte(val)
}

func dtoaD2b(d float64, e *int32, bits *int32) *big.Int {
	var dblBuf [8]byte
	var dbl_bits []byte
	var i, k, y, z, de int32
	dBits := jutil.DoubleToLongBits(d)
	d0 := int32(jutil.Ushr64(dBits, 32))
	d1 := int32(dBits)

	z = d0 & DToA_Frac_mask
	d0 &= 0x7fffffff /* clear sign bit, which we ignore */

	de = jutil.Ushr32(d0, int64(DToA_Exp_shift))
	if de != 0 {
		z |= DToA_Exp_msk1
	}

	y = d1
	if y != 0 {
		dbl_bits = dblBuf[:8]
		k = DToA_Lo0bits(y)
		y = jutil.Ushr32(y, int64(k))
		if k != 0 {
			DToA_StuffBits(dbl_bits, 4, y|jutil.Shl32(z, int64(32-k)))
			z = jutil.Shr32(z, int64(k))
		} else {
			DToA_StuffBits(dbl_bits, 4, y)
		}
		DToA_StuffBits(dbl_bits, 0, z)
		if z != 0 {
			i = 2
		} else {
			i = 1
		}
	} else {
		dbl_bits = dblBuf[:4]
		k = DToA_Lo0bits(z)
		z = jutil.Ushr32(z, int64(k))
		DToA_StuffBits(dbl_bits, 0, z)
		k += 32
		i = 1
	}
	if de != 0 {
		*e = de - DToA_Bias - (DToA_P - 1) + k
		*bits = DToA_P - k
	} else {
		*e = de - DToA_Bias - (DToA_P - 1) + 1 + k
		*bits = 32*i - DToA_Hi0bits(z)
	}
	return jutil.NewBigIntegerFromBytes(dbl_bits)
}

// DToA_D2b converts d into the form b*2^e, where b is an odd integer. b is
// the returned Bigint and e is the returned binary exponent. Return the
// number of significant bits in b in bits. d must be finite and nonzero.
//
// Java: DToA.d2b(double, int[], int[])
func DToA_D2b(d float64, e []int32, bits []int32) *big.Int {
	return dtoaD2b(d, &e[0], &bits[0])
}

/* dtoa for IEEE arithmetic (dmg): convert double to ASCII string.
 *
 * Inspired by "How to Print Floating-Point Numbers Accurately" by
 * Guy L. Steele, Jr. and Jon L. White [Proc. ACM SIGPLAN '90, pp. 92-101].
 *
 * Modifications:
 *  1. Rather than iterating, we use a simple numeric overestimate
 *     to determine k = floor(log10(d)).  We scale relevant
 *     quantities using O(log2(k)) rather than O(k) multiplications.
 *  2. For some modes > 2 (corresponding to ecvt and fcvt), we don't
 *     try to generate digits strictly left to right.  Instead, we
 *     compute with fewer bits and propagate the carry if necessary
 *     when rounding the final digit up.  This is often faster.
 *  3. Under the assumption that input will be rounded nearest,
 *     mode 0 renders 1e23 as 1e23 rather than 9.999999999999999e22.
 *     That is, we allow equality in stopping tests when the
 *     round-nearest rule will give the same floating-point value
 *     as would satisfaction of the stopping test with strict
 *     inequality.
 *  4. We remove common factors of powers of 2 from relevant
 *     quantities.
 *  5. When converting floating-point integers less than 1e16,
 *     we use floating-point arithmetic rather than resorting
 *     to multiple-precision integers.
 *  6. When asked to produce fewer than 15 digits, we first try
 *     to get by with floating-point arithmetic; we resort to
 *     multiple-precision integer arithmetic only if we cannot
 *     guarantee that the floating-point calculation has given
 *     the correctly rounded result.  For k requested digits and
 *     "uniformly" distributed input, the probability is
 *     something like 10^(k-15) that we must resort to the Long
 *     calculation.
 */

// Java: DToA.word0(double)
func DToA_Word0(d float64) int32 {
	dBits := jutil.DoubleToLongBits(d)
	return int32(dBits >> 32)
}

// Java: DToA.setWord0(double, int)
func DToA_SetWord0(d float64, i int32) float64 {
	dBits := jutil.DoubleToLongBits(d)
	dBits = (int64(i) << 32) | (dBits & 0x0FFFFFFFF)
	return jutil.LongBitsToDouble(dBits)
}

// Java: DToA.word1(double)
func DToA_Word1(d float64) int32 {
	dBits := jutil.DoubleToLongBits(d)
	return int32(dBits)
}

// DToA_Pow5mult returns b * 5^k. k must be nonnegative.
//
// Java: DToA.pow5mult(BigInteger, int)
func DToA_Pow5mult(b *big.Int, k int32) *big.Int {
	// XXXX the C version built a cache of these
	return jutil.BigIntegerMultiply(b, jutil.BigIntegerPow(big.NewInt(5), k))
}

func dtoaRoundOff(buf []byte) ([]byte, bool) {
	i := len(buf)
	for i != 0 {
		i--
		c := buf[i]
		if c != '9' {
			buf[i] = c + 1
			return buf[:i+1], false
		}
	}
	return buf[:0], true
}

// Java: DToA.roundOff(StringBuilder)
func DToA_RoundOff(buf *[]byte) bool {
	var carry bool
	*buf, carry = dtoaRoundOff(*buf)
	return carry
}

// bigShiftLeft is BigInteger.shiftLeft(n).
func bigShiftLeft(x *big.Int, n int32) *big.Int {
	return jutil.BigIntegerShiftLeft(x, n)
}

// bigMul10 is x.multiply(BigInteger.valueOf(10)).
func bigMul10(x *big.Int) *big.Int {
	return jutil.BigIntegerMultiply(x, big.NewInt(10))
}

// dtoaBumpUp is the inlined Java `bump_up` sequence: strip trailing '9's, then
// increment the last digit (or emit "1" and bump k when all were '9').
func dtoaBumpUp(buf []byte, k *int32) []byte {
	var lastCh byte
	for {
		lastCh = buf[len(buf)-1]
		buf = buf[:len(buf)-1]
		if lastCh != '9' {
			break
		}
		if len(buf) == 0 {
			*k++
			lastCh = '0'
			break
		}
	}
	return append(buf, lastCh+1)
}

func dtoaInsertByte(buf []byte, pos int, ch byte) []byte {
	buf = append(buf, 0)
	copy(buf[pos+1:], buf[pos:])
	buf[pos] = ch
	return buf
}

// DToA_JS_dtoa always emits at least one digit.
//
// If biasUp is set, then rounding in modes 2 and 3 will round away from zero
// when the number is exactly halfway between two representable values. For
// example, rounding 2.5 to zero digits after the decimal point will return 3
// and not 2. 2.49 will still round to 2, and 2.51 will still round to 3.
// bufsize should be at least 20 for modes 0 and 1. For the other modes,
// bufsize should be two greater than the maximum number of output characters
// expected.
//
// Go-only: language-adapter — the three inlined copies of the Java
// `bump_up` loop are the shared helper dtoaBumpUp; digits accumulate in a
// stack-backed ASCII []byte buffer without StringBuilder heap allocations.
//
// Java: DToA.JS_dtoa(double, int, boolean, int, boolean[], StringBuilder)
func DToA_JS_dtoa(d float64, mode int32, biasUp bool, ndigits int32, sign []bool, buf *[]byte) int32 {
	var decPt int32
	*buf, decPt = dtoaJSdtoa(d, mode, biasUp, ndigits, &sign[0], *buf)
	return decPt
}

func dtoaJSdtoa(d float64, mode int32, biasUp bool, ndigits int32, sign *bool, buf []byte) ([]byte, int32) {
	var b2, b5, i, ieps, ilim, ilim0, ilim1, j, j1, k, k0, m2, m5, s2, s5 int32
	var dig byte
	var L int64
	var x int64
	var b, b1, delta, mlo, mhi, S *big.Int
	var be, bbits int32
	var d2, ds, eps float64
	var spec_case, denorm, k_check, try_quick, leftright bool
	_ = b1

	if (DToA_Word0(d) & DToA_Sign_bit) != 0 {
		/* set sign for everything, including 0's and NaNs */
		*sign = true
		d = DToA_SetWord0(d, DToA_Word0(d)&^DToA_Sign_bit)
	} else {
		*sign = false
	}

	if (DToA_Word0(d) & DToA_Exp_mask) == DToA_Exp_mask {
		/* Infinity or NaN */
		if (DToA_Word1(d) == 0) && ((DToA_Word0(d) & DToA_Frac_mask) == 0) {
			buf = append(buf, "Infinity"...)
		} else {
			buf = append(buf, "NaN"...)
		}
		return buf, 9999
	}
	if d == 0 {
		buf = append(buf[:0], '0') /* copy "0" to buffer */
		return buf, 1
	}

	b = dtoaD2b(d, &be, &bbits)
	i = jutil.Ushr32(DToA_Word0(d), int64(DToA_Exp_shift1)) & (DToA_Exp_mask >> DToA_Exp_shift1)
	if i != 0 {
		d2 = DToA_SetWord0(d, (DToA_Word0(d)&DToA_Frac_mask1)|DToA_Exp_11)
		i -= DToA_Bias
		denorm = false
	} else {
		/* d is denormalized */
		i = bbits + be + (DToA_Bias + (DToA_P - 1) - 1)
		if i > 32 {
			x = jutil.Shl64(int64(DToA_Word0(d)), int64(64-i)) | int64(jutil.Ushr32(DToA_Word1(d), int64(i-32)))
		} else {
			x = jutil.Shl64(int64(DToA_Word1(d)), int64(32-i))
		}
		// Go-only: language-adapter — Java widens the long x to double when
		// passing it to word0/setWord0(double, ...).
		xd := float64(x)
		d2 = DToA_SetWord0(xd, DToA_Word0(xd)-31*DToA_Exp_msk1)
		i -= (DToA_Bias + (DToA_P - 1) - 1) + 1
		denorm = true
	}
	/* At this point d = f*2^i, where 1 <= f < 2.  d2 is an approximation of f. */
	ds = (d2-1.5)*0.289529654602168 + 0.1760912590558 + float64(i)*0.301029995663981
	k = jutil.D2I(ds)
	if ds < 0.0 && ds != float64(k) {
		k-- /* want k = floor(ds) */
	}
	k_check = true
	if k >= 0 && k <= DToA_Ten_pmax {
		if d < DToA_tens[k] {
			k--
		}
		k_check = false
	}
	/* At this point floor(log10(d)) <= k <= floor(log10(d))+1.
	If k_check is zero, we're guaranteed that k = floor(log10(d)). */
	j = bbits - i - 1
	/* At this point d = b/2^j, where b is an odd integer. */
	if j >= 0 {
		b2 = 0
		s2 = j
	} else {
		b2 = -j
		s2 = 0
	}
	if k >= 0 {
		b5 = 0
		s5 = k
		s2 += k
	} else {
		b2 -= k
		b5 = -k
		s5 = 0
	}
	/* At this point d/10^k = (b * 2^b2 * 5^b5) / (2^s2 * 5^s5), where b is an odd integer,
	b2 >= 0, b5 >= 0, s2 >= 0, and s5 >= 0. */
	if mode < 0 || mode > 9 {
		mode = 0
	}
	try_quick = true
	if mode > 5 {
		mode -= 4
		try_quick = false
	}
	leftright = true
	ilim = 0
	ilim1 = 0
	switch mode {
	case 0, 1:
		ilim = -1
		ilim1 = -1
		i = 18
		ndigits = 0
	case 2:
		leftright = false
		fallthrough
	case 4:
		if ndigits <= 0 {
			ndigits = 1
		}
		ilim = ndigits
		ilim1 = ndigits
		i = ndigits
	case 3:
		leftright = false
		fallthrough
	case 5:
		i = ndigits + k + 1
		ilim = i
		ilim1 = i - 1
		if i <= 0 {
			i = 1
		}
	}
	/* ilim is the maximum number of significant digits we want, based on k and ndigits. */
	/* ilim1 is the maximum number of significant digits we want, based on k and ndigits,
	when it turns out that k was computed too high by one. */

	fast_failed := false
	if ilim >= 0 && ilim <= DToA_Quick_max && try_quick {

		/* Try to get by with floating-point arithmetic. */

		i = 0
		d2 = d
		k0 = k
		ilim0 = ilim
		ieps = 2 /* conservative */
		/* Divide d by 10^k, keeping track of the roundoff error and avoiding overflows. */
		if k > 0 {
			ds = DToA_tens[k&0xf]
			j = k >> 4
			if (j & DToA_Bletch) != 0 {
				/* prevent overflows */
				j &= DToA_Bletch - 1
				d /= DToA_bigtens[DToA_n_bigtens-1]
				ieps++
			}
			for ; j != 0; j, i = j>>1, i+1 {
				if (j & 1) != 0 {
					ieps++
					ds *= DToA_bigtens[i]
				}
			}
			d /= ds
		} else {
			j1 = -k
			if j1 != 0 {
				d *= DToA_tens[j1&0xf]
				for j = j1 >> 4; j != 0; j, i = j>>1, i+1 {
					if (j & 1) != 0 {
						ieps++
						d *= DToA_bigtens[i]
					}
				}
			}
		}
		/* Check that k was computed correctly. */
		if k_check && d < 1.0 && ilim > 0 {
			if ilim1 <= 0 {
				fast_failed = true
			} else {
				ilim = ilim1
				k--
				d *= 10.
				ieps++
			}
		}
		/* eps bounds the cumulative error. */
		eps = float64(ieps)*d + 7.0
		eps = DToA_SetWord0(eps, DToA_Word0(eps)-(DToA_P-1)*DToA_Exp_msk1)
		if ilim == 0 {
			S = nil
			mhi = nil
			d -= 5.0
			if d > eps {
				buf = append(buf, '1')
				k++
				return buf, k + 1
			}
			if d < -eps {
				buf = append(buf[:0], '0') /* copy "0" to buffer */
				return buf, 1
			}
			fast_failed = true
		}
		if !fast_failed {
			fast_failed = true
			if leftright {
				/* Use Steele & White method of only
				 * generating digits needed.
				 */
				eps = 0.5/DToA_tens[ilim-1] - eps
				for i = 0; ; {
					L = jutil.D2L(d)
					d -= float64(L)
					buf = append(buf, byte('0'+L))
					if d < eps {
						return buf, k + 1
					}
					if 1.0-d < eps {
						buf = dtoaBumpUp(buf, &k)
						return buf, k + 1
					}
					i++
					if i >= ilim {
						break
					}
					eps *= 10.0
					d *= 10.0
				}
			} else {
				/* Generate ilim digits, then fix them up. */
				eps *= DToA_tens[ilim-1]
				for i = 1; ; i, d = i+1, d*10.0 {
					L = jutil.D2L(d)
					d -= float64(L)
					buf = append(buf, byte('0'+L))
					if i == ilim {
						if d > 0.5+eps {
							buf = dtoaBumpUp(buf, &k)
							return buf, k + 1
						} else if d < 0.5-eps {
							buf = dtoaStripTrailingZeroes(buf)
							return buf, k + 1
						}
						break
					}
				}
			}
		}
		if fast_failed {
			buf = buf[:0]
			d = d2
			k = k0
			ilim = ilim0
		}
	}

	/* Do we have a "small" integer? */

	if be >= 0 && k <= DToA_Int_max {
		/* Yes. */
		ds = DToA_tens[k]
		if ndigits < 0 && ilim <= 0 {
			S = nil
			mhi = nil
			if ilim < 0 || d < 5*ds || (!biasUp && d == 5*ds) {
				buf = append(buf[:0], '0') /* copy "0" to buffer */
				return buf, 1
			}
			buf = append(buf, '1')
			k++
			return buf, k + 1
		}
		for i = 1; ; i++ {
			L = jutil.D2L(d / ds)
			d -= float64(L) * ds
			buf = append(buf, byte('0'+L))
			if i == ilim {
				d += d
				if (d > ds) || (d == ds && (((L & 1) != 0) || biasUp)) {
					buf = dtoaBumpUp(buf, &k)
				}
				break
			}
			d *= 10.0
			if d == 0 {
				break
			}
		}
		return buf, k + 1
	}

	m2 = b2
	m5 = b5
	mhi = nil
	mlo = nil
	if leftright {
		if mode < 2 {
			if denorm {
				i = be + (DToA_Bias + (DToA_P - 1) - 1 + 1)
			} else {
				i = 1 + DToA_P - bbits
			}
			/* i is 1 plus the number of trailing zero bits in d's significand. Thus,
			(2^m2 * 5^m5) / (2^(s2+i) * 5^s5) = (1/2 lsb of d)/10^k. */
		} else {
			j = ilim - 1
			if m5 >= j {
				m5 -= j
			} else {
				j -= m5
				s5 += j
				b5 += j
				m5 = 0
			}
			i = ilim
			if i < 0 {
				m2 -= i
				i = 0
			}
			/* (2^m2 * 5^m5) / (2^(s2+i) * 5^s5) = (1/2 * 10^(1-ilim))/10^k. */
		}
		b2 += i
		s2 += i
		mhi = big.NewInt(1)
		/* (mhi * 2^m2 * 5^m5) / (2^s2 * 5^s5) = one-half of last printed (when mode >= 2) or
		input (when mode < 2) significant digit, divided by 10^k. */
	}
	/* We still have d/10^k = (b * 2^b2 * 5^b5) / (2^s2 * 5^s5).  Reduce common factors in
	b2, m2, and s2 without changing the equalities. */
	if m2 > 0 && s2 > 0 {
		if m2 < s2 {
			i = m2
		} else {
			i = s2
		}
		b2 -= i
		m2 -= i
		s2 -= i
	}

	/* Fold b5 into b and m5 into mhi. */
	if b5 > 0 {
		if leftright {
			if m5 > 0 {
				mhi = DToA_Pow5mult(mhi, m5)
				b1 = jutil.BigIntegerMultiply(mhi, b)
				b = b1
			}
			j = b5 - m5
			if j != 0 {
				b = DToA_Pow5mult(b, j)
			}
		} else {
			b = DToA_Pow5mult(b, b5)
		}
	}
	/* Now we have d/10^k = (b * 2^b2) / (2^s2 * 5^s5) and
	(mhi * 2^m2) / (2^s2 * 5^s5) = one-half of last printed or input significant digit, divided by 10^k. */

	S = big.NewInt(1)
	if s5 > 0 {
		S = DToA_Pow5mult(S, s5)
	}
	/* Now we have d/10^k = (b * 2^b2) / (S * 2^s2) and
	(mhi * 2^m2) / (S * 2^s2) = one-half of last printed or input significant digit, divided by 10^k. */

	/* Check for special case that d is a normalized power of 2. */
	spec_case = false
	if mode < 2 {
		if (DToA_Word1(d) == 0) &&
			((DToA_Word0(d) & DToA_Bndry_mask) == 0) &&
			((DToA_Word0(d) & (DToA_Exp_mask & jutil.Shl32(DToA_Exp_mask, 1))) != 0) {
			/* The special case.  Here we want to be within a quarter of the last input
			significant digit instead of one half of it when the decimal output string's value is less than d.  */
			b2 += DToA_Log2P
			s2 += DToA_Log2P
			spec_case = true
		}
	}

	/* Arrange for convenient computation of quotients:
	 * shift left if necessary so divisor has 4 leading 0 bits.
	 */
	S_bytes := jutil.BigIntegerToByteArray(S)
	var S_hiWord int32 = 0
	for idx := 0; idx < 4; idx++ {
		S_hiWord = (S_hiWord << 8)
		if idx < len(S_bytes) {
			S_hiWord |= int32(S_bytes[idx] & 0xFF)
		}
	}
	var hi int32
	if s5 != 0 {
		hi = 32 - DToA_Hi0bits(S_hiWord)
	} else {
		hi = 1
	}
	i = (hi + s2) & 0x1f
	if i != 0 {
		i = 32 - i
	}
	/* i is the number of leading zero bits in the most significant word of S*2^s2. */
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
	/* Now S*2^s2 has exactly four leading zero bits in its most significant word. */
	if b2 > 0 {
		b = bigShiftLeft(b, b2)
	}
	if s2 > 0 {
		S = bigShiftLeft(S, s2)
	}
	/* Now we have d/10^k = b/S and
	(mhi * 2^m2) / S = maximum acceptable error, divided by 10^k. */
	if k_check {
		if b.Cmp(S) < 0 {
			k--
			b = bigMul10(b) /* we botched the k estimate */
			if leftright {
				mhi = bigMul10(mhi)
			}
			ilim = ilim1
		}
	}
	/* At this point 1 <= d/10^k = b/S < 10. */

	if ilim <= 0 && mode > 2 {
		/* We're doing fixed-mode output and d is less than the minimum nonzero output in this mode.
		Output either zero or the minimum nonzero output depending on which is closer to d. */
		noDigits := ilim < 0
		if !noDigits {
			S = jutil.BigIntegerMultiply(S, big.NewInt(5))
			i = int32(b.Cmp(S))
			noDigits = i < 0 || (i == 0 && !biasUp)
		}
		if noDigits {
			/* Always emit at least one digit.  If the number appears to be zero
			using the current mode, then emit one '0' digit and set decpt to 1. */
			buf = append(buf[:0], '0') /* copy "0" to buffer */
			return buf, 1
		}
		buf = append(buf, '1')
		k++
		return buf, k + 1
	}
	if leftright {
		if m2 > 0 {
			mhi = bigShiftLeft(mhi, m2)
		}

		/* Compute mlo -- check for special case
		 * that d is a normalized power of 2.
		 */
		mlo = mhi
		if spec_case {
			mhi = mlo
			mhi = bigShiftLeft(mhi, DToA_Log2P)
		}
		/* mlo/S = maximum acceptable error, divided by 10^k, if the output is less than d. */
		/* mhi/S = maximum acceptable error, divided by 10^k, if the output is greater than d. */

		for i = 1; ; i++ {
			q, r := new(big.Int).QuoRem(b, S, new(big.Int))
			b = r
			dig = byte(jutil.BigIntegerIntValue(q) + '0')
			/* Do we yet have the shortest decimal string
			 * that will round to d?
			 */
			j = int32(b.Cmp(mlo))
			/* j is b/S compared with mlo/S. */
			delta = new(big.Int).Sub(S, mhi)
			if delta.Sign() <= 0 {
				j1 = 1
			} else {
				j1 = int32(b.Cmp(delta))
			}
			/* j1 is b/S compared with 1 - mhi/S. */
			if (j1 == 0) && (mode == 0) && ((DToA_Word1(d) & 1) == 0) {
				if dig == '9' {
					buf = append(buf, '9')
					var carry bool
					buf, carry = dtoaRoundOff(buf)
					if carry {
						k++
						buf = append(buf, '1')
					}
					return buf, k + 1
				}
				if j > 0 {
					dig++
				}
				buf = append(buf, dig)
				return buf, k + 1
			}
			if (j < 0) || ((j == 0) && (mode == 0) && ((DToA_Word1(d) & 1) == 0)) {
				if j1 > 0 {
					/* Either dig or dig+1 would work here as the least significant decimal digit.
					Use whichever would produce a decimal value closer to d. */
					b = bigShiftLeft(b, 1)
					j1 = int32(b.Cmp(S))
					if (j1 > 0) || (j1 == 0 && (((dig & 1) == 1) || biasUp)) {
						oldDig := dig
						dig++
						if oldDig == '9' {
							buf = append(buf, '9')
							var carry bool
							buf, carry = dtoaRoundOff(buf)
							if carry {
								k++
								buf = append(buf, '1')
							}
							return buf, k + 1
						}
					}
				}
				buf = append(buf, dig)
				return buf, k + 1
			}
			if j1 > 0 {
				if dig == '9' {
					/* possible if i == 1 */
					buf = append(buf, '9')
					var carry bool
					buf, carry = dtoaRoundOff(buf)
					if carry {
						k++
						buf = append(buf, '1')
					}
					return buf, k + 1
				}
				buf = append(buf, dig+1)
				return buf, k + 1
			}
			buf = append(buf, dig)
			if i == ilim {
				break
			}
			b = bigMul10(b)
			if mlo == mhi {
				mhi = bigMul10(mhi)
				mlo = mhi
			} else {
				mlo = bigMul10(mlo)
				mhi = bigMul10(mhi)
			}
		}
	} else {
		for i = 1; ; i++ {
			q, r := new(big.Int).QuoRem(b, S, new(big.Int))
			b = r
			dig = byte(jutil.BigIntegerIntValue(q) + '0')
			buf = append(buf, dig)
			if i >= ilim {
				break
			}
			b = bigMul10(b)
		}
	}

	/* Round off last digit */

	b = bigShiftLeft(b, 1)
	j = int32(b.Cmp(S))
	if (j > 0) || (j == 0 && (((dig & 1) == 1) || biasUp)) {
		var carry bool
		buf, carry = dtoaRoundOff(buf)
		if carry {
			k++
			buf = append(buf, '1')
			return buf, k + 1
		}
	} else {
		buf = dtoaStripTrailingZeroes(buf)
	}
	return buf, k + 1
}

func dtoaStripTrailingZeroes(buf []byte) []byte {
	bl := len(buf)
	for {
		bl--
		if !(bl >= 0 && buf[bl] == '0') {
			break
		}
	}
	return buf[:bl+1]
}

// Java: DToA.stripTrailingZeroes(StringBuilder)
func DToA_StripTrailingZeroes(buf *[]byte) {
	*buf = dtoaStripTrailingZeroes(*buf)
}

// DToA_dtoaModes maps JSDToStrMode -> JS_dtoa mode.
//
// Java: DToA.dtoaModes
var DToA_dtoaModes = []int32{
	0 /* DTOSTR_STANDARD */, 0 /* DTOSTR_STANDARD_EXPONENTIAL, */, 3 /* DTOSTR_FIXED, */, 2, /* DTOSTR_EXPONENTIAL, */
	2,
} /* DTOSTR_PRECISION */

func dtoaJSdtostr(buf []byte, mode int32, precision int32, d float64) []byte {
	if mode == DToA_DTOSTR_STANDARD && precision == 0 && d != 0 && d >= -9007199254740991 && d <= 9007199254740991 {
		if iv := int64(d); float64(iv) == d {
			return strconv.AppendInt(buf, iv, 10)
		}
	}

	var decPt int32   /* Position of decimal point relative to first digit returned by JS_dtoa */
	var sign bool     /* true if the sign bit was set in d */
	var nDigits int32 /* Number of significand digits returned by JS_dtoa */

	if mode == DToA_DTOSTR_FIXED && (d >= 1e21 || d <= -1e21) {
		mode = DToA_DTOSTR_STANDARD /* Change mode here rather than below because the buffer may not be large enough to hold a large integer. */
	}

	buf, decPt = dtoaJSdtoa(d, DToA_dtoaModes[mode], mode >= DToA_DTOSTR_FIXED, precision, &sign, buf)
	nDigits = int32(len(buf))

	/* If Infinity, -Infinity, or NaN, return the string regardless of the mode. */
	if decPt != 9999 {
		exponentialNotation := false
		var minNDigits int32 = 0 /* Minimum number of significand digits required by mode and precision */
		var p int32

		switch mode {
		case DToA_DTOSTR_STANDARD:
			if decPt < -5 || decPt > 21 {
				exponentialNotation = true
			} else {
				minNDigits = decPt
			}

		case DToA_DTOSTR_FIXED:
			if precision >= 0 {
				minNDigits = decPt + precision
			} else {
				minNDigits = decPt
			}

		case DToA_DTOSTR_EXPONENTIAL:
			minNDigits = precision
			fallthrough
		case DToA_DTOSTR_STANDARD_EXPONENTIAL:
			exponentialNotation = true

		case DToA_DTOSTR_PRECISION:
			minNDigits = precision
			if decPt < -5 || decPt > precision {
				exponentialNotation = true
			}
		}

		/* If the number has fewer than minNDigits, pad it with zeros at the end */
		if nDigits < minNDigits {
			p = minNDigits
			nDigits = minNDigits
			for {
				buf = append(buf, '0')
				if int32(len(buf)) == p {
					break
				}
			}
		}

		if exponentialNotation {
			/* Insert a decimal point if more than one significand digit */
			if nDigits != 1 {
				buf = dtoaInsertByte(buf, 1, '.')
			}
			buf = append(buf, 'e')
			if (decPt - 1) >= 0 {
				buf = append(buf, '+')
			}
			buf = strconv.AppendInt(buf, int64(decPt-1), 10)
		} else if decPt != nDigits {
			/* Some kind of a fraction in fixed notation */
			if decPt > 0 {
				/* dd...dd . dd...dd */
				buf = dtoaInsertByte(buf, int(decPt), '.')
			} else {
				/* 0 . 00...00dd...dd */
				nZeros := int(1 - decPt)
				oldLen := len(buf)
				for i := 0; i < nZeros+1; i++ {
					buf = append(buf, 0)
				}
				copy(buf[nZeros+1:], buf[:oldLen])
				buf[0] = '0'
				buf[1] = '.'
				for i := 2; i < nZeros+1; i++ {
					buf[i] = '0'
				}
			}
		}
	}

	/* If negative and neither -0.0 nor NaN, output a leading '-'. */
	if sign &&
		!(DToA_Word0(d) == DToA_Sign_bit && DToA_Word1(d) == 0) &&
		!((DToA_Word0(d)&DToA_Exp_mask) == DToA_Exp_mask &&
			((DToA_Word1(d) != 0) || ((DToA_Word0(d) & DToA_Frac_mask) != 0))) {
		buf = dtoaInsertByte(buf, 0, '-')
	}
	return buf
}

// Java: DToA.JS_dtostr(StringBuilder, int, int, double)
func DToA_JS_dtostr(buffer *[]byte, mode int32, precision int32, d float64) {
	*buffer = dtoaJSdtostr(*buffer, mode, precision, d)
}

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

	var scratch [64]byte
	return string(dtoaJSdtostr(scratch[:0], DToA_DTOSTR_STANDARD, 0, value))
}

// NewDToA is the private constructor.
//
// Java: DToA()
func NewDToA() *DToA { return &DToA{} }
