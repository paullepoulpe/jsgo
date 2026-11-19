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
// Port of com/google/javascript/rhino/jstype/TemplateTypeMap.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.TemplateTypeMap
type TemplateTypeMap struct {
	val0       JSType                 `gson:"val0"`
	val1       JSType                 `gson:"val1"`
	base       *TemplateTypeMap       `gson:"base"`
	key0       *TemplateType          `gson:"key0"`
	extra      *TemplateTypeMap_Extra `gson:"extra"`
	idHash_    atomic.Int32
	ownValKind int8 `gson:"ownValKind"`
	flags      int8 `gson:"flags"`
}

func (this_ *TemplateTypeMap) BindSelf_(s any) {
	_ = s
}
func (this_ *TemplateTypeMap) IsJavaNil_() bool { return this_ == nil }
func (*TemplateTypeMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.TemplateTypeMap", "TemplateTypeMap"
}
func (this_ *TemplateTypeMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TemplateTypeMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.TemplateTypeMap", "com.google.javascript.rhino.jstype.TemplateTypeMap", func(v_ any) bool {
		return j2gort.InstanceOf[*TemplateTypeMap](v_)
	})
}

var static_TemplateTypeMap_EMPTY_KEYS []*TemplateType

func StaticGet_TemplateTypeMap_EMPTY_KEYS() []*TemplateType {
	clinit_TemplateTypeMap()
	return static_TemplateTypeMap_EMPTY_KEYS
}
func StaticSet_TemplateTypeMap_EMPTY_KEYS(v []*TemplateType) {
	clinit_TemplateTypeMap()
	static_TemplateTypeMap_EMPTY_KEYS = v
}

var static_TemplateTypeMap_EMPTY_VALUES []JSType

func StaticGet_TemplateTypeMap_EMPTY_VALUES() []JSType {
	clinit_TemplateTypeMap()
	return static_TemplateTypeMap_EMPTY_VALUES
}
func StaticSet_TemplateTypeMap_EMPTY_VALUES(v []JSType) {
	clinit_TemplateTypeMap()
	static_TemplateTypeMap_EMPTY_VALUES = v
}

var static_TemplateTypeMap_EMPTY_KEY_LIST j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_TemplateTypeMap_EMPTY_KEY_LIST() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_TemplateTypeMap()
	return static_TemplateTypeMap_EMPTY_KEY_LIST
}
func StaticSet_TemplateTypeMap_EMPTY_KEY_LIST(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_TemplateTypeMap()
	static_TemplateTypeMap_EMPTY_KEY_LIST = v
}

var static_TemplateTypeMap_EMPTY_VALUE_LIST j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_TemplateTypeMap_EMPTY_VALUE_LIST() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_TemplateTypeMap()
	return static_TemplateTypeMap_EMPTY_VALUE_LIST
}
func StaticSet_TemplateTypeMap_EMPTY_VALUE_LIST(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_TemplateTypeMap()
	static_TemplateTypeMap_EMPTY_VALUE_LIST = v
}
func StaticGet_TemplateTypeMap_FLAG_SUBMAP_BIT0() int8         { return int8(1) }
func StaticGet_TemplateTypeMap_RESOLVED_SHIFT() int32          { return int32(1) }
func StaticGet_TemplateTypeMap_RESOLVED_UNCOMPUTED() int32     { return int32(0) }
func StaticGet_TemplateTypeMap_RESOLVED_SAME_AS_VALUES() int32 { return int32(1) }
func StaticGet_TemplateTypeMap_RESOLVED_IN_EXTRA() int32       { return int32(2) }

var clinitDone_TemplateTypeMap uint32
var clinitState_TemplateTypeMap j2gort.ClinitState

func clinit_TemplateTypeMap() {
	if atomic.LoadUint32(&clinitDone_TemplateTypeMap) != 0 {
		return
	}
	clinitSlow_TemplateTypeMap()
}

func clinitSlow_TemplateTypeMap() {
	if !j2gort.EnterClinit(&clinitState_TemplateTypeMap, &clinitDone_TemplateTypeMap) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TemplateTypeMap, &clinitDone_TemplateTypeMap)
	static_TemplateTypeMap_EMPTY_KEYS = make([]*TemplateType, int32(0))
	static_TemplateTypeMap_EMPTY_VALUES = make([]JSType, int32(0))
	static_TemplateTypeMap_EMPTY_KEY_LIST = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	static_TemplateTypeMap_EMPTY_VALUE_LIST = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
}

func (this_ *TemplateTypeMap) initFields_TemplateTypeMap() {
	_ = this_
}

// Java constructor: TemplateTypeMap(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewTemplateTypeMap_1(registry *JSTypeRegistry) *TemplateTypeMap {
	clinit_TemplateTypeMap()
	this_ := &TemplateTypeMap{}
	InitTemplateTypeMap_1(this_, registry)
	return this_
}

func InitTemplateTypeMap_1(this_ *TemplateTypeMap, registry *JSTypeRegistry) {
	_ = registry
	this_.initFields_TemplateTypeMap()
	this_.val0 = nil
	this_.val1 = nil
	this_.base = nil
	this_.key0 = nil
	this_.extra = NewTemplateTypeMap_Extra(registry, nil, nil, int32(0), nil)
	this_.ownValKind = int8(int32(0))
	this_.flags = int8(2)
}

// Java constructor: TemplateTypeMap(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable TemplateTypeMap base, com.google.javascript.rhino.jstype.TemplateType key0, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType val0, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType val1, byte ownValKind, byte flags, com.google.javascript.rhino.jstype.TemplateTypeMap.@org.jspecify.annotations.Nullable Extra extra)
func NewTemplateTypeMap_7(base *TemplateTypeMap, key0 *TemplateType, val0 JSType, val1 JSType, ownValKind int8, flags int8, extra *TemplateTypeMap_Extra) *TemplateTypeMap {
	clinit_TemplateTypeMap()
	this_ := &TemplateTypeMap{}
	InitTemplateTypeMap_7(this_, base, key0, val0, val1, ownValKind, flags, extra)
	return this_
}

func InitTemplateTypeMap_7(this_ *TemplateTypeMap, base *TemplateTypeMap, key0 *TemplateType, val0 JSType, val1 JSType, ownValKind int8, flags int8, extra *TemplateTypeMap_Extra) {
	_ = base
	_ = key0
	_ = val0
	_ = val1
	_ = ownValKind
	_ = flags
	_ = extra
	this_.initFields_TemplateTypeMap()
	this_.val0 = val0
	this_.val1 = val1
	this_.base = base
	this_.key0 = key0
	this_.extra = extra
	this_.ownValKind = ownValKind
	this_.flags = flags
}

// Java static method: TemplateTypeMap.createEmpty(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func TemplateTypeMap_CreateEmpty(registry *JSTypeRegistry) *TemplateTypeMap {
	clinit_TemplateTypeMap()
	_ = registry
	if !(registry.GetEmptyTemplateTypeMap() == nil) {
		j2gort.CheckArgument(false)
	}
	return NewTemplateTypeMap_1(registry)
}

// Java static method: TemplateTypeMap.createSingleKey(com.google.javascript.rhino.jstype.TemplateType key, int subMapStartsMask)
func TemplateTypeMap_CreateSingleKey(key *TemplateType, subMapStartsMask int32) *TemplateTypeMap {
	clinit_TemplateTypeMap()
	_ = key
	_ = subMapStartsMask
	var f int8 = int8((int32((func() int8 {
		if (subMapStartsMask & int32(1)) != int32(0) {
			return int8(1)
		}
		return int8(int32(0))
	}())) | (javalang.Shl32(int32(1), int(int32(1))))))
	_ = f
	return NewTemplateTypeMap_7(nil, key, nil, nil, int8(0), f, nil)
}

// Java static method: TemplateTypeMap.ofSingleBinding(com.google.javascript.rhino.jstype.TemplateType key, com.google.javascript.rhino.jstype.JSType val)
func TemplateTypeMap_OfSingleBinding(key *TemplateType, val JSType) *TemplateTypeMap {
	clinit_TemplateTypeMap()
	_ = key
	_ = val
	return TemplateTypeMap_CreateSingleKeyAndValue(j2gort.CheckNotNullPtr(key), j2gort.CheckNotNullIface(val), int32(1))
}

