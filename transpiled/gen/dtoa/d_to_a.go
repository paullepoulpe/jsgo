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
// Port of com/google/javascript/rhino/dtoa/DToA.java (generated mechanically by j2go).

package dtoa

import (
	"math"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.dtoa.DToA
type DToA struct {
	idHash_ atomic.Int32
}

func (this_ *DToA) BindSelf_(s any) {
	_ = s
}
func (this_ *DToA) IsJavaNil_() bool { return this_ == nil }
func (*DToA) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.dtoa.DToA", "DToA"
}
func (this_ *DToA) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DToA) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.dtoa.DToA", "com.google.javascript.rhino.dtoa.DToA", func(v_ any) bool {
		return j2gort.InstanceOf[*DToA](v_)
	})
}

func StaticGet_DToA_DTOSTR_STANDARD() int32             { return int32(0) }
func StaticGet_DToA_DTOSTR_STANDARD_EXPONENTIAL() int32 { return int32(1) }
func StaticGet_DToA_DTOSTR_FIXED() int32                { return int32(2) }
func StaticGet_DToA_DTOSTR_EXPONENTIAL() int32          { return int32(3) }
func StaticGet_DToA_DTOSTR_PRECISION() int32            { return int32(4) }
func StaticGet_DToA_Frac_mask() int32                   { return int32(1048575) }
func StaticGet_DToA_Exp_shift() int32                   { return int32(20) }
func StaticGet_DToA_Exp_msk1() int32                    { return int32(1048576) }
func StaticGet_DToA_Bias() int32                        { return int32(1023) }
func StaticGet_DToA_P() int32                           { return int32(53) }
func StaticGet_DToA_Exp_shift1() int32                  { return int32(20) }
func StaticGet_DToA_Exp_mask() int32                    { return int32(2146435072) }
func StaticGet_DToA_Bndry_mask() int32                  { return int32(1048575) }
func StaticGet_DToA_Log2P() int32                       { return int32(1) }
func StaticGet_DToA_Sign_bit() int32                    { return int32(-2147483648) }
func StaticGet_DToA_Exp_11() int32                      { return int32(1072693248) }
func StaticGet_DToA_Ten_pmax() int32                    { return int32(22) }
func StaticGet_DToA_Quick_max() int32                   { return int32(14) }
func StaticGet_DToA_Bletch() int32                      { return int32(16) }
func StaticGet_DToA_Frac_mask1() int32                  { return int32(1048575) }
func StaticGet_DToA_Int_max() int32                     { return int32(14) }
func StaticGet_DToA_N_bigtens() int32                   { return int32(5) }

var static_DToA_tens []float64

func StaticGet_DToA_Tens() []float64 {
	clinit_DToA()
	return static_DToA_tens
}
func StaticSet_DToA_Tens(v []float64) {
	clinit_DToA()
	static_DToA_tens = v
}

var static_DToA_bigtens []float64

func StaticGet_DToA_Bigtens() []float64 {
	clinit_DToA()
	return static_DToA_bigtens
}
func StaticSet_DToA_Bigtens(v []float64) {
	clinit_DToA()
	static_DToA_bigtens = v
}

var static_DToA_dtoaModes []int32

func StaticGet_DToA_DtoaModes() []int32 {
	clinit_DToA()
	return static_DToA_dtoaModes
}
func StaticSet_DToA_DtoaModes(v []int32) {
	clinit_DToA()
	static_DToA_dtoaModes = v
}

var clinitDone_DToA uint32
var clinitState_DToA j2gort.ClinitState

func clinit_DToA() {
	if atomic.LoadUint32(&clinitDone_DToA) != 0 {
		return
	}
	clinitSlow_DToA()
}

func clinitSlow_DToA() {
	if !j2gort.EnterClinit(&clinitState_DToA, &clinitDone_DToA) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DToA, &clinitDone_DToA)
	static_DToA_tens = []float64{float64(1.0), float64(10.0), float64(100.0), float64(1000.0), float64(10000.0), float64(100000.0), float64(1000000.0), float64(1.0e7), float64(1.0e8), float64(1.0e9), float64(1.0e10), float64(1.0e11), float64(1.0e12), float64(1.0e13), float64(1.0e14), float64(1.0e15), float64(1.0e16), float64(1.0e17), float64(1.0e18), float64(1.0e19), float64(1.0e20), float64(1.0e21), float64(1.0e22)}
	static_DToA_bigtens = []float64{float64(1.0e16), float64(1.0e32), float64(1.0e64), float64(1.0e128), float64(1.0e256)}
	static_DToA_dtoaModes = []int32{int32(0), int32(0), int32(3), int32(2), int32(2)}
}

func (this_ *DToA) initFields_DToA() {
	_ = this_
}

// Java constructor: DToA()
func NewDToA() *DToA {
	clinit_DToA()
	this_ := &DToA{}
	InitDToA(this_)
	return this_
}

func InitDToA(this_ *DToA) {
	this_.initFields_DToA()
}

