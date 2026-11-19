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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportDeclarationTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ImportDeclarationTree
type ImportDeclarationTree struct {
	ParseTreeBase
	defaultBindingIdentifier  Token                                              `gson:"defaultBindingIdentifier"`
	importSpecifierList       j2gort.Ext_com_google_common_collect_ImmutableList `gson:"importSpecifierList"`
	nameSpaceImportIdentifier Token                                              `gson:"nameSpaceImportIdentifier"`
	moduleSpecifier           Token                                              `gson:"moduleSpecifier"`
}

func (this_ *ImportDeclarationTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ImportDeclarationTree) IsJavaNil_() bool { return this_ == nil }
func (*ImportDeclarationTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ImportDeclarationTree", "ImportDeclarationTree"
}
func (*ImportDeclarationTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ImportDeclarationTree", "com.google.javascript.jscomp.parsing.parser.trees.ImportDeclarationTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ImportDeclarationTree](v_)
	})
}
func (*ImportDeclarationTree) Implements_parser_ParseTree() {}
func (this_ *ImportDeclarationTree) FieldGet_defaultBindingIdentifier() Token {
	return this_.defaultBindingIdentifier
}
func (this_ *ImportDeclarationTree) FieldGet_importSpecifierList() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.importSpecifierList
}
func (this_ *ImportDeclarationTree) FieldGet_nameSpaceImportIdentifier() Token {
	return this_.nameSpaceImportIdentifier
}
func (this_ *ImportDeclarationTree) FieldGet_moduleSpecifier() Token { return this_.moduleSpecifier }

func (this_ *ImportDeclarationTree) initFields_ImportDeclarationTree() {
	_ = this_
}

// Java constructor: ImportDeclarationTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token defaultBindingIdentifier, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> importSpecifierList, com.google.javascript.jscomp.parsing.parser.Token nameSpaceImportIdentifier, com.google.javascript.jscomp.parsing.parser.Token moduleSpecifier)
func NewImportDeclarationTree(location SourceRange, defaultBindingIdentifier Token, importSpecifierList j2gort.Ext_com_google_common_collect_ImmutableList, nameSpaceImportIdentifier Token, moduleSpecifier Token) *ImportDeclarationTree {
	this_ := &ImportDeclarationTree{}
	this_.BindSelf_(this_)
	InitImportDeclarationTree(this_, location, defaultBindingIdentifier, importSpecifierList, nameSpaceImportIdentifier, moduleSpecifier)
	return this_
}
func allocPooled_ImportDeclarationTree(slab_ *j2gort.PooledSlab[ImportDeclarationTree], location SourceRange, defaultBindingIdentifier Token, importSpecifierList j2gort.Ext_com_google_common_collect_ImmutableList, nameSpaceImportIdentifier Token, moduleSpecifier Token) *ImportDeclarationTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitImportDeclarationTree(this_, location, defaultBindingIdentifier, importSpecifierList, nameSpaceImportIdentifier, moduleSpecifier)
	return this_
}

func InitImportDeclarationTree(this_ *ImportDeclarationTree, location SourceRange, defaultBindingIdentifier Token, importSpecifierList j2gort.Ext_com_google_common_collect_ImmutableList, nameSpaceImportIdentifier Token, moduleSpecifier Token) {
	_ = location
	_ = defaultBindingIdentifier
	_ = importSpecifierList
	_ = nameSpaceImportIdentifier
	_ = moduleSpecifier
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_IMPORT_DECLARATION(), location)
	this_.initFields_ImportDeclarationTree()
	this_.defaultBindingIdentifier = defaultBindingIdentifier
	this_.importSpecifierList = importSpecifierList
	this_.nameSpaceImportIdentifier = nameSpaceImportIdentifier
	this_.moduleSpecifier = moduleSpecifier
}
