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
// Port of com/google/javascript/jscomp/parsing/parser/trees/BreakStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.BreakStatementTree
type BreakStatementTree struct {
	ParseTreeBase
	name Token `gson:"name"`
}

func (this_ *BreakStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *BreakStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*BreakStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.BreakStatementTree", "BreakStatementTree"
}
func (*BreakStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.BreakStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.BreakStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*BreakStatementTree](v_)
	})
}
func (*BreakStatementTree) Implements_parser_ParseTree() {}
func (this_ *BreakStatementTree) FieldGet_name() Token   { return this_.name }

func (this_ *BreakStatementTree) initFields_BreakStatementTree() {
	_ = this_
}

// Java constructor: BreakStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.@org.jspecify.annotations.Nullable Token name)
func NewBreakStatementTree(location SourceRange, name Token) *BreakStatementTree {
	this_ := &BreakStatementTree{}
	this_.BindSelf_(this_)
	InitBreakStatementTree(this_, location, name)
	return this_
}
func allocPooled_BreakStatementTree(slab_ *j2gort.PooledSlab[BreakStatementTree], location SourceRange, name Token) *BreakStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitBreakStatementTree(this_, location, name)
	return this_
}

func InitBreakStatementTree(this_ *BreakStatementTree, location SourceRange, name Token) {
	_ = location
	_ = name
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_BREAK_STATEMENT(), location)
	this_.initFields_BreakStatementTree()
	this_.name = name
}

// Java method: BreakStatementTree.hasLabel()
func (this_ *BreakStatementTree) HasLabel() bool {
	return (this_.name != (Token{}))
}
