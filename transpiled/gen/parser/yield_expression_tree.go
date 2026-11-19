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
// Port of com/google/javascript/jscomp/parsing/parser/trees/YieldExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.YieldExpressionTree
type YieldExpressionTree struct {
	ParseTreeBase
	expression ParseTree `gson:"expression"`
	isYieldAll bool      `gson:"isYieldAll"`
}

func (this_ *YieldExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *YieldExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*YieldExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.YieldExpressionTree", "YieldExpressionTree"
}
func (*YieldExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.YieldExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.YieldExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*YieldExpressionTree](v_)
	})
}
func (*YieldExpressionTree) Implements_parser_ParseTree()         {}
func (this_ *YieldExpressionTree) FieldGet_expression() ParseTree { return this_.expression }
func (this_ *YieldExpressionTree) FieldGet_isYieldAll() bool      { return this_.isYieldAll }

func (this_ *YieldExpressionTree) initFields_YieldExpressionTree() {
	_ = this_
}

// Java constructor: YieldExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, boolean isYieldAll, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression)
func NewYieldExpressionTree(location SourceRange, isYieldAll bool, expression ParseTree) *YieldExpressionTree {
	this_ := &YieldExpressionTree{}
	this_.BindSelf_(this_)
	InitYieldExpressionTree(this_, location, isYieldAll, expression)
	return this_
}
func allocPooled_YieldExpressionTree(slab_ *j2gort.PooledSlab[YieldExpressionTree], location SourceRange, isYieldAll bool, expression ParseTree) *YieldExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitYieldExpressionTree(this_, location, isYieldAll, expression)
	return this_
}

func InitYieldExpressionTree(this_ *YieldExpressionTree, location SourceRange, isYieldAll bool, expression ParseTree) {
	_ = location
	_ = isYieldAll
	_ = expression
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_YIELD_EXPRESSION(), location)
	this_.initFields_YieldExpressionTree()
	this_.isYieldAll = isYieldAll
	this_.expression = expression
}
