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
// Port of com/google/javascript/jscomp/colors/Color.java (generated mechanically by j2go).

package colors

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.colors.Color
type Color struct {
	prototypes                 []*Color   `gson:"prototypes"`
	instanceColors             []*Color   `gson:"instanceColors"`
	ownProperties              []string   `gson:"ownProperties"`
	unionElements              []*Color   `gson:"unionElements"`
	subtractNullOrVoid         *Color     `gson:"-"`
	id                         ColorId    `gson:"id"`
	boxId                      ColorId    `gson:"boxId"`
	_syncMu                    sync.Mutex `gson:"-"`
	_tcHandle                  atomic.Uint32
	hashCode                   int32 `gson:"-"`
	invalidating               bool  `gson:"invalidating"`
	propertiesKeepOriginalName bool  `gson:"propertiesKeepOriginalName"`
	constructor                bool  `gson:"constructor"`
	closureAssert              bool  `gson:"closureAssert"`
}

func (this_ *Color) BindSelf_(s any) {
	_ = s
}
func (this_ *Color) IsJavaNil_() bool { return this_ == nil }
func (*Color) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.colors.Color", "Color"
}
func (this_ *Color) TCHandlePtr_() *atomic.Uint32 { return &this_._tcHandle }
func (*Color) Implements_colors_Color()           {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.colors.Color", "com.google.javascript.jscomp.colors.Color", func(v_ any) bool {
		return j2gort.InstanceOf[*Color](v_)
	})
}

var static_Color_EMPTY_COLORS []*Color

func StaticGet_Color_EMPTY_COLORS() []*Color {
	clinit_Color()
	return static_Color_EMPTY_COLORS
}
func StaticSet_Color_EMPTY_COLORS(v []*Color) {
	clinit_Color()
	static_Color_EMPTY_COLORS = v
}

var static_Color_EMPTY_STRINGS []string

func StaticGet_Color_EMPTY_STRINGS() []string {
	clinit_Color()
	return static_Color_EMPTY_STRINGS
}
func StaticSet_Color_EMPTY_STRINGS(v []string) {
	clinit_Color()
	static_Color_EMPTY_STRINGS = v
}

var clinitDone_Color uint32
var clinitState_Color j2gort.ClinitState

func clinit_Color() {
	if atomic.LoadUint32(&clinitDone_Color) != 0 {
		return
	}
	clinitSlow_Color()
}

func clinitSlow_Color() {
	if !j2gort.EnterClinit(&clinitState_Color, &clinitDone_Color) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Color, &clinitDone_Color)
	static_Color_EMPTY_COLORS = make([]*Color, int32(0))
	static_Color_EMPTY_STRINGS = j2gort.NewStringArray(int32(0))
}

func (this_ *Color) initFields_Color() {
	_ = this_
}

// Java constructor: Color(com.google.javascript.jscomp.colors.ColorId id, com.google.javascript.jscomp.colors.Color[] prototypes, com.google.javascript.jscomp.colors.Color[] instanceColors, boolean invalidating, boolean propertiesKeepOriginalName, boolean constructor, java.lang.String[] ownProperties, com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable ColorId boxId, boolean closureAssert, com.google.javascript.jscomp.colors.Color[] unionElements)
func NewColor(id ColorId, prototypes []*Color, instanceColors []*Color, invalidating bool, propertiesKeepOriginalName bool, constructor bool, ownProperties []string, boxId ColorId, closureAssert bool, unionElements []*Color) *Color {
	clinit_Color()
	this_ := &Color{}
	InitColor(this_, id, prototypes, instanceColors, invalidating, propertiesKeepOriginalName, constructor, ownProperties, boxId, closureAssert, unionElements)
	return this_
}

func InitColor(this_ *Color, id ColorId, prototypes []*Color, instanceColors []*Color, invalidating bool, propertiesKeepOriginalName bool, constructor bool, ownProperties []string, boxId ColorId, closureAssert bool, unionElements []*Color) {
	_ = id
	_ = prototypes
	_ = instanceColors
	_ = invalidating
	_ = propertiesKeepOriginalName
	_ = constructor
	_ = ownProperties
	_ = boxId
	_ = closureAssert
	_ = unionElements
	this_.initFields_Color()
	if id == (ColorId{}) {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null id")))
	}
	this_.id = id
	if prototypes == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null prototypes")))
	}
	this_.prototypes = prototypes
	if instanceColors == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null instanceColors")))
	}
	this_.instanceColors = instanceColors
	this_.invalidating = invalidating
	this_.propertiesKeepOriginalName = propertiesKeepOriginalName
	this_.constructor = constructor
	if ownProperties == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null ownProperties")))
	}
	this_.ownProperties = ownProperties
	this_.boxId = boxId
	this_.closureAssert = closureAssert
	if unionElements == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null unionElements")))
	}
	this_.unionElements = unionElements
}

// Java method: Color.getId()
func (this_ *Color) GetId() ColorId {
	return this_.id
}

// Java method: Color.getPrototypes()
func (this_ *Color) GetPrototypes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return Color_ToImmutableColorSet(this_.prototypes)
}

// Java method: Color.getPrototypesArray()
func (this_ *Color) GetPrototypesArray() []*Color {
	return this_.prototypes
}

// Java method: Color.getInstanceColors()
func (this_ *Color) GetInstanceColors() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return Color_ToImmutableColorSet(this_.instanceColors)
}

// Java method: Color.getInstanceColorsArray()
func (this_ *Color) GetInstanceColorsArray() []*Color {
	return this_.instanceColors
}

// Java method: Color.hasInstanceColors()
func (this_ *Color) HasInstanceColors() bool {
	return (int32(len(this_.instanceColors)) > int32(0))
}

// Java method: Color.isInvalidating()
func (this_ *Color) IsInvalidating() bool {
	return this_.invalidating
}

// Java method: Color.getPropertiesKeepOriginalName()
func (this_ *Color) GetPropertiesKeepOriginalName() bool {
	return this_.propertiesKeepOriginalName
}

// Java method: Color.isConstructor()
func (this_ *Color) IsConstructor() bool {
	return this_.constructor
}

// Java method: Color.getOwnProperties()
func (this_ *Color) GetOwnProperties() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return Color_ToImmutableStringSet(this_.ownProperties)
}

// Java method: Color.getOwnPropertiesArray()
func (this_ *Color) GetOwnPropertiesArray() []string {
	return this_.ownProperties
}

// Java method: Color.getBoxId()
func (this_ *Color) GetBoxId() ColorId {
	return this_.boxId
}

// Java method: Color.isClosureAssert()
func (this_ *Color) IsClosureAssert() bool {
	return this_.closureAssert
}

// Java method: Color.getUnionElements()
func (this_ *Color) GetUnionElements() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return Color_ToImmutableColorSet(this_.unionElements)
}

// Java method: Color.getUnionElementsArray()
func (this_ *Color) GetUnionElementsArray() []*Color {
	return this_.unionElements
}

// Java static method: Color.toImmutableColorSet(com.google.javascript.jscomp.colors.Color[] arr)
func Color_ToImmutableColorSet(arr []*Color) j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_Color()
	_ = arr
	var len_ int32 = int32(len(arr))
	_ = len_
	if len_ == int32(0) {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	}
	if len_ == int32(1) {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_1(j2gort.PtrToAny(arr[int32(0)]))
	}
	return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_slice_any(func() []any {
		arr_1 := arr
		if arr_1 == nil {
			return nil
		}
		out_2 := make([]any, len(arr_1))
		for i_, el_3 := range arr_1 {
			out_2[i_] = j2gort.PtrToAny(el_3)
		}
		return out_2
	}())
}

// Java static method: Color.toImmutableStringSet(java.lang.String[] arr)
func Color_ToImmutableStringSet(arr []string) j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_Color()
	_ = arr
	var len_ int32 = int32(len(arr))
	_ = len_
	if len_ == int32(0) {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	}
	if len_ == int32(1) {
		return j2gort.ImmSetOfStr1(arr[int32(0)])
	}
	return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_slice_any(func() []any {
		arr_4 := arr
		if arr_4 == nil {
			return nil
		}
		out_5 := make([]any, len(arr_4))
		for i_, el_6 := range arr_4 {
			out_5[i_] = j2gort.BoxString(el_6)
		}
		return out_5
	}())
}

