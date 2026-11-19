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
// Port of com/google/javascript/rhino/jstype/SubtypeChecker.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.SubtypeChecker
type SubtypeChecker struct {
	initialSupertype         JSType          `gson:"initialSupertype"`
	initialSubtype           JSType          `gson:"initialSubtype"`
	registry                 *JSTypeRegistry `gson:"registry"`
	cacheSub                 []JSType        `gson:"cacheSub"`
	cacheSup                 []JSType        `gson:"cacheSup"`
	cacheHash                []int32         `gson:"cacheHash"`
	cacheStatus              []int8          `gson:"cacheStatus"`
	cacheTable               []int32         `gson:"cacheTable"`
	cacheSize                int32           `gson:"cacheSize"`
	recursionDepth           int32           `gson:"recursionDepth"`
	idHash_                  atomic.Int32
	isUsingStructuralTyping  bool                 `gson:"isUsingStructuralTyping"`
	hasUsingStructuralTyping bool                 `gson:"hasUsingStructuralTyping"`
	subtypingMode            JSType_SubtypingMode `gson:"subtypingMode"`
	hasRun                   bool                 `gson:"hasRun"`
}

func (this_ *SubtypeChecker) BindSelf_(s any) {
	_ = s
}
func (this_ *SubtypeChecker) IsJavaNil_() bool { return this_ == nil }
func (*SubtypeChecker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.SubtypeChecker", "SubtypeChecker"
}
func (this_ *SubtypeChecker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SubtypeChecker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.SubtypeChecker", "com.google.javascript.rhino.jstype.SubtypeChecker", func(v_ any) bool {
		return j2gort.InstanceOf[*SubtypeChecker](v_)
	})
}

var static_SubtypeChecker_BIVARIANT_TYPES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_SubtypeChecker_BIVARIANT_TYPES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_SubtypeChecker()
	return static_SubtypeChecker_BIVARIANT_TYPES
}
func StaticSet_SubtypeChecker_BIVARIANT_TYPES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_SubtypeChecker()
	static_SubtypeChecker_BIVARIANT_TYPES = v
}
func StaticGet_SubtypeChecker_POTENTIALLY_CYCLIC_RECURSION_DEPTH() int32 { return int32(20) }
func StaticGet_SubtypeChecker_STATUS_PROCESSING() int8                   { return int8(1) }
func StaticGet_SubtypeChecker_STATUS_MATCH() int8                        { return int8(2) }
func StaticGet_SubtypeChecker_STATUS_NOT_MATCH() int8                    { return int8(3) }

var clinitDone_SubtypeChecker uint32
var clinitState_SubtypeChecker j2gort.ClinitState

func clinit_SubtypeChecker() {
	if atomic.LoadUint32(&clinitDone_SubtypeChecker) != 0 {
		return
	}
	clinitSlow_SubtypeChecker()
}

func clinitSlow_SubtypeChecker() {
	if !j2gort.EnterClinit(&clinitState_SubtypeChecker, &clinitDone_SubtypeChecker) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SubtypeChecker, &clinitDone_SubtypeChecker)
	static_SubtypeChecker_BIVARIANT_TYPES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_3("Object", "IArrayLike", "Array")
}

func (this_ *SubtypeChecker) initFields_SubtypeChecker() {
	_ = this_

}

// Java constructor: SubtypeChecker(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewSubtypeChecker_1(registry *JSTypeRegistry) *SubtypeChecker {
	this_ := &SubtypeChecker{}
	InitSubtypeChecker_1(this_, registry)
	return this_
}

func InitSubtypeChecker_1(this_ *SubtypeChecker, registry *JSTypeRegistry) {
	_ = registry
	this_.initFields_SubtypeChecker()
	this_.registry = registry
}

// Java constructor: SubtypeChecker(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType subtype, com.google.javascript.rhino.jstype.JSType supertype, boolean isUsingStructuralTyping, com.google.javascript.rhino.jstype.JSType.SubtypingMode subtypingMode)
func NewSubtypeChecker_5(registry *JSTypeRegistry, subtype JSType, supertype JSType, isUsingStructuralTyping bool, subtypingMode JSType_SubtypingMode) *SubtypeChecker {
	this_ := &SubtypeChecker{}
	InitSubtypeChecker_5(this_, registry, subtype, supertype, isUsingStructuralTyping, subtypingMode)
	return this_
}

func InitSubtypeChecker_5(this_ *SubtypeChecker, registry *JSTypeRegistry, subtype JSType, supertype JSType, isUsingStructuralTyping bool, subtypingMode JSType_SubtypingMode) {
	_ = registry
	_ = subtype
	_ = supertype
	_ = isUsingStructuralTyping
	_ = subtypingMode
	this_.initFields_SubtypeChecker()
	this_.registry = registry
	this_.initialSubtype = j2gort.CheckNotNullIface(subtype)
	this_.initialSupertype = j2gort.CheckNotNullIface(supertype)
	this_.isUsingStructuralTyping = isUsingStructuralTyping
	this_.hasUsingStructuralTyping = true
	this_.subtypingMode = j2gort.CheckNotNullEnum(subtypingMode)
	this_.hasRun = true
}

// Java method: SubtypeChecker.setSupertype(com.google.javascript.rhino.jstype.JSType value)
func (this_ *SubtypeChecker) SetSupertype(value JSType) *SubtypeChecker {
	_ = value
	this_.priv_SubtypeChecker_checkHasNotRun_9502()
	if this_.initialSupertype != nil && !this_.initialSupertype.IsJavaNil_() {
		j2gort.CheckState(false)
	}
	this_.initialSupertype = j2gort.CheckNotNullIface(value)
	return this_
}

// Java method: SubtypeChecker.setSubtype(com.google.javascript.rhino.jstype.JSType value)
func (this_ *SubtypeChecker) SetSubtype(value JSType) *SubtypeChecker {
	_ = value
	this_.priv_SubtypeChecker_checkHasNotRun_9502()
	if this_.initialSubtype != nil && !this_.initialSubtype.IsJavaNil_() {
		j2gort.CheckState(false)
	}
	this_.initialSubtype = j2gort.CheckNotNullIface(value)
	return this_
}

// Java method: SubtypeChecker.setUsingStructuralSubtyping(boolean value)
func (this_ *SubtypeChecker) SetUsingStructuralSubtyping(value bool) *SubtypeChecker {
	_ = value
	this_.priv_SubtypeChecker_checkHasNotRun_9502()
	if !(!(this_.hasUsingStructuralTyping)) {
		j2gort.CheckState(false)
	}
	this_.hasUsingStructuralTyping = true
	this_.isUsingStructuralTyping = value
	return this_
}

// Java method: SubtypeChecker.setSubtypingMode(com.google.javascript.rhino.jstype.JSType.SubtypingMode value)
func (this_ *SubtypeChecker) SetSubtypingMode(value JSType_SubtypingMode) *SubtypeChecker {
	_ = value
	this_.priv_SubtypeChecker_checkHasNotRun_9502()
	if !(this_.subtypingMode == 0) {
		j2gort.CheckState(false)
	}
	this_.subtypingMode = j2gort.CheckNotNullEnum(value)
	return this_
}

// Java method: SubtypeChecker.checkHasNotRun()
func (this_ *SubtypeChecker) priv_SubtypeChecker_checkHasNotRun_9502() {
	if !(!(this_.hasRun)) {
		j2gort.CheckState(false)
	}
}

