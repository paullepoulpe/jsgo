// Port of com/google/javascript/jscomp/parsing/parser/trees/FieldDeclarationTree.java @ closure-compiler v20261015

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

// Tree for public property fields
type FieldDeclarationTree struct {
	ParseTreeBase

	Name        *IdentifierToken
	IsStatic    bool
	Initializer ParseTree
}

// NewFieldDeclarationTree constructs a FieldDeclarationTree.
//
// Java: FieldDeclarationTree(SourceRange, IdentifierToken, boolean, ParseTree)
func NewFieldDeclarationTree(location *SourceRange, name *IdentifierToken, isStatic bool, initializer ParseTree) *FieldDeclarationTree {
	t := &FieldDeclarationTree{}
	t.InitParseTreeBase(t, ParseTreeType_FIELD_DECLARATION, location)
	t.Name = name
	t.IsStatic = isStatic
	t.Initializer = initializer
	return t
}