// Java static method: Color.toUniqueColorArray(java.util.Collection<com.google.javascript.jscomp.colors.Color> colors)
func Color_ToUniqueColorArray(colors j2gort.Ext_java_util_Collection) []*Color {
	clinit_Color()
	_ = colors
	var size int32 = colors.Size()
	_ = size
	if size == int32(0) {
		return StaticGet_Color_EMPTY_COLORS()
	}
	if j2gort.InstanceOf[j2gort.Ext_java_util_Set](colors) {
		var out []*Color = make([]*Color, size)
		_ = out
		var i int32 = int32(0)
		_ = i
		{
			cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(colors))
			for cur_7.HasNext() {
				rawElem_8 := cur_7.Next()
				c := j2gort.CastPtr[Color](rawElem_8)
				_ = c
				j2gort.CheckNotNullPtr(c)
				{
					a_9 := out
					i_10 := func() int32 { old_11 := i; nw_12 := old_11 + 1; i = nw_12; return old_11 }()
					a_9[i_10] = c
				}
			}
		}
		return out
	}
	var unique *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(size)
	_ = unique
	{
		cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(colors))
		for cur_13.HasNext() {
			rawElem_14 := cur_13.Next()
			c_2 := j2gort.CastPtr[Color](rawElem_14)
			_ = c_2
			j2gort.CheckNotNullPtr(c_2)
			unique.Add_1(j2gort.PtrToAny(c_2))
		}
	}
	return func() []*Color {
		arr_18 := unique.ToArray_1_slice_any(func() []any {
			arr_15 := StaticGet_Color_EMPTY_COLORS()
			if arr_15 == nil {
				return nil
			}
			out_16 := make([]any, len(arr_15))
			for i_, el_17 := range arr_15 {
				out_16[i_] = j2gort.PtrToAny(el_17)
			}
			return out_16
		}())
		if arr_18 == nil {
			return nil
		}
		out_19 := make([]*Color, len(arr_18))
		for i_, el_20 := range arr_18 {
			out_19[i_] = j2gort.CastPtr[Color](el_20)
		}
		return out_19
	}()
}

// Java static method: Color.toUniqueStringArray(java.util.Collection<java.lang.String> strings)
func Color_ToUniqueStringArray(strings j2gort.Ext_java_util_Collection) []string {
	clinit_Color()
	_ = strings
	var size int32 = strings.Size()
	_ = size
	if size == int32(0) {
		return StaticGet_Color_EMPTY_STRINGS()
	}
	if j2gort.InstanceOf[j2gort.Ext_java_util_Set](strings) {
		var out []string = j2gort.NewStringArray(size)
		_ = out
		var i int32 = int32(0)
		_ = i
		{
			cur_21 := j2gort.NewIterCursor(j2gort.CheckNotNull(strings))
			for cur_21.HasNext() {
				s := cur_21.NextStr()
				_ = s
				j2gort.CheckNotNullString(s)
				{
					a_22 := out
					i_23 := func() int32 { old_24 := i; nw_25 := old_24 + 1; i = nw_25; return old_24 }()
					a_22[i_23] = s
				}
			}
		}
		return out
	}
	var unique *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(size)
	_ = unique
	{
		cur_26 := j2gort.NewIterCursor(j2gort.CheckNotNull(strings))
		for cur_26.HasNext() {
			s_2 := cur_26.NextStr()
			_ = s_2
			j2gort.CheckNotNullString(s_2)
			j2gort.SetAddStr(unique, s_2)
		}
	}
	return func() []string {
		arr_30 := unique.ToArray_1_slice_any(func() []any {
			arr_27 := StaticGet_Color_EMPTY_STRINGS()
			if arr_27 == nil {
				return nil
			}
			out_28 := make([]any, len(arr_27))
			for i_, el_29 := range arr_27 {
				out_28[i_] = j2gort.BoxString(el_29)
			}
			return out_28
		}())
		if arr_30 == nil {
			return nil
		}
		out_31 := make([]string, len(arr_30))
		for i_, el_32 := range arr_30 {
			out_31[i_] = j2gort.CastString(el_32)
		}
		return out_31
	}()
}

// Java static method: Color.containsColor(com.google.javascript.jscomp.colors.Color[] arr, int limit, com.google.javascript.jscomp.colors.Color target)
func Color_ContainsColor(arr []*Color, limit int32, target *Color) bool {
	clinit_Color()
	_ = arr
	_ = limit
	_ = target
	{
		var i int32 = int32(0)
		_ = i
		for ; i < limit; i = i + 1 {
			if arr[i].Equals(j2gort.PtrToAny(target)) {
				return true
			}
		}
	}
	return false
}

// Java static method: Color.containsString(java.lang.String[] arr, int limit, java.lang.String target)
func Color_ContainsString(arr []string, limit int32, target string) bool {
	clinit_Color()
	_ = arr
	_ = limit
	_ = target
	{
		var i int32 = int32(0)
		_ = i
		for ; i < limit; i = i + 1 {
			if j2gort.CheckNotNullString(arr[i]) == target {
				return true
			}
		}
	}
	return false
}

// Java static method: Color.mergeColorArrays(com.google.javascript.jscomp.colors.Color[] a, com.google.javascript.jscomp.colors.Color[] b)
func Color_MergeColorArrays(a []*Color, b []*Color) []*Color {
	clinit_Color()
	_ = a
	_ = b
	var aLen int32 = int32(len(a))
	_ = aLen
	if aLen == int32(0) {
		return b
	}
	var bLen int32 = int32(len(b))
	_ = bLen
	if (bLen == int32(0)) || j2gort.SliceRefEquals(a, b) {
		return a
	}
	if (aLen == int32(1)) && (bLen == int32(1)) {
		if a[int32(0)].Equals(j2gort.PtrToAny(b[int32(0)])) {
			return a
		} else {
			return []*Color{a[int32(0)], b[int32(0)]}
		}
	}
	if (aLen + bLen) <= int32(8) {
		var extra int32 = int32(0)
		_ = extra
		{
			var i int32 = int32(0)
			_ = i
			for ; i < bLen; i = i + 1 {
				if !(Color_ContainsColor(a, aLen, b[i])) {
					extra = extra + 1
				}
			}
		}
		if extra == int32(0) {
			return a
		}
		var merged []*Color = make([]*Color, (aLen + extra))
		_ = merged
		func() {
			acSrc_33 := a
			acSP_34 := int32(0)
			acDst_35 := merged
			acDP_36 := int32(0)
			acLen_37 := aLen
			j2gort.CheckNotNullSlice(acSrc_33)
			j2gort.CheckNotNullSlice(acDst_35)
			if acSP_34 < 0 || acDP_36 < 0 || acLen_37 < 0 || int64(acSP_34)+int64(acLen_37) > int64(len(acSrc_33)) || int64(acDP_36)+int64(acLen_37) > int64(len(acDst_35)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_35[int(acDP_36):int(acDP_36)+int(acLen_37)], acSrc_33[int(acSP_34):int(acSP_34)+int(acLen_37)])
		}()
		var pos int32 = aLen
		_ = pos
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < bLen; i_2 = i_2 + 1 {
				var c *Color = b[i_2]
				_ = c
				if !(Color_ContainsColor(a, aLen, c)) {
					{
						a_38 := merged
						i_39 := func() int32 { old_40 := pos; nw_41 := old_40 + 1; pos = nw_41; return old_40 }()
						a_38[i_39] = c
					}
				}
			}
		}
		return merged
	}
	var set *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32((aLen + bLen))
	_ = set
	{
		var i_3 int32 = int32(0)
		_ = i_3
		for ; i_3 < aLen; i_3 = i_3 + 1 {
			set.Add_1(j2gort.PtrToAny(a[i_3]))
		}
	}
	var addedNew bool = false
	_ = addedNew
	{
		var i_4 int32 = int32(0)
		_ = i_4
		for ; i_4 < bLen; i_4 = i_4 + 1 {
			if set.Add_1(j2gort.PtrToAny(b[i_4])) {
				addedNew = true
			}
		}
	}
	if !(addedNew) {
		return a
	}
	return func() []*Color {
		arr_45 := set.ToArray_1_slice_any(func() []any {
			arr_42 := StaticGet_Color_EMPTY_COLORS()
			if arr_42 == nil {
				return nil
			}
			out_43 := make([]any, len(arr_42))
			for i_, el_44 := range arr_42 {
				out_43[i_] = j2gort.PtrToAny(el_44)
			}
			return out_43
		}())
		if arr_45 == nil {
			return nil
		}
		out_46 := make([]*Color, len(arr_45))
		for i_, el_47 := range arr_45 {
			out_46[i_] = j2gort.CastPtr[Color](el_47)
		}
		return out_46
	}()
}

