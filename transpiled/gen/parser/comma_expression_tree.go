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
// Port of com/google/javascript/jscomp/parsing/parser/trees/CommaExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.CommaExpressionTree
type CommaExpressionTree struct {
	ParseTreeBase
	expressions j2gort.Ext_com_google_common_collect_ImmutableList `gson:"expressions"`
}

func (this_ *CommaExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *CommaExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*CommaExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.CommaExpressionTree", "CommaExpressionTree"
}
func (*CommaExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.CommaExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.CommaExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*CommaExpressionTree](v_)
	})
}
func (*CommaExpressionTree) Implements_parser_ParseTree() {}
func (this_ *CommaExpressionTree) FieldGet_expressions() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.expressions
}

func (this_ *CommaExpressionTree) initFields_CommaExpressionTree() {
	_ = this_
}

// Java constructor: CommaExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> expressions)
func NewCommaExpressionTree(location SourceRange, expressions j2gort.Ext_com_google_common_collect_ImmutableList) *CommaExpressionTree {
	this_ := &CommaExpressionTree{}
	this_.BindSelf_(this_)
	InitCommaExpressionTree(this_, location, expressions)
	return this_
}
func allocPooled_CommaExpressionTree(slab_ *j2gort.PooledSlab[CommaExpressionTree], location SourceRange, expressions j2gort.Ext_com_google_common_collect_ImmutableList) *CommaExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitCommaExpressionTree(this_, location, expressions)
	return this_
}

func InitCommaExpressionTree(this_ *CommaExpressionTree, location SourceRange, expressions j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = expressions
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_COMMA_EXPRESSION(), location)
	this_.initFields_CommaExpressionTree()
	this_.expressions = expressions
}
