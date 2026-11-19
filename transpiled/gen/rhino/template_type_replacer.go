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
 *   John Lenz
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/TemplateTypeReplacer.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.TemplateTypeReplacer
type TemplateTypeReplacer struct {
	registry                   *JSTypeRegistry                            `gson:"registry"`
	bindings                   *TemplateTypeMap                           `gson:"bindings"`
	keyType                    *TemplateType                              `gson:"keyType"`
	seenStack                  []JSType                                   `gson:"seenStack"`
	visitedObjectTypes         *j2gort.IdentityOrderedMap[JSType, JSType] `gson:"visitedObjectTypes"`
	ownSubMapBoundary          int32                                      `gson:"ownSubMapBoundary"`
	seenDepth                  int32                                      `gson:"seenDepth"`
	idHash_                    atomic.Int32
	visitProperties            bool `gson:"visitProperties"`
	useUnknownForMissingKeys   bool `gson:"useUnknownForMissingKeys"`
	useUnknownForMissingValues bool `gson:"useUnknownForMissingValues"`
	hasMadeReplacement         bool `gson:"hasMadeReplacement"`
}

func (this_ *TemplateTypeReplacer) BindSelf_(s any) {
	_ = s
}
func (this_ *TemplateTypeReplacer) IsJavaNil_() bool { return this_ == nil }
func (*TemplateTypeReplacer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.TemplateTypeReplacer", "TemplateTypeReplacer"
}
func (this_ *TemplateTypeReplacer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TemplateTypeReplacer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.TemplateTypeReplacer", "com.google.javascript.rhino.jstype.TemplateTypeReplacer", func(v_ any) bool {
		return j2gort.InstanceOf[*TemplateTypeReplacer](v_)
	})
}
func (*TemplateTypeReplacer) Implements_rhino_Visitor() {}

func StaticGet_TemplateTypeReplacer_PARTIAL_REPLACER_ACTIVE_MASK() int32 { return int32(12517) }
func (this_ *TemplateTypeReplacer) initFields_TemplateTypeReplacer() {
	_ = this_

	this_.ownSubMapBoundary = int32(-1)

}

// Java constructor: TemplateTypeReplacer(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.TemplateTypeMap bindings, boolean visitProperties, boolean useUnknownForMissingKeys, boolean useUnknownForMissingValues)
func NewTemplateTypeReplacer(registry *JSTypeRegistry, bindings *TemplateTypeMap, visitProperties bool, useUnknownForMissingKeys bool, useUnknownForMissingValues bool) *TemplateTypeReplacer {
	this_ := &TemplateTypeReplacer{}
	InitTemplateTypeReplacer(this_, registry, bindings, visitProperties, useUnknownForMissingKeys, useUnknownForMissingValues)
	return this_
}

func InitTemplateTypeReplacer(this_ *TemplateTypeReplacer, registry *JSTypeRegistry, bindings *TemplateTypeMap, visitProperties bool, useUnknownForMissingKeys bool, useUnknownForMissingValues bool) {
	_ = registry
	_ = bindings
	_ = visitProperties
	_ = useUnknownForMissingKeys
	_ = useUnknownForMissingValues
	this_.initFields_TemplateTypeReplacer()
	this_.registry = registry
	this_.bindings = bindings
	this_.visitProperties = visitProperties
	this_.useUnknownForMissingKeys = useUnknownForMissingKeys
	this_.useUnknownForMissingValues = useUnknownForMissingValues
}

// Java static method: TemplateTypeReplacer.mayBeAffectedByPartialReplacement(com.google.javascript.rhino.jstype.JSType type)
func TemplateTypeReplacer_MayBeAffectedByPartialReplacement(type_ JSType) bool {
	_ = type_
	var mask int32 = type_.GetJSTypeBase().typeKindMask
	_ = mask
	return (((mask & int32(12517)) != int32(0)) && (mask != (int32(1) | int32(1024))))
}

// Java static method: TemplateTypeReplacer.createPooled(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func TemplateTypeReplacer_CreatePooled(registry *JSTypeRegistry) *TemplateTypeReplacer {
	_ = registry
	return NewTemplateTypeReplacer(registry, nil, false, false, false)
}