// Java static method: TemplateTypeMap.fromInferenceMap(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.util.Map<com.google.javascript.rhino.jstype.TemplateType,com.google.javascript.rhino.jstype.JSType> bindings)
func TemplateTypeMap_FromInferenceMap(registry *JSTypeRegistry, bindings j2gort.Ext_java_util_Map) *TemplateTypeMap {
	clinit_TemplateTypeMap()
	_ = registry
	_ = bindings
	var size int32 = bindings.Size()
	_ = size
	if size == int32(0) {
		return registry.GetEmptyTemplateTypeMap()
	}
	var unknown JSType = registry.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = unknown
	if size == int32(1) {
		var e j2gort.Ext_java_util_Map_Entry = j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.AsIface[j2gort.Ext_java_util_Set](bindings.EntrySet()).Iterator()).Next())
		_ = e
		var val JSType = j2gort.AsIface[JSType](e.GetValue())
		_ = val
		return TemplateTypeMap_CreateSingleKeyAndValue(j2gort.CheckNotNullPtr(j2gort.CastPtr[TemplateType](e.GetKey())), func() JSType {
			if val != nil && !val.IsJavaNil_() {
				return val
			}
			return unknown
		}(), int32(1))
	}
	var keys []*TemplateType = make([]*TemplateType, size)
	_ = keys
	var vals []JSType = make([]JSType, size)
	_ = vals
	var i int32 = int32(0)
	_ = i
	{
		cur_1 := j2gort.NewMapEntryCursor(bindings)
		for cur_1.NextEntry() {
			{
				a_2 := keys
				i_3 := i
				a_2[i_3] = j2gort.CheckNotNullPtr(j2gort.CastPtr[TemplateType](cur_1.GetKey()))
			}
			var val_2 JSType = j2gort.AsIface[JSType](cur_1.GetValue())
			_ = val_2
			{
				a_4 := vals
				i_5 := i
				a_4[i_5] = func() JSType {
					if val_2 != nil && !val_2.IsJavaNil_() {
						return val_2
					}
					return unknown
				}()
			}
			i = i + 1
		}
	}
	return TemplateTypeMap_CreateFromArrays(registry, keys, vals, int32(1), nil)
}

// Java static method: TemplateTypeMap.createSingleKeyAndValue(com.google.javascript.rhino.jstype.TemplateType key, com.google.javascript.rhino.jstype.JSType val, int subMapStartsMask)
func TemplateTypeMap_CreateSingleKeyAndValue(key *TemplateType, val JSType, subMapStartsMask int32) *TemplateTypeMap {
	clinit_TemplateTypeMap()
	_ = key
	_ = val
	_ = subMapStartsMask
	var f int8
	_ = f
	if (subMapStartsMask & int32(1)) != int32(0) {
		f = int8(1)
	} else {
		f = int8(int32(0))
	}
	return NewTemplateTypeMap_7(nil, key, val, nil, int8(1), f, nil)
}

// Java static method: TemplateTypeMap.createFromArrays(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.TemplateType[] templateKeys, com.google.javascript.rhino.jstype.JSType[] templateValues, int subMapStartsMask, java.util.@org.jspecify.annotations.Nullable BitSet subMapStartsOverflow)
func TemplateTypeMap_CreateFromArrays(registry *JSTypeRegistry, templateKeys []*TemplateType, templateValues []JSType, subMapStartsMask int32, subMapStartsOverflow *j2gort.Ext_java_util_BitSet) *TemplateTypeMap {
	clinit_TemplateTypeMap()
	_ = registry
	_ = templateKeys
	_ = templateValues
	_ = subMapStartsMask
	_ = subMapStartsOverflow
	var nKeys int32 = int32(len(templateKeys))
	_ = nKeys
	var nVals int32 = int32(len(templateValues))
	_ = nVals
	if !(nVals <= nKeys) {
		j2gort.CheckArgument(false)
	}
	if nKeys == int32(0) {
		var empty *TemplateTypeMap = registry.GetEmptyTemplateTypeMap()
		_ = empty
		if empty != nil {
			return empty
		} else {
			return NewTemplateTypeMap_1(registry)
		}
	}
	if (nKeys == int32(1)) && (subMapStartsOverflow == nil) {
		if nVals == int32(0) {
			return TemplateTypeMap_CreateSingleKey(templateKeys[int32(0)], subMapStartsMask)
		}
		return TemplateTypeMap_CreateSingleKeyAndValue(templateKeys[int32(0)], templateValues[int32(0)], subMapStartsMask)
	}
	var v0 JSType = nil
	_ = v0
	var v1 JSType = nil
	_ = v1
	var extraVals []JSType = nil
	_ = extraVals
	var valKind int8 = int8(0)
	_ = valKind
	var f int8 = int8(int32(0))
	_ = f
	if nVals == int32(0) {
		valKind = int8(int32(0))
		f = int8(2)
	} else if nVals == int32(1) {
		valKind = int8(int32(1))
		v0 = templateValues[int32(0)]
	} else if nVals == int32(2) {
		valKind = int8(int32(2))
		v0 = templateValues[int32(0)]
		v1 = templateValues[int32(1)]
	} else {
		valKind = int8(int32(3))
		extraVals = templateValues
	}
	var val_6 TemplateTypeMap_Extra
	InitTemplateTypeMap_Extra(&val_6, nil, templateKeys, extraVals, subMapStartsMask, subMapStartsOverflow)
	var extra *TemplateTypeMap_Extra = &val_6
	_ = extra
	return NewTemplateTypeMap_7(nil, templateKeys[int32(0)], v0, v1, valKind, f, extra)
}

// Java method: TemplateTypeMap.registry()
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_registry_9516() *JSTypeRegistry {
	if this_.key0 != nil {
		return this_.key0.GetJSTypeBase().registry
	} else {
		return this_.extra.emptyMapRegistry
	}
}

// Java static method: TemplateTypeMap.copyKeysListToArray(com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.TemplateType> keys)
func TemplateTypeMap_CopyKeysListToArray(keys j2gort.Ext_com_google_common_collect_ImmutableList) []*TemplateType {
	clinit_TemplateTypeMap()
	_ = keys
	var n int32 = keys.Size()
	_ = n
	if n == int32(0) {
		return StaticGet_TemplateTypeMap_EMPTY_KEYS()
	}
	var out []*TemplateType = make([]*TemplateType, n)
	_ = out
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			{
				a_7 := out
				i_8 := i
				a_7[i_8] = j2gort.CheckNotNullPtr(j2gort.CastPtr[TemplateType](keys.Get(i)))
			}
		}
	}
	return out
}

// Java static method: TemplateTypeMap.copyValuesListToArray(com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.JSType> values)
func TemplateTypeMap_CopyValuesListToArray(values j2gort.Ext_com_google_common_collect_ImmutableList) []JSType {
	clinit_TemplateTypeMap()
	_ = values
	var n int32 = values.Size()
	_ = n
	if n == int32(0) {
		return StaticGet_TemplateTypeMap_EMPTY_VALUES()
	}
	var out []JSType = make([]JSType, n)
	_ = out
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			{
				a_9 := out
				i_10 := i
				a_9[i_10] = j2gort.CheckNotNullIface(j2gort.AsIface[JSType](values.Get(i)))
			}
		}
	}
	return out
}

// Java method: TemplateTypeMap.subMapStartsMask32()
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_subMapStartsMask32_9517() int32 {
	if this_.base != nil {
		return this_.base.priv_TemplateTypeMap_subMapStartsMask32_9517()
	}
	if this_.key0 == nil {
		return int32(0)
	}
	var e *TemplateTypeMap_Extra = this_.extra
	_ = e
	if (e != nil) && (e.keys != nil) {
		return e.subMapStartsMask
	}
	if (int32(this_.flags) & int32(int8(1))) != int32(0) {
		return int32(1)
	} else {
		return int32(0)
	}
}

// Java method: TemplateTypeMap.subMapStartsOverflowOrNull()
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_subMapStartsOverflowOrNull_9518() *j2gort.Ext_java_util_BitSet {
	if this_.base != nil {
		return this_.base.priv_TemplateTypeMap_subMapStartsOverflowOrNull_9518()
	}
	var e *TemplateTypeMap_Extra = this_.extra
	_ = e
	if e != nil {
		return e.subMapStartsOverflow
	} else {
		return nil
	}
}

// Java method: TemplateTypeMap.isSubMapStart(int index)
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_isSubMapStart_9519(index int32) bool {
	_ = index
	if this_.base != nil {
		return this_.base.priv_TemplateTypeMap_isSubMapStart_9519(index)
	}
	if this_.key0 == nil {
		return false
	}
	var e *TemplateTypeMap_Extra = this_.extra
	_ = e
	if (e == nil) || (e.keys == nil) {
		return ((index == int32(0)) && ((int32(this_.flags) & int32(int8(1))) != int32(0)))
	}
	if index < int32(32) {
		return ((e.subMapStartsMask & (javalang.Shl32(int32(1), int(index)))) != int32(0))
	}
	return ((e.subMapStartsOverflow != nil) && e.subMapStartsOverflow.Get_1(index))
}

