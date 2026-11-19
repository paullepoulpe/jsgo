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
// Port of com/google/javascript/jscomp/parsing/parser/trees/NewExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.NewExpressionTree
type NewExpressionTree struct {
	ParseTreeBase
	operand          ParseTree         `gson:"operand"`
	arguments        *ArgumentListTree `gson:"arguments"`
	hasTrailingComma bool              `gson:"hasTrailingComma"`
}

func (this_ *NewExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *NewExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*NewExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.NewExpressionTree", "NewExpressionTree"
}
func (*NewExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.NewExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.NewExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*NewExpressionTree](v_)
	})
}
func (*NewExpressionTree) Implements_parser_ParseTree()                {}
func (this_ *NewExpressionTree) FieldGet_operand() ParseTree           { return this_.operand }
func (this_ *NewExpressionTree) FieldGet_arguments() *ArgumentListTree { return this_.arguments }
func (this_ *NewExpressionTree) FieldGet_hasTrailingComma() bool       { return this_.hasTrailingComma }

func (this_ *NewExpressionTree) initFields_NewExpressionTree() {
	_ = this_
}

// Java constructor: NewExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand, com.google.javascript.jscomp.parsing.parser.trees.ArgumentListTree arguments, boolean hasTrailingComma)
func NewNewExpressionTree(location SourceRange, operand ParseTree, arguments *ArgumentListTree, hasTrailingComma bool) *NewExpressionTree {
	this_ := &NewExpressionTree{}
	this_.BindSelf_(this_)
	InitNewExpressionTree(this_, location, operand, arguments, hasTrailingComma)
	return this_
}
func allocPooled_NewExpressionTree(slab_ *j2gort.PooledSlab[NewExpressionTree], location SourceRange, operand ParseTree, arguments *ArgumentListTree, hasTrailingComma bool) *NewExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitNewExpressionTree(this_, location, operand, arguments, hasTrailingComma)
	return this_
}

func InitNewExpressionTree(this_ *NewExpressionTree, location SourceRange, operand ParseTree, arguments *ArgumentListTree, hasTrailingComma bool) {
	_ = location
	_ = operand
	_ = arguments
	_ = hasTrailingComma
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_NEW_EXPRESSION(), location)
	this_.initFields_NewExpressionTree()
	this_.operand = operand
	this_.arguments = arguments
	this_.hasTrailingComma = hasTrailingComma
}
