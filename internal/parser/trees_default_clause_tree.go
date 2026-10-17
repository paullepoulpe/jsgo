// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultClauseTree.java @ closure-compiler v20260415

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

// DefaultClauseTree is the parse tree for ParseTreeType DEFAULT_CLAUSE.
type DefaultClauseTree struct {
	ParseTreeBase

	Statements []ParseTree
}

// NewDefaultClauseTree constructs a DefaultClauseTree.
//
// Java: DefaultClauseTree(SourceRange, ImmutableList)
func NewDefaultClauseTree(location *SourceRange, statements []ParseTree) *DefaultClauseTree {
	t := &DefaultClauseTree{}
	t.InitParseTreeBase(t, ParseTreeType_DEFAULT_CLAUSE, location)
	t.Statements = statements
	return t
}
