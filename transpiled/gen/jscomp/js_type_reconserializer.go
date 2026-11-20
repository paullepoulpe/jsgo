/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/serialization/JSTypeReconserializer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_js_type_reconserializer_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_js_type_reconserializer_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_js_type_reconserializer_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_js_type_reconserializer_4 j2gort.Ext_java_util_function_Function
var lambdaSingleton_js_type_reconserializer_5 j2gort.Ext_java_util_function_Consumer

func init() {
	lambdaSingleton_js_type_reconserializer_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_6 any) any {
		r := j2gort.CastPtr[JSTypeReconserializer_SeenTypeRecord](lp_6)
		_ = r
		return j2gort.BoxNullableVal(r.colorId)
	})
	lambdaSingleton_js_type_reconserializer_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_7 any) any {
		r := j2gort.CastPtr[JSTypeReconserializer_SeenTypeRecord](lp_7)
		_ = r
		return j2gort.BoxNullableVal(r.colorId)
	})
	lambdaSingleton_js_type_reconserializer_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_11 any) any {
		r := j2gort.CastPtr[JSTypeReconserializer_SeenTypeRecord](lp_11)
		_ = r
		return j2gort.BoxNullableVal(r.colorId)
	})
	lambdaSingleton_js_type_reconserializer_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_12 any) any {
		r := j2gort.CastPtr[JSTypeReconserializer_SeenTypeRecord](lp_12)
		_ = r
		return j2gort.BoxNullableVal(r.colorId)
	})
	lambdaSingleton_js_type_reconserializer_5 = j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_51 any) {
		t := j2gort.AsIface[pkg_rhino.JSType](lp_51)
		_ = t
		if !(!(t.GetJSTypeBase().IsUnionType())) {
			j2gort.CheckStateMsg(false, t)
		}
	})
}

// Java class: com.google.javascript.jscomp.serialization.JSTypeReconserializer
type JSTypeReconserializer struct {
	registry                    *pkg_rhino.JSTypeRegistry                                                                `gson:"registry"`
	serializationMode           *SerializationOptions                                                                    `gson:"serializationMode"`
	invalidatingTypes           *InvalidatingTypes                                                                       `gson:"invalidatingTypes"`
	stringPoolBuilder           *StringPool_Builder                                                                      `gson:"stringPoolBuilder"`
	hasher                      *JSTypeColorIdHasher                                                                     `gson:"hasher"`
	shouldPropagatePropertyName j2gort.Ext_java_util_function_Predicate                                                  `gson:"shouldPropagatePropertyName"`
	astPropertySet              *SerializeTypesToPointers_AstPropertySet                                                 `gson:"astPropertySet"`
	unknownRecord               *JSTypeReconserializer_SeenTypeRecord                                                    `gson:"unknownRecord"`
	topFunctionRecord           *JSTypeReconserializer_SeenTypeRecord                                                    `gson:"topFunctionRecord"`
	typeToRecordById            []*JSTypeReconserializer_SeenTypeRecord                                                  `gson:"typeToRecordById"`
	fallbackTypeToRecordCache   *pkg_base.LinkedIdentityHashMap[pkg_rhino.JSType, *JSTypeReconserializer_SeenTypeRecord] `gson:"fallbackTypeToRecordCache"`
	seenTypeRecords             j2gort.OrderedMap[pkg_colors.ColorId, *JSTypeReconserializer_SeenTypeRecord]             `gson:"seenTypeRecords"`
	disambiguateEdges           j2gort.Ext_com_google_common_collect_SetMultimap                                         `gson:"disambiguateEdges"`
	scratchInstancePointers     *j2gort.Ext_java_util_LinkedHashSet                                                      `gson:"scratchInstancePointers"`
	scratchPrototypePointers    *j2gort.Ext_java_util_LinkedHashSet                                                      `gson:"scratchPrototypePointers"`
	scratchOwnProperties        *j2gort.Ext_java_util_LinkedHashSet                                                      `gson:"scratchOwnProperties"`
	scratchOwnPropNames         []string                                                                                 `gson:"scratchOwnPropNames"`
	idHash_                     atomic.Int32
	state                       JSTypeReconserializer_State `gson:"state"`
}

func (this_ *JSTypeReconserializer) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeReconserializer) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeReconserializer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.JSTypeReconserializer", "JSTypeReconserializer"
}
func (this_ *JSTypeReconserializer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeReconserializer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.JSTypeReconserializer", "com.google.javascript.jscomp.serialization.JSTypeReconserializer", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeReconserializer](v_)
	})
}

var static_JSTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_JSTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_JSTypeReconserializer()
	return static_JSTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP
}
func StaticSet_JSTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_JSTypeReconserializer()
	static_JSTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP = v
}
func StaticGet_JSTypeReconserializer_FORWARD_CANDIDATE_MASK() int32 { return int32(2102657) }

var clinitDone_JSTypeReconserializer uint32
var clinitState_JSTypeReconserializer j2gort.ClinitState

func clinit_JSTypeReconserializer() {
	if atomic.LoadUint32(&clinitDone_JSTypeReconserializer) != 0 {
		return
	}
	clinitSlow_JSTypeReconserializer()
}

func clinitSlow_JSTypeReconserializer() {
	if !j2gort.EnterClinit(&clinitState_JSTypeReconserializer, &clinitDone_JSTypeReconserializer) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSTypeReconserializer, &clinitDone_JSTypeReconserializer)
	static_JSTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder().Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_ALL_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_UNKNOWN()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_UNKNOWN()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_NO_OBJECT_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_UNKNOWN()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_NO_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_UNKNOWN()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_UNKNOWN()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_BIGINT_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_BIGINT()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_BOOLEAN()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_GBIGINT_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_GBIGINT()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_NUMBER()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_STRING()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_SYMBOL_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_SYMBOL()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_FUNCTION_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_TOP_FUNCTION()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_TOP_FUNCTION()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_PROTOTYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_OBJECT_FUNCTION_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_OBJECT_PROTOTYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))).BuildOrThrow())
}

func (this_ *JSTypeReconserializer) initFields_JSTypeReconserializer() {
	_ = this_
	this_.seenTypeRecords = j2gort.NewOrderedMapCap[pkg_colors.ColorId, *JSTypeReconserializer_SeenTypeRecord](int32(4096))
	this_.disambiguateEdges = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0())
	this_.scratchInstancePointers = j2gort.NewExt_java_util_LinkedHashSet()
	this_.scratchPrototypePointers = j2gort.NewExt_java_util_LinkedHashSet()
	this_.scratchOwnProperties = j2gort.NewExt_java_util_LinkedHashSet()
	this_.scratchOwnPropNames = j2gort.NewStringArray(int32(32))
	this_.state = StaticGet_JSTypeReconserializer_State_COLLECTING_TYPES()
}

