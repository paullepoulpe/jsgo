/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Port of com/google/javascript/jscomp/regex/CharRanges.java (generated mechanically by j2go).

package regex

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.regex.CharRanges
type CharRanges struct {
	ranges []int32 `gson:"ranges"`
}

func (this_ *CharRanges) BindSelf_(s any) {
	_ = s
}
func (this_ *CharRanges) IsJavaNil_() bool { return this_ == nil }
func (*CharRanges) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.regex.CharRanges", "CharRanges"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.regex.CharRanges", "com.google.javascript.jscomp.regex.CharRanges", func(v_ any) bool {
		return j2gort.InstanceOf[*CharRanges](v_)
	})
}

var static_CharRanges_EMPTY *CharRanges

func StaticGet_CharRanges_EMPTY() *CharRanges {
	clinit_CharRanges()
	return static_CharRanges_EMPTY
}
func StaticSet_CharRanges_EMPTY(v *CharRanges) {
	clinit_CharRanges()
	static_CharRanges_EMPTY = v
}

var static_CharRanges_ALL_CODE_UNITS *CharRanges

func StaticGet_CharRanges_ALL_CODE_UNITS() *CharRanges {
	clinit_CharRanges()
	return static_CharRanges_ALL_CODE_UNITS
}
func StaticSet_CharRanges_ALL_CODE_UNITS(v *CharRanges) {
	clinit_CharRanges()
	static_CharRanges_ALL_CODE_UNITS = v
}

var clinitDone_CharRanges uint32
var clinitState_CharRanges j2gort.ClinitState

func clinit_CharRanges() {
	if atomic.LoadUint32(&clinitDone_CharRanges) != 0 {
		return
	}
	clinitSlow_CharRanges()
}

func clinitSlow_CharRanges() {
	if !j2gort.EnterClinit(&clinitState_CharRanges, &clinitDone_CharRanges) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CharRanges, &clinitDone_CharRanges)
	static_CharRanges_EMPTY = NewCharRanges(make([]int32, int32(0)))
	static_CharRanges_ALL_CODE_UNITS = NewCharRanges([]int32{int32(0), int32(65536)})
}

func (this_ *CharRanges) initFields_CharRanges() {
	_ = this_
}

// Java constructor: CharRanges(int[] ranges)
func NewCharRanges(ranges []int32) *CharRanges {
	clinit_CharRanges()
	this_ := &CharRanges{}
	InitCharRanges(this_, ranges)
	return this_
}

func InitCharRanges(this_ *CharRanges, ranges []int32) {
	_ = ranges
	this_.initFields_CharRanges()
	this_.ranges = ranges
}

// Java static method: CharRanges.inclusive(int start, int end)
func CharRanges_Inclusive(start int32, end int32) *CharRanges {
	clinit_CharRanges()
	_ = start
	_ = end
	if start > end {
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.BoxString((j2gort.StrValInt32(start) + " > " + j2gort.StrValInt32(end))))))
	}
	return NewCharRanges([]int32{start, (end + int32(1))})
}

// Java static method: CharRanges.withMembers(int[] members)
func CharRanges_WithMembers(members []int32) *CharRanges {
	clinit_CharRanges()
	_ = members
	return NewCharRanges(CharRanges_IntArrayToRanges(members))
}

// Java static method: CharRanges.withRanges(int[] ranges)
func CharRanges_WithRanges(ranges []int32) *CharRanges {
	clinit_CharRanges()
	_ = ranges
	if (int32(len(ranges)) & int32(1)) != int32(0) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException")))
	}
	{
		var i int32 = int32(1)
		_ = i
		for ; i < int32(len(ranges)); i = i + 1 {
			if ranges[i] <= ranges[(i-int32(1))] {
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.BoxString((j2gort.StrValInt32(ranges[i]) + " > " + j2gort.StrValInt32(ranges[(i-int32(1))]))))))
			}
		}
	}
	return NewCharRanges(ranges)
}

