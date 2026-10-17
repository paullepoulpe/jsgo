// Port of com/google/javascript/jscomp/parsing/parser/trees/OptChainCallExpressionTree.java @ closure-compiler v20261015

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

// Creates a OptChainCallExpressionTree that represents a call expression within an optional chain.
// Has an isOptionalChainStart field to indicate whether it is the start of an optional chain.
type OptChainCallExpressionTree struct {
	ParseTreeBase

	Operand                ParseTree
	Arguments              *ArgumentListTree
	IsStartOfOptionalChain bool
	HasTrailingComma       bool
}

// NewOptChainCallExpressionTree constructs a OptChainCallExpressionTree.
//
// Java: OptChainCallExpressionTree(SourceRange, ParseTree, ArgumentListTree, boolean, boolean)
func NewOptChainCallExpressionTree(location *SourceRange, operand ParseTree, arguments *ArgumentListTree, isStartOfOptionalChain bool, hasTrailingComma bool) *OptChainCallExpressionTree {
	t := &OptChainCallExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_OPT_CHAIN_CALL_EXPRESSION, location)
	t.Operand = operand
	t.Arguments = arguments
	t.IsStartOfOptionalChain = isStartOfOptionalChain
	t.HasTrailingComma = hasTrailingComma
	return t
}
