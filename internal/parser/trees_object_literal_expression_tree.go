// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectLiteralExpressionTree.java @ closure-compiler v20261015

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

// ObjectLiteralExpressionTree is the parse tree for ParseTreeType OBJECT_LITERAL_EXPRESSION.
type ObjectLiteralExpressionTree struct {
	ParseTreeBase

	PropertyNameAndValues []ParseTree
	HasTrailingComma      bool
}

// NewObjectLiteralExpressionTree constructs a ObjectLiteralExpressionTree.
//
// Java: ObjectLiteralExpressionTree(SourceRange, ImmutableList, boolean)
func NewObjectLiteralExpressionTree(location *SourceRange, propertyNameAndValues []ParseTree, hasTrailingComma bool) *ObjectLiteralExpressionTree {
	t := &ObjectLiteralExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_OBJECT_LITERAL_EXPRESSION, location)
	t.PropertyNameAndValues = propertyNameAndValues
	t.HasTrailingComma = hasTrailingComma
	return t
}