// Java static method: CharRanges.intArrayToRanges(int[] members)
func CharRanges_IntArrayToRanges(members []int32) []int32 {
	clinit_CharRanges()
	_ = members
	var nMembers int32 = int32(len(members))
	_ = nMembers
	if nMembers == int32(0) {
		return make([]int32, int32(0))
	}
	j2gort.Ext_java_util_Arrays_Sort_1_slice_int32(members)
	var nRuns int32 = int32(1)
	_ = nRuns
	{
		var i int32 = int32(1)
		_ = i
		for ; i < nMembers; i = i + 1 {
			var current int32 = members[i]
			_ = current
			var last int32 = members[(i - int32(1))]
			_ = last
			if current == last {
				continue
			}
			if current != (last + int32(1)) {
				nRuns = nRuns + 1
			}
		}
	}
	var ranges []int32 = make([]int32, (nRuns * int32(2)))
	_ = ranges
	{
		a_1 := ranges
		i_2 := int32(0)
		a_1[i_2] = members[int32(0)]
	}
	var k int32 = int32(0)
	_ = k
	{
		var i_2 int32 = int32(1)
		_ = i_2
		for ; (k + int32(2)) < int32(len(ranges)); i_2 = i_2 + 1 {
			var current_2 int32 = members[i_2]
			_ = current_2
			var last_2 int32 = members[(i_2 - int32(1))]
			_ = last_2
			if current_2 == last_2 {
				continue
			}
			if current_2 != (last_2 + int32(1)) {
				{
					a_3 := ranges
					i_4 := func() int32 { old_5 := k; nw_6 := old_5 + 1; k = nw_6; return nw_6 }()
					a_3[i_4] = (last_2 + int32(1))
				}
				{
					a_7 := ranges
					i_8 := func() int32 { old_9 := k; nw_10 := old_9 + 1; k = nw_10; return nw_10 }()
					a_7[i_8] = current_2
				}
			}
		}
	}
	{
		a_11 := ranges
		i_12 := func() int32 { old_13 := k; nw_14 := old_13 + 1; k = nw_14; return nw_14 }()
		a_11[i_12] = (members[(nMembers-int32(1))] + int32(1))
	}
	return ranges
}

// Java method: CharRanges.contains(int bit)
func (this_ *CharRanges) Contains(bit int32) bool {
	_ = bit
	return ((j2gort.Ext_java_util_Arrays_BinarySearch_2_slice_int32_int32(this_.ranges, bit) & int32(1)) == int32(0))
}

// Java method: CharRanges.isEmpty()
func (this_ *CharRanges) IsEmpty() bool {
	return (int32(len(this_.ranges)) == int32(0))
}

// Java method: CharRanges.getNumRanges()
func (this_ *CharRanges) GetNumRanges() int32 {
	return javalang.Shr32(int32(len(this_.ranges)), int(int32(1)))
}

// Java method: CharRanges.start(int i)
func (this_ *CharRanges) Start(i int32) int32 {
	_ = i
	return this_.ranges[javalang.Shl32(i, int(int32(1)))]
}

// Java method: CharRanges.end(int i)
func (this_ *CharRanges) End(i int32) int32 {
	_ = i
	return this_.ranges[((javalang.Shl32(i, int(int32(1)))) | int32(1))]
}

