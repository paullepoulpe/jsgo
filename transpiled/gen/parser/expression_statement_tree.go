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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ExpressionStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ExpressionStatementTree
type ExpressionStatementTree struct {
	ParseTreeBase
	expression ParseTree `gson:"expression"`
}

func (this_ *ExpressionStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ExpressionStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*ExpressionStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ExpressionStatementTree", "ExpressionStatementTree"
}
func (*ExpressionStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ExpressionStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.ExpressionStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ExpressionStatementTree](v_)
	})
}
func (*ExpressionStatementTree) Implements_parser_ParseTree()         {}
func (this_ *ExpressionStatementTree) FieldGet_expression() ParseTree { return this_.expression }

func (this_ *ExpressionStatementTree) initFields_ExpressionStatementTree() {
	_ = this_
}

// Java constructor: ExpressionStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression)
func NewExpressionStatementTree(location SourceRange, expression ParseTree) *ExpressionStatementTree {
	this_ := &ExpressionStatementTree{}
	this_.BindSelf_(this_)
	InitExpressionStatementTree(this_, location, expression)
	return this_
}
func allocPooled_ExpressionStatementTree(slab_ *j2gort.PooledSlab[ExpressionStatementTree], location SourceRange, expression ParseTree) *ExpressionStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitExpressionStatementTree(this_, location, expression)
	return this_
}

func InitExpressionStatementTree(this_ *ExpressionStatementTree, location SourceRange, expression ParseTree) {
	_ = location
	_ = expression
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_EXPRESSION_STATEMENT(), location)
	this_.initFields_ExpressionStatementTree()
	this_.expression = expression
}
