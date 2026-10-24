// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyDefinitionTree.java @ closure-compiler v20261021

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

// Represents a computed property in an object literal.
type ComputedPropertyDefinitionTree struct {
	ParseTreeBase

	Property ParseTree
	Value    ParseTree
}

// NewComputedPropertyDefinitionTree constructs a ComputedPropertyDefinitionTree.
//
// Java: ComputedPropertyDefinitionTree(SourceRange, ParseTree, ParseTree)
func NewComputedPropertyDefinitionTree(location SourceRange, property ParseTree, value ParseTree) *ComputedPropertyDefinitionTree {
	t := &ComputedPropertyDefinitionTree{}
	t.InitParseTreeBase(ParseTreeType_COMPUTED_PROPERTY_DEFINITION, location)
	t.Property = property
	t.Value = value
	return t
}