// Java static method: SubtypeChecker.check(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType subtype, com.google.javascript.rhino.jstype.JSType supertype, boolean isUsingStructuralTyping, com.google.javascript.rhino.jstype.JSType.SubtypingMode subtypingMode)
func SubtypeChecker_Check(registry *JSTypeRegistry, subtype JSType, supertype JSType, isUsingStructuralTyping bool, subtypingMode JSType_SubtypingMode) bool {
	_ = registry
	_ = subtype
	_ = supertype
	_ = isUsingStructuralTyping
	_ = subtypingMode
	var val_1 SubtypeChecker
	InitSubtypeChecker_5(&val_1, registry, subtype, supertype, isUsingStructuralTyping, subtypingMode)
	var checker *SubtypeChecker = &val_1
	_ = checker
	return checker.priv_SubtypeChecker_isSubtypeCaching_9504(checker.initialSubtype, checker.initialSupertype)
}

// Java method: SubtypeChecker.check()
func (this_ *SubtypeChecker) Check() bool {
	this_.priv_SubtypeChecker_checkHasNotRun_9502()
	this_.hasRun = true
	return this_.priv_SubtypeChecker_isSubtypeCaching_9504(this_.initialSubtype, this_.initialSupertype)
}

// Java method: SubtypeChecker.isSubtypeCaching(com.google.javascript.rhino.jstype.JSType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isSubtypeCaching_9504(subtype JSType, supertype JSType) bool {
	_ = subtype
	_ = supertype
	j2gort.CheckNotNullIface(subtype)
	j2gort.CheckNotNullIface(supertype)
	if (this_.recursionDepth > int32(20)) && (this_.cacheSub == nil) {
		this_.cacheSub = make([]JSType, int32(8))
		this_.cacheSup = make([]JSType, int32(8))
		this_.cacheHash = make([]int32, int32(8))
		this_.cacheStatus = make([]int8, int32(8))
	}
	if this_.cacheSub == nil {
		this_.recursionDepth = this_.recursionDepth + 1
		var result bool = this_.priv_SubtypeChecker_isSubtypeDispatching_9509(subtype, supertype)
		_ = result
		this_.recursionDepth = this_.recursionDepth - 1
		return result
	}
	var h int32 = ((int32(31) * subtype.GetJSTypeBase().HashCode()) + supertype.GetJSTypeBase().HashCode())
	_ = h
	var existingIdx int32 = this_.priv_SubtypeChecker_findCachedEntry_9506(subtype, supertype, h)
	_ = existingIdx
	if existingIdx >= int32(0) {
		var st int8 = this_.cacheStatus[existingIdx]
		_ = st
		if int32(st) == int32(int8(1)) {
			{
				a_2 := this_.cacheStatus
				i_3 := existingIdx
				a_2[i_3] = int8(2)
			}
			return true
		}
		return (int32(st) == int32(int8(2)))
	}
	var idx int32 = this_.priv_SubtypeChecker_insertCachedEntry_9507(subtype, supertype, h)
	_ = idx
	var result_2 bool = this_.priv_SubtypeChecker_isSubtypeDispatching_9509(subtype, supertype)
	_ = result_2
	{
		a_4 := this_.cacheStatus
		i_5 := idx
		a_4[i_5] = func() int8 {
			if result_2 {
				return int8(2)
			}
			return int8(3)
		}()
	}
	return result_2
}

// Java method: SubtypeChecker.matchesCachedKey(int idx, com.google.javascript.rhino.jstype.JSType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_matchesCachedKey_9505(idx int32, subtype JSType, supertype JSType) bool {
	_ = idx
	_ = subtype
	_ = supertype
	var left JSType = this_.cacheSub[idx]
	_ = left
	var right JSType = this_.cacheSup[idx]
	_ = right
	return (((left == subtype || ((left == nil || left.IsJavaNil_()) && (subtype == nil || subtype.IsJavaNil_()))) && (right == supertype || ((right == nil || right.IsJavaNil_()) && (supertype == nil || supertype.IsJavaNil_())))) || (subtype.GetJSTypeBase().EqualsJSType(left) && supertype.GetJSTypeBase().EqualsJSType(right)))
}

// Java method: SubtypeChecker.findCachedEntry(com.google.javascript.rhino.jstype.JSType subtype, com.google.javascript.rhino.jstype.JSType supertype, int h)
func (this_ *SubtypeChecker) priv_SubtypeChecker_findCachedEntry_9506(subtype JSType, supertype JSType, h int32) int32 {
	_ = subtype
	_ = supertype
	_ = h
	var table []int32 = this_.cacheTable
	_ = table
	if table == nil {
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = this_.cacheSize
			_ = n
			for ; i < n; i = i + 1 {
				if (this_.cacheHash[i] == h) && this_.priv_SubtypeChecker_matchesCachedKey_9505(i, subtype, supertype) {
					return i
				}
			}
		}
		return int32(-1)
	}
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
		if (this_.cacheHash[i_2] == h) && this_.priv_SubtypeChecker_matchesCachedKey_9505(i_2, subtype, supertype) {
			return i_2
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: SubtypeChecker.insertCachedEntry(com.google.javascript.rhino.jstype.JSType subtype, com.google.javascript.rhino.jstype.JSType supertype, int h)
func (this_ *SubtypeChecker) priv_SubtypeChecker_insertCachedEntry_9507(subtype JSType, supertype JSType, h int32) int32 {
	_ = subtype
	_ = supertype
	_ = h
	var idx int32 = this_.cacheSize
	_ = idx
	if idx == int32(len(this_.cacheSub)) {
		var newCap int32 = (idx * int32(2))
		_ = newCap
		this_.cacheSub = func() []JSType {
			arrSrc_6 := this_.cacheSub
			arrLen_7 := newCap
			j2gort.CheckNotNullSlice(arrSrc_6)
			if arrLen_7 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_8 := make([]JSType, arrLen_7)
			copy(arrOut_8, arrSrc_6)
			return arrOut_8
		}()
		this_.cacheSup = func() []JSType {
			arrSrc_9 := this_.cacheSup
			arrLen_10 := newCap
			j2gort.CheckNotNullSlice(arrSrc_9)
			if arrLen_10 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_11 := make([]JSType, arrLen_10)
			copy(arrOut_11, arrSrc_9)
			return arrOut_11
		}()
		this_.cacheHash = func() []int32 {
			arrSrc_12 := this_.cacheHash
			arrLen_13 := newCap
			j2gort.CheckNotNullSlice(arrSrc_12)
			if arrLen_13 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_14 := make([]int32, arrLen_13)
			copy(arrOut_14, arrSrc_12)
			return arrOut_14
		}()
		this_.cacheStatus = func() []int8 {
			arrSrc_15 := this_.cacheStatus
			arrLen_16 := newCap
			j2gort.CheckNotNullSlice(arrSrc_15)
			if arrLen_16 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_17 := make([]int8, arrLen_16)
			copy(arrOut_17, arrSrc_15)
			return arrOut_17
		}()
	}
	{
		a_18 := this_.cacheSub
		i_19 := idx
		a_18[i_19] = subtype
	}
	{
		a_20 := this_.cacheSup
		i_21 := idx
		a_20[i_21] = supertype
	}
	{
		a_22 := this_.cacheHash
		i_23 := idx
		a_22[i_23] = h
	}
	{
		a_24 := this_.cacheStatus
		i_25 := idx
		a_24[i_25] = int8(1)
	}
	this_.cacheSize = (idx + int32(1))
	if this_.cacheSize > int32(8) {
		if this_.cacheTable == nil {
			this_.priv_SubtypeChecker_rehashCacheTable_9508(int32(32))
		} else if (this_.cacheSize * int32(4)) >= (int32(len(this_.cacheTable)) * int32(3)) {
			this_.priv_SubtypeChecker_rehashCacheTable_9508((int32(len(this_.cacheTable)) * int32(2)))
		} else {
			var mask int32 = (int32(len(this_.cacheTable)) - int32(1))
			_ = mask
			var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & mask)
			_ = slot
			for this_.cacheTable[slot] != int32(0) {
				slot = ((slot + int32(1)) & mask)
			}
			{
				a_26 := this_.cacheTable
				i_27 := slot
				a_26[i_27] = (idx + int32(1))
			}
		}
	}
	return idx
}