// Java static method: TemplateTypeMap.mayContainReplacerTarget(com.google.javascript.rhino.jstype.JSType type)
func TemplateTypeMap_MayContainReplacerTarget(type_ JSType) bool {
	clinit_TemplateTypeMap()
	_ = type_
	switch (type_.GetJSTypeBase().FastTypeClass()).Ordinal() {
	case int32(20) /* JSTypeClass.TEMPLATE */, int32(21) /* JSTypeClass.TEMPLATIZED */, int32(22) /* JSTypeClass.UNION */, int32(6) /* JSTypeClass.FUNCTION */, int32(15) /* JSTypeClass.PROXY_OBJECT */, int32(11) /* JSTypeClass.NO_RESOLVED */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: TemplateTypeMap.definitelyUnchangedByTotalReplacer(com.google.javascript.rhino.jstype.JSType type, boolean isReplacementHop)
func TemplateTypeMap_DefinitelyUnchangedByTotalReplacer(type_ JSType, isReplacementHop bool) bool {
	clinit_TemplateTypeMap()
	_ = type_
	_ = isReplacementHop
	if TemplateTypeMap_MayContainReplacerTarget(type_) {
		var union *UnionType
		_ = union
		var ok_11 bool
		union, ok_11 = j2gort.TryCast[*UnionType](type_)
		if ok_11 {
			{
				var i int32 = int32(0)
				_ = i
				var n int32 = union.GetAlternatesCount()
				_ = n
				for ; i < n; i = i + 1 {
					var alt JSType = union.GetAlternateAt(i)
					_ = alt
					if TemplateTypeMap_MayContainReplacerTarget(alt) || alt.GetJSTypeBase().IsCheckedUnknownType() {
						return false
					}
				}
			}
			return true
		}
		return false
	}
	return (!(isReplacementHop) || !(type_.GetJSTypeBase().IsCheckedUnknownType()))
}

// Java method: TemplateTypeMap.ensureResolvedValues()
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_ensureResolvedValues_9520() []JSType {
	var state int32 = ((javalang.Ushr32(int32(this_.flags), int(int32(1)))) & int32(3))
	_ = state
	if state == int32(1) {
		return nil
	}
	if state == int32(2) {
		return this_.extra.resolvedValues
	}
	var resolved []JSType = this_.priv_TemplateTypeMap_computeResolvedValues_9521()
	_ = resolved
	if resolved == nil {
		this_.flags = int8(((int32(this_.flags) & ^(javalang.Shl32(int32(3), int(int32(1))))) | (javalang.Shl32(int32(1), int(int32(1))))))
		return nil
	}
	var e *TemplateTypeMap_Extra = this_.extra
	_ = e
	if e == nil {
		e = NewTemplateTypeMap_Extra(nil, nil, nil, int32(0), nil)
		this_.extra = e
	}
	{
		r_12 := e
		r_12.resolvedValues = resolved
	}
	this_.flags = int8(((int32(this_.flags) & ^(javalang.Shl32(int32(3), int(int32(1))))) | (javalang.Shl32(int32(2), int(int32(1))))))
	return resolved
}

// Java method: TemplateTypeMap.computeResolvedValues()
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_computeResolvedValues_9521() []JSType {
	var nValues int32 = this_.GetTemplateValuesSize()
	_ = nValues
	if nValues == int32(0) {
		return nil
	}
	var resolvedValues []JSType = nil
	_ = resolvedValues
	var replacer *TemplateTypeReplacer = nil
	_ = replacer
	var reg *JSTypeRegistry = this_.priv_TemplateTypeMap_registry_9516()
	_ = reg
	defer func() {
		if replacer != nil {
			reg.ReleaseTemplateTypeReplacer(replacer)
		}
	}()
	{
		{
			var i int32 = int32(0)
			_ = i
			for ; i < nValues; i = i + 1 {
				var templateValue JSType = this_.GetTemplateValueAt(i)
				_ = templateValue
				var resolved JSType = nil
				_ = resolved
				if TemplateTypeMap_DefinitelyUnchangedByTotalReplacer(templateValue, false) {
					resolved = templateValue
				} else {
					if replacer == nil {
						replacer = reg.AcquireTemplateTypeReplacer(this_, false, false, true)
					}
					var nextSubMap int32 = this_.priv_TemplateTypeMap_nextSubMapStart_9523((i + int32(1)))
					_ = nextSubMap
					replacer.SetKeyType(this_.GetTemplateKeyAt(i), nextSubMap)
					resolved = replacer.Replace(templateValue)
				}
				if resolvedValues != nil {
					{
						a_13 := resolvedValues
						i_14 := i
						a_13[i_14] = resolved
					}
				} else if !(resolved == templateValue || ((resolved == nil || resolved.IsJavaNil_()) && (templateValue == nil || templateValue.IsJavaNil_()))) {
					resolvedValues = make([]JSType, nValues)
					{
						var j int32 = int32(0)
						_ = j
						for ; j < i; j = j + 1 {
							{
								a_15 := resolvedValues
								i_16 := j
								a_15[i_16] = this_.GetTemplateValueAt(j)
							}
						}
					}
					{
						a_17 := resolvedValues
						i_18 := i
						a_17[i_18] = resolved
					}
				}
			}
		}
		return resolvedValues
	}
	panic("unreachable")
}

// Java method: TemplateTypeMap.getResolvedTemplateTypeAt(int index)
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_getResolvedTemplateTypeAt_9522(index int32) JSType {
	_ = index
	var state int32 = ((javalang.Ushr32(int32(this_.flags), int(int32(1)))) & int32(3))
	_ = state
	if state == int32(1) {
		return this_.GetTemplateValueAt(index)
	}
	if state == int32(2) {
		return this_.extra.resolvedValues[index]
	}
	var rawVal JSType = this_.GetTemplateValueAt(index)
	_ = rawVal
	if TemplateTypeMap_DefinitelyUnchangedByTotalReplacer(rawVal, false) {
		return rawVal
	}
	if rawVal.GetJSTypeBase().FastTypeClass() == StaticGet_JSTypeClass_TEMPLATE() {
		var keyType *TemplateType = this_.GetTemplateKeyAt(index)
		_ = keyType
		var boundary int32 = this_.priv_TemplateTypeMap_nextSubMapStart_9523((index + int32(1)))
		_ = boundary
		var seen0 *TemplateType = nil
		_ = seen0
		var seen1 *TemplateType = nil
		_ = seen1
		var seen2 *TemplateType = nil
		_ = seen2
		var seen3 *TemplateType = nil
		_ = seen3
		var cur JSType = rawVal
		_ = cur
		{
			var hop int32 = int32(0)
			_ = hop
			for ; hop < int32(4); hop = hop + 1 {
				var tt *TemplateType = j2gort.CastPtr[TemplateType](cur)
				_ = tt
				if !(tt == keyType) && !(this_.HasTemplateKey_2(tt, boundary)) {
					return ToIface_JSType(tt)
				}
				if (((tt == seen0) || (tt == seen1)) || (tt == seen2)) || (tt == seen3) {
					return ToIface_JSType(tt)
				}
				var targetIdx int32 = this_.priv_TemplateTypeMap_getTemplateTypeIndex_9551(tt)
				_ = targetIdx
				if targetIdx == int32(-1) {
					return this_.priv_TemplateTypeMap_registry_9516().GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
				}
				var next JSType = this_.GetTemplateValueAt(targetIdx)
				_ = next
				if (((((next == keyType || ((next == nil || next.IsJavaNil_()) && (keyType == nil || keyType.IsJavaNil_()))) || (next == tt || ((next == nil || next.IsJavaNil_()) && (tt == nil || tt.IsJavaNil_())))) || (next == seen0 || ((next == nil || next.IsJavaNil_()) && (seen0 == nil || seen0.IsJavaNil_())))) || (next == seen1 || ((next == nil || next.IsJavaNil_()) && (seen1 == nil || seen1.IsJavaNil_())))) || (next == seen2 || ((next == nil || next.IsJavaNil_()) && (seen2 == nil || seen2.IsJavaNil_())))) || (next == seen3 || ((next == nil || next.IsJavaNil_()) && (seen3 == nil || seen3.IsJavaNil_()))) {
					return ToIface_JSType(tt)
				}
				if TemplateTypeMap_DefinitelyUnchangedByTotalReplacer(next, true) {
					return next
				}
				if next.GetJSTypeBase().FastTypeClass() != StaticGet_JSTypeClass_TEMPLATE() {
					break
				}
				if hop == int32(0) {
					seen0 = tt
				} else if hop == int32(1) {
					seen1 = tt
				} else if hop == int32(2) {
					seen2 = tt
				} else {
					seen3 = tt
				}
				cur = next
			}
		}
	}
	var resolved []JSType = this_.priv_TemplateTypeMap_ensureResolvedValues_9520()
	_ = resolved
	if resolved != nil {
		return resolved[index]
	} else {
		return rawVal
	}
}

// Java method: TemplateTypeMap.nextSubMapStart(int index)
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_nextSubMapStart_9523(index int32) int32 {
	_ = index
	var nKeys int32 = this_.Size()
	_ = nKeys
	for (index < nKeys) && !(this_.priv_TemplateTypeMap_isSubMapStart_9519(index)) {
		index = index + 1
	}
	return index
}

// Java method: TemplateTypeMap.copyKeysInto(com.google.javascript.rhino.jstype.TemplateType[] dst, int dstOffset)
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_copyKeysInto_9524(dst []*TemplateType, dstOffset int32) {
	_ = dst
	_ = dstOffset
	if this_.base != nil {
		this_.base.priv_TemplateTypeMap_copyKeysInto_9524(dst, dstOffset)
		return
	}
	if this_.key0 == nil {
		return
	}
	var e *TemplateTypeMap_Extra = this_.extra
	_ = e
	if (e != nil) && (e.keys != nil) {
		func() {
			acSrc_19 := e.keys
			acSP_20 := int32(0)
			acDst_21 := dst
			acDP_22 := dstOffset
			acLen_23 := int32(len(e.keys))
			j2gort.CheckNotNullSlice(acSrc_19)
			j2gort.CheckNotNullSlice(acDst_21)
			if acSP_20 < 0 || acDP_22 < 0 || acLen_23 < 0 || int64(acSP_20)+int64(acLen_23) > int64(len(acSrc_19)) || int64(acDP_22)+int64(acLen_23) > int64(len(acDst_21)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_21[int(acDP_22):int(acDP_22)+int(acLen_23)], acSrc_19[int(acSP_20):int(acSP_20)+int(acLen_23)])
		}()
	} else {
		{
			a_24 := dst
			i_25 := dstOffset
			a_24[i_25] = this_.key0
		}
	}
}

// Java method: TemplateTypeMap.copyValuesInto(com.google.javascript.rhino.jstype.JSType[] dst, int dstOffset)
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_copyValuesInto_9525(dst []JSType, dstOffset int32) {
	_ = dst
	_ = dstOffset
	var pos int32 = dstOffset
	_ = pos
	if this_.base != nil {
		this_.base.priv_TemplateTypeMap_copyValuesInto_9525(dst, pos)
		pos = (pos + this_.base.GetTemplateValuesSize())
	}
	switch this_.ownValKind {
	case int8(int32(0)):
	case int8(int32(1)):
		{
			a_26 := dst
			i_27 := pos
			a_26[i_27] = this_.val0
		}
	case int8(int32(2)):
		{
			a_28 := dst
			i_29 := pos
			a_28[i_29] = this_.val0
		}
		{
			a_30 := dst
			i_31 := (pos + int32(1))
			a_30[i_31] = this_.val1
		}
	default:
		var own []JSType = this_.extra.ownValues
		_ = own
		func() {
			acSrc_32 := own
			acSP_33 := int32(0)
			acDst_34 := dst
			acDP_35 := pos
			acLen_36 := int32(len(own))
			j2gort.CheckNotNullSlice(acSrc_32)
			j2gort.CheckNotNullSlice(acDst_34)
			if acSP_33 < 0 || acDP_35 < 0 || acLen_36 < 0 || int64(acSP_33)+int64(acLen_36) > int64(len(acSrc_32)) || int64(acDP_35)+int64(acLen_36) > int64(len(acDst_34)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_34[int(acDP_35):int(acDP_35)+int(acLen_36)], acSrc_32[int(acSP_33):int(acSP_33)+int(acLen_36)])
		}()
	}
}

// Java method: TemplateTypeMap.copyFilledWithValues(com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.JSType> values)
func (this_ *TemplateTypeMap) CopyFilledWithValues(values j2gort.Ext_com_google_common_collect_ImmutableList) *TemplateTypeMap {
	_ = values
	j2gort.CheckNotNull(values)
	var unfilled int32 = this_.priv_TemplateTypeMap_numUnfilledTemplateKeys_9547()
	_ = unfilled
	var valuesSize int32 = values.Size()
	_ = valuesSize
	var requiredUnknownCount int32 = (unfilled - valuesSize)
	_ = requiredUnknownCount
	if !(requiredUnknownCount >= int32(0)) {
		j2gort.CheckArgumentMsg(false, j2gort.ValueOfInt(requiredUnknownCount))
	}
	if unfilled == int32(0) {
		return this_
	}
	if !(this_.base == nil) {
		j2gort.CheckState(false)
	}
	var existingSize int32 = this_.GetTemplateValuesSize()
	_ = existingSize
	if unfilled == int32(1) {
		var v0 JSType
		_ = v0
		if valuesSize == int32(1) {
			v0 = j2gort.CheckNotNullIface(j2gort.AsIface[JSType](values.Get(int32(0))))
		} else {
			v0 = this_.DefaultValueType(this_.GetTemplateKeyAt(existingSize))
		}
		return NewTemplateTypeMap_7(this_, this_.key0, v0, nil, int8(1), int8(0), nil)
	}
	if unfilled == int32(2) {
		var v0_2 JSType
		_ = v0_2
		if valuesSize >= int32(1) {
			v0_2 = j2gort.CheckNotNullIface(j2gort.AsIface[JSType](values.Get(int32(0))))
		} else {
			v0_2 = this_.DefaultValueType(this_.GetTemplateKeyAt(existingSize))
		}
		var v1 JSType
		_ = v1
		if valuesSize >= int32(2) {
			v1 = j2gort.CheckNotNullIface(j2gort.AsIface[JSType](values.Get(int32(1))))
		} else {
			v1 = this_.DefaultValueType(this_.GetTemplateKeyAt((existingSize + int32(1))))
		}
		return NewTemplateTypeMap_7(this_, this_.key0, v0_2, v1, int8(2), int8(0), nil)
	}
	var addedValues []JSType = make([]JSType, unfilled)
	_ = addedValues
	{
		var i int32 = int32(0)
		_ = i
		for ; i < valuesSize; i = i + 1 {
			{
				a_37 := addedValues
				i_38 := i
				a_37[i_38] = j2gort.CheckNotNullIface(j2gort.AsIface[JSType](values.Get(i)))
			}
		}
	}
	{
		var i_2 int32 = valuesSize
		_ = i_2
		for ; i_2 < unfilled; i_2 = i_2 + 1 {
			{
				a_39 := addedValues
				i_40 := i_2
				a_39[i_40] = this_.DefaultValueType(this_.GetTemplateKeyAt((existingSize + i_2)))
			}
		}
	}
	var val_41 TemplateTypeMap_Extra
	InitTemplateTypeMap_Extra(&val_41, nil, nil, addedValues, int32(0), nil)
	var extra *TemplateTypeMap_Extra = &val_41
	_ = extra
	return NewTemplateTypeMap_7(this_, this_.key0, nil, nil, int8(3), int8(0), extra)
}

// Java method: TemplateTypeMap.asKeysOnlySingleSubMap()
func (this_ *TemplateTypeMap) AsKeysOnlySingleSubMap() *TemplateTypeMap {
	var n int32 = this_.Size()
	_ = n
	if n == int32(0) {
		return this_
	}
	if (this_.base == nil) && (int32(this_.ownValKind) == int32(0)) {
		if (n == int32(1)) && ((int32(this_.flags) & int32(int8(1))) != int32(0)) {
			return this_
		}
		var e *TemplateTypeMap_Extra = this_.extra
		_ = e
		if (((n <= int32(32)) && (e != nil)) && (e.subMapStartsMask == int32(1))) && (e.subMapStartsOverflow == nil) {
			return this_
		}
	}
	if n == int32(1) {
		return TemplateTypeMap_CreateSingleKey(this_.key0, int32(1))
	}
	var copiedKeys []*TemplateType = make([]*TemplateType, n)
	_ = copiedKeys
	this_.priv_TemplateTypeMap_copyKeysInto_9524(copiedKeys, int32(0))
	return TemplateTypeMap_CreateFromArrays(this_.priv_TemplateTypeMap_registry_9516(), copiedKeys, StaticGet_TemplateTypeMap_EMPTY_VALUES(), int32(1), nil)
}

// Java method: TemplateTypeMap.copyWithExtension(com.google.javascript.rhino.jstype.TemplateTypeMap extension)
func (this_ *TemplateTypeMap) CopyWithExtension_1(extension *TemplateTypeMap) *TemplateTypeMap {
	_ = extension
	if extension.IsEmpty() && (this_.priv_TemplateTypeMap_numUnfilledTemplateKeys_9547() == int32(0)) {
		return this_
	}
	if this_.IsEmpty() && (this_.priv_TemplateTypeMap_registry_9516() == extension.priv_TemplateTypeMap_registry_9516()) {
		return extension
	}
	var extKeyCount int32 = extension.Size()
	_ = extKeyCount
	var extValCount int32 = extension.GetTemplateValuesSize()
	_ = extValCount
	var extKeys []*TemplateType = nil
	_ = extKeys
	if extKeyCount == int32(0) {
		extKeys = StaticGet_TemplateTypeMap_EMPTY_KEYS()
	} else if ((extension.base == nil) && (extension.extra != nil)) && (extension.extra.keys != nil) {
		extKeys = extension.extra.keys
	} else {
		extKeys = make([]*TemplateType, extKeyCount)
		extension.priv_TemplateTypeMap_copyKeysInto_9524(extKeys, int32(0))
	}
	var extVals []JSType = nil
	_ = extVals
	if extValCount == int32(0) {
		extVals = StaticGet_TemplateTypeMap_EMPTY_VALUES()
	} else if (((extension.base == nil) && (int32(extension.ownValKind) == int32(3))) && (extension.extra != nil)) && (extension.extra.ownValues != nil) {
		extVals = extension.extra.ownValues
	} else {
		extVals = make([]JSType, extValCount)
		extension.priv_TemplateTypeMap_copyValuesInto_9525(extVals, int32(0))
	}
	var extMask int32
	_ = extMask
	if extKeyCount <= int32(32) {
		extMask = extension.priv_TemplateTypeMap_subMapStartsMask32_9517()
	} else {
		extMask = int32(0)
	}
	var extOverflow *j2gort.Ext_java_util_BitSet = extension.priv_TemplateTypeMap_subMapStartsOverflowOrNull_9518()
	_ = extOverflow
	if (extKeyCount > int32(32)) && (extOverflow == nil) {
		extOverflow = j2gort.NewExt_java_util_BitSet_1(extKeyCount)
		{
			var i int32 = int32(0)
			_ = i
			for ; i < extKeyCount; i = i + 1 {
				if extension.priv_TemplateTypeMap_isSubMapStart_9519(i) {
					if i < int32(32) {
						extMask = (extMask | (javalang.Shl32(int32(1), int(i))))
					} else {
						extOverflow.Set_1(i)
					}
				}
			}
		}
	}
	return this_.priv_TemplateTypeMap_copyWithExtensionArrays_9531(extKeys, extVals, extMask, extOverflow)
}

// Java method: TemplateTypeMap.copyWithExtensionKeysOnly(com.google.javascript.rhino.jstype.TemplateType[] keys)
func (this_ *TemplateTypeMap) CopyWithExtensionKeysOnly(keys []*TemplateType) *TemplateTypeMap {
	_ = keys
	j2gort.CheckNotNullSlice(keys)
	var newKeyCount int32 = int32(len(keys))
	_ = newKeyCount
	if newKeyCount == int32(0) {
		if this_.priv_TemplateTypeMap_numUnfilledTemplateKeys_9547() == int32(0) {
			return this_
		}
		return this_.priv_TemplateTypeMap_copyWithExtensionArrays_9531(StaticGet_TemplateTypeMap_EMPTY_KEYS(), StaticGet_TemplateTypeMap_EMPTY_VALUES(), int32(0), nil)
	}
	if this_.IsEmpty() && (newKeyCount == int32(1)) {
		return TemplateTypeMap_CreateSingleKey(j2gort.CheckNotNullPtr(keys[int32(0)]), int32(1))
	}
	var safeKeys []*TemplateType = make([]*TemplateType, newKeyCount)
	_ = safeKeys
	{
		var i int32 = int32(0)
		_ = i
		for ; i < newKeyCount; i = i + 1 {
			{
				a_42 := safeKeys
				i_43 := i
				a_42[i_43] = j2gort.CheckNotNullPtr(keys[i])
			}
		}
	}
	if this_.IsEmpty() {
		return TemplateTypeMap_CreateFromArrays(this_.priv_TemplateTypeMap_registry_9516(), safeKeys, StaticGet_TemplateTypeMap_EMPTY_VALUES(), int32(1), nil)
	}
	return this_.priv_TemplateTypeMap_copyWithExtensionArrays_9531(safeKeys, StaticGet_TemplateTypeMap_EMPTY_VALUES(), int32(1), nil)
}

// Java method: TemplateTypeMap.copyWithExtension(com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.TemplateType> keys, com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.JSType> values)
func (this_ *TemplateTypeMap) CopyWithExtension_2(keys j2gort.Ext_com_google_common_collect_ImmutableList, values j2gort.Ext_com_google_common_collect_ImmutableList) *TemplateTypeMap {
	_ = keys
	_ = values
	j2gort.CheckNotNull(keys)
	j2gort.CheckNotNull(values)
	var newKeyCount int32 = keys.Size()
	_ = newKeyCount
	var newValCount int32 = values.Size()
	_ = newValCount
	var extendedUnfilledCount int32 = (newKeyCount - newValCount)
	_ = extendedUnfilledCount
	if !(extendedUnfilledCount >= int32(0)) {
		j2gort.CheckArgumentMsg(false, j2gort.ValueOfInt(extendedUnfilledCount))
	}
	if newKeyCount == int32(0) {
		if this_.priv_TemplateTypeMap_numUnfilledTemplateKeys_9547() == int32(0) {
			return this_
		}
		return this_.priv_TemplateTypeMap_copyWithExtensionArrays_9531(StaticGet_TemplateTypeMap_EMPTY_KEYS(), StaticGet_TemplateTypeMap_EMPTY_VALUES(), int32(0), nil)
	}
	if this_.IsEmpty() && (newKeyCount == int32(1)) {
		var k0 *TemplateType = j2gort.CheckNotNullPtr(j2gort.CastPtr[TemplateType](keys.Get(int32(0))))
		_ = k0
		if newValCount == int32(0) {
			return TemplateTypeMap_CreateSingleKey(k0, int32(1))
		}
		return TemplateTypeMap_CreateSingleKeyAndValue(k0, j2gort.CheckNotNullIface(j2gort.AsIface[JSType](values.Get(int32(0)))), int32(1))
	}
	var keyArr []*TemplateType = TemplateTypeMap_CopyKeysListToArray(keys)
	_ = keyArr
	var valArr []JSType = TemplateTypeMap_CopyValuesListToArray(values)
	_ = valArr
	if this_.IsEmpty() {
		return TemplateTypeMap_CreateFromArrays(this_.priv_TemplateTypeMap_registry_9516(), keyArr, valArr, int32(1), nil)
	}
	return this_.priv_TemplateTypeMap_copyWithExtensionArrays_9531(keyArr, valArr, int32(1), nil)
}

// Java method: TemplateTypeMap.copyWithExtensionArrays(com.google.javascript.rhino.jstype.TemplateType[] keys, com.google.javascript.rhino.jstype.JSType[] values, int newStartsMask, java.util.@org.jspecify.annotations.Nullable BitSet newStartsOverflow)
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_copyWithExtensionArrays_9531(keys []*TemplateType, values []JSType, newStartsMask int32, newStartsOverflow *j2gort.Ext_java_util_BitSet) *TemplateTypeMap {
	_ = keys
	_ = values
	_ = newStartsMask
	_ = newStartsOverflow
	var newKeyCount int32 = int32(len(keys))
	_ = newKeyCount
	var newValCount int32 = int32(len(values))
	_ = newValCount
	var extendedUnfilledCount int32 = (newKeyCount - newValCount)
	_ = extendedUnfilledCount
	if !(extendedUnfilledCount >= int32(0)) {
		j2gort.CheckArgumentMsg(false, j2gort.ValueOfInt(extendedUnfilledCount))
	}
	if (this_.priv_TemplateTypeMap_numUnfilledTemplateKeys_9547() == int32(0)) && (newKeyCount == int32(0)) {
		return this_
	}
	var oldKeyCount int32 = this_.Size()
	_ = oldKeyCount
	var oldValCount int32 = this_.GetTemplateValuesSize()
	_ = oldValCount
	var totalKeys int32 = (oldKeyCount + newKeyCount)
	_ = totalKeys
	var extendedKeys []*TemplateType = nil
	_ = extendedKeys
	if (((newKeyCount == int32(0)) && (this_.base == nil)) && (this_.extra != nil)) && (this_.extra.keys != nil) {
		extendedKeys = this_.extra.keys
	} else if oldKeyCount == int32(0) {
		extendedKeys = keys
	} else {
		extendedKeys = make([]*TemplateType, totalKeys)
		this_.priv_TemplateTypeMap_copyKeysInto_9524(extendedKeys, int32(0))
		if newKeyCount > int32(0) {
			func() {
				acSrc_44 := keys
				acSP_45 := int32(0)
				acDst_46 := extendedKeys
				acDP_47 := oldKeyCount
				acLen_48 := newKeyCount
				j2gort.CheckNotNullSlice(acSrc_44)
				j2gort.CheckNotNullSlice(acDst_46)
				if acSP_45 < 0 || acDP_47 < 0 || acLen_48 < 0 || int64(acSP_45)+int64(acLen_48) > int64(len(acSrc_44)) || int64(acDP_47)+int64(acLen_48) > int64(len(acDst_46)) {
					panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
				}
				copy(acDst_46[int(acDP_47):int(acDP_47)+int(acLen_48)], acSrc_44[int(acSP_45):int(acSP_45)+int(acLen_48)])
			}()
		}
	}
	var extendedMask int32 = int32(0)
	_ = extendedMask
	var extendedOverflow *j2gort.Ext_java_util_BitSet = nil
	_ = extendedOverflow
	if totalKeys <= int32(32) {
		extendedOverflow = nil
		var oldMask int32 = this_.priv_TemplateTypeMap_subMapStartsMask32_9517()
		_ = oldMask
		if newKeyCount == int32(0) {
			extendedMask = (oldMask & (func() int32 {
				if oldKeyCount == int32(32) {
					return int32(-1)
				}
				return ((javalang.Shl32(int32(1), int(oldKeyCount))) - int32(1))
			}()))
		} else if oldKeyCount == int32(0) {
			extendedMask = (newStartsMask & (func() int32 {
				if newKeyCount == int32(32) {
					return int32(-1)
				}
				return ((javalang.Shl32(int32(1), int(newKeyCount))) - int32(1))
			}()))
		} else {
			var oldClean int32 = (oldMask & ((javalang.Shl32(int32(1), int(oldKeyCount))) - int32(1)))
			_ = oldClean
			var newClean int32 = (newStartsMask & ((javalang.Shl32(int32(1), int(newKeyCount))) - int32(1)))
			_ = newClean
			extendedMask = (oldClean | (javalang.Shl32(newClean, int(oldKeyCount))))
		}
	} else {
		extendedMask = int32(0)
		extendedOverflow = j2gort.NewExt_java_util_BitSet_1(totalKeys)
		{
			var i int32 = int32(0)
			_ = i
			for ; i < oldKeyCount; i = i + 1 {
				if this_.priv_TemplateTypeMap_isSubMapStart_9519(i) {
					if i < int32(32) {
						extendedMask = (extendedMask | (javalang.Shl32(int32(1), int(i))))
					} else {
						extendedOverflow.Set_1(i)
					}
				}
			}
		}
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < newKeyCount; i_2 = i_2 + 1 {
				var start bool
				_ = start
				if i_2 < int32(32) {
					start = ((newStartsMask & (javalang.Shl32(int32(1), int(i_2)))) != int32(0))
				} else {
					start = ((newStartsOverflow != nil) && newStartsOverflow.Get_1(i_2))
				}
				if start {
					var dst int32 = (oldKeyCount + i_2)
					_ = dst
					if dst < int32(32) {
						extendedMask = (extendedMask | (javalang.Shl32(int32(1), int(dst))))
					} else {
						extendedOverflow.Set_1(dst)
					}
				}
			}
		}
	}
	var totalValues int32 = (oldKeyCount + newValCount)
	_ = totalValues
	var reg *JSTypeRegistry = this_.priv_TemplateTypeMap_registry_9516()
	_ = reg
	if totalValues == int32(0) {
		return TemplateTypeMap_CreateFromArrays(reg, extendedKeys, StaticGet_TemplateTypeMap_EMPTY_VALUES(), extendedMask, extendedOverflow)
	}
	var extendedValues []JSType = make([]JSType, totalValues)
	_ = extendedValues
	if oldValCount > int32(0) {
		this_.priv_TemplateTypeMap_copyValuesInto_9525(extendedValues, int32(0))
	}
	{
		var i_3 int32 = oldValCount
		_ = i_3
		for ; i_3 < oldKeyCount; i_3 = i_3 + 1 {
			{
				a_49 := extendedValues
				i_50 := i_3
				a_49[i_50] = this_.DefaultValueType(this_.GetTemplateKeyAt(i_3))
			}
		}
	}
	if newValCount > int32(0) {
		func() {
			acSrc_51 := values
			acSP_52 := int32(0)
			acDst_53 := extendedValues
			acDP_54 := oldKeyCount
			acLen_55 := newValCount
			j2gort.CheckNotNullSlice(acSrc_51)
			j2gort.CheckNotNullSlice(acDst_53)
			if acSP_52 < 0 || acDP_54 < 0 || acLen_55 < 0 || int64(acSP_52)+int64(acLen_55) > int64(len(acSrc_51)) || int64(acDP_54)+int64(acLen_55) > int64(len(acDst_53)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_53[int(acDP_54):int(acDP_54)+int(acLen_55)], acSrc_51[int(acSP_52):int(acSP_52)+int(acLen_55)])
		}()
	}
	return TemplateTypeMap_CreateFromArrays(reg, extendedKeys, extendedValues, extendedMask, extendedOverflow)
}

// Java method: TemplateTypeMap.copyWithoutKeys(java.util.Set<com.google.javascript.rhino.jstype.TemplateType> removals)
func (this_ *TemplateTypeMap) CopyWithoutKeys(removals j2gort.Ext_java_util_Set) *TemplateTypeMap {
	_ = removals
	if removals.IsEmpty() {
		return this_
	}
	var valCount int32 = this_.GetTemplateValuesSize()
	_ = valCount
	var keyCount int32 = this_.Size()
	_ = keyCount
	var removedCount int32 = int32(0)
	_ = removedCount
	{
		var i int32 = valCount
		_ = i
		for ; i < keyCount; i = i + 1 {
			if removals.Contains(j2gort.PtrToAny(this_.GetTemplateKeyAt(i))) {
				removedCount = removedCount + 1
			}
		}
	}
	if removedCount == int32(0) {
		return this_
	}
	var targetSize int32 = (keyCount - removedCount)
	_ = targetSize
	var reg *JSTypeRegistry = this_.priv_TemplateTypeMap_registry_9516()
	_ = reg
	if targetSize == int32(0) {
		var empty *TemplateTypeMap = reg.GetEmptyTemplateTypeMap()
		_ = empty
		if empty != nil {
			return empty
		} else {
			return NewTemplateTypeMap_1(reg)
		}
	}
	var newKeys []*TemplateType = make([]*TemplateType, targetSize)
	_ = newKeys
	var newMask int32 = int32(0)
	_ = newMask
	var newOverflow *j2gort.Ext_java_util_BitSet
	_ = newOverflow
	if targetSize > int32(32) {
		newOverflow = j2gort.NewExt_java_util_BitSet_1(targetSize)
	} else {
		newOverflow = nil
	}
	if valCount > int32(0) {
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < valCount; i_2 = i_2 + 1 {
				{
					a_56 := newKeys
					i_57 := i_2
					a_56[i_57] = this_.GetTemplateKeyAt(i_2)
				}
				if this_.priv_TemplateTypeMap_isSubMapStart_9519(i_2) {
					if i_2 < int32(32) {
						newMask = (newMask | (javalang.Shl32(int32(1), int(i_2))))
					} else {
						newOverflow.Set_1(i_2)
					}
				}
			}
		}
	}
	var newKeysSize int32 = int32(0)
	_ = newKeysSize
	var inNewSubmap bool = false
	_ = inNewSubmap
	{
		var i_3 int32 = valCount
		_ = i_3
		for ; i_3 < keyCount; i_3 = i_3 + 1 {
			var key *TemplateType = this_.GetTemplateKeyAt(i_3)
			_ = key
			inNewSubmap = (inNewSubmap || this_.priv_TemplateTypeMap_isSubMapStart_9519(i_3))
			if !(removals.Contains(j2gort.PtrToAny(key))) {
				var dst int32 = (valCount + func() int32 { old_58 := newKeysSize; nw_59 := old_58 + 1; newKeysSize = nw_59; return old_58 }())
				_ = dst
				{
					a_60 := newKeys
					i_61 := dst
					a_60[i_61] = key
				}
				if inNewSubmap {
					if dst < int32(32) {
						newMask = (newMask | (javalang.Shl32(int32(1), int(dst))))
					} else {
						newOverflow.Set_1(dst)
					}
				}
				inNewSubmap = false
			}
		}
	}
	if valCount == int32(0) {
		return TemplateTypeMap_CreateFromArrays(reg, newKeys, StaticGet_TemplateTypeMap_EMPTY_VALUES(), newMask, newOverflow)
	}
	var copiedVals []JSType = nil
	_ = copiedVals
	if (((this_.base == nil) && (int32(this_.ownValKind) == int32(3))) && (this_.extra != nil)) && (this_.extra.ownValues != nil) {
		copiedVals = this_.extra.ownValues
	} else {
		copiedVals = make([]JSType, valCount)
		this_.priv_TemplateTypeMap_copyValuesInto_9525(copiedVals, int32(0))
	}
	return TemplateTypeMap_CreateFromArrays(reg, newKeys, copiedVals, newMask, newOverflow)
}

