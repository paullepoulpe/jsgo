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
// Port of com/google/javascript/jscomp/parsing/parser/trees/CaseClauseTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.CaseClauseTree
type CaseClauseTree struct {
	ParseTreeBase
	expression ParseTree                                          `gson:"expression"`
	statements j2gort.Ext_com_google_common_collect_ImmutableList `gson:"statements"`
}

func (this_ *CaseClauseTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *CaseClauseTree) IsJavaNil_() bool { return this_ == nil }
func (*CaseClauseTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.CaseClauseTree", "CaseClauseTree"
}
func (*CaseClauseTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.CaseClauseTree", "com.google.javascript.jscomp.parsing.parser.trees.CaseClauseTree", func(v_ any) bool {
		return j2gort.InstanceOf[*CaseClauseTree](v_)
	})
}
func (*CaseClauseTree) Implements_parser_ParseTree()         {}
func (this_ *CaseClauseTree) FieldGet_expression() ParseTree { return this_.expression }
func (this_ *CaseClauseTree) FieldGet_statements() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.statements
}

func (this_ *CaseClauseTree) initFields_CaseClauseTree() {
	_ = this_
}

// Java constructor: CaseClauseTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> statements)
func NewCaseClauseTree(location SourceRange, expression ParseTree, statements j2gort.Ext_com_google_common_collect_ImmutableList) *CaseClauseTree {
	this_ := &CaseClauseTree{}
	this_.BindSelf_(this_)
	InitCaseClauseTree(this_, location, expression, statements)
	return this_
}
func allocPooled_CaseClauseTree(slab_ *j2gort.PooledSlab[CaseClauseTree], location SourceRange, expression ParseTree, statements j2gort.Ext_com_google_common_collect_ImmutableList) *CaseClauseTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitCaseClauseTree(this_, location, expression, statements)
	return this_
}

func InitCaseClauseTree(this_ *CaseClauseTree, location SourceRange, expression ParseTree, statements j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = expression
	_ = statements
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_CASE_CLAUSE(), location)
	this_.initFields_CaseClauseTree()
	this_.expression = expression
	this_.statements = statements
}
