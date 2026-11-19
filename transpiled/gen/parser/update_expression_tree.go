/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/trees/UpdateExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree
type UpdateExpressionTree struct {
	ParseTreeBase
	operator         Token                                 `gson:"operator"`
	operand          ParseTree                             `gson:"operand"`
	operatorPosition UpdateExpressionTree_OperatorPosition `gson:"operatorPosition"`
}

func (this_ *UpdateExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *UpdateExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*UpdateExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree", "UpdateExpressionTree"
}
func (*UpdateExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*UpdateExpressionTree](v_)
	})
}
func (*UpdateExpressionTree) Implements_parser_ParseTree()   {}
func (this_ *UpdateExpressionTree) FieldGet_operator() Token { return this_.operator }
func (this_ *UpdateExpressionTree) FieldGet_operatorPosition() UpdateExpressionTree_OperatorPosition {
	return this_.operatorPosition
}
func (this_ *UpdateExpressionTree) FieldGet_operand() ParseTree { return this_.operand }

func (this_ *UpdateExpressionTree) initFields_UpdateExpressionTree() {
	_ = this_
}

// Java constructor: UpdateExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token operator, com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree.OperatorPosition operatorPosition, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand)
func NewUpdateExpressionTree(location SourceRange, operator Token, operatorPosition UpdateExpressionTree_OperatorPosition, operand ParseTree) *UpdateExpressionTree {
	this_ := &UpdateExpressionTree{}
	this_.BindSelf_(this_)
	InitUpdateExpressionTree(this_, location, operator, operatorPosition, operand)
	return this_
}
func allocPooled_UpdateExpressionTree(slab_ *j2gort.PooledSlab[UpdateExpressionTree], location SourceRange, operator Token, operatorPosition UpdateExpressionTree_OperatorPosition, operand ParseTree) *UpdateExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitUpdateExpressionTree(this_, location, operator, operatorPosition, operand)
	return this_
}

func InitUpdateExpressionTree(this_ *UpdateExpressionTree, location SourceRange, operator Token, operatorPosition UpdateExpressionTree_OperatorPosition, operand ParseTree) {
	_ = location
	_ = operator
	_ = operatorPosition
	_ = operand
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_UPDATE_EXPRESSION(), location)
	this_.initFields_UpdateExpressionTree()
	this_.operator = operator
	this_.operatorPosition = operatorPosition
	this_.operand = operand
}

// Java static method: UpdateExpressionTree.prefix(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token operator, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand)
func UpdateExpressionTree_Prefix(location SourceRange, operator Token, operand ParseTree) *UpdateExpressionTree {
	_ = location
	_ = operator
	_ = operand
	return NewUpdateExpressionTree(location, operator, StaticGet_UpdateExpressionTree_OperatorPosition_PREFIX(), operand)
}
func allocPooledVia_UpdateExpressionTree_Prefix(slab_ *j2gort.PooledSlab[UpdateExpressionTree], location SourceRange, operator Token, operand ParseTree) *UpdateExpressionTree {
	_ = location
	_ = operator
	_ = operand
	return allocPooled_UpdateExpressionTree(slab_, location, operator, StaticGet_UpdateExpressionTree_OperatorPosition_PREFIX(), operand)
}

// Java static method: UpdateExpressionTree.postfix(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token operator, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand)
func UpdateExpressionTree_Postfix(location SourceRange, operator Token, operand ParseTree) *UpdateExpressionTree {
	_ = location
	_ = operator
	_ = operand
	return NewUpdateExpressionTree(location, operator, StaticGet_UpdateExpressionTree_OperatorPosition_POSTFIX(), operand)
}
func allocPooledVia_UpdateExpressionTree_Postfix(slab_ *j2gort.PooledSlab[UpdateExpressionTree], location SourceRange, operator Token, operand ParseTree) *UpdateExpressionTree {
	_ = location
	_ = operator
	_ = operand
	return allocPooled_UpdateExpressionTree(slab_, location, operator, StaticGet_UpdateExpressionTree_OperatorPosition_POSTFIX(), operand)
}