// Java method: TemplateTypeMap.copyWithKeysFromAndValues(com.google.javascript.rhino.jstype.TemplateTypeMap keysSource, com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.JSType> values)
func (this_ *TemplateTypeMap) CopyWithKeysFromAndValues(keysSource *TemplateTypeMap, values j2gort.Ext_com_google_common_collect_ImmutableList) *TemplateTypeMap {
	_ = keysSource
	_ = values
	j2gort.CheckNotNullPtr(keysSource)
	j2gort.CheckNotNull(values)
	var newKeyCount int32 = keysSource.Size()
	_ = newKeyCount
	var newValCount int32 = values.Size()
	_ = newValCount
	var extendedUnfilledCount int32 = (newKeyCount - newValCount)
	_ = extendedUnfilledCount
	if !(extendedUnfilledCount >= int32(0)) {
		j2gort.CheckArgumentMsg(false, j2gort.ValueOfInt(extendedUnfilledCount))
	}
	if newKeyCount == int32(0) {
		if this_.priv_TemplateTypeMap_numUnfilledTemplateKeys_9547() == int32(0) {
			return this_
		}
		return this_.priv_TemplateTypeMap_copyWithExtensionArrays_9531(StaticGet_TemplateTypeMap_EMPTY_KEYS(), StaticGet_TemplateTypeMap_EMPTY_VALUES(), int32(0), nil)
	}
	if this_.IsEmpty() && (newKeyCount == int32(1)) {
		var k0 *TemplateType = keysSource.GetTemplateKeyAt(int32(0))
		_ = k0
		if newValCount == int32(0) {
			return TemplateTypeMap_CreateSingleKey(k0, int32(1))
		}
		return TemplateTypeMap_CreateSingleKeyAndValue(k0, j2gort.CheckNotNullIface(j2gort.AsIface[JSType](values.Get(int32(0)))), int32(1))
	}
	var keyArr []*TemplateType = nil
	_ = keyArr
	if ((keysSource.base == nil) && (keysSource.extra != nil)) && (keysSource.extra.keys != nil) {
		keyArr = keysSource.extra.keys
	} else {
		keyArr = make([]*TemplateType, newKeyCount)
		keysSource.priv_TemplateTypeMap_copyKeysInto_9524(keyArr, int32(0))
	}
	var valArr []JSType = TemplateTypeMap_CopyValuesListToArray(values)
	_ = valArr
	if this_.IsEmpty() {
		return TemplateTypeMap_CreateFromArrays(this_.priv_TemplateTypeMap_registry_9516(), keyArr, valArr, int32(1), nil)
	}
	return this_.priv_TemplateTypeMap_copyWithExtensionArrays_9531(keyArr, valArr, int32(1), nil)
}