// Java method: TemplateTypeReplacer.initForPool(com.google.javascript.rhino.jstype.TemplateTypeMap bindings, boolean visitProperties, boolean useUnknownForMissingKeys, boolean useUnknownForMissingValues)
func (this_ *TemplateTypeReplacer) InitForPool(bindings *TemplateTypeMap, visitProperties bool, useUnknownForMissingKeys bool, useUnknownForMissingValues bool) {
	_ = bindings
	_ = visitProperties
	_ = useUnknownForMissingKeys
	_ = useUnknownForMissingValues
	this_.bindings = bindings
	this_.visitProperties = visitProperties
	this_.useUnknownForMissingKeys = useUnknownForMissingKeys
	this_.useUnknownForMissingValues = useUnknownForMissingValues
	this_.hasMadeReplacement = false
	this_.keyType = nil
	this_.ownSubMapBoundary = int32(-1)
}

// Java method: TemplateTypeReplacer.resetForPool()
func (this_ *TemplateTypeReplacer) ResetForPool() {
	this_.bindings = nil
	this_.keyType = nil
	this_.ownSubMapBoundary = int32(-1)
	this_.hasMadeReplacement = false
	var depth int32 = this_.seenDepth
	_ = depth
	if (depth > int32(0)) && (this_.seenStack != nil) {
		j2gort.ArraysFillRangeZero(this_.seenStack, int32(0), depth)
		this_.seenDepth = int32(0)
	}
	if (this_.seenStack != nil) && (int32(len(this_.seenStack)) > int32(32)) {
		this_.seenStack = nil
	}
	if (this_.visitedObjectTypes != nil) && !(this_.visitedObjectTypes.IsEmpty()) {
		if this_.visitedObjectTypes.Size() > int32(64) {
			this_.visitedObjectTypes = nil
		} else {
			this_.visitedObjectTypes.Clear()
		}
	}
}

// Java static method: TemplateTypeReplacer.forInference(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.util.Map<com.google.javascript.rhino.jstype.TemplateType,com.google.javascript.rhino.jstype.JSType> bindings)
func TemplateTypeReplacer_ForInference(registry *JSTypeRegistry, bindings j2gort.Ext_java_util_Map) *TemplateTypeReplacer {
	_ = registry
	_ = bindings
	var map_ *TemplateTypeMap = TemplateTypeMap_FromInferenceMap(registry, bindings)
	_ = map_
	return NewTemplateTypeReplacer(registry, map_, true, true, true)
}

// Java static method: TemplateTypeReplacer.forTotalReplacement(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.TemplateTypeMap bindings)
func TemplateTypeReplacer_ForTotalReplacement(registry *JSTypeRegistry, bindings *TemplateTypeMap) *TemplateTypeReplacer {
	_ = registry
	_ = bindings
	return NewTemplateTypeReplacer(registry, bindings, false, false, true)
}

// Java static method: TemplateTypeReplacer.forPartialReplacement(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.TemplateTypeMap bindings)
func TemplateTypeReplacer_ForPartialReplacement(registry *JSTypeRegistry, bindings *TemplateTypeMap) *TemplateTypeReplacer {
	_ = registry
	_ = bindings
	return NewTemplateTypeReplacer(registry, bindings, false, false, false)
}

// Java method: TemplateTypeReplacer.hasMadeReplacement()
func (this_ *TemplateTypeReplacer) HasMadeReplacement() bool {
	return this_.hasMadeReplacement
}

// Java method: TemplateTypeReplacer.isSeenType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_isSeenType_9534(type_ JSType) bool {
	_ = type_
	var stack []JSType = this_.seenStack
	_ = stack
	if stack == nil {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
		var depth int32 = this_.seenDepth
		_ = depth
		for ; i < depth; i = i + 1 {
			if j2gort.IfaceEquals[JSType](stack[i], type_) {
				return true
			}
		}
	}
	return false
}

// Java method: TemplateTypeReplacer.enterSeenType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_enterSeenType_9535(type_ JSType) bool {
	_ = type_
	var stack []JSType = this_.seenStack
	_ = stack
	var depth int32 = this_.seenDepth
	_ = depth
	if stack == nil {
		stack = make([]JSType, int32(8))
		this_.seenStack = stack
	} else {
		{
			var i int32 = int32(0)
			_ = i
			for ; i < depth; i = i + 1 {
				if j2gort.IfaceEquals[JSType](stack[i], type_) {
					return false
				}
			}
		}
		if depth == int32(len(stack)) {
			stack = func() []JSType {
				arrSrc_1 := stack
				arrLen_2 := (depth * int32(2))
				j2gort.CheckNotNullSlice(arrSrc_1)
				if arrLen_2 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_3 := make([]JSType, arrLen_2)
				copy(arrOut_3, arrSrc_1)
				return arrOut_3
			}()
			this_.seenStack = stack
		}
	}
	{
		a_4 := stack
		i_5 := depth
		a_4[i_5] = type_
	}
	this_.seenDepth = (depth + int32(1))
	return true
}

