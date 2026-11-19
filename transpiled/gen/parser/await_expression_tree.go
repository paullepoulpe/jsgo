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
// Port of com/google/javascript/jscomp/parsing/parser/trees/AwaitExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.AwaitExpressionTree
type AwaitExpressionTree struct {
	ParseTreeBase
	expression ParseTree `gson:"expression"`
}

func (this_ *AwaitExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *AwaitExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*AwaitExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.AwaitExpressionTree", "AwaitExpressionTree"
}
func (*AwaitExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.AwaitExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.AwaitExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*AwaitExpressionTree](v_)
	})
}
func (*AwaitExpressionTree) Implements_parser_ParseTree()         {}
func (this_ *AwaitExpressionTree) FieldGet_expression() ParseTree { return this_.expression }

func (this_ *AwaitExpressionTree) initFields_AwaitExpressionTree() {
	_ = this_
}

// Java constructor: AwaitExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression)
func NewAwaitExpressionTree(location SourceRange, expression ParseTree) *AwaitExpressionTree {
	this_ := &AwaitExpressionTree{}
	this_.BindSelf_(this_)
	InitAwaitExpressionTree(this_, location, expression)
	return this_
}
func allocPooled_AwaitExpressionTree(slab_ *j2gort.PooledSlab[AwaitExpressionTree], location SourceRange, expression ParseTree) *AwaitExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitAwaitExpressionTree(this_, location, expression)
	return this_
}

func InitAwaitExpressionTree(this_ *AwaitExpressionTree, location SourceRange, expression ParseTree) {
	_ = location
	_ = expression
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_AWAIT_EXPRESSION(), location)
	this_.initFields_AwaitExpressionTree()
	this_.expression = expression
}