// Java method: CharRanges.union(com.google.javascript.jscomp.regex.CharRanges other)
func (this_ *CharRanges) Union(other *CharRanges) *CharRanges {
	_ = other
	var q []int32 = this_.ranges
	_ = q
	var r []int32 = other.ranges
	_ = r
	var m int32 = int32(len(q))
	_ = m
	var n int32 = int32(len(r))
	_ = n
	if m == int32(0) {
		return other
	}
	if n == int32(0) {
		return this_
	}
	var out []int32 = make([]int32, (m + n))
	_ = out
	var i int32 = int32(0)
	_ = i
	var j int32 = int32(0)
	_ = j
	var k int32 = int32(0)
	_ = k
	for (i < m) && (j < n) {
		var a0 int32 = q[i]
		_ = a0
		var a1 int32 = q[(i + int32(1))]
		_ = a1
		var b0 int32 = r[j]
		_ = b0
		var b1 int32 = r[(j + int32(1))]
		_ = b1
		if a1 < b0 {
			{
				a_15 := out
				i_16 := func() int32 { old_17 := k; nw_18 := old_17 + 1; k = nw_18; return old_17 }()
				a_15[i_16] = a0
			}
			{
				a_19 := out
				i_20 := func() int32 { old_21 := k; nw_22 := old_21 + 1; k = nw_22; return old_21 }()
				a_19[i_20] = a1
			}
			i = (i + int32(2))
		} else if b1 < a0 {
			{
				a_23 := out
				i_24 := func() int32 { old_25 := k; nw_26 := old_25 + 1; k = nw_26; return old_25 }()
				a_23[i_24] = b0
			}
			{
				a_27 := out
				i_28 := func() int32 { old_29 := k; nw_30 := old_29 + 1; k = nw_30; return old_29 }()
				a_27[i_28] = b1
			}
			j = (j + int32(2))
		} else {
			var start int32 = min(a0, b0)
			_ = start
			var end int32 = max(a1, b1)
			_ = end
			i = (i + int32(2))
			j = (j + int32(2))
			for (i < m) || (j < n) {
				if (i < m) && (q[i] <= end) {
					end = max(end, q[(i+int32(1))])
					i = (i + int32(2))
				} else if (j < n) && (r[j] <= end) {
					end = max(end, r[(j+int32(1))])
					j = (j + int32(2))
				} else {
					break
				}
			}
			{
				a_31 := out
				i_32 := func() int32 { old_33 := k; nw_34 := old_33 + 1; k = nw_34; return old_33 }()
				a_31[i_32] = start
			}
			{
				a_35 := out
				i_36 := func() int32 { old_37 := k; nw_38 := old_37 + 1; k = nw_38; return old_37 }()
				a_35[i_36] = end
			}
		}
	}
	if i < m {
		func() {
			acSrc_39 := q
			acSP_40 := i
			acDst_41 := out
			acDP_42 := k
			acLen_43 := (m - i)
			j2gort.CheckNotNullSlice(acSrc_39)
			j2gort.CheckNotNullSlice(acDst_41)
			if acSP_40 < 0 || acDP_42 < 0 || acLen_43 < 0 || int64(acSP_40)+int64(acLen_43) > int64(len(acSrc_39)) || int64(acDP_42)+int64(acLen_43) > int64(len(acDst_41)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_41[int(acDP_42):int(acDP_42)+int(acLen_43)], acSrc_39[int(acSP_40):int(acSP_40)+int(acLen_43)])
		}()
		k = (k + (m - i))
	} else if j < n {
		func() {
			acSrc_44 := r
			acSP_45 := j
			acDst_46 := out
			acDP_47 := k
			acLen_48 := (n - j)
			j2gort.CheckNotNullSlice(acSrc_44)
			j2gort.CheckNotNullSlice(acDst_46)
			if acSP_45 < 0 || acDP_47 < 0 || acLen_48 < 0 || int64(acSP_45)+int64(acLen_48) > int64(len(acSrc_44)) || int64(acDP_47)+int64(acLen_48) > int64(len(acDst_46)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_46[int(acDP_47):int(acDP_47)+int(acLen_48)], acSrc_44[int(acSP_45):int(acSP_45)+int(acLen_48)])
		}()
		k = (k + (n - j))
	}
	if k != int32(len(out)) {
		var clipped []int32 = func() []int32 {
			arrSrc_49 := out
			arrLen_50 := k
			j2gort.CheckNotNullSlice(arrSrc_49)
			if arrLen_50 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_51 := make([]int32, arrLen_50)
			copy(arrOut_51, arrSrc_49)
			return arrOut_51
		}()
		_ = clipped
		out = clipped
	}
	return NewCharRanges(out)
}

