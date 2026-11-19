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
// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberLookupExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.OptionalMemberLookupExpressionTree
type OptionalMemberLookupExpressionTree struct {
	ParseTreeBase
	operand                ParseTree `gson:"operand"`
	memberExpression       ParseTree `gson:"memberExpression"`
	isStartOfOptionalChain bool      `gson:"isStartOfOptionalChain"`
}

func (this_ *OptionalMemberLookupExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *OptionalMemberLookupExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*OptionalMemberLookupExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.OptionalMemberLookupExpressionTree", "OptionalMemberLookupExpressionTree"
}
func (*OptionalMemberLookupExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.OptionalMemberLookupExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.OptionalMemberLookupExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*OptionalMemberLookupExpressionTree](v_)
	})
}
func (*OptionalMemberLookupExpressionTree) Implements_parser_ParseTree()      {}
func (this_ *OptionalMemberLookupExpressionTree) FieldGet_operand() ParseTree { return this_.operand }
func (this_ *OptionalMemberLookupExpressionTree) FieldGet_memberExpression() ParseTree {
	return this_.memberExpression
}
func (this_ *OptionalMemberLookupExpressionTree) FieldGet_isStartOfOptionalChain() bool {
	return this_.isStartOfOptionalChain
}

func (this_ *OptionalMemberLookupExpressionTree) initFields_OptionalMemberLookupExpressionTree() {
	_ = this_
}

// Java constructor: OptionalMemberLookupExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand, com.google.javascript.jscomp.parsing.parser.trees.ParseTree memberExpression, boolean isStartOfOptionalChain)
func NewOptionalMemberLookupExpressionTree(location SourceRange, operand ParseTree, memberExpression ParseTree, isStartOfOptionalChain bool) *OptionalMemberLookupExpressionTree {
	this_ := &OptionalMemberLookupExpressionTree{}
	this_.BindSelf_(this_)
	InitOptionalMemberLookupExpressionTree(this_, location, operand, memberExpression, isStartOfOptionalChain)
	return this_
}
func allocPooled_OptionalMemberLookupExpressionTree(slab_ *j2gort.PooledSlab[OptionalMemberLookupExpressionTree], location SourceRange, operand ParseTree, memberExpression ParseTree, isStartOfOptionalChain bool) *OptionalMemberLookupExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitOptionalMemberLookupExpressionTree(this_, location, operand, memberExpression, isStartOfOptionalChain)
	return this_
}

func InitOptionalMemberLookupExpressionTree(this_ *OptionalMemberLookupExpressionTree, location SourceRange, operand ParseTree, memberExpression ParseTree, isStartOfOptionalChain bool) {
	_ = location
	_ = operand
	_ = memberExpression
	_ = isStartOfOptionalChain
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION(), location)
	this_.initFields_OptionalMemberLookupExpressionTree()
	this_.operand = operand
	this_.memberExpression = memberExpression
	this_.isStartOfOptionalChain = isStartOfOptionalChain
}