// Java method: SubtypeChecker.rehashCacheTable(int newTableLen)
func (this_ *SubtypeChecker) priv_SubtypeChecker_rehashCacheTable_9508(newTableLen int32) {
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
				a_28 := newTable
				i_29 := slot
				a_28[i_29] = (i + int32(1))
			}
		}
	}
	this_.cacheTable = newTable
}

// Java method: SubtypeChecker.isSubtypeDispatching(com.google.javascript.rhino.jstype.JSType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isSubtypeDispatching_9509(subtype JSType, supertype JSType) bool {
	_ = subtype
	_ = supertype
	switch (subtype.GetJSTypeBase().FastTypeClass()).Ordinal() {
	case int32(5) /* JSTypeClass.ENUM_ELEMENT */ :
		return this_.priv_SubtypeChecker_isEnumElementSubtype_9521(j2gort.CastPtr[EnumElementType](subtype), supertype)
	case int32(4) /* JSTypeClass.ENUM */ :
		return this_.priv_SubtypeChecker_isEnumSubtype_9520(j2gort.CastPtr[EnumType](subtype), supertype)
	case int32(10) /* JSTypeClass.NO_OBJECT */, int32(9) /* JSTypeClass.NO */ :
		return this_.priv_SubtypeChecker_isVariousBottomsSubtype_9518(subtype, supertype)
	case int32(6) /* JSTypeClass.FUNCTION */ :
		return this_.priv_SubtypeChecker_isFunctionSubtype_9517(j2gort.CastPtr[FunctionType](subtype), supertype)
	case int32(20) /* JSTypeClass.TEMPLATE */ :
		return this_.priv_SubtypeChecker_isTemplateSubtype_9523(j2gort.CastPtr[TemplateType](subtype), supertype)
	case int32(15) /* JSTypeClass.PROXY_OBJECT */ :
		return this_.priv_SubtypeChecker_isProxyObjectSubtype_9522(j2gort.AsIface[ProxyObjectType](subtype), supertype)
	case int32(19) /* JSTypeClass.WELL_KNOWN_SYMBOL */ :
		return this_.priv_SubtypeChecker_isWellKnownSymbolSubtype_9510(j2gort.CastPtr[KnownSymbolType](subtype), supertype)
	default:
		return this_.priv_SubtypeChecker_isSubtypeHelper_9511(subtype, supertype)
	}
	panic("unreachable switch expression")
}

// Java method: SubtypeChecker.isWellKnownSymbolSubtype(com.google.javascript.rhino.jstype.KnownSymbolType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isWellKnownSymbolSubtype_9510(subtype *KnownSymbolType, supertype JSType) bool {
	_ = subtype
	_ = supertype
	base_supertype_ := supertype.GetJSTypeBase()
	if base_supertype_.IsKnownSymbolValueType() {
		return subtype.EqualsJSType(supertype)
	}
	if base_supertype_.IsSymbolValueType() {
		return true
	}
	return this_.priv_SubtypeChecker_isSubtypeHelper_9511(ToIface_JSType(subtype), supertype)
}

// Java method: SubtypeChecker.isSubtypeHelper(com.google.javascript.rhino.jstype.JSType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isSubtypeHelper_9511(subtype JSType, supertype JSType) bool {
	_ = subtype
	_ = supertype
	if subtype == supertype || ((subtype == nil || subtype.IsJavaNil_()) && (supertype == nil || supertype.IsJavaNil_())) {
		return true
	}
	base_subtype_ := subtype.GetJSTypeBase()
	var subMask int32 = base_subtype_.typeKindMask
	_ = subMask
	base_supertype_ := supertype.GetJSTypeBase()
	var supMask int32 = base_supertype_.typeKindMask
	_ = supMask
	if ((base_supertype_.IsUnknownType() || (((supMask & (int32(16) | int32(1))) != int32(0)) && base_supertype_.IsAllType())) || base_subtype_.IsUnknownType()) || (((subMask & (int32(2) | int32(1))) != int32(0)) && base_subtype_.IsNoType()) {
		return true
	}
	if (this_.subtypingMode == StaticGet_JSType_SubtypingMode_IGNORE_NULL_UNDEFINED()) && ((subMask & (int32(131072) | int32(262144))) != int32(0)) {
		return true
	}
	if EqualityChecker_CheckNonIdenticalIdentity(subtype, supertype) {
		return true
	}
	if (subMask & int32(1024)) != int32(0) {
		return this_.priv_SubtypeChecker_isSubtypeDispatching_9509(subtype.ToMaybeNamedType().GetReferencedType(), supertype)
	} else if (supMask & int32(1024)) != int32(0) {
		return this_.priv_SubtypeChecker_isSubtypeDispatching_9509(subtype, supertype.ToMaybeNamedType().GetReferencedType())
	}
	if ((subMask & (int32(64) | int32(1))) != int32(0)) && base_subtype_.IsUnionType() {
		var subUnion *UnionType = subtype.ToMaybeUnionType()
		_ = subUnion
		var n int32 = subUnion.GetAlternatesCount()
		_ = n
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				if !(this_.priv_SubtypeChecker_isSubtypeCaching_9504(subUnion.GetAlternateAt(i), supertype)) {
					return false
				}
			}
		}
		return true
	} else if ((supMask & (int32(64) | int32(1))) != int32(0)) && base_supertype_.IsUnionType() {
		var superUnion *UnionType = supertype.ToMaybeUnionType()
		_ = superUnion
		var n_2 int32 = superUnion.GetAlternatesCount()
		_ = n_2
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < n_2; i_2 = i_2 + 1 {
				if this_.priv_SubtypeChecker_isSubtypeCaching_9504(subtype, superUnion.GetAlternateAt(i_2)) {
					return true
				}
			}
		}
		return false
	}
	if !(base_subtype_.IsObjectType()) || !(base_supertype_.IsObjectType()) {
		return false
	}
	return this_.priv_SubtypeChecker_isObjectSubtypeHelper_9512(base_subtype_.AssertObjectType(), base_supertype_.AssertObjectType())
}

