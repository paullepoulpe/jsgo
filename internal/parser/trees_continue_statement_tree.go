// Port of com/google/javascript/jscomp/parsing/parser/trees/ContinueStatementTree.java @ closure-compiler v20261021

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

// ContinueStatementTree is the parse tree for ParseTreeType CONTINUE_STATEMENT.
type ContinueStatementTree struct {
	ParseTreeBase

	Name *IdentifierToken
}

// NewContinueStatementTree constructs a ContinueStatementTree.
//
// Java: ContinueStatementTree(SourceRange, IdentifierToken)
func NewContinueStatementTree(location SourceRange, name *IdentifierToken) *ContinueStatementTree {
	t := &ContinueStatementTree{}
	t.InitParseTreeBase(ParseTreeType_CONTINUE_STATEMENT, location)
	t.Name = name
	return t
}

// HasLabel reports whether the statement names a label.
//
// Java: ContinueStatementTree.hasLabel()
func (t *ContinueStatementTree) HasLabel() bool {
	return t.Name != nil
}