// Java method: TemplateTypeMap.getTemplateKeysCount()
func (this_ *TemplateTypeMap) GetTemplateKeysCount() int32 {
	return this_.Size()
}

// Java method: TemplateTypeMap.size()
func (this_ *TemplateTypeMap) Size() int32 {
	if this_.base != nil {
		return this_.base.Size()
	}
	if this_.key0 == nil {
		return int32(0)
	}
	var e *TemplateTypeMap_Extra = this_.extra
	_ = e
	if (e != nil) && (e.keys != nil) {
		return int32(len(e.keys))
	} else {
		return int32(1)
	}
}

// Java method: TemplateTypeMap.isEmpty()
func (this_ *TemplateTypeMap) IsEmpty() bool {
	return (this_.key0 == nil)
}

// Java method: TemplateTypeMap.getTemplateKeyAt(int index)
func (this_ *TemplateTypeMap) GetTemplateKeyAt(index int32) *TemplateType {
	_ = index
	if this_.base != nil {
		return this_.base.GetTemplateKeyAt(index)
	}
	var e *TemplateTypeMap_Extra = this_.extra
	_ = e
	if (e != nil) && (e.keys != nil) {
		return e.keys[index]
	}
	if (index == int32(0)) && (this_.key0 != nil) {
		return this_.key0
	}
	panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
}