// Java method: TemplateTypeReplacer.exitSeenType()
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_exitSeenType_9536() {
	var nextDepth int32 = (this_.seenDepth - int32(1))
	_ = nextDepth
	{
		a_6 := this_.seenStack
		i_7 := nextDepth
		a_6[i_7] = nil
	}
	this_.seenDepth = nextDepth
}

// Java method: TemplateTypeReplacer.replace(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TemplateTypeReplacer) Replace(type_ JSType) JSType {
	_ = type_
	if !(this_.visitProperties) && !(TemplateTypeReplacer_MayBeAffectedByPartialReplacement(type_)) {
		return type_
	}
	switch (type_.GetTypeClass()).Ordinal() {
	case int32(20) /* JSTypeClass.TEMPLATE */ :
		return j2gort.AsIface[JSType](this_.CaseTemplateType(j2gort.CastPtr[TemplateType](type_)))
	case int32(22) /* JSTypeClass.UNION */ :
		return j2gort.AsIface[JSType](this_.CaseUnionType(j2gort.CastPtr[UnionType](type_)))
	case int32(21) /* JSTypeClass.TEMPLATIZED */ :
		return j2gort.AsIface[JSType](this_.CaseTemplatizedType(j2gort.CastPtr[TemplatizedType](type_)))
	case int32(6) /* JSTypeClass.FUNCTION */ :
		return j2gort.AsIface[JSType](this_.CaseFunctionType(j2gort.CastPtr[FunctionType](type_)))
	case int32(16) /* JSTypeClass.RECORD */ :
		return j2gort.AsIface[JSType](this_.CaseObjectType(j2gort.AsIface[ObjectType](type_)))
	case int32(15) /* JSTypeClass.PROXY_OBJECT */ :
		return j2gort.AsIface[JSType](this_.CaseProxyObjectType(j2gort.AsIface[ProxyObjectType](type_)))
	case int32(23) /* JSTypeClass.UNKNOWN */ :
		return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_UNKNOWN_TYPE())
	case int32(10) /* JSTypeClass.NO_OBJECT */ :
		return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_NO_OBJECT_TYPE())
	case int32(0) /* JSTypeClass.ALL */ :
		return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_ALL_TYPE())
	case int32(2) /* JSTypeClass.BOOLEAN */ :
		return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_BOOLEAN_TYPE())
	case int32(12) /* JSTypeClass.NULL */ :
		return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_NULL_TYPE())
	case int32(13) /* JSTypeClass.NUMBER */ :
		return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_NUMBER_TYPE())
	case int32(3) /* JSTypeClass.BIGINT */ :
		return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_BIGINT_TYPE())
	case int32(17) /* JSTypeClass.STRING */ :
		return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_STRING_TYPE())
	case int32(24) /* JSTypeClass.VOID */ :
		return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_VOID_TYPE())
	case int32(9) /* JSTypeClass.NO */, int32(5) /* JSTypeClass.ENUM_ELEMENT */, int32(4) /* JSTypeClass.ENUM */, int32(7) /* JSTypeClass.INSTANCE_OBJECT */, int32(14) /* JSTypeClass.PROTOTYPE_OBJECT */, int32(11) /* JSTypeClass.NO_RESOLVED */, int32(8) /* JSTypeClass.NAMED */, int32(18) /* JSTypeClass.SYMBOL */, int32(19) /* JSTypeClass.WELL_KNOWN_SYMBOL */ :
		return type_
	case int32(1) /* JSTypeClass.ARROW */ :
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	panic("unreachable switch expression")
}

// Java method: TemplateTypeReplacer.caseNoType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *TemplateTypeReplacer) CaseNoType(type_ *FunctionType) any {
	_ = type_
	return j2gort.PtrToAny(type_)
}

// Java method: TemplateTypeReplacer.caseEnumElementType(com.google.javascript.rhino.jstype.EnumElementType type)
func (this_ *TemplateTypeReplacer) CaseEnumElementType(type_ *EnumElementType) any {
	_ = type_
	return j2gort.PtrToAny(type_)
}

