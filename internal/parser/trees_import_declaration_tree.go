// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportDeclarationTree.java @ closure-compiler v20261015

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

// ImportDeclarationTree is the parse tree for ParseTreeType IMPORT_DECLARATION.
type ImportDeclarationTree struct {
	ParseTreeBase

	// The identifier for the "default" import from the module.
	DefaultBindingIdentifier *IdentifierToken
	// The list of names imported from the module when using the
	// `import {...} from 'specifier'` form.
	ImportSpecifierList []ParseTree
	// The namespace into which imported names are put when using the
	// `import * as nameSpaceImportIdentifier from 'specifier'` form.
	NameSpaceImportIdentifier *IdentifierToken
	// The string identifying the module.
	ModuleSpecifier LiteralToken
}

// NewImportDeclarationTree constructs a ImportDeclarationTree.
//
// Java: ImportDeclarationTree(SourceRange, IdentifierToken, ImmutableList, IdentifierToken, LiteralToken)
func NewImportDeclarationTree(location *SourceRange, defaultBindingIdentifier *IdentifierToken, importSpecifierList []ParseTree, nameSpaceImportIdentifier *IdentifierToken, moduleSpecifier LiteralToken) *ImportDeclarationTree {
	t := &ImportDeclarationTree{}
	t.InitParseTreeBase(t, ParseTreeType_IMPORT_DECLARATION, location)
	t.DefaultBindingIdentifier = defaultBindingIdentifier
	t.ImportSpecifierList = importSpecifierList
	t.NameSpaceImportIdentifier = nameSpaceImportIdentifier
	t.ModuleSpecifier = moduleSpecifier
	return t
}
