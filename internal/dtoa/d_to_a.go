// Port of com/google/javascript/rhino/dtoa/DToA.java @ closure-compiler v20260415

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

// DToA_D2b converts d into the form b*2^e, where b is an odd integer. b is
// the returned Bigint and e is the returned binary exponent. Return the
// number of significant bits in b in bits. d must be finite and nonzero.
//
// Java: DToA.d2b(double, int[], int[])
func DToA_D2b(d float64, e []int32, bits []int32) *big.Int {
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
		dbl_bits = make([]byte, 8)
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
		//        JS_ASSERT(z);
		dbl_bits = make([]byte, 4)
		k = DToA_Lo0bits(z)
		z = jutil.Ushr32(z, int64(k))
		DToA_StuffBits(dbl_bits, 0, z)
		k += 32
		i = 1
	}
	if de != 0 {
		e[0] = de - DToA_Bias - (DToA_P - 1) + k
		bits[0] = DToA_P - k
	} else {
		e[0] = de - DToA_Bias - (DToA_P - 1) + 1 + k
		bits[0] = 32*i - DToA_Hi0bits(z)
	}
	return jutil.NewBigIntegerFromBytes(dbl_bits)
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

// Java: DToA.roundOff(StringBuilder)
func DToA_RoundOff(buf *jutil.StringBuilder) bool {
	i := buf.Length()
	for i != 0 {
		i--
		c := buf.CharAt(i)
		if c != '9' {
			buf.SetCharAt(i, c+1)
			buf.SetLength(i + 1)
			return false
		}
	}
	buf.SetLength(0)
	return true
}

// bigShiftLeft is BigInteger.shiftLeft(n).
func bigShiftLeft(x *big.Int, n int32) *big.Int {
	return jutil.BigIntegerShiftLeft(x, n)
}

// bigMul10 is x.multiply(BigInteger.valueOf(10)).
func bigMul10(x *big.Int) *big.Int {
	return jutil.BigIntegerMultiply(x, big.NewInt(10))
}

