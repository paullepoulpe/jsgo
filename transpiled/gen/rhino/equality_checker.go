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
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/EqualityChecker.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.EqualityChecker
type EqualityChecker struct {
	cacheLeft      []JSType `gson:"cacheLeft"`
	cacheRight     []JSType `gson:"cacheRight"`
	cacheStatus    []int8   `gson:"cacheStatus"`
	cacheHash      []int32  `gson:"cacheHash"`
	cacheTable     []int32  `gson:"cacheTable"`
	cacheSize      int32    `gson:"cacheSize"`
	recursionDepth int32    `gson:"recursionDepth"`
	idHash_        atomic.Int32
	eqMethod       EqualityChecker_EqMethod `gson:"eqMethod"`
	hasRun         bool                     `gson:"hasRun"`
}

func (this_ *EqualityChecker) BindSelf_(s any) {
	_ = s
}
func (this_ *EqualityChecker) IsJavaNil_() bool { return this_ == nil }
func (*EqualityChecker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.EqualityChecker", "EqualityChecker"
}
func (this_ *EqualityChecker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*EqualityChecker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.EqualityChecker", "com.google.javascript.rhino.jstype.EqualityChecker", func(v_ any) bool {
		return j2gort.InstanceOf[*EqualityChecker](v_)
	})
}

func StaticGet_EqualityChecker_POTENTIALLY_CYCLIC_RECURSION_DEPTH() int32 { return int32(20) }
func StaticGet_EqualityChecker_STATUS_PROCESSING() int8                   { return int8(1) }
func StaticGet_EqualityChecker_STATUS_MATCH() int8                        { return int8(2) }
func StaticGet_EqualityChecker_STATUS_NOT_MATCH() int8                    { return int8(3) }
func (this_ *EqualityChecker) initFields_EqualityChecker() {
	_ = this_

}

// Java constructor: EqualityChecker()
func NewEqualityChecker() *EqualityChecker {
	this_ := &EqualityChecker{}
	InitEqualityChecker(this_)
	return this_
}

func InitEqualityChecker(this_ *EqualityChecker) {
	this_.initFields_EqualityChecker()
}

// Java constructor: EqualityChecker(com.google.javascript.rhino.jstype.EqualityChecker.EqMethod eqMethod)
func NewEqualityChecker_1(eqMethod EqualityChecker_EqMethod) *EqualityChecker {
	this_ := &EqualityChecker{}
	InitEqualityChecker_1(this_, eqMethod)
	return this_
}

func InitEqualityChecker_1(this_ *EqualityChecker, eqMethod EqualityChecker_EqMethod) {
	_ = eqMethod
	this_.initFields_EqualityChecker()
	this_.eqMethod = eqMethod
	this_.hasRun = true
}

// Java method: EqualityChecker.setEqMethod(com.google.javascript.rhino.jstype.EqualityChecker.EqMethod x)
func (this_ *EqualityChecker) SetEqMethod(x EqualityChecker_EqMethod) *EqualityChecker {
	_ = x
	this_.priv_EqualityChecker_checkHasNotRun_8863()
	if !(this_.eqMethod == 0) {
		j2gort.CheckState(false)
	}
	this_.eqMethod = x
	return this_
}

// Java method: EqualityChecker.checkHasNotRun()
func (this_ *EqualityChecker) priv_EqualityChecker_checkHasNotRun_8863() {
	if !(!(this_.hasRun)) {
		j2gort.CheckState(false)
	}
}

