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

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.serialization.JSTypeColorIdHasher
type JSTypeColorIdHasher struct {
	standardColorObjectIds *pkg_base.LinkedIdentityHashMap[pkg_rhino.ObjectType, pkg_colors.ColorId] `gson:"standardColorObjectIds"`
	scratchBytes           []int8                                                                    `gson:"scratchBytes"`
	scratchLen             int32                                                                     `gson:"scratchLen"`
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

func StaticGet_JSTypeColorIdHasher_FARM_K0() int64 { return int64(-4348849565147123417) }
func StaticGet_JSTypeColorIdHasher_FARM_K1() int64 { return int64(-5435081209227447693) }
func StaticGet_JSTypeColorIdHasher_FARM_K2() int64 { return int64(-7286425919675154353) }

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
	static_JSTypeColorIdHasher_NATIVE_TYPE_TO_ID = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder().Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_ARGUMENTS_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_ARGUMENTS_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_ARRAY_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_READONLY_ARRAY_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_READONLY_ARRAY_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_ASYNC_ITERATOR_ITERABLE_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_GENERATOR_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_GENERATOR_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_I_TEMPLATE_ARRAY_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_I_TEMPLATE_ARRAY_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_ITERATOR_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_ITERATOR_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_PROMISE_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_BIGINT_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_BIGINT_OBJECT_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_BOOLEAN_OBJECT_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_NUMBER_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_NUMBER_OBJECT_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_STRING_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_STRING_OBJECT_ID()))).Put_2(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_SYMBOL_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_SYMBOL_OBJECT_ID()))).BuildOrThrow())
}

func (this_ *JSTypeColorIdHasher) initFields_JSTypeColorIdHasher() {
	_ = this_
	this_.scratchBytes = make([]int8, int32(256))

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
	var ids *pkg_base.LinkedIdentityHashMap[pkg_rhino.ObjectType, pkg_colors.ColorId] = pkg_base.NewLinkedIdentityHashMap[pkg_rhino.ObjectType, pkg_colors.ColorId]()
	_ = ids
	{
		cur_1 := j2gort.NewMapEntryCursor(StaticGet_JSTypeColorIdHasher_NATIVE_TYPE_TO_ID())
		for cur_1.NextEntry() {
			ids.Put(registry.GetNativeObjectType(j2gort.Cast[pkg_rhino.JSTypeNative](cur_1.GetKey())), j2gort.Cast[pkg_colors.ColorId](cur_1.GetValue()))
		}
	}
	this_.standardColorObjectIds = ids
}

// Java method: JSTypeColorIdHasher.ensureCapacity(int additional)
func (this_ *JSTypeColorIdHasher) priv_JSTypeColorIdHasher_ensureCapacity_7774(additional int32) {
	_ = additional
	var needed int32 = (this_.scratchLen + additional)
	_ = needed
	var buf []int8 = this_.scratchBytes
	_ = buf
	if needed > int32(len(buf)) {
		var newCap int32 = (int32(len(buf)) * int32(2))
		_ = newCap
		for newCap < needed {
			newCap = (newCap * int32(2))
		}
		this_.scratchBytes = func() []int8 {
			arrSrc_2 := buf
			arrLen_3 := newCap
			j2gort.CheckNotNullSlice(arrSrc_2)
			if arrLen_3 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_4 := make([]int8, arrLen_3)
			copy(arrOut_4, arrSrc_2)
			return arrOut_4
		}()
	}
}

// Java method: JSTypeColorIdHasher.putInt(int v)
func (this_ *JSTypeColorIdHasher) priv_JSTypeColorIdHasher_putInt_7775(v int32) {
	_ = v
	this_.priv_JSTypeColorIdHasher_ensureCapacity_7774(int32(4))
	var buf []int8 = this_.scratchBytes
	_ = buf
	var pos int32 = this_.scratchLen
	_ = pos
	{
		a_5 := buf
		i_6 := pos
		a_5[i_6] = int8(v)
	}
	{
		a_7 := buf
		i_8 := (pos + int32(1))
		a_7[i_8] = int8((javalang.Ushr32(v, int(int32(8)))))
	}
	{
		a_9 := buf
		i_10 := (pos + int32(2))
		a_9[i_10] = int8((javalang.Ushr32(v, int(int32(16)))))
	}
	{
		a_11 := buf
		i_12 := (pos + int32(3))
		a_11[i_12] = int8((javalang.Ushr32(v, int(int32(24)))))
	}
	this_.scratchLen = (pos + int32(4))
}

