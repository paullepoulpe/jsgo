// Port of com/google/javascript/jscomp/parsing/parser/trees/SwitchStatementTree.java @ closure-compiler v20261015

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

// SwitchStatementTree is the parse tree for ParseTreeType SWITCH_STATEMENT.
type SwitchStatementTree struct {
	ParseTreeBase

	Expression  ParseTree
	CaseClauses []ParseTree
}

// NewSwitchStatementTree constructs a SwitchStatementTree.
//
// Java: SwitchStatementTree(SourceRange, ParseTree, ImmutableList)
func NewSwitchStatementTree(location SourceRange, expression ParseTree, caseClauses []ParseTree) *SwitchStatementTree {
	t := &SwitchStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_SWITCH_STATEMENT, location)
	t.Expression = expression
	t.CaseClauses = caseClauses
	return t
}
