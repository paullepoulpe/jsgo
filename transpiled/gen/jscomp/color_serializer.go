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
	serializationMode       *SerializationOptions                        `gson:"serializationMode"`
	getStringPoolIndexFn    ColorSerializer_StringPoolIndexFn            `gson:"getStringPoolIndexFn"`
	propertyFilter          ColorSerializer_PropertyFilter               `gson:"propertyFilter"`
	colorIdToTypePointer    j2gort.OrderedMap[pkg_colors.ColorId, int32] `gson:"colorIdToTypePointer"`
	colorsInSerializedOrder []*pkg_colors.Color                          `gson:"colorsInSerializedOrder"`
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
	this_.colorIdToTypePointer = j2gort.NewOrderedMapCap[pkg_colors.ColorId, int32](int32(4096))
	this_.colorsInSerializedOrder = make([]*pkg_colors.Color, 0, int32(4096))
}

// Java constructor: ColorSerializer(com.google.javascript.jscomp.serialization.SerializationOptions serializationMode, com.google.javascript.jscomp.serialization.ColorSerializer.StringPoolIndexFn stringPoolIndexFn, com.google.javascript.jscomp.serialization.ColorSerializer.PropertyFilter propertyFilter)
func NewColorSerializer(serializationMode *SerializationOptions, stringPoolIndexFn ColorSerializer_StringPoolIndexFn, propertyFilter ColorSerializer_PropertyFilter) *ColorSerializer {
	this_ := &ColorSerializer{}
	InitColorSerializer(this_, serializationMode, stringPoolIndexFn, propertyFilter)
	return this_
}

func InitColorSerializer(this_ *ColorSerializer, serializationMode *SerializationOptions, stringPoolIndexFn ColorSerializer_StringPoolIndexFn, propertyFilter ColorSerializer_PropertyFilter) {
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
	existing, existingOk_ := this_.colorIdToTypePointer.GetOk(id)
	_ = existingOk_
	_ = existing
	if existingOk_ {
		return func() int32 {
			if !existingOk_ {
				j2gort.ThrowNPE()
			}
			return existing
		}()
	}
	var index int32 = int32(len(this_.colorsInSerializedOrder))
	_ = index
	this_.colorsInSerializedOrder = append(this_.colorsInSerializedOrder, color)
	this_.colorIdToTypePointer.Put(id, index)
	return index
}

// Java method: ColorSerializer.generateTypePool(com.google.javascript.jscomp.serialization.ColorSerializer.DisambiguationSupertypesFn getDisambiguationSupertypesFn, java.util.function.@org.jspecify.annotations.Nullable Function<com.google.javascript.jscomp.colors.Color,com.google.common.collect.ImmutableSet<java.lang.String>> getMismatchSourceRefsFn)
func (this_ *ColorSerializer) GenerateTypePool(getDisambiguationSupertypesFn ColorSerializer_DisambiguationSupertypesFn, getMismatchSourceRefsFn j2gort.Ext_java_util_function_Function) *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool {
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
			typePointer, typePointerOk_ := this_.colorIdToTypePointer.GetOk(color.GetId())
			_ = typePointerOk_
			_ = typePointer
			typePoolBuilder.AddType_1_Ext_com_google_javascript_jscomp_serialization_TypeProto(this_.priv_ColorSerializer_generateTypeProto_7706(color))
			var supertypes []*pkg_colors.Color = getDisambiguationSupertypesFn.Apply(color)
			_ = supertypes
			{
				var s int32 = int32(0)
				_ = s
				var sLen int32 = int32(len(supertypes))
				_ = sLen
				for ; s < sLen; s = s + 1 {
					typePoolBuilder.AddDisambiguationEdgesBuilder_0().SetSubtype(func() int32 {
						if !typePointerOk_ {
							j2gort.ThrowNPE()
						}
						return typePointer
					}()).SetSupertype(this_.AddColor(supertypes[s]))
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
			typePointer_2, typePointer_2Ok_ := this_.colorIdToTypePointer.GetOk(color_2.GetId())
			_ = typePointer_2Ok_
			_ = typePointer_2
			{
				cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](getMismatchSourceRefsFn.Apply(j2gort.PtrToAny(color_2)))))
				for cur_7.HasNext() {
					srcRef := cur_7.NextStr()
					_ = srcRef
					var typePointerList j2gort.Ext_java_util_ArrayList = j2gort.OrderedMapComputeIfAbsentFn(&srcRefToTypePointerList, srcRef, nil, lambdaSingleton_color_serializer_1)

					_ = typePointerList
					typePointerList.Add_1(func() any {
						if !typePointer_2Ok_ {
							return nil
						}
						return j2gort.ValueOfInt(typePointer_2)
					}())
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
func (this_ *ColorSerializer) priv_ColorSerializer_generateTypeProto_7706(color *pkg_colors.Color) *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto {
	_ = color
	var typeProtoBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto_Builder = j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto_NewBuilder_0()
	_ = typeProtoBuilder
	if color.IsUnion() {
		typeProtoBuilder.GetUnionBuilder().AddAllUnionMember(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.AddColors_1_Color(color.GetUnionElementsArray())))
	} else {
		var objectTypeProtoBuilder *j2gort.Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto_Builder = typeProtoBuilder.GetObjectBuilder()
		_ = objectTypeProtoBuilder
		objectTypeProtoBuilder.SetIsInvalidating(color.IsInvalidating()).SetUuid(color.GetId().AsByteString()).SetPropertiesKeepOriginalName(color.GetPropertiesKeepOriginalName()).AddAllInstanceType(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.AddColors_1_Color(color.GetInstanceColorsArray()))).AddAllPrototype(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.AddColors_1_Color(color.GetPrototypesArray()))).SetMarkedConstructor(color.IsConstructor()).AddAllOwnProperty(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.priv_ColorSerializer_getOwnPropertyStringPoolOffsets_7707(color))).SetClosureAssert(color.IsClosureAssert())
	}
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](typeProtoBuilder.Build())
}

