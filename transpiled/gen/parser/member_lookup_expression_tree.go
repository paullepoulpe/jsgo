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
// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberLookupExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.MemberLookupExpressionTree
type MemberLookupExpressionTree struct {
	ParseTreeBase
	operand          ParseTree `gson:"operand"`
	memberExpression ParseTree `gson:"memberExpression"`
}

func (this_ *MemberLookupExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *MemberLookupExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*MemberLookupExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.MemberLookupExpressionTree", "MemberLookupExpressionTree"
}
func (*MemberLookupExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.MemberLookupExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.MemberLookupExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*MemberLookupExpressionTree](v_)
	})
}
func (*MemberLookupExpressionTree) Implements_parser_ParseTree()      {}
func (this_ *MemberLookupExpressionTree) FieldGet_operand() ParseTree { return this_.operand }
func (this_ *MemberLookupExpressionTree) FieldGet_memberExpression() ParseTree {
	return this_.memberExpression
}

func (this_ *MemberLookupExpressionTree) initFields_MemberLookupExpressionTree() {
	_ = this_
}

// Java constructor: MemberLookupExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand, com.google.javascript.jscomp.parsing.parser.trees.ParseTree memberExpression)
func NewMemberLookupExpressionTree(location SourceRange, operand ParseTree, memberExpression ParseTree) *MemberLookupExpressionTree {
	this_ := &MemberLookupExpressionTree{}
	this_.BindSelf_(this_)
	InitMemberLookupExpressionTree(this_, location, operand, memberExpression)
	return this_
}
func allocPooled_MemberLookupExpressionTree(slab_ *j2gort.PooledSlab[MemberLookupExpressionTree], location SourceRange, operand ParseTree, memberExpression ParseTree) *MemberLookupExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitMemberLookupExpressionTree(this_, location, operand, memberExpression)
	return this_
}

func InitMemberLookupExpressionTree(this_ *MemberLookupExpressionTree, location SourceRange, operand ParseTree, memberExpression ParseTree) {
	_ = location
	_ = operand
	_ = memberExpression
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_MEMBER_LOOKUP_EXPRESSION(), location)
	this_.initFields_MemberLookupExpressionTree()
	this_.operand = operand
	this_.memberExpression = memberExpression
}
