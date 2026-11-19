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
// Port of com/google/javascript/jscomp/parsing/parser/trees/SwitchStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.SwitchStatementTree
type SwitchStatementTree struct {
	ParseTreeBase
	expression  ParseTree                                          `gson:"expression"`
	caseClauses j2gort.Ext_com_google_common_collect_ImmutableList `gson:"caseClauses"`
}

func (this_ *SwitchStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *SwitchStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*SwitchStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.SwitchStatementTree", "SwitchStatementTree"
}
func (*SwitchStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.SwitchStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.SwitchStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*SwitchStatementTree](v_)
	})
}
func (*SwitchStatementTree) Implements_parser_ParseTree()         {}
func (this_ *SwitchStatementTree) FieldGet_expression() ParseTree { return this_.expression }
func (this_ *SwitchStatementTree) FieldGet_caseClauses() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.caseClauses
}

func (this_ *SwitchStatementTree) initFields_SwitchStatementTree() {
	_ = this_
}

// Java constructor: SwitchStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> caseClauses)
func NewSwitchStatementTree(location SourceRange, expression ParseTree, caseClauses j2gort.Ext_com_google_common_collect_ImmutableList) *SwitchStatementTree {
	this_ := &SwitchStatementTree{}
	this_.BindSelf_(this_)
	InitSwitchStatementTree(this_, location, expression, caseClauses)
	return this_
}
func allocPooled_SwitchStatementTree(slab_ *j2gort.PooledSlab[SwitchStatementTree], location SourceRange, expression ParseTree, caseClauses j2gort.Ext_com_google_common_collect_ImmutableList) *SwitchStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitSwitchStatementTree(this_, location, expression, caseClauses)
	return this_
}

func InitSwitchStatementTree(this_ *SwitchStatementTree, location SourceRange, expression ParseTree, caseClauses j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = expression
	_ = caseClauses
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_SWITCH_STATEMENT(), location)
	this_.initFields_SwitchStatementTree()
	this_.expression = expression
	this_.caseClauses = caseClauses
}