// Java method: SubtypeChecker.isObjectSubtypeHelper(com.google.javascript.rhino.jstype.ObjectType subtype, com.google.javascript.rhino.jstype.ObjectType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isObjectSubtypeHelper_9512(subtype ObjectType, supertype ObjectType) bool {
	_ = subtype
	_ = supertype
	var subAndSuperAreSameBaseType bool = false
	_ = subAndSuperAreSameBaseType
	if subtype.GetObjectTypeBase().IsTemplatizedType() && supertype.GetObjectTypeBase().IsTemplatizedType() {
		var templatizedSubType *TemplatizedType = subtype.ToMaybeTemplatizedType()
		_ = templatizedSubType
		var templatizedSuperType *TemplatizedType = supertype.ToMaybeTemplatizedType()
		_ = templatizedSuperType
		var superTypeReferencedType ObjectType = templatizedSuperType.GetReferencedType()
		_ = superTypeReferencedType
		var subTypeReferencedType ObjectType = templatizedSubType.GetReferencedType()
		_ = subTypeReferencedType
		if subTypeReferencedType.GetObjectTypeBase().EqualsJSType(Upcast_JSType(superTypeReferencedType)) {
			subAndSuperAreSameBaseType = true
		}
	}
	var bivarantMatch bool = false
	_ = bivarantMatch
	if SubtypeChecker_IsBivariantType(supertype) {
		var key *TemplateType = subtype.GetJSTypeBase().registry.GetObjectElementKey()
		_ = key
		var thisElement JSType = subtype.GetTemplateTypeMap().GetResolvedTemplateType(key)
		_ = thisElement
		var thatElement JSType = supertype.GetTemplateTypeMap().GetResolvedTemplateType(key)
		_ = thatElement
		if !(this_.priv_SubtypeChecker_meetVarianceConstraint_9525(StaticGet_SubtypeChecker_Variance_BIVARIANT(), thisElement, thatElement)) {
			return false
		}
		bivarantMatch = true
	}
	var covariantKey *TemplateType
	_ = covariantKey
	if bivarantMatch {
		covariantKey = nil
	} else {
		covariantKey = SubtypeChecker_GetTemplateKeyIfCovariantType(supertype)
	}
	if (!(subAndSuperAreSameBaseType) && this_.isUsingStructuralTyping) && supertype.IsStructuralType() {
		if (covariantKey == this_.registry.GetReadonlyArrayElementKey()) && this_.priv_SubtypeChecker_isNativeArrayType_9516(subtype) {
			var thisElement_2 JSType = subtype.GetTemplateTypeMap().GetResolvedTemplateType(covariantKey)
			_ = thisElement_2
			var thatElement_2 JSType = supertype.GetTemplateTypeMap().GetResolvedTemplateType(covariantKey)
			_ = thatElement_2
			return this_.priv_SubtypeChecker_meetVarianceConstraint_9525(StaticGet_SubtypeChecker_Variance_COVARIANT(), thisElement_2, thatElement_2)
		}
		return this_.priv_SubtypeChecker_isStructuralSubtypeHelper_9513(subtype, supertype, StaticGet_SubtypeChecker_PropertyOptionality_VOIDABLE_PROPS_ARE_OPTIONAL())
	}
	if !(bivarantMatch) {
		if covariantKey != nil {
			var thisElement_3 JSType = subtype.GetTemplateTypeMap().GetResolvedTemplateType(covariantKey)
			_ = thisElement_3
			var thatElement_3 JSType = supertype.GetTemplateTypeMap().GetResolvedTemplateType(covariantKey)
			_ = thatElement_3
			if !(this_.priv_SubtypeChecker_meetVarianceConstraint_9525(StaticGet_SubtypeChecker_Variance_COVARIANT(), thisElement_3, thatElement_3)) {
				return false
			}
		} else if !(this_.priv_SubtypeChecker_isTypeMapSubmap_9524(subtype, supertype)) {
			return false
		}
	}
	var subtypeCtor *FunctionType = subtype.GetConstructor()
	_ = subtypeCtor
	var supertypeCtor *FunctionType = supertype.GetConstructor()
	_ = supertypeCtor
	if (subtypeCtor != nil) && subtypeCtor.IsInterface() {
		{
			cur_30 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](subtype.GetCtorExtendedInterfaces())))
			for cur_30.HasNext() {
				rawElem_31 := cur_30.Next()
				subtypeInterface := j2gort.AsIface[ObjectType](rawElem_31)
				_ = subtypeInterface
				if this_.priv_SubtypeChecker_isSubtypeCaching_9504(Upcast_JSType(subtypeInterface), Upcast_JSType(supertype)) {
					return true
				}
			}
		}
	} else if (supertypeCtor != nil) && supertypeCtor.IsInterface() {
		{
			cur_32 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](subtype.GetCtorImplementedInterfaces())))
			for cur_32.HasNext() {
				rawElem_33 := cur_32.Next()
				subtypeInterface_2 := j2gort.AsIface[ObjectType](rawElem_33)
				_ = subtypeInterface_2
				if this_.priv_SubtypeChecker_isSubtypeCaching_9504(Upcast_JSType(subtypeInterface_2), Upcast_JSType(supertype)) {
					return true
				}
			}
		}
	}
	return supertype.GetObjectTypeBase().IsImplicitPrototypeOf(subtype)
}

// Java method: SubtypeChecker.isStructuralSubtypeHelper(com.google.javascript.rhino.jstype.ObjectType subtype, com.google.javascript.rhino.jstype.ObjectType supertype, com.google.javascript.rhino.jstype.SubtypeChecker.PropertyOptionality optionality)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isStructuralSubtypeHelper_9513(subtype ObjectType, supertype ObjectType, optionality SubtypeChecker_PropertyOptionality) bool {
	_ = subtype
	_ = supertype
	_ = optionality
	if supertype.GetObjectTypeBase().IsRecordType() {
		var superProto PrototypeObjectType = ObjectType_UnwrapToPrototypeObjectType(supertype)
		_ = superProto
		var ownMap *PropertyMap
		_ = ownMap
		if superProto != nil && !superProto.IsJavaNil_() {
			ownMap = superProto.GetPrototypeObjectTypeBase().GetOwnPropertyMapOrNull()
		} else {
			ownMap = nil
		}
		if ownMap == nil {
			return true
		}
		if ownMap.HasFlatProps() {
			{
				var i int32 = int32(0)
				_ = i
				var n int32 = ownMap.GetOwnPropertiesCount()
				_ = n
				for ; i < n; i = i + 1 {
					if !(this_.priv_SubtypeChecker_checkPropertyPresence_9514(ownMap.GetFlatPropNameAt(i), subtype, supertype, optionality)) {
						return false
					}
				}
			}
		} else {
			{
				cur_34 := j2gort.NewIterCursor(j2gort.CheckNotNull(ownMap.GetOwnPropertyNames()))
				for cur_34.HasNext() {
					property := cur_34.NextStr()
					_ = property
					if !(this_.priv_SubtypeChecker_checkPropertyPresence_9514(property, subtype, supertype, optionality)) {
						return false
					}
				}
			}
		}
		if ownMap.HasOwnKnownSymbols() {
			{
				cur_35 := j2gort.NewIterCursor(j2gort.CheckNotNull(ownMap.GetOwnKnownSymbols()))
				for cur_35.HasNext() {
					rawElem_36 := cur_35.Next()
					knownSymbolType := j2gort.CastPtr[KnownSymbolType](rawElem_36)
					_ = knownSymbolType
					var knownSymbolTypeKey Property_Key = NewProperty_SymbolKey(knownSymbolType)
					_ = knownSymbolTypeKey
					if !(this_.priv_SubtypeChecker_checkPropertyPresence_9515(knownSymbolTypeKey, subtype, supertype, optionality)) {
						return false
					}
				}
			}
		}
		return true
	}
	var keys *PropertyMap_AllKeys = supertype.GetPropertyMap().GetAllKeys()
	_ = keys
	{
		cur_37 := j2gort.NewIterCursor(j2gort.CheckNotNull(keys.StringKeys()))
		for cur_37.HasNext() {
			property_2 := cur_37.NextStr()
			_ = property_2
			if !(this_.priv_SubtypeChecker_checkPropertyPresence_9514(property_2, subtype, supertype, optionality)) {
				return false
			}
		}
	}
	{
		cur_38 := j2gort.NewIterCursor(j2gort.CheckNotNull(keys.KnownSymbolKeys()))
		for cur_38.HasNext() {
			rawElem_39 := cur_38.Next()
			knownSymbolType_2 := j2gort.CastPtr[KnownSymbolType](rawElem_39)
			_ = knownSymbolType_2
			var knownSymbolTypeKey_2 Property_Key = NewProperty_SymbolKey(knownSymbolType_2)
			_ = knownSymbolTypeKey_2
			if !(this_.priv_SubtypeChecker_checkPropertyPresence_9515(knownSymbolTypeKey_2, subtype, supertype, optionality)) {
				return false
			}
		}
	}
	return true
}