// Java static method: Color.mergeStringArrays(java.lang.String[] a, java.lang.String[] b)
func Color_MergeStringArrays(a []string, b []string) []string {
	clinit_Color()
	_ = a
	_ = b
	var aLen int32 = int32(len(a))
	_ = aLen
	if aLen == int32(0) {
		return b
	}
	var bLen int32 = int32(len(b))
	_ = bLen
	if (bLen == int32(0)) || j2gort.SliceRefEquals(a, b) {
		return a
	}
	if (aLen == int32(1)) && (bLen == int32(1)) {
		{
			var iifeCond_1 bool
			{

				sr_48 := a[int32(0)]
				se_49 := b[int32(0)]
				iifeCond_1 = j2gort.CheckNotNullString(sr_48) == se_49
			}
			if iifeCond_1 {
				return a
			} else {
				return []string{a[int32(0)], b[int32(0)]}
			}
		}
	}
	if (aLen + bLen) <= int32(8) {
		var extra int32 = int32(0)
		_ = extra
		{
			var i int32 = int32(0)
			_ = i
			for ; i < bLen; i = i + 1 {
				if !(Color_ContainsString(a, aLen, b[i])) {
					extra = extra + 1
				}
			}
		}
		if extra == int32(0) {
			return a
		}
		var merged []string = j2gort.NewStringArray((aLen + extra))
		_ = merged
		func() {
			acSrc_50 := a
			acSP_51 := int32(0)
			acDst_52 := merged
			acDP_53 := int32(0)
			acLen_54 := aLen
			j2gort.CheckNotNullSlice(acSrc_50)
			j2gort.CheckNotNullSlice(acDst_52)
			if acSP_51 < 0 || acDP_53 < 0 || acLen_54 < 0 || int64(acSP_51)+int64(acLen_54) > int64(len(acSrc_50)) || int64(acDP_53)+int64(acLen_54) > int64(len(acDst_52)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_52[int(acDP_53):int(acDP_53)+int(acLen_54)], acSrc_50[int(acSP_51):int(acSP_51)+int(acLen_54)])
		}()
		var pos int32 = aLen
		_ = pos
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < bLen; i_2 = i_2 + 1 {
				var s string = b[i_2]
				_ = s
				if !(Color_ContainsString(a, aLen, s)) {
					{
						a_55 := merged
						i_56 := func() int32 { old_57 := pos; nw_58 := old_57 + 1; pos = nw_58; return old_57 }()
						a_55[i_56] = s
					}
				}
			}
		}
		return merged
	}
	var set *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32((aLen + bLen))
	_ = set
	{
		var i_3 int32 = int32(0)
		_ = i_3
		for ; i_3 < aLen; i_3 = i_3 + 1 {
			j2gort.SetAddStr(set, a[i_3])
		}
	}
	var addedNew bool = false
	_ = addedNew
	{
		var i_4 int32 = int32(0)
		_ = i_4
		for ; i_4 < bLen; i_4 = i_4 + 1 {
			if j2gort.SetAddStr(set, b[i_4]) {
				addedNew = true
			}
		}
	}
	if !(addedNew) {
		return a
	}
	return func() []string {
		arr_62 := set.ToArray_1_slice_any(func() []any {
			arr_59 := StaticGet_Color_EMPTY_STRINGS()
			if arr_59 == nil {
				return nil
			}
			out_60 := make([]any, len(arr_59))
			for i_, el_61 := range arr_59 {
				out_60[i_] = j2gort.BoxString(el_61)
			}
			return out_60
		}())
		if arr_62 == nil {
			return nil
		}
		out_63 := make([]string, len(arr_62))
		for i_, el_64 := range arr_62 {
			out_63[i_] = j2gort.CastString(el_64)
		}
		return out_63
	}()
}

// Java static method: Color.singleBuilder()
func Color_SingleBuilder() *Color_Builder {
	clinit_Color()
	var val_65 Color_Builder
	InitColor_Builder(&val_65)
	var b *Color_Builder = &val_65
	_ = b
	{
		r_66 := b
		r_66.closureAssert = false
	}
	{
		r_67 := b
		r_67.constructor = false
	}
	{
		r_68 := b
		r_68.instanceColors = StaticGet_Color_EMPTY_COLORS()
	}
	{
		r_69 := b
		r_69.invalidating = false
	}
	{
		r_70 := b
		r_70.ownProperties = StaticGet_Color_EMPTY_STRINGS()
	}
	{
		r_71 := b
		r_71.propertiesKeepOriginalName = false
	}
	{
		r_72 := b
		r_72.prototypes = StaticGet_Color_EMPTY_COLORS()
	}
	{
		r_73 := b
		r_73.unionElements = StaticGet_Color_EMPTY_COLORS()
	}
	{
		r_74 := b
		r_74.set_d0 = int8(int32(15))
	}
	return b
}

// Java static method: Color.createUnion2(com.google.javascript.jscomp.colors.Color first, com.google.javascript.jscomp.colors.Color second)
func Color_CreateUnion2(first *Color, second *Color) *Color {
	clinit_Color()
	_ = first
	_ = second
	j2gort.CheckNotNullPtr(first)
	j2gort.CheckNotNullPtr(second)
	if first.Equals(j2gort.PtrToAny(second)) {
		return first
	}
	if !(first.IsUnion()) && !(second.IsUnion()) {
		return NewColor(ColorId_Union2(first.id, second.id), Color_MergeColorArrays(first.prototypes, second.prototypes), Color_MergeColorArrays(first.instanceColors, second.instanceColors), (first.invalidating || second.invalidating), (first.propertiesKeepOriginalName || second.propertiesKeepOriginalName), (first.constructor && second.constructor), Color_MergeStringArrays(first.ownProperties, second.ownProperties), (ColorId{}), (first.closureAssert && second.closureAssert), []*Color{first, second})
	}
	return Color_CreateUnionFromArrayInternal([]*Color{first, second})
}

// Java static method: Color.createUnion(com.google.javascript.jscomp.colors.Color[] elements)
func Color_CreateUnion_1_Color(elements []*Color) *Color {
	clinit_Color()
	_ = elements
	var size int32 = int32(len(elements))
	_ = size
	if size == int32(0) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
	}
	if size == int32(1) {
		return j2gort.CheckNotNullPtr(elements[int32(0)])
	}
	if size == int32(2) {
		return Color_CreateUnion2(elements[int32(0)], elements[int32(1)])
	}
	return Color_CreateUnionFromArrayInternal(elements)
}

