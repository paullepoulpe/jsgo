// Port of com/google/javascript/jscomp/parsing/parser/trees/MemberLookupExpressionTree.java @ closure-compiler v20261015

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

package parser

// MemberLookupExpressionTree is the parse tree for ParseTreeType MEMBER_LOOKUP_EXPRESSION.
type MemberLookupExpressionTree struct {
	ParseTreeBase

	Operand          ParseTree
	MemberExpression ParseTree
}

// NewMemberLookupExpressionTree constructs a MemberLookupExpressionTree.
//
// Java: MemberLookupExpressionTree(SourceRange, ParseTree, ParseTree)
func NewMemberLookupExpressionTree(location SourceRange, operand ParseTree, memberExpression ParseTree) *MemberLookupExpressionTree {
	t := &MemberLookupExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_MEMBER_LOOKUP_EXPRESSION, location)
	t.Operand = operand
	t.MemberExpression = memberExpression
	return t
}
