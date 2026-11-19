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
// Port of com/google/javascript/jscomp/colors/StandardColors.java (generated mechanically by j2go).

package colors

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.colors.StandardColors
type StandardColors struct {
	idHash_ atomic.Int32
}

func (this_ *StandardColors) BindSelf_(s any) {
	_ = s
}
func (this_ *StandardColors) IsJavaNil_() bool { return this_ == nil }
func (*StandardColors) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.colors.StandardColors", "StandardColors" }
func (this_ *StandardColors) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*StandardColors) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.colors.StandardColors", "com.google.javascript.jscomp.colors.StandardColors", func(v_ any) bool {
		return j2gort.InstanceOf[*StandardColors](v_)
	})
}

var static_StandardColors_ARGUMENTS_ID ColorId
func StaticGet_StandardColors_ARGUMENTS_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_ARGUMENTS_ID
}
func StaticSet_StandardColors_ARGUMENTS_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_ARGUMENTS_ID = v
}
var static_StandardColors_ARRAY_ID ColorId
func StaticGet_StandardColors_ARRAY_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_ARRAY_ID
}
func StaticSet_StandardColors_ARRAY_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_ARRAY_ID = v
}
var static_StandardColors_READONLY_ARRAY_ID ColorId
func StaticGet_StandardColors_READONLY_ARRAY_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_READONLY_ARRAY_ID
}
func StaticSet_StandardColors_READONLY_ARRAY_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_READONLY_ARRAY_ID = v
}
var static_StandardColors_GENERATOR_ID ColorId
func StaticGet_StandardColors_GENERATOR_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_GENERATOR_ID
}
func StaticSet_StandardColors_GENERATOR_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_GENERATOR_ID = v
}
var static_StandardColors_I_TEMPLATE_ARRAY_ID ColorId
func StaticGet_StandardColors_I_TEMPLATE_ARRAY_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_I_TEMPLATE_ARRAY_ID
}
func StaticSet_StandardColors_I_TEMPLATE_ARRAY_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_I_TEMPLATE_ARRAY_ID = v
}
var static_StandardColors_ITERATOR_ID ColorId
func StaticGet_StandardColors_ITERATOR_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_ITERATOR_ID
}
func StaticSet_StandardColors_ITERATOR_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_ITERATOR_ID = v
}
var static_StandardColors_ASYNC_ITERATOR_ITERABLE_ID ColorId
func StaticGet_StandardColors_ASYNC_ITERATOR_ITERABLE_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_ASYNC_ITERATOR_ITERABLE_ID
}
func StaticSet_StandardColors_ASYNC_ITERATOR_ITERABLE_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_ASYNC_ITERATOR_ITERABLE_ID = v
}
var static_StandardColors_PROMISE_ID ColorId
func StaticGet_StandardColors_PROMISE_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_PROMISE_ID
}
func StaticSet_StandardColors_PROMISE_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_PROMISE_ID = v
}
var static_StandardColors_ARRAY_PROTOTYPE_ID ColorId
func StaticGet_StandardColors_ARRAY_PROTOTYPE_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_ARRAY_PROTOTYPE_ID
}
func StaticSet_StandardColors_ARRAY_PROTOTYPE_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_ARRAY_PROTOTYPE_ID = v
}
var static_StandardColors_READONLY_ARRAY_PROTOTYPE_ID ColorId
func StaticGet_StandardColors_READONLY_ARRAY_PROTOTYPE_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_READONLY_ARRAY_PROTOTYPE_ID
}
func StaticSet_StandardColors_READONLY_ARRAY_PROTOTYPE_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_READONLY_ARRAY_PROTOTYPE_ID = v
}
var static_StandardColors_I_TEMPLATE_ARRAY_PROTOTYPE_ID ColorId
func StaticGet_StandardColors_I_TEMPLATE_ARRAY_PROTOTYPE_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_I_TEMPLATE_ARRAY_PROTOTYPE_ID
}
func StaticSet_StandardColors_I_TEMPLATE_ARRAY_PROTOTYPE_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_I_TEMPLATE_ARRAY_PROTOTYPE_ID = v
}
var static_StandardColors_MAP_ID ColorId
func StaticGet_StandardColors_MAP_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_MAP_ID
}
func StaticSet_StandardColors_MAP_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_MAP_ID = v
}
var static_StandardColors_MAP_PROTOTYPE_ID ColorId
func StaticGet_StandardColors_MAP_PROTOTYPE_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_MAP_PROTOTYPE_ID
}
func StaticSet_StandardColors_MAP_PROTOTYPE_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_MAP_PROTOTYPE_ID = v
}
var static_StandardColors_READONLY_MAP_ID ColorId
func StaticGet_StandardColors_READONLY_MAP_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_READONLY_MAP_ID
}
func StaticSet_StandardColors_READONLY_MAP_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_READONLY_MAP_ID = v
}
var static_StandardColors_READONLY_MAP_PROTOTYPE_ID ColorId
func StaticGet_StandardColors_READONLY_MAP_PROTOTYPE_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_READONLY_MAP_PROTOTYPE_ID
}
func StaticSet_StandardColors_READONLY_MAP_PROTOTYPE_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_READONLY_MAP_PROTOTYPE_ID = v
}
var static_StandardColors_SET_ID ColorId
func StaticGet_StandardColors_SET_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_SET_ID
}
func StaticSet_StandardColors_SET_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_SET_ID = v
}
var static_StandardColors_SET_PROTOTYPE_ID ColorId
func StaticGet_StandardColors_SET_PROTOTYPE_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_SET_PROTOTYPE_ID
}
func StaticSet_StandardColors_SET_PROTOTYPE_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_SET_PROTOTYPE_ID = v
}
var static_StandardColors_READONLY_SET_ID ColorId
func StaticGet_StandardColors_READONLY_SET_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_READONLY_SET_ID
}
func StaticSet_StandardColors_READONLY_SET_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_READONLY_SET_ID = v
}
var static_StandardColors_READONLY_SET_PROTOTYPE_ID ColorId
func StaticGet_StandardColors_READONLY_SET_PROTOTYPE_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_READONLY_SET_PROTOTYPE_ID
}
func StaticSet_StandardColors_READONLY_SET_PROTOTYPE_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_READONLY_SET_PROTOTYPE_ID = v
}
var static_StandardColors_BIGINT_OBJECT_ID ColorId
func StaticGet_StandardColors_BIGINT_OBJECT_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_BIGINT_OBJECT_ID
}
func StaticSet_StandardColors_BIGINT_OBJECT_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_BIGINT_OBJECT_ID = v
}
var static_StandardColors_BOOLEAN_OBJECT_ID ColorId
func StaticGet_StandardColors_BOOLEAN_OBJECT_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_BOOLEAN_OBJECT_ID
}
func StaticSet_StandardColors_BOOLEAN_OBJECT_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_BOOLEAN_OBJECT_ID = v
}
var static_StandardColors_NUMBER_OBJECT_ID ColorId
func StaticGet_StandardColors_NUMBER_OBJECT_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_NUMBER_OBJECT_ID
}
func StaticSet_StandardColors_NUMBER_OBJECT_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_NUMBER_OBJECT_ID = v
}
var static_StandardColors_STRING_OBJECT_ID ColorId
func StaticGet_StandardColors_STRING_OBJECT_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_STRING_OBJECT_ID
}
func StaticSet_StandardColors_STRING_OBJECT_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_STRING_OBJECT_ID = v
}
var static_StandardColors_SYMBOL_OBJECT_ID ColorId
func StaticGet_StandardColors_SYMBOL_OBJECT_ID() ColorId {
	clinit_StandardColors()
	return static_StandardColors_SYMBOL_OBJECT_ID
}
func StaticSet_StandardColors_SYMBOL_OBJECT_ID(v ColorId) {
	clinit_StandardColors()
	static_StandardColors_SYMBOL_OBJECT_ID = v
}
var static_StandardColors_BIGINT *Color
func StaticGet_StandardColors_BIGINT() *Color {
	clinit_StandardColors()
	return static_StandardColors_BIGINT
}
func StaticSet_StandardColors_BIGINT(v *Color) {
	clinit_StandardColors()
	static_StandardColors_BIGINT = v
}
var static_StandardColors_BOOLEAN *Color
func StaticGet_StandardColors_BOOLEAN() *Color {
	clinit_StandardColors()
	return static_StandardColors_BOOLEAN
}
func StaticSet_StandardColors_BOOLEAN(v *Color) {
	clinit_StandardColors()
	static_StandardColors_BOOLEAN = v
}
var static_StandardColors_NULL_OR_VOID *Color
func StaticGet_StandardColors_NULL_OR_VOID() *Color {
	clinit_StandardColors()
	return static_StandardColors_NULL_OR_VOID
}
func StaticSet_StandardColors_NULL_OR_VOID(v *Color) {
	clinit_StandardColors()
	static_StandardColors_NULL_OR_VOID = v
}
var static_StandardColors_NUMBER *Color
func StaticGet_StandardColors_NUMBER() *Color {
	clinit_StandardColors()
	return static_StandardColors_NUMBER
}
func StaticSet_StandardColors_NUMBER(v *Color) {
	clinit_StandardColors()
	static_StandardColors_NUMBER = v
}
var static_StandardColors_STRING *Color
func StaticGet_StandardColors_STRING() *Color {
	clinit_StandardColors()
	return static_StandardColors_STRING
}
func StaticSet_StandardColors_STRING(v *Color) {
	clinit_StandardColors()
	static_StandardColors_STRING = v
}
var static_StandardColors_SYMBOL *Color
func StaticGet_StandardColors_SYMBOL() *Color {
	clinit_StandardColors()
	return static_StandardColors_SYMBOL
}
func StaticSet_StandardColors_SYMBOL(v *Color) {
	clinit_StandardColors()
	static_StandardColors_SYMBOL = v
}
var static_StandardColors_TOP_OBJECT *Color
func StaticGet_StandardColors_TOP_OBJECT() *Color {
	clinit_StandardColors()
	return static_StandardColors_TOP_OBJECT
}
func StaticSet_StandardColors_TOP_OBJECT(v *Color) {
	clinit_StandardColors()
	static_StandardColors_TOP_OBJECT = v
}
var static_StandardColors_TOP_FUNCTION *Color
func StaticGet_StandardColors_TOP_FUNCTION() *Color {
	clinit_StandardColors()
	return static_StandardColors_TOP_FUNCTION
}
func StaticSet_StandardColors_TOP_FUNCTION(v *Color) {
	clinit_StandardColors()
	static_StandardColors_TOP_FUNCTION = v
}
var static_StandardColors_UNKNOWN *Color
func StaticGet_StandardColors_UNKNOWN() *Color {
	clinit_StandardColors()
	return static_StandardColors_UNKNOWN
}
func StaticSet_StandardColors_UNKNOWN(v *Color) {
	clinit_StandardColors()
	static_StandardColors_UNKNOWN = v
}
var static_StandardColors_GBIGINT *Color
func StaticGet_StandardColors_GBIGINT() *Color {
	clinit_StandardColors()
	return static_StandardColors_GBIGINT
}
func StaticSet_StandardColors_GBIGINT(v *Color) {
	clinit_StandardColors()
	static_StandardColors_GBIGINT = v
}
var static_StandardColors_AXIOMATIC_COLORS j2gort.Ext_com_google_common_collect_ImmutableMap
func StaticGet_StandardColors_AXIOMATIC_COLORS() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_StandardColors()
	return static_StandardColors_AXIOMATIC_COLORS
}
func StaticSet_StandardColors_AXIOMATIC_COLORS(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_StandardColors()
	static_StandardColors_AXIOMATIC_COLORS = v
}
var static_StandardColors_PRIMITIVE_COLORS j2gort.Ext_com_google_common_collect_ImmutableMap
func StaticGet_StandardColors_PRIMITIVE_COLORS() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_StandardColors()
	return static_StandardColors_PRIMITIVE_COLORS
}
func StaticSet_StandardColors_PRIMITIVE_COLORS(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_StandardColors()
	static_StandardColors_PRIMITIVE_COLORS = v
}
var static_StandardColors_PRIMITIVE_BOX_IDS j2gort.Ext_com_google_common_collect_ImmutableSet
func StaticGet_StandardColors_PRIMITIVE_BOX_IDS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_StandardColors()
	return static_StandardColors_PRIMITIVE_BOX_IDS
}
func StaticSet_StandardColors_PRIMITIVE_BOX_IDS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_StandardColors()
	static_StandardColors_PRIMITIVE_BOX_IDS = v
}
var static_StandardColors_STANDARD_OBJECT_IDS j2gort.Ext_com_google_common_collect_ImmutableSet
func StaticGet_StandardColors_STANDARD_OBJECT_IDS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_StandardColors()
	return static_StandardColors_STANDARD_OBJECT_IDS
}
func StaticSet_StandardColors_STANDARD_OBJECT_IDS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_StandardColors()
	static_StandardColors_STANDARD_OBJECT_IDS = v
}
var clinitDone_StandardColors uint32
var clinitState_StandardColors j2gort.ClinitState
func clinit_StandardColors() {
	if atomic.LoadUint32(&clinitDone_StandardColors) != 0 {
		return
	}
	clinitSlow_StandardColors()
}

