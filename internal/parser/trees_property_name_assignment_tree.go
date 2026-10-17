// Port of com/google/javascript/jscomp/parsing/parser/trees/PropertyNameAssignmentTree.java @ closure-compiler v20261015

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

// PropertyNameAssignmentTree is the parse tree for ParseTreeType PROPERTY_NAME_ASSIGNMENT.
type PropertyNameAssignmentTree struct {
	ParseTreeBase

	Name Token
	// May be null in ES6 and above.
	Value ParseTree
}

// NewPropertyNameAssignmentTree constructs a PropertyNameAssignmentTree.
//
// Java: PropertyNameAssignmentTree(SourceRange, Token, ParseTree)
func NewPropertyNameAssignmentTree(location *SourceRange, name Token, value ParseTree) *PropertyNameAssignmentTree {
	t := &PropertyNameAssignmentTree{}
	t.InitParseTreeBase(t, ParseTreeType_PROPERTY_NAME_ASSIGNMENT, location)
	t.Name = name
	t.Value = value
	return t
}