// Java static method: Color.createUnion(java.util.Set<com.google.javascript.jscomp.colors.Color> elements)
func Color_CreateUnion_1_Ext_java_util_Set(elements j2gort.Ext_java_util_Set) *Color {
	clinit_Color()
	_ = elements
	var size int32 = elements.Size()
	_ = size
	if size == int32(0) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
	}
	if size == int32(1) {
		return j2gort.CastPtr[Color](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](elements)))
	}
	if size == int32(2) {
		var it j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](elements.Iterator())
		_ = it
		var first *Color = j2gort.CastPtr[Color](it.Next())
		_ = first
		var second *Color = j2gort.CastPtr[Color](it.Next())
		_ = second
		return Color_CreateUnion2(first, second)
	}
	var singleInstanceColors []*Color = StaticGet_Color_EMPTY_COLORS()
	_ = singleInstanceColors
	var instanceColorsBuilder *j2gort.Ext_java_util_LinkedHashSet = nil
	_ = instanceColorsBuilder
	var singlePrototypes []*Color = StaticGet_Color_EMPTY_COLORS()
	_ = singlePrototypes
	var prototypesBuilder *j2gort.Ext_java_util_LinkedHashSet = nil
	_ = prototypesBuilder
	var singleOwnProperties []string = StaticGet_Color_EMPTY_STRINGS()
	_ = singleOwnProperties
	var ownPropertiesBuilder *j2gort.Ext_java_util_LinkedHashSet = nil
	_ = ownPropertiesBuilder
	var newElements *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(size)
	_ = newElements
	var isClosureAssert bool = true
	_ = isClosureAssert
	var isConstructor bool = true
	_ = isConstructor
	var isInvalidating bool = false
	_ = isInvalidating
	var propertiesKeepOriginalName bool = false
	_ = propertiesKeepOriginalName
	{
		cur_75 := j2gort.NewIterCursor(j2gort.CheckNotNull(elements))
		for cur_75.HasNext() {
			rawElem_76 := cur_75.Next()
			element := j2gort.CastPtr[Color](rawElem_76)
			_ = element
			if element.IsUnion() {
				for _, el_77 := range j2gort.CheckNotNullSlice(element.unionElements) {
					nestedElement := el_77
					_ = nestedElement
					newElements.Add_1(j2gort.PtrToAny(nestedElement))
				}
			} else {
				newElements.Add_1(j2gort.PtrToAny(element))
			}
			var elemInstanceColors []*Color = element.instanceColors
			_ = elemInstanceColors
			if int32(len(elemInstanceColors)) > int32(0) {
				if instanceColorsBuilder != nil {
					for _, el_78 := range j2gort.CheckNotNullSlice(elemInstanceColors) {
						c := el_78
						_ = c
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c))
					}
				} else if int32(len(singleInstanceColors)) == int32(0) {
					singleInstanceColors = elemInstanceColors
				} else if !(Color_SetEqualsColors(singleInstanceColors, elemInstanceColors)) {
					instanceColorsBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singleInstanceColors)) + int32(len(elemInstanceColors))))
					for _, el_79 := range j2gort.CheckNotNullSlice(singleInstanceColors) {
						c_2 := el_79
						_ = c_2
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c_2))
					}
					for _, el_80 := range j2gort.CheckNotNullSlice(elemInstanceColors) {
						c_3 := el_80
						_ = c_3
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c_3))
					}
				}
			}
			isClosureAssert = func() bool { ba_81 := isClosureAssert; bb_82 := element.closureAssert; return ba_81 && bb_82 }()
			isConstructor = func() bool { ba_83 := isConstructor; bb_84 := element.constructor; return ba_83 && bb_84 }()
			isInvalidating = func() bool { ba_85 := isInvalidating; bb_86 := element.invalidating; return ba_85 || bb_86 }()
			var elemOwnProps []string = element.ownProperties
			_ = elemOwnProps
			if int32(len(elemOwnProps)) > int32(0) {
				if ownPropertiesBuilder != nil {
					for _, el_87 := range j2gort.CheckNotNullSlice(elemOwnProps) {
						p := el_87
						_ = p
						j2gort.SetAddStr(ownPropertiesBuilder, p)
					}
				} else if int32(len(singleOwnProperties)) == int32(0) {
					singleOwnProperties = elemOwnProps
				} else if !(Color_SetEqualsStrings(singleOwnProperties, elemOwnProps)) {
					ownPropertiesBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singleOwnProperties)) + int32(len(elemOwnProps))))
					for _, el_88 := range j2gort.CheckNotNullSlice(singleOwnProperties) {
						p_2 := el_88
						_ = p_2
						j2gort.SetAddStr(ownPropertiesBuilder, p_2)
					}
					for _, el_89 := range j2gort.CheckNotNullSlice(elemOwnProps) {
						p_3 := el_89
						_ = p_3
						j2gort.SetAddStr(ownPropertiesBuilder, p_3)
					}
				}
			}
			propertiesKeepOriginalName = func() bool {
				ba_90 := propertiesKeepOriginalName
				bb_91 := element.propertiesKeepOriginalName
				return ba_90 || bb_91
			}()
			var elemPrototypes []*Color = element.prototypes
			_ = elemPrototypes
			if int32(len(elemPrototypes)) > int32(0) {
				if prototypesBuilder != nil {
					for _, el_92 := range j2gort.CheckNotNullSlice(elemPrototypes) {
						c_4 := el_92
						_ = c_4
						prototypesBuilder.Add_1(j2gort.PtrToAny(c_4))
					}
				} else if int32(len(singlePrototypes)) == int32(0) {
					singlePrototypes = elemPrototypes
				} else if !(Color_SetEqualsColors(singlePrototypes, elemPrototypes)) {
					prototypesBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singlePrototypes)) + int32(len(elemPrototypes))))
					for _, el_93 := range j2gort.CheckNotNullSlice(singlePrototypes) {
						c_5 := el_93
						_ = c_5
						prototypesBuilder.Add_1(j2gort.PtrToAny(c_5))
					}
					for _, el_94 := range j2gort.CheckNotNullSlice(elemPrototypes) {
						c_6 := el_94
						_ = c_6
						prototypesBuilder.Add_1(j2gort.PtrToAny(c_6))
					}
				}
			}
		}
	}
	if newElements.Size() == int32(1) {
		return j2gort.CastPtr[Color](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(newElements))
	}
	var builtUnionElements []*Color = func() []*Color {
		arr_98 := newElements.ToArray_1_slice_any(func() []any {
			arr_95 := StaticGet_Color_EMPTY_COLORS()
			if arr_95 == nil {
				return nil
			}
			out_96 := make([]any, len(arr_95))
			for i_, el_97 := range arr_95 {
				out_96[i_] = j2gort.PtrToAny(el_97)
			}
			return out_96
		}())
		if arr_98 == nil {
			return nil
		}
		out_99 := make([]*Color, len(arr_98))
		for i_, el_100 := range arr_98 {
			out_99[i_] = j2gort.CastPtr[Color](el_100)
		}
		return out_99
	}()
	_ = builtUnionElements
	var result *Color = NewColor(ColorId_UnionFromColors(builtUnionElements), func() []*Color {
		if prototypesBuilder != nil {
			return func() []*Color {
				arr_104 := prototypesBuilder.ToArray_1_slice_any(func() []any {
					arr_101 := StaticGet_Color_EMPTY_COLORS()
					if arr_101 == nil {
						return nil
					}
					out_102 := make([]any, len(arr_101))
					for i_, el_103 := range arr_101 {
						out_102[i_] = j2gort.PtrToAny(el_103)
					}
					return out_102
				}())
				if arr_104 == nil {
					return nil
				}
				out_105 := make([]*Color, len(arr_104))
				for i_, el_106 := range arr_104 {
					out_105[i_] = j2gort.CastPtr[Color](el_106)
				}
				return out_105
			}()
		}
		return singlePrototypes
	}(), func() []*Color {
		if instanceColorsBuilder != nil {
			return func() []*Color {
				arr_110 := instanceColorsBuilder.ToArray_1_slice_any(func() []any {
					arr_107 := StaticGet_Color_EMPTY_COLORS()
					if arr_107 == nil {
						return nil
					}
					out_108 := make([]any, len(arr_107))
					for i_, el_109 := range arr_107 {
						out_108[i_] = j2gort.PtrToAny(el_109)
					}
					return out_108
				}())
				if arr_110 == nil {
					return nil
				}
				out_111 := make([]*Color, len(arr_110))
				for i_, el_112 := range arr_110 {
					out_111[i_] = j2gort.CastPtr[Color](el_112)
				}
				return out_111
			}()
		}
		return singleInstanceColors
	}(), isInvalidating, propertiesKeepOriginalName, isConstructor, func() []string {
		if ownPropertiesBuilder != nil {
			return func() []string {
				arr_116 := ownPropertiesBuilder.ToArray_1_slice_any(func() []any {
					arr_113 := StaticGet_Color_EMPTY_STRINGS()
					if arr_113 == nil {
						return nil
					}
					out_114 := make([]any, len(arr_113))
					for i_, el_115 := range arr_113 {
						out_114[i_] = j2gort.BoxString(el_115)
					}
					return out_114
				}())
				if arr_116 == nil {
					return nil
				}
				out_117 := make([]string, len(arr_116))
				for i_, el_118 := range arr_116 {
					out_117[i_] = j2gort.CastString(el_118)
				}
				return out_117
			}()
		}
		return singleOwnProperties
	}(), (ColorId{}), isClosureAssert, builtUnionElements)
	_ = result
	if !(int32(len(result.unionElements)) > int32(1)) {
		j2gort.CheckStateMsg(false, result)
	}
	if !(result.boxId == (ColorId{})) {
		j2gort.CheckStateMsg(false, result)
	}
	return result
}

