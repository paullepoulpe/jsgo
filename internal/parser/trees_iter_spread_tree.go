// Port of com/google/javascript/jscomp/parsing/parser/trees/IterSpreadTree.java @ closure-compiler v20261021

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

// IterSpreadTree is the parse tree for ParseTreeType ITER_SPREAD.
type IterSpreadTree struct {
	ParseTreeBase

	Expression ParseTree
}

// NewIterSpreadTree constructs a IterSpreadTree.
//
// Java: IterSpreadTree(SourceRange, ParseTree)
func NewIterSpreadTree(location SourceRange, expression ParseTree) *IterSpreadTree {
	t := &IterSpreadTree{}
	t.InitParseTreeBase(ParseTreeType_ITER_SPREAD, location)
	t.Expression = expression
	return t
}