// Java method: SubtypeChecker.checkPropertyPresence(java.lang.String propertyName, com.google.javascript.rhino.jstype.ObjectType subtype, com.google.javascript.rhino.jstype.ObjectType supertype, com.google.javascript.rhino.jstype.SubtypeChecker.PropertyOptionality optionality)
func (this_ *SubtypeChecker) priv_SubtypeChecker_checkPropertyPresence_9514(propertyName string, subtype ObjectType, supertype ObjectType, optionality SubtypeChecker_PropertyOptionality) bool {
	_ = propertyName
	_ = subtype
	_ = supertype
	_ = optionality
	var supertypeProp JSType = supertype.GetObjectTypeBase().GetPropertyType_1_string(propertyName)
	_ = supertypeProp
	base_subtype_ := subtype.GetObjectTypeBase()
	if base_subtype_.HasProperty_1_string(propertyName) {
		var subtypeProp JSType = base_subtype_.GetPropertyType_1_string(propertyName)
		_ = subtypeProp
		if !(this_.priv_SubtypeChecker_isSubtypeCaching_9504(subtypeProp, supertypeProp)) {
			return false
		}
	} else if !(optionality.IsOptional(supertypeProp)) {
		return false
	}
	return true
}

// Java method: SubtypeChecker.checkPropertyPresence(com.google.javascript.rhino.jstype.Property.Key propertyKey, com.google.javascript.rhino.jstype.ObjectType subtype, com.google.javascript.rhino.jstype.ObjectType supertype, com.google.javascript.rhino.jstype.SubtypeChecker.PropertyOptionality optionality)
func (this_ *SubtypeChecker) priv_SubtypeChecker_checkPropertyPresence_9515(propertyKey Property_Key, subtype ObjectType, supertype ObjectType, optionality SubtypeChecker_PropertyOptionality) bool {
	_ = propertyKey
	_ = subtype
	_ = supertype
	_ = optionality
	var supertypeProp JSType = supertype.GetPropertyType_1_Property_Key(propertyKey)
	_ = supertypeProp
	if subtype.GetObjectTypeBase().HasProperty_1_Property_Key(propertyKey) {
		var subtypeProp JSType = subtype.GetPropertyType_1_Property_Key(propertyKey)
		_ = subtypeProp
		if !(this_.priv_SubtypeChecker_isSubtypeCaching_9504(subtypeProp, supertypeProp)) {
			return false
		}
	} else if !(optionality.IsOptional(supertypeProp)) {
		return false
	}
	return true
}

// Java method: SubtypeChecker.isNativeArrayType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isNativeArrayType_9516(type_ ObjectType) bool {
	_ = type_
	var unwrapped ObjectType = SubtypeChecker_GetObjectTypeIfNative(type_)
	_ = unwrapped
	return ((unwrapped != nil && !unwrapped.IsJavaNil_()) && unwrapped.GetObjectTypeBase().EqualsJSType(this_.registry.GetNativeType(StaticGet_JSTypeNative_ARRAY_TYPE())))
}

// Java method: SubtypeChecker.isFunctionSubtype(com.google.javascript.rhino.jstype.FunctionType subtype, com.google.javascript.rhino.jstype.JSType nonFunctionSupertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isFunctionSubtype_9517(subtype *FunctionType, nonFunctionSupertype JSType) bool {
	_ = subtype
	_ = nonFunctionSupertype
	if this_.priv_SubtypeChecker_isSubtypeHelper_9511(ToIface_JSType(subtype), nonFunctionSupertype) {
		return true
	}
	if nonFunctionSupertype.GetJSTypeBase().IsFunctionType() {
		var supertype *FunctionType = nonFunctionSupertype.ToMaybeFunctionType()
		_ = supertype
		if supertype.IsInterface() {
			return true
		}
		if subtype.IsInterface() {
			return false
		}
		return (this_.priv_SubtypeChecker_shouldTreatThisTypesAsCovariant_9526(subtype, supertype) && this_.priv_SubtypeChecker_isArrowTypeSubtype_9519(subtype, supertype))
	}
	return this_.priv_SubtypeChecker_isSubtypeCaching_9504(this_.registry.GetNativeType(StaticGet_JSTypeNative_FUNCTION_PROTOTYPE()), nonFunctionSupertype)
}

// Java method: SubtypeChecker.isVariousBottomsSubtype(com.google.javascript.rhino.jstype.JSType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isVariousBottomsSubtype_9518(subtype JSType, supertype JSType) bool {
	_ = subtype
	_ = supertype
	if this_.priv_SubtypeChecker_isSubtypeHelper_9511(subtype, supertype) {
		return true
	}
	return (supertype.IsObject() && !(supertype.GetJSTypeBase().IsNoType()))
}

