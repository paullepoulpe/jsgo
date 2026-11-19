/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/disambiguate/ColorGraphNode.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.disambiguate.ColorGraphNode
type ColorGraphNode struct {
	color           *pkg_colors.Color                                                      `gson:"color"`
	associatedProps j2gort.OrderedMap[*PropertyClustering, ColorGraphNode_PropAssociation] `gson:"associatedProps"`
	subtypeIndices  *j2gort.Ext_java_util_BitSet                                           `gson:"subtypeIndices"`
	index           int32                                                                  `gson:"index"`
	idHash_         atomic.Int32
}

func (this_ *ColorGraphNode) BindSelf_(s any) {
	_ = s
}
func (this_ *ColorGraphNode) IsJavaNil_() bool { return this_ == nil }
func (*ColorGraphNode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.ColorGraphNode", "ColorGraphNode"
}
func (this_ *ColorGraphNode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ColorGraphNode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.ColorGraphNode", "com.google.javascript.jscomp.disambiguate.ColorGraphNode", func(v_ any) bool {
		return j2gort.InstanceOf[*ColorGraphNode](v_)
	})
}

func (this_ *ColorGraphNode) initFields_ColorGraphNode() {
	_ = this_
	this_.associatedProps = j2gort.NewOrderedMap[*PropertyClustering, ColorGraphNode_PropAssociation]()
	this_.subtypeIndices = j2gort.NewExt_java_util_BitSet()
}

// Java constructor: ColorGraphNode(com.google.javascript.jscomp.colors.Color single, int index)
func NewColorGraphNode(single *pkg_colors.Color, index int32) *ColorGraphNode {
	this_ := &ColorGraphNode{}
	InitColorGraphNode(this_, single, index)
	return this_
}

func InitColorGraphNode(this_ *ColorGraphNode, single *pkg_colors.Color, index int32) {
	_ = single
	_ = index
	this_.initFields_ColorGraphNode()
	j2gort.CheckNotNullPtr(single)
	this_.index = index
	this_.color = single
}

// Java static method: ColorGraphNode.create(com.google.javascript.jscomp.colors.Color single, int index)
func ColorGraphNode_Create(single *pkg_colors.Color, index int32) *ColorGraphNode {
	_ = single
	_ = index
	if !(index >= int32(0)) {
		j2gort.CheckArgument(false)
	}
	return NewColorGraphNode(single, index)
}

// Java static method: ColorGraphNode.createForTesting(int index)
func ColorGraphNode_CreateForTesting_1(index int32) *ColorGraphNode {
	_ = index
	if !(index < int32(0)) {
		j2gort.CheckArgument(false)
	}
	return NewColorGraphNode(pkg_colors.StaticGet_StandardColors_UNKNOWN(), index)
}

// Java static method: ColorGraphNode.createForTesting(com.google.javascript.jscomp.colors.Color single, int index)
func ColorGraphNode_CreateForTesting_2(single *pkg_colors.Color, index int32) *ColorGraphNode {
	_ = single
	_ = index
	if !(index < int32(0)) {
		j2gort.CheckArgument(false)
	}
	return NewColorGraphNode(single, index)
}

// Java method: ColorGraphNode.getColor()
func (this_ *ColorGraphNode) GetColor() *pkg_colors.Color {
	return j2gort.CheckNotNullPtr(this_.color)
}

// Java method: ColorGraphNode.getIndex()
func (this_ *ColorGraphNode) GetIndex() int32 {
	return this_.index
}

// Java method: ColorGraphNode.getAssociatedProps()
func (this_ *ColorGraphNode) GetAssociatedProps() *j2gort.OrderedMap[*PropertyClustering, ColorGraphNode_PropAssociation] {
	return &this_.associatedProps

}

// Java method: ColorGraphNode.getSubtypeIndices()
func (this_ *ColorGraphNode) GetSubtypeIndices() *j2gort.Ext_java_util_BitSet {
	return this_.subtypeIndices
}

// Java method: ColorGraphNode.toString()
func (this_ *ColorGraphNode) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_int32("index", this_.index).Add_2_string_any("color.id", j2gort.BoxNullableVal(this_.color.GetId())).Add_2_string_any("color.ownProperties", j2gort.BoxString(j2gort.SliceFormatJavaString(this_.color.GetOwnPropertiesArray()))).String()
}

