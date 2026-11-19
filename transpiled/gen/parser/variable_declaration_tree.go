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
// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationTree
type VariableDeclarationTree struct {
	ParseTreeBase
	lvalue      ParseTree `gson:"lvalue"`
	initializer ParseTree `gson:"initializer"`
}

func (this_ *VariableDeclarationTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *VariableDeclarationTree) IsJavaNil_() bool { return this_ == nil }
func (*VariableDeclarationTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationTree", "VariableDeclarationTree"
}
func (*VariableDeclarationTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationTree", "com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationTree", func(v_ any) bool {
		return j2gort.InstanceOf[*VariableDeclarationTree](v_)
	})
}
func (*VariableDeclarationTree) Implements_parser_ParseTree()          {}
func (this_ *VariableDeclarationTree) FieldGet_lvalue() ParseTree      { return this_.lvalue }
func (this_ *VariableDeclarationTree) FieldGet_initializer() ParseTree { return this_.initializer }

func (this_ *VariableDeclarationTree) initFields_VariableDeclarationTree() {
	_ = this_
}

// Java constructor: VariableDeclarationTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree lvalue, com.google.javascript.jscomp.parsing.parser.trees.ParseTree initializer)
func NewVariableDeclarationTree(location SourceRange, lvalue ParseTree, initializer ParseTree) *VariableDeclarationTree {
	this_ := &VariableDeclarationTree{}
	this_.BindSelf_(this_)
	InitVariableDeclarationTree(this_, location, lvalue, initializer)
	return this_
}
func allocPooled_VariableDeclarationTree(slab_ *j2gort.PooledSlab[VariableDeclarationTree], location SourceRange, lvalue ParseTree, initializer ParseTree) *VariableDeclarationTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitVariableDeclarationTree(this_, location, lvalue, initializer)
	return this_
}

func InitVariableDeclarationTree(this_ *VariableDeclarationTree, location SourceRange, lvalue ParseTree, initializer ParseTree) {
	_ = location
	_ = lvalue
	_ = initializer
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_VARIABLE_DECLARATION(), location)
	this_.initFields_VariableDeclarationTree()
	this_.lvalue = lvalue
	this_.initializer = initializer
}
