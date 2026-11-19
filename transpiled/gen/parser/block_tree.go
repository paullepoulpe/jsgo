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
// Port of com/google/javascript/jscomp/parsing/parser/trees/BlockTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.BlockTree
type BlockTree struct {
	ParseTreeBase
	statements j2gort.Ext_com_google_common_collect_ImmutableList `gson:"statements"`
}

func (this_ *BlockTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *BlockTree) IsJavaNil_() bool { return this_ == nil }
func (*BlockTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.BlockTree", "BlockTree"
}
func (*BlockTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.BlockTree", "com.google.javascript.jscomp.parsing.parser.trees.BlockTree", func(v_ any) bool {
		return j2gort.InstanceOf[*BlockTree](v_)
	})
}
func (*BlockTree) Implements_parser_ParseTree() {}
func (this_ *BlockTree) FieldGet_statements() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.statements
}

func (this_ *BlockTree) initFields_BlockTree() {
	_ = this_
}

// Java constructor: BlockTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> statements)
func NewBlockTree(location SourceRange, statements j2gort.Ext_com_google_common_collect_ImmutableList) *BlockTree {
	this_ := &BlockTree{}
	this_.BindSelf_(this_)
	InitBlockTree(this_, location, statements)
	return this_
}
func allocPooled_BlockTree(slab_ *j2gort.PooledSlab[BlockTree], location SourceRange, statements j2gort.Ext_com_google_common_collect_ImmutableList) *BlockTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitBlockTree(this_, location, statements)
	return this_
}

func InitBlockTree(this_ *BlockTree, location SourceRange, statements j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = statements
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_BLOCK(), location)
	this_.initFields_BlockTree()
	this_.statements = statements
}
