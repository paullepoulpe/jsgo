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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ContinueStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ContinueStatementTree
type ContinueStatementTree struct {
	ParseTreeBase
	name Token `gson:"name"`
}

func (this_ *ContinueStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ContinueStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*ContinueStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ContinueStatementTree", "ContinueStatementTree"
}
func (*ContinueStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ContinueStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.ContinueStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ContinueStatementTree](v_)
	})
}
func (*ContinueStatementTree) Implements_parser_ParseTree() {}
func (this_ *ContinueStatementTree) FieldGet_name() Token   { return this_.name }

func (this_ *ContinueStatementTree) initFields_ContinueStatementTree() {
	_ = this_
}

// Java constructor: ContinueStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.@org.jspecify.annotations.Nullable Token name)
func NewContinueStatementTree(location SourceRange, name Token) *ContinueStatementTree {
	this_ := &ContinueStatementTree{}
	this_.BindSelf_(this_)
	InitContinueStatementTree(this_, location, name)
	return this_
}
func allocPooled_ContinueStatementTree(slab_ *j2gort.PooledSlab[ContinueStatementTree], location SourceRange, name Token) *ContinueStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitContinueStatementTree(this_, location, name)
	return this_
}

func InitContinueStatementTree(this_ *ContinueStatementTree, location SourceRange, name Token) {
	_ = location
	_ = name
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_CONTINUE_STATEMENT(), location)
	this_.initFields_ContinueStatementTree()
	this_.name = name
}

// Java method: ContinueStatementTree.hasLabel()
func (this_ *ContinueStatementTree) HasLabel() bool {
	return !(this_.name).IsJavaNil_()
}