// Java constructor: JSTypeReconserializer(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.jscomp.InvalidatingTypes invalidatingTypes, com.google.javascript.jscomp.serialization.StringPool.Builder stringPoolBuilder, java.util.function.@org.jspecify.annotations.Nullable Predicate<java.lang.String> shouldPropagatePropertyName, com.google.javascript.jscomp.serialization.SerializeTypesToPointers.@org.jspecify.annotations.Nullable AstPropertySet astPropertySet, com.google.javascript.jscomp.serialization.SerializationOptions serializationMode)
func NewJSTypeReconserializer(registry *pkg_rhino.JSTypeRegistry, invalidatingTypes *InvalidatingTypes, stringPoolBuilder *StringPool_Builder, shouldPropagatePropertyName j2gort.Ext_java_util_function_Predicate, astPropertySet *SerializeTypesToPointers_AstPropertySet, serializationMode *SerializationOptions) *JSTypeReconserializer {
	clinit_JSTypeReconserializer()
	this_ := &JSTypeReconserializer{}
	InitJSTypeReconserializer(this_, registry, invalidatingTypes, stringPoolBuilder, shouldPropagatePropertyName, astPropertySet, serializationMode)
	return this_
}

func InitJSTypeReconserializer(this_ *JSTypeReconserializer, registry *pkg_rhino.JSTypeRegistry, invalidatingTypes *InvalidatingTypes, stringPoolBuilder *StringPool_Builder, shouldPropagatePropertyName j2gort.Ext_java_util_function_Predicate, astPropertySet *SerializeTypesToPointers_AstPropertySet, serializationMode *SerializationOptions) {
	_ = registry
	_ = invalidatingTypes
	_ = stringPoolBuilder
	_ = shouldPropagatePropertyName
	_ = astPropertySet
	_ = serializationMode
	this_.initFields_JSTypeReconserializer()
	this_.registry = registry
	this_.hasher = NewJSTypeColorIdHasher(registry)
	this_.invalidatingTypes = invalidatingTypes
	this_.stringPoolBuilder = stringPoolBuilder
	this_.shouldPropagatePropertyName = shouldPropagatePropertyName
	this_.astPropertySet = astPropertySet
	this_.serializationMode = serializationMode
	this_.typeToRecordById = make([]*JSTypeReconserializer_SeenTypeRecord, max(int32(16384), (registry.GetNextTypeId()+int32(64))))
	this_.priv_JSTypeReconserializer_seedCachesWithAxiomaticTypes_7791()
	this_.unknownRecord = this_.seenTypeRecords.Get(pkg_colors.StaticGet_StandardColors_UNKNOWN().GetId())
	this_.topFunctionRecord = this_.seenTypeRecords.Get(pkg_colors.StaticGet_StandardColors_TOP_FUNCTION().GetId())
}

// Java method: JSTypeReconserializer.shouldPropagateProperty(java.lang.String ownProperty)
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_shouldPropagateProperty_7780(ownProperty string) bool {
	_ = ownProperty
	if this_.astPropertySet != nil {
		return this_.astPropertySet.Contains(ownProperty)
	}
	return this_.shouldPropagatePropertyName.Test(j2gort.BoxString(ownProperty))
}

// Java static method: JSTypeReconserializer.create(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.jscomp.InvalidatingTypes invalidatingTypes, com.google.javascript.jscomp.serialization.StringPool.Builder stringPoolBuilder, java.util.function.Predicate<java.lang.String> shouldPropagatePropertyName, com.google.javascript.jscomp.serialization.SerializationOptions serializationMode)
func JSTypeReconserializer_Create(registry *pkg_rhino.JSTypeRegistry, invalidatingTypes *InvalidatingTypes, stringPoolBuilder *StringPool_Builder, shouldPropagatePropertyName j2gort.Ext_java_util_function_Predicate, serializationMode *SerializationOptions) *JSTypeReconserializer {
	clinit_JSTypeReconserializer()
	_ = registry
	_ = invalidatingTypes
	_ = stringPoolBuilder
	_ = shouldPropagatePropertyName
	_ = serializationMode
	clinit_JSTypeReconserializer()
	var val_1 JSTypeReconserializer
	InitJSTypeReconserializer(&val_1, registry, invalidatingTypes, stringPoolBuilder, shouldPropagatePropertyName, nil, serializationMode)
	var serializer *JSTypeReconserializer = &val_1
	_ = serializer
	serializer.priv_JSTypeReconserializer_checkValidLinearTime_7792()
	return serializer
}

// Java static method: JSTypeReconserializer.createWithPropertySet(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.jscomp.InvalidatingTypes invalidatingTypes, com.google.javascript.jscomp.serialization.StringPool.Builder stringPoolBuilder, com.google.javascript.jscomp.serialization.SerializeTypesToPointers.AstPropertySet astPropertySet, com.google.javascript.jscomp.serialization.SerializationOptions serializationMode)
func JSTypeReconserializer_CreateWithPropertySet(registry *pkg_rhino.JSTypeRegistry, invalidatingTypes *InvalidatingTypes, stringPoolBuilder *StringPool_Builder, astPropertySet *SerializeTypesToPointers_AstPropertySet, serializationMode *SerializationOptions) *JSTypeReconserializer {
	clinit_JSTypeReconserializer()
	_ = registry
	_ = invalidatingTypes
	_ = stringPoolBuilder
	_ = astPropertySet
	_ = serializationMode
	clinit_JSTypeReconserializer()
	var val_2 JSTypeReconserializer
	InitJSTypeReconserializer(&val_2, registry, invalidatingTypes, stringPoolBuilder, nil, astPropertySet, serializationMode)
	var serializer *JSTypeReconserializer = &val_2
	_ = serializer
	serializer.priv_JSTypeReconserializer_checkValidLinearTime_7792()
	return serializer
}

// Java method: JSTypeReconserializer.serializeType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeReconserializer) SerializeType(type_ pkg_rhino.JSType) int32 {
	_ = type_
	var record *JSTypeReconserializer_SeenTypeRecord = this_.priv_JSTypeReconserializer_recordType_7782(type_)
	_ = record
	return record.pointer
}