// Java method: TemplateTypeMap.ownValuesSize()
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_ownValuesSize_9538() int32 {
	var kind int32 = int32(this_.ownValKind)
	_ = kind
	if kind < int32(3) {
		return kind
	} else {
		return int32(len(this_.extra.ownValues))
	}
}

// Java method: TemplateTypeMap.getTemplateValuesSize()
func (this_ *TemplateTypeMap) GetTemplateValuesSize() int32 {
	var own int32 = this_.priv_TemplateTypeMap_ownValuesSize_9538()
	_ = own
	if this_.base != nil {
		return (this_.base.GetTemplateValuesSize() + own)
	} else {
		return own
	}
}

// Java method: TemplateTypeMap.getTemplateValueAt(int index)
func (this_ *TemplateTypeMap) GetTemplateValueAt(index int32) JSType {
	_ = index
	var localIdx int32 = index
	_ = localIdx
	if this_.base != nil {
		var prefix int32 = this_.base.GetTemplateValuesSize()
		_ = prefix
		if index < prefix {
			return this_.base.GetTemplateValueAt(index)
		}
		localIdx = (index - prefix)
	}
	var kind int32 = int32(this_.ownValKind)
	_ = kind
	if kind < int32(3) {
		if (localIdx == int32(0)) && (kind >= int32(1)) {
			return this_.val0
		}
		if (localIdx == int32(1)) && (kind == int32(2)) {
			return this_.val1
		}
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
	}
	return this_.extra.ownValues[localIdx]
}

