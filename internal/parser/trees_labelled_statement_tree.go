// Port of com/google/javascript/jscomp/parsing/parser/trees/LabelledStatementTree.java @ closure-compiler v20261015

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

// LabelledStatementTree is the parse tree for ParseTreeType LABELLED_STATEMENT.
type LabelledStatementTree struct {
	ParseTreeBase

	Name      *IdentifierToken
	Statement ParseTree
}

// NewLabelledStatementTree constructs a LabelledStatementTree.
//
// Java: LabelledStatementTree(SourceRange, IdentifierToken, ParseTree)
func NewLabelledStatementTree(location SourceRange, name *IdentifierToken, statement ParseTree) *LabelledStatementTree {
	t := &LabelledStatementTree{}
	t.InitParseTreeBase(ParseTreeType_LABELLED_STATEMENT, location)
	t.Name = name
	t.Statement = statement
	return t
}