// Java static method: DToA.lo0bits(int y)
func DToA_Lo0bits(y int32) int32 {
	clinit_DToA()
	_ = y
	var k int32 = int32(0)
	_ = k
	var x int32 = y
	_ = x
	if (x & int32(7)) != int32(0) {
		if (x & int32(1)) != int32(0) {
			return int32(0)
		}
		if (x & int32(2)) != int32(0) {
			return int32(1)
		}
		return int32(2)
	}
	k = int32(0)
	if (x & int32(65535)) == int32(0) {
		k = int32(16)
		x = javalang.Ushr32(x, int(int32(16)))
	}
	if (x & int32(255)) == int32(0) {
		k = (k + int32(8))
		x = javalang.Ushr32(x, int(int32(8)))
	}
	if (x & int32(15)) == int32(0) {
		k = (k + int32(4))
		x = javalang.Ushr32(x, int(int32(4)))
	}
	if (x & int32(3)) == int32(0) {
		k = (k + int32(2))
		x = javalang.Ushr32(x, int(int32(2)))
	}
	if (x & int32(1)) == int32(0) {
		k = k + 1
		x = javalang.Ushr32(x, int(int32(1)))
		if (x & int32(1)) == int32(0) {
			return int32(32)
		}
	}
	return k
}

// Java static method: DToA.hi0bits(int x)
func DToA_Hi0bits(x int32) int32 {
	clinit_DToA()
	_ = x
	var k int32 = int32(0)
	_ = k
	if (x & int32(-65536)) == int32(0) {
		k = int32(16)
		x = javalang.Shl32(x, int(int32(16)))
	}
	if (x & int32(-16777216)) == int32(0) {
		k = (k + int32(8))
		x = javalang.Shl32(x, int(int32(8)))
	}
	if (x & int32(-268435456)) == int32(0) {
		k = (k + int32(4))
		x = javalang.Shl32(x, int(int32(4)))
	}
	if (x & int32(-1073741824)) == int32(0) {
		k = (k + int32(2))
		x = javalang.Shl32(x, int(int32(2)))
	}
	if (x & int32(-2147483648)) == int32(0) {
		k = k + 1
		if (x & int32(1073741824)) == int32(0) {
			return int32(32)
		}
	}
	return k
}

// Java static method: DToA.stuffBits(byte[] bits, int offset, int val)
func DToA_StuffBits(bits []int8, offset int32, val int32) {
	clinit_DToA()
	_ = bits
	_ = offset
	_ = val
	{
		a_1 := bits
		i_2 := offset
		a_1[i_2] = int8((javalang.Shr32(val, int(int32(24)))))
	}
	{
		a_3 := bits
		i_4 := (offset + int32(1))
		a_3[i_4] = int8((javalang.Shr32(val, int(int32(16)))))
	}
	{
		a_5 := bits
		i_6 := (offset + int32(2))
		a_5[i_6] = int8((javalang.Shr32(val, int(int32(8)))))
	}
	{
		a_7 := bits
		i_8 := (offset + int32(3))
		a_7[i_8] = int8((val))
	}
}

// Java static method: DToA.d2b(double d, int[] e, int[] bits)
func DToA_D2b(d float64, e []int32, bits []int32) *j2gort.Ext_java_math_BigInteger {
	clinit_DToA()
	_ = d
	_ = e
	_ = bits
	var dbl_bits []int8 = nil
	_ = dbl_bits
	var i int32 = int32(0)
	_ = i
	var k int32 = int32(0)
	_ = k
	var y int32 = int32(0)
	_ = y
	var z int32 = int32(0)
	_ = z
	var de int32 = int32(0)
	_ = de
	var dBits int64 = int64(math.Float64bits(d))
	_ = dBits
	var d0 int32 = int32((javalang.Ushr64(dBits, int(int32(32)))))
	_ = d0
	var d1 int32 = int32((dBits))
	_ = d1
	z = (d0 & int32(1048575))
	d0 = (d0 & int32(2147483647))
	{
		var iifeCond_1 int32
		{

			var asgn_9 int32 = (javalang.Ushr32(d0, int(int32(20))))
			de = asgn_9
			iifeCond_1 = asgn_9
		}
		if (iifeCond_1) != int32(0) {
			z = (z | int32(1048576))
		}
	}
	{
		var iifeCond_1 int32
		{

			var asgn_10 int32 = d1
			y = asgn_10
			iifeCond_1 = asgn_10
		}
		if (iifeCond_1) != int32(0) {
			dbl_bits = make([]int8, int32(8))
			k = DToA_Lo0bits(y)
			y = javalang.Ushr32(y, int(k))
			if k != int32(0) {
				DToA_StuffBits(dbl_bits, int32(4), (y | javalang.Shl32(z, int((int32(32)-k)))))
				z = javalang.Shr32(z, int(k))
			} else {
				DToA_StuffBits(dbl_bits, int32(4), y)
			}
			DToA_StuffBits(dbl_bits, int32(0), z)
			if z != int32(0) {
				i = int32(2)
			} else {
				i = int32(1)
			}
		} else {
			dbl_bits = make([]int8, int32(4))
			k = DToA_Lo0bits(z)
			z = javalang.Ushr32(z, int(k))
			DToA_StuffBits(dbl_bits, int32(0), z)
			k = (k + int32(32))
			i = int32(1)
		}
	}
	if de != int32(0) {
		{
			a_11 := e
			i_12 := int32(0)
			a_11[i_12] = (((de - int32(1023)) - (int32(53) - int32(1))) + k)
		}
		{
			a_13 := bits
			i_14 := int32(0)
			a_13[i_14] = (int32(53) - k)
		}
	} else {
		{
			a_15 := e
			i_16 := int32(0)
			a_15[i_16] = ((((de - int32(1023)) - (int32(53) - int32(1))) + int32(1)) + k)
		}
		{
			a_17 := bits
			i_18 := int32(0)
			a_17[i_18] = ((int32(32) * i) - DToA_Hi0bits(z))
		}
	}
	return j2gort.NewExt_java_math_BigInteger_1_slice_int8(dbl_bits)
}