// Java method: JSTypeReconserializer.recordType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_recordType_7782(type_ pkg_rhino.JSType) *JSTypeReconserializer_SeenTypeRecord {
	_ = type_
	for {
		var mask int32 = type_.GetJSTypeBase().GetKindMask()
		_ = mask
		if (mask & int32(2102657)) != int32(0) {
			var forwardedType pkg_rhino.JSType = nil
			_ = forwardedType
			if (mask & int32(1)) == int32(0) {
				if (mask & int32(256)) != int32(0) {
					forwardedType = type_.ToMaybeEnumElementType().GetPrimitiveType()
				} else if (mask & int32(2097152)) != int32(0) {
					forwardedType = this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_SYMBOL_TYPE())
				} else {
					forwardedType = pkg_rhino.ToIface_JSType(type_.ToMaybeFunctionType().GetCanonicalRepresentation())
				}
			} else if (mask & int32(1024)) != int32(0) {
				forwardedType = type_.ToMaybeNamedType().GetReferencedType()
			} else if type_.GetJSTypeBase().IsEnumElementType() {
				forwardedType = type_.ToMaybeEnumElementType().GetPrimitiveType()
			} else if ((mask & int32(4096)) != int32(0)) || type_.GetJSTypeBase().IsTemplatizedType() {
				forwardedType = pkg_rhino.Upcast_JSType(type_.ToMaybeTemplatizedType().GetReferencedType())
			} else if type_.GetJSTypeBase().IsFunctionType() {
				forwardedType = pkg_rhino.ToIface_JSType(type_.ToMaybeFunctionType().GetCanonicalRepresentation())
			} else {
				forwardedType = nil
			}
			if forwardedType != nil && !forwardedType.IsJavaNil_() {
				type_ = forwardedType
				continue
			}
		}
		if type_.GetJSTypeBase().IsUnknownType() {
			return this_.unknownRecord
		}
		if (mask & int32(1)) == int32(0) {
			if (mask & int32(8)) != int32(0) {
				return this_.unknownRecord
			}
			if (mask & int32(128)) != int32(0) {
				var fn *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
				_ = fn
				if !(fn.HasInstanceType()) && (fn.GetClosurePrimitive() == 0) {
					return this_.topFunctionRecord
				}
			}
		} else {
			if (((mask & int32(8192)) != int32(0)) || type_.GetJSTypeBase().IsNoResolvedType()) || type_.GetJSTypeBase().IsTemplateType() {
				return this_.unknownRecord
			}
			if type_.GetJSTypeBase().IsFunctionType() {
				var fn_2 *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
				_ = fn_2
				if !(fn_2.HasInstanceType()) && (fn_2.GetClosurePrimitive() == 0) {
					return this_.topFunctionRecord
				}
			}
		}
		var typeId int32 = type_.GetJSTypeBase().GetTypeId()
		_ = typeId
		if typeId > int32(0) {
			var cache []*JSTypeReconserializer_SeenTypeRecord = this_.typeToRecordById
			_ = cache
			if typeId < int32(len(cache)) {
				var jstypeRecord *JSTypeReconserializer_SeenTypeRecord = cache[typeId]
				_ = jstypeRecord
				if jstypeRecord != nil {
					return jstypeRecord
				}
			}
		} else if this_.fallbackTypeToRecordCache != nil {
			var jstypeRecord_2 *JSTypeReconserializer_SeenTypeRecord = this_.fallbackTypeToRecordCache.Get(type_)
			_ = jstypeRecord_2
			if jstypeRecord_2 != nil {
				return jstypeRecord_2
			}
		}
		if (mask & int32(1)) == int32(0) {
			if (mask & int32(64)) != int32(0) {
				return this_.priv_JSTypeReconserializer_recordUnionType_7783(type_.ToMaybeUnionType())
			}
		} else if type_.GetJSTypeBase().IsUnionType() {
			return this_.priv_JSTypeReconserializer_recordUnionType_7783(type_.ToMaybeUnionType())
		}
		if type_.GetJSTypeBase().IsObjectType() {
			return this_.priv_JSTypeReconserializer_recordObjectType_7784(type_.GetJSTypeBase().ToMaybeObjectType())
		}
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", type_)))
	}
	panic("unreachable")
}

// Java method: JSTypeReconserializer.recordUnionType(com.google.javascript.rhino.jstype.UnionType type)
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_recordUnionType_7783(type_ *pkg_rhino.UnionType) *JSTypeReconserializer_SeenTypeRecord {
	_ = type_
	j2gort.CheckNotNullPtr(type_)
	var r0 *JSTypeReconserializer_SeenTypeRecord = nil
	_ = r0
	var r1 *JSTypeReconserializer_SeenTypeRecord = nil
	_ = r1
	var r2 *JSTypeReconserializer_SeenTypeRecord = nil
	_ = r2
	var r3 *JSTypeReconserializer_SeenTypeRecord = nil
	_ = r3
	var count int32 = int32(0)
	_ = count
	var altRecords *j2gort.Ext_java_util_LinkedHashSet = nil
	_ = altRecords
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = type_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			var alt *JSTypeReconserializer_SeenTypeRecord = this_.priv_JSTypeReconserializer_recordType_7782(type_.GetAlternateAt(i))
			_ = alt
			if alt.unionMembers == nil {
				if altRecords != nil {
					altRecords.Add_1(j2gort.PtrToAny(alt))
				} else if (((alt != r0) && (alt != r1)) && (alt != r2)) && (alt != r3) {
					switch count {
					case int32(0):
						r0 = alt
						count = int32(1)
					case int32(1):
						r1 = alt
						count = int32(2)
					case int32(2):
						r2 = alt
						count = int32(3)
					case int32(3):
						r3 = alt
						count = int32(4)
					default:
						altRecords = j2gort.NewExt_java_util_LinkedHashSet_1_int32(int32(8))
						altRecords.Add_1(j2gort.PtrToAny(r0))
						altRecords.Add_1(j2gort.PtrToAny(r1))
						altRecords.Add_1(j2gort.PtrToAny(r2))
						altRecords.Add_1(j2gort.PtrToAny(r3))
						altRecords.Add_1(j2gort.PtrToAny(alt))
					}
				}
			} else {
				{
					cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(alt.unionMembers))
					for cur_3.HasNext() {
						rawElem_4 := cur_3.Next()
						member := j2gort.CastPtr[JSTypeReconserializer_SeenTypeRecord](rawElem_4)
						_ = member
						if altRecords != nil {
							altRecords.Add_1(j2gort.PtrToAny(member))
						} else if (((member != r0) && (member != r1)) && (member != r2)) && (member != r3) {
							switch count {
							case int32(0):
								r0 = member
								count = int32(1)
							case int32(1):
								r1 = member
								count = int32(2)
							case int32(2):
								r2 = member
								count = int32(3)
							case int32(3):
								r3 = member
								count = int32(4)
							default:
								altRecords = j2gort.NewExt_java_util_LinkedHashSet_1_int32(int32(8))
								altRecords.Add_1(j2gort.PtrToAny(r0))
								altRecords.Add_1(j2gort.PtrToAny(r1))
								altRecords.Add_1(j2gort.PtrToAny(r2))
								altRecords.Add_1(j2gort.PtrToAny(r3))
								altRecords.Add_1(j2gort.PtrToAny(member))
							}
						}
					}
				}
			}
		}
	}
	if altRecords == nil {
		if count == int32(1) {
			return r0
		}
		var unionId pkg_colors.ColorId = (pkg_colors.ColorId{})
		_ = unionId
		switch count {
		case int32(2):
			unionId = pkg_colors.ColorId_Union2(r0.colorId, r1.colorId)
		case int32(3):
			unionId = pkg_colors.ColorId_Union3(r0.colorId, r1.colorId, r2.colorId)
		case int32(4):
			unionId = pkg_colors.ColorId_Union4(r0.colorId, r1.colorId, r2.colorId, r3.colorId)
		default:
			panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.ValueOfInt(count))))
		}
		var record *JSTypeReconserializer_SeenTypeRecord = this_.priv_JSTypeReconserializer_getOrCreateRecord_7786(unionId, pkg_rhino.ToIface_JSType(type_))
		_ = record
		if record.unionMembers == nil {
			{
				r_5 := record
				r_5.unionMembers = func() j2gort.Ext_com_google_common_collect_ImmutableSet {
					switch count {
					case int32(2):
						return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_2(j2gort.PtrToAny(r0), j2gort.PtrToAny(r1))
					case int32(3):
						return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_3(j2gort.PtrToAny(r0), j2gort.PtrToAny(r1), j2gort.PtrToAny(r2))
					case int32(4):
						return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4(j2gort.PtrToAny(r0), j2gort.PtrToAny(r1), j2gort.PtrToAny(r2), j2gort.PtrToAny(r3))
					default:
						panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.ValueOfInt(count))))
					}
					panic("unreachable switch expression")
				}()
			}
		} else if this_.serializationMode.RunValidation() {
			var expected j2gort.Ext_com_google_common_collect_ImmutableSet = nil
			_ = expected
			switch count {
			case int32(2):
				expected = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_2(j2gort.PtrToAny(r0), j2gort.PtrToAny(r1))
			case int32(3):
				expected = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_3(j2gort.PtrToAny(r0), j2gort.PtrToAny(r1), j2gort.PtrToAny(r2))
			case int32(4):
				expected = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4(j2gort.PtrToAny(r0), j2gort.PtrToAny(r1), j2gort.PtrToAny(r2), j2gort.PtrToAny(r3))
			default:
				panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.ValueOfInt(count))))
			}
			j2gort.CheckStateFmt(expected.Equals(record.unionMembers), "Unions with same ID must have same members: %s => %s == %s", j2gort.BoxNullableVal(unionId), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](expected.Stream().Map(lambdaSingleton_js_type_reconserializer_1).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet())), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](record.unionMembers.Stream().Map(lambdaSingleton_js_type_reconserializer_2).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet())))
		}
		return record
	}
	var alternateIds j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_BuilderWithExpectedSize(altRecords.Size())
	_ = alternateIds
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(altRecords))
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			altRecord := j2gort.CastPtr[JSTypeReconserializer_SeenTypeRecord](rawElem_9)
			_ = altRecord
			alternateIds.Add_1_any(j2gort.BoxNullableVal(altRecord.colorId))
		}
	}
	var unionId_2 pkg_colors.ColorId = pkg_colors.ColorId_Union(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](alternateIds.Build()))
	_ = unionId_2
	var record_2 *JSTypeReconserializer_SeenTypeRecord = this_.priv_JSTypeReconserializer_getOrCreateRecord_7786(unionId_2, pkg_rhino.ToIface_JSType(type_))
	_ = record_2
	if record_2.unionMembers == nil {
		{
			r_10 := record_2
			r_10.unionMembers = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](altRecords))
		}
	} else if this_.serializationMode.RunValidation() {
		j2gort.CheckStateFmt(altRecords.Equals(record_2.unionMembers), "Unions with same ID must have same members: %s => %s == %s", j2gort.BoxNullableVal(unionId_2), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](altRecords.Stream().Map(lambdaSingleton_js_type_reconserializer_3).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet())), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](record_2.unionMembers.Stream().Map(lambdaSingleton_js_type_reconserializer_4).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet())))
	}
	return record_2
}

