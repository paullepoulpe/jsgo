// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyMethodTree.java @ closure-compiler v20261015

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

// ComputedPropertyMethodTree is the parse tree for ParseTreeType COMPUTED_PROPERTY_METHOD.
type ComputedPropertyMethodTree struct {
	ParseTreeBase

	Property ParseTree
	Method   ParseTree
}

// NewComputedPropertyMethodTree constructs a ComputedPropertyMethodTree.
//
// Java: ComputedPropertyMethodTree(SourceRange, ParseTree, ParseTree)
func NewComputedPropertyMethodTree(location SourceRange, property ParseTree, method ParseTree) *ComputedPropertyMethodTree {
	t := &ComputedPropertyMethodTree{}
	t.InitParseTreeBase(ParseTreeType_COMPUTED_PROPERTY_METHOD, location)
	t.Property = property
	t.Method = method
	return t
}
