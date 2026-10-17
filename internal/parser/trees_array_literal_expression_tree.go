// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayLiteralExpressionTree.java @ closure-compiler v20260415

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

// ArrayLiteralExpressionTree is the parse tree for ParseTreeType ARRAY_LITERAL_EXPRESSION.
type ArrayLiteralExpressionTree struct {
	ParseTreeBase

	Elements         []ParseTree
	HasTrailingComma bool
}

// NewArrayLiteralExpressionTree constructs a ArrayLiteralExpressionTree.
//
// Java: ArrayLiteralExpressionTree(SourceRange, ImmutableList, boolean)
func NewArrayLiteralExpressionTree(location *SourceRange, elements []ParseTree, hasTrailingComma bool) *ArrayLiteralExpressionTree {
	t := &ArrayLiteralExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_ARRAY_LITERAL_EXPRESSION, location)
	t.Elements = elements
	t.HasTrailingComma = hasTrailingComma
	return t
}