// Java method: TemplateTypeMap.getTemplateKeysSlice(int fromIndex, int toIndex)
func (this_ *TemplateTypeMap) GetTemplateKeysSlice(fromIndex int32, toIndex int32) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = fromIndex
	_ = toIndex
	var nKeys int32 = this_.Size()
	_ = nKeys
	if !(((fromIndex >= int32(0)) && (toIndex <= nKeys)) && (fromIndex <= toIndex)) {
		j2gort.CheckArgument(false)
	}
	var len_ int32 = (toIndex - fromIndex)
	_ = len_
	if len_ == int32(0) {
		return StaticGet_TemplateTypeMap_EMPTY_KEY_LIST()
	}
	if len_ == int32(1) {
		return j2gort.ImmListOfPtr(this_.GetTemplateKeyAt(fromIndex))
	}
	if len_ == int32(2) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(this_.GetTemplateKeyAt(fromIndex)), j2gort.PtrToAny(this_.GetTemplateKeyAt((fromIndex + int32(1)))))
	}
	if len_ == int32(3) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_3(j2gort.PtrToAny(this_.GetTemplateKeyAt(fromIndex)), j2gort.PtrToAny(this_.GetTemplateKeyAt((fromIndex + int32(1)))), j2gort.PtrToAny(this_.GetTemplateKeyAt((fromIndex + int32(2)))))
	}
	var b *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, len_)
	_ = b
	{
		var i int32 = fromIndex
		_ = i
		for ; i < toIndex; i = i + 1 {
			b.AddVoid_(j2gort.PtrToAny(this_.GetTemplateKeyAt(i)))
		}
	}
	return b.BuildImmList_()
}

// Java method: TemplateTypeMap.getTemplateKeys()
func (this_ *TemplateTypeMap) GetTemplateKeys() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.GetTemplateKeysSlice(int32(0), this_.Size())
}

// Java method: TemplateTypeMap.getTemplateValues()
func (this_ *TemplateTypeMap) GetTemplateValues() j2gort.Ext_com_google_common_collect_ImmutableList {
	var len_ int32 = this_.GetTemplateValuesSize()
	_ = len_
	if len_ == int32(0) {
		return StaticGet_TemplateTypeMap_EMPTY_VALUE_LIST()
	}
	if len_ == int32(1) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(this_.GetTemplateValueAt(int32(0)))
	}
	if len_ == int32(2) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(this_.GetTemplateValueAt(int32(0)), this_.GetTemplateValueAt(int32(1)))
	}
	var b *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, len_)
	_ = b
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			b.AddVoid_(this_.GetTemplateValueAt(i))
		}
	}
	return b.BuildImmList_()
}

// Java method: TemplateTypeMap.hasTemplateKey(com.google.javascript.rhino.jstype.TemplateType templateKey)
func (this_ *TemplateTypeMap) HasTemplateKey_1(templateKey *TemplateType) bool {
	_ = templateKey
	return this_.HasTemplateKey_2(templateKey, int32(0))
}

// Java method: TemplateTypeMap.hasTemplateKey(com.google.javascript.rhino.jstype.TemplateType templateKey, int nextSubMap)
func (this_ *TemplateTypeMap) HasTemplateKey_2(templateKey *TemplateType, nextSubMap int32) bool {
	_ = templateKey
	_ = nextSubMap
	if this_.base != nil {
		return this_.base.HasTemplateKey_2(templateKey, nextSubMap)
	}
	if this_.key0 == nil {
		return false
	}
	var start int32 = max(int32(0), nextSubMap)
	_ = start
	var e *TemplateTypeMap_Extra = this_.extra
	_ = e
	if (e != nil) && (e.keys != nil) {
		var keys []*TemplateType = e.keys
		_ = keys
		var keyCount int32 = int32(len(keys))
		_ = keyCount
		{
			var i int32 = start
			_ = i
			for ; i < keyCount; i = i + 1 {
				if templateKey == keys[i] {
					return true
				}
			}
		}
		return false
	}
	return ((start == int32(0)) && (templateKey == this_.key0))
}

// Java method: TemplateTypeMap.getTemplateKeyCountThisShouldAlwaysBeOneOrZeroButIsnt(com.google.javascript.rhino.jstype.TemplateType templateKey)
func (this_ *TemplateTypeMap) GetTemplateKeyCountThisShouldAlwaysBeOneOrZeroButIsnt(templateKey *TemplateType) int32 {
	_ = templateKey
	if this_.base != nil {
		return this_.base.GetTemplateKeyCountThisShouldAlwaysBeOneOrZeroButIsnt(templateKey)
	}
	if this_.key0 == nil {
		return int32(0)
	}
	var e *TemplateTypeMap_Extra = this_.extra
	_ = e
	if (e != nil) && (e.keys != nil) {
		var matches int32 = int32(0)
		_ = matches
		var keys []*TemplateType = e.keys
		_ = keys
		var keyCount int32 = int32(len(keys))
		_ = keyCount
		{
			var i int32 = int32(0)
			_ = i
			for ; i < keyCount; i = i + 1 {
				if templateKey == keys[i] {
					matches = matches + 1
				}
			}
		}
		return matches
	}
	if templateKey == this_.key0 {
		return int32(1)
	} else {
		return int32(0)
	}
}

// Java method: TemplateTypeMap.numUnfilledTemplateKeys()
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_numUnfilledTemplateKeys_9547() int32 {
	return (this_.Size() - this_.GetTemplateValuesSize())
}

// Java method: TemplateTypeMap.hasTemplateType(com.google.javascript.rhino.jstype.TemplateType key)
func (this_ *TemplateTypeMap) HasTemplateType(key *TemplateType) bool {
	_ = key
	return (this_.priv_TemplateTypeMap_getTemplateTypeIndex_9551(key) != int32(-1))
}

// Java method: TemplateTypeMap.getUnresolvedOriginalTemplateType(com.google.javascript.rhino.jstype.TemplateType key)
func (this_ *TemplateTypeMap) GetUnresolvedOriginalTemplateType(key *TemplateType) JSType {
	_ = key
	var index int32 = this_.priv_TemplateTypeMap_getTemplateTypeIndex_9551(key)
	_ = index
	if index == int32(-1) {
		return this_.priv_TemplateTypeMap_registry_9516().GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
	} else {
		return this_.GetTemplateValueAt(index)
	}
}

