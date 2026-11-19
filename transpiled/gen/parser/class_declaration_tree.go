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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ClassDeclarationTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ClassDeclarationTree
type ClassDeclarationTree struct {
	ParseTreeBase
	name       Token                                              `gson:"name"`
	superClass ParseTree                                          `gson:"superClass"`
	elements   j2gort.Ext_com_google_common_collect_ImmutableList `gson:"elements"`
}

func (this_ *ClassDeclarationTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ClassDeclarationTree) IsJavaNil_() bool { return this_ == nil }
func (*ClassDeclarationTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ClassDeclarationTree", "ClassDeclarationTree"
}
func (*ClassDeclarationTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ClassDeclarationTree", "com.google.javascript.jscomp.parsing.parser.trees.ClassDeclarationTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ClassDeclarationTree](v_)
	})
}
func (*ClassDeclarationTree) Implements_parser_ParseTree()         {}
func (this_ *ClassDeclarationTree) FieldGet_name() Token           { return this_.name }
func (this_ *ClassDeclarationTree) FieldGet_superClass() ParseTree { return this_.superClass }
func (this_ *ClassDeclarationTree) FieldGet_elements() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.elements
}

func (this_ *ClassDeclarationTree) initFields_ClassDeclarationTree() {
	_ = this_
}

// Java constructor: ClassDeclarationTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token name, com.google.javascript.jscomp.parsing.parser.trees.ParseTree superClass, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> elements)
func NewClassDeclarationTree(location SourceRange, name Token, superClass ParseTree, elements j2gort.Ext_com_google_common_collect_ImmutableList) *ClassDeclarationTree {
	this_ := &ClassDeclarationTree{}
	this_.BindSelf_(this_)
	InitClassDeclarationTree(this_, location, name, superClass, elements)
	return this_
}
func allocPooled_ClassDeclarationTree(slab_ *j2gort.PooledSlab[ClassDeclarationTree], location SourceRange, name Token, superClass ParseTree, elements j2gort.Ext_com_google_common_collect_ImmutableList) *ClassDeclarationTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitClassDeclarationTree(this_, location, name, superClass, elements)
	return this_
}

func InitClassDeclarationTree(this_ *ClassDeclarationTree, location SourceRange, name Token, superClass ParseTree, elements j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = name
	_ = superClass
	_ = elements
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_CLASS_DECLARATION(), location)
	this_.initFields_ClassDeclarationTree()
	this_.name = name
	this_.superClass = superClass
	this_.elements = elements
}
