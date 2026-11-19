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
// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationListTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationListTree
type VariableDeclarationListTree struct {
	ParseTreeBase
	declarations    []*VariableDeclarationTree `gson:"declarations"`
	declarationType TokenType                  `gson:"declarationType"`
}

func (this_ *VariableDeclarationListTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *VariableDeclarationListTree) IsJavaNil_() bool { return this_ == nil }
func (*VariableDeclarationListTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationListTree", "VariableDeclarationListTree"
}
func (*VariableDeclarationListTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationListTree", "com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationListTree", func(v_ any) bool {
		return j2gort.InstanceOf[*VariableDeclarationListTree](v_)
	})
}
func (*VariableDeclarationListTree) Implements_parser_ParseTree() {}
func (this_ *VariableDeclarationListTree) FieldGet_declarationType() TokenType {
	return this_.declarationType
}
func (this_ *VariableDeclarationListTree) FieldGet_declarations() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(this_.declarations)
}

func (this_ *VariableDeclarationListTree) initFields_VariableDeclarationListTree() {
	_ = this_
}

// Java constructor: VariableDeclarationListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.TokenType declarationType, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationTree> declarations)
func NewVariableDeclarationListTree(location SourceRange, declarationType TokenType, declarations j2gort.Ext_com_google_common_collect_ImmutableList) *VariableDeclarationListTree {
	this_ := &VariableDeclarationListTree{}
	this_.BindSelf_(this_)
	InitVariableDeclarationListTree(this_, location, declarationType, declarations)
	return this_
}
func allocPooled_VariableDeclarationListTree(slab_ *j2gort.PooledSlab[VariableDeclarationListTree], location SourceRange, declarationType TokenType, declarations j2gort.Ext_com_google_common_collect_ImmutableList) *VariableDeclarationListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitVariableDeclarationListTree(this_, location, declarationType, declarations)
	return this_
}

func InitVariableDeclarationListTree(this_ *VariableDeclarationListTree, location SourceRange, declarationType TokenType, declarations j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = declarationType
	_ = declarations
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_VARIABLE_DECLARATION_LIST(), location)
	this_.initFields_VariableDeclarationListTree()
	this_.declarationType = declarationType
	this_.declarations = j2gort.ImmListToSliceNullable[*VariableDeclarationTree](declarations)
}
func (this_ *VariableDeclarationListTree) Slice_FieldGet_declarations() []*VariableDeclarationTree {
	return this_.declarations
}
