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
	if (id).IsJavaNil_() {
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
					var i_10 int32
					{

						old_11 := i
						nw_12 := old_11 + 1
						i = nw_12
						i_10 = old_11
					}
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
	return j2gort.ColToPtrArray[Color](unique, StaticGet_Color_EMPTY_COLORS())
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
			cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(strings))
			for cur_15.HasNext() {
				s := cur_15.NextStr()
				_ = s
				j2gort.CheckNotNullString(s)
				{
					a_16 := out
					var i_17 int32
					{

						old_18 := i
						nw_19 := old_18 + 1
						i = nw_19
						i_17 = old_18
					}
					a_16[i_17] = s
				}
			}
		}
		return out
	}
	var unique *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(size)
	_ = unique
	{
		cur_20 := j2gort.NewIterCursor(j2gort.CheckNotNull(strings))
		for cur_20.HasNext() {
			s_2 := cur_20.NextStr()
			_ = s_2
			j2gort.CheckNotNullString(s_2)
			j2gort.SetAddStr(unique, s_2)
		}
	}
	return j2gort.ColToStringArray(unique, StaticGet_Color_EMPTY_STRINGS())
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
			acSrc_21 := a
			acSP_22 := int32(0)
			acDst_23 := merged
			acDP_24 := int32(0)
			acLen_25 := aLen
			j2gort.CheckNotNullSlice(acSrc_21)
			j2gort.CheckNotNullSlice(acDst_23)
			if acSP_22 < 0 || acDP_24 < 0 || acLen_25 < 0 || int64(acSP_22)+int64(acLen_25) > int64(len(acSrc_21)) || int64(acDP_24)+int64(acLen_25) > int64(len(acDst_23)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_23[int(acDP_24):int(acDP_24)+int(acLen_25)], acSrc_21[int(acSP_22):int(acSP_22)+int(acLen_25)])
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
						a_26 := merged
						var i_27 int32
						{

							old_28 := pos
							nw_29 := old_28 + 1
							pos = nw_29
							i_27 = old_28
						}
						a_26[i_27] = c
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
	return j2gort.ColToPtrArray[Color](set, StaticGet_Color_EMPTY_COLORS())
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

				sr_30 := a[int32(0)]
				se_31 := b[int32(0)]
				iifeCond_1 = j2gort.CheckNotNullString(sr_30) == se_31
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
			acSrc_32 := a
			acSP_33 := int32(0)
			acDst_34 := merged
			acDP_35 := int32(0)
			acLen_36 := aLen
			j2gort.CheckNotNullSlice(acSrc_32)
			j2gort.CheckNotNullSlice(acDst_34)
			if acSP_33 < 0 || acDP_35 < 0 || acLen_36 < 0 || int64(acSP_33)+int64(acLen_36) > int64(len(acSrc_32)) || int64(acDP_35)+int64(acLen_36) > int64(len(acDst_34)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_34[int(acDP_35):int(acDP_35)+int(acLen_36)], acSrc_32[int(acSP_33):int(acSP_33)+int(acLen_36)])
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
						a_37 := merged
						var i_38 int32
						{

							old_39 := pos
							nw_40 := old_39 + 1
							pos = nw_40
							i_38 = old_39
						}
						a_37[i_38] = s
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
	return j2gort.ColToStringArray(set, StaticGet_Color_EMPTY_STRINGS())
}

