// Port of com/google/javascript/jscomp/parsing/parser/trees/BreakStatementTree.java @ closure-compiler v20261015

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

// BreakStatementTree is the parse tree for ParseTreeType BREAK_STATEMENT.
type BreakStatementTree struct {
	ParseTreeBase

	Name *IdentifierToken
}

// NewBreakStatementTree constructs a BreakStatementTree.
//
// Java: BreakStatementTree(SourceRange, IdentifierToken)
func NewBreakStatementTree(location *SourceRange, name *IdentifierToken) *BreakStatementTree {
	t := &BreakStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_BREAK_STATEMENT, location)
	t.Name = name
	return t
}

// HasLabel reports whether the statement names a label.
//
// Java: BreakStatementTree.hasLabel()
func (t *BreakStatementTree) HasLabel() bool {
	return t.Name != nil
}
