// Port of com/google/javascript/jscomp/parsing/parser/trees/ParenExpressionTree.java @ closure-compiler v20261021

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

// ParenExpressionTree is the parse tree for ParseTreeType PAREN_EXPRESSION.
type ParenExpressionTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewParenExpressionTree constructs a ParenExpressionTree.
//
// Java: ParenExpressionTree(SourceRange, ParseTree)
func NewParenExpressionTree(location SourceRange, expression ParseTree) *ParenExpressionTree {
	t := &ParenExpressionTree{}
	t.InitParseTreeBase(ParseTreeType_PAREN_EXPRESSION, location)
	t.Expression = expression
	return t
}

// IsPattern looks through the parentheses.
//
// Java: ParseTree.isPattern()
func (t *ParenExpressionTree) IsPattern() bool {
	return t.Expression.IsPattern()
}

// IsValidNonVanillaAssignmentTarget looks through the parentheses.
//
// Java: ParseTree.isValidNonVanillaAssignmentTarget()
func (t *ParenExpressionTree) IsValidNonVanillaAssignmentTarget() bool {
	return t.Expression.IsValidNonVanillaAssignmentTarget()
}

// IsValidAssignmentTarget looks through the parentheses.
//
// Java: ParseTree.isValidAssignmentTarget()
func (t *ParenExpressionTree) IsValidAssignmentTarget() bool {
	return t.Expression.IsValidAssignmentTarget()
}
