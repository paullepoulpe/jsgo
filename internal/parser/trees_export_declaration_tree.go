// Port of com/google/javascript/jscomp/parsing/parser/trees/ExportDeclarationTree.java @ closure-compiler v20261015

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

// ExportDeclarationTree is the parse tree for ParseTreeType EXPORT_DECLARATION.
type ExportDeclarationTree struct {
	ParseTreeBase

	// VariableStatement, Declaration (let, const, function or class),
	// or default with assignmentExpression
	IsDefault           bool
	IsExportAll         bool
	Declaration         ParseTree
	ExportSpecifierList []ParseTree
	From                LiteralToken
}

// NewExportDeclarationTree constructs a ExportDeclarationTree.
//
// Java: ExportDeclarationTree(SourceRange, boolean, boolean, ParseTree, ImmutableList, LiteralToken)
func NewExportDeclarationTree(location *SourceRange, isDefault bool, isExportAll bool, declaration ParseTree, exportSpecifierList []ParseTree, from LiteralToken) *ExportDeclarationTree {
	t := &ExportDeclarationTree{}
	t.InitParseTreeBase(t, ParseTreeType_EXPORT_DECLARATION, location)
	t.IsDefault = isDefault
	t.IsExportAll = isExportAll
	t.Declaration = declaration
	t.ExportSpecifierList = exportSpecifierList
	t.From = from
	return t
}
