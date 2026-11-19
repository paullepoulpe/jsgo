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
// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.VariableStatementTree
type VariableStatementTree struct {
	ParseTreeBase
	declarations *VariableDeclarationListTree `gson:"declarations"`
}

func (this_ *VariableStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *VariableStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*VariableStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.VariableStatementTree", "VariableStatementTree"
}
func (*VariableStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.VariableStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.VariableStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*VariableStatementTree](v_)
	})
}
func (*VariableStatementTree) Implements_parser_ParseTree() {}
func (this_ *VariableStatementTree) FieldGet_declarations() *VariableDeclarationListTree {
	return this_.declarations
}

func (this_ *VariableStatementTree) initFields_VariableStatementTree() {
	_ = this_
}

// Java constructor: VariableStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationListTree declarations)
func NewVariableStatementTree(location SourceRange, declarations *VariableDeclarationListTree) *VariableStatementTree {
	this_ := &VariableStatementTree{}
	this_.BindSelf_(this_)
	InitVariableStatementTree(this_, location, declarations)
	return this_
}
func allocPooled_VariableStatementTree(slab_ *j2gort.PooledSlab[VariableStatementTree], location SourceRange, declarations *VariableDeclarationListTree) *VariableStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitVariableStatementTree(this_, location, declarations)
	return this_
}

func InitVariableStatementTree(this_ *VariableStatementTree, location SourceRange, declarations *VariableDeclarationListTree) {
	_ = location
	_ = declarations
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_VARIABLE_STATEMENT(), location)
	this_.initFields_VariableStatementTree()
	this_.declarations = declarations
}
