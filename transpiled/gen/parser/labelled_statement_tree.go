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
// Port of com/google/javascript/jscomp/parsing/parser/trees/LabelledStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.LabelledStatementTree
type LabelledStatementTree struct {
	ParseTreeBase
	name      Token     `gson:"name"`
	statement ParseTree `gson:"statement"`
}

func (this_ *LabelledStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *LabelledStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*LabelledStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.LabelledStatementTree", "LabelledStatementTree"
}
func (*LabelledStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.LabelledStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.LabelledStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*LabelledStatementTree](v_)
	})
}
func (*LabelledStatementTree) Implements_parser_ParseTree()        {}
func (this_ *LabelledStatementTree) FieldGet_name() Token          { return this_.name }
func (this_ *LabelledStatementTree) FieldGet_statement() ParseTree { return this_.statement }

func (this_ *LabelledStatementTree) initFields_LabelledStatementTree() {
	_ = this_
}

// Java constructor: LabelledStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token name, com.google.javascript.jscomp.parsing.parser.trees.ParseTree statement)
func NewLabelledStatementTree(location SourceRange, name Token, statement ParseTree) *LabelledStatementTree {
	this_ := &LabelledStatementTree{}
	this_.BindSelf_(this_)
	InitLabelledStatementTree(this_, location, name, statement)
	return this_
}
func allocPooled_LabelledStatementTree(slab_ *j2gort.PooledSlab[LabelledStatementTree], location SourceRange, name Token, statement ParseTree) *LabelledStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitLabelledStatementTree(this_, location, name, statement)
	return this_
}

func InitLabelledStatementTree(this_ *LabelledStatementTree, location SourceRange, name Token, statement ParseTree) {
	_ = location
	_ = name
	_ = statement
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_LABELLED_STATEMENT(), location)
	this_.initFields_LabelledStatementTree()
	this_.name = name
	this_.statement = statement
}