// Java static method: Color.createUnionFromArrayInternal(com.google.javascript.jscomp.colors.Color[] elements)
func Color_CreateUnionFromArrayInternal(elements []*Color) *Color {
	clinit_Color()
	_ = elements
	var size int32 = int32(len(elements))
	_ = size
	var singleInstanceColors []*Color = StaticGet_Color_EMPTY_COLORS()
	_ = singleInstanceColors
	var instanceColorsBuilder *j2gort.Ext_java_util_LinkedHashSet = nil
	_ = instanceColorsBuilder
	var singlePrototypes []*Color = StaticGet_Color_EMPTY_COLORS()
	_ = singlePrototypes
	var prototypesBuilder *j2gort.Ext_java_util_LinkedHashSet = nil
	_ = prototypesBuilder
	var singleOwnProperties []string = StaticGet_Color_EMPTY_STRINGS()
	_ = singleOwnProperties
	var ownPropertiesBuilder *j2gort.Ext_java_util_LinkedHashSet = nil
	_ = ownPropertiesBuilder
	var newElements *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(size)
	_ = newElements
	var isClosureAssert bool = true
	_ = isClosureAssert
	var isConstructor bool = true
	_ = isConstructor
	var isInvalidating bool = false
	_ = isInvalidating
	var propertiesKeepOriginalName bool = false
	_ = propertiesKeepOriginalName
	{
		var idx int32 = int32(0)
		_ = idx
		for ; idx < size; idx = idx + 1 {
			var element *Color = j2gort.CheckNotNullPtr(elements[idx])
			_ = element
			if element.IsUnion() {
				for _, el_119 := range j2gort.CheckNotNullSlice(element.unionElements) {
					nestedElement := el_119
					_ = nestedElement
					newElements.Add_1(j2gort.PtrToAny(nestedElement))
				}
			} else {
				newElements.Add_1(element)
			}
			var elemInstanceColors []*Color = element.instanceColors
			_ = elemInstanceColors
			if int32(len(elemInstanceColors)) > int32(0) {
				if instanceColorsBuilder != nil {
					for _, el_120 := range j2gort.CheckNotNullSlice(elemInstanceColors) {
						c := el_120
						_ = c
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c))
					}
				} else if int32(len(singleInstanceColors)) == int32(0) {
					singleInstanceColors = elemInstanceColors
				} else if !(Color_SetEqualsColors(singleInstanceColors, elemInstanceColors)) {
					instanceColorsBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singleInstanceColors)) + int32(len(elemInstanceColors))))
					for _, el_121 := range j2gort.CheckNotNullSlice(singleInstanceColors) {
						c_2 := el_121
						_ = c_2
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c_2))
					}
					for _, el_122 := range j2gort.CheckNotNullSlice(elemInstanceColors) {
						c_3 := el_122
						_ = c_3
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c_3))
					}
				}
			}
			isClosureAssert = func() bool { ba_123 := isClosureAssert; bb_124 := element.closureAssert; return ba_123 && bb_124 }()
			isConstructor = func() bool { ba_125 := isConstructor; bb_126 := element.constructor; return ba_125 && bb_126 }()
			isInvalidating = func() bool { ba_127 := isInvalidating; bb_128 := element.invalidating; return ba_127 || bb_128 }()
			var elemOwnProps []string = element.ownProperties
			_ = elemOwnProps
			if int32(len(elemOwnProps)) > int32(0) {
				if ownPropertiesBuilder != nil {
					for _, el_129 := range j2gort.CheckNotNullSlice(elemOwnProps) {
						p := el_129
						_ = p
						j2gort.SetAddStr(ownPropertiesBuilder, p)
					}
				} else if int32(len(singleOwnProperties)) == int32(0) {
					singleOwnProperties = elemOwnProps
				} else if !(Color_SetEqualsStrings(singleOwnProperties, elemOwnProps)) {
					ownPropertiesBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singleOwnProperties)) + int32(len(elemOwnProps))))
					for _, el_130 := range j2gort.CheckNotNullSlice(singleOwnProperties) {
						p_2 := el_130
						_ = p_2
						j2gort.SetAddStr(ownPropertiesBuilder, p_2)
					}
					for _, el_131 := range j2gort.CheckNotNullSlice(elemOwnProps) {
						p_3 := el_131
						_ = p_3
						j2gort.SetAddStr(ownPropertiesBuilder, p_3)
					}
				}
			}
			propertiesKeepOriginalName = func() bool {
				ba_132 := propertiesKeepOriginalName
				bb_133 := element.propertiesKeepOriginalName
				return ba_132 || bb_133
			}()
			var elemPrototypes []*Color = element.prototypes
			_ = elemPrototypes
			if int32(len(elemPrototypes)) > int32(0) {
				if prototypesBuilder != nil {
					for _, el_134 := range j2gort.CheckNotNullSlice(elemPrototypes) {
						c_4 := el_134
						_ = c_4
						prototypesBuilder.Add_1(j2gort.PtrToAny(c_4))
					}
				} else if int32(len(singlePrototypes)) == int32(0) {
					singlePrototypes = elemPrototypes
				} else if !(Color_SetEqualsColors(singlePrototypes, elemPrototypes)) {
					prototypesBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singlePrototypes)) + int32(len(elemPrototypes))))
					for _, el_135 := range j2gort.CheckNotNullSlice(singlePrototypes) {
						c_5 := el_135
						_ = c_5
						prototypesBuilder.Add_1(j2gort.PtrToAny(c_5))
					}
					for _, el_136 := range j2gort.CheckNotNullSlice(elemPrototypes) {
						c_6 := el_136
						_ = c_6
						prototypesBuilder.Add_1(j2gort.PtrToAny(c_6))
					}
				}
			}
		}
	}
	if newElements.Size() == int32(1) {
		return j2gort.CastPtr[Color](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(newElements))
	}
	var builtUnionElements []*Color = func() []*Color {
		arr_140 := newElements.ToArray_1_slice_any(func() []any {
			arr_137 := StaticGet_Color_EMPTY_COLORS()
			if arr_137 == nil {
				return nil
			}
			out_138 := make([]any, len(arr_137))
			for i_, el_139 := range arr_137 {
				out_138[i_] = j2gort.PtrToAny(el_139)
			}
			return out_138
		}())
		if arr_140 == nil {
			return nil
		}
		out_141 := make([]*Color, len(arr_140))
		for i_, el_142 := range arr_140 {
			out_141[i_] = j2gort.CastPtr[Color](el_142)
		}
		return out_141
	}()
	_ = builtUnionElements
	var result *Color = NewColor(ColorId_UnionFromColors(builtUnionElements), func() []*Color {
		if prototypesBuilder != nil {
			return func() []*Color {
				arr_146 := prototypesBuilder.ToArray_1_slice_any(func() []any {
					arr_143 := StaticGet_Color_EMPTY_COLORS()
					if arr_143 == nil {
						return nil
					}
					out_144 := make([]any, len(arr_143))
					for i_, el_145 := range arr_143 {
						out_144[i_] = j2gort.PtrToAny(el_145)
					}
					return out_144
				}())
				if arr_146 == nil {
					return nil
				}
				out_147 := make([]*Color, len(arr_146))
				for i_, el_148 := range arr_146 {
					out_147[i_] = j2gort.CastPtr[Color](el_148)
				}
				return out_147
			}()
		}
		return singlePrototypes
	}(), func() []*Color {
		if instanceColorsBuilder != nil {
			return func() []*Color {
				arr_152 := instanceColorsBuilder.ToArray_1_slice_any(func() []any {
					arr_149 := StaticGet_Color_EMPTY_COLORS()
					if arr_149 == nil {
						return nil
					}
					out_150 := make([]any, len(arr_149))
					for i_, el_151 := range arr_149 {
						out_150[i_] = j2gort.PtrToAny(el_151)
					}
					return out_150
				}())
				if arr_152 == nil {
					return nil
				}
				out_153 := make([]*Color, len(arr_152))
				for i_, el_154 := range arr_152 {
					out_153[i_] = j2gort.CastPtr[Color](el_154)
				}
				return out_153
			}()
		}
		return singleInstanceColors
	}(), isInvalidating, propertiesKeepOriginalName, isConstructor, func() []string {
		if ownPropertiesBuilder != nil {
			return func() []string {
				arr_158 := ownPropertiesBuilder.ToArray_1_slice_any(func() []any {
					arr_155 := StaticGet_Color_EMPTY_STRINGS()
					if arr_155 == nil {
						return nil
					}
					out_156 := make([]any, len(arr_155))
					for i_, el_157 := range arr_155 {
						out_156[i_] = j2gort.BoxString(el_157)
					}
					return out_156
				}())
				if arr_158 == nil {
					return nil
				}
				out_159 := make([]string, len(arr_158))
				for i_, el_160 := range arr_158 {
					out_159[i_] = j2gort.CastString(el_160)
				}
				return out_159
			}()
		}
		return singleOwnProperties
	}(), (ColorId{}), isClosureAssert, builtUnionElements)
	_ = result
	if !(int32(len(result.unionElements)) > int32(1)) {
		j2gort.CheckStateMsg(false, result)
	}
	if !(result.boxId == (ColorId{})) {
		j2gort.CheckStateMsg(false, result)
	}
	return result
}

