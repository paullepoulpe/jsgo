// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayPatternTree.java @ closure-compiler v20261015

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

// ArrayPatternTree is the parse tree for ParseTreeType ARRAY_PATTERN.
type ArrayPatternTree struct {
	ParseTreeBase

	Elements []ParseTree
}

// NewArrayPatternTree constructs a ArrayPatternTree.
//
// Java: ArrayPatternTree(SourceRange, ImmutableList)
func NewArrayPatternTree(location *SourceRange, elements []ParseTree) *ArrayPatternTree {
	t := &ArrayPatternTree{}
	t.InitParseTreeBase(t, ParseTreeType_ARRAY_PATTERN, location)
	t.Elements = elements
	return t
}
