/*
 * Copyright 2021 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/serialization/ColorSerializer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

var lambdaSingleton_color_serializer_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_color_serializer_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_12 any) any {
		key := j2gort.CastString(lp_12)
		_ = key
		return j2gort.NewExt_java_util_ArrayList()
	})
}

// Java class: com.google.javascript.jscomp.serialization.ColorSerializer
type ColorSerializer struct {
	serializationMode       *SerializationOptions                                 `gson:"serializationMode"`
	getStringPoolIndexFn    j2gort.Ext_java_util_function_Function                `gson:"getStringPoolIndexFn"`
	propertyFilter          j2gort.Ext_java_util_function_Predicate               `gson:"propertyFilter"`
	colorIdToTypePointer    j2gort.OrderedMap[pkg_colors.ColorId, *j2gort.BoxInt] `gson:"colorIdToTypePointer"`
	colorsInSerializedOrder []*pkg_colors.Color                                   `gson:"colorsInSerializedOrder"`
	idHash_                 atomic.Int32
}

func (this_ *ColorSerializer) BindSelf_(s any) {
	_ = s
}
func (this_ *ColorSerializer) IsJavaNil_() bool { return this_ == nil }
func (*ColorSerializer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.ColorSerializer", "ColorSerializer"
}
func (this_ *ColorSerializer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ColorSerializer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ColorSerializer", "com.google.javascript.jscomp.serialization.ColorSerializer", func(v_ any) bool {
		return j2gort.InstanceOf[*ColorSerializer](v_)
	})
}

func (this_ *ColorSerializer) initFields_ColorSerializer() {
	_ = this_
	this_.colorIdToTypePointer = j2gort.NewOrderedMapCap[pkg_colors.ColorId, *j2gort.BoxInt](int32(4096))
	this_.colorsInSerializedOrder = make([]*pkg_colors.Color, 0, int32(4096))
}

// Java constructor: ColorSerializer(com.google.javascript.jscomp.serialization.SerializationOptions serializationMode, java.util.function.Function<java.lang.String,java.lang.Integer> stringPoolIndexFn, java.util.function.Predicate<java.lang.String> propertyFilter)
func NewColorSerializer(serializationMode *SerializationOptions, stringPoolIndexFn j2gort.Ext_java_util_function_Function, propertyFilter j2gort.Ext_java_util_function_Predicate) *ColorSerializer {
	this_ := &ColorSerializer{}
	InitColorSerializer(this_, serializationMode, stringPoolIndexFn, propertyFilter)
	return this_
}

func InitColorSerializer(this_ *ColorSerializer, serializationMode *SerializationOptions, stringPoolIndexFn j2gort.Ext_java_util_function_Function, propertyFilter j2gort.Ext_java_util_function_Predicate) {
	_ = serializationMode
	_ = stringPoolIndexFn
	_ = propertyFilter
	this_.initFields_ColorSerializer()
	this_.serializationMode = serializationMode
	this_.getStringPoolIndexFn = stringPoolIndexFn
	this_.propertyFilter = propertyFilter
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_TypePointers_OFFSET_TO_AXIOMATIC_COLOR()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			color := j2gort.CastPtr[pkg_colors.Color](rawElem_2)
			_ = color
			this_.AddColor(color)
		}
	}
}

// Java method: ColorSerializer.addColors(java.util.Collection<com.google.javascript.jscomp.colors.Color> colors)
func (this_ *ColorSerializer) AddColors_1_Ext_java_util_Collection(colors j2gort.Ext_java_util_Collection) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = colors
	var builder []int32 = []int32{}
	_ = builder
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(colors))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			color := j2gort.CastPtr[pkg_colors.Color](rawElem_4)
			_ = color
			builder = append(builder, this_.AddColor(color))
		}
	}
	return j2gort.ImmListFromTypedSlice(builder)
}

// Java method: ColorSerializer.addColors(com.google.javascript.jscomp.colors.Color[] colors)
func (this_ *ColorSerializer) AddColors_1_Color(colors []*pkg_colors.Color) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = colors
	var len_ int32 = int32(len(colors))
	_ = len_
	if len_ == int32(0) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	if len_ == int32(1) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.ValueOfInt(this_.AddColor(colors[int32(0)])))
	}
	var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, len_)
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			builder.AddVoid_(j2gort.ValueOfInt(this_.AddColor(colors[i])))
		}
	}
	return builder.BuildImmList_()
}

// Java method: ColorSerializer.addColor(com.google.javascript.jscomp.colors.Color color)
func (this_ *ColorSerializer) AddColor(color *pkg_colors.Color) int32 {
	_ = color
	var id pkg_colors.ColorId = color.GetId()
	_ = id
	var existing *j2gort.BoxInt = this_.colorIdToTypePointer.Get(id)
	_ = existing
	if existing != nil {
		return j2gort.UnboxInt(existing)
	}
	var index int32 = int32(len(this_.colorsInSerializedOrder))
	_ = index
	this_.colorsInSerializedOrder = append(this_.colorsInSerializedOrder, color)
	this_.colorIdToTypePointer.Put(id, j2gort.ValueOfInt(index))
	return index
}

// Java method: ColorSerializer.generateTypePool(java.util.function.Function<com.google.javascript.jscomp.colors.Color,com.google.javascript.jscomp.colors.Color[]> getDisambiguationSupertypesFn, java.util.function.@org.jspecify.annotations.Nullable Function<com.google.javascript.jscomp.colors.Color,com.google.common.collect.ImmutableSet<java.lang.String>> getMismatchSourceRefsFn)
func (this_ *ColorSerializer) GenerateTypePool(getDisambiguationSupertypesFn j2gort.Ext_java_util_function_Function, getMismatchSourceRefsFn j2gort.Ext_java_util_function_Function) *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool {
	_ = getDisambiguationSupertypesFn
	_ = getMismatchSourceRefsFn
	var typePoolBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool_NewBuilder_0()
	_ = typePoolBuilder
	{
		var i int32 = TypePointers_UntrimOffset(int32(0))
		_ = i
		for ; i < int32(len(this_.colorsInSerializedOrder)); i = i + 1 {
			var color *pkg_colors.Color = this_.colorsInSerializedOrder[i]
			_ = color
			var typePointer *j2gort.BoxInt = this_.colorIdToTypePointer.Get(color.GetId())
			_ = typePointer
			typePoolBuilder.AddType_1_Ext_com_google_javascript_jscomp_serialization_TypeProto(this_.priv_ColorSerializer_generateTypeProto_7689(color))
			var supertypes []*pkg_colors.Color = j2gort.Cast[[]*pkg_colors.Color](getDisambiguationSupertypesFn.Apply(j2gort.PtrToAny(color)))
			_ = supertypes
			{
				var s int32 = int32(0)
				_ = s
				var sLen int32 = int32(len(supertypes))
				_ = sLen
				for ; s < sLen; s = s + 1 {
					typePoolBuilder.AddDisambiguationEdgesBuilder_0().SetSubtype(j2gort.UnboxInt(typePointer)).SetSupertype(this_.AddColor(supertypes[s]))
				}
			}
		}
	}
	if this_.serializationMode.IncludeDebugInfo() {
		j2gort.CheckNotNull(getMismatchSourceRefsFn)
		var debugInfoBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Builder = typePoolBuilder.GetDebugInfoBuilder()
		_ = debugInfoBuilder
		var srcRefToTypePointerList j2gort.OrderedMap[string, j2gort.Ext_java_util_ArrayList] = j2gort.NewOrderedMap[string, j2gort.Ext_java_util_ArrayList]()
		_ = srcRefToTypePointerList
		for _, color_2 := range this_.colorsInSerializedOrder {
			_ = color_2
			var typePointer_2 *j2gort.BoxInt = this_.colorIdToTypePointer.Get(color_2.GetId())
			_ = typePointer_2
			{
				cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](getMismatchSourceRefsFn.Apply(j2gort.PtrToAny(color_2)))))
				for cur_7.HasNext() {
					srcRef := cur_7.NextStr()
					_ = srcRef
					var typePointerList j2gort.Ext_java_util_ArrayList = j2gort.OrderedMapComputeIfAbsentFn(&srcRefToTypePointerList, srcRef, nil, lambdaSingleton_color_serializer_1)

					_ = typePointerList
					typePointerList.Add_1(j2gort.PtrToAny(typePointer_2))
				}
			}
		}

		{
			cur_13 := srcRefToTypePointerList.IterEntries()
			for cur_13.NextEntry() {
				debugInfoBuilder.AddMismatchBuilder_0().SetSourceRef(cur_13.Key()).AddAllInvolvedColor(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](cur_13.Value()))
			}
		}
	}
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool](typePoolBuilder.Build())
}

// Java method: ColorSerializer.generateTypeProto(com.google.javascript.jscomp.colors.Color color)
func (this_ *ColorSerializer) priv_ColorSerializer_generateTypeProto_7689(color *pkg_colors.Color) *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto {
	_ = color
	var typeProtoBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto_NewBuilder_0()
	_ = typeProtoBuilder
	if color.IsUnion() {
		typeProtoBuilder.GetUnionBuilder().AddAllUnionMember(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.AddColors_1_Color(color.GetUnionElementsArray())))
	} else {
		var objectTypeProtoBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto_Builder = typeProtoBuilder.GetObjectBuilder()
		_ = objectTypeProtoBuilder
		objectTypeProtoBuilder.SetIsInvalidating(color.IsInvalidating()).SetUuid(color.GetId().AsByteString()).SetPropertiesKeepOriginalName(color.GetPropertiesKeepOriginalName()).AddAllInstanceType(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.AddColors_1_Color(color.GetInstanceColorsArray()))).AddAllPrototype(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.AddColors_1_Color(color.GetPrototypesArray()))).SetMarkedConstructor(color.IsConstructor()).AddAllOwnProperty(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.priv_ColorSerializer_getOwnPropertyStringPoolOffsets_7690(color))).SetClosureAssert(color.IsClosureAssert())
	}
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](typeProtoBuilder.Build())
}

// Java method: ColorSerializer.getOwnPropertyStringPoolOffsets(com.google.javascript.jscomp.colors.Color color)
func (this_ *ColorSerializer) priv_ColorSerializer_getOwnPropertyStringPoolOffsets_7690(color *pkg_colors.Color) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = color
	var ownProperties []string = color.GetOwnPropertiesArray()
	_ = ownProperties
	var len_ int32 = int32(len(ownProperties))
	_ = len_
	if len_ == int32(0) {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = nil
	_ = builder
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			var ownProperty string = ownProperties[i]
			_ = ownProperty
			if this_.propertyFilter.Test(j2gort.BoxString(ownProperty)) {
				if builder == nil {
					builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
				}
				builder.AddVoid_(j2gort.PtrToAny(j2gort.CastPtr[j2gort.BoxInt](this_.getStringPoolIndexFn.Apply(j2gort.BoxString(ownProperty)))))
			}
		}
	}
	if builder != nil {
		return builder.BuildImmList_()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
}