// Java method: TemplateTypeReplacer.caseAllType()
func (this_ *TemplateTypeReplacer) CaseAllType() any {
	return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_ALL_TYPE())
}

// Java method: TemplateTypeReplacer.caseBooleanType()
func (this_ *TemplateTypeReplacer) CaseBooleanType() any {
	return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_BOOLEAN_TYPE())
}

// Java method: TemplateTypeReplacer.caseNoObjectType()
func (this_ *TemplateTypeReplacer) CaseNoObjectType() any {
	return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_NO_OBJECT_TYPE())
}

// Java method: TemplateTypeReplacer.caseFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *TemplateTypeReplacer) CaseFunctionType(type_ *FunctionType) any {
	_ = type_
	if this_.priv_TemplateTypeReplacer_isNativeFunctionType_9544(type_) || (!(type_.IsOrdinaryFunction()) && !(type_.IsConstructor())) {
		return j2gort.PtrToAny(type_)
	}
	if !(this_.priv_TemplateTypeReplacer_enterSeenType_9535(ToIface_JSType(type_))) {
		return j2gort.PtrToAny(type_)
	}
	defer func() {
		this_.priv_TemplateTypeReplacer_exitSeenType_9536()
	}()
	{
		return this_.priv_TemplateTypeReplacer_caseFunctionTypeUnguarded_9538(type_)
	}
	panic("unreachable")
}

// Java method: TemplateTypeReplacer.caseFunctionTypeUnguarded(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_caseFunctionTypeUnguarded_9538(type_ *FunctionType) JSType {
	_ = type_
	var changed bool = false
	_ = changed
	var beforeThis JSType = type_.GetTypeOfThis()
	_ = beforeThis
	var afterThis JSType = this_.priv_TemplateTypeReplacer_coerseToThisType_9539(this_.Replace(beforeThis))
	_ = afterThis
	if !(beforeThis == afterThis || ((beforeThis == nil || beforeThis.IsJavaNil_()) && (afterThis == nil || afterThis.IsJavaNil_()))) {
		changed = true
	}
	var beforeReturn JSType = type_.GetReturnType()
	_ = beforeReturn
	var afterReturn JSType = this_.Replace(beforeReturn)
	_ = afterReturn
	if !(beforeReturn == afterReturn || ((beforeReturn == nil || beforeReturn.IsJavaNil_()) && (afterReturn == nil || afterReturn.IsJavaNil_()))) {
		changed = true
	}
	var origParams []FunctionType_Parameter = type_.GetParameterArray()
	_ = origParams
	var numParams int32 = int32(len(origParams))
	_ = numParams
	var newParams []FunctionType_Parameter = nil
	_ = newParams
	{
		var i int32 = int32(0)
		_ = i
		for ; i < numParams; i = i + 1 {
			var parameter FunctionType_Parameter = origParams[i]
			_ = parameter
			var beforeParamType JSType = parameter.GetJSType()
			_ = beforeParamType
			var afterParamType JSType = this_.Replace(beforeParamType)
			_ = afterParamType
			if !(beforeParamType == afterParamType || ((beforeParamType == nil || beforeParamType.IsJavaNil_()) && (afterParamType == nil || afterParamType.IsJavaNil_()))) {
				changed = true
				if newParams == nil {
					newParams = make([]FunctionType_Parameter, numParams)
					if i > int32(0) {
						func() {
							acSrc_8 := origParams
							acSP_9 := int32(0)
							acDst_10 := newParams
							acDP_11 := int32(0)
							acLen_12 := i
							j2gort.CheckNotNullSlice(acSrc_8)
							j2gort.CheckNotNullSlice(acDst_10)
							if acSP_9 < 0 || acDP_11 < 0 || acLen_12 < 0 || int64(acSP_9)+int64(acLen_12) > int64(len(acSrc_8)) || int64(acDP_11)+int64(acLen_12) > int64(len(acDst_10)) {
								panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
							}
							copy(acDst_10[int(acDP_11):int(acDP_11)+int(acLen_12)], acSrc_8[int(acSP_9):int(acSP_9)+int(acLen_12)])
						}()
					}
				}
				var adjustedType JSType
				_ = adjustedType
				if parameter.IsOptional() {
					adjustedType = this_.registry.CreateOptionalType(afterParamType)
				} else {
					adjustedType = afterParamType
				}
				{
					a_13 := newParams
					i_14 := i
					a_13[i_14] = FunctionType_Parameter_Create(adjustedType, parameter.IsOptional(), parameter.IsVariadic())
				}
			} else if newParams != nil {
				{
					a_15 := newParams
					i_16 := i
					a_15[i_16] = parameter
				}
			}
		}
	}
	if changed {
		var fb_ FunctionType_Builder
		type_.InitBuilder_(&fb_)
		(&fb_).WithParameterArray(func() []FunctionType_Parameter {
			if newParams != nil {
				return newParams
			}
			return origParams
		}())
		(&fb_).WithReturnType_1(afterReturn)
		(&fb_).WithTypeOfThis(afterThis)
		(&fb_).WithIsAbstract(false)
		return ToIface_JSType(j2gort.CastPtr[FunctionType]((&fb_).Build()))
	}
	return ToIface_JSType(type_)
}

