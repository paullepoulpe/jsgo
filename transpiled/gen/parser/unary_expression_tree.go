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
// Port of com/google/javascript/jscomp/parsing/parser/trees/UnaryExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.UnaryExpressionTree
type UnaryExpressionTree struct {
	ParseTreeBase
	operator Token     `gson:"operator"`
	operand  ParseTree `gson:"operand"`
}

func (this_ *UnaryExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *UnaryExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*UnaryExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.UnaryExpressionTree", "UnaryExpressionTree"
}
func (*UnaryExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.UnaryExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.UnaryExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*UnaryExpressionTree](v_)
	})
}
func (*UnaryExpressionTree) Implements_parser_ParseTree()      {}
func (this_ *UnaryExpressionTree) FieldGet_operator() Token    { return this_.operator }
func (this_ *UnaryExpressionTree) FieldGet_operand() ParseTree { return this_.operand }

func (this_ *UnaryExpressionTree) initFields_UnaryExpressionTree() {
	_ = this_
}

// Java constructor: UnaryExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token operator, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand)
func NewUnaryExpressionTree(location SourceRange, operator Token, operand ParseTree) *UnaryExpressionTree {
	this_ := &UnaryExpressionTree{}
	this_.BindSelf_(this_)
	InitUnaryExpressionTree(this_, location, operator, operand)
	return this_
}
func allocPooled_UnaryExpressionTree(slab_ *j2gort.PooledSlab[UnaryExpressionTree], location SourceRange, operator Token, operand ParseTree) *UnaryExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitUnaryExpressionTree(this_, location, operator, operand)
	return this_
}

func InitUnaryExpressionTree(this_ *UnaryExpressionTree, location SourceRange, operator Token, operand ParseTree) {
	_ = location
	_ = operator
	_ = operand
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_UNARY_EXPRESSION(), location)
	this_.initFields_UnaryExpressionTree()
	this_.operator = operator
	this_.operand = operand
}