// Java static method: DToA.word0(double d)
func DToA_Word0(d float64) int32 {
	clinit_DToA()
	_ = d
	var dBits int64 = int64(math.Float64bits(d))
	_ = dBits
	return int32((javalang.Shr64(dBits, int(int32(32)))))
}

// Java static method: DToA.setWord0(double d, int i)
func DToA_SetWord0(d float64, i int32) float64 {
	clinit_DToA()
	_ = d
	_ = i
	var dBits int64 = int64(math.Float64bits(d))
	_ = dBits
	dBits = ((javalang.Shl64(int64(i), int(int32(32)))) | (dBits & int64(4294967295)))
	return math.Float64frombits(uint64(dBits))
}

// Java static method: DToA.word1(double d)
func DToA_Word1(d float64) int32 {
	clinit_DToA()
	_ = d
	var dBits int64 = int64(math.Float64bits(d))
	_ = dBits
	return int32((dBits))
}

// Java static method: DToA.pow5mult(java.math.BigInteger b, int k)
func DToA_Pow5mult(b *j2gort.Ext_java_math_BigInteger, k int32) *j2gort.Ext_java_math_BigInteger {
	clinit_DToA()
	_ = b
	_ = k
	return b.Multiply(j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(5))).Pow(k))
}

// Java static method: DToA.roundOff(java.lang.StringBuilder buf)
func DToA_RoundOff(buf *j2gort.Ext_java_lang_StringBuilder) bool {
	clinit_DToA()
	_ = buf
	var i int32 = buf.Length()
	_ = i
	for i != int32(0) {
		i = i - 1
		var c uint16 = buf.CharAt(i)
		_ = c
		if int32(c) != int32(uint16(57)) {
			buf.SetCharAt(i, uint16((int32(c) + int32(1))))
			buf.SetLength((i + int32(1)))
			return false
		}
	}
	buf.SetLength(int32(0))
	return true
}