// Java method: CharRanges.intersection(com.google.javascript.jscomp.regex.CharRanges other)
func (this_ *CharRanges) Intersection(other *CharRanges) *CharRanges {
	_ = other
	var aRanges []int32 = this_.ranges
	_ = aRanges
	var bRanges []int32 = other.ranges
	_ = bRanges
	var aLen int32 = int32(len(aRanges))
	_ = aLen
	var bLen int32 = int32(len(bRanges))
	_ = bLen
	if aLen == int32(0) {
		return this_
	}
	if bLen == int32(0) {
		return other
	}
	var aIdx int32 = int32(0)
	_ = aIdx
	var bIdx int32 = int32(0)
	_ = bIdx
	var intersection []int32 = make([]int32, min(aLen, bLen))
	_ = intersection
	var intersectionIdx int32 = int32(0)
	_ = intersectionIdx
	var pos int32 = min(aRanges[int32(0)], bRanges[int32(0)])
	_ = pos
	for (aIdx < aLen) && (bIdx < bLen) {
		if aRanges[(aIdx+int32(1))] <= pos {
			aIdx = (aIdx + int32(2))
		} else if bRanges[(bIdx+int32(1))] <= pos {
			bIdx = (bIdx + int32(2))
		} else {
			var start int32 = max(aRanges[aIdx], bRanges[bIdx])
			_ = start
			if pos < start {
				pos = start
			} else {
				var end int32 = min(aRanges[(aIdx+int32(1))], bRanges[(bIdx+int32(1))])
				_ = end
				if (intersectionIdx != int32(0)) && (pos == intersection[(intersectionIdx-int32(1))]) {
					{
						a_52 := intersection
						i_53 := (intersectionIdx - int32(1))
						a_52[i_53] = end
					}
				} else {
					if intersectionIdx == int32(len(intersection)) {
						var newArr []int32 = make([]int32, (intersectionIdx * int32(2)))
						_ = newArr
						func() {
							acSrc_54 := intersection
							acSP_55 := int32(0)
							acDst_56 := newArr
							acDP_57 := int32(0)
							acLen_58 := intersectionIdx
							j2gort.CheckNotNullSlice(acSrc_54)
							j2gort.CheckNotNullSlice(acDst_56)
							if acSP_55 < 0 || acDP_57 < 0 || acLen_58 < 0 || int64(acSP_55)+int64(acLen_58) > int64(len(acSrc_54)) || int64(acDP_57)+int64(acLen_58) > int64(len(acDst_56)) {
								panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
							}
							copy(acDst_56[int(acDP_57):int(acDP_57)+int(acLen_58)], acSrc_54[int(acSP_55):int(acSP_55)+int(acLen_58)])
						}()
						intersection = newArr
					}
					{
						a_59 := intersection
						i_60 := func() int32 { old_61 := intersectionIdx; nw_62 := old_61 + 1; intersectionIdx = nw_62; return old_61 }()
						a_59[i_60] = pos
					}
					{
						a_63 := intersection
						i_64 := func() int32 { old_65 := intersectionIdx; nw_66 := old_65 + 1; intersectionIdx = nw_66; return old_65 }()
						a_63[i_64] = end
					}
				}
				pos = end
			}
		}
	}
	if intersectionIdx != int32(len(intersection)) {
		var newArr_2 []int32 = func() []int32 {
			arrSrc_67 := intersection
			arrLen_68 := intersectionIdx
			j2gort.CheckNotNullSlice(arrSrc_67)
			if arrLen_68 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_69 := make([]int32, arrLen_68)
			copy(arrOut_69, arrSrc_67)
			return arrOut_69
		}()
		_ = newArr_2
		intersection = newArr_2
	}
	return NewCharRanges(intersection)
}