// Java method: JSTypeColorIdHasher.putUtf8(java.lang.String s)
func (this_ *JSTypeColorIdHasher) priv_JSTypeColorIdHasher_putUtf8_7776(s string) {
	_ = s
	var sLen int32 = j2gort.JStrLen(s)
	_ = sLen
	if sLen == int32(0) {
		return
	}
	this_.priv_JSTypeColorIdHasher_ensureCapacity_7774((sLen * int32(3)))
	var buf []int8 = this_.scratchBytes
	_ = buf
	var pos int32 = this_.scratchLen
	_ = pos
	{
		var i int32 = int32(0)
		_ = i
		for ; i < sLen; i = i + 1 {
			var c uint16 = j2gort.JStrCharAt(s, i)
			_ = c
			if int32(c) < int32(128) {
				{
					a_13 := buf
					var i_14 int32
					{

						old_15 := pos
						nw_16 := old_15 + 1
						pos = nw_16
						i_14 = old_15
					}
					a_13[i_14] = int8(c)
				}
			} else if int32(c) < int32(2048) {
				{
					a_17 := buf
					var i_18 int32
					{

						old_19 := pos
						nw_20 := old_19 + 1
						pos = nw_20
						i_18 = old_19
					}
					a_17[i_18] = int8((int32(192) | (javalang.Ushr32(int32(c), int(int32(6))))))
				}
				{
					a_21 := buf
					var i_22 int32
					{

						old_23 := pos
						nw_24 := old_23 + 1
						pos = nw_24
						i_22 = old_23
					}
					a_21[i_22] = int8((int32(128) | (int32(c) & int32(63))))
				}
			} else if j2gort.BoxChar_IsSurrogate(c) {
				if j2gort.BoxChar_IsHighSurrogate(c) && ((i + int32(1)) < sLen) {
					var low uint16 = j2gort.JStrCharAt(s, (i + int32(1)))
					_ = low
					if j2gort.BoxChar_IsLowSurrogate(low) {
						var cp int32 = j2gort.BoxChar_ToCodePoint(c, low)
						_ = cp
						i = i + 1
						{
							a_25 := buf
							var i_26 int32
							{

								old_27 := pos
								nw_28 := old_27 + 1
								pos = nw_28
								i_26 = old_27
							}
							a_25[i_26] = int8((int32(240) | (javalang.Ushr32(cp, int(int32(18))))))
						}
						{
							a_29 := buf
							var i_30 int32
							{

								old_31 := pos
								nw_32 := old_31 + 1
								pos = nw_32
								i_30 = old_31
							}
							a_29[i_30] = int8((int32(128) | ((javalang.Ushr32(cp, int(int32(12)))) & int32(63))))
						}
						{
							a_33 := buf
							var i_34 int32
							{

								old_35 := pos
								nw_36 := old_35 + 1
								pos = nw_36
								i_34 = old_35
							}
							a_33[i_34] = int8((int32(128) | ((javalang.Ushr32(cp, int(int32(6)))) & int32(63))))
						}
						{
							a_37 := buf
							var i_38 int32
							{

								old_39 := pos
								nw_40 := old_39 + 1
								pos = nw_40
								i_38 = old_39
							}
							a_37[i_38] = int8((int32(128) | (cp & int32(63))))
						}
						continue
					}
				}
				{
					a_41 := buf
					var i_42 int32
					{

						old_43 := pos
						nw_44 := old_43 + 1
						pos = nw_44
						i_42 = old_43
					}
					a_41[i_42] = int8(63)
				}
			} else {
				{
					a_45 := buf
					var i_46 int32
					{

						old_47 := pos
						nw_48 := old_47 + 1
						pos = nw_48
						i_46 = old_47
					}
					a_45[i_46] = int8((int32(224) | (javalang.Ushr32(int32(c), int(int32(12))))))
				}
				{
					a_49 := buf
					var i_50 int32
					{

						old_51 := pos
						nw_52 := old_51 + 1
						pos = nw_52
						i_50 = old_51
					}
					a_49[i_50] = int8((int32(128) | ((javalang.Ushr32(int32(c), int(int32(6)))) & int32(63))))
				}
				{
					a_53 := buf
					var i_54 int32
					{

						old_55 := pos
						nw_56 := old_55 + 1
						pos = nw_56
						i_54 = old_55
					}
					a_53[i_54] = int8((int32(128) | (int32(c) & int32(63))))
				}
			}
		}
	}
	this_.scratchLen = pos
}