// bumpUp is the inlined Java `bump_up` sequence: strip trailing '9's, then
// increment the last digit (or emit "1" and bump k when all were '9').
func bumpUp(buf *jutil.StringBuilder, k *int32) {
	var lastCh uint16
	for {
		lastCh = buf.CharAt(buf.Length() - 1)
		buf.SetLength(buf.Length() - 1)
		if lastCh != '9' {
			break
		}
		if buf.Length() == 0 {
			*k++
			lastCh = '0'
			break
		}
	}
	buf.AppendChar(lastCh + 1)
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
// `bump_up` loop are the shared helper bumpUp.
//
// Java: DToA.JS_dtoa(double, int, boolean, int, boolean[], StringBuilder)
func DToA_JS_dtoa(d float64, mode int32, biasUp bool, ndigits int32, sign []bool, buf *jutil.StringBuilder) int32 {
	/*  Arguments ndigits, decpt, sign are similar to those
	    of ecvt and fcvt; trailing zeros are suppressed from
	    the returned string.  If not null, *rve is set to point
	    to the end of the return value.  If d is +-Infinity or NaN,
	    then *decpt is set to 9999.

	    mode:
	    0 ==> shortest string that yields d when read in
	    and rounded to nearest.
	    1 ==> like 0, but with Steele & White stopping rule;
	    e.g. with IEEE P754 arithmetic , mode 0 gives
	    1e23 whereas mode 1 gives 9.999999999999999e22.
	    2 ==> max(1,ndigits) significant digits.  This gives a
	    return value similar to that of ecvt, except
	    that trailing zeros are suppressed.
	    3 ==> through ndigits past the decimal point.  This
	    gives a return value similar to that from fcvt,
	    except that trailing zeros are suppressed, and
	    ndigits can be negative.
	    4-9 should give the same return values as 2-3, i.e.,
	    4 <= mode <= 9 ==> same return as mode
	    2 + (mode & 1).  These modes are mainly for
	    debugging; often they run slower but sometimes
	    faster than modes 2-3.
	    4,5,8,9 ==> left-to-right digit generation.
	    6-9 ==> don't try fast floating-point estimate
	    (if applicable).

	    Values of mode other than 0-9 are treated as mode 0.

	    Sufficient space is allocated to the return value
	    to hold the suppressed trailing zeros.
	*/

	var b2, b5, i, ieps, ilim, ilim0, ilim1, j, j1, k, k0, m2, m5, s2, s5 int32
	var dig uint16
	var L int64
	var x int64
	var b, b1, delta, mlo, mhi, S *big.Int
	be := make([]int32, 1)
	bbits := make([]int32, 1)
	var d2, ds, eps float64
	var spec_case, denorm, k_check, try_quick, leftright bool
	_ = b1

	if (DToA_Word0(d) & DToA_Sign_bit) != 0 {
		/* set sign for everything, including 0's and NaNs */
		sign[0] = true
		// word0(d) &= ~Sign_bit;  /* clear sign bit */
		d = DToA_SetWord0(d, DToA_Word0(d)&^DToA_Sign_bit)
	} else {
		sign[0] = false
	}

	if (DToA_Word0(d) & DToA_Exp_mask) == DToA_Exp_mask {
		/* Infinity or NaN */
		if (DToA_Word1(d) == 0) && ((DToA_Word0(d) & DToA_Frac_mask) == 0) {
			buf.Append("Infinity")
		} else {
			buf.Append("NaN")
		}
		return 9999
	}
	if d == 0 {
		//          no_digits:
		buf.SetLength(0)
		buf.AppendChar('0') /* copy "0" to buffer */
		return 1
	}

	b = DToA_D2b(d, be, bbits)
	i = jutil.Ushr32(DToA_Word0(d), int64(DToA_Exp_shift1)) & (DToA_Exp_mask >> DToA_Exp_shift1)
	if i != 0 {
		d2 = DToA_SetWord0(d, (DToA_Word0(d)&DToA_Frac_mask1)|DToA_Exp_11)
		/* log(x)   ~=~ log(1.5) + (x-1.5)/1.5
		 * log10(x)  =  log(x) / log(10)
		 *      ~=~ log(1.5)/log(10) + (x-1.5)/(1.5*log(10))
		 * log10(d) = (i-Bias)*log(2)/log(10) + log10(d2)
		 *
		 * This suggests computing an approximation k to log10(d) by
		 *
		 * k = (i - Bias)*0.301029995663981
		 *  + ( (d2-1.5)*0.289529654602168 + 0.176091259055681 );
		 *
		 * We want k to be too large rather than too small.
		 * The error in the first-order Taylor series approximation
		 * is in our favor, so we just round up the constant enough
		 * to compensate for any error in the multiplication of
		 * (i - Bias) by 0.301029995663981; since |i - Bias| <= 1077,
		 * and 1077 * 0.30103 * 2^-52 ~=~ 7.2e-14,
		 * adding 1e-13 to the constant term more than suffices.
		 * Hence we adjust the constant term to 0.1760912590558.
		 * (We could get a more accurate k by invoking log10,
		 *  but this is probably not worthwhile.)
		 */
		i -= DToA_Bias
		denorm = false
	} else {
		/* d is denormalized */
		i = bbits[0] + be[0] + (DToA_Bias + (DToA_P - 1) - 1)
		if i > 32 {
			x = jutil.Shl64(int64(DToA_Word0(d)), int64(64-i)) | int64(jutil.Ushr32(DToA_Word1(d), int64(i-32)))
		} else {
			x = jutil.Shl64(int64(DToA_Word1(d)), int64(32-i))
		}
		//            d2 = x;
		//            word0(d2) -= 31*Exp_msk1; /* adjust exponent */
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
	j = bbits[0] - i - 1
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
		//            eps = ieps*d + 7.0;
		//            word0(eps) -= (P-1)*Exp_msk1;
		eps = float64(ieps)*d + 7.0
		eps = DToA_SetWord0(eps, DToA_Word0(eps)-(DToA_P-1)*DToA_Exp_msk1)
		if ilim == 0 {
			S = nil
			mhi = nil
			d -= 5.0
			if d > eps {
				buf.AppendChar('1')
				k++
				return k + 1
			}
			if d < -eps {
				buf.SetLength(0)
				buf.AppendChar('0') /* copy "0" to buffer */
				return 1
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
					buf.AppendChar(uint16('0' + L))
					if d < eps {
						return k + 1
					}
					if 1.0-d < eps {
						//                            goto bump_up;
						bumpUp(buf, &k)
						return k + 1
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
					buf.AppendChar(uint16('0' + L))
					if i == ilim {
						if d > 0.5+eps {
							//                                goto bump_up;
							bumpUp(buf, &k)
							return k + 1
						} else if d < 0.5-eps {
							DToA_StripTrailingZeroes(buf)
							//                                    while(*--s == '0') ;
							//                                    s++;
							return k + 1
						}
						break
					}
				}
			}
		}
		if fast_failed {
			buf.SetLength(0)
			d = d2
			k = k0
			ilim = ilim0
		}
	}

	/* Do we have a "small" integer? */

	if be[0] >= 0 && k <= DToA_Int_max {
		/* Yes. */
		ds = DToA_tens[k]
		if ndigits < 0 && ilim <= 0 {
			S = nil
			mhi = nil
			if ilim < 0 || d < 5*ds || (!biasUp && d == 5*ds) {
				buf.SetLength(0)
				buf.AppendChar('0') /* copy "0" to buffer */
				return 1
			}
			buf.AppendChar('1')
			k++
			return k + 1
		}
		for i = 1; ; i++ {
			L = jutil.D2L(d / ds)
			d -= float64(L) * ds
			buf.AppendChar(uint16('0' + L))
			if i == ilim {
				d += d
				if (d > ds) || (d == ds && (((L & 1) != 0) || biasUp)) {
					//                    bump_up:
					//                        while(*--s == '9')
					//                            if (s == buf) {
					//                                k++;
					//                                *s = '0';
					//                                break;
					//                            }
					//                        ++*s++;
					bumpUp(buf, &k)
				}
				break
			}
			d *= 10.0
			if d == 0 {
				break
			}
		}
		return k + 1
	}

	m2 = b2
	m5 = b5
	mhi = nil
	mlo = nil
	if leftright {
		if mode < 2 {
			if denorm {
				i = be[0] + (DToA_Bias + (DToA_P - 1) - 1 + 1)
			} else {
				i = 1 + DToA_P - bbits[0]
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
	 *
	 * Perhaps we should just compute leading 28 bits of S once
	 * and for all and pass them and a shift to quorem, so it
	 * can do shifts and ors to compute the numerator for q.
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
			/*no_digits:
			k = -1 - ndigits;
			goto ret; */
			buf.SetLength(0)
			buf.AppendChar('0') /* copy "0" to buffer */
			return 1
			//                goto no_digits;
		}
		//        one_digit:
		buf.AppendChar('1')
		k++
		return k + 1
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
			dig = uint16(jutil.BigIntegerIntValue(q) + '0')
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
					buf.AppendChar('9')
					if DToA_RoundOff(buf) {
						k++
						buf.AppendChar('1')
					}
					return k + 1
					//                        goto round_9_up;
				}
				if j > 0 {
					dig++
				}
				buf.AppendChar(dig)
				return k + 1
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
							buf.AppendChar('9')
							if DToA_RoundOff(buf) {
								k++
								buf.AppendChar('1')
							}
							return k + 1
							//                                goto round_9_up;
						}
					}
				}
				buf.AppendChar(dig)
				return k + 1
			}
			if j1 > 0 {
				if dig == '9' {
					/* possible if i == 1 */
					//                    round_9_up:
					//                        *s++ = '9';
					//                        goto roundoff;
					buf.AppendChar('9')
					if DToA_RoundOff(buf) {
						k++
						buf.AppendChar('1')
					}
					return k + 1
				}
				buf.AppendChar(dig + 1)
				return k + 1
			}
			buf.AppendChar(dig)
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
			//                (char)(dig = quorem(b,S) + '0');
			q, r := new(big.Int).QuoRem(b, S, new(big.Int))
			b = r
			dig = uint16(jutil.BigIntegerIntValue(q) + '0')
			buf.AppendChar(dig)
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
		//        roundoff:
		//            while(*--s == '9')
		//                if (s == buf) {
		//                    k++;
		//                    *s++ = '1';
		//                    goto ret;
		//                }
		//            ++*s++;
		if DToA_RoundOff(buf) {
			k++
			buf.AppendChar('1')
			return k + 1
		}
	} else {
		DToA_StripTrailingZeroes(buf)
		//            while(*--s == '0') ;
		//            s++;
	}
	//      ret:
	//        Bfree(S);
	//        if (mhi) {
	//            if (mlo && mlo != mhi)
	//                Bfree(mlo);
	//            Bfree(mhi);
	//        }
	//      ret1:
	//        Bfree(b);
	//        JS_ASSERT(s < buf + bufsize);
	return k + 1
}