// Java method: TemplateTypeReplacer.coerseToThisType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_coerseToThisType_9539(type_ JSType) JSType {
	_ = type_
	if type_ != nil && !type_.IsJavaNil_() {
		return type_
	} else {
		return Upcast_JSType(this_.registry.GetNativeObjectType(StaticGet_JSTypeNative_UNKNOWN_TYPE()))
	}
}

// Java method: TemplateTypeReplacer.caseObjectType(com.google.javascript.rhino.jstype.ObjectType objType)
func (this_ *TemplateTypeReplacer) CaseObjectType(objType ObjectType) any {
	_ = objType
	if ((!(this_.visitProperties) || objType.IsNominalType()) || j2gort.InstanceOf[ProxyObjectType](objType)) || !(objType.GetObjectTypeBase().IsRecordType()) {
		return objType
	}
	if !(this_.priv_TemplateTypeReplacer_enterSeenType_9535(Upcast_JSType(objType))) {
		return objType
	}
	defer func() {
		this_.priv_TemplateTypeReplacer_exitSeenType_9536()
	}()
	{
		return this_.priv_TemplateTypeReplacer_caseObjectTypeUnguarded_9540(objType)
	}
	panic("unreachable")
}

// Java method: TemplateTypeReplacer.caseObjectTypeUnguarded(com.google.javascript.rhino.jstype.ObjectType objType)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_caseObjectTypeUnguarded_9540(objType ObjectType) JSType {
	_ = objType
	var cached JSType = this_.priv_TemplateTypeReplacer_getVisitedObjectTypeOrNull_9547(Upcast_JSType(objType))
	_ = cached
	if cached != nil && !cached.IsJavaNil_() {
		return cached
	}
	var changed bool = false
	_ = changed
	var val_17 RecordTypeBuilder
	InitRecordTypeBuilder(&val_17, this_.registry)
	var builder *RecordTypeBuilder = &val_17
	_ = builder
	var proto PrototypeObjectType
	_ = proto
	var ownMap *PropertyMap
	_ = ownMap
	var ok_18 bool
	proto, ok_18 = j2gort.TryCast[PrototypeObjectType](objType)
	if ok_18 {
		ownMap = proto.GetPrototypeObjectTypeBase().GetOwnPropertyMapOrNull()
	} else {
		ownMap = nil
	}
	if (ownMap != nil) && ownMap.HasFlatProps() {
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = ownMap.GetOwnPropertiesCount()
			_ = n
			for ; i < n; i = i + 1 {
				var slot *Property = ownMap.GetFlatPropAt(i)
				_ = slot
				var prop string = slot.GetName()
				_ = prop
				var propertyNode Node = slot.GetNode()
				_ = propertyNode
				var beforeType JSType = objType.GetObjectTypeBase().GetPropertyType_1_string(prop)
				_ = beforeType
				var afterType JSType = this_.Replace(beforeType)
				_ = afterType
				if !(beforeType == afterType || ((beforeType == nil || beforeType.IsJavaNil_()) && (afterType == nil || afterType.IsJavaNil_()))) {
					changed = true
				}
				builder.AddProperty(prop, afterType, propertyNode)
			}
		}
	} else if (ownMap != nil) || !(j2gort.InstanceOf[PrototypeObjectType](objType)) {
		{
			cur_19 := j2gort.NewIterCursor(j2gort.CheckNotNull(objType.GetOwnPropertyNames()))
			for cur_19.HasNext() {
				prop_2 := cur_19.NextStr()
				_ = prop_2
				var propertyNode_2 Node = objType.GetObjectTypeBase().GetPropertyNode_1_string(prop_2)
				_ = propertyNode_2
				var beforeType_2 JSType = objType.GetObjectTypeBase().GetPropertyType_1_string(prop_2)
				_ = beforeType_2
				var afterType_2 JSType = this_.Replace(beforeType_2)
				_ = afterType_2
				if !(beforeType_2 == afterType_2 || ((beforeType_2 == nil || beforeType_2.IsJavaNil_()) && (afterType_2 == nil || afterType_2.IsJavaNil_()))) {
					changed = true
				}
				builder.AddProperty(prop_2, afterType_2, propertyNode_2)
			}
		}
	}
	var result JSType
	_ = result
	if changed {
		result = builder.Build()
	} else {
		result = Upcast_JSType(objType)
	}
	this_.visitedObjectTypes.Put(objType, result)
	return result
}