// Java method: JSTypeColorIdHasher.hashObjectType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *JSTypeColorIdHasher) HashObjectType(type_ pkg_rhino.ObjectType) pkg_colors.ColorId {
	_ = type_
	if !((((((type_ != nil && !type_.IsJavaNil_()) && !(type_.GetObjectTypeBase().IsEnumElementType())) && !(type_.GetObjectTypeBase().IsNoResolvedType())) && !(type_.GetObjectTypeBase().IsTemplatizedType())) && !(type_.GetObjectTypeBase().IsUnknownType())) && !(type_.GetObjectTypeBase().IsTemplateType())) {
		j2gort.CheckStateMsg(false, type_)
	}
	var boxId pkg_colors.ColorId = this_.standardColorObjectIds.Get(type_)
	_ = boxId
	if !(boxId).IsJavaNil_() {
		return boxId
	}
	this_.scratchLen = int32(0)
	base_type__ := type_.GetObjectTypeBase()
	if base_type__.IsFunctionType() {
		var fnType *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
		_ = fnType
		if fnType.HasInstanceType() {
			this_.priv_JSTypeColorIdHasher_putInt_7775(int32(-1937018661))
		}
	}
	if base_type__.HasReferenceName() {
		this_.priv_JSTypeColorIdHasher_putReferenceNameType_7778(type_)
	} else {
		var ownProps j2gort.Ext_java_util_Set = type_.GetOwnPropertyNames()
		_ = ownProps
		var propCount int32 = ownProps.Size()
		_ = propCount
		if (propCount <= int32(1)) || base_type__.IsRecordType() {
			{
				cur_57 := j2gort.NewIterCursor(j2gort.CheckNotNull(ownProps))
				for cur_57.HasNext() {
					prop := cur_57.NextStr()
					_ = prop
					this_.priv_JSTypeColorIdHasher_putUtf8_7776(prop)
				}
			}
		} else {
			var sortedProps []string = j2gort.ColToStringArray(ownProps, j2gort.NewStringArray(propCount))
			_ = sortedProps
			j2gort.ArraysSortStrings(sortedProps)
			{
				var i int32 = int32(0)
				_ = i
				for ; i < int32(len(sortedProps)); i = i + 1 {
					this_.priv_JSTypeColorIdHasher_putUtf8_7776(sortedProps[i])
				}
			}
		}
	}
	return pkg_colors.ColorId_FromUnsigned_1_int64(JSTypeColorIdHasher_FarmHash64(this_.scratchBytes, this_.scratchLen))
}

// Java method: JSTypeColorIdHasher.putReferenceNameType(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *JSTypeColorIdHasher) priv_JSTypeColorIdHasher_putReferenceNameType_7778(type_ pkg_rhino.ObjectType) {
	_ = type_
	this_.priv_JSTypeColorIdHasher_putUtf8_7776(type_.GetReferenceName())
	var sourceRef pkg_rhino.JSType_WithSourceRef = JSTypeColorIdHasher_SourceRefFor(type_)
	_ = sourceRef
	if sourceRef == nil || sourceRef.IsJavaNil_() {
		this_.priv_JSTypeColorIdHasher_putInt_7775(int32(1712056194))
		return
	}
	var moduleId string = sourceRef.GetGoogModuleId()
	_ = moduleId
	if moduleId == j2gort.NullString {
		this_.priv_JSTypeColorIdHasher_putInt_7775(int32(630441471))
	} else {
		this_.priv_JSTypeColorIdHasher_putUtf8_7776(moduleId)
	}
}

