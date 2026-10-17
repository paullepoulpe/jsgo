// Port of com/google/javascript/jscomp/parsing/parser/trees/ForOfStatementTree.java @ closure-compiler v20260415

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

// ForOfStatementTree is the parse tree for ParseTreeType FOR_OF_STATEMENT.
type ForOfStatementTree struct {
	ParseTreeBase

	Initializer ParseTree
	Collection  ParseTree
	Body        ParseTree
}

// NewForOfStatementTree constructs a ForOfStatementTree.
//
// Java: ForOfStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewForOfStatementTree(location *SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForOfStatementTree {
	t := &ForOfStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_FOR_OF_STATEMENT, location)
	t.Initializer = initializer
	t.Collection = collection
	t.Body = body
	return t
}
