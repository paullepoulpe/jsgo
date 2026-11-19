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
// Port of com/google/javascript/jscomp/parsing/parser/trees/OptChainCallExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.OptChainCallExpressionTree
type OptChainCallExpressionTree struct {
	ParseTreeBase
	operand                ParseTree         `gson:"operand"`
	arguments              *ArgumentListTree `gson:"arguments"`
	isStartOfOptionalChain bool              `gson:"isStartOfOptionalChain"`
	hasTrailingComma       bool              `gson:"hasTrailingComma"`
}

func (this_ *OptChainCallExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *OptChainCallExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*OptChainCallExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.OptChainCallExpressionTree", "OptChainCallExpressionTree"
}
func (*OptChainCallExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.OptChainCallExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.OptChainCallExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*OptChainCallExpressionTree](v_)
	})
}
func (*OptChainCallExpressionTree) Implements_parser_ParseTree()      {}
func (this_ *OptChainCallExpressionTree) FieldGet_operand() ParseTree { return this_.operand }
func (this_ *OptChainCallExpressionTree) FieldGet_arguments() *ArgumentListTree {
	return this_.arguments
}
func (this_ *OptChainCallExpressionTree) FieldGet_isStartOfOptionalChain() bool {
	return this_.isStartOfOptionalChain
}
func (this_ *OptChainCallExpressionTree) FieldGet_hasTrailingComma() bool {
	return this_.hasTrailingComma
}

func (this_ *OptChainCallExpressionTree) initFields_OptChainCallExpressionTree() {
	_ = this_
}

// Java constructor: OptChainCallExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand, com.google.javascript.jscomp.parsing.parser.trees.ArgumentListTree arguments, boolean isStartOfOptionalChain, boolean hasTrailingComma)
func NewOptChainCallExpressionTree(location SourceRange, operand ParseTree, arguments *ArgumentListTree, isStartOfOptionalChain bool, hasTrailingComma bool) *OptChainCallExpressionTree {
	this_ := &OptChainCallExpressionTree{}
	this_.BindSelf_(this_)
	InitOptChainCallExpressionTree(this_, location, operand, arguments, isStartOfOptionalChain, hasTrailingComma)
	return this_
}
func allocPooled_OptChainCallExpressionTree(slab_ *j2gort.PooledSlab[OptChainCallExpressionTree], location SourceRange, operand ParseTree, arguments *ArgumentListTree, isStartOfOptionalChain bool, hasTrailingComma bool) *OptChainCallExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitOptChainCallExpressionTree(this_, location, operand, arguments, isStartOfOptionalChain, hasTrailingComma)
	return this_
}

func InitOptChainCallExpressionTree(this_ *OptChainCallExpressionTree, location SourceRange, operand ParseTree, arguments *ArgumentListTree, isStartOfOptionalChain bool, hasTrailingComma bool) {
	_ = location
	_ = operand
	_ = arguments
	_ = isStartOfOptionalChain
	_ = hasTrailingComma
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_OPT_CHAIN_CALL_EXPRESSION(), location)
	this_.initFields_OptChainCallExpressionTree()
	this_.operand = operand
	this_.arguments = arguments
	this_.isStartOfOptionalChain = isStartOfOptionalChain
	this_.hasTrailingComma = hasTrailingComma
}
