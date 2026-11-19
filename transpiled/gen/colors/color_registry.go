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
// Port of com/google/javascript/jscomp/colors/ColorRegistry.java (generated mechanically by j2go).

package colors

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.colors.ColorRegistry
type ColorRegistry struct {
	nativeColors                        j2gort.Ext_com_google_common_collect_ImmutableMap          `gson:"nativeColors"`
	colorToDisambiguationSupertypeGraph *j2gort.Ext_java_util_LinkedHashMap                        `gson:"colorToDisambiguationSupertypeGraph"`
	mismatchLocations                   *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap `gson:"mismatchLocations"`
	idHash_                             atomic.Int32
}

func (this_ *ColorRegistry) BindSelf_(s any) {
	_ = s
}
func (this_ *ColorRegistry) IsJavaNil_() bool { return this_ == nil }
func (*ColorRegistry) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.colors.ColorRegistry", "ColorRegistry"
}
func (this_ *ColorRegistry) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ColorRegistry) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.colors.ColorRegistry", "com.google.javascript.jscomp.colors.ColorRegistry", func(v_ any) bool {
		return j2gort.InstanceOf[*ColorRegistry](v_)
	})
}

var static_ColorRegistry_REQUIRED_IDS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ColorRegistry_REQUIRED_IDS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ColorRegistry()
	return static_ColorRegistry_REQUIRED_IDS
}
func StaticSet_ColorRegistry_REQUIRED_IDS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ColorRegistry()
	static_ColorRegistry_REQUIRED_IDS = v
}

var clinitDone_ColorRegistry uint32
var clinitState_ColorRegistry j2gort.ClinitState

func clinit_ColorRegistry() {
	if atomic.LoadUint32(&clinitDone_ColorRegistry) != 0 {
		return
	}
	clinitSlow_ColorRegistry()
}

func clinitSlow_ColorRegistry() {
	if !j2gort.EnterClinit(&clinitState_ColorRegistry, &clinitDone_ColorRegistry) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ColorRegistry, &clinitDone_ColorRegistry)
	static_ColorRegistry_REQUIRED_IDS = StaticGet_StandardColors_STANDARD_OBJECT_IDS()
}

func (this_ *ColorRegistry) initFields_ColorRegistry() {
	_ = this_
}

// Java constructor: ColorRegistry(com.google.javascript.jscomp.colors.ColorRegistry.Builder builder)
func NewColorRegistry(builder *ColorRegistry_Builder) *ColorRegistry {
	clinit_ColorRegistry()
	this_ := &ColorRegistry{}
	InitColorRegistry(this_, builder)
	return this_
}

func InitColorRegistry(this_ *ColorRegistry, builder *ColorRegistry_Builder) {
	_ = builder
	this_.initFields_ColorRegistry()
	this_.nativeColors = builder.nativeColors.ToImmutableMap()
	var rawGraph *j2gort.Ext_java_util_LinkedHashMap = builder.colorToDisambiguationSupertypeGraph
	_ = rawGraph
	var builtGraph *j2gort.Ext_java_util_LinkedHashMap = j2gort.NewExt_java_util_LinkedHashMap_1_int32(rawGraph.Size())
	_ = builtGraph
	{
		cur_1 := j2gort.NewMapEntryCursor(rawGraph)
		for cur_1.NextEntry() {
			var val any = cur_1.GetValue()
			_ = val
			var single *Color
			_ = single
			var ok_2 bool
			single, ok_2 = j2gort.TryCast[*Color](val)
			if ok_2 {
				builtGraph.Put_2(j2gort.PtrToAny(j2gort.CastPtr[Color](cur_1.GetKey())), []*Color{single})
			} else {
				var set *j2gort.Ext_java_util_LinkedHashSet = j2gort.CastPtr[j2gort.Ext_java_util_LinkedHashSet](val)
				_ = set
				builtGraph.Put_2(j2gort.PtrToAny(j2gort.CastPtr[Color](cur_1.GetKey())), j2gort.SliceToAny(j2gort.ColToPtrArray[Color](set, StaticGet_Color_EMPTY_COLORS())))
			}
		}
	}
	this_.colorToDisambiguationSupertypeGraph = builtGraph
	if builder.mismatchLocations == nil {
		this_.mismatchLocations = j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_Of_0()
	} else {
		this_.mismatchLocations = j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_CopyOf_1_Ext_com_google_common_collect_Multimap(j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](builder.mismatchLocations))
	}
	if !(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.nativeColors.KeySet_0()).Equals(StaticGet_ColorRegistry_REQUIRED_IDS())) {
		j2gort.CheckState(false)
	}
}