// Java method: JSTypeReconserializer.recordObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_recordObjectType_7784(type_ pkg_rhino.ObjectType) *JSTypeReconserializer_SeenTypeRecord {
	_ = type_
	j2gort.CheckNotNullIface(type_)
	var id pkg_colors.ColorId = this_.hasher.HashObjectType(type_)
	_ = id
	var record *JSTypeReconserializer_SeenTypeRecord = this_.priv_JSTypeReconserializer_getOrCreateRecord_7786(id, pkg_rhino.Upcast_JSType(type_))
	_ = record
	this_.priv_JSTypeReconserializer_addSupertypeEdges_7785(type_, j2gort.ValueOfInt(record.pointer))
	if type_.GetObjectTypeBase().IsFunctionType() {
		var fnType *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
		_ = fnType
		if fnType.HasInstanceType() && !(j2gort.IfaceEquals[pkg_rhino.ObjectType](fnType.GetInstanceType(), nil)) {
			this_.SerializeType(pkg_rhino.Upcast_JSType(fnType.GetInstanceType()))
			this_.SerializeType(pkg_rhino.Upcast_JSType(fnType.GetPrototype()))
		}
	}
	return record
}

// Java method: JSTypeReconserializer.addSupertypeEdges(com.google.javascript.rhino.jstype.ObjectType subtype, java.lang.Integer serializedSubtype)
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_addSupertypeEdges_7785(subtype pkg_rhino.ObjectType, serializedSubtype *j2gort.BoxInt) {
	_ = subtype
	_ = serializedSubtype
	this_.disambiguateEdges.PutAll_2(j2gort.PtrToAny(serializedSubtype), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.priv_JSTypeReconserializer_ownAncestorInterfacesOf_7790(subtype)))
	if !(j2gort.IfaceEquals[pkg_rhino.ObjectType](subtype.GetImplicitPrototype(), nil)) {
		var supertype *j2gort.BoxInt = j2gort.ValueOfInt(this_.SerializeType(pkg_rhino.Upcast_JSType(subtype.GetImplicitPrototype())))
		_ = supertype
		this_.disambiguateEdges.Put(j2gort.PtrToAny(serializedSubtype), supertype)
	}
}

// Java method: JSTypeReconserializer.getOrCreateRecord(com.google.javascript.jscomp.colors.ColorId id, com.google.javascript.rhino.jstype.JSType jstype)
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_getOrCreateRecord_7786(id pkg_colors.ColorId, jstype pkg_rhino.JSType) *JSTypeReconserializer_SeenTypeRecord {
	_ = id
	_ = jstype
	j2gort.CheckNotNullIface(jstype)
	if !((StaticGet_JSTypeReconserializer_State_COLLECTING_TYPES() == this_.state) || (StaticGet_JSTypeReconserializer_State_GENERATING_POOL() == this_.state)) {
		j2gort.CheckState(false)
	}
	var record *JSTypeReconserializer_SeenTypeRecord = this_.seenTypeRecords.Get(id)
	_ = record
	if record == nil {
		var pointer int32 = this_.seenTypeRecords.Size()
		_ = pointer
		record = NewJSTypeReconserializer_SeenTypeRecord(id, pointer)
		this_.seenTypeRecords.Put(id, record)
	}
	if this_.priv_JSTypeReconserializer_putCachedTypeRecord_7787(jstype, record) == nil {
		record.AddJstype(jstype)
	}
	return record
}