// Java method: SubtypeChecker.isArrowTypeSubtype(com.google.javascript.rhino.jstype.FunctionType subtype, com.google.javascript.rhino.jstype.FunctionType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isArrowTypeSubtype_9519(subtype *FunctionType, supertype *FunctionType) bool {
	_ = subtype
	_ = supertype
	if !(this_.priv_SubtypeChecker_isSubtypeCaching_9504(subtype.GetReturnType(), supertype.GetReturnType())) {
		return false
	}
	var subtypeParameters []FunctionType_Parameter = subtype.GetParameterArray()
	_ = subtypeParameters
	var supertypeParameters []FunctionType_Parameter = supertype.GetParameterArray()
	_ = supertypeParameters
	var subSize int32 = int32(len(subtypeParameters))
	_ = subSize
	var superSize int32 = int32(len(supertypeParameters))
	_ = superSize
	var subIdx int32 = int32(0)
	_ = subIdx
	var superIdx int32 = int32(0)
	_ = superIdx
	var subtypeParam FunctionType_Parameter
	_ = subtypeParam
	if subIdx < subSize {
		subtypeParam = subtypeParameters[func() int32 { old_40 := subIdx; nw_41 := old_40 + 1; subIdx = nw_41; return old_40 }()]
	} else {
		subtypeParam = (FunctionType_Parameter{})
	}
	var supertypeParam FunctionType_Parameter
	_ = supertypeParam
	if superIdx < superSize {
		supertypeParam = supertypeParameters[func() int32 { old_42 := superIdx; nw_43 := old_42 + 1; superIdx = nw_43; return old_42 }()]
	} else {
		supertypeParam = (FunctionType_Parameter{})
	}
	for (subtypeParam != (FunctionType_Parameter{})) && (supertypeParam != (FunctionType_Parameter{})) {
		var subtypeParamType JSType = subtypeParam.GetJSType()
		_ = subtypeParamType
		var supertypeParamType JSType = supertypeParam.GetJSType()
		_ = supertypeParamType
		if subtypeParamType != nil && !subtypeParamType.IsJavaNil_() {
			if (supertypeParamType == nil || supertypeParamType.IsJavaNil_()) || !(this_.priv_SubtypeChecker_isSubtypeCaching_9504(supertypeParamType, subtypeParamType)) {
				return false
			}
		}
		var thisIsVarArgs bool = subtypeParam.IsVariadic()
		_ = thisIsVarArgs
		var thatIsVarArgs bool = supertypeParam.IsVariadic()
		_ = thatIsVarArgs
		var thisIsOptional bool = (thisIsVarArgs || subtypeParam.IsOptional())
		_ = thisIsOptional
		var thatIsOptional bool = (thatIsVarArgs || supertypeParam.IsOptional())
		_ = thatIsOptional
		if !(thisIsOptional) && thatIsOptional {
			var isTopFunction bool = (thatIsVarArgs && (((supertypeParamType == nil || supertypeParamType.IsJavaNil_()) || supertypeParamType.GetJSTypeBase().IsUnknownType()) || supertypeParamType.GetJSTypeBase().IsNoType()))
			_ = isTopFunction
			if !(isTopFunction) {
				return false
			}
		}
		if !(thisIsVarArgs) {
			if subIdx < subSize {
				subtypeParam = subtypeParameters[func() int32 { old_44 := subIdx; nw_45 := old_44 + 1; subIdx = nw_45; return old_44 }()]
			} else {
				subtypeParam = (FunctionType_Parameter{})
			}
		}
		if !(thatIsVarArgs) {
			if superIdx < superSize {
				supertypeParam = supertypeParameters[func() int32 { old_46 := superIdx; nw_47 := old_46 + 1; superIdx = nw_47; return old_46 }()]
			} else {
				supertypeParam = (FunctionType_Parameter{})
			}
		}
		if thisIsVarArgs && thatIsVarArgs {
			break
		}
	}
	return ((((subtypeParam == (FunctionType_Parameter{})) || subtypeParam.IsOptional()) || subtypeParam.IsVariadic()) || (supertypeParam != (FunctionType_Parameter{})))
}

// Java method: SubtypeChecker.isEnumSubtype(com.google.javascript.rhino.jstype.EnumType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isEnumSubtype_9520(subtype *EnumType, supertype JSType) bool {
	_ = subtype
	_ = supertype
	return ((supertype.GetJSTypeBase().EqualsJSType(this_.registry.GetNativeType(StaticGet_JSTypeNative_OBJECT_TYPE())) || subtype.EqualsJSType(this_.registry.GetNativeType(StaticGet_JSTypeNative_OBJECT_PROTOTYPE()))) || this_.priv_SubtypeChecker_isSubtypeHelper_9511(ToIface_JSType(subtype), supertype))
}

// Java method: SubtypeChecker.isEnumElementSubtype(com.google.javascript.rhino.jstype.EnumElementType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isEnumElementSubtype_9521(subtype *EnumElementType, supertype JSType) bool {
	_ = subtype
	_ = supertype
	if supertype.GetJSTypeBase().IsEnumElementType() && (subtype.GetEnumType() == supertype.ToMaybeEnumElementType().GetEnumType()) {
		return this_.priv_SubtypeChecker_isSubtypeCaching_9504(subtype.GetPrimitiveType(), supertype.ToMaybeEnumElementType().GetPrimitiveType())
	}
	if this_.priv_SubtypeChecker_isSubtypeHelper_9511(ToIface_JSType(subtype), supertype) {
		return true
	} else {
		return this_.priv_SubtypeChecker_isSubtypeCaching_9504(subtype.GetPrimitiveType(), supertype)
	}
}

// Java method: SubtypeChecker.isProxyObjectSubtype(com.google.javascript.rhino.jstype.ProxyObjectType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isProxyObjectSubtype_9522(subtype ProxyObjectType, supertype JSType) bool {
	_ = subtype
	_ = supertype
	return this_.priv_SubtypeChecker_isSubtypeDispatching_9509(subtype.GetProxyObjectTypeBase().GetReferencedTypeInternal(), supertype)
}

// Java method: SubtypeChecker.isTemplateSubtype(com.google.javascript.rhino.jstype.TemplateType subtype, com.google.javascript.rhino.jstype.JSType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isTemplateSubtype_9523(subtype *TemplateType, supertype JSType) bool {
	_ = subtype
	_ = supertype
	if (!(subtype.GetBound().GetJSTypeBase().IsUnknownType()) && supertype.GetJSTypeBase().IsTemplateType()) && !(supertype.ToMaybeTemplateType().GetBound().GetJSTypeBase().IsUnknownType()) {
		return (j2gort.Cast[ContainsUpperBoundSuperTypeVisitor_Result](subtype.Visit_1(NewContainsUpperBoundSuperTypeVisitor(supertype))) == StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_PRESENT())
	} else {
		return this_.priv_SubtypeChecker_isProxyObjectSubtype_9522(ToIface_ProxyObjectType(subtype), supertype)
	}
}

// Java method: SubtypeChecker.isTypeMapSubmap(com.google.javascript.rhino.jstype.ObjectType subtype, com.google.javascript.rhino.jstype.ObjectType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_isTypeMapSubmap_9524(subtype ObjectType, supertype ObjectType) bool {
	_ = subtype
	_ = supertype
	if ((supertype.GetJSTypeBase().typeKindMask & int32(67108865)) == int32(0)) || subtype.GetObjectTypeBase().IsFunctionPrototypeType() {
		return true
	}
	var submap *TemplateTypeMap = subtype.GetTemplateTypeMap()
	_ = submap
	var supermap *TemplateTypeMap = supertype.GetTemplateTypeMap()
	_ = supermap
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = supermap.Size()
		_ = n
		for ; i < n; i = i + 1 {
			var key *TemplateType = supermap.GetTemplateKeyAt(i)
			_ = key
			var keyCount int32 = submap.GetTemplateKeyCountThisShouldAlwaysBeOneOrZeroButIsnt(key)
			_ = keyCount
			if keyCount == int32(0) {
				if subtype.LoosenTypecheckingDueToForwardReferencedSupertype() {
					continue
				}
				return false
			} else if keyCount > int32(1) {
				continue
			}
			var subvalue JSType = submap.GetResolvedTemplateType(key)
			_ = subvalue
			var supervalue JSType = supermap.GetResolvedTemplateType(key)
			_ = supervalue
			if !(this_.priv_SubtypeChecker_meetVarianceConstraint_9525(StaticGet_SubtypeChecker_Variance_INVARIANT(), subvalue, supervalue)) {
				return false
			}
		}
	}
	return true
}

