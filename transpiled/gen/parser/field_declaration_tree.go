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
// Port of com/google/javascript/jscomp/parsing/parser/trees/FieldDeclarationTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.FieldDeclarationTree
type FieldDeclarationTree struct {
	ParseTreeBase
	name        Token     `gson:"name"`
	initializer ParseTree `gson:"initializer"`
	isStatic    bool      `gson:"isStatic"`
}

func (this_ *FieldDeclarationTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *FieldDeclarationTree) IsJavaNil_() bool { return this_ == nil }
func (*FieldDeclarationTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.FieldDeclarationTree", "FieldDeclarationTree"
}
func (*FieldDeclarationTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.FieldDeclarationTree", "com.google.javascript.jscomp.parsing.parser.trees.FieldDeclarationTree", func(v_ any) bool {
		return j2gort.InstanceOf[*FieldDeclarationTree](v_)
	})
}
func (*FieldDeclarationTree) Implements_parser_ParseTree()          {}
func (this_ *FieldDeclarationTree) FieldGet_name() Token            { return this_.name }
func (this_ *FieldDeclarationTree) FieldGet_isStatic() bool         { return this_.isStatic }
func (this_ *FieldDeclarationTree) FieldGet_initializer() ParseTree { return this_.initializer }

func (this_ *FieldDeclarationTree) initFields_FieldDeclarationTree() {
	_ = this_
}

// Java constructor: FieldDeclarationTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token name, boolean isStatic, com.google.javascript.jscomp.parsing.parser.trees.ParseTree initializer)
func NewFieldDeclarationTree(location SourceRange, name Token, isStatic bool, initializer ParseTree) *FieldDeclarationTree {
	this_ := &FieldDeclarationTree{}
	this_.BindSelf_(this_)
	InitFieldDeclarationTree(this_, location, name, isStatic, initializer)
	return this_
}
func allocPooled_FieldDeclarationTree(slab_ *j2gort.PooledSlab[FieldDeclarationTree], location SourceRange, name Token, isStatic bool, initializer ParseTree) *FieldDeclarationTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitFieldDeclarationTree(this_, location, name, isStatic, initializer)
	return this_
}

func InitFieldDeclarationTree(this_ *FieldDeclarationTree, location SourceRange, name Token, isStatic bool, initializer ParseTree) {
	_ = location
	_ = name
	_ = isStatic
	_ = initializer
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FIELD_DECLARATION(), location)
	this_.initFields_FieldDeclarationTree()
	this_.name = name
	this_.isStatic = isStatic
	this_.initializer = initializer
}