// Java static method: DToA.JS_dtoa(double d, int mode, boolean biasUp, int ndigits, boolean[] sign, java.lang.StringBuilder buf)
func DToA_JS_dtoa(d float64, mode int32, biasUp bool, ndigits int32, sign []bool, buf *j2gort.Ext_java_lang_StringBuilder) int32 {
	clinit_DToA()
	_ = d
	_ = mode
	_ = biasUp
	_ = ndigits
	_ = sign
	_ = buf
	var b2 int32 = int32(0)
	_ = b2
	var b5 int32 = int32(0)
	_ = b5
	var i int32 = int32(0)
	_ = i
	var ieps int32 = int32(0)
	_ = ieps
	var ilim int32 = int32(0)
	_ = ilim
	var ilim0 int32 = int32(0)
	_ = ilim0
	var ilim1 int32 = int32(0)
	_ = ilim1
	var j int32 = int32(0)
	_ = j
	var j1 int32 = int32(0)
	_ = j1
	var k int32 = int32(0)
	_ = k
	var k0 int32 = int32(0)
	_ = k0
	var m2 int32 = int32(0)
	_ = m2
	var m5 int32 = int32(0)
	_ = m5
	var s2 int32 = int32(0)
	_ = s2
	var s5 int32 = int32(0)
	_ = s5
	var dig uint16 = uint16(0)
	_ = dig
	var l_L int64 = int64(0)
	_ = l_L
	var x int64 = int64(0)
	_ = x
	var b *j2gort.Ext_java_math_BigInteger = nil
	_ = b
	var b1 *j2gort.Ext_java_math_BigInteger = nil
	_ = b1
	var delta *j2gort.Ext_java_math_BigInteger = nil
	_ = delta
	var mlo *j2gort.Ext_java_math_BigInteger = nil
	_ = mlo
	var mhi *j2gort.Ext_java_math_BigInteger = nil
	_ = mhi
	var l_S *j2gort.Ext_java_math_BigInteger = nil
	_ = l_S
	var be []int32 = make([]int32, int32(1))
	_ = be
	var bbits []int32 = make([]int32, int32(1))
	_ = bbits
	var d2 float64 = float64(0)
	_ = d2
	var ds float64 = float64(0)
	_ = ds
	var eps float64 = float64(0)
	_ = eps
	var spec_case bool = false
	_ = spec_case
	var denorm bool = false
	_ = denorm
	var k_check bool = false
	_ = k_check
	var try_quick bool = false
	_ = try_quick
	var leftright bool = false
	_ = leftright
	if (DToA_Word0(d) & int32(-2147483648)) != int32(0) {
		{
			a_19 := sign
			i_20 := int32(0)
			a_19[i_20] = true
		}
		d = DToA_SetWord0(d, (DToA_Word0(d) & ^(int32(-2147483648))))
	} else {
		{
			a_21 := sign
			i_22 := int32(0)
			a_21[i_22] = false
		}
	}
	if (DToA_Word0(d) & int32(2146435072)) == int32(2146435072) {
		var iifeArg_2 string
	iifeSwitch_1:
		switch {
		default:

			if (DToA_Word1(d) == int32(0)) && ((DToA_Word0(d) & int32(1048575)) == int32(0)) {
				iifeArg_2 = "Infinity"
				break iifeSwitch_1
			}
			iifeArg_2 = "NaN"
			break iifeSwitch_1
		}
		buf.Append_1_string(iifeArg_2)
		return int32(9999)
	}
	if d == float64(int32(0)) {
		buf.SetLength(int32(0))
		buf.Append_1_uint16(uint16(48))
		return int32(1)
	}
	b = DToA_D2b(d, be, bbits)
	{
		var iifeCond_1 int32
		{

			var asgn_23 int32 = (javalang.Ushr32(DToA_Word0(d), int(int32(20))) & (javalang.Shr32(int32(2146435072), int(int32(20)))))
			i = asgn_23
			iifeCond_1 = asgn_23
		}
		if (iifeCond_1) != int32(0) {
			d2 = DToA_SetWord0(d, ((DToA_Word0(d) & int32(1048575)) | int32(1072693248)))
			i = (i - int32(1023))
			denorm = false
		} else {
			i = ((bbits[int32(0)] + be[int32(0)]) + ((int32(1023) + (int32(53) - int32(1))) - int32(1)))
			if i > int32(32) {
				x = (javalang.Shl64((int64(DToA_Word0(d))), int((int32(64)-i))) | int64(javalang.Ushr32(DToA_Word1(d), int((i-int32(32))))))
			} else {
				x = javalang.Shl64((int64(DToA_Word1(d))), int((int32(32) - i)))
			}
			d2 = DToA_SetWord0(float64(x), (DToA_Word0(float64(x)) - (int32(31) * int32(1048576))))
			i = (i - (((int32(1023) + (int32(53) - int32(1))) - int32(1)) + int32(1)))
			denorm = true
		}
	}
	ds = ((((d2 - float64(1.5)) * float64(0.289529654602168)) + float64(0.1760912590558)) + (float64(i) * float64(0.301029995663981)))
	k = javalang.D2I(ds)
	if (ds < float64(0.0)) && (ds != float64(k)) {
		k = k - 1
	}
	k_check = true
	if (k >= int32(0)) && (k <= int32(22)) {
		if d < StaticGet_DToA_Tens()[k] {
			k = k - 1
		}
		k_check = false
	}
	j = ((bbits[int32(0)] - i) - int32(1))
	if j >= int32(0) {
		b2 = int32(0)
		s2 = j
	} else {
		b2 = -(j)
		s2 = int32(0)
	}
	if k >= int32(0) {
		b5 = int32(0)
		s5 = k
		s2 = (s2 + k)
	} else {
		b2 = (b2 - k)
		b5 = -(k)
		s5 = int32(0)
	}
	if (mode < int32(0)) || (mode > int32(9)) {
		mode = int32(0)
	}
	try_quick = true
	if mode > int32(5) {
		mode = (mode - int32(4))
		try_quick = false
	}
	leftright = true
	{
		var asgn_24 int32 = int32(0)
		ilim1 = asgn_24
		ilim = asgn_24
	}
	switch mode {
	case int32(0), int32(1):
		{
			var asgn_25 int32 = int32(-1)
			ilim1 = asgn_25
			ilim = asgn_25
		}
		i = int32(18)
		ndigits = int32(0)
	case int32(2):
		leftright = false
		fallthrough
	case int32(4):
		if ndigits <= int32(0) {
			ndigits = int32(1)
		}
		{
			var asgn_26 int32
			{
				var asgn_27 int32 = ndigits
				i = asgn_27
				asgn_26 = asgn_27
			}
			ilim1 = asgn_26
			ilim = asgn_26
		}
	case int32(3):
		leftright = false
		fallthrough
	case int32(5):
		i = ((ndigits + k) + int32(1))
		ilim = i
		ilim1 = (i - int32(1))
		if i <= int32(0) {
			i = int32(1)
		}
	}
	var fast_failed bool = false
	_ = fast_failed
	if ((ilim >= int32(0)) && (ilim <= int32(14))) && try_quick {
		i = int32(0)
		d2 = d
		k0 = k
		ilim0 = ilim
		ieps = int32(2)
		if k > int32(0) {
			ds = StaticGet_DToA_Tens()[(k & int32(15))]
			j = javalang.Shr32(k, int(int32(4)))
			if (j & int32(16)) != int32(0) {
				j = (j & (int32(16) - int32(1)))
				d = (d / StaticGet_DToA_Bigtens()[(int32(5)-int32(1))])
				ieps = ieps + 1
			}
			for ; j != int32(0); func() { j = javalang.Shr32(j, int(int32(1))); i = i + 1 }() {
				if (j & int32(1)) != int32(0) {
					ieps = ieps + 1
					ds = (ds * StaticGet_DToA_Bigtens()[i])
				}
			}
			d = (d / ds)
		} else {
			var iifeCond_1 int32
			{

				var asgn_28 int32 = -(k)
				j1 = asgn_28
				iifeCond_1 = asgn_28
			}
			if (iifeCond_1) != int32(0) {
				d = (d * StaticGet_DToA_Tens()[(j1&int32(15))])
				{
					j = javalang.Shr32(j1, int(int32(4)))
					for ; j != int32(0); func() { j = javalang.Shr32(j, int(int32(1))); i = i + 1 }() {
						if (j & int32(1)) != int32(0) {
							ieps = ieps + 1
							d = (d * StaticGet_DToA_Bigtens()[i])
						}
					}
				}
			}
		}
		if (k_check && (d < float64(1.0))) && (ilim > int32(0)) {
			if ilim1 <= int32(0) {
				fast_failed = true
			} else {
				ilim = ilim1
				k = k - 1
				d = (d * float64(10.0))
				ieps = ieps + 1
			}
		}
		eps = ((float64(ieps) * d) + float64(7.0))
		eps = DToA_SetWord0(eps, (DToA_Word0(eps) - ((int32(53) - int32(1)) * int32(1048576))))
		if ilim == int32(0) {
			{
				var asgn_29 *j2gort.Ext_java_math_BigInteger = nil
				mhi = asgn_29
				l_S = asgn_29
			}
			d = (d - float64(5.0))
			if d > eps {
				buf.Append_1_uint16(uint16(49))
				k = k + 1
				return (k + int32(1))
			}
			if d < -(eps) {
				buf.SetLength(int32(0))
				buf.Append_1_uint16(uint16(48))
				return int32(1)
			}
			fast_failed = true
		}
		if !(fast_failed) {
			fast_failed = true
			if leftright {
				eps = ((float64(0.5) / StaticGet_DToA_Tens()[(ilim-int32(1))]) - eps)
				{
					i = int32(0)
					for {
						l_L = javalang.D2L(d)
						d = (d - float64(l_L))
						buf.Append_1_uint16(uint16((int64(uint16(48)) + l_L)))
						if d < eps {
							return (k + int32(1))
						}
						if (float64(1.0) - d) < eps {
							var lastCh uint16 = uint16(0)
							_ = lastCh
							for {
								lastCh = buf.CharAt((buf.Length() - int32(1)))
								buf.SetLength((buf.Length() - int32(1)))
								if int32(lastCh) != int32(uint16(57)) {
									break
								}
								if buf.Length() == int32(0) {
									k = k + 1
									lastCh = uint16(48)
									break
								}
							}
							buf.Append_1_uint16(uint16((int32(lastCh) + int32(1))))
							return (k + int32(1))
						}
						{
							var iifeCond_1 int32
							{

								old_30 := i
								nw_31 := old_30 + 1
								i = nw_31
								iifeCond_1 = nw_31
							}
							if iifeCond_1 >= ilim {
								break
							}
						}
						eps = (eps * float64(10.0))
						d = (d * float64(10.0))
					}
				}
			} else {
				eps = (eps * StaticGet_DToA_Tens()[(ilim-int32(1))])
				{
					i = int32(1)
					for ; ; func() { i = i + 1; d = (d * float64(10.0)) }() {
						l_L = javalang.D2L(d)
						d = (d - float64(l_L))
						buf.Append_1_uint16(uint16((int64(uint16(48)) + l_L)))
						if i == ilim {
							if d > (float64(0.5) + eps) {
								var lastCh_2 uint16 = uint16(0)
								_ = lastCh_2
								for {
									lastCh_2 = buf.CharAt((buf.Length() - int32(1)))
									buf.SetLength((buf.Length() - int32(1)))
									if int32(lastCh_2) != int32(uint16(57)) {
										break
									}
									if buf.Length() == int32(0) {
										k = k + 1
										lastCh_2 = uint16(48)
										break
									}
								}
								buf.Append_1_uint16(uint16((int32(lastCh_2) + int32(1))))
								return (k + int32(1))
							} else if d < (float64(0.5) - eps) {
								DToA_StripTrailingZeroes(buf)
								return (k + int32(1))
							}
							break
						}
					}
				}
			}
		}
		if fast_failed {
			buf.SetLength(int32(0))
			d = d2
			k = k0
			ilim = ilim0
		}
	}
	if (be[int32(0)] >= int32(0)) && (k <= int32(14)) {
		ds = StaticGet_DToA_Tens()[k]
		if (ndigits < int32(0)) && (ilim <= int32(0)) {
			{
				var asgn_32 *j2gort.Ext_java_math_BigInteger = nil
				mhi = asgn_32
				l_S = asgn_32
			}
			if ((ilim < int32(0)) || (d < (float64(int32(5)) * ds))) || (!(biasUp) && (d == (float64(int32(5)) * ds))) {
				buf.SetLength(int32(0))
				buf.Append_1_uint16(uint16(48))
				return int32(1)
			}
			buf.Append_1_uint16(uint16(49))
			k = k + 1
			return (k + int32(1))
		}
		{
			i = int32(1)
			for ; ; i = i + 1 {
				l_L = javalang.D2L((d / ds))
				d = (d - (float64(l_L) * ds))
				buf.Append_1_uint16(uint16((int64(uint16(48)) + l_L)))
				if i == ilim {
					d = (d + d)
					if (d > ds) || ((d == ds) && (((l_L & int64(int32(1))) != int64(int32(0))) || biasUp)) {
						var lastCh_3 uint16 = uint16(0)
						_ = lastCh_3
						for {
							lastCh_3 = buf.CharAt((buf.Length() - int32(1)))
							buf.SetLength((buf.Length() - int32(1)))
							if int32(lastCh_3) != int32(uint16(57)) {
								break
							}
							if buf.Length() == int32(0) {
								k = k + 1
								lastCh_3 = uint16(48)
								break
							}
						}
						buf.Append_1_uint16(uint16((int32(lastCh_3) + int32(1))))
					}
					break
				}
				d = (d * float64(10.0))
				if d == float64(int32(0)) {
					break
				}
			}
		}
		return (k + int32(1))
	}
	m2 = b2
	m5 = b5
	{
		var asgn_33 *j2gort.Ext_java_math_BigInteger = nil
		mlo = asgn_33
		mhi = asgn_33
	}
	if leftright {
		if mode < int32(2) {
			if denorm {
				i = (be[int32(0)] + (((int32(1023) + (int32(53) - int32(1))) - int32(1)) + int32(1)))
			} else {
				i = ((int32(1) + int32(53)) - bbits[int32(0)])
			}
		} else {
			j = (ilim - int32(1))
			if m5 >= j {
				m5 = (m5 - j)
			} else {
				s5 = (s5 + func() int32 { cv_34 := (j - m5); j = cv_34; return cv_34 }())
				b5 = (b5 + j)
				m5 = int32(0)
			}
			{
				var iifeCond_1 int32
				{

					var asgn_35 int32 = ilim
					i = asgn_35
					iifeCond_1 = asgn_35
				}
				if (iifeCond_1) < int32(0) {
					m2 = (m2 - i)
					i = int32(0)
				}
			}
		}
		b2 = (b2 + i)
		s2 = (s2 + i)
		mhi = j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(1)))
	}
	if (m2 > int32(0)) && (s2 > int32(0)) {
		if m2 < s2 {
			i = m2
		} else {
			i = s2
		}
		b2 = (b2 - i)
		m2 = (m2 - i)
		s2 = (s2 - i)
	}
	if b5 > int32(0) {
		if leftright {
			if m5 > int32(0) {
				mhi = DToA_Pow5mult(mhi, m5)
				b1 = mhi.Multiply(b)
				b = b1
			}
			{
				var iifeCond_1 int32
				{

					var asgn_36 int32 = (b5 - m5)
					j = asgn_36
					iifeCond_1 = asgn_36
				}
				if (iifeCond_1) != int32(0) {
					b = DToA_Pow5mult(b, j)
				}
			}
		} else {
			b = DToA_Pow5mult(b, b5)
		}
	}
	l_S = j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(1)))
	if s5 > int32(0) {
		l_S = DToA_Pow5mult(l_S, s5)
	}
	spec_case = false
	if mode < int32(2) {
		if ((DToA_Word1(d) == int32(0)) && ((DToA_Word0(d) & int32(1048575)) == int32(0))) && ((DToA_Word0(d) & (int32(2146435072) & javalang.Shl32(int32(2146435072), int(int32(1))))) != int32(0)) {
			b2 = (b2 + int32(1))
			s2 = (s2 + int32(1))
			spec_case = true
		}
	}
	var l_S_bytes []int8 = l_S.ToByteArray()
	_ = l_S_bytes
	var l_S_hiWord int32 = int32(0)
	_ = l_S_hiWord
	{
		var idx int32 = int32(0)
		_ = idx
		for ; idx < int32(4); idx = idx + 1 {
			l_S_hiWord = (javalang.Shl32(l_S_hiWord, int(int32(8))))
			if idx < int32(len(l_S_bytes)) {
				l_S_hiWord = (l_S_hiWord | (int32(l_S_bytes[idx]) & int32(255)))
			}
		}
	}
	if (func() int32 {
		var asgn_37 int32 = (((func() int32 {
			if s5 != int32(0) {
				return (int32(32) - DToA_Hi0bits(l_S_hiWord))
			}
			return int32(1)
		}()) + s2) & int32(31))
		i = asgn_37
		return asgn_37
	}()) != int32(0) {
		i = (int32(32) - i)
	}
	if i > int32(4) {
		i = (i - int32(4))
		b2 = (b2 + i)
		m2 = (m2 + i)
		s2 = (s2 + i)
	} else if i < int32(4) {
		i = (i + int32(28))
		b2 = (b2 + i)
		m2 = (m2 + i)
		s2 = (s2 + i)
	}
	if b2 > int32(0) {
		b = b.ShiftLeft(b2)
	}
	if s2 > int32(0) {
		l_S = l_S.ShiftLeft(s2)
	}
	if k_check {
		if b.CompareTo(j2gort.PtrToAny(l_S)) < int32(0) {
			k = k - 1
			b = b.Multiply(j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(10))))
			if leftright {
				mhi = mhi.Multiply(j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(10))))
			}
			ilim = ilim1
		}
	}
	if (ilim <= int32(0)) && (mode > int32(2)) {
		if ((ilim < int32(0)) || ((func() int32 {
			var iifeArg_1 *j2gort.
				Ext_java_math_BigInteger
			{

				var asgn_39 *j2gort.Ext_java_math_BigInteger = l_S.Multiply(j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(5))))
				l_S = asgn_39
				iifeArg_1 = asgn_39
			}
			var asgn_38 int32 = b.CompareTo(j2gort.PtrToAny(iifeArg_1))
			i = asgn_38
			return asgn_38
		}()) < int32(0))) || ((i == int32(0)) && !(biasUp)) {
			buf.SetLength(int32(0))
			buf.Append_1_uint16(uint16(48))
			return int32(1)
		}
		buf.Append_1_uint16(uint16(49))
		k = k + 1
		return (k + int32(1))
	}
	if leftright {
		if m2 > int32(0) {
			mhi = mhi.ShiftLeft(m2)
		}
		mlo = mhi
		if spec_case {
			mhi = mlo
			mhi = mhi.ShiftLeft(int32(1))
		}
		{
			i = int32(1)
			for ; ; i = i + 1 {
				var divResult []*j2gort.Ext_java_math_BigInteger = b.DivideAndRemainder(l_S)
				_ = divResult
				b = divResult[int32(1)]
				dig = uint16((divResult[int32(0)].IntValue() + int32(uint16(48))))
				j = b.CompareTo(j2gort.PtrToAny(mlo))
				delta = l_S.Subtract(mhi)
				if delta.Signum() <= int32(0) {
					j1 = int32(1)
				} else {
					j1 = b.CompareTo(j2gort.PtrToAny(delta))
				}
				if ((j1 == int32(0)) && (mode == int32(0))) && ((DToA_Word1(d) & int32(1)) == int32(0)) {
					if int32(dig) == int32(uint16(57)) {
						buf.Append_1_uint16(uint16(57))
						if DToA_RoundOff(buf) {
							k = k + 1
							buf.Append_1_uint16(uint16(49))
						}
						return (k + int32(1))
					}
					if j > int32(0) {
						dig = dig + 1
					}
					buf.Append_1_uint16(dig)
					return (k + int32(1))
				}
				if (j < int32(0)) || (((j == int32(0)) && (mode == int32(0))) && ((DToA_Word1(d) & int32(1)) == int32(0))) {
					if j1 > int32(0) {
						b = b.ShiftLeft(int32(1))
						j1 = b.CompareTo(j2gort.PtrToAny(l_S))
						if ((j1 > int32(0)) || ((j1 == int32(0)) && (((int32(dig) & int32(1)) == int32(1)) || biasUp))) && (int32(func() uint16 { old_40 := dig; nw_41 := old_40 + 1; dig = nw_41; return old_40 }()) == int32(uint16(57))) {
							buf.Append_1_uint16(uint16(57))
							if DToA_RoundOff(buf) {
								k = k + 1
								buf.Append_1_uint16(uint16(49))
							}
							return (k + int32(1))
						}
					}
					buf.Append_1_uint16(dig)
					return (k + int32(1))
				}
				if j1 > int32(0) {
					if int32(dig) == int32(uint16(57)) {
						buf.Append_1_uint16(uint16(57))
						if DToA_RoundOff(buf) {
							k = k + 1
							buf.Append_1_uint16(uint16(49))
						}
						return (k + int32(1))
					}
					buf.Append_1_uint16(uint16((int32(dig) + int32(1))))
					return (k + int32(1))
				}
				buf.Append_1_uint16(dig)
				if i == ilim {
					break
				}
				b = b.Multiply(j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(10))))
				if mlo == mhi {
					{
						var asgn_42 *j2gort.Ext_java_math_BigInteger = mhi.Multiply(j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(10))))
						mhi = asgn_42
						mlo = asgn_42
					}
				} else {
					mlo = mlo.Multiply(j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(10))))
					mhi = mhi.Multiply(j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(10))))
				}
			}
		}
	} else {
		{
			i = int32(1)
			for ; ; i = i + 1 {
				var divResult_2 []*j2gort.Ext_java_math_BigInteger = b.DivideAndRemainder(l_S)
				_ = divResult_2
				b = divResult_2[int32(1)]
				dig = uint16((divResult_2[int32(0)].IntValue() + int32(uint16(48))))
				buf.Append_1_uint16(dig)
				if i >= ilim {
					break
				}
				b = b.Multiply(j2gort.Ext_java_math_BigInteger_ValueOf(int64(int32(10))))
			}
		}
	}
	b = b.ShiftLeft(int32(1))
	j = b.CompareTo(j2gort.PtrToAny(l_S))
	if (j > int32(0)) || ((j == int32(0)) && (((int32(dig) & int32(1)) == int32(1)) || biasUp)) {
		if DToA_RoundOff(buf) {
			k = k + 1
			buf.Append_1_uint16(uint16(49))
			return (k + int32(1))
		}
	} else {
		DToA_StripTrailingZeroes(buf)
	}
	return (k + int32(1))
}