// Java static method: JSTypeColorIdHasher.sourceRefFor(com.google.javascript.rhino.jstype.ObjectType type)
func JSTypeColorIdHasher_SourceRefFor(type_ pkg_rhino.ObjectType) pkg_rhino.JSType_WithSourceRef {
	clinit_JSTypeColorIdHasher()
	_ = type_
	base_type__ := type_.GetObjectTypeBase()
	if base_type__.IsEnumType() {
		return pkg_rhino.ToIface_JSType_WithSourceRef(type_.ToMaybeEnumType())
	} else if base_type__.IsEnumElementType() {
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", type_)))
	}
	if base_type__.IsFunctionType() {
		return pkg_rhino.ToIface_JSType_WithSourceRef(type_.ToMaybeFunctionType())
	} else if base_type__.IsFunctionPrototypeType() {
		return pkg_rhino.ToIface_JSType_WithSourceRef(type_.GetOwnerFunction())
	} else if type_.GetConstructor() != nil {
		return pkg_rhino.ToIface_JSType_WithSourceRef(type_.GetConstructor())
	}
	return nil
}

// Java static method: JSTypeColorIdHasher.load64(byte[] bytes, int offset)
func JSTypeColorIdHasher_Load64(bytes []int8, offset int32) int64 {
	clinit_JSTypeColorIdHasher()
	_ = bytes
	_ = offset
	return ((((((((int64(bytes[offset]) & int64(255)) | (javalang.Shl64((int64(bytes[(offset+int32(1))]) & int64(255)), int(int32(8))))) | (javalang.Shl64((int64(bytes[(offset+int32(2))]) & int64(255)), int(int32(16))))) | (javalang.Shl64((int64(bytes[(offset+int32(3))]) & int64(255)), int(int32(24))))) | (javalang.Shl64((int64(bytes[(offset+int32(4))]) & int64(255)), int(int32(32))))) | (javalang.Shl64((int64(bytes[(offset+int32(5))]) & int64(255)), int(int32(40))))) | (javalang.Shl64((int64(bytes[(offset+int32(6))]) & int64(255)), int(int32(48))))) | (javalang.Shl64((int64(bytes[(offset+int32(7))]) & int64(255)), int(int32(56)))))
}

// Java static method: JSTypeColorIdHasher.load32(byte[] bytes, int offset)
func JSTypeColorIdHasher_Load32(bytes []int8, offset int32) int64 {
	clinit_JSTypeColorIdHasher()
	_ = bytes
	_ = offset
	return ((((int64(bytes[offset]) & int64(255)) | (javalang.Shl64((int64(bytes[(offset+int32(1))]) & int64(255)), int(int32(8))))) | (javalang.Shl64((int64(bytes[(offset+int32(2))]) & int64(255)), int(int32(16))))) | (javalang.Shl64((int64(bytes[(offset+int32(3))]) & int64(255)), int(int32(24)))))
}

// Java static method: JSTypeColorIdHasher.shiftMix(long val)
func JSTypeColorIdHasher_ShiftMix(val int64) int64 {
	clinit_JSTypeColorIdHasher()
	_ = val
	return (val ^ (javalang.Ushr64(val, int(int32(47)))))
}

// Java static method: JSTypeColorIdHasher.hashLength16(long u, long v, long mul)
func JSTypeColorIdHasher_HashLength16(u int64, v int64, mul int64) int64 {
	clinit_JSTypeColorIdHasher()
	_ = u
	_ = v
	_ = mul
	var a int64 = ((u ^ v) * mul)
	_ = a
	a = (a ^ (javalang.Ushr64(a, int(int32(47)))))
	var b int64 = ((v ^ a) * mul)
	_ = b
	b = (b ^ (javalang.Ushr64(b, int(int32(47)))))
	b = (b * mul)
	return b
}

// Java static method: JSTypeColorIdHasher.farmHash64(byte[] bytes, int length)
func JSTypeColorIdHasher_FarmHash64(bytes []int8, length int32) int64 {
	clinit_JSTypeColorIdHasher()
	_ = bytes
	_ = length
	if length <= int32(16) {
		return JSTypeColorIdHasher_HashLength0to16(bytes, length)
	}
	if length <= int32(32) {
		return JSTypeColorIdHasher_HashLength17to32(bytes, length)
	}
	if length <= int32(64) {
		return JSTypeColorIdHasher_HashLength33To64(bytes, length)
	}
	return JSTypeColorIdHasher_HashLength65Plus(bytes, length)
}

