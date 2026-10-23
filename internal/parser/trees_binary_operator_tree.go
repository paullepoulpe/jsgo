// Port of com/google/javascript/jscomp/parsing/parser/trees/BinaryOperatorTree.java @ closure-compiler v20261015

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

// BinaryOperatorTree is the parse tree for ParseTreeType BINARY_OPERATOR.
type BinaryOperatorTree struct {
	ParseTreeBase

	Left     ParseTree
	Operator Token
	Right    ParseTree
}

// NewBinaryOperatorTree constructs a BinaryOperatorTree.
//
// Java: BinaryOperatorTree(SourceRange, ParseTree, Token, ParseTree)
func NewBinaryOperatorTree(location SourceRange, left ParseTree, operator Token, right ParseTree) *BinaryOperatorTree {
	t := &BinaryOperatorTree{}
	t.InitParseTreeBase(ParseTreeType_BINARY_OPERATOR, location)
	t.Left = left
	t.Operator = operator
	t.Right = right
	return t
}
