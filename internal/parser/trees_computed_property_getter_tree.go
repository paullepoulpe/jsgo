// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyGetterTree.java @ closure-compiler v20261015

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

// ComputedPropertyGetterTree is the parse tree for ParseTreeType COMPUTED_PROPERTY_GETTER.
type ComputedPropertyGetterTree struct {
	ParseTreeBase

	Property ParseTree
	IsStatic bool
	Body     *BlockTree
}

// NewComputedPropertyGetterTree constructs a ComputedPropertyGetterTree.
//
// Java: ComputedPropertyGetterTree(SourceRange, ParseTree, boolean, BlockTree)
func NewComputedPropertyGetterTree(location SourceRange, property ParseTree, isStatic bool, body *BlockTree) *ComputedPropertyGetterTree {
	t := &ComputedPropertyGetterTree{}
	t.InitParseTreeBase(ParseTreeType_COMPUTED_PROPERTY_GETTER, location)
	t.Property = property
	t.IsStatic = isStatic
	t.Body = body
	return t
}