// Java static method: JSTypeColorIdHasher.hashLength0to16(byte[] bytes, int length)
func JSTypeColorIdHasher_HashLength0to16(bytes []int8, length int32) int64 {
	clinit_JSTypeColorIdHasher()
	_ = bytes
	_ = length
	if length >= int32(8) {
		var mul int64 = (int64(-7286425919675154353) + (int64(length) * int64(2)))
		_ = mul
		var a int64 = (JSTypeColorIdHasher_Load64(bytes, int32(0)) + int64(-7286425919675154353))
		_ = a
		var b int64 = JSTypeColorIdHasher_Load64(bytes, (length - int32(8)))
		_ = b
		var c int64 = ((j2gort.BoxLong_RotateRight(b, int32(37)) * mul) + a)
		_ = c
		var d int64 = ((j2gort.BoxLong_RotateRight(a, int32(25)) + b) * mul)
		_ = d
		return JSTypeColorIdHasher_HashLength16(c, d, mul)
	}
	if length >= int32(4) {
		var mul_2 int64 = (int64(-7286425919675154353) + (int64(length) * int64(2)))
		_ = mul_2
		var a_2 int64 = JSTypeColorIdHasher_Load32(bytes, int32(0))
		_ = a_2
		return JSTypeColorIdHasher_HashLength16((int64(length) + (javalang.Shl64(a_2, int(int32(3))))), JSTypeColorIdHasher_Load32(bytes, (length-int32(4))), mul_2)
	}
	if length > int32(0) {
		var a_3 int8 = bytes[int32(0)]
		_ = a_3
		var b_2 int8 = bytes[javalang.Shr32(length, int(int32(1)))]
		_ = b_2
		var c_2 int8 = bytes[(length - int32(1))]
		_ = c_2
		var y int32 = ((int32(a_3) & int32(255)) + (javalang.Shl32((int32(b_2) & int32(255)), int(int32(8)))))
		_ = y
		var z int32 = (length + (javalang.Shl32((int32(c_2) & int32(255)), int(int32(2)))))
		_ = z
		return (JSTypeColorIdHasher_ShiftMix(((int64(y) * int64(-7286425919675154353)) ^ (int64(z) * int64(-4348849565147123417)))) * int64(-7286425919675154353))
	}
	return int64(-7286425919675154353)
}

// Java static method: JSTypeColorIdHasher.hashLength17to32(byte[] bytes, int length)
func JSTypeColorIdHasher_HashLength17to32(bytes []int8, length int32) int64 {
	clinit_JSTypeColorIdHasher()
	_ = bytes
	_ = length
	var mul int64 = (int64(-7286425919675154353) + (int64(length) * int64(2)))
	_ = mul
	var a int64 = (JSTypeColorIdHasher_Load64(bytes, int32(0)) * int64(-5435081209227447693))
	_ = a
	var b int64 = JSTypeColorIdHasher_Load64(bytes, int32(8))
	_ = b
	var c int64 = (JSTypeColorIdHasher_Load64(bytes, (length-int32(8))) * mul)
	_ = c
	var d int64 = (JSTypeColorIdHasher_Load64(bytes, (length-int32(16))) * int64(-7286425919675154353))
	_ = d
	return JSTypeColorIdHasher_HashLength16(((j2gort.BoxLong_RotateRight((a+b), int32(43)) + j2gort.BoxLong_RotateRight(c, int32(30))) + d), ((a + j2gort.BoxLong_RotateRight((b+int64(-7286425919675154353)), int32(18))) + c), mul)
}