// Java method: TemplateTypeReplacer.caseTemplatizedType(com.google.javascript.rhino.jstype.TemplatizedType type)
func (this_ *TemplateTypeReplacer) CaseTemplatizedType(type_ *TemplatizedType) any {
	_ = type_
	if !(this_.priv_TemplateTypeReplacer_enterSeenType_9535(ToIface_JSType(type_))) {
		return j2gort.PtrToAny(type_)
	}
	defer func() {
		this_.priv_TemplateTypeReplacer_exitSeenType_9536()
	}()
	{
		return this_.priv_TemplateTypeReplacer_caseTemplatizedTypeUnguarded_9541(type_)
	}
	panic("unreachable")
}

// Java method: TemplateTypeReplacer.caseTemplatizedTypeUnguarded(com.google.javascript.rhino.jstype.TemplatizedType type)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_caseTemplatizedTypeUnguarded_9541(type_ *TemplatizedType) JSType {
	_ = type_
	var changed bool = false
	_ = changed
	var beforeBaseType ObjectType = type_.GetReferencedType()
	_ = beforeBaseType
	var afterBaseType ObjectType = ObjectType_Cast(this_.Replace(Upcast_JSType(beforeBaseType)))
	_ = afterBaseType
	if !(beforeBaseType == afterBaseType || ((beforeBaseType == nil || beforeBaseType.IsJavaNil_()) && (afterBaseType == nil || afterBaseType.IsJavaNil_()))) {
		changed = true
	}
	var origTemplateTypes []JSType = j2gort.CheckNotNullPtr(type_).Slice_GetTemplateTypes()
	_ = origTemplateTypes
	var templateCount int32 = int32(len(origTemplateTypes))
	_ = templateCount
	var newTemplateTypes []JSType = nil
	_ = newTemplateTypes
	{
		var i int32 = int32(0)
		_ = i
		for ; i < templateCount; i = i + 1 {
			var beforeTemplateType JSType = origTemplateTypes[i]
			_ = beforeTemplateType
			var afterTemplateType JSType = this_.Replace(beforeTemplateType)
			_ = afterTemplateType
			if !(beforeTemplateType == afterTemplateType || ((beforeTemplateType == nil || beforeTemplateType.IsJavaNil_()) && (afterTemplateType == nil || afterTemplateType.IsJavaNil_()))) {
				changed = true
				if newTemplateTypes == nil {
					newTemplateTypes = make([]JSType, templateCount)
					{
						var j int32 = int32(0)
						_ = j
						for ; j < i; j = j + 1 {
							{
								a_20 := newTemplateTypes
								i_21 := j
								a_20[i_21] = origTemplateTypes[j]
							}
						}
					}
				}
			}
			if newTemplateTypes != nil {
				{
					a_22 := newTemplateTypes
					i_23 := i
					a_22[i_23] = afterTemplateType
				}
			}
		}
	}
	if changed {
		var updatedArgs []JSType
		_ = updatedArgs
		if newTemplateTypes != nil {
			updatedArgs = j2gort.ImmListCopyArrayToSlice[JSType](func() []any {
				arr_24 := newTemplateTypes
				if arr_24 == nil {
					return nil
				}
				out_25 := make([]any, len(arr_24))
				for i_, el_26 := range arr_24 {
					out_25[i_] = el_26
				}
				return out_25
			}())
		} else {
			updatedArgs = origTemplateTypes
		}
		type_ = this_.registry.CreateTemplatizedType_2_ObjectType_Ext_com_google_common_collect_ImmutableList(afterBaseType, j2gort.ImmListFromTypedSlice(updatedArgs))
	}
	return ToIface_JSType(type_)
}