// Java method: SubtypeChecker.meetVarianceConstraint(com.google.javascript.rhino.jstype.SubtypeChecker.Variance variance, com.google.javascript.rhino.jstype.JSType override, com.google.javascript.rhino.jstype.JSType reference)
func (this_ *SubtypeChecker) priv_SubtypeChecker_meetVarianceConstraint_9525(variance SubtypeChecker_Variance, override JSType, reference JSType) bool {
	_ = variance
	_ = override
	_ = reference
	switch (variance).Ordinal() {
	case int32(0) /* Variance.COVARIANT */ :
		return this_.priv_SubtypeChecker_isSubtypeCaching_9504(override, reference)
	case int32(1) /* Variance.CONTRAVARIANT */ :
		return this_.priv_SubtypeChecker_isSubtypeCaching_9504(reference, override)
	case int32(2) /* Variance.BIVARIANT */ :
		return (this_.priv_SubtypeChecker_meetVarianceConstraint_9525(StaticGet_SubtypeChecker_Variance_COVARIANT(), reference, override) || this_.priv_SubtypeChecker_meetVarianceConstraint_9525(StaticGet_SubtypeChecker_Variance_CONTRAVARIANT(), reference, override))
	case int32(3) /* Variance.INVARIANT */ :
		return (this_.priv_SubtypeChecker_meetVarianceConstraint_9525(StaticGet_SubtypeChecker_Variance_COVARIANT(), reference, override) && this_.priv_SubtypeChecker_meetVarianceConstraint_9525(StaticGet_SubtypeChecker_Variance_CONTRAVARIANT(), reference, override))
	}
	panic("unreachable switch expression")
}

// Java static method: SubtypeChecker.isBivariantType(com.google.javascript.rhino.jstype.ObjectType type)
func SubtypeChecker_IsBivariantType(type_ ObjectType) bool {
	_ = type_
	var unwrapped ObjectType = SubtypeChecker_GetObjectTypeIfNative(type_)
	_ = unwrapped
	return ((unwrapped != nil && !unwrapped.IsJavaNil_()) && j2gort.SetContainsStr(StaticGet_SubtypeChecker_BIVARIANT_TYPES(), unwrapped.GetReferenceName()))
}

// Java static method: SubtypeChecker.getTemplateKeyIfCovariantType(com.google.javascript.rhino.jstype.ObjectType type)
func SubtypeChecker_GetTemplateKeyIfCovariantType(type_ ObjectType) *TemplateType {
	_ = type_
	var mask int32 = type_.GetJSTypeBase().typeKindMask
	_ = mask
	if (mask & ((int32(4096) | int32(33554432)) | int32(1))) == int32(0) {
		return nil
	}
	if ((mask & (int32(4096) | int32(1))) != int32(0)) && type_.GetObjectTypeBase().IsTemplatizedType() {
		var ttype *TemplatizedType = type_.ToMaybeTemplatizedType()
		_ = ttype
		if ttype.GetTemplateTypeMap().HasTemplateKey_1(ttype.GetJSTypeBase().registry.GetIThenableTemplate()) {
			return ttype.GetJSTypeBase().registry.GetIThenableTemplate()
		}
	}
	var unwrapped ObjectType = SubtypeChecker_GetObjectTypeIfNative(type_)
	_ = unwrapped
	var unwrappedTypeName string
	_ = unwrappedTypeName
	if unwrapped == nil || unwrapped.IsJavaNil_() {
		unwrappedTypeName = j2gort.NullString
	} else {
		unwrappedTypeName = unwrapped.GetReferenceName()
	}
	if unwrappedTypeName == j2gort.NullString {
		return nil
	}
	switch j2gort.CheckNotNullString(unwrappedTypeName) {
	case "ReadonlyArray":
		return unwrapped.GetJSTypeBase().registry.GetReadonlyArrayElementKey()
	case "IteratorLike":
		return unwrapped.GetJSTypeBase().registry.GetIteratorLikeValueTemplate()
	case "Iterator":
		return unwrapped.GetJSTypeBase().registry.GetIteratorValueTemplate()
	case "Generator":
		return unwrapped.GetJSTypeBase().registry.GetGeneratorValueTemplate()
	case "AsyncIterator":
		return unwrapped.GetJSTypeBase().registry.GetAsyncIteratorValueTemplate()
	case "Iterable":
		return unwrapped.GetJSTypeBase().registry.GetIterableValueTemplate()
	case "IteratorIterable":
		return unwrapped.GetJSTypeBase().registry.GetIteratorIterableValueTemplate()
	case "IIterableResult":
		return unwrapped.GetJSTypeBase().registry.GetIIterableResultValueTemplate()
	case "AsyncIterable":
		return unwrapped.GetJSTypeBase().registry.GetAsyncIterableValueTemplate()
	default:
		return nil
	}
	panic("unreachable switch expression")
}

// Java method: SubtypeChecker.shouldTreatThisTypesAsCovariant(com.google.javascript.rhino.jstype.FunctionType subtype, com.google.javascript.rhino.jstype.FunctionType supertype)
func (this_ *SubtypeChecker) priv_SubtypeChecker_shouldTreatThisTypesAsCovariant_9526(subtype *FunctionType, supertype *FunctionType) bool {
	_ = subtype
	_ = supertype
	if (!(j2gort.IfaceEquals[ObjectType](supertype.GetTypeOfThis().GetJSTypeBase().ToObjectType(), nil)) && (supertype.GetTypeOfThis().GetJSTypeBase().ToObjectType().GetConstructor() != nil)) && supertype.GetTypeOfThis().GetJSTypeBase().ToObjectType().GetConstructor().IsInterface() {
		return true
	}
	var originalMode JSType_SubtypingMode = this_.subtypingMode
	_ = originalMode
	this_.subtypingMode = StaticGet_JSType_SubtypingMode_NORMAL()
	var result bool = (this_.priv_SubtypeChecker_isSubtypeCaching_9504(supertype.GetTypeOfThis(), subtype.GetTypeOfThis()) || this_.priv_SubtypeChecker_isSubtypeCaching_9504(subtype.GetTypeOfThis(), supertype.GetTypeOfThis()))
	_ = result
	this_.subtypingMode = originalMode
	return result
}

// Java static method: SubtypeChecker.getObjectTypeIfNative(com.google.javascript.rhino.jstype.ObjectType type)
func SubtypeChecker_GetObjectTypeIfNative(type_ ObjectType) ObjectType {
	_ = type_
	var mask int32 = type_.GetJSTypeBase().typeKindMask
	_ = mask
	if (mask & (int32(33554432) | int32(1))) == int32(0) {
		return nil
	}
	var unwrapped ObjectType
	_ = unwrapped
	if (mask & int32(1)) == int32(0) {
		unwrapped = type_
	} else {
		unwrapped = ObjectType_DeeplyUnwrap(type_)
	}
	if (unwrapped != nil && !unwrapped.IsJavaNil_()) && unwrapped.GetObjectTypeBase().IsNativeObjectType() {
		return unwrapped
	} else {
		return nil
	}
}

// Java enum: com.google.javascript.rhino.jstype.SubtypeChecker.PropertyOptionality
type SubtypeChecker_PropertyOptionality uint8

const (
	enum_SubtypeChecker_PropertyOptionality_VOIDABLE_PROPS_ARE_OPTIONAL SubtypeChecker_PropertyOptionality = 1
	enum_SubtypeChecker_PropertyOptionality_ALL_PROPS_ARE_REQUIRED      SubtypeChecker_PropertyOptionality = 2
)

var enumBaseTable_SubtypeChecker_PropertyOptionality = [...]j2gort.EnumBase{
	{Name_: "VOIDABLE_PROPS_ARE_OPTIONAL", Ordinal_: 0},
	{Name_: "ALL_PROPS_ARE_REQUIRED", Ordinal_: 1},
}
var enumValues_SubtypeChecker_PropertyOptionality = []SubtypeChecker_PropertyOptionality{1, 2}

