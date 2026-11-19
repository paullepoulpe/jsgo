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
// Port of com/google/javascript/jscomp/disambiguate/ColorGraphNodeFactory.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.disambiguate.ColorGraphNodeFactory
type ColorGraphNodeFactory struct {
	typeIndex *j2gort.CustomOrderedMap[*pkg_colors.Color, *ColorGraphNode] `gson:"typeIndex"`
	registry  *pkg_colors.ColorRegistry                                    `gson:"registry"`
	idHash_   atomic.Int32
}

func (this_ *ColorGraphNodeFactory) BindSelf_(s any) {
	_ = s
}
func (this_ *ColorGraphNodeFactory) IsJavaNil_() bool { return this_ == nil }
func (*ColorGraphNodeFactory) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.ColorGraphNodeFactory", "ColorGraphNodeFactory"
}
func (this_ *ColorGraphNodeFactory) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ColorGraphNodeFactory) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.ColorGraphNodeFactory", "com.google.javascript.jscomp.disambiguate.ColorGraphNodeFactory", func(v_ any) bool {
		return j2gort.InstanceOf[*ColorGraphNodeFactory](v_)
	})
}

func (this_ *ColorGraphNodeFactory) initFields_ColorGraphNodeFactory() {
	_ = this_
}

// Java constructor: ColorGraphNodeFactory(java.util.LinkedHashMap<com.google.javascript.jscomp.colors.Color,com.google.javascript.jscomp.disambiguate.ColorGraphNode> initialTypeIndex, com.google.javascript.jscomp.colors.ColorRegistry registry)
func NewColorGraphNodeFactory(initialTypeIndex *j2gort.CustomOrderedMap[*pkg_colors.Color, *ColorGraphNode], registry *pkg_colors.ColorRegistry) *ColorGraphNodeFactory {
	this_ := &ColorGraphNodeFactory{}
	InitColorGraphNodeFactory(this_, initialTypeIndex, registry)
	return this_
}

func InitColorGraphNodeFactory(this_ *ColorGraphNodeFactory, initialTypeIndex *j2gort.CustomOrderedMap[*pkg_colors.Color, *ColorGraphNode], registry *pkg_colors.ColorRegistry) {
	_ = initialTypeIndex
	_ = registry
	this_.initFields_ColorGraphNodeFactory()
	this_.typeIndex = initialTypeIndex
	this_.registry = registry
}

// Java static method: ColorGraphNodeFactory.createFactory(com.google.javascript.jscomp.colors.ColorRegistry colorRegistry)
func ColorGraphNodeFactory_CreateFactory(colorRegistry *pkg_colors.ColorRegistry) *ColorGraphNodeFactory {
	_ = colorRegistry
	var typeIndex j2gort.CustomOrderedMap[*pkg_colors.Color, *ColorGraphNode] = j2gort.NewCustomOrderedMap[*pkg_colors.Color, *ColorGraphNode]()
	_ = typeIndex
	var unknownColorNode *ColorGraphNode = ColorGraphNode_Create(pkg_colors.StaticGet_StandardColors_UNKNOWN(), int32(0))
	_ = unknownColorNode
	typeIndex.Put(pkg_colors.StaticGet_StandardColors_UNKNOWN(), unknownColorNode)
	return NewColorGraphNodeFactory(&typeIndex, colorRegistry)
}

// Java method: ColorGraphNodeFactory.createNode(com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable Color color)
func (this_ *ColorGraphNodeFactory) CreateNode(color *pkg_colors.Color) *ColorGraphNode {
	_ = color
	var key *pkg_colors.Color = this_.priv_ColorGraphNodeFactory_simplifyColor_6212(color)
	_ = key
	return func() *ColorGraphNode {
		bound_1 := this_
		return j2gort.CustomOrderedMapComputeIfAbsent(this_.typeIndex, key, nil, func(mrp0_ *pkg_colors.Color) *ColorGraphNode {
			return bound_1.priv_ColorGraphNodeFactory_newColorGraphNode_6211(mrp0_)
		})
	}()
}

// Java method: ColorGraphNodeFactory.getAllKnownTypes()
func (this_ *ColorGraphNodeFactory) GetAllKnownTypes() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.typeIndex.LiveValues()))
}

// Java method: ColorGraphNodeFactory.newColorGraphNode(com.google.javascript.jscomp.colors.Color key)
func (this_ *ColorGraphNodeFactory) priv_ColorGraphNodeFactory_newColorGraphNode_6211(key *pkg_colors.Color) *ColorGraphNode {
	_ = key
	var id int32 = this_.typeIndex.Size()
	_ = id
	return ColorGraphNode_Create(key, id)
}

// Java method: ColorGraphNodeFactory.simplifyColor(com.google.javascript.jscomp.colors.@org.jspecify.annotations.Nullable Color type)
func (this_ *ColorGraphNodeFactory) priv_ColorGraphNodeFactory_simplifyColor_6212(type_ *pkg_colors.Color) *pkg_colors.Color {
	_ = type_
	if type_ == nil {
		return pkg_colors.StaticGet_StandardColors_UNKNOWN()
	}
	if type_.IsUnion() {
		type_ = type_.SubtractNullOrVoid()
		if !(type_.IsUnion()) {
			return this_.priv_ColorGraphNodeFactory_simplifyColor_6212(type_)
		}
		var unionElems []*pkg_colors.Color = type_.GetUnionElementsArray()
		_ = unionElems
		if int32(len(unionElems)) == int32(2) {
			return pkg_colors.Color_CreateUnion2(this_.priv_ColorGraphNodeFactory_simplifyColor_6212(unionElems[int32(0)]), this_.priv_ColorGraphNodeFactory_simplifyColor_6212(unionElems[int32(1)]))
		}
		var simplified *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(int32(len(unionElems)))
		_ = simplified
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = int32(len(unionElems))
			_ = n
			for ; i < n; i = i + 1 {
				simplified.Add_1(j2gort.PtrToAny(this_.priv_ColorGraphNodeFactory_simplifyColor_6212(unionElems[i])))
			}
		}
		return pkg_colors.Color_CreateUnion_1_Ext_java_util_Set(simplified)
	} else if type_.GetBoxId() != (pkg_colors.ColorId{}) {
		return this_.registry.Get(type_.GetBoxId())
	} else if type_.Equals(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID())) {
		return pkg_colors.StaticGet_StandardColors_UNKNOWN()
	} else {
		return type_
	}
}
