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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportMetaExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ImportMetaExpressionTree
type ImportMetaExpressionTree struct {
	ParseTreeBase
}

func (this_ *ImportMetaExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ImportMetaExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*ImportMetaExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ImportMetaExpressionTree", "ImportMetaExpressionTree"
}
func (*ImportMetaExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ImportMetaExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.ImportMetaExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ImportMetaExpressionTree](v_)
	})
}
func (*ImportMetaExpressionTree) Implements_parser_ParseTree() {}

func (this_ *ImportMetaExpressionTree) initFields_ImportMetaExpressionTree() {
	_ = this_
}

// Java constructor: ImportMetaExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewImportMetaExpressionTree(location SourceRange) *ImportMetaExpressionTree {
	this_ := &ImportMetaExpressionTree{}
	this_.BindSelf_(this_)
	InitImportMetaExpressionTree(this_, location)
	return this_
}
func allocPooled_ImportMetaExpressionTree(slab_ *j2gort.PooledSlab[ImportMetaExpressionTree], location SourceRange) *ImportMetaExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitImportMetaExpressionTree(this_, location)
	return this_
}

func InitImportMetaExpressionTree(this_ *ImportMetaExpressionTree, location SourceRange) {
	_ = location
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_IMPORT_META_EXPRESSION(), location)
	this_.initFields_ImportMetaExpressionTree()
}
