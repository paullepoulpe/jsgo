// Port of com/google/javascript/jscomp/parsing/parser/trees/ConditionalExpressionTree.java @ closure-compiler v20261015

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

// ConditionalExpressionTree is the parse tree for ParseTreeType CONDITIONAL_EXPRESSION.
type ConditionalExpressionTree struct {
	ParseTreeBase

	Condition ParseTree
	Left      ParseTree
	Right     ParseTree
}

// NewConditionalExpressionTree constructs a ConditionalExpressionTree.
//
// Java: ConditionalExpressionTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewConditionalExpressionTree(location *SourceRange, condition ParseTree, left ParseTree, right ParseTree) *ConditionalExpressionTree {
	t := &ConditionalExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_CONDITIONAL_EXPRESSION, location)
	t.Condition = condition
	t.Left = left
	t.Right = right
	return t
}
