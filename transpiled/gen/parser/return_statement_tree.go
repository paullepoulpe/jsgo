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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ReturnStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ReturnStatementTree
type ReturnStatementTree struct {
	ParseTreeBase
	expression ParseTree `gson:"expression"`
}

func (this_ *ReturnStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ReturnStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*ReturnStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ReturnStatementTree", "ReturnStatementTree"
}
func (*ReturnStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ReturnStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.ReturnStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ReturnStatementTree](v_)
	})
}
func (*ReturnStatementTree) Implements_parser_ParseTree()         {}
func (this_ *ReturnStatementTree) FieldGet_expression() ParseTree { return this_.expression }

func (this_ *ReturnStatementTree) initFields_ReturnStatementTree() {
	_ = this_
}

// Java constructor: ReturnStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression)
func NewReturnStatementTree(location SourceRange, expression ParseTree) *ReturnStatementTree {
	this_ := &ReturnStatementTree{}
	this_.BindSelf_(this_)
	InitReturnStatementTree(this_, location, expression)
	return this_
}
func allocPooled_ReturnStatementTree(slab_ *j2gort.PooledSlab[ReturnStatementTree], location SourceRange, expression ParseTree) *ReturnStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitReturnStatementTree(this_, location, expression)
	return this_
}

func InitReturnStatementTree(this_ *ReturnStatementTree, location SourceRange, expression ParseTree) {
	_ = location
	_ = expression
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_RETURN_STATEMENT(), location)
	this_.initFields_ReturnStatementTree()
	this_.expression = expression
}
