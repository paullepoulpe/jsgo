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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportSpecifierTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ExportSpecifierTree
type ExportSpecifierTree struct {
	ParseTreeBase
	importedName    Token `gson:"importedName"`
	destinationName Token `gson:"destinationName"`
}

func (this_ *ExportSpecifierTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ExportSpecifierTree) IsJavaNil_() bool { return this_ == nil }
func (*ExportSpecifierTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ExportSpecifierTree", "ExportSpecifierTree"
}
func (*ExportSpecifierTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ExportSpecifierTree", "com.google.javascript.jscomp.parsing.parser.trees.ExportSpecifierTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ExportSpecifierTree](v_)
	})
}
func (*ExportSpecifierTree) Implements_parser_ParseTree()          {}
func (this_ *ExportSpecifierTree) FieldGet_importedName() Token    { return this_.importedName }
func (this_ *ExportSpecifierTree) FieldGet_destinationName() Token { return this_.destinationName }

func (this_ *ExportSpecifierTree) initFields_ExportSpecifierTree() {
	_ = this_
}

// Java constructor: ExportSpecifierTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token importedName, com.google.javascript.jscomp.parsing.parser.Token destinationName)
func NewExportSpecifierTree(location SourceRange, importedName Token, destinationName Token) *ExportSpecifierTree {
	this_ := &ExportSpecifierTree{}
	this_.BindSelf_(this_)
	InitExportSpecifierTree(this_, location, importedName, destinationName)
	return this_
}
func allocPooled_ExportSpecifierTree(slab_ *j2gort.PooledSlab[ExportSpecifierTree], location SourceRange, importedName Token, destinationName Token) *ExportSpecifierTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitExportSpecifierTree(this_, location, importedName, destinationName)
	return this_
}

func InitExportSpecifierTree(this_ *ExportSpecifierTree, location SourceRange, importedName Token, destinationName Token) {
	_ = location
	_ = importedName
	_ = destinationName
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_EXPORT_SPECIFIER(), location)
	this_.initFields_ExportSpecifierTree()
	this_.importedName = importedName
	this_.destinationName = destinationName
}
