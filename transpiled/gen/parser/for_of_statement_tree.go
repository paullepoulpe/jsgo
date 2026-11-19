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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ForOfStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ForOfStatementTree
type ForOfStatementTree struct {
	ParseTreeBase
	initializer ParseTree `gson:"initializer"`
	collection  ParseTree `gson:"collection"`
	body        ParseTree `gson:"body"`
}

func (this_ *ForOfStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ForOfStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*ForOfStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ForOfStatementTree", "ForOfStatementTree"
}
func (*ForOfStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ForOfStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.ForOfStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ForOfStatementTree](v_)
	})
}
func (*ForOfStatementTree) Implements_parser_ParseTree()          {}
func (this_ *ForOfStatementTree) FieldGet_initializer() ParseTree { return this_.initializer }
func (this_ *ForOfStatementTree) FieldGet_collection() ParseTree  { return this_.collection }
func (this_ *ForOfStatementTree) FieldGet_body() ParseTree        { return this_.body }

func (this_ *ForOfStatementTree) initFields_ForOfStatementTree() {
	_ = this_
}

// Java constructor: ForOfStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree initializer, com.google.javascript.jscomp.parsing.parser.trees.ParseTree collection, com.google.javascript.jscomp.parsing.parser.trees.ParseTree body)
func NewForOfStatementTree(location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForOfStatementTree {
	this_ := &ForOfStatementTree{}
	this_.BindSelf_(this_)
	InitForOfStatementTree(this_, location, initializer, collection, body)
	return this_
}
func allocPooled_ForOfStatementTree(slab_ *j2gort.PooledSlab[ForOfStatementTree], location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForOfStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitForOfStatementTree(this_, location, initializer, collection, body)
	return this_
}

func InitForOfStatementTree(this_ *ForOfStatementTree, location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) {
	_ = location
	_ = initializer
	_ = collection
	_ = body
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FOR_OF_STATEMENT(), location)
	this_.initFields_ForOfStatementTree()
	this_.initializer = initializer
	this_.collection = collection
	this_.body = body
}