// Java method: JSTypeReconserializer.putCachedTypeRecord(com.google.javascript.rhino.jstype.JSType jstype, com.google.javascript.jscomp.serialization.JSTypeReconserializer.SeenTypeRecord record)
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_putCachedTypeRecord_7787(jstype pkg_rhino.JSType, record *JSTypeReconserializer_SeenTypeRecord) *JSTypeReconserializer_SeenTypeRecord {
	_ = jstype
	_ = record
	var typeId int32 = jstype.GetJSTypeBase().GetTypeId()
	_ = typeId
	if typeId > int32(0) {
		var cache []*JSTypeReconserializer_SeenTypeRecord = this_.typeToRecordById
		_ = cache
		if typeId >= int32(len(cache)) {
			var nextCap int32 = (int32(len(cache)) * int32(2))
			_ = nextCap
			for nextCap <= typeId {
				nextCap = (nextCap * int32(2))
			}
			{

				arrSrc_13 := cache
				arrLen_14 := nextCap
				j2gort.CheckNotNullSlice(arrSrc_13)
				if arrLen_14 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_15 := make([]*JSTypeReconserializer_SeenTypeRecord, arrLen_14)
				copy(arrOut_15, arrSrc_13)
				cache = arrOut_15
			}

			this_.typeToRecordById = cache
		}
		var prev *JSTypeReconserializer_SeenTypeRecord = cache[typeId]
		_ = prev
		{
			a_16 := cache
			i_17 := typeId
			a_16[i_17] = record
		}
		return prev
	}
	var fallback *pkg_base.LinkedIdentityHashMap[pkg_rhino.JSType, *JSTypeReconserializer_SeenTypeRecord] = this_.fallbackTypeToRecordCache
	_ = fallback
	if fallback == nil {
		fallback = pkg_base.NewLinkedIdentityHashMap[pkg_rhino.JSType, *JSTypeReconserializer_SeenTypeRecord]()
		this_.fallbackTypeToRecordCache = fallback
	}
	return fallback.Put(jstype, record)
}

// Java method: JSTypeReconserializer.reconcileUnionTypes(com.google.javascript.jscomp.serialization.JSTypeReconserializer.SeenTypeRecord seen)
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_reconcileUnionTypes_7788(seen *JSTypeReconserializer_SeenTypeRecord) *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto {
	_ = seen
	var sortedMembers j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList_1_int32(seen.unionMembers.Size())
	_ = sortedMembers
	{
		cur_18 := j2gort.NewIterCursor(j2gort.CheckNotNull(seen.unionMembers))
		for cur_18.HasNext() {
			rawElem_19 := cur_18.Next()
			member := j2gort.CastPtr[JSTypeReconserializer_SeenTypeRecord](rawElem_19)
			_ = member
			sortedMembers.Add_1(j2gort.ValueOfInt(member.pointer))
		}
	}
	j2gort.Ext_java_util_Collections_Sort_1(sortedMembers)
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto_NewBuilder_0().SetUnion_1_Ext_com_google_javascript_jscomp_serialization_UnionTypeProto(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_UnionTypeProto](j2gort.Ext_com_google_javascript_jscomp_serialization_UnionTypeProto_NewBuilder_0().AddAllUnionMember(sortedMembers).Build())).Build())
}

// Java method: JSTypeReconserializer.reconcileObjectTypes(com.google.javascript.jscomp.serialization.JSTypeReconserializer.SeenTypeRecord seen)
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_reconcileObjectTypes_7789(seen *JSTypeReconserializer_SeenTypeRecord) *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto {
	_ = seen
	var instancePointers *j2gort.Ext_java_util_LinkedHashSet = this_.scratchInstancePointers
	_ = instancePointers
	var prototypePointers *j2gort.Ext_java_util_LinkedHashSet = this_.scratchPrototypePointers
	_ = prototypePointers
	var ownProperties *j2gort.Ext_java_util_LinkedHashSet = this_.scratchOwnProperties
	_ = ownProperties
	instancePointers.Clear()
	prototypePointers.Clear()
	ownProperties.Clear()
	var isClosureAssert bool = false
	_ = isClosureAssert
	var isConstructor bool = false
	_ = isConstructor
	var isInvalidating bool = false
	_ = isInvalidating
	var propertiesKeepOriginalName bool = false
	_ = propertiesKeepOriginalName
	{
		var j int32 = int32(0)
		_ = j
		var jstypeCount int32 = seen.JstypeCount()
		_ = jstypeCount
		for ; j < jstypeCount; j = j + 1 {
			var type_ pkg_rhino.JSType = seen.JstypeAt(j)
			_ = type_
			var objType pkg_rhino.ObjectType = j2gort.CheckNotNullIface(type_.GetJSTypeBase().ToMaybeObjectType())
			_ = objType
			if objType.GetObjectTypeBase().IsFunctionType() {
				var fnType *pkg_rhino.FunctionType = objType.ToMaybeFunctionType()
				_ = fnType
				if fnType.HasInstanceType() && !(j2gort.IfaceEquals[pkg_rhino.ObjectType](fnType.GetInstanceType(), nil)) {
					instancePointers.Add_1(j2gort.ValueOfInt(this_.SerializeType(pkg_rhino.Upcast_JSType(fnType.GetInstanceType()))))
					prototypePointers.Add_1(j2gort.ValueOfInt(this_.SerializeType(pkg_rhino.Upcast_JSType(fnType.GetPrototype()))))
					{
						ba_20 := isConstructor
						bb_21 := fnType.IsConstructor()
						isConstructor = ba_20 || bb_21
					}
				}
				{

					ba_22 := isClosureAssert
					bb_23 := JSTypeReconserializer_IsClosureAssert(fnType.GetClosurePrimitive())
					isClosureAssert = ba_22 || bb_23
				}

			}
			var propBound int32 = objType.GetObjectTypeBase().GetOwnPropertyNamesUpperBound()
			_ = propBound
			if propBound > int32(0) {
				var propBuf []string = this_.scratchOwnPropNames
				_ = propBuf
				if int32(len(propBuf)) < propBound {
					var nextCap int32 = int32(len(propBuf))
					_ = nextCap
					for nextCap < propBound {
						nextCap = (nextCap * int32(2))
					}
					propBuf = j2gort.NewStringArray(nextCap)
					this_.scratchOwnPropNames = propBuf
				}
				var propCount int32 = objType.GetObjectTypeBase().CollectOwnPropertyNames(propBuf)
				_ = propCount
				if (propCount > int32(1)) && !(objType.GetObjectTypeBase().IsRecordType()) {
					JSTypeReconserializer_SortStrings(propBuf, propCount)
				}
				{
					var i int32 = int32(0)
					_ = i
					for ; i < propCount; i = i + 1 {
						var ownProperty string = propBuf[i]
						_ = ownProperty
						{
							a_24 := propBuf
							i_25 := i
							a_24[i_25] = j2gort.NullString
						}
						if this_.priv_JSTypeReconserializer_shouldPropagateProperty_7780(ownProperty) {
							ownProperties.Add_1(j2gort.ValueOfInt(this_.stringPoolBuilder.Put(ownProperty)))
						}
					}
				}
			}
			{

				ba_26 := isInvalidating
				bb_27 := this_.invalidatingTypes.IsInvalidating(pkg_rhino.Upcast_JSType(objType))
				isInvalidating = ba_26 || bb_27
			}
			{

				ba_28 := propertiesKeepOriginalName
				bb_29 := objType.GetObjectTypeBase().IsEnumType()
				propertiesKeepOriginalName = ba_28 || bb_29
			}

		}
	}
	var protoBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto_NewBuilder_0()
	_ = protoBuilder
	if !(instancePointers.IsEmpty()) {
		protoBuilder.AddAllInstanceType(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](instancePointers))
	}
	if !(ownProperties.IsEmpty()) {
		protoBuilder.AddAllOwnProperty(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](ownProperties))
	}
	if !(prototypePointers.IsEmpty()) {
		protoBuilder.AddAllPrototype(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](prototypePointers))
	}
	var objectProto *j2gort.Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto](protoBuilder.SetClosureAssert(isClosureAssert).SetIsInvalidating(isInvalidating).SetMarkedConstructor(isConstructor).SetPropertiesKeepOriginalName(propertiesKeepOriginalName).SetUuid(seen.colorId.AsByteString()).Build())
	_ = objectProto
	instancePointers.Clear()
	prototypePointers.Clear()
	ownProperties.Clear()
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto_NewBuilder_0().SetObject_1_Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto(objectProto).Build())
}

