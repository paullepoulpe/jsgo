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
// Port of com/google/javascript/jscomp/parsing/parser/trees/CallExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.CallExpressionTree
type CallExpressionTree struct {
	ParseTreeBase
	operand   ParseTree         `gson:"operand"`
	arguments *ArgumentListTree `gson:"arguments"`
}

func (this_ *CallExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *CallExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*CallExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.CallExpressionTree", "CallExpressionTree"
}
func (*CallExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.CallExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.CallExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*CallExpressionTree](v_)
	})
}
func (*CallExpressionTree) Implements_parser_ParseTree()                {}
func (this_ *CallExpressionTree) FieldGet_operand() ParseTree           { return this_.operand }
func (this_ *CallExpressionTree) FieldGet_arguments() *ArgumentListTree { return this_.arguments }

func (this_ *CallExpressionTree) initFields_CallExpressionTree() {
	_ = this_
}

// Java constructor: CallExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand, com.google.javascript.jscomp.parsing.parser.trees.ArgumentListTree arguments)
func NewCallExpressionTree(location SourceRange, operand ParseTree, arguments *ArgumentListTree) *CallExpressionTree {
	this_ := &CallExpressionTree{}
	this_.BindSelf_(this_)
	InitCallExpressionTree(this_, location, operand, arguments)
	return this_
}
func allocPooled_CallExpressionTree(slab_ *j2gort.PooledSlab[CallExpressionTree], location SourceRange, operand ParseTree, arguments *ArgumentListTree) *CallExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitCallExpressionTree(this_, location, operand, arguments)
	return this_
}

func InitCallExpressionTree(this_ *CallExpressionTree, location SourceRange, operand ParseTree, arguments *ArgumentListTree) {
	_ = location
	_ = operand
	_ = arguments
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_CALL_EXPRESSION(), location)
	this_.initFields_CallExpressionTree()
	this_.operand = operand
	this_.arguments = arguments
}
