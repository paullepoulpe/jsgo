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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ParenExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ParenExpressionTree
type ParenExpressionTree struct {
	ParseTreeBase
	expression ParseTree `gson:"expression"`
}

func (this_ *ParenExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ParenExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*ParenExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ParenExpressionTree", "ParenExpressionTree"
}
func (*ParenExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ParenExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.ParenExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ParenExpressionTree](v_)
	})
}
func (*ParenExpressionTree) Implements_parser_ParseTree()         {}
func (this_ *ParenExpressionTree) FieldGet_expression() ParseTree { return this_.expression }

func (this_ *ParenExpressionTree) initFields_ParenExpressionTree() {
	_ = this_
}

// Java constructor: ParenExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression)
func NewParenExpressionTree(location SourceRange, expression ParseTree) *ParenExpressionTree {
	this_ := &ParenExpressionTree{}
	this_.BindSelf_(this_)
	InitParenExpressionTree(this_, location, expression)
	return this_
}
func allocPooled_ParenExpressionTree(slab_ *j2gort.PooledSlab[ParenExpressionTree], location SourceRange, expression ParseTree) *ParenExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitParenExpressionTree(this_, location, expression)
	return this_
}

func InitParenExpressionTree(this_ *ParenExpressionTree, location SourceRange, expression ParseTree) {
	_ = location
	_ = expression
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_PAREN_EXPRESSION(), location)
	this_.initFields_ParenExpressionTree()
	this_.expression = expression
}