// Java static method: JSTypeReconserializer.sortStrings(java.lang.String[] a, int n)
func JSTypeReconserializer_SortStrings(a []string, n int32) {
	clinit_JSTypeReconserializer()
	_ = a
	_ = n
	if n <= int32(16) {
		{
			var i int32 = int32(1)
			_ = i
			for ; i < n; i = i + 1 {
				var cur string = a[i]
				_ = cur
				var j int32 = (i - int32(1))
				_ = j
				for (j >= int32(0)) && (j2gort.JStrCompareTo(a[j], cur) > int32(0)) {
					{
						a_30 := a
						i_31 := (j + int32(1))
						a_30[i_31] = a[j]
					}
					j = j - 1
				}
				{
					a_32 := a
					i_33 := (j + int32(1))
					a_32[i_33] = cur
				}
			}
		}
		return
	}
	{
		var gap int32 = javalang.Ushr32(n, int(int32(1)))
		_ = gap
		for ; gap > int32(0); gap = javalang.Ushr32(gap, int(int32(1))) {
			{
				var i_2 int32 = gap
				_ = i_2
				for ; i_2 < n; i_2 = i_2 + 1 {
					var cur_2 string = a[i_2]
					_ = cur_2
					var j_2 int32 = i_2
					_ = j_2
					for (j_2 >= gap) && (j2gort.JStrCompareTo(a[(j_2-gap)], cur_2) > int32(0)) {
						{
							a_34 := a
							i_35 := j_2
							a_34[i_35] = a[(j_2 - gap)]
						}
						j_2 = (j_2 - gap)
					}
					{
						a_36 := a
						i_37 := j_2
						a_36[i_37] = cur_2
					}
				}
			}
		}
	}
}

// Java method: JSTypeReconserializer.ownAncestorInterfacesOf(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_ownAncestorInterfacesOf_7790(type_ pkg_rhino.ObjectType) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = type_
	var ctorType *pkg_rhino.FunctionType = type_.GetConstructor()
	_ = ctorType
	if ctorType == nil {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var extended j2gort.Ext_java_lang_Iterable = j2gort.AsIface[j2gort.Ext_java_lang_Iterable](ctorType.GetExtendedInterfaces())
	_ = extended
	var implemented j2gort.Ext_java_lang_Iterable = j2gort.AsIface[j2gort.Ext_java_lang_Iterable](ctorType.GetOwnImplementedInterfaces())
	_ = implemented
	if j2gort.Ext_com_google_common_collect_Iterables_IsEmpty(extended) && j2gort.Ext_com_google_common_collect_Iterables_IsEmpty(implemented) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var ancestors []int32 = []int32{}
	_ = ancestors
	{
		cur_38 := j2gort.NewIterCursor(j2gort.CheckNotNull(extended))
		for cur_38.HasNext() {
			rawElem_39 := cur_38.Next()
			ancestor := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_39)
			_ = ancestor
			ancestors = append(ancestors, this_.SerializeType(pkg_rhino.Upcast_JSType(ancestor)))
		}
	}
	{
		cur_40 := j2gort.NewIterCursor(j2gort.CheckNotNull(implemented))
		for cur_40.HasNext() {
			rawElem_41 := cur_40.Next()
			ancestor_2 := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_41)
			_ = ancestor_2
			ancestors = append(ancestors, this_.SerializeType(pkg_rhino.Upcast_JSType(ancestor_2)))
		}
	}
	return j2gort.ImmListFromTypedSlice(ancestors)
}

// Java method: JSTypeReconserializer.seedCachesWithAxiomaticTypes()
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_seedCachesWithAxiomaticTypes_7791() {
	if !(this_.seenTypeRecords.IsEmpty()) {
		j2gort.CheckState(false)
	}
	{
		cur_42 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_TypePointers_OFFSET_TO_AXIOMATIC_COLOR()))
		for cur_42.HasNext() {
			rawElem_43 := cur_42.Next()
			axiomatic := j2gort.CastPtr[pkg_colors.Color](rawElem_43)
			_ = axiomatic
			var index int32 = this_.seenTypeRecords.Size()
			_ = index
			var record *JSTypeReconserializer_SeenTypeRecord = NewJSTypeReconserializer_SeenTypeRecord(axiomatic.GetId(), index)
			_ = record
			this_.seenTypeRecords.Put(axiomatic.GetId(), record)
		}
	}
	StaticGet_JSTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP().ForEach_1(j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(lp_44 any, lp_45 any) {
		jstypeNative := j2gort.Cast[pkg_rhino.JSTypeNative](lp_44)
		_ = jstypeNative
		axiomatic_2 := j2gort.CastPtr[pkg_colors.Color](lp_45)
		_ = axiomatic_2
		this_.priv_JSTypeReconserializer_getOrCreateRecord_7786(axiomatic_2.GetId(), this_.registry.GetNativeType(jstypeNative))
	}))
	if !(this_.seenTypeRecords.Size() == StaticGet_TypePointers_OFFSET_TO_AXIOMATIC_COLOR().Size()) {
		j2gort.CheckState(false)
	}
}

// Java method: JSTypeReconserializer.checkValidLinearTime()
func (this_ *JSTypeReconserializer) priv_JSTypeReconserializer_checkValidLinearTime_7792() {
	if !(this_.serializationMode.RunValidation()) {
		return
	}
	var totalTypeCount int32 = this_.seenTypeRecords.Size()
	_ = totalTypeCount
	{
		cur_46 := this_.seenTypeRecords.IterValues()
		for cur_46.HasNext() {

			seen := cur_46.
				Next()
			_ = seen
			var offset int32 = seen.pointer
			_ = offset
			if !(offset >= int32(0)) {
				j2gort.CheckState(false)
			}
			if !(offset <= totalTypeCount) {
				j2gort.CheckStateFmt(false, "Found invalid pointer %s, out of a total of %s user-defined types", j2gort.ValueOfInt(offset), j2gort.ValueOfInt(totalTypeCount))
			}
		}
	}
}

