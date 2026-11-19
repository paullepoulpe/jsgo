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
// Port of com/google/javascript/jscomp/parsing/parser/trees/WithStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.WithStatementTree
type WithStatementTree struct {
	ParseTreeBase
	expression ParseTree `gson:"expression"`
	body       ParseTree `gson:"body"`
}

func (this_ *WithStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *WithStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*WithStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.WithStatementTree", "WithStatementTree"
}
func (*WithStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.WithStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.WithStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*WithStatementTree](v_)
	})
}
func (*WithStatementTree) Implements_parser_ParseTree()         {}
func (this_ *WithStatementTree) FieldGet_expression() ParseTree { return this_.expression }
func (this_ *WithStatementTree) FieldGet_body() ParseTree       { return this_.body }

func (this_ *WithStatementTree) initFields_WithStatementTree() {
	_ = this_
}

// Java constructor: WithStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression, com.google.javascript.jscomp.parsing.parser.trees.ParseTree body)
func NewWithStatementTree(location SourceRange, expression ParseTree, body ParseTree) *WithStatementTree {
	this_ := &WithStatementTree{}
	this_.BindSelf_(this_)
	InitWithStatementTree(this_, location, expression, body)
	return this_
}
func allocPooled_WithStatementTree(slab_ *j2gort.PooledSlab[WithStatementTree], location SourceRange, expression ParseTree, body ParseTree) *WithStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitWithStatementTree(this_, location, expression, body)
	return this_
}

func InitWithStatementTree(this_ *WithStatementTree, location SourceRange, expression ParseTree, body ParseTree) {
	_ = location
	_ = expression
	_ = body
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_WITH_STATEMENT(), location)
	this_.initFields_WithStatementTree()
	this_.expression = expression
	this_.body = body
}
