/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.TemplateLiteralExpressionTree
type TemplateLiteralExpressionTree struct {
	ParseTreeBase
	operand  ParseTree                                          `gson:"operand"`
	elements j2gort.Ext_com_google_common_collect_ImmutableList `gson:"elements"`
}

func (this_ *TemplateLiteralExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *TemplateLiteralExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*TemplateLiteralExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.TemplateLiteralExpressionTree", "TemplateLiteralExpressionTree"
}
func (*TemplateLiteralExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.TemplateLiteralExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.TemplateLiteralExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*TemplateLiteralExpressionTree](v_)
	})
}
func (*TemplateLiteralExpressionTree) Implements_parser_ParseTree()      {}
func (this_ *TemplateLiteralExpressionTree) FieldGet_operand() ParseTree { return this_.operand }
func (this_ *TemplateLiteralExpressionTree) FieldGet_elements() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.elements
}

func (this_ *TemplateLiteralExpressionTree) initFields_TemplateLiteralExpressionTree() {
	_ = this_
}

// Java constructor: TemplateLiteralExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> elements)
func NewTemplateLiteralExpressionTree(location SourceRange, operand ParseTree, elements j2gort.Ext_com_google_common_collect_ImmutableList) *TemplateLiteralExpressionTree {
	this_ := &TemplateLiteralExpressionTree{}
	this_.BindSelf_(this_)
	InitTemplateLiteralExpressionTree(this_, location, operand, elements)
	return this_
}
func allocPooled_TemplateLiteralExpressionTree(slab_ *j2gort.PooledSlab[TemplateLiteralExpressionTree], location SourceRange, operand ParseTree, elements j2gort.Ext_com_google_common_collect_ImmutableList) *TemplateLiteralExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitTemplateLiteralExpressionTree(this_, location, operand, elements)
	return this_
}

func InitTemplateLiteralExpressionTree(this_ *TemplateLiteralExpressionTree, location SourceRange, operand ParseTree, elements j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = operand
	_ = elements
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_TEMPLATE_LITERAL_EXPRESSION(), location)
	this_.initFields_TemplateLiteralExpressionTree()
	this_.operand = operand
	this_.elements = elements
}