// Java: DToA.stripTrailingZeroes(StringBuilder)
func DToA_StripTrailingZeroes(buf *jutil.StringBuilder) {
	//      while(*--s == '0') ;
	//      s++;
	bl := buf.Length()
	for {
		bl--
		if !(bl >= 0 && buf.CharAt(bl) == '0') {
			break
		}
		// empty
	}
	buf.SetLength(bl + 1)
}

// DToA_dtoaModes maps JSDToStrMode -> JS_dtoa mode.
//
// Java: DToA.dtoaModes
var DToA_dtoaModes = []int32{
	0 /* DTOSTR_STANDARD */, 0 /* DTOSTR_STANDARD_EXPONENTIAL, */, 3 /* DTOSTR_FIXED, */, 2, /* DTOSTR_EXPONENTIAL, */
	2,
} /* DTOSTR_PRECISION */

// Java: DToA.JS_dtostr(StringBuilder, int, int, double)
func DToA_JS_dtostr(buffer *jutil.StringBuilder, mode int32, precision int32, d float64) {
	var decPt int32         /* Position of decimal point relative to first digit returned by JS_dtoa */
	sign := make([]bool, 1) /* true if the sign bit was set in d */
	var nDigits int32       /* Number of significand digits returned by JS_dtoa */

	//        JS_ASSERT(bufferSize >= (size_t)(mode <= DTOSTR_STANDARD_EXPONENTIAL ? DTOSTR_STANDARD_BUFFER_SIZE :
	//                DTOSTR_VARIABLE_BUFFER_SIZE(precision)));

	if mode == DToA_DTOSTR_FIXED && (d >= 1e21 || d <= -1e21) {
		mode = DToA_DTOSTR_STANDARD /* Change mode here rather than below because the buffer may not be large enough to hold a large integer. */
	}

	decPt = DToA_JS_dtoa(d, DToA_dtoaModes[mode], mode >= DToA_DTOSTR_FIXED, precision, sign, buffer)
	nDigits = int32(buffer.Length())

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
			//                    JS_ASSERT(precision > 0);
			minNDigits = precision
			/* Fall through */
			fallthrough
		case DToA_DTOSTR_STANDARD_EXPONENTIAL:
			exponentialNotation = true

		case DToA_DTOSTR_PRECISION:
			//                    JS_ASSERT(precision > 0);
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
				buffer.AppendChar('0')
				if int32(buffer.Length()) == p {
					break
				}
			}
		}

		if exponentialNotation {
			/* Insert a decimal point if more than one significand digit */
			if nDigits != 1 {
				buffer.InsertChar(1, '.')
			}
			buffer.AppendChar('e')
			if (decPt - 1) >= 0 {
				buffer.AppendChar('+')
			}
			buffer.AppendInt(int64(decPt - 1))
			//                JS_snprintf(numEnd, bufferSize - (numEnd - buffer), "e%+d", decPt-1);
		} else if decPt != nDigits {
			/* Some kind of a fraction in fixed notation */
			//                JS_ASSERT(decPt <= nDigits);
			if decPt > 0 {
				/* dd...dd . dd...dd */
				buffer.InsertChar(int(decPt), '.')
			} else {
				/* 0 . 00...00dd...dd */
				for i := int32(0); i < 1-decPt; i++ {
					buffer.InsertChar(0, '0')
				}
				buffer.InsertChar(1, '.')
			}
		}
	}

	/* If negative and neither -0.0 nor NaN, output a leading '-'. */
	if sign[0] &&
		!(DToA_Word0(d) == DToA_Sign_bit && DToA_Word1(d) == 0) &&
		!((DToA_Word0(d)&DToA_Exp_mask) == DToA_Exp_mask &&
			((DToA_Word1(d) != 0) || ((DToA_Word0(d) & DToA_Frac_mask) != 0))) {
		buffer.InsertChar(0, '-')
	}
}

// DToA_NumberToString is JavaScript's Number.prototype.toString().
//
// Java: DToA.numberToString(double)
func DToA_NumberToString(value float64) string {
	if math.IsNaN(value) {
		return "NaN"
	}
	if value == 0.0 {
		return "0"
	}
	if value < 0 {
		return "-" + DToA_NumberToString(-value)
	}
	if math.IsInf(value, 0) {
		return "Infinity"
	}

	buffer := jutil.NewStringBuilder()
	DToA_JS_dtostr(buffer, DToA_DTOSTR_STANDARD, 0, value)
	return buffer.String()
}

// NewDToA is the private constructor.
//
// Java: DToA()
func NewDToA() *DToA { return &DToA{} }
