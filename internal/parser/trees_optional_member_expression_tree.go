// Port of com/google/javascript/jscomp/parsing/parser/trees/OptionalMemberExpressionTree.java @ closure-compiler v20261015

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

// Creates a OptionalMemberExpressionTree that represents a getprop expression within an optional
// chain. Has an isOptionalChainStart field to indicate whether it is the start of an optional
// chain.
type OptionalMemberExpressionTree struct {
	ParseTreeBase

	Operand                ParseTree
	MemberName             *IdentifierToken
	IsStartOfOptionalChain bool
}

// NewOptionalMemberExpressionTree constructs a OptionalMemberExpressionTree.
//
// Java: OptionalMemberExpressionTree(SourceRange, ParseTree, IdentifierToken, boolean)
func NewOptionalMemberExpressionTree(location SourceRange, operand ParseTree, memberName *IdentifierToken, isStartOfOptionalChain bool) *OptionalMemberExpressionTree {
	t := &OptionalMemberExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_OPT_CHAIN_MEMBER_EXPRESSION, location)
	t.Operand = operand
	t.MemberName = memberName
	t.IsStartOfOptionalChain = isStartOfOptionalChain
	return t
}