// Java method: CharRanges.difference(com.google.javascript.jscomp.regex.CharRanges subtrahendRanges)
func (this_ *CharRanges) Difference(subtrahendRanges *CharRanges) *CharRanges {
	_ = subtrahendRanges
	var minuend []int32 = this_.ranges
	_ = minuend
	var subtrahend []int32 = subtrahendRanges.ranges
	_ = subtrahend
	var mn int32 = int32(len(minuend))
	_ = mn
	var sn int32 = int32(len(subtrahend))
	_ = sn
	if (mn == int32(0)) || (sn == int32(0)) {
		return this_
	}
	var difference []int32 = make([]int32, int32(len(minuend)))
	_ = difference
	var mIdx int32 = int32(0)
	_ = mIdx
	var sIdx int32 = int32(0)
	_ = sIdx
	var dIdx int32 = int32(0)
	_ = dIdx
	var pos int32 = minuend[int32(0)]
	_ = pos
	for mIdx < mn {
		if pos >= minuend[(mIdx+int32(1))] {
			mIdx = (mIdx + int32(2))
		} else if pos < minuend[mIdx] {
			pos = minuend[mIdx]
		} else if (sIdx < sn) && (pos >= subtrahend[sIdx]) {
			pos = subtrahend[(sIdx + int32(1))]
			sIdx = (sIdx + int32(2))
		} else {
			var end int32
			_ = end
			if sIdx < sn {
				end = min(minuend[(mIdx+int32(1))], subtrahend[sIdx])
			} else {
				end = minuend[(mIdx + int32(1))]
			}
			if (dIdx != int32(0)) && (difference[(dIdx-int32(1))] == pos) {
				{
					a_70 := difference
					i_71 := (dIdx - int32(1))
					a_70[i_71] = pos
				}
			} else {
				if dIdx == int32(len(difference)) {
					var newArr []int32 = make([]int32, (dIdx * int32(2)))
					_ = newArr
					func() {
						acSrc_72 := difference
						acSP_73 := int32(0)
						acDst_74 := newArr
						acDP_75 := int32(0)
						acLen_76 := dIdx
						j2gort.CheckNotNullSlice(acSrc_72)
						j2gort.CheckNotNullSlice(acDst_74)
						if acSP_73 < 0 || acDP_75 < 0 || acLen_76 < 0 || int64(acSP_73)+int64(acLen_76) > int64(len(acSrc_72)) || int64(acDP_75)+int64(acLen_76) > int64(len(acDst_74)) {
							panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
						}
						copy(acDst_74[int(acDP_75):int(acDP_75)+int(acLen_76)], acSrc_72[int(acSP_73):int(acSP_73)+int(acLen_76)])
					}()
					difference = newArr
				}
				{
					a_77 := difference
					i_78 := func() int32 { old_79 := dIdx; nw_80 := old_79 + 1; dIdx = nw_80; return old_79 }()
					a_77[i_78] = pos
				}
				{
					a_81 := difference
					i_82 := func() int32 { old_83 := dIdx; nw_84 := old_83 + 1; dIdx = nw_84; return old_83 }()
					a_81[i_82] = end
				}
			}
			pos = end
		}
	}
	if dIdx != int32(len(difference)) {
		var newArr_2 []int32 = func() []int32 {
			arrSrc_85 := difference
			arrLen_86 := dIdx
			j2gort.CheckNotNullSlice(arrSrc_85)
			if arrLen_86 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_87 := make([]int32, arrLen_86)
			copy(arrOut_87, arrSrc_85)
			return arrOut_87
		}()
		_ = newArr_2
		difference = newArr_2
	}
	return NewCharRanges(difference)
}