// Java enum: com.google.javascript.jscomp.disambiguate.ColorGraphNode.PropAssociation
type ColorGraphNode_PropAssociation uint8

const (
	enum_ColorGraphNode_PropAssociation_AST         ColorGraphNode_PropAssociation = 1
	enum_ColorGraphNode_PropAssociation_TYPE_SYSTEM ColorGraphNode_PropAssociation = 2
	enum_ColorGraphNode_PropAssociation_SUPERTYPE   ColorGraphNode_PropAssociation = 3
)

var enumBaseTable_ColorGraphNode_PropAssociation = [...]j2gort.EnumBase{
	{Name_: "AST", Ordinal_: 0},
	{Name_: "TYPE_SYSTEM", Ordinal_: 1},
	{Name_: "SUPERTYPE", Ordinal_: 2},
}
var enumValues_ColorGraphNode_PropAssociation = []ColorGraphNode_PropAssociation{1, 2, 3}

func StaticGet_ColorGraphNode_PropAssociation_AST() ColorGraphNode_PropAssociation {
	return enum_ColorGraphNode_PropAssociation_AST
}
func StaticGet_ColorGraphNode_PropAssociation_TYPE_SYSTEM() ColorGraphNode_PropAssociation {
	return enum_ColorGraphNode_PropAssociation_TYPE_SYSTEM
}
func StaticGet_ColorGraphNode_PropAssociation_SUPERTYPE() ColorGraphNode_PropAssociation {
	return enum_ColorGraphNode_PropAssociation_SUPERTYPE
}

func ColorGraphNode_PropAssociation_Values() []ColorGraphNode_PropAssociation {
	return j2gort.CloneSlice(enumValues_ColorGraphNode_PropAssociation)
}

func ColorGraphNode_PropAssociation_ValueOf(name string) ColorGraphNode_PropAssociation {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ColorGraphNode_PropAssociation {
		if enumBaseTable_ColorGraphNode_PropAssociation[i_].Name_ == name {
			return ColorGraphNode_PropAssociation(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.disambiguate.ColorGraphNode.PropAssociation." + name))
}

func (this_ ColorGraphNode_PropAssociation) IsJavaNil_() bool                              { return this_ == 0 }
func (ColorGraphNode_PropAssociation) HasIdentityEquals_()                                 {}
func (ColorGraphNode_PropAssociation) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ColorGraphNode_PropAssociation) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ColorGraphNode_PropAssociation) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ColorGraphNode_PropAssociation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.ColorGraphNode$PropAssociation", "PropAssociation"
}
func (this_ ColorGraphNode_PropAssociation) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ColorGraphNode_PropAssociation[this_-1]
}
func (this_ ColorGraphNode_PropAssociation) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ColorGraphNode_PropAssociation) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ColorGraphNode_PropAssociation) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ColorGraphNode_PropAssociation[this_-1].Name_
}
func (this_ ColorGraphNode_PropAssociation) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ColorGraphNode_PropAssociation[this_-1].Name_
}
func (this_ ColorGraphNode_PropAssociation) Equals(other_ any) bool {
	o_, ok_ := other_.(ColorGraphNode_PropAssociation)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ColorGraphNode_PropAssociation) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ColorGraphNode_PropAssociation[this_-1].HashCode()
}
func (this_ ColorGraphNode_PropAssociation) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ColorGraphNode_PropAssociation[this_-1].CompareTo(other_)
}
func (this_ ColorGraphNode_PropAssociation) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ColorGraphNode_PropAssociation[this_-1].Name_)
}
func init() {
	vals_ := []any{ColorGraphNode_PropAssociation(1), ColorGraphNode_PropAssociation(2), ColorGraphNode_PropAssociation(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.disambiguate.ColorGraphNode.PropAssociation", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.disambiguate.ColorGraphNode$PropAssociation", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.ColorGraphNode.PropAssociation", "com.google.javascript.jscomp.disambiguate.ColorGraphNode$PropAssociation", func(v_ any) bool {
		return j2gort.InstanceOf[ColorGraphNode_PropAssociation](v_)
	})
}
