/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/trees/DynamicImportTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.DynamicImportTree
type DynamicImportTree struct {
	ParseTreeBase
	argument ParseTree `gson:"argument"`
}

func (this_ *DynamicImportTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *DynamicImportTree) IsJavaNil_() bool { return this_ == nil }
func (*DynamicImportTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.DynamicImportTree", "DynamicImportTree"
}
func (*DynamicImportTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.DynamicImportTree", "com.google.javascript.jscomp.parsing.parser.trees.DynamicImportTree", func(v_ any) bool {
		return j2gort.InstanceOf[*DynamicImportTree](v_)
	})
}
func (*DynamicImportTree) Implements_parser_ParseTree()       {}
func (this_ *DynamicImportTree) FieldGet_argument() ParseTree { return this_.argument }

func (this_ *DynamicImportTree) initFields_DynamicImportTree() {
	_ = this_
}

// Java constructor: DynamicImportTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree argument)
func NewDynamicImportTree(location SourceRange, argument ParseTree) *DynamicImportTree {
	this_ := &DynamicImportTree{}
	this_.BindSelf_(this_)
	InitDynamicImportTree(this_, location, argument)
	return this_
}
func allocPooled_DynamicImportTree(slab_ *j2gort.PooledSlab[DynamicImportTree], location SourceRange, argument ParseTree) *DynamicImportTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitDynamicImportTree(this_, location, argument)
	return this_
}

func InitDynamicImportTree(this_ *DynamicImportTree, location SourceRange, argument ParseTree) {
	_ = location
	_ = argument
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_DYNAMIC_IMPORT_EXPRESSION(), location)
	this_.initFields_DynamicImportTree()
	this_.argument = argument
}
