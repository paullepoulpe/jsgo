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
// Port of com/google/javascript/jscomp/parsing/parser/trees/LiteralExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree
type LiteralExpressionTree struct {
	ParseTreeBase
	literalToken Token `gson:"literalToken"`
}

func (this_ *LiteralExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *LiteralExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*LiteralExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree", "LiteralExpressionTree"
}
func (*LiteralExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*LiteralExpressionTree](v_)
	})
}
func (*LiteralExpressionTree) Implements_parser_ParseTree()       {}
func (this_ *LiteralExpressionTree) FieldGet_literalToken() Token { return this_.literalToken }

func (this_ *LiteralExpressionTree) initFields_LiteralExpressionTree() {
	_ = this_
}

// Java constructor: LiteralExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token literalToken)
func NewLiteralExpressionTree(location SourceRange, literalToken Token) *LiteralExpressionTree {
	this_ := &LiteralExpressionTree{}
	this_.BindSelf_(this_)
	InitLiteralExpressionTree(this_, location, literalToken)
	return this_
}
func allocPooled_LiteralExpressionTree(slab_ *j2gort.PooledSlab[LiteralExpressionTree], location SourceRange, literalToken Token) *LiteralExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitLiteralExpressionTree(this_, location, literalToken)
	return this_
}

func InitLiteralExpressionTree(this_ *LiteralExpressionTree, location SourceRange, literalToken Token) {
	_ = location
	_ = literalToken
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_LITERAL_EXPRESSION(), location)
	this_.initFields_LiteralExpressionTree()
	this_.literalToken = literalToken
}