// Java static method: Color.singleBuilder()
func Color_SingleBuilder() *Color_Builder {
	clinit_Color()
	var val_41 Color_Builder
	InitColor_Builder(&val_41)
	var b *Color_Builder = &val_41
	_ = b
	{
		r_42 := b
		r_42.closureAssert = false
	}
	{
		r_43 := b
		r_43.constructor = false
	}
	{
		r_44 := b
		r_44.instanceColors = StaticGet_Color_EMPTY_COLORS()
	}
	{
		r_45 := b
		r_45.invalidating = false
	}
	{
		r_46 := b
		r_46.ownProperties = StaticGet_Color_EMPTY_STRINGS()
	}
	{
		r_47 := b
		r_47.propertiesKeepOriginalName = false
	}
	{
		r_48 := b
		r_48.prototypes = StaticGet_Color_EMPTY_COLORS()
	}
	{
		r_49 := b
		r_49.unionElements = StaticGet_Color_EMPTY_COLORS()
	}
	{
		r_50 := b
		r_50.set_d0 = int8(int32(15))
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
		cur_51 := j2gort.NewIterCursor(j2gort.CheckNotNull(elements))
		for cur_51.HasNext() {
			rawElem_52 := cur_51.Next()
			element := j2gort.CastPtr[Color](rawElem_52)
			_ = element
			if element.IsUnion() {
				for _, el_53 := range j2gort.CheckNotNullSlice(element.unionElements) {
					nestedElement := el_53
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
					for _, el_54 := range j2gort.CheckNotNullSlice(elemInstanceColors) {
						c := el_54
						_ = c
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c))
					}
				} else if int32(len(singleInstanceColors)) == int32(0) {
					singleInstanceColors = elemInstanceColors
				} else if !(Color_SetEqualsColors(singleInstanceColors, elemInstanceColors)) {
					instanceColorsBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singleInstanceColors)) + int32(len(elemInstanceColors))))
					for _, el_55 := range j2gort.CheckNotNullSlice(singleInstanceColors) {
						c_2 := el_55
						_ = c_2
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c_2))
					}
					for _, el_56 := range j2gort.CheckNotNullSlice(elemInstanceColors) {
						c_3 := el_56
						_ = c_3
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c_3))
					}
				}
			}
			{
				ba_57 := isClosureAssert
				bb_58 := element.closureAssert
				isClosureAssert = ba_57 && bb_58
			}
			{
				ba_59 := isConstructor
				bb_60 := element.constructor
				isConstructor = ba_59 && bb_60
			}
			{
				ba_61 := isInvalidating
				bb_62 := element.invalidating
				isInvalidating = ba_61 || bb_62
			}
			var elemOwnProps []string = element.ownProperties
			_ = elemOwnProps
			if int32(len(elemOwnProps)) > int32(0) {
				if ownPropertiesBuilder != nil {
					for _, el_63 := range j2gort.CheckNotNullSlice(elemOwnProps) {
						p := el_63
						_ = p
						j2gort.SetAddStr(ownPropertiesBuilder, p)
					}
				} else if int32(len(singleOwnProperties)) == int32(0) {
					singleOwnProperties = elemOwnProps
				} else if !(Color_SetEqualsStrings(singleOwnProperties, elemOwnProps)) {
					ownPropertiesBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singleOwnProperties)) + int32(len(elemOwnProps))))
					for _, el_64 := range j2gort.CheckNotNullSlice(singleOwnProperties) {
						p_2 := el_64
						_ = p_2
						j2gort.SetAddStr(ownPropertiesBuilder, p_2)
					}
					for _, el_65 := range j2gort.CheckNotNullSlice(elemOwnProps) {
						p_3 := el_65
						_ = p_3
						j2gort.SetAddStr(ownPropertiesBuilder, p_3)
					}
				}
			}
			propertiesKeepOriginalName = func() bool {
				ba_66 := propertiesKeepOriginalName
				bb_67 := element.propertiesKeepOriginalName
				return ba_66 || bb_67
			}()
			var elemPrototypes []*Color = element.prototypes
			_ = elemPrototypes
			if int32(len(elemPrototypes)) > int32(0) {
				if prototypesBuilder != nil {
					for _, el_68 := range j2gort.CheckNotNullSlice(elemPrototypes) {
						c_4 := el_68
						_ = c_4
						prototypesBuilder.Add_1(j2gort.PtrToAny(c_4))
					}
				} else if int32(len(singlePrototypes)) == int32(0) {
					singlePrototypes = elemPrototypes
				} else if !(Color_SetEqualsColors(singlePrototypes, elemPrototypes)) {
					prototypesBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singlePrototypes)) + int32(len(elemPrototypes))))
					for _, el_69 := range j2gort.CheckNotNullSlice(singlePrototypes) {
						c_5 := el_69
						_ = c_5
						prototypesBuilder.Add_1(j2gort.PtrToAny(c_5))
					}
					for _, el_70 := range j2gort.CheckNotNullSlice(elemPrototypes) {
						c_6 := el_70
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
	var builtUnionElements []*Color = j2gort.ColToPtrArray[Color](newElements, StaticGet_Color_EMPTY_COLORS())
	_ = builtUnionElements
	var result *Color = NewColor(ColorId_UnionFromColors(builtUnionElements), func() []*Color {
		if prototypesBuilder != nil {
			return j2gort.ColToPtrArray[Color](prototypesBuilder, StaticGet_Color_EMPTY_COLORS())
		}
		return singlePrototypes
	}(), func() []*Color {
		if instanceColorsBuilder != nil {
			return j2gort.ColToPtrArray[Color](instanceColorsBuilder, StaticGet_Color_EMPTY_COLORS())
		}
		return singleInstanceColors
	}(), isInvalidating, propertiesKeepOriginalName, isConstructor, func() []string {
		if ownPropertiesBuilder != nil {
			return j2gort.ColToStringArray(ownPropertiesBuilder, StaticGet_Color_EMPTY_STRINGS())
		}
		return singleOwnProperties
	}(), (ColorId{}), isClosureAssert, builtUnionElements)
	_ = result
	if !(int32(len(result.unionElements)) > int32(1)) {
		j2gort.CheckStateMsg(false, result)
	}
	if !((result.boxId).IsJavaNil_()) {
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
				for _, el_71 := range j2gort.CheckNotNullSlice(element.unionElements) {
					nestedElement := el_71
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
					for _, el_72 := range j2gort.CheckNotNullSlice(elemInstanceColors) {
						c := el_72
						_ = c
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c))
					}
				} else if int32(len(singleInstanceColors)) == int32(0) {
					singleInstanceColors = elemInstanceColors
				} else if !(Color_SetEqualsColors(singleInstanceColors, elemInstanceColors)) {
					instanceColorsBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singleInstanceColors)) + int32(len(elemInstanceColors))))
					for _, el_73 := range j2gort.CheckNotNullSlice(singleInstanceColors) {
						c_2 := el_73
						_ = c_2
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c_2))
					}
					for _, el_74 := range j2gort.CheckNotNullSlice(elemInstanceColors) {
						c_3 := el_74
						_ = c_3
						instanceColorsBuilder.Add_1(j2gort.PtrToAny(c_3))
					}
				}
			}
			{
				ba_75 := isClosureAssert
				bb_76 := element.closureAssert
				isClosureAssert = ba_75 && bb_76
			}
			{
				ba_77 := isConstructor
				bb_78 := element.constructor
				isConstructor = ba_77 && bb_78
			}
			{
				ba_79 := isInvalidating
				bb_80 := element.invalidating
				isInvalidating = ba_79 || bb_80
			}
			var elemOwnProps []string = element.ownProperties
			_ = elemOwnProps
			if int32(len(elemOwnProps)) > int32(0) {
				if ownPropertiesBuilder != nil {
					for _, el_81 := range j2gort.CheckNotNullSlice(elemOwnProps) {
						p := el_81
						_ = p
						j2gort.SetAddStr(ownPropertiesBuilder, p)
					}
				} else if int32(len(singleOwnProperties)) == int32(0) {
					singleOwnProperties = elemOwnProps
				} else if !(Color_SetEqualsStrings(singleOwnProperties, elemOwnProps)) {
					ownPropertiesBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singleOwnProperties)) + int32(len(elemOwnProps))))
					for _, el_82 := range j2gort.CheckNotNullSlice(singleOwnProperties) {
						p_2 := el_82
						_ = p_2
						j2gort.SetAddStr(ownPropertiesBuilder, p_2)
					}
					for _, el_83 := range j2gort.CheckNotNullSlice(elemOwnProps) {
						p_3 := el_83
						_ = p_3
						j2gort.SetAddStr(ownPropertiesBuilder, p_3)
					}
				}
			}
			propertiesKeepOriginalName = func() bool {
				ba_84 := propertiesKeepOriginalName
				bb_85 := element.propertiesKeepOriginalName
				return ba_84 || bb_85
			}()
			var elemPrototypes []*Color = element.prototypes
			_ = elemPrototypes
			if int32(len(elemPrototypes)) > int32(0) {
				if prototypesBuilder != nil {
					for _, el_86 := range j2gort.CheckNotNullSlice(elemPrototypes) {
						c_4 := el_86
						_ = c_4
						prototypesBuilder.Add_1(j2gort.PtrToAny(c_4))
					}
				} else if int32(len(singlePrototypes)) == int32(0) {
					singlePrototypes = elemPrototypes
				} else if !(Color_SetEqualsColors(singlePrototypes, elemPrototypes)) {
					prototypesBuilder = j2gort.NewExt_java_util_LinkedHashSet_1_int32((int32(len(singlePrototypes)) + int32(len(elemPrototypes))))
					for _, el_87 := range j2gort.CheckNotNullSlice(singlePrototypes) {
						c_5 := el_87
						_ = c_5
						prototypesBuilder.Add_1(j2gort.PtrToAny(c_5))
					}
					for _, el_88 := range j2gort.CheckNotNullSlice(elemPrototypes) {
						c_6 := el_88
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
	var builtUnionElements []*Color = j2gort.ColToPtrArray[Color](newElements, StaticGet_Color_EMPTY_COLORS())
	_ = builtUnionElements
	var result *Color = NewColor(ColorId_UnionFromColors(builtUnionElements), func() []*Color {
		if prototypesBuilder != nil {
			return j2gort.ColToPtrArray[Color](prototypesBuilder, StaticGet_Color_EMPTY_COLORS())
		}
		return singlePrototypes
	}(), func() []*Color {
		if instanceColorsBuilder != nil {
			return j2gort.ColToPtrArray[Color](instanceColorsBuilder, StaticGet_Color_EMPTY_COLORS())
		}
		return singleInstanceColors
	}(), isInvalidating, propertiesKeepOriginalName, isConstructor, func() []string {
		if ownPropertiesBuilder != nil {
			return j2gort.ColToStringArray(ownPropertiesBuilder, StaticGet_Color_EMPTY_STRINGS())
		}
		return singleOwnProperties
	}(), (ColorId{}), isClosureAssert, builtUnionElements)
	_ = result
	if !(int32(len(result.unionElements)) > int32(1)) {
		j2gort.CheckStateMsg(false, result)
	}
	if !((result.boxId).IsJavaNil_()) {
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
		return this_.priv_Color_computeSubtractNullOrVoid_5979()
	}
	return cached
}

// Java method: Color.computeSubtractNullOrVoid()
func (this_ *Color) priv_Color_computeSubtractNullOrVoid_5979() *Color {
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
						acSrc_89 := elems
						acSP_90 := int32(0)
						acDst_91 := remaining
						acDP_92 := int32(0)
						acLen_93 := nullOrVoidIndex
						j2gort.CheckNotNullSlice(acSrc_89)
						j2gort.CheckNotNullSlice(acDst_91)
						if acSP_90 < 0 || acDP_92 < 0 || acLen_93 < 0 || int64(acSP_90)+int64(acLen_93) > int64(len(acSrc_89)) || int64(acDP_92)+int64(acLen_93) > int64(len(acDst_91)) {
							panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
						}
						copy(acDst_91[int(acDP_92):int(acDP_92)+int(acLen_93)], acSrc_89[int(acSP_90):int(acSP_90)+int(acLen_93)])
					}()
				}
				if nullOrVoidIndex < (len_ - int32(1)) {
					func() {
						acSrc_94 := elems
						acSP_95 := (nullOrVoidIndex + int32(1))
						acDst_96 := remaining
						acDP_97 := nullOrVoidIndex
						acLen_98 := ((len_ - nullOrVoidIndex) - int32(1))
						j2gort.CheckNotNullSlice(acSrc_94)
						j2gort.CheckNotNullSlice(acDst_96)
						if acSP_95 < 0 || acDP_97 < 0 || acLen_98 < 0 || int64(acSP_95)+int64(acLen_98) > int64(len(acSrc_94)) || int64(acDP_97)+int64(acLen_98) > int64(len(acDst_96)) {
							panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
						}
						copy(acDst_96[int(acDP_97):int(acDP_97)+int(acLen_98)], acSrc_94[int(acSP_95):int(acSP_95)+int(acLen_98)])
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
		{

			sr_99 := a[int32(0)]
			se_100 := b[int32(0)]
			return j2gort.CheckNotNullString(sr_99) == se_100
		}

	}
	if len_ == int32(2) {
		return ((func() bool {
			sr_101 := a[int32(0)]
			se_102 := b[int32(0)]
			return j2gort.CheckNotNullString(sr_101) == se_102
		}() && func() bool {
			sr_103 := a[int32(1)]
			se_104 := b[int32(1)]
			return j2gort.CheckNotNullString(sr_103) == se_104
		}()) || (func() bool {
			sr_105 := a[int32(0)]
			se_106 := b[int32(1)]
			return j2gort.CheckNotNullString(sr_105) == se_106
		}() && func() bool {
			sr_107 := a[int32(1)]
			se_108 := b[int32(0)]
			return j2gort.CheckNotNullString(sr_107) == se_108
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

					sr_109 := a[i]
					se_110 := b[i]
					iifeCond_1 = j2gort.CheckNotNullString(sr_109) == se_110
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
	var bSet j2gort.PlainSet[string] = j2gort.NewPlainSetCap[string]((len_ * int32(2)))
	_ = bSet
	{
		var i_3 int32 = int32(0)
		_ = i_3
		for ; i_3 < len_; i_3 = i_3 + 1 {
			bSet.AddVoid(b[i_3])
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
func (this_ *Color) Equals(arg_111 any) bool {
	o := arg_111
	_ = o
	if j2gort.RefEquals(o, this_) {
		return true
	}
	if j2gort.InstanceOf[*Color](o) {
		var that *Color = j2gort.CastPtr[Color](o)
		_ = that
		return ((((((((((func() bool {
			rv_112 := this_.id
			av_113 := that.id
			_ = j2gort.CheckNotNullVal(rv_112)
			return rv_112 == av_113
		}() && (this_.invalidating == that.invalidating)) && (this_.propertiesKeepOriginalName == that.propertiesKeepOriginalName)) && (this_.constructor == that.constructor)) && (this_.closureAssert == that.closureAssert)) && (func() bool {
			if (this_.boxId).IsJavaNil_() {
				return (that.boxId).IsJavaNil_()
			}
			{

				rv_114 := this_.boxId
				av_115 := that.boxId
				_ = j2gort.CheckNotNullVal(rv_114)
				return rv_114 == av_115
			}

		}())) && (this_.HashCode() == that.HashCode())) && Color_SetEqualsColors(this_.prototypes, that.prototypes)) && Color_SetEqualsColors(this_.instanceColors, that.instanceColors)) && Color_SetEqualsStrings(this_.ownProperties, that.ownProperties)) && Color_SetEqualsColors(this_.unionElements, that.unionElements))
	}
	return false
}

// Java method: Color.hashCode()
func (this_ *Color) HashCode() int32 {
	var h int32 = this_.hashCode
	_ = h
	if h == int32(0) {
		h = this_.priv_Color_computeHashCode_5983()
		this_.hashCode = h
	}
	return h
}

// Java method: Color.computeHashCode()
func (this_ *Color) priv_Color_computeHashCode_5983() int32 {
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
		if (this_.boxId).IsJavaNil_() {
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
	if (id).IsJavaNil_() {
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
	if (((((int32(this_.set_d0) != int32(15)) || (this_.id).IsJavaNil_()) || (this_.prototypes == nil)) || (this_.instanceColors == nil)) || (this_.ownProperties == nil)) || (this_.unionElements == nil) {
		var missing *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = missing
		if (this_.id).IsJavaNil_() {
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
	if !((result.boxId).IsJavaNil_()) {
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
