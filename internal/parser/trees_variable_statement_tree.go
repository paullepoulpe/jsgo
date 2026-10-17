// Port of com/google/javascript/jscomp/parsing/parser/trees/VariableStatementTree.java @ closure-compiler v20261015

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

// VariableStatementTree is the parse tree for ParseTreeType VARIABLE_STATEMENT.
type VariableStatementTree struct {
	ParseTreeBase

	Declarations *VariableDeclarationListTree
}

// NewVariableStatementTree constructs a VariableStatementTree.
//
// Java: VariableStatementTree(SourceRange, VariableDeclarationListTree)
func NewVariableStatementTree(location *SourceRange, declarations *VariableDeclarationListTree) *VariableStatementTree {
	t := &VariableStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_VARIABLE_STATEMENT, location)
	t.Declarations = declarations
	return t
}
