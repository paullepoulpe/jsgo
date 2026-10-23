// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableDeclarationListTree.java @ closure-compiler v20261015

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

// VariableDeclarationListTree is the parse tree for ParseTreeType VARIABLE_DECLARATION_LIST.
type VariableDeclarationListTree struct {
	ParseTreeBase

	DeclarationType TokenType // const, var, let
	Declarations    []*VariableDeclarationTree
}

// NewVariableDeclarationListTree constructs a VariableDeclarationListTree.
//
// Java: VariableDeclarationListTree(SourceRange, TokenType, ImmutableList)
func NewVariableDeclarationListTree(location SourceRange, declarationType TokenType, declarations []*VariableDeclarationTree) *VariableDeclarationListTree {
	t := &VariableDeclarationListTree{}
	t.InitParseTreeBase(ParseTreeType_VARIABLE_DECLARATION_LIST, location)
	t.DeclarationType = declarationType
	t.Declarations = declarations
	return t
}
