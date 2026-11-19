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
// Port of com/google/javascript/jscomp/serialization/JSTypeColorIdHasher.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_js_type_color_id_hasher_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_js_type_color_id_hasher_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](mrp0_).GetValue() })
}

// Java class: com.google.javascript.jscomp.serialization.JSTypeColorIdHasher
type JSTypeColorIdHasher struct {
	standardColorObjectIds j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"standardColorObjectIds"`
	idHash_                atomic.Int32
}

func (this_ *JSTypeColorIdHasher) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeColorIdHasher) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeColorIdHasher) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.JSTypeColorIdHasher", "JSTypeColorIdHasher"
}
func (this_ *JSTypeColorIdHasher) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeColorIdHasher) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.JSTypeColorIdHasher", "com.google.javascript.jscomp.serialization.JSTypeColorIdHasher", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeColorIdHasher](v_)
	})
}

var static_JSTypeColorIdHasher_FARM_64 j2gort.Ext_com_google_common_hash_HashFunction

func StaticGet_JSTypeColorIdHasher_FARM_64() j2gort.Ext_com_google_common_hash_HashFunction {
	clinit_JSTypeColorIdHasher()
	return static_JSTypeColorIdHasher_FARM_64
}
func StaticSet_JSTypeColorIdHasher_FARM_64(v j2gort.Ext_com_google_common_hash_HashFunction) {
	clinit_JSTypeColorIdHasher()
	static_JSTypeColorIdHasher_FARM_64 = v
}

var static_JSTypeColorIdHasher_NATIVE_TYPE_TO_ID j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_JSTypeColorIdHasher_NATIVE_TYPE_TO_ID() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_JSTypeColorIdHasher()
	return static_JSTypeColorIdHasher_NATIVE_TYPE_TO_ID
}
func StaticSet_JSTypeColorIdHasher_NATIVE_TYPE_TO_ID(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_JSTypeColorIdHasher()
	static_JSTypeColorIdHasher_NATIVE_TYPE_TO_ID = v
}

var clinitDone_JSTypeColorIdHasher uint32
var clinitState_JSTypeColorIdHasher j2gort.ClinitState

func clinit_JSTypeColorIdHasher() {
	if atomic.LoadUint32(&clinitDone_JSTypeColorIdHasher) != 0 {
		return
	}
	clinitSlow_JSTypeColorIdHasher()
}

func clinitSlow_JSTypeColorIdHasher() {
	if !j2gort.EnterClinit(&clinitState_JSTypeColorIdHasher, &clinitDone_JSTypeColorIdHasher) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JSTypeColorIdHasher, &clinitDone_JSTypeColorIdHasher)
	static_JSTypeColorIdHasher_FARM_64 = j2gort.Ext_com_google_common_hash_Hashing_FarmHashFingerprint64()
	static_JSTypeColorIdHasher_NATIVE_TYPE_TO_ID = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder().Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_ARGUMENTS_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_ARGUMENTS_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_ARRAY_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_READONLY_ARRAY_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_READONLY_ARRAY_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_ASYNC_ITERATOR_ITERABLE_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_GENERATOR_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_GENERATOR_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_I_TEMPLATE_ARRAY_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_I_TEMPLATE_ARRAY_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_ITERATOR_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_ITERATOR_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_PROMISE_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_BIGINT_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_BIGINT_OBJECT_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_BOOLEAN_OBJECT_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_NUMBER_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_NUMBER_OBJECT_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_STRING_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_STRING_OBJECT_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_SYMBOL_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_SYMBOL_OBJECT_ID()))).BuildOrThrow())
}

func (this_ *JSTypeColorIdHasher) initFields_JSTypeColorIdHasher() {
	_ = this_
}

// Java constructor: JSTypeColorIdHasher(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewJSTypeColorIdHasher(registry *pkg_rhino.JSTypeRegistry) *JSTypeColorIdHasher {
	clinit_JSTypeColorIdHasher()
	this_ := &JSTypeColorIdHasher{}
	InitJSTypeColorIdHasher(this_, registry)
	return this_
}

func InitJSTypeColorIdHasher(this_ *JSTypeColorIdHasher, registry *pkg_rhino.JSTypeRegistry) {
	_ = registry
	this_.initFields_JSTypeColorIdHasher()
	this_.standardColorObjectIds = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](StaticGet_JSTypeColorIdHasher_NATIVE_TYPE_TO_ID().EntrySet()).Stream().Collect_1(j2gort.Ext_com_google_common_collect_ImmutableMap_ToImmutableMap_2(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
		e := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_1)
		_ = e
		return registry.GetNativeObjectType(j2gort.Cast[pkg_rhino.JSTypeNative](e.GetKey()))
	}), lambdaSingleton_js_type_color_id_hasher_1)))
}

// Java method: JSTypeColorIdHasher.hashObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *JSTypeColorIdHasher) HashObjectType(type_ pkg_rhino.ObjectType) pkg_colors.ColorId {
	_ = type_
	if !((((((type_ != nil && !type_.IsJavaNil_()) && !(type_.IsEnumElementType())) && !(type_.IsNoResolvedType())) && !(type_.IsTemplatizedType())) && !(type_.IsUnknownType())) && !(type_.IsTemplateType())) {
		j2gort.CheckStateMsg(false, type_)
	}
	var boxId pkg_colors.ColorId = j2gort.Cast[pkg_colors.ColorId](this_.standardColorObjectIds.Get(type_))
	_ = boxId
	if boxId != (pkg_colors.ColorId{}) {
		return boxId
	}
	var hasher j2gort.Ext_com_google_common_hash_Hasher = StaticGet_JSTypeColorIdHasher_FARM_64().NewHasher_0()
	_ = hasher
	if type_.IsFunctionType() {
		var fnType *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
		_ = fnType
		if fnType.HasInstanceType() {
			hasher.PutInt(int32(-1937018661))
		}
	}
	if type_.HasReferenceName() {
		this_.priv_JSTypeColorIdHasher_putReferenceNameType_7692(hasher, type_)
	} else {
		var ownProps j2gort.Ext_java_util_Set = type_.GetOwnPropertyNames()
		_ = ownProps
		if (ownProps.Size() <= int32(1)) || type_.IsRecordType() {
			{
				cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(ownProps))
				for cur_2.HasNext() {
					prop := cur_2.NextStr()
					_ = prop
					JSTypeColorIdHasher_PutUtf8(hasher, prop)
				}
			}
		} else {
			var sortedProps []string = j2gort.CopyStrColToSlice(j2gort.AsIface[j2gort.Ext_java_util_Collection](ownProps))
			_ = sortedProps
			j2gort.SliceSortNatural(sortedProps)
			{
				var i int32 = int32(0)
				_ = i
				for ; i < int32(len(sortedProps)); i = i + 1 {
					JSTypeColorIdHasher_PutUtf8(hasher, sortedProps[i])
				}
			}
		}
	}
	return pkg_colors.ColorId_FromUnsigned_1_int64(hasher.Hash().AsLong())
}

// Java method: JSTypeColorIdHasher.putReferenceNameType(com.google.common.hash.Hasher hasher, com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *JSTypeColorIdHasher) priv_JSTypeColorIdHasher_putReferenceNameType_7692(hasher j2gort.Ext_com_google_common_hash_Hasher, type_ pkg_rhino.ObjectType) {
	_ = hasher
	_ = type_
	JSTypeColorIdHasher_PutUtf8(hasher, type_.GetReferenceName())
	var sourceRef pkg_rhino.JSType_WithSourceRef = JSTypeColorIdHasher_SourceRefFor(type_)
	_ = sourceRef
	if sourceRef == nil || sourceRef.IsJavaNil_() {
		hasher.PutInt(int32(1712056194))
		return
	}
	var moduleId string = sourceRef.GetGoogModuleId()
	_ = moduleId
	if moduleId == j2gort.NullString {
		hasher.PutInt(int32(630441471))
	} else {
		JSTypeColorIdHasher_PutUtf8(hasher, moduleId)
	}
}

// Java static method: JSTypeColorIdHasher.sourceRefFor(com.google.javascript.rhino.jstype.ObjectType type)
func JSTypeColorIdHasher_SourceRefFor(type_ pkg_rhino.ObjectType) pkg_rhino.JSType_WithSourceRef {
	clinit_JSTypeColorIdHasher()
	_ = type_
	if type_.IsEnumType() {
		return pkg_rhino.ToIface_JSType_WithSourceRef(type_.ToMaybeEnumType())
	} else if type_.IsEnumElementType() {
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", type_)))
	}
	if type_.IsFunctionType() {
		return pkg_rhino.ToIface_JSType_WithSourceRef(type_.ToMaybeFunctionType())
	} else if type_.IsFunctionPrototypeType() {
		return pkg_rhino.ToIface_JSType_WithSourceRef(type_.GetOwnerFunction())
	} else if type_.GetConstructor() != nil {
		return pkg_rhino.ToIface_JSType_WithSourceRef(type_.GetConstructor())
	}
	return nil
}

// Java static method: JSTypeColorIdHasher.putUtf8(com.google.common.hash.Hasher hasher, java.lang.String s)
func JSTypeColorIdHasher_PutUtf8(hasher j2gort.Ext_com_google_common_hash_Hasher, s string) {
	clinit_JSTypeColorIdHasher()
	_ = hasher
	_ = s
	hasher.PutString(j2gort.BoxString(s), j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8)
}

// Java class: com.google.javascript.jscomp.serialization.JSTypeColorIdHasher.Marker
type JSTypeColorIdHasher_Marker struct {
	idHash_ atomic.Int32
}

func (this_ *JSTypeColorIdHasher_Marker) BindSelf_(s any) {
	_ = s
}
func (this_ *JSTypeColorIdHasher_Marker) IsJavaNil_() bool { return this_ == nil }
func (*JSTypeColorIdHasher_Marker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.JSTypeColorIdHasher$Marker", "Marker"
}
func (this_ *JSTypeColorIdHasher_Marker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSTypeColorIdHasher_Marker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.JSTypeColorIdHasher.Marker", "com.google.javascript.jscomp.serialization.JSTypeColorIdHasher$Marker", func(v_ any) bool {
		return j2gort.InstanceOf[*JSTypeColorIdHasher_Marker](v_)
	})
}

func StaticGet_JSTypeColorIdHasher_Marker_HAS_INSTANCE_TYPE() int32 { return int32(-1937018661) }
func StaticGet_JSTypeColorIdHasher_Marker_NO_GOOG_MODULE_ID() int32 { return int32(630441471) }
func StaticGet_JSTypeColorIdHasher_Marker_UNKNOWN_SOURCEREF() int32 { return int32(1712056194) }
func (this_ *JSTypeColorIdHasher_Marker) initFields_JSTypeColorIdHasher_Marker() {
	_ = this_
}

// Java constructor: Marker()
func NewJSTypeColorIdHasher_Marker() *JSTypeColorIdHasher_Marker {
	this_ := &JSTypeColorIdHasher_Marker{}
	InitJSTypeColorIdHasher_Marker(this_)
	return this_
}

func InitJSTypeColorIdHasher_Marker(this_ *JSTypeColorIdHasher_Marker) {
	this_.initFields_JSTypeColorIdHasher_Marker()
}
