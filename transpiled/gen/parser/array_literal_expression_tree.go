/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayLiteralExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ArrayLiteralExpressionTree
type ArrayLiteralExpressionTree struct {
	ParseTreeBase
	elements         j2gort.Ext_com_google_common_collect_ImmutableList `gson:"elements"`
	hasTrailingComma bool                                               `gson:"hasTrailingComma"`
}

func (this_ *ArrayLiteralExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ArrayLiteralExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*ArrayLiteralExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ArrayLiteralExpressionTree", "ArrayLiteralExpressionTree"
}
func (*ArrayLiteralExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ArrayLiteralExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.ArrayLiteralExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ArrayLiteralExpressionTree](v_)
	})
}
func (*ArrayLiteralExpressionTree) Implements_parser_ParseTree() {}
func (this_ *ArrayLiteralExpressionTree) FieldGet_elements() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.elements
}
func (this_ *ArrayLiteralExpressionTree) FieldGet_hasTrailingComma() bool {
	return this_.hasTrailingComma
}

func (this_ *ArrayLiteralExpressionTree) initFields_ArrayLiteralExpressionTree() {
	_ = this_
}

// Java constructor: ArrayLiteralExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> elements, boolean hasTrailingComma)
func NewArrayLiteralExpressionTree(location SourceRange, elements j2gort.Ext_com_google_common_collect_ImmutableList, hasTrailingComma bool) *ArrayLiteralExpressionTree {
	this_ := &ArrayLiteralExpressionTree{}
	this_.BindSelf_(this_)
	InitArrayLiteralExpressionTree(this_, location, elements, hasTrailingComma)
	return this_
}
func allocPooled_ArrayLiteralExpressionTree(slab_ *j2gort.PooledSlab[ArrayLiteralExpressionTree], location SourceRange, elements j2gort.Ext_com_google_common_collect_ImmutableList, hasTrailingComma bool) *ArrayLiteralExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitArrayLiteralExpressionTree(this_, location, elements, hasTrailingComma)
	return this_
}

func InitArrayLiteralExpressionTree(this_ *ArrayLiteralExpressionTree, location SourceRange, elements j2gort.Ext_com_google_common_collect_ImmutableList, hasTrailingComma bool) {
	_ = location
	_ = elements
	_ = hasTrailingComma
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_ARRAY_LITERAL_EXPRESSION(), location)
	this_.initFields_ArrayLiteralExpressionTree()
	this_.elements = elements
	this_.hasTrailingComma = hasTrailingComma
}