func clinitSlow_StandardColors() {
	if !j2gort.EnterClinit(&clinitState_StandardColors, &clinitDone_StandardColors) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_StandardColors, &clinitDone_StandardColors)
	static_StandardColors_ARGUMENTS_ID = ColorId_FromUnsigned_1_int32(int32(423208557))
	static_StandardColors_ARRAY_ID = ColorId_FromUnsigned_1_int32(int32(2043979267))
	static_StandardColors_READONLY_ARRAY_ID = ColorId_FromUnsigned_1_int32(int32(1445461936))
	static_StandardColors_GENERATOR_ID = ColorId_FromUnsigned_1_int32(int32(-1682886593))
	static_StandardColors_I_TEMPLATE_ARRAY_ID = ColorId_FromUnsigned_1_int32(int32(1185627918))
	static_StandardColors_ITERATOR_ID = ColorId_FromUnsigned_1_int32(int32(1098830507))
	static_StandardColors_ASYNC_ITERATOR_ITERABLE_ID = ColorId_FromUnsigned_1_int32(int32(-885510646))
	static_StandardColors_PROMISE_ID = ColorId_FromUnsigned_1_int32(int32(962075327))
	static_StandardColors_ARRAY_PROTOTYPE_ID = ColorId_FromUnsigned_1_int64(int64(7260434188007790327))
	static_StandardColors_READONLY_ARRAY_PROTOTYPE_ID = ColorId_FromUnsigned_1_int64(int64(3857846951871685759))
	static_StandardColors_I_TEMPLATE_ARRAY_PROTOTYPE_ID = ColorId_FromUnsigned_1_int64(int64(8809266781907489176))
	static_StandardColors_MAP_ID = ColorId_FromUnsigned_1_int64(int64(3720784123213224612))
	static_StandardColors_MAP_PROTOTYPE_ID = ColorId_FromUnsigned_1_int64(int64(7854789731312742870))
	static_StandardColors_READONLY_MAP_ID = ColorId_FromUnsigned_1_int64(int64(263136405622606104))
	static_StandardColors_READONLY_MAP_PROTOTYPE_ID = ColorId_FromUnsigned_1_int64(int64(1552746670790408862))
	static_StandardColors_SET_ID = ColorId_FromUnsigned_1_int64(int64(5640978925949894333))
	static_StandardColors_SET_PROTOTYPE_ID = ColorId_FromUnsigned_1_int64(int64(-106931852624587719))
	static_StandardColors_READONLY_SET_ID = ColorId_FromUnsigned_1_int64(int64(3496494622634221469))
	static_StandardColors_READONLY_SET_PROTOTYPE_ID = ColorId_FromUnsigned_1_int64(int64(-8282285649609606861))
	static_StandardColors_BIGINT_OBJECT_ID = ColorId_FromUnsigned_1_int32(int32(-1445352083))
	static_StandardColors_BOOLEAN_OBJECT_ID = ColorId_FromUnsigned_1_int32(int32(-1845109754))
	static_StandardColors_NUMBER_OBJECT_ID = ColorId_FromUnsigned_1_int32(int32(884617137))
	static_StandardColors_STRING_OBJECT_ID = ColorId_FromUnsigned_1_int32(int32(408946857))
	static_StandardColors_SYMBOL_OBJECT_ID = ColorId_FromUnsigned_1_int32(int32(1582387070))
	static_StandardColors_BIGINT = Color_SingleBuilder().SetId(ColorId_FromUnsigned_1_int32(int32(592360986))).SetBoxId(StaticGet_StandardColors_BIGINT_OBJECT_ID()).SetInvalidating(false).BuildAxiomatic()
	static_StandardColors_BOOLEAN = Color_SingleBuilder().SetId(ColorId_FromUnsigned_1_int32(int32(308810478))).SetBoxId(StaticGet_StandardColors_BOOLEAN_OBJECT_ID()).SetInvalidating(false).BuildAxiomatic()
	static_StandardColors_NULL_OR_VOID = Color_SingleBuilder().SetBoxId((ColorId{})).SetId(ColorId_FromUnsigned_1_int32(int32(582262633))).SetInvalidating(false).BuildAxiomatic()
	static_StandardColors_NUMBER = Color_SingleBuilder().SetId(ColorId_FromUnsigned_1_int32(int32(-796822996))).SetBoxId(StaticGet_StandardColors_NUMBER_OBJECT_ID()).SetInvalidating(false).BuildAxiomatic()
	static_StandardColors_STRING = Color_SingleBuilder().SetId(ColorId_FromUnsigned_1_int32(int32(-1941074075))).SetBoxId(StaticGet_StandardColors_STRING_OBJECT_ID()).SetInvalidating(false).BuildAxiomatic()
	static_StandardColors_SYMBOL = Color_SingleBuilder().SetId(ColorId_FromUnsigned_1_int32(int32(1973362790))).SetBoxId(StaticGet_StandardColors_SYMBOL_OBJECT_ID()).SetInvalidating(false).BuildAxiomatic()
	static_StandardColors_TOP_OBJECT = Color_SingleBuilder().SetId(ColorId_FromUnsigned_1_int32(int32(-2003081160))).SetBoxId((ColorId{})).SetInvalidating(true).BuildAxiomatic()
	static_StandardColors_TOP_FUNCTION = Color_SingleBuilder().SetId(ColorId_FromUnsigned_1_int32(int32(674977848))).SetBoxId((ColorId{})).SetInvalidating(true).BuildAxiomatic()
	static_StandardColors_UNKNOWN = Color_SingleBuilder().SetBoxId((ColorId{})).SetId(ColorId_FromUnsigned_1_int32(int32(0))).SetInvalidating(true).BuildAxiomatic()
	static_StandardColors_GBIGINT = Color_SingleBuilder().SetId(ColorId_FromUnsigned_1_int32(int32(-1431423368))).SetBoxId((ColorId{})).SetInvalidating(false).BuildAxiomatic()
	static_StandardColors_AXIOMATIC_COLORS = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder().Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_BIGINT().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_BIGINT()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_BOOLEAN().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_BOOLEAN()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_NULL_OR_VOID().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_NULL_OR_VOID()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_NUMBER().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_NUMBER()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_STRING().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_STRING()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_SYMBOL().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_SYMBOL()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_TOP_OBJECT().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_TOP_OBJECT()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_UNKNOWN().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_UNKNOWN()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_GBIGINT().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_GBIGINT()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_TOP_FUNCTION().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_TOP_FUNCTION()))).BuildOrThrow())
	static_StandardColors_PRIMITIVE_COLORS = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap_Builder](j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder().Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_BIGINT().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_BIGINT()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_BOOLEAN().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_BOOLEAN()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_NULL_OR_VOID().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_NULL_OR_VOID()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_NUMBER().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_NUMBER()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_STRING().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_STRING()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_SYMBOL().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_SYMBOL()))).Put_2(j2gort.BoxNullableVal(StaticGet_StandardColors_GBIGINT().GetId()), j2gort.PtrToAny(StaticGet_StandardColors_GBIGINT()))).BuildOrThrow())
	static_StandardColors_PRIMITIVE_BOX_IDS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_5(j2gort.BoxNullableVal(StaticGet_StandardColors_BIGINT_OBJECT_ID()), j2gort.BoxNullableVal(StaticGet_StandardColors_BOOLEAN_OBJECT_ID()), j2gort.BoxNullableVal(StaticGet_StandardColors_NUMBER_OBJECT_ID()), j2gort.BoxNullableVal(StaticGet_StandardColors_STRING_OBJECT_ID()), j2gort.BoxNullableVal(StaticGet_StandardColors_SYMBOL_OBJECT_ID()))
	static_StandardColors_STANDARD_OBJECT_IDS = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_StandardColors_PRIMITIVE_BOX_IDS()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_ARRAY_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_ARRAY_PROTOTYPE_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_READONLY_ARRAY_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_READONLY_ARRAY_PROTOTYPE_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_ARGUMENTS_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_ASYNC_ITERATOR_ITERABLE_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_GENERATOR_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_I_TEMPLATE_ARRAY_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_I_TEMPLATE_ARRAY_PROTOTYPE_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_ITERATOR_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_MAP_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_MAP_PROTOTYPE_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_READONLY_MAP_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_READONLY_MAP_PROTOTYPE_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_SET_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_SET_PROTOTYPE_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_READONLY_SET_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_READONLY_SET_PROTOTYPE_ID()))).Add_1_any(j2gort.BoxNullableVal(StaticGet_StandardColors_PROMISE_ID()))).Build())
}

func (this_ *StandardColors) initFields_StandardColors() {
	_ = this_
}

// Java constructor: StandardColors()
func NewStandardColors() *StandardColors {
	clinit_StandardColors()
	this_ := &StandardColors{}
	InitStandardColors(this_)
	return this_
}

func InitStandardColors(this_ *StandardColors) {
	this_.initFields_StandardColors()
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

