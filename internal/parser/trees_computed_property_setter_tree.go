// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertySetterTree.java @ closure-compiler v20260415

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

// ComputedPropertySetterTree is the parse tree for ParseTreeType COMPUTED_PROPERTY_SETTER.
type ComputedPropertySetterTree struct {
	ParseTreeBase

	Property  ParseTree
	Parameter *FormalParameterListTree
	IsStatic  bool
	Body      *BlockTree
}

// NewComputedPropertySetterTree constructs a ComputedPropertySetterTree.
//
// Java: ComputedPropertySetterTree(SourceRange, ParseTree, boolean, FormalParameterListTree, BlockTree)
func NewComputedPropertySetterTree(location *SourceRange, property ParseTree, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *ComputedPropertySetterTree {
	t := &ComputedPropertySetterTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPUTED_PROPERTY_SETTER, location)
	t.Property = property
	t.IsStatic = isStatic
	t.Parameter = parameter
	t.Body = body
	return t
}
