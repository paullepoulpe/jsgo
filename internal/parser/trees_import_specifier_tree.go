// Port of com/google/javascript/jscomp/parsing/parser/trees/ImportSpecifierTree.java @ closure-compiler v20261015

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

// ImportSpecifierTree is the parse tree for ParseTreeType IMPORT_SPECIFIER.
type ImportSpecifierTree struct {
	ParseTreeBase

	ImportedName    *IdentifierToken
	DestinationName *IdentifierToken
}

// NewImportSpecifierTree constructs a ImportSpecifierTree.
//
// Java: ImportSpecifierTree(SourceRange, IdentifierToken, IdentifierToken)
func NewImportSpecifierTree(location SourceRange, importedName *IdentifierToken, destinationName *IdentifierToken) *ImportSpecifierTree {
	t := &ImportSpecifierTree{}
	t.InitParseTreeBase(t, ParseTreeType_IMPORT_SPECIFIER, location)
	t.ImportedName = importedName
	t.DestinationName = destinationName
	return t
}
