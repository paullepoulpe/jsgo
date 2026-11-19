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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ForAwaitOfStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ForAwaitOfStatementTree
type ForAwaitOfStatementTree struct {
	ParseTreeBase
	initializer ParseTree `gson:"initializer"`
	collection  ParseTree `gson:"collection"`
	body        ParseTree `gson:"body"`
}

func (this_ *ForAwaitOfStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ForAwaitOfStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*ForAwaitOfStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ForAwaitOfStatementTree", "ForAwaitOfStatementTree"
}
func (*ForAwaitOfStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ForAwaitOfStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.ForAwaitOfStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ForAwaitOfStatementTree](v_)
	})
}
func (*ForAwaitOfStatementTree) Implements_parser_ParseTree()          {}
func (this_ *ForAwaitOfStatementTree) FieldGet_initializer() ParseTree { return this_.initializer }
func (this_ *ForAwaitOfStatementTree) FieldGet_collection() ParseTree  { return this_.collection }
func (this_ *ForAwaitOfStatementTree) FieldGet_body() ParseTree        { return this_.body }

func (this_ *ForAwaitOfStatementTree) initFields_ForAwaitOfStatementTree() {
	_ = this_
}

// Java constructor: ForAwaitOfStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree initializer, com.google.javascript.jscomp.parsing.parser.trees.ParseTree collection, com.google.javascript.jscomp.parsing.parser.trees.ParseTree body)
func NewForAwaitOfStatementTree(location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForAwaitOfStatementTree {
	this_ := &ForAwaitOfStatementTree{}
	this_.BindSelf_(this_)
	InitForAwaitOfStatementTree(this_, location, initializer, collection, body)
	return this_
}
func allocPooled_ForAwaitOfStatementTree(slab_ *j2gort.PooledSlab[ForAwaitOfStatementTree], location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForAwaitOfStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitForAwaitOfStatementTree(this_, location, initializer, collection, body)
	return this_
}

func InitForAwaitOfStatementTree(this_ *ForAwaitOfStatementTree, location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) {
	_ = location
	_ = initializer
	_ = collection
	_ = body
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FOR_AWAIT_OF_STATEMENT(), location)
	this_.initFields_ForAwaitOfStatementTree()
	this_.initializer = initializer
	this_.collection = collection
	this_.body = body
}
