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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ConditionalExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ConditionalExpressionTree
type ConditionalExpressionTree struct {
	ParseTreeBase
	condition ParseTree `gson:"condition"`
	left      ParseTree `gson:"left"`
	right     ParseTree `gson:"right"`
}

func (this_ *ConditionalExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ConditionalExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*ConditionalExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ConditionalExpressionTree", "ConditionalExpressionTree"
}
func (*ConditionalExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ConditionalExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.ConditionalExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ConditionalExpressionTree](v_)
	})
}
func (*ConditionalExpressionTree) Implements_parser_ParseTree()        {}
func (this_ *ConditionalExpressionTree) FieldGet_condition() ParseTree { return this_.condition }
func (this_ *ConditionalExpressionTree) FieldGet_left() ParseTree      { return this_.left }
func (this_ *ConditionalExpressionTree) FieldGet_right() ParseTree     { return this_.right }

func (this_ *ConditionalExpressionTree) initFields_ConditionalExpressionTree() {
	_ = this_
}

// Java constructor: ConditionalExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree condition, com.google.javascript.jscomp.parsing.parser.trees.ParseTree left, com.google.javascript.jscomp.parsing.parser.trees.ParseTree right)
func NewConditionalExpressionTree(location SourceRange, condition ParseTree, left ParseTree, right ParseTree) *ConditionalExpressionTree {
	this_ := &ConditionalExpressionTree{}
	this_.BindSelf_(this_)
	InitConditionalExpressionTree(this_, location, condition, left, right)
	return this_
}
func allocPooled_ConditionalExpressionTree(slab_ *j2gort.PooledSlab[ConditionalExpressionTree], location SourceRange, condition ParseTree, left ParseTree, right ParseTree) *ConditionalExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitConditionalExpressionTree(this_, location, condition, left, right)
	return this_
}

func InitConditionalExpressionTree(this_ *ConditionalExpressionTree, location SourceRange, condition ParseTree, left ParseTree, right ParseTree) {
	_ = location
	_ = condition
	_ = left
	_ = right
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_CONDITIONAL_EXPRESSION(), location)
	this_.initFields_ConditionalExpressionTree()
	this_.condition = condition
	this_.left = left
	this_.right = right
}