// Java method: ColorRegistry.get(com.google.javascript.jscomp.colors.ColorId id)
func (this_ *ColorRegistry) Get(id ColorId) *Color {
	_ = id
	j2gort.CheckNotNullVal(id)
	return j2gort.CheckNotNullPtr(j2gort.CastPtr[Color](j2gort.MapGetEqVal(this_.nativeColors, id)))
}

// Java method: ColorRegistry.getDisambiguationSupertypes(com.google.javascript.jscomp.colors.Color x)
func (this_ *ColorRegistry) GetDisambiguationSupertypes(x *Color) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = x
	var supertypes []*Color = this_.GetDisambiguationSupertypesArray(x)
	_ = supertypes
	var len_ int32 = int32(len(supertypes))
	_ = len_
	if len_ == int32(0) {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	}
	if len_ == int32(1) {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_1(j2gort.PtrToAny(supertypes[int32(0)]))
	}
	return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_slice_any(func() []any {
		arr_3 := supertypes
		if arr_3 == nil {
			return nil
		}
		out_4 := make([]any, len(arr_3))
		for i_, el_5 := range arr_3 {
			out_4[i_] = j2gort.PtrToAny(el_5)
		}
		return out_4
	}())
}

// Java method: ColorRegistry.getDisambiguationSupertypesArray(com.google.javascript.jscomp.colors.Color x)
func (this_ *ColorRegistry) GetDisambiguationSupertypesArray(x *Color) []*Color {
	_ = x
	j2gort.CheckNotNullPtr(x)
	var supertypes []*Color = j2gort.Cast[[]*Color](this_.colorToDisambiguationSupertypeGraph.Get(j2gort.PtrToAny(x)))
	_ = supertypes
	if supertypes != nil {
		return supertypes
	} else {
		return StaticGet_Color_EMPTY_COLORS()
	}
}

// Java method: ColorRegistry.getMismatchLocationsForDebugging()
func (this_ *ColorRegistry) GetMismatchLocationsForDebugging() *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap {
	return this_.mismatchLocations
}

// Java static method: ColorRegistry.builder()
func ColorRegistry_S_Builder() *ColorRegistry_Builder {
	clinit_ColorRegistry()
	return NewColorRegistry_Builder()
}

// Java class: com.google.javascript.jscomp.colors.ColorRegistry.Builder
type ColorRegistry_Builder struct {
	nativeColors                        j2gort.OrderedMap[ColorId, *Color]               `gson:"nativeColors"`
	colorToDisambiguationSupertypeGraph *j2gort.Ext_java_util_LinkedHashMap              `gson:"colorToDisambiguationSupertypeGraph"`
	mismatchLocations                   j2gort.Ext_com_google_common_collect_SetMultimap `gson:"mismatchLocations"`
	idHash_                             atomic.Int32
}

func (this_ *ColorRegistry_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *ColorRegistry_Builder) IsJavaNil_() bool { return this_ == nil }
func (*ColorRegistry_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.colors.ColorRegistry$Builder", "Builder"
}
func (this_ *ColorRegistry_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ColorRegistry_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.colors.ColorRegistry.Builder", "com.google.javascript.jscomp.colors.ColorRegistry$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*ColorRegistry_Builder](v_)
	})
}