// Java static method: DToA.stripTrailingZeroes(java.lang.StringBuilder buf)
func DToA_StripTrailingZeroes(buf *j2gort.Ext_java_lang_StringBuilder) {
	clinit_DToA()
	_ = buf
	var bl int32 = buf.Length()
	_ = bl
	for (func() int32 { old_43 := bl; nw_44 := old_43 - 1; bl = nw_44; return old_43 }() > int32(0)) && (int32(buf.CharAt(bl)) == int32(uint16(48))) {
	}
	buf.SetLength((bl + int32(1)))
}

// Java static method: DToA.JS_dtostr(java.lang.StringBuilder buffer, int mode, int precision, double d)
func DToA_JS_dtostr(buffer *j2gort.Ext_java_lang_StringBuilder, mode int32, precision int32, d float64) {
	clinit_DToA()
	_ = buffer
	_ = mode
	_ = precision
	_ = d
	var decPt int32 = int32(0)
	_ = decPt
	var sign []bool = make([]bool, int32(1))
	_ = sign
	var nDigits int32 = int32(0)
	_ = nDigits
	if (mode == int32(2)) && ((d >= float64(1.0e21)) || (d <= -(float64(1.0e21)))) {
		mode = int32(0)
	}
	decPt = DToA_JS_dtoa(d, StaticGet_DToA_DtoaModes()[mode], (mode >= int32(2)), precision, sign, buffer)
	nDigits = buffer.Length()
	if decPt != int32(9999) {
		var exponentialNotation bool = false
		_ = exponentialNotation
		var minNDigits int32 = int32(0)
		_ = minNDigits
		var p int32 = int32(0)
		_ = p
		switch mode {
		case int32(0):
			if (decPt < int32(-5)) || (decPt > int32(21)) {
				exponentialNotation = true
			} else {
				minNDigits = decPt
			}
		case int32(2):
			if precision >= int32(0) {
				minNDigits = (decPt + precision)
			} else {
				minNDigits = decPt
			}
		case int32(3):
			minNDigits = precision
			fallthrough
		case int32(1):
			exponentialNotation = true
		case int32(4):
			minNDigits = precision
			if (decPt < int32(-5)) || (decPt > precision) {
				exponentialNotation = true
			}
		}
		if nDigits < minNDigits {
			p = minNDigits
			nDigits = minNDigits
			for doFirst_45 := true; doFirst_45 || (buffer.Length() != p); doFirst_45 = false {
				buffer.Append_1_uint16(uint16(48))
			}
		}
		if exponentialNotation {
			if nDigits != int32(1) {
				buffer.Insert_2_int32_uint16(int32(1), uint16(46))
			}
			buffer.Append_1_uint16(uint16(101))
			if (decPt - int32(1)) >= int32(0) {
				buffer.Append_1_uint16(uint16(43))
			}
			buffer.Append_1_int32((decPt - int32(1)))
		} else if decPt != nDigits {
			if decPt > int32(0) {
				buffer.Insert_2_int32_uint16(decPt, uint16(46))
			} else {
				{
					var i int32 = int32(0)
					_ = i
					for ; i < (int32(1) - decPt); i = i + 1 {
						buffer.Insert_2_int32_uint16(int32(0), uint16(48))
					}
				}
				buffer.Insert_2_int32_uint16(int32(1), uint16(46))
			}
		}
	}
	if (sign[int32(0)] && !((DToA_Word0(d) == int32(-2147483648)) && (DToA_Word1(d) == int32(0)))) && !(((DToA_Word0(d) & int32(2146435072)) == int32(2146435072)) && ((DToA_Word1(d) != int32(0)) || ((DToA_Word0(d) & int32(1048575)) != int32(0)))) {
		buffer.Insert_2_int32_uint16(int32(0), uint16(45))
	}
}

// Java static method: DToA.numberToString(double value)
func DToA_NumberToString(value float64) string {
	clinit_DToA()
	_ = value
	if math.IsNaN(value) {
		return "NaN"
	}
	if value == float64(0.0) {
		return "0"
	}
	if value < float64(int32(0)) {
		return ("-" + j2gort.StrValString(DToA_NumberToString(-(value))))
	}
	if math.IsInf(value, 0) {
		return "Infinity"
	}
	var buffer *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = buffer
	DToA_JS_dtostr(buffer, int32(0), int32(0), value)
	return buffer.String()
}