// Java method: CharRanges.containsAll(com.google.javascript.jscomp.regex.CharRanges sub)
func (this_ *CharRanges) ContainsAll(sub *CharRanges) bool {
	_ = sub
	var superRanges []int32 = this_.ranges
	_ = superRanges
	var subRanges []int32 = sub.ranges
	_ = subRanges
	var superIdx int32 = int32(0)
	_ = superIdx
	var subIdx int32 = int32(0)
	_ = subIdx
	var superLen int32 = int32(len(superRanges))
	_ = superLen
	var subLen int32 = int32(len(subRanges))
	_ = subLen
	for subIdx < subLen {
		if superIdx == superLen {
			return false
		}
		if superRanges[(superIdx+int32(1))] <= subRanges[subIdx] {
			superIdx = (superIdx + int32(2))
		} else if superRanges[superIdx] > subRanges[subIdx] {
			return false
		} else if superRanges[(superIdx+int32(1))] >= subRanges[(subIdx+int32(1))] {
			subIdx = (subIdx + int32(2))
		} else {
			return false
		}
	}
	return (subIdx == subLen)
}

// Java method: CharRanges.shift(int delta)
func (this_ *CharRanges) Shift(delta int32) *CharRanges {
	_ = delta
	var n int32 = int32(len(this_.ranges))
	_ = n
	if (delta == int32(0)) || (n == int32(0)) {
		return this_
	}
	if delta < int32(0) {
		var lmin int64 = int64((this_.ranges[int32(0)] + delta))
		_ = lmin
		if lmin < int64(int32(-2147483648)) {
			panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException")))
		}
	} else {
		var lmax int64 = int64((this_.ranges[(n-int32(1))] + delta))
		_ = lmax
		if lmax > int64(int32(2147483647)) {
			panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException")))
		}
	}
	var shiftedRanges []int32 = make([]int32, n)
	_ = shiftedRanges
	{
		var i int32 = n
		_ = i
		for func() int32 { old_88 := i; nw_89 := old_88 - 1; i = nw_89; return nw_89 }() >= int32(0) {
			{
				a_90 := shiftedRanges
				i_91 := i
				a_90[i_91] = (this_.ranges[i] + delta)
			}
		}
	}
	return NewCharRanges(shiftedRanges)
}

// Java method: CharRanges.toString()
func (this_ *CharRanges) String() string {
	var sb j2gort.LocalStringBuilder

	sb.AppendChar(uint16(91))
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(this_.ranges)); i = i + 1 {
			if ((i & int32(1)) != int32(0)) && (this_.ranges[i] == (this_.ranges[(i-int32(1))] + int32(1))) {
				continue
			}
			if i != int32(0) {
				sb.AppendChar(func() uint16 {
					if (i & int32(1)) == int32(0) {
						return uint16(32)
					}
					return uint16(45)
				}())
			}
			sb.AppendStr("0x").AppendStr(j2gort.BoxInt_ToString_2((this_.ranges[i] - (i & int32(1))), int32(16)))
		}
	}
	sb.AppendChar(uint16(93))
	return sb.String()
}

// Java method: CharRanges.equals(java.lang.Object o)
func (this_ *CharRanges) Equals(arg_92 any) bool {
	o := arg_92
	_ = o
	var charRanges *CharRanges
	_ = charRanges
	var ok_93 bool
	charRanges, ok_93 = j2gort.TryCast[*CharRanges](o)
	if !ok_93 {
		return false
	}
	return j2gort.Ext_java_util_Arrays_Equals_2_slice_int32_slice_int32(this_.ranges, charRanges.ranges)
}

// Java method: CharRanges.hashCode()
func (this_ *CharRanges) HashCode() int32 {
	var hc int32 = int32(0)
	_ = hc
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = min(int32(16), int32(len(this_.ranges)))
		_ = n
		for ; i < n; i = i + 1 {
			hc = ((javalang.Shl32(hc, int(int32(2)))) + this_.ranges[i])
		}
	}
	return hc
}