func (this_ *ColorRegistry_Builder) initFields_ColorRegistry_Builder() {
	_ = this_
	this_.nativeColors = j2gort.NewOrderedMap[ColorId, *Color]()
	this_.colorToDisambiguationSupertypeGraph = j2gort.NewExt_java_util_LinkedHashMap()

}

// Java constructor: Builder()
func NewColorRegistry_Builder() *ColorRegistry_Builder {
	this_ := &ColorRegistry_Builder{}
	InitColorRegistry_Builder(this_)
	return this_
}

func InitColorRegistry_Builder(this_ *ColorRegistry_Builder) {
	this_.initFields_ColorRegistry_Builder()
}

// Java method: Builder.setNativeColor(com.google.javascript.jscomp.colors.Color x)
func (this_ *ColorRegistry_Builder) SetNativeColor(x *Color) *ColorRegistry_Builder {
	_ = x
	if !(j2gort.SetContainsEqVal(StaticGet_ColorRegistry_REQUIRED_IDS(), x.GetId())) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(x))
	}
	this_.nativeColors.Put(x.GetId(), x)
	return this_
}

// Java method: Builder.addDisambiguationEdge(com.google.javascript.jscomp.colors.Color subtype, com.google.javascript.jscomp.colors.Color supertype)
func (this_ *ColorRegistry_Builder) AddDisambiguationEdge(subtype *Color, supertype *Color) *ColorRegistry_Builder {
	_ = subtype
	_ = supertype
	j2gort.CheckNotNullPtr(subtype)
	j2gort.CheckNotNullPtr(supertype)
	var existing any = this_.colorToDisambiguationSupertypeGraph.Get(j2gort.PtrToAny(subtype))
	_ = existing
	var first *Color
	_ = first
	if existing == nil {
		this_.colorToDisambiguationSupertypeGraph.Put_2(j2gort.PtrToAny(subtype), j2gort.PtrToAny(supertype))
	} else {
		var ok_6 bool
		first, ok_6 = j2gort.TryCast[*Color](existing)
		if ok_6 {
			if !(first.Equals(j2gort.PtrToAny(supertype))) {
				var set *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(int32(4))
				_ = set
				set.Add_1(j2gort.PtrToAny(first))
				set.Add_1(j2gort.PtrToAny(supertype))
				this_.colorToDisambiguationSupertypeGraph.Put_2(j2gort.PtrToAny(subtype), set)
			}
		} else {
			_ = (j2gort.CastPtr[j2gort.Ext_java_util_LinkedHashSet](existing)).Add_1(j2gort.PtrToAny(supertype))
		}
	}
	return this_
}

// Java method: Builder.addMismatchLocation(com.google.javascript.jscomp.colors.ColorId id, java.lang.String location)
func (this_ *ColorRegistry_Builder) AddMismatchLocation(id ColorId, location string) *ColorRegistry_Builder {
	_ = id
	_ = location
	j2gort.CheckNotNullVal(id)
	j2gort.CheckNotNullString(location)
	if this_.mismatchLocations == nil {
		this_.mismatchLocations = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_MultimapBuilder_LinkedHashKeys_0().LinkedHashSetValues_0().Build_0())
	}
	this_.mismatchLocations.Put(j2gort.BoxNullableVal(id), j2gort.BoxString(location))
	return this_
}

// Java method: Builder.setDefaultNativeColorsForTesting()
func (this_ *ColorRegistry_Builder) SetDefaultNativeColorsForTesting() *ColorRegistry_Builder {
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_ColorRegistry_REQUIRED_IDS()))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			id := j2gort.Cast[ColorId](rawElem_8)
			_ = id
			this_.SetNativeColor(Color_SingleBuilder().SetId(id).Build())
		}
	}
	return this_
}

// Java method: Builder.build()
func (this_ *ColorRegistry_Builder) Build() *ColorRegistry {
	return NewColorRegistry(this_)
}