// Java method: TemplateTypeReplacer.caseUnknownType()
func (this_ *TemplateTypeReplacer) CaseUnknownType() any {
	return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_UNKNOWN_TYPE())
}

// Java method: TemplateTypeReplacer.caseNullType()
func (this_ *TemplateTypeReplacer) CaseNullType() any {
	return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_NULL_TYPE())
}

// Java method: TemplateTypeReplacer.caseNumberType()
func (this_ *TemplateTypeReplacer) CaseNumberType() any {
	return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_NUMBER_TYPE())
}

// Java method: TemplateTypeReplacer.caseBigIntType()
func (this_ *TemplateTypeReplacer) CaseBigIntType() any {
	return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_BIGINT_TYPE())
}

// Java method: TemplateTypeReplacer.caseStringType()
func (this_ *TemplateTypeReplacer) CaseStringType() any {
	return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_STRING_TYPE())
}

// Java method: TemplateTypeReplacer.caseSymbolType(com.google.javascript.rhino.jstype.SymbolType type)
func (this_ *TemplateTypeReplacer) CaseSymbolType(type_ SymbolType) any {
	_ = type_
	return type_
}

// Java method: TemplateTypeReplacer.caseVoidType()
func (this_ *TemplateTypeReplacer) CaseVoidType() any {
	return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_VOID_TYPE())
}

// Java method: TemplateTypeReplacer.caseUnionType(com.google.javascript.rhino.jstype.UnionType type)
func (this_ *TemplateTypeReplacer) CaseUnionType(type_ *UnionType) any {
	_ = type_
	if !(this_.priv_TemplateTypeReplacer_enterSeenType_9535(ToIface_JSType(type_))) {
		return j2gort.PtrToAny(type_)
	}
	defer func() {
		this_.priv_TemplateTypeReplacer_exitSeenType_9536()
	}()
	{
		return this_.priv_TemplateTypeReplacer_caseUnionTypeUnguarded_9542(type_)
	}
	panic("unreachable")
}

// Java method: TemplateTypeReplacer.caseUnionTypeUnguarded(com.google.javascript.rhino.jstype.UnionType type)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_caseUnionTypeUnguarded_9542(type_ *UnionType) JSType {
	_ = type_
	var results *UnionType_Builder = nil
	_ = results
	defer func() {
		if results != nil {
			this_.registry.ReleaseUnionTypeBuilder(results)
		}
	}()
	{
		var alternateCount int32 = type_.GetAlternatesCount()
		_ = alternateCount
		{
			var i int32 = int32(0)
			_ = i
			for ; i < alternateCount; i = i + 1 {
				var alternative JSType = type_.GetAlternateAt(i)
				_ = alternative
				var replacement JSType = this_.Replace(alternative)
				_ = replacement
				if (results == nil) && !(replacement == alternative || ((replacement == nil || replacement.IsJavaNil_()) && (alternative == nil || alternative.IsJavaNil_()))) {
					results = this_.registry.AcquireUnionTypeBuilder()
					{
						var j int32 = int32(0)
						_ = j
						for ; j < i; j = j + 1 {
							results.AddAlternate(type_.GetAlternateAt(j))
						}
					}
				}
				if results != nil {
					results.AddAlternate(replacement)
				}
			}
		}
		if results != nil {
			return results.Build()
		}
		return ToIface_JSType(type_)
	}
	panic("unreachable")
}

