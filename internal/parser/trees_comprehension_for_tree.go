// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionForTree.java @ closure-compiler v20261015

/*
 * Copyright 2014 The Closure Compiler Authors.
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

// ComprehensionForTree is the parse tree for ParseTreeType COMPREHENSION_FOR.
type ComprehensionForTree struct {
	ParseTreeBase

	Initializer ParseTree
	Collection  ParseTree
}

// NewComprehensionForTree constructs a ComprehensionForTree.
//
// Java: ComprehensionForTree(SourceRange, ParseTree, ParseTree)
func NewComprehensionForTree(location SourceRange, initializer ParseTree, collection ParseTree) *ComprehensionForTree {
	t := &ComprehensionForTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPREHENSION_FOR, location)
	t.Initializer = initializer
	t.Collection = collection
	return t
}
