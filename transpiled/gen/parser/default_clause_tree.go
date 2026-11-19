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
// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultClauseTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.DefaultClauseTree
type DefaultClauseTree struct {
	ParseTreeBase
	statements j2gort.Ext_com_google_common_collect_ImmutableList `gson:"statements"`
}

func (this_ *DefaultClauseTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *DefaultClauseTree) IsJavaNil_() bool { return this_ == nil }
func (*DefaultClauseTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.DefaultClauseTree", "DefaultClauseTree"
}
func (*DefaultClauseTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.DefaultClauseTree", "com.google.javascript.jscomp.parsing.parser.trees.DefaultClauseTree", func(v_ any) bool {
		return j2gort.InstanceOf[*DefaultClauseTree](v_)
	})
}
func (*DefaultClauseTree) Implements_parser_ParseTree() {}
func (this_ *DefaultClauseTree) FieldGet_statements() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.statements
}

func (this_ *DefaultClauseTree) initFields_DefaultClauseTree() {
	_ = this_
}

// Java constructor: DefaultClauseTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> statements)
func NewDefaultClauseTree(location SourceRange, statements j2gort.Ext_com_google_common_collect_ImmutableList) *DefaultClauseTree {
	this_ := &DefaultClauseTree{}
	this_.BindSelf_(this_)
	InitDefaultClauseTree(this_, location, statements)
	return this_
}
func allocPooled_DefaultClauseTree(slab_ *j2gort.PooledSlab[DefaultClauseTree], location SourceRange, statements j2gort.Ext_com_google_common_collect_ImmutableList) *DefaultClauseTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitDefaultClauseTree(this_, location, statements)
	return this_
}

func InitDefaultClauseTree(this_ *DefaultClauseTree, location SourceRange, statements j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = statements
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_DEFAULT_CLAUSE(), location)
	this_.initFields_DefaultClauseTree()
	this_.statements = statements
}
