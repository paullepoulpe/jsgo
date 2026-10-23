// Port of com/google/javascript/jscomp/parsing/parser/trees/ArgumentListTree.java @ closure-compiler v20261015

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

// ArgumentListTree is the parse tree for ParseTreeType ARGUMENT_LIST.
type ArgumentListTree struct {
	ParseTreeBase

	Arguments        []ParseTree
	HasTrailingComma bool
	CommaPositions   []SourcePosition
}

// NewArgumentListTree constructs a ArgumentListTree.
//
// Java: ArgumentListTree(SourceRange, ImmutableList, boolean, ImmutableList)
func NewArgumentListTree(location SourceRange, arguments []ParseTree, hasTrailingComma bool, commaPositions []SourcePosition) *ArgumentListTree {
	t := &ArgumentListTree{}
	t.InitParseTreeBase(t, ParseTreeType_ARGUMENT_LIST, location)
	t.Arguments = arguments
	t.HasTrailingComma = hasTrailingComma
	t.CommaPositions = commaPositions
	return t
}
