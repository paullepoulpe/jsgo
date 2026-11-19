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
// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.OptionalMemberExpressionTree
type OptionalMemberExpressionTree struct {
	ParseTreeBase
	operand                ParseTree `gson:"operand"`
	memberName             Token     `gson:"memberName"`
	isStartOfOptionalChain bool      `gson:"isStartOfOptionalChain"`
}

func (this_ *OptionalMemberExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *OptionalMemberExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*OptionalMemberExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.OptionalMemberExpressionTree", "OptionalMemberExpressionTree"
}
func (*OptionalMemberExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.OptionalMemberExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.OptionalMemberExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*OptionalMemberExpressionTree](v_)
	})
}
func (*OptionalMemberExpressionTree) Implements_parser_ParseTree()      {}
func (this_ *OptionalMemberExpressionTree) FieldGet_operand() ParseTree { return this_.operand }
func (this_ *OptionalMemberExpressionTree) FieldGet_memberName() Token  { return this_.memberName }
func (this_ *OptionalMemberExpressionTree) FieldGet_isStartOfOptionalChain() bool {
	return this_.isStartOfOptionalChain
}

func (this_ *OptionalMemberExpressionTree) initFields_OptionalMemberExpressionTree() {
	_ = this_
}

// Java constructor: OptionalMemberExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand, com.google.javascript.jscomp.parsing.parser.Token memberName, boolean isStartOfOptionalChain)
func NewOptionalMemberExpressionTree(location SourceRange, operand ParseTree, memberName Token, isStartOfOptionalChain bool) *OptionalMemberExpressionTree {
	this_ := &OptionalMemberExpressionTree{}
	this_.BindSelf_(this_)
	InitOptionalMemberExpressionTree(this_, location, operand, memberName, isStartOfOptionalChain)
	return this_
}
func allocPooled_OptionalMemberExpressionTree(slab_ *j2gort.PooledSlab[OptionalMemberExpressionTree], location SourceRange, operand ParseTree, memberName Token, isStartOfOptionalChain bool) *OptionalMemberExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitOptionalMemberExpressionTree(this_, location, operand, memberName, isStartOfOptionalChain)
	return this_
}

func InitOptionalMemberExpressionTree(this_ *OptionalMemberExpressionTree, location SourceRange, operand ParseTree, memberName Token, isStartOfOptionalChain bool) {
	_ = location
	_ = operand
	_ = memberName
	_ = isStartOfOptionalChain
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_OPT_CHAIN_MEMBER_EXPRESSION(), location)
	this_.initFields_OptionalMemberExpressionTree()
	this_.operand = operand
	this_.memberName = memberName
	this_.isStartOfOptionalChain = isStartOfOptionalChain
}
