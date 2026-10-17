// Port of com/google/javascript/jscomp/parsing/parser/trees/ClassDeclarationTree.java @ closure-compiler v20261015

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

// ClassDeclarationTree is the parse tree for ParseTreeType CLASS_DECLARATION.
type ClassDeclarationTree struct {
	ParseTreeBase

	Name       *IdentifierToken
	SuperClass ParseTree
	Elements   []ParseTree
}

// NewClassDeclarationTree constructs a ClassDeclarationTree.
//
// Java: ClassDeclarationTree(SourceRange, IdentifierToken, ParseTree, ImmutableList)
func NewClassDeclarationTree(location *SourceRange, name *IdentifierToken, superClass ParseTree, elements []ParseTree) *ClassDeclarationTree {
	t := &ClassDeclarationTree{}
	t.InitParseTreeBase(t, ParseTreeType_CLASS_DECLARATION, location)
	t.Name = name
	t.SuperClass = superClass
	t.Elements = elements
	return t
}
