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
// Port of com/google/javascript/jscomp/parsing/parser/trees/TryStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.TryStatementTree
type TryStatementTree struct {
	ParseTreeBase
	body         ParseTree `gson:"body"`
	catchBlock   ParseTree `gson:"catchBlock"`
	finallyBlock ParseTree `gson:"finallyBlock"`
}

func (this_ *TryStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *TryStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*TryStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.TryStatementTree", "TryStatementTree"
}
func (*TryStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.TryStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.TryStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*TryStatementTree](v_)
	})
}
func (*TryStatementTree) Implements_parser_ParseTree()           {}
func (this_ *TryStatementTree) FieldGet_body() ParseTree         { return this_.body }
func (this_ *TryStatementTree) FieldGet_catchBlock() ParseTree   { return this_.catchBlock }
func (this_ *TryStatementTree) FieldGet_finallyBlock() ParseTree { return this_.finallyBlock }

func (this_ *TryStatementTree) initFields_TryStatementTree() {
	_ = this_
}

// Java constructor: TryStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree body, com.google.javascript.jscomp.parsing.parser.trees.ParseTree catchBlock, com.google.javascript.jscomp.parsing.parser.trees.ParseTree finallyBlock)
func NewTryStatementTree(location SourceRange, body ParseTree, catchBlock ParseTree, finallyBlock ParseTree) *TryStatementTree {
	this_ := &TryStatementTree{}
	this_.BindSelf_(this_)
	InitTryStatementTree(this_, location, body, catchBlock, finallyBlock)
	return this_
}
func allocPooled_TryStatementTree(slab_ *j2gort.PooledSlab[TryStatementTree], location SourceRange, body ParseTree, catchBlock ParseTree, finallyBlock ParseTree) *TryStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitTryStatementTree(this_, location, body, catchBlock, finallyBlock)
	return this_
}

func InitTryStatementTree(this_ *TryStatementTree, location SourceRange, body ParseTree, catchBlock ParseTree, finallyBlock ParseTree) {
	_ = location
	_ = body
	_ = catchBlock
	_ = finallyBlock
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_TRY_STATEMENT(), location)
	this_.initFields_TryStatementTree()
	this_.body = body
	this_.catchBlock = catchBlock
	this_.finallyBlock = finallyBlock
}
