// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyFieldTree.java @ closure-compiler v20261015

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

// Tree for public comptuted property fields
type ComputedPropertyFieldTree struct {
	ParseTreeBase

	Property    ParseTree
	IsStatic    bool
	Initializer ParseTree
}

// NewComputedPropertyFieldTree constructs a ComputedPropertyFieldTree.
//
// Java: ComputedPropertyFieldTree(SourceRange, ParseTree, boolean, ParseTree)
func NewComputedPropertyFieldTree(location *SourceRange, property ParseTree, isStatic bool, initializer ParseTree) *ComputedPropertyFieldTree {
	t := &ComputedPropertyFieldTree{}
	t.InitParseTreeBase(t, ParseTreeType_COMPUTED_PROPERTY_FIELD, location)
	t.Property = property
	t.IsStatic = isStatic
	t.Initializer = initializer
	return t
}