// Java static method: EqualityChecker.isWildcardOrProxyOrUnion(com.google.javascript.rhino.jstype.JSTypeClass c)
func EqualityChecker_IsWildcardOrProxyOrUnion(c JSTypeClass) bool {
	_ = c
	switch (c).Ordinal() {
	case int32(15) /* JSTypeClass.PROXY_OBJECT */, int32(8) /* JSTypeClass.NAMED */, int32(21) /* JSTypeClass.TEMPLATIZED */, int32(20) /* JSTypeClass.TEMPLATE */, int32(22) /* JSTypeClass.UNION */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: EqualityChecker.canNonIdenticalTypesEqual(com.google.javascript.rhino.jstype.JSType left, com.google.javascript.rhino.jstype.JSType right)
func EqualityChecker_CanNonIdenticalTypesEqual(left JSType, right JSType) bool {
	_ = left
	_ = right
	var lc JSTypeClass = left.GetTypeClass()
	_ = lc
	var rc JSTypeClass = right.GetTypeClass()
	_ = rc
	if EqualityChecker_IsWildcardOrProxyOrUnion(lc) || EqualityChecker_IsWildcardOrProxyOrUnion(rc) {
		return true
	}
	switch (lc).Ordinal() {
	case int32(11) /* JSTypeClass.NO_RESOLVED */ :
		return (rc == StaticGet_JSTypeClass_NO_RESOLVED())
	case int32(16) /* JSTypeClass.RECORD */ :
		return (rc == StaticGet_JSTypeClass_RECORD())
	case int32(6) /* JSTypeClass.FUNCTION */ :
		return (((rc == StaticGet_JSTypeClass_FUNCTION()) && ((j2gort.CastPtr[FunctionType](left)).GetKind() == StaticGet_FunctionType_Kind_ORDINARY())) && ((j2gort.CastPtr[FunctionType](right)).GetKind() == StaticGet_FunctionType_Kind_ORDINARY()))
	case int32(7) /* JSTypeClass.INSTANCE_OBJECT */, int32(5) /* JSTypeClass.ENUM_ELEMENT */ :
		return (((rc == StaticGet_JSTypeClass_INSTANCE_OBJECT()) || (rc == StaticGet_JSTypeClass_ENUM_ELEMENT())) && (!(left.IsResolved()) || !(right.IsResolved())))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: EqualityChecker.check(com.google.javascript.rhino.jstype.EqualityChecker.EqMethod eqMethod, com.google.javascript.rhino.jstype.JSType left, com.google.javascript.rhino.jstype.JSType right)
func EqualityChecker_Check(eqMethod EqualityChecker_EqMethod, left JSType, right JSType) bool {
	_ = eqMethod
	_ = left
	_ = right
	if left == right || ((left == nil || left.IsJavaNil_()) && (right == nil || right.IsJavaNil_())) {
		return true
	}
	if (left == nil || left.IsJavaNil_()) || (right == nil || right.IsJavaNil_()) {
		return false
	}
	if eqMethod == StaticGet_EqualityChecker_EqMethod_DATA_FLOW() {
		var leftUnknown bool = left.IsUnknownType()
		_ = leftUnknown
		var rightUnknown bool = right.IsUnknownType()
		_ = rightUnknown
		if leftUnknown || rightUnknown {
			return (leftUnknown && rightUnknown)
		}
	}
	if !(EqualityChecker_CanNonIdenticalTypesEqual(left, right)) {
		return false
	}
	var val_1 EqualityChecker
	InitEqualityChecker_1(&val_1, eqMethod)
	var checker *EqualityChecker = &val_1
	_ = checker
	return checker.priv_EqualityChecker_areEqualCaching_8866(left, right)
}

// Java static method: EqualityChecker.checkCallType(com.google.javascript.rhino.jstype.EqualityChecker.EqMethod eqMethod, com.google.javascript.rhino.jstype.FunctionType left, com.google.javascript.rhino.jstype.FunctionType right)
func EqualityChecker_CheckCallType(eqMethod EqualityChecker_EqMethod, left *FunctionType, right *FunctionType) bool {
	_ = eqMethod
	_ = left
	_ = right
	if left == right {
		return true
	}
	if (left == nil) || (right == nil) {
		return false
	}
	var val_2 EqualityChecker
	InitEqualityChecker_1(&val_2, eqMethod)
	var checker *EqualityChecker = &val_2
	_ = checker
	r_3 := checker
	r_3.recursionDepth = r_3.recursionDepth + 1
	var result bool = checker.priv_EqualityChecker_areCallTypeEqual_8874(left, right)
	_ = result
	r_4 := checker
	r_4.recursionDepth = r_4.recursionDepth - 1
	return result
}

// Java static method: EqualityChecker.checkParameters(com.google.javascript.rhino.jstype.EqualityChecker.EqMethod eqMethod, com.google.javascript.rhino.jstype.FunctionType left, com.google.javascript.rhino.jstype.FunctionType right)
func EqualityChecker_CheckParameters(eqMethod EqualityChecker_EqMethod, left *FunctionType, right *FunctionType) bool {
	_ = eqMethod
	_ = left
	_ = right
	if left == right {
		return true
	}
	var val_5 EqualityChecker
	InitEqualityChecker_1(&val_5, eqMethod)
	var checker *EqualityChecker = &val_5
	_ = checker
	return checker.priv_EqualityChecker_areParametersEqual_8875(left, right)
}

// Java method: EqualityChecker.check(com.google.javascript.rhino.jstype.JSType left, com.google.javascript.rhino.jstype.JSType right)
func (this_ *EqualityChecker) Check(left JSType, right JSType) bool {
	_ = left
	_ = right
	this_.priv_EqualityChecker_checkHasNotRun_8863()
	this_.hasRun = true
	if left == right || ((left == nil || left.IsJavaNil_()) && (right == nil || right.IsJavaNil_())) {
		return true
	}
	if (left == nil || left.IsJavaNil_()) || (right == nil || right.IsJavaNil_()) {
		return false
	}
	if this_.eqMethod == StaticGet_EqualityChecker_EqMethod_DATA_FLOW() {
		var leftUnknown bool = left.IsUnknownType()
		_ = leftUnknown
		var rightUnknown bool = right.IsUnknownType()
		_ = rightUnknown
		if leftUnknown || rightUnknown {
			return (leftUnknown && rightUnknown)
		}
	}
	if !(EqualityChecker_CanNonIdenticalTypesEqual(left, right)) {
		return false
	}
	return this_.priv_EqualityChecker_areEqualCaching_8866(left, right)
}

// Java method: EqualityChecker.checkParameters(com.google.javascript.rhino.jstype.FunctionType left, com.google.javascript.rhino.jstype.FunctionType right)
func (this_ *EqualityChecker) CheckParameters(left *FunctionType, right *FunctionType) bool {
	_ = left
	_ = right
	this_.priv_EqualityChecker_checkHasNotRun_8863()
	this_.hasRun = true
	return ((left == right) || this_.priv_EqualityChecker_areParametersEqual_8875(left, right))
}

// Java method: EqualityChecker.areEqualCaching(com.google.javascript.rhino.jstype.JSType left, com.google.javascript.rhino.jstype.JSType right)
func (this_ *EqualityChecker) priv_EqualityChecker_areEqualCaching_8866(left JSType, right JSType) bool {
	_ = left
	_ = right
	if left == right || ((left == nil || left.IsJavaNil_()) && (right == nil || right.IsJavaNil_())) {
		return true
	}
	if (left == nil || left.IsJavaNil_()) || (right == nil || right.IsJavaNil_()) {
		return false
	}
	if (this_.recursionDepth > int32(20)) && (this_.cacheLeft == nil) {
		this_.cacheLeft = make([]JSType, int32(8))
		this_.cacheRight = make([]JSType, int32(8))
		this_.cacheStatus = make([]int8, int32(8))
	}
	if this_.cacheLeft == nil {
		this_.recursionDepth = this_.recursionDepth + 1
		var result bool = this_.priv_EqualityChecker_areEqualInternal_8871(left, right)
		_ = result
		this_.recursionDepth = this_.recursionDepth - 1
		return result
	}
	var existingIdx int32 = this_.priv_EqualityChecker_findCachedEntry_8868(left, right)
	_ = existingIdx
	if existingIdx >= int32(0) {
		var st int8 = this_.cacheStatus[existingIdx]
		_ = st
		if int32(st) == int32(int8(1)) {
			{
				a_6 := this_.cacheStatus
				i_7 := existingIdx
				a_6[i_7] = int8(2)
			}
			return true
		}
		return (int32(st) == int32(int8(2)))
	}
	var idx int32 = this_.priv_EqualityChecker_insertCachedEntry_8869(left, right)
	_ = idx
	var result_2 bool = this_.priv_EqualityChecker_areEqualInternal_8871(left, right)
	_ = result_2
	{
		a_8 := this_.cacheStatus
		i_9 := idx
		a_8[i_9] = func() int8 {
			if result_2 {
				return int8(2)
			}
			return int8(3)
		}()
	}
	return result_2
}

// Java method: EqualityChecker.matchesCachedPair(int idx, com.google.javascript.rhino.jstype.JSType left, com.google.javascript.rhino.jstype.JSType right)
func (this_ *EqualityChecker) priv_EqualityChecker_matchesCachedPair_8867(idx int32, left JSType, right JSType) bool {
	_ = idx
	_ = left
	_ = right
	var aL JSType = this_.cacheLeft[idx]
	_ = aL
	var aR JSType = this_.cacheRight[idx]
	_ = aR
	return (((aL == left || ((aL == nil || aL.IsJavaNil_()) && (left == nil || left.IsJavaNil_()))) && (aR == right || ((aR == nil || aR.IsJavaNil_()) && (right == nil || right.IsJavaNil_())))) || ((aL == right || ((aL == nil || aL.IsJavaNil_()) && (right == nil || right.IsJavaNil_()))) && (aR == left || ((aR == nil || aR.IsJavaNil_()) && (left ==

		// Java method: EqualityChecker.findCachedEntry(com.google.javascript.rhino.jstype.JSType left, com.google.javascript.rhino.jstype.JSType right)
		nil || left.IsJavaNil_())))))
}

func (this_ *EqualityChecker) priv_EqualityChecker_findCachedEntry_8868(left JSType, right JSType) int32 {
	_ = left
	_ = right
	var table []int32 = this_.cacheTable
	_ = table
	if table == nil {
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = this_.cacheSize
			_ = n
			for ; i < n; i = i + 1 {
				if this_.priv_EqualityChecker_matchesCachedPair_8867(i, left, right) {
					return i
				}
			}
		}
		return int32(-1)
	}
	var h int32 = (j2gort.IdentityHash(left) ^ j2gort.IdentityHash(right))
	_ = h
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & mask)
	_ = slot
	for {
		var entry int32 = table[slot]
		_ = entry
		if entry == int32(0) {
			return int32(-1)
		}
		var i_2 int32 = (entry - int32(1))
		_ = i_2
		if (this_.cacheHash[i_2] == h) && this_.priv_EqualityChecker_matchesCachedPair_8867(i_2, left, right) {
			return i_2
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: EqualityChecker.insertCachedEntry(com.google.javascript.rhino.jstype.JSType left, com.google.javascript.rhino.jstype.JSType right)
func (this_ *EqualityChecker) priv_EqualityChecker_insertCachedEntry_8869(left JSType, right JSType) int32 {
	_ = left
	_ = right
	var idx int32 = this_.cacheSize
	_ = idx
	if idx == int32(len(this_.cacheLeft)) {
		var newCap int32 = (idx * int32(2))
		_ = newCap
		this_.cacheLeft = func() []JSType {
			arrSrc_10 := this_.cacheLeft
			arrLen_11 := newCap
			j2gort.CheckNotNullSlice(arrSrc_10)
			if arrLen_11 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_12 := make([]JSType, arrLen_11)
			copy(arrOut_12, arrSrc_10)
			return arrOut_12
		}()
		this_.cacheRight = func() []JSType {
			arrSrc_13 := this_.cacheRight
			arrLen_14 := newCap
			j2gort.CheckNotNullSlice(arrSrc_13)
			if arrLen_14 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_15 := make([]JSType, arrLen_14)
			copy(arrOut_15, arrSrc_13)
			return arrOut_15
		}()
		this_.cacheStatus = func() []int8 {
			arrSrc_16 := this_.cacheStatus
			arrLen_17 := newCap
			j2gort.CheckNotNullSlice(arrSrc_16)
			if arrLen_17 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_18 := make([]int8, arrLen_17)
			copy(arrOut_18, arrSrc_16)
			return arrOut_18
		}()
		if this_.cacheHash != nil {
			this_.cacheHash = func() []int32 {
				arrSrc_19 := this_.cacheHash
				arrLen_20 := newCap
				j2gort.CheckNotNullSlice(arrSrc_19)
				if arrLen_20 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_21 := make([]int32, arrLen_20)
				copy(arrOut_21, arrSrc_19)
				return arrOut_21
			}()
		}
	}
	{
		a_22 := this_.cacheLeft
		i_23 := idx
		a_22[i_23] = left
	}
	{
		a_24 := this_.cacheRight
		i_25 := idx
		a_24[i_25] = right
	}
	{
		a_26 := this_.cacheStatus
		i_27 := idx
		a_26[i_27] = int8(1)
	}
	this_.cacheSize = (idx + int32(1))
	if this_.cacheSize > int32(8) {
		if this_.cacheTable == nil {
			this_.cacheHash = make([]int32, int32(len(this_.cacheLeft)))
			{
				var i int32 = int32(0)
				_ = i
				for ; i < this_.cacheSize; i = i + 1 {
					{
						a_28 := this_.cacheHash
						i_29 := i
						a_28[i_29] = (j2gort.IdentityHash(this_.cacheLeft[i]) ^ j2gort.IdentityHash(this_.cacheRight[i]))
					}
				}
			}
			this_.priv_EqualityChecker_rehashCacheTable_8870(int32(32))
		} else {
			var h int32 = (j2gort.IdentityHash(left) ^ j2gort.IdentityHash(right))
			_ = h
			{
				a_30 := this_.cacheHash
				i_31 := idx
				a_30[i_31] = h
			}
			if (this_.cacheSize * int32(4)) >= (int32(len(this_.cacheTable)) * int32(3)) {
				this_.priv_EqualityChecker_rehashCacheTable_8870((int32(len(this_.cacheTable)) * int32(2)))
			} else {
				var mask int32 = (int32(len(this_.cacheTable)) - int32(1))
				_ = mask
				var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & mask)
				_ = slot
				for this_.cacheTable[slot] != int32(0) {
					slot = ((slot + int32(1)) & mask)
				}
				{
					a_32 := this_.cacheTable
					i_33 := slot
					a_32[i_33] = (idx + int32(1))
				}
			}
		}
	}
	return idx
}

// Java method: EqualityChecker.rehashCacheTable(int newTableLen)
func (this_ *EqualityChecker) priv_EqualityChecker_rehashCacheTable_8870(newTableLen int32) {
	_ = newTableLen
	var newTable []int32 = make([]int32, newTableLen)
	_ = newTable
	var mask int32 = (newTableLen - int32(1))
	_ = mask
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = this_.cacheSize
		_ = n
		for ; i < n; i = i + 1 {
			var h int32 = this_.cacheHash[i]
			_ = h
			var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & mask)
			_ = slot
			for newTable[slot] != int32(0) {
				slot = ((slot + int32(1)) & mask)
			}
			{
				a_34 := newTable
				i_35 := slot
				a_34[i_35] = (i + int32(1))
			}
		}
	}
	this_.cacheTable = newTable
}

// Java method: EqualityChecker.areEqualInternal(com.google.javascript.rhino.jstype.JSType left, com.google.javascript.rhino.jstype.JSType right)
func (this_ *EqualityChecker) priv_EqualityChecker_areEqualInternal_8871(left JSType, right JSType) bool {
	_ = left
	_ = right
	if left == right || ((left == nil || left.IsJavaNil_()) && (right == nil || right.IsJavaNil_())) {
		return true
	} else if (left == nil || left.IsJavaNil_()) || (right == nil || right.IsJavaNil_()) {
		return false
	}
	if left.IsNoResolvedType() && right.IsNoResolvedType() {
		return (left.ToObjectType().GetReferenceName() == right.ToObjectType().GetReferenceName())
	}
	var leftUnknown bool = left.IsUnknownType()
	_ = leftUnknown
	var rightUnknown bool = right.IsUnknownType()
	_ = rightUnknown
	if leftUnknown || rightUnknown {
		if this_.eqMethod == StaticGet_EqualityChecker_EqMethod_DATA_FLOW() {
			return (leftUnknown && rightUnknown)
		} else if (leftUnknown && rightUnknown) && ((left.IsNominalType()) != (right.IsNominalType())) {
			return false
		}
	}
	if !(EqualityChecker_CanNonIdenticalTypesEqual(left, right)) {
		return false
	}
	if left.IsUnionType() && right.IsUnionType() {
		return this_.priv_EqualityChecker_areUnionEqual_8872(left.ToMaybeUnionType(), right.ToMaybeUnionType())
	} else if left.IsUnionType() {
		var leftUnion *UnionType = left.ToMaybeUnionType()
		_ = leftUnion
		if leftUnion.GetAlternatesCount() == int32(1) {
			return this_.priv_EqualityChecker_areEqualInternal_8871(leftUnion.GetAlternateAt(int32(0)), right)
		}
	} else if right.IsUnionType() {
		var rightUnion *UnionType = right.ToMaybeUnionType()
		_ = rightUnion
		if rightUnion.GetAlternatesCount() == int32(1) {
			return this_.priv_EqualityChecker_areEqualInternal_8871(left, rightUnion.GetAlternateAt(int32(0)))
		}
	}
	if left.IsFunctionType() && right.IsFunctionType() {
		return this_.priv_EqualityChecker_areFunctionEqual_8873(left.ToMaybeFunctionType(), right.ToMaybeFunctionType())
	}
	if !(this_.priv_EqualityChecker_areTypeMapEqual_8877(left.GetTemplateTypeMap(), right.GetTemplateTypeMap())) {
		return false
	}
	if left.IsRecordType() && right.IsRecordType() {
		return this_.priv_EqualityChecker_areRecordEqual_8876(left.ToMaybeRecordType(), right.ToMaybeRecordType())
	}
	if left.IsNominalType() && right.IsNominalType() {
		var leftUnwrapped ObjectType = EqualityChecker_UnwrapNominalTypeProxies(left.ToObjectType())
		_ = leftUnwrapped
		var rightUnwrapped ObjectType = EqualityChecker_UnwrapNominalTypeProxies(right.ToObjectType())
		_ = rightUnwrapped
		if !(leftUnwrapped.IsNominalType() && rightUnwrapped.IsNominalType()) {
			j2gort.CheckState(false)
		}
		if left.IsResolved() && right.IsResolved() {
			return leftUnwrapped == rightUnwrapped || ((leftUnwrapped == nil || leftUnwrapped.IsJavaNil_()) && (rightUnwrapped == nil || rightUnwrapped.IsJavaNil_()))
		} else {
			var nameOfleft string = j2gort.CheckNotNullString(leftUnwrapped.GetReferenceName())
			_ = nameOfleft
			var nameOfright string = j2gort.CheckNotNullString(rightUnwrapped.GetReferenceName())
			_ = nameOfright
			var googModuleIdOfLeft string = EqualityChecker_GetGoogModuleId(leftUnwrapped)
			_ = googModuleIdOfLeft
			var googModuleIdOfRight string = EqualityChecker_GetGoogModuleId(rightUnwrapped)
			_ = googModuleIdOfRight
			return ((nameOfleft == nameOfright) && (googModuleIdOfLeft == googModuleIdOfRight))
		}
	}
	var proxyObjectType ProxyObjectType
	_ = proxyObjectType
	var ok_36 bool
	proxyObjectType, ok_36 = j2gort.TryCast[ProxyObjectType](left)
	if ok_36 && (!(j2gort.InstanceOf[*TemplateType](left))) {
		return this_.priv_EqualityChecker_areEqualCaching_8866(proxyObjectType.GetReferencedTypeInternal(), right)
	}
	var proxyObjectType_2 ProxyObjectType
	_ = proxyObjectType_2
	var ok_37 bool
	proxyObjectType_2, ok_37 = j2gort.TryCast[ProxyObjectType](right)
	if ok_37 && (!(j2gort.InstanceOf[*TemplateType](right))) {
		return this_.priv_EqualityChecker_areEqualCaching_8866(left, proxyObjectType_2.GetReferencedTypeInternal())
	}
	return false
}

// Java method: EqualityChecker.areUnionEqual(com.google.javascript.rhino.jstype.UnionType left, com.google.javascript.rhino.jstype.UnionType right)
func (this_ *EqualityChecker) priv_EqualityChecker_areUnionEqual_8872(left *UnionType, right *UnionType) bool {
	_ = left
	_ = right
	var leftSize int32 = left.GetAlternatesCount()
	_ = leftSize
	var rightSize int32 = right.GetAlternatesCount()
	_ = rightSize
	if (this_.eqMethod == StaticGet_EqualityChecker_EqMethod_IDENTITY()) && (leftSize != rightSize) {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
	outer:
		for ; i < rightSize; i = i + 1 {
			var rightAlt JSType = right.GetAlternateAt(i)
			_ = rightAlt
			{
				var k int32 = int32(0)
				_ = k
				for ; k < leftSize; k = k + 1 {
					var leftAlt JSType = left.GetAlternateAt(k)
					_ = leftAlt
					if this_.priv_EqualityChecker_areEqualCaching_8866(leftAlt, rightAlt) {
						continue outer
					}
				}
			}
			return false
		}
	}
	return true
}

// Java static method: EqualityChecker.getGoogModuleId(com.google.javascript.rhino.jstype.ObjectType type)
func EqualityChecker_GetGoogModuleId(type_ ObjectType) string {
	_ = type_
	var functionType *FunctionType
	_ = functionType
	var ok_38 bool
	functionType, ok_38 = j2gort.TryCast[*FunctionType](type_)
	if ok_38 {
		return functionType.GetGoogModuleId()
	}
	if type_.GetConstructor() != nil {
		return type_.GetConstructor().GetGoogModuleId()
	}
	return j2gort.NullString
}

// Java method: EqualityChecker.areFunctionEqual(com.google.javascript.rhino.jstype.FunctionType left, com.google.javascript.rhino.jstype.FunctionType right)
func (this_ *EqualityChecker) priv_EqualityChecker_areFunctionEqual_8873(left *FunctionType, right *FunctionType) bool {
	_ = left
	_ = right
	if left == right {
		return true
	}
	if !(left.GetKind() == right.GetKind()) {
		return false
	}
	switch (left.GetKind()).Ordinal() {
	case int32(1) /* Kind.CONSTRUCTOR */, int32(2) /* Kind.INTERFACE */ :
		return false
	case int32(0) /* Kind.ORDINARY */ :
		return ((this_.priv_EqualityChecker_areEqualCaching_8866(left.GetTypeOfThis(), right.GetTypeOfThis()) && this_.priv_EqualityChecker_areCallTypeEqual_8874(left, right)) && (left.GetClosurePrimitive() == right.GetClosurePrimitive()))
	default:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	panic("unreachable switch expression")
}

// Java method: EqualityChecker.areCallTypeEqual(com.google.javascript.rhino.jstype.FunctionType left, com.google.javascript.rhino.jstype.FunctionType right)
func (this_ *EqualityChecker) priv_EqualityChecker_areCallTypeEqual_8874(left *FunctionType, right *FunctionType) bool {
	_ = left
	_ = right
	return (this_.priv_EqualityChecker_areEqualCaching_8866(left.GetReturnType(), right.GetReturnType()) && this_.priv_EqualityChecker_areParametersEqual_8875(left, right))
}

// Java method: EqualityChecker.areParametersEqual(com.google.javascript.rhino.jstype.FunctionType left, com.google.javascript.rhino.jstype.FunctionType right)
func (this_ *EqualityChecker) priv_EqualityChecker_areParametersEqual_8875(left *FunctionType, right *FunctionType) bool {
	_ = left
	_ = right
	var leftParams []FunctionType_Parameter = left.GetParameterArray()
	_ = leftParams
	var rightParams []FunctionType_Parameter = right.GetParameterArray()
	_ = rightParams
	var size int32 = int32(len(leftParams))
	_ = size
	if size != int32(len(rightParams)) {
		return false
	}
	if j2gort.SliceRefEquals(leftParams, rightParams) {
		return true
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < size; i = i + 1 {
			var leftParam FunctionType_Parameter = leftParams[i]
			_ = leftParam
			var rightParam FunctionType_Parameter = rightParams[i]
			_ = rightParam
			var leftParamType JSType = leftParam.GetJSType()
			_ = leftParamType
			var rightParamType JSType = rightParam.GetJSType()
			_ = rightParamType
			if leftParamType != nil && !leftParamType.IsJavaNil_() {
				if (rightParamType != nil && !rightParamType.IsJavaNil_()) && !(this_.priv_EqualityChecker_areEqualCaching_8866(leftParamType, rightParamType)) {
					return false
				}
			} else {
				if rightParamType != nil && !rightParamType.IsJavaNil_() {
					return false
				}
			}
			if leftParam.IsOptional() != rightParam.IsOptional() {
				return false
			}
			if leftParam.IsVariadic() != rightParam.IsVariadic() {
				return false
			}
		}
	}
	return true
}

// Java method: EqualityChecker.areRecordEqual(com.google.javascript.rhino.jstype.RecordType left, com.google.javascript.rhino.jstype.RecordType right)
func (this_ *EqualityChecker) priv_EqualityChecker_areRecordEqual_8876(left *RecordType, right *RecordType) bool {
	_ = left
	_ = right
	var leftMap *PropertyMap = left.GetOwnPropertyMapOrNull()
	_ = leftMap
	var rightMap *PropertyMap = right.GetOwnPropertyMapOrNull()
	_ = rightMap
	var leftCount int32
	_ = leftCount
	if leftMap != nil {
		leftCount = leftMap.GetOwnPropertiesCount()
	} else {
		leftCount = int32(0)
	}
	var rightCount int32
	_ = rightCount
	if rightMap != nil {
		rightCount = rightMap.GetOwnPropertiesCount()
	} else {
		rightCount = int32(0)
	}
	if leftCount != rightCount {
		return false
	}
	if leftCount == int32(0) {
		return true
	}
	if leftCount == int32(1) {
		var lp0 *Property = leftMap.GetInlineProp0()
		_ = lp0
		var rp0 *Property = rightMap.GetInlineProp0()
		_ = rp0
		if (lp0 != nil) && (rp0 != nil) {
			var k string = lp0.GetName()
			_ = k
			return (func() bool { sr_39 := k; se_40 := rp0.GetName(); return j2gort.CheckNotNullString(sr_39) == se_40 }() && this_.priv_EqualityChecker_areEqualCaching_8866(left.GetPropertyType_1_string(k), right.GetPropertyType_1_string(k)))
		}
	} else if leftCount == int32(2) {
		var lp0_2 *Property = leftMap.GetInlineProp0()
		_ = lp0_2
		var lp1 *Property = leftMap.GetInlineProp1()
		_ = lp1
		var rp0_2 *Property = rightMap.GetInlineProp0()
		_ = rp0_2
		var rp1 *Property = rightMap.GetInlineProp1()
		_ = rp1
		if (((lp0_2 != nil) && (lp1 != nil)) && (rp0_2 != nil)) && (rp1 != nil) {
			var lk0 string = lp0_2.GetName()
			_ = lk0
			var lk1 string = lp1.GetName()
			_ = lk1
			var rk0 string = rp0_2.GetName()
			_ = rk0
			var rk1 string = rp1.GetName()
			_ = rk1
			if ((j2gort.CheckNotNullString(lk0) == rk0) && (j2gort.CheckNotNullString(lk1) == rk1)) || ((j2gort.CheckNotNullString(lk0) == rk1) && (j2gort.CheckNotNullString(lk1) == rk0)) {
				return (this_.priv_EqualityChecker_areEqualCaching_8866(left.GetPropertyType_1_string(lk0), right.GetPropertyType_1_string(lk0)) && this_.priv_EqualityChecker_areEqualCaching_8866(left.GetPropertyType_1_string(lk1), right.GetPropertyType_1_string(lk1)))
			}
			return false
		}
	}
	if leftMap.HasFlatProps() && rightMap.HasFlatProps() {
		{
			var i int32 = int32(0)
			_ = i
			for ; i < leftCount; i = i + 1 {
				if rightMap.GetOwnProperty_1_string(leftMap.GetFlatPropAt(i).GetName()) == nil {
					return false
				}
			}
		}
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < leftCount; i_2 = i_2 + 1 {
				var key string = leftMap.GetFlatPropAt(i_2).GetName()
				_ = key
				if !(this_.priv_EqualityChecker_areEqualCaching_8866(left.GetPropertyType_1_string(key), right.GetPropertyType_1_string(key))) {
					return false
				}
			}
		}
		return true
	}
	var leftKeys j2gort.Ext_java_util_Set = left.GetOwnPropertyNames()
	_ = leftKeys
	var rightKeys j2gort.Ext_java_util_Set = right.GetOwnPropertyNames()
	_ = rightKeys
	if !(rightKeys.Equals(leftKeys)) {
		return false
	}
	{
		cur_41 := j2gort.NewIterCursor(j2gort.CheckNotNull(leftKeys))
		for cur_41.HasNext() {
			key_2 := cur_41.NextStr()
			_ = key_2
			if !(this_.priv_EqualityChecker_areEqualCaching_8866(left.GetPropertyType_1_string(key_2), right.GetPropertyType_1_string(key_2))) {
				return false
			}
		}
	}
	return true
}

// Java method: EqualityChecker.areTypeMapEqual(com.google.javascript.rhino.jstype.TemplateTypeMap left, com.google.javascript.rhino.jstype.TemplateTypeMap right)
func (this_ *EqualityChecker) priv_EqualityChecker_areTypeMapEqual_8877(left *TemplateTypeMap, right *TemplateTypeMap) bool {
	_ = left
	_ = right
	var leftKeysLen int32 = left.Size()
	_ = leftKeysLen
	var rightKeysLen int32 = right.Size()
	_ = rightKeysLen
	var leftValsLen int32 = left.GetTemplateValuesSize()
	_ = leftValsLen
	var rightValsLen int32 = right.GetTemplateValuesSize()
	_ = rightValsLen
	{
		var i int32 = int32(0)
		_ = i
	outer:
		for ; i < leftKeysLen; i = i + 1 {
			var leftKey *TemplateType = left.GetTemplateKeyAt(i)
			_ = leftKey
			{
				var j int32 = int32(0)
				_ = j
			inner:
				for ; j < rightKeysLen; j = j + 1 {
					var rightKey *TemplateType = right.GetTemplateKeyAt(j)
					_ = rightKey
					if !(leftKey == rightKey) {
						continue inner
					}
					var leftValue JSType
					_ = leftValue
					if i < leftValsLen {
						leftValue = left.GetTemplateValueAt(i)
					} else {
						leftValue = left.DefaultValueType(leftKey)
					}
					var rightValue JSType
					_ = rightValue
					if j < rightValsLen {
						rightValue = right.GetTemplateValueAt(j)
					} else {
						rightValue = right.DefaultValueType(rightKey)
					}
					if this_.priv_EqualityChecker_areEqualCaching_8866(leftValue, rightValue) {
						continue outer
					}
				}
			}
			return false
		}
	}
	return true
}

// Java static method: EqualityChecker.unwrapNominalTypeProxies(com.google.javascript.rhino.jstype.ObjectType objType)
func EqualityChecker_UnwrapNominalTypeProxies(objType ObjectType) ObjectType {
	_ = objType
	if !(objType.IsResolved()) || (!(objType.IsNamedType()) && !(objType.IsTemplatizedType())) {
		return objType
	}
	var internal ObjectType
	_ = internal
	if objType.IsNamedType() {
		internal = objType.ToMaybeNamedType().GetReferencedObjTypeInternal()
	} else {
		internal = objType.ToMaybeTemplatizedType().GetReferencedObjTypeInternal()
	}
	return EqualityChecker_UnwrapNominalTypeProxies(internal)
}

// Java enum: com.google.javascript.rhino.jstype.EqualityChecker.EqMethod
type EqualityChecker_EqMethod uint8

const (
	enum_EqualityChecker_EqMethod_IDENTITY  EqualityChecker_EqMethod = 1
	enum_EqualityChecker_EqMethod_DATA_FLOW EqualityChecker_EqMethod = 2
)

var enumBaseTable_EqualityChecker_EqMethod = [...]j2gort.EnumBase{
	{Name_: "IDENTITY", Ordinal_: 0},
	{Name_: "DATA_FLOW", Ordinal_: 1},
}
var enumValues_EqualityChecker_EqMethod = []EqualityChecker_EqMethod{1, 2}

func StaticGet_EqualityChecker_EqMethod_IDENTITY() EqualityChecker_EqMethod {
	return enum_EqualityChecker_EqMethod_IDENTITY
}
func StaticGet_EqualityChecker_EqMethod_DATA_FLOW() EqualityChecker_EqMethod {
	return enum_EqualityChecker_EqMethod_DATA_FLOW
}

func EqualityChecker_EqMethod_Values() []EqualityChecker_EqMethod {
	return j2gort.CloneSlice(enumValues_EqualityChecker_EqMethod)
}

func EqualityChecker_EqMethod_ValueOf(name string) EqualityChecker_EqMethod {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_EqualityChecker_EqMethod {
		if enumBaseTable_EqualityChecker_EqMethod[i_].Name_ == name {
			return EqualityChecker_EqMethod(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.EqualityChecker.EqMethod." + name))
}

func (this_ EqualityChecker_EqMethod) IsJavaNil_() bool                              { return this_ == 0 }
func (EqualityChecker_EqMethod) HasIdentityEquals_()                                 {}
func (EqualityChecker_EqMethod) Implements_j2gort_Ext_java_lang_Enum()               {}
func (EqualityChecker_EqMethod) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ EqualityChecker_EqMethod) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (EqualityChecker_EqMethod) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.EqualityChecker$EqMethod", "EqMethod"
}
func (this_ EqualityChecker_EqMethod) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_EqualityChecker_EqMethod[this_-1]
}
func (this_ EqualityChecker_EqMethod) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ EqualityChecker_EqMethod) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ EqualityChecker_EqMethod) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_EqualityChecker_EqMethod[this_-1].Name_
}
func (this_ EqualityChecker_EqMethod) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_EqualityChecker_EqMethod[this_-1].Name_
}
func (this_ EqualityChecker_EqMethod) Equals(other_ any) bool {
	o_, ok_ := other_.(EqualityChecker_EqMethod)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ EqualityChecker_EqMethod) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_EqualityChecker_EqMethod[this_-1].HashCode()
}
func (this_ EqualityChecker_EqMethod) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_EqualityChecker_EqMethod[this_-1].CompareTo(other_)
}
func (this_ EqualityChecker_EqMethod) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_EqualityChecker_EqMethod[this_-1].Name_)
}
func init() {
	vals_ := []any{EqualityChecker_EqMethod(1), EqualityChecker_EqMethod(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.EqualityChecker.EqMethod", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.EqualityChecker$EqMethod", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.EqualityChecker.EqMethod", "com.google.javascript.rhino.jstype.EqualityChecker$EqMethod", func(v_ any) bool {
		return j2gort.InstanceOf[EqualityChecker_EqMethod](v_)
	})
}
