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
// Port of com/google/javascript/jscomp/parsing/parser/trees/WhileStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.WhileStatementTree
type WhileStatementTree struct {
	ParseTreeBase
	condition ParseTree `gson:"condition"`
	body      ParseTree `gson:"body"`
}

func (this_ *WhileStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *WhileStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*WhileStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.WhileStatementTree", "WhileStatementTree"
}
func (*WhileStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.WhileStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.WhileStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*WhileStatementTree](v_)
	})
}
func (*WhileStatementTree) Implements_parser_ParseTree()        {}
func (this_ *WhileStatementTree) FieldGet_condition() ParseTree { return this_.condition }
func (this_ *WhileStatementTree) FieldGet_body() ParseTree      { return this_.body }

func (this_ *WhileStatementTree) initFields_WhileStatementTree() {
	_ = this_
}

// Java constructor: WhileStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree condition, com.google.javascript.jscomp.parsing.parser.trees.ParseTree body)
func NewWhileStatementTree(location SourceRange, condition ParseTree, body ParseTree) *WhileStatementTree {
	this_ := &WhileStatementTree{}
	this_.BindSelf_(this_)
	InitWhileStatementTree(this_, location, condition, body)
	return this_
}
func allocPooled_WhileStatementTree(slab_ *j2gort.PooledSlab[WhileStatementTree], location SourceRange, condition ParseTree, body ParseTree) *WhileStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitWhileStatementTree(this_, location, condition, body)
	return this_
}

func InitWhileStatementTree(this_ *WhileStatementTree, location SourceRange, condition ParseTree, body ParseTree) {
	_ = location
	_ = condition
	_ = body
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_WHILE_STATEMENT(), location)
	this_.initFields_WhileStatementTree()
	this_.condition = condition
	this_.body = body
}