// Java method: ColorSerializer.getOwnPropertyStringPoolOffsets(com.google.javascript.jscomp.colors.Color color)
func (this_ *ColorSerializer) priv_ColorSerializer_getOwnPropertyStringPoolOffsets_7707(color *pkg_colors.Color) j2gort.Ext_com_google_common_collect_ImmutableList {
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
			if this_.propertyFilter.Test(ownProperty) {
				if builder == nil {
					builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
				}
				builder.AddVoid_(j2gort.ValueOfInt(this_.getStringPoolIndexFn.Apply(ownProperty)))
			}
		}
	}
	if builder != nil {
		return builder.BuildImmList_()
	} else {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
}

// Java interface: com.google.javascript.jscomp.serialization.ColorSerializer.StringPoolIndexFn
type ColorSerializer_StringPoolIndexFn interface {
	Implements_jscomp_ColorSerializer_StringPoolIndexFn()
	IsJavaNil_() bool
	Apply(string) int32
}

func ToIface_ColorSerializer_StringPoolIndexFn[T interface {
	ColorSerializer_StringPoolIndexFn
	comparable
}](v T) ColorSerializer_StringPoolIndexFn {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ColorSerializer.StringPoolIndexFn", "com.google.javascript.jscomp.serialization.ColorSerializer$StringPoolIndexFn", func(v_ any) bool {
		return j2gort.InstanceOf[ColorSerializer_StringPoolIndexFn](v_)
	})
}

type lambda_ColorSerializer_StringPoolIndexFn struct {
	fn func(string) int32
}

func (l *lambda_ColorSerializer_StringPoolIndexFn) IsJavaNil_() bool { return l == nil }
func NewLambda_ColorSerializer_StringPoolIndexFn(fn func(string) int32) ColorSerializer_StringPoolIndexFn {
	if fn == nil {
		return nil
	}
	return &lambda_ColorSerializer_StringPoolIndexFn{fn: fn}
}

func (*lambda_ColorSerializer_StringPoolIndexFn) Implements_jscomp_ColorSerializer_StringPoolIndexFn() {
}
func (l *lambda_ColorSerializer_StringPoolIndexFn) Apply(p0_ string) int32 {
	return l.fn(p0_)
}

// Java interface: com.google.javascript.jscomp.serialization.ColorSerializer.PropertyFilter
type ColorSerializer_PropertyFilter interface {
	Implements_jscomp_ColorSerializer_PropertyFilter()
	IsJavaNil_() bool
	Test(string) bool
}

func ToIface_ColorSerializer_PropertyFilter[T interface {
	ColorSerializer_PropertyFilter
	comparable
}](v T) ColorSerializer_PropertyFilter {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ColorSerializer.PropertyFilter", "com.google.javascript.jscomp.serialization.ColorSerializer$PropertyFilter", func(v_ any) bool {
		return j2gort.InstanceOf[ColorSerializer_PropertyFilter](v_)
	})
}

type lambda_ColorSerializer_PropertyFilter struct {
	fn func(string) bool
}

func (l *lambda_ColorSerializer_PropertyFilter) IsJavaNil_() bool { return l == nil }
func NewLambda_ColorSerializer_PropertyFilter(fn func(string) bool) ColorSerializer_PropertyFilter {
	if fn == nil {
		return nil
	}
	return &lambda_ColorSerializer_PropertyFilter{fn: fn}
}

func (*lambda_ColorSerializer_PropertyFilter) Implements_jscomp_ColorSerializer_PropertyFilter() {}
func (l *lambda_ColorSerializer_PropertyFilter) Test(p0_ string) bool {
	return l.fn(p0_)
}

// Java interface: com.google.javascript.jscomp.serialization.ColorSerializer.DisambiguationSupertypesFn
type ColorSerializer_DisambiguationSupertypesFn interface {
	Implements_jscomp_ColorSerializer_DisambiguationSupertypesFn()
	IsJavaNil_() bool
	Apply(*pkg_colors.Color) []*pkg_colors.Color
}

func ToIface_ColorSerializer_DisambiguationSupertypesFn[T interface {
	ColorSerializer_DisambiguationSupertypesFn
	comparable
}](v T) ColorSerializer_DisambiguationSupertypesFn {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ColorSerializer.DisambiguationSupertypesFn", "com.google.javascript.jscomp.serialization.ColorSerializer$DisambiguationSupertypesFn", func(v_ any) bool {
		return j2gort.InstanceOf[ColorSerializer_DisambiguationSupertypesFn](v_)
	})
}

type lambda_ColorSerializer_DisambiguationSupertypesFn struct {
	fn func(*pkg_colors.Color) []*pkg_colors.Color
}

func (l *lambda_ColorSerializer_DisambiguationSupertypesFn) IsJavaNil_() bool { return l == nil }
func NewLambda_ColorSerializer_DisambiguationSupertypesFn(fn func(*pkg_colors.Color) []*pkg_colors.Color) ColorSerializer_DisambiguationSupertypesFn {
	if fn == nil {
		return nil
	}
	return &lambda_ColorSerializer_DisambiguationSupertypesFn{fn: fn}
}

func (*lambda_ColorSerializer_DisambiguationSupertypesFn) Implements_jscomp_ColorSerializer_DisambiguationSupertypesFn() {
}
func (l *lambda_ColorSerializer_DisambiguationSupertypesFn) Apply(p0_ *pkg_colors.Color) []*pkg_colors.Color {
	return l.fn(p0_)
}