// Java method: TemplateTypeMap.getLastTemplateTypeKeyByName(java.lang.String keyName)
func (this_ *TemplateTypeMap) GetLastTemplateTypeKeyByName(keyName string) *TemplateType {
	_ = keyName
	if this_.base != nil {
		return this_.base.GetLastTemplateTypeKeyByName(keyName)
	}
	if this_.key0 == nil {
		return nil
	}
	var e *TemplateTypeMap_Extra = this_.extra
	_ = e
	if (e != nil) && (e.keys != nil) {
		var keys []*TemplateType = e.keys
		_ = keys
		{
			var i int32 = (int32(len(keys)) - int32(1))
			_ = i
			for ; i >= int32(0); i = i - 1 {
				var key *TemplateType = keys[i]
				_ = key
				if j2gort.CheckNotNullString(key.GetReferenceName()) == keyName {
					return key
				}
			}
		}
		return nil
	}
	if j2gort.CheckNotNullString(this_.key0.GetReferenceName()) == keyName {
		return this_.key0
	} else {
		return nil
	}
}

// Java method: TemplateTypeMap.getTemplateTypeIndex(com.google.javascript.rhino.jstype.TemplateType key)
func (this_ *TemplateTypeMap) priv_TemplateTypeMap_getTemplateTypeIndex_9551(key *TemplateType) int32 {
	_ = key
	if this_.key0 == nil {
		return int32(-1)
	}
	var maxIndex int32 = min(this_.Size(), this_.GetTemplateValuesSize())
	_ = maxIndex
	if maxIndex <= int32(0) {
		return int32(-1)
	}
	var keySource *TemplateTypeMap
	_ = keySource
	if this_.base != nil {
		keySource = this_.base
	} else {
		keySource = this_
	}
	var ke *TemplateTypeMap_Extra = keySource.extra
	_ = ke
	if (ke != nil) && (ke.keys != nil) {
		var keys []*TemplateType = ke.keys
		_ = keys
		{
			var i int32 = (maxIndex - int32(1))
			_ = i
			for ; i >= int32(0); i = i - 1 {
				if keys[i] == key {
					return i
				}
			}
		}
		return int32(-1)
	}
	if keySource.key0 == key {
		return int32(0)
	} else {
		return int32(-1)
	}
}

// Java method: TemplateTypeMap.getResolvedTemplateType(com.google.javascript.rhino.jstype.TemplateType key)
func (this_ *TemplateTypeMap) GetResolvedTemplateType(key *TemplateType) JSType {
	_ = key
	var index int32 = this_.priv_TemplateTypeMap_getTemplateTypeIndex_9551(key)
	_ = index
	if index == int32(-1) {
		return this_.DefaultValueType(key)
	}
	return this_.priv_TemplateTypeMap_getResolvedTemplateTypeAt_9522(index)
}

// Java method: TemplateTypeMap.hasAnyTemplateTypesInternal()
func (this_ *TemplateTypeMap) HasAnyTemplateTypesInternal() bool {
	var nValues int32 = this_.GetTemplateValuesSize()
	_ = nValues
	{
		var i int32 = int32(0)
		_ = i
		for ; i < nValues; i = i + 1 {
			if this_.priv_TemplateTypeMap_getResolvedTemplateTypeAt_9522(i).GetJSTypeBase().HasAnyTemplateTypes() {
				return true
			}
		}
	}
	var nKeys int32 = this_.Size()
	_ = nKeys
	{
		var i_2 int32 = nValues
		_ = i_2
		for ; i_2 < nKeys; i_2 = i_2 + 1 {
			if this_.GetTemplateKeyAt(i_2).GetBound().GetJSTypeBase().HasAnyTemplateTypes() {
				return true
			}
		}
	}
	return false
}

// Java method: TemplateTypeMap.toString()
func (this_ *TemplateTypeMap) String() string {
	var s string = ""
	_ = s
	var len_ int32 = this_.Size()
	_ = len_
	var nValues int32 = this_.GetTemplateValuesSize()
	_ = nValues
	s = (j2gort.StrValString(s) + "{ ")
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			var k *TemplateType = this_.GetTemplateKeyAt(i)
			_ = k
			s = (j2gort.StrValString(s) + "(")
			s = (j2gort.StrValString(s) + j2gort.StrVal(k))
			s = (j2gort.StrValString(s) + ",")
			s = (j2gort.StrValString(s) + j2gort.StrVal(func() any {
				if i < nValues {
					return this_.GetTemplateValueAt(i)
				}
				return ""
			}()))
			s = (j2gort.StrValString(s) + ",")
			if i < nValues {
				s = (j2gort.StrValString(s) + j2gort.StrVal(this_.priv_TemplateTypeMap_getResolvedTemplateTypeAt_9522(i)))
			} else {
				s = (j2gort.StrValString(s) + j2gort.StrVal(k.GetBound()))
			}
			s = (j2gort.StrValString(s) + ") ")
		}
	}
	s = (j2gort.StrValString(s) + "}")
	return s
}

// Java method: TemplateTypeMap.defaultValueType(com.google.javascript.rhino.jstype.TemplateType type)
func (this_ *TemplateTypeMap) DefaultValueType(type_ *TemplateType) JSType {
	_ = type_
	if type_.GetBound().GetJSTypeBase().IsUnknownType() {
		return this_.priv_TemplateTypeMap_registry_9516().GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
	} else {
		return ToIface_JSType(type_)
	}
}

// Java method: TemplateTypeMap.getIndicesOfSubmapsForTesting()
func (this_ *TemplateTypeMap) GetIndicesOfSubmapsForTesting() j2gort.Ext_com_google_common_collect_ImmutableList {
	var indices []int32 = []int32{}
	_ = indices
	var n int32 = this_.Size()
	_ = n
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			if this_.priv_TemplateTypeMap_isSubMapStart_9519(i) {
				indices = append(indices, i)
			}
		}
	}
	return j2gort.ImmListFromTypedSlice(indices)
}

// Java class: com.google.javascript.rhino.jstype.TemplateTypeMap.Extra
type TemplateTypeMap_Extra struct {
	emptyMapRegistry     *JSTypeRegistry              `gson:"emptyMapRegistry"`
	keys                 []*TemplateType              `gson:"keys"`
	ownValues            []JSType                     `gson:"ownValues"`
	resolvedValues       []JSType                     `gson:"resolvedValues"`
	subMapStartsOverflow *j2gort.Ext_java_util_BitSet `gson:"subMapStartsOverflow"`
	subMapStartsMask     int32                        `gson:"subMapStartsMask"`
	idHash_              atomic.Int32
}

func (this_ *TemplateTypeMap_Extra) BindSelf_(s any) {
	_ = s
}
func (this_ *TemplateTypeMap_Extra) IsJavaNil_() bool { return this_ == nil }
func (*TemplateTypeMap_Extra) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.TemplateTypeMap$Extra", "Extra"
}
func (this_ *TemplateTypeMap_Extra) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TemplateTypeMap_Extra) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.TemplateTypeMap.Extra", "com.google.javascript.rhino.jstype.TemplateTypeMap$Extra", func(v_ any) bool {
		return j2gort.InstanceOf[*TemplateTypeMap_Extra](v_)
	})
}

func (this_ *TemplateTypeMap_Extra) initFields_TemplateTypeMap_Extra() {
	_ = this_
}

// Java constructor: Extra(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSTypeRegistry emptyMapRegistry, com.google.javascript.rhino.jstype.TemplateType @org.jspecify.annotations.Nullable [] keys, com.google.javascript.rhino.jstype.JSType @org.jspecify.annotations.Nullable [] ownValues, int subMapStartsMask, java.util.@org.jspecify.annotations.Nullable BitSet subMapStartsOverflow)
func NewTemplateTypeMap_Extra(emptyMapRegistry *JSTypeRegistry, keys []*TemplateType, ownValues []JSType, subMapStartsMask int32, subMapStartsOverflow *j2gort.Ext_java_util_BitSet) *TemplateTypeMap_Extra {
	this_ := &TemplateTypeMap_Extra{}
	InitTemplateTypeMap_Extra(this_, emptyMapRegistry, keys, ownValues, subMapStartsMask, subMapStartsOverflow)
	return this_
}

func InitTemplateTypeMap_Extra(this_ *TemplateTypeMap_Extra, emptyMapRegistry *JSTypeRegistry, keys []*TemplateType, ownValues []JSType, subMapStartsMask int32, subMapStartsOverflow *j2gort.Ext_java_util_BitSet) {
	_ = emptyMapRegistry
	_ = keys
	_ = ownValues
	_ = subMapStartsMask
	_ = subMapStartsOverflow
	this_.initFields_TemplateTypeMap_Extra()
	this_.emptyMapRegistry = emptyMapRegistry
	this_.keys = keys
	this_.ownValues = ownValues
	this_.subMapStartsMask = subMapStartsMask
	this_.subMapStartsOverflow = subMapStartsOverflow
}
