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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportDeclarationTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ExportDeclarationTree
type ExportDeclarationTree struct {
	ParseTreeBase
	declaration         ParseTree                                          `gson:"declaration"`
	exportSpecifierList j2gort.Ext_com_google_common_collect_ImmutableList `gson:"exportSpecifierList"`
	from                Token                                              `gson:"from"`
	isDefault           bool                                               `gson:"isDefault"`
	isExportAll         bool                                               `gson:"isExportAll"`
}

func (this_ *ExportDeclarationTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ExportDeclarationTree) IsJavaNil_() bool { return this_ == nil }
func (*ExportDeclarationTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ExportDeclarationTree", "ExportDeclarationTree"
}
func (*ExportDeclarationTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ExportDeclarationTree", "com.google.javascript.jscomp.parsing.parser.trees.ExportDeclarationTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ExportDeclarationTree](v_)
	})
}
func (*ExportDeclarationTree) Implements_parser_ParseTree()          {}
func (this_ *ExportDeclarationTree) FieldGet_isDefault() bool        { return this_.isDefault }
func (this_ *ExportDeclarationTree) FieldGet_isExportAll() bool      { return this_.isExportAll }
func (this_ *ExportDeclarationTree) FieldGet_declaration() ParseTree { return this_.declaration }
func (this_ *ExportDeclarationTree) FieldGet_exportSpecifierList() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.exportSpecifierList
}
func (this_ *ExportDeclarationTree) FieldGet_from() Token { return this_.from }

func (this_ *ExportDeclarationTree) initFields_ExportDeclarationTree() {
	_ = this_
}

// Java constructor: ExportDeclarationTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, boolean isDefault, boolean isExportAll, com.google.javascript.jscomp.parsing.parser.trees.ParseTree declaration, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> exportSpecifierList, com.google.javascript.jscomp.parsing.parser.Token from)
func NewExportDeclarationTree(location SourceRange, isDefault bool, isExportAll bool, declaration ParseTree, exportSpecifierList j2gort.Ext_com_google_common_collect_ImmutableList, from Token) *ExportDeclarationTree {
	this_ := &ExportDeclarationTree{}
	this_.BindSelf_(this_)
	InitExportDeclarationTree(this_, location, isDefault, isExportAll, declaration, exportSpecifierList, from)
	return this_
}
func allocPooled_ExportDeclarationTree(slab_ *j2gort.PooledSlab[ExportDeclarationTree], location SourceRange, isDefault bool, isExportAll bool, declaration ParseTree, exportSpecifierList j2gort.Ext_com_google_common_collect_ImmutableList, from Token) *ExportDeclarationTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitExportDeclarationTree(this_, location, isDefault, isExportAll, declaration, exportSpecifierList, from)
	return this_
}

func InitExportDeclarationTree(this_ *ExportDeclarationTree, location SourceRange, isDefault bool, isExportAll bool, declaration ParseTree, exportSpecifierList j2gort.Ext_com_google_common_collect_ImmutableList, from Token) {
	_ = location
	_ = isDefault
	_ = isExportAll
	_ = declaration
	_ = exportSpecifierList
	_ = from
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_EXPORT_DECLARATION(), location)
	this_.initFields_ExportDeclarationTree()
	this_.isDefault = isDefault
	this_.isExportAll = isExportAll
	this_.declaration = declaration
	this_.exportSpecifierList = exportSpecifierList
	this_.from = from
}