// Java static method: JSTypeColorIdHasher.hashLength33To64(byte[] bytes, int length)
func JSTypeColorIdHasher_HashLength33To64(bytes []int8, length int32) int64 {
	clinit_JSTypeColorIdHasher()
	_ = bytes
	_ = length
	var mul int64 = (int64(-7286425919675154353) + (int64(length) * int64(2)))
	_ = mul
	var a int64 = (JSTypeColorIdHasher_Load64(bytes, int32(0)) * int64(-7286425919675154353))
	_ = a
	var b int64 = JSTypeColorIdHasher_Load64(bytes, int32(8))
	_ = b
	var c int64 = (JSTypeColorIdHasher_Load64(bytes, (length-int32(8))) * mul)
	_ = c
	var d int64 = (JSTypeColorIdHasher_Load64(bytes, (length-int32(16))) * int64(-7286425919675154353))
	_ = d
	var y int64 = ((j2gort.BoxLong_RotateRight((a+b), int32(43)) + j2gort.BoxLong_RotateRight(c, int32(30))) + d)
	_ = y
	var z int64 = JSTypeColorIdHasher_HashLength16(y, ((a + j2gort.BoxLong_RotateRight((b+int64(-7286425919675154353)), int32(18))) + c), mul)
	_ = z
	var e int64 = (JSTypeColorIdHasher_Load64(bytes, int32(16)) * mul)
	_ = e
	var f int64 = JSTypeColorIdHasher_Load64(bytes, int32(24))
	_ = f
	var g int64 = ((y + JSTypeColorIdHasher_Load64(bytes, (length-int32(32)))) * mul)
	_ = g
	var h int64 = ((z + JSTypeColorIdHasher_Load64(bytes, (length-int32(24)))) * mul)
	_ = h
	return JSTypeColorIdHasher_HashLength16(((j2gort.BoxLong_RotateRight((e+f), int32(43)) + j2gort.BoxLong_RotateRight(g, int32(30))) + h), ((e + j2gort.BoxLong_RotateRight((f+a), int32(18))) + g), mul)
}

