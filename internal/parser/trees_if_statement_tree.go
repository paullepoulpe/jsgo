// Port of com/google/javascript/jscomp/parsing/parser/trees/IfStatementTree.java @ closure-compiler v20261015

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

// IfStatementTree is the parse tree for ParseTreeType IF_STATEMENT.
type IfStatementTree struct {
	ParseTreeBase

	Condition  ParseTree
	IfClause   ParseTree
	ElseClause ParseTree
}

// NewIfStatementTree constructs a IfStatementTree.
//
// Java: IfStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewIfStatementTree(location SourceRange, condition ParseTree, ifClause ParseTree, elseClause ParseTree) *IfStatementTree {
	t := &IfStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_IF_STATEMENT, location)
	t.Condition = condition
	t.IfClause = ifClause
	t.ElseClause = elseClause
	return t
}
