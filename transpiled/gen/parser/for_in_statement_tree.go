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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ForInStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ForInStatementTree
type ForInStatementTree struct {
	ParseTreeBase
	initializer ParseTree `gson:"initializer"`
	collection  ParseTree `gson:"collection"`
	body        ParseTree `gson:"body"`
}

func (this_ *ForInStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ForInStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*ForInStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ForInStatementTree", "ForInStatementTree"
}
func (*ForInStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ForInStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.ForInStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ForInStatementTree](v_)
	})
}
func (*ForInStatementTree) Implements_parser_ParseTree()          {}
func (this_ *ForInStatementTree) FieldGet_initializer() ParseTree { return this_.initializer }
func (this_ *ForInStatementTree) FieldGet_collection() ParseTree  { return this_.collection }
func (this_ *ForInStatementTree) FieldGet_body() ParseTree        { return this_.body }

func (this_ *ForInStatementTree) initFields_ForInStatementTree() {
	_ = this_
}

// Java constructor: ForInStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree initializer, com.google.javascript.jscomp.parsing.parser.trees.ParseTree collection, com.google.javascript.jscomp.parsing.parser.trees.ParseTree body)
func NewForInStatementTree(location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForInStatementTree {
	this_ := &ForInStatementTree{}
	this_.BindSelf_(this_)
	InitForInStatementTree(this_, location, initializer, collection, body)
	return this_
}
func allocPooled_ForInStatementTree(slab_ *j2gort.PooledSlab[ForInStatementTree], location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForInStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitForInStatementTree(this_, location, initializer, collection, body)
	return this_
}

func InitForInStatementTree(this_ *ForInStatementTree, location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) {
	_ = location
	_ = initializer
	_ = collection
	_ = body
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FOR_IN_STATEMENT(), location)
	this_.initFields_ForInStatementTree()
	this_.initializer = initializer
	this_.collection = collection
	this_.body = body
}
