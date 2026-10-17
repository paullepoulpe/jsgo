// Port of com/google/javascript/jscomp/parsing/parser/trees/DoWhileStatementTree.java @ closure-compiler v20261015

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

// DoWhileStatementTree is the parse tree for ParseTreeType DO_WHILE_STATEMENT.
type DoWhileStatementTree struct {
	ParseTreeBase

	Body      ParseTree
	Condition ParseTree
}

// NewDoWhileStatementTree constructs a DoWhileStatementTree.
//
// Java: DoWhileStatementTree(SourceRange, ParseTree, ParseTree)
func NewDoWhileStatementTree(location *SourceRange, body ParseTree, condition ParseTree) *DoWhileStatementTree {
	t := &DoWhileStatementTree{}
	t.InitParseTreeBase(t, ParseTreeType_DO_WHILE_STATEMENT, location)
	t.Body = body
	t.Condition = condition
	return t
}
