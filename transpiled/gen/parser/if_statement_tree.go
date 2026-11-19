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
// Port of com/google/javascript/jscomp/parsing/parser/trees/IfStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.IfStatementTree
type IfStatementTree struct {
	ParseTreeBase
	condition  ParseTree `gson:"condition"`
	ifClause   ParseTree `gson:"ifClause"`
	elseClause ParseTree `gson:"elseClause"`
}

func (this_ *IfStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *IfStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*IfStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.IfStatementTree", "IfStatementTree"
}
func (*IfStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.IfStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.IfStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*IfStatementTree](v_)
	})
}
func (*IfStatementTree) Implements_parser_ParseTree()         {}
func (this_ *IfStatementTree) FieldGet_condition() ParseTree  { return this_.condition }
func (this_ *IfStatementTree) FieldGet_ifClause() ParseTree   { return this_.ifClause }
func (this_ *IfStatementTree) FieldGet_elseClause() ParseTree { return this_.elseClause }

func (this_ *IfStatementTree) initFields_IfStatementTree() {
	_ = this_
}

// Java constructor: IfStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree condition, com.google.javascript.jscomp.parsing.parser.trees.ParseTree ifClause, com.google.javascript.jscomp.parsing.parser.trees.ParseTree elseClause)
func NewIfStatementTree(location SourceRange, condition ParseTree, ifClause ParseTree, elseClause ParseTree) *IfStatementTree {
	this_ := &IfStatementTree{}
	this_.BindSelf_(this_)
	InitIfStatementTree(this_, location, condition, ifClause, elseClause)
	return this_
}
func allocPooled_IfStatementTree(slab_ *j2gort.PooledSlab[IfStatementTree], location SourceRange, condition ParseTree, ifClause ParseTree, elseClause ParseTree) *IfStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitIfStatementTree(this_, location, condition, ifClause, elseClause)
	return this_
}

func InitIfStatementTree(this_ *IfStatementTree, location SourceRange, condition ParseTree, ifClause ParseTree, elseClause ParseTree) {
	_ = location
	_ = condition
	_ = ifClause
	_ = elseClause
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_IF_STATEMENT(), location)
	this_.initFields_IfStatementTree()
	this_.condition = condition
	this_.ifClause = ifClause
	this_.elseClause = elseClause
}