// Java enum: com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree.OperatorPosition
type UpdateExpressionTree_OperatorPosition uint8

const (
	enum_UpdateExpressionTree_OperatorPosition_PREFIX  UpdateExpressionTree_OperatorPosition = 1
	enum_UpdateExpressionTree_OperatorPosition_POSTFIX UpdateExpressionTree_OperatorPosition = 2
)

var enumBaseTable_UpdateExpressionTree_OperatorPosition = [...]j2gort.EnumBase{
	{Name_: "PREFIX", Ordinal_: 0},
	{Name_: "POSTFIX", Ordinal_: 1},
}
var enumValues_UpdateExpressionTree_OperatorPosition = []UpdateExpressionTree_OperatorPosition{1, 2}

func StaticGet_UpdateExpressionTree_OperatorPosition_PREFIX() UpdateExpressionTree_OperatorPosition {
	return enum_UpdateExpressionTree_OperatorPosition_PREFIX
}
func StaticGet_UpdateExpressionTree_OperatorPosition_POSTFIX() UpdateExpressionTree_OperatorPosition {
	return enum_UpdateExpressionTree_OperatorPosition_POSTFIX
}

func UpdateExpressionTree_OperatorPosition_Values() []UpdateExpressionTree_OperatorPosition {
	return j2gort.CloneSlice(enumValues_UpdateExpressionTree_OperatorPosition)
}

func UpdateExpressionTree_OperatorPosition_ValueOf(name string) UpdateExpressionTree_OperatorPosition {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_UpdateExpressionTree_OperatorPosition {
		if enumBaseTable_UpdateExpressionTree_OperatorPosition[i_].Name_ == name {
			return UpdateExpressionTree_OperatorPosition(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree.OperatorPosition." + name))
}

func (this_ UpdateExpressionTree_OperatorPosition) IsJavaNil_() bool                              { return this_ == 0 }
func (UpdateExpressionTree_OperatorPosition) HasIdentityEquals_()                                 {}
func (UpdateExpressionTree_OperatorPosition) Implements_j2gort_Ext_java_lang_Enum()               {}
func (UpdateExpressionTree_OperatorPosition) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ UpdateExpressionTree_OperatorPosition) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (UpdateExpressionTree_OperatorPosition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree$OperatorPosition", "OperatorPosition"
}
func (this_ UpdateExpressionTree_OperatorPosition) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_UpdateExpressionTree_OperatorPosition[this_-1]
}
func (this_ UpdateExpressionTree_OperatorPosition) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ UpdateExpressionTree_OperatorPosition) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ UpdateExpressionTree_OperatorPosition) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_UpdateExpressionTree_OperatorPosition[this_-1].Name_
}
func (this_ UpdateExpressionTree_OperatorPosition) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_UpdateExpressionTree_OperatorPosition[this_-1].Name_
}
func (this_ UpdateExpressionTree_OperatorPosition) Equals(other_ any) bool {
	o_, ok_ := other_.(UpdateExpressionTree_OperatorPosition)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ UpdateExpressionTree_OperatorPosition) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_UpdateExpressionTree_OperatorPosition[this_-1].HashCode()
}
func (this_ UpdateExpressionTree_OperatorPosition) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_UpdateExpressionTree_OperatorPosition[this_-1].CompareTo(other_)
}
func (this_ UpdateExpressionTree_OperatorPosition) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_UpdateExpressionTree_OperatorPosition[this_-1].Name_)
}
func init() {
	vals_ := []any{UpdateExpressionTree_OperatorPosition(1), UpdateExpressionTree_OperatorPosition(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree.OperatorPosition", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree$OperatorPosition", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree.OperatorPosition", "com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree$OperatorPosition", func(v_ any) bool {
		return j2gort.InstanceOf[UpdateExpressionTree_OperatorPosition](v_)
	})
}
