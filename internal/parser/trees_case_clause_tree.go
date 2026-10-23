// Port of com/google/javascript/jscomp/parsing/parser/trees/CaseClauseTree.java @ closure-compiler v20261015

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

// CaseClauseTree is the parse tree for ParseTreeType CASE_CLAUSE.
type CaseClauseTree struct {
	ParseTreeBase

	Expression ParseTree
	Statements []ParseTree
}

// NewCaseClauseTree constructs a CaseClauseTree.
//
// Java: CaseClauseTree(SourceRange, ParseTree, ImmutableList)
func NewCaseClauseTree(location SourceRange, expression ParseTree, statements []ParseTree) *CaseClauseTree {
	t := &CaseClauseTree{}
	t.InitParseTreeBase(t, ParseTreeType_CASE_CLAUSE, location)
	t.Expression = expression
	t.Statements = statements
	return t
}