// Java method: Color.isPrimitive()
func (this_ *Color) IsPrimitive() bool {
	if !(!(this_.IsUnion())) {
		j2gort.CheckStateMsg(false, this_)
	}
	return j2gort.MapContainsEqVal(StaticGet_StandardColors_PRIMITIVE_COLORS(), this_.id)
}

// Java method: Color.isUnion()
func (this_ *Color) IsUnion() bool {
	return (int32(len(this_.unionElements)) > int32(0))
}

// Java method: Color.subtractNullOrVoid()
func (this_ *Color) SubtractNullOrVoid() *Color {
	var cached *Color = this_.subtractNullOrVoid
	_ = cached
	if cached == nil {
		return this_.priv_Color_computeSubtractNullOrVoid_5967()
	}
	return cached
}

// Java method: Color.computeSubtractNullOrVoid()
func (this_ *Color) priv_Color_computeSubtractNullOrVoid_5967() *Color {
	{
		if this_ == nil {
			j2gort.ThrowNPE()
		}
		this_._syncMu.Lock()
		defer this_._syncMu.Unlock()
		var cached *Color = this_.subtractNullOrVoid
		_ = cached
		if cached == nil {
			if !(this_.IsUnion()) {
				j2gort.CheckState(false)
			}
			var elems []*Color = this_.unionElements
			_ = elems
			var len_ int32 = int32(len(elems))
			_ = len_
			var nullOrVoidIndex int32 = int32(-1)
			_ = nullOrVoidIndex
			{
				var i int32 = int32(0)
				_ = i
				for ; i < len_; i = i + 1 {
					if elems[i].Equals(j2gort.PtrToAny(StaticGet_StandardColors_NULL_OR_VOID())) {
						nullOrVoidIndex = i
						break
					}
				}
			}
			if nullOrVoidIndex < int32(0) {
				cached = this_
			} else if len_ == int32(2) {
				cached = elems[(int32(1) - nullOrVoidIndex)]
			} else if len_ == int32(3) {
				var first *Color
				_ = first
				if nullOrVoidIndex == int32(0) {
					first = elems[int32(1)]
				} else {
					first = elems[int32(0)]
				}
				var second *Color
				_ = second
				if nullOrVoidIndex == int32(2) {
					second = elems[int32(1)]
				} else {
					second = elems[int32(2)]
				}
				cached = Color_CreateUnion2(first, second)
			} else {
				var remaining []*Color = make([]*Color, (len_ - int32(1)))
				_ = remaining
				if nullOrVoidIndex > int32(0) {
					func() {
						acSrc_161 := elems
						acSP_162 := int32(0)
						acDst_163 := remaining
						acDP_164 := int32(0)
						acLen_165 := nullOrVoidIndex
						j2gort.CheckNotNullSlice(acSrc_161)
						j2gort.CheckNotNullSlice(acDst_163)
						if acSP_162 < 0 || acDP_164 < 0 || acLen_165 < 0 || int64(acSP_162)+int64(acLen_165) > int64(len(acSrc_161)) || int64(acDP_164)+int64(acLen_165) > int64(len(acDst_163)) {
							panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
						}
						copy(acDst_163[int(acDP_164):int(acDP_164)+int(acLen_165)], acSrc_161[int(acSP_162):int(acSP_162)+int(acLen_165)])
					}()
				}
				if nullOrVoidIndex < (len_ - int32(1)) {
					func() {
						acSrc_166 := elems
						acSP_167 := (nullOrVoidIndex + int32(1))
						acDst_168 := remaining
						acDP_169 := nullOrVoidIndex
						acLen_170 := ((len_ - nullOrVoidIndex) - int32(1))
						j2gort.CheckNotNullSlice(acSrc_166)
						j2gort.CheckNotNullSlice(acDst_168)
						if acSP_167 < 0 || acDP_169 < 0 || acLen_170 < 0 || int64(acSP_167)+int64(acLen_170) > int64(len(acSrc_166)) || int64(acDP_169)+int64(acLen_170) > int64(len(acDst_168)) {
							panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
						}
						copy(acDst_168[int(acDP_169):int(acDP_169)+int(acLen_170)], acSrc_166[int(acSP_167):int(acSP_167)+int(acLen_170)])
					}()
				}
				cached = Color_CreateUnionFromArrayInternal(remaining)
			}
			if cached == nil {
				panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "subtractNullOrVoid() cannot return null")))
			}
			this_.subtractNullOrVoid = cached
		}
		return cached
	}
	panic("unreachable")
}

// Java method: Color.toString()
func (this_ *Color) String() string {
	return ("Color{id=" + j2gort.StrVal(this_.id) + ", prototypes=" + j2gort.StrValString(j2gort.SliceFormatJavaString(this_.prototypes)) + ", instanceColors=" + j2gort.StrValString(j2gort.SliceFormatJavaString(this_.instanceColors)) + ", invalidating=" + j2gort.StrValBool(this_.invalidating) + ", propertiesKeepOriginalName=" + j2gort.StrValBool(this_.propertiesKeepOriginalName) + ", constructor=" + j2gort.StrValBool(this_.constructor) + ", ownProperties=" + j2gort.StrValString(j2gort.SliceFormatJavaString(this_.ownProperties)) + ", boxId=" + j2gort.StrVal(this_.boxId) + ", closureAssert=" + j2gort.StrValBool(this_.closureAssert) + ", unionElements=" + j2gort.StrValString(j2gort.SliceFormatJavaString(this_.unionElements)) + "}")
}

// Java static method: Color.setEqualsColors(com.google.javascript.jscomp.colors.Color[] a, com.google.javascript.jscomp.colors.Color[] b)
func Color_SetEqualsColors(a []*Color, b []*Color) bool {
	clinit_Color()
	_ = a
	_ = b
	if j2gort.SliceRefEquals(a, b) {
		return true
	}
	var len_ int32 = int32(len(a))
	_ = len_
	if len_ != int32(len(b)) {
		return false
	}
	if len_ == int32(1) {
		return a[int32(0)].Equals(j2gort.PtrToAny(b[int32(0)]))
	}
	if len_ == int32(2) {
		return ((a[int32(0)].Equals(j2gort.PtrToAny(b[int32(0)])) && a[int32(1)].Equals(j2gort.PtrToAny(b[int32(1)]))) || (a[int32(0)].Equals(j2gort.PtrToAny(b[int32(1)])) && a[int32(1)].Equals(j2gort.PtrToAny(b[int32(0)]))))
	}
	var sameOrder bool = true
	_ = sameOrder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			if !(a[i].Equals(j2gort.PtrToAny(b[i]))) {
				sameOrder = false
				break
			}
		}
	}
	if sameOrder {
		return true
	}
	if len_ <= int32(16) {
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < len_; i_2 = i_2 + 1 {
				if !(Color_ContainsColor(b, len_, a[i_2])) {
					return false
				}
			}
		}
		return true
	}
	var bSet j2gort.CustomOrderedSet[*Color] = j2gort.NewCustomOrderedSetCap[*Color]((len_ * int32(2)))
	_ = bSet
	{
		var i_3 int32 = int32(0)
		_ = i_3
		for ; i_3 < len_; i_3 = i_3 + 1 {
			bSet.Add(b[i_3])
		}
	}
	{
		var i_4 int32 = int32(0)
		_ = i_4
		for ; i_4 < len_; i_4 = i_4 + 1 {
			if !(bSet.Contains(a[i_4])) {
				return false
			}
		}
	}
	return true
}

