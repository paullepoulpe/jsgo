// Port of com/google/javascript/jscomp/parsing/parser/trees/WhileStatementTree.java @ closure-compiler v20261021

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

// WhileStatementTree is the parse tree for ParseTreeType WHILE_STATEMENT.
type WhileStatementTree struct {
	ParseTreeBase

	Condition ParseTree
	Body      ParseTree
}

// NewWhileStatementTree constructs a WhileStatementTree.
//
// Java: WhileStatementTree(SourceRange, ParseTree, ParseTree)
func NewWhileStatementTree(location SourceRange, condition ParseTree, body ParseTree) *WhileStatementTree {
	t := &WhileStatementTree{}
	t.InitParseTreeBase(ParseTreeType_WHILE_STATEMENT, location)
	t.Condition = condition
	t.Body = body
	return t
}