// Java static method: JSTypeColorIdHasher.hashLength65Plus(byte[] bytes, int length)
func JSTypeColorIdHasher_HashLength65Plus(bytes []int8, length int32) int64 {
	clinit_JSTypeColorIdHasher()
	_ = bytes
	_ = length
	var seed int32 = int32(81)
	_ = seed
	var x int64 = int64(seed)
	_ = x
	var y int64 = ((int64(seed) * int64(-5435081209227447693)) + int64(int32(113)))
	_ = y
	var z int64 = (JSTypeColorIdHasher_ShiftMix(((y * int64(-7286425919675154353)) + int64(int32(113)))) * int64(-7286425919675154353))
	_ = z
	var v0 int64 = int64(0)
	_ = v0
	var v1 int64 = int64(0)
	_ = v1
	var w0 int64 = int64(0)
	_ = w0
	var w1 int64 = int64(0)
	_ = w1
	x = ((x * int64(-7286425919675154353)) + JSTypeColorIdHasher_Load64(bytes, int32(0)))
	var offset int32 = int32(0)
	_ = offset
	var end int32 = (((length - int32(1)) / int32(64)) * int32(64))
	_ = end
	var last64offset int32 = ((end + ((length - int32(1)) & int32(63))) - int32(63))
	_ = last64offset
	for doFirst_58 := true; doFirst_58 || (offset != end); doFirst_58 = false {
		x = (j2gort.BoxLong_RotateRight((((x+y)+v0)+JSTypeColorIdHasher_Load64(bytes, (offset+int32(8)))), int32(37)) * int64(-5435081209227447693))
		y = (j2gort.BoxLong_RotateRight(((y+v1)+JSTypeColorIdHasher_Load64(bytes, (offset+int32(48)))), int32(42)) * int64(-5435081209227447693))
		x = (x ^ w1)
		y = (y + (v0 + JSTypeColorIdHasher_Load64(bytes, (offset+int32(40)))))
		z = (j2gort.BoxLong_RotateRight((z+w0), int32(33)) * int64(-5435081209227447693))
		{
			var p1 int64 = JSTypeColorIdHasher_Load64(bytes, offset)
			_ = p1
			var p2 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(8)))
			_ = p2
			var p3 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(16)))
			_ = p3
			var p4 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(24)))
			_ = p4
			var sa int64 = ((v1 * int64(-5435081209227447693)) + p1)
			_ = sa
			var sb int64 = j2gort.BoxLong_RotateRight((((x + w0) + sa) + p4), int32(21))
			_ = sb
			var sc int64 = sa
			_ = sc
			sa = (sa + (p2 + p3))
			sb = (sb + j2gort.BoxLong_RotateRight(sa, int32(44)))
			v0 = (sa + p4)
			v1 = (sb + sc)
		}
		{
			var p1_2 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(32)))
			_ = p1_2
			var p2_2 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(40)))
			_ = p2_2
			var p3_2 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(48)))
			_ = p3_2
			var p4_2 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(56)))
			_ = p4_2
			var sa_2 int64 = ((z + w1) + p1_2)
			_ = sa_2
			var sb_2 int64 = j2gort.BoxLong_RotateRight((((y + JSTypeColorIdHasher_Load64(bytes, (offset+int32(16)))) + sa_2) + p4_2), int32(21))
			_ = sb_2
			var sc_2 int64 = sa_2
			_ = sc_2
			sa_2 = (sa_2 + (p2_2 + p3_2))
			sb_2 = (sb_2 + j2gort.BoxLong_RotateRight(sa_2, int32(44)))
			w0 = (sa_2 + p4_2)
			w1 = (sb_2 + sc_2)
		}
		var tmp int64 = x
		_ = tmp
		x = z
		z = tmp
		offset = (offset + int32(64))
	}
	var mul int64 = (int64(-5435081209227447693) + (javalang.Shl64((z & int64(int32(255))), int(int32(1)))))
	_ = mul
	offset = last64offset
	w0 = (w0 + int64(((length - int32(1)) & int32(63))))
	v0 = (v0 + w0)
	w0 = (w0 + v0)
	x = (j2gort.BoxLong_RotateRight((((x+y)+v0)+JSTypeColorIdHasher_Load64(bytes, (offset+int32(8)))), int32(37)) * mul)
	y = (j2gort.BoxLong_RotateRight(((y+v1)+JSTypeColorIdHasher_Load64(bytes, (offset+int32(48)))), int32(42)) * mul)
	x = (x ^ (w1 * int64(int32(9))))
	y = (y + ((v0 * int64(int32(9))) + JSTypeColorIdHasher_Load64(bytes, (offset+int32(40)))))
	z = (j2gort.BoxLong_RotateRight((z+w0), int32(33)) * mul)
	{
		var p1_3 int64 = JSTypeColorIdHasher_Load64(bytes, offset)
		_ = p1_3
		var p2_3 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(8)))
		_ = p2_3
		var p3_3 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(16)))
		_ = p3_3
		var p4_3 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(24)))
		_ = p4_3
		var sa_3 int64 = ((v1 * mul) + p1_3)
		_ = sa_3
		var sb_3 int64 = j2gort.BoxLong_RotateRight((((x + w0) + sa_3) + p4_3), int32(21))
		_ = sb_3
		var sc_3 int64 = sa_3
		_ = sc_3
		sa_3 = (sa_3 + (p2_3 + p3_3))
		sb_3 = (sb_3 + j2gort.BoxLong_RotateRight(sa_3, int32(44)))
		v0 = (sa_3 + p4_3)
		v1 = (sb_3 + sc_3)
	}
	{
		var p1_4 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(32)))
		_ = p1_4
		var p2_4 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(40)))
		_ = p2_4
		var p3_4 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(48)))
		_ = p3_4
		var p4_4 int64 = JSTypeColorIdHasher_Load64(bytes, (offset + int32(56)))
		_ = p4_4
		var sa_4 int64 = ((z + w1) + p1_4)
		_ = sa_4
		var sb_4 int64 = j2gort.BoxLong_RotateRight((((y + JSTypeColorIdHasher_Load64(bytes, (offset+int32(16)))) + sa_4) + p4_4), int32(21))
		_ = sb_4
		var sc_4 int64 = sa_4
		_ = sc_4
		sa_4 = (sa_4 + (p2_4 + p3_4))
		sb_4 = (sb_4 + j2gort.BoxLong_RotateRight(sa_4, int32(44)))
		w0 = (sa_4 + p4_4)
		w1 = (sb_4 + sc_4)
	}
	return JSTypeColorIdHasher_HashLength16(((JSTypeColorIdHasher_HashLength16(v0, w0, mul) + (JSTypeColorIdHasher_ShiftMix(y) * int64(-4348849565147123417))) + x), (JSTypeColorIdHasher_HashLength16(v1, w1, mul) + z), mul)
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