// Java static method: Color.setEqualsStrings(java.lang.String[] a, java.lang.String[] b)
func Color_SetEqualsStrings(a []string, b []string) bool {
	clinit_Color()
	_ = a
	_ = b
	if j2gort.SliceRefEquals(a, b) {
		return true
	}
	var len_ int32 = int32(len(a))
	_ = len_
	if len_ != int32(len(b)) {
		return false
	}
	if len_ == int32(1) {
		return func() bool {
			sr_171 := a[int32(0)]
			se_172 := b[int32(0)]
			return j2gort.CheckNotNullString(sr_171) == se_172
		}()
	}
	if len_ == int32(2) {
		return ((func() bool {
			sr_173 := a[int32(0)]
			se_174 := b[int32(0)]
			return j2gort.CheckNotNullString(sr_173) == se_174
		}() && func() bool {
			sr_175 := a[int32(1)]
			se_176 := b[int32(1)]
			return j2gort.CheckNotNullString(sr_175) == se_176
		}()) || (func() bool {
			sr_177 := a[int32(0)]
			se_178 := b[int32(1)]
			return j2gort.CheckNotNullString(sr_177) == se_178
		}() && func() bool {
			sr_179 := a[int32(1)]
			se_180 := b[int32(0)]
			return j2gort.CheckNotNullString(sr_179) == se_180
		}()))
	}
	var sameOrder bool = true
	_ = sameOrder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			{
				var iifeCond_1 bool
				{

					sr_181 := a[i]
					se_182 := b[i]
					iifeCond_1 = j2gort.CheckNotNullString(sr_181) == se_182
				}
				if !(iifeCond_1) {
					sameOrder = false
					break
				}
			}
		}
	}
	if sameOrder {
		return true
	}
	if len_ <= int32(16) {
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < len_; i_2 = i_2 + 1 {
				if !(Color_ContainsString(b, len_, a[i_2])) {
					return false
				}
			}
		}
		return true
	}
	var bSet j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSetCap[string]((len_ * int32(2)))
	_ = bSet
	{
		var i_3 int32 = int32(0)
		_ = i_3
		for ; i_3 < len_; i_3 = i_3 + 1 {
			bSet.Add(b[i_3])
		}
	}
	{
		var i_4 int32 = int32(0)
		_ = i_4
		for ; i_4 < len_; i_4 = i_4 + 1 {
			if !(bSet.Contains(a[i_4])) {
				return false
			}
		}
	}
	return true
}

// Java method: Color.equals(java.lang.@org.jspecify.annotations.Nullable Object o)
func (this_ *Color) Equals(arg_183 any) bool {
	o := arg_183
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if j2gort.InstanceOf[*Color](o) {
		var that *Color = j2gort.CastPtr[Color](o)
		_ = that
		return ((((((((((func() bool {
			rv_184 := this_.id
			av_185 := that.id
			_ = j2gort.CheckNotNullVal(rv_184)
			return rv_184 == av_185
		}() && (this_.invalidating == that.invalidating)) && (this_.propertiesKeepOriginalName == that.propertiesKeepOriginalName)) && (this_.constructor == that.constructor)) && (this_.closureAssert == that.closureAssert)) && (func() bool {
			if this_.boxId == (ColorId{}) {
				return (that.boxId == (ColorId{}))
			}
			return func() bool {
				rv_186 := this_.boxId
				av_187 := that.boxId
				_ = j2gort.CheckNotNullVal(rv_186)
				return rv_186 == av_187
			}()
		}())) && (this_.HashCode() == that.HashCode())) && Color_SetEqualsColors(this_.prototypes, that.prototypes)) && Color_SetEqualsColors(this_.instanceColors, that.instanceColors)) && Color_SetEqualsStrings(this_.ownProperties, that.ownProperties)) && Color_SetEqualsColors(this_.unionElements, that.unionElements))
	}
	return false
}

// Java method: Color.hashCode()
func (this_ *Color) HashCode() int32 {
	var h int32 = this_.hashCode
	_ = h
	if h == int32(0) {
		h = this_.priv_Color_computeHashCode_5971()
		this_.hashCode = h
	}
	return h
}

// Java method: Color.computeHashCode()
func (this_ *Color) priv_Color_computeHashCode_5971() int32 {
	var protoHash int32 = int32(0)
	_ = protoHash
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(this_.prototypes))
		_ = n
		for ; i < n; i = i + 1 {
			protoHash = (protoHash + this_.prototypes[i].HashCode())
		}
	}
	var instHash int32 = int32(0)
	_ = instHash
	{
		var i_2 int32 = int32(0)
		_ = i_2
		var n_2 int32 = int32(len(this_.instanceColors))
		_ = n_2
		for ; i_2 < n_2; i_2 = i_2 + 1 {
			instHash = (instHash + this_.instanceColors[i_2].HashCode())
		}
	}
	var ownPropHash int32 = int32(0)
	_ = ownPropHash
	{
		var i_3 int32 = int32(0)
		_ = i_3
		var n_3 int32 = int32(len(this_.ownProperties))
		_ = n_3
		for ; i_3 < n_3; i_3 = i_3 + 1 {
			ownPropHash = (ownPropHash + javalang.StringHashCode(j2gort.CheckNotNullString(this_.ownProperties[i_3])))
		}
	}
	var unionHash int32 = int32(0)
	_ = unionHash
	{
		var i_4 int32 = int32(0)
		_ = i_4
		var n_4 int32 = int32(len(this_.unionElements))
		_ = n_4
		for ; i_4 < n_4; i_4 = i_4 + 1 {
			unionHash = (unionHash + this_.unionElements[i_4].HashCode())
		}
	}
	var h_d int32 = int32(1)
	_ = h_d
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ this_.id.HashCode())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ protoHash)
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ instHash)
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.invalidating {
			return int32(1231)
		}
		return int32(1237)
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.propertiesKeepOriginalName {
			return int32(1231)
		}
		return int32(1237)
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.constructor {
			return int32(1231)
		}
		return int32(1237)
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ ownPropHash)
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.boxId == (ColorId{}) {
			return int32(0)
		}
		return this_.boxId.HashCode()
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ func() int32 {
		if this_.closureAssert {
			return int32(1231)
		}
		return int32(1237)
	}())
	h_d = (h_d * int32(1000003))
	h_d = (h_d ^ unionHash)
	return h_d
}

// Java class: com.google.javascript.jscomp.colors.Color.Builder
type Color_Builder struct {
	prototypes                 []*Color `gson:"prototypes"`
	instanceColors             []*Color `gson:"instanceColors"`
	ownProperties              []string `gson:"ownProperties"`
	unionElements              []*Color `gson:"unionElements"`
	id                         ColorId  `gson:"id"`
	boxId                      ColorId  `gson:"boxId"`
	idHash_                    atomic.Int32
	invalidating               bool `gson:"invalidating"`
	propertiesKeepOriginalName bool `gson:"propertiesKeepOriginalName"`
	constructor                bool `gson:"constructor"`
	closureAssert              bool `gson:"closureAssert"`
	set_d0                     int8 `gson:"set$0"`
}

func (this_ *Color_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *Color_Builder) IsJavaNil_() bool { return this_ == nil }
func (*Color_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.colors.Color$Builder", "Builder"
}
func (this_ *Color_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Color_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.colors.Color.Builder", "com.google.javascript.jscomp.colors.Color$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*Color_Builder](v_)
	})
}

func (this_ *Color_Builder) initFields_Color_Builder() {
	_ = this_
}

