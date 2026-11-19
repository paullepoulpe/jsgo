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
// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.MemberExpressionTree
type MemberExpressionTree struct {
	ParseTreeBase
	operand    ParseTree `gson:"operand"`
	memberName Token     `gson:"memberName"`
}

func (this_ *MemberExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *MemberExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*MemberExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.MemberExpressionTree", "MemberExpressionTree"
}
func (*MemberExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.MemberExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.MemberExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*MemberExpressionTree](v_)
	})
}
func (*MemberExpressionTree) Implements_parser_ParseTree()      {}
func (this_ *MemberExpressionTree) FieldGet_operand() ParseTree { return this_.operand }
func (this_ *MemberExpressionTree) FieldGet_memberName() Token  { return this_.memberName }

func (this_ *MemberExpressionTree) initFields_MemberExpressionTree() {
	_ = this_
}

// Java constructor: MemberExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree operand, com.google.javascript.jscomp.parsing.parser.Token memberName)
func NewMemberExpressionTree(location SourceRange, operand ParseTree, memberName Token) *MemberExpressionTree {
	this_ := &MemberExpressionTree{}
	this_.BindSelf_(this_)
	InitMemberExpressionTree(this_, location, operand, memberName)
	return this_
}
func allocPooled_MemberExpressionTree(slab_ *j2gort.PooledSlab[MemberExpressionTree], location SourceRange, operand ParseTree, memberName Token) *MemberExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitMemberExpressionTree(this_, location, operand, memberName)
	return this_
}

func InitMemberExpressionTree(this_ *MemberExpressionTree, location SourceRange, operand ParseTree, memberName Token) {
	_ = location
	_ = operand
	_ = memberName
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_MEMBER_EXPRESSION(), location)
	this_.initFields_MemberExpressionTree()
	this_.operand = operand
	this_.memberName = memberName
}