// Java method: JSTypeReconserializer.generateTypePool()
func (this_ *JSTypeReconserializer) GenerateTypePool() *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool {
	if !(this_.state == StaticGet_JSTypeReconserializer_State_COLLECTING_TYPES()) {
		j2gort.CheckState(false)
	}
	this_.priv_JSTypeReconserializer_checkValidLinearTime_7792()
	var builder *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool_NewBuilder_0()
	_ = builder
	if this_.serializationMode.IncludeDebugInfo() {
		var debugInfo *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Builder = builder.GetDebugInfoBuilder()
		_ = debugInfo
		j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSetMultimap](this_.invalidatingTypes.GetMismatchLocations().Inverse()).AsMap()).ForEach_1(j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(lp_48 any, lp_49 any) {
			location := j2gort.CastNode[pkg_rhino.Node](lp_48)
			_ = location
			types := j2gort.AsIface[j2gort.Ext_java_util_Collection](lp_49)
			_ = types
			debugInfo.AddMismatchBuilder_0().SetSourceRef(location.GetLocation()).AddAllInvolvedColor(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](types.Stream().Peek(lambdaSingleton_js_type_reconserializer_5).Map(func() j2gort.Ext_java_util_function_Function {
				bound_50 := this_
				return j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
					return j2gort.ValueOfInt(bound_50.SerializeType(j2gort.AsIface[pkg_rhino.JSType](mrp0_)))
				})
			}()).Distinct().Sorted_0().Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList())))
		}))
	}
	this_.state = StaticGet_JSTypeReconserializer_State_GENERATING_POOL()
	{
		cur_52 := this_.seenTypeRecords.IterValues()
		for cur_52.HasNext() {

			seen := cur_52.
				Next()
			_ = seen
			if j2gort.MapContainsEqVal(pkg_colors.StaticGet_StandardColors_AXIOMATIC_COLORS(), seen.colorId) {
				if !(TypePointers_IsAxiomatic(seen.pointer)) {
					j2gort.CheckStateFmt(false, "Missing .type for SeenTypeRecord %s", j2gort.PtrToAny(seen))
				}
				continue
			}
			var iifeArg_1 *j2gort.
				Ext_com_google_javascript_jscomp_serialization_TypeProto
		iifeSwitch_1:
			switch {
			default:

				if seen.unionMembers != nil {
					iifeArg_1 = this_.priv_JSTypeReconserializer_reconcileUnionTypes_7788(seen)
					break iifeSwitch_1
				}
				iifeArg_1 = this_.priv_JSTypeReconserializer_reconcileObjectTypes_7789(seen)
				break iifeSwitch_1
			}
			builder.AddType_1_Ext_com_google_javascript_jscomp_serialization_TypeProto(iifeArg_1)
		}
	}
	{
		cur_54 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.disambiguateEdges.KeySet())))
		for cur_54.HasNext() {
			rawElem_55 := cur_54.Next()
			subtype := j2gort.CastPtr[j2gort.BoxInt](rawElem_55)
			_ = subtype
			{
				cur_56 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.disambiguateEdges.Get(j2gort.PtrToAny(subtype)))))
				for cur_56.HasNext() {
					rawElem_57 := cur_56.Next()
					supertype := j2gort.CastPtr[j2gort.BoxInt](rawElem_57)
					_ = supertype
					builder.AddDisambiguationEdges_1_Ext_com_google_javascript_jscomp_serialization_SubtypingEdge_Builder(j2gort.Ext_com_google_javascript_jscomp_serialization_SubtypingEdge_NewBuilder_0().SetSubtype(j2gort.UnboxInt(subtype)).SetSupertype(j2gort.UnboxInt(supertype)))
				}
			}
		}
	}
	this_.state = StaticGet_JSTypeReconserializer_State_FINISHED()
	this_.priv_JSTypeReconserializer_checkValidLinearTime_7792()
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool](builder.Build())
}

// Java method: JSTypeReconserializer.getColorIdToJSTypeMapForDebugging()
func (this_ *JSTypeReconserializer) GetColorIdToJSTypeMapForDebugging() j2gort.Ext_com_google_common_collect_ImmutableMultimap {
	var colorIdToTypes j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder](j2gort.Ext_com_google_common_collect_ImmutableMultimap_S_Builder().OrderKeysBy(j2gort.Ext_java_util_Comparator_NaturalOrder()))
	_ = colorIdToTypes
	{
		cur_58 := this_.seenTypeRecords.IterValues()
		for cur_58.HasNext() {

			record := cur_58.
				Next()
			_ = record
			var jstypeCount int32 = record.JstypeCount()
			_ = jstypeCount
			if jstypeCount > int32(0) {
				var key string = record.colorId.String()
				_ = key
				{
					var i int32 = int32(0)
					_ = i
					for ; i < jstypeCount; i = i + 1 {
						j2gort.MultimapBuilderPutStr(colorIdToTypes, key, record.JstypeAt(i))
					}
				}
			}
		}
	}
	if this_.fallbackTypeToRecordCache != nil {
		this_.fallbackTypeToRecordCache.ForEach(j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(lp_60 any, lp_61 any) {
			jstype := j2gort.AsIface[pkg_rhino.JSType](lp_60)
			_ = jstype
			record_2 := j2gort.CastPtr[JSTypeReconserializer_SeenTypeRecord](lp_61)
			_ = record_2
			if !(record_2.ContainsJstype(jstype)) {
				j2gort.MultimapBuilderPutStr(colorIdToTypes, record_2.colorId.String(), jstype)
			}
		}))
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap](colorIdToTypes.Build())
}

// Java static method: JSTypeReconserializer.isClosureAssert(com.google.javascript.rhino.@org.jspecify.annotations.Nullable ClosurePrimitive primitive)
func JSTypeReconserializer_IsClosureAssert(primitive pkg_rhino.ClosurePrimitive) bool {
	clinit_JSTypeReconserializer()
	_ = primitive
	if primitive == 0 {
		return false
	}
	switch (primitive).Ordinal() {
	case int32(2) /* ClosurePrimitive.ASSERTS_TRUTHY */, int32(1) /* ClosurePrimitive.ASSERTS_MATCHES_RETURN */ :
		return true
	case int32(0) /* ClosurePrimitive.ASSERTS_FAIL */ :
		return false
	}
	panic("unreachable switch expression")
}

// Java enum: com.google.javascript.jscomp.serialization.JSTypeReconserializer.State
type JSTypeReconserializer_State uint8

const (
	enum_JSTypeReconserializer_State_COLLECTING_TYPES JSTypeReconserializer_State = 1
	enum_JSTypeReconserializer_State_GENERATING_POOL  JSTypeReconserializer_State = 2
	enum_JSTypeReconserializer_State_FINISHED         JSTypeReconserializer_State = 3
)

var enumBaseTable_JSTypeReconserializer_State = [...]j2gort.EnumBase{
	{Name_: "COLLECTING_TYPES", Ordinal_: 0},
	{Name_: "GENERATING_POOL", Ordinal_: 1},
	{Name_: "FINISHED", Ordinal_: 2},
}
var enumValues_JSTypeReconserializer_State = []JSTypeReconserializer_State{1, 2, 3}

func StaticGet_JSTypeReconserializer_State_COLLECTING_TYPES() JSTypeReconserializer_State {
	return enum_JSTypeReconserializer_State_COLLECTING_TYPES
}
func StaticGet_JSTypeReconserializer_State_GENERATING_POOL() JSTypeReconserializer_State {
	return enum_JSTypeReconserializer_State_GENERATING_POOL
}
func StaticGet_JSTypeReconserializer_State_FINISHED() JSTypeReconserializer_State {
	return enum_JSTypeReconserializer_State_FINISHED
}