// Java constructor: Builder()
func NewColor_Builder() *Color_Builder {
	this_ := &Color_Builder{}
	InitColor_Builder(this_)
	return this_
}

func InitColor_Builder(this_ *Color_Builder) {
	this_.initFields_Color_Builder()
}

// Java method: Builder.setId(com.google.javascript.jscomp.colors.ColorId id)
func (this_ *Color_Builder) SetId(id ColorId) *Color_Builder {
	_ = id
	if id == (ColorId{}) {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null id")))
	}
	this_.id = id
	return this_
}

// Java method: Builder.setPrototypes(com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.colors.Color> prototypes)
func (this_ *Color_Builder) SetPrototypes(prototypes j2gort.Ext_com_google_common_collect_ImmutableSet) *Color_Builder {
	_ = prototypes
	if prototypes == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null prototypes")))
	}
	this_.prototypes = Color_ToUniqueColorArray(j2gort.AsIface[j2gort.Ext_java_util_Collection](prototypes))
	return this_
}

// Java method: Builder.setPrototypesArray(com.google.javascript.jscomp.colors.Color[] prototypes)
func (this_ *Color_Builder) SetPrototypesArray(prototypes []*Color) *Color_Builder {
	_ = prototypes
	if prototypes == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null prototypes")))
	}
	this_.prototypes = prototypes
	return this_
}

// Java method: Builder.setInstanceColors(com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.colors.Color> instanceColors)
func (this_ *Color_Builder) SetInstanceColors(instanceColors j2gort.Ext_com_google_common_collect_ImmutableSet) *Color_Builder {
	_ = instanceColors
	if instanceColors == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null instanceColors")))
	}
	this_.instanceColors = Color_ToUniqueColorArray(j2gort.AsIface[j2gort.Ext_java_util_Collection](instanceColors))
	return this_
}

// Java method: Builder.setInstanceColorsArray(com.google.javascript.jscomp.colors.Color[] instanceColors)
func (this_ *Color_Builder) SetInstanceColorsArray(instanceColors []*Color) *Color_Builder {
	_ = instanceColors
	if instanceColors == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null instanceColors")))
	}
	this_.instanceColors = instanceColors
	return this_
}

// Java method: Builder.setInvalidating(boolean invalidating)
func (this_ *Color_Builder) SetInvalidating(invalidating bool) *Color_Builder {
	_ = invalidating
	this_.invalidating = invalidating
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(1))))
	return this_
}

// Java method: Builder.setPropertiesKeepOriginalName(boolean propertiesKeepOriginalName)
func (this_ *Color_Builder) SetPropertiesKeepOriginalName(propertiesKeepOriginalName bool) *Color_Builder {
	_ = propertiesKeepOriginalName
	this_.propertiesKeepOriginalName = propertiesKeepOriginalName
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(2))))
	return this_
}

// Java method: Builder.setConstructor(boolean constructor)
func (this_ *Color_Builder) SetConstructor(constructor bool) *Color_Builder {
	_ = constructor
	this_.constructor = constructor
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(4))))
	return this_
}

// Java method: Builder.setOwnProperties(com.google.common.collect.ImmutableSet<java.lang.String> ownProperties)
func (this_ *Color_Builder) SetOwnProperties(ownProperties j2gort.Ext_com_google_common_collect_ImmutableSet) *Color_Builder {
	_ = ownProperties
	if ownProperties == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null ownProperties")))
	}
	this_.ownProperties = Color_ToUniqueStringArray(j2gort.AsIface[j2gort.Ext_java_util_Collection](ownProperties))
	return this_
}

// Java method: Builder.setOwnPropertiesArray(java.lang.String[] ownProperties)
func (this_ *Color_Builder) SetOwnPropertiesArray(ownProperties []string) *Color_Builder {
	_ = ownProperties
	if ownProperties == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null ownProperties")))
	}
	this_.ownProperties = ownProperties
	return this_
}

// Java method: Builder.setBoxId(com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable ColorId boxId)
func (this_ *Color_Builder) SetBoxId(boxId ColorId) *Color_Builder {
	_ = boxId
	this_.boxId = boxId
	return this_
}

// Java method: Builder.setClosureAssert(boolean closureAssert)
func (this_ *Color_Builder) SetClosureAssert(closureAssert bool) *Color_Builder {
	_ = closureAssert
	this_.closureAssert = closureAssert
	this_.set_d0 = int8((int32(this_.set_d0) | int32(int8(8))))
	return this_
}

// Java method: Builder.setUnionElements(com.google.common.collect.ImmutableSet<com.google.javascript.jscomp.colors.Color> unionElements)
func (this_ *Color_Builder) SetUnionElements(unionElements j2gort.Ext_com_google_common_collect_ImmutableSet) *Color_Builder {
	_ = unionElements
	if unionElements == nil {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "Null unionElements")))
	}
	this_.unionElements = Color_ToUniqueColorArray(j2gort.AsIface[j2gort.Ext_java_util_Collection](unionElements))
	return this_
}

// Java method: Builder.setPrototype(com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable Color x)
func (this_ *Color_Builder) SetPrototype(x *Color) *Color_Builder {
	_ = x
	if x == nil {
		this_.prototypes = StaticGet_Color_EMPTY_COLORS()
	} else {
		this_.prototypes = []*Color{x}
	}
	return this_
}

// Java method: Builder.setInstanceColor(com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable Color x)
func (this_ *Color_Builder) SetInstanceColor(x *Color) *Color_Builder {
	_ = x
	if x == nil {
		this_.instanceColors = StaticGet_Color_EMPTY_COLORS()
	} else {
		this_.instanceColors = []*Color{x}
	}
	return this_
}

// Java method: Builder.buildInternal()
func (this_ *Color_Builder) BuildInternal() *Color {
	if (((((int32(this_.set_d0) != int32(15)) || (this_.id == (ColorId{}))) || (this_.prototypes == nil)) || (this_.instanceColors == nil)) || (this_.ownProperties == nil)) || (this_.unionElements == nil) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if this_.id == (ColorId{}) {
			missing.Append_1_string(" id")
		}
		if this_.prototypes == nil {
			missing.Append_1_string(" prototypes")
		}
		if this_.instanceColors == nil {
			missing.Append_1_string(" instanceColors")
		}
		if (int32(this_.set_d0) & int32(1)) == int32(0) {
			missing.Append_1_string(" invalidating")
		}
		if (int32(this_.set_d0) & int32(2)) == int32(0) {
			missing.Append_1_string(" propertiesKeepOriginalName")
		}
		if (int32(this_.set_d0) & int32(4)) == int32(0) {
			missing.Append_1_string(" constructor")
		}
		if this_.ownProperties == nil {
			missing.Append_1_string(" ownProperties")
		}
		if (int32(this_.set_d0) & int32(8)) == int32(0) {
			missing.Append_1_string(" closureAssert")
		}
		if this_.unionElements == nil {
			missing.Append_1_string(" unionElements")
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Missing required properties:"+j2gort.StrVal(missing))))
	}
	return NewColor(this_.id, this_.prototypes, this_.instanceColors, this_.invalidating, this_.propertiesKeepOriginalName, this_.constructor, this_.ownProperties, this_.boxId, this_.closureAssert, this_.unionElements)
}

// Java method: Builder.build()
func (this_ *Color_Builder) Build() *Color {
	var result *Color = this_.BuildInternal()
	_ = result
	if !(int32(len(result.unionElements)) == int32(0)) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(result))
	}
	if !(result.boxId == (ColorId{})) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(result))
	}
	if !(!(j2gort.MapContainsEqVal(StaticGet_StandardColors_AXIOMATIC_COLORS(), result.id))) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(result))
	}
	return result
}

// Java method: Builder.buildAxiomatic()
func (this_ *Color_Builder) BuildAxiomatic() *Color {
	if !(StaticGet_StandardColors_AXIOMATIC_COLORS() == nil) {
		j2gort.CheckStateStr(false, "StandardColors are all defined")
	}
	var result *Color = this_.BuildInternal()
	_ = result
	if !(int32(len(result.unionElements)) == int32(0)) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(result))
	}
	var unused int32 = result.HashCode()
	_ = unused
	return result
}