func StaticGet_SubtypeChecker_PropertyOptionality_VOIDABLE_PROPS_ARE_OPTIONAL() SubtypeChecker_PropertyOptionality {
	return enum_SubtypeChecker_PropertyOptionality_VOIDABLE_PROPS_ARE_OPTIONAL
}
func StaticGet_SubtypeChecker_PropertyOptionality_ALL_PROPS_ARE_REQUIRED() SubtypeChecker_PropertyOptionality {
	return enum_SubtypeChecker_PropertyOptionality_ALL_PROPS_ARE_REQUIRED
}

func SubtypeChecker_PropertyOptionality_Values() []SubtypeChecker_PropertyOptionality {
	return j2gort.CloneSlice(enumValues_SubtypeChecker_PropertyOptionality)
}

func SubtypeChecker_PropertyOptionality_ValueOf(name string) SubtypeChecker_PropertyOptionality {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_SubtypeChecker_PropertyOptionality {
		if enumBaseTable_SubtypeChecker_PropertyOptionality[i_].Name_ == name {
			return SubtypeChecker_PropertyOptionality(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.SubtypeChecker.PropertyOptionality." + name))
}

// Java method: PropertyOptionality.isOptional(com.google.javascript.rhino.jstype.JSType propType)
func (this_ SubtypeChecker_PropertyOptionality) IsOptional(propType JSType) bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = propType
	return ((this_ == StaticGet_SubtypeChecker_PropertyOptionality_VOIDABLE_PROPS_ARE_OPTIONAL()) && propType.IsExplicitlyVoidable())
}

func (this_ SubtypeChecker_PropertyOptionality) IsJavaNil_() bool                              { return this_ == 0 }
func (SubtypeChecker_PropertyOptionality) HasIdentityEquals_()                                 {}
func (SubtypeChecker_PropertyOptionality) Implements_j2gort_Ext_java_lang_Enum()               {}
func (SubtypeChecker_PropertyOptionality) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ SubtypeChecker_PropertyOptionality) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (SubtypeChecker_PropertyOptionality) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.SubtypeChecker$PropertyOptionality", "PropertyOptionality"
}
func (this_ SubtypeChecker_PropertyOptionality) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_SubtypeChecker_PropertyOptionality[this_-1]
}
func (this_ SubtypeChecker_PropertyOptionality) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ SubtypeChecker_PropertyOptionality) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ SubtypeChecker_PropertyOptionality) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_SubtypeChecker_PropertyOptionality[this_-1].Name_
}
func (this_ SubtypeChecker_PropertyOptionality) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_SubtypeChecker_PropertyOptionality[this_-1].Name_
}
func (this_ SubtypeChecker_PropertyOptionality) Equals(other_ any) bool {
	o_, ok_ := other_.(SubtypeChecker_PropertyOptionality)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ SubtypeChecker_PropertyOptionality) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_SubtypeChecker_PropertyOptionality[this_-1].HashCode()
}
func (this_ SubtypeChecker_PropertyOptionality) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_SubtypeChecker_PropertyOptionality[this_-1].CompareTo(other_)
}
func (this_ SubtypeChecker_PropertyOptionality) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_SubtypeChecker_PropertyOptionality[this_-1].Name_)
}
func init() {
	vals_ := []any{SubtypeChecker_PropertyOptionality(1), SubtypeChecker_PropertyOptionality(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.SubtypeChecker.PropertyOptionality", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.SubtypeChecker$PropertyOptionality", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.SubtypeChecker.PropertyOptionality", "com.google.javascript.rhino.jstype.SubtypeChecker$PropertyOptionality", func(v_ any) bool {
		return j2gort.InstanceOf[SubtypeChecker_PropertyOptionality](v_)
	})
}

// Java enum: com.google.javascript.rhino.jstype.SubtypeChecker.Variance
type SubtypeChecker_Variance uint8

const (
	enum_SubtypeChecker_Variance_COVARIANT     SubtypeChecker_Variance = 1
	enum_SubtypeChecker_Variance_CONTRAVARIANT SubtypeChecker_Variance = 2
	enum_SubtypeChecker_Variance_BIVARIANT     SubtypeChecker_Variance = 3
	enum_SubtypeChecker_Variance_INVARIANT     SubtypeChecker_Variance = 4
)

var enumBaseTable_SubtypeChecker_Variance = [...]j2gort.EnumBase{
	{Name_: "COVARIANT", Ordinal_: 0},
	{Name_: "CONTRAVARIANT", Ordinal_: 1},
	{Name_: "BIVARIANT", Ordinal_: 2},
	{Name_: "INVARIANT", Ordinal_: 3},
}
var enumValues_SubtypeChecker_Variance = []SubtypeChecker_Variance{1, 2, 3, 4}

func StaticGet_SubtypeChecker_Variance_COVARIANT() SubtypeChecker_Variance {
	return enum_SubtypeChecker_Variance_COVARIANT
}
func StaticGet_SubtypeChecker_Variance_CONTRAVARIANT() SubtypeChecker_Variance {
	return enum_SubtypeChecker_Variance_CONTRAVARIANT
}
func StaticGet_SubtypeChecker_Variance_BIVARIANT() SubtypeChecker_Variance {
	return enum_SubtypeChecker_Variance_BIVARIANT
}
func StaticGet_SubtypeChecker_Variance_INVARIANT() SubtypeChecker_Variance {
	return enum_SubtypeChecker_Variance_INVARIANT
}

func SubtypeChecker_Variance_Values() []SubtypeChecker_Variance {
	return j2gort.CloneSlice(enumValues_SubtypeChecker_Variance)
}

func SubtypeChecker_Variance_ValueOf(name string) SubtypeChecker_Variance {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_SubtypeChecker_Variance {
		if enumBaseTable_SubtypeChecker_Variance[i_].Name_ == name {
			return SubtypeChecker_Variance(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.SubtypeChecker.Variance." + name))
}

func (this_ SubtypeChecker_Variance) IsJavaNil_() bool                              { return this_ == 0 }
func (SubtypeChecker_Variance) HasIdentityEquals_()                                 {}
func (SubtypeChecker_Variance) Implements_j2gort_Ext_java_lang_Enum()               {}
func (SubtypeChecker_Variance) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ SubtypeChecker_Variance) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (SubtypeChecker_Variance) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.SubtypeChecker$Variance", "Variance"
}
func (this_ SubtypeChecker_Variance) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_SubtypeChecker_Variance[this_-1]
}
func (this_ SubtypeChecker_Variance) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ SubtypeChecker_Variance) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ SubtypeChecker_Variance) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_SubtypeChecker_Variance[this_-1].Name_
}
func (this_ SubtypeChecker_Variance) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_SubtypeChecker_Variance[this_-1].Name_
}
func (this_ SubtypeChecker_Variance) Equals(other_ any) bool {
	o_, ok_ := other_.(SubtypeChecker_Variance)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ SubtypeChecker_Variance) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_SubtypeChecker_Variance[this_-1].HashCode()
}
func (this_ SubtypeChecker_Variance) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_SubtypeChecker_Variance[this_-1].CompareTo(other_)
}
func (this_ SubtypeChecker_Variance) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_SubtypeChecker_Variance[this_-1].Name_)
}
func init() {
	vals_ := []any{SubtypeChecker_Variance(1), SubtypeChecker_Variance(2), SubtypeChecker_Variance(3), SubtypeChecker_Variance(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.SubtypeChecker.Variance", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.SubtypeChecker$Variance", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.SubtypeChecker.Variance", "com.google.javascript.rhino.jstype.SubtypeChecker$Variance", func(v_ any) bool {
		return j2gort.InstanceOf[SubtypeChecker_Variance](v_)
	})
}
