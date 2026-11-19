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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ForStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ForStatementTree
type ForStatementTree struct {
	ParseTreeBase
	initializer ParseTree `gson:"initializer"`
	condition   ParseTree `gson:"condition"`
	increment   ParseTree `gson:"increment"`
	body        ParseTree `gson:"body"`
}

func (this_ *ForStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ForStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*ForStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ForStatementTree", "ForStatementTree"
}
func (*ForStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ForStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.ForStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ForStatementTree](v_)
	})
}
func (*ForStatementTree) Implements_parser_ParseTree()          {}
func (this_ *ForStatementTree) FieldGet_initializer() ParseTree { return this_.initializer }
func (this_ *ForStatementTree) FieldGet_condition() ParseTree   { return this_.condition }
func (this_ *ForStatementTree) FieldGet_increment() ParseTree   { return this_.increment }
func (this_ *ForStatementTree) FieldGet_body() ParseTree        { return this_.body }

func (this_ *ForStatementTree) initFields_ForStatementTree() {
	_ = this_
}

// Java constructor: ForStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree variables, com.google.javascript.jscomp.parsing.parser.trees.ParseTree condition, com.google.javascript.jscomp.parsing.parser.trees.ParseTree increment, com.google.javascript.jscomp.parsing.parser.trees.ParseTree body)
func NewForStatementTree(location SourceRange, variables ParseTree, condition ParseTree, increment ParseTree, body ParseTree) *ForStatementTree {
	this_ := &ForStatementTree{}
	this_.BindSelf_(this_)
	InitForStatementTree(this_, location, variables, condition, increment, body)
	return this_
}
func allocPooled_ForStatementTree(slab_ *j2gort.PooledSlab[ForStatementTree], location SourceRange, variables ParseTree, condition ParseTree, increment ParseTree, body ParseTree) *ForStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitForStatementTree(this_, location, variables, condition, increment, body)
	return this_
}

func InitForStatementTree(this_ *ForStatementTree, location SourceRange, variables ParseTree, condition ParseTree, increment ParseTree, body ParseTree) {
	_ = location
	_ = variables
	_ = condition
	_ = increment
	_ = body
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FOR_STATEMENT(), location)
	this_.initFields_ForStatementTree()
	this_.initializer = variables
	this_.condition = condition
	this_.increment = increment
	this_.body = body
}
