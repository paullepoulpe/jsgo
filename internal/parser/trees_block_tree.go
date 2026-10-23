// Port of com/google/javascript/jscomp/parsing/parser/trees/BlockTree.java @ closure-compiler v20261015

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

// BlockTree is the parse tree for ParseTreeType BLOCK.
type BlockTree struct {
	ParseTreeBase

	Statements []ParseTree
}

// NewBlockTree constructs a BlockTree.
//
// Java: BlockTree(SourceRange, ImmutableList)
func NewBlockTree(location SourceRange, statements []ParseTree) *BlockTree {
	t := &BlockTree{}
	t.InitParseTreeBase(t, ParseTreeType_BLOCK, location)
	t.Statements = statements
	return t
}
