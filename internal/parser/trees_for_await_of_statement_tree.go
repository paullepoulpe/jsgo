// Port of com/google/javascript/jscomp/parsing/parser/trees/ForAwaitOfStatementTree.java @ closure-compiler v20261015

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

// Represents a for-await-of statement, which was added in ES2018
type ForAwaitOfStatementTree struct {
	ParseTreeBase

	Initializer ParseTree
	Collection  ParseTree
	Body        ParseTree
}

// NewForAwaitOfStatementTree constructs a ForAwaitOfStatementTree.
//
// Java: ForAwaitOfStatementTree(SourceRange, ParseTree, ParseTree, ParseTree)
func NewForAwaitOfStatementTree(location SourceRange, initializer ParseTree, collection ParseTree, body ParseTree) *ForAwaitOfStatementTree {
	t := &ForAwaitOfStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_FOR_AWAIT_OF_STATEMENT, location)
	t.Initializer = initializer
	t.Collection = collection
	t.Body = body
	return t
}
