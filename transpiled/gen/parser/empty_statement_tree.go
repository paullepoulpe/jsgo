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
// Port of com/google/javascript/jscomp/parsing/parser/trees/EmptyStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.EmptyStatementTree
type EmptyStatementTree struct {
	ParseTreeBase
}

func (this_ *EmptyStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *EmptyStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*EmptyStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.EmptyStatementTree", "EmptyStatementTree"
}
func (*EmptyStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.EmptyStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.EmptyStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*EmptyStatementTree](v_)
	})
}
func (*EmptyStatementTree) Implements_parser_ParseTree() {}

func (this_ *EmptyStatementTree) initFields_EmptyStatementTree() {
	_ = this_
}

// Java constructor: EmptyStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewEmptyStatementTree(location SourceRange) *EmptyStatementTree {
	this_ := &EmptyStatementTree{}
	this_.BindSelf_(this_)
	InitEmptyStatementTree(this_, location)
	return this_
}
func allocPooled_EmptyStatementTree(slab_ *j2gort.PooledSlab[EmptyStatementTree], location SourceRange) *EmptyStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitEmptyStatementTree(this_, location)
	return this_
}

func InitEmptyStatementTree(this_ *EmptyStatementTree, location SourceRange) {
	_ = location
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_EMPTY_STATEMENT(), location)
	this_.initFields_EmptyStatementTree()
}