// Java method: TemplateTypeReplacer.caseTemplateType(com.google.javascript.rhino.jstype.TemplateType type)
func (this_ *TemplateTypeReplacer) CaseTemplateType(type_ *TemplateType) any {
	_ = type_
	this_.hasMadeReplacement = true
	if !(type_ == this_.keyType) && !(this_.bindings.HasTemplateKey_2(type_, this_.ownSubMapBoundary)) {
		if this_.useUnknownForMissingKeys {
			return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_UNKNOWN_TYPE())
		} else {
			return ToIface_JSType(type_)
		}
	}
	if this_.priv_TemplateTypeReplacer_isSeenType_9534(ToIface_JSType(type_)) {
		return j2gort.PtrToAny(type_)
	} else if !(this_.bindings.HasTemplateType(type_)) {
		if this_.useUnknownForMissingValues {
			return this_.priv_TemplateTypeReplacer_getNativeType_9543(StaticGet_JSTypeNative_UNKNOWN_TYPE())
		} else {
			return ToIface_JSType(type_)
		}
	}
	var replacement JSType = this_.bindings.GetUnresolvedOriginalTemplateType(type_)
	_ = replacement
	if (replacement == this_.keyType || ((replacement == nil || replacement.IsJavaNil_()) && (this_.keyType == nil || this_.keyType.IsJavaNil_()))) || this_.priv_TemplateTypeReplacer_isSeenType_9534(replacement) {
		return j2gort.PtrToAny(type_)
	}
	this_.priv_TemplateTypeReplacer_enterSeenType_9535(ToIface_JSType(type_))
	var visitedReplacement JSType = nil
	_ = visitedReplacement
	func() {
		defer func() {
			this_.priv_TemplateTypeReplacer_exitSeenType_9536()
		}()
		visitedReplacement = this_.Replace(replacement)
	}()
	if !(!(visitedReplacement == this_.keyType || ((visitedReplacement == nil || visitedReplacement.IsJavaNil_()) && (this_.keyType == nil || this_.keyType.IsJavaNil_())))) {
		j2gort.CheckStateFmt(false, "Trying to replace key %s with the same value", j2gort.PtrToAny(this_.keyType))
	}
	return visitedReplacement
}

// Java method: TemplateTypeReplacer.getNativeType(com.google.javascript.rhino.jstype.JSTypeNative nativeType)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_getNativeType_9543(nativeType JSTypeNative) JSType {
	_ = nativeType
	return this_.registry.GetNativeType(nativeType)
}

// Java method: TemplateTypeReplacer.isNativeFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_isNativeFunctionType_9544(type_ *FunctionType) bool {
	_ = type_
	return type_.IsNativeObjectType()
}

// Java method: TemplateTypeReplacer.caseNamedType(com.google.javascript.rhino.jstype.NamedType type)
func (this_ *TemplateTypeReplacer) CaseNamedType(type_ *NamedType) any {
	_ = type_
	return j2gort.PtrToAny(type_)
}

// Java method: TemplateTypeReplacer.caseProxyObjectType(com.google.javascript.rhino.jstype.ProxyObjectType type)
func (this_ *TemplateTypeReplacer) CaseProxyObjectType(type_ ProxyObjectType) any {
	_ = type_
	if !(this_.priv_TemplateTypeReplacer_enterSeenType_9535(Upcast_JSType(type_))) {
		return type_
	}
	defer func() {
		this_.priv_TemplateTypeReplacer_exitSeenType_9536()
	}()
	{
		return this_.priv_TemplateTypeReplacer_caseProxyObjectTypeUnguarded_9545(type_)
	}
	panic("unreachable")
}

// Java method: TemplateTypeReplacer.caseProxyObjectTypeUnguarded(com.google.javascript.rhino.jstype.ProxyObjectType type)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_caseProxyObjectTypeUnguarded_9545(type_ ProxyObjectType) JSType {
	_ = type_
	var beforeType JSType = type_.GetProxyObjectTypeBase().GetReferencedTypeInternal()
	_ = beforeType
	var replacement JSType = this_.Replace(beforeType)
	_ = replacement
	if !(replacement == beforeType || ((replacement == nil || replacement.IsJavaNil_()) && (beforeType == nil || beforeType.IsJavaNil_(

	// Java method: TemplateTypeReplacer.setKeyType(com.google.javascript.rhino.jstype.TemplateType keyType, int ownSubMapBoundary)
	)))) {
		return replacement
	}
	return Upcast_JSType(type_)
}

func (this_ *TemplateTypeReplacer) SetKeyType(keyType *TemplateType, ownSubMapBoundary int32) {
	_ = keyType
	_ = ownSubMapBoundary
	this_.keyType = keyType
	this_.ownSubMapBoundary = ownSubMapBoundary
}

// Java method: TemplateTypeReplacer.getVisitedObjectTypeOrNull(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_getVisitedObjectTypeOrNull_9547(type_ JSType) JSType {
	_ = type_
	if this_.visitedObjectTypes == nil {
		this_.visitedObjectTypes = &j2gort.IdentityOrderedMap[JSType, JSType]{}
		return nil
	}
	return this_.visitedObjectTypes.Get(type_)
}