func JSTypeReconserializer_State_Values() []JSTypeReconserializer_State {
	return j2gort.CloneSlice(enumValues_JSTypeReconserializer_State)
}

func JSTypeReconserializer_State_ValueOf(name string) JSTypeReconserializer_State {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSTypeReconserializer_State {
		if enumBaseTable_JSTypeReconserializer_State[i_].Name_ == name {
			return JSTypeReconserializer_State(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.serialization.JSTypeReconserializer.State." + name))
}

func (this_ JSTypeReconserializer_State) IsJavaNil_() bool                              { return this_ == 0 }
func (JSTypeReconserializer_State) HasIdentityEquals_()                                 {}
func (JSTypeReconserializer_State) Implements_j2gort_Ext_java_lang_Enum()               {}
func (JSTypeReconserializer_State) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSTypeReconserializer_State) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (JSTypeReconserializer_State) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.JSTypeReconserializer$State", "State"
}
func (this_ JSTypeReconserializer_State) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSTypeReconserializer_State[this_-1]
}
func (this_ JSTypeReconserializer_State) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSTypeReconserializer_State) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSTypeReconserializer_State) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSTypeReconserializer_State[this_-1].Name_
}
func (this_ JSTypeReconserializer_State) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSTypeReconserializer_State[this_-1].Name_
}
func (this_ JSTypeReconserializer_State) Equals(other_ any) bool {
	o_, ok_ := other_.(JSTypeReconserializer_State)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSTypeReconserializer_State) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSTypeReconserializer_State[this_-1].HashCode()
}
func (this_ JSTypeReconserializer_State) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSTypeReconserializer_State[this_-1].CompareTo(other_)
}
func (this_ JSTypeReconserializer_State) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSTypeReconserializer_State[this_-1].Name_)
}
func init() {
	vals_ := []any{JSTypeReconserializer_State(1), JSTypeReconserializer_State(2), JSTypeReconserializer_State(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.serialization.JSTypeReconserializer.State", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.serialization.JSTypeReconserializer$State", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.JSTypeReconserializer.State", "com.google.javascript.jscomp.serialization.JSTypeReconserializer$State", func(v_ any) bool {
		return j2gort.InstanceOf[JSTypeReconserializer_State](v_)
	})
}

// Java class: com.google.javascript.jscomp.serialization.JSTypeReconserializer.SeenTypeRecord
type JSTypeReconserializer_SeenTypeRecord struct {
	firstJstype      pkg_rhino.JSType                                  `gson:"firstJstype"`
	extraJstypes     []pkg_rhino.JSType                                `gson:"extraJstypes"`
	unionMembers     j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"unionMembers"`
	colorId          pkg_colors.ColorId                                `gson:"colorId"`
	pointer          int32                                             `gson:"pointer"`
	extraJstypeCount int32                                             `gson:"extraJstypeCount"`
	idHash_          atomic.Int32
}

func (this_ *JSTypeReconserializer_SeenTypeRecord) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeReconserializer_SeenTypeRecord) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeReconserializer_SeenTypeRecord) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.JSTypeReconserializer$SeenTypeRecord", "SeenTypeRecord"
}
func (this_ *JSTypeReconserializer_SeenTypeRecord) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeReconserializer_SeenTypeRecord) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.JSTypeReconserializer.SeenTypeRecord", "com.google.javascript.jscomp.serialization.JSTypeReconserializer$SeenTypeRecord", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeReconserializer_SeenTypeRecord](v_)
	})
}

func (this_ *JSTypeReconserializer_SeenTypeRecord) initFields_JSTypeReconserializer_SeenTypeRecord() {
	_ = this_

}

// Java constructor: SeenTypeRecord(com.google.javascript.jscomp.colors.ColorId colorId, int pointer)
func NewJSTypeReconserializer_SeenTypeRecord(colorId pkg_colors.ColorId, pointer int32) *JSTypeReconserializer_SeenTypeRecord {
	this_ := &JSTypeReconserializer_SeenTypeRecord{}
	InitJSTypeReconserializer_SeenTypeRecord(this_, colorId, pointer)
	return this_
}

func InitJSTypeReconserializer_SeenTypeRecord(this_ *JSTypeReconserializer_SeenTypeRecord, colorId pkg_colors.ColorId, pointer int32) {
	_ = colorId
	_ = pointer
	this_.initFields_JSTypeReconserializer_SeenTypeRecord()
	this_.colorId = colorId
	this_.pointer = pointer
}

// Java method: SeenTypeRecord.addJstype(com.google.javascript.rhino.jstype.JSType jstype)
func (this_ *JSTypeReconserializer_SeenTypeRecord) AddJstype(jstype pkg_rhino.JSType) {
	_ = jstype
	if this_.firstJstype == nil || this_.firstJstype.IsJavaNil_() {
		this_.firstJstype = jstype
		return
	}
	var extra []pkg_rhino.JSType = this_.extraJstypes
	_ = extra
	var count int32 = this_.extraJstypeCount
	_ = count
	if extra == nil {
		extra = make([]pkg_rhino.JSType, int32(4))
		this_.extraJstypes = extra
	} else if count == int32(len(extra)) {
		{

			arrSrc_62 := extra
			arrLen_63 := (count * int32(2))
			j2gort.CheckNotNullSlice(arrSrc_62)
			if arrLen_63 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_64 := make([]pkg_rhino.JSType, arrLen_63)
			copy(arrOut_64, arrSrc_62)
			extra = arrOut_64
		}

		this_.extraJstypes = extra
	}
	{
		a_65 := extra
		i_66 := count
		a_65[i_66] = jstype
	}
	this_.extraJstypeCount = (count + int32(1))
}

// Java method: SeenTypeRecord.jstypeCount()
func (this_ *JSTypeReconserializer_SeenTypeRecord) JstypeCount() int32 {
	if this_.firstJstype == nil || this_.firstJstype.IsJavaNil_() {
		return int32(0)
	} else {
		return (int32(1) + this_.extraJstypeCount)
	}
}

// Java method: SeenTypeRecord.jstypeAt(int index)
func (this_ *JSTypeReconserializer_SeenTypeRecord) JstypeAt(index int32) pkg_rhino.JSType {
	_ = index
	if index == int32(0) {
		return this_.firstJstype
	} else {
		return this_.extraJstypes[(index - int32(1))]
	}
}

// Java method: SeenTypeRecord.containsJstype(com.google.javascript.rhino.jstype.JSType target)
func (this_ *JSTypeReconserializer_SeenTypeRecord) ContainsJstype(target pkg_rhino.JSType) bool {
	_ = target
	if this_.firstJstype == nil || this_.firstJstype.IsJavaNil_() {
		return false
	}
	if this_.firstJstype.GetJSTypeBase().EqualsJSType(target) {
		return true
	}
	var extra []pkg_rhino.JSType = this_.extraJstypes
	_ = extra
	var count int32 = this_.extraJstypeCount
	_ = count
	{
		var i int32 = int32(0)
		_ = i
		for ; i < count; i = i + 1 {
			if extra[i].GetJSTypeBase().EqualsJSType(target) {
				return true
			}
		}
	}
	return false
}
